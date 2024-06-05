package serviceregister

import (
	"github.com/mick-roper/rdfox-cli/rdfox"
	v7 "github.com/mick-roper/rdfox-cli/rdfox/v7"
)

func RetrieveCreateConnectionFunc(_ int) (rdfox.CreateConnection, error) {
	return v7.CreateConnection, nil
}

func RetrieveDeleteConnectionFunc(_ int) (rdfox.DeleteConnection, error) {
	return v7.DeleteConnection, nil
}
