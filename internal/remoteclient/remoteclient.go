package remoteclient

import (
	"context"
	"fmt"
	"game-server/internal/client"
	pb "game-server/internal/msgs/msg"
	"game-server/internal/remoteserver"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const DefaultServerHost = "localhost"

type RemoteClient interface {
	JoinRemoteServer()
}

type RemoteRandomClient struct {
	RandomClient *client.RandomClient
	ServerConfig remoteserver.RemoteRandomGameServerConfig
}

func New() *RemoteRandomClient {
	return &RemoteRandomClient{
		RandomClient: client.New(),
		ServerConfig: remoteserver.RemoteRandomGameServerConfig{
			Host: DefaultServerHost,
			Port: remoteserver.DefaultPort,
		},
	}
}

func (rrc *RemoteRandomClient) Join() error {
	conn, err := grpc.Dial(
		rrc.grpcDialTarget(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewGameClient(conn)

	return rrc.play(client)
}

func (rrc *RemoteRandomClient) grpcDialTarget() string {
	return fmt.Sprintf("%s:%d", rrc.ServerConfig.Host, rrc.ServerConfig.Port)
}

func (rrc *RemoteRandomClient) play(client pb.GameClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	stream, err := client.Play(ctx)
	if err != nil {
		return fmt.Errorf("client.Play failed: %v", err)
	}

	waitc := make(chan struct{})

	req := &pb.PlayRequest{
		Player: &pb.Player{
			Id:    rrc.RandomClient.Player.ID,
			Level: rrc.RandomClient.Player.Level,
		},
	}

	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Printf("client.Play failed: %v", err)
				break
			}
			log.Println(in.Message)
		}
	}()

	if err := stream.Send(req); err != nil {
		return fmt.Errorf("client.Play: stream.Send(%v) failed: %v", req, err)
	}

	if err := stream.CloseSend(); err != nil {
		return err
	}

	<-waitc

	return nil
}
