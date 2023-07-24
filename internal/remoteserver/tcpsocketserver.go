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
		go func(conn net.Conn) {
			defer wg.Done()
			defer conn.Close()

			if err := rrgs.Play(conn); err != nil {
				log.Println(err)
			}
		}(conn)
	}
}

func (rrgs *TcpSocketRandomGameServer) Play(conn net.Conn) error {
	req, err := rrgs.recvRequest(conn)
	if err != nil {
		return err
	}

	player := &player.Player{
		GameOverCh: make(chan struct{}),
		ID:         req.Player.Id,
		Level:      req.Player.Level,
		MessagesCh: make(chan string),
	}

	err = rrgs.RandomGameServer.Join(player)
	if err != nil {
		// an error in joining means we are not accepting new players, so send the message to the client
		if sendErr := rrgs.sendResponse(conn, err.Error()); sendErr != nil {
			return sendErr
		}

		return err
	}

	for {
		select {
		case logMsg := <-player.MessagesCh:
			if err := rrgs.sendResponse(conn, logMsg); err != nil {
				return err
			}
		case <-player.GameOverCh:
			if err := rrgs.sendResponse(conn, "game over!"); err != nil {
				return err
			}

			return conn.Close()
		}
	}
}

func (rrgs *TcpSocketRandomGameServer) recvRequest(conn net.Conn) (*msgs.PlayRequest, error) {
	netData, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return nil, err
	}

	return rrgs.deserializeRequest(netData)
}

func (rrgs *TcpSocketRandomGameServer) deserializeRequest(message string) (*msgs.PlayRequest, error) {
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

func (rrgs *TcpSocketRandomGameServer) sendResponse(conn net.Conn, msg string) error {
	resp := &msgs.PlayReply{Message: msg}

	bytes, err := rrgs.serializeResponse(resp)
	if err != nil {
		return err
	}

	_, err = conn.Write(bytes)

	return err
}

func (rrgs *TcpSocketRandomGameServer) serializeResponse(resp *msgs.PlayReply) ([]byte, error) {
	marshaled, err := proto.Marshal(resp)
	if err != nil {
		return []byte{}, err
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(marshaled))

	bytes := []byte(fmt.Sprintf("%s\n", encoded))

	return bytes, nil
}
