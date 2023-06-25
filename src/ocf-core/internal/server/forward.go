package server

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"ocfcore/internal/server/p2p"

	"github.com/gin-gonic/gin"
	p2phttp "github.com/libp2p/go-libp2p-http"
)

func ErrorHandler(res http.ResponseWriter, req *http.Request, err error) {
	res.Write([]byte(fmt.Sprintf("ERROR: %s", err.Error())))
}

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
	// print peers
	for _, p := range node.Peerstore().Peers() {
		fmt.Println("peer: ", p)
	}
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
