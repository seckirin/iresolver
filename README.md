# iresolver

[README.md](./README.md) | [README_zh-CN.md](./README_zh-CN.md)

iresolver is a DNS resolver written in Go language. It can compare whether the resolution results of different DNS servers are the same as those of the baseline DNS servers.

# Background

Originally, I was using [dnsvalidator](https://github.com/vortexau/dnsvalidator), but it could never run on my machine, and I encountered problems with the Python version when installing it on the server. Therefore, I decided to refer to the functions of dnsvalidator and implement a new program in Go language.

# Installation

You can directly install this program using go install after configuring the GOBIN directory.

```bash
go install github.com/yuukisec/iresolver/cmd/iresolver@latest
iresolver -h
```

You can also clone the project to your local machine and compile it yourself.

```bash
git clone https://github.com/yuukisec/iresolver
cd iresolver
go build -o iresolver cmd/iresolver/main.go
./iresolver -h
```

# Features

- **Specify DNS Server List:** You can use the -target parameter to specify a file or URL containing a list of DNS servers. The file from the URL must be in plain text format.
- **Silent Mode:** If you only want to output the DNS servers that have been successfully resolved, you can use the -silent parameter.
- **Specify Baseline DNS Servers:** The program’s built-in default baseline DNS servers are 1.1.1.1 and 8.8.8.8. You can use the -dns parameter to specify other baseline DNS servers.
- **Specify Baseline Domains:** The default baseline domains are qq.com and tencent.com. You can use the -domain parameter to specify other baseline domains.

For more parameter explanations, please refer to [Complete Parameters](#complete-parameters).

# Examples

**Example 1:** Use the built-in baseline DNS servers and domains of the tool to filter available DNS servers from Public DNS, specify threads as 200

```bash
iresolver -target https://public-dns.info/nameservers.txt -threads 200
```

**Example 2:** Use the built-in baseline DNS servers and domains of the tool to filter available DNS servers from the target.txt file

```bash
iresolver -target target.txt -threads 200
```

**Example 3:**

- Use custom baseline DNS servers and domains to filter available DNS servers from Public DNS
- Specify threads as 200, timeout as 10 seconds, retry times as 3
- Stop running the program when the number of available DNS servers reaches 1000
- Save the results to the resolvers.txt file

```bash
iresolver \
  -target https://public-dns.info/nameservers.txt \
  -dns 1.1.1.1,8.8.8.8 \
  -domain qq.com,tencent.com \
  -threads 200 \
  - timeout 10 \
  -retry 3 \
  -count 20
```

# DNS Server List

Here are some DNS server lists that can be used for testing:

- https://public-dns.info/nameservers.txt
- https://raw.githubusercontent.com/trickest/resolvers/main/resolvers.txt

# Complete Parameters

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
    	Specify the number of running threads (default 5)
  -timeout int
    	Specify the timeout for DNS requests (default 3)
```

# Additional Notes

- The program runs in multithreaded mode. The quantity specified by the -count parameter may not be completely accurate. The results will only be more, not less.