/*
Copyright  2022, YashanDB and/or its affiliates. All rights reserved.
YashanDB Driver for golang is licensed under the terms of the mulan PSL v2.0

License: 	http://license.coscl.org.cn/MulanPSL2
Home page: 	https://www.yashandb.com/
*/

package yasdb

/*
#cgo CFLAGS: -I./yacapi/include -I./yacapi/src

#include "yacapi.h"
#include "yapi_inc.h"
#include <stdio.h>
#include <stdlib.h>
*/
import "C"

import (
	"context"
	"database/sql/driver"
	"fmt"
	"strconv"
	"sync"
	"time"
	"unsafe"
)

type serverStatus uint8

const (
	_DefaultNcharsetRatio = 4

	SS_UNKNOWN     serverStatus = 255
	SS_UNCONNECTED serverStatus = 0
	SS_NORMAL      serverStatus = 1
)

type YasConn struct {
	Env               *C.YapiEnv
	Conn              *C.YapiConnect
	closed            bool
	charsetRatio      uint32 // 最大CHARSET膨胀比率
	ncharsetRatio     uint32 // 最大NCHARSET膨胀比率
	numberAsString    bool   // YashanDB的number类型返回为golang的string类型，默认返回float64类型
	cliPrepare        bool
	autocommit        bool
	timestampFormat   string
	timestampTzFormat string
	dateFormat        string
	timeFormat        string
	dsIntervalFormat  string
	ymIntervalFormat  string
	mu                sync.Mutex
}

func (conn *YasConn) Prepare(query string) (driver.Stmt, error) {
	return PrepareContext(conn, context.Background(), query)
}

func (conn *YasConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	return PrepareContext(conn, ctx, query)
}

func (conn *YasConn) Begin() (driver.Tx, error) {
	return conn.BeginTx(context.Background(), driver.TxOptions{})
}

func (conn *YasConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	_ = conn.setAutoCommit(false)
	return &YasTx{Conn: conn}, nil
}

func (conn *YasConn) Close() error {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	if conn.closed {
		return nil
	}

	conn.closed = true
	connErr := releaseConn(conn.Conn)
	envErr := releaseEnv(conn.Env)
	if envErr != nil && connErr != nil {
		return fmt.Errorf("release env err: %s ; release conn err: %s", envErr, connErr)
	} else if envErr != nil {
		return envErr
	} else if connErr != nil {
		return connErr
	}
	conn.Conn = nil
	conn.Env = nil
	return nil
}

func (conn *YasConn) Ping(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if conn.Conn == nil {
		return ErrNoConnect()
	}
	if err := yapiPing(conn.Conn, -1); err != nil {
		// c driver is unsupport yacPingWithTimeout, support in 23.4.4
		if isLoadSymbolErr(err) {
			return nil
		}
		return err
	}
	return nil
}

func (conn *YasConn) getConnAttr() error {
	if err := conn.getCharsetRatio(); err != nil {
		return err
	}

	return conn.getNcharsetRatio()
}

func (conn *YasConn) getCharsetRatio() error {
	var ratio C.uint32_t
	size := C.int32_t(unsafe.Sizeof(ratio))
	if err := conn.yapiGetConnAttr(C.YAPI_ATTR_MAX_CHARSET_RATIO, unsafe.Pointer(&ratio), size); err != nil {
		return err
	}
	conn.charsetRatio = uint32(ratio)
	return nil
}

func (conn *YasConn) getNcharsetRatio() error {
	var (
		ratio     C.uint32_t
		stringLen C.int32_t
	)
	size := C.int32_t(unsafe.Sizeof(ratio))
	if existYasError(C.yapiGetConnAttr(conn.Conn, C.YAPI_ATTR_MAX_NCHARSET_RATIO, unsafe.Pointer(&ratio), size, &stringLen)) {
		conn.ncharsetRatio = _DefaultNcharsetRatio
	} else {
		conn.ncharsetRatio = uint32(ratio)
	}
	return nil
}

func (conn *YasConn) setAutoCommit(auto bool) error {
	var a C.int32_t = 0
	if auto {
		a = 1
	}
	size := C.int32_t(unsafe.Sizeof(a))
	if err := conn.yapiSetConnAttr(C.YAPI_ATTR_AUTOCOMMIT, unsafe.Pointer(&a), size); err != nil {
		return err
	}
	return nil
}

func (conn *YasConn) setHeartbeatEnable(enable bool) error {
	if !enable {
		return nil
	}
	var a C.bool = true
	size := C.int32_t(unsafe.Sizeof(a))
	if err := conn.yapiSetConnAttr(C.YAPI_ATTR_HEARTBEAT_ENABLED, unsafe.Pointer(&a), size); err != nil {
		if !isUnknownAttributeIdErr(err) {
			return err
		}
	}
	return nil
}

func (conn *YasConn) getServerStatus() serverStatus {
	var status C.uint32_t
	size := C.int32_t(unsafe.Sizeof(status))
	err := conn.yapiGetConnAttr(C.YAPI_ATTR_SERVER_STATUS, unsafe.Pointer(&status), size)
	if err != nil {
		if isUnknownAttributeIdErr(err) {
			return SS_UNKNOWN
		}
		return SS_UNCONNECTED
	}
	if status == 1 {
		return SS_NORMAL
	}
	return SS_UNCONNECTED
}

func (conn *YasConn) setCompatVector(compatVector string) error {
	if compatVector == "" || compatVector == "null" {
		return nil
	}

	stmt, err := PrepareContext(conn, context.Background(), (fmt.Sprintf("alter session set compat_vector=%s", compatVector)))
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.yacExecute()
}

func (conn *YasConn) yapiSetConnAttr(attr C.YapiConnAttr, value unsafe.Pointer, bufLength C.int32_t) error {
	return yapiSetConnAttr(conn.Conn, attr, value, bufLength)
}

func (conn *YasConn) yapiGetConnAttr(attr C.YapiConnAttr, value unsafe.Pointer, bufLength C.int32_t) error {
	return yapiGetConnAttr(conn.Conn, attr, value, bufLength)
}

func (conn *YasConn) yacCommit() error {
	defer conn.setAutoCommit(conn.autocommit)
	return yapiCommit(conn.Conn)
}

func (conn *YasConn) yacRollback() error {
	defer conn.setAutoCommit(conn.autocommit)
	return yapiRollback(conn.Conn)
}

func (conn *YasConn) lobRead(lobLocator *C.YapiLobLocator) ([]byte, error) {
	lobLen, err := conn.yacLobGetLength(lobLocator)
	if err != nil {
		return nil, err
	}
	data, err := conn.yacLobRead(lobLocator, lobLen)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (conn *YasConn) yacLobRead(lobLocator *C.YapiLobLocator, lobLen uint64) ([]byte, error) {
	if lobLen == 0 {
		return []byte{}, nil
	}
	data := make([]byte, 0, lobLen)
	bytes := C.uint64_t(_LobBufLen)
	for {
		readBuffer := byteBufferPool.Get().([]byte)
		buf := (*C.uint8_t)((unsafe.Pointer)(&readBuffer[0]))
		if err := yapiLobRead(conn.Conn, lobLocator, &bytes, buf); err != nil {
			return nil, err
		}
		data = append(data, readBuffer[:uint64(bytes)]...)
		if uint64(bytes) <= 0 {
			break
		}
	}
	return data, nil
}

func (conn *YasConn) yacLobGetLength(lobLocator *C.YapiLobLocator) (uint64, error) {
	var lobLen C.uint64_t
	if err := yapiLobGetLength(conn.Conn, lobLocator, &lobLen); err != nil {
		return 0, err
	}
	return uint64(lobLen), nil
}

func (conn *YasConn) lobWrite(yacType C.YapiType, data []byte) (*unsafe.Pointer, error) {
	desc, err := conn.yacLobDescAlloc(yacType)
	if err != nil {
		return nil, err
	}
	lobLocator := (**C.YapiLobLocator)(unsafe.Pointer(desc))
	if err := conn.yacLobCreateTemporary(*lobLocator); err != nil {
		C.yapiLobDescFree(unsafe.Pointer(*lobLocator), yacType)
		return nil, err
	}
	if err := conn.yacLobWrite(*lobLocator, data); err != nil {
		C.yapiLobFreeTemporary(conn.Conn, *lobLocator)
		C.yapiLobDescFree(unsafe.Pointer(*lobLocator), yacType)
		return nil, err
	}
	return desc, nil
}

func (conn *YasConn) yacLobDescAlloc(yacType C.YapiType) (*unsafe.Pointer, error) {
	desc := new(unsafe.Pointer)
	if err := yapiLobDescAlloc(conn.Conn, yacType, desc); err != nil {
		return nil, err
	}
	return desc, nil
}

func (conn *YasConn) yacLobCreateTemporary(lobLocator *C.YapiLobLocator) error {
	if err := yapiLobCreateTemporary(conn.Conn, lobLocator); err != nil {
		return err
	}
	return nil
}

func (conn *YasConn) yacLobWrite(lobLocator *C.YapiLobLocator, data []byte) error {
	if len(data) == 0 || data == nil {
		return nil
	}
	bufLen := uint64(_LobBufLen)
	dataLen := uint64(len(data))
	writeBuffer := byteBufferPool.Get().([]byte)
	if _LobBufLen > dataLen {
		bufLen = dataLen
		copy(writeBuffer, data)
	} else {
		copy(writeBuffer, data[0:_LobBufLen])
	}
	buf := (*C.uint8_t)((unsafe.Pointer)(&writeBuffer[0]))
	count := uint64(0)
	for {
		if err := yapiLobWrite(conn.Conn, lobLocator, buf, C.uint64_t(bufLen)); err != nil {
			return nil
		}
		count += bufLen
		if count >= dataLen {
			break
		}
		if count+bufLen < dataLen {
			copy(writeBuffer, data[count:count+bufLen])
		} else {
			copy(writeBuffer, data[count:])
			bufLen = dataLen - count
		}
	}
	return nil
}

func (conn *YasConn) lobFree(yacType C.YapiType, lobLocator *C.YapiLobLocator) {
	if yacType != C.YAPI_TYPE_BLOB && yacType != C.YAPI_TYPE_CLOB {
		return
	}
	C.yapiLobFreeTemporary(conn.Conn, lobLocator)
	C.yapiLobDescFree(unsafe.Pointer(lobLocator), yacType)
}

func (conn *YasConn) handleYacCancel(ctx context.Context, done <-chan struct{}) {
	select {
	case <-done:
	case <-ctx.Done():
		select {
		case <-done:
		default:
			_ = conn.yacCancel()
		}
	}
}

func (conn *YasConn) yacCancel() error {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	if conn.closed {
		return nil
	}
	return yapiCancel(conn.Conn)
}

func (conn *YasConn) ResetSession(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if conn == nil || conn.closed {
		return driver.ErrBadConn
	}
	status := conn.getServerStatus()
	if status == SS_NORMAL || status == SS_UNKNOWN {
		return nil
	}
	return driver.ErrBadConn
}

func (conn *YasConn) yapiTimestampToTime(dateTime *C.YapiTimestamp, zone bool) (*time.Time, error) {
	// get date
	var (
		year     C.int16_t
		month    C.uint8_t
		day      C.uint8_t
		hour     C.uint8_t
		minute   C.uint8_t
		second   C.uint8_t
		fraction C.uint32_t
	)

	if err := yapiTimestampGetTimestamp(
		*dateTime,
		&year,
		&month,
		&day,
		&hour,
		&minute,
		&second,
		&fraction,
	); err != nil {
		return nil, err
	}

	if !zone {
		aTime := time.Date(int(year), time.Month(month), int(day), int(hour), int(minute), int(second), int(fraction), time.UTC)
		return &aTime, nil
	}

	var timeZoneHour C.int8_t
	var timeZoneMin C.int8_t

	if err := yapiDateTimeGetTimeZoneOffset(conn.Env, *dateTime, &timeZoneHour, &timeZoneMin); err != nil {
		return nil, err
	}

	aTime := time.Date(int(year), time.Month(month), int(day), int(hour), int(minute), int(second), int(fraction),
		timezoneToLocation(int64(timeZoneHour), int64(timeZoneMin)))

	return &aTime, nil
}

func (conn *YasConn) timeToYapiTimestamp(dest *time.Time) (*C.YapiTimestamp, error) {
	var timestamp C.YapiTimestamp
	p := C.malloc(C.size_t(unsafe.Sizeof(timestamp)))

	tpointer := (*C.YapiTimestamp)(p)
	year := C.int16_t(dest.Year())
	month := C.uint8_t(dest.Month())
	day := C.uint8_t(dest.Day())
	hour := C.uint8_t(dest.Hour())
	mintue := C.uint8_t(dest.Minute())
	second := C.uint8_t(dest.Second())
	fraction := C.uint32_t(dest.Nanosecond())

	if err := yapiTimestampSetTimestamp(tpointer, year, month, day, hour, mintue, second, fraction); err != nil {
		C.free(p)
		return nil, err
	}
	return tpointer, nil

}

func (conn *YasConn) stringToYapiDSInterval(dest *string) (*C.YapiDSInterval, error) {
	var dsInterval C.YapiDSInterval
	p := C.malloc(C.size_t(unsafe.Sizeof(dsInterval)))
	dsPointer := (*C.YapiDSInterval)(p)
	cstr := C.CString(*dest)
	defer C.free(unsafe.Pointer(cstr))
	length := C.uint32_t(len((*dest)))

	if err := yapiDSIntervalFromText(conn.Env, dsPointer, cstr, length); err != nil {
		C.free(p)
		return nil, err
	}
	return dsPointer, nil
}

func (conn *YasConn) yapiDSIntervalToString(interval *C.YapiDSInterval) (string, error) {

	var (
		day      C.int32_t
		hour     C.int32_t
		mintue   C.int32_t
		second   C.int32_t
		fraction C.int32_t
	)
	if err := yapiDSIntervalGetDaySecond(
		*interval,
		&day,
		&hour,
		&mintue,
		&second,
		&fraction,
	); err != nil {
		return "", err
	}

	t := time.Date(0, 0, int(day), int(hour), int(mintue), int(second), int(fraction), time.UTC)
	return FormatTime(conn.dsIntervalFormat, t), nil
}

func (conn *YasConn) yapiYMIntervalToString(interval *C.YapiYMInterval) (string, error) {

	var (
		year  C.int32_t
		month C.int32_t
	)
	if err := yapiYMIntervalGetYearMonth(
		*interval,
		&year,
		&month,
	); err != nil {
		return "", err
	}

	return FormatYMInterval(conn.ymIntervalFormat, int32(year), int32(month)), nil
}

func (conn *YasConn) stringToYapiYMInterval(dest *string) (*C.YapiYMInterval, error) {

	var dsInterval C.YapiYMInterval
	p := C.malloc(C.size_t(unsafe.Sizeof(dsInterval)))
	dsPointer := (*C.YapiYMInterval)(p)
	cstr := C.CString(*dest)
	defer C.free(unsafe.Pointer(cstr))
	length := C.uint32_t(len((*dest)))
	if len(*dest) == 0 {
		return dsPointer, nil
	}

	if err := yapiYMIntervalFromText(conn.Env, dsPointer, cstr, length); err != nil {
		C.free(p)
		return nil, err
	}
	return dsPointer, nil
}

func (conn *YasConn) float64ToYapiNumber(dest *float64) (*C.YapiNumber, error) {

	var number C.YapiNumber
	p := C.malloc(C.size_t(unsafe.Sizeof(number)))
	np := (*C.YapiNumber)(p)

	yp := C.YapiPointer(unsafe.Pointer(dest))
	length := C.uint32_t(unsafe.Sizeof(*dest))

	if err := yapiNumberFromReal(yp, length, np); err != nil {
		C.free(p)
		return nil, err
	}
	return np, nil
}

func (conn *YasConn) yapiNumberToFloat64(number *C.YapiNumber) (float64, error) {

	var res float64

	length := C.uint32_t(unsafe.Sizeof(res))

	if err := yapiNumberToReal(number, length, C.YapiPointer(unsafe.Pointer(&res))); err != nil {
		return 0, err
	}
	return res, nil
}

func PrepareContext(conn *YasConn, ctx context.Context, query string) (*YasStmt, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var stmt *C.YapiStmt
	nQuery, cst := tryRmSemicolon(query)
	if conn.cliPrepare && (cst == CST_INSERT || cst == CST_DELETE || cst == CST_UPDATE || cst == CST_SELECT) {
		if err := yapiStmtCreate(conn.Conn, &stmt); err != nil {
			return nil, err
		}
		// insert 直接走 DirectExecute
		return &YasStmt{
			Conn:     conn,
			Stmt:     stmt,
			SqlType:  (uint32)(cst),
			Sqlstr:   nQuery,
			prepared: false,
		}, nil
	}
	queryP := C.CString(nQuery)
	defer C.free(unsafe.Pointer(queryP))
	sqlLength := C.int32_t(len(nQuery))
	if err := yapiPrepare(conn.Conn, queryP, sqlLength, &stmt); err != nil {
		return nil, err
	}

	var sqlType C.uint32_t
	sqlSize := C.int32_t(unsafe.Sizeof(sqlType))
	if err := yapiGetStmtAttr(
		stmt,
		C.YAPI_ATTR_SQLTYPE,
		unsafe.Pointer(&sqlType),
		sqlSize,
		sqlLength); err != nil {
		return nil, err
	}

	yasStmt := &YasStmt{
		Conn:     conn,
		Stmt:     stmt,
		SqlType:  (uint32)(sqlType),
		Sqlstr:   nQuery,
		prepared: true,
	}

	return yasStmt, nil
}

func timezoneToLocation(hour int64, minute int64) *time.Location {
	if minute != 0 || hour > 14 || hour < -12 {
		// create location with FixedZone
		var name string
		if hour < 0 {
			name = strconv.FormatInt(hour, 10) + ":"
		} else {
			name = "+" + strconv.FormatInt(hour, 10) + ":"
		}
		if minute == 0 {
			name += "00"
		} else {
			if minute < 10 {
				name += "0"
			}
			name += strconv.FormatInt(minute, 10)
		}
		return time.FixedZone(name, (3600*int(hour))+(60*int(minute)))
	}

	// use location from timeLocations cache
	return timeLocations[12+hour]
}
