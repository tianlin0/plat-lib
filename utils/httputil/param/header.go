package param

import (
	"net/http"
)

// getAllHeaders 取得所有请求的header列表
func getAllHeaders(r *http.Request) http.Header {
	headers := http.Header{}
	if r == nil {
		return headers
	}
	if r.Header != nil && len(r.Header) > 0 {
		headers = r.Header.Clone()
	}
	return headers
}
