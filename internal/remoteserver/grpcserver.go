package remoteserver

import (
	"fmt"
	"game-server/internal/msgs"
	"game-server/internal/player"
	"game-server/internal/server"
	"log"
	"net"

	"google.golang.org/grpc"
)

type GrpcRandomGameServer struct {
	msgs.UnimplementedGameServer
	RandomGameServer *server.RandomGameServer
	Config           RemoteServerConfig
}

func NewGrpcRandomGameServer() *GrpcRandomGameServer {
	return &GrpcRandomGameServer{
		RandomGameServer: server.New(),
		Config: RemoteServerConfig{
			Port: DefaultPort,
		},
	}
}

func (rrgs *GrpcRandomGameServer) Serve() error {
	address := fmt.Sprintf(":%d", rrgs.Config.Port)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	msgs.RegisterGameServer(server, rrgs)

	log.Printf("gRPC server listening at %v", lis.Addr())

	rrgs.RandomGameServer.Loop()

	if err := server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

func (rrgs *GrpcRandomGameServer) Play(req *msgs.PlayRequest, stream msgs.Game_PlayServer) error {
	player := &player.Player{
		GameOverCh: make(chan bool),
		ID:         req.Player.Id,
		Level:      req.Player.Level,
		MessagesCh: make(chan string),
	}

	rrgs.RandomGameServer.Join(player)

	gameOver := false
	for !gameOver {
		select {
		case logMsg := <-player.MessagesCh:
			resp := &msgs.PlayReply{Message: logMsg}

			if err := stream.Send(resp); err != nil {
				return err
			}
		case <-player.GameOverCh:
			resp := &msgs.PlayReply{Message: "game over!"}

			if err := stream.Send(resp); err != nil {
				return err
			}
			gameOver = true

			close(player.GameOverCh)
			close(player.MessagesCh)
		}
	}

	return nil
}
