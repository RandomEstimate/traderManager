package handler

// StrategyRequest /*
/*
StrategyRequest 定义 cli访问传参结构体
*/
type StrategyRequest struct {
	Name string
}

// StrategyResponse /*
/*
StrategyResponse 定义cli访问返回结果
*/
type StrategyResponse struct {
	Code    StrategyResponseCode `json:"code"`
	Message string               `json:"message"`
}

type StrategyResponseCode int

const (
	STRATEGY_SUCCESS StrategyResponseCode = iota
	STRATEGY_PARAM_ERROR
	STRATEGY_EXIST
	STRATEGY_DB_ERROR
)
