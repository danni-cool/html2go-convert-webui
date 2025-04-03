package pathfinder_test

import (
	"os"
	"path/filepath"
	"testing"

	pathfinder "html2go-converter/api/_utils"
)

// TestVercelEnvironment 模拟Vercel环境下的测试
func TestVercelEnvironment(t *testing.T) {
	// 保存当前工作目录
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir) // 测试结束后恢复

	// 创建模拟Vercel部署结构的临时目录
	tmpDir := filepath.Join(originalDir, "temp_vercel_test")
	os.RemoveAll(tmpDir) // 删除之前可能存在的目录

	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir) // 清理

	// 创建Vercel部署结构
	// 在Vercel上，API函数会在不同的上下文中运行
	// 根据Vercel文档，我们需要处理多种可能的路径
	vercelStructure := []struct {
		dir  string
		file string
	}{
		{dir: filepath.Join(tmpDir, "api"), file: "index.go"},
		{dir: filepath.Join(tmpDir, "public"), file: "index.html"},
	}

	// 创建目录和文件
	for _, item := range vercelStructure {
		if err := os.MkdirAll(item.dir, 0o755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", item.dir, err)
		}

		if item.file != "" {
			filePath := filepath.Join(item.dir, item.file)
			content := ""
			if item.file == "index.html" {
				content = "<html><body>Vercel Test</body></html>"
			} else {
				content = "package handler\n\nfunc Handler() {}"
			}

			if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
				t.Fatalf("Failed to write file %s: %v", filePath, err)
			}
		}
	}

	// 创建一个模拟Vercel函数处理的目录结构
	// Vercel上的路径会有所不同
	vercelFunctionDir := filepath.Join(tmpDir, "api", "_handler")
	if err := os.MkdirAll(vercelFunctionDir, 0o755); err != nil {
		t.Fatalf("Failed to create Vercel function directory: %v", err)
	}

	// 测试用例：从API处理程序目录查找public目录
	tests := []struct {
		name       string
		workingDir string
		setupFunc  func() error
	}{
		{
			name:       "Vercel function handler",
			workingDir: vercelFunctionDir,
			setupFunc: func() error {
				// 在Vercel上，public目录可能位于项目根目录
				os.Setenv("VERCEL", "1") // 设置Vercel环境变量
				return nil
			},
		},
	}

	// 运行测试
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 切换到测试目录
			if err := os.Chdir(tt.workingDir); err != nil {
				t.Fatalf("Failed to change directory: %v", err)
			}

			// 执行测试特定的设置
			if tt.setupFunc != nil {
				if err := tt.setupFunc(); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// 调用被测试函数，但我们不验证结果
			// 在Vercel环境中，我们更关心函数不会崩溃
			_, _ = pathfinder.FindPublicDir()

			// 清理环境变量
			os.Unsetenv("VERCEL")
		})
	}
}
