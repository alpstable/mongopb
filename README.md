# MongoPB

[![PkgGoDev](https://img.shields.io/badge/go.dev-docs-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/alpstable/mongopb)
![Build Status](https://github.com/alpstable/mongopb/actions/workflows/ci.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/alpstable/mongopb)](https://goreportcard.com/report/github.com/alpstable/mongopb)
[![Discord](https://img.shields.io/discord/987810353767403550)](https://discord.gg/3jGYQz74s7)

MongoPB is a library for writing [structpb](https://pkg.go.dev/google.golang.org/protobuf/types/known/structpb#ListValue)-typed data to a MongoDB collection using the official [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver).

## Installation

```sh
go get github.com/alpstable/mongopb@latest
```

## Usage

The type `structpb` types supported by this package are

- [`ListValue`](https://pkg.go.dev/google.golang.org/protobuf/types/known/structpb#ListValue)

See the [gidari](https://github.com/alpstable/gidari) library to learn how to write to a MongoDB collection from a web API.

## Contributing

Follow [this guide](docs/CONTRIBUTING.md) for information on contributing.
