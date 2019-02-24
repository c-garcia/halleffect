// +build !integration
// +build !service

package concourse

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBuild_Finished(t *testing.T) {
	instant1 := time.Now().Add(-24 * time.Hour)
	instant2 := instant1.Add(10 * time.Minute)
	var instant0 time.Time
	type fields struct {
		StartTime time.Time
		EndTime   time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"both fields are greater than zero", fields{instant1, instant2}, true},
		{"finished field is zero", fields{instant1, instant0}, false},
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
	type fields struct {
		StartTime time.Time
		EndTime   time.Time
	}
	tests := []struct {
		Name   string
		Fields fields
		Want   time.Duration
	}{
		{"0 duration", fields{time1, time1plus1sec}, 1 * time.Second},
		{">0 duration", fields{time1, time1}, 0 * time.Second},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			startTime := test.Fields.StartTime
			endTime := test.Fields.EndTime
			assert.Equal(t, test.Want, Build{StartTime: startTime, EndTime: endTime}.Duration())
		})
	}
}

func TestBuild_Succeeded(t *testing.T) {
	var instant0 time.Time
	instant1 := time.Now().Add(-24 * time.Hour)
	instant2 := instant1.Add(10 * time.Minute)
	type fields struct {
		Status    string
		StartTime time.Time
		EndTime   time.Time
	}
	tests := []struct {
		Name   string
		Fields fields
		Want   bool
	}{
		{"succeeded", fields{"succeeded", instant1, instant2}, true},
		{"failed", fields{"failed", instant1, instant2}, false},
		{"errored", fields{"errored", instant1, instant2}, false},
		{"started", fields{"started", instant1, instant0}, false},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, test.Want, Build{Status: test.Fields.Status, StartTime: test.Fields.StartTime, EndTime: test.Fields.EndTime}.Succeeded())
		})
	}

}
