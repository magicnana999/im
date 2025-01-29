package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/magicnana999/im/broker/pb"
	"github.com/magicnana999/im/broker/protocol"
	"github.com/magicnana999/im/logger"
	"github.com/magicnana999/im/util"
	"google.golang.org/protobuf/proto"
	"net"
	"os"
	"sync"
	"time"
)

const (
	serverAddress     = "127.0.0.1:7539" // IM 服务器地址和端口
	heartbeatInterval = 1 * time.Second  // 心跳间隔
)

func main() {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		logger.FatalF("无法连接到服务器: %v", err)
	}
	defer conn.Close()

	var wg sync.WaitGroup
	wg.Add(3)

	ctx, cancel := context.WithCancel(context.Background())

	go sendHeartbeat(ctx, conn, &wg)

	//go readMessages(ctx, conn, &wg)

	go sendMessage(ctx, cancel, conn, &wg)

	wg.Wait()
	fmt.Println("OK")
}

func sendMessage(ctx context.Context, cancel context.CancelFunc, conn net.Conn, wg *sync.WaitGroup) {

	defer wg.Done()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Print("请输入消息: ")
			if !scanner.Scan() {
				continue
			}
			message := scanner.Text()
			if message != "exit" {
				cancel()
			} else {
				writeMessage(conn, message)
			}
		}
	}
}

func writeMessage(conn net.Conn, message string) {
	text := protocol.TextContent{
		Text: "message",
	}

	body := protocol.MessageBody{
		MType:    protocol.MText,
		CId:      "sdfsdf",
		To:       "sdfsdf",
		GroupId:  "",
		TType:    protocol.TSingle,
		Sequence: 100,
		Content:  text,
		At:       nil,
		Refer:    nil,
	}

	packet := &protocol.Packet{
		Id:      util.GenerateXId(),
		AppId:   "STARTSPACE",
		UserId:  "sdifejrjersdf",
		Flow:    protocol.FlowRequest,
		NeedAck: protocol.YES,
		Type:    protocol.TypeMessage,
		CTime:   12123123,
		STime:   12123123,
		Body:    body,
	}

	p, e := pb.ConvertPacket(packet)
	if e != nil {
		panic(e)
	}

	b, ee := proto.Marshal(p)
	if ee != nil {
		panic(e)
	}

	conn.Write(b)
}

func sendHeartbeat(ctx context.Context, conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:

			buffer := new(bytes.Buffer)

			binary.Write(buffer, binary.BigEndian, uint32(4))

			binary.Write(buffer, binary.BigEndian, uint32(12))

			_, err := conn.Write(buffer.Bytes())
			if err != nil {
				fmt.Println(err)
			}
			//fmt.Println("发送心跳 4 12")
		case <-ctx.Done():
			return
		}
	}
}

func readMessages(ctx context.Context, conn net.Conn, wg *sync.WaitGroup) {

	defer wg.Done()

	reader := bufio.NewReader(conn)
	for {

		select {
		case <-ctx.Done():
			return
		default:

			bs, err := reader.Peek(4)
			reader.Discard(4)
			if err != nil {
				panic(err)
			}

			length := int(binary.BigEndian.Uint32(bs))

			if length == 4 && reader.Buffered() >= length {
				b, e := reader.Peek(length)
				reader.Discard(length)
				if e != nil {
					panic(e)
				}

				heartbeat := binary.BigEndian.Uint32(b)
				fmt.Println("Read heartbeat:", heartbeat)
			}

			if length > 4 && reader.Buffered() >= length {
				bb, ee := reader.Peek(length)
				reader.Discard(length)
				if ee != nil {
					panic(ee)
				}

				var p pb.Packet
				if e4 := proto.Unmarshal(bb, &p); err != nil {
					panic(e4)
				}

				packet, eee := pb.RevertPacket(&p)
				if eee != nil {
					panic(eee)
				}

				js, eeee := json.Marshal(packet)
				if eeee != nil {
					panic(eeee)
				}

				fmt.Println(string(js))
			}
		}
	}
}
