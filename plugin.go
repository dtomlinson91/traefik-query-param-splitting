package traefik_query_param_splitting_middleware

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type Config struct {
	Delimiter  string `yaml:"delimiter"`
	ParamRegex string `yaml:"paramRegex"`
}

type Response struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func CreateConfig() *Config {
	return &Config{
		Delimiter:  "|",
		ParamRegex: ".*",
	}
}

type QueryParam struct {
	next       http.Handler
	delimiter  string
	paramRegex string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	// check only one delimiter is given
	if len(config.Delimiter) != 1 {
		return nil, fmt.Errorf("only one delimiter character can be specified")
	}

	return &QueryParam{
		next:       next,
		delimiter:  config.Delimiter,
		paramRegex: config.ParamRegex,
	}, nil
}

func (qp *QueryParam) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	u := req.URL.Query()
	l := make([]string, 0)

	for qryp, qryv := range u {
		// parse paramRegex as regex
		re, err := regexp.Compile(qp.paramRegex)
		if err != nil {
			http.Error(rw, fmt.Sprintf("invalid regex: %v", err), http.StatusInternalServerError)
			return
		}

		// check if query param matches the regex
		if re.MatchString(qryp) {
			// for each value
			for _, qryvRaw := range qryv {
				// split by delimiter
				s := strings.Split(qryvRaw, qp.delimiter)
				if len(s) > 1 {
					// if delimiter found, clear query param and set individual value
					u.Del(qryp)
					u[qryp] = append(l, s...)
				}
			}
		}
	}
	req.URL.RawQuery = u.Encode()
	req.RequestURI = req.URL.RequestURI()
	qp.next.ServeHTTP(rw, req)
}
