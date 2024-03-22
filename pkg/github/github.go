package github

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Client struct {
	*github.Client
}

func NewClient(token string) (client *Client) {

	if token != "" {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)
		return &Client{github.NewClient(tc)}
	} else {
		return &Client{github.NewClient(&http.Client{})}
	}
}

func (c Client) ListReleaseAssets(owner string, repo string, id int64) (allAssets []*github.ReleaseAsset) {

	// option for ListReleaseAssets
	opt := &github.ListOptions{
		Page:    0,
		PerPage: 100, // max assets by page
	}

	for {
		assets, resp, err := c.Repositories.ListReleaseAssets(context.Background(), owner, repo, id, opt)
		if err != nil {
			fmt.Printf("Erreur lors de la récupération des actifs de la version: %v\n", err)
		}

		allAssets = append(allAssets, assets...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage

	}
	return allAssets
}

func FilterAssets(assets []*github.ReleaseAsset) []*github.ReleaseAsset {
	assets = filterAssetsByArch(assets)
	return filterAssetsByOs(assets)
}

func filterAssetsByArch(assets []*github.ReleaseAsset) []*github.ReleaseAsset {
	return filterAssetsBy(assets, runtime.GOARCH)
}

func filterAssetsByOs(assets []*github.ReleaseAsset) []*github.ReleaseAsset {
	return filterAssetsBy(assets, runtime.GOOS)
}

// FilterAssets filters assets based on the provided name.
func filterAssetsBy(assets []*github.ReleaseAsset, filter string) []*github.ReleaseAsset {
	var filteredAssets []*github.ReleaseAsset
	for _, asset := range assets {
		if strings.Contains(*asset.Name, filter) {
			filteredAssets = append(filteredAssets, asset)
		}
	}
	return filteredAssets
}

// convertList filters assets based on the provided name.
func ConvertList(assets []*github.ReleaseAsset) (listAssets []string) {
	for _, asset := range assets {
		listAssets = append(listAssets, *asset.Name)
	}
	return
}

func (c Client) ListReleases(owner string, repo string) (allReleases []*github.RepositoryRelease) {

	// option for ListReleaseAssets
	opt := &github.ListOptions{
		Page:    0,
		PerPage: 100, // max assets by page
	}

	for {
		releases, resp, err := c.Repositories.ListReleases(context.Background(), owner, repo, opt)
		if err != nil {
			fmt.Printf("error to list all version of repo: %v\n", err)
		}

		allReleases = append(allReleases, releases...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage

	}
	return allReleases
}

// GetRelease version bases on version.
func GetRelease(assets []*github.RepositoryRelease, version string) (*github.RepositoryRelease, error) {
	for _, asset := range assets {
		if strings.Contains(*asset.TagName, version) {
			return asset, nil
		}
	}
	return nil, errors.New("no release found")
}

func DownloadFile(url string, filepath string) error {
	// Fetch the content from the URL
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Create a new file on the file system
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy the content from the HTTP response into the local file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	fmt.Printf("file downloaded: %s\n", filepath)
	return nil
}

func CheckSuffixArchive(assetFileName string) (bool, string) {
	if strings.HasSuffix(assetFileName, ".tar.gz") || strings.HasSuffix(assetFileName, ".tgz") {
		return true, "tar"
	}
	if strings.HasSuffix(assetFileName, ".zip") {
		return true, "zip"
	}
	return false, ""
}
