package testhelper

import (
	"fmt"
	"net"
)

// GetFreePort asks the kernel for a free open port that is ready to use.
func GetFreePort(host string, preferredPort uint32) (int, error) {
	address := host + ":" + fmt.Sprint(preferredPort)
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// AllocatePort returns a port that is available, given host and a preferred port
// if none of the preferred ports are available, it will keep searching by adding 1 to the port number
func AllocatePort(host string, preferredPort uint32) uint32 {
	preferredPortStr := fmt.Sprint(preferredPort)
	allocatedPort, err := GetFreePort(host, preferredPort)
	for err != nil {
		preferredPort = preferredPort + 1
		allocatedPort, err = GetFreePort(host, preferredPort)
		if err != nil {
			fmt.Println("Failed to connect to", preferredPortStr)
		}
	}
	fmt.Println("Allocated port", allocatedPort)
	return uint32(allocatedPort)
}
