package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/oklog/run"
)

const (
	default_port    string = ":48200"
	default_content string = "/var/tmp/rest-stub"
)

var (
	content = flag.String("content", default_content, "Path to directory containing stub content files")
	port    = flag.String("port", default_port, "TCP port of service")
)

func main() {
	var g run.Group
	var logger log.Logger

	flag.Parse()

	logger = log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	{
		// Handle system interrupts
		cancel := make(chan struct{})
		g.Add(
			func() error {
				c := make(chan os.Signal, 1)
				signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
				select {
				case sig := <-c:
					return fmt.Errorf("received signal %s", sig)
				case <-cancel:
					return nil
				}
			},
			func(error) {
				close(cancel)
			})
	}

	{
		// HTTP listener
		listener, err := net.Listen("tcp", *port)
		if err != nil {
			logger.Log("during", "Listen", "error", err)
			os.Exit(1)
		}
		g.Add(
			func() error {
				logger.Log("event", "start", "port", *port, "content-dir", *content)
				handler := MiddlewareChain(
					RequestId(), // Enables identification of each request
					LogRequest(logger),
					Delay(logger),
					MimeType(logger),
					ResponseContent(),
					Charset(),
				).Then(RequestHandler(*content, logger))
				return http.Serve(listener, handler)
			},
			func(error) {
				listener.Close()
			})
	}

	logger.Log("event", "stop", "exit", g.Run())
}
