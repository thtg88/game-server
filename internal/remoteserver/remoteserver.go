package remoteserver

type RemoteServer interface {
	Serve() error
}

type RemoteServerConfig struct {
	Host string
	Port uint16
}

// GAME
const DefaultPort = 4263
