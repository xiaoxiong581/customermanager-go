package resultcode

const (
    Success                 = "0"
    SystemInternalException = "102"
    UserAuthFail            = "103"
    RequestIllegal          = "104"
    CustomerNotExist        = "10000"
    UserNameOrPasswordError = "10001"
)

var ResultMessage = map[string]string{
    Success:                 "success",
    SystemInternalException: "system internal exception",
    UserAuthFail:            "user auth fail",
    RequestIllegal:          "request illegal",
    CustomerNotExist:        "customer not exist",
    UserNameOrPasswordError: "userName or password error",
}
