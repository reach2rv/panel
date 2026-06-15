package request

import (
	"net/http"
)

type Paginate struct {
	Page  uint `json:"page" form:"page" query:"page" validate:"required|min:1"`
	Limit uint `json:"limit" form:"limit" query:"limit" validate:"required|min:1|max:10000"`
}

func (r *Paginate) Messages(_ *http.Request) map[string]string {
	return map[string]string{
		"Page.gte":       "Page must be greater than or equal to 1",
		"Limit.gte":      "Limit must be greater than or equal to 1",
		"Limit.lte":      "Limit must be less than or equal to 10000",
		"Page.number":    "Page must be a number",
		"Limit.number":   "Limit must be a number",
		"Page.required":  "Page is required",
		"Limit.required": "Limit is required",
	}
}

func (r *Paginate) Prepare(_ *http.Request) error {
	if r.Page == 0 {
		r.Page = 1
	}
	if r.Limit == 0 {
		r.Limit = 10
	}
	return nil
}
