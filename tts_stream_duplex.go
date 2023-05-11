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

func (v *DuplexStream) Read() ([]byte, error) {
	messageType, audioBytes, err := v.conn.Read(v.ctx)
	v.markStreamActivity()

	// check for normal closure
	{
		closeErr := &websocket.CloseError{}
		ok := errors.As(err, closeErr)

		if ok && closeErr.Code == websocket.StatusNormalClosure {
			v.metadata.mu.Lock()
			defer v.metadata.mu.Unlock()
			v.metadata.isClosed = true
			return nil, io.EOF
		}
	}

	if messageType != websocket.MessageBinary {
		return nil, ErrStreamUnknownMessageRecv
	}

	return audioBytes, err
}

func (v *DuplexStream) Convert(text string) error {
	v.markStreamActivity()

	if text == "" {
		return ErrEmptyText
	}

	if v.metadata.isClosed || v.metadata.isCloseRequested {
		return ErrStreamClosed
	}

	err := v.conn.Write(v.ctx, websocket.MessageText, []byte(fmt.Sprintf(`{ "type": "convert", "text": "%s" }`, text)))
	if err != nil {
		return err
	}

	return nil
}

func (v *DuplexStream) Wait() {
	ticker := time.NewTicker(500 * time.Millisecond)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for {
			select {
			case <-ticker.C:
				if !v.isStreamActive() {
					ticker.Stop()
					wg.Done()
				}
			}
		}
	}()

	wg.Wait()
}

func (v *DuplexStream) Close() error {
	v.metadata.mu.Lock()
	defer v.metadata.mu.Unlock()

	if v.metadata.isClosed || v.metadata.isCloseRequested {
		return ErrStreamClosed
	}

	v.metadata.isCloseRequested = true

	return v.conn.Write(v.ctx, websocket.MessageText, []byte(`{ "type": "close" }`))
}

func (v *DuplexStream) Terminate() error {
	v.metadata.mu.Lock()
	defer v.metadata.mu.Unlock()

	if v.metadata.isClosed || v.metadata.isCloseRequested {
		return ErrStreamClosed
	}

	v.metadata.isClosed = true

	return v.conn.Close(websocket.StatusNormalClosure, "client terminated the connection")
}

func (v *DuplexStream) authenticate(clientKey string, apiKey string) error {
	v.markStreamActivity()

	if v.metadata.isClosed || v.metadata.isCloseRequested {
		return ErrStreamClosed
	}

	err := v.conn.Write(v.ctx, websocket.MessageText, []byte(fmt.Sprintf(`{ "type": "authApiKey", "clientKey": "%s", "apiKey": "%s" }`, clientKey, apiKey)))
	if err != nil {
		return err
	}

	return nil
}

func (v *DuplexStream) markStreamActivity() {
	v.metadata.mu.Lock()
	defer v.metadata.mu.Unlock()

	v.metadata.lastActiveAt = time.Now()
}

func (v *DuplexStream) msSinceStreamLastActive() int64 {
	v.metadata.mu.Lock()
	defer v.metadata.mu.Unlock()

	diff := time.Now().Sub(v.metadata.lastActiveAt)

	return diff.Milliseconds()
}

func (v *DuplexStream) isStreamActive() bool {
	return v.msSinceStreamLastActive() < streamInactivityThresholdMs
}
