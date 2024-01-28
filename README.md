# XSpider

XSpider是一款高效的站点爬虫工具，采用深度遍历策略，持续发现并访问页面中的新链接。

# Screenshots
![xspider](png/xspider.png)

# 安装


# 使用指南
## 二进制文件运行
在[releases](https://github.com/Ithrael/XSpider/releases)下载二进制文件
```./xspider-v0.0.1-darwin-arm64 -url https://www.apple.com```
![help](png/xspider.png)

## 源码运行
```make help```
![help](png/help.png)

```make run ARGS='-file test/test_urls.txt'```
![run_file](png/run_file.png)

```make run ARGS='-url https://www.apple.com'```
![run_url](png/run_url.png)



