package gdbc

import "time"

type Statement interface {
	ExecuteQuery() (result ResultSet, error DbError)

	ExecuteUpdate() (result int, error DbError)

	Execute() (result bool, error DbError)

	Close() DbError

	GetMaxFieldSize() (max int, error DbError)

	SetMaxFieldSize(max int) DbError

	GetMaxRows() (max int, error DbError)

	SetMaxRows(max int) DbError

	Cancel() DbError

	ExecuteSql(sql string) (result int, error DbError)

	GetResultSet() (result ResultSet, error DbError)

	GetUpdateCount() (result int, error DbError)

	GetConnection() (connect Connection, error DbError)

	IsClose() (isClose bool, error DbError)
}

type PrepareStatement interface {
	ExecuteQuery() (result ResultSet, error DbError)

	ExecuteUpdate() (result int, error DbError)

	Execute() (result bool, error DbError)

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

	SetBytes(parameterIndex int, x []byte) DbError

	SetTime(parameterIndex int, x time.Time) DbError
}
