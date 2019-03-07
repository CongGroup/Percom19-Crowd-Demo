package appClient

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/websocket"
	"log"
	"math/big"
	"zcrypto"
)

func (h *Handler) writeMessage(messageType int, data []byte) error {
	h.mutex.Lock()
	err:=h.w.WriteMessage(messageType, data)
	defer h.mutex.Unlock()
	return err
}

func (h *Handler) Send(data []byte) error {
	h.mutex.Lock()
	err:=h.w.WriteMessage(websocket.TextMessage,data)
	defer h.mutex.Unlock()
	return err
}

func (h *Handler) wrapperAndSend(gcuid int, res interface{}) error {
	resWrapper, err := json.Marshal(res)
	if err != nil {
		log.Println(err.Error())
		h.errorHandler(gcuid,UNMARSHAL_JSON_ERROR,err)
		return err
	}
	err =h.writeMessage(websocket.TextMessage, resWrapper)
	return err
}

func (h *Handler) sendError(gcuid int, errorCode int, reason string) error{
	reqError := &Error{
		Status: FAIL,
		Gcuid:  gcuid,
		Code: errorCode,
		Reason: reason,
	}
	reqErrorWrapper, err := json.Marshal(reqError)
	if err != nil {
		log.Println(err.Error())
		h.writeMessage(websocket.TextMessage, []byte(MSG_UNMARSHAL_JSON_ERROR))
		return err
	}
	err =h.writeMessage(websocket.TextMessage, reqErrorWrapper)
	return err
}


func derivePrivateKey(rawString string) (*ecdsa.PrivateKey,error) {
	pk,err:= crypto.HexToECDSA(rawString)
	if err!=nil {
		return nil,err
	}
	return pk,nil
}

func (h *Handler) errorHandler(gcuid int, errorId int, err error) {
	switch errorId {
	case UNMARSHAL_JSON_ERROR:
		h.sendError(gcuid,CLIENT_ERROR_CODE,MSG_UNMARSHAL_JSON_ERROR)
		break;
	case DATA_FORMAT_ERROR:
		h.sendError(gcuid,CLIENT_ERROR_CODE,MSG_DATA_FORMAT_ERROR)
		break;
	case KEY_FORMAT_ERROR:
		h.sendError(gcuid,CLIENT_ERROR_CODE,MSG_KEY_FORMAT_ERROR)
		break;
	case TRANSACTION_ERROR:
		h.sendError(gcuid,CLIENT_ERROR_CODE,MSG_TRANSACTION_ERROR)
		break;
	case CALL_TRANSACTION_ERROR:
		h.sendError(gcuid,CLIENT_ERROR_CODE,MSG_CALL_TRANSACTION_ERROR)
		break;
	case ENCRYPTION_ERROR:
		h.sendError(gcuid,CLIENT_ERROR_CODE,MSG_ENCRYPTION_ERROR)
		break;
	}
	return
}

func PackSubmitPayload(amount *big.Int) ([]byte,error) {
	var submitProof string

	if amount.Cmp(big.NewInt(MIN_RANGE))==-1 || amount.Cmp(big.NewInt(MAX_RANGE))!=-1 {
		log.Println("out of range, generate random proof");
		data:= make([]byte,1248)
		submitProof="0x"+hex.EncodeToString(data)
	} else {
		rpV := zcrypto.RPProve(amount);
		submitProof="0x"+hex.EncodeToString(rpV.Bytes())
	}

	cipher,err:=zcrypto.PubKey.Encrypt(amount)
	if err!=nil {
		log.Println(err.Error())
		return nil,err
	}

	submitData := "0x"+cipher.C.Text(16)
	//log.Println("proof:",submitProof)

	submitPayload,err := json.Marshal(&SubmitPayload{
		SubmitData:submitData,
		SubmitProof:submitProof,
	})
	if err!=nil {
		log.Println(err.Error())
		return nil,err
	}
	return submitPayload,nil
}

func GenBulletProof(amount *big.Int) ([]byte) {
	var proof []byte
	if amount.Cmp(big.NewInt(MIN_RANGE))==-1 || amount.Cmp(big.NewInt(MAX_RANGE))!=-1 {
		log.Println("out of range, generate random proof");
		proof = make([]byte,1248)
	} else {
		rpV := zcrypto.RPProve(amount);
		proof = rpV.Bytes()
	}
	return proof
}

func GenEncryption(amount *big.Int) (*big.Int,error) {
	cipher,err:=zcrypto.PubKey.Encrypt(amount)
	return cipher.C,err
}


func SetBigIntBytes(v * big.Int)  []byte{
	d:=v.Bytes()
	if len(d)==32 {
		return d
	}
	pad:=make([]byte,32-len(d))
	pad = append(pad,d...)
	return pad
}