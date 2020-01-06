package health

import (
	"fmt"
	"strings"
)

// CheckFun is a that is called in order to evaluate the check
// It is called to obtain the health state of the Check.
// It should return nil if the check is healthy.
// In case the check is not healthy the according error should be returned
type CheckFun func() error

// NewSimpleCheck creates a Check based on a given CheckFun and a name.
// Instead of implementing the Check interface you can use this way instead.
func NewSimpleCheck(name string, check CheckFun) (Check, error) {

	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return nil, fmt.Errorf("Can't create a Check with an empty name")
	}

	if check == nil {
		return nil, fmt.Errorf("Can't create a Check whose CheckFun is nil")
	}

	return simpleCheck{
		name:  name,
		check: check,
	}, nil
}

type simpleCheck struct {
	name  string
	check CheckFun
}

func (s simpleCheck) IsHealthy() error {
	return s.check()
}

func (s simpleCheck) Name() string {
	return s.name
}
