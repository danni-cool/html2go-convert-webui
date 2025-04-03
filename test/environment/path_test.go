package environment

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// 测试在不同环境下查找public目录的功能
func TestFindPublicPath(t *testing.T) {
	// 保存当前工作目录
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir) // 测试结束后恢复工作目录

	// 测试用例：从不同目录寻找public目录
	testCases := []struct {
		name           string
		setupFunc      func() error // 设置环境的函数
		expectedResult bool         // 是否期望找到public目录
	}{
		{
			name: "Current directory contains public",
			setupFunc: func() error {
				// 当前目录下应该已经有public目录
				return nil
			},
			expectedResult: true,
		},
		{
			name: "One level up from api directory",
			setupFunc: func() error {
				// 创建临时的api目录并切换到那里
				err := os.MkdirAll("temp_test/api", 0o755)
				if err != nil {
					return err
				}
				return os.Chdir("temp_test/api")
			},
			expectedResult: true,
		},
		{
			name: "Vercel deployment simulation",
			setupFunc: func() error {
				// 创建模拟Vercel环境的目录结构
				// 在临时目录中创建类似Vercel部署的结构
				err := os.MkdirAll("temp_test/vercel/api", 0o755)
				if err != nil {
					return err
				}
				// 切换到api目录
				return os.Chdir("temp_test/vercel/api")
			},
			expectedResult: true,
		},
	}

	// 尝试在不同位置查找public目录的函数
	findPublicDir := func() (string, bool) {
		publicPaths := []string{
			"./public",             // 本地开发环境
			"public",               // 相对路径
			"../public",            // 上一级目录
			"../../public",         // 上两级目录
			"../../../public",      // 上三级目录
			"/var/task/public",     // Vercel 可能的路径
			"/home/site/public",    // 另一个可能的 Vercel 路径
			"/vercel/path0/public", // 另一个可能的 Vercel 路径
		}

		for _, dir := range publicPaths {
			absPath, _ := filepath.Abs(dir)
			fmt.Printf("Checking path: %s (abs: %s)\n", dir, absPath)

			if _, err := os.Stat(dir); !os.IsNotExist(err) {
				// 如果目录存在，确认它包含index.html
				testPath := filepath.Join(dir, "index.html")
				if _, err := os.Stat(testPath); !os.IsNotExist(err) {
					return dir, true
				}
				fmt.Printf("Directory exists but no index.html: %s\n", dir)
			}
		}
		return "", false
	}

	// 执行测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 清理之前测试可能创建的临时目录
			os.RemoveAll("temp_test")

			// 设置测试环境
			err := tc.setupFunc()
			if err != nil {
				t.Fatalf("Failed to setup test environment: %v", err)
			}

			// 输出当前工作目录，帮助调试
			currentDir, _ := os.Getwd()
			fmt.Printf("Current working directory: %s\n", currentDir)

			// 在本地尝试找到public目录
			path, found := findPublicDir()

			if found != tc.expectedResult {
				if tc.expectedResult {
					t.Errorf("Expected to find public directory, but didn't. Current dir: %s", currentDir)
				} else {
					t.Errorf("Expected to not find public directory, but found at: %s. Current dir: %s", path, currentDir)
				}
			} else if found {
				fmt.Printf("Found public directory at: %s\n", path)
			}
		})

		// 回到原始目录，以便下一个测试
		os.Chdir(originalDir)
		// 清理测试创建的临时目录
		os.RemoveAll("temp_test")
	}
}
