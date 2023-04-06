package charactr

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func getVoices(ctx context.Context, url string, credentials *Credentials) ([]Voice, error) {
	var result []Voice

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Client-Key", credentials.ClientKey)
	req.Header.Set("X-API-Key", credentials.APIKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusOK {
		err = json.Unmarshal(body, &result)
		if err != nil {
			return nil, err
		}

		return result, nil
	}

	var errRes errResponse
	err = json.Unmarshal(body, &errRes)
	if err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("CharactrAPI request has failed with code %d: %s", res.StatusCode, errRes.Msg)
}
