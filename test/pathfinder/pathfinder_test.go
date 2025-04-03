package pathfinder_test

import (
	"os"
	"path/filepath"
	"testing"

	pathfinder "html2go-converter/api_utils"
)

// TestFindPublicDir 测试查找public目录的功能
func TestFindPublicDir(t *testing.T) {
	// 保存当前工作目录
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir) // 测试结束后恢复

	// 创建临时测试目录
	tmpDir := filepath.Join(originalDir, "temp_test")
	os.RemoveAll(tmpDir) // 删除之前可能存在的目录

	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir) // 清理

	// 创建测试目录结构
	testDirs := []string{
		filepath.Join(tmpDir, "public"),
		filepath.Join(tmpDir, "api"),
	}

	for _, dir := range testDirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// 创建测试文件
	indexPath := filepath.Join(tmpDir, "public", "index.html")
	if err := os.WriteFile(indexPath, []byte("<html><body>Test</body></html>"), 0o644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// 测试用例
	tests := []struct {
		name           string
		workingDir     string
		expectedResult bool
	}{
		{
			name:           "Base directory",
			workingDir:     tmpDir,
			expectedResult: true,
		},
		{
			name:           "API directory",
			workingDir:     filepath.Join(tmpDir, "api"),
			expectedResult: true,
		},
	}

	// 运行测试
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 切换到测试目录
			if err := os.Chdir(tt.workingDir); err != nil {
				t.Fatalf("Failed to change directory: %v", err)
			}

			// 调用被测试函数
			_, found := pathfinder.FindPublicDir()

			// 验证结果
			if found != tt.expectedResult {
				t.Errorf("FindPublicDir() in %s = %v, want %v", tt.workingDir, found, tt.expectedResult)
			}
		})
	}
}
