package gdbc

import "time"

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

	IsLast() (result bool, error DbError)

	First() (result bool, error DbError)

	Last() (result bool, error DbError)

	GetRow() (result int, error DbError)

	Absolute(row int) (result bool, error DbError)

	Relative(row int) (result bool, error DbError)

	Previous() (result bool, error DbError)
}
