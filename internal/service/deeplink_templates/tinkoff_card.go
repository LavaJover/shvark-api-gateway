package deeplink_templates

import "html/template"

func init() {
    tmpl := template.Must(template.New("tinkoff_card").Parse(`
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Tinkoff Bank C2C Card Transfer</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            background: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .payment-info {
            background: #f0f7ff;
            padding: 15px;
            border-radius: 5px;
            margin: 20px 0;
            border-left: 4px solid #ffdd2d;
        }
        .btn {
            background: #ffdd2d;
            color: #333;
            border: none;
            padding: 15px 30px;
            font-size: 16px;
            border-radius: 8px;
            cursor: pointer;
            width: 100%;
            margin: 10px 0;
            font-weight: bold;
        }
        .btn:hover {
            background: #f5d21c;
        }
        .btn-secondary {
            background: #6c757d;
            color: white;
        }
        .btn-secondary:hover {
            background: #5a6268;
        }
        .log {
            background: #1a1a1a;
            color: #00ff00;
            padding: 15px;
            border-radius: 5px;
            margin-top: 20px;
            font-family: 'Courier New', monospace;
            font-size: 12px;
            max-height: 300px;
            overflow-y: auto;
        }
        .status {
            padding: 10px;
            border-radius: 5px;
            margin: 10px 0;
            text-align: center;
        }
        .success { background: #d4edda; color: #155724; }
        .error { background: #f8d7da; color: #721c24; }
        .info { background: #d1ecf1; color: #0c5460; }
        .scheme-grid {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
            gap: 10px;
            margin-top: 15px;
        }
        .scheme-btn {
            background: #e9ecef;
            border: 1px solid #dee2e6;
            padding: 10px;
            border-radius: 5px;
            cursor: pointer;
            text-align: center;
            font-size: 12px;
        }
        .scheme-btn:hover {
            background: #d1ecf1;
        }
        .card-number {
            font-family: 'Courier New', monospace;
            letter-spacing: 1px;
            background: #f8f9fa;
            padding: 5px 10px;
            border-radius: 3px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>💳 Tinkoff Bank C2C Card Transfer</h1>
        
		<div class="payment-info">
		    <h3>Данные перевода:</h3>
		    <p><strong>Номер карты:</strong> <span class="card-number">{{.CardNumber}}</span></p> <!-- Убрали маскировку -->
		    <p><strong>Сумма:</strong> {{.Amount}} ₽</p>
		</div>

        <div id="status"></div>

        <button class="btn" id="tryAllBtn">
            🔍 Найти работающий deeplink Tinkoff
        </button>

        <button class="btn btn-secondary" id="manualBtn">
            🎯 Вручную протестировать схемы
        </button>

        <div class="log" id="log"></div>

        <div id="manualTest" style="display: none; margin-top: 20px;">
            <h3>Ручное тестирование схем Tinkoff (C2C):</h3>
            <div class="scheme-grid" id="schemeButtons"></div>
        </div>
    </div>

    <script>
        // Параметры для Tinkoff C2C перевода по карте
        const params = {
            cardNumber: '{{.CardNumber}}',
            amount: '{{.Amount}}'
        };

        // Все схемы Tinkoff
        const schemes = [
            'catch',
            'freelancecase',
            'yourmoney',
            'tinkoffbank',  // Основная схема Tinkoff
            'tbank',
            'wheels',
            'clanstrix',
            'feedaways',
            'toffice',
            'tguard',
            'shuttersmart',
            'petraise',
            'mobtrs',
            'goaloriented',
            'tmydocs',
            'tfinstudy',
            'tsplit',
            'tfinskills',
            'bank100000000004',
            'tassets',
            'tdata',
            'smarthome',
            'divevector',
            'framedit',
            'outpharmas',
            'yellowt',
            'invault',
            'ressinside',
            'youreporter',
            'plantu',
            'temperology',
            'logapp'
        ];

        let currentIndex = 0;
        let workingSchemes = [];
        let isTesting = false;

        // Элементы DOM
        const logElement = document.getElementById('log');
        const statusElement = document.getElementById('status');
        const tryAllBtn = document.getElementById('tryAllBtn');
        const manualBtn = document.getElementById('manualBtn');
        const manualTest = document.getElementById('manualTest');
        const schemeButtons = document.getElementById('schemeButtons');

        // Логирование
        function log(message, type = 'info') {
            const timestamp = new Date().toLocaleTimeString();
            const logEntry = document.createElement('div');
            logEntry.innerHTML = '<span style="color: #888">[' + timestamp + ']</span> ' + message;
            logElement.appendChild(logEntry);
            logElement.scrollTop = logElement.scrollHeight;
            console.log(message);
        }

        // Обновление статуса
        function updateStatus(message, type = 'info') {
            statusElement.innerHTML = '<div class="status ' + type + '">' + message + '</div>';
        }

        // Создание deeplink из схемы
        function createDeeplink(scheme) {
            return scheme + '://Main/Pay/C2C?amount=' + params.amount + '&targetCardNumber=' + params.cardNumber + '&numberCard=' + params.cardNumber;
        }

        // Попытка открыть deeplink
        function tryOpenDeeplink(deepLink, schemeName) {
            return new Promise((resolve) => {
                log('Попытка: ' + schemeName, 'info');
                log('Ссылка: ' + deepLink, 'info');

                let appOpened = false;
                const timeout = 2000;

                // Слушаем потерю фокуса (признак того, что приложение открылось)
                window.addEventListener('blur', function onBlur() {
                    appOpened = true;
                    window.removeEventListener('blur', onBlur);
                    log('✅ Сработало! Окно потеряло фокус - приложение открылось', 'success');
                    resolve(true);
                });

                // Пытаемся открыть deeplink
                try {
                    window.location.href = deepLink;
                } catch (error) {
                    log('❌ Ошибка: ' + error.message, 'error');
                    resolve(false);
                }

                // Проверяем через таймаут
                setTimeout(() => {
                    if (!appOpened) {
                        log('⏰ Таймаут - приложение не ответило', 'error');
                        resolve(false);
                    }
                }, timeout);
            });
        }

        // Автоматический перебор всех схем
        async function tryAllSchemes() {
            if (isTesting) return;
            
            isTesting = true;
            currentIndex = 0;
            workingSchemes = [];
            logElement.innerHTML = '';
            
            updateStatus('🔍 Начинаю поиск работающих схем Tinkoff для C2C...', 'info');
            tryAllBtn.disabled = true;
            manualBtn.disabled = true;

            for (let i = 0; i < schemes.length; i++) {
                const scheme = schemes[i];
                const deepLink = createDeeplink(scheme);

                log('[' + (i + 1) + '/' + schemes.length + '] Тестирую: ' + scheme, 'info');

                const success = await tryOpenDeeplink(deepLink, scheme);
                
                if (success) {
                    workingSchemes.push(scheme);
                    updateStatus('✅ Найдена работающая схема: ' + scheme, 'success');
                    break;
                }

                // Задержка между попытками
                await new Promise(resolve => setTimeout(resolve, 1000));
            }

            if (workingSchemes.length === 0) {
                updateStatus('❌ Ни одна схема не сработала', 'error');
                log('💡 Совет: Убедитесь, что приложение Тинькофф установлено', 'info');
                log('💡 Возможно, нужны другие параметры для C2C перевода', 'info');
            } else {
                updateStatus('🎉 Найдено работающих схем: ' + workingSchemes.join(', '), 'success');
            }

            tryAllBtn.disabled = false;
            manualBtn.disabled = false;
            isTesting = false;
        }

        // Создание кнопок для ручного тестирования
        function createManualTestButtons() {
            schemeButtons.innerHTML = '';
            schemes.forEach((scheme, index) => {
                const deepLink = createDeeplink(scheme);
                
                const button = document.createElement('div');
                button.className = 'scheme-btn';
                button.innerHTML = '<div><strong>' + (index + 1) + '. ' + scheme + '</strong></div><div style="font-size: 10px; color: #666; margin-top: 5px;">Нажмите для теста</div>';
                button.onclick = function(e) {
                    e.preventDefault();
                    log('🧪 Ручной тест: ' + scheme, 'info');
                    log('📋 Полная ссылка: ' + deepLink, 'info');
                    window.location.href = deepLink;
                };
                
                schemeButtons.appendChild(button);
            });
        }

        // Альтернативные форматы deeplink для Tinkoff C2C
        function testAlternativeFormats() {
            log('🔄 Тестирование альтернативных форматов deeplink для C2C...', 'info');
            
            const alternativeFormats = [
                'tinkoffbank://transfer/card?card=' + params.cardNumber + '&amount=' + params.amount,
                'tinkoffbank://pay/card?number=' + params.cardNumber + '&sum=' + params.amount,
                'tinkoffbank://payment/card?card=' + params.cardNumber + '&amount=' + params.amount,
                'tbank://transfer/card?card=' + params.cardNumber + '&amount=' + params.amount,
                'tinkoffbank://start?screen=transfer&card=' + params.cardNumber + '&amount=' + params.amount,
                'tinkoffbank://c2c?card=' + params.cardNumber + '&amount=' + params.amount,
                'tinkoffbank://transfer?type=card&to=' + params.cardNumber + '&amount=' + params.amount
            ];

            // Создаем кнопки для альтернативных форматов
            const altContainer = document.createElement('div');
            altContainer.style.marginTop = '20px';
            altContainer.innerHTML = '<h4>Альтернативные форматы C2C:</h4>';
            
            const altGrid = document.createElement('div');
            altGrid.className = 'scheme-grid';
            
            alternativeFormats.forEach((format, index) => {
                const btn = document.createElement('div');
                btn.className = 'scheme-btn';
                btn.style.background = '#fff3cd';
                btn.innerHTML = '<div><strong>Альт. ' + (index + 1) + '</strong></div><div style="font-size: 9px; color: #666; margin-top: 3px; word-break: break-all;">' + format.substring(0, 50) + '...</div>';
                btn.onclick = function() {
                    log('🔧 Альтернативный формат C2C ' + (index + 1), 'info');
                    log('📋 Ссылка: ' + format, 'info');
                    window.location.href = format;
                };
                altGrid.appendChild(btn);
            });

            altContainer.appendChild(altGrid);
            manualTest.appendChild(altContainer);
        }

        // Инициализация
        function init() {
            log('Страница инициализирована. Нажмите кнопку для тестирования C2C deeplink Tinkoff.', 'info');
            log('Параметры перевода: карта ' + params.cardNumber + ', сумма ' + params.amount, 'info');
            
            tryAllBtn.addEventListener('click', function(e) {
                e.preventDefault();
                tryAllSchemes();
            });

            manualBtn.addEventListener('click', function(e) {
                e.preventDefault();
                manualTest.style.display = 'block';
                createManualTestButtons();
                testAlternativeFormats();
                log('Режим ручного тестирования Tinkoff C2C активирован', 'info');
            });
        }

        // Запуск при загрузке
        window.addEventListener('load', init);
    </script>
</body>
</html>
`))

    RegisterTemplate(BankTemplateConfig{
        BankCode:         "tinkoff_card",
        BankName:         "Тинькофф (C2C по карте)",
        Template:         tmpl,
        SupportedSystems: []string{"C2C", "TINKOFF"},
        TransferType:     "card",
    })
}