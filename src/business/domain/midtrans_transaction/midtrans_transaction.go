package midtranstransaction

import (
	"go-clean/src/business/entity"

	"gorm.io/gorm"
)

type Interface interface {
	Create(midtransTransaction entity.MidtransTransaction) (entity.MidtransTransaction, error)
	Get(param entity.MidtransTransactionParam) (entity.MidtransTransaction, error)
	Update(selectParam entity.MidtransTransaction, updateParam entity.UpdateMidtransTransactionParam) error
}

type midtransTransaction struct {
	db *gorm.DB
}

func Init(db *gorm.DB) Interface {
	mt := &midtransTransaction{
		db: db,
	}

	return mt
}

func (mt *midtransTransaction) Create(midtransTransaction entity.MidtransTransaction) (entity.MidtransTransaction, error) {
	if err := mt.db.Create(&midtransTransaction).Error; err != nil {
		return midtransTransaction, err
	}

	return midtransTransaction, nil
}

func (mt *midtransTransaction) Get(param entity.MidtransTransactionParam) (entity.MidtransTransaction, error) {
	res := entity.MidtransTransaction{}
	if err := mt.db.Where(param).First(&res).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (mt *midtransTransaction) Update(selectParam entity.MidtransTransaction, updateParam entity.UpdateMidtransTransactionParam) error {
	if err := mt.db.Model(&selectParam).Updates(entity.MidtransTransaction{
		Status: updateParam.Status,
	}).Error; err != nil {
		return err
	}

	return nil
}
