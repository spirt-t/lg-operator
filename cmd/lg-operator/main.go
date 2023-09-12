package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	lgo "github.com/spirt-t/lg-operator/internal/app/api/lg-operator"
	"github.com/spirt-t/lg-operator/internal/cleaner"
	"github.com/spirt-t/lg-operator/internal/config"
	"github.com/spirt-t/lg-operator/internal/k8s"
	"github.com/spirt-t/lg-operator/internal/logger"
	desc "github.com/spirt-t/lg-operator/pkg/lg-operator"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

var (
	configFilePath = flag.String("cfg", "./config.yaml", "path to config file")
)

const (
	httpPortKey = "service.ports.http"
	grpcPortKey = "service.ports.grpc"
)

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	flag.Parse()

	cfgManager, err := config.NewManager(*configFilePath)
	if err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}

	lg, err := logger.NewLogger(cfgManager)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	k8sClient := k8s.NewClient(cfgManager, lg)
	if err = k8sClient.Init(ctx); err != nil {
		return fmt.Errorf("failed to make new k8s client: %w", err)
	}

	k8sManager, err := k8s.NewManager(k8sClient, cfgManager, lg)
	if err != nil {
		return fmt.Errorf("failed to make new k8s manager: %w", err)
	}

	cleaners := []lgo.Cleaner{
		cleaner.NewCompletedLGCleaner(cfgManager, k8sManager, lg),
		cleaner.NewOutdatedLGCleaner(cfgManager, k8sManager, lg),
	}

	service := lgo.NewService(k8sManager, cfgManager, lg, cleaners)
	service.RunCleaning(ctx)

	// serve
	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		listener, err := listen(cfgManager, grpcPortKey)
		if err != nil {
			return fmt.Errorf("failed to listen grpc port: %w", err)
		}
		return serveGRPC(gctx, listener, service)
	})

	g.Go(func() error {
		listener, err := listen(cfgManager, httpPortKey)
		if err != nil {
			return fmt.Errorf("failed to listen http port: %w", err)
		}
		return serveHTTP(gctx, listener, service)
	})

	return g.Wait()
}

func listen(cfgManager config.Manager, portKey string) (net.Listener, error) {
	var port int

	err := cfgManager.UnmarshalKey(portKey, &port)
	if err != nil {
		return nil, err
	}

	addr := ":" + strconv.Itoa(port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return listener, nil
}

func serveHTTP(ctx context.Context, listener net.Listener, service *lgo.Service) error {
	mux := runtime.NewServeMux()

	err := desc.RegisterLoadGeneratorOperatorServiceHandlerServer(ctx, mux, service)
	if err != nil {
		return fmt.Errorf("fail to register grpc-gateway handler: %w", err)
	}

	srv := &http.Server{Handler: mux}
	go func() {
		<-ctx.Done()
		_ = srv.Shutdown(ctx)
	}()

	log.Printf("Serving http address %s", listener.Addr())

	return srv.Serve(listener)
}

func serveGRPC(ctx context.Context, listener net.Listener, service *lgo.Service) error {
	baseGrpcServer := grpc.NewServer()
	desc.RegisterLoadGeneratorOperatorServiceServer(baseGrpcServer, service)

	log.Printf("Serving grpc address %s", listener.Addr())

	go func() {
		<-ctx.Done()
		baseGrpcServer.GracefulStop()
	}()

	return baseGrpcServer.Serve(listener)
}
