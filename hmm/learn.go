package hmm

type updatedHmmModel struct {
	StartProbability      []float64
	TransitionNumerator   [][]float64
	TransitionDenominator []float64
	EmissionNumerator     []map[string]float64
	EmissionDenominator   []float64
}

func newUpdatedHmmModel(stateNum int, observations []string) *updatedHmmModel {
	model := &updatedHmmModel{}
	model.StartProbability = make([]float64, stateNum)
	model.TransitionNumerator = make([][]float64, stateNum)
	model.TransitionDenominator = make([]float64, stateNum)
	for i := 0; i < stateNum; i++ {
		model.TransitionNumerator[i] = make([]float64, stateNum)
	}
	model.EmissionNumerator = make([]map[string]float64, stateNum)
	model.EmissionDenominator = make([]float64, stateNum)
	for i := 0; i < stateNum; i++ {
		model.EmissionNumerator[i] = make(map[string]float64)
	}
	return model
}

func forwardBackward(model *HmmModel, observations []string, labels []string) *updatedHmmModel {
	forwardMatrix := forward(model, observations)
	backwardMatrix := backward(model, observations)

	statePerTime := make([][]float64, len(observations))
	for t := 0; t < len(observations); t++ {
		statePerTime[t] = make([]float64, len(model.States))
		total := 0.0
		for i := 0; i < len(model.States); i++ {
			statePerTime[t][i] = forwardMatrix[t][i] * backwardMatrix[t][i]
			total += statePerTime[t][i]
		}
		for i := 0; i < len(model.States); i++ {
			statePerTime[t][i] /= total
		}
	}
	transitionPerTime := make([][][]float64, len(observations))
	for t := 0; t < len(observations); t++ {
		transitionPerTime[t] = make([][]float64, len(model.States))
		total := 0.0
		for i, stateI := range model.States {
			transitionPerTime[t][i] = make([]float64, len(model.States))
			for j, stateJ := range model.States {
				transitionPerTime[t][i][j] = forwardMatrix[t][i] * model.TransitionProbability[stateI][stateJ] *
					model.EmissionProbability[stateJ][observations[t]] * backwardMatrix[t][j]
				total += transitionPerTime[t][i][j]
			}
		}
		for i, _ := range model.States {
			for j, _ := range model.States {
				transitionPerTime[t][i][j] /= total
			}
		}
	}

	updatedModel := newUpdatedHmmModel(len(model.States), observations)
	for i, _ := range model.States {
		updatedModel.StartProbability[i] = statePerTime[0][i]

		lastIndex := len(observations) - 1
		for j, _ := range model.States {
			updatedModel.TransitionNumerator[i][j] = 0.0
			for t := 0; t < lastIndex; t++ {
				updatedModel.TransitionNumerator[i][j] += transitionPerTime[t][i][j]
			}
		}

		for t, labelState := range labels {
			if _, ok := updatedModel.EmissionNumerator[i][labelState]; !ok {
				updatedModel.EmissionNumerator[i][labelState] = 0.0
			}
			updatedModel.EmissionNumerator[i][labelState] += statePerTime[t][i]
		}

		updatedModel.TransitionDenominator[i] = 0.0
		for t := 0; t < lastIndex; i++ {
			updatedModel.TransitionDenominator[i] += statePerTime[t][i]
		}
		updatedModel.EmissionDenominator[i] = updatedModel.TransitionDenominator[i] + statePerTime[lastIndex][i]
	}
	return updatedModel
}

func backward(model *HmmModel, observations []string) [][]float64 {
	backwardMatrix := make([][]float64, len(observations))
	for i := 0; i < len(observations); i++ {
		backwardMatrix[i] = make([]float64, len(model.States))
	}
	lastIndex := len(observations) - 1
	for i := 0; i < len(model.States); i++ {
		backwardMatrix[lastIndex][i] = 1.0
	}

	for i := lastIndex - 1; i >= 0; i-- {
		for j, state := range model.States {
			backwardMatrix[i][j] = 0.0
			for nextJ, nextState := range model.States {
				backwardMatrix[i][j] += model.TransitionProbability[state][nextState] * backwardMatrix[i+1][nextJ] *
					model.EmissionProbability[nextState][observations[i+1]]
			}
		}
	}

	return backwardMatrix
}
