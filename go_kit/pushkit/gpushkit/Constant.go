package gpushkit

//push消息地址
const (
	CONSTANT_PUSH_URL = "https://restapi.getui.com/v1/"
)

//返回结果，如果result为sign_error或者not_auth,就应该重新鉴权
const (
	CONSTANT_OK                = "ok"
	CONSTANT_NO_MSG            = "no_msg"
	CONSTANT_ALIAS_ERROR       = "alias_error"
	CONSTANT_BLACK_IP          = "black_ip"
	CONSTANT_SIGN_ERROR        = "sign_error"
	CONSTANT_PUSHNUM_OVERLIMIT = "pushnum_overlimit"
	CONSTANT_NOT_AUTH          = "not_auth"
)

const (
	GTPUSHTYPE = 11
)
