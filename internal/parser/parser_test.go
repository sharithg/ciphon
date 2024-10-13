package parser

import (
	"testing"
)

// TestParseConfig tests the ParseConfig function with different YAML inputs.
func TestParseConfig(t *testing.T) {
	tests := []struct {
		name              string
		yamlData          string
		expectedError     bool
		expectedJobs      int
		expectedWorkflows int
	}{
		{
			name: "Valid configuration with all steps",
			yamlData: `
version: 0.1

# asd
jobs:
    build-test-lint:
        docker: node
        node: siphon
        steps:
            - checkout
            - restore_cache:
                  name: Restore pnpm Package Cache
                  keys:
                      - pnpm-packages-{{ checksum "pnpm-lock.yaml" }}
            - run:
                  name: Install pnpm package manager
                  command: |
                      npm i -g pnpm
                      corepack enable
                      corepack prepare pnpm@latest-8 --activate
            - run:
                  name: Install Dependencies
                  command: |
                      pnpm install
            - save_cache:
                  name: Save pnpm Package Cache
                  key: pnpm-packages-{{ checksum "pnpm-lock.yaml" }}
                  paths:
                      - node_modules
                      - /home/circleci/.cache/Cypress
            - run:
                  name: Check format
                  command: pnpm format
            - run:
                  name: Run build
                  command: pnpm build
            - run:
                  name: Run lint
                  command: pnpm lint
            - run:
                  name: Run test
                  command: |
                      pnpm test

    test-flow:
        docker: node
        node: siphon
        steps:
            - checkout
            - run:
                  name: Install pnpm package manager
                  command: ls

workflows:
    ci:
        jobs:
            - build-test-lint
            - test-flow:
                  requires: build-test-lint
`,
			expectedError:     false,
			expectedJobs:      2,
			expectedWorkflows: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := ParseConfig(tt.yamlData)

			if (err != nil) != tt.expectedError {
				t.Errorf("ParseConfig() error = %v, expectedError %v", err, tt.expectedError)
				return
			}

			if config != nil && len(config.Jobs) != tt.expectedJobs {
				t.Errorf("ParseConfig() jobs = %v, expectedJobs %v", len(config.Jobs), tt.expectedJobs)
			}

			if config != nil && len(config.Workflows) != tt.expectedWorkflows {
				t.Errorf("ParseConfig() jobs = %v, expectedWorkflows %v", len(config.Workflows), tt.expectedWorkflows)
			}
		})
	}
}
