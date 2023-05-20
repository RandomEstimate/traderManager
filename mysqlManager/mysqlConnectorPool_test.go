package mysqlManager

import "testing"

func TestMySQLConnectionPool(t *testing.T) {
	mysqlConfig := MySQLConfig{
		Host:     "124.71.84.193",
		Port:     3306,
		User:     "root",
		Password: "751037790qQ!",
		Database: "test",
	}

	mySQLConnectionPool, err := NewMySQLConnectionPool(&mysqlConfig)
	if err != nil {
		t.Error(err)
	}

	_, err = mySQLConnectionPool.GetDB()
	if err != nil {
		t.Error(err)
	}

}
