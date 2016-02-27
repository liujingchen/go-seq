package hmm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewHmmModel(t *testing.T) {
	model, err := NewHmmModel([]string{"S1", "S2"}, []string{"A", "B"})
	require.NoError(t, err)
	require.Equal(t, []string{"S1", "S2"}, model.States)
	require.Equal(t, []string{"A", "B"}, model.Observations)
	require.NotNil(t, model.TransitionProbability[0])
	require.NotNil(t, model.TransitionProbability[1])
	require.NotNil(t, model.EmissionProbability[0])
	require.NotNil(t, model.EmissionProbability[1])
}

func TestNewHmmModelError(t *testing.T) {
	_, err := NewHmmModel(nil, []string{"A", "B"})
	require.Error(t, err)
	_, err = NewHmmModel([]string{"S1", "S2"}, nil)
	require.Error(t, err)
	_, err = NewHmmModel([]string{"S1"}, []string{"A", "B"})
	require.Error(t, err)
	_, err = NewHmmModel([]string{"S1", "S2"}, []string{"A"})
	require.Error(t, err)
}
