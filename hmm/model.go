package hmm

import "errors"

type HmmModel struct {
	States                []string
	Observations          []string
	StartProbability      []float64
	TransitionProbability [][]float64
	EmissionProbability   []map[string]float64
}

func NewHmmModel(states, observations []string) (*HmmModel, error) {
	if states == nil || observations == nil || len(states) <= 1 || len(observations) <= 1 {
		return nil, errors.New("States and observations must have more than 1 elements.")
	}
	model := &HmmModel{States: states, Observations: observations}
	model.StartProbability = make([]float64, len(states))
	model.TransitionProbability = make([][]float64, len(states))
	for i, _ := range states {
		model.TransitionProbability[i] = make([]float64, len(states))
	}
	model.EmissionProbability = make([]map[string]float64, len(states))
	for i, _ := range states {
		model.EmissionProbability[i] = make(map[string]float64)
	}
	return model, nil
}
