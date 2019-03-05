package main

import (
	"appClient"
	"contract"
	"encoding/json"
	"fmt"
	"github.com/deroproject/derosuite/crypto"
	"github.com/deroproject/derosuite/crypto/ringct"
	"github.com/ethereum/go-ethereum/common"
	h "github.com/gorilla/handlers"
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

func testBulletProof() {
	s1:=*(crypto.RandomScalar())
	bp:=ringct.BULLETPROOF_Prove_Amount(14,&s1)
	if !bp.BULLETPROOF_Verify_ultrafast() {
		log.Println("bullet proof test fail")
	} else {
		log.Println("past test")
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

func getSubmitValues(agg *contract.Agg)func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("get submit values")
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
	cipher,err:=zcrypto.PubKey.Encrypt(big.NewInt(amount))
	if err!=nil {
		log.Println(err.Error())
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	submitData,err:= json.Marshal("0x"+cipher.C.Text(16));
	if err!=nil {
		log.Println(err.Error())
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	w.Write(submitData)
}

func main() {
	//testPallier()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//testBulletProof()
	r := mux.NewRouter()

	manager:= appClient.NewHandlerManager()
	agg:= contract.NewAgg(contract.GETH_PORT,contract.CONTRACT_ABI,common.HexToAddress(contract.CONTRACT_ADDRESS))
	token:= contract.NewERC20(contract.GETH_PORT,contract.ERC20_ABI,common.HexToAddress(contract.ERC20_ADDRESS))

	go manager.Start()
	go manager.SubScriptContractEvent(agg)

	r.HandleFunc("/", requestParserWrapper(manager,agg,token)).Methods("GET")
	r.HandleFunc("/nonce/{user}",getNonce(agg)).Methods("GET")
	r.HandleFunc("/chainId",getChainId(agg)).Methods("GET")
	r.HandleFunc("/getSubmitValues/{taskId}",getSubmitValues(agg)).Methods("GET")
	r.HandleFunc("/encryptedData",getEncryptedData).Methods("POST")

	fmt.Println("Running http server")
	http.ListenAndServe(
		"0.0.0.0:4000",
		h.CORS(
			h.AllowedMethods([]string{"get", "options", "post", "put", "head"}),
			h.AllowedOrigins([]string{"*"}),
			h.AllowedHeaders([]string{"Content-Type"}),
		)(r),
	)
}
