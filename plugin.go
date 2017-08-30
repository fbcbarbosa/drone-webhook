package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/drone/drone-go/drone"
)

type (

	// Plugin defines the Webhook plugin parameters.
	Plugin struct {
		URLs        []string          `json:"urls"`
		SkipVerify  bool              `json:"skip_verify"`
		Debug       bool              `json:"debug"`
		Auth        Auth              `json:"auth"`
		Headers     map[string]string `json:"header"`
		Method      string            `json:"method"`
		Template    string            `json:"template"`
		ContentType string            `json:"content_type"`
		Payload     Payload           `json:"payload"`
	}

	// Auth represents a basic HTTP authentication username and password.
	Auth struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Payload defines the full payload set by the Webhook.
	Payload struct {
		Repo  drone.Repo  `json:"repo"`
		Build drone.Build `json:"build"`
	}
)

const (
	respFormat      = "Webhook %d\n  URL: %s\n  RESPONSE STATUS: %s\n  RESPONSE BODY: %s\n"
	debugRespFormat = "Webhook %d\n  URL: %s\n  METHOD: %s\n  HEADERS: %s\n  REQUEST BODY: %s\n  RESPONSE STATUS: %s\n  RESPONSE BODY: %s\n"
)

var (
	buildCommit string
)

// Exec runs the plugin
func (p *Plugin) Exec() error {
	if p.Debug {
		log.SetLevel(log.DebugLevel)
	}

	log.WithFields(log.Fields{
		"auth":    p.Auth,
		"headers": p.Headers,
	}).Debug()

	if p.Method == "" {
		p.Method = "POST"
	}

	if p.ContentType == "" {
		p.ContentType = "application/json"
	}

	var b []byte
	if p.Template == "" {
		buf, err := json.Marshal(&p.Payload)
		if err != nil {
			return fmt.Errorf("failed to encode JSON payload. %s", err)
		}
		b = buf
	} else {
		msg, err := RenderTrim(p.Template, &p.Payload)
		if err != nil {
			return fmt.Errorf("failed to execute the content template. %s", err)
		}
		b = []byte(msg)
	}

	// build and execute a request for each url.
	// all auth, headers, method, template (payload),
	// and content_type values will be applied to
	// every webhook request.

	for i, rawurl := range p.URLs {
		uri, err := url.Parse(rawurl)

		if err != nil {
			return fmt.Errorf("error: Failed to parse the hook URL. %s", err)
		}

		r := bytes.NewReader(b)

		req, err := http.NewRequest(p.Method, uri.String(), r)

		if err != nil {
			return fmt.Errorf("error: Failed to create the HTTP request. %s", err)
		}

		req.Header.Set("Content-Type", p.ContentType)

		for key, value := range p.Headers {
			req.Header.Set(key, value)
		}

		if p.Auth.Username != "" {
			req.SetBasicAuth(p.Auth.Username, p.Auth.Password)
		}

		client := http.DefaultClient
		if p.SkipVerify {
			client = &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
			}
		}
		resp, err := client.Do(req)

		if err != nil {
			return fmt.Errorf("failed to execute the HTTP request. %s", err)
		}

		defer resp.Body.Close()

		if p.Debug || resp.StatusCode >= http.StatusBadRequest {
			body, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				log.Infof("failed to read the HTTP response body. %s", err)
			}

			if p.Debug {
				log.Infof(
					debugRespFormat,
					i+1,
					req.URL,
					req.Method,
					req.Header,
					string(b),
					resp.Status,
					string(body),
				)
			} else {
				log.Infof(
					respFormat,
					i+1,
					req.URL,
					string(body),
				)
			}
		}
	}
	return nil
}
