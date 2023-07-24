package remoteserver

import (
	"fmt"
	"game-server/internal/msgs"
	"game-server/internal/player"
	"game-server/internal/server"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		GameOverCh: make(chan struct{}),
		ID:         req.Player.Id,
		Level:      req.Player.Level,
		MessagesCh: make(chan string),
	}

	err := rrgs.RandomGameServer.Join(player)
	if err != nil {
		return status.Error(codes.Unavailable, err.Error())
	}

	gameOver := false
	for !gameOver {
		select {
		case logMsg := <-player.MessagesCh:
			if err := rrgs.send(stream, logMsg); err != nil {
				return err
			}
		case <-player.GameOverCh:
			if err := rrgs.send(stream, "game over!"); err != nil {
				return err
			}

			gameOver = true
		}
	}

	return nil
}

func (rrgs *GrpcRandomGameServer) send(stream msgs.Game_PlayServer, message string) error {
	resp := &msgs.PlayReply{Message: message}

	return stream.Send(resp)
}
