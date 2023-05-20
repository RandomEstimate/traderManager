package util

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2/futures"
	"strconv"
	"traderManager/mysqlManager"
)

type Stat struct {
	db                 *mysqlManager.MySQLConnectionPool
	OrderTableName     string
	StatTableName      string
	TestOrderTableName string
	TestStatTableName  string
}

func NewStat(db *mysqlManager.MySQLConnectionPool, orderTableName string, statTableName string, testOrderTableName string, testStatTableName string) *Stat {
	return &Stat{db: db, OrderTableName: orderTableName, StatTableName: statTableName, TestOrderTableName: testOrderTableName, TestStatTableName: testStatTableName}
}

func (a *Stat) PositionStat(insertTable string, orderTable string) error {

	db, err := a.db.GetDB()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	// 创建一个临时表
	_, err = tx.Exec("CREATE TEMPORARY TABLE temp_price_table (symbol VARCHAR(100), price float)")
	if err != nil {
		tx.Rollback()
		return err
	}

	// 将价格进行存储进入
	resp, err := futures.NewClient("", "").NewListBookTickersService().Do(context.Background())
	if err != nil {
		tx.Rollback()
		return err
	}

	stmx, err := tx.Prepare("INSERT INTO temp_price_table(symbol,price) VALUES(?,?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, info := range resp {
		askPrice, _ := strconv.ParseFloat(info.AskPrice, 10)
		bidPrice, _ := strconv.ParseFloat(info.BidPrice, 10)
		_, err := stmx.Exec(info.Symbol, (askPrice+bidPrice)/2)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	stmx.Close()

	// 进行价格拼接
	_, err = tx.Exec(fmt.Sprintf("TRUNCATE TABLE %s", insertTable))
	if err != nil {
		tx.Rollback()
		return err
	}
	joinSql := fmt.Sprintf(`
	INSERT INTO %s (strategy,symbol,profit) 
	SELECT strategy,symbol,SUM(side * (price - open_price)) as profit FROM (
		SELECT t1.strategy as strategy,t1.symbol as symbol,t1.side as side,t1.qty as qty,t1.price as open_price,t2.price as price 
		FROM %s as t1 
		LEFT JOIN  %s as t2 ON t1.symbol = t2.symbol
		WHERE t2.price IS NOT NULL 
	) as tmp
	GROUP BY strategy,symbol
	`, insertTable, orderTable, "temp_price_table")
	_, err = tx.Exec(joinSql)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(fmt.Sprintf("DROP TABLE %s", "temp_price_table"))
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil

}
