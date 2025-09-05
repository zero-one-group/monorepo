package product

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/zero-one-group/go-modulith/internal/database"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateProduct(ctx context.Context, product *Product) error {
	ctx, span := otel.Tracer("product").Start(ctx, "repository.create_product")
	defer span.End()

	if err := r.db.WithContext(ctx).Create(product).Error; err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (r *Repository) GetProducts(ctx context.Context, filters ProductFilters) ([]Product, int64, error) {
	ctx, span := otel.Tracer("product").Start(ctx, "repository.get_products")
	defer span.End()

	var products []Product
	var total int64

	query := r.db.WithContext(ctx).Model(&Product{})

	if filters.Search != "" {
		searchPattern := "%" + strings.ToLower(filters.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchPattern, searchPattern)
	}

	if filters.CategoryID != nil {
		query = query.Where("category_id = ?", *filters.CategoryID)
	}

	if filters.MinPrice != nil {
		query = query.Where("price >= ?", *filters.MinPrice)
	}

	if filters.MaxPrice != nil {
		query = query.Where("price <= ?", *filters.MaxPrice)
	}

	if err := query.Count(&total).Error; err != nil {
		span.RecordError(err)
		return nil, 0, err
	}

	orderClause := fmt.Sprintf("%s %s", filters.Sort, strings.ToUpper(filters.Order))
	if err := query.Order(orderClause).
		Offset(filters.Offset()).
		Limit(filters.Limit).
		Find(&products).Error; err != nil {
		span.RecordError(err)
		return nil, 0, err
	}

	return products, total, nil
}

func (r *Repository) GetProductByID(ctx context.Context, id uuid.UUID) (*Product, error) {
	ctx, span := otel.Tracer("product").Start(ctx, "repository.get_product_by_id")
	defer span.End()

	var product Product
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		span.RecordError(err)
		return nil, err
	}

	return &product, nil
}

func (r *Repository) UpdateProduct(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*Product, error) {
	ctx, span := otel.Tracer("product").Start(ctx, "repository.update_product")
	defer span.End()

	var product Product
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		span.RecordError(err)
		return nil, err
	}

	if err := r.db.WithContext(ctx).Model(&product).Updates(updates).Error; err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &product, nil
}

func (r *Repository) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	ctx, span := otel.Tracer("product").Start(ctx, "repository.delete_product")
	defer span.End()

	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&Product{})
	if result.Error != nil {
		span.RecordError(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *Repository) GetCategoryByID(ctx context.Context, id uuid.UUID) (*Category, error) {
	ctx, span := otel.Tracer("product").Start(ctx, "repository.get_category_by_id")
	defer span.End()

	var category Category
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		span.RecordError(err)
		return nil, err
	}

	return &category, nil
}