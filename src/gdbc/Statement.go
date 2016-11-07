package gdbc

import "reflect"

//import "time"

type Field struct {
	Name      string
	Type      reflect.Kind
	Size      int
	Precision int
}

type Statement interface {
	Prepare(sql string) DbError

	Execute() DbError

	CreateResultSet() (result ResultSet, error DbError)

	Close() DbError

	SetValueByName(name string, value interface{})

	SetValue(values ...interface{}) DbError

	GetFields() (fields []Field, error DbError)

	GetUpdateRowNum() (int, DbError)
	//
	Query(sql string, values ...interface{}) (result ResultSet, error DbError)

	Update(sql string, values ...interface{}) (int, DbError)
}
