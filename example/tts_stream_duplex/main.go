package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"golang.org/x/sync/errgroup"

	"github.com/charactr-platform/charactr-api-sdk-go/v2"
	"github.com/charactr-platform/charactr-api-sdk-go/v2/example"
)

func main() {
	ctx := context.TODO()

	sdk := charactr.New(&example.Credentials)

	voices, err := sdk.TTS.GetVoices(ctx)
	if err != nil {
		panic(err)
	}

	var bytesWritten = 0

	file, err := os.Create("./result_tts_stream_duplex.wav")
	if err != nil {
		panic(err)
	}

	stream, err := sdk.TTS.StartDuplexStream(ctx, voices[0].ID)
	if err != nil {
		panic(err)
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		for {
			audioBytes, err := stream.Read()
			if err != nil {
				if err == io.EOF { // normal closure
					return nil
				}
				return err
			}

			n, err := file.Write(audioBytes)
			if err != nil {
				return err
			}

			bytesWritten += n
		}
	})

	stream.Convert("Hello world from the charactr TTS Duplex Streaming.")
	stream.Convert("You can send as much text as you want asynchronously.")

	err = stream.Close()
	if err != nil {
		panic(err)
	}

	err = g.Wait()
	fileErr := file.Close()
	if fileErr != nil {
		panic(fileErr)
	}
	if err != nil {
		panic(err)
	}

	fmt.Println("result_tts_stream_duplex.wav has been saved.")
	fmt.Println("Size: ", bytesWritten, "bytes")
}
