package accrual_test

import (
	"testing"
	"time"

	accrual2 "github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/model/accrual"
	"github.com/stretchr/testify/assert"
)

func TestNewOrderExt(t *testing.T) {
	now := time.Now()
	number := "123"
	status := accrual2.OrderStatusNew
	var accrual float32
	uploadedAt := now
	username := "user1"
	want := &accrual2.OrderExt{
		Username: username, Number: number, Status: status, Accrual: accrual, UploadedAt: now,
	}

	got := accrual2.NewOrderExt(number, status, accrual, uploadedAt, username)
	assert.NotNil(t, got)
	assert.Equal(t, want, got)
}

func TestGetOrderExt(t *testing.T) {
	now := time.Now()
	number := "123"
	status := accrual2.OrderStatusNew
	var accrual float32 = 42
	uploadedAt := now
	username := "user1"

	orderAcc := accrual2.OrderAccrual{Order: number, Status: status, Accrual: accrual}
	want := &accrual2.OrderExt{
		Username: username, Number: number, Status: status, Accrual: accrual, UploadedAt: uploadedAt,
	}

	got := orderAcc.GetOrderExt(username, uploadedAt)
	assert.NotNil(t, got)
	assert.Equal(t, want, got)
}
