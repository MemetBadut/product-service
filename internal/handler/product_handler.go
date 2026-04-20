package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/MemetBadut/product-service/internal/model"
	"github.com/MemetBadut/product-service/internal/service"
)

type ProductHandler struct {
	svc service.ProductService
}

func NewProductHandler(svc service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

// RegisterRoutes mendaftarkan semua route product ke router Fiber
func (h *ProductHandler) RegisterRoutes(router fiber.Router) {
	products := router.Group("/products")
	products.Get("/", h.GetAll)
	products.Get("/:id", h.GetByID)
	products.Post("/", h.Create)
	products.Put("/:id", h.Update)
	products.Delete("/:id", h.Delete)
}

func (h *ProductHandler) GetAll(c fiber.Ctx) error {
	products, err := h.svc.GetAllProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": products, "count": len(products)})
}

func (h *ProductHandler) GetByID(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID tidak valid"})
	}
	product, err := h.svc.GetProductByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": product})
}

func (h *ProductHandler) Create(c fiber.Ctx) error {
	var req model.CreateProductRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "request body tidak valid"})
	}
	product, err := h.svc.CreateProduct(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": product, "message": "produk berhasil dibuat"})
}

func (h *ProductHandler) Update(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID tidak valid"})
	}
	var req model.CreateProductRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "request body tidak valid"})
	}
	product, err := h.svc.UpdateProduct(uint(id), &req)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": product, "message": "produk berhasil diupdate"})
}

func (h *ProductHandler) Delete(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID tidak valid"})
	}
	if err := h.svc.DeleteProduct(uint(id)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "produk berhasil dihapus"})
}
