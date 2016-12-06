package path

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Matcher is a function that matches given path and returns map[string]string
// with path parameters
type Matcher func(string) map[string]string

// NewMatcher creates a new matcher for given pattern.
// Pattern has the following form: /fix_text/:parameter_name1/:parameter_name2
func NewMatcher(pattern string) (Matcher, error) {

	if len(pattern) == 0 || pattern[0] != '/' {
		return nil, errors.New("Pattern must start with '/'")
	}

	parts := strings.Split(pattern, "/")

	pat := ""

	for i, part := range parts {
		if i > 0 {
			pat += "/"
		}
		if len(part) > 0 {
			if part[0] == ':' {
				ptn := fmt.Sprintf("(?P<%s>[^/]*)", regexp.QuoteMeta(part[1:]))
				pat += ptn
				continue
			}
			pat += regexp.QuoteMeta(part)
		}

	}

	exp, err := regexp.Compile("^" + pat + "$")
	if err != nil {
		return nil, err
	}
	return func(path string) map[string]string {
		m := exp.FindStringSubmatch(path)
		n1 := exp.SubexpNames()

		if m == nil {
			return nil
		}

		result := map[string]string{}

		for i, v := range m {
			if i > 0 {
				result[n1[i]] = v
			}
		}
		return result
	}, nil

}
