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



## 功能列表  

1. fofa、quake、hunter、zoomeye网络资产测绘，自动识别ip和域名，分别采用对应语法
2. 自动识别ip是否为网上已知cdn范围
3. 自动提取主域名，并打印输出（不在文件输出中），方便提取根域名进行其他操作
4. 文件输出时为：http(s):// + ip（域名）+端口号格式，方便指纹识别
5. 可识别URL，自动提取域名、ip
6. 支持读文件，并采用多线程方式进行处理
7. 批量发送网络请求部分，如zoomeye，采用协程方式进行处理
8. 自动去重(http与https不去重）
9. 新增vt搜索引擎



## 将要更新

1. 增加shodan网络资产测绘
2. 重新调整google hacking模块
3. 采用go语言重构本程序
