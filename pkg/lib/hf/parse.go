// Copyright 2024 The KitOps Authors.
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
//
// SPDX-License-Identifier: Apache-2.0

package hf

import (
	"fmt"
	"net/url"
	"strings"
)

// ParseHuggingFaceRepo parses a HuggingFace repository path or URL and returns
// the normalized repository name and type.
//
// Supported formats:
//   - Full URL: https://huggingface.co/org/repo
//   - Full URL (dataset): https://huggingface.co/datasets/org/repo
//   - Short form: org/repo (defaults to model)
//   - Dataset short form: datasets/org/repo
func ParseHuggingFaceRepo(path string) (repo string, repoType RepositoryType, err error) {
	// Handle "datasets/" prefix (short form)
	if strings.HasPrefix(path, "datasets/") {
		repo = strings.TrimPrefix(path, "datasets/")
		if !isValidRepoFormat(repo) {
			return "", RepoTypeUnknown, fmt.Errorf("invalid dataset repository format: %s", path)
		}
		return repo, RepoTypeDataset, nil
	}

	// Check if this looks like a URL (contains :// or starts with domain)
	if strings.Contains(path, "://") || strings.HasPrefix(path, "huggingface.co/") {
		// Normalize URL by adding scheme if missing
		urlStr := path
		if !strings.Contains(path, "://") {
			urlStr = "https://" + path
		}

		parsedURL, err := url.Parse(urlStr)
		if err != nil {
			return "", RepoTypeUnknown, fmt.Errorf("invalid URL format: %w", err)
		}

		// Security: validate hostname is exactly huggingface.co
		if parsedURL.Hostname() != "huggingface.co" {
			return "", RepoTypeUnknown, fmt.Errorf("unsupported hostname: %s (only huggingface.co is supported)", parsedURL.Hostname())
		}

		// Parse path segments
		pathSegments := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")

		// Check for dataset URLs
		if len(pathSegments) >= 1 && pathSegments[0] == "datasets" {
			// Dataset URLs must have format: datasets/org/repo
			if len(pathSegments) >= 3 {
				repo = strings.Join(pathSegments[1:3], "/")
				return repo, RepoTypeDataset, nil
			}
			// Invalid dataset URL (not enough segments)
			return "", RepoTypeUnknown, fmt.Errorf("invalid dataset URL: expected datasets/org/repo, got '%s'", parsedURL.Path)
		}

		// Model URLs: should be exactly 2 segments (org/repo)
		if len(pathSegments) == 2 {
			repo = strings.Join(pathSegments, "/")
			return repo, RepoTypeModel, nil
		}

		return "", RepoTypeUnknown, fmt.Errorf("unrecognized HuggingFace URL pattern: %s", path)
	}

	// Handle short form (org/repo)
	if isValidRepoFormat(path) {
		return path, RepoTypeModel, nil
	}

	return "", RepoTypeUnknown, fmt.Errorf("invalid repository format: %s", path)
}

// isValidRepoFormat checks if a string is in "org/repo" format
func isValidRepoFormat(s string) bool {
	parts := strings.Split(s, "/")
	return len(parts) == 2 && parts[0] != "" && parts[1] != ""
}
