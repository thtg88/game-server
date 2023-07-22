# Game Server

This Go program simulates a naïve game server, where random players of different levels join a server for a random amount of time, through clients.

Additionally, it provides gRPC and TCP socket wrappers, as a way to test networking, Protobuf, gRPC, TCP sockets, and Go channels capabilities.

## Requirements

Go v1.20, you can get it using [gvm](https://github.com/moovweb/gvm) with:

```bash
gvm install go1.20
```

## Development

You can test the program using:

```bash
go run main.go
```
