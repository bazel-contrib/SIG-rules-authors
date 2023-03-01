package main

import (
	"context"
	"flag"
	"io"
	"strings"

	"github.com/golang/glog"
	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

var (
	appID         int64 = 290065
	privateKey          = "appventure-test.2023-02-06.private-key.pem"
	personalToken       = ""
	repoOwner           = ""
	repoName            = ""
)

func init() {
	flag.Int64Var(&appID, "app_id", appID, "The Github App ID.")
	flag.StringVar(&privateKey, "private_key", privateKey, "The path to the Github App private key file.")
	flag.StringVar(&personalToken, "personal_token", personalToken, "The personal token to use for the Github API. When set, instead of using App credential to fetch all installed repos, you must specify the owner and repo to update metrics.")
	flag.StringVar(&repoOwner, "owner", repoOwner, "The owner of the repo to update metrics. Must be specified when using personal token.")
	flag.StringVar(&repoName, "repo", repoName, "The name of the repo to update metrics. Must be specified when using personal token.")
}

var releaseContentType = map[string]bool{
	"application/x-gzip": true,
	"application/gzip":   true,
	"application/zip":    true,
}

func parse(repoPath string) (owner, repo string) {
	idx := strings.LastIndexByte(repoPath, '/')
	return repoPath[:idx], repoPath[idx+1:]
}

func main() {
	flag.Parse()

	ctx := context.Background()
	ctx = context.WithValue(ctx, oauth2.HTTPClient, newLoggedHTTPClient())

	if personalToken != "" {
		if repoOwner == "" || repoName == "" {
			glog.Exitf("owner and repo must be specified when using personal token")
		}
		r := &repo{
			cli:   github.NewTokenClient(ctx, personalToken),
			owner: repoOwner,
			name:  repoName,
		}
		if err := r.Update(ctx); err != nil {
			glog.Exit(err)
		}
		if err := r.UpdateMetrics(ctx); err != nil {
			glog.Exit(err)
		}
		return
	}

	app, err := newApp(appID, privateKey)
	if err != nil {
		glog.Exit(err)
	}
	// TODO(@ashi009): flatten the iterator.
	for it := app.ListInstallations(ctx); ; {
		inst, err := it.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			glog.Exit(err)
		}
		for rit := inst.ListRepos(ctx); ; {
			r, err := rit.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				glog.Exit(err)
			}
			if err := r.UpdateMetrics(ctx); err != nil {
				glog.Exit(err)
			}
			if err := r.UpdateVersions(ctx); err != nil {
				glog.Exit(err)
			}
		}
	}
}

func (r *repo) Update(ctx context.Context) error {
	if err := r.UpdateVersions(ctx); err != nil {
		return err
	}
	if err := r.UpdateMetrics(ctx); err != nil {
		return err
	}
	return nil
}

func (r *repo) UpdateMetrics(ctx context.Context) error {
	glog.Infof("Updating metrics for %s/%s", r.owner, r.name)
	if err := r.UpdateRepoStats(ctx); err != nil {
		return err
	}
	if err := r.UpdateCommitActivity(ctx); err != nil {
		return err
	}
	if err := r.UpdateIssueActivity(ctx); err != nil {
		return err
	}
	if err := r.UpdateComunityHealth(ctx); err != nil {
		return err
	}
	if err := r.UpdateTraffic(ctx); err != nil {
		return err
	}
	return nil
}
