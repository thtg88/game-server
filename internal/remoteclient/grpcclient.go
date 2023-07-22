package remoteclient

import (
	"context"
	"fmt"
	"game-server/internal/client"
	"game-server/internal/msgs"
	"game-server/internal/remoteserver"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type GrpcRandomClient struct {
	RandomClient *client.RandomClient
	ServerConfig remoteserver.RemoteServerConfig
}

func NewGrpcRandomClient() *GrpcRandomClient {
	return &GrpcRandomClient{
		RandomClient: client.New(),
		ServerConfig: remoteserver.RemoteServerConfig{
			Host: DefaultServerHost,
			Port: remoteserver.DefaultPort,
		},
	}
}

func (rrc *GrpcRandomClient) Join() error {
	kacp := keepalive.ClientParameters{
		Time:    5 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout: time.Second,     // wait 1 second for ping ack before considering the connection dead
		// PermitWithoutStream: true,             // send pings even without active streams
	}

	conn, err := grpc.Dial(
		rrc.grpcDialTarget(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(kacp),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := msgs.NewGameClient(conn)

	return rrc.play(client)
}

func (rrc *GrpcRandomClient) grpcDialTarget() string {
	return fmt.Sprintf("%s:%d", rrc.ServerConfig.Host, rrc.ServerConfig.Port)
}

func (rrc *GrpcRandomClient) play(client msgs.GameClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &msgs.PlayRequest{
		Player: &msgs.Player{
			Id:    rrc.RandomClient.Player.ID,
			Level: rrc.RandomClient.Player.Level,
		},
	}

	stream, err := client.Play(ctx, req)
	if err != nil {
		return fmt.Errorf("client.Play failed: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("client.Play failed: %v", err)
		}

		log.Println(resp.Message)
	}

	return nil
}
