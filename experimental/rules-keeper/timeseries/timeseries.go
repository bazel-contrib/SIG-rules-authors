package timeseries

import (
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/gocarina/gocsv"
	intervalpb "google.golang.org/genproto/googleapis/type/interval"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/bazel-contrib/SIG-rules-authors/experimental/rules-keeper/timeseries/proto"
)

// Point is a point in a time series. The struct that implements this interface
// should be a valid CSV row that gocsv can marshal/unmarshal.
type Point interface {
	Timestamp() time.Time
}

// AlignFunc aligns a timestamp to a given time unit.
type AlignFunc func(time.Time) time.Time

// Descriptor describes a time series.
type Descriptor[P Point] struct {
	// Align aligns a timestamp to a given time unit.
	Align AlignFunc
	// NewPoint creates a new point for a given timestamp.
	NewPoint func(time.Time) P
	// Retention is the maximum age of a point in the time series.
	Retention time.Duration
}

// Store is a time series store. Points in the store are distinguised by their
// timestamp, which is saying the point doesn't have labels. To store timeseries
// with labels, use multiple stores.
type Store[P Point] struct {
	name     string
	desc     *Descriptor[P]
	metadata pb.Metadata
	idx      map[int64]P
	points   []P
}

// Load loads a time series from disk.
func Load[P Point](desc *Descriptor[P], name string) (*Store[P], error) {
	ts := &Store[P]{
		name: name,
		desc: desc,
		idx:  make(map[int64]P),
	}
	b, err := os.ReadFile(ts.name + ".csv")
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if len(b) > 0 {
		if err := gocsv.UnmarshalBytes(b, &ts.points); err != nil {
			return nil, err
		}
		for _, p := range ts.points {
			// This assumes all timestamps are already aligned.
			ts.idx[p.Timestamp().Unix()] = p
		}
	}
	b, err = os.ReadFile(ts.name + ".metadata")
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if err := prototext.Unmarshal(b, &ts.metadata); err != nil {
		return nil, err
	}
	return ts, nil
}

// Flush writes the time series to disk.
func (ts *Store[P]) Flush() error {
	sort.Slice(ts.points, func(i, j int) bool {
		return ts.points[i].Timestamp().Before(ts.points[j].Timestamp())
	})
	if err := os.MkdirAll(filepath.Dir(ts.name), 0755); err != nil {
		return err
	}
	b, err := gocsv.MarshalBytes(ts.points)
	if err != nil {
		return err
	}
	if err := os.WriteFile(ts.name+".csv", b, 0644); err != nil {
		return err
	}
	b, err = prototext.MarshalOptions{Multiline: true}.Marshal(&ts.metadata)
	if err != nil {
		return err
	}
	return os.WriteFile(ts.name+".metadata", b, 0644)
}

// Window returns the start and end time of the time series window.
func (ts *Store[P]) Window() (start, end time.Time) {
	if ts.metadata.Window == nil {
		return time.Time{}, time.Now().Add(-ts.desc.Retention)
	}
	return ts.metadata.Window.StartTime.AsTime(), ts.metadata.Window.EndTime.AsTime()
}

// ShiftWindow shifts the time series window to the given end time.
func (ts *Store[P]) ShiftWindow(end time.Time) {
	start := end.Add(-ts.desc.Retention)
	ts.metadata.Window = &intervalpb.Interval{
		StartTime: timestamppb.New(start),
		EndTime:   timestamppb.New(end),
	}
	ts.removeStatePoints(start)
}

// removeStalePoints removes points that are before notBefore.
func (ts *Store[P]) removeStatePoints(start time.Time) {
	for len(ts.points) > 0 {
		t := ts.points[0].Timestamp()
		if !t.Before(start) {
			break
		}
		ts.points = ts.points[1:]
		delete(ts.idx, t.Unix())
	}
}

// GetOrCreatePointAt returns the point at the given timestamp, creating it if
// it doesn't exist.
func (ts *Store[P]) GetOrCreatePointAt(t time.Time) P {
	t = ts.desc.Align(t)
	p, ok := ts.idx[t.Unix()]
	if !ok {
		p = ts.desc.NewPoint(t)
		ts.idx[t.Unix()] = p
		ts.points = append(ts.points, p)
	}
	return p
}

// Date is a time that represents an UTC date that is suitable for CSV
// marshalling.
//
// N.B. It would be nice to add Timestamp method to this, and make any struct
// that embeds it conform to Point interface. But gocsv doesn't support
// embedding.
type Date struct{ time.Time }

func (t *Date) MarshalCSV() (string, error) {
	return t.Time.UTC().Format(time.DateOnly), nil
}

func (t *Date) UnmarshalCSV(s string) (err error) {
	t.Time, err = time.ParseInLocation(time.DateOnly, s, time.UTC)
	return
}

// Date is a time that represents an UTC date with time at second precision that
// is suitable for CSV marshalling.
//
// N.B. See comments for Date.
type DateTime struct{ time.Time }

func (t *DateTime) MarshalCSV() (string, error) {
	return t.Time.UTC().Format(time.DateTime), nil
}

func (t *DateTime) UnmarshalCSV(s string) (err error) {
	t.Time, err = time.ParseInLocation(time.DateTime, s, time.UTC)
	return
}

// AlignToDayInUTC aligns a timestamp to the start of the day.
func AlignToDayInUTC(t time.Time) time.Time {
	y, m, d := t.UTC().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

// AlignToISO8601WeekStartDayInUTC aligns a timestamp to the start of the week
// by the definition of ISO 8601.
func AlignToISO8601WeekStartDayInUTC(t time.Time) time.Time {
	// A UTC week starts on monday.
	t = t.UTC()
	wd := int(t.Weekday() - time.Monday)
	if wd < 0 { // Sunday
		wd += 7
	}
	t = t.AddDate(0, 0, -wd)
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

// AlignToSecond aligns a timestamp to the start of the second.
func AlignToSecond(t time.Time) time.Time {
	return t.Truncate(time.Second)
}

// AlignToDuration returns a AlignFunc that aligns a timestamp to a given
// duration. See Time.Truncate for details.
func AlignToDuration(d time.Duration) AlignFunc {
	return func(t time.Time) time.Time {
		return t.Truncate(d)
	}
}
