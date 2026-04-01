// SPDX-FileCopyrightText: Copyright 2026 The Minder Authors
// SPDX-License-Identifier: Apache-2.0

// Package testing provides shared test helpers for the pkg/testing test suite.
package testing

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"gopkg.in/yaml.v3"
)

// fixture holds the contents of a rule-test YAML fixture file.
type fixture struct {
	Version   string     `yaml:"version"`
	RuleName  string     `yaml:"rule_name"`
	TestCases []TestCase `yaml:"test_cases"`
}

// TestCase describes one scenario inside a fixture file.
type TestCase struct {
	Name       string             `yaml:"name"`
	Expect     string             `yaml:"expect"`
	SkipReason string             `yaml:"skip_reason"`
	MockData   ProviderMockConfig `yaml:"mock_data"`
}

// ProviderMockConfig holds mock data for all provider types.
type ProviderMockConfig struct {
	GitFiles            map[string]string           `yaml:"git_files"`
	HTTPResponses       map[string]HTTPResponseMock `yaml:"http_responses"`
	DataSourceResponses map[string]HTTPResponseMock `yaml:"data_source_responses"`
}

// HTTPResponseMock defines a canned HTTP response for a given URL.
type HTTPResponseMock struct {
	StatusCode int    `yaml:"status_code"`
	Body       string `yaml:"body"`
}

// Mocks holds the constructed mock objects for a single test case.
type Mocks struct {
	GitFilesystem    billy.Filesystem
	HTTPClient       *http.Client
	DataSourceClient *http.Client
}

// Result records the outcome of a single test case after DryRun.
type Result struct {
	Name       string
	Err        error
	Skipped    bool
	SkipReason string
}

type mockRoundTripper struct {
	responses map[string]HTTPResponseMock
}

// NewMockRoundTripper creates a mockRoundTripper from the given map.
func NewMockRoundTripper(responses map[string]HTTPResponseMock) *mockRoundTripper {
	return &mockRoundTripper{responses: responses}
}

// RoundTrip implements http.RoundTripper by serving canned responses.
func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	key := req.URL.String()
	if mock, ok := m.responses[key]; ok {
		return &http.Response{
			StatusCode: mock.StatusCode,
			Body:       io.NopCloser(strings.NewReader(mock.Body)),
			Header:     make(http.Header),
		}, nil
	}
	return &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(strings.NewReader(fmt.Sprintf("no mock for URL: %s", key))),
		Header:     make(http.Header),
	}, nil
}

// Parse reads a fixture YAML file from disk and returns the parsed fixture.
func Parse(path string) (*fixture, error) {
	//nolint:gosec // path is provided by test fixtures, not user input
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading fixture %s: %w", path, err)
	}
	var f fixture
	if err := yaml.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("parsing fixture %s: %w", path, err)
	}
	return &f, nil
}

// NewMockBillyFS creates an in-memory billy filesystem pre-populated with files.
func NewMockBillyFS(files map[string]string) (billy.Filesystem, error) {
	fs := memfs.New()
	for path, content := range files {
		f, err := fs.Create(path)
		if err != nil {
			return nil, fmt.Errorf("creating %q: %w", path, err)
		}
		if _, writeErr := f.Write([]byte(content)); writeErr != nil {
			f.Close()
			return nil, fmt.Errorf("writing %q: %w", path, writeErr)
		}
		if err := f.Close(); err != nil {
			return nil, fmt.Errorf("closing %q: %w", path, err)
		}
	}
	return fs, nil
}

// BuildMocks constructs the full set of mock objects for a test case.
func BuildMocks(tc TestCase) (*Mocks, error) {
	m := &Mocks{}
	if len(tc.MockData.GitFiles) > 0 {
		fs, err := NewMockBillyFS(tc.MockData.GitFiles)
		if err != nil {
			return nil, fmt.Errorf("case %q: building git filesystem: %w", tc.Name, err)
		}
		m.GitFilesystem = fs
	} else {
		m.GitFilesystem = memfs.New()
	}
	if len(tc.MockData.HTTPResponses) > 0 {
		m.HTTPClient = &http.Client{Transport: NewMockRoundTripper(tc.MockData.HTTPResponses)}
	} else {
		m.HTTPClient = &http.Client{Transport: NewMockRoundTripper(nil)}
	}
	if len(tc.MockData.DataSourceResponses) > 0 {
		m.DataSourceClient = &http.Client{Transport: NewMockRoundTripper(tc.MockData.DataSourceResponses)}
	} else {
		m.DataSourceClient = &http.Client{Transport: NewMockRoundTripper(nil)}
	}
	return m, nil
}

// DryRun parses a fixture file and validates each test case by calling BuildMocks.
func DryRun(path string) ([]Result, error) {
	fx, err := Parse(path)
	if err != nil {
		return nil, fmt.Errorf("dry-run: %w", err)
	}
	results := make([]Result, 0, len(fx.TestCases))
	for _, tc := range fx.TestCases {
		r := Result{Name: tc.Name}
		if tc.SkipReason != "" {
			r.Skipped = true
			r.SkipReason = tc.SkipReason
			results = append(results, r)
			continue
		}
		_, buildErr := BuildMocks(tc)
		if buildErr != nil {
			r.Err = buildErr
			results = append(results, r)
			continue
		}
		results = append(results, r)
	}
	return results, nil
}

func mustParseURL(t *testing.T, target string) *url.URL {
	t.Helper()
	u, err := url.Parse(target)
	if err != nil {
		t.Fatalf("could not parse URL %q: %v", target, err)
	}
	return u
}
