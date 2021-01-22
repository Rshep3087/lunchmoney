package lunchmoney

import "time"

// CategoriesResponse is the response we get from requesting categories.
type CategoriesResponse struct {
	Categories []*Category `json:"categories"`
	Error      string      `json:"error"`
}

// Category is a single LM category.
type Category struct {
	ID                int64     `json:"id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	IsIncome          bool      `json:"is_income"`
	ExcludeFromBudget bool      `json:"exclude_from_budget"`
	ExcludeFromTotals bool      `json:"exclude_from_totals"`
	UpdatedAt         time.Time `json:"updated_at"`
	CreatedAt         time.Time `json:"created_at"`
	IsGroup           bool      `json:"is_group"`
	GroupID           int64     `json:"group_id"`
}
