package nndatasetmaker

import (
	"fmt"

	"github.com/donyori/cqa/common/json"
)

type ExampleNumbersPerLabelPair struct {
	TrainDataset int32 `json:"train_dataset"`
	EvalDataset  int32 `json:"eval_dataset"`
}

type Settings struct {
	MaxExampleNumbersPerLabel   map[string]*ExampleNumbersPerLabelPair `json:"max_example_numbers_per_label"`
	DoesContainNoLabelQuestions bool                                   `json:"does_contain_no_label_questions"`
	LogStep                     int                                    `json:"log_step"`
}

const SettingsFilename string = "../settings/tool/nndatasetmaker.json"

var GlobalSettings Settings

func init() {
	// Default values:
	GlobalSettings.MaxExampleNumbersPerLabel = make(
		map[string]*ExampleNumbersPerLabelPair)
	GlobalSettings.DoesContainNoLabelQuestions = false
	GlobalSettings.LogStep = 1000
	epls := [...]*ExampleNumbersPerLabelPair{
		&ExampleNumbersPerLabelPair{TrainDataset: 5000, EvalDataset: 500},
		&ExampleNumbersPerLabelPair{TrainDataset: 1000, EvalDataset: 250},
		&ExampleNumbersPerLabelPair{TrainDataset: 500, EvalDataset: 150},
		&ExampleNumbersPerLabelPair{TrainDataset: 256, EvalDataset: 128},
	}
	for _, epl := range epls {
		GlobalSettings.MaxExampleNumbersPerLabel[fmt.Sprintf(
			"t%de%dpl", epl.TrainDataset, epl.EvalDataset)] = epl
	}

	_, err := json.DecodeJsonFromFileIfExist(
		SettingsFilename, &GlobalSettings)
	if err != nil {
		panic(err)
	}
}
