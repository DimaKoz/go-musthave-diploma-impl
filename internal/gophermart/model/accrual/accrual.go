package accrual

import "time"

const (
	OrderStatusNew        = "NEW"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusInvalid    = "INVALID"
	OrderStatusProcessed  = "PROCESSED"
)

type OrderAccrual struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual"`
}

type OrderExt struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float32   `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"` //nolint:tagliatelle
	Username   string    `json:"-"`
}

func NewOrderExt(number string, status string, accrual float32, uploadedAt time.Time, username string) *OrderExt {
	return &OrderExt{
		Number:     number,
		Status:     status,
		Accrual:    accrual,
		UploadedAt: uploadedAt,
		Username:   username,
	}
}

func (orInternal *OrderAccrual) GetOrderExt(username string, time time.Time) *OrderExt {
	if orInternal == nil {
		return nil
	}

	return NewOrderExt(orInternal.Order, orInternal.Status, orInternal.Accrual, time, username)
}
