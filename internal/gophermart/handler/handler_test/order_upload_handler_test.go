package handler_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/config"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/handler"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/sqldb"
	"github.com/labstack/echo/v4"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testArgsUploadOrder struct {
	time       time.Time
	orderRow   string
	orderQuery string
	user       string
	accrual    float32
	status     string
}

func setupMock(mock pgxmock.PgxConnIface, args *testArgsUploadOrder, err error) {
	rows := pgxmock.NewRows([]string{"number", "status", "accrual", "username", "uploaded_at"})
	if args != nil {
		rows.AddRow(args.orderRow, args.status, args.accrual, args.user, args.time)
	}

	eQuery := mock.ExpectQuery(
		"SELECT number, status, accrual, username, uploaded_at FROM orders WHERE number=\\$1")
	if args != nil {
		eQuery.WithArgs(args.orderQuery)
	}
	if err == nil {
		eQuery.WillReturnRows(rows)
	} else {
		eQuery.WillReturnError(err)
	}
}

func getEchoMockCtxUploadHandler(
	echoFr *echo.Echo,
	order string,
	user string,
) (*echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(echo.POST, "http://localhost:1323/", strings.NewReader(order))
	req.Header.Add("Authorization",
		fmt.Sprintf("Authorization:[%s]", user))

	rec := httptest.NewRecorder()
	ctx := echoFr.NewContext(req, rec)

	return &ctx, rec
}

func TestOrderUploadHandlerStatus200(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	defer closeMockDB(t, mock)

	arg := &testArgsUploadOrder{ //nolint:exhaustruct
		orderRow:   "79927398713",
		orderQuery: "79927398713",
		user:       "user1",
		time:       time.Now(),
		status:     "NEW",
	}
	setupMock(mock, arg, nil)

	var pgConn sqldb.PgxIface = mock

	// Mock echo
	echoFr := echo.New()
	defer echoFr.Close()

	ctx, rec := getEchoMockCtxUploadHandler(echoFr, "79927398713", "user1")

	cfg := config.NewConfig()
	baseH := handler.NewBaseHandler(&pgConn, *cfg)

	err = baseH.OrderUploadHandler(*ctx)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	got := rec.Result()
	defer got.Body.Close()

	wantStatusCode := http.StatusOK
	assert.Equal(t, wantStatusCode, got.StatusCode, "StatusCode got: %v, want: %v", got.StatusCode, wantStatusCode)
}

func closeMockDB(t *testing.T, mock pgxmock.PgxConnIface) {
	t.Helper()
	mock.ExpectClose()
	err := mock.Close(context.Background())
	require.NoError(t, err)
}

func TestOrderUploadHandlerStatusConflict(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	defer closeMockDB(t, mock)

	arg := &testArgsUploadOrder{ //nolint:exhaustruct
		orderRow:   "2377225624",
		orderQuery: "2377225624",
		user:       "user1",
		time:       time.Now(),
		status:     "NEW",
	}
	setupMock(mock, arg, nil)

	var pgConn sqldb.PgxIface = mock

	// Mock echo
	echoFr := echo.New()
	defer echoFr.Close()

	bodyStr := "2377225624" // or 7142326
	req := httptest.NewRequest(echo.POST, "http://localhost:1323/", strings.NewReader(bodyStr))
	req.Header.Add("Authorization",
		fmt.Sprintf("Authorization:[%s]", "user2"))

	rec := httptest.NewRecorder()
	ctx := echoFr.NewContext(req, rec)

	cfg := config.NewConfig()
	baseH := handler.NewBaseHandler(&pgConn, *cfg)

	err = baseH.OrderUploadHandler(ctx)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	got := rec.Result()
	defer got.Body.Close()

	wantStatusCode := http.StatusConflict
	assert.Equal(t, wantStatusCode, got.StatusCode, "StatusCode got: %v, want: %v", got.StatusCode, wantStatusCode)
}

func TestOrderUploadHandlerStatusInternalErr(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	defer closeMockDB(t, mock)

	rows := pgxmock.NewRows([]string{"number", "status", "accrual", "username", "uploaded_at"})

	mock.ExpectQuery(
		"SELECT number, status, accrual, username, uploaded_at FROM orders WHERE number=\\$1").
		WithArgs("79927398713").
		WillReturnRows(rows)

	var pgConn sqldb.PgxIface = mock

	// Mock echo
	echoFr := echo.New()
	defer echoFr.Close()

	ctx, rec := getEchoMockCtxUploadHandler(echoFr, "79927398713", "user1")

	cfg := config.NewConfig()
	baseH := handler.NewBaseHandler(&pgConn, *cfg)

	err = baseH.OrderUploadHandler(*ctx)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	got := rec.Result()
	defer got.Body.Close()

	wantStatusCode := http.StatusInternalServerError
	assert.Equal(t, wantStatusCode, got.StatusCode, "StatusCode got: %v, want: %v", got.StatusCode, wantStatusCode)
}

func TestOrderUploadHandlerStatusAccepted(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	defer closeMockDB(t, mock)

	rows := pgxmock.NewRows([]string{"number", "status", "accrual", "username", "uploaded_at"})

	mock.ExpectQuery(
		"SELECT number, status, accrual, username, uploaded_at FROM orders WHERE number=\\$1").
		WithArgs("79927398713").
		WillReturnRows(rows)

	mock.ExpectExec("insert into orders").
		WithArgs("79927398713", "NEW", float32(0), "user1", pgxmock.AnyArg()).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	var pgConn sqldb.PgxIface = mock

	// Mock echo
	echoFr := echo.New()
	defer echoFr.Close()

	ctx, rec := getEchoMockCtxUploadHandler(echoFr, "79927398713", "user1")

	cfg := config.NewConfig()
	baseH := handler.NewBaseHandler(&pgConn, *cfg)

	err = baseH.OrderUploadHandler(*ctx)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	got := rec.Result()
	defer got.Body.Close()

	wantStatusCode := http.StatusAccepted
	assert.Equal(t, wantStatusCode, got.StatusCode, "StatusCode got: %v, want: %v", got.StatusCode, wantStatusCode)
}
