package gdbc

import (
	"sync"
	"reflect"
)

type DriverManager struct {
	driverList map[string]Driver
	mutex  sync.Mutex
}

var driverManager DriverManager

func RegisterDriver(driver Driver) bool{
	driverName := reflect.TypeOf(driver).Name()

	driverManager.mutex.Lock()
	defer driverManager.mutex.Unlock()

	if _, ok := driverManager.driverList[driverName]; ok {
		//WriteInfoDbLogf("Driver %s exist!", driverName)
		return false
	}

	driverManager.driverList[driverName] = driver

	return true
}

func DeregisterDriver(driver Driver) {
	driverName := reflect.TypeOf(driver).Name()

	driverManager.mutex.Lock()
	defer driverManager.mutex.Unlock()

	delete(driverManager.driverList, driverName)
}

func checkDriver(driverName string) bool {
	driverManager.mutex.Lock()
	defer driverManager.mutex.Unlock()

	_, ok := driverManager.driverList[driverName]

	return ok
}

func GetConnection(url string) (Connection, bool) {

	info := make(map[string]string)

	return GetConnectionByProperties(url, info)
}

func GetConnectionByProperties(url string, info map[string]string) (conn Connection, ok bool) {

	driverManager.mutex.Lock()
	defer driverManager.mutex.Unlock()
	for _, driver := range driverManager.driverList {
		if conn, ok = driver.Connect(url, info); ok {
			return conn, ok
		}
	}

	return conn, false
}

func GetConnectionByUser(url string, userName string, password string) (Connection, bool) {
	info := make(map[string]string)

	info["UserName"] = userName
	info["Password"] = password

	return GetConnectionByProperties(url, info)
}





