## 这个是酱梨专用的校园网连接脚本

### GoVer使用说明
- 环境要求：
  - Go 1.21.3 或以上版本
  - fasthttp v1.55.0

### 功能特点
- 自动连接校园网
- 支持断线重连
- 快速响应，低资源占用
- 支持命令行参数配置

### 使用方法
1. 克隆仓库到本地
2. 在命令行中运行 `go run main.go -a <username> -p <password> -t <telecom/移动/联通> -i <interval_ms> -l <log_file_path>`
