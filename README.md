# charactr-api-sdk-go

Go SDK to interact with the charactr API.

## Terminology
**VC** - *Voice conversion* - converting one voice from audio input to another voice.

**TTS** - *Text to speech* - converting text to voice audio.

## Features

- making TTS requests
- making VC requests
- managing Voice Clones
- getting lists of available voices

## Installation
```bash
go get github.com/charactr-platform/charactr-api-sdk-go/v2
```

## Usage

For the detailed SDK usage, please refer to the [SDK reference](https://docs.api.gemelo.ai/reference/go) or the `./example` directory.

### Quick start
```go
import "github.com/charactr-platform/charactr-api-sdk-go/v2"

func main() {
    sdk := charactr.New(&charactr.Credentials{ClientKey: "", APIKey: ""})
}
```

## How to run examples

#### Clone the SDK locally
```bash
git clone https://github.com/charactr-platform/charactr-api-sdk-go/v2
```

#### Provide credentials
Open `./example/credentials.go` and provide your credentials. You can find them in your [Client Panel](https://app.gemelo.ai) account.

#### Use TTS
```bash
go run example/tts/main.go
```

#### Use TTS Simplex streaming
```bash
go run example/tts_stream_simplex/main.go
```

#### Use TTS Duplex streaming
```bash
go run example/tts_stream_duplex/main.go
```

#### Use VC
```bash
go run example/vc/main.go
```

#### Use Voice Clones
```bash
go run example/clones/main.go
```
