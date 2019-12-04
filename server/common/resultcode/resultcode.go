package resultcode

const (
	Success                  = "0"
	SystemInternalException  = "102"
	UserAuthFail             = "103"
	RequestIllegal           = "104"
	CustomerNotExist         = "10000"
	UserNameOrPasswordError  = "10001"
	CustomerNameAlreadyExist = "10002"
	EmailAlreadyExist        = "10003"
	UserStatusError          = "10004"
)

var ResultMessage = map[string]string{
	Success:                  "success",
	SystemInternalException:  "system internal exception",
	UserAuthFail:             "user auth fail",
	RequestIllegal:           "request illegal, error: %s",
	CustomerNotExist:         "customer not exist",
	UserNameOrPasswordError:  "userName or password error",
	CustomerNameAlreadyExist: "customerName %s already exist",
	EmailAlreadyExist:        "email %s already exist",
	UserStatusError:          "user status is not normal",
}
