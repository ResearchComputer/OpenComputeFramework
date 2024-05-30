package server

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"ocf/internal/common"
	"ocf/internal/protocol"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	p2phttp "github.com/libp2p/go-libp2p-http"
)

func ErrorHandler(res http.ResponseWriter, req *http.Request, err error) {
	res.Write([]byte(fmt.Sprintf("ERROR: %s", err.Error())))
}

// P2P handler for forwarding requests to other peers
func P2PForwardHandler(c *gin.Context) {
	requestPeer := c.Param("peerId")
	requestPath := c.Param("path")
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tr := &http.Transport{}
	node, _ := protocol.GetP2PNode(nil)
	tr.RegisterProtocol("libp2p", p2phttp.NewTransport(node))
	target := url.URL{
		Scheme: "libp2p",
		Host:   requestPeer,
		Path:   requestPath,
	}
	common.Logger.Info("Forwarding request to %s", target.String())
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

// ServiceHandler
func ServiceForwardHandler(c *gin.Context) {
	serviceName := c.Param("service")
	requestPath := c.Param("path")
	service, err := protocol.GetService(serviceName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	target := url.URL{
		Scheme: "http",
		Host:   service.Host + ":" + service.Port,
		Path:   requestPath,
	}
	director := func(req *http.Request) {
		req.Host = target.Host
		req.URL.Host = req.Host
		req.URL.Scheme = target.Scheme
		req.URL.Path = target.Path
		req.Body = io.NopCloser(bytes.NewBuffer(body))
	}
	proxy := httputil.NewSingleHostReverseProxy(&target)
	proxy.Director = director
	proxy.ServeHTTP(c.Writer, c.Request)
}

// in case of global service, we need to forward the request to the service, identified by the service name and identity group
func GlobalServiceForwardHandler(c *gin.Context) {
	serviceName := c.Param("service")
	requestPath := c.Param("path")
	providers, err := protocol.GetAllProviders(serviceName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	// find proper service that are within the same identity group
	// first filter by service name, then iterative over the identity groups
	// always find all the services that are in the same identity group
	var candidates []string
	for _, provider := range providers {
		selected := false
		for _, service := range provider.Service {
			if service.Name == serviceName {
				selected = true
				// check if the service is in the same identity group
				if len(service.IdentityGroup) > 0 {
					for _, ig := range service.IdentityGroup {
						igGroup := strings.Split(ig, "=")
						igKey := igGroup[0]
						igValue := igGroup[1]
						requestGroup, err := jsonparser.GetString(body, igKey)
						if err != nil {
							selected = false
						}
						if requestGroup != igValue {
							selected = false
						}
					}
				}
				// append the service to the candidates
				if selected {
					candidates = append(candidates, provider.ID)
				}
			}
		}
	}
	// randomly select one of the candidates
	// here's where we can implement a load balancing algorithm
	randomIndex := rand.Intn(len(candidates))
	tr := &http.Transport{}
	node, _ := protocol.GetP2PNode(nil)
	tr.RegisterProtocol("libp2p", p2phttp.NewTransport(node))
	targetPeer := candidates[randomIndex]
	// replace the request path with the _service path
	requestPath = "/v1/_service/" + serviceName + requestPath
	common.Logger.Info("Forwarding request to: ", targetPeer)
	common.Logger.Info("Forwarding path to: ", requestPath)
	target := url.URL{
		Scheme: "libp2p",
		Host:   targetPeer,
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
