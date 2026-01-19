package handlers

import (
	"time"

	orderpb "github.com/LavaJover/shvark-order-service/proto/gen/order"
)

// HTTP Models для Store Service

// StoreTrafficType - строковое представление enum
type StoreTrafficType string

const (
	StoreTrafficTypeUnspecified StoreTrafficType = "STORE_TRAFFIC_TYPE_UNSPECIFIED"
	StoreTrafficTypePayIn       StoreTrafficType = "STORE_TRAFFIC_TYPE_PAYIN"
	StoreTrafficTypePayOut      StoreTrafficType = "STORE_TRAFFIC_TYPE_PAYOUT"
)

func (s StoreTrafficType) ToProto() orderpb.StoreTrafficType {
	switch s {
	case StoreTrafficTypePayIn:
		return orderpb.StoreTrafficType_STORE_TRAFFIC_TYPE_PAYIN
	case StoreTrafficTypePayOut:
		return orderpb.StoreTrafficType_STORE_TRAFFIC_TYPE_PAYOUT
	default:
		return orderpb.StoreTrafficType_STORE_TRAFFIC_TYPE_UNSPECIFIED
	}
}

func FromProtoTrafficType(t orderpb.StoreTrafficType) StoreTrafficType {
	switch t {
	case orderpb.StoreTrafficType_STORE_TRAFFIC_TYPE_PAYIN:
		return StoreTrafficTypePayIn
	case orderpb.StoreTrafficType_STORE_TRAFFIC_TYPE_PAYOUT:
		return StoreTrafficTypePayOut
	default:
		return StoreTrafficTypeUnspecified
	}
}

// Store HTTP Models

type StoreHTTP struct {
	ID           string                 `json:"id"`
	MerchantID   string                 `json:"merchant_id"`
	Enabled      bool                   `json:"enabled"`
	MerchantInfo MerchantUserInfoHTTP   `json:"merchant_info"`
	BusinessParams StoreBusinessParamsHTTP `json:"business_params"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	DeletedAt    *time.Time             `json:"deleted_at,omitempty"`
}

type MerchantUserInfoHTTP struct {
	MerchantID       string `json:"merchant_id"`
	MerchantUsername string `json:"merchant_username"`
}

type StoreBusinessParamsHTTP struct {
	TrafficType           StoreTrafficType `json:"traffic_type"`
	DealPendingDuration   string           `json:"deal_pending_duration"` // e.g., "10s", "1h"
	DealCreatedDuration   string           `json:"deal_created_duration"`
	PlatformFee           float64          `json:"platform_fee"`
	Name                  string           `json:"name"`
}

// Pagination
type PaginationParamsHTTP struct {
	Page  int32 `json:"page" form:"page" binding:"min=1"`
	Limit int32 `json:"limit" form:"limit" binding:"min=1,max=100"`
}

type PaginationInfoHTTP struct {
	CurrentPage int32 `json:"current_page"`
	Limit       int32 `json:"limit"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int64 `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
	HasPrev     bool  `json:"has_prev"`
}

// Filters
type StoreFiltersHTTP struct {
	MerchantID  string           `json:"merchant_id" form:"merchant_id"`
	Enabled     *bool            `json:"enabled,omitempty" form:"enabled"`
	TrafficType StoreTrafficType `json:"traffic_type" form:"traffic_type"`
	Name        string           `json:"name" form:"name"`
	CreatedFrom *time.Time       `json:"created_from,omitempty" form:"created_from"`
	CreatedTo   *time.Time       `json:"created_to,omitempty" form:"created_to"`
	StoreIDs    []string         `json:"store_ids,omitempty" form:"store_ids"`
}

// Metrics
type StoreMetricsHTTP struct {
	StoreID               string     `json:"store_id"`
	TotalTraffics         int64      `json:"total_traffics"`
	ActiveTraffics        int64      `json:"active_traffics"`
	AverageTraderReward   float64    `json:"average_trader_reward"`
	UnlockedTrafficsCount int64      `json:"unlocked_traffics_count"`
	LockedTrafficsCount   int64      `json:"locked_traffics_count"`
	CreatedAt             time.Time  `json:"created_at"`
	LastActivityAt        time.Time  `json:"last_activity_at"`
	TotalRevenue          float64    `json:"total_revenue"`
	AverageDealAmount     float64    `json:"average_deal_amount"`
	TotalDeals            int64      `json:"total_deals"`
	SuccessfulDeals       int64      `json:"successful_deals"`
	FailedDeals           int64      `json:"failed_deals"`
}

// Request/Response Models

type CreateStoreRequestHTTP struct {
	MerchantID          string          `json:"merchant_id" binding:"required"`
	Enabled             bool            `json:"enabled"`
	MerchantUsername    string          `json:"merchant_username" binding:"required"`
	TrafficType         StoreTrafficType `json:"traffic_type" binding:"required"`
	DealPendingDuration string          `json:"deal_pending_duration"` // e.g., "10s", "1h"
	DealCreatedDuration string          `json:"deal_created_duration"`
	PlatformFee         float64         `json:"platform_fee" binding:"min=0,max=100"`
	Name                string          `json:"name" binding:"required,min=3,max=100"`
}

type CreateStoreResponseHTTP struct {
	Store StoreHTTP `json:"store"`
}

type GetStoreRequestHTTP struct {
	StoreID string `json:"store_id" uri:"store_id" binding:"required"`
}

type GetStoreResponseHTTP struct {
	Store StoreHTTP `json:"store"`
}

type UpdateStoreRequestHTTP struct {
	StoreID             string           `json:"store_id" uri:"store_id" binding:"required"`
	Enabled             *bool            `json:"enabled,omitempty"`
	MerchantUsername    string           `json:"merchant_username,omitempty"`
	TrafficType         *StoreTrafficType `json:"traffic_type,omitempty"`
	DealPendingDuration *string          `json:"deal_pending_duration,omitempty"`
	DealCreatedDuration *string          `json:"deal_created_duration,omitempty"`
	PlatformFee         *float64         `json:"platform_fee,omitempty" binding:"omitempty,min=0,max=1"`
	Name                string           `json:"name,omitempty"`
}

type UpdateStoreResponseHTTP struct {
	Store StoreHTTP `json:"store"`
}

type DeleteStoreRequestHTTP struct {
	StoreID string `json:"store_id" uri:"store_id" binding:"required"`
	Force   bool   `json:"force" form:"force"`
}

type DeleteStoreResponseHTTP struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ListStoresRequestHTTP struct {
	Pagination PaginationParamsHTTP `json:"pagination"`
	Filters    StoreFiltersHTTP     `json:"filters"`
}

type ListStoresResponseHTTP struct {
	Stores        []StoreHTTP        `json:"stores"`
	PaginationInfo PaginationInfoHTTP `json:"pagination_info"`
}

type GetStoreWithTrafficsRequestHTTP struct {
	StoreID                string `json:"store_id" uri:"store_id" binding:"required"`
	Page                   int32  `json:"page" form:"page" binding:"min=1"`
	Limit                  int32  `json:"limit" form:"limit" binding:"min=1,max=100"`
	TraderFilter           string `json:"trader_filter" form:"trader_filter"`
	IncludeTrafficDetails  bool   `json:"include_traffic_details" form:"include_traffic_details"`
}

type GetStoreWithTrafficsResponseHTTP struct {
	Store          StoreHTTP          `json:"store"`
	TrafficRecords []TrafficRecordHTTP `json:"traffic_records"`
	PaginationInfo PaginationInfoHTTP `json:"pagination_info"`
	TotalTraffics  int64              `json:"total_traffics"`
	TraderFilter   string             `json:"trader_filter"`
}

type GetStoreByTrafficIdRequestHTTP struct {
	TrafficID string `json:"traffic_id" uri:"traffic_id" binding:"required"`
}

type GetStoreByTrafficIdResponseHTTP struct {
	Store StoreHTTP `json:"store"`
}

type CheckStoreNameUniqueRequestHTTP struct {
	MerchantID string `json:"merchant_id" form:"merchant_id" binding:"required"`
	StoreName  string `json:"store_name" form:"store_name" binding:"required"`
}

type CheckStoreNameUniqueResponseHTTP struct {
	IsUnique   bool   `json:"is_unique"`
	Message    string `json:"message"`
	MerchantID string `json:"merchant_id"`
	StoreName  string `json:"store_name"`
}

type ValidateStoreForTrafficRequestHTTP struct {
	StoreID string `json:"store_id" uri:"store_id" binding:"required"`
}

type ValidateStoreForTrafficResponseHTTP struct {
	Valid   bool     `json:"valid"`
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

type ToggleStoreStatusRequestHTTP struct {
	StoreID string `json:"store_id" uri:"store_id" binding:"required"`
	Enabled bool   `json:"enabled" binding:"required"`
}

type ToggleStoreStatusResponseHTTP struct {
	Store         StoreHTTP `json:"store"`
	PreviousStatus string   `json:"previous_status"`
	NewStatus      string   `json:"new_status"`
}

type EnableStoreRequestHTTP struct {
	StoreID string `json:"store_id" uri:"store_id" binding:"required"`
}

type EnableStoreResponseHTTP struct {
	Store StoreHTTP `json:"store"`
}

type DisableStoreRequestHTTP struct {
	StoreID string `json:"store_id" uri:"store_id" binding:"required"`
	Force   bool   `json:"force" form:"force"`
}

type DisableStoreResponseHTTP struct {
	Store StoreHTTP `json:"store"`
}

type BulkUpdateStoresStatusRequestHTTP struct {
	StoreIDs []string `json:"store_ids" binding:"required,min=1"`
	Enabled  bool     `json:"enabled" `
}

type BulkUpdateStoresStatusResponseHTTP struct {
	UpdatedCount    int32    `json:"updated_count"`
	FailedCount     int32    `json:"failed_count"`
	FailedStoreIDs  []string `json:"failed_store_ids"`
	Errors          []string `json:"errors"`
}

type SearchStoresRequestHTTP struct {
	Pagination PaginationParamsHTTP `json:"pagination"`
	Filters    StoreFiltersHTTP     `json:"filters"`
	SearchQuery string              `json:"search_query"`
	SortBy      []string            `json:"sort_by"`
	SortDesc    bool                `json:"sort_desc"`
}

type SearchStoresResponseHTTP struct {
	Stores         []StoreHTTP        `json:"stores"`
	PaginationInfo PaginationInfoHTTP `json:"pagination_info"`
	AppliedFilters StoreFiltersHTTP   `json:"applied_filters"`
	SearchQuery    string             `json:"search_query"`
}

type GetStoresByMerchantRequestHTTP struct {
	MerchantID string `json:"merchant_id" uri:"merchant_id" binding:"required"`
	Page       int32  `json:"page" form:"page" binding:"min=1"`
	Limit      int32  `json:"limit" form:"limit" binding:"min=1,max=100"`
	OnlyActive bool   `json:"only_active" form:"only_active"`
}

type GetStoresByMerchantResponseHTTP struct {
	Stores       []StoreHTTP        `json:"stores"`
	PaginationInfo PaginationInfoHTTP `json:"pagination_info"`
	MerchantID   string             `json:"merchant_id"`
	TotalStores  int32              `json:"total_stores"`
	ActiveStores int32              `json:"active_stores"`
}

type GetActiveStoresRequestHTTP struct {
	MerchantID string `json:"merchant_id" form:"merchant_id" binding:"required"`
}

type GetActiveStoresResponseHTTP struct {
	Stores     []StoreHTTP `json:"stores"`
	MerchantID string      `json:"merchant_id"`
	Count      int32       `json:"count"`
}

type GetStoreMetricsRequestHTTP struct {
	StoreID    string     `json:"store_id" uri:"store_id" binding:"required"`
	PeriodFrom *time.Time `json:"period_from" form:"period_from"`
	PeriodTo   *time.Time `json:"period_to" form:"period_to"`
}

type GetStoreMetricsResponseHTTP struct {
	Metrics StoreMetricsHTTP `json:"metrics"`
}

type CalculateStoreMetricsRequestHTTP struct {
	StoreID string `json:"store_id" uri:"store_id" binding:"required"`
}

type CalculateStoreMetricsResponseHTTP struct {
	Metrics      StoreMetricsHTTP `json:"metrics"`
	CalculatedAt time.Time        `json:"calculated_at"`
}

type BatchGetStoresRequestHTTP struct {
	StoreIDs []string `json:"store_ids" binding:"required,min=1"`
}

type BatchGetStoresResponseHTTP struct {
	Stores      []StoreHTTP `json:"stores"`
	NotFoundIDs []string    `json:"not_found_ids"`
}

type HealthCheckResponseHTTP struct {
	Status    string            `json:"status"`
	Version   string            `json:"version"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

// Вспомогательные модели
type TrafficRecordHTTP struct {
	ID                  string                     `json:"id"`
	StoreID             string                     `json:"store_id"`
	TraderID            string                     `json:"trader_id"`
	TraderRewardPercent float64                    `json:"trader_reward_percent"`
	TraderPriority      float64                    `json:"trader_priority"`
	ActivityParams      TrafficActivityParamsHTTP  `json:"activity_params"`
	AntifraudParams     TrafficAntifraudParamsHTTP `json:"antifraud_params"`
	BusinessParams      TrafficBusinessParamsHTTP  `json:"business_params"`
	CreatedAt           time.Time                  `json:"created_at"`
	UpdatedAt           time.Time                  `json:"updated_at"`
}

type TrafficActivityParamsHTTP struct {
	MerchantUnlocked  bool `json:"merchant_unlocked"`
	TraderUnlocked    bool `json:"trader_unlocked"`
	AntifraudUnlocked bool `json:"antifraud_unlocked"`
	ManuallyUnlocked  bool `json:"manually_unlocked"`
}

type TrafficAntifraudParamsHTTP struct {
	AntifraudRequired bool `json:"antifraud_required"`
}

type TrafficBusinessParamsHTTP struct {
	MerchantDealsDuration string `json:"merchant_deals_duration"`
}

// Error Response
type ErrorResponseHTTP struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}