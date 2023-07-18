package accrual

import "time"

const (
	OrderStatusNew        = "NEW"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusInvalid    = "INVALID"
	OrderStatusProcessed  = "PROCESSED"
)

type BalanceExt struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}

type OrderAccrual struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual"`
}

type WithdrawAccrual struct {
	Order string  `json:"order"`
	Sum   float32 `json:"sum"`
}

type WithdrawExt struct {
	Order       string    `json:"order"`
	Sum         float32   `json:"accrual"`
	ProcessedAt time.Time `json:"processed_at"` //nolint:tagliatelle
	Username    string    `json:"-"`
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

func NewWithdrawExt(number string, sum float32, processedAt time.Time, username string) *WithdrawExt {
	return &WithdrawExt{
		Order:       number,
		Sum:         sum,
		ProcessedAt: processedAt,
		Username:    username,
	}
}

func (wInternal *WithdrawAccrual) GetWithdrawExt(username string, time time.Time) *WithdrawExt {
	if wInternal == nil {
		return nil
	}

	return NewWithdrawExt(wInternal.Order, wInternal.Sum, time, username)
}
