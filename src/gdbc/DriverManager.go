package gdbc

import "sync"

type DriverManager struct {
	driverList map[string]Driver
	mutex  sync.Mutex
}

var driverManager DriverManager

func RegisterDriver(driver Driver) bool{
	return true
}

func DeregisterDriver(driver Driver) bool {

}

func GetConnection(url string) (Connection, bool) {

}

func GetConnectionByProperties(url string, info map[string]string) (Connection, bool) {

}

func GetConnectionByUser(url string, userName string, password string) (Connection, bool) {

}

func GetDriver(url string) Driver {

}




