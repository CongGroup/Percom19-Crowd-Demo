package main

import (
	"appClient"
	"contract"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"zcrypto"
)


func testPallier() {
	pub,pri,err:=zcrypto.NewPallier(256);
	cypher,err:=pub.Encrypt(big.NewInt(232131230))
	if err!=nil {
		log.Println()
		return
	}
	log.Println("N: ",pri.N," G: ",pri.G, "lambda", pri.Lambda, "mu:", pri.Mu," cipher:",cypher.C)
	m:= pri.Decrypt(cypher)
	log.Println(m)
}

func testBulletProof2() {
	rp:=zcrypto.RPProve(big.NewInt(25))
	result := zcrypto.RPVerify(rp);
	log.Println(len(rp.IPP.L),len(rp.IPP.R),len(rp.IPP.Challenges))
	if result {
		log.Println("pass test")
	} else {
		log.Println("test fail")
	}

}

func requestParserWrapper(manager *appClient.HandlerManager, agg *contract.Agg, token *contract.ERC20) func(w http.ResponseWriter, r *http.Request){
	return func (w http.ResponseWriter, r *http.Request){
		log.Println("receive a reqeust")
		var upgrader = websocket.Upgrader{}
		upgrader.CheckOrigin = func(rq *http.Request) bool { return true }
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}
		handler:= appHandler(c,agg,token)
		manager.Register <- handler
		handler.HandleRequest()
		manager.Unregister <- handler
		defer c.Close()
	}
}

func appHandler(c *websocket.Conn, agg *contract.Agg, token *contract.ERC20) *appClient.Handler {
	return appClient.NewHandler(c,agg,token)
}

func getNonce(agg *contract.Agg) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("get nonce")
		vars:=mux.Vars(r)
		user:= vars["user"]
		nonce,err:= agg.GetNonce(common.HexToAddress(user))
		if err!=nil {
			log.Println(err.Error())
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}
		log.Println("user:",user)
		log.Println("nonce:",nonce)
		nonceWrapper, err:= json.Marshal(nonce)
		if err!=nil {
			log.Println(err.Error())
			log.Println(appClient.UNMARSHAL_JSON_ERROR)
			return
		}
		w.Write(nonceWrapper)
	}
}

func getChainId(agg *contract.Agg) func(w http.ResponseWriter, r* http.Request) {
	return func(w http.ResponseWriter,r *http.Request){
		log.Println("get chain id")

		chainId,err:= agg.GetChainId()
		if err!=nil {
			log.Println(err.Error())
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}

		chainIdWrapper, err:= json.Marshal(chainId)
		if err!=nil {
			log.Println(err.Error())
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}
		w.Write(chainIdWrapper)
	}
}

func getEncryptedData(w http.ResponseWriter, r *http.Request) {
	log.Println("get encrypted data")
	var amount int64
	data, err:= ioutil.ReadAll(r.Body)
	if err!=nil {
		log.Println(err.Error())
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(data,&amount)
	if err!=nil {
		log.Println(err.Error())
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	log.Println("data to encrypt:",amount)



	var submitProof string

	if amount<=appClient.MIN_RANGE || amount>= appClient.MAX_RANGE {
		log.Println("out of range, generate random proof");
		data:= make([]byte,1248)
		submitProof="0x"+hex.EncodeToString(data)
	} else {
		rpV := zcrypto.RPProve(big.NewInt(amount));
		submitProof="0x"+hex.EncodeToString(rpV.Bytes())
	}

	cipher,err:=zcrypto.PubKey.Encrypt(big.NewInt(amount))
	if err!=nil {
		log.Println(err.Error())
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	submitData := "0x"+cipher.C.Text(16)
	//log.Println("proof:",submitProof)

	submitPayload,err := json.Marshal(&appClient.SubmitPayload{
		SubmitData:submitData,
		SubmitProof:submitProof,
	})
	if err!=nil {
		log.Println(err.Error())
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	w.Write(submitPayload)
}

func getStatistics(agg *contract.Agg) func(w http.ResponseWriter, r* http.Request) {
	return func(w http.ResponseWriter,r *http.Request){
		log.Println("get statistics")
		vars:=mux.Vars(r)
		taskId,_:= new(big.Int).SetString(vars["taskId"],10)

		countByte,err:= agg.Call(contract.FUNCTION_GET_REGISTER_NUMBER_OF_TASK,taskId)
		if err!=nil {
			log.Println(err.Error())
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}
		count:=new(big.Int).SetBytes(countByte)
		submitValues:=make([]*big.Int,count.Int64(),count.Int64())
		for i:=0; i< int(count.Int64()); i++ {
			submitDataByte,err:= agg.Call(contract.FUNCTION_GET_SUBMIT_DATA_OF_TASK, taskId,big.NewInt(int64(i)))
			submitDataByte=submitDataByte[64:]

			if err!=nil {
				log.Println(err.Error())
				http.Error(w,err.Error(),http.StatusInternalServerError)
				return
			}
			encryptedSubmitData:=new(big.Int).SetBytes(submitDataByte)

			submitData:=zcrypto.PriKey.Decrypt(&zcrypto.Cypher {
				C: encryptedSubmitData,
			})
			submitValues[i] = submitData
		}

		type payload struct {
			SubmitValues []*big.Int `json:"submitValues"`
		}
		submitValuesWrapper,err :=json.Marshal(&payload{
			submitValues,
		})
		if err!=nil {
			log.Println(err.Error())
			http.Error(w,err.Error(),http.StatusInternalServerError)
		}
		w.Write(submitValuesWrapper)
	}
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func start() {
	r := mux.NewRouter()

	manager:= appClient.NewHandlerManager()
	agg:= contract.NewAgg(contract.GETH_PORT,contract.CONTRACT_ABI,common.HexToAddress(contract.CONTRACT_ADDRESS))
	token:= contract.NewERC20(contract.GETH_PORT,contract.ERC20_ABI,common.HexToAddress(contract.ERC20_ADDRESS))

	go manager.Start()
	go manager.SubScriptContractEvent(agg)

	r.HandleFunc("/", requestParserWrapper(manager,agg,token)).Methods("GET")
	r.HandleFunc("/nonce/{user}",getNonce(agg)).Methods("GET")
	r.HandleFunc("/chainId",getChainId(agg)).Methods("GET")
	r.HandleFunc("/encryptedData",getEncryptedData).Methods("POST")
	r.HandleFunc("/statistics/{taskId}",getStatistics(agg)).Methods("GET")

	fmt.Println("Running http server")
	http.ListenAndServe(
		"0.0.0.0:4000",
		handlers.CORS(
			handlers.AllowedMethods([]string{"get", "options", "post", "put", "head"}),
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedHeaders([]string{"Content-Type"}),
		)(r),
	)
}

func main() {
	//testPallier()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//testBulletProof()
	//testBulletProof2()

	start()
}
