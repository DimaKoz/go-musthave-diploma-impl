package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/sqldb"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAddAuthHeaders(t *testing.T) {
	echoFr := echo.New()
	req := httptest.NewRequest(echo.GET, "http://localhost:1323/admin/user_points/settings", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	ctx := echoFr.NewContext(req, rec)
	defer echoFr.Close()

	want := "Authorization:[asdf]"
	message := "asdf"

	AddAuthHeaders(ctx, message)

	gotAuth1 := ctx.Response().Header().Get("Authorization")
	gotAuth2 := ctx.Response().Header().Get("Set-Cookie")
	assert.Equal(t, want, gotAuth1)
	assert.Equal(t, want, gotAuth2)
}

func TestNewBaseHandler(t *testing.T) {
	type args struct {
		dbConn *sqldb.PgxIface
	}
	tests := []struct {
		name string
		args args
		want *BaseHandler
	}{
		{
			name: "nil dbConn",
			args: args{dbConn: nil},
			want: NewBaseHandler(nil),
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			got := NewBaseHandler(test.args.dbConn)
			assert.NotNil(t, got)
			assert.Equal(t, test.want, got)
		})
	}
}

var errIn1 = fmt.Errorf("error")

func TestWrapHandlerErr(t *testing.T) {
	echoFr := echo.New()
	req := httptest.NewRequest(echo.GET, "http://localhost:1323/admin/user_points/settings", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	ctx := echoFr.NewContext(req, rec)
	defer echoFr.Close()
	type args struct {
		ctx        echo.Context
		statusCode int
		msg        string
		errIn      error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test 500",
			args: args{
				statusCode: http.StatusInternalServerError,
				ctx:        ctx,
				msg:        "error",
				errIn:      errIn1,
			},
			wantErr: false,
		},
	}
	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			err := WrapHandlerErr(tt.args.ctx, tt.args.statusCode, tt.args.msg, tt.args.errIn)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
