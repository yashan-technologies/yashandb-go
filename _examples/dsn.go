package example

import (
	"encoding/json"
	"flag"
	"io/ioutil"
)

var dsnFile = "example/defaultDsn.json"

type dsnReadStruct struct {
	DefaultDsn string `json:"defaultDsn"`
}

func GetDsn() string {
	defaultDsn := getDefaultDsn()
	if defaultDsn == "" {
		defaultDsn = "regress/regress@127.0.0.1:1688"
	}
	dsn := flag.String("dsn", defaultDsn, "input you dsn(DataSourceName, format:username/password@host:port[?param1=value1&...&paramN=valueN]) to connect yashandb.")
	flag.Parse()
	return *dsn
}

func getDefaultDsn() string {
	data, _ := ioutil.ReadFile(dsnFile)
	var dsn dsnReadStruct
	json.Unmarshal(data, &dsn)
	return dsn.DefaultDsn
}
