package server

import (
	"net/http"
	"ocf/internal/common"
	"ocf/internal/protocol"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func StartServer() {
	var wg sync.WaitGroup
	_, cancelCtx := protocol.GetCRDTStore()
	defer cancelCtx()
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(corsHeader())
	r.Use(gin.Recovery())

	v1 := r.Group("/v1")
	{
		v1.GET("/health", healthStatusCheck)
		crdtGroup := v1.Group("/dnt")
		{
			crdtGroup.GET("/table", getDNT)
			crdtGroup.GET("/peers", listPeers)
			crdtGroup.POST("/_node", updateLocal)
		}
		proxyGroup := v1.Group("/proxy")
		{
			proxyGroup.PATCH("/:peerId/*path", ForwardHandler)
			proxyGroup.POST("/:peerId/*path", ForwardHandler)
			proxyGroup.GET("/:peerId/*path", ForwardHandler)
		}
	}
	p2plistener := P2PListener()
	go func() {
		err := http.Serve(p2plistener, r)
		if err != nil {
			common.Logger.Error("http.Serve: %s", err)
		}
	}()
	wg.Wait()
	err := r.Run("0.0.0.0:" + viper.GetString("port"))
	if err != nil {
		panic(err)
	}
}
