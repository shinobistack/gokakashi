package experiment

import "strings"

// Experiments is a comma-separated experiment flags string, set via CLI flag in root.go
var Experiments string

type Experiment string

const (
	V2Agents Experiment = "v2_agents"
)

// knownExperiments contains all valid experiments
var knownExperiments = map[Experiment]struct{}{
	V2Agents: {},
}

// userChosenExperiments returns the experiments as a slice of strings
func userChosenExperiments() map[Experiment]struct{} {
	expSet := make(map[Experiment]struct{})
	for _, exp := range strings.Split(Experiments, ",") {
		exp = strings.TrimSpace(exp)
		e := Experiment(exp)

		if _, ok := knownExperiments[e]; ok {
			expSet[e] = struct{}{}
		}
	}

	return expSet
}

// Enabled returns true if the experiment is enabled by the user
func Enabled(exp Experiment) bool {
	expSet := userChosenExperiments()
	_, ok := expSet[exp]
	return ok
}
