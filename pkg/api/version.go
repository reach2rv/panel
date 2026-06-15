package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/libtnb/utils/env"
)

type VersionDownload struct {
	URL      string `json:"url"`
	Arch     string `json:"arch"`
	Checksum string `json:"checksum"`
}

type Version struct {
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Type        string            `json:"type"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Downloads   []VersionDownload `json:"downloads"`
}

type Versions []Version

type GitHubReleaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type GitHubRelease struct {
	TagName   string               `json:"tag_name"`
	Body      string               `json:"body"`
	Assets    []GitHubReleaseAsset `json:"assets"`
	CreatedAt time.Time            `json:"created_at"`
}

// LatestVersion 返回最新版本
func (r *API) LatestVersion(channel string) (*Version, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/reach2rv/panel/releases/latest")
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch latest release from GitHub: status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	arch := "amd64"
	if env.IsArm() {
		arch = "arm64"
	}

	var downloadURL string
	targetAsset := fmt.Sprintf("ornaverse-panel_linux_%s.zip", arch)
	for _, asset := range release.Assets {
		if asset.Name == targetAsset {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}
	if downloadURL == "" {
		for _, asset := range release.Assets {
			if strings.Contains(asset.Name, "linux") && strings.Contains(asset.Name, arch) && strings.HasSuffix(asset.Name, ".zip") {
				downloadURL = asset.BrowserDownloadURL
				break
			}
		}
	}

	if downloadURL == "" {
		return nil, fmt.Errorf("failed to find release asset for %s", targetAsset)
	}

	var checksumURL string
	for _, asset := range release.Assets {
		if strings.Contains(asset.Name, "checksums.txt") {
			checksumURL = asset.BrowserDownloadURL
			break
		}
	}
	if checksumURL == "" {
		checksumURL = downloadURL + ".sha256"
	}

	version := &Version{
		Version:     release.TagName,
		Description: release.Body,
		Downloads: []VersionDownload{
			{
				URL:      downloadURL,
				Arch:     arch,
				Checksum: checksumURL,
			},
		},
	}

	return version, nil
}

// IntermediateVersions 返回当前版本之后的所有版本
func (r *API) IntermediateVersions(channel string) (*Versions, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/reach2rv/panel/releases")
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch releases from GitHub: status %d", resp.StatusCode)
	}

	var releases []GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, err
	}

	var list Versions
	for _, release := range releases {
		if release.TagName == r.panelVersion {
			break
		}
		list = append(list, Version{
			CreatedAt:   release.CreatedAt,
			Version:     release.TagName,
			Description: release.Body,
		})
	}

	return &list, nil
}
