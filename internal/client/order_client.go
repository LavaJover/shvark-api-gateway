package client

// import (
// 	"context"
// 	"time"

// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// 	"google.golang.org/protobuf/types/known/durationpb"
// 	"google.golang.org/protobuf/types/known/timestamppb"

// 	orderpb "github.com/LavaJover/shvark-order-service/proto/gen/order"
// )

// type OrderClient struct {
// 	conn *grpc.ClientConn
// 	service orderpb.OrderServiceClient
// 	trafficService orderpb.TrafficServiceClient
// 	bankDetailService orderpb.BankDetailServiceClient
// 	teamRelationsService orderpb.TeamRelationsServiceClient
// 	deviceService orderpb.DeviceServiceClient
// 	antifraudService orderpb.AntiFraudServiceClient
// }

// func NewOrderClient(addr string) (*OrderClient, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	conn, err := grpc.DialContext(
// 		ctx,
// 		addr,
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 		grpc.WithBlock(),
// 		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &OrderClient{
// 		conn: conn,
// 		service: orderpb.NewOrderServiceClient(conn),
// 		trafficService: orderpb.NewTrafficServiceClient(conn),
// 		bankDetailService: orderpb.NewBankDetailServiceClient(conn),
// 		teamRelationsService: orderpb.NewTeamRelationsServiceClient(conn),
// 		deviceService: orderpb.NewDeviceServiceClient(conn),
// 		antifraudService: orderpb.NewAntiFraudServiceClient(conn),
// 	}, nil
// }

// func (c *OrderClient) CreatePayInOrder(orderRequest *orderpb.CreatePayInOrderRequest) (*orderpb.CreatePayInOrderResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	return c.service.CreatePayInOrder(
// 		ctx,
// 		orderRequest,
// 	)
// }

// func (c *OrderClient) CreatePayOutOrder(orderRequest *orderpb.CreatePayOutOrderRequest) (*orderpb.CreatePayOutOrderResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	return c.service.CreatePayOutOrder(
// 		ctx,
// 		orderRequest,
// 	)
// }

// func (c *OrderClient) GetOrderByID(orderID string) (*orderpb.GetOrderByIDResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	return c.service.GetOrderByID(
// 		ctx,
// 		&orderpb.GetOrderByIDRequest{
// 			OrderId: orderID,
// 		},
// 	)
// }

// func (c *OrderClient) GetOrderByMerchantOrderID(merchantOrderID string) (*orderpb.GetOrderByMerchantOrderIDResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	return c.service.GetOrderByMerchantOrderID(
// 		ctx,
// 		&orderpb.GetOrderByMerchantOrderIDRequest{
// 			MerchantOrderId: merchantOrderID,
// 		},
// 	)
// }

// func (c *OrderClient) GetOrdersByTraderID(request *orderpb.GetOrdersByTraderIDRequest) (*orderpb.GetOrdersByTraderIDResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	return c.service.GetOrdersByTraderID(
// 		ctx,
// 		request,
// 	)
// }

// func (c *OrderClient) ApproveOrder(orderID string) (*orderpb.ApproveOrderResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	return c.service.ApproveOrder(
// 		ctx,
// 		&orderpb.ApproveOrderRequest{
// 			OrderId: orderID,
// 		},
// 	)
// }

// func (c *OrderClient) CancelOrder(orderID string) (*orderpb.CancelOrderResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	return c.service.CancelOrder(
// 		ctx,
// 		&orderpb.CancelOrderRequest{
// 			OrderId: orderID,
// 		},
// 	)
// }

// func (c *OrderClient) SetTraderLockTrafficStatus(r *orderpb.SetTraderLockTrafficStatusRequest) (*orderpb.SetTraderLockTrafficStatusResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	return c.trafficService.SetTraderLockTrafficStatus(
// 		ctx,
// 		r,
// 	)
// }

// func (c *OrderClient) SetMerchantLockTrafficStatus(r *orderpb.SetMerchantLockTrafficStatusRequest) (*orderpb.SetMerchantLockTrafficStatusResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	return c.trafficService.SetMerchantLockTrafficStatus(
// 		ctx,
// 		r,
// 	)
// }

// func (c *OrderClient) SetManuallyLockTrafficStatus(r *orderpb.SetManuallyLockTrafficStatusRequest) (*orderpb.SetManuallyLockTrafficStatusResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	return c.trafficService.SetManuallyLockTrafficStatus(
// 		ctx,
// 		r,
// 	)
// }

// func (c *OrderClient) SetAntifraudLockTrafficStatus(r *orderpb.SetAntifraudLockTrafficStatusRequest) (*orderpb.SetAntifraudLockTrafficStatusResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	return c.trafficService.SetAntifraudLockTrafficStatus(
// 		ctx,
// 		r,
// 	)
// }

// func (c *OrderClient) AddTraffic(r *orderpb.AddTrafficRequest) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	_, err := c.trafficService.AddTraffic(
// 		ctx,
// 		r,
// 	)
// 	return err
// }

// func (c *OrderClient) EditTraffic(r *orderpb.EditTrafficRequest) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	_, err := c.trafficService.EditTraffic(
// 		ctx,
// 		r,
// 	)

// 	return err
// }

// func (c *OrderClient) GetTrafficRecords(page, limit int32) ([]*orderpb.Traffic, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	trafficResponse, err := c.trafficService.GetTrafficRecords(
// 		ctx,
// 		&orderpb.GetTrafficRecordsRequest{
// 			Page: page,
// 			Limit: limit,
// 		},
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return trafficResponse.TrafficRecords, nil
// }

// func (c *OrderClient) DeleteTraffic(trafficID string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	_, err := c.trafficService.DeleteTraffic(
// 		ctx,
// 		&orderpb.DeleteTrafficRequest{
// 			TrafficId: trafficID,
// 		},
// 	)

// 	return err
// }

// func (c *OrderClient) CreateDispute(
// 	orderID, proofUrl, disputeReason string,
// 	ttl time.Duration,
// 	disputeAmountFiat float64,
// ) (string, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	disputeResponse, err := c.service.CreateOrderDispute(
// 		ctx,
// 		&orderpb.CreateOrderDisputeRequest{
// 			OrderId: orderID,
// 			ProofUrl: proofUrl,
// 			DisputeReason: disputeReason,
// 			Ttl: durationpb.New(ttl),
// 			DisputeAmountFiat: disputeAmountFiat,
// 		},
// 	)

// 	if err != nil {
// 		return "", err
// 	}

// 	return disputeResponse.DisputeId, nil
// }

// func (c *OrderClient) AcceptDispute(
// 	disputeID string,
// ) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	_, err := c.service.AcceptOrderDispute(
// 		ctx,
// 		&orderpb.AcceptOrderDisputeRequest{
// 			DisputeId: disputeID,
// 		},
// 	)

// 	return err
// }

// func (c *OrderClient) RejectDispute(
// 	disputeID string,
// ) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	_, err := c.service.RejectOrderDispute(
// 		ctx,
// 		&orderpb.RejectOrderDisputeRequest{
// 			DisputeId: disputeID,
// 		},
// 	)
// 	return err
// }

// type Dispute struct {
// 	DisputeID 	  string
// 	OrderID 	  string
// 	ProofUrl 	  string
// 	DisputeReason string
// 	DisputeStatus string
// }

// func (c *OrderClient) GetDisputeInfo(
// 	disputeID string,
// ) (*Dispute, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	disputeResponse, err := c.service.GetOrderDisputeInfo(
// 		ctx,
// 		&orderpb.GetOrderDisputeInfoRequest{
// 			DisputeId: disputeID,
// 		},
// 	)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &Dispute{
// 		DisputeID: disputeResponse.Dispute.DisputeId,
// 		OrderID: disputeResponse.Dispute.OrderId,
// 		ProofUrl: disputeResponse.Dispute.ProofUrl,
// 		DisputeReason: disputeResponse.Dispute.DisputeReason,
// 		DisputeStatus: disputeResponse.Dispute.DisputeStatus,
// 	}, nil
// }

// func (c *OrderClient) FreeezeDispute(disputeID string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	_, err := c.service.FreezeOrderDispute(
// 		ctx,
// 		&orderpb.FreezeOrderDisputeRequest{
// 			DisputeId: disputeID,
// 		},
// 	)

// 	return err
// }

// func (c *OrderClient) CreateBankDetail(createBankDetailRequest *orderpb.CreateBankDetailRequest) (*orderpb.CreateBankDetailResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	return c.bankDetailService.CreateBankDetail(
// 		ctx,
// 		createBankDetailRequest,
// 	)
// }

// func (c *OrderClient) EditBankDetail(editBankDetailRequest *orderpb.UpdateBankDetailRequest) (*orderpb.UpdateBankDetailResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	return c.bankDetailService.UpdateBankDetail(
// 		ctx,
// 		editBankDetailRequest,
// 	)
// }

// func (c *OrderClient) DeleteBankDetail(deleteBankDetailRequest *orderpb.DeleteBankDetailRequest) (*orderpb.DeleteBankDetailResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	return c.bankDetailService.DeleteBankDetail(
// 		ctx,
// 		deleteBankDetailRequest,
// 	)
// }

// func (c *OrderClient) GetBankDetailsByTraderID(getBankDetailsRequest *orderpb.GetBankDetailsByTraderIDRequest) (*orderpb.GetBankDetailsByTraderIDResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	return c.bankDetailService.GetBankDetailsByTraderID(
// 		ctx,
// 		getBankDetailsRequest,
// 	)
// }

// func (c *OrderClient) GetBankDetailByID(getbankDetailRequest *orderpb.GetBankDetailByIDRequest) (*orderpb.GetBankDetailByIDResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	return c.bankDetailService.GetBankDetailByID(
// 		ctx,
// 		getbankDetailRequest,
// 	)
// }

// func (c *OrderClient) GetBankDetailsStatsByTraderID(getStatsRequest *orderpb.GetBankDetailsStatsByTraderIDRequest) (*orderpb.GetBankDetailsStatsByTraderIDResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	return c.bankDetailService.GetBankDetailsStatsByTraderID(
// 		ctx,
// 		getStatsRequest,
// 	)
// }


// func (c *OrderClient) GetOrderDisputes(r *orderpb.GetOrderDisputesRequest) (*orderpb.GetOrderDisputesResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	return c.service.GetOrderDisputes(
// 		ctx,
// 		r,
// 	)
// }

// type OrderStats struct {
// 	TotalOrders 			int64 	
// 	SucceedOrders 			int64 	
// 	CanceledOrders 			int64 	
// 	ProcessedAmountFiat 	float64 
// 	ProcessedAmountCrypto 	float64 
// 	CanceledAmountFiat 		float64 
// 	CanceledAmountCrypto 	float64 
// 	IncomeCrypto 			float64 
// }

// func (c *OrderClient) GetOrderStats(
// 	traderID string,
// 	dateFrom, dateTo time.Time,
// ) (*orderpb.GetOrderStatisticsResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	return c.service.GetOrderStatistics(
// 		ctx,
// 		&orderpb.GetOrderStatisticsRequest{
// 			TraderId: traderID,
// 			DateFrom: timestamppb.New(dateFrom),
// 			DateTo: timestamppb.New(dateTo),
// 		},
// 	)
// }

// func (c *OrderClient) GetOrders(r *orderpb.GetOrdersRequest) (*orderpb.GetOrdersResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	return c.service.GetOrders(
// 		ctx,
// 		r,
// 	)
// }

// func (c *OrderClient) CreateTeamRelation(r *orderpb.CreateTeamRelationRequest) (*orderpb.CreateTeamRelationResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	return c.teamRelationsService.CreateTeamRelation(
// 		ctx,
// 		r,
// 	)
// }

// func (c *OrderClient) UpdateTeamRelationParams(r *orderpb.UpdateRelationParamsRequest) (*orderpb.UpdateRelationParamsResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	return c.teamRelationsService.UpdateRelationParams(
// 		ctx,
// 		r,
// 	)
// }

// func (c *OrderClient) GetTeamRelationsByTeamLeadID(r *orderpb.GetRelationsByTeamLeadIDRequest) (*orderpb.GetRelationsByTeamLeadIDResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	return c.teamRelationsService.GetRelationsByTeamLeadID(
// 		ctx,
// 		r,
// 	)
// }

// func (c *OrderClient) DeleteTeamRelationship(r *orderpb.DeleteTeamRelationshipRequest) (*orderpb.DeleteTeamRelationshipResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	return c.teamRelationsService.DeleteTeamRelationship(
// 		ctx,
// 		r,
// 	)
// }

// func (c *OrderClient) GetAllOrders(r *orderpb.GetAllOrdersRequest) (*orderpb.GetAllOrdersResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	return c.service.GetAllOrders(
// 		ctx,
// 		r,
// 	)
// }

// func (c *OrderClient) CreateDevice(r *orderpb.CreateDeviceRequest) (*orderpb.CreateDeviceResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	return c.deviceService.CreateDevice(
// 		ctx,
// 		r,
// 	)
// }

// func (c *OrderClient) DeleteDevice(r *orderpb.DeleteDeviceRequest) (*orderpb.DeleteDeviceResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	return c.deviceService.DeleteDevice(
// 		ctx,
// 		r,
// 	)
// }

// func (c *OrderClient) EditeDevice(r *orderpb.EditDeviceRequest) (*orderpb.EditDeviceResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	return c.deviceService.EditDevice(
// 		ctx,
// 		r,
// 	)
// }

// func (c *OrderClient) GetTraderDevices(r *orderpb.GetTraderDevicesRequest) (*orderpb.GetTraderDevicesResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	return c.deviceService.GetTraderDevices(
// 		ctx,
// 		r,
// 	)
// }

// func (c *OrderClient) GetBankDetails(r *orderpb.GetBankDetailsRequest) (*orderpb.GetBankDetailsResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	return c.bankDetailService.GetBankDetails(
// 		ctx,
// 		r,
// 	)
// }

// func (c *OrderClient) ProcessAutomaticPayment(ctx context.Context, grpcReq *orderpb.ProcessAutomaticPaymentRequest) (*orderpb.ProcessAutomaticPaymentResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	return c.service.ProcessAutomaticPayment(
// 		ctx,
// 		grpcReq,
// 	)
// }

// GetTrafficLockStatuses получает статусы блокировки трафика
// func (c *OrderClient) GetTrafficLockStatuses(r *orderpb.GetTrafficLockStatusesRequest) (*orderpb.GetTrafficLockStatusesResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	return c.trafficService.GetTrafficLockStatuses(ctx, r)
// }

// // CheckTrafficUnlocked проверяет, разблокирован ли трафик
// func (c *OrderClient) CheckTrafficUnlocked(r *orderpb.CheckTrafficUnlockedRequest) (*orderpb.CheckTrafficUnlockedResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	return c.trafficService.CheckTrafficUnlocked(ctx, r)
// }

// ============= АНТИФРОД =============

// // CheckTrader проверяет трейдера по антифрод правилам
// func (c *OrderClient) CheckTrader(r *orderpb.CheckTraderRequest) (*orderpb.CheckTraderResponse, error) {
//     ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//     defer cancel()

//     return c.antifraudService.CheckTrader(ctx, r)
// }

// // ProcessTraderCheck проверяет трейдера и обновляет статус трафика
// func (c *OrderClient) ProcessTraderCheck(r *orderpb.ProcessTraderCheckRequest) (*orderpb.ProcessTraderCheckResponse, error) {
//     ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//     defer cancel()

//     return c.antifraudService.ProcessTraderCheck(ctx, r)
// }

// // CreateAntiFraudRule создает новое правило антифрода
// func (c *OrderClient) CreateAntiFraudRule(r *orderpb.CreateRuleRequest) (*orderpb.CreateRuleResponse, error) {
//     ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//     defer cancel()

//     return c.antifraudService.CreateRule(ctx, r)
// }

// // UpdateAntiFraudRule обновляет правило антифрода
// func (c *OrderClient) UpdateAntiFraudRule(r *orderpb.UpdateRuleRequest) (*orderpb.UpdateRuleResponse, error) {
//     ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//     defer cancel()

//     return c.antifraudService.UpdateRule(ctx, r)
// }

// // GetAntiFraudRules получает список правил антифрода
// func (c *OrderClient) GetAntiFraudRules(r *orderpb.GetRulesRequest) (*orderpb.GetRulesResponse, error) {
//     ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//     defer cancel()

//     return c.antifraudService.GetRules(ctx, r)
// }

// // GetAntiFraudRule получает конкретное правило антифрода
// func (c *OrderClient) GetAntiFraudRule(r *orderpb.GetRuleRequest) (*orderpb.GetRuleResponse, error) {
//     ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//     defer cancel()

//     return c.antifraudService.GetRule(ctx, r)
// }

// // DeleteAntiFraudRule удаляет правило антифрода
// func (c *OrderClient) DeleteAntiFraudRule(r *orderpb.DeleteRuleRequest) (*orderpb.DeleteRuleResponse, error) {
//     ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//     defer cancel()

//     return c.antifraudService.DeleteRule(ctx, r)
// }

// // GetAuditLogs получает логи аудита антифрода
// func (c *OrderClient) GetAuditLogs(r *orderpb.GetAuditLogsRequest) (*orderpb.GetAuditLogsResponse, error) {
//     ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//     defer cancel()

//     return c.antifraudService.GetAuditLogs(ctx, r)
// }

// // GetTraderAuditHistory получает историю проверок трейдера
// func (c *OrderClient) GetTraderAuditHistory(r *orderpb.GetTraderAuditHistoryRequest) (*orderpb.GetTraderAuditHistoryResponse, error) {
//     ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//     defer cancel()

//     return c.antifraudService.GetTraderAuditHistory(ctx, r)
// }

// // ============= АНТИФРОД - Manual Unlock =============

// // ManualUnlock вручную разблокирует трейдера с грейс-периодом
// func (c *OrderClient) ManualUnlock(r *orderpb.ManualUnlockRequest) (*orderpb.ManualUnlockResponse, error) {
//     ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//     defer cancel()

//     return c.antifraudService.ManualUnlock(ctx, r)
// }

// // ResetGracePeriod сбрасывает грейс-период трейдера
// func (c *OrderClient) ResetGracePeriod(r *orderpb.ResetGracePeriodRequest) (*orderpb.ResetGracePeriodResponse, error) {
//     ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//     defer cancel()

//     return c.antifraudService.ResetGracePeriod(ctx, r)
// }

// // GetUnlockHistory получает историю разблокировок трейдера
// func (c *OrderClient) GetUnlockHistory(r *orderpb.GetUnlockHistoryRequest) (*orderpb.GetUnlockHistoryResponse, error) {
//     ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//     defer cancel()

//     return c.antifraudService.GetUnlockHistory(ctx, r)
// }

// // GetTraderTraffic получает все записи трафика для трейдера
// func (c *OrderClient) GetTraderTraffic(r *orderpb.GetTraderTrafficRequest) (*orderpb.GetTraderTrafficResponse, error) {
//     ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//     defer cancel()

//     return c.trafficService.GetTraderTraffic(ctx, r)
// }

// // ============= АВТОМАТИКА =============

// // ProcessAutomaticPayment обрабатывает автоматический платёж
// func (c *OrderClient) ProcessAutomaticPayment(ctx context.Context, grpcReq *orderpb.ProcessAutomaticPaymentRequest) (*orderpb.ProcessAutomaticPaymentResponse, error) {
//     timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
//     defer cancel()

//     return c.service.ProcessAutomaticPayment(timeoutCtx, grpcReq)
// }

// // GetAutomaticLogs получает логи автоматики с фильтрацией
// func (c *OrderClient) GetAutomaticLogs(ctx context.Context, req *orderpb.GetAutomaticLogsRequest) (*orderpb.GetAutomaticLogsResponse, error) {
//     timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
//     defer cancel()

//     return c.service.GetAutomaticLogs(timeoutCtx, req)
// }

// ============= СТАТУС УСТРОЙСТВ =============

// UpdateDeviceLiveness обновляет статус онлайн устройства
// func (c *OrderClient) UpdateDeviceLiveness(ctx context.Context, req *orderpb.UpdateDeviceLivenessRequest) (*orderpb.UpdateDeviceLivenessResponse, error) {
//     timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
//     defer cancel()

//     return c.deviceService.UpdateDeviceLiveness(timeoutCtx, req)
// }

// // GetDeviceStatus получает статус устройства
// func (c *OrderClient) GetDeviceStatus(ctx context.Context, req *orderpb.GetDeviceStatusRequest) (*orderpb.GetDeviceStatusResponse, error) {
//     timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
//     defer cancel()

//     return c.deviceService.GetDeviceStatus(timeoutCtx, req)
// }

// // GetTraderDevicesStatus получает статусы всех устройств трейдера
// func (c *OrderClient) GetTraderDevicesStatus(ctx context.Context, req *orderpb.GetTraderDevicesStatusRequest) (*orderpb.GetTraderDevicesStatusResponse, error) {
//     timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
//     defer cancel()

//     return c.deviceService.GetTraderDevicesStatus(timeoutCtx, req)
// }

// // GetAutomaticStats получает статистику автоматики
// func (c *OrderClient) GetAutomaticStats(ctx context.Context, req *orderpb.GetAutomaticStatsRequest) (*orderpb.GetAutomaticStatsResponse, error) {
//     ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
//     defer cancel()

//     return c.service.GetAutomaticStats(ctx, req)
// }