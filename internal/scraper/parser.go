package scraper

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	reRefAssign = regexp.MustCompile(`\$R\[\d+\]\s*=\s*`)
	reNewDate   = regexp.MustCompile(`new Date\("([^"]+)"\)`)
)

func normalizeJS(js string) string {
	s := js

	s = reRefAssign.ReplaceAllString(s, "")

	s = reNewDate.ReplaceAllString(s, `"$1"`)

	s = strings.ReplaceAll(s, "!0", "true")
	s = strings.ReplaceAll(s, "!1", "false")

	s = quoteBareKeys(s)

	return s
}

func quoteBareKeys(s string) string {
	var result strings.Builder
	result.Grow(len(s))

	i := 0
	for i < len(s) {
		if s[i] == '"' {
			result.WriteByte('"')
			i++
			for i < len(s) && s[i] != '"' {
				if s[i] == '\\' && i+1 < len(s) {
					result.WriteByte(s[i])
					i++
				}
				result.WriteByte(s[i])
				i++
			}
			if i < len(s) {
				result.WriteByte('"')
				i++
			}
			continue
		}

		if isIdentStart(s[i]) && (i == 0 || !isIdentChar(s[i-1])) {
			end := i + 1
			for end < len(s) && isIdentChar(s[end]) {
				end++
			}
			word := s[i:end]

			colonIdx := end
			for colonIdx < len(s) && s[colonIdx] == ' ' {
				colonIdx++
			}

			if colonIdx < len(s) && s[colonIdx] == ':' && word != "true" && word != "false" && word != "null" {
				result.WriteByte('"')
				result.WriteString(word)
				result.WriteByte('"')
				i = end
				continue
			}
		}

		result.WriteByte(s[i])
		i++
	}
	return result.String()
}

func isIdentStart(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isIdentChar(c byte) bool {
	return isIdentStart(c) || (c >= '0' && c <= '9')
}

func findBalanced(s string, startIdx int, open, close byte) (string, error) {
	if startIdx >= len(s) || s[startIdx] != open {
		return "", fmt.Errorf("expected '%c' at index %d", open, startIdx)
	}

	depth := 0
	inString := false
	for i := startIdx; i < len(s); i++ {
		ch := s[i]

		if inString {
			if ch == '\\' {
				i++
				continue
			}
			if ch == '"' {
				inString = false
			}
			continue
		}

		if ch == '"' {
			inString = true
			continue
		}
		if ch == open {
			depth++
		} else if ch == close {
			depth--
			if depth == 0 {
				return s[startIdx : i+1], nil
			}
		}
	}
	return "", fmt.Errorf("unbalanced '%c'...'%c' starting at index %d", open, close, startIdx)
}

func extractInlineScripts(html string) []string {
	var scripts []string
	lower := strings.ToLower(html)
	searchFrom := 0

	for {
		tagStart := strings.Index(lower[searchFrom:], "<script")
		if tagStart == -1 {
			break
		}
		tagStart += searchFrom

		tagEnd := strings.Index(html[tagStart:], ">")
		if tagEnd == -1 {
			break
		}
		tagEnd += tagStart

		tagContent := strings.ToLower(html[tagStart : tagEnd+1])
		if strings.Contains(tagContent, "src=") {
			searchFrom = tagEnd + 1
			continue
		}

		closeTag := strings.Index(lower[tagEnd+1:], "</script>")
		if closeTag == -1 {
			break
		}
		closeTag += tagEnd + 1

		scriptBody := html[tagEnd+1 : closeTag]
		if strings.Contains(scriptBody, "$R") {
			scripts = append(scripts, scriptBody)
		}

		searchFrom = closeTag + 9
	}

	return scripts
}
