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

func (rrc *RemoteRandomClient) Join() {
	conn, err := grpc.Dial(
		rrc.grpcDialTarget(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewGameClient(conn)

	rrc.play(client)
}

func (rrc *RemoteRandomClient) grpcDialTarget() string {
	return fmt.Sprintf("%s:%d", rrc.ServerConfig.Host, rrc.ServerConfig.Port)
}

func (rrc *RemoteRandomClient) play(client pb.GameClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	stream, err := client.Play(ctx)
	if err != nil {
		log.Fatalf("client.Play failed: %v", err)
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
				log.Fatalf("client.Play failed: %v", err)
			}
			log.Println(in.Message)
		}
	}()

	if err := stream.Send(req); err != nil {
		log.Fatalf("client.Play: stream.Send(%v) failed: %v", req, err)
	}

	stream.CloseSend()

	<-waitc
}
