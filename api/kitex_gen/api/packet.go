package api

import "google.golang.org/protobuf/proto"

func (p *Packet) GetPacketId() string {
	switch p.Type {
	case TypeCommand:
		return p.GetCommand().GetCommandId()
	case TypeMessage:
		return p.GetMessage().GetMessageId()
	default:
		return ""
	}
}

func (p *Packet) IsHeartbeat() bool {
	return p.Type == TypeHeartbeat
}

func (p *Packet) IsCommand() bool {
	return p.Type == TypeCommand
}

func (p *Packet) IsMessage() bool {
	return p.Type == TypeMessage
}

func (p *Packet) Failure(e error) *Packet {
	switch p.Type {
	case TypeHeartbeat:
		return nil
	case TypeCommand:
		return p.GetCommand().Failure(e).Wrap()
	case TypeMessage:
		return p.GetMessage().Failure(e).Wrap()
	default:
		return nil
	}
}

func (p *Packet) Success(c proto.Message) *Packet {
	switch p.Type {
	case TypeHeartbeat:
		return nil
	case TypeCommand:
		return p.GetCommand().Success(c).Wrap()
	case TypeMessage:
		return p.GetMessage().Success(c).Wrap()
	default:
		return nil
	}
}
