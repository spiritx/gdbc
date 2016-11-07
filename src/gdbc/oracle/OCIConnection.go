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

const (
	OCICONN_UNINITAL = 0
	OCICONN_NORMAL   = 1
)

type OciConnection struct {
	ociServer  unsafe.Pointer
	ociSvcCtx  unsafe.Pointer
	ociSession unsafe.Pointer
	ociError   unsafe.Pointer

	UserName string
	Password string
	Dblink   string
	env      *OCIEnv

	status int
}

func (conn *OciConnection) AllocHandle() DbError {
	if conn == nil || conn.env.envhpp == nil {
		ErrorLog("conn is nil or conn.env is nil")
		return NewOCIError(-1, "conn is nil or conn.env is nil")
	}

	oraCode := int(C.OCIHandleAlloc(conn.env.envhpp, &conn.ociError, OCI_HTYPE_ERROR, 0, nil))
	if OCIIsFailure(oraCode) {
		ErrorLog("OCIHandleAlloc OCI_HTYPE_ERROR fail! oraCode = ", oraCode)
		return NewOCIError(oraCode, "OCIHandleAlloc OCI_HTYPE_ERROR fail!")
	}

	oraCode = int(C.OCIHandleAlloc(conn.env.envhpp, &conn.ociServer, OCI_HTYPE_SERVER, 0, nil))
	if OCIIsFailure(oraCode) {
		ErrorLog("OCIHandleAlloc OCI_HTYPE_SERVER fail! oraCode = ", oraCode)
		return NewOCIError(oraCode, "OCIHandleAlloc OCI_HTYPE_SERVER fail!")
	}

	oraCode = int(C.OCIHandleAlloc(conn.env.envhpp, &conn.ociSession, OCI_HTYPE_SESSION, 0, nil))
	if OCIIsFailure(oraCode) {
		ErrorLog("OCIHandleAlloc OCI_HTYPE_SESSION fail! oraCode = ", oraCode)
		return NewOCIError(oraCode, "OCIHandleAlloc OCI_HTYPE_SESSION fail!")
	}

	oraCode = int(C.OCIHandleAlloc(conn.env.envhpp, &conn.ociSvcCtx, OCI_HTYPE_SVCCTX, 0, nil))
	if OCIIsFailure(oraCode) {
		ErrorLog("OCIHandleAlloc OCI_HTYPE_SVCCTX fail! oraCode = ", oraCode)
		return NewOCIError(oraCode, "OCIHandleAlloc OCI_HTYPE_SVCCTX fail!")
	}

	return nil
}

func (conn *OciConnection) FreeHandle() {
	if conn == nil {
		DebugLog("connection is nil")
		return
	}

	if conn.ociSvcCtx != nil {
		if oraCode := C.OCIHandleFree(conn.ociSvcCtx, OCI_HTYPE_SVCCTX); oraCode != 0 {
			ErrorLog("OCIHandleFree OCI_HTYPE_SVCCTX fail! oraCode = ", oraCode)
		}
		conn.ociSvcCtx = nil
	}

	if conn.ociSession != nil {
		if oraCode := C.OCIHandleFree(conn.ociSession, OCI_HTYPE_SESSION); oraCode != 0 {
			ErrorLog("OCIHandleFree OCI_HTYPE_SESSION fail! oraCode = ", oraCode)
		}
		conn.ociSession = nil
	}

	if conn.ociServer != nil {
		if oraCode := C.OCIHandleFree(conn.ociServer, OCI_HTYPE_SERVER); oraCode != 0 {
			ErrorLog("OCIHandleFree OCI_HTYPE_SERVER fail! oraCode = ", oraCode)
		}
		conn.ociServer = nil
	}

	if conn.ociError != nil {
		if oraCode := C.OCIHandleFree(conn.ociError, OCI_HTYPE_ERROR); oraCode != 0 {
			ErrorLog("OCIHandleFree OCI_HTYPE_ERROR fail! oraCode = ", oraCode)
		}
		conn.ociError = nil
	}
}

func (conn *OciConnection) ServerDetach() (error DbError) {
	var ociError unsafe.Pointer
	oraCode := int(C.OCIHandleAlloc(conn.env.envhpp, &ociError, OCI_HTYPE_ERROR, 0, nil))
	if OCIIsFailure(oraCode) {
		error = MakeOCIError(ociError)
		ErrorLogf("OCIHandleAlloc OCI_HTYPE_ERROR fail! oraCode = %d, %d, %s", oraCode, error.Code(), error.Error())
		return
	}

	if conn.ociSvcCtx != nil && conn.ociSession != nil {
		oraCode = int(C.OCISessionEnd(conn.ociSvcCtx, ociError, conn.ociSession, C.ub4(0)))
		if OCIIsFailure(oraCode) {
			error = MakeOCIError(ociError)
			ErrorLogf("OCISessionEnd fail.oraCode =%d, %d, %s", oraCode, error.Code(), error.Error())
		}
	}

	if conn.ociServer != nil {
		oraCode = int(C.OCIServerDetach(conn.ociServer, ociError, OCI_DEFAULT))
		if OCIIsFailure(oraCode) {
			error = MakeOCIError(ociError)
			ErrorLogf("OCIServerDetach fail.oraCode =%d, %d, %s", oraCode, error.Code(), error.Error())
		}
	}

	C.OCIHandleFree(ociError, OCI_HTYPE_ERROR)

	return nil
}

func (conn *OciConnection) ServerAttach() (error DbError) {

	oraCode := int(C.OCIServerAttach(conn.ociServer, conn.ociError,
		(*C.OraText)((unsafe.Pointer)(C.CString(conn.Dblink))), C.sb4(len(conn.Dblink)), OCI_DEFAULT))
	if OCIIsFailure(oraCode) {
		error = MakeOCIError(conn.ociError)
		ErrorLogf("OCIServerAttach failed! dblink=%s, %d,%s", conn.Dblink, error.Code(), error.Error())
		goto ERROR_OUT
	}

	oraCode = int(C.OCIAttrSet(conn.ociSvcCtx, OCI_HTYPE_SVCCTX, conn.ociServer, 0, OCI_ATTR_SERVER, conn.ociError))
	if OCIIsFailure(oraCode) {
		error = MakeOCIError(conn.ociError)
		ErrorLogf("OCIAttrSet OCI_HTYPE_SVCCTX failed! dblink=%s, %d,%s", conn.Dblink, error.Code(), error.Error())
		goto ERROR_OUT
	}

	oraCode = int(C.OCIAttrSet(conn.ociSession, OCI_HTYPE_SESSION,
		unsafe.Pointer(C.CString(conn.UserName)), C.ub4(len(conn.UserName)), OCI_ATTR_USERNAME, conn.ociError))
	if OCIIsFailure(oraCode) {
		error = MakeOCIError(conn.ociError)
		ErrorLogf("OCIAttrSet OCI_HTYPE_SESSION failed! dblink=%s, userName =%s, %d,%s",
			conn.Dblink, conn.UserName, error.Code(), error.Error())
		goto ERROR_OUT
	}

	oraCode = int(C.OCIAttrSet(conn.ociSession, OCI_HTYPE_SESSION,
		unsafe.Pointer(C.CString(conn.Password)), C.ub4(len(conn.Password)), OCI_ATTR_PASSWORD, conn.ociError))
	if OCIIsFailure(oraCode) {
		error = MakeOCIError(conn.ociError)
		ErrorLogf("OCIAttrSet OCI_HTYPE_SESSION failed! dblink=%s, userName =%s, password=%s, %d,%s",
			conn.Dblink, conn.UserName, conn.Password, error.Code(), error.Error())
		goto ERROR_OUT
	}

	oraCode = int(C.OCISessionBegin(conn.ociSvcCtx, conn.ociError, conn.ociSession, OCI_CRED_RDBMS, OCI_DEFAULT))
	if OCIIsFailure(oraCode) {
		error = MakeOCIError(conn.ociError)
		ErrorLogf("OCISessionBegin failed! dblink=%s, userName =%s, password=%s, %d,%s",
			conn.Dblink, conn.UserName, conn.Password, error.Code(), error.Error())
		goto ERROR_OUT
	}

	oraCode = int(C.OCIAttrSet(conn.ociSvcCtx, OCI_HTYPE_SVCCTX, conn.ociSession, 0, OCI_ATTR_SESSION, conn.ociError))
	if OCIIsFailure(oraCode) {
		error = MakeOCIError(conn.ociError)
		ErrorLogf("OCIAttrSet OCI_HTYPE_SVCCTX OCI_ATTR_SESSION failed! dblink=%s, userName =%s, password=%s, %d,%s",
			conn.Dblink, conn.UserName, conn.Password, error.Code(), error.Error())
		goto ERROR_OUT
	}

	return nil

ERROR_OUT:
	//conn.ServerDetach()

	return error
}

func (conn *OciConnection) Connect(env *OCIEnv, userName string, password string, dblink string) (error DbError) {
	if conn.status != OCICONN_UNINITAL {
		conn.ServerDetach()
		conn.FreeHandle()
		conn.status = OCICONN_UNINITAL
	}

	conn.env = env
	conn.UserName = userName
	conn.Password = password
	conn.Dblink = dblink

	if error = conn.AllocHandle(); error != nil {
		conn.FreeHandle()
		return
	}

	if error = conn.ServerAttach(); error != nil {
		conn.ServerDetach()
		conn.FreeHandle()
		return
	}

	conn.status = OCICONN_NORMAL
	return nil
}

func (conn *OciConnection) Disconnect() (error DbError) {
	error = conn.ServerDetach()
	conn.FreeHandle()
	conn.status = OCICONN_UNINITAL
	return
}

func (conn *OciConnection) Close() (error DbError) {
	return conn.Disconnect()
}

func (conn *OciConnection) CreateStatement() (statement Statement, error DbError) {
	if conn.status != OCICONN_NORMAL {
		error = NewOCIError(-1, "connection is invalid!")
		ErrorLog(error.Error())
		return nil, error
	}
	statement, error = NewOCIStatement(conn)
	return statement, error
}

func (conn *OciConnection) SetAutoCommit(autoCommit bool) DbError {
	if conn.status != OCICONN_NORMAL {
		error := NewOCIError(-1, "connection is invalid!")
		ErrorLog(error.Error())
		return error
	}

	return nil
}

func (conn *OciConnection) GetAutoCommit() (autoCommit bool, error DbError) {
	if conn.status != OCICONN_NORMAL {
		error = NewOCIError(-1, "connection is invalid!")
		ErrorLog(error.Error())
		return autoCommit, error
	}

	return autoCommit, nil
}

func (conn *OciConnection) Commit() DbError {
	if conn.status != OCICONN_NORMAL {
		error := NewOCIError(-1, "connection is invalid!")
		ErrorLog(error.Error())
		return error
	}

	oraCode := int(C.OCITransCommit(conn.ociSvcCtx, conn.ociError, OCI_DEFAULT))
	if OCIIsFailure(oraCode) {
		error := MakeOCIError(conn.ociError)
		ErrorLogf("OCITransCommit failed! dblink=%s, %d,%s",
			conn.Dblink, error.Code(), error.Error())
		return error
	}

	return nil
}

func (conn *OciConnection) Rollback() DbError {
	if conn.status != OCICONN_NORMAL {
		error := NewOCIError(-1, "connection is invalid!")
		ErrorLog(error.Error())
		return error
	}

	oraCode := int(C.OCITransRollback(conn.ociSvcCtx, conn.ociError, OCI_DEFAULT))
	if OCIIsFailure(oraCode) {
		error := MakeOCIError(conn.ociError)
		ErrorLogf("OCITransCommit failed! dblink=%s, %d,%s",
			conn.Dblink, error.Code(), error.Error())
		return error
	}

	return nil
}

func (conn *OciConnection) GetStatus() int {
	return conn.status
}

func (conn *OciConnection) IsClose() bool {
	return conn.status != OCICONN_NORMAL
}
