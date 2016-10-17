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
	"unsafe"
)

type OCIError struct {
	ociCode  int
	ociMsg   string
	ociError unsafe.Pointer
}

func NewOCIError(ociCode int, v ...interface{}) *OCIError {
	if ociCode == OCI_SUCCESS {
		return nil
	}

	return &OCIError{ociCode: ociCode, ociMsg: fmt.Sprint(v)}
}

func NewOCIErrorf(ociCode int, format string, v ...interface{}) *OCIError {
	if ociCode == OCI_SUCCESS {
		return nil
	}

	return &OCIError{ociCode: ociCode, ociMsg: fmt.Sprintf(format, v)}
}

func MakeOCIError(ociError unsafe.Pointer) *OCIError {
	if ociError == nil {
		return nil
	}

	var errorCode = (C.sb4)(0)
	var maxlen = (C.ub4)(C.OCI_ERROR_MAXMSG_SIZE2)
	var msgC = C.malloc(C.OCI_ERROR_MAXMSG_SIZE2 + 1)
	var msg = (*C.OraText)((unsafe.Pointer(msgC)))

	C.OCIErrorGet(ociError, (C.ub4)(1), nil, &errorCode, msg, maxlen, OCI_HTYPE_ERROR)

	message := C.GoString((*C.char)(unsafe.Pointer(msg)))
	C.free(msgC)

	return &OCIError{ociCode: int(errorCode), ociMsg: message, ociError: ociError}
}

func OCIIsFailure(ociCode int) bool {
	switch ociCode {
	case OCI_SUCCESS:
		return false
	case OCI_SUCCESS_WITH_INFO:
		InfoLog("OCI_SUCCESS_WITH_INFO")
		return false
	case OCI_NO_DATA:
		InfoLog("OCI_NO_DATA")
		return false
	default:
		return true
	}
}

func (error *OCIError) Code() int {
	return error.ociCode
}

func (error *OCIError) Error() string {
	return error.ociMsg
}
