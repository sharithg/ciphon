package parser

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version   string              `yaml:"version"`
	Jobs      map[string]Job      `yaml:"jobs"`
	Workflows map[string]Workflow `yaml:"workflows"`
}

type Job struct {
	Docker string        `yaml:"docker"`
	Node   string        `yaml:"node"`
	Steps  []StepWrapper `yaml:"steps"`
}

type Step struct {
	Type         string            `yaml:"type,omitempty"`
	Name         string            `yaml:"name,omitempty"`
	Command      string            `yaml:"command,omitempty"`
	Checkout     bool              `yaml:"checkout,omitempty"`
	RestoreCache *RestoreCacheStep `yaml:"restore_cache,omitempty"`
	SaveCache    *SaveCacheStep    `yaml:"save_cache,omitempty"`
	Run          *RunStep          `yaml:"run,omitempty"`
}

type StepWrapper struct {
	Step Step
}

func (w *StepWrapper) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Try unmarshaling as a string first
	var stepName string
	if err := unmarshal(&stepName); err == nil {
		w.Step = Step{Type: stepName}
		return nil
	}

	// Try unmarshaling as a complex step object
	var stepMap map[string]interface{}
	if err := unmarshal(&stepMap); err == nil {
		for key, value := range stepMap {
			w.Step.Type = key

			// Use type assertions to handle specific step types
			switch key {
			case "restore_cache":
				var restoreCache RestoreCacheStep
				mapData, _ := value.(map[string]interface{})
				if err := decodeMapToStruct(mapData, &restoreCache); err == nil {
					w.Step.RestoreCache = &restoreCache
				}
				w.Step.Name = restoreCache.Name
				w.Step.RestoreCache = &restoreCache

			case "save_cache":
				var saveCache SaveCacheStep
				mapData, _ := value.(map[string]interface{})
				if err := decodeMapToStruct(mapData, &saveCache); err == nil {
					w.Step.SaveCache = &saveCache
				}
				w.Step.Name = saveCache.Name
				w.Step.SaveCache = &saveCache

			case "run":
				var runStep RunStep
				mapData, _ := value.(map[string]interface{})
				if err := decodeMapToStruct(mapData, &runStep); err == nil {
					w.Step.Run = &runStep
				}
				w.Step.Name = runStep.Name
				w.Step.Command = runStep.Command

			default:
				// For any unhandled complex type, directly unmarshal into Step
				if err := unmarshal(&w.Step); err != nil {
					return err
				}
			}
		}
		return nil
	}

	return unmarshal(&w.Step)
}

func (w *JobWithRequires) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var jobName string
	if err := unmarshal(&jobName); err == nil {
		w.Name = jobName
		return nil
	}

	var jobMap map[string]interface{}
	if err := unmarshal(&jobMap); err == nil {
		for key, value := range jobMap {
			w.Name = key

			if requires, ok := value.(map[string]interface{})["requires"]; ok {
				switch v := requires.(type) {
				case []interface{}:
					for _, req := range v {
						w.Requires = append(w.Requires, req.(string))
					}
				case string:
					w.Requires = append(w.Requires, v)
				default:
					return fmt.Errorf("unexpected type for requires field")
				}
			}
		}
		return nil
	}

	return fmt.Errorf("failed to unmarshal job")
}

func decodeMapToStruct(data map[string]interface{}, out interface{}) error {
	bytes, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(bytes, out)
}

type RestoreCacheStep struct {
	Name string   `yaml:"name"`
	Keys []string `yaml:"keys"`
}

type SaveCacheStep struct {
	Name  string   `yaml:"name"`
	Key   string   `yaml:"key"`
	Paths []string `yaml:"paths"`
}

type RunStep struct {
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
}

type Workflow struct {
	Jobs []JobWithRequires `yaml:"jobs"`
}

type JobWithRequires struct {
	Name     string   `yaml:"name"`
	Requires []string `yaml:"requires,omitempty"`
}
