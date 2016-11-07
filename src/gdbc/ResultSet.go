package gdbc

//import "time"

type ResultSet interface {
	Next() (isEnd bool, error DbError)

	CLose()

	CloseWithStatement()

	GetString(index int) (value string, error DbError)

	GetInt64(index int) (value int64, error DbError)

	GetInt(index int) (value int, error DbError)

	GetFloat32(index int) (value float32, error DbError)

	GetFloat64(index int) (value float64, error DbError)
}
