package service

import (
    "bytes"
    "fmt"
    "html/template"
    "log"
    "time"

    "github.com/LavaJover/shvark-api-gateway/internal/client"
    "github.com/LavaJover/shvark-api-gateway/internal/domain"
    "github.com/LavaJover/shvark-api-gateway/internal/service/deeplink_templates"
    orderpb "github.com/LavaJover/shvark-order-service/proto/gen/order"
)

type DeeplinkService struct {
    orderClient *client.OrderClient
}

func NewDeeplinkService(orderClient *client.OrderClient) *DeeplinkService {
    return &DeeplinkService{
        orderClient: orderClient,
    }
}

// GenerateBankSelectionPage генерирует страницу выбора банков
func (ds *DeeplinkService) GenerateBankSelectionPage(orderID string) (*domain.DeeplinkData, error) {
    order, err := ds.orderClient.GetOrderByID(orderID)
    if err != nil {
        return nil, fmt.Errorf("failed to get order: %w", err)
    }

    paymentSystem := "all"
    if order.Order.BankDetail != nil {
        paymentSystem = order.Order.BankDetail.PaymentSystem
    }

    // Получаем доступные шаблоны для платежной системы
    availableTemplates := deeplink_templates.GetTemplatesForSystem(paymentSystem)
    
    // Если нет специфичных шаблонов, используем все
    if len(availableTemplates) == 0 {
        availableTemplates = deeplink_templates.GetAllTemplates()
    }

    templateData := ds.prepareTemplateData(order, nil, "bank_selection")
    templateData["AvailableBanks"] = availableTemplates
    templateData["PaymentSystem"] = paymentSystem

    htmlContent, err := ds.renderBankSelectionTemplate(templateData)
    if err != nil {
        return nil, fmt.Errorf("failed to render bank selection template: %w", err)
    }

    return &domain.DeeplinkData{
        HTMLContent:  htmlContent,
        DeeplinkType: "bank_selection",
        BankCode:     "multiple",
        OrderID:      orderID,
    }, nil
}

// GenerateSpecificDeeplink генерирует диплинк для конкретного банка
func (ds *DeeplinkService) GenerateSpecificDeeplink(orderID, bankCode string, phoneNumber *string) (*domain.DeeplinkData, error) {
    order, err := ds.orderClient.GetOrderByID(orderID)
    if err != nil {
        return nil, fmt.Errorf("failed to get order: %w", err)
    }

    // Получаем конфигурацию шаблона
    templateConfig, exists := deeplink_templates.BankTemplates[bankCode]
    if !exists {
        return nil, fmt.Errorf("bank template not found: %s", bankCode)
    }

    templateData := ds.prepareTemplateData(order, phoneNumber, bankCode)
    htmlContent, err := ds.renderSpecificTemplate(templateConfig.Template, templateData)
    if err != nil {
        return nil, fmt.Errorf("failed to render template for bank %s: %w", bankCode, err)
    }

    return &domain.DeeplinkData{
        HTMLContent:  htmlContent,
        DeeplinkType: bankCode,
        BankCode:     bankCode,
        OrderID:      orderID,
    }, nil
}

func (ds *DeeplinkService) prepareTemplateData(order *orderpb.GetOrderByIDResponse, phoneNumber *string, deeplinkType string) map[string]interface{} {
    data := map[string]interface{}{
        "Amount":        fmt.Sprintf("%.2f", order.Order.AmountFiat),
        "OrderID":       order.Order.OrderId,
        "PhoneNumber":   "",
        "Timestamp":     time.Now().Format("2006-01-02 15:04:05"),
        "PaymentSystem": "",
        "CardNumber":    "",
        "MaskedCardNumber": "", // Теперь тоже будет содержать полный номер
    }

    if order.Order.BankDetail != nil {
        data["PaymentSystem"] = order.Order.BankDetail.PaymentSystem
        
        // Обрабатываем номер карты - БЕЗ МАСКИРОВКИ
        if order.Order.BankDetail.CardNumber != "" {
            cardNumber := order.Order.BankDetail.CardNumber
            data["CardNumber"] = cardNumber
            data["MaskedCardNumber"] = cardNumber // Убираем маскировку - показываем полный номер
        }
        
        // Обрабатываем телефон
        if order.Order.BankDetail.Phone != "" {
            data["PhoneNumber"] = order.Order.BankDetail.Phone
        }
    }

    if phoneNumber != nil {
        data["PhoneNumber"] = *phoneNumber
    }

    log.Printf("Template data prepared - Card: %s, Phone: %s, Amount: %s", 
        data["CardNumber"], data["PhoneNumber"], data["Amount"])

    return data
}

func (ds *DeeplinkService) renderBankSelectionTemplate(data map[string]interface{}) (string, error) {
    tmpl := template.Must(template.New("bank_selection").Parse(`
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Выберите способ оплаты</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
        .container { background: white; padding: 30px; border-radius: 15px; box-shadow: 0 10px 30px rgba(0,0,0,0.2); }
        .header { text-align: center; margin-bottom: 30px; }
        .payment-info { background: #f8f9fa; padding: 20px; border-radius: 10px; margin: 20px 0; border-left: 5px solid #667eea; }
        .bank-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 15px; margin: 30px 0; }
        .bank-card { background: white; border: 2px solid #e9ecef; border-radius: 10px; padding: 20px; text-align: center; cursor: pointer; transition: all 0.3s ease; }
        .bank-card:hover { transform: translateY(-5px); box-shadow: 0 5px 15px rgba(0,0,0,0.1); border-color: #667eea; }
        .bank-icon { font-size: 2em; margin-bottom: 10px; }
        .bank-name { font-weight: bold; margin: 10px 0; color: #333; }
        .amount { font-size: 1.5em; font-weight: bold; color: #28a745; margin: 10px 0; }
        .info-text { color: #6c757d; font-size: 0.9em; }
        .recommended { border-color: #28a745; background: #f8fff9; }
        .recommended-badge { background: #28a745; color: white; padding: 2px 8px; border-radius: 10px; font-size: 0.8em; margin-left: 5px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎯 Выберите способ оплаты</h1>
            <p class="info-text">Выберите банк для быстрой оплаты через мобильное приложение</p>
        </div>

        <div class="payment-info">
            <h3>Данные платежа:</h3>
            {{if .MaskedCardNumber}}
            <p><strong>Номер карты:</strong> {{.CardNumber}}</p>
            {{end}}
            {{if .PhoneNumber}}
            <p><strong>Телефон:</strong> {{.PhoneNumber}}</p>
            {{end}}
            <p class="amount">{{.Amount}} ₽</p>
            <p class="info-text">Order ID: {{.OrderID}}</p>
            {{if eq .PaymentSystem "C2C"}}
            <p class="info-text" style="color: #28a745; margin-top: 10px;">
                💡 <strong>Рекомендуется Tinkoff</strong> - оптимальный выбор для C2C переводов
            </p>
            {{end}}
        </div>

        <div class="bank-grid">
            {{range .AvailableBanks}}
            <div class="bank-card {{if and (eq .BankCode "tinkoff_card") (eq $.PaymentSystem "C2C")}}recommended{{end}}" onclick="selectBank('{{.BankCode}}')">
                <div class="bank-icon">
                    {{if eq .BankCode "sberbank"}}🏦
                    {{else if eq .BankCode "tinkoff_card"}}💳
                    {{else if eq .BankCode "tinkoff_phone"}}📱
                    {{else if eq .BankCode "vtb"}}🔵
                    {{else}}🏦{{end}}
                </div>
                <div class="bank-name">
                    {{.BankName}}
                    {{if and (eq .BankCode "tinkoff_card") (eq $.PaymentSystem "C2C")}}
                    <span class="recommended-badge">рекомендуется</span>
                    {{end}}
                </div>
                <div class="info-text">Нажмите для оплаты</div>
            </div>
            {{end}}
        </div>

        <div style="text-align: center; margin-top: 20px;">
            <p class="info-text">После выбора банка откроется приложение для завершения платежа</p>
        </div>
    </div>

    <script>
        function selectBank(bankCode) {
            // Показываем индикатор загрузки
            const card = event.currentTarget;
            const originalContent = card.innerHTML;
            card.style.background = '#f8f9fa';
            card.innerHTML = '<div style="padding: 20px;">⏳ Загрузка...</div>';
            
            // Перенаправляем на конкретный диплинк
            window.location.href = '/api/v1/payments/deeplink/specific?order_id={{.OrderID}}&bank=' + bankCode;
            
            // В случае ошибки возвращаем оригинальный контент
            setTimeout(() => {
                if (!document.hidden) {
                    card.innerHTML = originalContent;
                    card.style.background = 'white';
                    alert('Не удалось открыть приложение. Убедитесь, что приложение банка установлено.');
                }
            }, 3000);
        }

        // УБИРАЕМ автоматический выбор для C2C - пользователь должен выбирать сам
    </script>
</body>
</html>
`))

    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, data); err != nil {
        return "", err
    }
    return buf.String(), nil
}

func (ds *DeeplinkService) renderSpecificTemplate(tmpl *template.Template, data map[string]interface{}) (string, error) {
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, data); err != nil {
        return "", err
    }
    return buf.String(), nil
}