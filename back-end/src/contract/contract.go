package contract

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

const (
	GETH_PORT    = "ws://localhost:8650"
	CONTRACT_ADDRESS = "0xd573ff992c921a9d3597e8c17b5267a50f257def"
	CONTRACT_ABI     = `[
	{
		"constant": false,
		"inputs": [
			{
				"name": "task_id",
				"type": "uint256"
			},
			{
				"name": "aggregation",
				"type": "bytes"
			},
			{
				"name": "qualifiedSets",
				"type": "bytes"
			},
			{
				"name": "share",
				"type": "bytes"
			},
			{
				"name": "attestatino",
				"type": "bytes"
			}
		],
		"name": "aggregate",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "task_id",
				"type": "uint256"
			}
		],
		"name": "approve",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "erc20",
				"type": "address"
			}
		],
		"name": "changeERC20Contract",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "task_id",
				"type": "uint256"
			}
		],
		"name": "claim",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "task_id",
				"type": "uint256"
			}
		],
		"name": "register",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "task_id",
				"type": "uint256"
			},
			{
				"name": "data",
				"type": "bytes"
			},
			{
				"name": "proof",
				"type": "bytes"
			}
		],
		"name": "registerAndSubmit",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "data_fee",
				"type": "uint256"
			},
			{
				"name": "service_fee",
				"type": "uint256"
			},
			{
				"name": "service_provider",
				"type": "address"
			},
			{
				"name": "target",
				"type": "uint256"
			}
		],
		"name": "solicit",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "task_id",
				"type": "uint256"
			}
		],
		"name": "stopRegisterAndSubmit",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "task_id",
				"type": "uint256"
			},
			{
				"name": "data",
				"type": "bytes"
			},
			{
				"name": "proof",
				"type": "bytes"
			}
		],
		"name": "submit",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"name": "erc20",
				"type": "address"
			}
		],
		"payable": true,
		"stateMutability": "payable",
		"type": "constructor"
	},
	{
		"payable": true,
		"stateMutability": "payable",
		"type": "fallback"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"name": "taskId",
				"type": "uint256"
			},
			{
				"indexed": false,
				"name": "newStage",
				"type": "uint256"
			},
			{
				"indexed": false,
				"name": "oldStage",
				"type": "uint256"
			}
		],
		"name": "StageTransfer",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"name": "taskId",
				"type": "uint256"
			}
		],
		"name": "Solicit",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"name": "taskId",
				"type": "uint256"
			},
			{
				"indexed": false,
				"name": "registerNumber",
				"type": "uint256"
			}
		],
		"name": "Register",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"name": "taskId",
				"type": "uint256"
			},
			{
				"indexed": false,
				"name": "submitNumber",
				"type": "uint256"
			}
		],
		"name": "Submit",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"name": "taskId",
				"type": "uint256"
			},
			{
				"indexed": false,
				"name": "aggregateResult",
				"type": "bytes"
			}
		],
		"name": "Aggregate",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"name": "taskId",
				"type": "uint256"
			}
		],
		"name": "Approve",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"name": "taskId",
				"type": "uint256"
			},
			{
				"indexed": false,
				"name": "claimNumber",
				"type": "uint256"
			}
		],
		"name": "Claim",
		"type": "event"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "task_id",
				"type": "uint256"
			}
		],
		"name": "getAggregationResultOfTask",
		"outputs": [
			{
				"name": "",
				"type": "bytes"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "task_id",
				"type": "uint256"
			}
		],
		"name": "getClaimNumberOfTask",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "task_id",
				"type": "uint256"
			}
		],
		"name": "getQualifiedNumberOfTask",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "task_id",
				"type": "uint256"
			}
		],
		"name": "getRegisterNumberOfTask",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "task_id",
				"type": "uint256"
			}
		],
		"name": "getSolicitInfoOfTask",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			},
			{
				"name": "",
				"type": "uint256"
			},
			{
				"name": "",
				"type": "address"
			},
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "task_id",
				"type": "uint256"
			}
		],
		"name": "getStageOfTask",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "task_id",
				"type": "uint256"
			}
		],
		"name": "getSubmitCountOfTask",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "task_id",
				"type": "uint256"
			},
			{
				"name": "index",
				"type": "uint256"
			}
		],
		"name": "getSubmitDataOfTask",
		"outputs": [
			{
				"name": "",
				"type": "bytes"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "task_id",
				"type": "uint256"
			},
			{
				"name": "index",
				"type": "uint256"
			}
		],
		"name": "getSubmitProofOfTask",
		"outputs": [
			{
				"name": "",
				"type": "bytes"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "task_id",
				"type": "uint256"
			},
			{
				"name": "user",
				"type": "address"
			}
		],
		"name": "isQualifiedProviderForTask",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "lastest_task",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "MAX_TASK",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "token",
		"outputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	}
]`
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

