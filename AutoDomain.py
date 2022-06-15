#!/usr/bin/env python
# -*- coding:utf-8 -*-
#author:Sh4d0w_小白

import argparse
from re import search, sub
from os import system, _exit
import requests
import json
import base64
from urllib.parse import quote
from time import sleep

rs = [] #存放主域名结果
rs2 = [] #存放fofa、quake查询结果
keyword = "" #存放关键词
ap = argparse.ArgumentParser()
group = ap.add_mutually_exclusive_group()
group.add_argument("-u", "--url", help = "Input IP/DOMAIN/URL", metavar = "www.baidu.com")
group.add_argument("-f", "--file", help = "Input FILENAME", metavar = "1.txt")
ap.add_argument("-m", "--mode", help = "Mode is fofa、quake、hunter、all", metavar = "all", default = "all")


#配置#
fmail = "" #填写fofa邮箱
fkey = "" #填写fofa key
qkey = "" #填写quake key
hkey = "" #填写hunter key
header = {
	"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4621.0 Safari/537.36"
}	

def Scan(mode):
	global keyword
	_url = ""
	if mode == "fofa":
		keyword = base64.urlsafe_b64encode(keyword.encode()).decode()
		url = "https://fofa.info/api/v1/search/all?email={0}&key={1}&qbase64={2}&full=false&fields=protocol,host&size=100".format(
		fmail, fkey, keyword)
		try:
			response = requests.get(url, timeout = 10, headers = header )
			datas = json.loads(response.text)
			if "results" in datas.keys():
				for data in datas["results"]:
					if "http" in data[1] or "https" in data[1]:
						_url = data[1]
					elif "http" == data[0] or "https" == data[0]:
						_url = "{0}://{1}".format(data[0], data[1])
					elif "" == data[0]:
						_url = "{0}://{1}".format("http", data[1])
					if _url and _url not in rs2:
						rs2.append(_url)
		except Exception as err:
			print(err)

	if mode == "quake":
		data = {
            "query": keyword,
            "start": 0,
            "size": 100
        }
		try:
			response = requests.post(url = "https://quake.360.cn/api/v3/search/quake_service", headers = {"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4621.0 Safari/537.36","X-QuakeToken":qkey},json = data, timeout = 10)
			datas = json.loads(response.text)
			if len(datas['data']) >= 1 and datas['code'] == 0:
				for data in datas['data']:
					port = "" if data['port'] == 80 or data["port"] == 443 else ":{}".format(str(data['port']))
					if 'http/ssl' == data['service']['name']:
						_url = 'https://' + data['service']['http']['host'] + port
					elif 'http' == data['service']['name']:
						_url = 'http://' + data['service']['http']['host'] + port
					if _url and _url not in rs2:
						rs2.append(_url)
		except Exception as err:
			print(err)

	if mode == "hunter":
		keyword = base64.urlsafe_b64encode(keyword.encode()).decode()
		url = "https://hunter.qianxin.com/openApi/search?api-key={0}&search={1}&page=1&page_size=100&is_web=3".format(
		hkey, keyword)
		try:
			response = requests.get(url, timeout = 10, headers = header)
			datas = json.loads(response.text)
			for i in  range(len(datas["data"]["arr"])):
				_url = datas["data"]["arr"][i]["url"]
				if _url and _url not in rs2:
						rs2.append(_url)
		except Exception as err:
			print(err)
			
	keyword = ""

def Generate(mode):
	global keyword
	if mode == "fofa" or mode == "hunter":
		grammar = "domain="
	elif mode == "quake":
		grammar = "domain:"
	else:
		print("参数错误！")
		_exit(0)

	for i in rs:
		keyword = keyword + grammar + i + " || "
	keyword = keyword.rstrip(" || ")

	Scan(mode)

def Match(url):
	ip = search(r"\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}", url)
	if ip and ip.group() not in rs:
		rs.append(ip.group())


	
	if(search(r"(http|https)\:\/\/", url)): # 当输入URL时提取出域名
	    url = sub(r"(http|https)\:\/\/", "", url)
	    if (search(r"(\/|\\).*", url)):
	        url = sub(r"(\/|\\).*", "", url)
	domain = search(r"^([a-zA-Z0-9]([a-zA-Z0-9-_]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,11}$", url) #检测是否为域名
	if domain:
		domain = search(r"([a-z0-9][a-z0-9\-]*?\.(?:com|cn|net|org|gov|info|la|cc|co|jp|net|edu|org|top|tk)(?:\.(?:cn|jp))?)$", domain[0])
		if domain and domain[0] not in rs:
			rs.append(domain[0])

if __name__ == '__main__':
	args = ap.parse_args()
	target = args.url or args.file
	mode = args.mode

	print("开始进行扫描：")
	if args.file:
		for i in open(target):
			Match(i.strip())
	else:
		Match(target)

	if mode == "all":
		Generate("fofa")
		sleep(0.1)
		Generate("quake")
		sleep(0.1)
		Generate("hunter")
	else:
		Generate(mode)

	with open("result.txt","w+",encoding='utf8') as f:
		print("主域名：")
		for i in rs:
			f.write(i + "\n")
			print(i)
		print("子域名：")	
		for i in rs2:
			f.write(i + "\n")
			print(i)

	print("扫描结束！")
	print("已保存到result.txt文件中，按回车键退出程序！")
	input()
	_exit(0)

