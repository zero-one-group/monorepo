package product

import (
	"context"
	"math"

	"github.com/google/uuid"
	"github.com/zero-one-group/go-modulith/internal/errors"
	"go.opentelemetry.io/otel"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateProduct(ctx context.Context, req CreateProductRequest, createdBy uuid.UUID) (*ProductResponse, error) {
	ctx, span := otel.Tracer("product").Start(ctx, "service.create_product")
	defer span.End()

	if req.CategoryID != nil {
		category, err := s.repo.GetCategoryByID(ctx, *req.CategoryID)
		if err != nil {
			span.RecordError(err)
			return nil, errors.ErrInternal
		}
		if category == nil {
			return nil, errors.ErrBadRequest.WithDetails(map[string]string{
				"category_id": "Category not found",
			})
		}
	}

	product := &Product{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		CategoryID:  req.CategoryID,
		CreatedBy:   createdBy,
	}

	if err := s.repo.CreateProduct(ctx, product); err != nil {
		span.RecordError(err)
		return nil, errors.ErrInternal
	}

	response := product.ToResponse()

	if req.CategoryID != nil {
		category, _ := s.repo.GetCategoryByID(ctx, *req.CategoryID)
		if category != nil {
			categoryResponse := category.ToResponse()
			response.Category = &categoryResponse
		}
	}

	return &response, nil
}

func (s *Service) GetProducts(ctx context.Context, filters ProductFilters) (*PaginatedProductsResponse, error) {
	ctx, span := otel.Tracer("product").Start(ctx, "service.get_products")
	defer span.End()

	filters.SetDefaults()

	products, total, err := s.repo.GetProducts(ctx, filters)
	if err != nil {
		span.RecordError(err)
		return nil, errors.ErrInternal
	}

	productResponses := make([]ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = product.ToResponse()

		if product.CategoryID != nil {
			category, _ := s.repo.GetCategoryByID(ctx, *product.CategoryID)
			if category != nil {
				categoryResponse := category.ToResponse()
				productResponses[i].Category = &categoryResponse
			}
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(filters.Limit)))

	return &PaginatedProductsResponse{
		Products:   productResponses,
		Total:      total,
		Page:       filters.Page,
		Limit:      filters.Limit,
		TotalPages: totalPages,
	}, nil
}

func (s *Service) GetProductByID(ctx context.Context, id uuid.UUID) (*ProductResponse, error) {
	ctx, span := otel.Tracer("product").Start(ctx, "service.get_product_by_id")
	defer span.End()

	product, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, errors.ErrInternal
	}

	if product == nil {
		return nil, errors.ErrNotFound
	}

	response := product.ToResponse()

	if product.CategoryID != nil {
		category, _ := s.repo.GetCategoryByID(ctx, *product.CategoryID)
		if category != nil {
			categoryResponse := category.ToResponse()
			response.Category = &categoryResponse
		}
	}

	return &response, nil
}

func (s *Service) UpdateProduct(ctx context.Context, id uuid.UUID, req UpdateProductRequest, userID uuid.UUID) (*ProductResponse, error) {
	ctx, span := otel.Tracer("product").Start(ctx, "service.update_product")
	defer span.End()

	existingProduct, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, errors.ErrInternal
	}

	if existingProduct == nil {
		return nil, errors.ErrNotFound
	}

	if existingProduct.CreatedBy != userID {
		return nil, errors.ErrForbidden.WithDetails(map[string]string{
			"message": "You can only update your own products",
		})
	}

	if req.CategoryID != nil {
		category, err := s.repo.GetCategoryByID(ctx, *req.CategoryID)
		if err != nil {
			span.RecordError(err)
			return nil, errors.ErrInternal
		}
		if category == nil {
			return nil, errors.ErrBadRequest.WithDetails(map[string]string{
				"category_id": "Category not found",
			})
		}
	}

	updates := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"price":       req.Price,
		"category_id": req.CategoryID,
	}

	product, err := s.repo.UpdateProduct(ctx, id, updates)
	if err != nil {
		span.RecordError(err)
		return nil, errors.ErrInternal
	}

	response := product.ToResponse()

	if product.CategoryID != nil {
		category, _ := s.repo.GetCategoryByID(ctx, *product.CategoryID)
		if category != nil {
			categoryResponse := category.ToResponse()
			response.Category = &categoryResponse
		}
	}

	return &response, nil
}

func (s *Service) DeleteProduct(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	ctx, span := otel.Tracer("product").Start(ctx, "service.delete_product")
	defer span.End()

	existingProduct, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return errors.ErrInternal
	}

	if existingProduct == nil {
		return errors.ErrNotFound
	}

	if existingProduct.CreatedBy != userID {
		return errors.ErrForbidden.WithDetails(map[string]string{
			"message": "You can only delete your own products",
		})
	}

	if err := s.repo.DeleteProduct(ctx, id); err != nil {
		if err.Error() == "record not found" {
			return errors.ErrNotFound
		}
		span.RecordError(err)
		return errors.ErrInternal
	}

	return nil
}