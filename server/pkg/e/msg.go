package e

var MsgFlags = map[int]string{
	SUCCESS:        "ok",
	ERROR:          "fail",
	INVALID_PARAMS: "请求参数错误",

	ERROR_AUTH:               "Auth错误",
	ERROR_AUTH_TOKEN_FAIL:    "Token错误",
	ERROR_AUTH_TOKEN_ILLEGAL: "Token不合法",
	ERROR_AUTH_TOKEN_EXPIRED: "Token过期",
	ERROR_AUTH_CREATE_FAIL:   "注册失败",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if !ok {
		return MsgFlags[ERROR]
	}
	return msg
}
