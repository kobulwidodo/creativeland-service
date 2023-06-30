package entity

import (
	"go-clean/src/lib/midtrans"

	"gorm.io/gorm"
)

const (
	StatusChallange = "challange"
	StatusSuccess   = "success"
	StatusDeny      = "deny"
	StatusFailure   = "failure"
	StatusPending   = "pending"
)

type MidtransTransaction struct {
	gorm.Model
	TransactionID uint
	MidtransID    string
	OrderID       string
	PaymentType   int
	GrossAmount   int
	Status        string
	PaymentData   string
}

type PaymentData struct {
	Key string `json:"key"`
	Qr  string `json:"qr"`
}

type MidtransTransactionParam struct {
	ID            uint   `json:"id"`
	TransactionID uint   `uri:"transaction_id" json:"transaction_id"`
	OrderID       string `json:"order_id"`
}

type MidtransTransactionPaymentDetail struct {
	Status      string      `json:"status"`
	MidtransID  string      `json:"midtrans_id"`
	PaymentData PaymentData `json:"payment_data"`
}

type UpdateMidtransTransactionParam struct {
	Status string `json:"string"`
}

func (mt *MidtransTransaction) GetPaymentType() string {
	result := ""

	switch mt.PaymentType {
	case midtrans.GopayPayment:
		result = "Gopay"
	}

	return result
}
