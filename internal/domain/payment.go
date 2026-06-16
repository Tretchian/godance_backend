package domain

// PaymentRepository — персистентность платежей (записи Payment).
// Будет наполнена в срезе payments.
type PaymentRepository interface{}

// PaymentGateway абстрагирует платёжный шлюз с двухфазной моделью hold/capture.
// Реальная интеграция появится позже; на данном этапе используется заглушка.
type PaymentGateway interface {
	// CreateHold инициирует удержание средств (escrow) и возвращает URL оплаты.
	CreateHold(relatedID uint, amount float64) (paymentURL string, err error)
	// Capture списывает ранее удержанные средства в пользу исполнителя.
	Capture(relatedID uint) error
	// Release возвращает удержание плательщику.
	Release(relatedID uint) error
}
