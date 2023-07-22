package remoteserver

import (
	"fmt"
	"game-server/internal/msgs"
	"game-server/internal/player"
	"game-server/internal/server"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
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

	kaep := keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
		PermitWithoutStream: true,            // Allow pings even when there are no active streams
	}

	kasp := keepalive.ServerParameters{
		MaxConnectionIdle:     20 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
		MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
		MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
		Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
		Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
	}

	server := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(kaep), grpc.KeepaliveParams(kasp))
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

	rrgs.RandomGameServer.Join(player)

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
