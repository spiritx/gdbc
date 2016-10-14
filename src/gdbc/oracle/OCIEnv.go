package oracle

//#cgo CFLAGS: -I /Users/xiebo/oracle/source/sdk/include
//#cgo LDFLAGS: -L /Users/xiebo/oracle/lib -lclntsh
//#include "stdlib.h"
//#include "oci.h"
//#include "stdio.h"
import "C"
import (
	. "gdbc"
	"strings"
	"unsafe"
)

const (
	OCIENV_UNINITAL = 0
	OCIENV_NORMAL   = 1
)

type OCIEnv struct {
	envhpp unsafe.Pointer

	status int
	info   map[string]string
}

func TransferOCIEnvMode(modeString string) uint32 {
	var mode uint32 = 0
	if strings.Contains(modeString, "OCI_DEFAULT") {
		mode = mode | C.OCI_DEFAULT
	}

	if strings.Contains(modeString, "OCI_THREADED") {
		mode = mode | C.OCI_THREADED
	}

	if strings.Contains(modeString, "OCI_OBJECT") {
		mode = mode | C.OCI_OBJECT
	}

	if strings.Contains(modeString, "OCI_EVENTS") {
		mode = mode | C.OCI_EVENTS
	}

	return mode
}

func (env *OCIEnv) Init(info map[string]string) (error DbError) {
	if env.status == OCIENV_NORMAL {
		env.Free()
		env.info = nil
	}
	env.info = make(map[string]string)
	for key, value := range info {
		info[key] = value
	}
	var mode uint32 = C.OCI_DEFAULT | C.OCI_THREADED | C.OCI_OBJECT
	if modeString, ok := info["OCIENV_MODE"]; ok {
		mode = TransferOCIEnvMode(modeString)
	}

	error = env.Create(mode)
	if error != nil && error.IsOk() {
		env.status = OCIENV_NORMAL
	}
	return
}

func (env *OCIEnv) GetConnection(info map[string]string) (Connection, DbError) {

	dblink := info[SID]
	if host, ok := info[HOST]; ok {
		dblink = host + ":" + info[PORT] + "/" + info[SID]
	}

	conn := OciConnection{}
	error := conn.Connect(env, info["UserName"], info["Password"], dblink)
	return conn, error
}

func (env *OCIEnv) Create(mode uint32) (error DbError) {
	var envhpp = (*C.OCIEnv)(env.envhpp)

	oraCode := C.OCIEnvCreate(&envhpp,
		C.ub4(mode), nil,
		nil, nil, nil, 0, nil)

	if OCIIsFailure(oraCode) {
		error = NewOCIError(int(oraCode), "OCIEnvCreate fail!")
		FatalLog(error.Error())
		return error
	}

	return nil
}

func (env *OCIEnv) Free() {
	C.OCIHandleFree(env.envhpp, OCI_HTYPE_ENV)
}

func (env *OCIEnv) Uninit() {
	env.Free()
	env.envhpp = nil
	env.info = nil
}
