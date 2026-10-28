package ports

import (
	"net"
	"strconv"

	gnet "github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

func CheckFreePorts(port int) bool {
	addr := net.JoinHostPort("", strconv.Itoa(port))
	l,err := net.Listen("tcp",addr)
	
	if err != nil {
		return false
	}

	_ = l.Close()

	return true
}

func FirstFree(start,end int) {
	for port := start; port < end; port++ {
		if CheckFreePorts(port) {
			return port, nil
		}
	}

	return 0,ErrNoFreePort
}