package remoteclient

const DefaultServerHost = "localhost"

type RemoteClient interface {
	Join() error
}
