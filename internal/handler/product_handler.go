package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rnd/mysql-replication/internal/model"
	"github.com/rnd/mysql-replication/internal/repository"
	"github.com/rs/zerolog/log"
)

type ProductHandler struct {
	repo *repository.ProductRepository
}

func NewProductHandler(repo *repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

// GET /api/products
func (h *ProductHandler) GetAll(c *fiber.Ctx) error {
	products, err := h.repo.FindAll(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("GetAll: failed to fetch products")
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: "Failed to fetch products",
		})
	}
	return c.JSON(model.Response{
		Success: true,
		Message: "Products retrieved successfully",
		Data:    products,
	})
}

// GET /api/products/:id
func (h *ProductHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Success: false,
			Message: "Invalid product ID",
		})
	}

	product, err := h.repo.FindByID(c.Context(), uint(id))
	if err != nil {
		log.Error().Err(err).Uint64("product_id", id).Msg("GetByID: failed to fetch product")
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: "Failed to fetch product",
		})
	}
	if product == nil {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Success: false,
			Message: "Product not found",
		})
	}

	return c.JSON(model.Response{
		Success: true,
		Message: "Product retrieved successfully",
		Data:    product,
	})
}

// POST /api/products
func (h *ProductHandler) Create(c *fiber.Ctx) error {
	var req model.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Success: false,
			Message: "Invalid request body: " + err.Error(),
		})
	}

	if err := validateCreate(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(model.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	product, err := h.repo.Create(c.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("Create: failed to create product")
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: "Failed to create product",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.Response{
		Success: true,
		Message: "Product created successfully",
		Data:    product,
	})
}

// PUT /api/products/:id
func (h *ProductHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Success: false,
			Message: "Invalid product ID",
		})
	}

	var req model.UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Success: false,
			Message: "Invalid request body: " + err.Error(),
		})
	}

	if err := validateUpdate(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(model.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	product, err := h.repo.Update(c.Context(), uint(id), req)
	if err != nil {
		log.Error().Err(err).Uint64("product_id", id).Msg("Update: failed to update product")
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: "Failed to update product",
		})
	}
	if product == nil {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Success: false,
			Message: "Product not found",
		})
	}

	return c.JSON(model.Response{
		Success: true,
		Message: "Product updated successfully",
		Data:    product,
	})
}

// DELETE /api/products/:id
func (h *ProductHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Success: false,
			Message: "Invalid product ID",
		})
	}

	deleted, err := h.repo.Delete(c.Context(), uint(id))
	if err != nil {
		log.Error().Err(err).Uint64("product_id", id).Msg("Delete: failed to delete product")
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: "Failed to delete product",
		})
	}
	if !deleted {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Success: false,
			Message: "Product not found",
		})
	}

	return c.JSON(model.Response{
		Success: true,
		Message: "Product deleted successfully",
	})
}

func validateCreate(req model.CreateProductRequest) error {
	if req.Name == "" {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "name is required")
	}
	if req.Price <= 0 {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "price must be greater than 0")
	}
	if req.Stock < 0 {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "stock must be >= 0")
	}
	return nil
}

func validateUpdate(req model.UpdateProductRequest) error {
	if req.Name == "" {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "name is required")
	}
	if req.Price <= 0 {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "price must be greater than 0")
	}
	if req.Stock < 0 {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "stock must be >= 0")
	}
	return nil
}
