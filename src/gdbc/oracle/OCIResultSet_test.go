package oracle

import (
	"fmt"
	. "gdbc"
	_ "gdbc/oracle"
	"testing"
)

var DBURL = "gdbc:oracle:oci:@192.168.105.128:1521/tiger10"
var USERNAME = "dev"
var PASSWORD = "dev"

func TestOCIResultSet_GetInt(t *testing.T) {
	conn, error := GetConnectionByUser(DBURL, USERNAME, PASSWORD)
	if error != nil {
		fmt.Println("GetConnection error!", error.Error())
		t.Fail()
		return
	}
	if conn == nil {
		fmt.Println("conn is nil")
		t.Fail()
		return
	}
	fmt.Println("conn is ok")
	defer conn.Close()

	stmt, error := conn.CreateStatement()
	if error != nil {
		fmt.Println(error.Error())
		t.Fail()
		return
	}
	defer stmt.Close()

	if error = stmt.Prepare("select * from T_LOTTERY"); error != nil {
		fmt.Println(error.Error())
		t.Fail()
		return
	}

	if error = stmt.Execute(); error != nil {
		fmt.Println(error.Error())
		t.Fail()
		return
	}

	result, error := stmt.CreateResultSet()
	if error != nil {
		fmt.Println(error.Error())
		t.Fail()
		return
	}

	bEnd, error := result.Next()
	if error != nil {
		fmt.Println("next()", error.Error())
		t.Fail()
		return
	}

	if bEnd {
		t.Error("end")
		t.Fail()
	}

	v, error := result.GetInt(1)
	fmt.Println("v:", v)

}
