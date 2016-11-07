package oracle

//#cgo CFLAGS: -I /Users/xiebo/oracle/source/sdk/include
//#cgo LDFLAGS: -L /Users/xiebo/oracle/lib -lclntsh
//#include "stdlib.h"
//#include "oci.h"
//#include "stdio.h"
import "C"

import (
	. "gdbc"
	"strconv"
	"unsafe"
)

type COutValue struct {
	Name           string
	Size           int
	Precision      int
	Scale          int
	OCITypeCode    int
	OutOCITypeCode int
	CTypeName      string
	offset         int
	Capacity       int
}

const ALIGN_BYTES int = 8

const (
	CSTRING   = "char*"
	CLONG     = "long long"
	CDOUBLE   = "double"
	CDATETIME = "OCIDate"
)

func MakeCOutValueList(fields []OCIField) (valueList []COutValue, capacity int) {
	fieldsLen := len(fields)
	valueList = make([]COutValue, fieldsLen)
	indicatorSize := fieldsLen * C.sizeof_ub2
	offset := (indicatorSize + ALIGN_BYTES - 1) / ALIGN_BYTES * ALIGN_BYTES

	for i := 0; i < fieldsLen; i++ {
		valueList[i].Name = fields[i].Name
		valueList[i].Size = fields[i].Size
		valueList[i].Precision = fields[i].Precision
		valueList[i].Scale = fields[i].Scale
		valueList[i].OCITypeCode = fields[i].Type
		valueList[i].offset = offset

		switch valueList[i].OCITypeCode {
		case C.SQLT_CHR, C.SQLT_AFC, C.SQLT_STR, C.SQLT_VCS, C.SQLT_RID: //VARCHAR2(n)
			valueList[i].CTypeName = CSTRING
			valueList[i].OutOCITypeCode = SQLT_STR
			valueList[i].Capacity = ((valueList[i].Size + 1 + ALIGN_BYTES - 1) / ALIGN_BYTES) * ALIGN_BYTES
		case C.SQLT_NUM, C.SQLT_INT, C.SQLT_LNG: //NUMBER
			if valueList[i].Precision > 0 {
				valueList[i].CTypeName = CDOUBLE
				valueList[i].OutOCITypeCode = SQLT_FLT
				valueList[i].Capacity = ((C.sizeof_double + ALIGN_BYTES - 1) / ALIGN_BYTES) * ALIGN_BYTES
			} else if valueList[i].Size > 18 {
				valueList[i].CTypeName = CSTRING
				valueList[i].OutOCITypeCode = SQLT_STR
				valueList[i].Capacity = ((valueList[i].Size + ALIGN_BYTES - 1) / ALIGN_BYTES) * ALIGN_BYTES
			} else {
				valueList[i].CTypeName = CLONG
				valueList[i].OutOCITypeCode = SQLT_INT
				valueList[i].Capacity = ((C.sizeof_longlong + ALIGN_BYTES - 1) / ALIGN_BYTES) * ALIGN_BYTES
			}
		case C.SQLT_FLT: //
			valueList[i].CTypeName = CDOUBLE
			valueList[i].OutOCITypeCode = SQLT_FLT
			valueList[i].Capacity = ((C.sizeof_double + ALIGN_BYTES - 1) / ALIGN_BYTES) * ALIGN_BYTES
		case C.SQLT_DAT: //DATE
			valueList[i].CTypeName = CDATETIME
			valueList[i].OutOCITypeCode = SQLT_DAT
			valueList[i].Capacity = ((C.sizeof_struct_OCIDate + ALIGN_BYTES - 1) / ALIGN_BYTES) * ALIGN_BYTES
		}
		offset += valueList[i].Capacity
	}
	return valueList, offset
}

type COutValueBuffer struct {
	buffer    *CBuffer //C buffer
	valueList []COutValue
	capacity  int
}

func NewCValueBuffer(fields []OCIField) *COutValueBuffer {
	buffer := &COutValueBuffer{}

	buffer.valueList, buffer.capacity = MakeCOutValueList(fields)
	buffer.buffer = NewCBuffer(buffer.capacity)

	return buffer
}

func (buffer *COutValueBuffer) Reset() {
	if buffer.buffer != nil {
		buffer.buffer.Reset()
	}
}

func (buffer *COutValueBuffer) getPointByIndex(index int) unsafe.Pointer {
	if buffer.buffer == nil {
		return nil
	}

	return buffer.buffer.Offset(buffer.valueList[index].offset)
}

func (buffer *COutValueBuffer) getIndicatorByIndex(index int) unsafe.Pointer {
	if buffer.buffer == nil {
		return nil
	}

	return buffer.buffer.Offset(index * C.sizeof_ub2)
}

func (buffer *COutValueBuffer) bind(stmt *OCIStatement) DbError {
	fieldsCount := len(buffer.valueList)

	for i := 0; i < fieldsCount; i++ {
		var ociDefine *C.OCIDefine = nil
		oraCode := int(C.OCIDefineByPos(stmt.ociStmt,
			&ociDefine,
			stmt.conn.ociError,
			C.ub4(i+1),
			buffer.getPointByIndex(i),
			C.sb4(buffer.valueList[i].Capacity),
			C.ub2(buffer.valueList[i].OutOCITypeCode),
			buffer.getIndicatorByIndex(i),
			nil, nil, OCI_DEFAULT))

		if OCIIsFailure(oraCode) {
			error := MakeOCIError(stmt.conn.ociError)
			ErrorLogf("OCIDefineByPos  fail! oraCode = %d, %s",
				oraCode, error.Error())
			return error
		}

	}
	return nil
}

func (buffer *COutValueBuffer) GetValue(index int) (interface{}, DbError) {
	if index > len(buffer.valueList) || index < 1 {
		error := NewOCIErrorf(-1, "index out range: %d >= %d or < 1", index, len(buffer.valueList))
		ErrorLog("GetValue", error.Error())
		return 0, error
	}

	fieldIndex := index - 1

	if buffer.valueList[fieldIndex].OutOCITypeCode == SQLT_INT {
		return buffer.buffer.GetInt64ByOffset(buffer.valueList[fieldIndex].offset), nil
	} else if buffer.valueList[fieldIndex].OutOCITypeCode == SQLT_FLT {
		return buffer.buffer.GetFloatOffset(buffer.valueList[fieldIndex].offset), nil
	} else if buffer.valueList[fieldIndex].OutOCITypeCode == SQLT_STR {
		return buffer.buffer.GetStringByOffset(buffer.valueList[fieldIndex].offset), nil
	}

	error := NewOCIErrorf(-1, "type mismatch error!interface{} <> %s", buffer.valueList[fieldIndex].CTypeName)
	ErrorLog("GetValue", error.Error())
	return 0, error
}

func (buffer *COutValueBuffer) GetInt(index int) (int, DbError) {
	if index > len(buffer.valueList) || index < 1 {
		error := NewOCIErrorf(-1, "index out range: %d >= %d or < 1", index, len(buffer.valueList))
		return 0, error
	}

	fieldIndex := index - 1

	if buffer.valueList[fieldIndex].OutOCITypeCode == SQLT_INT {
		return int(buffer.buffer.GetInt64ByOffset(buffer.valueList[fieldIndex].offset)), nil
	} else if buffer.valueList[fieldIndex].OutOCITypeCode == SQLT_FLT {
		return int(buffer.buffer.GetFloatOffset(buffer.valueList[fieldIndex].offset)), nil
	} else if buffer.valueList[fieldIndex].OutOCITypeCode == SQLT_STR {
		valueString := buffer.buffer.GetStringByOffset(buffer.valueList[fieldIndex].offset)
		if value, err := strconv.ParseInt(valueString, 10, 32); err == nil {
			return int(value), nil
		} else {
			DebugLogf("valueString:[%s]%s", valueString, err.Error())
		}
	}

	error := NewOCIErrorf(-1, "type mismatch error!int <> %s,%d",
		buffer.valueList[fieldIndex].CTypeName,
		buffer.valueList[fieldIndex].OCITypeCode)
	ErrorLog(error.Error())
	return 0, error
}

func (buffer *COutValueBuffer) GetInt64(index int) (value int64, error DbError) {
	if index > len(buffer.valueList) || index < 1 {
		error = NewOCIErrorf(-1, "index out range: %d >= %d or < 1", index, len(buffer.valueList))
		return 0, error
	}

	fieldIndex := index - 1

	if buffer.valueList[fieldIndex].OutOCITypeCode == SQLT_INT {
		return buffer.buffer.GetInt64ByOffset(buffer.valueList[fieldIndex].offset), nil
	} else if buffer.valueList[fieldIndex].OutOCITypeCode == SQLT_FLT {
		return int64(buffer.buffer.GetFloatOffset(buffer.valueList[fieldIndex].offset)), nil
	}

	error = NewOCIErrorf(-1, "type mismatch error!int64 <> %s", buffer.valueList[fieldIndex].CTypeName)
	ErrorLog(error.Error())
	return 0, error
}

func (buffer *COutValueBuffer) GetFloat32(index int) (value float32, error DbError) {
	if index > len(buffer.valueList) || index < 1 {
		error = NewOCIErrorf(-1, "index out range: %d >= %d or < 1", index, len(buffer.valueList))
		return 0, error
	}

	fieldIndex := index - 1

	if buffer.valueList[fieldIndex].OutOCITypeCode == SQLT_INT {
		return float32(buffer.buffer.GetInt64ByOffset(buffer.valueList[fieldIndex].offset)), nil
	} else if buffer.valueList[fieldIndex].OutOCITypeCode == SQLT_FLT {
		return float32(buffer.buffer.GetFloatOffset(buffer.valueList[fieldIndex].offset)), nil
	}

	error = NewOCIErrorf(-1, "type mismatch error!float32 <> %s", buffer.valueList[fieldIndex].CTypeName)
	ErrorLog(error.Error())
	return 0, error
}

func (buffer *COutValueBuffer) GetFloat64(index int) (value float64, error DbError) {
	if index > len(buffer.valueList) || index < 1 {
		error = NewOCIErrorf(-1, "index out range: %d >= %d or < 1", index, len(buffer.valueList))
		return 0, error
	}

	fieldIndex := index - 1

	if buffer.valueList[fieldIndex].OutOCITypeCode == SQLT_INT {
		return float64(buffer.buffer.GetInt64ByOffset(buffer.valueList[fieldIndex].offset)), nil
	} else if buffer.valueList[fieldIndex].OutOCITypeCode == SQLT_FLT {
		return buffer.buffer.GetFloatOffset(buffer.valueList[fieldIndex].offset), nil
	}

	error = NewOCIErrorf(-1, "type mismatch error!float64 <> %s", buffer.valueList[fieldIndex].CTypeName)
	ErrorLog(error.Error())
	return 0, error
}

func (buffer *COutValueBuffer) GetString(index int) (value string, error DbError) {
	if index > len(buffer.valueList) || index < 1 {
		error = NewOCIErrorf(-1, "index out range: %d >= %d or < 1", index, len(buffer.valueList))
		return "", error
	}

	fieldIndex := index - 1

	if buffer.valueList[fieldIndex].OutOCITypeCode == SQLT_INT {
		return strconv.FormatInt(buffer.buffer.GetInt64ByOffset(buffer.valueList[fieldIndex].offset), 10), nil
	} else if buffer.valueList[fieldIndex].OutOCITypeCode == SQLT_FLT {
		return strconv.FormatFloat(buffer.buffer.GetFloatOffset(buffer.valueList[fieldIndex].offset), 'f', -1, 64), nil
	} else if buffer.valueList[fieldIndex].OutOCITypeCode == SQLT_STR {
		return buffer.buffer.GetStringByOffset(buffer.valueList[fieldIndex].offset), nil
	}

	error = NewOCIErrorf(-1, "type mismatch error!string <> %s", buffer.valueList[fieldIndex].CTypeName)
	ErrorLog(error.Error())
	return "", error
}

func (buffer *COutValueBuffer) Free() {
	buffer.buffer.Free()
	buffer.capacity = 0
}
