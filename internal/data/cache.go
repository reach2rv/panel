package data

import (
	"encoding/json"
	"errors"
	"io"
	"slices"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"resty.dev/v3"

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

	dotnetVersions := map[string]string{
		"6.0": "6.0.428",
		"8.0": "8.0.112",
		"9.0": "9.0.102",
	}

	type DotNetReleasesIndex struct {
		ReleasesIndex []struct {
			ChannelVersion string `json:"channel-version"`
			LatestSdk      string `json:"latest-sdk"`
		} `json:"releases-index"`
	}

	client := resty.New()
	client.SetTimeout(5 * time.Second)
	resp, err := client.R().Get("https://dotnetcli.blob.core.windows.net/dotnet/release-metadata/releases-index.json")
	if err == nil && resp.IsStatusSuccess() {
		defer func() { _ = resp.Body.Close() }()
		bodyBytes, err := io.ReadAll(resp.Body)
		if err == nil {
			var index DotNetReleasesIndex
			if err := json.Unmarshal(bodyBytes, &index); err == nil {
				for _, release := range index.ReleasesIndex {
					if release.ChannelVersion != "" && release.LatestSdk != "" {
						dotnetVersions[release.ChannelVersion] = release.LatestSdk
					}
				}
			}
		}
	}

	hasNet := make(map[string]bool)
	for _, env := range *environments {
		if env.Type == "dotnet" {
			hasNet[env.Slug] = true
		}
	}

	for channel, version := range dotnetVersions {
		if channel < "6.0" {
			continue
		}
		if !hasNet[channel] {
			*environments = append(*environments, &api.Environment{
				Type:        "dotnet",
				Slug:        channel,
				Name:        ".NET " + channel,
				Version:     version,
				Description: ".NET " + channel + " Runtime",
			})
		} else {
			for _, env := range *environments {
				if env.Type == "dotnet" && env.Slug == channel {
					env.Version = version
				}
			}
		}
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
