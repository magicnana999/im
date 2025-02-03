package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/magicnana999/im/common/pb"
	protocol2 "github.com/magicnana999/im/common/protocol"
	"github.com/magicnana999/im/logger"
	"github.com/magicnana999/im/util/id"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
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
			if message == "exit" {
				cancel()
			} else if message == "login" {
				fmt.Print("请输入UserSig: ")
				if !scanner.Scan() {
					continue
				}
				userSig := scanner.Text()

				writeLogin(conn, userSig)

			} else {
				writeMessage(conn, message)
			}
		}
	}
}

func writeLogin(conn net.Conn, userSig string) {
	loginContent := pb.LoginContent{
		AppId:        "19860220",
		UserSig:      userSig,
		Version:      "1.0.0",
		Os:           pb.OSType_OSIos,
		PushDeviceId: id.GenerateXId(),
	}

	c, _ := anypb.New(&loginContent)

	commandBody := pb.CommandBody{
		MType:   protocol2.MUserLogin,
		Content: c,
	}

	b, _ := anypb.New(&commandBody)
	packet := pb.Packet{
		Id:    id.GenerateXId(),
		AppId: loginContent.AppId,
		Type:  protocol2.TypeCommand,
		CTime: time.Now().UnixMilli(),
		Body:  b,
	}

	body, _ := proto.Marshal(&packet)

	buffer1 := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer1, uint32(len(body)))

	buffer := new(bytes.Buffer)
	buffer.Write(buffer1)
	buffer.Write(body)

	conn.Write(buffer.Bytes())
}

func writeMessage(conn net.Conn, message string) {
	text := protocol2.TextContent{
		Text: message,
	}

	body := protocol2.MessageBody{
		MType:    protocol2.MText,
		CId:      "sdfsdf",
		To:       "sdfsdf",
		GroupId:  "",
		TType:    protocol2.TSingle,
		Sequence: 100,
		Content:  text,
		At:       nil,
		Refer:    nil,
	}

	packet := &protocol2.Packet{
		Id:      id.GenerateXId(),
		AppId:   "STARTSPACE",
		UserId:  10012,
		Flow:    protocol2.FlowRequest,
		NeedAck: protocol2.YES,
		Type:    protocol2.TypeMessage,
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

	buffer1 := make([]byte, 4)
	buffer2 := bytes.NewBuffer(b)

	binary.BigEndian.PutUint32(buffer1, uint32(len(b)))

	buffer := new(bytes.Buffer)
	buffer.Write(buffer1)
	buffer.Write(buffer2.Bytes())

	conn.Write(buffer.Bytes())

	var pbp pb.Packet
	if e4 := proto.Unmarshal(b, &pbp); e4 != nil {
		panic(e4)
	}

	ret, e5 := pb.RevertPacket(&pbp)
	if e5 != nil {
		panic(e5)
	}

	js, e7 := json.Marshal(ret)
	if e7 != nil {
		panic(e7)
	}

	fmt.Println(string(js))
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
