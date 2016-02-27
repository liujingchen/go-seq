package hmm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func createExampleModelEvaluate() *HmmModel {
	// Example in http://www.cnblogs.com/tornadomeet/archive/2012/03/24/2415583.html
	model, _ := NewHmmModel([]string{"S1", "S2", "S3"}, []string{"A", "B"})
	model.StartProbability[0] = 1.0
	model.TransitionProbability[0][0] = 0.4
	model.TransitionProbability[0][1] = 0.6
	model.TransitionProbability[1][1] = 0.8
	model.TransitionProbability[1][2] = 0.2
	model.TransitionProbability[2][2] = 1.0
	model.EmissionProbability[0]["A"] = 0.7
	model.EmissionProbability[0]["B"] = 0.3
	model.EmissionProbability[1]["A"] = 0.4
	model.EmissionProbability[1]["B"] = 0.6
	model.EmissionProbability[2]["A"] = 0.8
	model.EmissionProbability[2]["B"] = 0.2
	return model
}

func TestEvaluate(t *testing.T) {
	model := createExampleModelEvaluate()
	result, err := Evaluate(model, []string{"A", "B", "A", "B"})
	require.NoError(t, err)
	require.InDelta(t, 0.0717696, result, 0.000000001)
}

func TestEvaluateLenOne(t *testing.T) {
	model := createExampleModelEvaluate()
	result, err := Evaluate(model, []string{"A"})
	require.NoError(t, err)
	require.InDelta(t, 0.7, result, 0.000000001)
}

func TestEvaluateError(t *testing.T) {
	model := createExampleModelEvaluate()
	_, err := Evaluate(model, []string{})
	require.Error(t, err)
}
