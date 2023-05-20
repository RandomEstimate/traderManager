package handler

import (
	"fmt"
	"github.com/RandomEstimate/traderManager/mysqlManager"
	"net/http"
)

// StrategyHandler /*
/*
StrategyHandler 主要实现策略的创建和删除任务

RegisterTable           生产 策略名称存储表
OrderTable   			生产 订单名称存储表
TestRegisterTable       测试 策略名称存储表
TestOrderTable          测试 订单名称存储表

StrategyHandler 对外暴露的接口
-- StrategyRegister 策略注册
-- StrategyDelete 策略删除

*/
type StrategyHandler struct {
	Db                *mysqlManager.MySQLConnectionPool
	RegisterTable     string
	OrderTable        string
	TestRegisterTable string
	TestOrderTable    string
}

func (a *StrategyHandler) StrategyRegister(w http.ResponseWriter, r *http.Request) {
	// 解析传递参数
	values := r.URL.Query()
	param := values.Get("StrategyName")

	if len(param) <= 4 {
		response(w, StrategyResponse{
			Code:    STRATEGY_PARAM_ERROR,
			Message: fmt.Sprintf("策略名称小于等于4个字符无法被解析\n"),
		})
		return
	}

	if param[:4] == "TEST" {
		// 注册为测试策略 写入测试策略数据库
		strategyResponseCode, err := a.strategyRegister(a.TestRegisterTable, param)
		if err != nil {
			response(w, StrategyResponse{
				Code:    strategyResponseCode,
				Message: fmt.Sprintf("注册策略错误：%v\n", err),
			})
			return
		}
	} else {
		// 注册为正式策略 写入正式策略数据库
		strategyResponseCode, err := a.strategyRegister(a.RegisterTable, param)
		if err != nil {
			response(w, StrategyResponse{
				Code:    strategyResponseCode,
				Message: fmt.Sprintf("注册策略错误：%v\n", err),
			})
			return
		}
	}

	response(w, StrategyResponse{
		Code:    STRATEGY_SUCCESS,
		Message: fmt.Sprintf("注册成功\n"),
	})
	return

}

func (a *StrategyHandler) strategyRegister(tableName string, strategyName string) (StrategyResponseCode, error) {
	db, err := a.Db.GetDB()
	if err != nil {
		return STRATEGY_DB_ERROR, err
	}

	rows, err := db.Query(fmt.Sprintf("SELECT strategy FROM %s WHERE strategy = ?", tableName), strategyName)
	if err != nil {
		return STRATEGY_DB_ERROR, err
	}

	count := 0
	for rows.Next() {
		count += 1
	}

	if count != 0 {
		return STRATEGY_EXIST, fmt.Errorf("策略已经被注册，无需再次注册")
	}

	_, err = db.Exec(fmt.Sprintf("INSERT INTO %s (strategy) values (?)", tableName), strategyName)
	if err != nil {
		return STRATEGY_DB_ERROR, fmt.Errorf("策略注册写入数据库失败：%v", err)
	}

	return STRATEGY_SUCCESS, nil

}

func (a *StrategyHandler) StrategyDelete(w http.ResponseWriter, r *http.Request) {
	// 解析传递参数
	values := r.URL.Query()
	param := values.Get("StrategyName")

	if len(param) <= 4 {
		response(w, StrategyResponse{
			Code:    STRATEGY_PARAM_ERROR,
			Message: fmt.Sprintf("策略名称小于等于4个字符无法被解析\n"),
		})
		return
	}

	if param[:4] == "TEST" {
		// 注册为测试策略 写入测试策略数据库
		strategyResponseCode, err := a.strategyDelete(a.TestRegisterTable, a.TestOrderTable, param)
		if err != nil {
			response(w, StrategyResponse{
				Code:    strategyResponseCode,
				Message: fmt.Sprintf("策略删除错误：%v\n", err),
			})
			return
		}
	} else {
		// 注册为正式策略 写入正式策略数据库
		strategyResponseCode, err := a.strategyDelete(a.RegisterTable, a.OrderTable, param)
		if err != nil {
			response(w, StrategyResponse{
				Code:    strategyResponseCode,
				Message: fmt.Sprintf("策略删除错误：%v\n", err),
			})
			return
		}
	}

	response(w, StrategyResponse{
		Code:    STRATEGY_SUCCESS,
		Message: fmt.Sprintf("删除成功\n"),
	})

}

func (a *StrategyHandler) strategyDelete(tableName string, orderTableName string, strategyName string) (StrategyResponseCode, error) {
	db, err := a.Db.GetDB()
	if err != nil {
		return STRATEGY_DB_ERROR, err
	}

	tx, err := db.Begin()
	if err != nil {
		return STRATEGY_DB_ERROR, err
	}

	_, err = tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE strategy = ?", tableName), strategyName)
	if err != nil {
		tx.Rollback()
		return STRATEGY_DB_ERROR, err
	}
	_, err = tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE strategy = ?", orderTableName), strategyName)
	if err != nil {
		tx.Rollback()
		return STRATEGY_DB_ERROR, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return STRATEGY_DB_ERROR, err
	}
	return STRATEGY_SUCCESS, err
}
