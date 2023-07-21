package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/config"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/sqldb"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	cfg := config.NewConfig()
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
			want: NewBaseHandler(nil, *cfg),
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			got := NewBaseHandler(test.args.dbConn, *cfg)
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

type testresponsewriter struct {
	w http.ResponseWriter
	// buf  bytes.Buffer
	code int
}

func (rw *testresponsewriter) Header() http.Header {
	return rw.w.Header()
}

func (rw *testresponsewriter) WriteHeader(statusCode int) {
	rw.code = statusCode
}

func (rw *testresponsewriter) Write( /*data*/ _ []byte) (int, error) {
	return 0, io.EOF

	/*return rw.buf.Write(data)*/
}

func (rw *testresponsewriter) Done() (int64, error) {
	return 0, io.EOF

	/*	if rw.code > 0 {
			rw.w.WriteHeader(rw.code)
		}
		return io.Copy(rw.w, &rw.buf)
	*/
}

func TestWrapHandlerErrIO(t *testing.T) {
	echoFr := echo.New()
	req := httptest.NewRequest(echo.POST, "http://localhost:1323/admin/user_points/settings", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoFr.NewContext(req, rec)
	ctx.Response().Writer = &testresponsewriter{w: rec} //nolint:exhaustruct

	type args struct {
		ctx        echo.Context
		statusCode int
		msg        string
		errIn      error
	}

	arguments := args{
		statusCode: http.StatusInternalServerError,
		ctx:        ctx,
		msg:        "error",
		errIn:      errIn1,
	}

	err := WrapHandlerErr(arguments.ctx, arguments.statusCode, arguments.msg, arguments.errIn)
	assert.Error(t, err)
	assert.ErrorIs(t, err, io.EOF)
}

func TestGetAuthFromCtxEmpty(t *testing.T) {
	echoFr := echo.New()
	defer func(echoFramework *echo.Echo) {
		err := echoFramework.Close()
		require.NoError(t, err)
	}(echoFr)

	req := httptest.NewRequest(echo.GET, "http://localhost:1323", nil)
	rec := httptest.NewRecorder()
	ctx := echoFr.NewContext(req, rec)
	got := GetAuthFromCtx(ctx)

	assert.Empty(t, got)

	ctx.Request().Header.Add("Authorization", "123")

	got = GetAuthFromCtx(ctx)
	assert.Empty(t, got)
}

func TestIsAuthorizedFalse(t *testing.T) {
	echoFr := echo.New()
	defer func(echoFramework *echo.Echo) {
		err := echoFramework.Close()
		require.NoError(t, err)
	}(echoFr)

	req := httptest.NewRequest(echo.GET, "http://localhost:1323", nil)
	rec := httptest.NewRecorder()
	ctx := echoFr.NewContext(req, rec)

	got := IsAuthorized(ctx, nil)
	assert.False(t, got)
}
