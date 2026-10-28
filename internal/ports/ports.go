package ports

import (
	"errors"
	"fmt"
	"net"
	"strconv"
)

func CheckFreePorts(port int) bool {
	addr := net.JoinHostPort("", strconv.Itoa(port))
	l, err := net.Listen("tcp", addr)

	if err != nil {
		return false
	}

	_ = l.Close()

	return true
}

func FirstFree(start, end int) (int, error) {
	if start < 1 || end > 65535 || start > end {
		return 0, fmt.Errorf("invalid port range %d-%d", start, end)
	}

	for port := start; port < end; port++ {
		if CheckFreePorts(port) {
			return port, nil
		}
	}

	return 0, ErrNoFreePort
}

var ErrNoFreePort = errors.New("no free port in range")
