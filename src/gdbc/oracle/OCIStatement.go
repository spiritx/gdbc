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
	varlist []interface{}
}

func NewOCIStatement(conn *OciConnection) (Statement, DbError) {
	stmt := OCIStatement{conn: conn}
	error := stmt.AllocHandle()
	if error != nil {
		return nil, error
	}

	return stmt, nil
}

func (stmt *OCIStatement) AllocHandle() DbError {
	if stmt.ociStmt != nil {
		stmt.FreeHandle()
	}

	oraCode := int(C.OCIHandleAlloc(stmt.conn.env.envhpp, &stmt.ociStmt, OCI_HTYPE_STMT, 0, nil))
	if OCIIsFailure(oraCode) {
		error := NewOCIError(oraCode, "OCIHandleAlloc OCI_HTYPE_STMT fail!")
		FatalLog(error.ociCode, error.ociError)
		return error
	}

	return nil
}

func (stmt *OCIStatement) FreeHandle() DbError {
	if stmt == nil || stmt.ociStmt == nil {
		error := NewOCIError(-1, "stmt is nil")
		//ErrorLog(error.Error())
		return error
	}

	C.OCIHandleFree(stmt.ociStmt, OCI_HTYPE_STMT)
	stmt.ociStmt = nil
	stmt.conn = nil
	stmt.varlist = nil

	return nil
}

func (stmt *OCIStatement) Prepare(sql string) DbError {
	if stmt == nil || stmt.ociStmt == nil {
		error := NewOCIError(-1, "stmt is nil")
		ErrorLog(error.Error())
		return error
	}
	stmt.sql = sql

	oraCode := int(C.OCIStmtPrepare(stmt.ociStmt, stmt.conn.ociError,
		(*C.OraText)(unsafe.Pointer(C.CString(sql))), C.ub4(len(sql)), OCI_NTV_SYNTAX, OCI_DEFAULT))
	if OCIIsFailure(oraCode) {
		error := MakeOCIError(stmt.conn.ociError)
		FatalLog("OCIStmtPrepare :+ %s +fail! oraCode = ", error.Code(), error.Error())
		return error
	}

	return nil
}

func (stmt *OCIStatement) Execute() DbError {
	if stmt == nil || stmt.ociStmt == nil {
		error := NewOCIError(-1, "stmt is nil")
		ErrorLog(error.Error())
		return error
	}

	oraCode := int(C.OCIStmtExecute(stmt.conn.ociSvcCtx, stmt.ociStmt, stmt.conn.ociError, 0, 0, nil, nil, OCI_DEFAULT))
	if OCIIsFailure(oraCode) {
		error := MakeOCIError(stmt.conn.ociError)
		FatalLog("OCIStmtExecute :+ %s +fail! sql:%s, oraCode = ", stmt.sql, error.Code(), error.Error())
		return error
	}

	return nil
}

//type OCIField struct {
//	Name      string
//	Type      int
//	Size      int
//	Precision int
//	Scale     int
//}
//
//func OCIGetFieldDescByIndex(stmt *OCIStatement, index int) (field OCIField, ok bool) {
//	var ociParam unsafe.Pointer
//	oraCode := C.OCIParamGet(stmt.ociStmt, OCI_HTYPE_STMT, stmt.conn.ociError, &ociParam, C.ub4(index+1))
//	if OCIIsFailure(oraCode) {
//		ErrorLogf("OCIParamGet OCI_HTYPE_STMT fail! oraCode = %d, %s", oraCode, GetOciErrorMsg(stmt.conn))
//		return field, false
//	}
//
//	//field type
//	var paramType C.ub2
//	oraCode = C.OCIAttrGet(ociParam, OCI_DTYPE_PARAM, unsafe.Pointer(&paramType),
//		nil, C.OCI_ATTR_DATA_TYPE, stmt.conn.ociError)
//	if OCIIsFailure(oraCode) {
//		ErrorLogf("OCIAttrGet OCI_DTYPE_PARAM  OCI_ATTR_DATA_TYPE fail! oraCode = %d, %s",
//			oraCode, GetOciErrorMsg(stmt.conn))
//		return field, false
//	}
//	field.Type = int(paramType)
//
//	var paramName *C.char
//	var paramNameLen C.ub4
//	oraCode = C.OCIAttrGet(ociParam, OCI_DTYPE_PARAM, unsafe.Pointer(&paramName),
//		&paramNameLen, C.OCI_ATTR_NAME, stmt.conn.ociError)
//	if OCIIsFailure(oraCode) {
//		ErrorLogf("OCIAttrGet OCI_DTYPE_PARAM  OCI_ATTR_NAME fail! oraCode = %d, %s",
//			oraCode, GetOciErrorMsg(stmt.conn))
//		return field, false
//	}
//	field.Name = C.GoString(paramName)
//	//C.free(paramName)
//
//	var paramSize C.ub2
//	oraCode = C.OCIAttrGet(ociParam, OCI_DTYPE_PARAM, unsafe.Pointer(&paramSize),
//		nil, C.OCI_ATTR_DATA_SIZE, stmt.conn.ociError)
//	if OCIIsFailure(oraCode) {
//		ErrorLogf("OCIAttrGet OCI_DTYPE_PARAM  OCI_ATTR_DATA_SIZE fail! oraCode = %d, %s",
//			oraCode, GetOciErrorMsg(stmt.conn))
//		return field, false
//	}
//	field.Size = int(paramSize)
//
//	var paramPrecision C.ub2
//	oraCode = C.OCIAttrGet(ociParam, OCI_DTYPE_PARAM, unsafe.Pointer(&paramPrecision),
//		nil, C.OCI_ATTR_PRECISION, stmt.conn.ociError)
//	if OCIIsFailure(oraCode) {
//		ErrorLogf("OCIAttrGet OCI_DTYPE_PARAM  OCI_ATTR_PRECISION fail! oraCode = %d, %s",
//			oraCode, GetOciErrorMsg(stmt.conn))
//		return field, false
//	}
//	field.Precision = int(paramPrecision)
//
//	var paramScale C.ub2
//	oraCode = C.OCIAttrGet(ociParam, OCI_DTYPE_PARAM, unsafe.Pointer(&paramScale),
//		nil, C.OCI_ATTR_SCALE, stmt.conn.ociError)
//	if OCIIsFailure(oraCode) {
//		ErrorLogf("OCIAttrGet OCI_DTYPE_PARAM  OCI_ATTR_SCALE fail! oraCode = %d, %s",
//			oraCode, GetOciErrorMsg(stmt.conn))
//		return field, false
//	}
//	field.Precision = int(paramScale)
//
//	return field, true
//}

//func OCIGetStmtParamCount(stmt *OCIStatement) (int, bool) {
//	var paramCount C.ub4
//	oraCode := C.OCIAttrGet(stmt.ociStmt, C.OCI_HTYPE_STMT, unsafe.Pointer(&paramCount), nil,
//		C.OCI_ATTR_PARAM_COUNT, stmt.conn.ociError)
//
//	if OCIIsFailure(oraCode) {
//		ErrorLogf("OCIAttrGet OCI_HTYPE_STMT  OCI_ATTR_PARAM_COUNT fail! oraCode = %d, %s",
//			oraCode, GetOciErrorMsg(stmt.conn))
//		return 0, false
//	}
//
//	return int(paramCount), true
//}
