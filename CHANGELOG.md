# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [v2.2.0] - 2024-03-21

- Update dependencies
- Fix v2 module name

## [v2.1.0] - 2023-12-20

- Added Voice Clone support
- Minor nil check fixes

## [v2.0.0] - 2023-08-02

- Changed API URL to [gemelo.ai](https://gemelo.ai)

## [v1.1.0] - 2023-07-26

- Changed User-Agent header to custom value to differentiate SDKs in the backend
- Exposed custom audio format and sample rate settings for TTS Streaming

## [v1.0.2] - 2023-06-07

- Fixed error handling when user passes invalid VoiceID in the TTS Streaming

## [v1.0.1] - 2023-05-26

- Added User-Agent header to the WebSocket handshake request

## [v1.0.0] - 2023-05-11

- Added TTS Simplex Streaming feature
- Added TTS Duplex Streaming feature

## [v0.0.1] - 2023-04-04

We implemented basic SDK features.

### Added

- TTS module
  - making TTS requests
  - fetching TTS voices
- VC module
  - making VC requests
  - fetching VC voices
