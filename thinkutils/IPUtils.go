package thinkutils

import "net"

type iputils struct {
}

func (this iputils) LocalIPv4s() ([]string, error) {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP.String())
		}
	}

	return ips, nil
}

func (this iputils) LocalIP() string {
	ips, _ := this.LocalIPv4s()
	return ips[0]
}
