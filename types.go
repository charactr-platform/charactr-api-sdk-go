package charactr

import "io"

type config struct {
	apiUrl string
}

type Credentials struct {
	ClientKey string
	APIKey    string
}

type VoiceLabel struct {
	Category string `json:"category"`
	Label    string `json:"label"`
}

type Voice struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	PreviewURL  string       `json:"previewUrl"`
	Labels      []VoiceLabel `json:"labels"`
}

type AudioResponse struct {
	DurationMs  int    `json:"durationMs"`
	SizeBytes   int    `json:"sizeBytes"`
	ContentType string `json:"contentType"`
	Audio       io.Reader
}

type errResponse struct {
	Msg string `json:"message"`
}
