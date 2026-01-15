package orderservice

import (
	"context"
	"time"

	orderpb "github.com/LavaJover/shvark-order-service/proto/gen/order"
)

// ============= АВТОМАТИКА =============

// ProcessAutomaticPayment обрабатывает автоматический платёж
func (c *OrderClient) ProcessAutomaticPayment(ctx context.Context, grpcReq *orderpb.ProcessAutomaticPaymentRequest) (*orderpb.ProcessAutomaticPaymentResponse, error) {
    timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    return c.service.ProcessAutomaticPayment(timeoutCtx, grpcReq)
}

// GetAutomaticLogs получает логи автоматики с фильтрацией
func (c *OrderClient) GetAutomaticLogs(ctx context.Context, req *orderpb.GetAutomaticLogsRequest) (*orderpb.GetAutomaticLogsResponse, error) {
    timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    return c.service.GetAutomaticLogs(timeoutCtx, req)
}

// GetAutomaticStats получает статистику автоматики
func (c *OrderClient) GetAutomaticStats(ctx context.Context, req *orderpb.GetAutomaticStatsRequest) (*orderpb.GetAutomaticStatsResponse, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    return c.service.GetAutomaticStats(ctx, req)
}