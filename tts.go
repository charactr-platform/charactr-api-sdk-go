package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type tts struct {
	credentials *Credentials
}

func newTTS(credentials *Credentials) *tts {
	return &tts{credentials: credentials}
}

func (v *tts) GetVoices() ([]Voice, error) {
	return getVoices(fmt.Sprintf("%s/v1/tts/voices", config.apiUrl), v.credentials)
}

func (v *tts) Convert(voiceID int, text string) (*AudioResponse, error) {
	type input struct {
		VoiceID int    `json:"voiceId"`
		Text    string `json:"text"`
	}

	reqBody, err := json.Marshal(input{
		VoiceID: voiceID,
		Text:    text,
	})

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/tts/convert", config.apiUrl), bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Client-Key", v.credentials.ClientKey)
	req.Header.Set("X-API-Key", v.credentials.APIKey)

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

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var errRes ErrResponse
	err = json.Unmarshal(body, &errRes)
	if err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("CharactrAPI request has failed with code %d: %s", res.StatusCode, errRes.Msg)
}
