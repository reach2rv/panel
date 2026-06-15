package api

import "fmt"

type Environment struct {
	Type        string `json:"type"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

type Environments []*Environment

// Environments 返回所有环境
func (r *API) Environments() (*Environments, error) {
	resp, err := r.client.R().SetResult(&Response{}).Get("/environments")
	if err != nil {
		return nil, err
	}
	if !resp.IsStatusSuccess() {
		return nil, fmt.Errorf("failed to get environments: %s", resp.String())
	}

	environments, err := getResponseData[Environments](resp)
	if err != nil {
		return nil, err
	}

	if r.locale == "en" {
		for _, env := range *environments {
			// Translate Name
			switch env.Name {
			case "Go 运行环境", "Go":
				env.Name = "Go"
			case "Java 运行环境", "Java":
				env.Name = "Java"
			case "Node.js 运行环境", "Node.js":
				env.Name = "Node.js"
			case "Python 运行环境", "Python":
				env.Name = "Python"
			}

			// Translate Description
			switch env.Description {
			case "Go 运行环境":
				env.Description = "Go runtime environment"
			case "Java 运行环境":
				env.Description = "Java runtime environment"
			case "Node.js 运行环境":
				env.Description = "Node.js runtime environment"
			case "Python 运行环境":
				env.Description = "Python runtime environment"
			}
			if env.Type == "php" {
				env.Name = fmt.Sprintf("PHP-%s", env.Slug)
				env.Description = fmt.Sprintf("PHP-%s runtime environment", env.Slug)
			} else if env.Type == "dotnet" {
				env.Name = fmt.Sprintf(".NET %s", env.Slug)
				env.Description = fmt.Sprintf(".NET %s runtime environment", env.Slug)
			} else if env.Type == "go" {
				env.Name = fmt.Sprintf("Go %s", env.Slug)
				env.Description = fmt.Sprintf("Go %s runtime environment", env.Slug)
			} else if env.Type == "java" {
				env.Name = fmt.Sprintf("Java %s", env.Slug)
				env.Description = fmt.Sprintf("Java %s runtime environment", env.Slug)
			} else if env.Type == "nodejs" {
				env.Name = fmt.Sprintf("Node.js %s", env.Slug)
				env.Description = fmt.Sprintf("Node.js %s runtime environment", env.Slug)
			} else if env.Type == "python" {
				env.Name = fmt.Sprintf("Python %s", env.Slug)
				env.Description = fmt.Sprintf("Python %s runtime environment", env.Slug)
			}
		}
	}

	return environments, nil
}

// EnvironmentCallback 环境下载回调
func (r *API) EnvironmentCallback(typ, slug string) error {
	resp, err := r.client.R().
		SetResult(&Response{}).
		Post(fmt.Sprintf("/environments/%s/%s/callback", typ, slug))
	if err != nil {
		return err
	}
	if !resp.IsStatusSuccess() {
		return fmt.Errorf("failed to callback environment: %s", resp.String())
	}

	return nil
}
