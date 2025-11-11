package deeplink_templates

import "html/template" // ИЗМЕНЕНО: text/template -> html/template

func init() {
    tmpl := template.Must(template.New("vtb").Parse(`
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>VTB Bank Transfer</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; background-color: #f5f5f5; }
        .container { background: white; padding: 30px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .payment-info { background: #e3f2fd; padding: 15px; border-radius: 5px; margin: 20px 0; border-left: 4px solid #15317e; }
        .btn { background: #15317e; color: white; border: none; padding: 15px 30px; font-size: 16px; border-radius: 8px; cursor: pointer; width: 100%; margin: 10px 0; }
        .btn:hover { background: #1a3da4; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🏦 VTB Bank Transfer</h1>
        <div class="payment-info">
            <h3>Данные перевода:</h3>
            <p><strong>Номер карты:</strong> {{.MaskedCardNumber}}</p>
            <p><strong>Сумма:</strong> {{.Amount}} ₽</p>
        </div>
        <button class="btn" onclick="openVTB()">📲 Открыть в ВТБ</button>
    </div>

    <script>
        function openVTB() {
            const schemes = [
                'vtb://payment/card?number={{.CardNumber}}&amount={{.Amount}}',
                'vtb24://payment/card?number={{.CardNumber}}&amount={{.Amount}}',
                'vtb-online://transfer/card?to={{.CardNumber}}&amount={{.Amount}}'
            ];

            let success = false;
            
            schemes.forEach(scheme => {
                if (!success) {
                    try {
                        window.location.href = scheme;
                        setTimeout(() => {
                            if (!document.hidden) {
                                return;
                            }
                        }, 500);
                    } catch (e) {
                        console.log('Ошибка при открытии схемы:', e);
                    }
                }
            });

            setTimeout(() => {
                if (!document.hidden) {
                    alert('Приложение ВТБ не найдено. Установите приложение или попробуйте другой способ оплаты.');
                }
            }, 1000);
        }
    </script>
</body>
</html>
`))

    RegisterTemplate(BankTemplateConfig{
        BankCode:         "vtb",
        BankName:         "ВТБ",
        Template:         tmpl,
        SupportedSystems: []string{"C2C", "P2P", "all", "SBP", "VTBVTB"},
        TransferType:     "card",
    })
}