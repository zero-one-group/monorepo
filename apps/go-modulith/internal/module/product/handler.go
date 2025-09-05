package product

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/zero-one-group/go-modulith/internal/errors"
	"github.com/zero-one-group/go-modulith/internal/middleware"
	"go.opentelemetry.io/otel"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(e *echo.Group, authMiddleware echo.MiddlewareFunc) {
	products := e.Group("/products")
	{
		products.GET("", h.GetProducts)
		products.GET("/:id", h.GetProductByID)
		products.POST("", h.CreateProduct, authMiddleware)
		products.PUT("/:id", h.UpdateProduct, authMiddleware)
		products.DELETE("/:id", h.DeleteProduct, authMiddleware)
	}
}

func (h *Handler) CreateProduct(c echo.Context) error {
	ctx, span := otel.Tracer("product").Start(c.Request().Context(), "handler.create_product")
	defer span.End()

	userIDStr := middleware.GetUserID(c)
	if userIDStr == "" {
		return errors.ErrUnauthorized
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return errors.ErrUnauthorized
	}

	var req CreateProductRequest
	if err := c.Bind(&req); err != nil {
		span.RecordError(err)
		return errors.ErrBadRequest
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	response, err := h.service.CreateProduct(ctx, req, userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *Handler) GetProducts(c echo.Context) error {
	ctx, span := otel.Tracer("product").Start(c.Request().Context(), "handler.get_products")
	defer span.End()

	var filters ProductFilters
	if err := c.Bind(&filters); err != nil {
		span.RecordError(err)
		return errors.ErrBadRequest
	}

	if err := c.Validate(&filters); err != nil {
		return err
	}

	response, err := h.service.GetProducts(ctx, filters)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) GetProductByID(c echo.Context) error {
	ctx, span := otel.Tracer("product").Start(c.Request().Context(), "handler.get_product_by_id")
	defer span.End()

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return errors.ErrBadRequest.WithDetails(map[string]string{
			"id": "Invalid UUID format",
		})
	}

	response, err := h.service.GetProductByID(ctx, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) UpdateProduct(c echo.Context) error {
	ctx, span := otel.Tracer("product").Start(c.Request().Context(), "handler.update_product")
	defer span.End()

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return errors.ErrBadRequest.WithDetails(map[string]string{
			"id": "Invalid UUID format",
		})
	}

	userIDStr := middleware.GetUserID(c)
	if userIDStr == "" {
		return errors.ErrUnauthorized
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return errors.ErrUnauthorized
	}

	var req UpdateProductRequest
	if err := c.Bind(&req); err != nil {
		span.RecordError(err)
		return errors.ErrBadRequest
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	response, err := h.service.UpdateProduct(ctx, id, req, userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) DeleteProduct(c echo.Context) error {
	ctx, span := otel.Tracer("product").Start(c.Request().Context(), "handler.delete_product")
	defer span.End()

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return errors.ErrBadRequest.WithDetails(map[string]string{
			"id": "Invalid UUID format",
		})
	}

	userIDStr := middleware.GetUserID(c)
	if userIDStr == "" {
		return errors.ErrUnauthorized
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return errors.ErrUnauthorized
	}

	if err := h.service.DeleteProduct(ctx, id, userID); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Product deleted successfully",
	})
}