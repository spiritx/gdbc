package oracle


//#cgo CFLAGS: -I /Users/xiebo/oracle/source/sdk/include
//#cgo LDFLAGS: -L /Users/xiebo/oracle/lib -lclntsh
//#include "stdlib.h"
//#include "oci.h"
//#include "stdio.h"
import "C"
import (
	"unsafe"
	. "gdbc"
	"strconv"
)

type OCIEnv struct {
	envhpp unsafe.Pointer
}

type OCIError struct {
	errorCode int
	errorMsg  string
}

func (error *OCIError)GetErrorCode() int {
	return error.errorCode
}

func (error *OCIError)GetErrorMessage() string {
	return error.errorMsg
}

func OciFail(ociCode C.sword) bool {
	switch ociCode {
	case OCI_SUCCESS:
		return false;
	case OCI_SUCCESS_WITH_INFO:
		WriteInfoDbLog("OCI_SUCCESS_WITH_INFO")
		return false
	case OCI_NO_DATA:
		WriteInfoDbLog("OCI_NO_DATA")
		return false
	default:
		return true
	}
}

func OCIEnvCreate(mode uint32) (env OCIEnv, errorCode int) {
	var envhpp = (*C.OCIEnv)(env.envhpp)

	oraCode := C.OCIEnvCreate(&envhpp,
		C.ub4(mode), nil,
		nil, nil, nil, 0, nil)

	if OciFail(oraCode) {
		errorCode = int(oraCode)
		WriteFatalDbLog("OCIEnvCreate fail!")
	} else {
		errorCode = 0
	}

	return env, errorCode
}

func OCIEnvFree(env OCIEnv) {
	C.OCIHandleFree(env.envhpp, OCI_HTYPE_ENV)
}

type OciConnection struct {
	ociServer  unsafe.Pointer
	ociSvcCtx  unsafe.Pointer
	ociSession unsafe.Pointer
	ociError   unsafe.Pointer
	userName   string
	password   string
	dblink     string
	env        OCIEnv
}

func OCIAllocConnection(env OCIEnv) (conn OciConnection, ok bool) {
	oraCode := C.OCIHandleAlloc(env.envhpp, &conn.ociError, OCI_HTYPE_ERROR, 0, nil)
	if OciFail(oraCode) {
		WriteFatalDbLog("OCIHandleAlloc OCI_HTYPE_ERROR fail! oraCode = ", oraCode)
		return conn, false
	}

	oraCode = C.OCIHandleAlloc(env.envhpp, &conn.ociServer, OCI_HTYPE_SERVER, 0, nil)
	if (OciFail(oraCode)) {
		WriteFatalDbLog("OCIHandleAlloc OCI_HTYPE_SERVER fail! oraCode = ", oraCode)
		return conn, false
	}

	oraCode = C.OCIHandleAlloc(env.envhpp, &conn.ociSession, OCI_HTYPE_SESSION, 0, nil)
	if (OciFail(oraCode)) {
		WriteFatalDbLog("OCIHandleAlloc OCI_HTYPE_SESSION fail! oraCode = ", oraCode)
		return conn, false
	}

	oraCode = C.OCIHandleAlloc(env.envhpp, &conn.ociSvcCtx, OCI_HTYPE_SVCCTX, 0, nil)
	if (OciFail(oraCode)) {
		WriteFatalDbLog("OCIHandleAlloc OCI_HTYPE_SVCCTX fail! oraCode = ", oraCode)
		return conn, false
	}

	conn.env = env

	return conn, true
}

func OCIFreeConnection(conn *OciConnection) {
	if conn == nil {
		WriteDebugDbLog("connection is nil")
		return
	}

	if conn.ociSvcCtx != nil {
		if oraCode := C.OCIHandleFree(conn.ociSvcCtx, OCI_HTYPE_SVCCTX); oraCode != 0 {
			WriteFatalDbLog("OCIHandleFree OCI_HTYPE_SVCCTX fail! oraCode = ", oraCode)
		}
		conn.ociSvcCtx = nil
	}

	if conn.ociSession != nil {
		if oraCode := C.OCIHandleFree(conn.ociSession, OCI_HTYPE_SESSION); oraCode != 0 {
			WriteFatalDbLog("OCIHandleFree OCI_HTYPE_SESSION fail! oraCode = ", oraCode)
		}
		conn.ociSession = nil
	}

	if conn.ociServer != nil {
		if oraCode := C.OCIHandleFree(conn.ociServer, OCI_HTYPE_SERVER); oraCode != 0 {
			WriteFatalDbLog("OCIHandleFree OCI_HTYPE_SERVER fail! oraCode = ", oraCode)
		}
		conn.ociServer = nil
	}

	if conn.ociError != nil {
		if oraCode := C.OCIHandleFree(conn.ociError, OCI_HTYPE_ERROR); oraCode != 0 {
			WriteFatalDbLog("OCIHandleFree OCI_HTYPE_ERROR fail! oraCode = ", oraCode)
		}
		conn.ociError = nil
	}
}

func getOciError(ociError unsafe.Pointer) (int, string) {
	if ociError == nil {
		return -1, "Unknown oracle error!"
	}

	var errorCode = (C.sb4)(0)
	var maxlen = (C.ub4)(C.OCI_ERROR_MAXMSG_SIZE2)
	var msgC = C.malloc(C.OCI_ERROR_MAXMSG_SIZE2 + 1)
	var msg = (*C.OraText)((unsafe.Pointer(msgC)))

	C.OCIErrorGet(ociError, (C.ub4)(1), nil, &errorCode, msg, maxlen, OCI_HTYPE_ERROR)

	message := C.GoString((*C.char)(unsafe.Pointer(msg)))
	C.free(msgC)

	return int(errorCode), message
}

func GetOciError(conn *OciConnection) (int, string) {
	if conn == nil {
		return -1, "connection is nil."
	}

	return getOciError(conn.ociError)
}

func GetOciErrorMsg(conn *OciConnection) string {
	if conn == nil {
		return "oraCode = -1, connection is nil."
	}

	oraCode, msg := getOciError(conn.ociError)

	return "oraCode = " + strconv.Itoa(oraCode) + "," + msg
}

func OCIServerDetach(conn *OciConnection) {
	var ociError unsafe.Pointer
	oraCode := C.OCIHandleAlloc(conn.env.envhpp, &ociError, OCI_HTYPE_ERROR, 0, nil)
	if OciFail(oraCode) {
		errorCode, msg := getOciError(ociError)
		WriteFatalDbLog("OCIHandleAlloc OCI_HTYPE_ERROR fail! oraCode = , %d, %s", oraCode, errorCode, msg)
		return
	}

	if conn.ociSvcCtx != nil && conn.ociSession != nil {
		oraCode = C.OCISessionEnd(conn.ociSvcCtx, ociError, conn.ociSession, C.ub4(0))
		if OciFail(oraCode) {
			errorCode, msg := getOciError(ociError)
			WriteFatalDbLog("OCISessionEnd fail.oraCode =%d, %d, %s", oraCode, errorCode, msg)
		}
	}

	if conn.ociServer != nil {
		oraCode = C.OCIServerDetach(conn.ociServer, ociError, OCI_DEFAULT)
		if OciFail(oraCode) {
			errorCode, msg := getOciError(ociError)
			WriteFatalDbLog("OCIServerDetach fail.oraCode =%d, %d, %s", oraCode, errorCode, msg)
		}
	}

	OCIFreeConnection(conn)
}

func OCIServerAttach(env OCIEnv, userName string, password string, dblink string) (conn OciConnection, ok bool) {

	conn, ok = OCIAllocConnection(env)
	if !ok {
		WriteFatalDbLog("OCIAllocConnection Fail")
		OCIFreeConnection(&conn)
		return conn, false
	}

	conn.userName = userName
	conn.password = password
	conn.dblink = dblink

	oraCode := C.OCIServerAttach(conn.ociServer, conn.ociError,
		(*C.OraText)((unsafe.Pointer)(C.CString(conn.dblink))), C.sb4(len(conn.dblink)), OCI_DEFAULT)
	if OciFail(oraCode) {
		WriteErrorDbLogf("OCIServerAttach failed! dblink=%s, %s", dblink, GetOciErrorMsg(&conn))
		goto OUT
	}

	oraCode = C.OCIAttrSet(conn.ociSvcCtx, OCI_HTYPE_SVCCTX, conn.ociServer, 0, OCI_ATTR_SERVER, conn.ociError)
	if OciFail(oraCode) {
		WriteErrorDbLogf("OCIAttrSet OCI_HTYPE_SVCCTX failed! dblink=%s, %s", dblink, GetOciErrorMsg(&conn))
		goto OUT
	}

	oraCode = C.OCIAttrSet(conn.ociSession, OCI_HTYPE_SESSION,
		unsafe.Pointer(C.CString(userName)), C.ub4(len(userName)), OCI_ATTR_USERNAME, conn.ociError)
	if OciFail(oraCode) {
		WriteErrorDbLogf("OCIAttrSet OCI_HTYPE_SESSION failed! dblink=%s, userName =%s, %s",
			dblink, userName, GetOciErrorMsg(&conn))
		goto OUT
	}

	oraCode = C.OCIAttrSet(conn.ociSession, OCI_HTYPE_SESSION,
		unsafe.Pointer(C.CString(password)), C.ub4(len(password)), OCI_ATTR_PASSWORD, conn.ociError)
	if OciFail(oraCode) {
		WriteErrorDbLogf("OCIAttrSet OCI_HTYPE_SESSION failed! dblink=%s, userName =%s, password=%s, %s",
			dblink, userName, password, GetOciErrorMsg(&conn))
		goto OUT
	}

	oraCode = C.OCISessionBegin(conn.ociSvcCtx, conn.ociError, conn.ociSession, OCI_CRED_RDBMS, OCI_DEFAULT)
	if OciFail(oraCode) {
		WriteErrorDbLogf("OCISessionBegin failed! dblink=%s, userName =%s, password=%s, %s",
			dblink, userName, password, GetOciErrorMsg(&conn))
		goto OUT
	}

	oraCode = C.OCIAttrSet(conn.ociSvcCtx, OCI_HTYPE_SVCCTX, conn.ociSession, 0, OCI_ATTR_SESSION, conn.ociError)
	if OciFail(oraCode) {
		WriteErrorDbLogf("OCIAttrSet OCI_HTYPE_SVCCTX OCI_ATTR_SESSION failed! dblink=%s, userName =%s, password=%s, %s",
			dblink, userName, password, GetOciErrorMsg(&conn))
		goto OUT
	}

	return conn, true

	OUT:
	OCIServerDetach(&conn)

	return conn, false
}

