package core

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/magicnana999/im/broker/pb"
	"github.com/magicnana999/im/broker/protocol"
	"github.com/magicnana999/im/logger"
	"github.com/panjf2000/gnet/v2"
	"google.golang.org/protobuf/proto"
	"math"
)

type Codec interface {
	Encode(p *protocol.Packet) ([][]byte, error)
	Decode(s *BrokerServer, c gnet.Conn) ([]*protocol.Packet, error)
}

var defaultCodec = &LengthFieldBasedFrameCodec{}

type LengthFieldBasedFrameCodec struct {
}

func (l LengthFieldBasedFrameCodec) Encode(p *protocol.Packet) ([][]byte, error) {
	if p.IsHeartbeat() {

		if p.Body.(uint32) <= 0 || p.Body.(uint32) >= math.MaxInt32 {
			return nil, errors.New("invalid heartbeat body")
		}

		bs := make([][]byte, 2)
		bs[0] = make([]byte, 4)
		bs[1] = make([]byte, 4)

		binary.BigEndian.PutUint32(bs[0], uint32(4))
		binary.BigEndian.PutUint32(bs[1], p.Body.(uint32))
		return bs, nil
	} else {
		pp, err := pb.ConvertPacket(p)
		if err != nil {
			return nil, err
		}

		bs := make([][]byte, 2)
		bs[0] = make([]byte, 4)
		bs[1], err = proto.Marshal(pp)
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

func (l LengthFieldBasedFrameCodec) Decode(s *BrokerServer, c gnet.Conn) ([]*protocol.Packet, error) {

	result := make([]*protocol.Packet, 0)

	user, err := currentUserConnection(c)
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
				packet := &protocol.Packet{
					Type: protocol.TypeHeartbeat,
					Body: heartbeat,
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

			packet, e5 := pb.RevertPacket(&p)
			if e5 != nil {

				logger.ErrorF("[%s#%s] Decode packet body,revert packet,buffer:%d,error:%v",
					c.RemoteAddr().String(),
					user.Label(),
					c.InboundBuffered(),
					e5)

				c.Discard(c.InboundBuffered())

				return nil, e5
				//continue
			}

			jb, _ := json.Marshal(packet)

			logger.InfoF("[%s#%s] Decode packet body,buffer:%d,value:%s",
				c.RemoteAddr().String(),
				user.Label(),
				c.InboundBuffered(),
				jb)

			result = append(result, packet)

		}

	}

	return result, nil
}
