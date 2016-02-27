package hmm

import (
	"log"
	"math"
)

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

func (m *updatedHmmModel) add(otherModel *updatedHmmModel) {
	for i, _ := range m.StartProbability {
		m.StartProbability[i] += otherModel.StartProbability[i]
		m.TransitionDenominator[i] += otherModel.TransitionDenominator[i]
		m.EmissionDenominator[i] += otherModel.EmissionDenominator[i]
		for j, _ := range m.StartProbability {
			m.TransitionNumerator[i][j] += otherModel.TransitionNumerator[i][j]
		}
		for observation, p := range otherModel.EmissionNumerator[i] {
			if _, ok := m.EmissionNumerator[i][observation]; !ok {
				m.EmissionNumerator[i][observation] = 0.0
			}
			m.EmissionNumerator[i][observation] += p
		}
	}
}

func (m *updatedHmmModel) toHmmModel(oldModel *HmmModel) *HmmModel {
	model, _ := newHmmModel(oldModel.States, oldModel.Observations)
	total := 0.0
	for i, _ := range model.States {
		total += m.StartProbability[i]
	}
	for i, _ := range model.States {
		model.StartProbability[i] = m.StartProbability[i] / total
	}
	for i, _ := range model.States {
		for j, _ := range model.States {
			model.TransitionProbability[i][j] = m.TransitionNumerator[i][j] / m.TransitionDenominator[i]
		}
	}
	for i, _ := range model.States {
		for observation, p := range m.EmissionNumerator[i] {
			model.EmissionProbability[i][observation] = p / m.EmissionDenominator[i]
		}
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
		for i, _ := range model.States {
			transitionPerTime[t][i] = make([]float64, len(model.States))
			for j, _ := range model.States {
				transitionPerTime[t][i][j] = forwardMatrix[t][i] * model.TransitionProbability[i][j] *
					model.EmissionProbability[j][observations[t]] * backwardMatrix[t][j]
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

		for t, observation := range observations {
			if _, ok := updatedModel.EmissionNumerator[i][observation]; !ok {
				updatedModel.EmissionNumerator[i][observation] = 0.0
			}
			updatedModel.EmissionNumerator[i][observation] += statePerTime[t][i]
		}

		updatedModel.TransitionDenominator[i] = 0.0
		for t := 0; t < lastIndex; t++ {
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
		for j, _ := range model.States {
			backwardMatrix[i][j] = 0.0
			for nextJ, _ := range model.States {
				backwardMatrix[i][j] += model.TransitionProbability[j][nextJ] * backwardMatrix[i+1][nextJ] *
					model.EmissionProbability[nextJ][observations[i+1]]
			}
		}
	}

	return backwardMatrix
}

func Train(trainningData, labels [][]string, gamma float64) (*HmmModel, error) {
	states := getAllTokens(labels)
	observations := getAllTokens(trainningData)
	model, err := newHmmModel(states, observations)
	if err != nil {
		return nil, err
	}
	model.initialize()
	currentEvaluation := make([]float64, len(trainningData))
	confirm := make(chan bool)
	for i, _ := range trainningData {
		go func(i int) {
			currentEvaluation[i], _ = Evaluate(model, trainningData[i])
			confirm <- true
		}(i)
	}
	for i := 0; i < len(trainningData); i++ {
		<-confirm
	}
	iterNum := 0
	for {
		sumUpdatedModel := newUpdatedHmmModel(len(model.States), observations)
		updates := make(chan *updatedHmmModel)
		for i, _ := range trainningData {
			go func(i int) {
				updates <- forwardBackward(model, trainningData[i], labels[i])
			}(i)
		}
		for i := 0; i < len(trainningData); i++ {
			updatedModel := <-updates
			sumUpdatedModel.add(updatedModel)
		}
		nextModel := sumUpdatedModel.toHmmModel(model)
		deltas := make(chan float64)
		for i, _ := range trainningData {
			go func(i int) {
				preEval := currentEvaluation[i]
				currentEvaluation[i], _ = Evaluate(nextModel, trainningData[i])
				deltas <- math.Abs((preEval - currentEvaluation[i]) / currentEvaluation[i])
			}(i)
		}
		delta := 0.0
		for i := 0; i < len(trainningData); i++ {
			delta += <-deltas
		}
		delta /= float64(len(trainningData))
		iterNum++
		log.Printf("Iteration %d, delta=%f", iterNum, delta)
		if delta < gamma {
			break
		}
		model = nextModel
	}
	return model, nil
}

func getAllTokens(data [][]string) []string {
	tokenMap := make(map[string]int)
	index := 0
	for _, tokenList := range data {
		for _, token := range tokenList {
			if _, ok := tokenMap[token]; !ok {
				tokenMap[token] = index
				index++
			}
		}
	}
	tokenList := make([]string, index)
	for token, index := range tokenMap {
		tokenList[index] = token
	}
	return tokenList
}
