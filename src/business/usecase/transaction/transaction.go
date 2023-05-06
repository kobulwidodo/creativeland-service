package transaction

import (
	"context"
	"encoding/json"
	"errors"
	cartDom "go-clean/src/business/domain/cart"
	menuDom "go-clean/src/business/domain/menu"
	midtransDom "go-clean/src/business/domain/midtrans"
	midtransTransactionDom "go-clean/src/business/domain/midtrans_transaction"
	transactionDom "go-clean/src/business/domain/transaction"
	"go-clean/src/business/entity"
	"go-clean/src/lib/auth"
	"go-clean/src/lib/midtrans"
	"strconv"

	"github.com/midtrans/midtrans-go/coreapi"
)

type Interface interface {
	Crete(ctx context.Context, param entity.CreateTransactionParam) (uint, error)
}

type transaction struct {
	transaction         transactionDom.Interface
	cart                cartDom.Interface
	auth                auth.Interface
	menu                menuDom.Interface
	midtrans            midtransDom.Interface
	midtransTransaction midtransTransactionDom.Interface
}

func Init(auth auth.Interface, td transactionDom.Interface, cd cartDom.Interface, md menuDom.Interface, mtd midtransDom.Interface, mtt midtransTransactionDom.Interface) Interface {
	t := &transaction{
		transaction:         td,
		cart:                cd,
		auth:                auth,
		menu:                md,
		midtrans:            mtd,
		midtransTransaction: mtt,
	}

	return t
}

func (t *transaction) Crete(ctx context.Context, param entity.CreateTransactionParam) (uint, error) {
	user, err := t.auth.GetUserAuthInfo(ctx)
	if err != nil {
		return 0, err
	}

	carts, err := t.cart.GetList(entity.CartParam{
		Status:  entity.StatusActive,
		GuestID: user.User.GuestID,
	})
	if err != nil {
		return 0, err
	}

	if len(carts) == 0 {
		return 0, errors.New("cart empty")
	}

	menuIDs := []int64{}
	for _, c := range carts {
		menuIDs = append(menuIDs, int64(c.ID))
	}

	menus, err := t.menu.GetListInByID(menuIDs)
	if err != nil {
		return 0, err
	}

	menusMap := make(map[uint]entity.Menu)
	for _, m := range menus {
		menusMap[m.ID] = m
	}

	grossAmount := 0
	for _, cart := range carts {
		grossAmount += cart.TotalPrice
	}

	transaction, err := t.transaction.Create(entity.Transaction{
		GuestID:   user.User.GuestID,
		BuyerName: param.BuyerName,
		Seat:      param.Seat,
		Notes:     param.Notes,
		Price:     grossAmount,
	})
	if err != nil {
		return 0, err
	}

	coreApiRes, err := t.midtrans.Create(midtrans.CreateOrderParam{
		OrderID:      transaction.ID,
		PaymentID:    param.PaymentID,
		GrossAmount:  int64(grossAmount),
		ItemsDetails: t.convertToItemsDetails(carts, menusMap),
		CustomerDetails: midtrans.CustomerDetails{
			Name:  param.BuyerName,
			Email: param.Email,
		},
	})
	if err != nil {
		return 0, err
	}

	paymentData, err := t.getPaymentData(param.PaymentID, coreApiRes)
	if err != nil {
		return 0, err
	}

	paymenDataMarshal, err := json.Marshal(paymentData)
	if err != nil {
		return 0, err
	}

	if err := t.cart.Update(entity.CartParam{
		Status:  entity.StatusActive,
		GuestID: user.User.GuestID,
	}, entity.UpdateCartParam{
		Status:        entity.StatusInactive,
		TransactionID: transaction.ID,
	}); err != nil {
		return 0, err
	}

	_, err = t.midtransTransaction.Create(entity.MidtransTransaction{
		TransactionID: transaction.ID,
		MidtransID:    coreApiRes.TransactionID,
		OrderID:       coreApiRes.OrderID,
		PaymentType:   param.PaymentID,
		GrossAmount:   grossAmount,
		Status:        "pending",
		PaymentData:   string(paymenDataMarshal),
	})
	if err != nil {
		return 0, err
	}

	return transaction.ID, nil
}

func (t *transaction) getPaymentData(paymentId int, coreApiRes *coreapi.ChargeResponse) (entity.PaymentData, error) {
	paymentData := entity.PaymentData{}
	if paymentId == midtrans.GopayPayment {
		paymentData.Key = coreApiRes.Actions[1].URL
		paymentData.Qr = coreApiRes.Actions[0].URL
	} else {
		return paymentData, errors.New("failed to get payment data")
	}

	return paymentData, nil
}

func (t *transaction) convertToItemsDetails(carts []entity.Cart, menus map[uint]entity.Menu) []midtrans.ItemsDetails {
	res := []midtrans.ItemsDetails{}
	for _, c := range carts {
		resTemp := midtrans.ItemsDetails{
			ID:    strconv.Itoa(int(c.ID)),
			Price: int64(c.PricePerItem),
			Qty:   c.Amount,
			Name:  menus[c.MenuID].Name,
		}
		res = append(res, resTemp)
	}

	return res
}
