package lunchmoney

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Rhymond/go-money"
	"github.com/go-playground/validator/v10"
)

// TransactionsResponse is the response we get from requesting transactions.
type TransactionsResponse struct {
	Transactions []*Transaction `json:"transactions"`
}

// Transaction is a single LM transaction.
type Transaction struct {
	ID             int64  `json:"id"`
	Date           string `json:"date" validate:"omitempty,datetime=2006-01-02"`
	Payee          string `json:"payee"`
	Amount         string `json:"amount"`
	Currency       string `json:"currency"`
	Notes          string `json:"notes"`
	CategoryID     int64  `json:"category_id"`
	RecurringID    int64  `json:"recurring_id"`
	AssetID        int64  `json:"asset_id"`
	PlaidAccountID int64  `json:"plaid_account_id"`
	Status         string `json:"status"`
	IsGroup        bool   `json:"is_group"`
	GroupID        int64  `json:"group_id"`
	ParentID       int64  `json:"parent_id"`
	ExternalID     int64  `json:"external_id"`
}

// ParsedAmount turns the currency from lunchmoney into a Go currency.
func (t *Transaction) ParsedAmount() (*money.Money, error) {
	return ParseCurrency(t.Amount, t.Currency)
}

// TransactionFilters are options to pass into the request for transactions.
type TransactionFilters struct {
	TagID           int64  `json:"tag_id"`
	RecurringID     int64  `json:"recurring_id"`
	PlaidAccountID  int64  `json:"plaid_account_id"`
	CategoryID      int64  `json:"category_id"`
	AssetID         int64  `json:"asset_id"`
	Offset          int64  `json:"offset"`
	Limit           int64  `json:"limit"`
	StartDate       string `json:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate         string `json:"end_date" validate:"omitempty,datetime=2006-01-02"`
	DebitAsNegative bool   `json:"debit_as_negative"`
}

// ToMap converts the filters to a string map to be sent with the request as
// GET parameters.
func (r *TransactionFilters) ToMap() (map[string]string, error) {
	ret := map[string]string{}
	b, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &ret); err != nil {
		return nil, err
	}

	return ret, nil
}

// GetTransactions gets all transactions filtered by the filters.
func (c *Client) GetTransactions(ctx context.Context, filters *TransactionFilters) ([]*Transaction, error) {
	validate := validator.New()
	options := map[string]string{}
	if filters != nil {
		if err := validate.Struct(filters); err != nil {
			return nil, err
		}

		maps, err := filters.ToMap()
		if err != nil {
			return nil, fmt.Errorf("convert filters to map: %w", err)
		}
		options = maps
	}

	body, err := c.Get(ctx, "/v1/transactions", options)
	if err != nil {
		return nil, fmt.Errorf("get transactions: %w", err)
	}

	resp := &TransactionsResponse{}
	if err := json.NewDecoder(body).Decode(resp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if err := validate.Struct(resp); err != nil {
		return nil, err
	}

	return resp.Transactions, nil
}

// GetTransaction gets a transaction by id.
func (c *Client) GetTransaction(ctx context.Context, id int64, filters *TransactionFilters) (*Transaction, error) {
	validate := validator.New()
	options := map[string]string{}
	if filters != nil {
		if err := validate.Struct(filters); err != nil {
			return nil, err
		}

		maps, err := filters.ToMap()
		if err != nil {
			return nil, err
		}
		options = maps
	}

	body, err := c.Get(ctx, fmt.Sprintf("/v1/transactions/%d", id), options)
	if err != nil {
		return nil, fmt.Errorf("get transaction %d: %w", id, err)
	}

	resp := &Transaction{}
	if err := json.NewDecoder(body).Decode(resp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if err := validate.Struct(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// UpdateTransaction is the transaction to update.
type UpdateTransaction struct {
	Date        *string `json:"date,omitempty" validate:"omitnil,datetime=2006-01-02"`
	CategoryID  *int    `json:"category_id,omitempty"`
	Payee       *string `json:"payee,omitempty"`
	Currency    *string `json:"currency,omitempty"`
	AssetID     *int    `json:"asset_id,omitempty"`
	RecurringID *int    `json:"recurring_id,omitempty"`
	Notes       *string `json:"notes,omitempty"`
	Status      *string `json:"status,omitempty" validate:"omitnil,oneof=cleared uncleared"`
	ExternalID  *string `json:"external_id,omitempty"`
}

// UpdateRequest is the request to update a transaction.
type UpdateRequest struct {
	Transaction *UpdateTransaction `json:"transaction"`
}

// UpdateTransactionResp is the response we get from updating a transaction.
type UpdateTransactionResp struct {
	Updated bool  `json:"updated"`
	Split   []int `json:"split"`
}

func (c *Client) UpdateTransaction(ctx context.Context, id int64, ut *UpdateTransaction) (*UpdateTransactionResp, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(ut); err != nil {
		return nil, err
	}

	body, err := c.Put(ctx, fmt.Sprintf("/v1/transactions/%d", id), &UpdateRequest{Transaction: ut})
	if err != nil {
		return nil, fmt.Errorf("update transaction %d: %w", id, err)
	}

	resp := &UpdateTransactionResp{}
	if err := json.NewDecoder(body).Decode(resp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return resp, nil
}
