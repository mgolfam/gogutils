package utils

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/mgolfam/gogutils/glog"
)

func GetLocalPAddressByInterfaceName(interfaceName string) ([]string, error) {
	// Retrieve the network interface by name
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return nil, err
	}

	// Retrieve the addresses associated with the interface
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}

	ips := make([]string, len(addrs))

	// Display the IP addresses associated with the interface
	for index, addr := range addrs {
		ips[index] = addr.String()

		if strings.Contains(ips[index], "/") {
			ips[index] = strings.Split(ips[index], "/")[0]
		}
	}

	return ips, nil
}

func GetLocalIPAddressByInterfaceNameContains(interfaceName string) ([]string, error) {
	// Retrieve all network interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		if strings.Contains(iface.Name, interfaceName) {
			return GetLocalPAddressByInterfaceName(iface.Name)
		}
	}

	return nil, errors.New("no interface contains " + interfaceName)
}

// Function 1: Get the local machine's IP addresses
func GetLocalIPAddresses() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	var ips []string
	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err == nil && ip.IsGlobalUnicast() {
			ips = append(ips, ip.String())
		}
	}

	return ips, nil
}

// Function 2: Resolve a domain to its IP addresses
func ResolveDNS(domain string) ([]string, error) {
	ips, err := net.LookupHost(domain)
	if err != nil {
		return nil, err
	}

	return ips, nil
}

// Function 3: Get the IP address of a domain's mail server
func ResolveMX(domain string) ([]*net.MX, error) {
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return nil, err
	}

	return mxRecords, nil
}

// Function 4: Check if a specific port is open on a remote host
func IsPortOpen(host string, port int) bool {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// Function 5: Get the MAC address of a network interface
func GetMACAddress(interfaceName string) (string, error) {
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return "", err
	}

	return iface.HardwareAddr.String(), nil
}

func TCPServer(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		glog.LogL(glog.DEBUG, "Error listening:", err)
		return
	}
	defer listener.Close()

	glog.LogL(glog.DEBUG, "Listening on port %d\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			glog.LogL(glog.DEBUG, "Error accepting connection:", err)
			continue
		}

		go handleTCPConnection(conn)
	}
}

func handleTCPConnection(conn net.Conn) {
	// Handle TCP connection logic here
	conn.Close()
}

func UDPSend(destination string, port int, data []byte) error {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", destination, port))
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func UDPReceive(port int) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		glog.LogL(glog.DEBUG, "Error resolving UDP address:", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		glog.LogL(glog.DEBUG, "Error listening for UDP packets:", err)
		return
	}
	defer conn.Close()

	glog.LogL(glog.DEBUG, "Listening for UDP packets on port %d\n", port)

	buffer := make([]byte, 1024)

	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			glog.LogL(glog.DEBUG, "Error reading UDP packet:", err)
			continue
		}

		glog.LogL(glog.DEBUG, "Received UDP packet: %s\n", buffer[:n])
	}
}

func GetPublicIPAddress() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}
