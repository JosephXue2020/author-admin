package e

var MsgFlags = map[int]string{
	SUCCESS:        "ok",
	ERROR:          "fail",
	INVALID_PARAMS: "请求参数错误",

	ERROR_TOKEN:         "Token错误",
	ERROR_TOKEN_FAIL:    "创建Token失败",
	ERROR_TOKEN_ILLEGAL: "Token不合法",
	ERROR_TOKEN_EXPIRED: "Token过期",

	ERROR_USER:                "用户错误",
	ERROR_USER_INVALID:        "用户不合法",
	ERROR_USER_CREATE_FAIL:    "创建用户失败",
	ERROR_USER_LACK_AUTHORITY: "用户权限不足",
	ERROR_USER_ALREADY_EXIST:  "用户已存在",
	ERROR_USER_NOT_EXIST:      "用户不存在",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if !ok {
		return MsgFlags[ERROR]
	}
	return msg
}
