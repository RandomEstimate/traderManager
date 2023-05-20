package main

import (
	"flag"
	"github.com/RandomEstimate/traderManager/util"
	"log"
	"net/http"
	"strconv"
	"time"
)
import "github.com/RandomEstimate/traderManager/mysqlManager"
import "github.com/RandomEstimate/traderManager/handler"

func main() {
	var path string
	flag.StringVar(&path, "path", "", "properties path")
	flag.Parse()

	properties, err := util.LoadProperties(path)
	if err != nil {
		panic(err)
	}

	port, _ := strconv.ParseFloat(properties["port"], 10)
	mysqlBaseConfig := mysqlManager.MySQLConfig{
		Host:     properties["host"],
		Port:     int(port),
		User:     properties["user"],
		Password: properties["password"],
		Database: properties["dbBase"],
	}
	mysqlPool, _ := mysqlManager.NewMySQLConnectionPool(&mysqlBaseConfig)

	strategyObj := handler.StrategyHandler{
		Db:                mysqlPool,
		RegisterTable:     properties["registerTable"],
		TestRegisterTable: properties["testRegisterTable"],
		OrderTable:        properties["orderTable"],
		TestOrderTable:    properties["testOrderTable"],
	}

	orderObj := handler.OrderHandler{
		Db:                mysqlPool,
		RegisterTable:     properties["registerTable"],
		TestRegisterTable: properties["testRegisterTable"],
		OrderTable:        properties["orderTable"],
		TestOrderTable:    properties["testOrderTable"],
	}

	mysqlStatConfig := mysqlManager.MySQLConfig{
		Host:     properties["host"],
		Port:     int(port),
		User:     properties["user"],
		Password: properties["password"],
		Database: properties["dbBase"],
	}
	mysqlStatPool, _ := mysqlManager.NewMySQLConnectionPool(&mysqlStatConfig)
	StatObj := util.NewStat(mysqlStatPool, properties["orderTableName"], properties["statTableName"], properties["testOrderTableName"], properties["testStatTableName"])

	go func() {
		for {
			time.Sleep(time.Minute * 5)
			err := StatObj.PositionStat(StatObj.TestStatTableName, StatObj.TestOrderTableName)
			if err != nil {
				log.Println(err)
				//return
			}
			err = StatObj.PositionStat(StatObj.StatTableName, StatObj.OrderTableName)
			if err != nil {
				log.Println(err)
				//return
			}

		}
	}()

	http.HandleFunc("/BatchOrder", orderObj.BatchOrder)
	http.HandleFunc("/StrategyDelete", strategyObj.StrategyDelete)
	http.HandleFunc("/StrategyRegister", strategyObj.StrategyRegister)
	err = http.ListenAndServe(":10000", nil)
	if err != nil {
		log.Fatal(err)
	}

	select {}
}
