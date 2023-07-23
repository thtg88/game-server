package remoteclient

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"game-server/internal/client"
	"game-server/internal/msgs"
	"game-server/internal/remoteserver"
	"io"
	"log"
	"net"

	"google.golang.org/protobuf/proto"
)

type TcpSocketRandomClient struct {
	RandomClient *client.RandomClient
	ServerConfig remoteserver.RemoteServerConfig
}

func NewTcpSocketRandomClient() *TcpSocketRandomClient {
	return &TcpSocketRandomClient{
		RandomClient: client.New(),
		ServerConfig: remoteserver.RemoteServerConfig{
			Host: DefaultServerHost,
			Port: remoteserver.DefaultPort,
		},
	}
}

func (rrc *TcpSocketRandomClient) Join() error {
	addr := fmt.Sprintf("%s:%d", rrc.ServerConfig.Host, rrc.ServerConfig.Port)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	return rrc.play(conn)
}

func (rrc *TcpSocketRandomClient) play(conn net.Conn) error {
	waitc := make(chan struct{})

	req := &msgs.PlayRequest{
		Player: &msgs.Player{
			Id:    rrc.RandomClient.Player.ID,
			Level: rrc.RandomClient.Player.Level,
		},
	}

	go func() {
		for {
			resp, err := rrc.recv(conn)
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Printf("client.Play failed: %v", err)
				close(waitc)
				return
			}

			log.Println(resp.Message)
		}
	}()

	if err := rrc.send(conn, req); err != nil {
		return err
	}

	<-waitc

	return nil
}

func (rrc *TcpSocketRandomClient) send(conn net.Conn, req *msgs.PlayRequest) error {
	data, err := rrc.serializeRequest(req)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(conn, data)

	return err
}

func (rrc *TcpSocketRandomClient) serializeRequest(req *msgs.PlayRequest) (string, error) {
	data, err := proto.Marshal(req)
	if err != nil {
		return "", err
	}

	text := base64.StdEncoding.EncodeToString([]byte(data))

	return fmt.Sprintf("%s\n", text), nil
}

func (rrc *TcpSocketRandomClient) recv(conn net.Conn) (*msgs.PlayReply, error) {
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return nil, err
	}

	return rrc.deserializeResponse(message)
}

func (rrc *TcpSocketRandomClient) deserializeResponse(message string) (*msgs.PlayReply, error) {
	bytes, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return nil, err
	}

	var resp msgs.PlayReply

	err = proto.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
