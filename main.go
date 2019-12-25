package main
import (
	"fmt"
	"net"
	"time"
)

func initUdpServer() (*UDPConn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8800")
	if err != nil {
		fmt.Println("ResolveUDPAddr err:", err)
		return nil, err
	}
	conn, err := net.ListenUDP("udp", udpAddr) //创建数据通信socket
	if err != nil {
		fmt.Println("ListenUDP err:", err)
		return nil, err
	}
	defer conn.Close()
	return conn, err
}

func main() {
	fmt.Println("Server Start...")
	conn, _ := initUdpServer()

	for i := 0; i < 100000000; i++ {
		time.Sleep(1)
		buf := make([]byte, 1024)
		n, raddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			return
		}
		fmt.Println("[Client]", string(buf[:n]))

		_, err = conn.WriteToUDP([]byte("OK"), raddr)
		if err != nil {
			fmt.Println("WriteToUDP err:", err)
			return
		}
	}
}
