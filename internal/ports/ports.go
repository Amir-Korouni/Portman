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
