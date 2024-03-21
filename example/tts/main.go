package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/charactr-platform/charactr-api-sdk-go/v2"
	"github.com/charactr-platform/charactr-api-sdk-go/v2/example"
)

func main() {
	sdk := charactr.New(&example.Credentials)

	voices, err := sdk.TTS.GetVoices(context.TODO())
	if err != nil {
		panic(err)
	}

	result, err := sdk.TTS.Convert(context.TODO(), voices[0].ID, "Hello world")
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
	fmt.Println("Type: ", result.ContentType)
	fmt.Println("Size: ", result.SizeBytes, "bytes")
	fmt.Println("Duration: ", result.DurationMs, "ms")
}
