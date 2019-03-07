package app

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"math/big"
	"math/rand"
	"sync"
	"time"
)

const (
	RECONNECT_TIME = 4 * time.Second
	TRANSACTION_RETRY_TIME = 4 * time.Second
	SEND_BUFFER = 256
	TASKID = 0
	SUBMIT_DATA_MAX = 65536
)

type Client struct {
	W *websocket.Conn
	Send chan []byte
	Url string
	Wg sync.WaitGroup
	Closed chan int
	Account *Wallet

	Id int
}

func NewClient(id int,url string, account *Wallet) *Client {
	// connect
	w,_,err := websocket.DefaultDialer.Dial(url,nil)
	if err!=nil {
		log.Println(err.Error())
		panic(err)
	}

	return &Client{
		Id: id,
		W: w,
		Send: make(chan []byte,SEND_BUFFER),
		Url: url,
		Wg: sync.WaitGroup{},
		Closed: make(chan int),
		Account:account,
	}
}

func (c *Client) Reader() {
	for {
		_, message, err:= c.W.ReadMessage()
		if err!=nil {
			log.Println("read:",err)
			c.Closed <- 1;
			break;
		}
		// do with message
		var kvs map[string]interface{}
		json.Unmarshal(message,&kvs)
		gid,ok:= kvs["gcuid"]
		if !ok {
			log.Println("unknown message")
		}
		gcuid:= int64(gid.(float64))

		switch gcuid {
		case GCUID_CURRENT_STAGE:
			go c.StageChangeHandler(GCUID_CURRENT_STAGE,message)
		case GCUID_CLAIM:
			go c.ClaimHandler(GCUID_CLAIM,message)
		case GCUID_REGISTER_AND_SUBMIT:
			go c.RegisterAndSubmitHandler(GCUID_REGISTER_AND_SUBMIT,message)
		default:
		}
	}
	close(c.Closed)
	c.Wg.Done()
}

func (c *Client) Reconnect() error {
	log.Println("reconnect ",c.Id)
	c.W.Close()
	for{
		select {
		case <-time.After(RECONNECT_TIME*time.Second):
			var err error
			c.W,_,err= websocket.DefaultDialer.Dial(c.Url,nil)
			if err!=nil {
				log.Println(err.Error())
			} else {
				return nil
			}
			c.Closed = make(chan int)
		}
	}
}

func (c *Client) Start() {
	for {
		c.Wg.Add(1)
		c.Wg.Add(1)
		go c.Reader()
		go c.Sender()
		c.Wg.Wait()
		c.Reconnect()
	}
	defer func () {
		c.W.Close()
		close(c.Send)
		close(c.Closed)
	}()
}

func (c *Client) Sender() {
	for {
		select {
		case <-c.Closed:
			break;
		case data:= <-c.Send:
			err:=c.W.WriteMessage(websocket.TextMessage,data)
			if err!=nil {
				log.Println(err)
				// re-queue
				c.Send<-data
				break
			}
		}
	}
	c.Wg.Done()
}


func (c *Client) StageChangeHandler(gcuid int,data []byte) {
	var stageResponse GetStageResponse
	err:=json.Unmarshal(data,&stageResponse)
	if err!=nil {
		log.Println(err.Error())
		return
	}
	stage:= stageResponse.Stage
	switch stage.Int64() {
	case REGISTER:
		log.Println("start register")
		c.registerAndSubmit(GCUID_REGISTER_AND_SUBMIT)
	case CLAIM:
		log.Println("begin to claim")
		c.claim(GCUID_CLAIM)
	default:
	}
}

func (c *Client) claim(gcuid int) {
	claimRequest:= &ClaimRequest{
		Gcuid: gcuid,
		TaskId:  big.NewInt(TASKID),
		PrivateKey: c.Account.PrivateKey.D.Text(16),
		Address: c.Account.Address.String(),
	}
	data,err:=json.Marshal(claimRequest)
	if err!=nil {
		log.Println(err.Error())
		return
	}
	c.Send <- data
}

func (c *Client) registerAndSubmit(gcuid int) {
	registerAndSubmitRequest:=&RegisterAndSubmitRequest{
		Gcuid: gcuid,
		TaskId: big.NewInt(TASKID),
		PrivateKey: c.Account.PrivateKey.D.Text(16),
		Address:c.Account.Address.String(),
		Value: big.NewInt(int64(rand.Intn(SUBMIT_DATA_MAX))),
	}
	data,err:=json.Marshal(registerAndSubmitRequest)
	if err!=nil {
		log.Println(err.Error())
		return
	}
	c.Send <- data
}

func (c *Client) ClaimHandler(gcuid int, data[]byte) {
	var claimResponse ClaimResponse
	err:=json.Unmarshal(data,&claimResponse)
	if err!=nil {
		log.Println(err.Error())
		return
	}
	if claimResponse.Status == SUCCESS {
	} else {
		var errorResponse Error
		err:=json.Unmarshal(data,&errorResponse)
		if err!=nil {
			log.Println(err.Error())
			return
		}
		log.Println("claim fail, reason:",errorResponse.Reason)
		//<-time.After(TRANSACTION_RETRY_TIME)
		//c.claim(gcuid)
	}
}

func (c *Client) RegisterAndSubmitHandler(gcuid int, data[]byte) {
	var registerAndSubmitResponse RegisterAndSubmitResponse
	err:=json.Unmarshal(data,&registerAndSubmitResponse)
	if err!=nil {
		log.Println(err.Error())
		return
	}
	if registerAndSubmitResponse.Status == SUCCESS {
	} else {
		var errorResponse Error
		err:=json.Unmarshal(data,&errorResponse)
		if err!=nil {
			log.Println(err.Error())
			return
		}
		log.Println("register and submit fail, reason:",errorResponse.Reason)
		//<-time.After(TRANSACTION_RETRY_TIME)
		//c.registerAndSubmit(gcuid)
	}
}