package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"traderManager/handler"
)

type Client struct {
	Host string //http://127.0.0.1
	c    *http.Client
}

func NewClient(host string) *Client {
	return &Client{
		Host: host,
		c:    &http.Client{},
	}
}

func (a *Client) StrategyRegister(param *handler.StrategyRegisterRequest) (*handler.StrategyRegisterResponse, error) {
	api := "/StrategyHandler"
	req, _ := http.NewRequest("GET", a.Host+api+fmt.Sprintf("?StrategyName=%s", param.Name), nil)

	resp, err := a.c.Do(req)
	if err != nil {
		return nil, err
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	d := &handler.StrategyRegisterResponse{}
	err = json.Unmarshal(buf, d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (a *Client) StrategyDelete(param *handler.StrategyRegisterRequest) (*handler.StrategyRegisterResponse, error) {
	api := "/StrategyDelete"
	req, _ := http.NewRequest("GET", a.Host+api+fmt.Sprintf("?StrategyName=%s", param.Name), nil)

	resp, err := a.c.Do(req)
	if err != nil {
		return nil, err
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	d := &handler.StrategyRegisterResponse{}
	err = json.Unmarshal(buf, d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (a *Client) OrderCommit(param *handler.OrderHandlerRequest) (*handler.OrderHandlerResponse, error) {
	api := "/BatchOrder"
	buf, err := json.Marshal(param)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("Cannot parse param. %v", err))
	}

	req, _ := http.NewRequest("POST", a.Host+api, strings.NewReader(string(buf)))
	resp, err := a.c.Do(req)
	if err != nil {
		return nil, err
	}

	buf, _ = ioutil.ReadAll(resp.Body)
	d := &handler.OrderHandlerResponse{}
	err = json.Unmarshal(buf, d)
	if err != nil {
		return nil, err
	}
	return d, nil
}
