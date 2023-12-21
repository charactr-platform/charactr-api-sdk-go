package charactr

type CharactrAPISDK struct {
	credentials *Credentials
	TTS         *tts
	VC          *vc
	VoiceClone  *vcc
}

func New(credentials *Credentials) *CharactrAPISDK {
	return &CharactrAPISDK{
		credentials: credentials,
		TTS:         newTTS(credentials),
		VC:          newVC(credentials),
		VoiceClone:  newVCC(credentials),
	}
}
