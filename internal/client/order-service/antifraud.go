package orderservice

import (
	"context"
	"time"

	orderpb "github.com/LavaJover/shvark-order-service/proto/gen/order"
)

// ============= АНТИФРОД =============

// CheckTrader проверяет трейдера по антифрод правилам
func (c *OrderClient) CheckTrader(r *orderpb.CheckTraderRequest) (*orderpb.CheckTraderResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    return c.antifraudService.CheckTrader(ctx, r)
}

// ProcessTraderCheck проверяет трейдера и обновляет статус трафика
func (c *OrderClient) ProcessTraderCheck(r *orderpb.ProcessTraderCheckRequest) (*orderpb.ProcessTraderCheckResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    return c.antifraudService.ProcessTraderCheck(ctx, r)
}

// CreateAntiFraudRule создает новое правило антифрода
func (c *OrderClient) CreateAntiFraudRule(r *orderpb.CreateRuleRequest) (*orderpb.CreateRuleResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    return c.antifraudService.CreateRule(ctx, r)
}

// UpdateAntiFraudRule обновляет правило антифрода
func (c *OrderClient) UpdateAntiFraudRule(r *orderpb.UpdateRuleRequest) (*orderpb.UpdateRuleResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    return c.antifraudService.UpdateRule(ctx, r)
}

// GetAntiFraudRules получает список правил антифрода
func (c *OrderClient) GetAntiFraudRules(r *orderpb.GetRulesRequest) (*orderpb.GetRulesResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    return c.antifraudService.GetRules(ctx, r)
}

// GetAntiFraudRule получает конкретное правило антифрода
func (c *OrderClient) GetAntiFraudRule(r *orderpb.GetRuleRequest) (*orderpb.GetRuleResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    return c.antifraudService.GetRule(ctx, r)
}

// DeleteAntiFraudRule удаляет правило антифрода
func (c *OrderClient) DeleteAntiFraudRule(r *orderpb.DeleteRuleRequest) (*orderpb.DeleteRuleResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    return c.antifraudService.DeleteRule(ctx, r)
}

// GetAuditLogs получает логи аудита антифрода
func (c *OrderClient) GetAuditLogs(r *orderpb.GetAuditLogsRequest) (*orderpb.GetAuditLogsResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    return c.antifraudService.GetAuditLogs(ctx, r)
}

// GetTraderAuditHistory получает историю проверок трейдера
func (c *OrderClient) GetTraderAuditHistory(r *orderpb.GetTraderAuditHistoryRequest) (*orderpb.GetTraderAuditHistoryResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    return c.antifraudService.GetTraderAuditHistory(ctx, r)
}

// ============= АНТИФРОД - Manual Unlock =============

// ManualUnlock вручную разблокирует трейдера с грейс-периодом
func (c *OrderClient) ManualUnlock(r *orderpb.ManualUnlockRequest) (*orderpb.ManualUnlockResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    return c.antifraudService.ManualUnlock(ctx, r)
}

// ResetGracePeriod сбрасывает грейс-период трейдера
func (c *OrderClient) ResetGracePeriod(r *orderpb.ResetGracePeriodRequest) (*orderpb.ResetGracePeriodResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    return c.antifraudService.ResetGracePeriod(ctx, r)
}

// GetUnlockHistory получает историю разблокировок трейдера
func (c *OrderClient) GetUnlockHistory(r *orderpb.GetUnlockHistoryRequest) (*orderpb.GetUnlockHistoryResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    return c.antifraudService.GetUnlockHistory(ctx, r)
}