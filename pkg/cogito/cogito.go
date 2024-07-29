package cogito

import (
	"context"
	"net/http"
	"time"
)

type Cogito struct {
	URL        string
	APIKey     string
	HTTPClient *http.Client
	Logger     func(ctx context.Context, caller string, request *http.Request, response *http.Response, elapse time.Duration)
}

func (m *Cogito) BaseUrl() string      { return m.URL }
func (m *Cogito) Key() string          { return m.APIKey }
func (m *Cogito) Client() *http.Client { return m.HTTPClient }
func (m *Cogito) Log(ctx context.Context, caller string, request *http.Request, response *http.Response, elapse time.Duration) {
	m.Logger(ctx, caller, request, response, elapse)
}
