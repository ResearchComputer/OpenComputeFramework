package server

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"ocfcore/internal/cluster"
	"ocfcore/internal/common/requests"
	"ocfcore/internal/common/structs"
	"ocfcore/internal/profiler"
	"ocfcore/internal/server/p2p"
	"ocfcore/internal/server/queue"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	p2phttp "github.com/libp2p/go-libp2p-http"
	"github.com/nakabonne/tstorage"
	"github.com/spf13/viper"
)

func GetWorkerHub(c *gin.Context) {
	c.JSON(http.StatusOK, workerHub)
}

func GetConnections(c *gin.Context) {
	conn, err := queue.GetQueueStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, conn)
}

func GetWorkloadTable(c *gin.Context) {
	c.JSON(http.StatusOK, GlobalWorkloadTable())
}

func GetWorkerTable(c *gin.Context) {
	c.JSON(http.StatusOK, queue.NewLocalNodeTable())
}

func GetWorkerStatus(c *gin.Context) {
	workerID := c.Param("workerId")
	metricName := c.Param("metric")
	start := c.Query("start")
	if start == "" {
		start = "0"
	}
	end := c.Query("end")
	if end == "" {
		end = strconv.FormatInt(time.Now().Unix(), 10)
	}
	// convert to int64
	start_stamp, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	end_stamp, err := strconv.ParseInt(end, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var points []*tstorage.DataPoint
	var response []WorkerStatusResponse
	if metricName == "all" {
		metrics := [4]string{"Power Usage", "GPU Utilization", "Used Memory", "Available Memory"}
		var points []*tstorage.DataPoint
		for _, m := range metrics {
			metricPoints := append(points, profiler.QueryPoints(start_stamp, end_stamp, m, workerID)...)
			for _, point := range metricPoints {
				response = append(response, WorkerStatusResponse{workerID, m, point.Timestamp, point.Value})
			}
		}
	} else {
		points = profiler.QueryPoints(start_stamp, end_stamp, metricName, workerID)
		for _, point := range points {
			response = append(response, WorkerStatusResponse{workerID, metricName, point.Timestamp, point.Value})
		}
	}
	c.JSON(http.StatusOK, response)
}

func healthStatusCheck(c *gin.Context) {
	peer, exist := c.GetQuery("peer")
	if exist && peer == "1" {
		peers := []string{}
		for _, peer := range p2p.GetP2PNode().Peerstore().Peers() {
			err := requests.CheckPeerStatus(peer.String())
			if err == nil {
				peers = append(peers, peer.String())
			}
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok", "peers": peers})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func matchmakingStatus(c *gin.Context) {
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
	c.JSON(http.StatusOK, matchingStatus)
}

func GetSummary(c *gin.Context) {
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
	c.JSON(http.StatusOK, status)
}

// Forward Handler
func ForwardHandler(c *gin.Context) {
	requestPeer := c.Param("peerId")
	requestPath := c.Param("path")
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tr := &http.Transport{}
	node := p2p.GetP2PNode()
	tr.RegisterProtocol("libp2p", p2phttp.NewTransport(node))

	target := url.URL{
		Scheme: "libp2p",
		Host:   requestPeer,
		Path:   requestPath,
	}
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Path = target.Path
		req.URL.Host = req.Host
		req.Host = target.Host
		req.Method = c.Request.Method
		req.Body = io.NopCloser(bytes.NewBuffer(body))
	}
	proxy := httputil.NewSingleHostReverseProxy(&target)
	proxy.Director = director
	proxy.Transport = tr
	proxy.ErrorHandler = ErrorHandler
	proxy.ModifyResponse = rewriteHeader()
	proxy.ServeHTTP(c.Writer, c.Request)
}

// controller
func LoadWorkload(c *gin.Context) {
	workerID := c.Param("workerId")
	var provisionInstructions structs.ProvisionModelsPlan
	err := c.BindJSON(&provisionInstructions)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if instructionsHub.Instructions == nil {
		instructionsHub.Instructions = make(map[string]structs.ProvisionModelsPlan)
	}
	// check if the worker is already connected
	for _, worker := range workerHub.Workers {
		if worker.WorkerID == workerID {
			// update the worker
			instructionsHub.Instructions[workerID] = provisionInstructions
			c.JSON(200, gin.H{"message": "ok"})
			return
		}
	}
	c.JSON(404, gin.H{"error": "worker not found"})
}

func GetWorkloadInstructions(c *gin.Context) {
	workerID := c.Param("workerId")
	c.JSON(200, instructionsHub.Instructions[workerID])
}

func AddClusterNode(c *gin.Context) {
	var acquirePayload map[string]string
	c.BindJSON(&acquirePayload)
	acquireMachinePayload := structs.AcquireMachinePayload{
		Script: viper.GetString("acquire_machine.script"),
		Params: acquirePayload,
	}
	cluster.NewSlurmClusterClient().AcquireMachine(acquireMachinePayload)
}
