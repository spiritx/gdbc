package gdbc

import (
	"testing"
)

func TestDbLog_LogWrite(t *testing.T) {
	log := NewDbLog("/Users/xiebo/IdeaProjects/Test/src/gdbc", "aa", "INFO|TRACE", 1024)

	log.WriteDbLog(LOG_DEBUG, "This is a test")
}

func TestDbLog_LogWriteObj(t *testing.T) {
	log := NewDbLog("/Users/xiebo/IdeaProjects/Test/src/gdbc", "aa", "DEBUG|TRACE", 1024)

	log.WriteDbLogf(LOG_DEBUG, "%+v", log)
}

func TestDbLog_LogWritef(t *testing.T) {
	log := NewDbLog("/Users/xiebo/IdeaProjects/Test/src/gdbc", "aa", "DEBUG|TRACE", 1024)

	log.WriteDbLogf(LOG_DEBUG, "This is 测试%d,%s", 100, "hello")
}

func BenchmarkDbLog_LogWrite(b *testing.B) {
	log := NewDbLog("/Users/xiebo/IdeaProjects/Test/src/gdbc", "aa", "DEBUG|TRACE", 1024)

	for i:=0; i<b.N; i++ {
		log.WriteDbLog(LOG_DEBUG, "This is a test")
	}
}

func TestDbLog_WriteDbLog(t *testing.T) {
	WriteDbLog(LOG_DEBUG, "This is a test", 100, 10.11, "测试中文")
}

func TestDbLog_WriteDbLogf(t *testing.T) {
	WriteDbLogf(LOG_DEBUG, "中文可变字符测试abc%s,%d,%s", "OK", 10, "中文参数")
}

func TestDbLog_WriteDbLog2(t *testing.T) {
	count := 10
	in := make(chan int, 1)
	out := make(chan int, 1)
	counter := make(chan int, count)

	for i := 0; i < count; i++ {
		in <- i
		go func() {
			id := <- in
			out <- 0
			for k := 0; k < 10000; k++ {
				//fmt.Println(id)
				WriteDbLogf(LOG_DEBUG, "中文可变字符测试abc%d,%d", id, k)
			}
			counter <- 1
		}()
		<- out
	}

	for i := 0; i < 10; i++ {
		<- counter
	}
}