package middleware

import (
	"net/http"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

// Interface guards
var (
	_ caddy.Module                = (*Delayer)(nil)
	_ caddy.Validator             = (*Delayer)(nil)
	_ caddy.Provisioner           = (*Delayer)(nil)
	_ caddyhttp.MiddlewareHandler = (*Delayer)(nil)
	_ caddyfile.Unmarshaler       = (*Delayer)(nil)
)

func init() {
	caddy.RegisterModule(Delayer{})
	httpcaddyfile.RegisterHandlerDirective("delayer", parseCaddyfile)
}

// Delayer is an example; put your own type here.
type Delayer struct {
	Duration string `json:"duration,omitempty"`
	duration time.Duration
	logger   *zap.Logger
	// post          func(ctx context.Context, url string, payload, unmarshal interface{}) error
	// InjectLoginBy bool
}

func (a *Delayer) Validate() error {
	a.logger.Debug("🟢 Validate")
	if d, err := time.ParseDuration(a.Duration); err != nil {
		return err
	} else {
		a.duration = d
	}
	return nil
}

func (a *Delayer) Provision(ctx caddy.Context) error {
	// a.post = newClient().post
	a.logger = ctx.Logger(a)
	a.logger.Debug("🟢 Provision", zap.Any("a", *a))
	return nil
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (a *Delayer) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		// fmt.Printf("%#v\n", d)
		if !d.Args(&a.Duration) {
			return d.ArgErr()
		}
	}
	return nil
}

// parseCaddyfile unmarshals tokens from h into a new Middleware.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var a Delayer
	err := a.UnmarshalCaddyfile(h.Dispenser)
	return a, err
}

// CaddyModule returns the Caddy module information.
func (Delayer) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.delayer",
		New: func() caddy.Module { return new(Delayer) },
	}
}

func (a Delayer) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	a.logger.Sugar().Infof("delaying for %s", a.duration)
	time.Sleep(a.duration)
	return next.ServeHTTP(w, r)
}
