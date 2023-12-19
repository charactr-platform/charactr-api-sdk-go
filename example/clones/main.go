package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/charactr-platform/charactr-api-sdk-go"
	"github.com/charactr-platform/charactr-api-sdk-go/example"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx := context.TODO()
	sdk := charactr.New(&example.Credentials)

	// create/update/load cloned voices
	var voiceID int
	{
		file, err := os.ReadFile("./example/clones/input.opus") // minimum 10 seconds
		if err != nil {
			panic(err)
		}
		voice, err := sdk.VoiceClone.CreateClonedVoice(ctx, "sdk-example-voice", bytes.NewReader(file))
		if err != nil {
			panic(err)
		}
		fmt.Printf("Voice created\n")

		voice, err = sdk.VoiceClone.UpdateClonedVoice(ctx, voice.ID, charactr.UpdateVoiceCloneDto{
			Name: "sdk-example-voice-renamed",
		})
		if err != nil {
			panic(err)
		}
		voiceID = voice.ID
		fmt.Printf("Voice updated\n")

		voices, err := sdk.VoiceClone.GetClonedVoices(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Number of cloned voices: %d\n", len(voices))
	}

	// tts
	{
		result, err := sdk.TTS.Convert(ctx, voiceID, "Hello world", &charactr.TTSRequestOptions{ClonedVoice: true})
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
	}
	// vc
	{
		file, err := os.ReadFile("./example/vc/input.wav")
		if err != nil {
			panic(err)
		}

		result, err := sdk.VC.Convert(ctx, voiceID, file, &charactr.VCRequestOptions{ClonedVoice: true})
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
	}
	// tts streaming
	{
		var bytesWritten = 0

		file, err := os.Create("./result_tts_stream_simplex.wav")
		if err != nil {
			panic(err)
		}

		stream, err := sdk.TTS.StartSimplexStream(ctx, voiceID, "Hello world from the charactr TTS Simplex Streaming.", &charactr.TTSStreamingOptions{ClonedVoice: true})
		if err != nil {
			panic(err)
		}

		g, _ := errgroup.WithContext(ctx)

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

	// delete test voice
	{
		err := sdk.VoiceClone.DeleteClonedVoice(ctx, voiceID)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Voice deleted\n")
	}
}
