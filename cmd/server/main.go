package main

import (
	"context"
	"log"
	"os"

	"github.com/MemetBadut/product-service/internal/config"
	grpcserver "github.com/MemetBadut/product-service/internal/grpc"
	"github.com/MemetBadut/product-service/internal/handler"
	"github.com/MemetBadut/product-service/internal/kafka"
	"github.com/MemetBadut/product-service/internal/model"
	"github.com/MemetBadut/product-service/internal/repository"
	"github.com/MemetBadut/product-service/internal/service"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Gagal load config: %v", err)
    }

    cfg.DB.AutoMigrate(&model.Product{})

    productRepo := repository.NewProductRepository(cfg.DB)
    productSvc  := service.NewProductService(productRepo)
    productHandler := handler.NewProductHandler(productSvc)

    kafkaBroker := os.Getenv("KAFKA_BROKER")
    consumer := kafka.NewOrderConsumer(kafkaBroker, productSvc)
    defer consumer.Close()

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    go consumer.StartConsuming(ctx)

    // Jalankan gRPC server di goroutine terpisah
    // sehingga tidak memblock HTTP server
    go func() {
        if err := grpcserver.StartGRPCServer(cfg.GRPCPort, productSvc); err != nil {
            log.Fatalf("gRPC server error: %v", err)
        }
    }()

    app := fiber.New(fiber.Config{AppName: "Product Service v1.0"})
    app.Use(logger.New())
    app.Use(recover.New())

    app.Get("/health", func(c fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "status":    "ok",
            "service":   "product-service",
            "grpc_port": cfg.GRPCPort,
        })
    })

    api := app.Group("/api/v1")
    productHandler.RegisterRoutes(api)

    log.Printf("HTTP server di port %s | gRPC server di port %s", cfg.AppPort, cfg.GRPCPort)
    log.Fatal(app.Listen(":" + cfg.AppPort))
}