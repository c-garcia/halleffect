package publisher

import (
	"github.com/stretchr/testify/assert"
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
func TestFromConcourseJobStatus(t *testing.T) {
	const Concourse = "test-concourse"
	const SamplingTime = 10000
	testCases := []struct {
		Desc   string
		Status concourse.JobStatus
		Metric JobStatusMetric
	}{
		{
			Desc:   "Succeeded Job",
			Status: concourse.JobStatus{Id: 1, TeamName: "main", JobName: "j1", PipelineName: "p1", Status: "succeeded"},
			Metric: JobStatusMetric{Concourse: Concourse, TeamName: "main", PipelineName: "p1", JobName: "j1", Status: "up", SamplingTime: SamplingTime},
		},
		{
			Desc:   "Failed job",
			Status: concourse.JobStatus{Id: 1, TeamName: "main", JobName: "j1", PipelineName: "p1", Status: "failed"},
			Metric: JobStatusMetric{Concourse: Concourse, TeamName: "main", PipelineName: "p1", JobName: "j1", Status: "down", SamplingTime: SamplingTime},
		},
		{
			Desc:   "Errored job",
			Status: concourse.JobStatus{Id: 1, TeamName: "main", JobName: "j1", PipelineName: "p1", Status: "errored"},
			Metric: JobStatusMetric{Concourse: Concourse, TeamName: "main", PipelineName: "p1", JobName: "j1", Status: "down", SamplingTime: SamplingTime},
		},
		{
			Desc:   "Aborted job",
			Status: concourse.JobStatus{Id: 1, TeamName: "main", JobName: "j1", PipelineName: "p1", Status: "aborted"},
			Metric: JobStatusMetric{Concourse: Concourse, TeamName: "main", PipelineName: "p1", JobName: "j1", Status: "down", SamplingTime: SamplingTime},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Desc, func(t *testing.T) {
			assert.Equal(t, testCase.Metric, FromConcourseJobStatus(Concourse, SamplingTime, testCase.Status))
		})
	}
}

func TestFromConcourseJobStatus_PanicsForPendingAndUnknownStatuses(t *testing.T) {
	const Concourse = "test-concourse"
	const SamplingTime = 200000
	testCases := []struct {
		Desc   string
		Status concourse.JobStatus
	}{
		{
			Desc:   "Succeeded Job",
			Status: concourse.JobStatus{Id: 1, TeamName: "main", JobName: "j1", PipelineName: "p1", Status: "pending"},
		},
		{
			Desc:   "Failed job",
			Status: concourse.JobStatus{Id: 1, TeamName: "main", JobName: "j1", PipelineName: "p1", Status: "unknown"},
		},
		{
			Desc:   "Errored job",
			Status: concourse.JobStatus{Id: 1, TeamName: "main", JobName: "j1", PipelineName: "p1", Status: "ssewqe"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Desc, func(t *testing.T) {
			assert.Panics(t, func() { FromConcourseJobStatus(Concourse, SamplingTime, testCase.Status) })
		})
	}
}
