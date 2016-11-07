package oracle

import . "gdbc"

//#cgo CFLAGS: -I /Users/xiebo/oracle/source/sdk/include
//#cgo LDFLAGS: -L /Users/xiebo/oracle/lib -lclntsh
//#include "oci.h"
import "C"

type OCIResultSet struct {
	stmt   *OCIStatement
	fields []OCIField
	buffer *COutValueBuffer
}

func NewOCIResultSet(stmt *OCIStatement) (*OCIResultSet, DbError) {
	var error DbError

	if stmt == nil {
		error = NewOCIError(-1, "stmt is nil")
		ErrorLog("NewOCIResultSet error:", error.Error())
		return nil, error
	}
	result := &OCIResultSet{stmt: stmt}
	result.fields, error = stmt.getOCIField()
	if error != nil {
		ErrorLog("NewOCIResultSet GetFields error:", error.Error())
		return nil, error
	}
	result.buffer = NewCValueBuffer(result.fields)
	error = result.buffer.bind(stmt)

	return result, error
}

func (result *OCIResultSet) Next() (isEnd bool, error DbError) {
	oraCode := C.OCIStmtFetch(result.stmt.ociStmt,
		result.stmt.conn.ociError,
		1, C.OCI_FETCH_NEXT, OCI_DEFAULT)
	switch oraCode {
	case OCI_SUCCESS, OCI_SUCCESS_WITH_INFO:
		return false, nil
	case OCI_NO_DATA:
		return true, nil
	}
	error = MakeOCIError(result.stmt.conn.ociError)
	return false, error
}

func (result *OCIResultSet) GetInt(index int) (value int, error DbError) {
	value, error = result.buffer.GetInt(index)
	return
}

func (result *OCIResultSet) GetFloat32(index int) (value float32, error DbError) {
	value, error = result.buffer.GetFloat32(index)
	return
}

func (result *OCIResultSet) GetFloat64(index int) (value float64, error DbError) {
	value, error = result.buffer.GetFloat64(index)
	return
}

func (result *OCIResultSet) GetInt64(index int) (value int64, error DbError) {
	value, error = result.buffer.GetInt64(index)
	return
}

func (result *OCIResultSet) GetString(index int) (value string, error DbError) {
	value, error = result.buffer.GetString(index)
	return
}

func (result *OCIResultSet) Close() {
	if result.buffer != nil {
		result.buffer.Free()
		result.buffer = nil
	}
}

func (result *OCIResultSet) CloseWithStatement() {
	result.Close()
	result.stmt.Close()
}
