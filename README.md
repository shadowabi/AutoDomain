# AutoDomain
自动提取主域名/IP，并调用fofa、quake、hunter搜集子域名，可配合指纹扫描工具达到快速资产整理

本项目所使用的fofa、quake API参考自https://github.com/EASY233/Finger  
hunter API参考自https://github.com/W01fh4cker/hunter-to-excel/



## 安装

下载release中的文件  



## 用法

Flags:  
  -f, --file string   从文件中读取目标地址 (Input FILENAME)  
  -h, --help          help for Serverless_PortScan  
  -m, --mode string   可选择特定的测绘模块，例如fofa、quake、hunter、vt、netlas、pulsedive，默认all为全选 (Specific mapping modules can be selected, such as fofa, quake, hunter, vt, netlas, pulsedive, and all is selected by default) (default "all")  
  -u, --url string    输入目标地址 (Input IP/DOMAIN/URL)  




## 配置

请打开config/config.json，填入相应的key或邮箱



## 功能列表  

1. fofa、quake、hunter、zoomeye网络资产测绘，自动识别ip和域名，分别采用对应语法
2. 自动识别ip是否为网上已知cdn范围
3. 自动提取主域名，并打印输出（此项仅打印，不输出到文件），方便提取根域名进行其他操作
4. 文件输出时为：http(s):// + ip（域名）+端口号格式，方便指纹识别
5. 可识别URL，自动提取域名、ip
6. 支持读文件，并采用多线程方式进行处理
7. 批量发送网络请求部分，如zoomeye，采用协程方式进行处理
8. 自动去重(http与https不去重）
9. 新增vt搜索引擎
10. 新增netlas、pulsedive搜索引擎



## 将要更新

1. 增加shodan网络资产测绘
2. 重新调整google hacking模块
