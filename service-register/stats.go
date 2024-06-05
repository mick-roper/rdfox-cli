package serviceregister

import (
	"fmt"

	"github.com/mick-roper/rdfox-cli/rdfox"
	v6 "github.com/mick-roper/rdfox-cli/rdfox/v6"
	v7 "github.com/mick-roper/rdfox-cli/rdfox/v7"
)

func RetrieveGetStatsFunc(version int) (rdfox.GetStats, error) {
	switch version {
	case 6:
		return v6.GetStats, nil
	case 7:
		return v7.GetStats, nil
	default:
		return nil, fmt.Errorf("unsupported version %d", version)
	}
}
