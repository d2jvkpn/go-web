package api

import (
	"regexp"
)

const (
	HTTP_MaxFileHeaderSize = 8 << 20 // 8M
	KEY_User               = "User"
)

var (
	_UserStr     = "^[0-9a-zA-Z-_]{5,32}$"
	_UserRE      = regexp.MustCompile(_UserStr)
	_PasswordStr = "^[a-zA-Z0-9-_.]{8,20}$"
	_PasswordRE  = regexp.MustCompile(_PasswordStr)
	_FilenameStr = "^[0-9a-zA-Z-_][0-9a-zA-Z-_.]{0,31}$"
	_FilenameRE  = regexp.MustCompile(_FilenameStr)
	_UsersData   = NewUsersData()
)
