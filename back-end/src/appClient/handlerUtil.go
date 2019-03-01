package appClient

import (
	"crypto/ecdsa"
	"encoding/json"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/websocket"
	"log"
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