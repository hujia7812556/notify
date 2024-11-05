# Notify Service

一个简单而强大的消息推送服务，支持通过 WxPusher 发送消息。

## 功能特点

- 消息发送
  - WxPusher 消息推送
  - QPS 限制（最大2 QPS）
- 消息分发和限流
  - 工作池模式处理消息
  - 可配置的工作池大小
- HTTP API 接口
  - RESTful API 设计
  - Token 认证保护
  - JSON 格式数据交互
- 健康检查
  - 定时健康检查
  - 自定义检查时间
- 日志记录
  - 结构化日志
  - 可配置的日志级别和输出

## 快速开始

### 1. 安装

```bash
# 克隆项目
git clone https://github.com/yourusername/notify.git
cd notify

# 安装依赖
go mod tidy
```

### 2. 配置

```bash
# 复制配置模板
cp config/config.template.yaml etc/config.yaml

# 编辑配置
vim etc/config.yaml
```

配置示例：
```yaml
server:
  port: 8080
  mode: "debug"  # debug or release
  read_timeout: 5s
  write_timeout: 10s

wechat:
  sender_type: "wxpusher"
  wxpusher:
    app_token: "你的WxPusher应用Token"
    topic_ids: 
      - "你的主题ID"
    qps: 2

log:
  level: "info"
  format: "json"
  output: "stdout"

healthcheck:
  enabled: true
  check_time: "08:00"
  timeout: 10s
```

### 3. 运行

```bash
# 直接运行
go run main.go

# 或者编译后运行
go build -o bin/notify
./bin/notify
```

## API 使用

### 发送消息

```bash
curl -X POST http://localhost:8080/api/v1/notify \
  -H "Content-Type: application/json" \
  -d '{
    "platform": "wechat",
    "content": "测试消息",
    "summary": "消息摘要",
    "extra": {
      "user_id": "UID_xxx"
    }
  }'
```

### 健康检查

```bash
curl http://localhost:8080/api/v1/health
```

## 部署

### GitHub Actions 自动部署

需要配置的 Secrets：

```
# 服务器配置
SERVER_HOST: 服务器地址
SERVER_USERNAME: SSH用户名
SERVER_SSH_KEY: SSH私钥
SERVER_PORT: 服务端口
SERVER_TARGET: 部署目标路径

# WxPusher配置
WXPUSHER_APP_TOKEN: WxPusher应用Token
WXPUSHER_TOPIC_ID: WxPusher主题ID
```

### 服务器配置

```bash
# 创建用户和目录
sudo useradd -r -s /bin/false notify
sudo mkdir -p /etc/notify /var/log/notify
sudo chown -R notify:notify /etc/notify /var/log/notify

# 安装服务
sudo cp deploy/notify.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable notify
```

## 目录结构

```
notify/
├── bin/              # 二进制文件
├── config/           # 配置模板
├── etc/             # 本地配置
├── internal/        # 内部包
├── pkg/             # 公共包
└── deploy/          # 部署相关文件
```

## License

MIT License

## API 文档

### 认证方式

所有 `/api/v1/notify` 接口的调用都需要在请求头中携带 token：