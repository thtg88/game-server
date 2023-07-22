package grpcserver

import (
	"fmt"
	pb "game-server/internal/msgs/msg"
	"game-server/internal/player"
	"game-server/internal/server"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

// GAME
const DefaultPort = 4263

type GrpcRandomGameServerConfig struct {
	Host string
	Port uint16
}

type GrpcRandomGameServer struct {
	pb.UnimplementedGameServer
	RandomGameServer *server.RandomGameServer
	Config           GrpcRandomGameServerConfig
}

func New() *GrpcRandomGameServer {
	return &GrpcRandomGameServer{
		RandomGameServer: server.New(),
		Config: GrpcRandomGameServerConfig{
			Port: DefaultPort,
		},
	}
}

func Serve() error {
	rrgs := New()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", rrgs.Config.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterGameServer(s, rrgs)

	log.Printf("server listening at %v", lis.Addr())

	rrgs.RandomGameServer.Loop()

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

func (rrgs *GrpcRandomGameServer) Play(stream pb.Game_PlayServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		player := &player.Player{
			GameOverCh: make(chan bool),
			ID:         in.Player.Id,
			Level:      in.Player.Level,
			MessagesCh: make(chan string),
		}

		rrgs.RandomGameServer.Join(player)

		gameOver := false
		for !gameOver {
			select {
			case logMsg := <-player.MessagesCh:
				resp := &pb.PlayReply{Message: logMsg}

				if err := stream.Send(resp); err != nil {
					return err
				}
			case <-player.GameOverCh:
				resp := &pb.PlayReply{Message: "game over!"}

				if err := stream.Send(resp); err != nil {
					return err
				}
				gameOver = true
			}
		}
	}
}
