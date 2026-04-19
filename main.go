package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	originUrlEnv     string = "GO_DOCS_ORIGIN"
	remoteUrlEnv     string = "GO_DOCS_REMOTE"
	defaultOriginUrl string = "http://127.0.0.1:8888"
	defaultRemoteUrl string = "https://pkg.go.dev"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	origin := urlFromEnv(originUrlEnv, defaultOriginUrl)
	remote := urlFromEnv(remoteUrlEnv, defaultRemoteUrl)

	router := gin.Default()
	router.GET("/github.com/private/*any", forwarder(origin))
	router.NoRoute(forwarder(remote))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	stop()

	// 5 second grace period on shutdown request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

func urlFromEnv(env string, defaultUrl string) *url.URL {
	urlString, ok := os.LookupEnv(env)
	if !ok {
		urlString = defaultUrl
	}
	url, err := url.Parse(urlString)
	if err != nil {
		panic(err)
	}
	return url
}

func forwarder(url *url.URL) gin.HandlerFunc {
	return func(c *gin.Context) {
		// pkgsite context gets stuck if context is cancelled before a certain point
		// of the pkgsite's go proxy client call, causing all subsequent requests to
		// pkgsite to return 500 without attempting the go proxy.
		//
		// manually overriding the context here to avoid that issue - preferable to
		// use a bit more resources than to break the page until we restart the pod.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		proxy := httputil.NewSingleHostReverseProxy(url)
		proxy.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = url.Host
			req.URL.Host = url.Host
			req.URL.Scheme = url.Scheme
		}
		proxy.ServeHTTP(c.Writer, c.Request.WithContext(ctx))
	}
}
