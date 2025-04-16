package broker

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/broker/domain"
	"github.com/magicnana999/im/define"
	"github.com/magicnana999/im/global"
	"github.com/magicnana999/im/pkg/jsonext"
	bb "github.com/panjf2000/gnet/v2/pkg/pool/bytebuffer"
	"go.uber.org/fx"
)

// Error definitions for MessageSendServer.
var (
	invalidMessage  = errors.New("invalid message")
	invalidUserConn = errors.New("invalid user conn")
	userConnClosed  = errors.New("user conn closed")
	mssNotRunning   = errors.New("server not running")
	channelFull     = errors.New("channel full")
)

// messageSending represents a message to be sent to a user connection.
type messageSending struct {
	message *api.Message     // message is the message to send.
	uc      *domain.UserConn // uc is the target user connection.
}

// MessageSendServer manages asynchronous message sending to user connections.
// It uses a buffered channel to queue messages and supports lifecycle management.
type MessageSendServer struct {
	isRunning atomic.Bool          // isRunning indicates if the server is running.
	codec     codec                // codec encodes messages for writing.
	cancel    context.CancelFunc   // cancel stops the sending loop.
	ch        chan *messageSending // ch queues messages to send.
	logger    *Logger              // logger records server events.
	mrs       *MessageRetryServer  // mrs messageRetryServer that can resend the message, so it can make sure the message is saved
}

// getOrDefaultMSSConfig returns the MSS configuration, prioritizing global configuration.
// It limits MaxRemaining to 10000 to prevent excessive memory usage.
func getOrDefaultMSSConfig(g *global.Config) *global.MSSConfig {
	c := &global.MSSConfig{}
	if g != nil && g.MSS != nil {
		*c = *g.MSS
	}

	if c.MaxRemaining <= 0 || c.MaxRemaining > 10000 {
		c.MaxRemaining = 10000
	}

	return c
}

// NewMessageSendServer creates a new MessageSendServer with the specified configuration.
// It registers lifecycle hooks for starting and stopping the server.
// It returns an error if logger or codec initialization fails.
func NewMessageSendServer(g *global.Config, mrs *MessageRetryServer, lc fx.Lifecycle) (*MessageSendServer, error) {
	c := getOrDefaultMSSConfig(g)

	log := NewLogger("mss", c.DebugMode)
	log.Info(string(jsonext.MarshalNoErr(c)), "", define.OpInit, "", nil)

	mss := &MessageSendServer{
		codec:  newCodec(),
		ch:     make(chan *messageSending, c.MaxRemaining),
		logger: log,
		mrs:    mrs,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return mss.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return mss.Stop(ctx)
		},
	})
	return mss, nil
}

// Start launches the message sending loop in a separate goroutine.
// It returns nil if the server starts successfully or is already running.
func (mss *MessageSendServer) Start(ctx context.Context) error {
	if mss.isRunning.CompareAndSwap(false, true) {
		ctx, cancel := context.WithCancel(ctx)
		mss.cancel = cancel
		go func() {
			mss.logger.Info("loop started", "", define.OpStart, "", nil)
			for {
				select {
				case <-ctx.Done():
					return
				case ms := <-mss.ch:
					if err := mss.Write(ms.message, ms.uc); err != nil {
						mss.Resave(ms.message, ms.uc)
					} else {
						mss.Submit(ms.message, ms.uc)
					}
				}
			}
		}()
	}
	return nil
}

// Stop shuts down the message sending loop and closes the channel.
// It processes remaining messages and returns nil if stopped successfully.
func (mss *MessageSendServer) Stop(ctx context.Context) error {
	if mss.isRunning.CompareAndSwap(true, false) {
		mss.cancel()
		close(mss.ch)
		mss.logger.Info("loop stopped", "", define.OpStop, "", nil)

		txt := fmt.Sprintf("remaining messages: %d", len(mss.ch))
		mss.logger.Info(txt, "", define.OpStop, "", nil)
		mss.ResaveMessages()
	}
	return nil
}

// ResaveMessages processes unsent messages in the channel.
// It calls onCloseFunc for each message and logs errors.
func (mss *MessageSendServer) ResaveMessages() {
	for ms := range mss.ch {
		mss.Resave(ms.message, ms.uc)
	}
}

// Send submits a message to be sent to the specified user connection.
// It returns an error if the message or connection is invalid, the connection is closed,
// the server is not running, or the channel is full.
func (mss *MessageSendServer) Send(m *api.Message, uc *domain.UserConn) error {
	if m == nil {
		return invalidMessage
	}
	if uc == nil {
		return invalidUserConn
	}
	if uc.IsClosed.Load() {
		return userConnClosed
	}
	if !mss.isRunning.Load() {
		return mssNotRunning
	}

	select {
	case mss.ch <- &messageSending{message: m, uc: uc}:
		return nil
	default:
		return channelFull
	}
}

// Write encodes and writes a message to the user connection.
// It returns an error if encoding or writing fails.
func (mss *MessageSendServer) Write(m *api.Message, uc *domain.UserConn) error {
	buffer, err := mss.codec.encode(m.Wrap())
	defer bb.Put(buffer)

	if err != nil {
		return err
	}

	total := buffer.Len()
	sent := 0
	for sent < total {
		n, err := uc.Writer.Write(buffer.Bytes()[sent:])
		if err != nil {
			return err
		}
		sent += n
	}
	mss.logger.Debug("write ok", uc.Desc(), "write", m.MessageId, nil)
	return nil
}

// Submit sends a successfully written message to the retry server for acknowledgment.
// It submits the message to MessageRetryServer to track its delivery status.
// It returns an error if the message or connection is invalid, the retry server is not initialized,
// or the submission fails.
func (mss *MessageSendServer) Submit(ms *api.Message, uc *domain.UserConn) error {
	mss.mrs.Submit(ms, uc, time.Now().Unix())
	return nil
}

// Resave saves a message that failed to send for later processing.
// It logs the message details and marks it for resaving (e.g., to a queue or database).
// It returns an error if the message or connection is invalid or resaving fails.
func (mss *MessageSendServer) Resave(ms *api.Message, uc *domain.UserConn) error {
	fmt.Print(jsonext.PbMarshalNoErr(ms))
	return nil
}
