package traefik_query_param_splitting_middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type Config struct {
	Delimiter string
}

type Response struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func CreateConfig() *Config {
	return &Config{
		Delimiter: "|",
	}
}

type QueryParam struct {
	next      http.Handler
	delimiter string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	// check only one delimiter is given
	if len(config.Delimiter) != 1 {
		return nil, fmt.Errorf("only one delimiter character can be specified")
	}

	return &QueryParam{
		next:      next,
		delimiter: config.Delimiter,
	}, nil
}

func (qp *QueryParam) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	u := req.URL.Query()
	l := make([]string, 0)

	for qp, qv := range u {
		// for each value
		for _, qva := range qv {
			// split by delimiter
			s := strings.Split(qryvRaw, qp.delimiter)
			if len(s) > 1 {
				// if delimiter found, clear query param and set individual value
				u.Del(qp)
				u[qp] = append(l, s...)
			}
		}
	}
	req.URL.RawQuery = u.Encode()
	fmt.Println(req.URL)
	fmt.Println(req.URL.Query())
	fmt.Println(req.URL.Query().Get("filter[bedrooms]"))
	fmt.Println(req.URL.Query().Get("fields[property]"))
	qp.next.ServeHTTP(rw, req)
}
