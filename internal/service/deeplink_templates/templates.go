package deeplink_templates

import "html/template"

// BankTemplateConfig конфигурация шаблона банка
type BankTemplateConfig struct {
    BankCode         string
    BankName         string
    Template         *template.Template
    SupportedSystems []string // Поддерживаемые платежные системы
    TransferType     string   // "card", "phone", "both"
}

// BankTemplates глобальная конфигурация всех шаблонов
var BankTemplates = make(map[string]BankTemplateConfig)

// RegisterTemplate регистрирует новый шаблон
func RegisterTemplate(config BankTemplateConfig) {
    BankTemplates[config.BankCode] = config
}

// GetTemplatesForSystem возвращает шаблоны для указанной платежной системы
func GetTemplatesForSystem(paymentSystem string) []BankTemplateConfig {
    var result []BankTemplateConfig
    for _, config := range BankTemplates {
        for _, system := range config.SupportedSystems {
            if system == paymentSystem || system == "all" {
                result = append(result, config)
                break
            }
        }
    }
    return result
}

// GetAllTemplates возвращает все зарегистрированные шаблоны
func GetAllTemplates() []BankTemplateConfig {
    var result []BankTemplateConfig
    for _, config := range BankTemplates {
        result = append(result, config)
    }
    return result
}