package main

import (
	"appClient"
	"contract"
	"fmt"
	"github.com/deroproject/derosuite/crypto"
	"github.com/deroproject/derosuite/crypto/ringct"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	h "github.com/gorilla/handlers"
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

func requestParserWrapper(manager *appClient.HandlerManager, agg *contract.Agg) func(w http.ResponseWriter, r *http.Request){
	return func (w http.ResponseWriter, r *http.Request){
		log.Println("receive a reqeust")
		var upgrader = websocket.Upgrader{}
		upgrader.CheckOrigin = func(rq *http.Request) bool { return true }
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}
		log.Println("connection,:",&c);
		handler:= appHandler(c,agg)
		log.Println("handler:",handler)
		manager.Register <- handler
		handler.HandleRequest()
		manager.Unregister <- handler
		defer c.Close()
	}
}

//func requestParser (w http.ResponseWriter, r *http.Request){
//	log.Println("receive a reqeust")
//	var upgrader = websocket.Upgrader{}
//	upgrader.CheckOrigin = func(rq *http.Request) bool { return true }
//	c, err := upgrader.Upgrade(w, r, nil)
//	if err != nil {
//		log.Println("upgrade:", err)
//		return
//	}
//
//	appHandler(c).HandleRequest()
//	defer c.Close()
//}

func appHandler(c *websocket.Conn, agg *contract.Agg) *appClient.Handler {
	return appClient.NewHandler(c,agg)
}

func main() {
	//testPallier()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//testBulletProof()
	r := mux.NewRouter()

	manager:= appClient.NewHandlerManager()
	agg:= contract.NewAgg(contract.GETH_PORT,contract.CONTRACT_ABI,common.HexToAddress(contract.CONTRACT_ADDRESS))
	go manager.Start()
	go manager.SubScriptContractEvent(agg)

	r.HandleFunc("/", requestParserWrapper(manager,agg)).Methods("GET")

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
