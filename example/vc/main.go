package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/charactr-platform/charactr-api-sdk-go"
	"github.com/charactr-platform/charactr-api-sdk-go/example"
)

func main() {
	sdk := charactr.New(&example.Credentials)

	voices, err := sdk.VC.GetVoices(context.TODO())
	if err != nil {
		panic(err)
	}

	file, err := os.ReadFile("./example/vc/input.wav")
	if err != nil {
		panic(err)
	}

	result, err := sdk.VC.Convert(context.TODO(), voices[0].ID, file)
	if err != nil {
		panic(err)
	}

	audio, err := io.ReadAll(result.Audio)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("./result_vc.wav", audio, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("result_vc.wav has been saved.")
	fmt.Println("Type: ", result.ContentType)
	fmt.Println("Size: ", result.SizeBytes, "bytes")
	fmt.Println("Duration: ", result.DurationMs, "ms")
}
