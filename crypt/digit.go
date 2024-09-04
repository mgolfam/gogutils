package crypt

import (
	"strconv"

	"github.com/mgolfam/gogutils/glog"
)

func EncodeBase36(number int64) string {
	return strconv.FormatInt(int64(number), 36)
}

func DecondeBase36(text string) int64 {
	decimalValue, success := strconv.ParseInt(text, 36, 0)
	if success == nil {
		glog.LogL(glog.DEBUG, "Conversion failed.")
		return -1
	}
	return decimalValue
}
