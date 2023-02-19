# Contributing to MongoPB

- [But Fixes and New Features](#bug-fixes-and-new-features)
- [Dependencies](#dependencies)
- [Testing](#testing)

Thank you for your interest in contributing to MongoPB! Please make sure to fork this repository before working through issues.

## Bug Fixes and New Features

1. Fork this repository
2. Create a pull request pointing to "main"
3. Add a reviewer

All pull requests are subject to the GitHub workflow CI defined in the Actions section of the repository.

## Dependencies

To develop locally you will need to install the following dependencies:

1. Go: https://go.dev/doc/install
2. Google protobuf compiler (protoc):

> ### Mac OS and Linux
>
> - http://google.github.io/proto-lens/installing-protoc.html

> ### Windows
>
> - Download the latest release (e.g., "protoc-21.8-win64.zip") under "Assets" https://github.com/protocolbuffers/protobuf/releases
>
> - Add to PATH by extracting to "C:\protoc-XX.X-winXX" (Be sure to replace 'X' with your appropriate release and system type)

3. protoc-gen-go: https://developers.google.com/protocol-buffers/docs/gotutorial#compiling-your-protocol-buffers
4. `gofumpt`: https://github.com/mvdan/gofumpt
5. `golangcli-lint`: https://github.com/golangci/golangci-lint#install-golangci-lint

## Testing

To test run

```
make tests
```
