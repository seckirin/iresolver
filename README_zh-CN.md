# iresolver

[README.md](./README.md) | [README_zh-CN.md](./README_zh-CN.md)

iresolver 是一个用 Go 语言编写的 DNS 解析器，它可以比较不同 DNS 服务器与基准 DNS 服务器的解析结果是否相同。

# 背景

原本我使用的是 [dnsvalidator](https://github.com/vortexau/dnsvalidator)，但是在我的机器上始终无法运行，而且在服务器上安装时遇到了 Python 版本不符的问题。因此，我决定参考 dnsvalidator 的功能，使用 Go 语言实现一个新的程序。

# 安装

可以在配置了 GOBIN 目录之后直接使用 go install 安装本程序。

```bash
go install github.com/yuukisec/iresolver/cmd/iresolver@latest
iresolver -h
```

也可以将项目克隆到本地后自行编译

```bash
git clone https://github.com/yuukisec/iresolver
cd iresolver
go build -o iresolver cmd/iresolver/main.go
./iresolver -h
```

# 功能

- **指定 DNS 服务器列表:** 你可以使用 `-target` 参数指定一个包含 DNS 服务器列表的文件或 URL。URL 的文件必须为纯文本形式。

- **静默模式:** 如果你只想输出成功解析的 DNS 服务器，你可以使用 `-silent` 参数。

- **指定基准 DNS 服务器:** 程序内置的默认基准 DNS 服务器是 1.1.1.1 和 8.8.8.8。你可以使用 `-dns` 参数来指定其他的基准 DNS 服务器。

- **指定基准域名:** 默认的基准域名是 qq.com, tencent.com 和 google.com。你可以通过 `-domain` 参数来指定其他的基准域名。

更多的参数说明请参考[完整参数](#完整参数)。

# 使用示例

**示例一:** 使用工具内置的基准 DNS 服务器和域名从 Public DNS 筛选可用的 DNS 服务器，指定线程为 200

```bash
iresolver -target https://public-dns.info/nameservers.txt -threads 200
```

**示例二:** 使用工具内置的基准 DNS 服务器和域名从 target.txt 文件中筛选可用的 DNS 服务器，指定线程为 200

```bash
iresolver -target target.txt -threads 200
```

**示例三:**

- 使用自定义基准 DNS 服务器和域名从 Public DNS 中筛选可用的 DNS 服务器
- 指定线程为 200、超时时间为 10 秒、重试次数为 3
- 当可用 DNS 服务器的数量达到 1000 时停止运行程序
- 将结果保存到 resolvers.txt 文件

```bash
iresolver \
  -target https://public-dns.info/nameservers.txt \
  -dns 1.1.1.1,8.8.8.8 \
  -domain qq.com,tencent.com \
  -threads 200 \
  -timeout 10 \
  -retry 3 \
  -count 1000
  -outptu resolvers.txt
```

# DNS 服务器列表

以下是一些可以用于测试的 DNS 服务器列表：

- https://public-dns.info/nameservers.txt
- https://raw.githubusercontent.com/trickest/resolvers/main/resolvers.txt

# 完整参数

```bash
Usage of iresolver
  -count int
    	Specify the program to stop running after writing how many qualified domain servers (default 65535)
  -dns string
    	Manually specify DNS servers, multiple servers are separated by commas (default "1.1.1.1,8.8.8.8")
  -domain string
    	Manually specify domains, multiple domains are separated by commas (default "qq.com,tencent.com")
  -output string
    	Specify the output file
  -retry int
    	Specify the number of retries when matching fails (default 2)
  -silent
    	Whether to only output successfully resolved DNS servers
  -target string
    	Specify the file or URL of the DNS server list to be checked
  -threads int
    	Specify the number of running threads (default 20)
  -timeout int
    	Specify the timeout for DNS requests (default 10)
```

# 其他说明

- 程序是多线程运行的，-count 参数指定的数量不一定完全准确。结果只会多不会少。
- 直接运行程序时将自动选择 [Public DNS](https://public-dns.info/nameservers.txt) 中的 DNS 服务器作为目标，可以通过 -target 参数修改