package usecase

import (
	"go-clean/src/business/domain"
	analytic "go-clean/src/business/usecase/Analytic"
	"go-clean/src/business/usecase/cart"
	"go-clean/src/business/usecase/menu"
	midtranstransaction "go-clean/src/business/usecase/midtrans_transaction"
	"go-clean/src/business/usecase/transaction"
	"go-clean/src/business/usecase/umkm"
	"go-clean/src/business/usecase/user"
	"go-clean/src/lib/auth"
)

type Usecase struct {
	User                user.Interface
	Umkm                umkm.Interface
	Menu                menu.Interface
	Cart                cart.Interface
	Transaction         transaction.Interface
	MidtransTransaction midtranstransaction.Interface
	Analytic            analytic.Interface
}

func Init(auth auth.Interface, d *domain.Domains) *Usecase {
	uc := &Usecase{
		User:                user.Init(d.User, auth, d.Cart),
		Umkm:                umkm.Init(d.Umkm),
		Menu:                menu.Init(d.Menu),
		Cart:                cart.Init(d.Cart, auth, d.Menu, d.Umkm),
		Transaction:         transaction.Init(auth, d.Transaction, d.Cart, d.Menu, d.Umkm, d.Midtrans, d.MidtransTransaction),
		MidtransTransaction: midtranstransaction.Init(d.MidtransTransaction, d.Midtrans, d.Cart),
		Analytic:            analytic.Init(d.Cart),
	}

	return uc
}
