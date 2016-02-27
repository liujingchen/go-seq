package hmm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func createExampleModelDecode() *HmmModel {
	// Example from https://en.wikipedia.org/wiki/Viterbi_algorithm#Example
	model, _ := NewHmmModel([]string{"Healthy", "Fever"}, []string{"normal", "cold", "dizzy"})
	model.StartProbability[0] = 0.6
	model.StartProbability[1] = 0.4
	model.TransitionProbability[0][0] = 0.7
	model.TransitionProbability[0][1] = 0.3
	model.TransitionProbability[1][0] = 0.4
	model.TransitionProbability[1][1] = 0.6
	model.EmissionProbability[0]["normal"] = 0.5
	model.EmissionProbability[0]["cold"] = 0.4
	model.EmissionProbability[0]["dizzy"] = 0.1
	model.EmissionProbability[1]["normal"] = 0.1
	model.EmissionProbability[1]["cold"] = 0.3
	model.EmissionProbability[1]["dizzy"] = 0.6
	return model
}

func TestDecode(t *testing.T) {
	model := createExampleModelDecode()
	states, err := Decode(model, []string{"normal", "cold", "dizzy"})
	require.NoError(t, err)
	require.Equal(t, []string{"Healthy", "Healthy", "Fever"}, states)
}

func TestDecodeLenghOne(t *testing.T) {
	model := createExampleModelDecode()
	states, err := Decode(model, []string{"normal"})
	require.NoError(t, err)
	require.Equal(t, []string{"Healthy"}, states)
}

func TestDecodeError(t *testing.T) {
	_, err := Decode(nil, []string{"normal", "cold", "dizzy"})
	require.Error(t, err)
	model := createExampleModelDecode()
	_, err = Decode(model, nil)
	require.Error(t, err)
	_, err = Decode(model, []string{})
	require.Error(t, err)
}
