package protocol

type Serializer interface {
	Serialize(packet *Packet) ([]byte, error)
	Deserialization(data []byte) (*Packet, error)
}
