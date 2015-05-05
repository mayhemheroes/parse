package parse // import "github.com/tdewolff/parse"

import (
	"encoding/base64"
	"net/url"
)

func Copy(src []byte) (dst []byte) {
	dst = make([]byte, len(src))
	copy(dst, src)
	return
}

func ToLower(src []byte) []byte {
	for i, c := range src {
		if c >= 'A' && c <= 'Z' {
			src[i] = c + ('a' - 'A')
		}
	}
	return src
}

func CopyToLower(src []byte) []byte {
	dst := Copy(src)
	for i, c := range dst {
		if c >= 'A' && c <= 'Z' {
			dst[i] = c + ('a' - 'A')
		}
	}
	return dst
}

func Equal(s, match []byte) bool {
	if len(s) != len(match) {
		return false
	}
	for i, c := range match {
		if s[i] != c {
			return false
		}
	}
	return true
}

func EqualCaseInsensitive(s, matchLower []byte) bool {
	if len(s) != len(matchLower) {
		return false
	}
	for i, c := range matchLower {
		if s[i] != c && (c < 'A' && c > 'Z' || s[i]+('a'-'A') != c) {
			return false
		}
	}
	return true
}

// IsWhitespace returns true for space, \n, \t, \f, \r.
func IsWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f'
}

func IsAllWhitespace(b []byte) bool {
	for _, c := range b {
		if !IsWhitespace(c) {
			return false
		}
	}
	return true
}

// Trim removes any character from the front and the end that matches the function.
func Trim(b []byte, f func(byte) bool) []byte {
	n := len(b)
	start := n
	for i := 0; i < n; i++ {
		if !f(b[i]) {
			start = i
			break
		}
	}
	end := n
	for i := n - 1; i >= start; i-- {
		if !f(b[i]) {
			end = i + 1
			break
		}
	}
	return b[start:end]
}

// ReplaceMultiple replaces any character serie that matches the function into a single character.
func ReplaceMultiple(b []byte, f func(byte) bool, r byte) []byte {
	j := 0
	start := 0
	prevMatch := false
	for i, c := range b {
		if f(c) {
			if !prevMatch {
				prevMatch = true
				b[i] = r
			} else {
				if start < i {
					if start != 0 {
						j += copy(b[j:], b[start:i])
					} else {
						j += i
					}
				}
				start = i + 1
			}
		} else {
			prevMatch = false
		}
	}
	if start != 0 {
		j += copy(b[j:], b[start:])
		return b[:j]
	}
	return b
}

func NormalizeContentType(b []byte) []byte {
	j := 0
	start := 0
	inString := false
	for i, c := range b {
		if !inString && IsWhitespace(c) {
			if start != 0 {
				j += copy(b[j:], b[start:i])
			} else {
				j += i
			}
			start = i + 1
		} else if c == '"' {
			inString = !inString
		}
	}
	if start != 0 {
		j += copy(b[j:], b[start:])
		return ToLower(b[:j])
	}
	return ToLower(b)
}

// SplitDataURI splits the given URLToken and returns the mediatype, data and ok.
func SplitDataURI(dataURI []byte) ([]byte, []byte, bool) {
	if len(dataURI) > 5 && Equal(dataURI[:5], []byte("data:")) {
		dataURI = dataURI[5:]
		inBase64 := false
		mediatype := []byte{}
		i := 0
		for j, c := range dataURI {
			if c == '=' || c == ';' || c == ',' {
				if c != '=' && Equal(Trim(dataURI[i:j], IsWhitespace), []byte("base64")) {
					if len(mediatype) > 0 {
						mediatype = mediatype[:len(mediatype)-1]
					}
					inBase64 = true
					i = j
				} else if c != ',' {
					mediatype = append(append(mediatype, Trim(dataURI[i:j], IsWhitespace)...), c)
					i = j + 1
				} else {
					mediatype = append(mediatype, Trim(dataURI[i:j], IsWhitespace)...)
				}
				if c == ',' {
					if len(mediatype) == 0 || mediatype[0] == ';' {
						mediatype = []byte("text/plain")
					}
					data := dataURI[j+1:]
					if inBase64 {
						decoded := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
						n, err := base64.StdEncoding.Decode(decoded, data)
						if err != nil {
							return []byte{}, []byte{}, false
						}
						data = decoded[:n]
					} else {
						unescaped, err := url.QueryUnescape(string(data))
						if err != nil {
							return []byte{}, []byte{}, false
						}
						data = []byte(unescaped)
					}
					return mediatype, data, true
				}
			}
		}
	}
	return []byte{}, []byte{}, false
}
