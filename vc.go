package charactr

import (
	"bytes"
	"context"
	"fmt"
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

func (v *vc) GetVoices(ctx context.Context) ([]Voice, error) {
	return getVoices(ctx, fmt.Sprintf("%s/v1/vc/voices", sdkConfig.apiUrl), v.credentials)
}

func (v *vc) Convert(ctx context.Context, voiceID int, inputAudio []byte) (*AudioResponse, error) {
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

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/v1/vc/convert?voiceId=%d", sdkConfig.apiUrl, voiceID), &body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Client-Key", v.credentials.ClientKey)
	req.Header.Set("X-API-Key", v.credentials.APIKey)
	req.Header.Set("Content-Type", mp.FormDataContentType())
	req.Header.Set("User-Agent", "sdk-go")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusOK {
		duration, err := strconv.Atoi(res.Header.Get(audioDurationHeader))
		if err != nil {
			return nil, err
		}

		size, err := strconv.Atoi(res.Header.Get(audioSizeHeader))
		if err != nil {
			return nil, err
		}

		return &AudioResponse{
			DurationMs:  duration,
			SizeBytes:   size,
			ContentType: res.Header.Get("Content-Type"),
			Audio:       res.Body,
		}, nil
	}

	return nil, getApiErr(res)
}
