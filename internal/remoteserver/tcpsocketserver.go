package remoteserver

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"game-server/internal/msgs"
	"game-server/internal/player"
	"game-server/internal/server"
	"log"
	"net"
	"sync"

	"google.golang.org/protobuf/proto"
)

type TcpSocketRandomGameServer struct {
	RandomGameServer *server.RandomGameServer
	Config           RemoteServerConfig
}

func NewTcpSocketRandomGameServer() *TcpSocketRandomGameServer {
	return &TcpSocketRandomGameServer{
		RandomGameServer: server.New(),
		Config: RemoteServerConfig{
			Port: DefaultPort,
		},
	}
}

func (rrgs *TcpSocketRandomGameServer) Serve() error {
	address := fmt.Sprintf(":%d", rrgs.Config.Port)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	defer lis.Close()

	log.Printf("TCP server listening at %v", lis.Addr())

	rrgs.RandomGameServer.Loop()

	var wg sync.WaitGroup

	defer wg.Wait()

	for {
		conn, err := lis.Accept()
		if err != nil {
			return fmt.Errorf("failed to accept: %v", err)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			rrgs.serve(conn)
		}()
	}
}

func (rrgs *TcpSocketRandomGameServer) serve(conn net.Conn) error {
	req, err := recv(conn)
	if err != nil {
		return err
	}

	player := &player.Player{
		GameOverCh: make(chan bool),
		ID:         req.Player.Id,
		Level:      req.Player.Level,
		MessagesCh: make(chan string),
	}

	rrgs.RandomGameServer.Join(player)

	for {
		select {
		case logMsg := <-player.MessagesCh:
			resp := &msgs.PlayReply{Message: logMsg}

			if err := send(conn, resp); err != nil {
				return err
			}
		case <-player.GameOverCh:
			resp := &msgs.PlayReply{Message: "game over!"}

			if err := send(conn, resp); err != nil {
				return err
			}

			close(player.GameOverCh)
			close(player.MessagesCh)

			return conn.Close()
		}
	}
}

func recv(conn net.Conn) (*msgs.PlayRequest, error) {
	netData, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return nil, err
	}

	return deserializeRequest(netData)
}

func deserializeRequest(message string) (*msgs.PlayRequest, error) {
	bytes, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return nil, err
	}

	var req msgs.PlayRequest

	err = proto.Unmarshal(bytes, &req)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func send(conn net.Conn, resp *msgs.PlayReply) error {
	bytes, err := serializeResponse(resp)
	if err != nil {
		return err
	}

	_, err = conn.Write(bytes)

	return err
}

func serializeResponse(resp *msgs.PlayReply) ([]byte, error) {
	marshaled, err := proto.Marshal(resp)
	if err != nil {
		return []byte{}, err
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(marshaled))

	bytes := []byte(fmt.Sprintf("%s\n", encoded))

	return bytes, nil
}
