package lsf

import (
	"os"
	"strings"
)

const hostsEnv = "LSB_HOSTS"

// Hosts returns LSB_HOSTS environment variable as a slice of the hosts it
// describes. If LSB_HOSTS is not set, you will get a nil slice.
func Hosts() []string {
	if hosts := os.Getenv(hostsEnv); hosts != "" {
		return strings.Fields(hosts)
	}
	return nil
}
