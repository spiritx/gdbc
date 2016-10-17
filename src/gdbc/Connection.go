package gdbc

type Connection interface {
	CreateStatement() (statement Statement, error DbError)

	SetAutoCommit(autoCommit bool) DbError

	GetAutoCommit() (autoCommit bool, error DbError)

	Commit() DbError

	Rollback() DbError

	GetStatus() int

	IsClose() bool

	Close() DbError
}
