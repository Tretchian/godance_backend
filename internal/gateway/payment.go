package gateway

import (
	"fmt"
	"os"

	"godance/internal/domain"
)

// StubPayment — заглушка платёжного шлюза: формирует payment_url, но реальных
// операций со средствами не выполняет. Переходы статусов в системе драйвятся
// вебхуками /payments/webhooks/*, которые в проде шлёт настоящий шлюз.
// Заменяется реальной интеграцией в отдельной итерации.
type StubPayment struct{}

func NewStubPayment() *StubPayment { return &StubPayment{} }

var _ domain.PaymentGateway = (*StubPayment)(nil)

func (g *StubPayment) CreateHold(relatedID uint, amount float64) (string, error) {
	return fmt.Sprintf("%s/pay/feedback/%d", os.Getenv("PAYMENT_GATEWAY_URL"), relatedID), nil
}

func (g *StubPayment) Capture(relatedID uint) error { return nil }

func (g *StubPayment) Release(relatedID uint) error { return nil }
