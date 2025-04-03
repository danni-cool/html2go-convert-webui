package environment

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestSetup 测试设置和清理，为模拟Vercel环境做准备
func TestSetup(t *testing.T) {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	fmt.Printf("Current working directory: %s\n", currentDir)

	// 确保测试临时目录存在
	tmpDir := filepath.Join(currentDir, "../../temp_test")
	os.RemoveAll(tmpDir) // 删除之前可能存在的临时目录

	err = os.MkdirAll(tmpDir, 0o755)
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// 创建模拟目录结构
	createTestStructure(t, tmpDir)

	// 清理
	defer os.RemoveAll(tmpDir)
}

// 创建模拟测试结构
func createTestStructure(t *testing.T, baseDir string) {
	// 创建模拟Vercel部署环境
	vercelPaths := []string{
		"api",
		"public",
	}

	// 创建目录
	for _, path := range vercelPaths {
		fullPath := filepath.Join(baseDir, path)
		err := os.MkdirAll(fullPath, 0o755)
		if err != nil {
			t.Fatalf("Failed to create directory %s: %v", fullPath, err)
		}
	}

	// 创建示例文件
	createTestFile(t, filepath.Join(baseDir, "public", "index.html"), "<html><body>Test</body></html>")
	createTestFile(t, filepath.Join(baseDir, "api", "index.go"), "package handler\n\nfunc Handler() {}")

	// 验证目录结构
	validateStructure(t, baseDir)
}

// 创建测试文件
func createTestFile(t *testing.T, path string, content string) {
	err := os.WriteFile(path, []byte(content), 0o644)
	if err != nil {
		t.Fatalf("Failed to create file %s: %v", path, err)
	}
}

// 验证测试结构
func validateStructure(t *testing.T, baseDir string) {
	// 检查目录是否正确创建
	paths := []string{
		filepath.Join(baseDir, "api"),
		filepath.Join(baseDir, "public"),
		filepath.Join(baseDir, "public", "index.html"),
		filepath.Join(baseDir, "api", "index.go"),
	}

	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Fatalf("Expected path does not exist: %s", path)
		}
	}

	fmt.Println("Test structure validated successfully!")
}
