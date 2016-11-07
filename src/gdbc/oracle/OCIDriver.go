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

type OCIDriver struct {
	env OCIEnv
}

//gdbc:oracle:oci:@<host>:<port>:<SID>
//gdbc:oracle:oci:@//<host>:<port>/ServiceName
func (driver *OCIDriver) ParseUrl(url string, info map[string]string) DbError {
	isSeparator := func(c rune) bool {
		switch c {
		case ':', '/':
			return true
		default:
			return false
		}
	}

	protocols := strings.FieldsFunc(url, isSeparator)
	if len(protocols) != 6 {
		return NewDefaultDbErrorf(-1, "URL:%s is invalid! format is invalid:%d", url, len(protocols))
	}

	if protocols[0] != PROTOCOL || protocols[1] != DBTYPE || protocols[2] != DBAPI {
		return NewDefaultDbError(-1, "URL:"+url+" is invalid! protocol is invalid!")
	}

	if info == nil {
		info = make(map[string]string)
	}

	if len(protocols[3]) > 0 {
		if len(protocols[3]) < 2 || protocols[3][0] != '@' {
			return NewDefaultDbErrorf(-1, "URL:%s is invalid! host(%s) is invalid!", url, protocols[3])
		}

		_, err := strconv.Atoi(protocols[4])
		if err != nil {
			return NewDefaultDbErrorf(-1, "URL:%s is invalid! port(%s) is invalid!%s", url, protocols[4], err.Error())
		}

		if len(protocols[3]) > 3 && protocols[3][0:3] == "@//" {
			info[HOST] = protocols[3][4:]
		} else {
			info[HOST] = protocols[3][1:]
		}

		info[PORT] = protocols[4]
	}

	info[SID] = protocols[5]

	return nil
}

func (driver *OCIDriver) Connect(url string, info map[string]string) (conn Connection,
	error DbError) {
	if error = driver.ParseUrl(url, info); error != nil {
		ErrorLogf("ParseUrl error!url:%s,%d,%s", url, error.Code(), error.Error())
		return
	}

	if error = driver.env.Init(info); error != nil {
		ErrorLogf("env init error!url:%s,%d,%s", url, error.Code(), error.Error())
		return
	}

	conn, error = driver.env.GetConnection(info)
	return conn, error
}

func init() {
	driver := OCIDriver{}
	RegisterDriver(&driver)
	DebugLog("load OCIDriver")
}
