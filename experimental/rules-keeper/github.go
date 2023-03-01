package main

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v50/github"
)

type fetchFunc[T any] func(*github.ListOptions) ([]T, *github.Response, error)

type iterator[T any] struct {
	fetch    fetchFunc[T]
	nextPage int
	buf      []T
}

func newIterator[T any](f fetchFunc[T]) *iterator[T] {
	return &iterator[T]{
		fetch:    f,
		nextPage: 1,
	}
}

func (it *iterator[T]) Next() (res T, _ error) {
	if len(it.buf) == 0 && it.nextPage > 0 {
		buf, resp, err := it.fetch(&github.ListOptions{
			Page: it.nextPage,
		})
		if err != nil {
			return res, err
		}
		it.nextPage = resp.NextPage
		it.buf = buf
	}
	if len(it.buf) == 0 {
		return res, io.EOF
	}
	res, it.buf = it.buf[0], it.buf[1:]
	return res, nil
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

func (a *app) ListInstallations(ctx context.Context) *iterator[*installation] {
	return newIterator(func(opts *github.ListOptions) ([]*installation, *github.Response, error) {
		ins, resp, err := a.cli.Apps.ListInstallations(ctx, opts)
		if err != nil {
			return nil, nil, err
		}
		res := make([]*installation, 0, len(ins))
		for _, i := range ins {
			it, err := newInstallation(i.GetAppID(), i.GetID(), privateKey)
			if err != nil {
				return nil, nil, err
			}
			res = append(res, it)
		}
		return res, resp, nil
	})
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

func (n *installation) ListRepos(ctx context.Context) *iterator[*repo] {
	return newIterator(func(opts *github.ListOptions) ([]*repo, *github.Response, error) {
		r, resp, err := n.cli.Apps.ListRepos(ctx, opts)
		if err != nil {
			return nil, nil, err
		}
		res := make([]*repo, 0, len(r.Repositories))
		for _, r := range r.Repositories {
			res = append(res, &repo{
				cli:   n.cli,
				owner: r.GetOwner().GetLogin(),
				name:  r.GetName(),
			})
		}
		return res, resp, nil
	})
}

type repo struct {
	cli   *github.Client
	owner string
	name  string
}

func (r *repo) ListIssues(ctx context.Context, since time.Time) *iterator[*github.Issue] {
	return newIterator(func(opts *github.ListOptions) ([]*github.Issue, *github.Response, error) {
		return r.cli.Issues.ListByRepo(ctx, r.owner, r.name, &github.IssueListByRepoOptions{
			State:       "all",
			Since:       since,
			ListOptions: *opts,
		})
	})
}

// ListIssueEvents returns all issue events for the repository. The iterator
// traverses events in reverse chronological order.
func (r *repo) ListIssueEvents(ctx context.Context) *iterator[*github.IssueEvent] {
	return newIterator(func(opts *github.ListOptions) ([]*github.IssueEvent, *github.Response, error) {
		return r.cli.Issues.ListRepositoryEvents(ctx, r.owner, r.name, opts)
	})
}

func (r *repo) GetLicense(ctx context.Context) (string, error) {
	l, _, err := r.cli.Repositories.License(ctx, r.owner, r.name)
	if err != nil {
		return "", err
	}
	return l.GetLicense().GetName(), nil
}

func (r *repo) GetReadme(ctx context.Context, ref string) (string, error) {
	c, _, err := r.cli.Repositories.GetReadme(ctx, r.owner, r.name, &github.RepositoryContentGetOptions{
		Ref: ref,
	})
	if err != nil {
		return "", err
	}
	return c.GetContent()

}

func (r *repo) ListTags(ctx context.Context) *iterator[*github.RepositoryTag] {
	return newIterator(func(opts *github.ListOptions) ([]*github.RepositoryTag, *github.Response, error) {
		return r.cli.Repositories.ListTags(ctx, r.owner, r.name, opts)
	})
}

func (r *repo) GetParticipations(ctx context.Context) (*github.RepositoryParticipation, error) {
	pns, _, err := r.cli.Repositories.ListParticipation(ctx, r.owner, r.name)
	return pns, err
}

// Note: it will return 404 for forked repos.
func (r *repo) GetCommunityHealthMetrics(ctx context.Context) (*github.CommunityHealthMetrics, error) {
	chm, _, err := r.cli.Repositories.GetCommunityHealthMetrics(ctx, r.owner, r.name)
	return chm, err
}

// ListReleases returns the releases for a repo.
func (r *repo) ListReleases(ctx context.Context) *iterator[*github.RepositoryRelease] {
	return newIterator(func(opts *github.ListOptions) ([]*github.RepositoryRelease, *github.Response, error) {
		return r.cli.Repositories.ListReleases(ctx, r.owner, r.name, opts)
	})
}
