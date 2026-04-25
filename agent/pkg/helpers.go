package pkg

import "net"

func GetLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        return "unknown"
    }
    defer conn.Close()
    return conn.LocalAddr().(*net.UDPAddr).IP.String()
}
