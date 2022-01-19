/*
 *  Copyright 2022 Ashok Pon Kumar
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/go-github/v39/github"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var (
	rootCmd       *cobra.Command
	token         string
	org           string
	githubBaseURL string
)

const (
	tokenC = "token"
	orgC   = "org"
)

func init() {
	rootCmd = &cobra.Command{
		Use:   "github-org-stats",
		Short: "Get the repository stats summary of any organization",
		Long:  `github-org-stats uses the github api to get a sense of the stats like stars and fork count of the organisation.`,
		PreRun: func(cmd *cobra.Command, args []string) {
			if token == "" {
				token = os.Getenv("GITHUB_TOKEN")
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			client, err := getClient(ctx, token, githubBaseURL)
			if err != nil {
				return err
			}
			repos, err := repos(ctx, client, org)
			if err != nil {
				logrus.Fatalf("Unable to get repositories : %s", err)
			}
			csvObj := [][]string{
				{"Name", "Stars", "Forks", "IsFork", "URL", "Tags"},
			}
			totalStars := 0
			totalForks := 0
			for _, repo := range repos {
				starCount := repo.GetStargazersCount()
				forkCount := repo.GetForksCount()
				isFork := repo.GetFork()
				csvObj = append(csvObj, []string{repo.GetName(), strconv.Itoa(starCount), strconv.Itoa(forkCount), strconv.FormatBool(isFork), repo.GetURL(), fmt.Sprintf("%+v", repo.Topics)})
				totalStars += starCount
				totalForks += forkCount
			}
			f, err := os.Create("github-org-stats.csv")
			if err != nil {
				logrus.Fatalf("failed to open file : %s", err)
			}
			defer f.Close()
			w := csv.NewWriter(f)
			err = w.WriteAll(csvObj) // calls Flush internally
			if err != nil {
				log.Fatal(err)
			}
			logrus.Infof("Organization : %s", org)
			logrus.Infof("No of stars : %d", totalStars)
			logrus.Infof("No of forks : %d", totalForks)
			return nil
		},
	}

	rootCmd.Flags().StringVarP(&token, tokenC, "t", "", "github personal access token (default $GITHUB_TOKEN)")
	rootCmd.Flags().StringVarP(&org, orgC, "o", "", "github organisation name")

	rootCmd.MarkFlagRequired(orgC)

	rootCmd.Flags().StringVar(&githubBaseURL, "github-base-url", "", "Github base url, if it is not github.com")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Errorf("%s", err)
		os.Exit(1)
	}
}

func getClient(ctx context.Context, token, githubBaseURL string) (*github.Client, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})

	if githubBaseURL == "" {
		return github.NewClient(oauth2.NewClient(ctx, ts)), nil
	}

	return github.NewEnterpriseClient(githubBaseURL, "", oauth2.NewClient(ctx, ts))
}

func repos(ctx context.Context, client *github.Client, org string) ([]*github.Repository, error) {
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, org, opt)
		if err, ok := err.(*github.RateLimitError); ok {
			s := err.Rate.Reset.UTC().Sub(time.Now().UTC())
			if s < 0 {
				s = 5 * time.Second
			}
			time.Sleep(s)
			continue
		}
		if err != nil {
			return allRepos, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}
	return allRepos, nil
}
