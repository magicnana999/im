package pb

import (
	"github.com/magicnana999/im/broker/protocol"
	"google.golang.org/protobuf/proto"
)

type ProtoSerializer struct {
}

var ProtoSerializerInstance = &ProtoSerializer{}

func (p ProtoSerializer) Serialize(packet *protocol.Packet) ([]byte, error) {
	ret, e := ConvertPacket(packet)
	if e != nil {
		return nil, e
	}
	return proto.Marshal(ret)
}

func (p ProtoSerializer) Deserialize(data []byte) (*protocol.Packet, error) {
	packet := new(Packet)
	if err := proto.Unmarshal(data, packet); err != nil {
		return nil, err
	}

	ret, e := RevertPacket(packet)
	if e != nil {
		return nil, e
	}
	return ret, nil
}
