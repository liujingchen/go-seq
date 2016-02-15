package hmm

import "errors"

type HmmModel struct {
	States                []string
	Observations          []string
	StartProbability      map[string]float64
	TransitionProbability map[string]map[string]float64
	EmissionProbability   map[string]map[string]float64
}

func NewHmmModel(states, observations []string) (*HmmModel, error) {
	if states == nil || observations == nil || len(states) <= 1 || len(observations) <= 1 {
		return nil, errors.New("States and observations must have more than 1 elements.")
	}
	model := &HmmModel{States: states, Observations: observations}
	model.StartProbability = make(map[string]float64)
	for _, state := range states {
		model.StartProbability[state] = 0.0
	}
	model.TransitionProbability = make(map[string]map[string]float64)
	for _, stateFrom := range states {
		model.TransitionProbability[stateFrom] = make(map[string]float64)
		for _, stateTo := range states {
			model.TransitionProbability[stateFrom][stateTo] = 0.0
		}
	}
	model.EmissionProbability = make(map[string]map[string]float64)
	for _, state := range states {
		model.EmissionProbability[state] = make(map[string]float64)
		for _, observation := range observations {
			model.EmissionProbability[state][observation] = 0.0
		}
	}
	return model, nil
}
