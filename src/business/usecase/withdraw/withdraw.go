package withdraw

import (
	"context"
	umkmDom "go-clean/src/business/domain/umkm"
	withdrawDom "go-clean/src/business/domain/withdraw"
	"go-clean/src/business/entity"
	"sort"
)

type Interface interface {
	Create(ctx context.Context, param entity.CreateWithdrawParam) (entity.Withdraw, error)
	GetList(ctx context.Context, param entity.WithdrawParam) ([]entity.Withdraw, error)
}

type withdraw struct {
	withdraw withdrawDom.Interface
	umkm     umkmDom.Interface
}

func Init(wd withdrawDom.Interface, ud umkmDom.Interface) Interface {
	w := &withdraw{
		withdraw: wd,
		umkm:     ud,
	}
	return w
}

func (w *withdraw) Create(ctx context.Context, param entity.CreateWithdrawParam) (entity.Withdraw, error) {
	wdExist, err := w.withdraw.Get(entity.WithdrawParam{
		Date:   param.Date,
		UmkmID: param.UmkmID,
	})
	if err == nil {
		if err := w.withdraw.Update(entity.WithdrawParam{
			Date:   param.Date,
			UmkmID: param.UmkmID,
		}, entity.UpdateWithdrawParam{
			Amount: wdExist.Amount + param.Amount,
		}); err != nil {
			return wdExist, err
		}
		return wdExist, nil
	}

	wd, err := w.withdraw.Create(entity.Withdraw{
		Date:   param.Date,
		Amount: param.Amount,
		UmkmID: param.UmkmID,
	})
	if err != nil {
		return wd, err
	}

	return wd, nil
}

func (w *withdraw) GetList(ctx context.Context, param entity.WithdrawParam) ([]entity.Withdraw, error) {
	wds, err := w.withdraw.GetList(entity.WithdrawParam{
		Date:    param.Date,
		UmkmID:  param.UmkmID,
		Limit:   param.Limit,
		Offset:  (param.Page - 1) * param.Limit,
		OrderBy: "date desc",
	})
	if err != nil {
		return wds, err
	}

	umkms, err := w.umkm.GetList(entity.UmkmParam{})
	if err != nil {
		return []entity.Withdraw{}, err
	}

	umkmsMap := make(map[uint]entity.Umkm)
	for _, u := range umkms {
		umkmsMap[u.ID] = u
	}

	for idx := range wds {
		wds[idx].UmkmName = umkmsMap[wds[idx].UmkmID].Name
	}

	sort.Slice(wds, func(i, j int) bool {
		return wds[i].Date > wds[j].Date
	})

	return wds, nil
}
