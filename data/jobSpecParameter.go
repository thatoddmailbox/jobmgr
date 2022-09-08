package data

import (
	"fmt"
	"strconv"
)

type JobSpecParameter struct {
	Name string
	Type JobSpecParameterType
}

type JobSpecParameterType string

const (
	JobSpecParameterTypeString JobSpecParameterType = "string"
	JobSpecParameterTypeInt    JobSpecParameterType = "int"
)

func (p *JobSpecParameter) CheckValidType() error {
	if !(p.Type == JobSpecParameterTypeString || p.Type == JobSpecParameterTypeInt) {
		return fmt.Errorf("type '%s' of parameter '%s' is not valid", p.Type, p.Name)
	}

	return nil
}

func (p *JobSpecParameter) ParseValue(value string) (string, bool) {
	if p.Type == JobSpecParameterTypeInt {
		result, err := strconv.Atoi(value)
		if err != nil {
			return "", false
		}

		return strconv.Itoa(result), true
	}

	return value, true
}
