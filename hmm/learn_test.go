package hmm

import (
	"bufio"
	"indeed/gophers/3rdparty/p/github.com/stretchr/testify/require"
	"log"
	"os"
	"strings"
	"testing"
)

func readBerkeleyTestData(filename string) ([]string, []string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	observations := make([]string, 0)
	labels := make([]string, 0)
	for {
		bytes, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		line := string(bytes)
		parts := strings.Split(line, ",")
		labels = append(labels, parts[0])
		observations = append(observations, parts[1])
	}
	return observations, labels
}

func TestLearnUsingBerkeleyData(t *testing.T) {
	trainingData := make([][]string, 2)
	labels := make([][]string, 2)
	trainingData[0], labels[0] = readBerkeleyTestData("weather-test1-1000.txt")
	trainingData[1], labels[1] = readBerkeleyTestData("weather-test2-1000.txt")

	model, err := Train(trainingData, labels, 0.0001)
	require.NoError(t, err)
	require.NotNil(t, model)
	log.Print(model.States)
	log.Print(model.Observations)
	log.Print(model.StartProbability)
	log.Print(model.TransitionProbability)
	log.Print(model.EmissionProbability)
	answer, err := Decode(model, []string{"no", "no", "no", "yes", "no", "no", "yes", "yes", "no", "yes"})
	require.NoError(t, err)
	log.Print(answer)
	//require.Equal(t, []string{}, answer)
}
