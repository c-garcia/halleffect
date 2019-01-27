package metrics

import (
	"reflect"
	"testing"

	"github.com/c-garcia/halleffect/internal/pkg/concourse"
)

func TestFromConcourseBuild(t *testing.T) {
	tests := []struct {
		name string
		args concourse.Build
		want Metric
	}{
		{"example 1", concourse.Build{StartTime: 100, EndTime: 101}, Metric{Timestamp: 100, EndTime: 101}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromConcourseBuild(tt.args); !reflect.DeepEqual(got, tt.want) {
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
			m := Metric{
				Timestamp: tt.fields.Timestamp,
				EndTime:   tt.fields.EndTime,
			}
			if got := m.Duration(); got != tt.want {
				t.Errorf("Metric.Duration() = %v, want %v", got, tt.want)
			}
		})
	}
}
