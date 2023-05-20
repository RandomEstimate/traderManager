package cli

import (
	"fmt"
	"github.com/adshao/go-binance/v2/futures"
	"testing"
	"traderManager/handler"
)

//const host = "http://47.57.95.94:10000"
const host = "http://127.0.0.1:10000"

func TestCliStrategyRegister(t *testing.T) {
	c := NewClient(host)

	req := handler.StrategyRegisterRequest{
		Name: "test02",
	}

	register, err := c.StrategyRegister(&req)
	if err != nil {
		return
	}
	fmt.Println(register)

}

func TestClient_StrategyDelete(t *testing.T) {
	c := NewClient(host)

	req := handler.StrategyRegisterRequest{
		Name: "TEST-ONT-NEO-STORJ-v1",
	}

	register, err := c.StrategyDelete(&req)
	if err != nil {
		return
	}
	fmt.Println(register)
}

func TestClient_BatchOrder(t *testing.T) {
	c := NewClient(host)

	l := make([]*handler.CreateOrderService, 0)
	l = append(l, &handler.CreateOrderService{
		Symbol:    "ONTUSDT",
		Side:      futures.SideTypeBuy,
		OrderType: "",
		Quantity:  fmt.Sprint(1000),
	})
	l = append(l, &handler.CreateOrderService{
		Symbol:    "NEOUSDT",
		Side:      futures.SideTypeSell,
		OrderType: "",
		Quantity:  fmt.Sprint(1000),
	})
	l = append(l, &handler.CreateOrderService{
		Symbol:    "STORJUSDT",
		Side:      futures.SideTypeSell,
		OrderType: "",
		Quantity:  fmt.Sprint(1000),
	})
	o := handler.OrderHandlerRequest{
		AccessToken:  "",
		SecretToken:  "",
		StrategyName: "TEST-ONT-NEO-STORJ-v1",
		OrderList:    l,
	}

	OrderCommit, err := c.OrderCommit(&o)
	if err != nil {
		return
	}
	fmt.Println(OrderCommit)
}
