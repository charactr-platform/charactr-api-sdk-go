package charactr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type vcc struct {
	credentials *Credentials
}

func newVCC(credentials *Credentials) *vcc {
	return &vcc{credentials: credentials}
}

// GetClonedVoices returns list of cloned voices, that client has created
func (v *vcc) GetClonedVoices(ctx context.Context) ([]Voice, error) {
	return getVoicesPaginated(ctx, fmt.Sprintf("%s/v1/cloned-voices?limit=500", sdkConfig.apiUrl), v.credentials) // TODO: remove 500 limit in v2
}

// CreateClonedVoice creates a new cloned voice with provided name, that sounds like provided audio
func (v *vcc) CreateClonedVoice(ctx context.Context, name string, audio io.Reader) (*Voice, error) {
	r, ct, err := newVoiceCloneBody(name, audio)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/v1/cloned-voices", sdkConfig.apiUrl), r)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Client-Key", v.credentials.ClientKey)
	req.Header.Set("X-API-Key", v.credentials.APIKey)
	req.Header.Set("User-Agent", "sdk-go")
	req.Header.Set("Content-Type", ct)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		return parseVoiceCloneResp(res)
	}

	return nil, getApiErr(res)
}

type UpdateVoiceCloneDto struct {
	Name string `json:"name"`
}

// UpdateClonedVoice updates the name of the clone
func (v *vcc) UpdateClonedVoice(ctx context.Context, id int, upd UpdateVoiceCloneDto) (*Voice, error) {
	reqBody, err := json.Marshal(&upd)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", fmt.Sprintf("%s/v1/cloned-voices/%d", sdkConfig.apiUrl, id), bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Client-Key", v.credentials.ClientKey)
	req.Header.Set("X-API-Key", v.credentials.APIKey)
	req.Header.Set("User-Agent", "sdk-go")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		return parseVoiceCloneResp(res)
	}

	return nil, getApiErr(res)
}

// DeleteClonedVoice removes the clone of given ID
func (v *vcc) DeleteClonedVoice(ctx context.Context, id int) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/v1/cloned-voices/%d", sdkConfig.apiUrl, id), nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Client-Key", v.credentials.ClientKey)
	req.Header.Set("X-API-Key", v.credentials.APIKey)
	req.Header.Set("User-Agent", "sdk-go")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNoContent || res.StatusCode == http.StatusOK {
		return nil
	}

	return getApiErr(res)
}

// newVoiceCloneBody creates multipart body, that is used when creating new voice clones
func newVoiceCloneBody(name string, audio io.Reader) (r io.Reader, contentType string, err error) {
	var body bytes.Buffer
	mp := multipart.NewWriter(&body)

	err = mp.WriteField("name", name)
	if err != nil {
		return nil, "", err
	}

	part, err := mp.CreateFormFile("audio", "input.wav")
	if err != nil {
		return nil, "", err
	}

	_, err = io.Copy(part, audio)
	if err != nil {
		return nil, "", err
	}

	err = mp.Close()
	if err != nil {
		return nil, "", err
	}

	return &body, mp.FormDataContentType(), nil
}

// parseVoiceCloneResp parses the successful response from voice clone endpoints
func parseVoiceCloneResp(res *http.Response) (*Voice, error) {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	result := Voice{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
