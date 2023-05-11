package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"golang.org/x/sync/errgroup"

	"github.com/charactr-platform/charactr-api-sdk-go"
	"github.com/charactr-platform/charactr-api-sdk-go/example"
)

func main() {
	ctx := context.TODO()

	sdk := charactr.New(&example.Credentials)

	voices, err := sdk.TTS.GetVoices(ctx)
	if err != nil {
		panic(err)
	}

	var bytesWritten = 0

	file, err := os.Create("./result_tts_stream_simplex.wav")
	if err != nil {
		panic(err)
	}

	stream, err := sdk.TTS.StartSimplexStream(ctx, voices[0].ID, "Hello world from the charactr TTS Simplex Streaming.")
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

	err = g.Wait()
	fileErr := file.Close()
	if fileErr != nil {
		panic(fileErr)
	}
	if err != nil {
		panic(err)
	}

	fmt.Println("result_tts_stream_simplex.wav has been saved.")
	fmt.Println("Size: ", bytesWritten, "bytes")
}
