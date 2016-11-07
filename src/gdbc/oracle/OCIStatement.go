package oracle

//#cgo CFLAGS: -I /Users/xiebo/oracle/source/sdk/include
//#cgo LDFLAGS: -L /Users/xiebo/oracle/lib -lclntsh
//#include "stdlib.h"
//#include "oci.h"
//#include "stdio.h"
import "C"
import (
	"fmt"
	. "gdbc"
	"reflect"
	"strings"
	"unsafe"
)

const (
	OCISTMT_UNINITAL = iota
	OCISTMT_PREPARE
	OCISTMT_EXECUTE
	OCISTMT_CLOSED
)

type OCIStatement struct {
	ociStmt       unsafe.Pointer
	conn          *OciConnection
	sql           string
	status        int
	result        *OCIResultSet
	inValueByName *CInValueByNameBuffer
}

func NewOCIStatement(conn *OciConnection) (Statement, DbError) {
	stmt := OCIStatement{conn: conn, inValueByName: NewCInValueBuffer(), status: OCISTMT_UNINITAL}
	error := stmt.allocHandle()
	if error != nil {
		return nil, error
	}

	return &stmt, nil
}

func (stmt *OCIStatement) allocHandle() DbError {
	if stmt.ociStmt != nil {
		stmt.freeHandle()
	}

	oraCode := int(C.OCIHandleAlloc(stmt.conn.env.envhpp, &stmt.ociStmt, OCI_HTYPE_STMT, 0, nil))
	if OCIIsFailure(oraCode) {
		error := NewOCIError(oraCode, "OCIHandleAlloc OCI_HTYPE_STMT fail!")
		FatalLog(error.ociCode, error.ociError)
		return error
	}

	return nil
}

func (stmt *OCIStatement) Close() DbError {
	if stmt.inValueByName != nil {
		stmt.inValueByName.Free()
		stmt.inValueByName = nil
	}

	if stmt.result != nil {
		stmt.result.Close()
		stmt.result = nil
	}
	stmt.status = OCISTMT_CLOSED

	return stmt.freeHandle()
}

func (stmt *OCIStatement) freeHandle() DbError {
	C.OCIHandleFree(stmt.ociStmt, OCI_HTYPE_STMT)
	stmt.ociStmt = nil
	stmt.conn = nil
	return nil
}

func (stmt *OCIStatement) Prepare(sql string) DbError {
	if stmt == nil || stmt.ociStmt == nil {
		error := NewOCIErrorf(-1, "Statement [%s] is nil!", sql)
		ErrorLogf("Prepare Error:", error.Error())
		return error
	}

	if stmt.status == OCISTMT_CLOSED {
		error := NewOCIErrorf(-1, "Statement [%s] is closed!", sql)
		ErrorLogf("Prepare Error:", error.Error())
		return error
	}

	if stmt.status == OCISTMT_EXECUTE {
		if stmt.inValueByName != nil && stmt.inValueByName.IsExistValue() {
			stmt.inValueByName.Free()
			stmt.inValueByName = nil
		}
		if stmt.result != nil {
			stmt.result.Close()
			stmt.result = nil
		}
	}

	stmt.sql = sql

	oraCode := int(C.OCIStmtPrepare(stmt.ociStmt, stmt.conn.ociError,
		(*C.OraText)(unsafe.Pointer(C.CString(sql))), C.ub4(len(sql)), OCI_NTV_SYNTAX, OCI_DEFAULT))
	if OCIIsFailure(oraCode) {
		error := MakeOCIError(stmt.conn.ociError)
		ErrorLogf("OCIStmtPrepare :+ %s +fail! oraCode = %d", error.Code(), error.Error())
		return error
	}

	stmt.status = OCISTMT_PREPARE

	return nil
}

func (stmt *OCIStatement) Execute() DbError {
	if stmt.status != OCISTMT_PREPARE {
		error := NewOCIError(-1, "stmt is not prepare!")
		ErrorLogf("Execute error:", error.Error())
		return error
	}

	if stmt.inValueByName.IsExistValue() {
		if error := stmt.inValueByName.makeCInValueMap(); error != nil {
			return error
		}

		if error := stmt.inValueByName.BindByName(stmt); error != nil {
			return error
		}
	}

	var sqlType C.ub2
	oraCode := int(C.OCIAttrGet(stmt.ociStmt, OCI_HTYPE_STMT, unsafe.Pointer(&sqlType), nil,
		C.OCI_ATTR_STMT_TYPE, stmt.conn.ociError))
	if OCIIsFailure(oraCode) {
		error := MakeOCIError(stmt.conn.ociError)
		ErrorLogf("OCIAttrGet OCI_ATTR_STMT_TYPE fail! sql:%s, oraCode = %d, %s",
			stmt.sql, error.Code(), error.Error())
		return error
	}

	var iters C.ub4 = 1
	if sqlType == C.OCI_STMT_SELECT {
		iters = 0
	}

	oraCode = int(C.OCIStmtExecute(stmt.conn.ociSvcCtx,
		stmt.ociStmt, stmt.conn.ociError, iters, 0, nil, nil, OCI_DEFAULT))
	if OCIIsFailure(oraCode) {
		error := MakeOCIError(stmt.conn.ociError)
		ErrorLogf("OCIStmtExecute fail! sql:%s, oraCode = %d, %s", stmt.sql, error.Code(), error.Error())
		return error
	}

	stmt.status = OCISTMT_EXECUTE
	return nil
}

type OCIField struct {
	Name      string
	Type      int
	Size      int
	Precision int
	Scale     int
}

func (stmt *OCIStatement) getFieldDescByIndex(index int) (field OCIField, error DbError) {
	var ociParam unsafe.Pointer
	oraCode := int(C.OCIParamGet(stmt.ociStmt, OCI_HTYPE_STMT, stmt.conn.ociError, &ociParam, C.ub4(index)))
	if OCIIsFailure(oraCode) {
		error = MakeOCIError(stmt.conn.ociError)
		ErrorLogf("OCIParamGet OCI_HTYPE_STMT fail! oraCode = %d, %s", oraCode, error.Error())
		return field, error
	}
	defer C.OCIDescriptorFree(ociParam, OCI_HTYPE_STMT)

	//field type
	var paramType C.ub2
	oraCode = int(C.OCIAttrGet(ociParam, OCI_DTYPE_PARAM, unsafe.Pointer(&paramType),
		nil, C.OCI_ATTR_DATA_TYPE, stmt.conn.ociError))
	if OCIIsFailure(oraCode) {
		error = MakeOCIError(stmt.conn.ociError)
		ErrorLogf("OCIAttrGet OCI_DTYPE_PARAM  OCI_ATTR_DATA_TYPE fail! oraCode = %d, %s",
			oraCode, error.Error())
		return field, error
	}
	field.Type = int(paramType)

	var paramName = ""
	paramNameC := C.CString(paramName)
	defer C.free(unsafe.Pointer(paramNameC))
	var paramNameLen C.ub4 = 1024
	oraCode = int(C.OCIAttrGet(ociParam, OCI_DTYPE_PARAM, unsafe.Pointer(&paramNameC),
		&paramNameLen, C.OCI_ATTR_NAME, stmt.conn.ociError))
	if OCIIsFailure(oraCode) {
		error = MakeOCIError(stmt.conn.ociError)
		ErrorLogf("OCIAttrGet OCI_DTYPE_PARAM  OCI_ATTR_NAME fail! oraCode = %d, %s",
			oraCode, error.Error())
		return field, error
	}
	field.Name = C.GoString(paramNameC)

	var paramSize C.ub2
	oraCode = int(C.OCIAttrGet(ociParam, OCI_DTYPE_PARAM, unsafe.Pointer(&paramSize),
		nil, C.OCI_ATTR_DATA_SIZE, stmt.conn.ociError))
	if OCIIsFailure(oraCode) {
		error = MakeOCIError(stmt.conn.ociError)
		ErrorLogf("OCIAttrGet OCI_DTYPE_PARAM  OCI_ATTR_DATA_SIZE fail! oraCode = %d, %s",
			oraCode, error.Error())
		return field, error
	}
	field.Size = int(paramSize)

	var paramPrecision C.ub2
	oraCode = int(C.OCIAttrGet(ociParam, OCI_DTYPE_PARAM, unsafe.Pointer(&paramPrecision),
		nil, C.OCI_ATTR_PRECISION, stmt.conn.ociError))
	if OCIIsFailure(oraCode) {
		error = MakeOCIError(stmt.conn.ociError)
		ErrorLogf("OCIAttrGet OCI_DTYPE_PARAM  OCI_ATTR_PRECISION fail! oraCode = %d, %s",
			oraCode, error.Error())
		return field, error
	}
	field.Precision = int(paramPrecision)

	var paramScale C.ub2
	oraCode = int(C.OCIAttrGet(ociParam, OCI_DTYPE_PARAM, unsafe.Pointer(&paramScale),
		nil, C.OCI_ATTR_SCALE, stmt.conn.ociError))
	if OCIIsFailure(oraCode) {
		error = MakeOCIError(stmt.conn.ociError)
		ErrorLogf("OCIAttrGet OCI_DTYPE_PARAM  OCI_ATTR_SCALE fail! oraCode = %d, %s",
			oraCode, error.Error())
		return field, error
	}
	field.Precision = int(paramScale)

	return field, nil
}

func (stmt *OCIStatement) getStmtParamCount() (int, DbError) {
	var paramCount C.ub4
	oraCode := int(C.OCIAttrGet(stmt.ociStmt, C.OCI_HTYPE_STMT, unsafe.Pointer(&paramCount), nil,
		C.OCI_ATTR_PARAM_COUNT, stmt.conn.ociError))

	if OCIIsFailure(oraCode) {
		error := MakeOCIError(stmt.conn.ociError)
		ErrorLogf("OCIAttrGet OCI_HTYPE_STMT  OCI_ATTR_PARAM_COUNT fail! oraCode = %d, %s",
			oraCode, error.Error())
		return 0, error
	}

	return int(paramCount), nil
}

func (stmt *OCIStatement) getOCIField() (fields []OCIField, error DbError) {
	var fieldCount int
	fieldCount, error = stmt.getStmtParamCount()
	if error != nil {
		return
	}

	if fieldCount < 1 {
		error = NewOCIError(-1, "field count is 0!")
		ErrorLog("GetFields error :", error.Error())
		return
	}
	//DebugLog("fieldCount:", fieldCount)

	fields = make([]OCIField, fieldCount)

	for i := 1; i <= fieldCount; i++ {
		fields[i-1], error = stmt.getFieldDescByIndex(i)
		if error != nil {
			ErrorLogf("GetFieldDescByIndex %d error:%s", i, error.Error())
			return
		}
		//DebugLog(i, fields[i - 1])
	}
	return fields, nil
}

func (stmt *OCIStatement) GetFields() ([]Field, DbError) {
	if stmt.status != OCISTMT_EXECUTE {
		error := NewOCIError(-1, "stmt status is not OCISTMT_EXECUTE!")
		ErrorLog("GetFields error:", error.Error())
		return nil, error
	}

	ociFields, error := stmt.getOCIField()
	if error != nil {
		return nil, error
	}

	fieldLen := len(ociFields)
	fields := make([]Field, fieldLen)
	for i := 0; i < fieldLen; i++ {
		fields[i].Name = ociFields[i].Name
		fields[i].Precision = ociFields[i].Precision
		fields[i].Size = ociFields[i].Size
		switch ociFields[i].Type {
		case C.SQLT_CHR, C.SQLT_AFC, C.SQLT_STR, C.SQLT_VCS, C.SQLT_RID:
			fields[i].Type = reflect.String
		case C.SQLT_NUM, C.SQLT_INT, C.SQLT_LNG:
			if ociFields[i].Precision > 0 {
				fields[i].Type = reflect.Float64
			} else if ociFields[i].Size > 18 {
				fields[i].Type = reflect.String
			} else {
				fields[i].Type = reflect.Int64
			}
		case C.SQLT_FLT:
			fields[i].Type = reflect.Float64
		default:
			error := NewOCIErrorf(-1, "column %s type(%d) is invalid", ociFields[i].Name, ociFields[i].Type)
			return nil, error
		}
	}

	return fields, nil
}

func (stmt *OCIStatement) CreateResultSet() (result ResultSet, error DbError) {
	if stmt.status != OCISTMT_EXECUTE {
		error = NewOCIError(-1, "status is not OCISTMT_EXECUTE!")
		ErrorLog("CreateResultSet error:", error.Error())
		return nil, error
	}

	if stmt.result != nil {
		stmt.result.Close()
	}
	stmt.result, error = NewOCIResultSet(stmt)

	return stmt.result, error
}

func (stmt *OCIStatement) SetValueByName(name string, value interface{}) {
	stmt.inValueByName.SetValue(name, value)
}

func (stmt *OCIStatement) SetValue(values ...interface{}) DbError {
	if len(values) < 1 {
		return nil
	}

	if stmt.status != OCISTMT_PREPARE {
		error := NewOCIErrorf(-1, "statement status is invalid [%d]", stmt.status)
		ErrorLog("SetValue Error:", error.Error())
		return error
	}

	names := splitPlaceHolder(stmt.sql)
	if names == nil {
		error := NewOCIErrorf(-1, "not found placeholder in sql [%s]", stmt.sql)
		ErrorLog("SetValue Error:", error.Error())
		return error
	}

	if len(names) != len(values) {
		error := NewOCIErrorf(-1, " placeholder mismatch in sql [%s] %d != %d ", stmt.sql, len(names), len(values))
		ErrorLog("SetValue Error:", error.Error())
		return error
	}

	for i, name := range names {
		stmt.inValueByName.SetValue(name, values[i])
	}

	return nil
}

func splitPlaceHolder(sql string) []string {
	count := strings.Count(sql, ":")
	if count > 0 {
		names := make([]string, count)
		start := false
		startpos := 0
		index := 0
		for i, rune := range sql {
			if start {
				switch rune {
				case '\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0, '>', '<', '!', '=':
					names[index] = sql[startpos:i]
					index++
					fmt.Printf("index=%d, startpos=%d, i =%d, value=%s\n", index, startpos, i, sql[startpos:i])
					start = false
				}
			} else if rune == ':' {
				start = true
				startpos = i
				fmt.Printf("%d, %s\n", i, sql[i:])
			}
		}
		if start {
			names[index] = sql[startpos:]
		}
		return names
	}
	return nil
}

func (stmt *OCIStatement) Query(sql string, values ...interface{}) (result ResultSet, error DbError) {
	if error = stmt.Prepare(sql); error != nil {
		return nil, error
	}

	if error = stmt.SetValue(values...); error != nil {
		return nil, error
	}

	if error = stmt.Execute(); error != nil {
		return nil, error
	}

	if result, error = stmt.CreateResultSet(); error != nil {
		return nil, error
	}
	return result, nil
}

func (stmt *OCIStatement) GetUpdateRowNum() (int, DbError) {
	if stmt.status != OCISTMT_EXECUTE {
		error := NewOCIErrorf(-1, "statement status is invalid [%d]", stmt.status)
		ErrorLog("GetUpdateRowNum Error:", error.Error())
		return 0, error
	}

	var updateRowNum C.int = 0
	oraCode := int(C.OCIAttrGet(stmt.ociStmt,
		C.ub4(OCI_HTYPE_STMT),
		unsafe.Pointer(&updateRowNum),
		nil,
		C.OCI_ATTR_ROW_COUNT,
		stmt.conn.ociError))

	if OCIIsFailure(oraCode) {
		error := MakeOCIError(stmt.conn.ociError)
		ErrorLogf("OCIAttrGet OCI_HTYPE_STMT  OCI_ATTR_ROW_COUNT fail! oraCode = %d, %s",
			oraCode, error.Error())
		return 0, error
	}
	return int(updateRowNum), nil
}

func (stmt *OCIStatement) Update(sql string, values ...interface{}) (int, DbError) {
	if error := stmt.Prepare(sql); error != nil {
		return 0, error
	}

	if error := stmt.SetValue(values...); error != nil {
		return 0, error
	}

	if error := stmt.Execute(); error != nil {
		return 0, error
	}

	if updateRowCount, error := stmt.GetUpdateRowNum(); error != nil {
		return 0, error
	} else {
		return updateRowCount, nil
	}
}
