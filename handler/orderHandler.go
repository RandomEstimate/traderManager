package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/adshao/go-binance/v2/futures"
	"net/http"
	"strconv"
	"strings"
	"time"
	"traderManager/mysqlManager"
)

// OrderHandler /*
/*
OrderHandler 主要实现订单的创建和维护任务

RegisterTable           生产 策略名称存储表
OrderTable   			生产 订单名称存储表
TestRegisterTable       测试 策略名称存储表
TestOrderTable          测试 订单名称存储表

OrderHandler 对外暴露的接口
-- BatchOrder 批量下单

*/
type OrderHandler struct {
	Db                *mysqlManager.MySQLConnectionPool
	RegisterTable     string
	OrderTable        string
	TestRegisterTable string
	TestOrderTable    string
}

func (a *OrderHandler) BatchOrder(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		response(w, OrderHandlerResponse{
			Code:    ORDER_PARAM_ERROR,
			Message: fmt.Sprintf("无法解析参数：%v\n", err),
		})
		return
	}

	orderReq := &OrderHandlerRequest{}
	if err := json.NewDecoder(r.Body).Decode(orderReq); err != nil {
		response(w, OrderHandlerResponse{
			Code:    ORDER_PARAM_ERROR,
			Message: fmt.Sprintf("无法解析参数：%v\n", err),
		})
		return
	}

	if len(orderReq.StrategyName) <= 4 {
		response(w, OrderHandlerResponse{
			Code:    ORDER_PARAM_ERROR,
			Message: fmt.Sprintf("策略名称小于等于4个字符无法被解析\n"),
		})
		return
	}

	var (
		resp *futures.CreateBatchOrdersResponse
	)
	if orderReq.StrategyName[:4] == "TEST" {
		if !a.checkStrategyExist(a.TestRegisterTable, orderReq.StrategyName) {
			response(w, OrderHandlerResponse{
				Code:    ORDER_STRATEGY_NOT_EXIST,
				Message: fmt.Sprintf("策略没有被注册，请先注册策略\n"),
			})
			return
		}

		err := a.batchOrderToDb(orderReq, a.TestOrderTable)
		if err != nil {
			response(w, OrderHandlerResponse{
				Code:    ORDER_DB_ERROR,
				Message: fmt.Sprintf("写入数据库存在问题：%v\n", err),
			})
			return
		}

	} else {
		if !a.checkStrategyExist(a.RegisterTable, orderReq.StrategyName) {
			response(w, OrderHandlerResponse{
				Code:    ORDER_STRATEGY_NOT_EXIST,
				Message: fmt.Sprintf("策略没有被注册，请先注册策略\n"),
			})
			return
		}

		// 进行下单处理
		orderList := make([]*futures.CreateOrderService, 0)
		for _, v := range orderReq.OrderList {
			orderList = append(orderList, v.toCreateOrderService())
		}
		resp, err = futures.NewClient(orderReq.AccessToken, orderReq.SecretToken).NewCreateBatchOrdersService().OrderList(orderList).Do(context.Background())
		if err != nil {
			response(w, OrderHandlerResponse{
				Code:    ORDER_ERROR,
				Message: fmt.Sprintf("订单出现访问错误：%v \n", err),
			})
			return
		}

		// 写入数据库
		err := a.batchOrderToDb(orderReq, a.OrderTable)
		if err != nil {
			response(w, OrderHandlerResponse{
				Code:    ORDER_DB_ERROR,
				Message: fmt.Sprintf("写入数据库存在问题：%v\n", err),
			})
			return
		}
	}
	response(w, OrderHandlerResponse{
		Code: ORDER_SUCCESS,
		Message: fmt.Sprintf("下单成功：%v \n", func() string {
			if resp != nil {
				tmp := make([]string, 0)
				for _, order := range resp.Orders {
					buf, err := json.Marshal(*order)
					if err != nil {
						continue
					}
					tmp = append(tmp, string(buf))
				}
				return strings.Join(tmp, ",")
			}
			return "(无)"
		}()),
	})
	return
}

func (a *OrderHandler) batchOrderToDb(orderReq *OrderHandlerRequest, orderTableName string) error {
	// 检测策略是否被注册

	// 解析订单 并送入数据库
	db, err := a.Db.GetDB()
	if err != nil {
		return err
	}

	tickData, err := futures.NewClient("", "").NewListBookTickersService().Do(context.Background())
	if err != nil {
		return err
	}

	m := make(map[string]*futures.BookTicker, 0)
	for _, v := range tickData {
		m[v.Symbol] = v
	}

	for _, info := range orderReq.OrderList {
		var price float64
		if _, ok := m[info.Symbol]; !ok {
			return fmt.Errorf(fmt.Sprintf("%v 无法访问到标的物的价格", info.Symbol))
		}
		askPrice, _ := strconv.ParseFloat(m[info.Symbol].AskPrice, 10)
		bidPrice, _ := strconv.ParseFloat(m[info.Symbol].BidPrice, 10)
		price = (askPrice + bidPrice) / 2

		_, err := db.Exec(fmt.Sprintf("INSERT INTO %s (strategy,symbol,price,side,qty,time) VALUES(?,?,?,?,?,?)", orderTableName), orderReq.StrategyName, info.Symbol, price, func() int {
			if info.Side == futures.SideTypeBuy {
				return 1
			}
			return -1
		}(), info.Quantity, time.Now().UnixMilli())
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *OrderHandler) checkStrategyExist(tableName string, strategyName string) bool {
	db, err := a.Db.GetDB()
	if err != nil {
		return false
	}

	rows, err := db.Query(fmt.Sprintf("SELECT strategy FROM %s WHERE strategy = ?", tableName), strategyName)
	if err != nil {
		return false
	}
	count := 0
	for rows.Next() {
		count += 1
	}
	if count != 1 {
		return false
	}
	return true
}
