package main

import (
	"encoding/json"
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/drone/drone-go/drone"
	"github.com/urfave/cli"
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
			Name:   "repo",
			Usage:  "repo",
			EnvVar: "DRONE_REPO",
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

	payload := Payload{
		Repo: drone.Repo{
			Owner:  c.String("repo.owner"),
			Name:   c.String("repo.name"),
			Link:   c.String("repo.link"),
			Avatar: c.String("repo.avatar"),
			Branch: c.String("repo.branch"),
			Clone:  c.String("repo.clone"),
		},
		Build: drone.Build{
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

	plugin := Plugin{
		URLs:        c.StringSlice("urls"),
		SkipVerify:  c.Bool("skip_verify"),
		Debug:       c.Bool("debug"),
		Method:      c.String("method"),
		Template:    c.String("template"),
		ContentType: c.String("content_type"),
		Payload:     payload,
	}

	if c.String("auth") != "" {
		if err := json.Unmarshal([]byte(c.String("auth")), &plugin.Auth); err != nil {
			return err
		}
	}

	if c.String("headers") != "" {
		if err := json.Unmarshal([]byte(c.String("headers")), &plugin.Headers); err != nil {
			return err
		}
	}

	return plugin.Exec()
}
