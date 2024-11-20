package tglog

import (
	"bytes"
	"fmt"
	"github.com/tianlin0/plat-lib/goroutines"
	"github.com/tianlin0/plat-lib/logs/tglog/udpdnslb"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"encoding/hex"

	"go.uber.org/zap"
)

const (
	logQueueSize = 1024
	// datetimeLayout tgLog default datetime format
	datetimeLayout = "2006-01-02 15:04:05"
)

// TGLoggerV1 tgLogV1 client
type TGLoggerV1 struct {
	udpDNSLb *udpdnslb.UDPDnsLb
	UdpConn  *net.UDPConn
	queue    chan []byte
	qonce    sync.Once
	// hookFreeLog initialize the logger at the beginning to make sure no tgLog hook is appended to the logger
	hookFreeLog *zap.SugaredLogger
}

// NewTGLogV1 create tgLogV1 client
func NewTGLogV1(url string) (*TGLoggerV1, error) {
	tgLogger := new(TGLoggerV1)
	var domain string
	var port int
	var err error
	// setup udp connection
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("listen udp failed: %v", err))
	}
	tgLogger.UdpConn = conn
	// setup dns lb
	domainAndPort := strings.Split(url, ":")
	if len(domainAndPort) < 2 {
		return nil, fmt.Errorf(fmt.Sprintf("listen udp invalid: %v", url))
	}
	domain = domainAndPort[0]
	port, err = strconv.Atoi(domainAndPort[1])
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("tglog port failed: %v", err))
	}
	tgLogger.udpDNSLb = udpdnslb.NewUDPDnsLb(domain, port, 30*time.Second)
	tgLogger.queue = make(chan []byte, logQueueSize)

	if tgLogger.udpDNSLb != nil &&
		tgLogger.UdpConn != nil {

		tgLogger.registerQueueHandling()

		// hook-free lock
		if tgLogger.hookFreeLog == nil {
			tgLogger.hookFreeLog = zap.L().WithOptions(zap.Hooks()).Sugar()
		}

		return tgLogger, nil
	}

	return nil, fmt.Errorf(fmt.Sprintf("tglog return nil"))
}

// Log 同步自由格式log，默认使用同步方式发送
func (p *TGLoggerV1) Log(table string, fields ...interface{}) {
	p.LogSync(table, fields...)
}

// LogSync 同步自由格式log，日志会被同步地通过udp发送，性能较异步方式差
func (p *TGLoggerV1) LogSync(table string, fields ...interface{}) {
	table = FieldEscape(table)
	buf := bytes.NewBufferString(table)

	for _, f := range fields {
		buf.WriteString("|")
		buf.WriteString(strings.Replace(FormatField(f), "|", "$", -1))
	}
	buf.WriteString("\n")

	if p.hookFreeLog != nil {
		if p.UdpConn == nil || p.udpDNSLb == nil {
			p.hookFreeLog.With("log", buf.String()).Warn("tglog did not init")
			return
		}
	}

	udpAddr := p.udpDNSLb.GetUDPAddr()

	if p.hookFreeLog != nil {
		if _, err := p.UdpConn.WriteToUDP(buf.Bytes(), udpAddr); err != nil {
			p.hookFreeLog.With("tglog", buf.String()).Warnf("tglog send failed: %v", err)
		}
	}

}

// LogAsync 异步的自由格式log，高性能，但队列满后会丢弃日志
func (p *TGLoggerV1) LogAsync(table string, fields ...interface{}) {
	table = FieldEscape(table)
	buf := bytes.NewBufferString(table)

	for _, f := range fields {
		buf.WriteString("|")
		buf.WriteString(strings.Replace(FormatField(f), "|", "$", -1))
	}
	buf.WriteString("\n")

	// add to queue pending to send
	select {
	case p.queue <- buf.Bytes():
	default:
	}
}

// registerQueueHandling get logs from queue in the background and send them to the log server
func (p *TGLoggerV1) registerQueueHandling() {
	p.qonce.Do(func() {
		goroutines.GoAsyncHandler(func(params ...interface{}) {
			if p, ok := params[0].(*TGLoggerV1); ok {
				var bLog []byte
				for {
					select {
					case bLog = <-p.queue:
						if p.hookFreeLog != nil {
							if p.UdpConn == nil || p.udpDNSLb == nil {
								p.hookFreeLog.With("log", string(bLog)).Warn("tglog did not init")
								continue
							}
						}

						udpAddr := p.udpDNSLb.GetUDPAddr()
						if p.UdpConn != nil && p.hookFreeLog != nil {
							if _, err := p.UdpConn.WriteToUDP(bLog, udpAddr); err != nil {
								p.hookFreeLog.With("tglog", string(bLog)).Warnf("tglog send failed: %v", err)
							}
						}
					}
				}
			}
		}, nil, p)
	})
}

// FieldEscape escape "|" and "\n" to # and ##
func FieldEscape(old string) (result string) {
	result = strings.Replace(strings.Replace(old, "|", "#", -1), "\n", "##", -1)
	return
}

// FormatField 字段格式化
func FormatField(f interface{}) string {
	switch f.(type) {
	case string:
		return FieldEscape(f.(string))
	case []byte:
		return hex.EncodeToString(f.([]byte))
	case byte:
		return hex.EncodeToString([]byte{f.(byte)})
	case int:
		return strconv.Itoa(f.(int))
	case int64:
		return strconv.FormatInt(f.(int64), 10)
	case int32:
		return strconv.FormatInt(int64(f.(int32)), 10)
	case int16:
		return strconv.FormatInt(int64(f.(int16)), 10)
	case int8:
		return strconv.FormatInt(int64(f.(int8)), 10)
	case uint:
		return strconv.FormatUint(uint64(f.(uint)), 10)
	case uint64:
		return strconv.FormatUint(f.(uint64), 10)
	case uint32:
		return strconv.FormatUint(uint64(f.(uint32)), 10)
	case uint16:
		return strconv.FormatUint(uint64(f.(uint16)), 10)
	case float64:
		return strconv.FormatFloat(f.(float64), 'f', 3, 64)
	case float32:
		return strconv.FormatFloat(float64(f.(float32)), 'f', 3, 64)
	case bool:
		return strconv.FormatBool(f.(bool))
	case time.Time:
		return f.(time.Time).Format(datetimeLayout)
	case *time.Time:
		return f.(*time.Time).Format(datetimeLayout)
	case time.Duration:
		return f.(time.Duration).String()
	default:
		return FieldEscape(fmt.Sprintf("%v", f))
	}
}
