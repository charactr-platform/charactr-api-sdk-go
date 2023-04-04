# charactr-api-sdk-go

Go SDK to interact with the charactr API.

## Terminology
**VC** - *Voice conversion* - converting one voice from audio input to another voice.

**TTS** - *Text to speech* - converting text to voice audio.

## Features

- making TTS requests
- making VC requests
- getting lists of available voices

## Installation
```bash
$ go add github.com/charactr-platform/charactr-api-sdk-go
```

## Usage

For the detailed SDK usage, please refer to the [SDK reference](https://docs.api.charactr.com/reference/go) or the `./examples` directory.

## How to run examples

#### Clone the SDK locally
```bash
$ git clone https://github.com/charactr-platform/charactr-api-sdk-go
```

#### Use TTS
```bash
$ go run example/tts/main.go
```

#### Use VC
```bash
$ go run example/vc/main.go
```
