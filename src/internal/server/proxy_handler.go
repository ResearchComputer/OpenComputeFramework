package server

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"ocf/internal/common"
	"ocf/internal/protocol"

	"github.com/gin-gonic/gin"
	p2phttp "github.com/libp2p/go-libp2p-http"
)

func ErrorHandler(res http.ResponseWriter, req *http.Request, err error) {
	res.Write([]byte(fmt.Sprintf("ERROR: %s", err.Error())))
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
	node, _ := protocol.GetP2PNode(nil)
	tr.RegisterProtocol("libp2p", p2phttp.NewTransport(node))
	target := url.URL{
		Scheme: "libp2p",
		Host:   requestPeer + ":9000",
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
