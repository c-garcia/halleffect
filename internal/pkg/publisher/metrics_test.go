package publisher

import (
	"reflect"
	"testing"

	"github.com/c-garcia/halleffect/internal/pkg/concourse"
)

const CONCOURSE_HOST = "concourse"

func TestFromConcourseBuild(t *testing.T) {
	tests := []struct {
		name string
		args concourse.Build
		want JobDurationMetric
	}{
		{"example 1",
			concourse.Build{StartTime: 100, EndTime: 101, PipelineName: "p1", JobName: "j1", Status: "finished"},
			JobDurationMetric{Timestamp: 100, EndTime: 101, PipelineName: "p1", JobName: "j1", Status: "finished", Concourse: CONCOURSE_HOST},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromConcourseBuild(CONCOURSE_HOST, tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromConcourseBuild() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetric_Duration(t *testing.T) {
	type fields struct {
		Timestamp int
		EndTime   int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{"simple example", fields{100, 120}, 20},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := JobDurationMetric{
				Timestamp: tt.fields.Timestamp,
				EndTime:   tt.fields.EndTime,
			}
			if got := m.Duration(); got != tt.want {
				t.Errorf("JobDurationMetric.Duration() = %v, want %v", got, tt.want)
			}
		})
	}
}
