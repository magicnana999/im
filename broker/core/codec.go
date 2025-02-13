package core

import (
	"encoding/binary"
	"github.com/magicnana999/im/broker/state"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/logger"
	"github.com/magicnana999/im/pb"
	"github.com/panjf2000/gnet/v2"
	bb "github.com/panjf2000/gnet/v2/pkg/pool/bytebuffer"
	"google.golang.org/protobuf/proto"
)

type Codec interface {
	Encode(c gnet.Conn, p *pb.Packet) (*bb.ByteBuffer, error)
	Decode(c gnet.Conn) ([]*pb.Packet, error)
}

var defaultCodec = &LengthFieldBasedFrameCodec{}

type LengthFieldBasedFrameCodec struct {
}

func InitCodec() Codec {
	return defaultCodec
}

func (l LengthFieldBasedFrameCodec) Encode(c gnet.Conn, p *pb.Packet) (*bb.ByteBuffer, error) {

	buffer := bb.Get()

	user, err := state.CurrentUserFromConn(c)
	if err != nil {
		return nil, errors.EncodeError.Detail(err)
	}

	if p.IsHeartbeat() {

		hb := p.GetHeartbeatBody().Value

		binary.Write(buffer, binary.BigEndian, int32(4))
		binary.Write(buffer, binary.BigEndian, int32(hb))

		logger.DebugF("[%s#%s] Encode heartbeat,buffer:%d,length:%d,value:%d",
			c.RemoteAddr().String(),
			user.Label(),
			buffer.Len(),
			buffer.Len()-4,
			hb)

		return buffer, nil
	} else {

		bs, e := proto.Marshal(p)

		if e != nil {
			return nil, errors.EncodeError.Detail(e)
		}

		binary.Write(buffer, binary.BigEndian, int32(len(bs)))
		binary.Write(buffer, binary.BigEndian, bs)

		logger.DebugF("[%s#%s] Encode packet,buffer:%d,length:%d,id:%v",
			c.RemoteAddr().String(),
			user.Label(),
			buffer.Len(),
			len(bs),
			p)

		return buffer, nil
	}
}

func (l LengthFieldBasedFrameCodec) Decode(c gnet.Conn) ([]*pb.Packet, error) {

	result := make([]*pb.Packet, 0)

	user, err := state.CurrentUserFromConn(c)
	if err != nil {
		return nil, errors.DecodeError.Detail(err)

	}

	for c.InboundBuffered() >= 4 {

		var length int32
		binary.Read(c, binary.BigEndian, &length)

		if length == 4 && c.InboundBuffered() >= int(length) {

			var heartbeat int32
			binary.Read(c, binary.BigEndian, &heartbeat)

			logger.DebugF("[%s#%s] Decode heartbeat,buffer:%d,length:%d,value:%d",
				c.RemoteAddr().String(),
				user.Label(),
				c.InboundBuffered(),
				length,
				heartbeat)

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

			logger.DebugF("[%s#%s] Decode packet,buffer:%d,length:%d,packet:%v",
				c.RemoteAddr().String(),
				user.Label(),
				c.InboundBuffered(),
				length,
				p)

			result = append(result, &p)

		}

	}

	return result, nil
}
