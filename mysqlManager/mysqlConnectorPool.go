package mysqlManager

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"sync"
)

type MySQLConnectionPool struct {
	db     *sql.DB
	config *MySQLConfig
	mu     sync.Mutex
}

type MySQLConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func NewMySQLConnectionPool(config *MySQLConfig) (*MySQLConnectionPool, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.User, config.Password, config.Host, config.Port, config.Database)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	pool := &MySQLConnectionPool{
		db:     db,
		config: config,
	}

	return pool, nil
}

func (pool *MySQLConnectionPool) GetDB() (*sql.DB, error) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	if pool.db == nil {
		return nil, fmt.Errorf("connection pool is closed")
	}

	err := pool.db.Ping()
	if err != nil {
		// 如果连接不可用，则重新连接
		dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", pool.config.User, pool.config.Password, pool.config.Host, pool.config.Port, pool.config.Database)
		db, err := sql.Open("mysql", dataSourceName)
		if err != nil {
			return nil, err
		}
		pool.db = db
	}

	return pool.db, nil
}

func (pool *MySQLConnectionPool) Close() error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	if pool.db == nil {
		return nil
	}

	err := pool.db.Close()
	if err != nil {
		return err
	}

	pool.db = nil
	return nil
}
func (pool *MySQLConnectionPool) ExecWithParams(query string, args ...[]interface{}) (sql.Result, error) {

	db, err := pool.GetDB()
	if err != nil {
		return nil, err
	}
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Failed to begin transaction: ", err)
		return nil, err
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		fmt.Println("Failed to prepare statement: ", err)
		return nil, err
	}

	var result sql.Result
	if args == nil {
		result, err = stmt.Exec()
		if err != nil {
			return nil, err
		}
	} else {
		for _, arg := range args {
			_, err = stmt.Exec(arg...)
		}
	}

	if err != nil {
		fmt.Println("Failed to stmx transaction: ", err)
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("Failed to commit transaction: ", err)
		tx.Rollback()
		return nil, err
	}

	return result, nil
}

func (pool *MySQLConnectionPool) QueryWithParams(query string, args ...interface{}) (*sql.Rows, error) {
	db, err := pool.GetDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
