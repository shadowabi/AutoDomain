# AutoDomain
自动提取主域名/IP，并调用fofa、quake搜集子域名，可配合指纹扫描工具达到快速资产整理

本项目所使用的fofa、quake API参考自https://github.com/EASY233/Finger



## 安装

pip install -r requirements.txt



## 用法

python AutoDomain.py [-h] [-u www.baidu.com | -f 1.txt] [-m all]



## 配置

请打开AutoDomain.py文件，填写相应的fofa邮箱、key和quake的key
