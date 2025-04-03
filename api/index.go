package handler

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	pathfinder "html2go-converter/api/_utils"
)

// Index function for serving the index.html file
func Index(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 使用pathfinder包查找index.html
	indexPath, found := pathfinder.GetPublicIndexPath()

	// 如果没有找到有效路径
	if !found {
		log.Printf("Error: Could not find public directory in any of the expected locations")

		// 列出当前目录以帮助调试
		files, _ := filepath.Glob("*")
		log.Printf("Files in current directory: %v", files)

		// 尝试列出上级目录
		parentFiles, _ := filepath.Glob("../*")
		log.Printf("Files in parent directory: %v", parentFiles)

		// 列出环境变量用于调试
		log.Println("Environment variables:")
		for _, env := range os.Environ() {
			log.Println(env)
		}

		// 显示详细错误
		http.Error(w, "Unable to find static files. Please ensure public directory exists.", http.StatusInternalServerError)
		return
	}

	// 读取 index.html 文件
	content, err := os.ReadFile(indexPath)
	if err != nil {
		log.Printf("Error reading index.html at %s: %v", indexPath, err)
		http.Error(w, "Unable to read index file", http.StatusInternalServerError)
		return
	}

	// 设置内容类型并提供文件
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}
