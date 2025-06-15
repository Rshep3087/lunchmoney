package lunchmoney

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCrypto(t *testing.T) {
	tests := []struct {
		name        string
		response    string
		statusCode  int
		wantErr     bool
		errContains string
		want        []*Crypto
	}{
		{
			name: "successful response with mixed crypto assets",
			response: `{
				"crypto": [
					{
						"zabo_account_id": 544,
						"source": "synced",
						"created_at": "2020-07-27T11:53:02.722Z",
						"name": "Dogecoin",
						"display_name": null,
						"balance": "1.902383849000000000",
						"balance_as_of": "2021-05-21T00:05:36.000Z",
						"currency": "doge",
						"status": "active",
						"institution_name": "MetaMask",
						"to_base": 0.25
					},
					{
						"id": 152,
						"source": "manual",
						"created_at": "2021-04-03T04:16:48.230Z",
						"name": "Ether",
						"display_name": "BlockFi - ETH",
						"balance": "5.391445130000000000",
						"balance_as_of": "2021-05-20T16:57:00.000Z",
						"currency": "ETH",
						"status": "active",
						"institution_name": "BlockFi",
						"to_base": 12000.50
					}
				]
			}`,
			statusCode: http.StatusOK,
			want: []*Crypto{
				{
					ID:              nil,
					ZaboAccountID:   ptr(544),
					Source:          "synced",
					Name:            "Dogecoin",
					DisplayName:     nil,
					Balance:         "1.902383849000000000",
					BalanceAsOf:     time.Date(2021, 5, 21, 0, 5, 36, 0, time.UTC),
					Currency:        "doge",
					Status:          "active",
					InstitutionName: ptr("MetaMask"),
					CreatedAt:       time.Date(2020, 7, 27, 11, 53, 2, 722000000, time.UTC),
					ToBase:          ptr(0.25),
				},
				{
					ID:              ptr(152),
					ZaboAccountID:   nil,
					Source:          "manual",
					Name:            "Ether",
					DisplayName:     ptr("BlockFi - ETH"),
					Balance:         "5.391445130000000000",
					BalanceAsOf:     time.Date(2021, 5, 20, 16, 57, 0, 0, time.UTC),
					Currency:        "ETH",
					Status:          "active",
					InstitutionName: ptr("BlockFi"),
					CreatedAt:       time.Date(2021, 4, 3, 4, 16, 48, 230000000, time.UTC),
					ToBase:          ptr(12000.50),
				},
			},
		},
		{
			name: "empty response",
			response: `{
				"crypto": []
			}`,
			statusCode: http.StatusOK,
			want:       []*Crypto{},
		},
		{
			name:        "invalid JSON response",
			response:    `{"invalid": "json"`,
			statusCode:  http.StatusOK,
			wantErr:     true,
			errContains: "decode response",
		},
		{
			name:        "HTTP error",
			response:    `{"error": "Unauthorized"}`,
			statusCode:  http.StatusUnauthorized,
			wantErr:     true,
			errContains: "get crypto",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/v1/crypto", r.URL.Path)
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(tt.statusCode)
				_, err := w.Write([]byte(tt.response))
				require.NoError(t, err)
			}))
			defer server.Close()

			client, err := NewClient("test-token")
			require.NoError(t, err)
			client.Base, err = url.Parse(server.URL)
			require.NoError(t, err)

			got, err := client.GetCrypto(context.Background())
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestUpdateManualCrypto(t *testing.T) {
	tests := []struct {
		name        string
		id          int64
		update      *UpdateCrypto
		response    string
		statusCode  int
		wantErr     bool
		errContains string
		want        *Crypto
	}{
		{
			name: "successful update",
			id:   152,
			update: &UpdateCrypto{
				Name:        ptr("Ethereum"),
				DisplayName: ptr("Updated ETH"),
				Balance:     ptr("10.0"),
			},
			response: `{
				"id": 152,
				"source": "manual",
				"created_at": "2021-02-10T05:57:34.305Z",
				"name": "Ethereum",
				"display_name": "Updated ETH",
				"balance": "10.000000000000000000",
				"balance_as_of": "2021-05-20T16:57:00.000Z",
				"currency": "ETH",
				"status": "active",
				"institution_name": null,
				"to_base": 25000.0
			}`,
			statusCode: http.StatusOK,
			want: &Crypto{
				ID:              ptr(152),
				ZaboAccountID:   nil,
				Source:          "manual",
				Name:            "Ethereum",
				DisplayName:     ptr("Updated ETH"),
				Balance:         "10.000000000000000000",
				BalanceAsOf:     time.Date(2021, 5, 20, 16, 57, 0, 0, time.UTC),
				Currency:        "ETH",
				Status:          "active",
				InstitutionName: nil,
				CreatedAt:       time.Date(2021, 2, 10, 5, 57, 34, 305000000, time.UTC),
				ToBase:          ptr(25000.0),
			},
		},
		{
			name:        "invalid JSON response",
			id:          152,
			update:      &UpdateCrypto{Name: ptr("Bitcoin")},
			response:    `{"invalid": "json"`,
			statusCode:  http.StatusOK,
			wantErr:     true,
			errContains: "decode response",
		},
		{
			name:        "HTTP error",
			id:          152,
			update:      &UpdateCrypto{Name: ptr("Bitcoin")},
			response:    `{"errors": ["currency is invalid for crypto: fakecoin"]}`,
			statusCode:  http.StatusBadRequest,
			wantErr:     true,
			errContains: "put crypto",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/v1/crypto/manual/152", r.URL.Path)
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(tt.statusCode)
				_, err := w.Write([]byte(tt.response))
				require.NoError(t, err)
			}))
			defer server.Close()

			client, err := NewClient("test-token")
			require.NoError(t, err)
			client.Base, err = url.Parse(server.URL)
			require.NoError(t, err)

			got, err := client.UpdateManualCrypto(context.Background(), tt.id, tt.update)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestCryptoParsedAmount(t *testing.T) {
	tests := []struct {
		name     string
		crypto   *Crypto
		wantErr  bool
		wantCode string
	}{
		{
			name: "valid crypto amount",
			crypto: &Crypto{
				Balance:  "1.5",
				Currency: "BTC",
			},
			wantCode: "BTC",
		},
		{
			name: "invalid balance",
			crypto: &Crypto{
				Balance:  "invalid",
				Currency: "ETH",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.crypto.ParsedAmount()
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantCode, got.Currency().Code)
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}
