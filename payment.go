package qvo

//Payment struct to represent a qvo payment object.
type Payment struct {
	Amount        int64  `json:"amount"`
	Gateway       string `json:"gateway"`      //One of: webpay_plus, webpay_oneclick, olpays.
	PaymentType   string `json:"payment_type"` //credit or debit.
	Fee           int64  `json:"fee"`
	Installments  int32  `json:"installments"`
	PaymentMethod Card   `json:"payment_method"`
}
