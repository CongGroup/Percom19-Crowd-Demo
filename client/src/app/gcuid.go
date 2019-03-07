package app

import "math/big"

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
	SOLICIT = 0
	REGISTER = 1
	SUBMIT = 2
	AGGREGATE = 3
	APPROVE = 4
	CLAIM = 5
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
	SUCCESS = 0
	FAIL = 1
)

type GetStageResponse struct {
	Response
	Stage *big.Int `json:"stage"`
}

type Response struct {
	Gcuid int `json:"gcuid"`
	Status int `json:"status"`
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

type ClaimRequest struct {
	Gcuid int `json:"gcuid"`
	TaskId *big.Int `json:"taskId"`
	PrivateKey string `json:"privateKey"`
	Address string `json:"address"`
}

type ClaimResponse struct {
	Response
}

type Error struct {
	Status int `json:"status"`
	Gcuid int `json:"gcuid"`
	Code int `json:"code"`
	Reason string `json:"reason"`
}
