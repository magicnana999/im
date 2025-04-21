package broker

import (
	"encoding/binary"
	"errors"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/panjf2000/gnet/v2"
	bb "github.com/panjf2000/gnet/v2/pkg/pool/bytebuffer"
	"google.golang.org/protobuf/proto"
)

var (
	lengthError = errors.New("read length is less than field")
)

var defaultCodec = &Codec{}

type Codec struct {
}

func NewCodec() *Codec {
	return defaultCodec
}

func (l *Codec) Encode(p *api.Packet) (*bb.ByteBuffer, error) {

	buffer := bb.Get()

	if p.IsHeartbeat() {

		hb := p.GetHeartbeat().Value

		binary.Write(buffer, binary.BigEndian, int32(4))
		binary.Write(buffer, binary.BigEndian, int32(hb))

		return buffer, nil
	} else {

		bs, e := proto.Marshal(p)

		if e != nil {
			return nil, e
		}

		binary.Write(buffer, binary.BigEndian, int32(len(bs)))
		binary.Write(buffer, binary.BigEndian, bs)
		return buffer, nil
	}
}

func (l *Codec) Decode(c gnet.Conn) ([]*api.Packet, error) {

	result := make([]*api.Packet, 0)

	for c.InboundBuffered() >= 4 {

		var length int32
		if err := binary.Read(c, binary.BigEndian, &length); err != nil {
			return nil, err
		}

		if length == 4 && c.InboundBuffered() >= int(length) {

			var heartbeat int32
			if err := binary.Read(c, binary.BigEndian, &heartbeat); err != nil {
				return nil, err
			}

			packet := api.NewHeartbeat(int32(heartbeat)).Wrap()
			result = append(result, packet)
		}

		if length > 4 && c.InboundBuffered() >= int(length) {

			bs := make([]byte, int(length))
			n, e := c.Read(bs)
			if e != nil {
				return nil, e
			}

			if n != int(length) {
				return nil, lengthError
			}

			var p api.Packet
			if e4 := proto.Unmarshal(bs, &p); e4 != nil {
				return nil, e4
			}

			result = append(result, &p)

		}

	}

	return result, nil
}
