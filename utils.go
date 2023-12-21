package charactr

import (
	"context"
	"encoding/json"
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

	if res.StatusCode == http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(body, &result)
		if err != nil {
			return nil, err
		}

		return result, nil
	}

	return nil, getApiErr(res)
}

func getVoicesPaginated(ctx context.Context, url string, credentials *Credentials) ([]Voice, error) {
	type clonedVoicesRes struct {
		Items []Voice `json:"items"`
	}
	var result clonedVoicesRes

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

	if res.StatusCode == http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(body, &result)
		if err != nil {
			return nil, err
		}

		return result.Items, nil
	}

	return nil, getApiErr(res)
}
