package experiment

import (
	"reflect"
	"testing"
)

func TestExperimentsList(t *testing.T) {
	tests := []struct {
		experimentsInput string
		want             map[Experiment]struct{}
		name             string
	}{
		{
			experimentsInput: "",
			want:             map[Experiment]struct{}{},
			name:             "empty string",
		},
		{
			experimentsInput: "v2_agents",
			want:             map[Experiment]struct{}{V2Agents: {}},
			name:             "single experiment",
		},
		{
			experimentsInput: "v2_agents,foo,bar",
			want: map[Experiment]struct{}{
				V2Agents: {},
			},
			name: "multiple experiments",
		},
		{
			experimentsInput: " v2_agents , foo , bar ",
			want: map[Experiment]struct{}{
				V2Agents: {},
			},
			name: "experiments with spaces",
		},
		{
			experimentsInput: "v2_agents,,foo,",
			want: map[Experiment]struct{}{
				V2Agents: {},
			},
			name: "duplicates and empty entries",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Experiments = tt.experimentsInput
			got := userChosenExperiments()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExperimentsList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnabled(t *testing.T) {
	cases := []struct {
		name        string
		experiments string
		exp         Experiment
		wantEnabled bool
	}{
		{
			name:        "enabled experiment",
			experiments: "v2_agents",
			exp:         V2Agents,
			wantEnabled: true,
		},
		{
			name:        "not enabled experiment",
			experiments: "",
			exp:         V2Agents,
			wantEnabled: false,
		},
		{
			name:        "unknown experiment",
			experiments: "foo",
			exp:         V2Agents,
			wantEnabled: false,
		},
		{
			name:        "multiple experiments, known included",
			experiments: "foo,v2_agents,bar",
			exp:         V2Agents,
			wantEnabled: true,
		},
		{
			name:        "multiple experiments, known not included",
			experiments: "foo,bar",
			exp:         V2Agents,
			wantEnabled: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			Experiments = tc.experiments
			if got := Enabled(tc.exp); got != tc.wantEnabled {
				t.Errorf("Enabled(%q) with Experiments=%q = %v; want %v", tc.exp, tc.experiments, got, tc.wantEnabled)
			}
		})
	}
}
