package node

import (
	"fmt"
	"log"
	"net"
)

type Node struct {
	Ip string
}

func FetchAllNodes() error {
	ip, mask, class, err := getLocalIP()
	if err != nil {
		return err
	}
	
	// Now, scan the subnet
	// Just search in the subnet by changing the last section to till 255, starting from 0
	for i := 1; i <= 255; i++ {
		ipAddress := fmt.Sprintf("%d,%d,%d,%d", ip[0], ip[1], ip[2], i)
		fmt.Println(ipAddress)
	}
	return nil
}

func findIpClass(ip net.IP) string {
	firstOctet := ip[0]

	switch {
	case firstOctet >= 1 && firstOctet <= 127:
		return "A"
	case firstOctet >= 128 && firstOctet <= 191:
		return "B"
	case firstOctet >= 192 && firstOctet <= 223:
		return "C"
	default:
		log.Fatal("Invalid IP address or not a classful IP")
	}
	return ""
}

func getLocalIP() (net.IP, string, string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	ip := localAddress.IP

	// Check if IP is IPv4
	if ip4 := ip.To4(); ip4 != nil {
		// find the subnet
		ipNet := ip4.DefaultMask()

		// Got the IP, identify the ip class, for easy subnet scan
		class := findIpClass(ip4)
		return ip4, ipNet.String(), class, nil
	}
	log.Printf("Can't obtain ipv4 address")
	return nil, "", "", fmt.Errorf("Error")
}

func generateIPAddresses(iprange net.IP) {

}
