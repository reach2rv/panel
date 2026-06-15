package api

import (
	"fmt"
	"time"
)

type App struct {
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Slug        string    `json:"slug"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Categories  []string  `json:"categories"`
	Depends     string    `json:"depends"` // 依赖表达式
	Channels    []struct {
		Slug      string `json:"slug"`      // 渠道代号
		Name      string `json:"name"`      // 渠道名称
		Panel     string `json:"panel"`     // 最低支持面板版本
		Install   string `json:"install"`   // 安装脚本
		Uninstall string `json:"uninstall"` // 卸载脚本
		Update    string `json:"update"`    // 更新脚本
		Version   string `json:"version"`   // 版本号
		Log       string `json:"log"`       // 更新日志
	} `json:"channels"`
	Order int `json:"order"`
}

type Apps []*App

// Apps 返回所有应用
func (r *API) Apps() (*Apps, error) {
	resp, err := r.client.R().SetResult(&Response{}).Get("/apps")
	if err != nil {
		return nil, err
	}
	if !resp.IsStatusSuccess() {
		return nil, fmt.Errorf("failed to get apps: %s", resp.String())
	}

	apps, err := getResponseData[Apps](resp)
	if err != nil {
		return nil, err
	}

	if r.locale == "en" {
		for _, app := range *apps {
			if t, ok := AppTranslations[app.Slug]; ok {
				app.Name = t.Name
				app.Description = t.Description
			}
		}
	}

	return apps, nil
}

// AppBySlug 根据slug返回应用
func (r *API) AppBySlug(slug string) (*App, error) {
	resp, err := r.client.R().SetResult(&Response{}).Get(fmt.Sprintf("/apps/%s", slug))
	if err != nil {
		return nil, err
	}
	if !resp.IsStatusSuccess() {
		return nil, fmt.Errorf("failed to get app: %s", resp.String())
	}

	app, err := getResponseData[App](resp)
	if err != nil {
		return nil, err
	}

	if r.locale == "en" {
		if t, ok := AppTranslations[app.Slug]; ok {
			app.Name = t.Name
			app.Description = t.Description
		}
	}

	return app, nil
}

// AppCallback 应用下载回调
func (r *API) AppCallback(slug string) error {
	resp, err := r.client.R().
		SetResult(&Response{}).
		Post(fmt.Sprintf("/apps/%s/callback", slug))
	if err != nil {
		return err
	}
	if !resp.IsStatusSuccess() {
		return fmt.Errorf("failed to callback app: %s", resp.String())
	}

	return nil
}

var AppTranslations = map[string]struct {
	Name        string
	Description string
}{
	"nginx":         {Name: "Nginx", Description: "High-performance HTTP and reverse proxy web server"},
	"apache":        {Name: "Apache", Description: "Highly secure, efficient, and extensible web server"},
	"openresty":     {Name: "OpenResty", Description: "Web platform that integrates Nginx and LuaJIT"},
	"mysql":         {Name: "MySQL", Description: "Relational database management system"},
	"mariadb":       {Name: "MariaDB", Description: "Community-developed, commercially supported fork of MySQL"},
	"postgresql":    {Name: "PostgreSQL", Description: "Powerful, open-source object-relational database system"},
	"mongodb":       {Name: "MongoDB", Description: "Document-based, distributed database"},
	"redis":         {Name: "Redis", Description: "In-memory data structure store used as database, cache, and broker"},
	"valkey":        {Name: "Valkey", Description: "High-performance key-value store, a fork of Redis"},
	"memcached":     {Name: "Memcached", Description: "High-performance, distributed memory object caching system"},
	"pureftpd":      {Name: "Pure-FTPd", Description: "Free, secure, and production-quality FTP server"},
	"phpmyadmin":    {Name: "phpMyAdmin", Description: "Free software tool written in PHP for MySQL administration over the web"},
	"docker":        {Name: "Docker", Description: "Platform designed to help developers build, share, and run applications"},
	"codeserver":    {Name: "Code Server", Description: "Run VS Code on any machine and access it in the browser"},
	"gitea":         {Name: "Gitea", Description: "Painless self-hosted Git service"},
	"supervisor":    {Name: "Supervisor", Description: "Process manager to monitor and control processes on UNIX-like OS"},
	"frp":           {Name: "FRP", Description: "Fast reverse proxy to expose a local server behind a NAT to the internet"},
	"fail2ban":      {Name: "Fail2ban", Description: "Intrusion prevention software protecting servers from brute-force attacks"},
	"minio":         {Name: "MinIO", Description: "High-performance, Kubernetes-native object storage"},
	"prometheus":    {Name: "Prometheus", Description: "Open-source systems monitoring and alerting toolkit"},
	"grafana":       {Name: "Grafana", Description: "Multi-platform analytics and interactive visualization web application"},
	"clickhouse":    {Name: "ClickHouse", Description: "High-performance, column-oriented SQL database system for OLAP"},
	"elasticsearch": {Name: "Elasticsearch", Description: "Distributed, RESTful search and analytics engine"},
	"opensearch":    {Name: "OpenSearch", Description: "Flexible, scalable open-source search and analytics suite"},
	"kafka":         {Name: "Kafka", Description: "Distributed event store and stream-processing platform"},
	"rocketmq":      {Name: "RocketMQ", Description: "Distributed messaging and streaming platform"},
	"rsync":         {Name: "Rsync", Description: "Fast and extraordinarily versatile file-copying tool"},
	"s3fs":          {Name: "s3fs", Description: "FUSE-based file system backed by Amazon S3"},
	"percona":       {Name: "Percona", Description: "Production-ready, enterprise-grade open-source database server for MySQL"},
}
