package handler

import "github.com/adshao/go-binance/v2/futures"

// OrderHandlerRequest /*
/*
OrderHandlerRequest 定义 cli访问传参结构体
*/
type OrderHandlerRequest struct {
	AccessToken  string                `json:"accessToken"`
	SecretToken  string                `json:"secretToken"`
	StrategyName string                `json:"strategyName"`
	OrderList    []*CreateOrderService `json:"orderList"`
}

// CreateOrderService /*
/*
CreateOrderService 定义 OrderList传入结构体，目的在于方便进行manger管理
因为binance sdk 的参数无法被直接访问
*/
type CreateOrderService struct {
	Symbol           string                    `json:"symbol"`
	Side             futures.SideType          `json:"side"`
	PositionSide     *futures.PositionSideType `json:"positionSide"`
	OrderType        futures.OrderType         `json:"orderType"`
	TimeInForce      *futures.TimeInForceType  `json:"timeInForce"`
	Quantity         string                    `json:"quantity"`
	ReduceOnly       *bool                     `json:"reduceOnly"`
	Price            *string                   `json:"price"`
	NewClientOrderID *string                   `json:"newClientOrderID"`
	StopPrice        *string                   `json:"stopPrice"`
	WorkingType      *futures.WorkingType      `json:"workingType"`
	ActivationPrice  *string                   `json:"activationPrice"`
	CallbackRate     *string                   `json:"callbackRate"`
	PriceProtect     *bool                     `json:"priceProtect"`
	NewOrderRespType futures.NewOrderRespType  `json:"newOrderRespType"`
	ClosePosition    *bool                     `json:"closePosition"`
}

// toCreateOrderService /*
/*
将外部包装的CreateOrderService转化为sdk支持的futures.CreateOrderService
*/
func (c *CreateOrderService) toCreateOrderService() *futures.CreateOrderService {
	o := &futures.CreateOrderService{}
	o.Symbol(c.Symbol)
	o.Side(c.Side)
	if c.PositionSide != nil {
		o.PositionSide(*c.PositionSide)
	}
	o.Type(c.OrderType)
	if c.TimeInForce != nil {
		o.TimeInForce(*c.TimeInForce)
	}
	o.Quantity(c.Quantity)
	if c.ReduceOnly != nil {
		o.ReduceOnly(*c.ReduceOnly)
	}
	if c.Price != nil {
		o.Price(*c.Price)
	}
	if c.NewClientOrderID != nil {
		o.NewClientOrderID(*c.NewClientOrderID)
	}
	if c.StopPrice != nil {
		o.StopPrice(*c.StopPrice)
	}
	if c.WorkingType != nil {
		o.WorkingType(*c.WorkingType)
	}
	if c.ActivationPrice != nil {
		o.ActivationPrice(*c.ActivationPrice)
	}
	if c.CallbackRate != nil {
		o.CallbackRate(*c.CallbackRate)
	}
	if c.PriceProtect != nil {
		o.PriceProtect(*c.PriceProtect)
	}
	o.NewOrderResponseType(c.NewOrderRespType)
	if c.ClosePosition != nil {
		o.ClosePosition(*c.ClosePosition)
	}
	return o
}

// OrderHandlerResponse /*
/*
OrderHandlerResponse 定义cli访问返回结果
*/
type OrderHandlerResponse struct {
	Code    OrderResponseCode `json:"code"`
	Message string            `json:"message"`
}

type OrderResponseCode int

const (
	ORDER_SUCCESS OrderResponseCode = iota
	ORDER_PARAM_ERROR
	ORDER_STRATEGY_NOT_EXIST
	ORDER_ERROR
	ORDER_DB_ERROR
)
