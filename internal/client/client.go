package client

import (
	"encoding/json"
	"garantex/internal/models"
	"io"
	"net/http"
	"time"
)

const (
	url = "https://garantex.org/api/v2/depth?market=usdtrub"
)

type ApiClient struct{}

func (*ApiClient) DoRequest() (*models.PriceDepth, error) {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	t := &models.PriceDepth{}
	err = json.Unmarshal(b, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}
