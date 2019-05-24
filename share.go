// Copyright (C) 2019 tisnyo <tisnyo@gmail.com>.
//
// This file contains some common functions.
package gox

import (
	"container/list"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"io"
	"strconv"
	"strings"
)

// ConvertBoolFromInt converts int to bool.
func ConvertBoolFromInt(input int) bool {
	if input <= 0 {
		return false
	}
	return true
}

// List2Array converts list to array.
func List2Array(ls *list.List) []interface{} {
	if ls == nil {
		return nil
	}
	arr := make([]interface{}, ls.Len())
	index := 0
	for ele := ls.Front(); ele != nil; ele = ele.Next() {
		arr[index] = ele.Value
		index++
	}
	return arr
}

// ParseHostPortFromConnStr parses host and port from connection string.
func ParseHostPortFromConnStr(connStr string) (string, int) {
	host := strings.Split(connStr, ":")[0]
	port, _ := strconv.Atoi(strings.Split(connStr, ":")[1])
	return host, port
}

// TOperation simulates ternary operation.
func TOperation(condition bool, trueOperation func() interface{}, falseOperation func() interface{}) interface{} {
	if condition {
		if trueOperation == nil {
			return nil
		}
		return trueOperation()
	}
	if falseOperation == nil {
		return nil
	}
	return falseOperation()
}

// TValue ternary operation
func TValue(condition bool, trueValue interface{}, falseValue interface{}) interface{} {
	if condition {
		return trueValue
	}
	return falseValue
}

// WalkList walk a list.
// walker return value as break signal,
// if it is true, break walking
func WalkList(ls *list.List, walker func(item interface{}) bool) {
	if ls == nil {
		return
	}
	for ele := ls.Front(); ele != nil; ele = ele.Next() {
		breakWalk := walker(ele.Value)
		if breakWalk {
			break
		}
	}
}

// ConvertLength2Bytes converts an int64 value to a byte array.
func ConvertLength2Bytes(len int64, buffer *[]byte) *[]byte {
	binary.BigEndian.PutUint64(*buffer, uint64(len))
	return buffer
}

// ConvertBytes2Length converts a byte array to an int64 value.
func ConvertBytes2Length(ret *[]byte) int64 {
	return int64(binary.BigEndian.Uint64(*ret))
}

// Md5Sum calculates md5 value of some strings.
func Md5Sum(input ...string) string {
	h := md5.New()
	if input != nil {
		for _, v := range input {
			io.WriteString(h, v)
		}
	}
	sliceCipherStr := h.Sum(nil)
	sMd5 := hex.EncodeToString(sliceCipherStr)
	return sMd5
}

// ControlRange controls a variable's value in a value range.
// returns defaultValue if the range not provided or the value is not in the range.
func ControlRange(value interface{}, defaultValue interface{}, rangeValues ...interface{}) interface{} {
	if rangeValues == nil || len(rangeValues) == 0 {
		return defaultValue
	}
	in := false
	for _, v := range rangeValues {
		if v == value {
			in = true
			break
		}
	}
	if in {
		return value
	}
	return defaultValue
}
