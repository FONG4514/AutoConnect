package main

import (
	"errors"
	"flag"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

func main() {
	// 解析命令行参数
	arg1 := flag.String("a", "", "账户")
	arg2 := flag.String("p", "", "密码")
	arg3 := flag.String("t", "", "电信(telecom)/移动/联通")
	flag.Parse()

	if *arg1 == "" || *arg2 == "" || *arg3 == "" {
		log.Fatal("参数不足，请提供账户、密码和运营商类型")
	}

	for {
		if !isConnected() {
			log.Println("尝试重新连接校园网")
			req := loginToNetwork(*arg1, *arg2, *arg3)
			if isConnected() {
				log.Println("网络已连接")
			} else {
				if req != "" {
					log.Println("网络暂时未连接, req:", req)
				} else {
					log.Println("网络暂时未连接, req: 未知错误")
				}
			}
		} else {
			log.Println("网络未断开，请继续保持")
		}
		time.Sleep(500 * time.Millisecond)
	}
}

// 检查网络是否连接
func isConnected() bool {
	// 使用net包中的Dial函数尝试建立TCP连接
	conn, err := net.DialTimeout("tcp", "8.8.8.8:53", 1*time.Second)
	if err != nil {
		log.Println("网络断开:", err)
		return false
	}
	defer conn.Close()

	return true
}

// 执行登录操作
func loginToNetwork(account, password, networkType string) string {
	// 生成随机整数
	rand.Seed(time.Now().UnixNano())
	randomInt := rand.Intn(8889) + 1111 // 生成1111到9999之间的随机数

	// 获取JXUST-WLAN的IP地址
	wlanIP, err := getJXUSTWLANIP()
	if err != nil {
		log.Printf("获取JXUST-WLAN IP失败: %s\n", err)
		return ""
	}

	// 定义请求参数
	args := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(args)
	args.Add("callback", "dr1003")
	args.Add("login_method", "1")
	args.Add("user_account", account+"@"+networkType)
	args.Add("user_password", password)
	args.Add("wlan_user_ip", wlanIP)
	args.Add("wlan_user_ipv6", "")
	args.Add("wlan_user_mac", "000000000000")
	args.Add("wlan_ac_ip", "")
	args.Add("wlan_ac_name", "")
	args.Add("jsVersion", "4.1.3")
	args.Add("terminal_type", "1")
	args.Add("lang", "zh-cn")
	args.Add("v", strconv.Itoa(randomInt))
	args.Add("lang", "zh")

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
	err = fasthttp.Do(req, resp)
	if err != nil {
		log.Printf("请求失败: %s\n", err)
		return ""
	}

	return string(resp.Body())
}

// 获取JXUST-WLAN的IP地址
func getJXUSTWLANIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // 接口关闭
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // 忽略loopback接口
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // 不是IPv4地址
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("未找到JXUST-WLAN IP地址")
}
