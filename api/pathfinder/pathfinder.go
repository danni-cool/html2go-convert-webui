package pathfinder

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Handler Vercel需要的HTTP处理函数
func Handler(w http.ResponseWriter, r *http.Request) {
	// 返回pathfinder功能的信息或状态
	result := map[string]interface{}{
		"status":  "ok",
		"message": "Pathfinder utility is working",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// isVercelEnvironment 检查是否在Vercel环境中运行
func isVercelEnvironment() bool {
	_, exists := os.LookupEnv("VERCEL")
	return exists
}

// FindPublicDir 在不同可能的路径中查找public目录
// 返回找到的public目录路径和一个表示是否找到的布尔值
func FindPublicDir() (string, bool) {
	// 输出一些调试信息
	cwd, _ := os.Getwd()
	log.Printf("Current working directory: %s", cwd)

	// 检查是否在Vercel环境中
	inVercel := isVercelEnvironment()
	if inVercel {
		log.Println("Running in Vercel environment")
	} else {
		log.Println("Running in local environment")
	}

	// 列出当前目录内容
	files, _ := filepath.Glob("*")
	log.Printf("Files in current directory: %v", files)

	// 在Vercel环境中，直接复制public目录到当前执行目录
	if inVercel {
		err := copyPublicToCurrentDir(cwd)
		if err != nil {
			log.Printf("Error copying public directory: %v", err)
		} else {
			log.Printf("Successfully copied public directory to current directory")
		}
	}

	// 根据环境构建查找路径
	publicPaths := buildPublicPathsList(inVercel)

	// 尝试所有可能的路径
	for _, dir := range publicPaths {
		// 获取绝对路径进行日志记录
		absPath, _ := filepath.Abs(dir)
		log.Printf("Checking public directory: %s (abs: %s)", dir, absPath)

		if _, err := os.Stat(dir); !os.IsNotExist(err) {
			// 目录存在，现在检查它是否包含index.html
			indexPath := filepath.Join(dir, "index.html")
			if _, err := os.Stat(indexPath); !os.IsNotExist(err) {
				log.Printf("Found public directory with index.html at: %s", absPath)
				return dir, true
			}

			// 如果目录存在但没有index.html，记录这个情况
			log.Printf("Directory exists but no index.html: %s", absPath)
		}
	}

	// 如果没有找到公共目录，尝试基于当前目录结构推断项目根目录
	projectRoot := inferProjectRoot(cwd)
	if projectRoot != "" {
		log.Printf("Inferred project root: %s", projectRoot)
		publicDir := filepath.Join(projectRoot, "public")

		// 检查这个推断的目录
		if _, err := os.Stat(publicDir); !os.IsNotExist(err) {
			indexPath := filepath.Join(publicDir, "index.html")
			if _, err := os.Stat(indexPath); !os.IsNotExist(err) {
				log.Printf("Found public directory using inferred project root: %s", publicDir)
				return publicDir, true
			}
		}
	}

	// 尝试查找项目根目录的所有父目录
	rootDir := findRootGoingUp(cwd, 10) // 向上查找最多10层
	if rootDir != "" {
		log.Printf("Found project root by going up: %s", rootDir)
		publicDir := filepath.Join(rootDir, "public")
		if _, err := os.Stat(publicDir); !os.IsNotExist(err) {
			indexPath := filepath.Join(publicDir, "index.html")
			if _, err := os.Stat(indexPath); !os.IsNotExist(err) {
				log.Printf("Found public directory after going up: %s", publicDir)
				return publicDir, true
			}
		}
	}

	// Vercel特定逻辑：检查固定位置的缓存目录结构
	if inVercel {
		// 尝试特定的Vercel目录结构模式
		possibleRoots := []string{
			"/var/task",
			"/home/site",
			"/vercel/path0",
		}

		for _, root := range possibleRoots {
			publicDir := filepath.Join(root, "public")
			if _, err := os.Stat(publicDir); !os.IsNotExist(err) {
				indexPath := filepath.Join(publicDir, "index.html")
				if _, err := os.Stat(indexPath); !os.IsNotExist(err) {
					log.Printf("Found public directory in Vercel specific path: %s", publicDir)
					return publicDir, true
				}
			}
		}

		// 尝试从当前目录向上查找的绝对路径
		vercelProjectPaths := findVercelProjectPaths(cwd)
		for _, projectPath := range vercelProjectPaths {
			publicDir := filepath.Join(projectPath, "public")
			if _, err := os.Stat(publicDir); !os.IsNotExist(err) {
				indexPath := filepath.Join(publicDir, "index.html")
				if _, err := os.Stat(indexPath); !os.IsNotExist(err) {
					log.Printf("Found public directory in Vercel project path: %s", publicDir)
					return publicDir, true
				}
			}
		}
	}

	// 如果没有找到公共目录，尝试列出更多信息以进行调试
	log.Printf("Could not find public directory in any expected location")

	// 尝试列出上级目录
	parentFiles, _ := filepath.Glob("../*")
	log.Printf("Files in parent directory: %v", parentFiles)

	// 在Vercel环境中，尝试查找更多可能的位置
	if inVercel {
		// 检查Vercel特定的环境变量
		log.Println("Vercel environment variables:")
		for _, env := range os.Environ() {
			if strings.HasPrefix(env, "VERCEL_") {
				log.Println(env)
			}
		}
	}

	// 尝试查找public目录的特殊情况
	if isVercelEnvironment() {
		// 这是最后的尝试：直接返回项目根目录中的public
		userHome := os.Getenv("HOME")
		if userHome != "" {
			// 尝试在系统中查找项目目录
			projectPath := filepath.Join(userHome, "github", "html2go-convert-webui", "public")
			if _, err := os.Stat(projectPath); !os.IsNotExist(err) {
				log.Printf("Found public directory in HOME path: %s", projectPath)
				return projectPath, true
			}
		}
	}

	// 返回空结果
	return "", false
}

// copyPublicToCurrentDir 将项目根目录的public文件夹复制到当前目录
func copyPublicToCurrentDir(currentDir string) error {
	// 尝试找到项目根目录
	rootDir := findRootGoingUp(currentDir, 10)
	if rootDir == "" {
		return nil // 没有找到根目录，不执行操作
	}

	// 源目录和目标目录
	srcDir := filepath.Join(rootDir, "public")
	dstDir := filepath.Join(currentDir, "public")

	// 检查源目录是否存在
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		return nil // 源目录不存在，不执行操作
	}

	// 创建目标目录（如果不存在）
	if _, err := os.Stat(dstDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dstDir, 0o755); err != nil {
			return err
		}
	}

	// 复制文件
	return copyDir(srcDir, dstDir)
}

// copyDir 将文件从源目录复制到目标目录
func copyDir(src, dst string) error {
	// 检查源目录
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !srcInfo.IsDir() {
		return nil // 如果不是目录，跳过
	}

	// 创建目标目录
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	// 读取源目录内容
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// 复制每个文件/目录
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// 递归复制子目录
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// 复制文件
			srcFile, err := os.ReadFile(srcPath)
			if err != nil {
				return err
			}

			if err := os.WriteFile(dstPath, srcFile, 0o644); err != nil {
				return err
			}
		}
	}

	return nil
}

// buildPublicPathsList 根据当前环境构建可能的public目录路径列表
func buildPublicPathsList(inVercel bool) []string {
	// 基本路径列表，适用于所有环境
	paths := []string{
		"./public",           // 本地开发环境相对路径
		"public",             // 直接相对路径
		"../public",          // 上一级目录
		"../../public",       // 上两级目录
		"../../../public",    // 上三级目录
		"../../../../public", // 上四级目录
	}

	// Vercel特定的路径
	vercelPaths := []string{
		"/var/task/public",     // Vercel 部署路径
		"/home/site/public",    // 另一个可能的 Vercel 路径
		"/vercel/path0/public", // 另一个可能的 Vercel 路径
		"/public",              // 根目录
	}

	// 在Vercel环境中，添加更多可能的路径
	if inVercel {
		// 获取VERCEL_PROJECT_PATH环境变量（如果有）
		if projectPath := os.Getenv("VERCEL_PROJECT_PATH"); projectPath != "" {
			vercelPaths = append(vercelPaths, filepath.Join(projectPath, "public"))
		}

		// 检查当前目录是否在.vercel/cache下，如果是，添加项目根目录路径
		currentDir, _ := os.Getwd()
		if strings.Contains(strings.ToLower(currentDir), ".vercel/cache") {
			// 尝试查找项目根目录
			for i := 1; i <= 10; i++ {
				path := ""
				for j := 0; j < i; j++ {
					path += "../"
				}
				vercelPaths = append(vercelPaths, filepath.Join(path, "public"))
			}
		}
	}

	// 根据环境添加额外路径
	if inVercel {
		// 在Vercel环境中，优先考虑Vercel特定路径
		return append(vercelPaths, paths...)
	}

	// 在本地环境中，可以添加一些Vercel路径作为备选
	return append(paths, vercelPaths...)
}

// inferProjectRoot 尝试从当前目录推断项目根目录
func inferProjectRoot(currentDir string) string {
	// 如果当前目录包含package.json或go.mod，可能是项目根目录
	if isRootDirectory(currentDir) {
		return currentDir
	}

	// 向上遍历目录树，寻找可能的项目根目录
	parent := filepath.Dir(currentDir)
	for parent != currentDir {
		if isRootDirectory(parent) {
			return parent
		}
		currentDir = parent
		parent = filepath.Dir(currentDir)
	}

	return ""
}

// isRootDirectory 检查给定目录是否可能是项目根目录
func isRootDirectory(dir string) bool {
	rootIndicators := []string{
		"package.json",
		"go.mod",
		"vercel.json",
	}

	for _, file := range rootIndicators {
		path := filepath.Join(dir, file)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			return true
		}
	}

	return false
}

// findRootGoingUp 通过向上搜索查找项目根目录
func findRootGoingUp(startDir string, maxLevels int) string {
	currentDir := startDir
	for i := 0; i < maxLevels; i++ {
		if isRootDirectory(currentDir) {
			return currentDir
		}
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			break // 已经到达根目录
		}
		currentDir = parentDir
	}
	return ""
}

// findVercelProjectPaths 查找可能的Vercel项目路径
func findVercelProjectPaths(currentDir string) []string {
	// 根据典型的Vercel部署结构构建可能的路径
	paths := []string{}

	// 检查是否在.vercel/cache下
	if strings.Contains(strings.ToLower(currentDir), ".vercel/cache") {
		// 向上查找可能的项目根目录
		parts := strings.Split(currentDir, string(filepath.Separator))
		for i := len(parts) - 1; i >= 0; i-- {
			// 尝试找到.vercel的索引
			if i > 0 && parts[i] == ".vercel" {
				// 构建到此索引的路径
				projectPath := filepath.Join(parts[:i]...)
				if projectPath == "" {
					projectPath = "/"
				} else if !strings.HasPrefix(projectPath, "/") {
					projectPath = "/" + projectPath
				}
				paths = append(paths, projectPath)
				break
			}
		}
	}

	// 添加工作目录和父目录
	cwd, _ := os.Getwd()
	paths = append(paths, cwd)

	// 向上添加多达5层父目录
	dir := cwd
	for i := 0; i < 5; i++ {
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		paths = append(paths, parent)
		dir = parent
	}

	return paths
}

// GetPublicIndexPath 获取public目录中index.html的完整路径
func GetPublicIndexPath() (string, bool) {
	publicDir, found := FindPublicDir()
	if !found {
		return "", false
	}

	// 构建index.html的路径
	indexPath := filepath.Join(publicDir, "index.html")
	return indexPath, true
}
