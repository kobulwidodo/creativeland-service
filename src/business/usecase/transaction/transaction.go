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
	GetOrderDetail(ctx context.Context, param entity.TransactionParam) (entity.TransactionDetailResponse, error)
	GetTransactionListByUmkm(ctx context.Context, param entity.TransactionParam) ([]entity.TransactionDetailResponse, error)
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
		Status:  entity.StatusInCart,
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
		menuIDs = append(menuIDs, int64(c.MenuID))
	}

	menus, err := t.menu.GetListInByID(menuIDs)
	if err != nil {
		return 0, err
	}

	menusMap := make(map[int]entity.Menu)
	for _, m := range menus {
		menusMap[int(m.ID)] = m
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
		Status:  entity.StatusInCart,
		GuestID: user.User.GuestID,
	}, entity.UpdateCartParam{
		Status:        entity.StatusUnpaid,
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

func (t *transaction) convertToItemsDetails(carts []entity.Cart, menus map[int]entity.Menu) []midtrans.ItemsDetails {
	res := []midtrans.ItemsDetails{}
	for _, c := range carts {
		resTemp := midtrans.ItemsDetails{
			ID:    strconv.Itoa(int(c.ID)),
			Price: int64(c.PricePerItem),
			Qty:   c.Amount,
			Name:  menus[int(c.MenuID)].Name,
		}
		res = append(res, resTemp)
	}

	return res
}

func (t *transaction) GetOrderDetail(ctx context.Context, param entity.TransactionParam) (entity.TransactionDetailResponse, error) {
	result := entity.TransactionDetailResponse{}

	transaction, err := t.transaction.Get(param)
	if err != nil {
		return result, err
	}

	midtransTransaction, err := t.midtransTransaction.Get(entity.MidtransTransactionParam{
		TransactionID: transaction.ID,
	})
	if err != nil {
		return result, err
	}

	carts, err := t.cart.GetList(entity.CartParam{
		TransactionID: transaction.ID,
	})
	if err != nil {
		return result, err
	}

	menusID := []int64{}
	for _, c := range carts {
		menusID = append(menusID, int64(c.MenuID))
	}

	menus, err := t.menu.GetListInByID(menusID)
	if err != nil {
		return result, err
	}

	menuMap := make(map[uint]entity.Menu)
	for _, m := range menus {
		menuMap[m.ID] = m
	}

	result.ID = transaction.ID
	result.BuyerName = transaction.BuyerName
	result.Seat = transaction.Seat
	result.Notes = transaction.Notes
	result.Price = transaction.Price
	result.Status = midtransTransaction.Status
	result.PaymentType = midtransTransaction.GetPaymentType()
	result.ItemMenus = []entity.ItemMenu{}

	for _, c := range carts {
		itemMenu := entity.ItemMenu{
			Name:         menuMap[c.MenuID].Name,
			Price:        c.TotalPrice,
			Qty:          c.Amount,
			PricePerItem: c.PricePerItem,
		}
		result.ItemMenus = append(result.ItemMenus, itemMenu)
	}

	return result, nil
}

func (t *transaction) GetTransactionListByUmkm(ctx context.Context, param entity.TransactionParam) ([]entity.TransactionDetailResponse, error) {
	result := []entity.TransactionDetailResponse{}

	carts, err := t.cart.GetList(entity.CartParam{
		UmkmID: param.UmkmID,
		Status: param.Status,
	})
	if err != nil {
		return result, err
	}

	if len(carts) == 0 {
		return result, nil
	}

	cartsMap := make(map[uint][]entity.Cart)
	for _, c := range carts {
		cartsMap[c.TransactionID] = append(cartsMap[c.TransactionID], c)
	}

	menus, err := t.menu.GetAll(entity.MenuParam{
		UmkmID: param.UmkmID,
	})
	if err != nil {
		return result, err
	}
	menusMap := make(map[uint]entity.Menu)
	for _, m := range menus {
		menusMap[m.ID] = m
	}

	transactionIDs := []uint{}
	for k := range cartsMap {
		transactionIDs = append(transactionIDs, k)
	}

	transactions, err := t.transaction.GetListByIDs(transactionIDs)
	if err != nil {
		return result, err
	}

	for _, t := range transactions {
		transactionDetail := entity.TransactionDetailResponse{
			ID:        t.ID,
			BuyerName: t.BuyerName,
			Seat:      t.Seat,
			Notes:     t.Notes,
			Price:     t.Price,
			Status:    cartsMap[t.ID][0].Status,
		}
		itemMenus := []entity.ItemMenu{}
		for _, cm := range cartsMap[t.ID] {
			itemMenus = append(itemMenus, entity.ItemMenu{
				UmkmID:       cm.UmkmID,
				MenuID:       cm.MenuID,
				Name:         menusMap[cm.MenuID].Name,
				Price:        cm.TotalPrice,
				Qty:          cm.Amount,
				PricePerItem: cm.PricePerItem,
			})
		}
		transactionDetail.ItemMenus = itemMenus
		result = append(result, transactionDetail)
	}

	return result, nil
}
