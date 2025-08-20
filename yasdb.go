/*
Copyright  2022, YashanDB and/or its affiliates. All rights reserved.
YashanDB Driver for golang is licensed under the terms of the mulan PSL v2.0

License: 	http://license.coscl.org.cn/MulanPSL2
Home page: 	https://www.yashandb.com/
*/

package yasdb

import (
	"database/sql"
	"strconv"
	"time"
)

var (
	timeLocations []*time.Location
)

func init() {
	sql.Register("yasdb", &YasdbDriver{})

	timeLocationNames := []string{"Etc/GMT+12", "Pacific/Pago_Pago", // -12 to -11
		"Pacific/Honolulu", "Pacific/Gambier", "Pacific/Pitcairn", "America/Phoenix", "America/Costa_Rica", // -10 to -6
		"America/Panama", "America/Puerto_Rico", "America/Punta_Arenas", "America/Noronha", "Atlantic/Cape_Verde", // -5 to -1
		"GMT",                                                                         // 0
		"Africa/Lagos", "Africa/Cairo", "Europe/Moscow", "Asia/Dubai", "Asia/Karachi", // 1 to 5
		"Asia/Dhaka", "Asia/Jakarta", "Asia/Shanghai", "Asia/Tokyo", "Australia/Brisbane", // 6 to 10
		"Pacific/Noumea", "Asia/Anadyr", "Pacific/Enderbury", "Pacific/Kiritimati", // 11 to 14
	}

	var err error
	timeLocations = make([]*time.Location, len(timeLocationNames))
	for i := 0; i < len(timeLocations); i++ {
		timeLocations[i], err = time.LoadLocation(timeLocationNames[i])
		if err != nil {
			name := "GMT"
			if i < 12 {
				name += strconv.FormatInt(int64(i-12), 10)
			} else if i > 12 {
				name += "+" + strconv.FormatInt(int64(i-12), 10)
			}
			timeLocations[i] = time.FixedZone(name, 3600*(i-12))
		}
	}
}
