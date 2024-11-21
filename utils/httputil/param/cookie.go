package param

import (
	"fmt"
	"golang.org/x/net/http/httpguts"
	"net/http"
	"net/textproto"
	"strings"
)

var cookieNameSanitizer = strings.NewReplacer("\n", "-", "\r", "-")

func sanitizeCookieName(n string) string {
	return cookieNameSanitizer.Replace(n)
}
func validCookieValueByte(b byte) bool {
	return 0x20 <= b && b < 0x7f && b != '"' && b != ';' && b != '\\'
}
func sanitizeOrWarn(fieldName string, valid func(byte) bool, v string) string {
	ok := true
	for i := 0; i < len(v); i++ {
		if valid(v[i]) {
			continue
		}
		ok = false
		break
	}
	if ok {
		return v
	}
	buf := make([]byte, 0, len(v))
	for i := 0; i < len(v); i++ {
		if b := v[i]; valid(b) {
			buf = append(buf, b)
		}
	}
	return string(buf)
}
func sanitizeCookieValue(v string) string {
	v = sanitizeOrWarn("Cookie.Value", validCookieValueByte, v)
	if len(v) == 0 {
		return v
	}
	if strings.ContainsAny(v, " ,") {
		return `"` + v + `"`
	}
	return v
}

func parseCookieValue(raw string, allowDoubleQuote bool) (string, bool) {
	// Strip the quotes, if present.
	if allowDoubleQuote && len(raw) > 1 && raw[0] == '"' && raw[len(raw)-1] == '"' {
		raw = raw[1 : len(raw)-1]
	}
	for i := 0; i < len(raw); i++ {
		if !validCookieValueByte(raw[i]) {
			return "", false
		}
	}
	return raw, true
}

func isNotToken(r rune) bool {
	return !httpguts.IsTokenRune(r)
}

func isCookieNameValid(raw string) bool {
	if raw == "" {
		return false
	}
	return strings.IndexFunc(raw, isNotToken) < 0
}

func addOneCookie(h http.Header, c *http.Cookie) http.Header {
	s := fmt.Sprintf("%s=%s", sanitizeCookieName(c.Name), sanitizeCookieValue(c.Value))
	if c := h.Get("Cookie"); c != "" {
		h.Set("Cookie", c+"; "+s)
	} else {
		h.Set("Cookie", s)
	}
	return h
}

func readCookies(h http.Header, filter string) []*http.Cookie {
	lines := h["Cookie"]
	if len(lines) == 0 {
		return []*http.Cookie{}
	}

	cookies := make([]*http.Cookie, 0, len(lines)+strings.Count(lines[0], ";"))
	for _, line := range lines {
		line = textproto.TrimString(line)

		var part string
		for len(line) > 0 { // continue since we have rest
			part, line, _ = strings.Cut(line, ";")
			part = textproto.TrimString(part)
			if part == "" {
				continue
			}
			name, val, _ := strings.Cut(part, "=")
			if !isCookieNameValid(name) {
				continue
			}
			if filter != "" && filter != name {
				continue
			}
			val, ok := parseCookieValue(val, true)
			if !ok {
				continue
			}
			cookies = append(cookies, &http.Cookie{Name: name, Value: val})
		}
	}
	return cookies
}

// getAllCookies 获取所有cookie
func getAllCookies(r *http.Request) map[string]*http.Cookie {
	newCookies := make(map[string]*http.Cookie)
	if r == nil {
		return newCookies
	}
	cookies := r.Cookies()
	for _, one := range cookies {
		oneCookie, err := r.Cookie(one.Name)
		if err != nil {
			continue
		}
		newCookies[one.Name] = oneCookie
	}
	return newCookies
}

// AddCookies 给header中添加Cookie
func AddCookies(h http.Header, c ...*http.Cookie) http.Header {
	if h == nil {
		h = http.Header{} //新建一个header
	}

	if len(c) == 0 {
		return h
	}
	for _, oneCookie := range c {
		h = addOneCookie(h, oneCookie)
	}
	return h
}
