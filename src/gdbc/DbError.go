package gdbc

type DbError interface {
	GetErrorCode() int

	GetErrorMessage() string
}


