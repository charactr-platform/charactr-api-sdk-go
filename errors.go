package charactr

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrClient  = fmt.Errorf("charactr-api: client error")
	ErrServer  = fmt.Errorf("charactr-api: server error")
	ErrUnknown = fmt.Errorf("charactr-api: unknown error")
)

type errResponse struct {
	Msg string `json:"message"`
}

func getApiErr(res *http.Response) error {
	if res.StatusCode < 400 {
		return nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrUnknown, err)
	}

	var errRes errResponse
	err = json.Unmarshal(body, &errRes)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrUnknown, err)
	}

	if res.StatusCode < 500 {
		return fmt.Errorf("%w: [%d] %w", ErrClient, res.StatusCode, errRes.Msg)
	}

	return fmt.Errorf("%w: [%d] %w", ErrServer, res.StatusCode, errRes.Msg)
}
