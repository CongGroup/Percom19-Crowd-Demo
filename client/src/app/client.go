package app

import (
	"crypto/tls"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"math/big"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const (
	RECONNECT_TIME = 2 * time.Second
	TRANSACTION_RETRY_TIME = 4 * time.Second
	SEND_BUFFER = 256
	TASKID = 0
	SUBMIT_DATA_MAX = 65536
	SUBMIT_DATA_SPACE = 4096
	SUBMIT_DATA_MIN = 0
)

const (
	writeWait = 4 * time.Second
	pongWait = 6*time.Second
	pingPeroid = 4*time.Second
)

type Client struct {
	W *websocket.Conn
	Send chan []byte
	Url string
	HttpPath string
	Wg sync.WaitGroup
	Account *Wallet

	Id int
}

func NewClient(id int,url string, httpPath string, account *Wallet) *Client {

	return &Client{
		Id: id,
		Send: make(chan []byte,SEND_BUFFER),
		Url: url,
		Wg: sync.WaitGroup{},
		HttpPath:httpPath,
		Account:account,
	}
}

func (c *Client) Reader() {
	c.W.SetReadDeadline(time.Now().Add(pongWait))
	c.W.SetPongHandler( func(string) error {
		c.W.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	Outter:
	for {
		_, message, err:= c.W.ReadMessage()
		if err!=nil {
			log.Println("read:",err)
			break Outter
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
	c.Wg.Done()
}

func (c *Client) Reconnect() error {
	log.Println("reconnect ",c.Id)
	c.W.Close()
	Outter:
	for{
		select {
		case <-time.After(RECONNECT_TIME):
			var err error
			dialer:= websocket.Dialer{
				TLSClientConfig: &tls.Config {
					InsecureSkipVerify: true,
				},
			}
			c.W,_,err= dialer.Dial(c.Url,nil)
			if err!=nil {
				log.Println(err.Error())
			} else {
				break Outter;
			}
		}
	}
	log.Println("reconnected")
	c.GetCurrentStage(GCUID_CURRENT_STAGE)
	return nil
}

func (c *Client) Start() {
	// connect
	for{
		dialer:= websocket.Dialer{
			TLSClientConfig: &tls.Config {
				InsecureSkipVerify: true,
			},
		}
		w,_,err:= dialer.Dial(c.Url,nil)
		if err!=nil {
			log.Println(err.Error())
			<-time.After(RECONNECT_TIME)
			continue
		}
		c.W =w;
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		_,err=client.Get(c.HttpPath+"/requireEther/"+c.Account.Address.String())
		if err!=nil {
			log.Println(err.Error())
		}
		break;
	}

	for {
		c.Wg.Add(1)
		c.Wg.Add(1)
		go c.Reader()
		go c.Sender()
		c.GetCurrentStage(GCUID_CURRENT_STAGE)
		c.Wg.Wait()
		c.Reconnect()
	}
	defer func () {
		c.W.Close()
		close(c.Send)
	}()
}

func (c *Client) Sender() {
	ticker:= time.NewTicker(pingPeroid)
	Outter:
	for {
		select {
		case <-ticker.C:
			c.W.SetWriteDeadline(time.Now().Add(writeWait))
			err:=c.W.WriteMessage(websocket.PingMessage,nil)
			if err!=nil {
				log.Println(err.Error())
				break Outter
			}
		case data:= <-c.Send:
			c.W.SetWriteDeadline(time.Now().Add(writeWait))
			err:=c.W.WriteMessage(websocket.TextMessage,data)
			if err!=nil {
				log.Println(err)
				// re-queue
				c.Send <- data
				break Outter
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
	var submitData *big.Int
	if c.Id == 0  {
		submitData,_ = new(big.Int).SetString("2131232132141208903824928213",10)
	} else if c.Id == 1 {
		submitData,_ = new(big.Int).SetString("-65464334564",10)
	} else if c.Id == 2 {
		submitData,_ = new(big.Int).SetString("9546653",10)
	} else {
		var randData int
		randData = rand.Intn(SUBMIT_DATA_SPACE) + c.Id * SUBMIT_DATA_SPACE
		if randData >= SUBMIT_DATA_MAX {
			randData = rand.Intn(SUBMIT_DATA_MAX)
		}
		submitData = big.NewInt(int64(randData))
	}
	//submitData:= rand.Intn(SUBMIT_DATA_SPACE) + c.Id * SUBMIT_DATA_SPACE
	registerAndSubmitRequest:=&RegisterAndSubmitRequest{
		Gcuid: gcuid,
		TaskId: big.NewInt(TASKID),
		PrivateKey: c.Account.PrivateKey.D.Text(16),
		Address:c.Account.Address.String(),
		Value: submitData,
	}
	data,err:=json.Marshal(registerAndSubmitRequest)
	if err!=nil {
		log.Println(err.Error())
		return
	}
	c.Send <- data
}

func (c *Client) GetCurrentStage(gcuid int) {
	getStageRequest:= &GetStageRequest{
		Gcuid: gcuid,
		TaskId: big.NewInt(TASKID),
	}
	data,err:=json.Marshal(getStageRequest)
	if err!=nil {
		log.Println(err.Error())
		return
	}
	c.Send <-data
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