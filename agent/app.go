package agent

import (
	"fmt"
	"log"
	"net"

	"github.com/sharithg/siphon/agent/cli"
	"github.com/sharithg/siphon/agent/docker"
	"github.com/sharithg/siphon/agent/grpcserver"
	"github.com/sharithg/siphon/internal/config"
	pb "github.com/sharithg/siphon/internal/protogen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Application struct {
	Config      Config
	Cli         *cli.Cli
	Docker      *docker.Docker
	AgentConfig *config.AgentConfig
}

type Config struct {
	Addr string
	Env  string
}

func (app *Application) AuthInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

	if app.Config.Env == "local" {
		return handler(srv, ss)
	}
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return status.Error(codes.Unauthenticated, "missing metadata")
	}

	token, ok := md["authorization"]
	if !ok || len(token) == 0 {
		return status.Error(codes.Unauthenticated, "missing authorization token")
	}

	if token[0] != app.AgentConfig.Token {
		return status.Error(codes.Unauthenticated, "invalid token")
	}

	return handler(srv, ss)
}

func (app *Application) Run() error {

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0%s", app.Config.Addr))

	if err != nil {
		panic("error building server: " + err.Error())
	}

	s := grpc.NewServer(
		grpc.StreamInterceptor(app.AuthInterceptor),
	)
	pb.RegisterStreamingServiceServer(s, grpcserver.GrpcServer{})

	log.Printf("Server has started on %s, env: %s", app.Config.Addr, app.Config.Env)

	if err := s.Serve(listener); err != nil {
		panic("error building server: " + err.Error())
	}
	return nil
}
