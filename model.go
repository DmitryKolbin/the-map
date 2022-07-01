package themap

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

type ExtraFz54 struct {
	Cheque *ExtraFz54Cheque `json:"cheque,omitempty"`
	Goods  []ExtraFz54Good  `json:"goods"`
}
type ExtraFz54Cheque struct {
	AdditionalAttribute string `json:"additional_attribute,omitempty"`
	PenaltyAttribute    string `json:"penalty_attribute,omitempty"`
}

type ExtraFz54Good struct {
	Name               string             `json:"name"`
	Price              string             `json:"price"`
	Tax                *int               `json:"tax,omitempty"`
	PaymentSubjectType *int               `json:"payment_subject_type,omitempty"`
	PaymentMethodType  *int               `json:"payment_method_type,omitempty"`
	AgentType          *int               `json:"agent_type,omitempty"`
	Supplier           *ExtraFz54Supplier `json:"supplier,omitempty"`
}

type ExtraFz54Supplier struct {
	Name         string   `json:"name"`
	Inn          string   `json:"inn"`
	PhoneNumbers []string `json:"phone_numbers"`
}

type InitRequest struct {
	baseRequestKey
	*ExtraFz54

	AddCard         bool       `json:"add_card"`
	Type            string     `json:"type"`
	PaymentType     string     `json:"payment_type"`
	Lifetime        int        `json:"lifetime"`
	MerchantOrderId string     `json:"merchant_order_id"`
	Amount          int64      `json:"amount"`
	Credential      Credential `json:"credential,omitempty"`
	CustomParamsRaw string     `json:"custom_params_raw,omitempty"`
	Recurrent       bool       `json:"recurrent"`
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
	*ExtraFz54

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
	*ExtraFz54

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
	*ExtraFz54

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
