package oracle

import (
	. "gdbc"
	"strconv"
	"strings"
)

const (
	PROTOCOL = "gdbc"
	DBTYPE   = "oracle"
	DBAPI    = "oci"
	HOST     = "Host"
	PORT     = "Port"
	SID      = "SID"
)

type OracleDriver struct {
	env OCIEnv
}

//gdbc:oracle:oci:@<host>:<port>:<SID>
//gdbc:oracle:oci:@//<host>:<port>/ServiceName
func (driver *OracleDriver) ParseUrl(url string, info map[string]string) DbError {
	isSeparator := func(c rune) bool {
		switch c {
		case '.', '/':
			return true
		default:
			return false
		}
	}

	protocols := strings.FieldsFunc(url, isSeparator)
	if len(protocols) != 6 {
		return NewDefaultDbError(-1, "URL:"+url+" is invalid!")
	}

	if protocols[0] != PROTOCOL || protocols[2] != DBTYPE || protocols[3] != DBAPI {
		return NewDefaultDbError(-1, "URL:"+url+" is invalid!")
	}

	if info == nil {
		info = make(map[string]string)
	}

	if len(protocols[4] > 0) {
		if len(protocols[4]) < 2 || protocols[4][0] != '@' {
			return NewDefaultDbError(-1, "URL:"+url+" is invalid! host is invalid!")
		}

		_, err := strconv.Atoi(protocols[5])
		if err != nil {
			return NewDefaultDbError(-1, "URL:"+url+" is invalid! port is invalid!"+err.Error())
		}

		if len(protocols[4] > 3) && protocols[4][0:3] == "@//" {
			info[HOST] = protocols[4][4:]
		} else {
			info[HOST] = protocols[4][1:]
		}

		info[PORT] = protocols[5]
	}

	info[SID] = protocols[6]

	return nil
}

func (driver *OracleDriver) Connect(url string, info map[string]string) (Connection,
	error DbError) {
	if error = driver.ParseUrl(url, info); error.IsFailure() {
		ErrorLogf("ParseUrl error!url:%s,%d,%s", url, error.Code(), error.Error())
		return
	}

	if error = driver.env.Init(info); error.IsFailure() {
		ErrorLogf("env init error!url:%s,%d,%s", url, error.Code(), error.Error())
		return
	}

	return driver.env.GetConnection(info)
}
