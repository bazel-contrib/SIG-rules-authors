package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/bazel-contrib/SIG-rules-authors/experimental/rules-keeper/timeseries"

	"github.com/golang/glog"
	"github.com/google/go-github/v50/github"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/bazel-contrib/SIG-rules-authors/experimental/rules-keeper/proto"
)

func (r *repo) UpdateVersions(ctx context.Context) error {
	glog.Infof("Updating versions for %s/%s", r.owner, r.name)

	var vs []*Version
	m := make(map[string]*Version)

	repo, _, err := r.cli.Repositories.Get(ctx, r.owner, r.name)
	if err != nil {
		return err
	}
	v, err := r.getVersionWithStats(ctx, repo.GetDefaultBranch())
	if err != nil {
		return err
	}
	vs = append(vs, v)

	for it := r.ListTags(ctx); ; {
		tag, err := it.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		v, err := r.getVersionWithStats(ctx, tag.GetName())
		if err != nil {
			return err
		}
		vs = append(vs, v)
		m[tag.GetName()] = v
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
		v.Releases = append(v.Releases, &pb.Release{
			Title:       rel.GetName(),
			Description: rel.GetBody(),
			Preprelease: rel.GetPrerelease(),
			PublishTime: timestamppb.New(rel.GetPublishedAt().Time),
			Assets:      as,
		})
		v.stats.GetOrCreatePointAt(now).DowloadCount += dlc
		v.stats.SetUpdateTime(now)
	}

	for _, v := range vs {
		if err := v.Flush(); err != nil {
			return err
		}
	}

	return nil
}

type Version struct {
	owner string
	repo  string
	ref   string
	*pb.Version
	stats *timeseries.Store[*VersionStatsPoint]
}

// TODO(@ashi009): refactor this.
func (v *Version) Flush() error {
	b, err := prototext.MarshalOptions{Multiline: true}.Marshal(v.Version)
	if err != nil {
		return err
	}
	name := fmt.Sprintf("store/%s/%s/versions/%s/METADATA", v.owner, v.repo, v.ref)
	if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(name, b, 0644); err != nil {
		return err
	}
	return v.stats.Flush()
}

func (r *repo) getVersionWithStats(ctx context.Context, ref string) (*Version, error) {
	v, err := r.getVersion(ctx, ref)
	if err != nil {
		return nil, err
	}
	s, err := r.loadVersionStatsStore(ctx, ref)
	if err != nil {
		return nil, err
	}
	return &Version{
		owner:   r.owner,
		repo:    r.name,
		ref:     ref,
		Version: v,
		stats:   s,
	}, nil
}

func (r *repo) getVersion(ctx context.Context, ref string) (*pb.Version, error) {
	c, _, err := r.cli.Repositories.GetCommit(ctx, r.owner, r.name, ref, &github.ListOptions{})
	if err != nil {
		return nil, err
	}
	rf, _, err := r.cli.Repositories.GetReadme(ctx, r.owner, r.name, &github.RepositoryContentGetOptions{
		Ref: ref,
	})
	// TODO(@ashi009): handle 404 (no readme).
	if err != nil {
		return nil, err
	}
	rc, err := rf.GetContent()
	if err != nil {
		return nil, err
	}
	mf, err := r.getModuleFile(ctx, ref)
	if err != nil {
		return nil, err
	}
	return &pb.Version{
		Ref:        ref,
		Sha:        c.GetSHA(),
		Readme:     rc,
		ModuleFile: mf,
	}, nil
}

func (r *repo) getModuleFile(ctx context.Context, ref string) (*pb.ModuleFile, error) {
	fc, _, _, err := r.cli.Repositories.GetContents(ctx, r.owner, r.name, "MODULE.bazel", &github.RepositoryContentGetOptions{
		Ref: ref,
	})
	// It's ok if the module file doesn't exist.
	if err != nil {
		return nil, nil
	}
	rc, err := fc.GetContent()
	if err != nil {
		return nil, err
	}
	// TODO(@ashi009): parse the module file.
	_ = rc
	return &pb.ModuleFile{}, nil
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

func (r *repo) loadVersionStatsStore(ctx context.Context, tag string) (*timeseries.Store[*VersionStatsPoint], error) {
	return timeseries.Load(VersionStats, fmt.Sprintf("store/%s/%s/versions/%s/metrics/version_stats", r.owner, r.name, tag))
}
