package domain

import (
	"go-clean/src/business/domain/cart"
	"go-clean/src/business/domain/menu"
	"go-clean/src/business/domain/midtrans"
	midtranstransaction "go-clean/src/business/domain/midtrans_transaction"
	"go-clean/src/business/domain/transaction"
	"go-clean/src/business/domain/umkm"
	"go-clean/src/business/domain/user"
	midtransSdk "go-clean/src/lib/midtrans"

	"gorm.io/gorm"
)

type Domains struct {
	User                user.Interface
	Umkm                umkm.Interface
	Menu                menu.Interface
	Cart                cart.Interface
	Transaction         transaction.Interface
	Midtrans            midtrans.Interface
	MidtransTransaction midtranstransaction.Interface
}

func Init(db *gorm.DB, m midtransSdk.Interface) *Domains {
	d := &Domains{
		User:                user.Init(db),
		Umkm:                umkm.Init(db),
		Menu:                menu.Init(db),
		Cart:                cart.Init(db),
		Transaction:         transaction.Init(db),
		Midtrans:            midtrans.Init(m),
		MidtransTransaction: midtranstransaction.Init(db),
	}

	return d
}
