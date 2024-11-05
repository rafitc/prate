package node

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type Listener struct {
	Peer Node
}

type Node struct {
	Ip string
}

func FetchAllNodes() error {
	ip, mask, class, err := getLocalIP()
	if err != nil {
		return err
	}
	fmt.Printf("%s - %s ", mask, class)
	generateIps(ip, mask)
	// Now, scan the subnet
	// Just search in the subnet by changing the last section to till 255, starting from 0

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

func generateIps(ip net.IP, mask string) {
	// First identify the start and end ip range
	bits := hexToBinaryOfMask(mask)
	fmt.Print(bits)
	startIp, endIp, err := calculateIPRange(ip.String(), bits)
	if err != nil {
		log.Fatal("Something went wrong while calculating ip and netmask ")
	}

	// build all possible IP address
	var allIps []Node
	for s1 := int(startIp[0]); s1 <= int(endIp[0]); s1++ {
		for s2 := int(startIp[1]); s2 <= int(endIp[1]); s2++ {
			for s3 := int(startIp[2]); s3 <= int(endIp[2]); s3++ {
				for s4 := int(startIp[3]); s4 <= int(endIp[3]); s4++ {
					allIps = append(allIps, Node{fmt.Sprintf("%v.%v.%v.%v", s1, s2, s3, s4)})
				}
			}
		}
	}

	// spin a go routine for each ip range of 10, 192.168.1.0 - 192.168.1.10
	// this goroutine pass into nmap to check is there any listeners in this ip range, if then update the list using go channels
	// listeners := make(chan Listener)

}

func hexToBinaryOfMask(hex string) int {
	binaryString := ""
	bits := 0

	for _, hexDigit := range hex {
		decimalDigit, _ := strconv.ParseInt(string(hexDigit), 16, 64)
		binaryDigit := fmt.Sprintf("%04b", decimalDigit)
		binaryString += binaryDigit
		bits += strings.Count(binaryDigit, "1")
	}
	return bits
}

func calculateIPRange(ipStr string, bits int) (net.IP, net.IP, error) {
	// Parse IP address
	ip := net.ParseIP(ipStr).To4()
	if ip == nil {
		return nil, nil, fmt.Errorf("invalid IPv4 address")
	}

	// Calculate subnet mask
	var mask uint32 = ^(uint32(0xFFFFFFFF) >> bits)

	// Convert IP to uint32 for bitwise operations
	ipInt := binary.BigEndian.Uint32(ip)

	// Calculate start and end IPs in the range
	startIPInt := ipInt & mask
	endIPInt := ipInt | ^mask

	// Convert start and end IPs back to byte slices
	startIP := make(net.IP, 4)
	endIP := make(net.IP, 4)
	binary.BigEndian.PutUint32(startIP, startIPInt)
	binary.BigEndian.PutUint32(endIP, endIPInt)

	return startIP, endIP, nil
}
