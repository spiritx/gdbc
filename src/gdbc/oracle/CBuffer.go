package oracle

//#cgo CFLAGS: -I /Users/xiebo/oracle/source/sdk/include
//#cgo LDFLAGS: -L /Users/xiebo/oracle/lib -lclntsh
//#include "stdlib.h"
//#include "oci.h"
//#include "stdio.h"
/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
void* MallocBuffer(int size)
{
	char* p = malloc(size);
	memset(p, 0, size);
	return p;
}

void FreeBuffer(void* p)
{
	free(p);
}

void SetInt64ByOffset(void* p, int offset, long long a)
{
	//printf("SetInt64ByOffset=%lld\n", a);
	*((long long*)((char*)p + offset)) = a;
}

long long GetInt64ByOffset(void* p, int offset)
{
	return *((long long*)((char*)p + offset));
}

void SetDoubleByOffset(void* p, int offset, double a)
{
	*((double*)((char*)p + offset)) = a;
}

double GetDoubleByOffset(void* p, int offset)
{
	return *((double*)((char*)p + offset));
}

char* GetStringByOffset(void* p, int offset)
{
	//printf("GetStringByOffset=[%s]\n", (char*)p + offset);
	return (char*)p + offset;
}

void SetStringByOffset(void* p, int offset, char* str)
{
	//printf("SetStringByOffset=[%s]\n", str);
	strcpy((char*)p + offset, str);
	*((char*)p + offset + strlen(str)) = '\0';
	//printf("s SetStringByOffset=[%s]\n", (char*)p + offset);

}

void* Offset(void* p, int offset)
{
	return (char*)p + offset;
}

*/
import "C"
import (
	. "gdbc"
	"unsafe"
)

type CBuffer struct {
	head unsafe.Pointer
	size int
}

func NewCBuffer(size int) *CBuffer {
	buffer := &CBuffer{nil, size}
	buffer.head = C.MallocBuffer(C.int(size))
	return buffer
}

func (buffer *CBuffer) Free() {
	if buffer.head != nil {
		C.FreeBuffer(buffer.head)
		buffer.head = nil
		buffer.size = 0
	}
}

func (buffer *CBuffer) isInvalid() bool {
	if buffer.head == nil || buffer.size <= 0 {
		return true
	}
	return false
}

func (buffer *CBuffer) SetInt64ByOffset(offset int, value int64) {
	if buffer.isInvalid() {
		ErrorLog("SetInt64ByOffset buffer is invalid!")
		panic("buffer is invalid!")
	}
	C.SetInt64ByOffset(buffer.head, C.int(offset), C.longlong(value))
}

func (buffer *CBuffer) GetInt64ByOffset(offset int) int64 {
	if buffer.isInvalid() {
		ErrorLog("GetInt64ByOffset buffer is invalid!")
		panic("buffer is invalid!")
	}
	return int64(C.GetInt64ByOffset(buffer.head, C.int(offset)))
}

func (buffer *CBuffer) SetFloatOffset(offset int, value float64) {
	if buffer.isInvalid() {
		ErrorLog("SetFloatOffset buffer is invalid!")
		panic("buffer is invalid!")
	}
	C.SetDoubleByOffset(buffer.head, C.int(offset), C.double(value))
}

func (buffer *CBuffer) GetFloatOffset(offset int) float64 {
	if buffer.isInvalid() {
		ErrorLog("GetFloatOffset buffer is invalid!")
		panic("buffer is invalid!")
	}
	return float64(C.GetDoubleByOffset(buffer.head, C.int(offset)))
}

func (buffer *CBuffer) SetStringByOffset(offset int, value string) {
	if buffer.isInvalid() {
		ErrorLog("SetStringByOffset buffer is invalid!")
		panic("buffer is invalid!")
	}

	str := C.CString(value)
	defer C.free(unsafe.Pointer(str))
	C.SetStringByOffset(buffer.head, C.int(offset), str)
}

func (buffer *CBuffer) GetStringByOffset(offset int) string {
	if buffer.isInvalid() {
		ErrorLog("GetStringByOffset buffer is invalid!")
		panic("buffer is invalid!")
	}

	return C.GoString(C.GetStringByOffset(buffer.head, C.int(offset)))
}

func (buffer *CBuffer) Reset() {
	if !buffer.isInvalid() {
		C.memset(buffer.head, 0, C.size_t(buffer.size))
	}
}

func (buffer *CBuffer) Offset(offset int) unsafe.Pointer {
	if buffer.isInvalid() {
		ErrorLog("Offset buffer is invalid!")
		panic("buffer is invalid!")
	}

	return C.Offset(buffer.head, C.int(offset))
}
