package hmm

import "errors"

func Evaluate(model *HmmModel, observations []string) (float64, error) {
	if model == nil || observations == nil || len(observations) == 0 {
		return 0.0, errors.New("Model and observations can not be nil or empty.")
	}
	forwardMatrix := forward(model, observations)
	lastIndex := len(observations) - 1
	result := 0.0
	for i := 0; i < len(model.States); i++ {
		result += forwardMatrix[lastIndex][i]
	}
	return result, nil
}

func forward(model *HmmModel, observations []string) [][]float64 {
	forwardMatrix := make([][]float64, len(observations))
	for i := 0; i < len(observations); i++ {
		forwardMatrix[i] = make([]float64, len(model.States))
	}

	for i, _ := range model.States {
		forwardMatrix[0][i] = model.StartProbability[i] * model.EmissionProbability[i][observations[0]]
	}

	for i := 1; i < len(observations); i++ {
		for j, _ := range model.States {
			transitionPro := 0.0
			for k, _ := range model.States {
				transitionPro += forwardMatrix[i-1][k] * model.TransitionProbability[k][j]
			}
			forwardMatrix[i][j] = transitionPro * model.EmissionProbability[j][observations[i]]
		}
	}
	return forwardMatrix
}
