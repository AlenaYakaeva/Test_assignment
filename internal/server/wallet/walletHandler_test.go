package wallet

import (
	"encoding/json"
	"errors"
	walletDomain "exampleApp/internal/domain/wallet"
	"exampleApp/internal/mocks"
	"exampleApp/internal/repository/db"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	ErrNotFoundWallet = errors.New("wallet not found")
)

func TestGetBalanceByID(t *testing.T) {
	type want struct {
		statusCode int
		amount     float64
		err        error
	}

	type test struct {
		name   string
		method string
		url    string
		uuid   string
		want   want
	}

	tests := []test{
		{
			name:   "valid request",
			method: "GET",
			url:    "/api/v1/wallets/",
			uuid:   "0769c575-cc70-476d-881a-e401b28fef70",
			want: want{
				statusCode: http.StatusOK,
				amount:     1000.0,
				err:        nil,
			},
		},
		{
			name:   "invalid request",
			method: "GET",
			url:    "/api/v1/wallets/",
			uuid:   "0769c575-cc70-476d-881a-e401b28fef71",
			want: want{
				statusCode: http.StatusBadRequest,
				amount:     0.0,
				err:        ErrNotFoundWallet,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			walletServiceMock := mocks.NewMockWalletService(t)

			walletServiceMock.EXPECT().
				GetBalanceByID(tc.uuid).
				Return(tc.want.amount, tc.want.err)

			gin.SetMode(gin.ReleaseMode)
			router := gin.New()
			wh := NewWalletHandler(walletServiceMock)
			router.GET(tc.url+":id", wh.GetBalanceByID)

			httpSrv := httptest.NewServer(router)
			defer httpSrv.Close()

			req := resty.New().R()
			req.Method = tc.method
			req.URL = httpSrv.URL + tc.url + tc.uuid

			resp, err := req.Send()
			assert.NoError(t, err)

			assert.Equal(t, tc.want.statusCode, resp.StatusCode())

			if tc.want.err != nil {
				type response struct {
					Error string `json:"error"`
				}
				var respData response
				err = json.Unmarshal(resp.Body(), &respData)
				assert.NoError(t, err)

				assert.Equal(t, tc.want.err.Error(), respData.Error)
				return
			}

			type response struct {
				Amount float64 `json:"amount"`
			}
			var respData response
			err = json.Unmarshal(resp.Body(), &respData)
			assert.NoError(t, err)

			assert.Equal(t, tc.want.amount, respData.Amount)
		})
	}
}

func TestChageAmount(t *testing.T) {
	type want struct {
		statusCode int
		wallet     walletDomain.Wallet
		err        error
	}

	type test struct {
		name   string
		req    string
		method string
		url    string
		want   want
	}

	tests := []test{
		{
			name:   "deposit request",
			req:    `{"walletId":"0769c575-cc70-476d-881a-e401b28fef70","operationType":"DEPOSIT","amount":1000.0}`,
			method: "POST",
			url:    "/api/v1/wallet",
			want: want{
				statusCode: http.StatusOK,
				wallet: walletDomain.Wallet{
					WalletID: "0769c575-cc70-476d-881a-e401b28fef70",
					Amount:   2000.0,
				},
				err: nil,
			},
		},
		{
			name:   "create deposit request",
			req:    `{"walletId":"0769c575-cc70-476d-881a-e401b28fef71","operationType":"DEPOSIT","amount":1000.0}`,
			method: "POST",
			want: want{
				statusCode: http.StatusOK,
				wallet: walletDomain.Wallet{
					WalletID: "0769c575-cc70-476d-881a-e401b28fef86",
					Amount:   1000.0,
				},
				err: nil,
			},
		},
		{
			name:   "valid withdraw request",
			req:    `{"walletId":"0769c575-cc70-476d-881a-e401b28fef70","operationType":"WITHDRAW","amount":100.0}`,
			method: "POST",
			want: want{
				statusCode: http.StatusOK,
				wallet: walletDomain.Wallet{
					WalletID: "0769c575-cc70-476d-881a-e401b28fef70",
					Amount:   900.0,
				},
				err: nil,
			},
		},
		{
			name:   "invalid withdraw request not exists",
			req:    `{"walletId":"0769c575-cc70-476d-881a-e401b28fef71","operationType":"WITHDRAW","amount":100.0}`,
			method: "POST",
			want: want{
				statusCode: http.StatusBadRequest,
				wallet:     walletDomain.Wallet{},
				err:        db.ErrWithdrawWalletNotFound,
			},
		},
		{
			name:   "invalid withdraw request paiment faild",
			req:    `{"walletId":"0769c575-cc70-476d-881a-e401b28fef71","operationType":"WITHDRAW","amount":100.0}`,
			method: "POST",
			want: want{
				statusCode: http.StatusBadRequest,
				wallet:     walletDomain.Wallet{},
				err:        db.ErrPaymentFailed,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			walletServiceMock := mocks.NewMockWalletService(t)

			var reqObj walletDomain.WalletRequest
			err := json.Unmarshal([]byte(tc.req), &reqObj)
			assert.NoError(t, err)

			walletServiceMock.EXPECT().
				ChageAmount(mock.Anything, reqObj).
				Return(tc.want.wallet, tc.want.err)

			gin.SetMode(gin.ReleaseMode)
			router := gin.New()
			wh := NewWalletHandler(walletServiceMock)
			router.POST(tc.url, wh.ChageAmount)

			httpSrv := httptest.NewServer(router)
			defer httpSrv.Close()

			req := resty.New().R()
			req.Method = tc.method
			req.URL = httpSrv.URL + tc.url
			req.Body = tc.req

			resp, err := req.Send()
			assert.NoError(t, err)

			assert.Equal(t, tc.want.statusCode, resp.StatusCode())

			if tc.want.err != nil {
				type response struct {
					Error string `json:"error"`
				}
				var respData response
				err = json.Unmarshal(resp.Body(), &respData)
				assert.NoError(t, err)

				assert.Equal(t, tc.want.err.Error(), respData.Error)
				return
			}

			type response struct {
				WalletID string  `json:"walletID"`
				Amount   float64 `json:"amount"`
			}
			var respData response
			err = json.Unmarshal(resp.Body(), &respData)
			assert.NoError(t, err)

			respWallet := walletDomain.Wallet{
				WalletID: respData.WalletID,
				Amount:   respData.Amount,
			}

			assert.Equal(t, tc.want.wallet, respWallet)
		})
	}
}
