package midtranstransaction

import (
	"fmt"
	"go-clean/src/business/entity"

	"gorm.io/gorm"
)

type Interface interface {
	Create(midtransTransaction entity.MidtransTransaction) (entity.MidtransTransaction, error)
	Get(param entity.MidtransTransactionParam) (entity.MidtransTransaction, error)
	GetList(param entity.MidtransTransactionParam) ([]entity.MidtransTransaction, error)
	GetListByTrxIDs(ids []uint, param entity.MidtransTransactionParam) ([]entity.MidtransTransaction, error)
	Update(selectParam entity.MidtransTransactionParam, updateParam entity.UpdateMidtransTransactionParam) error
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

func (mt *midtransTransaction) GetList(param entity.MidtransTransactionParam) ([]entity.MidtransTransaction, error) {
	res := []entity.MidtransTransaction{}
	query := mt.db.Where(param).Where("order_id LIKE ?", fmt.Sprintf("%%%s%%", param.OrderIDLike))

	if param.CreatedAt != "" {
		query = query.Where("created_at LIKE ?", fmt.Sprintf("%%%s%%", param.CreatedAt))
	} else if !param.CreatedAtMoreThan.IsZero() {
		query = query.Where("created_at > ?", param.CreatedAtMoreThan)
	}

	if err := query.Limit(param.Limit).Offset(param.Offset).Order(param.OrderBy).Find(&res).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (mt *midtransTransaction) GetListByTrxIDs(ids []uint, param entity.MidtransTransactionParam) ([]entity.MidtransTransaction, error) {
	res := []entity.MidtransTransaction{}
	if err := mt.db.Where("transaction_id IN ? AND order_id LIKE ?", ids, fmt.Sprintf("%%%s%%", param.OrderIDLike)).Find(&res).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (mt *midtransTransaction) Update(selectParam entity.MidtransTransactionParam, updateParam entity.UpdateMidtransTransactionParam) error {
	if err := mt.db.Model(entity.MidtransTransaction{}).Where(selectParam).Updates(updateParam).Error; err != nil {
		return err
	}

	return nil
}
