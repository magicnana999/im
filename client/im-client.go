package main

import (
	"fmt"
	console "github.com/asynkron/goconsole"
	"sync"
	"time"
)

type IMClients struct {
	tcp *TcpServer
	hts *HeartbeatServer
	m   sync.Map
	mt  sync.Map
}

func (s *IMClients) Connect(user *User) {
	s.mt.Store(user.FD, user)
}

func (s *IMClients) DisConnect(user *User) {
	s.m.Delete(user.UserID)
}

func NewIMClients() *IMClients {

	im := &IMClients{}

	hts := NewHeartbeatServer(time.Second * 30)
	handler := NewPacketHandler()
	tcp := NewTcpServer(handler, hts, im)

	im.tcp = tcp
	im.hts = hts

	return im
}

func (s *IMClients) Start() {
	go s.tcp.Start()
	go s.hts.Start()
}

func (s *IMClients) Stop() {
	s.tcp.Stop()
	s.hts.Stop()
}

func (s *IMClients) Login() {

}
func main() {
	s := NewIMClients()
	defer s.Stop()

	c := console.NewConsole(func(s string) {
		fmt.Println(s)
	})

	c.Command("conn", func(text string) {
		s.tcp.Connect()
	})
	c.Run()
}
