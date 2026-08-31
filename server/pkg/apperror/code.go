package apperror

const (
	CodeSuccess         = 0
	CodeUserNotFound    = 10001
	CodeUserExists      = 10002
	CodePasswordWrong   = 10003
	CodeNotLogin        = 20001
	CodeNoPermission    = 20002
	CodeTokenExpired    = 20003
	CodeParamMissing    = 30001
	CodeParamFormat     = 30002
	CodeResourceNotFnd  = 40001
	CodeOrderInvalid    = 40101
	CodeOrderPaid       = 40102
	CodePlanExpired     = 40103
	CodeTrafficExhausted = 40104
	CodePaymentInvalid  = 40105
	CodeCouponInvalid   = 40201
	CodeBalanceNotEnoug = 40301
	CodeRateLimited     = 29999
	CodeDBError         = 90001
	CodeUpstreamError   = 90002
)
