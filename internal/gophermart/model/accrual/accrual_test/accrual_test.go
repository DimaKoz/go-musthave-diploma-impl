package accrual_test

import (
	"testing"
	"time"

	accrual2 "github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/model/accrual"
	"github.com/stretchr/testify/assert"
)

const userAccrualTesting = "user1"

func TestNewOrderExt(t *testing.T) {
	now := time.Now()
	number := "123"
	status := accrual2.OrderStatusNew
	var accrual float32
	uploadedAt := now
	username := userAccrualTesting
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
	username := userAccrualTesting

	orderAcc := &accrual2.OrderAccrual{Order: number, Status: status, Accrual: accrual}
	want := &accrual2.OrderExt{
		Username: username, Number: number, Status: status, Accrual: accrual, UploadedAt: uploadedAt,
	}

	got := orderAcc.GetOrderExt(username, uploadedAt)
	assert.NotNil(t, got)
	assert.Equal(t, want, got)

	orderAcc = nil
	got = orderAcc.GetOrderExt(username, uploadedAt)
	assert.Nil(t, got)
}

func TestNewWithdrawExt(t *testing.T) {
	now := time.Now()
	number := "123"
	var sum float32 = 94.3
	username := userAccrualTesting
	want := &accrual2.WithdrawExt{
		Username: username, Order: number, Sum: sum, ProcessedAt: now,
	}

	got := accrual2.NewWithdrawExt(number, sum, now, username)
	assert.NotNil(t, got)
	assert.Equal(t, want, got)
}

func TestGetWithdrawExt(t *testing.T) {
	now := time.Now()
	number := "123"
	var sum float32 = 94.3
	username := userAccrualTesting
	want := &accrual2.WithdrawExt{
		Username: username, Order: number, Sum: sum, ProcessedAt: now,
	}

	withdrawAcc := &accrual2.WithdrawAccrual{Order: number, Sum: sum}

	got := withdrawAcc.GetWithdrawExt(username, now)
	assert.NotNil(t, got)
	assert.Equal(t, want, got)

	withdrawAcc = nil
	got = withdrawAcc.GetWithdrawExt(username, now)
	assert.Nil(t, got)
}
