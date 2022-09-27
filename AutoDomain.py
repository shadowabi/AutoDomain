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
from time import sleep
import traceback
from config import *
import readline
import ipaddress
from queue import Queue
from threading import Thread
import urllib3
from bs4 import BeautifulSoup
from lxml import etree

urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

Drs = [] #存放主域名结果
Irs = [] #存放IP结果
rs2 = [] #存放资产测绘查询结果
keyword = "" #存放关键词
flag = 0 #区别IP和域名
zflag = 0 #为zoomeye区分域名和IP
zIp = [] #存放zoomeye特定url
dIp = [] #存放zoomeye特定url
q = Queue() #创建队列
modes = ["fofa","quake","hunter","zoomeye"]

ap = argparse.ArgumentParser()
group = ap.add_mutually_exclusive_group()
group.add_argument("-u", "--url", help = "Input IP/DOMAIN/URL", metavar = "www.baidu.com")
group.add_argument("-f", "--file", help = "Input FILENAME", metavar = "1.txt")
ap.add_argument("-m", "--mode", help = "Mode is fofa、quake、hunter、all", metavar = "all", default = "all")
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


def Scan(mode):
	global keyword,zflag
	_url = ""
	if mode == "fofa":
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

	if mode == "quake":
		data = {
            "query": keyword,
            "start": 0,
            "size": 100
        }
		try:
			response = s.post(url = "https://quake.360.net/api/v3/search/quake_service", headers = {"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4621.0 Safari/537.36","X-QuakeToken":qkey},json = data, timeout = 5)
			datas = json.loads(response.text)
			if len(datas['data']) >= 1 and datas['code'] == 0:
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

	if mode == "hunter":
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


	if mode == "zoomeye":
		grs = []
		header1 = {"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4621.0 Safari/537.36","API-KEY":zkey}
		if zflag - 2 >= 0:
			for i in Irs:
				zIp.append(("https://api.zoomeye.org/host/search?query=ip:{0}&facets=port").format(
				i))
			zflag -= 2
		if zflag - 1 >= 0:
			for i in Drs:
				dIp.append(("https://api.zoomeye.org/domain/search?q={0}&type=1").format(
				i))
			zflag -= 1
		try:
			if(len(zIp)):
				for i in zIp:
					grs.append(grequests.get(i, headers = header1, timeout = 5, verify = False))
			if (len(dIp)):
				for i in dIp:
					grs.append(grequests.get(i, headers = header1, timeout = 5, verify = False))

			for j in grequests.map(grs):
				if j != None and j.text != "null":
					datas = json.loads(j.text)
					if(datas['total'] != 0):
						for i in range(datas['total']):
							if("matches" in datas.keys()): #ip结果处理
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

							if ("list" in datas.keys()): #域名结果处理
								_url = "http://" + datas['list'][i]['name']
								if _url and _url not in rs2:
									rs2.append(_url.strip())

		except Exception as err:
			traceback.print_exc()

	# if mode == "google":
	# 	print(q.get())
	# 	url = "https://www.wuzhuiso.com/s?q=site:" + q.get()
	# 	grs = []
	# 	q2 = Queue()
	# 	for i in range(1,6):
	# 		url2 = url + "&pn=" + str(i)
	# 		try:
	# 			grs.append(grequests.get(url2, timeout = 5, verify = False))
	# 		except:
	# 			pass
	# 	for j in grequests.map(grs):

	# 		if j != None and j.text != "null":
	# 			q2.put(j.text)
	# 	def googlers():
	# 		while not q2.empty():
	# 			html = q2.get()
	# 			html = BeautifulSoup(html,'lxml')
	# 			html = str(html.select("cite"))
	# 			html = etree.HTML(text = html)
	# 			html = html.xpath('string(.)')
	# 			html = compile(r'\w{1,}\.\w{1,}\.\w{1,}').findall(html)
	# 			for _url in html:
	# 				if _url and _url not in rs2:
	# 					rs2.append(_url.strip())
	# 			q2.task_done()
	# 	for i in range(5):
	# 		t2 = Thread(target = googlers)
	# 		sleep(0.1)
	# 		t2.start()
	# 		t2.join(1)
	# 	q.task_done()

			
	keyword = ""

def Generate(mode):
	global keyword,zflag
	grammar = ""
	if mode in modes and mode != "zoomeye":
		if mode == "fofa" or mode == "hunter":
			grammar = "="
		elif mode == "quake":
			grammar = ":"

		if len(Drs):
			for i in Drs:
				keyword = keyword + "domain" + grammar + i.strip() + " || "
				q.put(i)
		if len(Irs):
			for i in Irs:
				keyword = keyword + "ip" + grammar + i.strip() + " || "
		keyword = keyword.rstrip(" || ")
		if keyword != "":
			Scan(mode)
		else:
			print("无效目标，退出程序！")
			_exit(0)

	elif mode == "zoomeye":
		if len(Drs):
			zflag +=1 #zflag为1时，target为域名
		if len(Irs):
			zflag +=2 #zflag>=2时，target存在IP
		if zflag == 0:
			print("无效目标，退出程序！")
			_exit(0)

		Scan(mode)

	else:
		print("参数错误！")
		_exit(0)





def Match(url):
	ips = ""
	ip = search(r"\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}(?:/\d{1,2}|)", url)
	# ip = search(r"\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}", url)
	if ip:
		if match(r"\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}/\d{1,2}", ip.group()):
			ips = ipaddress.ip_network(ip.group(),strict = False)
		if ips and ips not in Irs:
			good = IsCDN(ips, 2)
		elif ip and ip.group() not in Irs:
			good = IsCDN(ip.group())
		if good:
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
	if args.file:
		for i in open(target):
			Match(i.strip())
	else:
		Match(target)

	if mode == "all":
		for i in modes:
			t1 = Thread(target = Generate, args = [i])
			sleep(0.1)
			t1.start()
			t1.join()
	else:
		Generate(mode)

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

