package gdbc

import (
	"reflect"
	"sync"
)

type DriverManager struct {
	driverList map[string]Driver
	mutex      sync.Mutex
}

var driverManager DriverManager

func RegisterDriver(driver Driver) (error DbError) {
	driverName := reflect.TypeOf(driver).Name()

	driverManager.mutex.Lock()
	defer driverManager.mutex.Unlock()

	if _, ok := driverManager.driverList[driverName]; ok {
		InfoLogf("Driver %s exist!", driverName)
		return NewDefaultDbError(-1, "Driver "+driverName+" exist!")
	}

	driverManager.driverList[driverName] = driver

	return nil
}

func DeregisterDriver(driver Driver) {
	driverName := reflect.TypeOf(driver).Name()

	driverManager.mutex.Lock()
	defer driverManager.mutex.Unlock()

	delete(driverManager.driverList, driverName)
}

func CheckDriver(driverName string) bool {
	driverManager.mutex.Lock()
	defer driverManager.mutex.Unlock()

	_, ok := driverManager.driverList[driverName]

	return ok
}

func GetConnection(url string) (connection Connection, error DbError) {

	info := make(map[string]string)

	return GetConnectionByProperties(url, info)
}

func GetConnectionByProperties(url string, info map[string]string) (conn Connection,
	error DbError) {

	driverManager.mutex.Lock()
	defer driverManager.mutex.Unlock()
	for _, driver := range driverManager.driverList {
		if conn, error = driver.Connect(url, info); error != nil {
			return conn, error
		}
	}

	return conn, nil
}

func GetConnectionByUser(url string, userName string, password string) (connection Connection,
	error DbError) {
	info := make(map[string]string)

	info["UserName"] = userName
	info["Password"] = password

	return GetConnectionByProperties(url, info)
}
