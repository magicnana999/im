package broker

import (
	"encoding/binary"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/logger"
	"github.com/magicnana999/im/pb"
	"github.com/panjf2000/gnet/v2"
	bb "github.com/panjf2000/gnet/v2/pkg/pool/bytebuffer"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type codec interface {
	encode(c gnet.Conn, p *pb.Packet) (*bb.ByteBuffer, error)
	decode(c gnet.Conn) ([]*pb.Packet, error)
}

var defaultCodec = &lengthFieldBasedFrameCodec{}

type lengthFieldBasedFrameCodec struct {
}

func initCodec() codec {
	return defaultCodec
}

func (l *lengthFieldBasedFrameCodec) encode(c gnet.Conn, p *pb.Packet) (*bb.ByteBuffer, error) {

	buffer := bb.Get()

	user, err := currentUserFromConn(c)
	if err != nil {
		return nil, errors.EncodeError.Detail(err)
	}

	if p.IsHeartbeat() {

		hb := p.GetHeartbeatBody().Value

		binary.Write(buffer, binary.BigEndian, int32(4))
		binary.Write(buffer, binary.BigEndian, int32(hb))

		return buffer, nil
	} else {

		if logger.IsDebug() {
			bss, _ := protojson.Marshal(p)
			logger.DebugF("[%s#%s] encode %s", user.ClientAddr, user.Label(), string(bss))
		}

		bs, e := proto.Marshal(p)

		if e != nil {
			return nil, errors.EncodeError.Detail(e)
		}

		binary.Write(buffer, binary.BigEndian, int32(len(bs)))
		binary.Write(buffer, binary.BigEndian, bs)
		return buffer, nil
	}
}

func (l *lengthFieldBasedFrameCodec) decode(c gnet.Conn) ([]*pb.Packet, error) {

	result := make([]*pb.Packet, 0)

	user, err := currentUserFromConn(c)
	if err != nil {
		return nil, errors.DecodeError.Detail(err)

	}

	for c.InboundBuffered() >= 4 {

		var length int32
		binary.Read(c, binary.BigEndian, &length)

		if length == 4 && c.InboundBuffered() >= int(length) {

			var heartbeat int32
			binary.Read(c, binary.BigEndian, &heartbeat)

			packet := pb.NewHeartbeat(int32(heartbeat))
			result = append(result, packet)

		}

		if length > 4 && c.InboundBuffered() >= int(length) {

			bs := make([]byte, int(length))
			n, e := c.Read(bs)
			if e != nil {
				return nil, errors.DecodeError.Detail(e)
			}

			if n != int(length) {
				return nil, errors.DecodeError.DetailString("failed to read packet")
			}

			var p pb.Packet
			if e4 := proto.Unmarshal(bs, &p); e4 != nil {
				return nil, errors.DecodeError.Detail(e4)
			}

			if logger.IsDebug() {
				bss, _ := protojson.Marshal(&p)
				logger.DebugF("[%s#%s] decode %s", user.ClientAddr, user.Label(), string(bss))
			}

			result = append(result, &p)

		}

	}

	return result, nil
}
