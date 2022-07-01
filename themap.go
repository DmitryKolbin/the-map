package themap

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type MapPaymentInfo struct {
	Status         string
	OrderId        string
	PaymentOrderId string
	Amount         float64
}

type mapRequest interface {
	setKey(key string)
}

type mapResponse interface {
	getSuccess() bool
	getErrCode() string
	getErrMessage() string
}

type MapPaymentService interface {
	Init(request InitRequest) (InitResponse, error)
	CreatePaymentUrl(sessionId string) (string, error)
	Block(request BlockRequest) (BlockResponse, error)
	Block3DS(orderId string, pares string) (Block3DSResponse, error)
	Charge(orderId string, amount int64) (ChargeResponse, error)
	GetState(orderId string) (OrderStateResponse, error)
	Unblock(request UnblockRequest) (UnblockResponse, error)
	Refund(request RefundRequest) (RefundResponse, error)
	StoreCard(request StoreCardRequest) (StoreCardResponse, error)
	RemoveCard(request RemoveCardRequest) (RemoveCardResponse, error)
}

type mapPaymentService struct {
	url string
	key string
}

func NewMapPaymentService(url string, key string) MapPaymentService {
	return &mapPaymentService{
		url: url,
		key: key,
	}
}

func (ms *mapPaymentService) Init(request InitRequest) (InitResponse, error) {
	return makeRequestGeneric[InitResponse](*ms, http.MethodPost, "/Init", &request, InitResponse{})
}

func (ms *mapPaymentService) CreatePaymentUrl(sessionId string) (string, error) {
	u, err := url.Parse(ms.url)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, "createPayment")
	query := u.Query()
	query.Set("SessionID", sessionId)
	u.RawQuery = query.Encode()
	return u.String(), nil
}

func (ms *mapPaymentService) Block(request BlockRequest) (BlockResponse, error) {
	return makeRequestGeneric[BlockResponse](*ms, http.MethodPost, "/Block", &request, BlockResponse{})
}

func (ms *mapPaymentService) Block3DS(orderId string, pares string) (Block3DSResponse, error) {
	request := Block3DSRequest{
		MerchantOrderId: orderId,
		Pares:           pares,
	}
	return makeRequestGeneric[Block3DSResponse](*ms, http.MethodPost, "/Block3DS", &request, Block3DSResponse{})
}

func (ms *mapPaymentService) Charge(orderId string, amount int64) (ChargeResponse, error) {
	request := ChargeRequest{
		Amount:     amount,
		MapOrderId: orderId,
	}
	return makeRequestGeneric[ChargeResponse](*ms, http.MethodPost, "/Charge", &request, ChargeResponse{})
}

func (ms *mapPaymentService) GetState(orderId string) (OrderStateResponse, error) {
	request := OrderStateRequest{
		MerchantOrderId: orderId,
	}
	return makeRequestGeneric[OrderStateResponse](*ms, http.MethodPost, "/getState", &request, OrderStateResponse{})
}

func (ms *mapPaymentService) StoreCard(request StoreCardRequest) (StoreCardResponse, error) {
	return makeRequestGeneric[StoreCardResponse](*ms, http.MethodPost, "/storeCard", &request, StoreCardResponse{})
}

func (ms *mapPaymentService) RemoveCard(request RemoveCardRequest) (RemoveCardResponse, error) {
	return makeRequestGeneric[RemoveCardResponse](*ms, http.MethodPost, "/removeCard", &request, RemoveCardResponse{})
}

func (ms *mapPaymentService) Unblock(request UnblockRequest) (UnblockResponse, error) {
	return makeRequestGeneric[UnblockResponse](*ms, http.MethodPost, "/Unblock", &request, UnblockResponse{})
}
func (ms *mapPaymentService) Refund(request RefundRequest) (RefundResponse, error) {
	return makeRequestGeneric[RefundResponse](*ms, http.MethodPost, "/Refund", &request, RefundResponse{})
}

func makeRequestGeneric[res mapResponse](ms mapPaymentService, method, endpoint string, request mapRequest, response res) (res, error) {
	request.setKey(ms.key)

	data, err := json.Marshal(request)
	if err != nil {
		return response, err
	}

	httpRequest, _ := http.NewRequest(method, ms.url+endpoint, bytes.NewBuffer(data))
	httpRequest.Header.Add("Content-Type", "application/json")

	c := &http.Client{}

	httpResponse, err := c.Do(httpRequest)
	if err != nil {
		return response, err
	}

	defer httpResponse.Body.Close()
	body, _ := ioutil.ReadAll(httpResponse.Body)

	if httpResponse.StatusCode >= 200 && httpResponse.StatusCode < 300 {
		if err := json.Unmarshal(body, &response); err != nil {
			return response, err
		}
		if !response.getSuccess() {
			if len(response.getErrMessage()) == 0 {
				return response, errors.New(response.getErrCode())
			}
			return response, errors.New(response.getErrCode() + " " + response.getErrMessage())
		}

		return response, nil
	} else {
		return response, errors.New("response statusCode: " + strconv.Itoa(httpResponse.StatusCode))
	}
}
