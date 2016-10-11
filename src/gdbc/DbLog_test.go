package gdbc

import (
	"testing"
	"fmt"
	"runtime"
	"path/filepath"
	"os"
)

func TestDebugLog(t *testing.T) {
	if !DebugLog(" 中文测试", "this is a test") {
		t.Fail()
	}
}

func TestDebugLogf(t *testing.T) {
	if !DebugLogf("=%s,%d,%f", "aaaa", 10, 10.0) {
		t.Fail()
	}
}

func TestErrorLog(t *testing.T) {
	if !ErrorLog("error") {
		t.Fail()
	}
}

func TestErrorLogf(t *testing.T) {
	if !ErrorLogf(" 中%s文", "中文汉字，偏门值") {
		t.Fail()
	}
}

func TestFatalLog(t *testing.T) {
	if !FatalLog("严重错误") {
		t.Fail()
	}
}

func TestFatalLogf(t *testing.T) {
	if !FatalLogf("严重错误%s%+v", "error", t) {
		t.Fail()
	}
}

func TestInfoLog(t *testing.T) {
	if !InfoLog("提示信息", t) {
		t.Fail()
	}
}

func TestInfoLogf(t *testing.T) {
	if !InfoLogf("中文测试") {
		t.Fail()
	}
}

func TestDbLog_LogWrite(t *testing.T) {
	SetDbLog("./", "aa", "INFO|TRACE", 1024)

	if WriteDbLog(LOG_DEBUG, "This is a test") {
		t.Error("no LOG_DEBUG")
	}

	if !InfoLog("abc") {
		t.Error("INFO")
	}
}

func TestSetDbLog(t *testing.T) {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fmt.Println("dir:", dir)
	SetDbLog(dir, "bb", "FATAL", 1024*1024)
	if WriteDbLog(LOG_DEBUG, "This is a test") {
		t.Error("no LOG_DEBUG")
	}

	if !FatalLog("abc") {
		t.Error("FATAL")
	}
}

func TestSetDbLogLevel(t *testing.T) {
	SetDbLogLevel("DEBUG|TRACE")
	if !WriteDbLog(LOG_DEBUG, "!!!!!TestSetDbLogLevel") {
		t.Error("no LOG_DEBUG")
	}

	if InfoLog("abc") {
		t.Error("INFO")
	}
}

func TestWarnLog(t *testing.T) {
	SetDbLogLevel("TRACE|WARN")
	if !WarnLog("TestWarnLog") {
		t.Error("WARN")
	}
	if DebugLog("abc") {
		t.Error("no DEBUG")
	}
}

func TestWarnLogf(t *testing.T) {
	SetDbLogLevel("WARN|TRACE")

	if !WarnLogf("warn:%s,%d", "aa", 10) {
		t.Fail()
	}
}

func TestWriteDbLog(t *testing.T) {
	SetDbLogLevel("ERROR|TRACE")
	if !WriteDbLog(LOG_ERROR, "ERROR", 10, 10.1, true) {
		t.Fail()
	}
}

func TestWriteDbLogf(t *testing.T) {
	SetDbLogLevel("ERROR|TRACE")
	if !WriteDbLogf(LOG_ERROR, "ERROR%d,%f,%t", 10, 10.1, true) {
		t.Fail()
	}

}

func BenchmarkDebugLog(b *testing.B) {
	SetDbLog("./", "aa", "DEBUG|TRACE", 1024)

	for i := 0; i < b.N; i++ {
		if !WriteDbLog(LOG_DEBUG, "This is a test") {
			b.Fail()
		}
	}

}





func BenchmarkDbLog_WriteDbLog2(b *testing.B) {
	count := 10
	in := make(chan int, 1)
	out := make(chan int, 1)
	counter := make(chan int, count)

	fmt.Println("GOMAXPROCS=", runtime.GOMAXPROCS(4))
	SetDbLog("./", "aa", "DEBUG|TRACE", 1024*1024)
	for i := 0; i < count; i++ {
		in <- i
		go func() {
			id := <-in
			out <- 0
			for k := 0; k < b.N; k++ {
				//fmt.Println(id)
				if !WriteDbLogf(LOG_DEBUG, "中文可变字符测试abc%d,%d", id, k) {
					b.Fail()
				}
				runtime.Gosched()
			}
			counter <- 1
		}()
		<-out
	}

	for i := 0; i < 10; i++ {
		<-counter
	}
}
