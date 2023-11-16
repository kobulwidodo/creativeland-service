package transaction

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	cartDom "go-clean/src/business/domain/cart"
	menuDom "go-clean/src/business/domain/menu"
	midtransDom "go-clean/src/business/domain/midtrans"
	midtransTransactionDom "go-clean/src/business/domain/midtrans_transaction"
	transactionDom "go-clean/src/business/domain/transaction"
	umkmDom "go-clean/src/business/domain/umkm"
	"go-clean/src/business/entity"
	"go-clean/src/lib/auth"
	"go-clean/src/lib/midtrans"
	"go-clean/src/lib/timeutils"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/xuri/excelize/v2"
)

type Interface interface {
	Create(ctx context.Context, param entity.CreateTransactionParam) (uint, error)
	GetOrderDetail(ctx context.Context, param entity.TransactionParam) (entity.TransactionDetailResponse, error)
	GetTransactionListByUmkm(ctx context.Context, param entity.TransactionParam) ([]entity.TransactionDetailResponse, error)
	GetTransactionList(ctx context.Context, param entity.TransactionParam) ([]entity.TransactionDetailResponse, error)
	GetMyTransaction(ctx context.Context, param entity.TransactionParam) ([]entity.TransactionDetailResponse, error)
	GetRecapSalesList(ctx context.Context, param entity.TransactionParam) ([]entity.SalesRecapResponse, error)
	GenerateExcel(ctx context.Context, param entity.TransactionParam) (*excelize.File, string, error)
	CompleteOrder(ctx context.Context, param entity.TransactionParam) error
	CancelOrder(ctx context.Context, param entity.TransactionParam) error
}

type transaction struct {
	transaction         transactionDom.Interface
	cart                cartDom.Interface
	auth                auth.Interface
	menu                menuDom.Interface
	umkm                umkmDom.Interface
	midtrans            midtransDom.Interface
	midtransTransaction midtransTransactionDom.Interface
}

func Init(auth auth.Interface, td transactionDom.Interface, cd cartDom.Interface, md menuDom.Interface, ud umkmDom.Interface, mtd midtransDom.Interface, mtt midtransTransactionDom.Interface) Interface {
	t := &transaction{
		transaction:         td,
		cart:                cd,
		auth:                auth,
		menu:                md,
		umkm:                ud,
		midtrans:            mtd,
		midtransTransaction: mtt,
	}

	return t
}

func (t *transaction) Create(ctx context.Context, param entity.CreateTransactionParam) (uint, error) {
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

	coreApiRes := &coreapi.ChargeResponse{}
	if param.PaymentID == midtrans.Cash {
		coreApiRes.TransactionID = "0"
		coreApiRes.OrderID = fmt.Sprintf("%s-%d-%d", "CL", transaction.ID, time.Now().Unix())
	} else {
		coreApiRes, err = t.midtrans.Create(midtrans.CreateOrderParam{
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
	} else if paymentId == midtrans.Cash {
		return paymentData, nil
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

	transaction, err := t.transaction.Get(entity.TransactionParam{
		ID: param.ID,
	})
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

	umkms, err := t.umkm.GetList(entity.UmkmParam{})
	if err != nil {
		return result, err
	}

	umkmMap := make(map[uint]entity.Umkm)
	for _, u := range umkms {
		umkmMap[u.ID] = u
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
			UmkmName:     umkmMap[c.UmkmID].Name,
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

	carts, err := t.cart.GetListInByStatus(param.Statuses, entity.CartParam{
		UmkmID: param.UmkmID,
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

	midtransTransactions, err := t.midtransTransaction.GetListByTrxIDs(transactionIDs, entity.MidtransTransactionParam{
		OrderIDLike: param.MidtransOrderID,
	})
	if err != nil {
		return result, err
	}

	midtransTransactionMap := make(map[uint]entity.MidtransTransaction)
	for _, mt := range midtransTransactions {
		midtransTransactionMap[mt.TransactionID] = mt
	}

	log.Println(len(midtransTransactionMap))

	for _, t := range transactions {
		if _, ok := midtransTransactionMap[t.ID]; ok {
			transactionDetail := entity.TransactionDetailResponse{
				ID:              t.ID,
				BuyerName:       t.BuyerName,
				Seat:            t.Seat,
				Notes:           t.Notes,
				Price:           t.Price,
				Status:          cartsMap[t.ID][0].Status,
				MidtransOrderID: midtransTransactionMap[t.ID].OrderID,
				CreatedAt:       timeutils.DiffForHumans(t.CreatedAt),
			}
			itemMenus := []entity.ItemMenu{}
			for _, cm := range cartsMap[t.ID] {
				itemMenus = append(itemMenus, entity.ItemMenu{
					Name:         menusMap[cm.MenuID].Name,
					Price:        cm.TotalPrice,
					Qty:          cm.Amount,
					PricePerItem: cm.PricePerItem,
				})
			}
			transactionDetail.ItemMenus = itemMenus
			result = append(result, transactionDetail)
		}
	}

	return result, nil
}

func (t *transaction) GetTransactionList(ctx context.Context, param entity.TransactionParam) ([]entity.TransactionDetailResponse, error) {
	result := []entity.TransactionDetailResponse{}

	midtransTransaction, err := t.midtransTransaction.GetList(entity.MidtransTransactionParam{
		OrderIDLike: param.MidtransOrderID,
		Limit:       param.Limit,
		Offset:      (param.Page - 1) * param.Limit,
		OrderBy:     "id desc",
	})
	if err != nil {
		return result, err
	}

	if len(midtransTransaction) == 0 {
		return result, errors.New("data not found")
	}

	midtransTransactionMap := make(map[uint]entity.MidtransTransaction)
	for _, mt := range midtransTransaction {
		midtransTransactionMap[mt.TransactionID] = mt
	}

	transactionIDs := []uint{}
	for _, mt := range midtransTransaction {
		transactionIDs = append(transactionIDs, mt.TransactionID)
	}

	transactions, err := t.transaction.GetListByIDs(transactionIDs)
	if err != nil {
		return result, err
	}

	carts, err := t.cart.GetListInByTransactionID(transactionIDs)
	if err != nil {
		return result, err
	}
	cartsMap := make(map[uint][]entity.Cart)
	for _, c := range carts {
		cartsMap[c.TransactionID] = append(cartsMap[c.TransactionID], c)
	}

	menuIDsMap := make(map[uint]bool)
	menuIDs := []int64{}
	for _, c := range carts {
		if _, ok := menuIDsMap[c.MenuID]; !ok {
			menuIDsMap[c.MenuID] = true
			menuIDs = append(menuIDs, int64(c.MenuID))
		}
	}

	menus, err := t.menu.GetListInByID(menuIDs)
	if err != nil {
		return result, err
	}

	menusMap := make(map[uint]entity.Menu)
	for _, m := range menus {
		menusMap[m.ID] = m
	}

	umkm, err := t.umkm.GetList(entity.UmkmParam{})
	if err != nil {
		return result, err
	}

	umkmsMap := make(map[uint]entity.Umkm)
	for _, u := range umkm {
		umkmsMap[u.ID] = u
	}

	for _, t := range transactions {
		transactionDetail := entity.TransactionDetailResponse{
			ID:              t.ID,
			BuyerName:       t.BuyerName,
			Seat:            t.Seat,
			Notes:           t.Notes,
			Price:           t.Price,
			Status:          midtransTransactionMap[t.ID].Status,
			MidtransOrderID: midtransTransactionMap[t.ID].OrderID,
		}
		itemMenus := []entity.ItemMenu{}
		for _, c := range cartsMap[t.ID] {
			itemMenus = append(itemMenus, entity.ItemMenu{
				UmkmName:     umkmsMap[c.UmkmID].Name,
				Name:         menusMap[c.MenuID].Name,
				Status:       c.Status,
				Price:        c.TotalPrice,
				Qty:          c.Amount,
				PricePerItem: c.PricePerItem,
			})
		}
		transactionDetail.ItemMenus = itemMenus
		result = append(result, transactionDetail)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID > result[j].ID
	})

	return result, nil
}

func (t *transaction) GetMyTransaction(ctx context.Context, param entity.TransactionParam) ([]entity.TransactionDetailResponse, error) {
	result := []entity.TransactionDetailResponse{}

	user, err := t.auth.GetUserAuthInfo(ctx)
	if err != nil {
		return result, err
	}

	carts, err := t.cart.GetListInByStatus([]string{entity.StatusPaid, entity.StatusUnpaid, entity.StatusDone}, entity.CartParam{
		GuestID: user.User.GuestID,
	})
	if err != nil {
		return result, err
	}

	if len(carts) == 0 {
		return result, nil
	}

	cartsMap := make(map[uint][]entity.Cart)
	menusMap := make(map[uint]entity.Menu)
	umkmsMap := make(map[uint]entity.Umkm)
	for _, c := range carts {
		cartsMap[c.TransactionID] = append(cartsMap[c.TransactionID], c)
		menusMap[c.MenuID] = entity.Menu{}
		umkmsMap[c.UmkmID] = entity.Umkm{}
	}

	menuIDs := []int64{}
	for k := range menusMap {
		menuIDs = append(menuIDs, int64(k))
	}
	menus, err := t.menu.GetListInByID(menuIDs)
	if err != nil {
		return result, err
	}
	for _, m := range menus {
		menusMap[m.ID] = m
	}

	umkmIDs := []uint{}
	for k := range umkmsMap {
		umkmIDs = append(umkmIDs, k)
	}
	umkms, err := t.umkm.GetListInByID(umkmIDs)
	if err != nil {
		return result, err
	}
	for _, u := range umkms {
		umkmsMap[u.ID] = u
	}

	transactionIDs := []uint{}
	for k := range cartsMap {
		transactionIDs = append(transactionIDs, k)
	}

	sort.Slice(transactionIDs, func(i, j int) bool {
		return transactionIDs[i] > transactionIDs[j]
	})

	if len(transactionIDs) > 10 {
		transactionIDs = transactionIDs[:10]
	}

	transactions, err := t.transaction.GetListByIDs(transactionIDs)
	if err != nil {
		return result, err
	}

	midtransTransactions, err := t.midtransTransaction.GetListByTrxIDs(transactionIDs, entity.MidtransTransactionParam{
		OrderIDLike: param.MidtransOrderID,
	})
	if err != nil {
		return result, err
	}

	midtransTransactionMap := make(map[uint]entity.MidtransTransaction)
	for _, mt := range midtransTransactions {
		midtransTransactionMap[mt.TransactionID] = mt
	}

	for _, t := range transactions {
		if _, ok := midtransTransactionMap[t.ID]; ok {
			mt := midtransTransactionMap[t.ID]
			transactionDetail := entity.TransactionDetailResponse{
				ID:              t.ID,
				BuyerName:       t.BuyerName,
				Seat:            t.Seat,
				Notes:           t.Notes,
				Price:           t.Price,
				Status:          mt.Status,
				MidtransOrderID: mt.OrderID,
				PaymentType:     mt.GetPaymentType(),
				CreatedAt:       timeutils.DiffForHumans(t.CreatedAt),
			}
			itemMenus := []entity.ItemMenu{}
			for _, cm := range cartsMap[t.ID] {
				itemMenus = append(itemMenus, entity.ItemMenu{
					UmkmName:     umkmsMap[cm.UmkmID].Name,
					Name:         menusMap[cm.MenuID].Name,
					Status:       cm.Status,
					Price:        cm.TotalPrice,
					Qty:          cm.Amount,
					PricePerItem: cm.PricePerItem,
				})
			}
			if transactionDetail.Status == entity.StatusUnpaid {
				paymentData := entity.PaymentData{}
				if err := json.Unmarshal([]byte(midtransTransactionMap[t.ID].PaymentData), &paymentData); err != nil {
					log.Println("failed to un marshal payment data")
					continue
				}
				transactionDetail.PaymentData = paymentData
			}
			transactionDetail.ItemMenus = itemMenus
			result = append(result, transactionDetail)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID > result[j].ID
	})

	return result, nil
}

func (t *transaction) GenerateExcel(ctx context.Context, param entity.TransactionParam) (*excelize.File, string, error) {
	recap := []entity.SalesRecapResponse{}

	carts, err := t.cart.GetList(entity.CartParam{
		Status:    entity.StatusDone,
		CreatedAt: param.Date,
	})
	if err != nil {
		return nil, "", err
	}

	if len(carts) == 0 {
		return nil, "", err
	}

	layout := "2006-01"
	timeFormat, err := time.Parse(layout, param.Date)
	if err != nil {
		return nil, "", err
	}

	startDate := time.Date(timeFormat.Year(), timeFormat.Month(), 1, 0, 0, 0, 0, timeFormat.Location())
	endDate := startDate.AddDate(0, 1, -1)

	dateTrxMap := make(map[string]entity.SalesRecapResponse)
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateFormatted := d.Format("2006-01-02")
		dateTrxMap[dateFormatted] = entity.SalesRecapResponse{
			Date:        dateFormatted,
			GrossAmount: 0,
			NetAmount:   0,
		}
	}

	for _, c := range carts {
		dateFormetted := c.CreatedAt.Format("2006-01-02")
		dateTrxMap[dateFormetted] = entity.SalesRecapResponse{
			Date:        dateFormetted,
			GrossAmount: c.TotalPrice,
			NetAmount:   c.TotalPrice * 17 / 100,
		}
	}

	for _, v := range dateTrxMap {
		recap = append(recap, v)
	}

	sort.Slice(recap, func(i, j int) bool {
		return recap[i].Date < recap[j].Date
	})

	f := excelize.NewFile()
	sheetname := fmt.Sprintf("%s-%d", param.Date, time.Now().Unix())
	index, err := f.NewSheet(sheetname)
	if err != nil {
		return nil, "", err
	}

	f.SetCellValue(sheetname, "A1", fmt.Sprintf("Recap Transaksi Bulan : %s", param.Date))
	f.SetCellValue(sheetname, "A2", "Tanggal")
	f.SetCellValue(sheetname, "B2", "Pendapatan Kotor")
	f.SetCellValue(sheetname, "C2", "Pendapatan Bersih")

	cellIndex := 0
	sumNet := 0
	sumGross := 0
	for i, r := range recap {
		cellIndex = i + 3
		f.SetCellValue(sheetname, fmt.Sprintf("A%d", cellIndex), r.Date)
		f.SetCellValue(sheetname, fmt.Sprintf("B%d", cellIndex), r.GrossAmount)
		f.SetCellValue(sheetname, fmt.Sprintf("C%d", cellIndex), r.NetAmount)
		sumNet += r.NetAmount
		sumGross += r.GrossAmount
	}

	f.SetCellValue(sheetname, fmt.Sprintf("A%d", cellIndex), "TOTAL")
	f.SetCellValue(sheetname, fmt.Sprintf("B%d", cellIndex), sumGross)
	f.SetCellValue(sheetname, fmt.Sprintf("C%d", cellIndex), sumNet)

	f.SetActiveSheet(index)

	return f, sheetname, nil
}

func (t *transaction) GetRecapSalesList(ctx context.Context, param entity.TransactionParam) ([]entity.SalesRecapResponse, error) {
	result := []entity.SalesRecapResponse{}

	now := time.Now()
	calcLastSevenDays := now.AddDate(0, 0, -7)

	carts, err := t.cart.GetList(entity.CartParam{
		Status:            entity.StatusDone,
		CreatedAt:         param.Date,
		CreatedAtMoreThan: calcLastSevenDays,
	})
	if err != nil {
		return result, err
	}

	umkms, err := t.umkm.GetList(entity.UmkmParam{})
	if err != nil {
		return result, err
	}

	umkmsMap := make(map[uint]entity.Umkm)
	for _, u := range umkms {
		umkmsMap[u.ID] = u
	}

	dateTrxMap := make(map[string]entity.SalesRecapResponse)

	var startDate, endDate time.Time

	if param.Date != "" {
		// Only one date is specified.
		parsedDate, err := time.Parse("2006-01-02", param.Date)
		if err != nil {
			return nil, err // handle error
		}
		startDate = parsedDate
		endDate = parsedDate
	} else {
		// Last seven days range is specified.
		startDate = calcLastSevenDays
		endDate = now
	}

	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateFormatted := d.Format("2006-01-02")
		dateTrxMap[dateFormatted] = entity.SalesRecapResponse{
			Date:        dateFormatted,
			GrossAmount: 0,
			NetAmount:   0,
			UmkmDetail:  []entity.UmkmDetailRecap{},
		}
	}

	umkmRecapMap := make(map[entity.KeyUmkmDetailRecap]entity.UmkmDetailRecap)
	for _, c := range carts {
		dateFormatted := c.CreatedAt.Format("2006-01-02")

		// Update sales recap for the date
		recap := dateTrxMap[dateFormatted]
		recap.NetAmount += c.TotalPrice * 17 / 100
		recap.GrossAmount += c.TotalPrice
		dateTrxMap[dateFormatted] = recap

		// Handle Umkm details
		key := entity.KeyUmkmDetailRecap{
			ID:          c.UmkmID,
			CreatedDate: dateFormatted,
		}
		if uRecap, ok := umkmRecapMap[key]; ok {
			uRecap.GrossAmount += c.TotalPrice
			uRecap.NetAmount += c.TotalPrice * 83 / 100
			uRecap.TotalOrder += c.Amount
			umkmRecapMap[key] = uRecap
		} else {
			umkmRecapMap[key] = entity.UmkmDetailRecap{
				ID:          c.UmkmID,
				UmkmName:    umkmsMap[c.UmkmID].Name,
				GrossAmount: c.TotalPrice,
				NetAmount:   c.TotalPrice * 83 / 100,
				TotalOrder:  c.Amount,
			}
		}
	}

	for key, recap := range umkmRecapMap {
		dateFormetted := key.CreatedDate
		if dateTrx, ok := dateTrxMap[dateFormetted]; ok {
			dateTrx.UmkmDetail = append(dateTrx.UmkmDetail, recap)
			dateTrxMap[dateFormetted] = dateTrx
		}
	}

	for _, sr := range dateTrxMap {
		result = append(result, sr)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Date > result[j].Date
	})

	return result, nil
}

func (t *transaction) CompleteOrder(ctx context.Context, param entity.TransactionParam) error {
	carts, err := t.cart.GetList(entity.CartParam{
		TransactionID: param.ID,
		UmkmID:        param.UmkmID,
		Status:        entity.StatusPaid,
	})
	if err != nil {
		return err
	}

	cartsID := []uint{}
	for _, c := range carts {
		cartsID = append(cartsID, c.ID)
	}

	if err := t.cart.UpdatesByIDs(cartsID, entity.UpdateCartParam{
		Status: entity.StatusDone,
	}); err != nil {
		return err
	}

	return nil
}

func (t *transaction) CancelOrder(ctx context.Context, param entity.TransactionParam) error {
	carts, err := t.cart.GetList(entity.CartParam{
		TransactionID: param.ID,
		UmkmID:        param.UmkmID,
	})
	if err != nil {
		return err
	}

	cartsID := []uint{}
	for _, c := range carts {
		cartsID = append(cartsID, c.ID)
	}

	if err := t.cart.UpdatesByIDs(cartsID, entity.UpdateCartParam{
		Status: entity.StatusCancel,
	}); err != nil {
		return err
	}

	return nil
}
