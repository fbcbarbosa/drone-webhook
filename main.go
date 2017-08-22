package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/template"
	"github.com/urfave/cli"
)

const (
	respFormat      = "Webhook %d\n  URL: %s\n  RESPONSE STATUS: %s\n  RESPONSE BODY: %s\n"
	debugRespFormat = "Webhook %d\n  URL: %s\n  METHOD: %s\n  HEADERS: %s\n  REQUEST BODY: %s\n  RESPONSE STATUS: %s\n  RESPONSE BODY: %s\n"
)

var (
	buildCommit string
)

func main() {
	log.Infof("Drone Webhook Plugin built from %s\n", buildCommit)

	app := cli.NewApp()
	app.Name = "webhook plugin"
	app.Usage = "webhook plugin"
	app.Action = run
	app.Version = fmt.Sprint(buildCommit)
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "webhook debug",
			EnvVar: "PLUGIN_DEBUG",
		},
		cli.BoolFlag{
			Name:   "skip_verify",
			Usage:  "webhook skip verify",
			EnvVar: "PLUGIN_SKIP_VERIFY",
		},
		cli.StringFlag{
			Name:   "method",
			Usage:  "webhook method",
			EnvVar: "PLUGIN_METHOD",
		},
		cli.StringSliceFlag{
			Name:   "urls",
			Usage:  "webhook urls",
			EnvVar: "PLUGIN_URLS",
		},
		cli.StringFlag{
			Name:   "content_type",
			Usage:  "webhook content type",
			EnvVar: "PLUGIN_CONTENT_TYPE",
		},
		cli.StringFlag{
			Name:   "headers",
			Usage:  "webhook content type",
			EnvVar: "PLUGIN_HEADERS",
		},
		cli.StringFlag{
			Name:   "auth",
			Usage:  "webhook auth",
			EnvVar: "PLUGIN_AUTH",
		},
		cli.StringFlag{
			Name:   "template",
			Usage:  "webhook template",
			EnvVar: "PLUGIN_TEMPLATE",
		},
		cli.StringFlag{
			Name:   "repo.owner",
			Usage:  "repo owner",
			EnvVar: "DRONE_REPO_OWNER",
		},
		cli.StringFlag{
			Name:   "repo.name",
			Usage:  "repo name",
			EnvVar: "DRONE_REPO_NAME",
		},
		cli.StringFlag{
			Name:   "repo.link",
			Usage:  "repo link",
			EnvVar: "DRONE_REPO_LINK",
		},
		cli.StringFlag{
			Name:   "repo.avatar",
			Usage:  "repo avatar",
			EnvVar: "DRONE_REPO_AVATAR",
		},
		cli.StringFlag{
			Name:   "repo.branch",
			Usage:  "repo branch",
			EnvVar: "DRONE_REPO_BRANCH",
		},
		cli.StringFlag{
			Name:   "repo.clone",
			Usage:  "repo clone",
			EnvVar: "DRONE_REMOTE_URL",
		},
		cli.StringFlag{
			Name:   "commit.sha",
			Usage:  "commit sha",
			EnvVar: "DRONE_COMMIT_SHA",
		},
		cli.StringFlag{
			Name:  "commit.ref",
			Usage: "commit ref",
		},
		cli.StringFlag{
			Name:   "commit.branch",
			Usage:  "commit branch",
			EnvVar: "DRONE_COMMIT_REF",
		},
		cli.StringFlag{
			Name:   "commit.link",
			Usage:  "commit link",
			EnvVar: "DRONE_COMMIT_LINK",
		},
		cli.StringFlag{
			Name:   "commit.message",
			Usage:  "commit message",
			EnvVar: "DRONE_COMMIT_MESSAGE",
		},
		cli.StringFlag{
			Name:   "commit.author.name",
			Usage:  "commit author name",
			EnvVar: "DRONE_COMMIT_AUTHOR",
		},
		cli.StringFlag{
			Name:   "commit.author.email",
			Usage:  "commit author email",
			EnvVar: "DRONE_COMMIT_AUTHOR_EMAIL",
		},
		cli.StringFlag{
			Name:   "commit.author.avatar",
			Usage:  "commit author avatar",
			EnvVar: "DRONE_COMMIT_AUTHOR_AVATAR",
		},
		cli.IntFlag{
			Name:   "build.number",
			Usage:  "build number",
			EnvVar: "DRONE_BUILD_NUMBER",
		},
		cli.StringFlag{
			Name:   "build.event",
			Usage:  "build event",
			EnvVar: "DRONE_BUILD_EVENT",
		},
		cli.StringFlag{
			Name:   "build.status",
			Usage:  "build status",
			EnvVar: "DRONE_BUILD_STATUS",
		},
		cli.StringFlag{
			Name:   "build.link",
			Usage:  "build link",
			EnvVar: "DRONE_BUILD_LINK",
		},
		cli.Int64Flag{
			Name:   "build.created",
			Usage:  "build created",
			EnvVar: "DRONE_BUILD_CREATED",
		},
		cli.Int64Flag{
			Name:   "build.started",
			Usage:  "build started",
			EnvVar: "DRONE_BUILD_STARTED",
		},
		cli.Int64Flag{
			Name:   "build.finished",
			Usage:  "build finished",
			EnvVar: "DRONE_BUILD_FINISHED",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	log.Info(c.Generic("headers"))

	payload := &drone.Payload{
		Repo: &drone.Repo{
			Owner:  c.String("repo.owner"),
			Name:   c.String("repo.name"),
			Link:   c.String("repo.link"),
			Avatar: c.String("repo.avatar"),
			Branch: c.String("repo.branch"),
			Clone:  c.String("repo.clone"),
		},
		Build: &drone.Build{
			Number:   c.Int("build.number"),
			Event:    c.String("build.event"),
			Status:   c.String("build.status"),
			Link:     c.String("build.link"),
			Created:  c.Int64("build.created"),
			Started:  c.Int64("build.started"),
			Finished: c.Int64("build.finished"),
			Commit:   c.String("commit.sha"),
			Ref:      c.String("commit.ref"),
			Branch:   c.String("commit.branch"),
			Message:  c.String("commit.message"),
			Author:   c.String("commit.author.name"),
			Avatar:   c.String("commit.author.avatar"),
			Email:    c.String("commit.author.email"),
		},
	}

	vargs := Params{
		URLs:        c.StringSlice("urls"),
		SkipVerify:  c.Bool("skip_verify"),
		Debug:       c.Bool("debug"),
		Method:      c.String("method"),
		Template:    c.String("template"),
		ContentType: c.String("content_type"),
	}

	if err := json.Unmarshal([]byte(c.String("auth")), &vargs.Auth); err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(c.String("headers")), &vargs.Headers); err != nil {
		return err
	}

	if vargs.Method == "" {
		vargs.Method = "POST"
	}

	if vargs.ContentType == "" {
		vargs.ContentType = "application/json"
	}

	var b []byte
	if vargs.Template == "" {
		buf, err := json.Marshal(&payload)
		if err != nil {
			return fmt.Errorf("failed to encode JSON payload. %s", err)
		}
		b = buf
	} else {
		msg, err := template.RenderTrim(vargs.Template, &payload)
		if err != nil {
			return fmt.Errorf("failed to execute the content template. %s", err)
		}
		b = []byte(msg)
	}

	if vargs.Debug {
		log.SetLevel(log.DebugLevel)
	}

	// build and execute a request for each url.
	// all auth, headers, method, template (payload),
	// and content_type values will be applied to
	// every webhook request.

	log.WithFields(log.Fields{
		"auth":    vargs.Auth,
		"headers": vargs.Headers,
	}).Debug()

	for i, rawurl := range vargs.URLs {
		uri, err := url.Parse(rawurl)

		if err != nil {
			return fmt.Errorf("error: Failed to parse the hook URL. %s", err)
		}

		r := bytes.NewReader(b)

		req, err := http.NewRequest(vargs.Method, uri.String(), r)

		if err != nil {
			return fmt.Errorf("error: Failed to create the HTTP request. %s", err)
		}

		req.Header.Set("Content-Type", vargs.ContentType)

		for key, value := range vargs.Headers {
			req.Header.Set(key, value)
		}

		if vargs.Auth.Username != "" {
			req.SetBasicAuth(vargs.Auth.Username, vargs.Auth.Password)
		}

		client := http.DefaultClient
		if vargs.SkipVerify {
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

		if vargs.Debug || resp.StatusCode >= http.StatusBadRequest {
			body, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				log.Infof("failed to read the HTTP response body. %s", err)
			}

			if vargs.Debug {
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
					resp.Status,
					string(body),
				)
			}
		}
	}
	return nil
}
