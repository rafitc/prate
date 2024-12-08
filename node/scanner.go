package node

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Listener struct {
	Peer Node
}

type Node struct {
	Ip string
}

var port uint16 = 88
var ConnWatcher = make(chan Listener)

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

	// build all possible IP address TODO only for first time, no need to redo till connection reset
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
	ConnWatcher = make(chan Listener, len(allIps))
	var blockOfNodes [][]string
	var block []string
	wg := &sync.WaitGroup{}

	for index, eachIp := range allIps {
		block = append(block, eachIp.Ip)

		// Append block when it reaches 10 items
		if (index+1)%10 == 0 {
			blockOfNodes = append(blockOfNodes, append([]string(nil), block...))
			block = block[:0] // Clear the block for the next set
		}
	}

	// Append any remaining IPs in the final block if not empty
	if len(block) > 0 {
		blockOfNodes = append(blockOfNodes, append([]string(nil), block...))
	}

	// Create a ticker to run the scan every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			// Scan IP blocks every 5 seconds
			for _, eachBlock := range blockOfNodes {
				wg.Add(1)
				go isNodeListening(eachBlock, port, ConnWatcher, wg)
			}
		}
	}()

	// Monitor worker to wait for goroutines to finish
	go monitorWorker(wg, ConnWatcher)

	// Receive values from the listeners channel
	for value := range ConnWatcher {
		fmt.Println(value)
	}
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

func isNodeListening(ips []string, port uint16, ch chan Listener, wg *sync.WaitGroup) {
	// ips, list of ip address so, one gocoroutine can check consecutive 10 ips together
	defer wg.Done()
	// Construct command arguments
	args := append([]string{"--open", "-Pn", fmt.Sprintf("-p%d", port), "-sS"}, ips...)
	cmd := exec.Command("nmap", args...)

	// Capture the command's output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	// execute
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("stderr:", stderr.String())
		return
	}
	ipAddresses := extractIPAddresses(stdout.String())
	for _, ip := range ipAddresses {
		ch <- Listener{Node{Ip: ip}}
	}
}

func extractIPAddresses(scanOutput string) []string {
	ipPattern := `\b(?:\d{1,3}\.){3}\d{1,3}\b`
	re := regexp.MustCompile(ipPattern)
	matches := re.FindAllString(scanOutput, -1)

	return matches
}

func monitorWorker(wg *sync.WaitGroup, cs chan Listener) {
	wg.Wait()
	close(cs)
}
