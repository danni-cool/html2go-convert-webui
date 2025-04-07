package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/zhangshanwen/html2go/parse"
)

// ConversionRequest represents the JSON request body for conversion
type ConversionRequest struct {
	HTML           string `json:"html"`
	PackagePrefix  string `json:"packagePrefix"`
	VuetifyPrefix  string `json:"vuetifyPrefix"`
	VuetifyXPrefix string `json:"vuetifyXPrefix"`
	Direction      string `json:"direction"`
	ChildrenMode   bool   `json:"childrenMode"`
}

// ConversionResponse represents the JSON response for conversion
type ConversionResponse struct {
	Code  string `json:"code,omitempty"`
	HTML  string `json:"html,omitempty"`
	Error string `json:"error,omitempty"`
}

// Handler is the API entry point for Vercel serverless functions
func Handler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req ConversionRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		sendJSONError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.HTML == "" {
		sendJSONError(w, "HTML content is required", http.StatusBadRequest)
		return
	}

	// Set default prefixes if not provided
	if req.VuetifyPrefix == "" {
		req.VuetifyPrefix = "v"
	}
	if req.VuetifyXPrefix == "" {
		req.VuetifyXPrefix = "vx"
	}

	// Process based on direction
	var response ConversionResponse
	switch req.Direction {
	case "html2go":
		code, err := convertHTMLToGo(req.HTML, req.PackagePrefix, req.VuetifyPrefix, req.VuetifyXPrefix, req.ChildrenMode)
		if err != nil {
			sendJSONError(w, fmt.Sprintf("HTML to Go conversion error: %v", err), http.StatusInternalServerError)
			return
		}
		response.Code = code
	case "go2html":
		// Not implemented yet - might be added in a future update
		sendJSONError(w, "Go to HTML conversion is not implemented yet", http.StatusNotImplemented)
		return
	default:
		sendJSONError(w, "Invalid conversion direction", http.StatusBadRequest)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func convertHTMLToGo(htmlContent, packagePrefix, vuetifyPrefix, vuetifyXPrefix string, childrenMode bool) (string, error) {
	// The function takes a reader, so we need to convert our string to a reader
	reader := strings.NewReader(htmlContent)

	// Using the Vuetify branch API
	// Generate HTML Go code with support for Vuetify components
	goCode := parse.GenerateHTMLGo(packagePrefix, vuetifyPrefix, vuetifyXPrefix, childrenMode, reader)

	// Process the generated code to extract important parts
	goCode = stripWrappers(goCode)

	return goCode, nil
}

// stripWrappers removes the package declaration and h.Body wrapper, and trims whitespace
func stripWrappers(code string) string {
	// 移除包声明
	if strings.Contains(code, "package hello") {
		// 找到包声明之后的第一个非空行
		lines := strings.Split(code, "\n")
		startLine := 0
		for i, line := range lines {
			if strings.Contains(line, "package hello") {
				startLine = i + 1
				break
			}
		}

		// 跳过空行
		for startLine < len(lines) && strings.TrimSpace(lines[startLine]) == "" {
			startLine++
		}

		if startLine < len(lines) {
			code = strings.Join(lines[startLine:], "\n")
		}
	}

	// 移除 var n = 声明
	if strings.Contains(code, "var n =") {
		// 移除var n =行
		lines := strings.Split(code, "\n")
		startLine := 0
		for i, line := range lines {
			if strings.Contains(line, "var n =") {
				startLine = i + 1
				// 如果这行不包含h.Body，可能h.Body在下一行
				if !strings.Contains(line, "h.Body(") {
					for j := startLine; j < len(lines); j++ {
						if strings.Contains(lines[j], "h.Body(") {
							startLine = j + 1
							break
						}
					}
				}
				break
			}
		}

		if startLine > 0 && startLine < len(lines) {
			// 拼接剩余内容
			code = strings.Join(lines[startLine-1:], "\n")
		}
	}

	// 尝试定位并提取h.Body的内容
	bodyStart := strings.Index(code, "h.Body(")
	if bodyStart >= 0 {
		// 找到第一个左括号后的内容
		openingIndex := bodyStart + 7
		// 找到匹配的最后一个闭合括号
		depth := 1
		closingIndex := -1

		for i := openingIndex; i < len(code); i++ {
			if code[i] == '(' {
				depth++
			} else if code[i] == ')' {
				depth--
				if depth == 0 {
					closingIndex = i
					break
				}
			}
		}

		if closingIndex > openingIndex {
			// 提取 h.Body(...) 内部的内容
			innerCode := code[openingIndex:closingIndex]
			code = innerCode
		}
	}

	// 移除最后可能存在的闭合括号
	code = removeTrailingParentheses(code)

	// 移除末尾的逗号
	code = removeTrailingCommas(code)

	// 清理代码：修剪前后的空白
	return strings.TrimSpace(code)
}

// removeTrailingParentheses 移除字符串末尾多余的闭合括号
func removeTrailingParentheses(s string) string {
	s = strings.TrimSpace(s)
	// 计算开括号和闭括号数量
	openCount := strings.Count(s, "(")
	closeCount := strings.Count(s, ")")

	// 如果闭括号比开括号多，移除末尾多余的闭括号
	if closeCount > openCount {
		// 从末尾开始寻找多余的闭括号
		excess := closeCount - openCount
		for i := 0; i < excess; i++ {
			if strings.HasSuffix(s, ")") {
				s = strings.TrimSuffix(s, ")")
				s = strings.TrimSpace(s)
			}
		}
	}
	return s
}

// removeTrailingCommas 移除代码末尾的逗号
func removeTrailingCommas(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, ",") {
		return strings.TrimSuffix(s, ",")
	}

	// 逐行检查，移除最后一行非空行末尾的逗号
	lines := strings.Split(s, "\n")
	lastNonEmptyLine := -1

	for i := len(lines) - 1; i >= 0; i-- {
		if strings.TrimSpace(lines[i]) != "" {
			lastNonEmptyLine = i
			break
		}
	}

	if lastNonEmptyLine >= 0 && strings.HasSuffix(strings.TrimSpace(lines[lastNonEmptyLine]), ",") {
		lines[lastNonEmptyLine] = strings.TrimSuffix(strings.TrimSpace(lines[lastNonEmptyLine]), ",")
		s = strings.Join(lines, "\n")
	}

	return s
}

func sendJSONError(w http.ResponseWriter, errMsg string, statusCode int) {
	response := ConversionResponse{
		Error: errMsg,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, errMsg, statusCode)
	}
}
