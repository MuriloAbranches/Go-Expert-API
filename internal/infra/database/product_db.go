package database

import (
	"github.com/MuriloAbranches/Go-Expert-API/internal/entity"
	"gorm.io/gorm"
)

type ProductDB struct {
	DB *gorm.DB
}

func NewProductDB(db *gorm.DB) *ProductDB {
	return &ProductDB{DB: db}
}

func (p *ProductDB) Create(product *entity.Product) error {
	return p.DB.Create(product).Error
}

func (p *ProductDB) FindAll(page, limit int, sort string) ([]entity.Product, error) {
	var products []entity.Product
	var err error

	if sort != "" && sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	if page != 0 && limit != 0 {
		err = p.DB.Limit(limit).Offset((page - 1) * limit).Order("created_at " + sort).Find(&products).Error
	} else {
		err = p.DB.Order("created_at " + sort).Find(&products).Error
	}

	return products, err
}

func (p *ProductDB) FindByID(id string) (*entity.Product, error) {
	var product entity.Product

	err := p.DB.First(&product, "id = ?", id).Error

	return &product, err
}

func (p *ProductDB) Update(product *entity.Product) error {
	_, err := p.FindByID(product.ID.String())
	if err != nil {
		return err
	}

	return p.DB.Save(product).Error
}

func (p *ProductDB) Delete(id string) error {
	product, err := p.FindByID(id)
	if err != nil {
		return err
	}

	return p.DB.Delete(product).Error
}
