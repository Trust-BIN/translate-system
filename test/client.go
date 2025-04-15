package main

import (
	"fmt"
	"net"
	"translate-system/Routes"
)

func main() {
	conn, err := net.Dial("tcp", "8000")
	if err != nil {
		fmt.Println("net.Dial err:", err)
	}
	Routes.SettingRoutes()
	defer conn.Close()
}
