package utils

import (
	"github.com/mgolfam/gogutils/enums"
)

func IsPasswordEqual(dbPass, inputPass string, passType int) bool {
	if enums.PASS_TYPE_PLAIN == passType {
		if dbPass == inputPass {
			return true
		}
	} else if enums.PASS_TYPE_MD5 == passType {
		if dbPass == HashMd5(inputPass) {
			return true
		}
	} else if enums.PASS_TYPE_SHA256 == passType {
		if dbPass == HashSha256(inputPass) {
			return true
		}
	}

	return false
}
