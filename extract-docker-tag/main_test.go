// Copyright 2021 FerretDB Inc.
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

package main

import (
	"testing"

	"github.com/sethvargo/go-githubactions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/FerretDB/github-actions/internal/testutil"
)

func TestExtract(t *testing.T) {
	t.Run("pull_request", func(t *testing.T) {
		getEnv := testutil.GetEnvFunc(t, map[string]string{
			"GITHUB_BASE_REF":   "main",
			"GITHUB_EVENT_NAME": "pull_request",
			"GITHUB_HEAD_REF":   "extract-docker-tag",
			"GITHUB_REF_NAME":   "1/merge",
			"GITHUB_REF_TYPE":   "branch",
			"GITHUB_REPOSITORY": "FerretDB/FerretDB",
		})

		action := githubactions.New(githubactions.WithGetenv(getEnv))
		result, err := extract(action)
		require.NoError(t, err)
		assert.Equal(t, "ferretdb", result.owner)
		assert.Equal(t, "ferretdb-dev", result.name)
		assert.Equal(t, []string{"pr-extract-docker-tag"}, result.tags)
		assert.Equal(t, "ghcr.io/ferretdb/ferretdb-dev:pr-extract-docker-tag", result.ghcr)
		assert.Equal(t, []string{"ghcr.io/ferretdb/ferretdb-dev:pr-extract-docker-tag"}, result.ghcrImages)
	})

	t.Run("pull_request_target", func(t *testing.T) {
		getEnv := testutil.GetEnvFunc(t, map[string]string{
			"GITHUB_BASE_REF":   "main",
			"GITHUB_EVENT_NAME": "pull_request_target",
			"GITHUB_HEAD_REF":   "extract-docker-tag",
			"GITHUB_REF_NAME":   "main",
			"GITHUB_REF_TYPE":   "branch",
			"GITHUB_REPOSITORY": "FerretDB/FerretDB",
		})

		action := githubactions.New(githubactions.WithGetenv(getEnv))
		result, err := extract(action)
		require.NoError(t, err)
		assert.Equal(t, "ferretdb", result.owner)
		assert.Equal(t, "ferretdb-dev", result.name)
		assert.Equal(t, []string{"pr-extract-docker-tag"}, result.tags)
		assert.Equal(t, "ghcr.io/ferretdb/ferretdb-dev:pr-extract-docker-tag", result.ghcr)
		assert.Equal(t, []string{"ghcr.io/ferretdb/ferretdb-dev:pr-extract-docker-tag"}, result.ghcrImages)
	})

	t.Run("pull_request/dependabot", func(t *testing.T) {
		getEnv := testutil.GetEnvFunc(t, map[string]string{
			"GITHUB_BASE_REF":   "main",
			"GITHUB_EVENT_NAME": "pull_request",
			"GITHUB_HEAD_REF":   "dependabot/submodules/tests/mongo-go-driver-29d768e",
			"GITHUB_REF_NAME":   "58/merge",
			"GITHUB_REF_TYPE":   "branch",
			"GITHUB_REPOSITORY": "FerretDB/FerretDB",
		})

		action := githubactions.New(githubactions.WithGetenv(getEnv))
		result, err := extract(action)
		require.NoError(t, err)
		assert.Equal(t, "ferretdb", result.owner)
		assert.Equal(t, "ferretdb-dev", result.name)
		assert.Equal(t, []string{"pr-mongo-go-driver-29d768e"}, result.tags)
		assert.Equal(t, "ghcr.io/ferretdb/ferretdb-dev:pr-mongo-go-driver-29d768e", result.ghcr)
		assert.Equal(t, []string{"ghcr.io/ferretdb/ferretdb-dev:pr-mongo-go-driver-29d768e"}, result.ghcrImages)
	})

	t.Run("push/main", func(t *testing.T) {
		getEnv := testutil.GetEnvFunc(t, map[string]string{
			"GITHUB_BASE_REF":   "",
			"GITHUB_EVENT_NAME": "push",
			"GITHUB_HEAD_REF":   "",
			"GITHUB_REF_NAME":   "main",
			"GITHUB_REF_TYPE":   "branch",
			"GITHUB_REPOSITORY": "FerretDB/FerretDB",
		})

		action := githubactions.New(githubactions.WithGetenv(getEnv))
		result, err := extract(action)
		require.NoError(t, err)
		assert.Equal(t, "ferretdb", result.owner)
		assert.Equal(t, "ferretdb-dev", result.name)
		assert.Equal(t, []string{"main"}, result.tags)
		assert.Equal(t, "ghcr.io/ferretdb/ferretdb-dev:main", result.ghcr)
		assert.Equal(t, []string{"ghcr.io/ferretdb/ferretdb-dev:main"}, result.ghcrImages)
	})

	t.Run("push/tag/beta", func(t *testing.T) {
		getEnv := testutil.GetEnvFunc(t, map[string]string{
			"GITHUB_BASE_REF":   "",
			"GITHUB_EVENT_NAME": "push",
			"GITHUB_HEAD_REF":   "",
			"GITHUB_REF_NAME":   "v0.1.0-beta",
			"GITHUB_REF_TYPE":   "tag",
			"GITHUB_REPOSITORY": "FerretDB/FerretDB",
		})

		action := githubactions.New(githubactions.WithGetenv(getEnv))
		result, err := extract(action)
		require.NoError(t, err)
		assert.Equal(t, "ferretdb", result.owner)
		assert.Equal(t, "ferretdb-dev", result.name)
		assert.Equal(t, []string{"0.1.0-beta", "latest"}, result.tags)
		assert.Equal(t, "ghcr.io/ferretdb/ferretdb-dev:0.1.0-beta", result.ghcr)
		assert.Equal(t, []string{
			"ghcr.io/ferretdb/ferretdb-dev:latest",
			"ghcr.io/ferretdb/ferretdb-dev:0.1.0-beta",
		}, result.ghcrImages)
	})

	t.Run("push/tag/release", func(t *testing.T) {
		getEnv := testutil.GetEnvFunc(t, map[string]string{
			"GITHUB_BASE_REF":   "",
			"GITHUB_EVENT_NAME": "push",
			"GITHUB_HEAD_REF":   "",
			"GITHUB_REF_NAME":   "v0.1.0",
			"GITHUB_REF_TYPE":   "tag",
			"GITHUB_REPOSITORY": "FerretDB/FerretDB",
		})

		action := githubactions.New(githubactions.WithGetenv(getEnv))
		result, err := extract(action)
		require.NoError(t, err)
		assert.Equal(t, "ferretdb", result.owner)
		assert.Equal(t, "ferretdb-dev", result.name)
		assert.Equal(t, []string{"0.1.0", "latest"}, result.tags)
		assert.Equal(t, "ghcr.io/ferretdb/ferretdb-dev:0.1.0", result.ghcr)
		assert.Equal(t, []string{
			"ghcr.io/ferretdb/ferretdb-dev:latest",
			"ghcr.io/ferretdb/ferretdb-dev:0.1.0",
		}, result.ghcrImages)
	})

	t.Run("push/tag/wrong", func(t *testing.T) {
		getEnv := testutil.GetEnvFunc(t, map[string]string{
			"GITHUB_BASE_REF":   "",
			"GITHUB_EVENT_NAME": "push",
			"GITHUB_HEAD_REF":   "",
			"GITHUB_REF_NAME":   "0.1.0", // no leading v
			"GITHUB_REF_TYPE":   "tag",
			"GITHUB_REPOSITORY": "FerretDB/FerretDB",
		})

		action := githubactions.New(githubactions.WithGetenv(getEnv))
		_, err := extract(action)
		require.Error(t, err)
	})

	t.Run("schedule", func(t *testing.T) {
		getEnv := testutil.GetEnvFunc(t, map[string]string{
			"GITHUB_BASE_REF":   "",
			"GITHUB_EVENT_NAME": "schedule",
			"GITHUB_HEAD_REF":   "",
			"GITHUB_REF_NAME":   "main",
			"GITHUB_REF_TYPE":   "branch",
			"GITHUB_REPOSITORY": "FerretDB/FerretDB",
		})

		action := githubactions.New(githubactions.WithGetenv(getEnv))
		result, err := extract(action)
		require.NoError(t, err)
		assert.Equal(t, "ferretdb", result.owner)
		assert.Equal(t, "ferretdb-dev", result.name)
		assert.Equal(t, []string{"main"}, result.tags)
		assert.Equal(t, "ghcr.io/ferretdb/ferretdb-dev:main", result.ghcr)
		assert.Equal(t, []string{"ghcr.io/ferretdb/ferretdb-dev:main"}, result.ghcrImages)
	})

	t.Run("workflow_run", func(t *testing.T) {
		getEnv := testutil.GetEnvFunc(t, map[string]string{
			"GITHUB_BASE_REF":   "",
			"GITHUB_EVENT_NAME": "workflow_run",
			"GITHUB_HEAD_REF":   "",
			"GITHUB_REF_NAME":   "main",
			"GITHUB_REF_TYPE":   "branch",
			"GITHUB_REPOSITORY": "FerretDB/FerretDB",
		})

		action := githubactions.New(githubactions.WithGetenv(getEnv))
		result, err := extract(action)
		require.NoError(t, err)
		assert.Equal(t, "ferretdb", result.owner)
		assert.Equal(t, "ferretdb-dev", result.name)
		assert.Equal(t, []string{"main"}, result.tags)
		assert.Equal(t, "ghcr.io/ferretdb/ferretdb-dev:main", result.ghcr)
		assert.Equal(t, []string{"ghcr.io/ferretdb/ferretdb-dev:main"}, result.ghcrImages)
	})
}
