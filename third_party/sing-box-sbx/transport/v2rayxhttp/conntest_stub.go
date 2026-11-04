package v2rayxhttp

import (
	"context"
	"net/http"
)

// xhttpTrace is a no-op on reF1nd (lx connection-test tracing is not present).
type xhttpTrace struct{}

func (c *Client) traceDial(context.Context) *xhttpTrace { return nil }

func (t *xhttpTrace) end(error) {}

func (t *xhttpTrace) response(*http.Response, error) {}

func (t *xhttpTrace) requestContext(parent context.Context) context.Context { return parent }

func (t *xhttpTrace) acquire(ctx context.Context, c *Client) (*xmuxClient, error) {
	return c.xmux.getContext(ctx)
}
