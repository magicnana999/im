package core

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/magicnana999/im/broker/state"
	"github.com/magicnana999/im/common/pb"
	"github.com/magicnana999/im/logger"
	"github.com/panjf2000/gnet/v2"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"math"
)

type Codec interface {
	Encode(p *pb.Packet) ([][]byte, error)
	Decode(s *BrokerServer, c gnet.Conn) ([]*pb.Packet, error)
}

var defaultCodec = &LengthFieldBasedFrameCodec{}

type LengthFieldBasedFrameCodec struct {
}

func (l LengthFieldBasedFrameCodec) Encode(p *pb.Packet) ([][]byte, error) {
	if pb.IsHeartbeat(p) {

		var hb wrapperspb.UInt32Value
		if err := p.Body.UnmarshalTo(&hb); err != nil {
			return nil, err
		}

		if hb.Value <= 0 || hb.Value >= math.MaxInt32 {
			return nil, errors.New("invalid heartbeat body")
		}

		bs := make([][]byte, 2)
		bs[0] = make([]byte, 4)
		bs[1] = make([]byte, 4)

		binary.BigEndian.PutUint32(bs[0], uint32(4))
		binary.BigEndian.PutUint32(bs[1], hb.Value)
		return bs, nil
	} else {
		var err error

		bs := make([][]byte, 2)
		bs[0] = make([]byte, 4)
		bs[1], err = proto.Marshal(p)

		if err != nil {
			return nil, err
		}

		if len(bs[1]) <= 0 || len(bs[1]) >= math.MaxInt32 {
			return nil, errors.New("invalid packet body")
		}

		binary.BigEndian.PutUint32(bs[0], uint32(len(bs[1])))

		return bs, nil
	}
}

func (l LengthFieldBasedFrameCodec) Decode(s *BrokerServer, c gnet.Conn) ([]*pb.Packet, error) {

	result := make([]*pb.Packet, 0)

	user, err := state.CurrentUserFromConn(c)
	if err != nil {
		logger.Error(err)
	}

	logger.InfoF("[%s#%s] Decode length field,buffer:%d",
		c.RemoteAddr().String(),
		user.Label(),
		c.InboundBuffered())

	for c.InboundBuffered() >= 4 {

		bs, e1 := c.Next(4)
		if e1 != nil {

			logger.ErrorF("[%s#%s] Decode length field,buffer:%d,error:%v",
				c.RemoteAddr().String(),
				user.Label(),
				c.InboundBuffered(),
				e1)

			c.Discard(c.InboundBuffered())
			//continue
			return nil, e1
		}

		l := binary.BigEndian.Uint32(bs)

		logger.InfoF("[%s#%s] Decode length field,buffer:%d,value:%d",
			c.RemoteAddr().String(),
			user.Label(),
			c.InboundBuffered(),
			l)

		if l <= 0 || l >= math.MaxInt32 {

			ee := errors.New("invalid decode length")
			logger.ErrorF("[%s#%s] Decode length field,buffer:%d,error:%v",
				c.RemoteAddr().String(),
				user.Label(),
				c.InboundBuffered(),
				ee)

			c.Discard(c.InboundBuffered())

			return nil, ee

		}

		length := int(l)

		if length == 4 && c.InboundBuffered() >= length {

			b, e2 := c.Next(length)

			if e2 != nil {

				logger.ErrorF("[%s#%s] Decode heartbeat body,buffer:%d,error:%v",
					c.RemoteAddr().String(),
					user.Label(),
					c.InboundBuffered(),
					e2)

				c.Discard(c.InboundBuffered())

				//continue
				return nil, e2
			}

			heartbeat := binary.BigEndian.Uint32(b)

			logger.InfoF("[%s#%s] Decode heartbeat body,buffer:%d,value:%d",
				c.RemoteAddr().String(),
				user.Label(),
				c.InboundBuffered(),
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

				logger.ErrorF("[%s#%s] Decode packet body,buffer:%d,error:%v",
					c.RemoteAddr().String(),
					user.Label(),
					c.InboundBuffered(),
					e3)

				c.Discard(c.InboundBuffered())

				//continue
				return nil, e3
			}

			var p pb.Packet
			if e4 := proto.Unmarshal(b, &p); e4 != nil {

				logger.ErrorF("[%s#%s] Decode packet body,unmarshal packet,buffer:%d,error:%v",
					c.RemoteAddr().String(),
					user.Label(),
					c.InboundBuffered(),
					e4)

				c.Discard(c.InboundBuffered())

				return nil, e4
				//continue
			}

			jb, _ := json.Marshal(p)

			logger.InfoF("[%s#%s] Decode packet body,buffer:%d,value:%s",
				c.RemoteAddr().String(),
				user.Label(),
				c.InboundBuffered(),
				jb)

			result = append(result, &p)

		}

	}

	return result, nil
}
