package requests

import (
	"encoding/json"
	"fmt"
	"io"
	"ocfcore/internal/common/structs"
	"strings"

	"github.com/spf13/viper"
)

func ForwardInferenceRequest(peerId string, req structs.InferenceStruct) (string, error) {
	remoteAddr := fmt.Sprintf("http://localhost:%s/api/v1/proxy/%s/api/v1/request/_inference", viper.GetString("port"), peerId)
	reqString, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	payload := strings.NewReader(string(reqString))
	resp, err := NewHTTPClient().Post(remoteAddr, "application/json", payload)
	if err != nil {
		return "", err
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
