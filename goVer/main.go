package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

func main() {
	// 解析命令行参数
	arg1 := flag.String("a", "", "账户")
	arg2 := flag.String("p", "", "密码")
	arg3 := flag.String("t", "", "电信(telecom)/移动/联通")
	arg4 := flag.Int("i", 5000, "检查间隔时间（毫秒），默认为5000毫秒")
	arg5 := flag.String("l", "", "日志文件路径")
	flag.Parse()

	if *arg1 == "" || *arg2 == "" {
		log.Fatalf("\nUsage: ./app -a <username> -p <password> [-t <telecom/移动/联通>] [-i <interval_ms>] [-l <log_file_path>]")
	}

	// 如果提供了日志文件路径，则配置日志输出到文件
	if *arg5 != "" {
		// 创建日志目录（如果不存在）
		logDir := filepath.Dir(*arg5)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			log.Fatalf("无法创建日志目录: %v", err)
		}

		// 配置日志文件
		logFile, err := os.OpenFile(*arg5, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("无法打开日志文件: %v", err)
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	} else {
		// 未提供日志文件路径，日志输出到标准输出
		log.SetOutput(os.Stdout)
	}

	// 将间隔时间转换为 Duration 类型
	initialInterval := time.Duration(*arg4) * time.Millisecond
	currentInterval := initialInterval
	maxInterval := 5 * time.Minute // 最大间隔时间，防止等待时间过长

	for {
		if !isConnected(*arg3) {
			log.Println("网络断开，尝试重新连接校园网...")
			req := loginToNetwork(*arg1, *arg2, *arg3)
			if isConnected(*arg3) {
				log.Println("网络已连接")
				currentInterval = initialInterval // 重置间隔时间
			} else {
				if req != "" {
					log.Println("重新连接失败，返回信息:", req)
				} else {
					log.Println("网络暂时未连接，请求返回未知错误")
				}
				// 增加间隔时间，防止过于频繁的尝试
				if currentInterval < maxInterval {
					currentInterval *= 2
					if currentInterval > maxInterval {
						currentInterval = maxInterval
					}
				}
			}
		} else {
			// 如果网络已连接，确保间隔时间为初始值
			currentInterval = initialInterval
		}
		// 使用动态的间隔时间再检查
		time.Sleep(currentInterval)
	}
}

// 检查网络是否连接
func isConnected(networkType string) bool {
	if networkType == "" {
		// 建立 TCP 连接到校园网地址
		conn, err := net.DialTimeout("tcp", "10.16.10.50:80", 5*time.Second)
		if err != nil {
			// 无法连接，认为未连接校园网
			return false
		}
		defer conn.Close()

		// 设置读写超时时间
		conn.SetDeadline(time.Now().Add(5 * time.Second))

		// 发送 HTTP GET 请求
		request := "GET / HTTP/1.1\r\nHost: 10.16.10.50\r\nConnection: close\r\n\r\n"
		_, err = conn.Write([]byte(request))
		if err != nil {
			return false
		}

		// 读取响应内容
		responseBytes, err := io.ReadAll(conn)
		if err != nil {
			// 读取超时，认为已连接校园网
			return os.IsTimeout(err)
		}

		response := string(responseBytes)

		// 检查响应内容是否包含特定字段，未连接校园网时会返回该字段
		return !strings.Contains(response, "noNamespaceSchemaLocation")
	} else {
		// 检查互联网连接
		conn, err := net.DialTimeout("tcp", "8.8.8.8:53", 1*time.Second)
		if err != nil {
			return false
		}
		defer conn.Close()
		return true
	}
}

// 执行登录操作
func loginToNetwork(account, password, networkType string) string {
	// 定义请求参数
	args := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(args)
	if networkType != "" {
		// 连接到指定运营商
		args.Add("user_account", account+"@"+networkType)
	} else {
		// 只连接到校园网
		args.Add("user_account", account)
	}
	args.Add("user_password", password)

	url := "http://10.17.8.18:801/eportal/portal/login?" + args.String()

	// 创建HTTP请求
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodGet)

	// 创建HTTP响应
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// 发送HTTP GET请求
	err := fasthttp.Do(req, resp)
	if err != nil {
		log.Printf("请求失败: %s\n", err)
		return ""
	}

	return string(resp.Body())
}
