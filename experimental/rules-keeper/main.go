package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/golang/glog"
	"github.com/google/go-github/v50/github"
)

var (
	appID      int64 = 290065
	privateKey       = "appventure-test.2023-02-06.private-key.pem"
)

func init() {
	flag.Int64Var(&appID, "app_id", appID, "The Github App ID.")
	flag.StringVar(&privateKey, "private_key", privateKey, "The path to the Github App private key file.")
}

const maxPageSize = 100 // max allowed page size.

func main() {
	flag.Parse()
	ctx := context.Background()
	app, err := newApp(appID, privateKey)
	if err != nil {
		glog.Exit(err)
	}
	insts, err := app.listInstallations(ctx)
	if err != nil {
		glog.Exit(err)
	}
	// TODO: break this into smaller pieces.
	for _, inst := range insts {
		fmt.Println(github.Stringify(inst))
		inst, err := newInstallation(inst.GetAppID(), inst.GetID(), privateKey)
		if err != nil {
			glog.Exit(err)
		}
		repos, err := inst.ListRepos(ctx)
		if err != nil {
			glog.Exit(err)
		}
		for _, repo := range repos {
			fmt.Println(github.Stringify(repo))
			v, err := inst.GetTrafficViews(ctx, repo)
			if err != nil {
				glog.Exit(err)
			}
			fmt.Println(github.Stringify(v))
		}
	}
}

type app struct {
	cli *github.Client
}

func newApp(appID int64, privateKey string) (*app, error) {
	// TODO: find a better way for handling the token source.
	itr, err := ghinstallation.NewAppsTransportKeyFromFile(http.DefaultTransport, appID, privateKey)
	if err != nil {
		return nil, err
	}
	return &app{github.NewClient(&http.Client{Transport: itr})}, nil
}

// TODO(@ashi009): write an itereator for listing APIS.
// Github is not using token based pagination, so the underlying data might
// change while loading the next page. So we need a way to remedy this, an
// iterator to dedup might be the solution.

func (a *app) listInstallations(ctx context.Context) ([]*github.Installation, error) {
	var res []*github.Installation
	for pg := 1; ; pg++ {
		ins, _, err := a.cli.Apps.ListInstallations(ctx, &github.ListOptions{
			Page:    pg,
			PerPage: maxPageSize,
		})
		if err != nil {
			return nil, err
		}
		res = append(res, ins...)
		if len(ins) < maxPageSize {
			break
		}
	}
	return res, nil
}

type installation struct {
	cli *github.Client
}

func newInstallation(appID, installID int64, privateKey string) (*installation, error) {
	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, appID, installID, privateKey)
	if err != nil {
		return nil, err
	}
	return &installation{github.NewClient(&http.Client{Transport: itr})}, nil
}

func (i *installation) ListRepos(ctx context.Context) ([]*github.Repository, error) {
	var res []*github.Repository
	for pg := 1; ; pg++ {
		repos, _, err := i.cli.Apps.ListRepos(ctx, &github.ListOptions{
			Page:    pg,
			PerPage: maxPageSize,
		})
		if err != nil {
			return nil, err
		}
		res = append(res, repos.Repositories...)
		if len(repos.Repositories) < maxPageSize || len(res) >= repos.GetTotalCount() {
			break
		}
	}
	return res, nil
}

func (i *installation) GetTrafficViews(ctx context.Context, repo *github.Repository) (*github.TrafficViews, error) {
	tvs, _, err := i.cli.Repositories.ListTrafficViews(ctx, repo.GetOwner().GetLogin(), repo.GetName(), &github.TrafficBreakdownOptions{
		Per: "day",
	})
	return tvs, err
}

func (i *installation) GetTrafficClones(ctx context.Context, repo *github.Repository) (*github.TrafficClones, error) {
	cls, _, err := i.cli.Repositories.ListTrafficClones(ctx, repo.GetOwner().GetLogin(), repo.GetName(), &github.TrafficBreakdownOptions{
		Per: "day",
	})
	return cls, err
}
