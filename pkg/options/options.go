// pkg/options/options.go
package options

import (
	"flag"
	"github.com/yuukisec/iresolver/pkg/except"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type Options struct {
	Target  string
	Output  string
	Silent  bool
	Dns     []string
	Domain  []string
	Threads int
	Timeout time.Duration
	Retry   int
	Count   int
}

func ParseOptions() Options {
	target := flag.String("target", "https://public-dns.info/nameservers.txt", "Specify the file or URL of the DNS server list to be checked")
	output := flag.String("output", "", "Specify the output file")
	silent := flag.Bool("silent", false, "Whether to only output successfully resolved DNS servers")
	dns := flag.String("dns", "1.1.1.1,8.8.8.8", "Manually specify DNS servers, multiple servers are separated by commas")
	domain := flag.String("domain", "qq.com,tencent.com", "Manually specify domains, multiple domains are separated by commas")
	threads := flag.Int("threads", 20, "Specify the number of running threads")
	timeout := flag.Int("timeout", 10, "Specify the timeout for DNS requests")
	retry := flag.Int("retry", 2, "Specify the number of retries when matching fails")
	count := flag.Int("count", 65535, "Specify the program to stop running after writing how many qualified domain servers")

	flag.Parse()

	opts := Options{
		Target:  *target,
		Output:  *output,
		Silent:  *silent,
		Dns:     strings.Split(*dns, ","),
		Domain:  strings.Split(*domain, ","),
		Threads: *threads,
		Timeout: time.Duration(*timeout) * time.Second,
		Retry:   *retry,
		Count:   *count,
	}

	return opts
}

//func GetTargetServers(target string) ([]string, error) {
//	if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
//		resp, err := http.Get(target)
//		if err != nil {
//			return nil, err
//		}
//		defer resp.Body.Close()
//
//		body, err := ioutil.ReadAll(resp.Body)
//		if err != nil {
//			return nil, err
//		}
//
//		return strings.Split(string(body), "\n"), nil
//	} else {
//		data, err := ioutil.ReadFile(target)
//		if err != nil {
//			return nil, err
//		}
//
//		return strings.Split(string(data), "\n"), nil
//	}
//}

func GetTargetServers(target string) ([]string, error) {
	if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
		// Check if the targetServer is reachable
		_, err := net.DialTimeout("tcp", strings.TrimPrefix(strings.TrimPrefix(target, "http://"), "https://")+":80", 10*time.Second)
		if err != nil {
			return nil, except.ErrServerUnreachable
		}

		resp, err := http.Get(target)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return strings.Split(string(body), "\n"), nil
	} else {
		// Check if the target file is accessible
		_, err := os.Stat(target)
		if err != nil {
			return nil, except.ErrFileUnreachable
		}

		data, err := ioutil.ReadFile(target)
		if err != nil {
			return nil, err
		}

		return strings.Split(string(data), "\n"), nil
	}
}
