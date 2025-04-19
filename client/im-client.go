package main

import "fmt"

func main() {
	tcp := NewTcpClientGroup()
	go tcp.Start()
	defer tcp.Stop()
	fmt.Println(tcp.NewClient())
}
