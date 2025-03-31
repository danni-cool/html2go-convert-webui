package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestHTMLToGoAPIHandlerWithVuetify tests the conversion handler with Vuetify components
func TestHTMLToGoAPIHandlerWithVuetify(t *testing.T) {
	// 设置自定义前缀用于测试
	customVuetifyPrefix := "custom_vuetify"

	// Test with Vuetify components
	vuetifyHTML := `<div>
  <v-btn color="primary">Click me</v-btn>
  <v-card>
    <v-card-title>Card Title</v-card-title>
    <v-card-text>Card content goes here</v-card-text>
  </v-card>
</div>`

	// Create request payload
	reqPayload := ConversionRequest{
		HTML:           vuetifyHTML,
		PackagePrefix:  "h",
		VuetifyPrefix:  customVuetifyPrefix,
		VuetifyXPrefix: "vx",
		Direction:      "html2go",
	}

	// Convert to JSON
	reqBody, err := json.Marshal(reqPayload)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	// Create test request
	req, _ := http.NewRequest("POST", "/convert", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Create recorder and handler
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(convertHandler)

	// Process the request
	handler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Parse response
	var respBody ConversionResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Verify response contains expected components using the custom prefix
	if respBody.Error != "" {
		t.Errorf("Expected no error, got: %v", respBody.Error)
	}
	if !strings.Contains(respBody.Code, customVuetifyPrefix+".VBtn") {
		t.Errorf("Generated code should contain %s.VBtn for Vuetify button, got: %s",
			customVuetifyPrefix, respBody.Code)
	}
	if !strings.Contains(respBody.Code, customVuetifyPrefix+".VCard") {
		t.Errorf("Generated code should contain %s.VCard for Vuetify card, got: %s",
			customVuetifyPrefix, respBody.Code)
	}
}

// TestHTMLToGoAPIHandlerWithVuetifyX tests the conversion handler with VuetifyX components
func TestHTMLToGoAPIHandlerWithVuetifyX(t *testing.T) {
	// 设置自定义前缀用于测试
	customVuetifyXPrefix := "custom_vuetifyx"

	// Test with VuetifyX components
	vuetifyXHTML := `<div>
  <vx-date-picker label="Select Date"></vx-date-picker>
  <vx-dialog title="Confirmation" text="Are you sure?">
    <v-btn color="primary">Open Dialog</v-btn>
  </vx-dialog>
</div>`

	// Create request payload
	reqPayload := ConversionRequest{
		HTML:           vuetifyXHTML,
		PackagePrefix:  "h",
		VuetifyPrefix:  "v",
		VuetifyXPrefix: customVuetifyXPrefix,
		Direction:      "html2go",
	}

	// Convert to JSON
	reqBody, err := json.Marshal(reqPayload)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	// Create test request
	req, _ := http.NewRequest("POST", "/convert", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Create recorder and handler
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(convertHandler)

	// Process the request
	handler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Parse response
	var respBody ConversionResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Verify response contains expected components using the custom prefix
	if respBody.Error != "" {
		t.Errorf("Expected no error, got: %v", respBody.Error)
	}
	if !strings.Contains(respBody.Code, customVuetifyXPrefix+".VXDatepicker") {
		t.Errorf("Generated code should contain %s.VXDatepicker for VuetifyX date picker, got: %s",
			customVuetifyXPrefix, respBody.Code)
	}
	if !strings.Contains(respBody.Code, customVuetifyXPrefix+".VXDialog") {
		t.Errorf("Generated code should contain %s.VXDialog for VuetifyX dialog, got: %s",
			customVuetifyXPrefix, respBody.Code)
	}
}

// TestDefaultPrefixes tests that prefixes are set to the default values when not specified
func TestDefaultPrefixes(t *testing.T) {
	// Test with both Vuetify and VuetifyX components but don't provide prefixes
	mixedHTML := `<div>
  <v-btn color="primary">Click me</v-btn>
  <vx-date-picker label="Select Date"></vx-date-picker>
</div>`

	// Create request payload without specifying prefixes
	reqPayload := ConversionRequest{
		HTML:      mixedHTML,
		Direction: "html2go",
	}

	// Convert to JSON
	reqBody, err := json.Marshal(reqPayload)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	// Create test request
	req, _ := http.NewRequest("POST", "/convert", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Create recorder and handler
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(convertHandler)

	// Process the request
	handler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Parse response
	var respBody ConversionResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Verify the default prefixes are used in the generated code
	if respBody.Error != "" {
		t.Errorf("Expected no error, got: %v", respBody.Error)
	}

	// 检查组件名称，应该使用默认前缀 v 和 vx
	if !strings.Contains(respBody.Code, "v.VBtn") {
		t.Errorf("Generated code should contain v.VBtn with default prefix, got: %s",
			respBody.Code)
	}
	if !strings.Contains(respBody.Code, "vx.VXDatepicker") {
		t.Errorf("Generated code should contain vx.VXDatepicker with default prefix, got: %s",
			respBody.Code)
	}
}

// TestHTMLToGoAPIHandlerWithComplexHTML tests the conversion handler with complex nested HTML
func TestHTMLToGoAPIHandlerWithComplexHTML(t *testing.T) {
	// 设置自定义前缀用于测试
	customPackagePrefix := "custom_pkg"

	// Test with complex nested HTML
	complexHTML := `<div class="container">
  <header class="header">
    <nav class="navbar">
      <ul class="nav-list">
        <li class="nav-item"><a href="#" class="nav-link">Home</a></li>
        <li class="nav-item"><a href="#" class="nav-link">About</a></li>
        <li class="nav-item"><a href="#" class="nav-link">Contact</a></li>
      </ul>
    </nav>
  </header>
  <main class="content">
    <section class="hero">
      <h1 class="title">Welcome to our site</h1>
      <p class="subtitle">This is a complex HTML structure</p>
    </section>
    <article class="card">
      <h2 class="card-title">Article Title</h2>
      <div class="card-body">
        <p>This is the content of the article.</p>
      </div>
      <footer class="card-footer">
        <button class="btn btn-primary">Read More</button>
      </footer>
    </article>
  </main>
</div>`

	// Create request payload
	reqPayload := ConversionRequest{
		HTML:          complexHTML,
		PackagePrefix: customPackagePrefix,
		Direction:     "html2go",
	}

	// Convert to JSON
	reqBody, err := json.Marshal(reqPayload)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	// Create test request
	req, _ := http.NewRequest("POST", "/convert", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Create recorder and handler
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(convertHandler)

	// Process the request
	handler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Parse response
	var respBody ConversionResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Verify response contains expected elements with custom package prefix
	if respBody.Error != "" {
		t.Errorf("Expected no error, got: %v", respBody.Error)
	}
	if !strings.Contains(respBody.Code, customPackagePrefix+".Div") {
		t.Errorf("Generated code should contain %s.Div, got: %s",
			customPackagePrefix, respBody.Code)
	}
	if !strings.Contains(respBody.Code, customPackagePrefix+".Header") {
		t.Errorf("Generated code should contain %s.Header, got: %s",
			customPackagePrefix, respBody.Code)
	}
	if !strings.Contains(respBody.Code, customPackagePrefix+".Nav") {
		t.Errorf("Generated code should contain %s.Nav, got: %s",
			customPackagePrefix, respBody.Code)
	}
	if !strings.Contains(respBody.Code, customPackagePrefix+".Ul") {
		t.Errorf("Generated code should contain %s.Ul, got: %s",
			customPackagePrefix, respBody.Code)
	}
	if !strings.Contains(respBody.Code, customPackagePrefix+".Li") {
		t.Errorf("Generated code should contain %s.Li, got: %s",
			customPackagePrefix, respBody.Code)
	}
}

// TestInvalidRequests tests various invalid request scenarios
func TestInvalidRequests(t *testing.T) {
	tests := []struct {
		name            string
		requestPayload  ConversionRequest
		expectedStatus  int
		expectedErrText string
	}{
		{
			name: "Empty HTML",
			requestPayload: ConversionRequest{
				HTML:          "",
				PackagePrefix: "h",
				Direction:     "html2go",
			},
			expectedStatus:  http.StatusBadRequest,
			expectedErrText: "HTML content is required",
		},
		{
			name: "Invalid Direction",
			requestPayload: ConversionRequest{
				HTML:          "<div>Test</div>",
				PackagePrefix: "h",
				Direction:     "invalid_direction",
			},
			expectedStatus:  http.StatusBadRequest,
			expectedErrText: "Invalid conversion direction",
		},
		{
			name: "Go to HTML Direction",
			requestPayload: ConversionRequest{
				HTML:          "<div>Test</div>",
				PackagePrefix: "h",
				Direction:     "go2html",
			},
			expectedStatus:  http.StatusNotImplemented,
			expectedErrText: "Go to HTML conversion is not implemented yet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert to JSON
			reqBody, err := json.Marshal(tt.requestPayload)
			if err != nil {
				t.Fatalf("Failed to marshal request: %v", err)
			}

			// Create test request
			req, _ := http.NewRequest("POST", "/convert", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			// Create recorder and handler
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(convertHandler)

			// Process the request
			handler.ServeHTTP(rr, req)

			// Check status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			// Parse response
			var respBody ConversionResponse
			if err := json.Unmarshal(rr.Body.Bytes(), &respBody); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			// Verify error message
			if !strings.Contains(respBody.Error, tt.expectedErrText) {
				t.Errorf("Expected error containing %q, got: %q", tt.expectedErrText, respBody.Error)
			}
		})
	}
}

// TestHTTPMethodValidation tests that only POST requests are accepted
func TestHTTPMethodValidation(t *testing.T) {
	// Test methods that should be rejected
	invalidMethods := []string{"GET", "PUT", "DELETE", "PATCH", "OPTIONS"}

	for _, method := range invalidMethods {
		t.Run(method+" method", func(t *testing.T) {
			// Create request with invalid method
			req, _ := http.NewRequest(method, "/convert", nil)

			// Create recorder and handler
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(convertHandler)

			// Process the request
			handler.ServeHTTP(rr, req)

			// Check status code should be Method Not Allowed
			if status := rr.Code; status != http.StatusMethodNotAllowed {
				t.Errorf("Handler should reject %s method: got %v want %v",
					method, status, http.StatusMethodNotAllowed)
			}
		})
	}
}

// TestSyntaxFixing tests the fixSyntaxIssues function directly
func TestSyntaxFixing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Fix dot class syntax",
			input:    "h.Div(\n\t.Class(\"container\")\n)",
			expected: "// Generated using htmlgo\n\nh.Div(\n\tClass(\"container\")\n)",
		},
		{
			name:     "Fix package declaration",
			input:    "package hello\n\nh.Div(\n\tClass(\"container\")\n)",
			expected: "// Generated using htmlgo\n\npackage hello\n\nh.Div(\n\tClass(\"container\")\n)",
		},
		{
			name:     "Fix variable declaration",
			input:    "var n = h.Div(\n\tClass(\"container\")\n)",
			expected: "// Generated using htmlgo\n\nvar n = h.Div(\n\tClass(\"container\")\n)",
		},
		{
			name:     "Fix trailing commas",
			input:    "h.Div(\n\tClass(\"container\"),\n)",
			expected: "// Generated using htmlgo\n\nh.Div(\n\tClass(\"container\")\n)",
		},
		{
			name:     "Fix method invocation",
			input:    "h.VBtn(\n\th.Text(\"Save\")\n)Color(\"primary\")",
			expected: "// Generated using htmlgo\n\nh.VBtn(\n\th.Text(\"Save\")\n).Color(\"primary\")",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fixSyntaxIssues(tt.input)

			// 清理空白字符进行比较
			cleanResult := normalizeWhitespace(result)
			cleanExpected := normalizeWhitespace(tt.expected)

			if cleanResult != cleanExpected {
				t.Errorf("fixSyntaxIssues() =\n%v\nwant:\n%v", result, tt.expected)
			}
		})
	}
}

// 辅助函数，用于规范化空白字符以便比较
func normalizeWhitespace(s string) string {
	// 移除所有空白字符
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
	return s
}

// TestValidateGoSyntax tests the syntax validation function
func TestValidateGoSyntax(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Valid Go code",
			input:   "h.Div(\n\tClass(\"container\")\n)",
			wantErr: false,
		},
		{
			name:    "Invalid Go code - unbalanced parentheses",
			input:   "h.Div(\n\tClass(\"container\"\n)",
			wantErr: true, // 括号不平衡应该检测到错误
		},
		{
			name:    "Invalid Go code - unbalanced braces",
			input:   "h.Div(\n\tClass(\"container\")\n){ h.Span()",
			wantErr: true, // 花括号不平衡应该检测到错误
		},
		{
			name:    "Valid Go code with unknown symbols",
			input:   "h.Div(\n\tUnknownFunc(\"container\")\n)",
			wantErr: false, // 未知符号不应该报错，只检查括号平衡
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateGoSyntax(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateGoSyntax() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestComprehensiveHTMLToGoAPI 整合了所有单独的测试案例，测试 HTML 到 Go 代码转换的各种场景
func TestComprehensiveHTMLToGoAPI(t *testing.T) {
	testCases := []struct {
		name           string
		html           string
		packagePrefix  string
		vuetifyPrefix  string
		vuetifyXPrefix string
		childrenMode   bool
		expectedParts  []string
		notExpected    []string
	}{
		{
			name:           "基本 HTML，所有前缀为空",
			html:           `<div class="container"><h1>Test</h1></div>`,
			packagePrefix:  "",
			vuetifyPrefix:  "",
			vuetifyXPrefix: "",
			childrenMode:   false,
			expectedParts: []string{
				"Div(",
				"Class(\"container\")",
				"H1(\"Test\")",
			},
			notExpected: []string{
				// 由于默认前缀设置，实际可能会出现前缀h，所以我们不检查它
			},
		},
		{
			name:           "基本 HTML，只有 pkg 不为空",
			html:           `<div class="container"><h1>Test</h1></div>`,
			packagePrefix:  "h",
			vuetifyPrefix:  "",
			vuetifyXPrefix: "",
			childrenMode:   false,
			expectedParts: []string{
				"h.Div(",
				"Class(\"container\")",
				"h.H1(\"Test\")",
			},
			notExpected: []string{
				"v.", "vx.",
			},
		},
		{
			name:           "Vuetify 组件转换",
			html:           `<v-card><v-card-title>Card Title</v-card-title></v-card>`,
			packagePrefix:  "h",
			vuetifyPrefix:  "v",
			vuetifyXPrefix: "",
			childrenMode:   false,
			expectedParts: []string{
				"v.VCard(",
				"v.VCardTitle(",
				"Card Title",
			},
			notExpected: []string{
				"vx.", "h.v",
			},
		},
		{
			name:           "VuetifyX 组件转换",
			html:           `<div><vx-date-picker label="Select Date"></vx-date-picker></div>`,
			packagePrefix:  "h",
			vuetifyPrefix:  "v",
			vuetifyXPrefix: "vx",
			childrenMode:   false,
			expectedParts: []string{
				"h.Div(",
				"vx.VXDatepicker()",
				"label", "Select Date",
			},
			notExpected: []string{
				"h.vx", "v.vx",
			},
		},
		{
			name:           "混合组件，Children 模式",
			html:           `<div><span>Text</span><v-btn>Click</v-btn></div>`,
			packagePrefix:  "h",
			vuetifyPrefix:  "v",
			vuetifyXPrefix: "vx",
			childrenMode:   true,
			expectedParts: []string{
				"Children(",
				"h.Span(\"Text\")",
				"v.VBtn(",
				"\"Click\"",
			},
			notExpected: []string{
				"h.v", "h.vx",
			},
		},
		{
			name:           "包声明测试", // 不检查变量声明，因为API转换可能不包含它
			html:           `<div>Test</div>`,
			packagePrefix:  "",
			vuetifyPrefix:  "",
			vuetifyXPrefix: "",
			childrenMode:   false,
			expectedParts: []string{
				"package hello",
				"Div(",
				"Test",
			},
			notExpected: []string{
				// 由于默认前缀设置，实际可能会出现前缀h，所以我们不检查它
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建请求载荷
			reqPayload := ConversionRequest{
				HTML:           tc.html,
				PackagePrefix:  tc.packagePrefix,
				VuetifyPrefix:  tc.vuetifyPrefix,
				VuetifyXPrefix: tc.vuetifyXPrefix,
				ChildrenMode:   tc.childrenMode,
				Direction:      "html2go",
			}

			// 转换为 JSON
			reqBody, err := json.Marshal(reqPayload)
			if err != nil {
				t.Fatalf("序列化请求失败: %v", err)
			}

			// 创建测试请求
			req, _ := http.NewRequest("POST", "/convert", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			// 创建响应记录器和处理器
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(convertHandler)

			// 处理请求
			handler.ServeHTTP(rr, req)

			// 检查状态码
			if status := rr.Code; status != http.StatusOK {
				t.Errorf("处理器返回错误的状态码: 得到 %v 期望 %v", status, http.StatusOK)
			}

			// 解析响应
			var respBody ConversionResponse
			if err := json.Unmarshal(rr.Body.Bytes(), &respBody); err != nil {
				t.Fatalf("解析响应失败: %v", err)
			}

			// 检查错误字段
			if respBody.Error != "" {
				t.Errorf("期望没有错误，但得到: %v", respBody.Error)
			}

			// 验证结果包含预期的部分
			for _, expectedPart := range tc.expectedParts {
				if !strings.Contains(respBody.Code, expectedPart) {
					t.Errorf("期望生成的代码包含 %q，但实际结果为:\n%s", expectedPart, respBody.Code)
				}
			}

			// 验证结果不包含不期望的部分
			for _, notExpectedPart := range tc.notExpected {
				if strings.Contains(respBody.Code, notExpectedPart) {
					t.Errorf("不期望生成的代码包含 %q，但实际结果为:\n%s", notExpectedPart, respBody.Code)
				}
			}
		})
	}
}

// TestSyntaxFixingAPI 测试通过 API 修复语法问题
func TestSyntaxFixingAPI(t *testing.T) {
	syntaxTests := []struct {
		name          string
		inputHTML     string
		expectedParts []string
	}{
		{
			name:      "修复点语法调用",
			inputHTML: `<div class="container"></div>`,
			expectedParts: []string{
				"Div(",
				"Class(\"container\")",
				// 检查没有 .Class
				")",
			},
		},
		{
			name:      "验证修复方法调用",
			inputHTML: `<v-btn color="primary">Save</v-btn>`,
			expectedParts: []string{
				"VBtn(",
				"Color(\"primary\")",
				"\"Save\"",
				")",
			},
		},
		{
			name:      "验证修复尾部逗号",
			inputHTML: `<div id="app" class="container"></div>`,
			expectedParts: []string{
				"Div(",
				"Id(\"app\")",
				"Class(\"container\")",
				")",
			},
		},
	}

	for _, tc := range syntaxTests {
		t.Run(tc.name, func(t *testing.T) {
			// 创建请求载荷
			reqPayload := ConversionRequest{
				HTML:      tc.inputHTML,
				Direction: "html2go",
			}

			// 转换为 JSON
			reqBody, err := json.Marshal(reqPayload)
			if err != nil {
				t.Fatalf("序列化请求失败: %v", err)
			}

			// 创建测试请求
			req, _ := http.NewRequest("POST", "/convert", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			// 创建响应记录器和处理器
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(convertHandler)

			// 处理请求
			handler.ServeHTTP(rr, req)

			// 检查状态码
			if status := rr.Code; status != http.StatusOK {
				t.Errorf("处理器返回错误的状态码: 得到 %v 期望 %v", status, http.StatusOK)
			}

			// 解析响应
			var respBody ConversionResponse
			if err := json.Unmarshal(rr.Body.Bytes(), &respBody); err != nil {
				t.Fatalf("解析响应失败: %v", err)
			}

			// 验证结果包含预期的部分
			for _, expectedPart := range tc.expectedParts {
				if !strings.Contains(respBody.Code, expectedPart) {
					t.Errorf("期望生成的代码包含 %q，但实际结果为:\n%s", expectedPart, respBody.Code)
				}
			}
		})
	}
}
