package hmm

import (
	"errors"
	"math/rand"
)

type HmmModel struct {
	States                []string
	Observations          []string
	StartProbability      []float64
	TransitionProbability [][]float64
	EmissionProbability   []map[string]float64
}

func (m *HmmModel) initialize() {
	m.StartProbability = randomFloats(len(m.States))
	for i, _ := range m.States {
		m.TransitionProbability[i] = randomFloats(len(m.States))
		randEmissions := randomFloats(len(m.Observations))
		for j, observation := range m.Observations {
			m.EmissionProbability[i][observation] = randEmissions[j]
		}
	}
}

func randomFloats(num int) []float64 {
	values := make([]float64, num)
	total := 0.0
	for i := 0; i < num; i++ {
		values[i] = rand.Float64()
		total += values[i]
	}
	for i := 0; i < num; i++ {
		values[i] /= total
	}
	return values
}

func newHmmModel(states, observations []string) (*HmmModel, error) {
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
