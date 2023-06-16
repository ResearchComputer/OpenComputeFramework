package server

import (
	"net/http"
	"ocfcore/internal/common"
	"ocfcore/internal/server/auth"
	"ocfcore/internal/server/p2p"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func StartServer() {
	var wg sync.WaitGroup
	r := gin.Default()
	r.Use(beforeResponse())
	r.Use(gin.Recovery())
	v1 := r.Group("/api/v1")
	{
		ocfcoreStatus := v1.Group("/status")
		{
			ocfcoreStatus.GET("/health", healthStatusCheck)
			ocfcoreStatus.GET("/worker/:workerId/:metric", GetWorkerStatus)
			ocfcoreStatus.GET("/workers", GetWorkerHub)
			ocfcoreStatus.GET("/matchmaking", matchmakingStatus)
			ocfcoreStatus.GET("/summary", GetSummary)
		}
		ocfcoreWs := v1.Group("/ws")
		{
			ocfcoreWs.GET("",
				gin.WrapH(NewRPCServer().WebsocketHandler([]string{"*"})))
		}
		ocfcoreProxy := v1.Group("/proxy")
		{
			ocfcoreProxy.PATCH("/:peerId/*path", ForwardHandler)
			ocfcoreProxy.POST("/:peerId/*path", ForwardHandler)
			ocfcoreProxy.GET("/:peerId/*path", ForwardHandler)
		}
		ocfcoreThrottle := v1.Group("/throttle")
		{
			ocfcoreThrottle.Use(auth.AuthorizeMiddleware())
			ocfcoreThrottle.POST("/instructions/:workerId", LoadWorkload)
			ocfcoreThrottle.GET("/instructions/:workerId", GetWorkloadInstructions)
			ocfcoreThrottle.POST("/cluster/nodes", AddClusterNode)
		}
	}
	p2plistener := p2p.P2PListener()
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
