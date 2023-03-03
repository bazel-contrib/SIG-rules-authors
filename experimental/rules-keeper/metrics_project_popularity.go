package main

import (
	"context"
	"fmt"
	"time"

	"github.com/bazel-contrib/SIG-rules-authors/experimental/rules-keeper/timeseries"

	"github.com/google/go-github/v50/github"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/bazel-contrib/SIG-rules-authors/experimental/rules-keeper/proto"
)

// GetProjectPopularity returns the popularity of a project.
func (r *repo) GetProjectPopularity(ctx context.Context) (*pb.ProjectPopularity, error) {
	now := time.Now().UTC()

	return &pb.ProjectPopularity{
		UpdateTime: timestamppb.New(now),
		// TODO(@ashi009): Aggregate ts here.
	}, nil
}

type RepoStatsPoint struct {
	Time       timeseries.DateTime `csv:"time (UTC)"`
	StarsCount int                 `csv:"# stars"`
	ForksCount int                 `csv:"# forks"`
}

func (p *RepoStatsPoint) Timestamp() time.Time {
	return p.Time.Time
}

var RepoStats = &timeseries.Descriptor[*RepoStatsPoint]{
	NewPoint: func(t time.Time) *RepoStatsPoint {
		return &RepoStatsPoint{Time: timeseries.DateTime{Time: t}}
	},
	Align:     timeseries.AlignToSecond,
	Retention: timeSeriesRetention,
}

func (r *repo) loadRepoStatStore() (*timeseries.Store[*RepoStatsPoint], error) {
	return timeseries.Load(RepoStats, fmt.Sprintf("store/%s/%s/metrics/repo_stats", r.owner, r.name))
}

func (r *repo) UpdateRepoStats(ctx context.Context) error {
	s, err := r.loadRepoStatStore()
	if err != nil {
		return err
	}

	now := time.Now()
	repo, _, err := r.cli.Repositories.Get(ctx, r.owner, r.name)
	if err != nil {
		return err
	}
	p := s.GetOrCreatePointAt(now)
	p.StarsCount = repo.GetStargazersCount()
	p.ForksCount = repo.GetForksCount()

	s.ShiftWindow(now)
	return s.Flush()
}

type TrafficPoint struct {
	Date              timeseries.Date `csv:"date (UTC)"`
	ViewsCount        int             `csv:"# views"`
	ViewsUniqueCount  int             `csv:"# unique views"`
	ClonesCount       int             `csv:"# clones"`
	ClonesUniqueCount int             `csv:"# unique clones"`
}

func (p *TrafficPoint) Timestamp() time.Time {
	return p.Date.Time
}

var Traffic = &timeseries.Descriptor[*TrafficPoint]{
	NewPoint: func(t time.Time) *TrafficPoint {
		return &TrafficPoint{Date: timeseries.Date{Time: t}}
	},
	Align:     timeseries.AlignToISO8601WeekStartDayInUTC,
	Retention: timeSeriesRetention,
}

func (r *repo) loadTrafficStore() (*timeseries.Store[*TrafficPoint], error) {
	return timeseries.Load(Traffic, fmt.Sprintf("store/%s/%s/metrics/traffic", r.owner, r.name))
}

func (r *repo) UpdateTraffic(ctx context.Context) error {
	s, err := r.loadTrafficStore()
	if err != nil {
		return err
	}

	_, lastUpdate := s.Window()
	now := time.Now()

	vs, _, err := r.cli.Repositories.ListTrafficViews(ctx, r.owner, r.name, &github.TrafficBreakdownOptions{
		Per: "week",
	})
	if err != nil {
		return err
	}
	for _, v := range vs.Views {
		t := v.GetTimestamp().Time
		if t.Before(lastUpdate) || !t.Before(now) {
			continue
		}
		p := s.GetOrCreatePointAt(t)
		p.ViewsCount = v.GetCount()
		p.ViewsUniqueCount = v.GetUniques()
	}

	cls, _, err := r.cli.Repositories.ListTrafficClones(ctx, r.owner, r.name, &github.TrafficBreakdownOptions{
		Per: "week",
	})
	if err != nil {
		return err
	}
	for _, c := range cls.Clones {
		t := c.GetTimestamp().Time
		if t.Before(lastUpdate) || !t.Before(now) {
			continue
		}
		p := s.GetOrCreatePointAt(t)
		p.ClonesCount = c.GetCount()
		p.ClonesUniqueCount = c.GetUniques()
	}

	s.ShiftWindow(now)
	return s.Flush()
}
