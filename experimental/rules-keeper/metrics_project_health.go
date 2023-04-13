package main

import (
	"context"
	"fmt"
	"time"

	"github.com/bazel-contrib/SIG-rules-authors/experimental/rules-keeper/timeseries"

	pb "github.com/bazel-contrib/SIG-rules-authors/experimental/rules-keeper/proto"
)

// GetProjectHealth returns the health of a project.
func (r *repo) GetProjectHealth(ctx context.Context) (*pb.ProjectHealth, error) {
	return &pb.ProjectHealth{
		// TODO(@ashi009): Aggregate ts here.
	}, nil
}

type CommunityHealthPoint struct {
	Time                      timeseries.DateTime `csv:"time (UTC)"`
	CommunityHealthPercentage int                 `csv:"community health percentage"`
}

func (p *CommunityHealthPoint) Timestamp() time.Time {
	return p.Time.Time
}

var CommunityHealth = &timeseries.Descriptor[*CommunityHealthPoint]{
	NewPoint: func(t time.Time) *CommunityHealthPoint {
		return &CommunityHealthPoint{Time: timeseries.DateTime{Time: t}}
	},
	Align:     timeseries.AlignToSecond,
	Retention: timeSeriesRetention,
}

func (r *repo) loadCommunityHealthStore() (*timeseries.Store[*CommunityHealthPoint], error) {
	return timeseries.Load(CommunityHealth, fmt.Sprintf("store/%s/%s/metrics/community_health", r.owner, r.name))
}

func (r *repo) UpdateComunityHealth(ctx context.Context) error {
	s, err := r.loadCommunityHealthStore()
	if err != nil {
		return err
	}

	now := time.Now()
	m, _, err := r.cli.Repositories.GetCommunityHealthMetrics(ctx, r.owner, r.name)
	if err != nil {
		return err
	}
	p := s.GetOrCreatePointAt(now)
	p.CommunityHealthPercentage = m.GetHealthPercentage()

	s.ShiftWindow(now)
	return s.Flush()
}
