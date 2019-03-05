package appClient

import (
	"contract"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gorilla/websocket"
	"log"
	"math/big"
	"sync"
	"user"
	"zcrypto"
)

type Handler struct {
	w *websocket.Conn
	agg *contract.Agg
	token *contract.ERC20
	mutex *sync.Mutex
}

func NewHandler(w *websocket.Conn, agg *contract.Agg, token *contract.ERC20) *Handler {
	return &Handler{
		w:w,
		mutex:&sync.Mutex{},
		agg:agg,
		token:token,
	}
}

func (h *Handler) solicitHandler(gcuid int,data []byte) {
	var payload SolicitRequest
	err:= json.Unmarshal(data,&payload)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,UNMARSHAL_JSON_ERROR,err)
		return
	}

	pk,err:= derivePrivateKey(payload.PrivateKey)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,KEY_FORMAT_ERROR,err)
		return
	}

	log.Println("data fee:",payload.DataFee);
	log.Println("service fee:",payload.ServiceFee);
	u:= user.NewUser(common.HexToAddress(payload.Address),pk,h.agg)
	err =u.Send(contract.FUNCTION_SOLICIT,payload.DataFee,payload.ServiceFee,common.HexToAddress(payload.ServiceProvider),payload.Target)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,TRANSACTION_ERROR,err)
		return
	}
	res:= &SolicitResponse{
		Response: Response {
			Gcuid:gcuid,
			Status:SUCCESS,
		},
	}

	h.wrapperAndSend(gcuid, res)
}

func (h *Handler) registerHandler(gcuid int ,data []byte) {
	var payload RegisterRequest
	err:= json.Unmarshal(data,&payload)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,UNMARSHAL_JSON_ERROR,err)
		return
	}

	pk,err:= derivePrivateKey(payload.PrivateKey)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,KEY_FORMAT_ERROR,err)
		return
	}
	u:= user.NewUser(common.HexToAddress(payload.Address),pk,h.agg)
	err =u.Send(contract.FUNCTION_REGISTER,payload.TaskId)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,TRANSACTION_ERROR,err)
		return
	}

	res:= &RegisterResponse{
		Response: Response {
			Gcuid:gcuid,
			Status:SUCCESS,
		},
	}

	h.wrapperAndSend(gcuid, res)
}

func (h *Handler) submitHandler(gcuid int, data []byte) {
	var payload SubmitRequest
	err:= json.Unmarshal(data,&payload)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,DATA_FORMAT_ERROR,err)
		return
	}

	pk,err:= derivePrivateKey(payload.PrivateKey)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,KEY_FORMAT_ERROR,err)
		return
	}

	log.Println("data to encrypt:",payload.Value)
	u:= user.NewUser(common.HexToAddress(payload.Address),pk,h.agg)
	cipher,err:=zcrypto.PubKey.Encrypt(payload.Value)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,ENCRYPTION_ERROR,err)
		return
	}

	submitData:=cipher.C.Bytes()
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,DATA_FORMAT_ERROR,err)
		return
	}

	err =u.Send(contract.FUNCTION_SUBMIT,payload.TaskId,submitData)

	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,TRANSACTION_ERROR,err)
		return
	}

	res:= &SubmitResponse{
		Response: Response {
			Gcuid:gcuid,
			Status:SUCCESS,
		},
	}

	h.wrapperAndSend(gcuid, res)
}

func (h *Handler) aggregateHandler(gcuid int, data []byte) {
	var payload AggregateResquest
	err:=json.Unmarshal(data,&payload)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,DATA_FORMAT_ERROR,err)
		return
	}

	pk,err:= derivePrivateKey(payload.PrivateKey)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,KEY_FORMAT_ERROR,err)
		return
	}
	u:= user.NewUser(common.HexToAddress(payload.Address),pk,h.agg)

	countByte,err:= h.agg.Call(contract.FUNCTION_GET_SUBMIT_COUNT_OF_TASK,payload.TaskId)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,CALL_TRANSACTION_ERROR,err)
		return
	}
	count:=new(big.Int).SetBytes(countByte)

	total:= big.NewInt(1)
	for i:=0; i< int(count.Int64()); i++ {
		submitDataByte,err:= h.agg.Call(contract.FUNCTION_GET_SUBMIT_DATA_OF_TASK,payload.TaskId,big.NewInt(int64(i)))
		submitDataByte=submitDataByte[64:]

		if err!=nil {
			log.Println(err.Error())
			h.errorHandler(gcuid,CALL_TRANSACTION_ERROR,err)
			return
		}
		submitData:=new(big.Int).SetBytes(submitDataByte)
		total.Mul(total,submitData)
		total.Mod(total,zcrypto.PubKey.GetNSquare())
	}

	m:=zcrypto.PriKey.Decrypt(&zcrypto.Cypher{C:total})
	log.Println("aggregate result = ",m);

	aggregation:=total.Bytes()
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,UNMARSHAL_JSON_ERROR,err)
		return
	}

	share:= []byte{}
	attestatino := []byte{}
	err =u.Send(contract.FUNCTION_AGGREGATE,payload.TaskId,aggregation,share,attestatino)

	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,TRANSACTION_ERROR,err)
		return
	}

	res:= &AggregateResponse{
		Response: Response {
			Gcuid:gcuid,
			Status:SUCCESS,
		},
	}

	h.wrapperAndSend(gcuid, res)
}

func (h *Handler) approveHandler(gcuid int, data []byte) {
	var payload ApproveRequest
	err:= json.Unmarshal(data,&payload)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid, DATA_FORMAT_ERROR, err)
		return
	}

	pk,err:= derivePrivateKey(payload.PrivateKey)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid, KEY_FORMAT_ERROR, err)
		return
	}
	u:= user.NewUser(common.HexToAddress(payload.Address),pk,h.agg)

	err =u.Send(contract.FUNCTION_APPROVE,payload.TaskId)

	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid, TRANSACTION_ERROR, err)
		return
	}


	countByte,err:= h.agg.Call(contract.FUNCTION_GET_SUBMIT_COUNT_OF_TASK,payload.TaskId)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,CALL_TRANSACTION_ERROR,err)
		return
	}
	count:=new(big.Int).SetBytes(countByte)

	submitValues:=make([]*big.Int,count.Int64(),count.Int64())
	for i:=0; i< int(count.Int64()); i++ {
		submitDataByte,err:= h.agg.Call(contract.FUNCTION_GET_SUBMIT_DATA_OF_TASK,payload.TaskId,big.NewInt(int64(i)))
		submitDataByte=submitDataByte[64:]

		if err!=nil {
			log.Println(err.Error())
			h.errorHandler(gcuid,CALL_TRANSACTION_ERROR,err)
			return
		}
		encryptedSubmitData:=new(big.Int).SetBytes(submitDataByte)

		submitData:=zcrypto.PriKey.Decrypt(&zcrypto.Cypher {
			C: encryptedSubmitData,
		})
		submitValues[i] = submitData
	}

	res:= &ApproveResponse{
		Response: Response {
			Gcuid:gcuid,
			Status:SUCCESS,
		},
		SubmitValues:submitValues,
	}

	h.wrapperAndSend(gcuid, res)
}

func (h *Handler) claimHandler(gcuid int, data []byte) {
	var payload ClaimRequest
	err:= json.Unmarshal(data,&payload)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid, DATA_FORMAT_ERROR, err)
		return
	}

	pk,err:= derivePrivateKey(payload.PrivateKey)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,KEY_FORMAT_ERROR,err)
		return
	}
	u:= user.NewUser(common.HexToAddress(payload.Address),pk,h.agg)

	err =u.Send(contract.FUNCTION_CLAIM,payload.TaskId)

	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,TRANSACTION_ERROR,err)
		return
	}


	res:= &ClaimResponse{
		Response: Response {
			Gcuid:gcuid,
			Status:SUCCESS,
		},
	}

	h.wrapperAndSend(gcuid, res)
}

func (h *Handler) registerAndSubmitHandler(gcuid int, data[]byte) {
	var payload RegisterAndSubmitRequest
	err:= json.Unmarshal(data,&payload)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,DATA_FORMAT_ERROR,err)
		return
	}

	pk,err:= derivePrivateKey(payload.PrivateKey)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,KEY_FORMAT_ERROR,err)
		return
	}

	log.Println("data to encrypt:",payload.Value)
	u:= user.NewUser(common.HexToAddress(payload.Address),pk,h.agg)
	cipher,err:=zcrypto.PubKey.Encrypt(payload.Value)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,ENCRYPTION_ERROR,err)
		return
	}

	submitData:=cipher.C.Bytes()
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,DATA_FORMAT_ERROR,err)
		return
	}

	err =u.Send(contract.FUNCTION_REGISTER_AND_SUBMIT,payload.TaskId,submitData)

	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,TRANSACTION_ERROR,err)
		return
	}

	res:= &RegisterAndSubmitResponse{
		Response: Response {
			Gcuid:gcuid,
			Status:SUCCESS,
		},
	}

	h.wrapperAndSend(gcuid, res)
}

func (h *Handler) stopRegisterAndSubmitHandler (gcuid int, data []byte) {
	var payload stopRegisterAndSubmitRequest
	err:= json.Unmarshal(data,&payload)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid, DATA_FORMAT_ERROR, err)
		return
	}

	pk,err:= derivePrivateKey(payload.PrivateKey)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,KEY_FORMAT_ERROR,err)
		return
	}
	u:= user.NewUser(common.HexToAddress(payload.Address),pk,h.agg)

	err =u.Send(contract.FUNCTION_STOP_REGISTER_AND_SUBMIT,payload.TaskId)

	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,TRANSACTION_ERROR,err)
		return
	}


	res:= &stopRegisterAndSubmitResponse{
		Response: Response {
			Gcuid:gcuid,
			Status:SUCCESS,
		},
	}

	h.wrapperAndSend(gcuid, res)
}

func (h *Handler) getEtherHandler(gcuid int, data []byte) {
	var payload GetEtherRequest
	err:= json.Unmarshal(data,&payload)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,DATA_FORMAT_ERROR,err)
		return
	}

	balance,err:= h.agg.GetEther(common.HexToAddress(payload.Address))
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,CALL_TRANSACTION_ERROR,err)
		return
	}

	res:= &GetEtherResponse{
		Response: Response {
			Gcuid:gcuid,
			Status:SUCCESS,
		},
		Amount:balance,
	}

	h.wrapperAndSend(gcuid, res)
}

func (h *Handler) getBalanceHandler(gcuid int, data[]byte) {
	var payload GetBalanceRequest
	err:= json.Unmarshal(data,&payload)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,DATA_FORMAT_ERROR,err)
		return
	}

	countByte,err:= h.token.Call(contract.FUNCTION_BALANCE_OF,common.HexToAddress(payload.Address))
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,CALL_TRANSACTION_ERROR,err)
		return
	}
	count:= new(big.Int).SetBytes(countByte)

	res:= &GetBalanceResponse{
		Response:Response{
			Gcuid:gcuid,
			Status:SUCCESS,
		},
		Amount:count,
		Address:payload.Address,
	}

	h.wrapperAndSend(gcuid, res)
}

func (h *Handler) getStageHandler(gcuid int, data []byte) {
	var payload GetStageRequest
	err:= json.Unmarshal(data,&payload)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,DATA_FORMAT_ERROR,err)
		return
	}

	stageByte,err:= h.agg.Call(contract.FUNCTION_GET_STAGE_OF_TASK,payload.TaskId)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,CALL_TRANSACTION_ERROR,err)
		return
	}
	stage:= new(big.Int).SetBytes(stageByte)

	res:= &GetStageResponse{
		Response:Response{
			Gcuid:gcuid,
			Status:SUCCESS,
		},
		Stage:stage,
	}
	h.wrapperAndSend(gcuid, res)
}

func (h* Handler) getRegisterNumberHandler(gcuid int, data []byte) {
	var payload GetRegisterNumberRequest
	err:= json.Unmarshal(data,&payload)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,DATA_FORMAT_ERROR,err)
		return
	}

	countByte,err:= h.agg.Call(contract.FUNCTION_GET_REGISTER_NUMBER_OF_TASK,payload.TaskId)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,CALL_TRANSACTION_ERROR,err)
		return
	}

	count:= new(big.Int).SetBytes(countByte)

	res:= &GetRegisterNumberResponse{
		Response:Response{
			Gcuid:gcuid,
			Status:SUCCESS,
		},
		Amount: count,
	}
	h.wrapperAndSend(gcuid, res)
}

func (h *Handler) getSubmissionNumberHandler(gcuid int, data []byte) {
	var payload GetSubmissionNumberRequest
	err:= json.Unmarshal(data,&payload)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,DATA_FORMAT_ERROR,err)
		return
	}

	countByte,err:= h.agg.Call(contract.FUNCTION_GET_SUBMIT_COUNT_OF_TASK,payload.TaskId)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,CALL_TRANSACTION_ERROR,err)
		return
	}
	count:= new(big.Int).SetBytes(countByte)

	res:= &GetSubmissionNumberResponse{
		Response:Response{
			Gcuid:gcuid,
			Status:SUCCESS,
		},
		Amount:count,
	}

	h.wrapperAndSend(gcuid, res)
}

func (h *Handler) getClaimNumberHandler(gcuid int, data[]byte) {
	var payload GetClaimNumberResquest
	err:= json.Unmarshal(data,&payload)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,DATA_FORMAT_ERROR,err)
		return
	}

	countByte,err:= h.agg.Call(contract.FUNCTION_GET_CLAIM_NUMBER_OF_TASK,payload.TaskId)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,CALL_TRANSACTION_ERROR,err)
		return
	}
	count:= new(big.Int).SetBytes(countByte)

	res:= &GetClaimNumberResponse{
		Response:Response{
			Gcuid:gcuid,
			Status:SUCCESS,
		},
		Amount:count,
	}

	h.wrapperAndSend(gcuid, res)
}

func (h *Handler) getAggregationResultHandler(gcuid int, data[]byte) {
	var payload GetAggregateResultRequest
	err:= json.Unmarshal(data,&payload)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,DATA_FORMAT_ERROR,err)
		return
	}
	aggregateResultByte,err:= h.agg.Call(contract.FUNCTION_GET_AGGREGATION_RESULT_OF_TASK, payload.TaskId)
	aggregateResultByte = aggregateResultByte[64:]
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,CALL_TRANSACTION_ERROR,err)
		return
	}

	aggregateResult:=zcrypto.PriKey.Decrypt(&zcrypto.Cypher {
		C: new(big.Int).SetBytes(aggregateResultByte),
	})
	log.Println("aggregation result:",aggregateResult)

	res:= &GetAggregateResultResponse{
		Response:Response{
			Gcuid:gcuid,
			Status:SUCCESS,
		},
		Amount:aggregateResult,
	}
	h.wrapperAndSend(gcuid,res)
}

func (h *Handler) getSolicitInfoHandler(gcuid int, data[]byte) {
	var payload GetSolicitInfoResquest
	err:= json.Unmarshal(data,&payload)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,DATA_FORMAT_ERROR,err)
		return
	}

	infoByte,err:= h.agg.Call(contract.FUNCTION_GET_SOLICITINFO_OF_TASK, payload.TaskId)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,CALL_TRANSACTION_ERROR,err)
		return
	}

	var info contract.SolicitInfo
	info.DataFee = new(big.Int).SetBytes(infoByte[:32])
	info.ServiceFee = new(big.Int).SetBytes(infoByte[32:64])
	info.ServiceProvider = common.HexToAddress("0x"+hex.EncodeToString(infoByte[76:96]))
	info.Target = new(big.Int).SetBytes(infoByte[96:128])


	res:= &GetSolicitInfoResponse{
		Response: Response{
			Gcuid:gcuid,
			Status:SUCCESS,
		},
		DataFee:info.DataFee,
		ServiceFee: info.ServiceFee,
		ServiceProvider: info.ServiceProvider.String(),
		Target: info.Target,
	}

	h.wrapperAndSend(gcuid,res)
}

func (h *Handler) sendTransactionHandler(gcuid int,data []byte) {
	var payload SendTransactionRequest
	err:= json.Unmarshal(data,&payload)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,DATA_FORMAT_ERROR,err)
		return
	}

	rawTx,err:= hex.DecodeString(payload.RawTransaction)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,DATA_FORMAT_ERROR,err)
		return
	}

	var tx types.Transaction

	rlp.DecodeBytes(rawTx,&tx)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,DATA_FORMAT_ERROR,err)
		return
	}

	_, err =h.agg.SendTransaction(&tx)
	if err!=nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,TRANSACTION_ERROR,err)
		return
	}

	res:= &SendTransactionResponse{
		Response: Response{
			Gcuid:gcuid,
			Status:SUCCESS,
		},
	}

	h.wrapperAndSend(gcuid,res)
}

func (h *Handler) HandleRequest () {
	for {
		_,data,err:= h.w.ReadMessage()
		var kvs map[string]interface{}
		if err != nil {
			log.Println(err.Error())
			return
		}
		json.Unmarshal(data, &kvs)
		gid, ok := kvs["gcuid"]
		if !ok {
			log.Println(errors.New("gcuid not exist"))
			continue
		}
		gcuid := int64(gid.(float64))
		log.Println("gcuid:", gcuid)

		switch gcuid {
		case GCUID_SOLICIT:
			go h.solicitHandler(GCUID_SOLICIT,data)
		case GCUID_REGISTER:
			go h.registerHandler(GCUID_REGISTER,data)
		case GCUID_SUBMIT:
			go h.submitHandler(GCUID_SUBMIT,data)
		case GCUID_AGGREGATION:
			go h.aggregateHandler(GCUID_AGGREGATION,data)
		case GCUID_APPROVE:
			go h.approveHandler(GCUID_APPROVE,data)
		case GCUID_CLAIM:
			go h.claimHandler(GCUID_CLAIM,data)
		case GCUID_REGISTER_AND_SUBMIT:
			go h.registerAndSubmitHandler(GCUID_REGISTER_AND_SUBMIT,data)
		case GCUID_STOP_REGISTER_AND_SUBMIT:
			go h.stopRegisterAndSubmitHandler(GCUID_STOP_REGISTER_AND_SUBMIT,data)
		case GCUID_ETHER:
			go h.getEtherHandler(GCUID_ETHER,data)
		case GCUID_CURRENT_STAGE:
			go h.getStageHandler(GCUID_CURRENT_STAGE,data)
		case GCUID_SUBMISSION_NUMBER:
			go h.getSubmissionNumberHandler(GCUID_SUBMISSION_NUMBER,data)
		case GCUID_REGISTER_NUMBER:
			go h.getRegisterNumberHandler(GCUID_REGISTER_NUMBER,data)
		case GCUID_SOLICIT_INFO:
			go h.getSolicitInfoHandler(GCUID_SOLICIT_INFO,data)
		case GCUID_AGGREGATE_RESULT:
			go h.getAggregationResultHandler(GCUID_AGGREGATE_RESULT,data)
		case GCUID_CLAIM_NUMBER:
			go h.getClaimNumberHandler(GCUID_CLAIM_NUMBER,data)
		case GCUID_BALANCE:
			go h.getBalanceHandler(GCUID_BALANCE,data)
		case GCUID_SEND_TRANSACTION:
			go h.sendTransactionHandler(GCUID_SEND_TRANSACTION,data)
		}
	}
}
