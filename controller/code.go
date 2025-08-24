package controller


type ResCode int64

const(
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy

	CodeNeedLogin
	CodeInvalidToken

)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:      	 "success",
	CodeInvalidParam: 	 "请求参数错误",
	CodeUserExist:  	 "用户已存在",
	CodeUserNotExist: 	 "用户不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy: 	 "服务器繁忙",
	CodeInvalidToken:     "无效的Token",
	CodeNeedLogin:       "需要登录",
}

func (c ResCode) Msg() string{                  // 接收者是 ResCode 类型  相当于绑定到这个类型作成员函数
	msg, ok := codeMsgMap[c]
	if !ok {
		return "unknown"
	}
	return msg
}