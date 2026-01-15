package handlers

import (
    "net/http"
    "strconv"
    "time"

    "github.com/LavaJover/shvark-api-gateway/internal/client/order-service"
    antifraudpb "github.com/LavaJover/shvark-order-service/proto/gen/order"
    "github.com/gin-gonic/gin"
    "google.golang.org/protobuf/types/known/structpb"
    "google.golang.org/protobuf/types/known/timestamppb"
)

type AntiFraudHandler struct {
    orderClient *orderservice.OrderClient
}

func NewAntiFraudHandler(orderClient *orderservice.OrderClient) *AntiFraudHandler {
    return &AntiFraudHandler{
        orderClient: orderClient,
    }
}

// ============= ПРОВЕРКА ТРЕЙДЕРА =============

// @Summary Check trader
// @Description Check trader against all active antifraud rules
// @Tags antifraud
// @Accept json
// @Produce json
// @Param traderID path string true "Trader ID"
// @Success 200 {object} CheckTraderResponse
// @Failure 400 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /antifraud/traders/{traderID}/check [post]
func (h *AntiFraudHandler) CheckTrader(c *gin.Context) {
    traderID := c.Param("traderID")
    if traderID == "" {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "traderID is required"})
        return
    }

    response, err := h.orderClient.CheckTrader(&antifraudpb.CheckTraderRequest{
        TraderId: traderID,
    })
    if err != nil {
        c.JSON(http.StatusBadGateway, ErrorResponse{Error: err.Error()})
        return
    }

    c.JSON(http.StatusOK, h.convertCheckTraderResponse(response))
}

// @Summary Process trader check
// @Description Check trader and update traffic status
// @Tags antifraud
// @Accept json
// @Produce json
// @Param traderID path string true "Trader ID"
// @Success 200 {object} ProcessTraderCheckResponse
// @Failure 400 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /antifraud/traders/{traderID}/process [post]
func (h *AntiFraudHandler) ProcessTraderCheck(c *gin.Context) {
    traderID := c.Param("traderID")
    if traderID == "" {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "traderID is required"})
        return
    }

    response, err := h.orderClient.ProcessTraderCheck(&antifraudpb.ProcessTraderCheckRequest{
        TraderId: traderID,
    })
    if err != nil {
        c.JSON(http.StatusBadGateway, ErrorResponse{Error: err.Error()})
        return
    }

    c.JSON(http.StatusOK, ProcessTraderCheckResponse{
        Success: response.Success,
        Message: response.Message,
    })
}

// ============= УПРАВЛЕНИЕ ПРАВИЛАМИ =============

// @Summary Create antifraud rule
// @Description Create a new antifraud rule
// @Tags antifraud
// @Accept json
// @Produce json
// @Param request body CreateRuleRequest true "Rule data"
// @Success 200 {object} CreateRuleResponse
// @Failure 400 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /antifraud/rules [post]
func (h *AntiFraudHandler) CreateRule(c *gin.Context) {
    var req CreateRuleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
        return
    }

    config, err := structpb.NewStruct(req.Config)
    if err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid config format"})
        return
    }

    response, err := h.orderClient.CreateAntiFraudRule(&antifraudpb.CreateRuleRequest{
        Name:     req.Name,
        Type:     req.Type,
        Config:   config,
        Priority: int32(req.Priority),
    })
    if err != nil {
        c.JSON(http.StatusBadGateway, ErrorResponse{Error: err.Error()})
        return
    }

    c.JSON(http.StatusOK, h.convertRuleResponse(response.Rule))
}

// @Summary Update antifraud rule
// @Description Update an existing antifraud rule
// @Tags antifraud
// @Accept json
// @Produce json
// @Param ruleID path string true "Rule ID"
// @Param request body UpdateRuleRequest true "Update data"
// @Success 200 {object} UpdateRuleResponse
// @Failure 400 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /antifraud/rules/{ruleID} [patch]
func (h *AntiFraudHandler) UpdateRule(c *gin.Context) {
    ruleID := c.Param("ruleID")
    if ruleID == "" {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "ruleID is required"})
        return
    }

    var req UpdateRuleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
        return
    }

    protoReq := &antifraudpb.UpdateRuleRequest{
        RuleId: ruleID,
    }

    if req.Config != nil {
        config, err := structpb.NewStruct(req.Config)
        if err != nil {
            c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid config format"})
            return
        }
        protoReq.Config = config
    }

    if req.IsActive != nil {
        protoReq.IsActive = req.IsActive
    }

    if req.Priority != nil {
        priority := int32(*req.Priority)
        protoReq.Priority = &priority
    }

    response, err := h.orderClient.UpdateAntiFraudRule(protoReq)
    if err != nil {
        c.JSON(http.StatusBadGateway, ErrorResponse{Error: err.Error()})
        return
    }

    c.JSON(http.StatusOK, UpdateRuleResponse{
        Success: response.Success,
        Message: response.Message,
    })
}

// @Summary Get antifraud rules
// @Description Get list of antifraud rules
// @Tags antifraud
// @Accept json
// @Produce json
// @Param active_only query bool false "Get only active rules"
// @Success 200 {object} GetRulesResponse
// @Failure 502 {object} ErrorResponse
// @Router /antifraud/rules [get]
func (h *AntiFraudHandler) GetRules(c *gin.Context) {
    activeOnlyStr := c.Query("active_only")
    activeOnly := false
    if activeOnlyStr != "" {
        activeOnly, _ = strconv.ParseBool(activeOnlyStr)
    }

    response, err := h.orderClient.GetAntiFraudRules(&antifraudpb.GetRulesRequest{
        ActiveOnly: activeOnly,
    })
    if err != nil {
        c.JSON(http.StatusBadGateway, ErrorResponse{Error: err.Error()})
        return
    }

    rules := make([]AntiFraudRuleResponse, 0, len(response.Rules))
    for _, rule := range response.Rules {
        rules = append(rules, h.convertRuleResponse(rule))
    }

    c.JSON(http.StatusOK, GetRulesResponse{
        Rules: rules,
    })
}

// @Summary Get antifraud rule
// @Description Get a specific antifraud rule by ID
// @Tags antifraud
// @Accept json
// @Produce json
// @Param ruleID path string true "Rule ID"
// @Success 200 {object} AntiFraudRuleResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /antifraud/rules/{ruleID} [get]
func (h *AntiFraudHandler) GetRule(c *gin.Context) {
    ruleID := c.Param("ruleID")
    if ruleID == "" {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "ruleID is required"})
        return
    }

    response, err := h.orderClient.GetAntiFraudRule(&antifraudpb.GetRuleRequest{
        RuleId: ruleID,
    })
    if err != nil {
        c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
        return
    }

    c.JSON(http.StatusOK, h.convertRuleResponse(response.Rule))
}

// @Summary Delete antifraud rule
// @Description Delete an antifraud rule
// @Tags antifraud
// @Accept json
// @Produce json
// @Param ruleID path string true "Rule ID"
// @Success 200 {object} DeleteRuleResponse
// @Failure 400 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /antifraud/rules/{ruleID} [delete]
func (h *AntiFraudHandler) DeleteRule(c *gin.Context) {
    ruleID := c.Param("ruleID")
    if ruleID == "" {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "ruleID is required"})
        return
    }

    response, err := h.orderClient.DeleteAntiFraudRule(&antifraudpb.DeleteRuleRequest{
        RuleId: ruleID,
    })
    if err != nil {
        c.JSON(http.StatusBadGateway, ErrorResponse{Error: err.Error()})
        return
    }

    c.JSON(http.StatusOK, DeleteRuleResponse{
        Success: response.Success,
        Message: response.Message,
    })
}

// ============= АУДИТ =============

// @Summary Get audit logs
// @Description Get antifraud audit logs with optional filters
// @Tags antifraud
// @Accept json
// @Produce json
// @Param trader_id query string false "Filter by trader ID"
// @Param from_date query string false "Filter from date (RFC3339)"
// @Param to_date query string false "Filter to date (RFC3339)"
// @Param only_failed query bool false "Show only failed checks"
// @Param limit query int false "Limit results" default(50)
// @Param offset query int false "Offset results" default(0)
// @Success 200 {object} GetAuditLogsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /antifraud/audit-logs [get]
func (h *AntiFraudHandler) GetAuditLogs(c *gin.Context) {
    req := &antifraudpb.GetAuditLogsRequest{}

    if traderID := c.Query("trader_id"); traderID != "" {
        req.TraderId = &traderID
    }

    if fromDateStr := c.Query("from_date"); fromDateStr != "" {
        fromDate, err := time.Parse(time.RFC3339, fromDateStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid from_date format"})
            return
        }
        req.FromDate = timestamppb.New(fromDate)
    }

    if toDateStr := c.Query("to_date"); toDateStr != "" {
        toDate, err := time.Parse(time.RFC3339, toDateStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid to_date format"})
            return
        }
        req.ToDate = timestamppb.New(toDate)
    }

    if onlyFailedStr := c.Query("only_failed"); onlyFailedStr != "" {
        onlyFailed, _ := strconv.ParseBool(onlyFailedStr)
        req.OnlyFailed = onlyFailed
    }

    limit := 50
    if limitStr := c.Query("limit"); limitStr != "" {
        if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
            limit = l
        }
    }
    req.Limit = int32(limit)

    offset := 0
    if offsetStr := c.Query("offset"); offsetStr != "" {
        if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
            offset = o
        }
    }
    req.Offset = int32(offset)

    response, err := h.orderClient.GetAuditLogs(req)
    if err != nil {
        c.JSON(http.StatusBadGateway, ErrorResponse{Error: err.Error()})
        return
    }

    logs := make([]AuditLogResponse, 0, len(response.Logs))
    for _, log := range response.Logs {
        logs = append(logs, h.convertAuditLogResponse(log))
    }

    c.JSON(http.StatusOK, GetAuditLogsResponse{
        Logs:  logs,
        Total: response.Total,
    })
}

// @Summary Get trader audit history
// @Description Get audit history for a specific trader
// @Tags antifraud
// @Accept json
// @Produce json
// @Param traderID path string true "Trader ID"
// @Param limit query int false "Limit results" default(10)
// @Success 200 {object} GetTraderAuditHistoryResponse
// @Failure 400 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /antifraud/traders/{traderID}/audit-history [get]
func (h *AntiFraudHandler) GetTraderAuditHistory(c *gin.Context) {
    traderID := c.Param("traderID")
    if traderID == "" {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "traderID is required"})
        return
    }

    limit := 10
    if limitStr := c.Query("limit"); limitStr != "" {
        if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
            limit = l
        }
    }

    response, err := h.orderClient.GetTraderAuditHistory(&antifraudpb.GetTraderAuditHistoryRequest{
        TraderId: traderID,
        Limit:    int32(limit),
    })
    if err != nil {
        c.JSON(http.StatusBadGateway, ErrorResponse{Error: err.Error()})
        return
    }

    logs := make([]AuditLogResponse, 0, len(response.Logs))
    for _, log := range response.Logs {
        logs = append(logs, h.convertAuditLogResponse(log))
    }

    c.JSON(http.StatusOK, GetTraderAuditHistoryResponse{
        Logs: logs,
    })
}

// ============= ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ =============

func (h *AntiFraudHandler) convertCheckTraderResponse(response *antifraudpb.CheckTraderResponse) CheckTraderResponse {
    results := make([]CheckResultResponse, 0, len(response.Results))
    for _, r := range response.Results {
        results = append(results, CheckResultResponse{
            RuleName: r.RuleName,
            Passed:   r.Passed,
            Message:  r.Message,
            Details:  r.Details.AsMap(),
        })
    }

    return CheckTraderResponse{
        TraderID:    response.TraderId,
        CheckedAt:   response.CheckedAt.AsTime(),
        AllPassed:   response.AllPassed,
        Results:     results,
        FailedRules: response.FailedRules,
    }
}

func (h *AntiFraudHandler) convertRuleResponse(rule *antifraudpb.AntiFraudRule) AntiFraudRuleResponse {
    return AntiFraudRuleResponse{
        ID:        rule.Id,
        Name:      rule.Name,
        Type:      rule.Type,
        Config:    rule.Config.AsMap(),
        IsActive:  rule.IsActive,
        Priority:  int(rule.Priority),
        CreatedAt: rule.CreatedAt.AsTime(),
        UpdatedAt: rule.UpdatedAt.AsTime(),
    }
}

func (h *AntiFraudHandler) convertAuditLogResponse(log *antifraudpb.AuditLog) AuditLogResponse {
    results := make([]CheckResultResponse, 0, len(log.Results))
    for _, r := range log.Results {
        results = append(results, CheckResultResponse{
            RuleName: r.RuleName,
            Passed:   r.Passed,
            Message:  r.Message,
            Details:  r.Details.AsMap(),
        })
    }

    return AuditLogResponse{
        ID:        log.Id,
        TraderID:  log.TraderId,
        CheckedAt: log.CheckedAt.AsTime(),
        AllPassed: log.AllPassed,
        Results:   results,
        CreatedAt: log.CreatedAt.AsTime(),
    }
}

// ============= СТРУКТУРЫ ОТВЕТОВ =============

type CheckTraderResponse struct {
    TraderID    string                `json:"trader_id"`
    CheckedAt   time.Time             `json:"checked_at"`
    AllPassed   bool                  `json:"all_passed"`
    Results     []CheckResultResponse `json:"results"`
    FailedRules []string              `json:"failed_rules,omitempty"`
}

type CheckResultResponse struct {
    RuleName string                 `json:"rule_name"`
    Passed   bool                   `json:"passed"`
    Message  string                 `json:"message"`
    Details  map[string]interface{} `json:"details,omitempty"`
}

type ProcessTraderCheckResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}

type CreateRuleRequest struct {
    Name     string                 `json:"name" binding:"required"`
    Type     string                 `json:"type" binding:"required"`
    Config   map[string]interface{} `json:"config" binding:"required"`
    Priority int                    `json:"priority"`
}

type CreateRuleResponse struct {
    Rule AntiFraudRuleResponse `json:"rule"`
}

type UpdateRuleRequest struct {
    Config   map[string]interface{} `json:"config,omitempty"`
    IsActive *bool                  `json:"is_active,omitempty"`
    Priority *int                   `json:"priority,omitempty"`
}

type UpdateRuleResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}

type GetRulesResponse struct {
    Rules []AntiFraudRuleResponse `json:"rules"`
}

type DeleteRuleResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}

type AntiFraudRuleResponse struct {
    ID        string                 `json:"id"`
    Name      string                 `json:"name"`
    Type      string                 `json:"type"`
    Config    map[string]interface{} `json:"config"`
    IsActive  bool                   `json:"is_active"`
    Priority  int                    `json:"priority"`
    CreatedAt time.Time              `json:"created_at"`
    UpdatedAt time.Time              `json:"updated_at"`
}

type GetAuditLogsResponse struct {
    Logs  []AuditLogResponse `json:"logs"`
    Total int32              `json:"total"`
}

type GetTraderAuditHistoryResponse struct {
    Logs []AuditLogResponse `json:"logs"`
}

type AuditLogResponse struct {
    ID        string                `json:"id"`
    TraderID  string                `json:"trader_id"`
    CheckedAt time.Time             `json:"checked_at"`
    AllPassed bool                  `json:"all_passed"`
    Results   []CheckResultResponse `json:"results"`
    CreatedAt time.Time             `json:"created_at"`
}

// ============= Manual Unlock =============

type ManualUnlockRequest struct {
    AdminID          string `json:"admin_id" binding:"required"`
    Reason           string `json:"reason" binding:"required"`
    GracePeriodHours int    `json:"grace_period_hours"`
}

type ManualUnlockResponse struct {
    Success          bool      `json:"success"`
    Message          string    `json:"message"`
    GracePeriodUntil time.Time `json:"grace_period_until"`
}

type ResetGracePeriodResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}

// @Summary Manual unlock trader
// @Description Manually unlock trader with grace period
// @Tags antifraud
// @Accept json
// @Produce json
// @Param traderID path string true "Trader ID"
// @Param request body ManualUnlockRequest true "Unlock data"
// @Success 200 {object} ManualUnlockResponse
// @Failure 400 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /antifraud/traders/{traderID}/manual-unlock [post]
func (h *AntiFraudHandler) ManualUnlock(c *gin.Context) {
    traderID := c.Param("traderID")
    if traderID == "" {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "traderID is required"})
        return
    }

    var req ManualUnlockRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
        return
    }

    if req.GracePeriodHours <= 0 {
        req.GracePeriodHours = 24
    }

    response, err := h.orderClient.ManualUnlock(&antifraudpb.ManualUnlockRequest{
        TraderId:         traderID,
        AdminId:          req.AdminID,
        Reason:           req.Reason,
        GracePeriodHours: int32(req.GracePeriodHours),
    })

    if err != nil {
        c.JSON(http.StatusBadGateway, ErrorResponse{Error: err.Error()})
        return
    }

    c.JSON(http.StatusOK, ManualUnlockResponse{
        Success:          response.Success,
        Message:          response.Message,
        GracePeriodUntil: response.GracePeriodUntil.AsTime(),
    })
}

// @Summary Reset grace period
// @Description Reset grace period for trader
// @Tags antifraud
// @Accept json
// @Produce json
// @Param traderID path string true "Trader ID"
// @Success 200 {object} ResetGracePeriodResponse
// @Failure 400 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /antifraud/traders/{traderID}/reset-grace-period [post]
func (h *AntiFraudHandler) ResetGracePeriod(c *gin.Context) {
    traderID := c.Param("traderID")
    if traderID == "" {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "traderID is required"})
        return
    }

    response, err := h.orderClient.ResetGracePeriod(&antifraudpb.ResetGracePeriodRequest{
        TraderId: traderID,
    })

    if err != nil {
        c.JSON(http.StatusBadGateway, ErrorResponse{Error: err.Error()})
        return
    }

    c.JSON(http.StatusOK, ResetGracePeriodResponse{
        Success: response.Success,
        Message: response.Message,
    })
}

type UnlockHistoryItem struct {
    ID               string    `json:"id"`
    TraderID         string    `json:"trader_id"`
    AdminID          string    `json:"admin_id"`
    Reason           string    `json:"reason"`
    GracePeriodHours int       `json:"grace_period_hours"`
    UnlockedAt       time.Time `json:"unlocked_at"`
    CreatedAt        time.Time `json:"created_at"`
}

type GetUnlockHistoryResponse struct {
    Items []UnlockHistoryItem `json:"items"`
}

// @Summary Get unlock history
// @Description Get history of manual unlocks for trader
// @Tags antifraud
// @Accept json
// @Produce json
// @Param traderID path string true "Trader ID"
// @Param limit query int false "Limit results" default(20)
// @Success 200 {object} GetUnlockHistoryResponse
// @Failure 400 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /antifraud/traders/{traderID}/unlock-history [get]
func (h *AntiFraudHandler) GetUnlockHistory(c *gin.Context) {
    traderID := c.Param("traderID")
    if traderID == "" {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "traderID is required"})
        return
    }

    limit := 20
    if limitStr := c.Query("limit"); limitStr != "" {
        if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
            limit = l
        }
    }

    response, err := h.orderClient.GetUnlockHistory(&antifraudpb.GetUnlockHistoryRequest{
        TraderId: traderID,
        Limit:    int32(limit),
    })

    if err != nil {
        c.JSON(http.StatusBadGateway, ErrorResponse{Error: err.Error()})
        return
    }

    items := make([]UnlockHistoryItem, 0, len(response.Items))
    for _, item := range response.Items {
        items = append(items, UnlockHistoryItem{
            ID:               item.Id,
            TraderID:         item.TraderId,
            AdminID:          item.AdminId,
            Reason:           item.Reason,
            GracePeriodHours: int(item.GracePeriodHours),
            UnlockedAt:       item.UnlockedAt.AsTime(),
            CreatedAt:        item.CreatedAt.AsTime(),
        })
    }

    c.JSON(http.StatusOK, GetUnlockHistoryResponse{
        Items: items,
    })
}