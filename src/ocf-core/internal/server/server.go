package server

import (
	"net/http"
	"ocfcore/internal/common"
	"ocfcore/internal/protocol/p2p"
	"ocfcore/internal/server/auth"
	"ocfcore/internal/server/queue"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func StartServer() {
	PrintWelcomeMessage()
	var wg sync.WaitGroup
	gin.SetMode(gin.ReleaseMode)
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
			ocfcoreStatus.GET("/connections", GetConnections)
			ocfcoreStatus.GET("/table", GetWorkloadTable)
			// ocfcoreStatus.POST("/table", UpdateWorkloadTable)
			ocfcoreStatus.GET("/peers", GetPeersInfo)
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
		ocfcoreThrottle := v1.Group("/controller")
		{
			ocfcoreThrottle.Use(auth.AuthorizeMiddleware())
			ocfcoreThrottle.POST("/instructions/:workerId", LoadWorkload)
			ocfcoreThrottle.GET("/instructions/:workerId", GetWorkloadInstructions)
			ocfcoreThrottle.POST("/cluster/nodes", AddClusterNode)
		}
		ocfcoreRequest := v1.Group("/request")
		{
			ocfcoreRequest.POST("/inference", AutoInferenceRequest)
			ocfcoreRequest.POST("/_inference", InferenceRequest)
		}
	}
	p2plistener := p2p.P2PListener()
	go func() {
		err := http.Serve(p2plistener, r)
		if err != nil {
			common.Logger.Error("http.Serve: %s", err)
		}
	}()
	queue.StartQueueServer()
	wg.Wait()
	err := r.Run("0.0.0.0:" + viper.GetString("port"))
	if err != nil {
		panic(err)
	}
}
