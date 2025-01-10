package net

import (
	"fmt"
	"net"
)

func ListenUDP(addr string) {
	ad, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	conn, err := net.ListenUDP("udp", ad)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	handleUDP(conn)
	//defer conn.Close()

	// 死循环，每当遇到连接时，调用 handle
	//for {
	//
	//	if err != nil {
	//		fmt.Println("Error accepting:", err)
	//		return
	//	}
	//	go handleUDP(conn)
	//}
}

func handleUDP(conn net.PacketConn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			//return
		}
		_, err = conn.WriteTo(buf[:n], addr)
		if err != nil {
			fmt.Println("Error writing to UDP:", err)
			//return
		}
	}
}

func forwardUDP(localConn, remoteConn *net.UDPConn) {
	defer localConn.Close()
	defer remoteConn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := localConn.Read(buf)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			return
		}
		_, err = remoteConn.Write(buf[:n])
		if err != nil {
			fmt.Println("Error writing to UDP:", err)
			return
		}
	}
}
