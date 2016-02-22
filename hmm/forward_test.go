package hmm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func createExampleModelForward() *HmmModel {
	// Example in http://www.cnblogs.com/tornadomeet/archive/2012/03/24/2415583.html
	model, _ := NewHmmModel([]string{"S1", "S2", "S3"}, []string{"A", "B"})
	model.StartProbability["S1"] = 1.0
	model.TransitionProbability["S1"]["S1"] = 0.4
	model.TransitionProbability["S1"]["S2"] = 0.6
	model.TransitionProbability["S2"]["S2"] = 0.8
	model.TransitionProbability["S2"]["S3"] = 0.2
	model.TransitionProbability["S3"]["S3"] = 1.0
	model.EmissionProbability["S1"]["A"] = 0.7
	model.EmissionProbability["S1"]["B"] = 0.3
	model.EmissionProbability["S2"]["A"] = 0.4
	model.EmissionProbability["S2"]["B"] = 0.6
	model.EmissionProbability["S3"]["A"] = 0.8
	model.EmissionProbability["S3"]["B"] = 0.2
	return model
}

func TestForward(t *testing.T) {
	model := createExampleModelForward()
	result, err := GetProbabilityOfObservation(model, []string{"A", "B", "A", "B"})
	require.NoError(t, err)
	require.InDelta(t, 0.0717696, result, 0.000000001)
}

func TestForwardLenOne(t *testing.T) {
	model := createExampleModelForward()
	result, err := GetProbabilityOfObservation(model, []string{"A"})
	require.NoError(t, err)
	require.InDelta(t, 0.7, result, 0.000000001)
}

func TestForwardError(t *testing.T) {
	model := createExampleModelForward()
	_, err := GetProbabilityOfObservation(model, []string{})
	require.Error(t, err)
}
