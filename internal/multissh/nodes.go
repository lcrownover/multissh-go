package multissh

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
)

// NodeList: Given a comma-separated list of nodes, return a valid list of nodes
func getNodeList(ns string) []string {
	var nl []string

	logger("nodes").Debug(fmt.Sprintf("parsing node string: %s", ns))

	l := strings.Split(ns, ",")

	for _, node := range l {
        // we dont care about empty strings
        n := strings.TrimSpace(node)
        if n == "" {
            continue
        }

		validateNodeString(n)
		nl = append(nl, n)
	}

	logger("nodes").Debug(fmt.Sprintf("node list: %v", nl))
	return nl
}

func validateNodeString(ns string) {
	var err error
	validate := validator.New()
	logger("nodes").Debug(fmt.Sprintf("validating node: %s", ns))

	// first try fqdn
	err = validate.Var(ns, "required,fqdn")
	if err == nil {
		return
	}
	// next try hostname
	err = validate.Var(ns, "required,hostname")
	if err == nil {
		return
	}
	// finally try ipaddr
	err = validate.Var(ns, "required,ip4_addr")
	if err == nil {
		return
	}
	logger("nodes").Error(fmt.Sprintf("invalid node: %s", ns))
	os.Exit(1)
}
