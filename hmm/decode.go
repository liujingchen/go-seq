package hmm

import "errors"

func Decode(model *HmmModel, observations []string) ([]string, error) {
	if model == nil || observations == nil || len(observations) == 0 {
		return nil, errors.New("Model and observations can not be empty.")
	}
	viterbiMatrix := viterbi(model, observations)
	result := make([]string, len(observations))
	for i := 0; i < len(observations); i++ {
		maxProbability := 0.0
		for j, state := range model.States {
			if viterbiMatrix[i][j] > maxProbability {
				result[i] = state
				maxProbability = viterbiMatrix[i][j]
			}
		}
	}
	return result, nil
}

func viterbi(model *HmmModel, observations []string) [][]float64 {
	viterbiMatrix := make([][]float64, len(observations))
	for i, _ := range viterbiMatrix {
		viterbiMatrix[i] = make([]float64, len(model.States))
	}
	for i, _ := range model.States {
		viterbiMatrix[0][i] = model.StartProbability[i] * model.EmissionProbability[i][observations[0]]
	}
	for i := 1; i < len(observations); i++ {
		for j, _ := range model.States {
			maxProbability := 0.0
			for preJ, _ := range model.States {
				probability := viterbiMatrix[i-1][preJ] * model.TransitionProbability[preJ][j]
				if maxProbability < probability {
					maxProbability = probability
				}
			}
			viterbiMatrix[i][j] = maxProbability * model.EmissionProbability[j][observations[i]]
		}
	}
	return viterbiMatrix
}
