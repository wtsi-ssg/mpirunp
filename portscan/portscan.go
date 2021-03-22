// Code largely taken from https://medium.com/@KentGruber/building-a-high-performance-port-scanner-with-golang-9976181ec39d

package portscan

import (
	"context"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

// PortScanner can be used to scan for available ports.
type PortScanner struct {
	ip      string
	lock    *semaphore.Weighted
	timeout time.Duration
}

// New makes a PortScanner that can scan ports on the given host, limited by
// the given semaphore and timeout.
func New(ip string, limiter *semaphore.Weighted, timeout time.Duration) *PortScanner {
	return &PortScanner{ip: ip, lock: limiter, timeout: timeout}
}

// PortAvailable checks to see if the given port on our host is available (that
// is to say, it would be possible to now open this port; returns true)
// or in use (returns false).
//
// NB: it doesn't know the difference between a port that is available and a
// port that is blocked by firewall rules.
func (ps *PortScanner) PortAvailable(port int) bool {
	target := net.JoinHostPort(ps.ip, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", target, ps.timeout)

	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "too many open files"):
			time.Sleep(ps.timeout)
			return ps.PortAvailable(port)
		case strings.Contains(msg, "refused"):
			return true
		default:
			return false
		}
	}

	conn.Close()
	return false
}

// AvailablePorts scans the given range of ports on our host, returning ports
// that are don't seem to be in use.
func (ps *PortScanner) AvailablePorts(min, max int) []int {
	ports := make(chan int, 20)

	wgRead := sync.WaitGroup{}
	wgRead.Add(1)
	var openPorts []int
	go func() {
		defer wgRead.Done()

		for port := range ports {
			openPorts = append(openPorts, port)
		}
	}()

	wg := sync.WaitGroup{}
	for port := min; port <= max; port++ {
		ps.lock.Acquire(context.TODO(), 1)
		wg.Add(1)

		go func(port int) {
			defer ps.lock.Release(1)
			defer wg.Done()

			if ps.PortAvailable(port) {
				ports <- port
			}
		}(port)
	}
	wg.Wait()
	close(ports)
	wgRead.Wait()

	return openPorts
}

// GetFreePort asks the kernel for free open ports that are ready to use.
func GetFreePorts(count int) ([]int, error) {
	var ports []int
	for i := 0; i < count; i++ {
		addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
		if err != nil {
			return nil, err
		}

		l, err := net.ListenTCP("tcp", addr)
		if err != nil {
			return nil, err
		}
		defer l.Close()
		ports = append(ports, l.Addr().(*net.TCPAddr).Port)
	}
	return ports, nil
}
