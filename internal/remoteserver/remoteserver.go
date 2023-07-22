package remoteserver

type RemoteServer interface {
	Serve() error
}

// GAME
const DefaultPort = 4263
