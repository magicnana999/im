package api

func NewHeartbeatPacket(v int32) *Packet {
	return NewHeartbeat(v).Wrap()
}

func NewHeartbeat(v int32) *Heartbeat {
	return &Heartbeat{Value: v}
}

func (mb *Heartbeat) Wrap() *Packet {
	return &Packet{
		Type: TypeHeartbeat,
		Body: &Packet_Heartbeat{Heartbeat: mb},
	}
}
