package charactr

import (
	"context"
	"errors"
	"fmt"
	"io"

	"nhooyr.io/websocket"
)

type SimplexStream struct {
	ctx  context.Context
	conn *websocket.Conn
}

func (v *SimplexStream) Read() ([]byte, error) {
	messageType, audioBytes, err := v.conn.Read(v.ctx)

	// check for normal closure
	{
		closeErr := &websocket.CloseError{}
		ok := errors.As(err, closeErr)

		if ok && closeErr.Code == websocket.StatusNormalClosure {
			return nil, io.EOF
		}
	}

	if messageType != websocket.MessageBinary {
		return nil, ErrStreamUnknownMessageRecv
	}

	return audioBytes, err
}

func (v *SimplexStream) authenticate(clientKey string, apiKey string) error {
	err := v.conn.Write(v.ctx, websocket.MessageText, []byte(fmt.Sprintf(`{ "type": "authApiKey", "clientKey": "%s", "apiKey": "%s" }`, clientKey, apiKey)))
	if err != nil {
		return err
	}

	return nil
}

func (v *SimplexStream) convert(text string) error {
	if text == "" {
		return ErrEmptyText
	}

	err := v.conn.Write(v.ctx, websocket.MessageText, []byte(fmt.Sprintf(`{ "type": "convert", "text": "%s" }`, text)))
	if err != nil {
		return err
	}

	return nil
}
