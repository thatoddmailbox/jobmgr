package data

import (
	"time"

	"github.com/BurntSushi/toml"
)

const defaultTimeout time.Duration = 10 * time.Second

type JobSpec struct {
	Command          string
	Arguments        []string
	WorkingDirectory string
	PreserveEnvVars  []string
	Timeout          duration

	Parameter []JobSpecParameter
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

func ParseJobSpec(filename string) (JobSpec, error) {
	jobspec := JobSpec{}

	_, err := toml.DecodeFile(filename, &jobspec)
	if err != nil {
		return JobSpec{}, err
	}

	if jobspec.Timeout.Duration == time.Duration(0) {
		jobspec.Timeout.Duration = defaultTimeout
	}

	for _, p := range jobspec.Parameter {
		err = p.CheckValidType()
		if err != nil {
			return JobSpec{}, err
		}
	}

	return jobspec, nil
}
