package main

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/bazel-contrib/SIG-rules-authors/experimental/rules-keeper/timeseries"
)

type CommitActivityPoint struct {
	Date       timeseries.Date `csv:"date (UTC)"`
	Count      int             `csv:"# commits"`
	OwnerCount int             `csv:"# commits by owner"`
}

func (p *CommitActivityPoint) Timestamp() time.Time {
	return p.Date.Time
}

var CommitActivity = &timeseries.Descriptor[*CommitActivityPoint]{
	NewPoint: func(t time.Time) *CommitActivityPoint {
		return &CommitActivityPoint{Date: timeseries.Date{Time: t}}
	},
	Align: timeseries.AlignToISO8601WeekStartDayInUTC,
}

func (r *repo) loadCommitActivityStore() (*timeseries.Store[*CommitActivityPoint], error) {
	return timeseries.Load(CommitActivity, fmt.Sprintf("store/%s/%s/metrics/commit_activity", r.owner, r.name))
}

func (r *repo) UpdateCommitActivity(ctx context.Context) error {
	s, err := r.loadCommitActivityStore()
	if err != nil {
		return err
	}

	now := time.Now()
	ps, _, err := r.cli.Repositories.ListParticipation(ctx, r.owner, r.name)
	if err != nil {
		return err
	}
	// ListParticipation returns the last 52 weeks of data, but it doesn't specify
	// when a week starts. We are using ISO 8601
	for i := range ps.All {
		t := now.AddDate(0, 0, -i*7)
		p := s.GetOrCreatePointAt(t)
		p.Count = ps.All[i]
		p.OwnerCount = ps.Owner[i]
	}

	s.ShiftWindow(now)
	return s.Flush()
}

type IssueActivityPoint struct {
	Date            timeseries.Date `csv:"date (UTC)"`
	OpenPRCount     int             `csv:"# open pr"`
	ClosePRCount    int             `csv:"# close pr"`
	OpenIssueCount  int             `csv:"# open issue"`
	CloseIssueCount int             `csv:"# close issue"`
}

func (p *IssueActivityPoint) Timestamp() time.Time {
	return p.Date.Time
}

var IssueActivity = &timeseries.Descriptor[*IssueActivityPoint]{
	NewPoint: func(t time.Time) *IssueActivityPoint {
		return &IssueActivityPoint{Date: timeseries.Date{Time: t}}
	},
	Align: timeseries.AlignToDayInUTC,
}

func (r *repo) loadIssueActivtyStore() (*timeseries.Store[*IssueActivityPoint], error) {
	return timeseries.Load(IssueActivity, fmt.Sprintf("store/%s/%s/metrics/issue_activity", r.owner, r.name))
}

func (r *repo) UpdateIssueActivity(ctx context.Context) error {
	s, err := r.loadIssueActivtyStore()
	if err != nil {
		return err
	}

	_, lastUpdate := s.Window()
	now := time.Now()

	// We want to add data from [lastUpdate, now), all new points need to check
	// against those boundaries.

	// List all repo events to check on close and reopen events. Interestingly
	// enough, ListIssueEvents doesn't include create events. We need to list
	// issues to find create events instead.
	for it := r.ListIssueEvents(ctx); ; {
		evt, err := it.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		ct := evt.GetCreatedAt().Time
		if !ct.Before(now) {
			continue
		}
		if ct.Before(lastUpdate) {
			break
		}
		switch evt.GetEvent() {
		case "closed":
			p := s.GetOrCreatePointAt(ct)
			if evt.GetIssue().PullRequestLinks == nil {
				p.CloseIssueCount++
			} else {
				p.ClosePRCount++
			}
		case "reopened":
			p := s.GetOrCreatePointAt(ct)
			if evt.GetIssue().PullRequestLinks == nil {
				p.OpenIssueCount++
			} else {
				p.OpenPRCount++
			}
		}
	}

	for it := r.ListIssues(ctx, lastUpdate); ; {
		iss, err := it.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		ct := iss.GetCreatedAt().Time
		// Note that, list issues returns all issues that has been updated since
		// lastUpdate, so It's possible that the create time are out of range. As we
		// only care about create events, we need to filter out those issues.
		if ct.Before(lastUpdate) || !ct.Before(now) {
			continue
		}
		p := s.GetOrCreatePointAt(ct)
		if iss.PullRequestLinks == nil {
			p.OpenIssueCount++
		} else {
			p.OpenPRCount++
		}
	}

	s.ShiftWindow(now)
	return s.Flush()
}
