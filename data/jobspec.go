package data

import "github.com/BurntSushi/toml"

type JobSpec struct {
	Command          string
	Arguments        []string
	WorkingDirectory string
}

func ParseJobSpec(filename string) (JobSpec, error) {
	jobspec := JobSpec{}

	_, err := toml.DecodeFile(filename, &jobspec)
	if err != nil {
		return JobSpec{}, err
	}

	return jobspec, nil
}
