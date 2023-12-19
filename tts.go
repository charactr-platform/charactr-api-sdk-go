package charactr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"nhooyr.io/websocket"
)

type tts struct {
	credentials *Credentials
}

type TTSStreamingOptions struct {
	Format      string
	SampleRate  int
	ClonedVoice bool
}

type TTSRequestOptions struct {
	ClonedVoice bool
}

func newTTS(credentials *Credentials) *tts {
	return &tts{credentials: credentials}
}

func (v *tts) GetVoices(ctx context.Context) ([]Voice, error) {
	return getVoices(ctx, fmt.Sprintf("%s/v1/tts/voices", sdkConfig.apiUrl), v.credentials)
}

func (v *tts) Convert(ctx context.Context, voiceID int, text string, options ...*TTSRequestOptions) (*AudioResponse, error) {
	type input struct {
		VoiceID   int    `json:"voiceId"`
		Text      string `json:"text"`
		VoiceType string `json:"voiceType"`
	}

	voiceType := "system"
	if len(options) > 0 && options[0] != nil && options[0].ClonedVoice {
		voiceType = "cloned"
	}

	in := input{
		VoiceID:   voiceID,
		Text:      text,
		VoiceType: voiceType,
	}

	reqBody, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/v1/tts/convert", sdkConfig.apiUrl), bytes.NewReader(reqBody))
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

func (v *tts) StartDuplexStream(ctx context.Context, voiceID int, options ...*TTSStreamingOptions) (*DuplexStream, error) {
	params := getTTSStreamingQueryParams(voiceID, options)

	ws, _, err := websocket.Dial(ctx, fmt.Sprintf("%s/v1/tts/stream/duplex/ws?%s", sdkConfig.wsApiUrl, params), &websocket.DialOptions{
		HTTPHeader: http.Header{"User-Agent": []string{"sdk-go"}},
	})
	if err != nil {
		return nil, err
	}

	stream := &DuplexStream{
		ctx:      ctx,
		conn:     ws,
		metadata: DuplexStreamMetadata{},
	}

	err = stream.authenticate(v.credentials.ClientKey, v.credentials.APIKey)
	if err != nil {
		return nil, err
	}

	return stream, nil
}

func (v *tts) StartSimplexStream(ctx context.Context, voiceID int, text string, options ...*TTSStreamingOptions) (*SimplexStream, error) {
	params := getTTSStreamingQueryParams(voiceID, options)

	ws, _, err := websocket.Dial(ctx, fmt.Sprintf("%s/v1/tts/stream/simplex/ws?%s", sdkConfig.wsApiUrl, params), &websocket.DialOptions{
		HTTPHeader: http.Header{"User-Agent": []string{"sdk-go"}},
	})
	if err != nil {
		return nil, err
	}

	stream := &SimplexStream{
		ctx:  ctx,
		conn: ws,
	}

	err = stream.authenticate(v.credentials.ClientKey, v.credentials.APIKey)
	if err != nil {
		return nil, err
	}

	err = stream.convert(text)
	if err != nil {
		return nil, err
	}

	return stream, nil
}

func getTTSStreamingQueryParams(voiceID int, options []*TTSStreamingOptions) string {
	params := url.Values{
		"voiceId": []string{strconv.Itoa(voiceID)},
	}

	if len(options) == 1 && options[0] != nil {
		opt := options[0]

		if opt.SampleRate != 0 && opt.Format == "" {
			opt.Format = "wav"
		}

		if opt.Format != "" {
			params.Set("format", opt.Format)
		}

		if opt.SampleRate != 0 {
			params.Set("sr", strconv.Itoa(opt.SampleRate))
		}

		if opt.ClonedVoice {
			params.Set("voiceType", "cloned")
		}
	}

	return params.Encode()
}
