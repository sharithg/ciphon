package parser

type Config struct {
	Version   string              `yaml:"version"`
	Jobs      map[string]Job      `yaml:"jobs"`
	Workflows map[string]Workflow `yaml:"workflows"`
}

type Job struct {
	Docker string `yaml:"docker"`
	Node   string `yaml:"node"`
	Steps  []Step `yaml:"steps"`
}

type Step struct {
	Name         string            `yaml:"name,omitempty"`
	Command      string            `yaml:"command,omitempty"`
	Checkout     bool              `yaml:"checkout,omitempty"`
	RestoreCache *RestoreCacheStep `yaml:"restore_cache,omitempty"`
	SaveCache    *SaveCacheStep    `yaml:"save_cache,omitempty"`
	Run          *RunStep          `yaml:"run,omitempty"`
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
	Jobs []string `yaml:"jobs"`
}
