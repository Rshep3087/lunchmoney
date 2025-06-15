package lunchmoney

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/go-playground/validator/v10"
)

// CryptoResponse is a response to a crypto lookup.
type CryptoResponse struct {
	Crypto []*Crypto `json:"crypto"`
}

// Crypto is a single LM crypto asset.
type Crypto struct {
	ID              *int      `json:"id"`               // Unique identifier for a manual crypto account (no ID for synced accounts)
	ZaboAccountID   *int      `json:"zabo_account_id"`  // Unique identifier for a synced crypto account (no ID for manual accounts)
	Source          string    `json:"source"`           // Either "synced" or "manual"
	Name            string    `json:"name"`             // Name of the crypto asset
	DisplayName     *string   `json:"display_name"`     // Display name of the crypto asset (as set by user)
	Balance         string    `json:"balance"`          // Current balance
	BalanceAsOf     time.Time `json:"balance_as_of"`    // Date/time the balance was last updated
	Currency        string    `json:"currency"`         // Abbreviation for the cryptocurrency
	Status          string    `json:"status"`           // Current status of the crypto account (active or error)
	InstitutionName *string   `json:"institution_name"` // Name of provider holding the asset
	CreatedAt       time.Time `json:"created_at"`       // Date/time the asset was created
	ToBase          *float64  `json:"to_base"`          // The balance converted to the user's primary currency
}

// ParsedAmount converts the crypto asset's balance and currency into a money.Money object.
// This provides a convenient way to work with the crypto asset's value using the go-money library's
// currency handling capabilities. Returns an error if the balance cannot be parsed.
func (c *Crypto) ParsedAmount() (*money.Money, error) {
	return ParseCurrency(c.Balance, c.Currency)
}

// GetCrypto retrieves all crypto assets from the Lunch Money API.
// It returns a slice of Crypto objects containing information about each crypto asset,
// including balance, institution, and status details. Returns an error if the request fails.
func (c *Client) GetCrypto(ctx context.Context) ([]*Crypto, error) {
	validate := validator.New()
	options := map[string]string{}

	body, err := c.Get(ctx, "/v1/crypto", options)
	if err != nil {
		return nil, fmt.Errorf("get crypto: %w", err)
	}

	resp := &CryptoResponse{}
	if err := json.NewDecoder(body).Decode(resp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if err := validate.Struct(resp); err != nil {
		return nil, err
	}

	return resp.Crypto, nil
}

// UpdateCrypto contains the fields that can be updated for an existing manual crypto asset.
// Only non-nil fields will be sent in the update request.
type UpdateCrypto struct {
	Name            *string `json:"name,omitempty"`             // Official or full name of the account. Max 45 characters
	DisplayName     *string `json:"display_name,omitempty"`     // Display name for the account. Max 25 characters
	InstitutionName *string `json:"institution_name,omitempty"` // Name of provider that holds the account. Max 50 characters
	Balance         *string `json:"balance,omitempty"`          // Numeric value of the current balance
	Currency        *string `json:"currency,omitempty"`         // Cryptocurrency that is supported for manual tracking
}

// UpdateManualCrypto modifies an existing manual crypto asset with the specified ID using the provided fields.
// It returns the updated crypto asset information or an error if the update fails.
// Only fields that are non-nil in the crypto parameter will be updated.
// This only works for manually-managed crypto assets (source: manual).
func (c *Client) UpdateManualCrypto(ctx context.Context, id int64, crypto *UpdateCrypto) (*Crypto, error) {
	validate := validator.New()
	if err := validate.Struct(crypto); err != nil {
		return nil, err
	}

	body, err := c.Put(ctx, fmt.Sprintf("/v1/crypto/manual/%d", id), crypto)
	if err != nil {
		return nil, fmt.Errorf("put crypto %d: %w", id, err)
	}

	resp := &Crypto{}
	if err := json.NewDecoder(body).Decode(resp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return resp, nil
}
