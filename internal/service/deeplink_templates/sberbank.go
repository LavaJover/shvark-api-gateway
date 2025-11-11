package deeplink_templates

import "html/template" // ИЗМЕНЕНО: text/template -> html/template

func init() {
    tmpl := template.Must(template.New("sberbank").Parse(`
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Sberbank P2P Payment</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; background-color: #f5f5f5; }
        .container { background: white; padding: 30px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .payment-info { background: #f8f9fa; padding: 15px; border-radius: 5px; margin: 20px 0; border-left: 4px solid #22a053; }
        .btn { background: #22a053; color: white; border: none; padding: 15px 30px; font-size: 16px; border-radius: 8px; cursor: pointer; width: 100%; margin: 10px 0; }
        .btn:hover { background: #1d8a47; }
        .log { background: #1a1a1a; color: #00ff00; padding: 15px; border-radius: 5px; margin-top: 20px; font-family: 'Courier New', monospace; font-size: 12px; max-height: 300px; overflow-y: auto; }
        .status { padding: 10px; border-radius: 5px; margin: 10px 0; text-align: center; }
        .success { background: #d4edda; color: #155724; }
        .error { background: #f8d7da; color: #721c24; }
    </style>
</head>
<body>
    <div class="container">
        <h1>💳 Sberbank P2P Payment</h1>
        <div class="payment-info">
            <h3>Данные перевода:</h3>
            <p><strong>Номер карты:</strong> {{.MaskedCardNumber}}</p>
            <p><strong>Сумма:</strong> {{.Amount}} ₽</p>
        </div>
        <div id="status"></div>
        <button class="btn" id="tryAllBtn">🔍 Найти работающий deeplink</button>
        <div class="log" id="log"></div>
    </div>

    <script>
        const params = { cardNumber: '{{.CardNumber}}', amount: '{{.Amount}}' };
        const schemes = [
            'sberbankonline://payments/p2p?type=card_number&requisiteNumber={cardNumber}&amount={amount}',
            'sbolonline://payments/p2p?type=card_number&requisiteNumber={cardNumber}&amount={amount}',
            'onlineappmobile://sbolonline/payments/p2p?type=card_number&requisiteNumber={cardNumber}&amount={amount}',
            'intent://ru.sberbankmobile/payments/p2p?type=card_number&requisiteNumber={cardNumber}&amount={amount}'
        ];

        let isTesting = false;
        const logElement = document.getElementById('log');
        const statusElement = document.getElementById('status');

        function log(message) {
            const timestamp = new Date().toLocaleTimeString();
            const logEntry = document.createElement('div');
            logEntry.innerHTML = '<span style="color: #888">[' + timestamp + ']</span> ' + message;
            logElement.appendChild(logEntry);
            logElement.scrollTop = logElement.scrollHeight;
        }

        function updateStatus(message, type = 'info') {
            statusElement.innerHTML = '<div class="status ' + type + '">' + message + '</div>';
        }

        function createDeeplink(schemeTemplate) {
            return schemeTemplate.replace('{cardNumber}', params.cardNumber).replace('{amount}', params.amount);
        }

        function tryOpenDeeplink(deepLink, schemeName) {
            return new Promise((resolve) => {
                log('Попытка: ' + schemeName);
                let appOpened = false;
                const timeout = 2000;

                window.addEventListener('blur', function onBlur() {
                    appOpened = true;
                    window.removeEventListener('blur', onBlur);
                    log('✅ Сработало! Приложение открылось');
                    resolve(true);
                });

                try {
                    window.location.href = deepLink;
                } catch (error) {
                    log('❌ Ошибка: ' + error.message);
                    resolve(false);
                }

                setTimeout(() => {
                    if (!appOpened) {
                        log('⏰ Таймаут - приложение не ответило');
                        resolve(false);
                    }
                }, timeout);
            });
        }

        async function tryAllSchemes() {
            if (isTesting) return;
            isTesting = true;
            logElement.innerHTML = '';
            updateStatus('🔍 Начинаю поиск работающих схем...', 'info');
            document.getElementById('tryAllBtn').disabled = true;

            for (let i = 0; i < schemes.length; i++) {
                const schemeTemplate = schemes[i];
                const schemeName = schemeTemplate.split('://')[0];
                const deepLink = createDeeplink(schemeTemplate);

                log('[' + (i + 1) + '/' + schemes.length + '] Тестирую: ' + schemeName);
                const success = await tryOpenDeeplink(deepLink, schemeName);
                
                if (success) {
                    updateStatus('✅ Найдена работающая схема: ' + schemeName, 'success');
                    break;
                }
                await new Promise(resolve => setTimeout(resolve, 1000));
            }

            document.getElementById('tryAllBtn').disabled = false;
            isTesting = false;
        }

        document.getElementById('tryAllBtn').addEventListener('click', tryAllSchemes);
        log('Страница инициализирована. Карта: ' + params.cardNumber + ', сумма: ' + params.amount);
    </script>
</body>
</html>
`))

    RegisterTemplate(BankTemplateConfig{
        BankCode:         "sberbank",
        BankName:         "Сбербанк",
        Template:         tmpl,
        SupportedSystems: []string{"C2C", "P2P"},
        TransferType:     "card",
    })
}