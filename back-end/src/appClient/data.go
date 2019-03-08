package appClient

import "math/big"

const (
	MAX_RANGE = 65536
	MIN_RANGE = 0
)

const (
	GCUID_SOLICIT = iota
	GCUID_REGISTER
	GCUID_SUBMIT
	GCUID_AGGREGATION
	GCUID_APPROVE
	GCUID_CLAIM

	GCUID_REGISTER_AND_SUBMIT // for demo
	GCUID_STOP_REGISTER_AND_SUBMIT // for demo
	GCUID_SEND_TRANSACTION
)

const (
	//stage
	SOLICIT = iota
	REGISTER
	SUBMIT
	AGGREGATION
	APPROVE
	CLAIM
)

const (
	GCUID_ETHER = 101+iota
	GCUID_CURRENT_STAGE
	GCUID_REGISTER_NUMBER
	GCUID_SUBMISSION_NUMBER
	GCUID_SOLICIT_INFO
	GCUID_AGGREGATE_RESULT
	GCUID_CLAIM_NUMBER
	GCUID_BALANCE

	GCUID_QUALIFIED_NUMBER
	GCUID_IS_QUALIFIED
)


const (
	GCUID_CAN_REGISTER_AND_SUBMIT_FOR_TASK = 800
	GCUID_CAN_CLAIM_FOR_TASK = 801
)

const (
	SUCCESS = 0
	FAIL = 1
)

const (
	// error id
	UNMARSHAL_JSON_ERROR = iota
	DATA_FORMAT_ERROR
	KEY_FORMAT_ERROR
	TRANSACTION_ERROR
	CALL_TRANSACTION_ERROR
	ENCRYPTION_ERROR
)

const (
	// error message
	MSG_UNMARSHAL_JSON_ERROR = "can not unmarshal json"
	MSG_DATA_FORMAT_ERROR = "wrong data format"
	MSG_KEY_FORMAT_ERROR = "wrong private key format"
	MSG_TRANSACTION_ERROR = "transaction revert"
	MSG_CALL_TRANSACTION_ERROR = "can not call transaction"
	MSG_ENCRYPTION_ERROR = "can not encrypt data"
)

const (
	SERVER_ERROR_CODE = 500
	CLIENT_ERROR_CODE = 400
)

type SolicitRequest struct {
	Gcuid int `json:"gcuid"`
	DataFee *big.Int `json:"dataFee"`
	ServiceFee *big.Int `json:"serviceFee"`
	ServiceProvider string `json:"serviceProvider"`
	Target *big.Int `json:"target"`
	PrivateKey string `json:"privateKey"`
	Address string `json:"address"`
}

type SolicitResponse struct {
	Response
}

type RegisterRequest struct {
	Gcuid int `json:"gcuid"`
	TaskId *big.Int `json:"taskId"`
	PrivateKey string `json:"privateKey"`
	Address string `json:"address"`
}

type RegisterResponse struct {
	Response
}

type SubmitRequest struct {
	Gcuid int `json:"gcuid"`
	TaskId *big.Int `json:"taskId"`
	Value *big.Int	`json:"value"`
	Address string `json:"address"`
	PrivateKey string `json:"privateKey"`
}

type SubmitResponse struct {
	Response
}

type AggregateResquest struct {
	Gcuid int `json:"gcuid"`
	TaskId *big.Int `json:"taskId"`
	PrivateKey string `json:"privateKey"`
	Address string `json:"address"`
}


type AggregateResponse struct {
	Response
}

type ApproveRequest struct {
	Gcuid int `json:"gcuid"`
	TaskId *big.Int `json:"taskId"`
	PrivateKey string `json:"privateKey"`
	Address string `json:"address"`
}

type ApproveResponse struct {
	Response
	SubmitValues []*big.Int `json:"submitValues"`
}

type ClaimRequest struct {
	Gcuid int `json:"gcuid"`
	TaskId *big.Int `json:"taskId"`
	PrivateKey string `json:"privateKey"`
	Address string `json:"address"`
}

type ClaimResponse struct {
	Response
}

type SendTransactionRequest struct {
	Gcuid int `json:"gcuid"`
	Txid int `json:"txid"`
	RawTransaction string `json:"rawTransaction"`
}

type SendTransactionResponse struct {
	Response
	Txid int `json:"txid"`
}


// for demo
type RegisterAndSubmitRequest struct {
	Gcuid int `json:"gcuid"`
	TaskId *big.Int `json:"taskId"`
	Value *big.Int	`json:"value"`
	Address string `json:"address"`
	PrivateKey string `json:"privateKey"`
}

// for demo
type RegisterAndSubmitResponse struct {
	Response
}

type stopRegisterAndSubmitRequest struct {
	Gcuid int `json:"gcuid"`
	TaskId *big.Int `json:"taskId"`
	PrivateKey string `json:"privateKey"`
	Address string `json:"address"`
}

type stopRegisterAndSubmitResponse struct {
	Response
}

type GetEtherRequest struct {
	Gcuid int `json:"gcuid"`
	Address string `json:"address"`
}

type GetBalanceRequest struct {
	Gcuid int `json:"gcuid"`
	Address string `json:"address"`
}

type GetBalanceResponse struct {
	Response
	Amount *big.Int `json:"amount"`
	Address string `json:"address"`
}

type GetEtherResponse struct {
	Response
	Amount *big.Int `json:"amount"`
}

type GetStageRequest struct {
	Gcuid int `json:"gcuid"`
	TaskId *big.Int `json:"taskId"`
}

type GetStageResponse struct {
	Response
	Stage *big.Int `json:"stage"`
}

type GetRegisterNumberRequest struct {
	Gcuid int `json:"gcuid"`
	TaskId *big.Int `json:"taskId"`
}

type GetRegisterNumberResponse struct {
	Response
	Amount *big.Int `json:"amount"`
}

type GetSubmissionNumberRequest struct {
	Gcuid int `json:"gcuid"`
	TaskId *big.Int `json:"taskId"`
}

type GetSubmissionNumberResponse struct {
	Response
	Amount *big.Int `json:"amount"`
}

type GetAggregateResultRequest struct {
	Gcuid int `json:"gcuid"`
	TaskId *big.Int `json:"taskId"`
}

type GetAggregateResultResponse struct {
	Response
	Amount *big.Int `json:"amount"`
	QualifiedNumber *big.Int `json:"qualifiedNumber"`
}

type GetSolicitInfoResquest struct {
	Gcuid int `json:"gcuid"`
	TaskId *big.Int `json:"taskId"`
}

type GetSolicitInfoResponse struct {
	Response
	DataFee *big.Int `json:"dataFee"`
	ServiceFee *big.Int `json:"serviceFee"`
	ServiceProvider string `json:"serviceProvider"`
	Target *big.Int `json:"target"`
}

type GetClaimNumberResquest struct {
	Gcuid int `json:"gcuid"`
	TaskId *big.Int `json:"taskId"`
}

type GetClaimNumberResponse struct {
	Response
	Amount *big.Int `json:"amount"`
}

type GetQualifiedNumberRequest struct {
	Gcuid int `json:"gcuid"`
	TaskId *big.Int `json:"taskId"`
}

type GetQualfiedNumberResponse struct {
	Response
	Amount *big.Int `json:"amount"`
}

type IsQualifiedRequest struct {
	Gcuid int `json:"gcuid"`
	TaskId *big.Int `json:"taskId"`
	Address string `json:"address"`
}

type IsQualifiedResponse struct {
	Response
	Qualified bool `json:"qualified"`
}

type Response struct {
	Gcuid int `json:"gcuid"`
	Status int `json:"status"`
}

type Error struct {
	Status int `json:"status"`
	Gcuid int `json:"gcuid"`
	Code int `json:"code"`
	Reason string `json:"reason"`
}

type SubmitPayload struct {
	SubmitData string `json:"submitData"`
	SubmitProof string `json:"submitProof"`
}

type CanRegisterAndSubmitRequest struct {
	Gcuid int `json:"gcuid"`
	TaskId *big.Int `json:"taskId"`
	Address string `json:"address"`
}

type CanRegisterAndSubmitResponse struct {
	Response
	CanRegister bool `json:"canRegister"`
}

type CanClaimRequest struct {
	Gcuid int `json:"gcuid"`
	TaskId *big.Int `json:"taskId"`
	Address string `json:"address"`
}

type CanClaimResponse struct {
	Response
	CanClaim bool `json:"canClaim"`
}

type SendTransactionErrorResponse struct {
	Error
	Txid int `json:"txid"`
}

