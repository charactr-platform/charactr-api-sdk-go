package charactr

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"nhooyr.io/websocket"
)

// streamInactivityThresholdMs defines how long after inactivity the stream will be considered inactive
const streamInactivityThresholdMs = 5000

type DuplexStream struct {
	ctx      context.Context
	conn     *websocket.Conn
	metadata DuplexStreamMetadata
}

type DuplexStreamMetadata struct {
	mu               sync.Mutex
	isClosed         bool
	isCloseRequested bool
	lastActiveAt     time.Time
}

// Read returns the raw audio data converted by the server
func (v *DuplexStream) Read() ([]byte, error) {
	messageType, audioBytes, err := v.conn.Read(v.ctx)
	v.metadata.markStreamActivity()

	// check for normal closure
	{
		closeErr := &websocket.CloseError{}
		ok := errors.As(err, closeErr)

		if ok && closeErr.Code == websocket.StatusNormalClosure {
			v.metadata.setStreamClosed()
			return nil, io.EOF
		} else if err != nil {
			return nil, err
		}
	}

	if messageType != websocket.MessageBinary {
		return nil, ErrStreamUnknownMessageRecv
	}

	return audioBytes, err
}

// Convert asynchronously feeds the stream with the text to be converted to audio
func (v *DuplexStream) Convert(text string) error {
	v.metadata.markStreamActivity()

	if text == "" {
		return ErrEmptyText
	}

	if v.metadata.isStreamClosed() || v.metadata.isStreamCloseRequested() {
		return ErrStreamClosed
	}

	err := v.conn.Write(v.ctx, websocket.MessageText, []byte(fmt.Sprintf(`{ "type": "convert", "text": %q }`, text)))
	if err != nil {
		return err
	}

	return nil
}

// Wait will block the execution until there was 5 seconds of stream inactivity
func (v *DuplexStream) Wait() {
	ticker := time.NewTicker(500 * time.Millisecond)

	for {
		select {
		case <-ticker.C:
			if !v.metadata.isStreamActive() || v.metadata.isStreamClosed() || v.metadata.isStreamCloseRequested() {
				ticker.Stop()
			}
		}
	}
}

// Close requests the server to close the connection gracefully
func (v *DuplexStream) Close() error {
	if v.metadata.isStreamClosed() || v.metadata.isStreamCloseRequested() {
		return ErrStreamClosed
	}

	v.metadata.setStreamCloseRequested()

	return v.conn.Write(v.ctx, websocket.MessageText, []byte(`{ "type": "close" }`))
}

// Terminate ends the stream immediately. In most cases we advise to use Close() instead.
func (v *DuplexStream) Terminate() error {
	if v.metadata.isStreamClosed() {
		return ErrStreamClosed
	}

	v.metadata.setStreamClosed()

	return v.conn.Close(websocket.StatusNormalClosure, "client terminated the connection")
}

func (v *DuplexStream) authenticate(clientKey string, apiKey string) error {
	v.metadata.markStreamActivity()

	if v.metadata.isStreamClosed() || v.metadata.isStreamCloseRequested() {
		return ErrStreamClosed
	}

	err := v.conn.Write(v.ctx, websocket.MessageText, []byte(fmt.Sprintf(`{ "type": "authApiKey", "clientKey": "%s", "apiKey": "%s" }`, clientKey, apiKey)))
	if err != nil {
		return err
	}

	return nil
}

func (v *DuplexStreamMetadata) markStreamActivity() {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.lastActiveAt = time.Now()
}

func (v *DuplexStreamMetadata) msSinceStreamLastActive() int64 {
	v.mu.Lock()
	defer v.mu.Unlock()

	diff := time.Since(v.lastActiveAt)

	return diff.Milliseconds()
}

func (v *DuplexStreamMetadata) isStreamActive() bool {
	return v.msSinceStreamLastActive() < streamInactivityThresholdMs
}

func (v *DuplexStreamMetadata) isStreamClosed() bool {
	v.mu.Lock()
	defer v.mu.Unlock()

	return v.isClosed
}

func (v *DuplexStreamMetadata) setStreamClosed() {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.isClosed = true
}

func (v *DuplexStreamMetadata) isStreamCloseRequested() bool {
	v.mu.Lock()
	defer v.mu.Unlock()

	return v.isCloseRequested
}

func (v *DuplexStreamMetadata) setStreamCloseRequested() {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.isCloseRequested = true
}
