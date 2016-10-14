package oracle

//#cgo CFLAGS: -I /Users/xiebo/oracle/source/sdk/include
//#cgo LDFLAGS: -L /Users/xiebo/oracle/lib -lclntsh
//#include "stdlib.h"
//#include "oci.h"
//#include "stdio.h"
import "C"
import (
	. "gdbc"
	"unsafe"
)

type OCIStatement struct {
	ociStmt unsafe.Pointer
	conn    *OciConnection
	sql     string
}

func OCIStmtRelease(stmt *OCIStatement) {
	if stmt == nil || stmt.ociStmt == nil {
		return
	}

	C.OCIHandleFree(stmt.ociStmt, OCI_HTYPE_STMT)
	stmt.ociStmt = nil
	stmt.conn = nil
}

func OCIStmtPrepare(conn OciConnection, sql string) (stmt OCIStatement, ok bool) {
	oraCode := C.OCIHandleAlloc(conn.env.envhpp, &stmt.ociStmt, OCI_HTYPE_STMT, 0, nil)
	if OCIFail(oraCode) {
		FatalLog("OCIHandleAlloc OCI_HTYPE_STMT fail! oraCode = ", oraCode)
		return stmt, false
	}

	stmt.conn = &conn
	stmt.sql = sql

	oraCode = C.OCIStmtPrepare(stmt.ociStmt, conn.ociError,
		(*C.OraText)(unsafe.Pointer(C.CString(sql))), C.ub4(len(sql)), OCI_NTV_SYNTAX, OCI_DEFAULT)

	if OCIFail(oraCode) {
		FatalLog("OCIStmtPrepare :+ %s +fail! oraCode = ", sql, oraCode)
		OCIStmtRelease(&stmt)
		return stmt, false
	}

	return stmt, true
}

type OCIField struct {
	Name      string
	Type      int
	Size      int
	Precision int
	Scale     int
}

func OCIGetFieldDescByIndex(stmt *OCIStatement, index int) (field OCIField, ok bool) {
	var ociParam unsafe.Pointer
	oraCode := C.OCIParamGet(stmt.ociStmt, OCI_HTYPE_STMT, stmt.conn.ociError, &ociParam, C.ub4(index+1))
	if OCIFail(oraCode) {
		ErrorLogf("OCIParamGet OCI_HTYPE_STMT fail! oraCode = %d, %s", oraCode, GetOciErrorMsg(stmt.conn))
		return field, false
	}

	//field type
	var paramType C.ub2
	oraCode = C.OCIAttrGet(ociParam, OCI_DTYPE_PARAM, unsafe.Pointer(&paramType),
		nil, C.OCI_ATTR_DATA_TYPE, stmt.conn.ociError)
	if OCIFail(oraCode) {
		ErrorLogf("OCIAttrGet OCI_DTYPE_PARAM  OCI_ATTR_DATA_TYPE fail! oraCode = %d, %s",
			oraCode, GetOciErrorMsg(stmt.conn))
		return field, false
	}
	field.Type = int(paramType)

	var paramName *C.char
	var paramNameLen C.ub4
	oraCode = C.OCIAttrGet(ociParam, OCI_DTYPE_PARAM, unsafe.Pointer(&paramName),
		&paramNameLen, C.OCI_ATTR_NAME, stmt.conn.ociError)
	if OCIFail(oraCode) {
		ErrorLogf("OCIAttrGet OCI_DTYPE_PARAM  OCI_ATTR_NAME fail! oraCode = %d, %s",
			oraCode, GetOciErrorMsg(stmt.conn))
		return field, false
	}
	field.Name = C.GoString(paramName)
	//C.free(paramName)

	var paramSize C.ub2
	oraCode = C.OCIAttrGet(ociParam, OCI_DTYPE_PARAM, unsafe.Pointer(&paramSize),
		nil, C.OCI_ATTR_DATA_SIZE, stmt.conn.ociError)
	if OCIFail(oraCode) {
		ErrorLogf("OCIAttrGet OCI_DTYPE_PARAM  OCI_ATTR_DATA_SIZE fail! oraCode = %d, %s",
			oraCode, GetOciErrorMsg(stmt.conn))
		return field, false
	}
	field.Size = int(paramSize)

	var paramPrecision C.ub2
	oraCode = C.OCIAttrGet(ociParam, OCI_DTYPE_PARAM, unsafe.Pointer(&paramPrecision),
		nil, C.OCI_ATTR_PRECISION, stmt.conn.ociError)
	if OCIFail(oraCode) {
		ErrorLogf("OCIAttrGet OCI_DTYPE_PARAM  OCI_ATTR_PRECISION fail! oraCode = %d, %s",
			oraCode, GetOciErrorMsg(stmt.conn))
		return field, false
	}
	field.Precision = int(paramPrecision)

	var paramScale C.ub2
	oraCode = C.OCIAttrGet(ociParam, OCI_DTYPE_PARAM, unsafe.Pointer(&paramScale),
		nil, C.OCI_ATTR_SCALE, stmt.conn.ociError)
	if OCIFail(oraCode) {
		ErrorLogf("OCIAttrGet OCI_DTYPE_PARAM  OCI_ATTR_SCALE fail! oraCode = %d, %s",
			oraCode, GetOciErrorMsg(stmt.conn))
		return field, false
	}
	field.Precision = int(paramScale)

	return field, true
}

func OCIGetStmtParamCount(stmt *OCIStatement) (int, bool) {
	var paramCount C.ub4
	oraCode := C.OCIAttrGet(stmt.ociStmt, C.OCI_HTYPE_STMT, unsafe.Pointer(&paramCount), nil,
		C.OCI_ATTR_PARAM_COUNT, stmt.conn.ociError)

	if OCIFail(oraCode) {
		ErrorLogf("OCIAttrGet OCI_HTYPE_STMT  OCI_ATTR_PARAM_COUNT fail! oraCode = %d, %s",
			oraCode, GetOciErrorMsg(stmt.conn))
		return 0, false
	}

	return int(paramCount), true
}
