package server

import (
	"net/http"
	"rccore/internal/common"
	"rccore/internal/server/auth"
	"rccore/internal/server/p2p"
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
		rccoreStatus := v1.Group("/status")
		{
			rccoreStatus.GET("/health", healthStatusCheck)
			rccoreStatus.GET("/worker/:workerId/:metric", GetWorkerStatus)
			rccoreStatus.GET("/workers", GetWorkerHub)
			rccoreStatus.GET("/matchmaking", matchmakingStatus)
			rccoreStatus.GET("/summary", GetSummary)
		}
		rccoreWs := v1.Group("/ws")
		{
			rccoreWs.GET("",
				gin.WrapH(NewRPCServer().WebsocketHandler([]string{"*"})))
		}
		rccoreProxy := v1.Group("/proxy")
		{
			rccoreProxy.PATCH("/:peerId/*path", ForwardHandler)
			rccoreProxy.POST("/:peerId/*path", ForwardHandler)
			rccoreProxy.GET("/:peerId/*path", ForwardHandler)
		}
		rccoreThrottle := v1.Group("/throttle")
		{
			rccoreThrottle.Use(auth.AuthorizeMiddleware())
			rccoreThrottle.POST("/instructions/:workerId", LoadWorkload)
			rccoreThrottle.GET("/instructions/:workerId", GetWorkloadInstructions)
			rccoreThrottle.POST("/cluster/nodes", AddClusterNode)
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
