package constants

import (
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
)

const (
	// 1400xxxx: GRPC

	GRPCInternalServerErrCode               = 14000000
	GRPCInternalServerErrMsg                = "internal server error"
	GRPCAllOfConnectionsCanNotBeUsedErrCode = 14000001
	GRPCAllOfConnectionsCanNotBeUsedErrMsg  = "all of connections can not be used"
)

var (
	// 1400xxxx: GRPC

	// 14000000(internal server error)
	GRPCInternalServerErr = tkErr.Error(GRPCInternalServerErrCode, GRPCInternalServerErrMsg)
	// 14000001(all of connections can not be used)
	GRPCAllOfConnectionsCanNotBeUsedErr = tkErr.Error(GRPCAllOfConnectionsCanNotBeUsedErrCode, GRPCAllOfConnectionsCanNotBeUsedErrMsg)
)
