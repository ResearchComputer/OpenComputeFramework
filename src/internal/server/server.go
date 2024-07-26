package server

import (
	"context"
	"net/http"
	"ocf/internal/common"
	"ocf/internal/protocol"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("gin-server")

func StartServer() {

	_, cancelCtx := protocol.GetCRDTStore()
	defer cancelCtx()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer stop()

	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(otelgin.Middleware(serviceName))
	r.Use(corsHeader())
	r.Use(gin.Recovery())
	go protocol.StartTicker()
	v1 := r.Group("/v1")
	{
		v1.GET("/health", healthStatusCheck)
		crdtGroup := v1.Group("/dnt")
		{
			crdtGroup.GET("/table", getDNT)
			crdtGroup.GET("/peers", listPeers)
			crdtGroup.POST("/_node", updateLocal)
			crdtGroup.DELETE("/_node", deleteLocal)
		}
		p2pGroup := v1.Group("/p2p")
		{
			p2pGroup.PATCH("/:peerId/*path", P2PForwardHandler)
			p2pGroup.POST("/:peerId/*path", P2PForwardHandler)
			p2pGroup.GET("/:peerId/*path", P2PForwardHandler)
		}
		globalServiceGroup := v1.Group("/service")
		{
			globalServiceGroup.GET("/:service/*path", GlobalServiceForwardHandler)
			globalServiceGroup.POST("/:service/*path", GlobalServiceForwardHandler)
			globalServiceGroup.PATCH("/:service/*path", GlobalServiceForwardHandler)
		}
		serviceGroup := v1.Group("/_service")
		{
			serviceGroup.GET("/:service/*path", ServiceForwardHandler)
			serviceGroup.POST("/:service/*path", ServiceForwardHandler)
			serviceGroup.PATCH("/:service/*path", ServiceForwardHandler)
		}
	}
	p2plistener := P2PListener()
	srv := &http.Server{
		Addr:    "0.0.0.0:" + viper.GetString("port"),
		Handler: r,
	}
	go func() {
		err := http.Serve(p2plistener, r)
		if err != nil {
			common.Logger.Error("http.Serve: %s", err)
		}
	}()
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			common.ReportError(err, "Server failed to start")
		}
	}()
	<-ctx.Done()
	// shutting down...
	protocol.DeleteNodeTable()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		common.ReportError(err, "Server shutdown failed")
	}
	common.Logger.Info("Server exiting")
}
