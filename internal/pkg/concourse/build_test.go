package concourse

import "testing"

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
