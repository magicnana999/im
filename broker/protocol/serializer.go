package protocol

type Serializer interface {
	Serialize(packet *Packet) ([]byte, error)
	Deserialize(data []byte) (*Packet, error)
}
