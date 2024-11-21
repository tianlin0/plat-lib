package gdplog

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/tianlin0/plat-lib/utils/httputil"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	ContextKeyRequestID = "log_request_id"
	ContextKeyTraceID   = "log_trace_id"
)

// SetRouterLogger 根据请求参数生成一个包含请求信息的Logger，并设置到上下文中
func SetRouterLogger(projectParam, paasParam string) gin.HandlerFunc {
	return SetRouterLoggerWithBusiness(projectParam, paasParam, "")
}

func SetRouterLoggerWithBusiness(projectParam, paasParam, businessParam string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-Gdp-Userid")
		requestID := c.GetHeader("X-Request-Id")
		traceID := c.GetHeader("traceparent")
		if requestID == "" {
			// 请求头X-Request-Id不存在，生成一个新的uuid
			requestUUID := httputil.GetLogId()
			requestID = fmt.Sprintf("g-%s", requestUUID) // 用于区分本地生成的Request-ID
			c.Request.Header.Set("X-Request-Id", requestID)
		}
		c.Set(ContextKeyRequestID, requestID)
		c.Set(ContextKeyTraceID, traceID)
		newCtx := context.WithValue(c.Request.Context(), ContextKeyRequestID, requestID)
		newCtx = context.WithValue(newCtx, ContextKeyTraceID, traceID)

		var project, paas, business string
		if projectParam != "" {
			project = c.Param(projectParam)
		}
		if paasParam != "" {
			paas = c.Param(paasParam)
		}
		if businessParam != "" {
			business = c.Param(businessParam)
		}

		target := ""
		if project != "" || paas != "" {
			target = fmt.Sprintf("%s/%s", project, paas)
		} else if business != "" {
			target = business
		}

		fields := logrus.Fields{}
		if requestID != "" {
			fields[RequestIDField] = requestID
		}
		if userID != "" {
			fields[UserIDField] = userID
		}
		if target != "" {
			fields[TargetField] = target
		}

		logger := defaultLogger.WithFields(fields)

		newCtx = context.WithValue(newCtx, contextKeyLogger, logger)
		c.Request = c.Request.WithContext(newCtx)
	}
}

// GinLogFormatter is the default log format function Logger middleware uses.
var GinLogFormatter = func(param gin.LogFormatterParams) string {
	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}
	var requestIDStr string
	requestID := param.Keys[ContextKeyRequestID]
	if requestID != nil {
		requestIDStr, _ = requestID.(string)
	}

	return fmt.Sprintf("%v [GIN] [%s] | %3d | %13v | %15s | %-7s %#v\n%s",
		param.TimeStamp.Format("2006/01/02 15:04:05.000"),
		requestIDStr,
		param.StatusCode,
		param.Latency,
		param.ClientIP,
		param.Method,
		param.Path,
		param.ErrorMessage,
	)
}

// WrappedGinWriter 封装Gin的ResponseWriter，同时提供Read接口
type WrappedGinWriter struct {
	responseBuffer *bytes.Buffer
	responseWriter gin.ResponseWriter
}

// WrapGinWriter 封装Gin的ResponseWriter，同时提供Read接口
func WrapGinWriter(responseWriter gin.ResponseWriter) *WrappedGinWriter {
	return &WrappedGinWriter{
		responseBuffer: &bytes.Buffer{},
		responseWriter: responseWriter,
	}
}

// Header 封装ResponseWriter的方法
func (w *WrappedGinWriter) Header() http.Header {
	return w.responseWriter.Header()
}

// Write 封装ResponseWriter的方法
func (w *WrappedGinWriter) Write(b []byte) (int, error) {
	// 仅将文本数据写入到Buffer中
	contentType := w.responseWriter.Header().Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") || strings.HasPrefix(contentType, "text/") {
		n, err := w.responseBuffer.Write(b)
		if err != nil {
			return n, err
		}
	}

	return w.responseWriter.Write(b)
}

// WriteHeader 封装ResponseWriter的方法
func (w *WrappedGinWriter) WriteHeader(statusCode int) {
	w.responseWriter.WriteHeader(statusCode)
}

// Hijack 封装ResponseWriter的方法
func (w *WrappedGinWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.responseWriter.Hijack()
}

// Flush 封装ResponseWriter的方法
func (w *WrappedGinWriter) Flush() {
	w.responseWriter.Flush()
}

// CloseNotify 封装ResponseWriter的方法
func (w *WrappedGinWriter) CloseNotify() <-chan bool {
	return w.responseWriter.CloseNotify()
}

// Status 封装ResponseWriter的方法
func (w *WrappedGinWriter) Status() int {
	return w.responseWriter.Status()
}

// Size 封装ResponseWriter的方法
func (w *WrappedGinWriter) Size() int {
	return w.responseWriter.Size()
}

// WriteString 封装ResponseWriter的方法
func (w *WrappedGinWriter) WriteString(str string) (int, error) {
	n, err := w.responseBuffer.WriteString(str)
	if err != nil {
		return n, err
	}
	return w.responseWriter.WriteString(str)
}

// Written 封装ResponseWriter的方法
func (w *WrappedGinWriter) Written() bool {
	return w.responseWriter.Written()
}

// WriteHeaderNow 封装ResponseWriter的方法
func (w *WrappedGinWriter) WriteHeaderNow() {
	w.responseWriter.WriteHeaderNow()
}

// Pusher 封装ResponseWriter的方法
func (w *WrappedGinWriter) Pusher() http.Pusher {
	return w.responseWriter.Pusher()
}

// ReadBodyAndRewind 读取Writer当中的Body，并将Reader的指向指回开头
func (w *WrappedGinWriter) ReadBodyAndRewind() []byte {
	if w.responseBuffer != nil {
		var newBuf bytes.Buffer
		tee := io.TeeReader(w.responseBuffer, &newBuf)
		b, _ := ioutil.ReadAll(tee)
		w.responseBuffer = &newBuf
		return b
	}

	return []byte{}
}

// String 读取Writer当中的Body，将其转化为可打印字符串，并将Reader的指向指回开头
func (w *WrappedGinWriter) String() string {
	bodyByes := w.ReadBodyAndRewind()
	if len(bodyByes) > 0 {
		return string(bodyByes)
	}
	return ""
}

// RequestAndResponseBodyPrinting 通过日志的方式打印接收到的请求与响应Body数据，不包含URL或Header信息
func RequestAndResponseBodyPrinting() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := Logger(ctx)

		// 打印请求Body，且仅读取文本数据
		if c.Request.Body != nil {
			var buf bytes.Buffer
			tee := io.TeeReader(c.Request.Body, &buf)
			b, _ := ioutil.ReadAll(tee)
			_ = c.Request.Body.Close()
			c.Request.Body = ioutil.NopCloser(&buf)
			if len(b) > 0 {
				logger.Infof("http request received: %s", string(b))
			}
		}

		savedWriter := c.Writer
		wrappedWriter := WrapGinWriter(c.Writer)
		c.Writer = wrappedWriter
		defer func() {
			c.Writer = savedWriter
		}()

		c.Next()

		// 打印响应Body
		responseBodyStr := wrappedWriter.String()
		if responseBodyStr != "" {
			logger.Infof("http response sent: %s", responseBodyStr)
		}
	}
}
