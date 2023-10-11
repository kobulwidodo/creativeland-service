package withdraw

import (
	"go-clean/src/business/entity"

	"gorm.io/gorm"
)

type Interface interface {
	Create(withdraw entity.Withdraw) (entity.Withdraw, error)
	Get(param entity.WithdrawParam) (entity.Withdraw, error)
	GetList(param entity.WithdrawParam) ([]entity.Withdraw, error)
	Update(selectParam entity.WithdrawParam, updateParam entity.UpdateWithdrawParam) error
}

type withdraw struct {
	db *gorm.DB
}

func Init(db *gorm.DB) Interface {
	w := &withdraw{
		db: db,
	}
	return w
}

func (w *withdraw) Create(withdraw entity.Withdraw) (entity.Withdraw, error) {
	if err := w.db.Create(&withdraw).Error; err != nil {
		return withdraw, err
	}

	return withdraw, nil
}

func (w *withdraw) Get(param entity.WithdrawParam) (entity.Withdraw, error) {
	withdraw := entity.Withdraw{}

	if err := w.db.Where(param).First(&withdraw).Error; err != nil {
		return withdraw, err
	}

	return withdraw, nil
}

func (w *withdraw) GetList(param entity.WithdrawParam) ([]entity.Withdraw, error) {
	withdraws := []entity.Withdraw{}

	if err := w.db.Where(param).Order(param.OrderBy).Limit(param.Limit).Offset(param.Offset).Find(&withdraws).Error; err != nil {
		return withdraws, err
	}

	return withdraws, nil
}

func (w *withdraw) Update(selectParam entity.WithdrawParam, updateParam entity.UpdateWithdrawParam) error {
	if err := w.db.Model(entity.Withdraw{}).Where(selectParam).Updates(updateParam).Error; err != nil {
		return err
	}

	return nil
}
