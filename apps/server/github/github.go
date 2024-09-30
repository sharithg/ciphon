// Copyright 2018 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package github

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gregjones/httpcache"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/rcrowley/go-metrics"
	"github.com/sharithg/siphon/config"
)

type Github struct {
	Handler http.Handler
}

func New() Github {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	configPath := filepath.Join(currentDir, "gh-config.yaml")
	config, err := config.ReadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatalln("error reading gh oauth config: %w", err)
	}
	metricsRegistry := metrics.DefaultRegistry

	cc, err := githubapp.NewDefaultCachingClientCreator(
		config.Github,
		githubapp.WithClientUserAgent("siphon-app/1.0.0"),
		githubapp.WithClientTimeout(3*time.Second),
		githubapp.WithClientCaching(false, func() httpcache.Cache { return httpcache.NewMemoryCache() }),
		githubapp.WithClientMiddleware(
			githubapp.ClientMetrics(metricsRegistry),
		),
	)

	if err != nil {
		log.Fatalln("error creating gh wh client: %w", err)
	}

	prCommentHandler := &CommitHandler{
		ClientCreator: cc,
		Preamble:      config.AppConfig.PullRequestPreamble,
	}

	webhookHandler := githubapp.NewDefaultEventDispatcher(config.Github, prCommentHandler)

	return Github{
		Handler: webhookHandler,
	}
}
