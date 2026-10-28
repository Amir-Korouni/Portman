package ports

import (
	"cmp"
	"context"
	"fmt"
	"slices"

	gnet "github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

type Entry struct {
	Port    int
	PID     int32  
	Process string 
	Address string
}

func Occupied(ctx context.Context, from, to int) ([]Entry, error) {
	conns, err := gnet.ConnectionsWithContext(ctx, "tcp")
	if err != nil {
		return nil, fmt.Errorf("list connections: %w", err)
	}

	seen := map[string]bool{} 
	var out []Entry

	for _, c := range conns {
		if c.Status != "LISTEN" {
			continue
		}
		port := int(c.Laddr.Port)
		if port < from || port > to {
			continue
		}
		key := fmt.Sprintf("%d:%d", port, c.Pid)
		if seen[key] {
			continue
		}
		seen[key] = true

		e := Entry{Port: port, PID: c.Pid, Address: c.Laddr.IP}
		if c.Pid > 0 {
			if p, err := process.NewProcess(c.Pid); err == nil {
				e.Process, _ = p.Name()
			}
		}
		out = append(out, e)
	}

	slices.SortFunc(out, func(a, b Entry) int { return cmp.Compare(a.Port, b.Port) })
	return out, nil
}