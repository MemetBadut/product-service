package main

import (
    "log"

    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/fiber/v3/middleware/logger"
    "github.com/gofiber/fiber/v3/middleware/recover"
    "github.com/MemetBadut/product-service/internal/config"
    "github.com/MemetBadut/product-service/internal/handler"
    "github.com/MemetBadut/product-service/internal/model"
    "github.com/MemetBadut/product-service/internal/repository"
    "github.com/MemetBadut/product-service/internal/service"
)

func main() {
    // 1. Load konfigurasi
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Gagal load config: %v", err)
    }

    // 2. Auto migrate database (buat tabel otomatis)
    if err := cfg.DB.AutoMigrate(&model.Product{}); err != nil {
        log.Fatalf("Gagal migrasi DB: %v", err)
    }
    log.Println("Database berhasil dimigrasi")

    // 3. Dependency injection (manual DI)
    productRepo := repository.NewProductRepository(cfg.DB)
    productSvc := service.NewProductService(productRepo)
    productHandler := handler.NewProductHandler(productSvc)

    // 4. Setup Fiber app
    app := fiber.New(fiber.Config{
        AppName: "Product Service v1.0",
    })

    // 5. Middleware
    app.Use(logger.New())  // Log setiap request
    app.Use(recover.New()) // Auto recover dari panic

    // 6. Health check route
    app.Get("/health", func(c fiber.Ctx) error {
        return c.JSON(fiber.Map{"status": "ok", "service": "product-service"})
    })

    // 7. Register API routes
    api := app.Group("/api/v1")
    productHandler.RegisterRoutes(api)

    // 8. Start server
    log.Printf("Product Service berjalan di port %s", cfg.AppPort)
    log.Fatal(app.Listen(":" + cfg.AppPort))
}