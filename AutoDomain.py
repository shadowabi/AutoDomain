#!/usr/bin/env python
# -*- coding:utf-8 -*-

import grequests
import argparse
from re import *
from os import system, _exit
import requests
from requests.adapters import HTTPAdapter
import json
import base64
from urllib.parse import quote
import traceback
from config import *
import readline
import ipaddress
from concurrent.futures import ThreadPoolExecutor
import urllib3
from bs4 import BeautifulSoup
from lxml import etree
from time import sleep

urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

Drs = [] #存放主域名结果
Irs = [] #存放IP结果
rs2 = [] #存放资产测绘查询结果
keyword = "" #存放关键词
flag = 0 #区别IP和域名
modes = ["fofa","quake","hunter","zoomeye","vt","netlas","pulsedive"]

ap = argparse.ArgumentParser()
group = ap.add_mutually_exclusive_group()
group.add_argument("-u", "--url", help = "Input IP/DOMAIN/URL", metavar = "www.baidu.com")
group.add_argument("-f", "--file", help = "Input FILENAME", metavar = "1.txt")
ap.add_argument("-m", "--mode", help = "Mode is fofa、quake、hunter、vt、netlas、pulsedive、all", metavar = "all", default = "all")
s = requests.Session()
s.mount('http://', HTTPAdapter(max_retries=5))
s.mount('https://', HTTPAdapter(max_retries=5))

header = {
	"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4621.0 Safari/537.36"
}


def IsCDN(ip,flag = 1):
	realip = 1
	with open("cdn_ip_cidr.json", 'r', encoding='utf-8') as f:
		cdns = json.load(f)
		if flag == 1:
			for cdn in cdns:
				if ipaddress.ip_address(ip) in ipaddress.ip_network(cdn, strict = False):
					realip = 0
		else:
			if ip in cdns:
				realip = 0
	if realip == 1:
		return str(ip)

def fofa():
	global keyword
	print("[+]fofa is working...")
	_url = ""
	keyword = base64.urlsafe_b64encode(keyword.encode()).decode()
	url = "https://fofa.info/api/v1/search/all?email={0}&key={1}&qbase64={2}&full=false&fields=protocol,host&size=1000".format(
	fmail, fkey, keyword)

	try:
		response = s.get(url, timeout = 5, headers = header )
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
					rs2.append(_url.strip())
	except Exception as err:
		traceback.print_exc()

def quake():
	global keyword
	_url = ""
	print("[+]quake is working...")
	data = {
        "query": keyword,
        "start": 0,
        "size": 100
    }
	header1 = {"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4621.0 Safari/537.36","X-QuakeToken":qkey}
	
	try:
		response = s.post(url = "https://quake.360.net/api/v3/search/quake_service", headers = header1, json = data, timeout = 5)

		datas = json.loads(response.text)
		if "data" in datas.keys() and datas['code'] == 0:
			for data in datas['data']:
				port = "" if data['port'] == 80 or data["port"] == 443 else ":{}".format(str(data['port']))
				if 'http/ssl' == data['service']['name'] and 'http' in data['service']:
					_url = 'https://' + data['service']['http']['host'] + port
				elif 'http' == data['service']['name']:
					_url = 'http://' + data['service']['http']['host'] + port
				if _url and _url not in rs2:
					rs2.append(_url.strip())

	except Exception as err:
		# traceback.print_exc()
		pass	

def hunter():
	global keyword
	print("[+]hunter is working...")
	_url = ""
	keyword = base64.urlsafe_b64encode(keyword.encode()).decode()
	url = "https://hunter.qianxin.com/openApi/search?api-key={0}&search={1}&page=1&page_size=100&is_web=3".format(
	hkey, keyword)

	try:
		response = s.get(url, timeout = 5, headers = header)
		datas = json.loads(response.text)
		if datas["data"]["arr"]:
			for i in range(len(datas["data"]["arr"])):
				_url = datas["data"]["arr"][i]["url"]
				if _url and _url not in rs2:
						rs2.append(_url.strip())

	except Exception as err:
		traceback.print_exc()

def zoomeye():
	print("[+]zoomeye is working...")
	_url = ""
	grs = [] #存放异步处理结果
	gIp = [] #存放需要异步请求的ip
	gDm = [] #存放需要异步请求的域名
	header1 = {"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4621.0 Safari/537.36","API-KEY":zkey}

	if len(Irs):
		for i in Irs:
			gIp.append(("https://api.zoomeye.org/host/search?query=ip:{0}&facets=port").format(i))
	if len(Drs):
		for i in Drs:
			gDm.append(("https://api.zoomeye.org/domain/search?q={0}&type=1").format(i))

	try:
		if len(gIp):
			for i in gIp:
				grs.append(grequests.get(i, headers = header1, timeout = 5, verify = False))
		if len(gDm):
			for i in gDm:
				grs.append(grequests.get(i, headers = header1, timeout = 5, verify = False))

		for j in grequests.map(grs):
			if j != None and j.text != "null" and j.status_code == 200:
				datas = json.loads(j.text)
				if datas['total'] != 0:
					for i in range(datas['total']):
						if "matches" in datas.keys(): #ip结果处理
							port = str(datas['matches'][i]['portinfo']['port'])
							if datas['matches'][i]['portinfo']['service'] == "http" or datas['matches'][i]['portinfo']['service'] == "https":
								protocol = datas['matches'][i]['portinfo']['service'] 
							else: 
								protocol = ""
							if protocol != "":
								_url = protocol + "://" + datas['matches'][0]['ip'] + ":" + port
							else:
								_url = "http://" + datas['matches'][0]['ip'] + ":" + port
							if _url and _url not in rs2:
								rs2.append(_url.strip())

						if "list" in datas.keys(): #域名结果处理
							_url = "http://" + datas['list'][i]['name']
							if _url and _url not in rs2:
								rs2.append(_url.strip())

	except Exception as err:
		traceback.print_exc()

def virustotal():
	print("[+]virustotal is working...")
	_url = ""
	grs = [] #存放异步处理结果
	gDm = [] #存放需要异步请求的域名
	header1 = {"X-Vt-Anti-Abuse-Header":"1","User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4621.0 Safari/537.36","Accept-Ianguage":"en-US,en;q=0.9,es;q=0.8","X-Tool":"vt-ui-main"}

	for i in Drs:
		gDm.append(("https://www.virustotal.com/ui/domains/{0}/subdomains?relationships=resolutions&cursor=eyJsaW1pdCI6IDIwMCwgIm9mZnNldCI6IDB9&limit=200").format(i))

	try:
		if len(gDm):
			for i in gDm:
				grs.append(grequests.get(i, headers = header1, timeout = 10, verify = False))

		for j in grequests.map(grs):
			if j != None and j.text != "null" and j.status_code == 200:
				datas = json.loads(j.text)
				if datas["data"]:
					for i in range(len(datas["data"])):
						_url = "http://" + datas["data"][i]["id"]
						if _url and _url not in rs2:
							rs2.append(_url.strip())
	except Exception as err:
		traceback.print_exc()	

def netlas():
	print("[+]netlas is working...")
	_url = ""
	grs = [] #存放异步处理结果
	gIp = [] #存放需要异步请求的ip
	gDm = [] #存放需要异步请求的域名
	if len(Irs):
		for i in Irs:
			gIp.append(("https://app.netlas.io/api/host/{0}/?source_type=include&fields=related_domains").format(i))
	if len(Drs):
		for i in Drs:
			gDm.append(("https://app.netlas.io/api/domains/?q=*.{0}&source_type=include&fields=domain").format(i))
			gDm.append(("https://app.netlas.io/api/host/{0}/?source_type=include&fields=related_domains").format(i))

			try:
				if len(gIp):
					for i in gIp:
						grs.append(grequests.get(i, headers = header, timeout = 10, verify = False))
						sleep(1)
				if len(gDm):
					for i in gDm:
						grs.append(grequests.get(i, headers = header, timeout = 10, verify = False))
						sleep(1)

				for j in grequests.map(grs):
					if j != None and j.text != "null" and j.status_code == 200:		
						datas = json.loads(j.text)
						if "related_domains" in datas.keys():
							for i in range(len(datas['related_domains'])):
								_url = "http://" + datas['related_domains'][i]
								if _url and _url not in rs2:
									rs2.append(_url.strip())
						if "items" in datas.keys():
							for i in range(len(datas["items"])):
								_url = "http://" + datas["items"][i]['data']['domain']
								if _url and _url not in rs2:
									rs2.append(_url.strip())

			except Exception as err:
				traceback.print_exc()

def pulsedive():
	print("[+]pulsedive is working...")
	_url = ""
	grs = [] #存放异步处理结果
	gDm = [] #存放需要异步请求的域名

	for i in Drs:
		gDm.append(("https://pulsedive.com/api/explore.php?q=ioc%3d*.{0}%20active%3dtrue").format(i))

	try:
		if len(gDm):
			for i in gDm:
				grs.append(grequests.get(i, headers = header, timeout = 10, verify = False))

		for j in grequests.map(grs):
			if j != None and j.text != "null" and j.status_code == 200:
				datas = json.loads(j.text)
				if "results" in datas.keys():
					for i in range(len(datas["results"])):
						_url = "http://" + datas["results"][i]["indicator"]
						if _url and _url not in rs2:
							rs2.append(_url.strip())
	except Exception as err:
		traceback.print_exc()		


def Scan(mode):
	global keyword
	_url = ""
	if mode == "fofa":
		fofa()

	if mode == "quake":
		quake()

	if mode == "hunter":
		hunter()

	if mode == "zoomeye":
		zoomeye()

	if mode == "vt":
		virustotal()

	if mode == "netlas":
		netlas()

	if mode == "pulsedive":
		pulsedive()

	keyword = "" #清空变量

def Generate(mode):
	global keyword
	grammar = ""
	if mode in modes:
		if mode == "fofa" or mode == "hunter":
			grammar = "="
		elif mode == "quake":
			grammar = ":"
		if grammar != "":
			if len(Drs):
				for i in Drs:
					keyword = keyword + "domain" + grammar + i.strip() + " || "
			if len(Irs):
				for i in Irs:
					keyword = keyword + "ip" + grammar + i.strip() + " || "
			keyword = keyword.rstrip(" || ")

		if len(Drs) == 0 and len(Irs) == 0:
			print("无效目标，退出程序！")
			_exit(0)

		Scan(mode)

	else:
		print("参数错误！")
		_exit(0)

def Match(url):
	ips = ""
	good = None
	ip = search(r"\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}(?:/\d{1,2}|)", url)
	if ip:
		if match(r"\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}/\d{1,2}", ip.group()):
			ips = ipaddress.ip_network(ip.group(),strict = False)
		if ips and ips not in Irs:
			good = IsCDN(ips, 2)
		elif ip and ip.group() not in Irs:
			good = IsCDN(ip.group())
		if good != None:
			Irs.append(good)
	
	if(search(r"(http|https)\:\/\/", url)): # 当输入URL时提取出域名
	    url = sub(r"(http|https)\:\/\/", "", url)
	    if (search(r"(\/|\\).*", url)):
	        url = sub(r"(\/|\\).*", "", url)
	domain = search(r"^([a-zA-Z0-9]([a-zA-Z0-9-_]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,11}$", url) #检测是否为域名
	if domain:
		domain = search(r"([a-z0-9][a-z0-9\-]*?\.(?:\w{2,4})(?:\.(?:cn|hk))?)$", domain[0])
		if domain and domain[0] not in Drs:
			Drs.append(domain[0])

if __name__ == '__main__':
	args = ap.parse_args()
	target = args.url or args.file
	mode = args.mode

	print("开始进行扫描：")
	with ThreadPoolExecutor(max_workers = 20) as executor:
		if args.file:
			with open(target) as f:
				lines = filter(lambda x: x.strip(), f)
				executor.map(Match, lines)
		else:
			executor.submit(Match, target.strip())

	with ThreadPoolExecutor(max_workers = 20) as executor:
		if mode == "all":
			executor.map(Generate, modes)
		else:
			executor.submit(Generate, mode)

	with open("result.txt", "w+", encoding = 'utf8') as f:
		print("主域名和IP：")
		for i in Drs:
			print(i.strip())
		for i in Irs:
			print(i.strip())
		print("子域名：")	
		for i in rs2:
			f.write(i.strip() + "\n")
			print(i.strip())

	print("扫描结束！")
	print("已保存到result.txt文件中，按回车键退出程序！")
	input()
	_exit(0)