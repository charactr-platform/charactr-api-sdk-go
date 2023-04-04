package main

import (
	"fmt"
	"io"
	"os"

	CharactrSDK "github.com/charactr-platform/charactr-api-sdk-go"
	"github.com/charactr-platform/charactr-api-sdk-go/example"
)

func main() {
	sdk := CharactrSDK.New(&example.Credentials)

	voices, err := sdk.TTS.GetVoices()
	if err != nil {
		panic(err)
	}

	result, err := sdk.TTS.Convert(voices[0].ID, "Hello world")
	if err != nil {
		panic(err)
	}

	audio, err := io.ReadAll(result.Audio)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("./result_tts.wav", audio, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("result_tts.wav has been saved.")
	fmt.Println("Type: ", result.Type)
	fmt.Println("Size: ", result.SizeBytes, "bytes")
	fmt.Println("Duration: ", result.DurationMs, "ms")
}
