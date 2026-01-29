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
	"testing"
)

func TestParseHuggingFaceRepo(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantRepo     string
		wantRepoType RepositoryType
		wantErr      bool
	}{
		// Valid model URLs
		{
			name:         "full HTTPS model URL",
			input:        "https://huggingface.co/myorg/mymodel",
			wantRepo:     "myorg/mymodel",
			wantRepoType: RepoTypeModel,
			wantErr:      false,
		},
		{
			name:         "full HTTP model URL",
			input:        "http://huggingface.co/myorg/mymodel",
			wantRepo:     "myorg/mymodel",
			wantRepoType: RepoTypeModel,
			wantErr:      false,
		},
		{
			name:         "scheme-less model URL",
			input:        "huggingface.co/myorg/mymodel",
			wantRepo:     "myorg/mymodel",
			wantRepoType: RepoTypeModel,
			wantErr:      false,
		},
		{
			name:         "short form model",
			input:        "myorg/mymodel",
			wantRepo:     "myorg/mymodel",
			wantRepoType: RepoTypeModel,
			wantErr:      false,
		},

		// Valid dataset URLs
		{
			name:         "full HTTPS dataset URL",
			input:        "https://huggingface.co/datasets/myorg/mydataset",
			wantRepo:     "myorg/mydataset",
			wantRepoType: RepoTypeDataset,
			wantErr:      false,
		},
		{
			name:         "full HTTP dataset URL",
			input:        "http://huggingface.co/datasets/myorg/mydataset",
			wantRepo:     "myorg/mydataset",
			wantRepoType: RepoTypeDataset,
			wantErr:      false,
		},
		{
			name:         "scheme-less dataset URL",
			input:        "huggingface.co/datasets/myorg/mydataset",
			wantRepo:     "myorg/mydataset",
			wantRepoType: RepoTypeDataset,
			wantErr:      false,
		},
		{
			name:         "short form dataset",
			input:        "datasets/myorg/mydataset",
			wantRepo:     "myorg/mydataset",
			wantRepoType: RepoTypeDataset,
			wantErr:      false,
		},

		// URLs with trailing slashes
		{
			name:         "model URL with trailing slash",
			input:        "https://huggingface.co/myorg/mymodel/",
			wantRepo:     "myorg/mymodel",
			wantRepoType: RepoTypeModel,
			wantErr:      false,
		},
		{
			name:         "dataset URL with trailing slash",
			input:        "https://huggingface.co/datasets/myorg/mydataset/",
			wantRepo:     "myorg/mydataset",
			wantRepoType: RepoTypeDataset,
			wantErr:      false,
		},

		// Security: Malicious URLs that should be rejected
		{
			name:         "subdomain attack",
			input:        "https://huggingface.co.evil.com/myorg/mymodel",
			wantRepo:     "",
			wantRepoType: RepoTypeUnknown,
			wantErr:      true,
		},
		{
			name:         "prefix attack",
			input:        "https://evilhuggingface.co/myorg/mymodel",
			wantRepo:     "",
			wantRepoType: RepoTypeUnknown,
			wantErr:      true,
		},
		{
			name:         "suffix attack",
			input:        "https://huggingface.co.attacker.com/myorg/mymodel",
			wantRepo:     "",
			wantRepoType: RepoTypeUnknown,
			wantErr:      true,
		},

		// Non-HuggingFace URLs
		{
			name:         "GitHub URL",
			input:        "https://github.com/myorg/myrepo",
			wantRepo:     "",
			wantRepoType: RepoTypeUnknown,
			wantErr:      true,
		},
		{
			name:         "GitLab URL",
			input:        "https://gitlab.com/myorg/myrepo",
			wantRepo:     "",
			wantRepoType: RepoTypeUnknown,
			wantErr:      true,
		},
		{
			name:         "random domain",
			input:        "https://example.com/myorg/myrepo",
			wantRepo:     "",
			wantRepoType: RepoTypeUnknown,
			wantErr:      true,
		},

		// Invalid formats
		{
			name:         "single segment",
			input:        "myrepo",
			wantRepo:     "",
			wantRepoType: RepoTypeUnknown,
			wantErr:      true,
		},
		{
			name:         "too many segments",
			input:        "org/repo/extra",
			wantRepo:     "",
			wantRepoType: RepoTypeUnknown,
			wantErr:      true,
		},
		{
			name:         "empty string",
			input:        "",
			wantRepo:     "",
			wantRepoType: RepoTypeUnknown,
			wantErr:      true,
		},
		{
			name:         "only slashes",
			input:        "///",
			wantRepo:     "",
			wantRepoType: RepoTypeUnknown,
			wantErr:      true,
		},
		{
			name:         "incomplete dataset short form",
			input:        "datasets/myorg",
			wantRepo:     "",
			wantRepoType: RepoTypeUnknown,
			wantErr:      true,
		},
		{
			name:         "incomplete dataset URL",
			input:        "https://huggingface.co/datasets/myorg",
			wantRepo:     "",
			wantRepoType: RepoTypeUnknown,
			wantErr:      true,
		},

		// Edge cases with special characters
		{
			name:         "repo with hyphens",
			input:        "my-org/my-model",
			wantRepo:     "my-org/my-model",
			wantRepoType: RepoTypeModel,
			wantErr:      false,
		},
		{
			name:         "repo with underscores",
			input:        "my_org/my_model",
			wantRepo:     "my_org/my_model",
			wantRepoType: RepoTypeModel,
			wantErr:      false,
		},
		{
			name:         "repo with numbers",
			input:        "org123/model456",
			wantRepo:     "org123/model456",
			wantRepoType: RepoTypeModel,
			wantErr:      false,
		},
		{
			name:         "repo with dots",
			input:        "my.org/my.model",
			wantRepo:     "my.org/my.model",
			wantRepoType: RepoTypeModel,
			wantErr:      false,
		},

		// Real-world examples
		{
			name:         "GPT-2 model",
			input:        "https://huggingface.co/gpt2",
			wantRepo:     "",
			wantRepoType: RepoTypeUnknown,
			wantErr:      true, // Single segment is invalid
		},
		{
			name:         "BERT base model",
			input:        "https://huggingface.co/bert-base-uncased",
			wantRepo:     "",
			wantRepoType: RepoTypeUnknown,
			wantErr:      true, // Single segment is invalid
		},
		{
			name:         "OpenAI GPT-2",
			input:        "openai/gpt2",
			wantRepo:     "openai/gpt2",
			wantRepoType: RepoTypeModel,
			wantErr:      false,
		},
		{
			name:         "SQuAD dataset",
			input:        "datasets/squad",
			wantRepo:     "",
			wantRepoType: RepoTypeUnknown,
			wantErr:      true, // Need org/repo format
		},
		{
			name:         "SQuAD dataset full form",
			input:        "datasets/rajpurkar/squad",
			wantRepo:     "rajpurkar/squad",
			wantRepoType: RepoTypeDataset,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRepo, gotRepoType, err := ParseHuggingFaceRepo(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseHuggingFaceRepo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotRepo != tt.wantRepo {
				t.Errorf("ParseHuggingFaceRepo() repo = %v, want %v", gotRepo, tt.wantRepo)
			}

			if gotRepoType != tt.wantRepoType {
				t.Errorf("ParseHuggingFaceRepo() repoType = %v, want %v", gotRepoType, tt.wantRepoType)
			}
		})
	}
}

func TestIsValidRepoFormat(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "valid org/repo",
			input: "myorg/myrepo",
			want:  true,
		},
		{
			name:  "valid with hyphens",
			input: "my-org/my-repo",
			want:  true,
		},
		{
			name:  "valid with underscores",
			input: "my_org/my_repo",
			want:  true,
		},
		{
			name:  "single segment",
			input: "myrepo",
			want:  false,
		},
		{
			name:  "three segments",
			input: "org/repo/extra",
			want:  false,
		},
		{
			name:  "empty string",
			input: "",
			want:  false,
		},
		{
			name:  "only slash",
			input: "/",
			want:  false,
		},
		{
			name:  "empty org",
			input: "/myrepo",
			want:  false,
		},
		{
			name:  "empty repo",
			input: "myorg/",
			want:  false,
		},
		{
			name:  "double slash",
			input: "myorg//myrepo",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidRepoFormat(tt.input)
			if got != tt.want {
				t.Errorf("isValidRepoFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}
