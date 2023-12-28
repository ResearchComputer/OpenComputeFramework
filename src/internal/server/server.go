package server

import (
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
	}
	wg.Wait()
	err := r.Run("0.0.0.0:" + viper.GetString("port"))
	if err != nil {
		panic(err)
	}
}
