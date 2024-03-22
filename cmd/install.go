package cmd

import (
	"fmt"
	"strings"

	"github.com/jplanckeel/kit/pkg/config"
	"github.com/jplanckeel/kit/pkg/github"
	"github.com/spf13/cobra"
)

var install = &cobra.Command{
	Use:   "install",
	Short: "install tools",
	Run: func(cmd *cobra.Command, args []string) {

		config, _ := config.NewSource("example.yaml")
		for _, tools := range config {
			for tool, versions := range tools {
				repo := strings.Split(tool, "/")
				fmt.Printf("repo=%s project=%s versions=%s\n", repo[0], repo[1], versions)

				// Get all releases for a tool
				gh := github.NewClient(flags.GithubToken)
				releases := gh.ListReleases(repo[0], repo[1])

				for _, version := range versions {
					release, err := github.GetRelease(releases, version)
					if err != nil {
						fmt.Println(err)
					}

					//tmp for dev
					if release != nil {
						assets := gh.ListReleaseAssets(repo[0], repo[1], *release.ID)
						listAssets := github.ConvertList(assets)

						// Filter assets with OS and Arch
						filteredAssets := github.FilterAssets(listAssets)

						for _, asset := range filteredAssets {
							fmt.Println(asset)
						}
					}
				}

			}
		}
		/*
			releases := github.ListReleases("terragrunt", "gruntwork-io")
			release, err := github.GetRelease(releases, "0.55.18")
			if err != nil {
				fmt.Println(err)
			}

			//tmp for dev
			if release != nil {
				assets := github.ListReleaseAssets("terragrunt", "gruntwork-io", *release.ID)
				listAssets := github.ConvertList(assets)

				// Filter assets with OS and Arch
				filteredAssets := github.FilterAssets(listAssets)

				for _, asset := range filteredAssets {
					fmt.Println(asset)
				}
			}*/

	},
}

func init() {
	rootCmd.AddCommand(install)
}
