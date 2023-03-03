package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/bazel-contrib/SIG-rules-authors/experimental/rules-keeper/timeseries"

	"github.com/golang/glog"
	"github.com/google/go-github/v50/github"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/bazel-contrib/SIG-rules-authors/experimental/rules-keeper/proto"
)

func (r *repo) UpdateVersions(ctx context.Context) error {
	glog.Infof("Updating versions for %s/%s", r.owner, r.name)

	m := make(map[string]*VersionStore)

	repo, _, err := r.cli.Repositories.Get(ctx, r.owner, r.name)
	if err != nil {
		return err
	}
	vs, err := r.loadVersionStore(ctx, repo.GetDefaultBranch())
	if err != nil {
		return err
	}
	v, err := r.getVersion(ctx, repo.GetDefaultBranch())
	if err != nil {
		return err
	}
	vs.Version.Reset()
	proto.Merge(&vs.Version, v)
	vs.Flush()

	for it := r.ListTags(ctx); ; {
		tag, err := it.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		vs, err := r.loadVersionStore(ctx, tag.GetName())
		if err != nil {
			return err
		}
		v, err := r.getVersion(ctx, repo.GetDefaultBranch())
		if err != nil {
			return err
		}
		vs.Version.Reset()
		proto.Merge(&vs.Version, v)
		m[tag.GetName()] = vs
	}

	now := time.Now()
	for it := r.ListReleases(ctx); ; {
		rel, err := it.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		v, ok := m[rel.GetTagName()]
		if !ok {
			return fmt.Errorf("release %s attaches to an unknown tag %s", rel.GetName(), rel.GetTagName())
		}
		as := make([]*pb.Release_Asset, 0, len(rel.Assets))
		dlc := 0
		for _, a := range rel.Assets {
			if !releaseContentType[a.GetContentType()] {
				continue
			}
			as = append(as, &pb.Release_Asset{
				Name: a.GetName(),
				Url:  a.GetBrowserDownloadURL(),
			})
			dlc += a.GetDownloadCount()
		}
		v.Release = &pb.Release{
			Title:       rel.GetName(),
			Description: rel.GetBody(),
			Preprelease: rel.GetPrerelease(),
			PublishTime: timestamppb.New(rel.GetPublishedAt().Time),
			Assets:      as,
		}
		v.stats.GetOrCreatePointAt(now).DowloadCount += dlc
		v.stats.ShiftWindow(now)
		if err := v.Flush(); err != nil {
			return err
		}
	}

	return nil
}

type VersionStore struct {
	name string
	pb.Version
	stats *timeseries.Store[*VersionStatsPoint]
}

func (r *repo) loadVersionStore(ctx context.Context, tag string) (*VersionStore, error) {
	name := fmt.Sprintf("store/%s/%s/versions/%s", r.owner, r.name, tag)
	b, err := os.ReadFile(name + "/METADATA")
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	s := &VersionStore{
		name: name,
	}
	if err := prototext.Unmarshal(b, &s.Version); err != nil {
		return nil, err
	}
	if s.stats, err = timeseries.Load(VersionStats, fmt.Sprintf("%s/metrics/version_stats", name)); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *VersionStore) Flush() error {
	if err := os.MkdirAll(s.name, 0755); err != nil {
		return err
	}
	b, err := prototext.MarshalOptions{Multiline: true}.Marshal(&s.Version)
	if err != nil {
		return err
	}
	if err := os.WriteFile(s.name+"/METADATA", b, 0644); err != nil {
		return err
	}
	return s.stats.Flush()
}

func (r *repo) getVersion(ctx context.Context, ref string) (*pb.Version, error) {
	c, _, err := r.cli.Repositories.GetCommit(ctx, r.owner, r.name, ref, &github.ListOptions{})
	if err != nil {
		return nil, err
	}
	readme, err := r.GetReadme(ctx, ref)
	if err != nil && !isNotFound(err) {
		return nil, err
	}
	m, err := r.getModuleFile(ctx, ref)
	if err != nil {
		return nil, err
	}
	return &pb.Version{
		Ref:        ref,
		Sha:        c.GetSHA(),
		Readme:     readme,
		ModuleFile: m,
	}, nil
}

func (r *repo) getModuleFile(ctx context.Context, ref string) (*pb.ModuleFile, error) {
	c, err := r.GetContent(ctx, ref, "MODULE.bazel")
	if err != nil {
		if isNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return parseModuleFile("MODULE.bazel", c)
}

type VersionStatsPoint struct {
	Time         timeseries.DateTime `csv:"time (UTC)"`
	DowloadCount int                 `csv:"# download"`
}

func (p *VersionStatsPoint) Timestamp() time.Time {
	return p.Time.Time
}

var VersionStats = &timeseries.Descriptor[*VersionStatsPoint]{
	NewPoint: func(t time.Time) *VersionStatsPoint {
		return &VersionStatsPoint{Time: timeseries.DateTime{Time: t}}
	},
	Align: timeseries.AlignToSecond,
}

func isNotFound(err error) bool {
	if err, ok := err.(*github.ErrorResponse); ok && err.Response.StatusCode == http.StatusNotFound {
		return true
	}
	return false
}
