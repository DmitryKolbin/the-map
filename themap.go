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

const (
	StateNew              = "New"
	StatePreAuthorized3DS = "PreAuthorized3DS"
	StateAuthorizing      = "Authorizing"
	StateAuthorized       = "Authorized"
	StateVoiding          = "Voiding"
	StateVoided           = "Voided"
	StateCharging         = "Charging"
	StateCharged          = "Charged"
	StateRefunding        = "Refunding"
	StateRefunded         = "Refunded"
	StateVerifying        = "Verifying"
	StatePayout           = "Payout"
	StatePaying           = "Paying"
	StatePaid             = "Paid"
	StateInProcess        = "InProcess"
	StateRejected         = "Rejected"

	SessionTypePay = "pay"
	SessionTypeAdd = "add"

	PaymentTypeOneStep = "OneStep"
	PaymentTypeTwoStep = "TwoStep"
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

type baseRequestKey struct {
	Key string `json:"key"`
}

func (r *baseRequestKey) setKey(key string) {
	r.Key = key
}

type baseResponse struct {
	Success    bool   `json:"Success"`
	ErrCode    string `json:"ErrCode"`
	ErrMessage string `json:"ErrMessage"`
}

func (r baseResponse) getSuccess() bool {
	return r.Success
}

func (r baseResponse) getErrCode() string {
	return r.ErrCode
}
func (r baseResponse) getErrMessage() string {
	return r.ErrMessage
}

type CardInfo struct {
	CardWithUid
	CardData
}

type CardWithUid struct {
	Uid string `json:"uid,omitempty"`
}

type CardData struct {
	Pan    string `json:"pan,omitempty"`
	EMonth int64  `json:"emonth,omitempty"`
	EYear  int64  `json:"eyear,omitempty"`
	Cvv    string `json:"cvv,omitempty"`
	Holder string `json:"holder,omitempty"`
}

type Credential struct {
	Login            string `json:"login"`
	Password         string `json:"password"`
	MerchantName     string `json:"merchant_name"`
	MerchantPassword string `json:"merchant_password"`
	TerminalPassword string `json:"terminal_password"`
}

type InitRequest struct {
	baseRequestKey

	AddCard         bool              `json:"add_card"`
	Type            string            `json:"type"`
	PaymentType     string            `json:"payment_type"`
	Lifetime        int               `json:"lifetime"`
	MerchantOrderId string            `json:"merchant_order_id"`
	Amount          int64             `json:"amount"`
	Credential      Credential        `json:"credential,omitempty"`
	CustomParamsRdy map[string]string `json:"custom_params_rdy,omitempty"`
	Recurrent       bool              `json:"recurrent"`
}

type InitResponse struct {
	baseResponse

	OrderId     string `json:"OrderId"`
	Amount      int64  `json:"Amount"`
	Type        string `json:"Type"`
	SessionGUID string `json:"SessionGUID"`
}

type BlockRequest struct {
	baseRequestKey
	Card            CardInfo          `json:"card"`
	MerchantOrderId string            `json:"merchant_order_id"`
	Amount          int64             `json:"amount"`
	Credential      Credential        `json:"credential,omitempty"`
	CustomParamsRdy map[string]string `json:"custom_params_rdy,omitempty"`
}

type BlockResponse struct {
	baseResponse

	OrderId    string `json:"OrderId"`
	Amount     int64  `json:"Amount"`
	ACSUrl     string `json:"ACSUrl"`
	PaReq      string `json:"PaReq"`
	ThreeDSKey string `json:"ThreeDSKey"`
	Is3DSVer1  bool   `json:"Is3DSVer1"`
}

type Block3DSRequest struct {
	baseRequestKey
	MerchantOrderId string `json:"merchant_order_id"`
	Pares           string `json:"pares"`
}

type Block3DSResponse struct {
	baseResponse

	OrderId string `json:"OrderId"`
	Amount  int64  `json:"Amount"`
}

type ChargeRequest struct {
	baseRequestKey
	MapOrderId string `json:"map_order_id"`
	Amount     int64  `json:"amount"`
}
type ChargeResponse struct {
	baseResponse

	OrderId string `json:"OrderId"`
	Key     string `json:"Key"`
	Amount  int64  `json:"Amount"`
}

type UnblockRequest struct {
	baseRequestKey

	MapOrderId string `json:"map_order_id"`
	Amount     int64  `json:"amount"`
}

type UnblockResponse struct {
	baseResponse

	OrderId   string `json:"OrderId"`
	NewAmount int64  `json:"NewAmount"`
}

type RefundRequest struct {
	baseRequestKey

	MapOrderId      string            `json:"map_order_id"`
	Amount          int64             `json:"amount"`
	CustomParamsRdy map[string]string `json:"custom_params_rdy"`
}

type RefundResponse struct {
	baseResponse

	OrderId   string `json:"OrderId"`
	NewAmount int64  `json:"NewAmount"`
}

type OrderStateRequest struct {
	baseRequestKey
	MerchantOrderId string `json:"merchant_order_id"`
}

type OrderStateResponse struct {
	baseResponse

	OrderId         string   `json:"OrderId"`
	Amount          int64    `json:"Amount"`
	State           string   `json:"State"`
	MerchantOrderId string   `json:"MerchantOrderId"`
	FeePercent      *float64 `json:"FeePercent"`
	CardType        string   `json:"CardType"`
	PanMask         string   `json:"PanMask"`
}

type StoreCardRequest struct {
	baseRequestKey
	MerchantOrderId string            `json:"merchant_order_id"`
	Amount          int64             `json:"amount"`
	Card            CardData          `json:"card"`
	Credential      Credential        `json:"Credential"`
	CustomParamsRdy map[string]string `json:"custom_params_rdy"`
}

type StoreCardResponse struct {
	baseResponse

	CardUId  string `json:"CardUId"`
	PANMask  string `json:"PANMask"`
	IsActive bool   `json:"IsActive"`
}

type RemoveCardRequest struct {
	baseRequestKey

	Credential Credential  `json:"credential"`
	Card       CardWithUid `json:"card"`
}

type RemoveCardResponse struct {
	baseResponse

	CardUId string `json:"CardUId"`
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
	return makeRequestGeneric[StoreCardResponse](*ms, http.MethodPost, "/storeCard ", &request, StoreCardResponse{})
}

func (ms *mapPaymentService) RemoveCard(request RemoveCardRequest) (RemoveCardResponse, error) {
	return makeRequestGeneric[RemoveCardResponse](*ms, http.MethodPost, "/removeCard ", &request, RemoveCardResponse{})
}

func (ms *mapPaymentService) Unblock(request UnblockRequest) (UnblockResponse, error) {
	return makeRequestGeneric[UnblockResponse](*ms, http.MethodPost, "/Unblock ", &request, UnblockResponse{})
}
func (ms *mapPaymentService) Refund(request RefundRequest) (RefundResponse, error) {
	return makeRequestGeneric[RefundResponse](*ms, http.MethodPost, "/Refund ", &request, RefundResponse{})
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
			return response, errors.New(response.getErrCode() + " " + response.getErrMessage())
		}

		return response, nil
	} else {
		return response, errors.New("response statusCode: " + strconv.Itoa(httpResponse.StatusCode))
	}
}
