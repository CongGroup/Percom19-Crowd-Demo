package contract

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

var (
	LogStageTransfer string
	LogSolicit string
	LogRegister string
	LogSubmit string
	LogAggregate string
	LogApprove string
	LogClaim string
)

type SolicitInfo struct {
	DataFee         *big.Int       `json:"dataFee"`
	ServiceFee      *big.Int       `json:"serviceFee"`
	ServiceProvider common.Address `json:"serviceProvider"`
	Target          *big.Int       `json:"target"`
}

const (
	//event list
	EVENT_STAGE_TRANSFER = "StageTransfer"
	EVENT_SOLICIT = "Solicit"
	EVENT_REGISTER = "Register"
	EVENT_SUBMIT = "Submit"
	EVENT_AGGREGATE = "Aggregate"
	EVENT_APPROVE = "Approve"
	EVENT_CLAIM = "Claim"
)

const (
	//function name
	FUNCTION_GET_STAGE_OF_TASK = "getStageOfTask"
	FUNCTION_GET_SUBMIT_DATA_OF_TASK = "getSubmitDataOfTask"
	FUNCTION_GET_SUBMIT_COUNT_OF_TASK = "getSubmitCountOfTask"
	FUNCTION_GET_REGISTER_NUMBER_OF_TASK = "getRegisterNumberOfTask"
	FUNCTION_GET_CLAIM_NUMBER_OF_TASK = "getClaimNumberOfTask"
	FUNCTION_GET_QUALIFIED_NUMBER_OF_TASK = "getQualifiedNumberOfTask"
	FUNCTION_IS_QUALIFIED_PROVIDER_FOR_TASK = "isQualifiedProviderForTask"
	FUNCTION_GET_SUBMIT_PROOF_OF_TASK = "getSubmitProofOfTask"

	FUNCTION_SOLICIT = "solicit"
	FUNCTION_REGISTER = "register"
	FUNCTION_SUBMIT = "submit"
	FUNCTION_AGGREGATE = "aggregate"
	FUNCTION_APPROVE = "approve"
	FUNCTION_CLAIM = "claim"

	FUNCTION_REGISTER_AND_SUBMIT = "registerAndSubmit"
	FUNCTION_STOP_REGISTER_AND_SUBMIT = "stopRegisterAndSubmit"

	FUNCTION_GET_SOLICITINFO_OF_TASK = "getSolicitInfoOfTask"
	FUNCTION_GET_AGGREGATION_RESULT_OF_TASK = "getAggregationResultOfTask"

	FUNCTION_IS_REGISTERED_OF_TASK = "isRegisteredOfTask"
	FUNCTION_IS_ClAIMMED_OF_TASK = "isClaimmedOfTask"
	FUNCTION_IS_SERVICE_PROVIDER_OF_TASK = "isServiceProviderOfTask"
)

type StageTransferEvent struct {
	TaskId *big.Int `json:"taskId"`
	NewStage *big.Int `json:"newStage"`
	OldStage *big.Int `json:"oldStage"`
}

type SolicitEvent struct {
	TaskId *big.Int `json:"taskId"`
}

type RegisterEvent struct {
	TaskId *big.Int `json:"taskId"`
	RegisterNumber *big.Int `json:"registerNumber"`
}

type SubmitEvent struct {
	TaskId *big.Int `json:"taskId"`
	SubmitNumber *big.Int `json:"submitNumber"`
}

type AggregateEvent struct {
	TaskId *big.Int `json:"taskId"`
	AggregateResult []byte `json:"aggregateResult"`
}

type ApproveEvent struct {
	TaskId *big.Int `json:"taskId"`
}

type ClaimEvent struct {
	TaskId *big.Int `json:"taskId"`
	ClaimNumber *big.Int `json:"claimNumber"`
}

func init() {
	stageTransferSig:= []byte("StageTransfer(uint256,uint256,uint256)")
	solicitSig:= []byte("Solicit(uint256)")
	registerSig:= []byte("Register(uint256,uint256)")
	submitSig:= []byte("Submit(uint256,uint256)")
	aggregateSig:= []byte("Aggregate(uint256,bytes)")
	approveSig:= []byte("Approve(uint256")
	claimSig:= []byte("Claim(uint256,uint256)")

	LogStageTransfer = crypto.Keccak256Hash(stageTransferSig).Hex()
	LogSolicit = crypto.Keccak256Hash(solicitSig).Hex()
	LogRegister = crypto.Keccak256Hash(registerSig).Hex()
	LogSubmit = crypto.Keccak256Hash(submitSig).Hex()
	LogAggregate = crypto.Keccak256Hash(aggregateSig).Hex()
	LogApprove = crypto.Keccak256Hash(approveSig).Hex()
	LogClaim = crypto.Keccak256Hash(claimSig).Hex()
}

type Agg struct {
	BaseContract
}

func NewAgg(port string, contractABI string ,contractAddress common.Address) *Agg {
	agg := new(Agg)
	agg.Connect(port)
	agg.Address = contractAddress
	agg.LoadABI(contractABI)
	return agg
}

