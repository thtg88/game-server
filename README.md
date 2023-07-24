# Game Server

This Go program simulates a naïve game server, where random players of different levels join a server for a random amount of time, through clients.

Additionally, it provides gRPC and TCP socket wrappers, as a way to test networking, Protobuf, gRPC, TCP sockets, and Go channels capabilities.

## Requirements

Go v1.20, you can get it using [gvm](https://github.com/moovweb/gvm) with:

```bash
gvm install go1.20
```

## Development

You can test the program locally without any remote calls with:

```bash
go run cmd/local-game-simulation/main.go
```

Otherwise you can test a client-server setup by running the server in one terminal with:

```bash
go run cmd/remote-server/main.go
```

And any other number of clients in different terminal windows with:

```bash
go run cmd/remote-client/main.go
```

Additionally a command to spawn 10,000 clients against the server is available with:

```bash
go run cmd/remote-client-spawn/main.go
```

Be careful though, as this is quite resource-intensive and, if you over-do it, you may run out of available TCP ports on your machine :D
