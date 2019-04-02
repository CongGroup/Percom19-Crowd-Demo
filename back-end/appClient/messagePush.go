package appClient

import (
	"Percome19-Crowd-Demo/back-end/contract"
	"Percome19-Crowd-Demo/back-end/zcrypto"
	"encoding/hex"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"math/big"
)

const CLIENT_BUFFER = 1024

type HandlerManager struct {
	Handlers   map[*Handler]bool
	Broadcast  chan []byte
	Register   chan *Handler
	Unregister chan *Handler
}

func NewHandlerManager() *HandlerManager {
	return &HandlerManager{
		Handlers:   make(map[*Handler]bool),
		Register:   make(chan *Handler, CLIENT_BUFFER),
		Unregister: make(chan *Handler, CLIENT_BUFFER),
		Broadcast:  make(chan []byte, CLIENT_BUFFER),
	}
}

func (this *HandlerManager) Start() {
	for {
		select {
		case conn := <-this.Register:
			this.Handlers[conn] = true
			log.Println("a new user register message push")
		case conn := <-this.Unregister:
			if _, ok := this.Handlers[conn]; ok {
				delete(this.Handlers, conn)
				log.Println("a user unregister message push")
			}
		case message := <-this.Broadcast:
			log.Println("broadcast a message:", string(message))
			for handler := range this.Handlers {
				go func(h *Handler) {
					err := h.Send(message)
					if err != nil {
						log.Println("push message to a close socket")
					}
				}(handler)
			}
		}
	}
}

func (this *HandlerManager) SubScriptContractEvent(c contract.Contract) {
	log.Println("start event watcher");
	logs, eventError := c.EventWatcher()

	for {
		select {
		case err := <-eventError:
			log.Println(err.Error())
		case vLog := <-logs:
			switch vLog.Topics[0].Hex() {
			case contract.LogStageTransfer:
				var stageTransferEvent contract.StageTransferEvent
				err := c.Unpack(&stageTransferEvent, contract.EVENT_STAGE_TRANSFER, vLog.Data)
				if err != nil {
					log.Println(err.Error())
					return
				}

				gcuid := GCUID_CURRENT_STAGE
				res := &GetStageResponse{
					Response: Response{
						Gcuid:  gcuid,
						Status: SUCCESS,
					},
					Stage: stageTransferEvent.NewStage,
				}
				data, err := json.Marshal(res)
				if err != nil {
					log.Println(err.Error())
					return
				}

				this.Broadcast <- data

			case contract.LogSolicit:
				var solicitEvent contract.SolicitEvent
				err := c.Unpack(&solicitEvent, contract.EVENT_SOLICIT, vLog.Data)
				if err != nil {
					log.Println(err.Error())
					return
				}

				infoByte, err := c.Call(contract.FUNCTION_GET_SOLICITINFO_OF_TASK, solicitEvent.TaskId)
				if err != nil {
					log.Println(err.Error())
					return
				}

				var info contract.SolicitInfo
				info.DataFee = new(big.Int).SetBytes(infoByte[:32])
				info.ServiceFee = new(big.Int).SetBytes(infoByte[32:64])
				info.ServiceProvider = common.HexToAddress("0x" + hex.EncodeToString(infoByte[76:96]))
				info.Target = new(big.Int).SetBytes(infoByte[96:128])

				gcuid := GCUID_SOLICIT_INFO
				res := &GetSolicitInfoResponse{
					Response: Response{
						Gcuid:  gcuid,
						Status: SUCCESS,
					},
					DataFee:         info.DataFee,
					ServiceFee:      info.ServiceFee,
					ServiceProvider: info.ServiceProvider.String(),
					Target:          info.Target,
				}
				data, err := json.Marshal(res)
				if err != nil {
					log.Println(err.Error())
					return
				}

				this.Broadcast <- data
			case contract.LogRegister:
				var registerEvent contract.RegisterEvent
				err := c.Unpack(&registerEvent, contract.EVENT_REGISTER, vLog.Data)
				if err != nil {
					log.Println(err.Error())
					return
				}

				gcuid := GCUID_REGISTER_NUMBER
				res := &GetRegisterNumberResponse{
					Response: Response{
						Gcuid:  gcuid,
						Status: SUCCESS,
					},
					Amount: registerEvent.RegisterNumber,
				}
				data, err := json.Marshal(res)
				if err != nil {
					log.Println(err.Error())
					return
				}
				this.Broadcast <- data
			case contract.LogSubmit:
				var submitEvent contract.SubmitEvent
				err := c.Unpack(&submitEvent, contract.EVENT_SUBMIT, vLog.Data)
				if err != nil {
					log.Println(err.Error())
					return
				}

				gcuid := GCUID_SUBMISSION_NUMBER
				res := &GetSubmissionNumberResponse{
					Response: Response{
						Gcuid:  gcuid,
						Status: SUCCESS,
					},
					Amount: submitEvent.SubmitNumber,
				}
				data, err := json.Marshal(res)
				if err != nil {
					log.Println(err.Error())
					return
				}
				this.Broadcast <- data
			case contract.LogAggregate:
				var aggregateEvent contract.AggregateEvent
				err := c.Unpack(&aggregateEvent, contract.EVENT_AGGREGATE, vLog.Data)
				if err != nil {
					log.Println(err.Error())
					return
				}

				gcuid := GCUID_AGGREGATE_RESULT

				aggregateResultByte, err := c.Call(contract.FUNCTION_GET_AGGREGATION_RESULT_OF_TASK, aggregateEvent.TaskId)
				if err != nil {
					log.Println(err.Error())
					return
				}
				aggregateResultLen := new(big.Int).SetBytes(aggregateResultByte[32:64]).Int64()
				aggregateResultByte = aggregateResultByte[64 : 64+aggregateResultLen]

				qualifiedNumberByte, err := c.Call(contract.FUNCTION_GET_QUALIFIED_NUMBER_OF_TASK, aggregateEvent.TaskId)
				if err != nil {
					log.Println(err.Error())
					return
				}

				aggregateResult := zcrypto.PriKey.Decrypt(&zcrypto.Cypher{
					C: new(big.Int).SetBytes(aggregateResultByte),
				})

				res := &GetAggregateResultResponse{
					Response: Response{
						Gcuid:  gcuid,
						Status: SUCCESS,
					},
					Amount:          aggregateResult,
					QualifiedNumber: new(big.Int).SetBytes(qualifiedNumberByte),
				}

				data, err := json.Marshal(res)
				if err != nil {
					log.Println(err.Error())
					return
				}
				this.Broadcast <- data
			case contract.LogClaim:
				var claimEvent contract.ClaimEvent
				err := c.Unpack(&claimEvent, contract.EVENT_CLAIM, vLog.Data)
				if err != nil {
					log.Println(err.Error())
					return
				}

				gcuid := GCUID_CLAIM_NUMBER
				res := &GetClaimNumberResponse{
					Response: Response{
						Gcuid:  gcuid,
						Status: SUCCESS,
					},
					Amount: claimEvent.ClaimNumber,
				}
				data, err := json.Marshal(res)
				if err != nil {
					log.Println(err.Error())
					return
				}
				this.Broadcast <- data
			}
		}
	}
}
