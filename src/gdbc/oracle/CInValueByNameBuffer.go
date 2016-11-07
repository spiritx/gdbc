package oracle

//#cgo CFLAGS: -I /Users/xiebo/oracle/source/sdk/include
//#cgo LDFLAGS: -L /Users/xiebo/oracle/lib -lclntsh
//#include "stdlib.h"
//#include "oci.h"
//#include "stdio.h"
import "C"
import "reflect"

import (
	. "gdbc"
	"unsafe"
)

type CInValueByName struct {
	Value       interface{}
	CTypeName   string
	OCITypeCode int
	NameOffset  int
	ValueOffset int
	Capacity    int
	CSize       int
	index       int //序号
}

type CInValueByNameBuffer struct {
	valueMap map[string]CInValueByName
	buffer   *CBuffer
	Capacity int
}

func NewCInValueBuffer() *CInValueByNameBuffer {
	buff := &CInValueByNameBuffer{}
	buff.valueMap = make(map[string]CInValueByName)
	buff.buffer = nil
	return buff
}

func (buffer *CInValueByNameBuffer) SetValue(name string, value interface{}) {
	buffer.valueMap[name] = CInValueByName{Value: value}
}

func (buffer *CInValueByNameBuffer) IsExistValue() bool {
	return len(buffer.valueMap) > 0
}

func (buffer *CInValueByNameBuffer) makeCInValueMap() DbError {
	indicatorSize := len(buffer.valueMap) * C.sizeof_ub2
	offset := (indicatorSize + ALIGN_BYTES - 1) / ALIGN_BYTES * ALIGN_BYTES
	index := 0
	for name, value := range buffer.valueMap {
		nameLen := len(name) + 1
		value.NameOffset = offset
		value.index = index
		index++
		offset += nameLen
		switch value.Value.(type) {
		case int, int8, int16, int32, int64:
			value.CTypeName = CLONG
			value.OCITypeCode = SQLT_INT
			value.ValueOffset = offset
			value.CSize = C.sizeof_longlong
		case string:
			value.CTypeName = CSTRING
			value.OCITypeCode = SQLT_STR
			value.ValueOffset = offset
			originalValue, _ := value.Value.(string)
			value.CSize = len(originalValue) + 1
		case float32, float64:
			value.CTypeName = CDOUBLE
			value.OCITypeCode = SQLT_FLT
			value.ValueOffset = offset
			value.CSize = C.sizeof_double
		default:
			error := NewOCIErrorf(-1, "type(%s) is invalid", reflect.TypeOf(value.Value).String())
			ErrorLog("MakeCInValueMap error:", error.Error())
			return error
		}
		value.Capacity = (value.CSize + ALIGN_BYTES - 1) / ALIGN_BYTES * ALIGN_BYTES
		offset += value.Capacity
		buffer.valueMap[name] = value
	}

	buffer.Capacity = offset
	if buffer.buffer != nil {
		buffer.buffer.Free()
		buffer.buffer = nil
	}

	buffer.buffer = NewCBuffer(buffer.Capacity)
	for name, value := range buffer.valueMap {
		buffer.buffer.SetStringByOffset(value.NameOffset, name)
		switch value.Value.(type) {
		case int:
			originalValue, _ := value.Value.(int)
			buffer.buffer.SetInt64ByOffset(value.ValueOffset, int64(originalValue))
		case int8:
			originalValue, _ := value.Value.(int8)
			buffer.buffer.SetInt64ByOffset(value.ValueOffset, int64(originalValue))
		case int16:
			originalValue, _ := value.Value.(int16)
			buffer.buffer.SetInt64ByOffset(value.ValueOffset, int64(originalValue))
		case int32:
			originalValue, _ := value.Value.(int32)
			buffer.buffer.SetInt64ByOffset(value.ValueOffset, int64(originalValue))
		case int64:
			originalValue, _ := value.Value.(int64)
			buffer.buffer.SetInt64ByOffset(value.ValueOffset, originalValue)
		case float32:
			originalValue, _ := value.Value.(float32)
			buffer.buffer.SetFloatOffset(value.ValueOffset, float64(originalValue))
		case float64:
			originalValue, _ := value.Value.(float64)
			buffer.buffer.SetFloatOffset(value.ValueOffset, originalValue)
		case string:
			originalValue, _ := value.Value.(string)
			buffer.buffer.SetStringByOffset(value.ValueOffset, originalValue)
		default:
			error := NewOCIErrorf(-1, "type(%s) is invalid", reflect.TypeOf(value.Value).String())
			ErrorLog("MakeCInValueMap error:", error.Error())
			return error
		}
	}
	return nil
}

func (buffer *CInValueByNameBuffer) getIndicatorPointer(name string) unsafe.Pointer {
	if buffer.buffer == nil {
		ErrorLogf("getNamePoint %s error: buffer is nil!", name)
		return nil
	}

	value, ok := buffer.valueMap[name]
	if ok {
		return buffer.buffer.Offset(value.index * C.sizeof_ub2)
	}

	ErrorLogf("getNamePointer %s error: placeholder not found!", name)
	return nil
}

func (buffer *CInValueByNameBuffer) getNamePointer(name string) unsafe.Pointer {
	if buffer.buffer == nil {
		ErrorLogf("getNamePoint %s error: buffer is nil!", name)
		return nil
	}

	value, ok := buffer.valueMap[name]
	if ok {
		DebugLog("name:", C.GoString((*C.char)(buffer.buffer.Offset(value.NameOffset))))
		DebugLogf("value:%+v", value)
		return buffer.buffer.Offset(value.NameOffset)
	}

	ErrorLogf("getNamePointer %s error: placeholder not found!", name)
	return nil
}

func (buffer *CInValueByNameBuffer) getValuePointer(name string) unsafe.Pointer {
	if buffer.buffer == nil {
		ErrorLogf("getValuePointer %s error: buffer is nil!", name)
		return nil
	}

	value, ok := buffer.valueMap[name]
	if ok {
		if value.CTypeName == CSTRING {
			DebugLog("value:", C.GoString((*C.char)(buffer.buffer.Offset(value.ValueOffset))))
		} else if value.CTypeName == CLONG {
			DebugLog("value:", *((*C.longlong)(buffer.buffer.Offset(value.ValueOffset))))
		}

		return buffer.buffer.Offset(value.ValueOffset)
	}

	ErrorLogf("getValuePointer %s error: placeholder not found!", name)
	return nil
}

func (buffer *CInValueByNameBuffer) BindByName(stmt *OCIStatement) DbError {
	for name, value := range buffer.valueMap {
		var ociBind *C.OCIBind = nil
		oraCode := int(C.OCIBindByName(stmt.ociStmt,
			&ociBind,
			stmt.conn.ociError,
			(*C.OraText)(buffer.getNamePointer(name)),
			C.sb4(len(name)),
			buffer.getValuePointer(name),
			C.sb4(value.CSize),
			C.ub2(value.OCITypeCode),
			buffer.getIndicatorPointer(name),
			nil, nil, 0, nil,
			OCI_DEFAULT))

		if OCIIsFailure(oraCode) {
			error := MakeOCIError(stmt.conn.ociError)
			ErrorLogf("OCIBindByName  %s fail! oraCode = %d, %s", name,
				oraCode, error.Error())
			return error
		}
	}

	return nil
}

func (buffer *CInValueByNameBuffer) Free() {
	if buffer.buffer != nil {
		buffer.buffer.Free()
		buffer.Capacity = 0
		buffer.buffer = nil
	}
}
