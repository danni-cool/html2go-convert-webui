#!/bin/bash
set -e # Exit immediately if a command exits with a non-zero status

# Color definitions for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 定义默认端口
SERVER_PORT=8080

# 日志函数
log_info() {
  echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
  echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
  echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
  echo -e "${RED}[ERROR]${NC} $1"
}

# 停止指定端口上运行的任何进程
stop_process_on_port() {
  local port=$1
  log_info "检查端口 $port 上的进程..."
  local pid=$(lsof -t -i:$port 2>/dev/null)

  if [ ! -z "$pid" ]; then
    log_warning "端口 $port 被进程 $pid 占用，尝试终止该进程..."
    kill -15 $pid 2>/dev/null || true
    sleep 1

    # 如果进程仍在运行，则强制终止
    if ps -p $pid >/dev/null 2>&1; then
      log_warning "进程未响应，强制终止..."
      kill -9 $pid 2>/dev/null || true
      sleep 1
    fi

    # 检查是否已终止
    if lsof -t -i:$port >/dev/null 2>&1; then
      log_error "无法释放端口 $port，尝试使用其他端口..."
      return 1
    else
      log_success "端口 $port 已成功释放"
      return 0
    fi
  else
    log_info "端口 $port 未被占用"
    return 0
  fi
}

# 清理函数 - 用于确保在脚本退出时关闭后台服务
cleanup() {
  log_info "清理资源..."
  if [ ! -z "$SERVER_PID" ]; then
    log_info "关闭后台运行的Go服务 (PID: $SERVER_PID)"
    kill -15 $SERVER_PID 2>/dev/null || kill -9 $SERVER_PID 2>/dev/null || true
    wait $SERVER_PID 2>/dev/null || true
  fi

  # 确保端口被释放
  stop_process_on_port $SERVER_PORT
}

# 设置 trap，确保清理函数在脚本退出时被调用
trap cleanup EXIT INT TERM

# 检查端口是否被占用，如果占用，则先尝试释放，如果不成功则寻找其他可用端口
check_port() {
  local port=$1

  # 首先尝试释放端口
  if lsof -i:$port >/dev/null 2>&1; then
    if stop_process_on_port $port; then
      return $port
    fi

    # 如果无法释放，尝试下一个端口
    while lsof -i:$port >/dev/null 2>&1; do
      log_warning "端口 $port 仍被占用，尝试下一个端口..."
      port=$((port + 1))
    done
  fi

  echo $port
}

# 获取可用端口
SERVER_PORT=$(check_port $SERVER_PORT)
log_info "将使用端口: $SERVER_PORT"

# 创建临时目录，用于存储PID文件等
TMP_DIR=$(mktemp -d)
log_info "创建临时目录: $TMP_DIR"
PID_FILE="$TMP_DIR/server.pid"

log_info "=== 开始API测试和服务初始化 ==="

# 运行后端测试(Go)
log_info "=== 运行后端Go API测试 ==="
if go test -v backend_test.go main.go; then
  log_success "后端测试通过!"
else
  log_error "后端测试失败! 请修复问题后再启动服务。"
  exit 1
fi

# 检查Node.js是否已安装
if ! command -v node &>/dev/null; then
  log_error "Node.js未安装。请安装Node.js以运行前端测试。"
  log_warning "跳过前端测试..."
else
  # 检查npm依赖是否已安装
  if [ ! -d "node_modules" ]; then
    log_info "安装npm依赖..."
    npm install
  fi

  # 检查test目录和测试文件是否存在
  if [ ! -f "test/frontend_test.js" ]; then
    log_error "前端测试文件不存在 (test/frontend_test.js)"
    log_info "创建测试目录..."
    mkdir -p test

    log_info "创建前端测试文件..."
    # 创建基本的测试文件内容
    cat >test/frontend_test.js <<'EOL'
// Frontend API Tests for HTML-Go Converter
const fetch = require('node-fetch');

// 从环境变量获取后端服务器的基础URL
const BASE_URL = process.env.SERVER_URL || 'http://localhost:8080';
console.log(`Using server at: ${BASE_URL}`);

describe('Frontend API Tests', () => {
  // Test HTML to Go conversion
  test('HTML to Go conversion API call succeeds', async () => {
    // 测试数据
    const testData = {
      html: '<div class="container"><h1 class="text-xl font-bold">Hello World</h1></div>',
      packagePrefix: 'h',
      vuetifyPrefix: 'v',
      vuetifyXPrefix: 'vx',
      direction: 'html2go'
    };

    // 发起实际请求到后端服务
    const response = await fetch(`${BASE_URL}/convert`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(testData),
    });

    // 验证响应
    expect(response.ok).toBeTruthy();
    expect(response.status).toBe(200);

    const data = await response.json();
    expect(data.code).toBeDefined();
    expect(data.error).toBeFalsy();

    // 验证生成的代码包含预期的元素
    expect(data.code).toContain('h.Div');
    expect(data.code).toContain('h.H1');
    expect(data.code).toContain('Class("container")');
    expect(data.code).toContain('Class("text-xl font-bold")');
  });
});
EOL
  fi

  # 确保端口未被占用
  stop_process_on_port $SERVER_PORT

  # 启动后端服务器以供前端测试使用
  log_info "启动后端服务以供前端测试使用..."
  go run main.go -port=$SERVER_PORT &
  SERVER_PID=$!
  echo $SERVER_PID >$PID_FILE
  log_info "后端服务已启动 (PID: $SERVER_PID) 在端口: $SERVER_PORT"

  # 等待服务器启动
  log_info "等待服务器完全启动..."
  sleep 3

  # 检查服务器是否正常运行
  MAX_RETRY=5
  RETRY=0
  while ! curl -s http://localhost:$SERVER_PORT/ >/dev/null && [ $RETRY -lt $MAX_RETRY ]; do
    log_warning "服务器尚未就绪，继续等待..."
    sleep 2
    RETRY=$((RETRY + 1))
  done

  if curl -s http://localhost:$SERVER_PORT/ >/dev/null; then
    log_success "服务器已成功启动，运行在 http://localhost:$SERVER_PORT/"

    # 设置环境变量传递给测试程序
    export SERVER_URL="http://localhost:$SERVER_PORT"

    # 运行前端测试
    log_info "=== 运行前端API测试 ==="
    if npm run test:frontend; then
      log_success "前端测试通过!"

      # 运行前缀测试
      log_info "=== 运行前缀专项测试 ==="
      if npm run test:prefix; then
        log_success "前缀测试通过! 确认Vuetify前缀为'v'，VuetifyX前缀为'vx'"
      else
        log_error "前缀测试失败! 请确保Vuetify前缀为'v'，VuetifyX前缀为'vx'"
        # 关闭测试服务器
        log_info "前缀测试失败，关闭测试服务器..."
        if [ ! -z "$SERVER_PID" ] && ps -p $SERVER_PID >/dev/null; then
          kill -15 $SERVER_PID && sleep 1
          # 确保进程已终止
          if ps -p $SERVER_PID >/dev/null; then
            kill -9 $SERVER_PID
          fi
          unset SERVER_PID
        fi

        # 确保端口被释放
        stop_process_on_port $SERVER_PORT

        exit 1
      fi

      # 关闭测试服务器
      log_info "前端测试完成，关闭测试服务器..."
      if [ ! -z "$SERVER_PID" ] && ps -p $SERVER_PID >/dev/null; then
        kill -15 $SERVER_PID && sleep 1
        # 确保进程已终止
        if ps -p $SERVER_PID >/dev/null; then
          kill -9 $SERVER_PID
        fi
        unset SERVER_PID
      fi

      # 确保端口被释放
      stop_process_on_port $SERVER_PORT

      # 重启服务器用于实际使用
      log_success "所有测试通过! 正在启动Go服务..."
      exec go run main.go -port=$SERVER_PORT
    else
      log_error "前端测试失败! 请修复前端代码中的问题。"
      # 在这里可以添加自动修复前端代码的逻辑
      log_info "关闭测试服务器..."
      if [ ! -z "$SERVER_PID" ] && ps -p $SERVER_PID >/dev/null; then
        kill -15 $SERVER_PID && sleep 1
        # 确保进程已终止
        if ps -p $SERVER_PID >/dev/null; then
          kill -9 $SERVER_PID
        fi
        unset SERVER_PID
      fi

      # 确保端口被释放
      stop_process_on_port $SERVER_PORT

      exit 1
    fi
  else
    log_error "服务器启动失败，无法连接到 http://localhost:$SERVER_PORT/"
    log_info "关闭测试服务器..."
    if [ ! -z "$SERVER_PID" ] && ps -p $SERVER_PID >/dev/null; then
      kill -15 $SERVER_PID && sleep 1
      # 确保进程已终止
      if ps -p $SERVER_PID >/dev/null; then
        kill -9 $SERVER_PID
      fi
      unset SERVER_PID
    fi

    # 确保端口被释放
    stop_process_on_port $SERVER_PORT

    exit 1
  fi
fi

# 如果代码执行到这里，表示没有Node.js，直接启动Go服务
log_info "直接启动Go服务..."
exec go run main.go -port=$SERVER_PORT
