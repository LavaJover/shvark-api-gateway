// internal/domain/deeplink.go
package domain

type DeeplinkData struct {
    HTMLContent string
    DeeplinkType string
    BankCode    string
    OrderID     string
}

type DeeplinkTemplate struct {
    BankCode string
    Template string
    Schemes  []DeeplinkScheme
}

type DeeplinkScheme struct {
    Name     string
    Template string
    OS       string // ios, android, universal
}