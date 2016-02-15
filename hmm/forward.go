package hmm

import "errors"

func GetProbabilityOfObservation(model *HmmModel, observations []string) (float64, error) {
	if model == nil || observations == nil || len(observations) == 0 {
		return 0.0, errors.New("Model and observations can not be nil or empty.")
	}
	a := make([][]float64, len(observations))
	for i := 0; i < len(observations); i++ {
		a[i] = make([]float64, len(model.States))
	}

	for i, state := range model.States {
		a[0][i] = model.StartProbability[state] * model.EmissionProbability[state][observations[0]]
	}

	for i := 1; i < len(observations); i++ {
		for j, state := range model.States {
			transitionPro := 0.0
			for k, preState := range model.States {
				transitionPro += a[i-1][k] * model.TransitionProbability[preState][state]
			}
			a[i][j] = transitionPro * model.EmissionProbability[state][observations[i]]
		}
	}
	lastIndex := len(observations) - 1
	result := 0.0
	for i := 0; i < len(model.States); i++ {
		result += a[lastIndex][i]
	}
	return result, nil
}
