package lanmsg

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	broadcastIP   = "192.168.1.255"
	broadcastPort = 8989
)

func StartMessaging(name string) {
	go receiveMessages()

	sendMessages(name)
}

func receiveMessages() {
	addr := &net.UDPAddr{
		Port: broadcastPort,
		IP:   net.IPv4zero, // 监听所有地址
	}

	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		fmt.Println("UDP监听失败:", err)
		return
	}
	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("接收失败:", err)
			continue
		}
		if n == 0 || remoteAddr == nil {
			continue
		}
		msg := strings.TrimSpace(string(buf[:n]))
		fmt.Printf("\n[来自 %s] %s\n你: ", remoteAddr.IP, msg)
	}
}

func sendMessages(name string) {
	broadcastAddr := &net.UDPAddr{
		IP:   net.ParseIP(broadcastIP),
		Port: broadcastPort,
	}

	conn, err := net.DialUDP("udp4", nil, broadcastAddr)
	if err != nil {
		fmt.Println("发送失败:", err)
		return
	}
	defer conn.Close()

	// ⚠️ 设置为允许广播
	conn.SetWriteBuffer(1024)
	conn.Write([]byte{}) // 强制触发连接？

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("你:")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)
		if message != "exit" {
			finalMessage := fmt.Sprintf("%s:%s", name, message)
			conn.Write([]byte(finalMessage))
		}
	}

}
