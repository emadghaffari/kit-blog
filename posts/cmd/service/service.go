package service

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	endpoint1 "github.com/go-kit/kit/endpoint"
	log "github.com/go-kit/kit/log"
	prometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/hashicorp/vault/api"
	group "github.com/oklog/oklog/pkg/group"
	opentracinggo "github.com/opentracing/opentracing-go"
	prometheus1 "github.com/prometheus/client_golang/prometheus"
	promhttp "github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	grpc1 "google.golang.org/grpc"

	"github.com/emadghaffari/kit-blog/posts/config"
	endpoint "github.com/emadghaffari/kit-blog/posts/pkg/endpoint"
	grpc "github.com/emadghaffari/kit-blog/posts/pkg/grpc"
	pb "github.com/emadghaffari/kit-blog/posts/pkg/grpc/pb"
	service "github.com/emadghaffari/kit-blog/posts/pkg/service"
)

var tracer opentracinggo.Tracer
var logger log.Logger

// Define our flags. Your service probably won't need to bind listeners for
// all* supported transports, but we do it here for demonstration purposes.
var fs = flag.NewFlagSet("posts", flag.ExitOnError)
var debugAddr = fs.String("debug.addr", ":9680", "Debug and metrics listen address")
var httpAddr = fs.String("http-addr", ":9681", "HTTP listen address")
var grpcAddr = fs.String("grpc-addr", ":9682", "gRPC listen address")
var thriftAddr = fs.String("thrift-addr", ":9683", "Thrift listen address")
var thriftProtocol = fs.String("thrift-protocol", "binary", "binary, compact, json, simplejson")
var thriftBuffer = fs.Int("thrift-buffer", 0, "0 for unbuffered")
var thriftFramed = fs.Bool("thrift-framed", false, "true to enable framing")
var zipkinURL = fs.String("zipkin-url", "", "Enable Zipkin tracing via a collector URL e.g. http://localhost:9411/api/v1/spans")
var lightstepToken = fs.String("lightstep-token", "", "Enable LightStep tracing via a LightStep access token")
var appdashAddr = fs.String("appdash-addr", "", "Enable Appdash tracing via an Appdash server host:port")

// Run func
func Run() {
	fs.Parse(os.Args[1:])

	// Create a single logger, which we'll use and give to other components.
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	//  Determine which tracer to use. We'll pass the tracer to all the
	// components that use it, as a dependency
	// Sample configuration for testing. Use constant sampling to sample every trace
	// and enable LogSpan to log every span via configured Logger.
	cfg := jaegercfg.Configuration{
		ServiceName: "posts",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "jaeger:6831",
		},
	}

	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	var closer io.Closer
	var err error
	tracer, closer, err = cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
		jaegercfg.ZipkinSharedRPCSpan(true),
	)
	if err != nil {
		logger.Log("during", "Listen", "jaeger", "err", err)
	}
	opentracinggo.SetGlobalTracer(tracer)
	defer closer.Close()

	// init vault
	if err := initVault(); err != nil {
		logger.Log(err)
		return
	}
	svc := service.New(getServiceMiddleware(logger))
	eps := endpoint.New(svc, getEndpointMiddleware(logger))
	g := createService(eps)
	initMetricsEndpoint(g)
	initCancelInterrupt(g)
	logger.Log("exit", g.Run())

}
func initGRPCHandler(endpoints endpoint.Endpoints, g *group.Group) {
	options := defaultGRPCOptions(logger, tracer)
	// Add your GRPC options here

	grpcServer := grpc.NewGRPCServer(endpoints, options)
	grpcListener, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		logger.Log("transport", "gRPC", "during", "Listen", "err", err)
	}
	g.Add(func() error {
		logger.Log("transport", "gRPC", "addr", *grpcAddr)
		baseServer := grpc1.NewServer()
		pb.RegisterPostsServer(baseServer, grpcServer)
		return baseServer.Serve(grpcListener)
	}, func(error) {
		grpcListener.Close()
	})

}
func getServiceMiddleware(logger log.Logger) (mw []service.Middleware) {
	mw = []service.Middleware{}
	mw = addDefaultServiceMiddleware(logger, mw)
	// Append your middleware here

	return
}
func getEndpointMiddleware(logger log.Logger) (mw map[string][]endpoint1.Middleware) {
	mw = map[string][]endpoint1.Middleware{}
	duration := prometheus.NewSummaryFrom(prometheus1.SummaryOpts{
		Help:      "Request duration in seconds.",
		Name:      "request_duration_seconds",
		Namespace: "example",
		Subsystem: "posts",
	}, []string{"method", "success"})
	addDefaultEndpointMiddleware(logger, duration, mw)
	// Add you endpoint middleware here

	return
}
func initMetricsEndpoint(g *group.Group) {
	http.DefaultServeMux.Handle("/metrics", promhttp.Handler())
	debugListener, err := net.Listen("tcp", *debugAddr)
	if err != nil {
		logger.Log("transport", "debug/HTTP", "during", "Listen", "err", err)
	}
	g.Add(func() error {
		logger.Log("transport", "debug/HTTP", "addr", *debugAddr)
		return http.Serve(debugListener, http.DefaultServeMux)
	}, func(error) {
		debugListener.Close()
	})
}
func initCancelInterrupt(g *group.Group) {
	cancelInterrupt := make(chan struct{})
	g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-cancelInterrupt:
			return nil
		}
	}, func(error) {
		close(cancelInterrupt)
	})
}

func initVault() error {
	config.Confs.Vault.Address = "http://vault:8200"
	config.Confs.Vault.Token = "s.4TnB3ozvkZYlQRTLZAteVivl"
	config.Confs.Posts.Path = "blog/posts"
	config.Confs.Posts.DebugAddr = *debugAddr
	config.Confs.Posts.HTTPAddr = *httpAddr
	config.Confs.Posts.GrpcAddr = *grpcAddr
	config.Confs.Posts.ThriftAddr = *thriftAddr
	config.Confs.Posts.Host = "localhost"
	config.Confs.Users.Path = "blog/users"

	confs := &api.Config{
		Address: config.Confs.Vault.Address,
	}
	client, err := api.NewClient(confs)
	if err != nil {
		logger.Log(err)
		return err
	}
	client.SetToken(config.Confs.Vault.Token)
	c := client.Logical()
	config.Confs.Vault.Logical = c

	// Read notif path
	users, err := c.Read(config.Confs.Users.Path)
	if err != nil {
		logger.Log(err)
		return err
	}
	config.Confs.Users.GrpcAddr = users.Data["grpc"].(string)

	// Write Posts Path
	_, err = c.Write(config.Confs.Posts.Path, map[string]interface{}{
		"debug":  config.Confs.Posts.Host + config.Confs.Posts.DebugAddr,
		"http":   config.Confs.Posts.Host + config.Confs.Posts.HTTPAddr,
		"grpc":   config.Confs.Posts.Host + config.Confs.Posts.GrpcAddr,
		"thrift": config.Confs.Posts.Host + config.Confs.Posts.ThriftAddr,
	})
	if err != nil {
		logger.Log(err)
		return err
	}

	return nil
}
