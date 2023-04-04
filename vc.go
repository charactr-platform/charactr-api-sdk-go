package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

type vc struct {
	credentials *Credentials
}

func newVC(credentials *Credentials) *vc {
	return &vc{credentials: credentials}
}

func (v *vc) GetVoices() ([]Voice, error) {
	return getVoices(fmt.Sprintf("%s/v1/vc/voices", config.apiUrl), v.credentials)
}

func (v *vc) Convert(voiceID int, inputAudio []byte) (*AudioResponse, error) {
	var body bytes.Buffer
	mp := multipart.NewWriter(&body)

	part, err := mp.CreateFormFile("file", "input.wav")
	if err != nil {
		return nil, err
	}

	_, err = part.Write(inputAudio)
	if err != nil {
		return nil, err
	}

	err = mp.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/vc/convert?voiceId=%d", config.apiUrl, voiceID), &body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Client-Key", v.credentials.ClientKey)
	req.Header.Set("X-API-Key", v.credentials.APIKey)
	req.Header.Set("Content-Type", mp.FormDataContentType())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusOK {
		duration, err := strconv.Atoi(res.Header.Get("Audio-Duration-Ms"))
		if err != nil {
			return nil, err
		}

		size, err := strconv.Atoi(res.Header.Get("Audio-Size-Bytes"))
		if err != nil {
			return nil, err
		}

		return &AudioResponse{
			DurationMs: duration,
			SizeBytes:  size,
			Type:       res.Header.Get("Content-Type"),
			Audio:      res.Body,
		}, nil
	}

	errBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var errRes ErrResponse
	err = json.Unmarshal(errBody, &errRes)
	if err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("CharactrAPI request has failed with code %d: %s", res.StatusCode, errRes.Msg)
}
