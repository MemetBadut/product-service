package grpcserver

import (
    "context"
    "fmt"
    "log"
    "net"

    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"

    pb "github.com/MemetBadut/product-service/proto/product"
    "github.com/MemetBadut/product-service/internal/service"
)

// ProductGRPCServer mengimplementasikan interface yang digenerate protoc
type ProductGRPCServer struct {
    pb.UnimplementedProductServiceServer // Wajib: untuk forward compatibility
    productSvc service.ProductService
}

func NewProductGRPCServer(svc service.ProductService) *ProductGRPCServer {
    return &ProductGRPCServer{productSvc: svc}
}

// CheckStock dipanggil oleh Order Service sebelum membuat order
func (s *ProductGRPCServer) CheckStock(ctx context.Context, req *pb.CheckStockRequest) (*pb.CheckStockResponse, error) {
    product, err := s.productSvc.GetProductByID(uint(req.ProductId))
    if err != nil {
        return nil, status.Error(codes.NotFound, "product not found")
    }

    available := product.Stock >= int(req.Quantity)
    msg := "Stok tersedia"
    if !available {
        msg = fmt.Sprintf("Stok tidak cukup. Tersedia: %d, diminta: %d", product.Stock, req.Quantity)
    }

    return &pb.CheckStockResponse{
        Available:    available,
        Message:      msg,
        CurrentStock: int32(product.Stock),
    }, nil
}

// DecreaseStock tidak dipanggil langsung dari Order Service.
// Pengurangan stok akan dipicu lewat event order.created di Hari 4.
func (s *ProductGRPCServer) UpdateStock(ctx context.Context, req *pb.UpdateStockRequest) (*pb.UpdateStockResponse, error) {
    err := s.productSvc.CheckAndUpdateStock(uint(req.ProductId), int(req.Quantity))
    if err != nil {
        return &pb.UpdateStockResponse{Success: false, Message: err.Error()}, nil
    }
    return &pb.UpdateStockResponse{Success: true, Message: "Stok berhasil diupdate"}, nil
}

// GetProduct mengembalikan detail produk via gRPC
func (s *ProductGRPCServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
    product, err := s.productSvc.GetProductByID(uint(req.ProductId))
    if err != nil {
        return nil, status.Error(codes.NotFound, "product not found")
    }

    return &pb.GetProductResponse{
        Id:    uint64(product.ID),
        Name:  product.Name,
        Price: product.Price,
        Stock: int32(product.Stock),
    }, nil
}

// StartGRPCServer menjalankan gRPC server di port terpisah
func StartGRPCServer(port string, svc service.ProductService) error {
    lis, err := net.Listen("tcp", ":"+port)
    if err != nil {
        return fmt.Errorf("gagal listen gRPC port: %w", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterProductServiceServer(grpcServer, NewProductGRPCServer(svc))

    log.Printf("gRPC Product Server berjalan di port %s", port)
    return grpcServer.Serve(lis)
}