package appClient

import (
	"contract"
	"encoding/json"
	"log"
)


const CLIENT_BUFFER = 1024

type HandlerManager struct {
	Handlers map[*Handler] bool
	Broadcast chan []byte
	Register chan *Handler
	Unregister chan *Handler
}

func NewHandlerManager() *HandlerManager {
	return &HandlerManager{
		Handlers: make(map[*Handler]bool),
		Register: make(chan *Handler, CLIENT_BUFFER),
		Unregister:make(chan *Handler,CLIENT_BUFFER),
		Broadcast: make(chan []byte, CLIENT_BUFFER),
	}
}

func (this *HandlerManager) Start() {
	for {
		select {
			case conn:= <-this.Register:
				this.Handlers[conn] = true
				log.Println("a new user register message push")
				case conn:= <- this.Unregister:
					if _,ok:=this.Handlers[conn];ok{
						delete(this.Handlers,conn)
						log.Println("a user unregister message push")
					}
				case message:= <- this.Broadcast:
					log.Println("broadcast a message:",string(message))
					for handler:= range this.Handlers {
						log.Println("")
						go func() {
							err:=handler.Send(message)
							if err!=nil {
								log.Println("push message to a close socket")
							}
						}()
					}
		}
	}
}

func (this *HandlerManager) SubScriptContractEvent(c contract.Contract) {
	log.Println("start event watcher");
	logs,eventError:= c.EventWatcher()

	for {
		select {
		case err:= <-eventError:
			log.Println(err.Error())
		case vLog:= <-logs:
			switch  vLog.Topics[0].Hex() {
			case contract.LogStageTransfer:
				var stageTransferEvent contract.StageTransferEvent
				err:=c.Unpack(&stageTransferEvent, contract.EVENT_STAGE_TRANSFER, vLog.Data)
				if err!=nil {
					log.Println(err.Error())
					return
				}

				gcuid:= GCUID_CURRENT_STAGE
				res:= &GetStageResponse{
					Response: Response{
						Gcuid:gcuid,
						Status: SUCCESS,
					},
					Stage: stageTransferEvent.NewStage,
				}
				data,err:=json.Marshal(res)
				if err!=nil {
					log.Println(err.Error())
					return
				}

				this.Broadcast<-data

			case contract.LogSolicit:
				// TODO
				//var res appClient.GetSolicitInfoResponse
				//gcuid:= appClient.GCUID_SOLICIT_INFO
			case contract.LogRegister:
				var registerEvent contract.RegisterEvent
				err:= c.Unpack(&registerEvent, contract.EVENT_REGISTER, vLog.Data)
				if err!=nil {
					log.Println(err.Error())
					return
				}

				gcuid:= GCUID_REGISTER_NUMBER
				res:= &GetRegisterNumberResponse{
					Response:Response{
						Gcuid: gcuid,
						Status: SUCCESS,
					},
					Amount: registerEvent.RegisterNumber,
				}
				data,err:=json.Marshal(res)
				if err!=nil {
					log.Println(err.Error())
					return
				}
				this.Broadcast <-data
			case contract.LogSubmit:
				var submitEvent contract.SubmitEvent
				err:= c.Unpack(&submitEvent,contract.EVENT_SUBMIT,vLog.Data)
				if err!=nil {
					log.Println(err.Error())
					return
				}

				gcuid:= GCUID_SUBMISSION_NUMBER
				res:= &GetSubmissionNumberResponse {
					Response: Response{
						Gcuid:gcuid,
						Status:SUCCESS,
					},
					Amount:submitEvent.SubmitNumber,
				}
				data,err:=json.Marshal(res)
				if err!=nil {
					log.Println(err.Error())
					return
				}
				this.Broadcast <-data
			case contract.LogAggregate:
				// TODO
				//var res appClient.GetAggregateResultResponse
				//gcuid:= appClient.GCUID_AGGREGATE_RESULT
				//case logApprove:
				//	gcuid:= appClient.GCUIDD_AGGREGATE_RESULT
			case contract.LogClaim:
				var claimEvent contract.ClaimEvent
				err:= c.Unpack(&claimEvent,contract.EVENT_CLAIM,vLog.Data)
				if err!=nil {
					log.Println(err.Error())
					return
				}

				gcuid:= GCUID_CLAIM_NUMBER
				res:= &GetClaimNumberResponse {
					Response: Response{
						Gcuid:gcuid,
						Status: SUCCESS,
					},
					Amount:claimEvent.ClaimNumber,
				}
				data,err:=json.Marshal(res)
				if err!=nil {
					log.Println(err.Error())
					return
				}
				this.Broadcast <-data
			}
		}
	}
}