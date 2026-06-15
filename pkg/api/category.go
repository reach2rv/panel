package api

import "fmt"

type Category struct {
	Slug  string `json:"slug"`
	Name  string `json:"name"`
	Order int    `json:"order"`
}

type Categories []*Category

// Categories 返回所有分类
func (r *API) Categories() (*Categories, error) {
	resp, err := r.client.R().SetResult(&Response{}).Get("/categories")
	if err != nil {
		return nil, err
	}
	if !resp.IsStatusSuccess() {
		return nil, fmt.Errorf("failed to get categories: %s", resp.String())
	}

	categories, err := getResponseData[Categories](resp)
	if err != nil {
		return nil, err
	}

	if r.locale == "en" {
		for _, cat := range *categories {
			if t, ok := CategoryTranslations[cat.Slug]; ok {
				cat.Name = t
			}
		}
	}

	return categories, nil
}

var CategoryTranslations = map[string]string{
	"website":    "Website",
	"database":   "Database",
	"tool":       "Tool",
	"runtime":    "Runtime",
	"app":        "Application",
	"monitor":    "Monitor",
	"deployment": "Deployment",
	"store":      "Store",
}
