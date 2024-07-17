#!/usr/bin/env python
# coding=utf-8
import subprocess
import logging
from time import sleep
from selenium import webdriver

# 配置日志
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')


def is_connected():
    try:
        subprocess.check_output("ping -n 1 36.152.44.95", shell=True, stderr=subprocess.STDOUT)
        return True
    except subprocess.CalledProcessError:
        logging.warning('网络断开')
        return False


def setup_edge_driver():
    options = {
        "browserName": "MicrosoftEdge",
        "version": "",
        "platform": "WINDOWS",
        "ms:edgeOptions": {
            "extensions": [], "args": ["--headless", "--disable-gpu", "--window-size=1920x1080"]
        }
    }
    # 如果环境变量配置没问题就不用指定executable_path，否则需要指定
    driver = webdriver.Edge(executable_path='E:\\WebDriverForEdge\\msedgedriver.exe', capabilities=options)
    # driver = webdriver.Edge()
    return driver


def login_to_network(driver, retry_times):
    while retry_times > 0:
        try:
            driver.get('http://10.17.8.18/')
            ele = driver.find_elements_by_xpath('//input[@name="DDDDD"]')[1]
            ele.send_keys("替换为你的账号")
            ele = driver.find_elements_by_xpath('//input[@name="upass"]')[1]
            ele.send_keys("替换为你的密码")
            ele = driver.find_elements_by_xpath('//input[@name="network"]')[2]
            ele.click()
            ele = driver.find_elements_by_xpath('//input[@class="edit_lobo_cell"]')[1]
            ele.click()

            flag = is_connected()
            if flag is not True:
                retry_times -= 1
                logging.info("错误发生,接下来还会尝试连接:"+str(retry_times)+"次")
            else:
                break
        except Exception as e:
            pass


def main():
    while True:
        if not is_connected():
            logging.info('尝试重新连接校园网')
            driver = setup_edge_driver()
            try:
                login_to_network(driver, 3)
            finally:
                driver.quit()
                logging.info('网络已连接')
        logging.info("网络未断开，请继续保持")
        sleep(5)


if __name__ == "__main__":
    main()
