package e

var MsgFlags = map[int]string{
	SUCCESS:        "ok",
	ERROR:          "fail",
	INVALID_PARAMS: "请求参数错误",

	ERROR_AUTH:               "Auth错误",
	ERROR_AUTH_TOKEN:         "Token错误",
	ERROR_AUTH_TOKEN_FAIL:    "Token鉴权失败",
	ERROR_AUTH_TOKEN_TIMEOUT: "Token已超时",
	ERROR_AUTH_CREATE_FAIL:   "用户注册失败",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if !ok {
		return MsgFlags[ERROR]
	}
	return msg
}
