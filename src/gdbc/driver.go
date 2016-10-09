package gdbc

import (
	"time"
)

type DbError interface {

	GetErrorCode() int

	GetErrorMessage() string
}

type Drive interface {
	Connect(driver string, properties map[string]string) (connection Connection, error DbError)
}

type Connection interface{
	CreateStatement() (statement Statement, error DbError)

	PrepareStatement(sql string) (prepareStatement PrepareStatement, error DbError)

	SetAutoCommit(autoCommit bool) DbError

	GetAutoCommit()(autoCommit bool, error DbError)

	Commit() DbError

	Rollback() DbError

	Close() DbError

	IsClose() (isClose bool, error DbError)

}


type Statement interface{
	ExecuteQuery() (result ResultSet, error DbError)

	ExecuteUpdate() (result int, error DbError)

	Execute()(result bool, error DbError)

	Close() DbError

	GetMaxFieldSize()(max int, error DbError)

	SetMaxFieldSize(max int) DbError

	GetMaxRows()(max int, error DbError)

	SetMaxRows(max int) DbError

	Cancel() DbError

	ExecuteSql(sql string) (result int, error DbError)

	GetResultSet() (result ResultSet, error DbError)

	GetUpdateCount() (result int, error DbError)

	GetConnection() (connect Connection, error DbError)

	IsClose()(isClose bool, error DbError)
}

type ResultSet interface {

	Next() (result bool, error DbError)

	CLose() DbError

	GetString(columIndex int) (result string, error DbError)

	GetBool(olumIndex int) (result bool, error DbError)

	GetByte(olumIndex int) (result byte, error DbError)

	GetBytes(olumIndex int) (result []byte, error DbError)

	GetInt8(olumIndex int) (result int8, error DbError)

	GetInt16(olumIndex int) (result int16, error DbError)

	GetInt32(olumIndex int) (result int32, error DbError)

	GetInt64(olumIndex int) (result int64, error DbError)

	GetInt(olumIndex int) (result int, error DbError)

	GetFloat32(olumIndex int) (result float32, error DbError)

	GetFloat64(olumIndex int) (result float64, error DbError)

	GetTime(olumIndex int) (result time.Time, error DbError)

	//
	GetStringByName(columIndex int) (result string, error DbError)

	GetBoolByName(columName string) (result bool, error DbError)

	GetByteByName(columName string) (result byte, error DbError)

	GetBytesByName(columName string) (result []byte, error DbError)

	GetInt8ByName(columName string) (result int8, error DbError)

	GetInt16ByName(columName string) (result int16, error DbError)

	GetInt32ByName(columName string) (result int32, error DbError)

	GetInt64ByName(columName string) (result int64, error DbError)

	GetIntByName(columName string) (result int, error DbError)

	GetFloat32ByName(columName string) (result float32, error DbError)

	GetFloat64ByName(columName string) (result float64, error DbError)

	GetTimeByName(columName string) (result time.Time, error DbError)

	//
	IsFirst() (result bool, error DbError)

	IsLast()(result bool, error DbError)

	First()(result bool, error DbError)

	Last()(result bool, error DbError)

	GetRow()(result int, error DbError)

	Absolute(row int)(result bool, error DbError)

	Relative(row int)(result bool, error DbError)

	Previous()(result bool, error DbError)

}

type PrepareStatement interface {

	ExecuteQuery() (result ResultSet, error DbError)

	ExecuteUpdate() (result int, error DbError)

	Execute()(result bool, error DbError)

	SetNull(parameterIndex int, sqlType int) DbError

	SetBoolean(parameterIndex int, x bool) DbError

	SetByte(parameterIndex int, x byte) DbError

	SetInt(parameterIndex int, x int) DbError

	Setint8(parameterIndex int, x int8) DbError

	Setint16(parameterIndex int, x int16) DbError

	SetInt32(parameterIndex int, x int32) DbError

	SetInt64(parameterIndex int, x int64) DbError

	SetFloat32(parameterIndex int, x float32) DbError

	SetFloat64(parameterIndex int, x float64) DbError

	SetString(parameterIndex int, x string) DbError

	SetBytes(parameterIndex int, x [] byte) DbError

	SetTime(parameterIndex int, x time.Time) DbError

}
