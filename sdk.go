package sdk

type CharactrAPISDK struct {
	credentials *Credentials
	TTS         *tts
	VC          *vc
}

func New(credentials *Credentials) *CharactrAPISDK {
	return &CharactrAPISDK{
		credentials: credentials,
		TTS:         newTTS(credentials),
		VC:          newVC(credentials),
	}
}
