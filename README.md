# AutoDomain
自动提取主域名/IP，并调用fofa、quake、hunter搜集子域名，可配合指纹扫描工具达到快速资产整理

本项目所使用的fofa、quake API参考自https://github.com/EASY233/Finger  
hunter API参考自https://github.com/W01fh4cker/hunter-to-excel/


## 安装

下载release中的二进制文件使用

或使用Makefile进行编译二进制文件后使用


## 配置

当首次运行AutoDomain时，会检测config.json文件是否存在，不存在则会自动创建

config.json的填写内容应该如下：  
```
{
  "FofaMail": "",
  "FofaKey": "",
  "QuakeKey": "",
  "HunterKey": "",
  "ZoomeyeKey": "",
  "PulsediveKey": ""
}
```
填入的对应内容可使用对应的指定模块


## 用法

```
Usage:

  AutoDomain [flags]
  AutoDomain [command]


Available Commands:

  all          search domain from all engine
  fofa         search domain from fofa
  help         Help about any command
  hunter       search domain from hunter
  netlas       search domain from netlas
  pulsedive    search domain from pulsedive
  quake        search domain from quake
  virustotal   search domain from virustotal
  zoomeye      search domain from zoomeye


Flags:

  -f, --file string       从文件中读取目标地址 (Input FILENAME)
  -h, --help              help for AutoDomain
      --logLevel string   设置日志等级 (Set log level) [trace|debug|info|warn|error|fatal|panic] (default "info")
  -o, --output string     输入结果文件输出的位置 (Enter the location of the scan result output) (default "./result.txt")
  -t, --timeout int       输入每个 http 请求的超时时间 (Enter the timeout period for every http request) (default 15)
  -u, --url string        输入目标地址 (Input [ip|domain|url])


Use "AutoDomain [command] --help" for more information about a command.
```


## 功能列表  

1. 多种网络资产测绘，自动识别ip和域名，分别采用对应语法
2. 文件输出时为：http(s):// + ip（域名）+端口号格式，方便指纹识别
3. 可识别URL，自动提取域名、ip
4. 自动去重(http与https不去重）

## 旧版本

python版本在old分支

旧的go版本在go-old
