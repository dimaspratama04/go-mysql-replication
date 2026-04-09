package repository

import (
	"context"
	"time"

	"github.com/rnd/mysql-replication/config"
	"github.com/rnd/mysql-replication/internal/model"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type ProductRepository struct {
	primary  *gorm.DB
	replicas *gorm.DB
}

func NewProductRepository(db *config.Database) *ProductRepository {
	return &ProductRepository{primary: db.Primary, replicas: db.Replicas}
}

func (r *ProductRepository) FindAll(ctx context.Context) ([]model.Product, error) {
	var products []model.Product
	result := r.replicas.WithContext(ctx).Order("id DESC").Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

func (r *ProductRepository) FindByID(ctx context.Context, id uint) (*model.Product, error) {
	var product model.Product
	result := r.replicas.WithContext(ctx).First(&product, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &product, nil
}

func (r *ProductRepository) Create(ctx context.Context, req model.CreateProductRequest) (*model.Product, error) {
	product := model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	result := r.primary.WithContext(ctx).Create(&product)
	if result.Error != nil {
		log.Error().Err(result.Error).Str("operation", "CREATE").Msg("Failed to create product")
		return nil, result.Error
	}

	log.Info().
		Str("operation", "CREATE").
		Uint("product_id", product.ID).
		Str("name", product.Name).
		Str("description", product.Description).
		Float64("price", product.Price).
		Int("stock", product.Stock).
		Msg("[CREATE] Product created successfully")

	return &product, nil
}

func (r *ProductRepository) Update(ctx context.Context, id uint, req model.UpdateProductRequest) (*model.Product, error) {
	product, err := r.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, nil
	}

	// Snapshot before state for logging
	beforeName := product.Name
	beforePrice := product.Price
	beforeStock := product.Stock

	result := r.primary.WithContext(ctx).Model(product).Updates(map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"price":       req.Price,
		"stock":       req.Stock,
	})
	if result.Error != nil {
		log.Error().Err(result.Error).Str("operation", "UPDATE").Uint("product_id", id).Msg("Failed to update product")
		return nil, result.Error
	}

	log.Info().
		Str("operation", "UPDATE").
		Uint("product_id", id).
		Str("before.name", beforeName).
		Float64("before.price", beforePrice).
		Int("before.stock", beforeStock).
		Str("after.name", req.Name).
		Float64("after.price", req.Price).
		Int("after.stock", req.Stock).
		Msg("[UPDATE] Product updated successfully")

	return product, nil
}

func (r *ProductRepository) Delete(ctx context.Context, id uint) (bool, error) {
	product, err := r.FindByID(ctx, id)
	if err != nil {
		return false, err
	}
	if product == nil {
		return false, nil
	}

	result := r.primary.WithContext(ctx).Delete(&model.Product{}, id)
	if result.Error != nil {
		log.Error().Err(result.Error).Str("operation", "DELETE").Uint("product_id", id).Msg("Failed to delete product")
		return false, result.Error
	}

	log.Info().
		Str("operation", "DELETE").
		Str("type", "soft_delete").
		Uint("product_id", id).
		Str("name", product.Name).
		Float64("price", product.Price).
		Int("stock", product.Stock).
		Time("deleted_at", time.Now().UTC()).
		Msg("[DELETE] Product soft-deleted successfully")

	return true, nil
}
