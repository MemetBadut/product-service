package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/MemetBadut/product-service/internal/service"
)

const TopicOrderEvents = "order-events"

type OrderEventPayload struct {
	EventName string `json:"event_name"`
	OrderID   uint   `json:"order_id"`
	ProductID uint   `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Status    string `json:"status"`
}

type OrderConsumer struct {
	reader     *kafka.Reader
	productSvc service.ProductService
}

func NewOrderConsumer(brokerAddr string, productSvc service.ProductService) *OrderConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{brokerAddr},
		Topic:    TopicOrderEvents,
		GroupID:  "product-service-group", // Consumer group untuk load balancing
		MinBytes: 10e3,                    // 10KB
		MaxBytes: 10e6,                    // 10MB
	})
	return &OrderConsumer{reader: reader, productSvc: productSvc}
}

// StartConsuming memulai loop konsumsi pesan dari Kafka.
// Harus dipanggil dalam goroutine terpisah.
func (c *OrderConsumer) StartConsuming(ctx context.Context) {
	log.Println("Kafka Consumer mulai listen topic: order-events")

	for {
		select {
		case <-ctx.Done():
			log.Println("Kafka Consumer berhenti")
			return
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error baca pesan kafka: %v", err)
				continue
			}

			// Parse payload
			var event OrderEventPayload
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Printf("Error parse pesan kafka: %v", err)
				continue
			}

			log.Printf("Terima event: %s OrderID=%d, ProductID=%d, Qty=%d, Status=%s",
				event.EventName, event.OrderID, event.ProductID, event.Quantity, event.Status)

			// Trigger pengurangan stok ditentukan oleh nama event.
			if event.EventName == "order.created" {
				if err := c.productSvc.CheckAndUpdateStock(event.ProductID, event.Quantity); err != nil {
					log.Printf("Gagal update stok untuk OrderID=%d: %v", event.OrderID, err)
					// Di production: kirim ke dead letter topic atau alert
				} else {
					log.Printf("Stok berhasil dikurangi %d untuk ProductID=%d",
						event.Quantity, event.ProductID)
				}
			}
		}
	}
}

func (c *OrderConsumer) Close() {
	c.reader.Close()
}
