package data

import (
	"encoding/json"
	"errors"
	"slices"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/acepanel/panel/v3/internal/app"
	"github.com/acepanel/panel/v3/internal/biz"
	"github.com/acepanel/panel/v3/pkg/api"
	"github.com/acepanel/panel/v3/pkg/apploader"
)

type cacheRepo struct {
	api *api.API
	db  *gorm.DB
}

func NewCacheRepo(db *gorm.DB) biz.CacheRepo {
	return &cacheRepo{
		api: api.NewAPI(app.Version, app.Locale),
		db:  db,
	}
}

func (r *cacheRepo) Get(key biz.CacheKey, defaultValue ...string) (string, error) {
	cache := new(biz.Cache)
	if err := r.db.Where("key = ?", key).First(cache).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", err
		}
	}

	if cache.Value == "" && len(defaultValue) > 0 {
		return defaultValue[0], nil
	}

	return cache.Value, nil
}

func (r *cacheRepo) Set(key biz.CacheKey, value string) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&biz.Cache{Key: key, Value: value}).Error
}

func (r *cacheRepo) UpdateCategories() error {
	categories, err := r.api.Categories()
	if err != nil {
		return err
	}

	encoded, err := json.Marshal(categories)
	if err != nil {
		return err
	}

	return r.Set(biz.CacheKeyCategories, string(encoded))
}

func (r *cacheRepo) UpdateApps() error {
	remote, err := r.api.Apps()
	if err != nil {
		return err
	}

	// 去除本地不存在的应用
	*remote = slices.Clip(slices.DeleteFunc(*remote, func(item *api.App) bool {
		return !slices.Contains(apploader.Slugs(), item.Slug)
	}))

	encoded, err := json.Marshal(remote)
	if err != nil {
		return err
	}

	return r.Set(biz.CacheKeyApps, string(encoded))
}

func (r *cacheRepo) UpdateEnvironments() error {
	environments, err := r.api.Environments()
	if err != nil {
		return err
	}

	hasNet6 := false
	hasNet8 := false
	hasNet9 := false
	for _, env := range *environments {
		if env.Type == "dotnet" {
			switch env.Slug {
			case "6.0":
				hasNet6 = true
			case "8.0":
				hasNet8 = true
			case "9.0":
				hasNet9 = true
			}
		}
	}

	if !hasNet6 {
		*environments = append(*environments, &api.Environment{
			Type:        "dotnet",
			Slug:        "6.0",
			Name:        ".NET 6.0",
			Version:     "6.0.36",
			Description: ".NET 6.0 Runtime",
		})
	}
	if !hasNet8 {
		*environments = append(*environments, &api.Environment{
			Type:        "dotnet",
			Slug:        "8.0",
			Name:        ".NET 8.0",
			Version:     "8.0.12",
			Description: ".NET 8.0 Runtime",
		})
	}
	if !hasNet9 {
		*environments = append(*environments, &api.Environment{
			Type:        "dotnet",
			Slug:        "9.0",
			Name:        ".NET 9.0",
			Version:     "9.0.2",
			Description: ".NET 9.0 Runtime",
		})
	}

	encoded, err := json.Marshal(environments)
	if err != nil {
		return err
	}

	return r.Set(biz.CacheKeyEnvironment, string(encoded))
}

func (r *cacheRepo) UpdateTemplates() error {
	templates, err := r.api.Templates()
	if err != nil {
		return err
	}

	encoded, err := json.Marshal(templates)
	if err != nil {
		return err
	}

	return r.Set(biz.CacheKeyTemplates, string(encoded))
}
