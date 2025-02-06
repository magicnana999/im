package core

import (
	"encoding/binary"
	"github.com/magicnana999/im/broker/state"
	"github.com/magicnana999/im/common/pb"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/logger"
	"github.com/panjf2000/gnet/v2"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"math"
)

type Codec interface {
	Encode(c gnet.Conn, p *pb.Packet) ([][]byte, error)
	Decode(c gnet.Conn) ([]*pb.Packet, error)
}

var DefaultCodec = &LengthFieldBasedFrameCodec{}

type LengthFieldBasedFrameCodec struct {
}

func (l LengthFieldBasedFrameCodec) Encode(c gnet.Conn, p *pb.Packet) ([][]byte, error) {

	user, err := state.CurrentUserFromConn(c)
	if err != nil {
		return nil, errors.ConnectionEncodeError.Fill("failed to get current user," + err.Error())
	}

	if pb.IsHeartbeat(p) {

		var hb wrapperspb.UInt32Value
		if err := p.Body.UnmarshalTo(&hb); err != nil {
			return nil, errors.ConnectionEncodeError.Fill("failed to unmarshal length field," + err.Error())
		}

		if hb.Value < 0 || hb.Value >= math.MaxInt32 {
			return nil, errors.ConnectionEncodeError.Fill("invalid heartbeat value")
		}

		bs := make([][]byte, 2)
		bs[0] = make([]byte, 4)
		bs[1] = make([]byte, 4)

		binary.BigEndian.PutUint32(bs[0], uint32(4))
		binary.BigEndian.PutUint32(bs[1], hb.Value)

		logger.DebugF("[%s#%s] Encode heartbeat,buffer:%d,length:%d,value:%s",
			c.RemoteAddr().String(),
			user.Label(),
			len(bs[0])+len(bs[1]),
			len(bs[1]),
			hb.Value)

		return bs, nil
	} else {

		var err error

		bs := make([][]byte, 2)
		bs[0] = make([]byte, 4)
		bs[1], err = proto.Marshal(p)

		if err != nil {
			return nil, errors.ConnectionDecodeError.Fill("failed to marshal packet," + err.Error())
		}

		if len(bs[1]) <= 0 || len(bs[1]) >= math.MaxInt32 {
			return nil, errors.ConnectionEncodeError.Fill("invalid packet length")
		}

		binary.BigEndian.PutUint32(bs[0], uint32(len(bs[1])))

		logger.DebugF("[%s#%s] Encode packet,buffer:%d,length:%d,id:%s",
			c.RemoteAddr().String(),
			user.Label(),
			len(bs[0])+len(bs[1]),
			len(bs[1]),
			p.Id)

		return bs, nil
	}
}

func (l LengthFieldBasedFrameCodec) Decode(c gnet.Conn) ([]*pb.Packet, error) {

	result := make([]*pb.Packet, 0)

	user, err := state.CurrentUserFromConn(c)
	if err != nil {
		return nil, errors.ConnectionDecodeError.Fill("failed to get current user," + err.Error())

	}

	for c.InboundBuffered() >= 4 {

		bs, e1 := c.Next(4)
		if e1 != nil {

			c.Discard(c.InboundBuffered())
			return nil, errors.ConnectionDecodeError.Fill("failed to read length field," + e1.Error())
		}

		l := binary.BigEndian.Uint32(bs)

		if l <= 0 || l >= math.MaxInt32 {

			c.Discard(c.InboundBuffered())
			ee := errors.ConnectionDecodeError.Fill("invalid packet length")
			return nil, ee

		}

		length := int(l)

		if length == 4 && c.InboundBuffered() >= length {

			b, e2 := c.Next(length)

			if e2 != nil {

				c.Discard(c.InboundBuffered())
				return nil, errors.ConnectionDecodeError.Fill("failed to read heartbeat," + e2.Error())
			}

			heartbeat := binary.BigEndian.Uint32(b)

			logger.InfoF("[%s#%s] Decode heartbeat,buffer:%d,length:%d,value:%d",
				c.RemoteAddr().String(),
				user.Label(),
				c.InboundBuffered(),
				length,
				heartbeat)

			if heartbeat > 0 && heartbeat < math.MaxInt32 {

				val := wrapperspb.UInt32(heartbeat)
				body, err2 := anypb.New(val)
				if err2 != nil {
					return nil, err2
				}

				packet := &pb.Packet{
					BType: pb.BTypeHeartbeat,
					Body:  body,
				}

				result = append(result, packet)
			}

		}

		if length > 4 && c.InboundBuffered() >= length {

			b, e3 := c.Next(length)

			if e3 != nil {

				c.Discard(c.InboundBuffered())
				return nil, errors.ConnectionDecodeError.Fill("failed to read packet," + e3.Error())
			}

			var p pb.Packet
			if e4 := proto.Unmarshal(b, &p); e4 != nil {

				c.Discard(c.InboundBuffered())
				return nil, errors.ConnectionDecodeError.Fill("failed to unmarshal packet," + e4.Error())
			}

			logger.InfoF("[%s#%s] Decode packet,buffer:%d,length:%d,id:%s",
				c.RemoteAddr().String(),
				user.Label(),
				c.InboundBuffered(),
				length,
				p.Id)

			result = append(result, &p)

		}

	}

	return result, nil
}
