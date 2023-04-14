package tlsclientconfig

import (
	"crypto/tls"
	"net/http"
)

type OwnRoundTripper struct {
	inner   http.RoundTripper
	options Options
}

type Options struct {
	AddMissingHeaders bool
	Headers           map[string]string
}

func AddRoundTripper(inner http.RoundTripper, options ...Options) http.RoundTripper {
	if trans, ok := inner.(*http.Transport); ok {
		trans.TLSClientConfig = getTLSConfiguration()
	}
	roundTripper := &OwnRoundTripper{
		inner: inner,
	}
	if options != nil {
		roundTripper.options = options[0]
	} else {
		roundTripper.options = GetDefaultOptions()
	}
	return roundTripper
}

func (ug *OwnRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	if ug.options.AddMissingHeaders {
		for header, value := range ug.options.Headers {
			if _, ok := r.Header[header]; !ok {
				r.Header.Set(header, value)
			}
		}
	}
	if ug.inner == nil {
		return (&http.Transport{
			TLSClientConfig: getTLSConfiguration(),
		}).RoundTrip(r)
	}
	return ug.inner.RoundTrip(r)
}

func getTLSConfiguration() *tls.Config {
	return &tls.Config{
		CipherSuites: []uint16{tls.TLS_AES_128_GCM_SHA256, tls.TLS_CHACHA20_POLY1305_SHA256, tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256, tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256, tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
		},
		CurvePreferences: []tls.CurveID{tls.CurveP256, tls.CurveP384, tls.CurveP521, tls.X25519},
	}
}

func GetDefaultOptions() Options {
	return Options{
		AddMissingHeaders: true,
		Headers: map[string]string{
			"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
			"Accept-Language": "en-US,en;q=0.5",
			"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:104.0) Gecko/20100101 Firefox/104.0",
		},
	}
}
