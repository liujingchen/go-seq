package hmm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func createExampleModelViterbi() *HmmModel {
	// Example from https://en.wikipedia.org/wiki/Viterbi_algorithm#Example
	model, _ := NewHmmModel([]string{"Healthy", "Fever"}, []string{"normal", "cold", "dizzy"})
	model.StartProbability["Healthy"] = 0.6
	model.StartProbability["Fever"] = 0.4
	model.TransitionProbability["Healthy"]["Healthy"] = 0.7
	model.TransitionProbability["Healthy"]["Fever"] = 0.3
	model.TransitionProbability["Fever"]["Healthy"] = 0.4
	model.TransitionProbability["Fever"]["Fever"] = 0.6
	model.EmissionProbability["Healthy"]["normal"] = 0.5
	model.EmissionProbability["Healthy"]["cold"] = 0.4
	model.EmissionProbability["Healthy"]["dizzy"] = 0.1
	model.EmissionProbability["Fever"]["normal"] = 0.1
	model.EmissionProbability["Fever"]["cold"] = 0.3
	model.EmissionProbability["Fever"]["dizzy"] = 0.6
	return model
}

func TestViterbi(t *testing.T) {
	model := createExampleModelViterbi()
	states, err := GetMostLikelyStates(model, []string{"normal", "cold", "dizzy"})
	require.NoError(t, err)
	require.Equal(t, []string{"Healthy", "Healthy", "Fever"}, states)
}

func TestViterbiLenghOne(t *testing.T) {
	model := createExampleModelViterbi()
	states, err := GetMostLikelyStates(model, []string{"normal"})
	require.NoError(t, err)
	require.Equal(t, []string{"Healthy"}, states)
}

func TestViterbiError(t *testing.T) {
	_, err := GetMostLikelyStates(nil, []string{"normal", "cold", "dizzy"})
	require.Error(t, err)
	model := createExampleModelViterbi()
	_, err = GetMostLikelyStates(model, nil)
	require.Error(t, err)
	_, err = GetMostLikelyStates(model, []string{})
	require.Error(t, err)
}
