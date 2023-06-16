package server

import (
	"net/http"
	"ocfcore/internal/common/requests"
	"ocfcore/internal/common/structs"
	"ocfcore/internal/profiler"
	"ocfcore/internal/server/p2p"
	"time"

	"github.com/gin-gonic/gin"
)

func healthStatusCheck(g *gin.Context) {
	peer, exist := g.GetQuery("peer")
	if exist && peer == "1" {
		peers := []string{}
		for _, peer := range p2p.GetP2PNode().Peerstore().Peers() {
			err := requests.CheckPeerStatus(peer.String())
			if err == nil {
				peers = append(peers, peer.String())
			}
		}
		g.JSON(http.StatusOK, gin.H{"status": "ok", "peers": peers})
	} else {
		g.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func matchmakingStatus(g *gin.Context) {
	workerStatus := make(map[string]structs.MatchingWorkerStatus)
	for _, worker := range workerHub.Workers {
		averageUtilization := profiler.AggregateAverageUtilization(worker.WorkerID, 30*time.Second)
		status := "IDLE"
		if averageUtilization > 0.3 {
			status = "BUSY"
		}
		workerStatus[worker.WorkerID] = structs.MatchingWorkerStatus{
			Accelerator: worker.GPUSpecifier,
			Status:      status,
		}
	}
	modelStatus := make(map[string]structs.MatchingModelStatus)
	matchingStatus := structs.MatchingStatus{
		Workers:   workerStatus,
		Models:    modelStatus,
		Timestamp: time.Now().Unix(),
	}
	g.JSON(http.StatusOK, matchingStatus)
}

func GetSummary(g *gin.Context) {
	status := make(map[string]structs.CardStatus)
	for _, worker := range workerHub.Workers {
		averageUtilization := profiler.AggregateAverageUtilization(worker.WorkerID, 30*time.Second)
		cardStatus := "IDLE"
		if averageUtilization > 0.3 {
			cardStatus = "BUSY"
		}
		cardMetrics, err := profiler.QueryCardSummary(worker.WorkerID)
		if err == nil {
			status[worker.WorkerID] = structs.CardStatus{
				CardID:          worker.WorkerID,
				Status:          cardStatus,
				Serving:         worker.Serving,
				PowerUsage:      cardMetrics.PowerUsage,
				GPUUtilization:  cardMetrics.GPUUtilization,
				UsedMemory:      cardMetrics.UsedMemory,
				AvailableMemory: cardMetrics.AvailableMemory,
				GPUSpecifier:    worker.GPUSpecifier,
			}
		}
	}
	g.JSON(http.StatusOK, status)
}
