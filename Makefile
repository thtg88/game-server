build/protobuf:
	protoc --go_out=internal/msgs --go_opt=paths=source_relative --go-grpc_out=internal/msgs --go-grpc_opt=paths=source_relative msg/game.proto
