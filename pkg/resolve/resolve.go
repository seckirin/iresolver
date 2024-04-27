package resolve

import (
	"context"
	"fmt"
	"github.com/yuukisec/iresolver/pkg/exporting"
	"github.com/yuukisec/iresolver/pkg/options"
	"github.com/yuukisec/iresolver/pkg/output"
	"net"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func ResolveAndCompare(opts options.Options) {
	baselineResults := getBaselineResults(opts)
	var wg sync.WaitGroup
	sem := make(chan struct{}, opts.Threads)

	targetServers, _ := options.GetTargetServers(opts.Target)
	successCount := int32(0)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, targetServer := range targetServers {
		wg.Add(1)
		sem <- struct{}{}
		go func(targetServer string) {
			defer wg.Done()
			success := compareResults(ctx, &successCount, targetServer, baselineResults, opts)
			if success {
				if atomic.AddInt32(&successCount, 1) >= int32(opts.Count) {
					cancel()
				}
			}
			<-sem
		}(targetServer)
	}
	wg.Wait()
}

func getBaselineResults(opts options.Options) map[string][]string {
	baselineResults := make(map[string][]string)
	var wg sync.WaitGroup
	sem := make(chan struct{}, opts.Threads)

	for _, domain := range opts.Domain {
		for _, dnsServer := range opts.Dns {
			wg.Add(1)
			sem <- struct{}{}
			go func(domain, dnsServer string) {
				defer wg.Done()
				var dnsResult []string
				var err error
				for i := 0; i < opts.Retry; i++ {
					dnsResult, err = ResolveDNS(domain, dnsServer, opts.Timeout, opts.Retry)
					if err == nil {
						break
					}
					if !opts.Silent {
						output.PrintError(err)
					}
				}
				if err != nil {
					return
				}
				sort.Strings(dnsResult)
				baselineResults[domain] = append(baselineResults[domain], strings.Join(dnsResult, ","))
				if !opts.Silent {
					output.PrintInfo(fmt.Sprintf("The DNS results of domain %s from baseline DNS server %s are %s", domain, dnsServer, strings.Join(dnsResult, ",")))
				}
				<-sem
			}(domain, dnsServer)
		}
	}
	wg.Wait()

	return baselineResults
}

func compareResults(ctx context.Context, successCount *int32, targetServer string, baselineResults map[string][]string, opts options.Options) bool {
	var successDomains []string
	var wg sync.WaitGroup
	sem := make(chan struct{}, opts.Threads)

	for _, domain := range opts.Domain {
		wg.Add(1)
		sem <- struct{}{}
		go func(domain string) {
			defer wg.Done()
			if atomic.LoadInt32(successCount) >= int32(opts.Count) {
				return
			}
			var targetResult []string
			var err error
			for i := 0; i < opts.Retry; i++ {
				targetResult, err = ResolveDNS(domain, targetServer, opts.Timeout, opts.Retry)
				if err == nil {
					break
				}
				if !opts.Silent {
					output.PrintError(fmt.Errorf("Error resolving domain %s on server %s: %w", domain, targetServer, err))
				}
			}
			if err != nil {
				return
			}
			sort.Strings(targetResult)
			target := strings.Join(targetResult, ",")
			if !opts.Silent {
				output.PrintInfo(fmt.Sprintf("The DNS results of domain %s from target DNS server %s are %s", domain, targetServer, target))
			}
			match := false
			for _, baseline := range baselineResults[domain] {
				if baseline == target {
					match = true
					break
				}
			}
			if match {
				successDomains = append(successDomains, domain)
			} else {
				output.PrintInfo(fmt.Sprintf("The DNS results of domain %s from target DNS server %s do not match the baseline results. This may indicate a DNS hijacking.", domain, targetServer))
			}
			<-sem
			if ctx.Err() != nil {
				return
			}
		}(domain)
	}
	wg.Wait()

	if len(successDomains) == len(opts.Domain) {
		if opts.Silent {
			output.PrintSilent(targetServer)
		} else {
			output.PrintSuccess(fmt.Sprintf("Server %s successfully resolved domains %s", targetServer, strings.Join(successDomains, ", ")))
			err := exporting.ExportToFile(opts.Output, []string{targetServer})
			if err != nil {
				output.PrintError(fmt.Errorf("Error exporting results to file: %w", err))
			}
		}
		return true
	}
	return false
}

func ResolveDNS(domain string, server string, timeout time.Duration, retry int) ([]string, error) {
	var err error
	var ips []net.IP
	for i := 0; i < retry; i++ {
		r := &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: timeout,
				}
				return d.DialContext(ctx, network, server+":53")
			},
		}

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		ips, err = r.LookupIP(ctx, "ip", domain)
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, err
	}

	var results []string
	for _, ip := range ips {
		results = append(results, ip.String())
	}

	return results, nil
}
