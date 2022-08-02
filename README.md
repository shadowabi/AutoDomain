# AutoDomain
自动提取主域名/IP，并调用fofa、quake、hunter搜集子域名，可配合指纹扫描工具达到快速资产整理

本项目所使用的fofa、quake API参考自https://github.com/EASY233/Finger  
hunter API参考自https://github.com/W01fh4cker/hunter-to-excel/



## 安装

pip install -r requirements.txt



## 用法

python AutoDomain.py [-h] [-u www.baidu.com | -f 1.txt] [-m all]  

**注意，本工具有一定的ip-cdn检测能力，但实战中大多数ip并没有被收录在网上公开的cdn范围中，所以强烈建议请勿将已知为CDN的IP放入该工具进行扫描，否则将会严重影响扫描结果**



## 配置

请打开config.py文件，在配置区域填写相关的key和邮箱等信息
