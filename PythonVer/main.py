#!/usr/bin/env python
# coding=utf-8
import subprocess
import logging
from time import sleep
import argparse
import requests
import random

# 配置日志
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

driver = None


def is_connected():
    try:
        subprocess.check_output("ping -n 1 -w 1000 81.71.3.20", shell=True, stderr=subprocess.STDOUT)
        return True
    except subprocess.CalledProcessError:
        logging.warning('网络断开')
        return False


def login_to_network(account, password, type):
    random_int = random.randint(1111, 9999)

    # 转化为字符串
    random_str = str(random_int)

    url = "http://10.17.8.18:801/eportal/portal/login"

    # 定义请求参数
    params = {
        "callback": "dr1003",
        "login_method": "1",
        "user_account": account + "@" + type,
        "user_password": password,
        "wlan_user_ip": "10.21.203.8",
        "wlan_user_ipv6": "",
        "wlan_user_mac": "000000000000",
        "wlan_ac_ip": "",
        "wlan_ac_name": "",
        "jsVersion": "4.1.3",
        "terminal_type": "1",
        "lang": "zh-cn",
        "v": random_str,
        "lang": "zh"
    }
    response = requests.get(url, params=params)
    return response


def main():
    parser = argparse.ArgumentParser(description="登录脚本")
    parser.add_argument('-a', '--arg1', type=str, required=True, help='账户')
    parser.add_argument('-p', '--arg2', type=str, required=True, help='密码')
    parser.add_argument('-t', '--arg3', type=str, required=True, help='电信(telecom)/移动/联通')

    args = parser.parse_args()
    while True:
        if not is_connected():
            logging.info('尝试重新连接校园网')
            try:
                req = login_to_network(args.arg1, args.arg2, args.arg3)
            finally:
                if is_connected():
                    logging.info('网络已连接')
                else:
                    if req is not None:
                        logging.info('网络暂时未连接,req: ' + req)
                    else:
                        logging.info("网络暂时未连接,req: 未知错误")
        else:
            logging.info("网络未断开，请继续保持")
        sleep(0.5)


if __name__ == "__main__":
    main()
