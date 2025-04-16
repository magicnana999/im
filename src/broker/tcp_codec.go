package broker

import (
	"encoding/binary"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/errors"
	"github.com/panjf2000/gnet/v2"
	bb "github.com/panjf2000/gnet/v2/pkg/pool/bytebuffer"
	"google.golang.org/protobuf/proto"
)

type codec interface {
	encode(p *api.Packet) (*bb.ByteBuffer, error)
	decode(c gnet.Conn) ([]*api.Packet, error)
}

var defaultCodec = &lengthFieldBasedFrameCodec{}

type lengthFieldBasedFrameCodec struct {
}

func newCodec() codec {
	return defaultCodec
}

func (l *lengthFieldBasedFrameCodec) encode(p *api.Packet) (*bb.ByteBuffer, error) {

	buffer := bb.Get()

	if p.IsHeartbeat() {

		hb := p.GetHeartbeat().Value

		binary.Write(buffer, binary.BigEndian, int32(4))
		binary.Write(buffer, binary.BigEndian, int32(hb))

		return buffer, nil
	} else {

		bs, e := proto.Marshal(p)

		if e != nil {
			return nil, errors.EncodeError.SetDetail(e)
		}

		binary.Write(buffer, binary.BigEndian, int32(len(bs)))
		binary.Write(buffer, binary.BigEndian, bs)
		return buffer, nil
	}
}

func (l *lengthFieldBasedFrameCodec) decode(c gnet.Conn) ([]*api.Packet, error) {

	result := make([]*api.Packet, 0)

	for c.InboundBuffered() >= 4 {

		var length int32
		binary.Read(c, binary.BigEndian, &length)

		if length == 4 && c.InboundBuffered() >= int(length) {

			var heartbeat int32
			binary.Read(c, binary.BigEndian, &heartbeat)

			packet := api.NewHeartbeat(int32(heartbeat)).Wrap()
			result = append(result, packet)

		}

		if length > 4 && c.InboundBuffered() >= int(length) {

			bs := make([]byte, int(length))
			n, e := c.Read(bs)
			if e != nil {
				return nil, errors.DecodeError.SetDetail(e)
			}

			if n != int(length) {
				return nil, errors.DecodeError.SetDetail("failed to read packet")
			}

			var p api.Packet
			if e4 := proto.Unmarshal(bs, &p); e4 != nil {
				return nil, errors.DecodeError.SetDetail(e4)
			}

			result = append(result, &p)

		}

	}

	return result, nil
}
