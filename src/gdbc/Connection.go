package gdbc

type Connection interface {
	CreateStatement() (statement Statement, error DbError)

	PrepareStatement(sql string) (prepareStatement PrepareStatement, error DbError)

	SetAutoCommit(autoCommit bool) DbError

	GetAutoCommit() (autoCommit bool, error DbError)

	Commit() DbError

	Rollback() DbError

	Close() DbError

	IsClose() (isClose bool, error DbError)
}

