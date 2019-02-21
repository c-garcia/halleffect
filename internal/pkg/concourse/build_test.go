//+build !integration, !service

package concourse

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBuild_Finished(t *testing.T) {
	type fields struct {
		StartTime int
		EndTime   int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"both fields are greater than zero", fields{1000, 1000}, false},
		{"finished field is zero", fields{1000, 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Build{
				StartTime: tt.fields.StartTime,
				EndTime:   tt.fields.EndTime,
			}
			if got := b.Finished(); got != tt.want {
				t.Errorf("Build.Finished() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuild_Duration(t *testing.T) {
	time1 := time.Now()
	time1plus1sec := time1.Add(1 * time.Second)
	time1plus1ms := time1.Add(1 * time.Millisecond)
	type fields struct {
		StartTime int
		EndTime   int
	}
	tests := []struct {
		Name   string
		Fields fields
		Want   time.Duration
	}{
		{"0 duration", fields{int(time1.Unix()), int(time1plus1sec.Unix())}, 1 * time.Second},
		{">0 duration", fields{int(time1.Unix()), int(time1.Unix())}, 0 * time.Second},
		{"second resolution", fields{int(time1.Unix()), int(time1plus1ms.Unix())}, 0 * time.Second},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, test.Want, Build{StartTime: test.Fields.StartTime, EndTime: test.Fields.EndTime}.Duration())
		})
	}
}

func TestBuild_Succeeded(t *testing.T) {
	type fields struct {
		Status    string
		StartTime int
		EndTime   int
	}
	tests := []struct {
		Name   string
		Fields fields
		Want   bool
	}{
		{"succeeded", fields{"succeeded", 1000, 1100}, true},
		{"failed", fields{"failed", 1000, 1100}, false},
		{"errored", fields{"errored", 1000, 1100}, false},
		{"started", fields{"started", 1000, 0}, false},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, test.Want, Build{Status: test.Fields.Status, StartTime: test.Fields.StartTime, EndTime: test.Fields.EndTime}.Succeeded())
		})
	}

}
