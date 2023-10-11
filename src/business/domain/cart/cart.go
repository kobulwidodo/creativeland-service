package cart

import (
	"fmt"
	"go-clean/src/business/entity"

	"gorm.io/gorm"
)

type Interface interface {
	Create(cart entity.Cart) (entity.Cart, error)
	GetList(param entity.CartParam) ([]entity.Cart, error)
	GetListInByID(ids []int64) ([]entity.Cart, error)
	GetListInByTransactionID(transaction_ids []uint) ([]entity.Cart, error)
	GetListInByStatus(status []string, param entity.CartParam) ([]entity.Cart, error)
	Get(param entity.CartParam) (entity.Cart, error)
	Update(selectParam entity.CartParam, updateParam entity.UpdateCartParam) error
	UpdatesByIDs(ids []uint, updateParam entity.UpdateCartParam) error
	Delete(param entity.CartParam) error
}

type cart struct {
	db *gorm.DB
}

func Init(db *gorm.DB) Interface {
	c := &cart{
		db: db,
	}

	return c
}

func (c *cart) Create(cart entity.Cart) (entity.Cart, error) {
	if err := c.db.Create(&cart).Error; err != nil {
		return cart, err
	}

	return cart, nil
}

func (c *cart) GetList(param entity.CartParam) ([]entity.Cart, error) {
	cart := []entity.Cart{}
	query := c.db.Where(param).Where("created_at LIKE ?", fmt.Sprintf("%%%s%%", param.CreatedAt))

	if param.CreatedAt != "" {
		query = query.Where("created_at LIKE ?", fmt.Sprintf("%%%s%%", param.CreatedAt))
	} else if !param.CreatedAtMoreThan.IsZero() {
		query = query.Where("created_at > ?", param.CreatedAtMoreThan)
	}

	if err := query.Find(&cart).Error; err != nil {
		return cart, err
	}

	return cart, nil
}

func (c *cart) GetListInByID(ids []int64) ([]entity.Cart, error) {
	carts := []entity.Cart{}

	if err := c.db.Where(ids).Find(&carts).Error; err != nil {
		return carts, err
	}

	return carts, nil
}

func (c *cart) GetListInByTransactionID(transaction_ids []uint) ([]entity.Cart, error) {
	carts := []entity.Cart{}

	if err := c.db.Where("transaction_id IN ?", transaction_ids).Find(&carts).Error; err != nil {
		return carts, err
	}

	return carts, nil
}

func (c *cart) GetListInByStatus(status []string, param entity.CartParam) ([]entity.Cart, error) {
	carts := []entity.Cart{}

	if err := c.db.Where("status IN ?", status).Where(param).Find(&carts).Error; err != nil {
		return carts, err
	}

	return carts, nil
}

func (c *cart) Get(param entity.CartParam) (entity.Cart, error) {
	cart := entity.Cart{}

	if err := c.db.Where(param).First(&cart).Error; err != nil {
		return cart, err
	}

	return cart, nil
}

func (c *cart) Update(selectParam entity.CartParam, updateParam entity.UpdateCartParam) error {
	if err := c.db.Model(entity.Cart{}).Where(selectParam).Updates(updateParam).Error; err != nil {
		return err
	}

	return nil
}

func (c *cart) UpdatesByIDs(ids []uint, updateParam entity.UpdateCartParam) error {
	if err := c.db.Model(entity.Cart{}).Where("id IN ?", ids).Updates(updateParam).Error; err != nil {
		return err
	}

	return nil
}

func (c *cart) Delete(param entity.CartParam) error {
	if err := c.db.Where(param).Delete(&entity.Cart{}).Error; err != nil {
		return err
	}

	return nil
}
