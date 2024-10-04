package parser

import (
	"testing"
)

// TestParseConfig tests the ParseConfig function with different YAML inputs.
func TestParseConfig(t *testing.T) {
	tests := []struct {
		name          string
		yamlData      string
		expectedError bool
		expectedJobs  int
	}{
		{
			name: "Valid configuration with all steps",
			yamlData: `
version: 0.1

jobs:
  build-test-lint:
    docker: cimg/node:18.20-browsers
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
            sudo corepack enable
            sudo corepack prepare pnpm@latest-8 --activate
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
            echo $RUNTIME_CONFIG_DEV | base64 --decode > .runtimeconfig.json
            pnpm test

workflows:
  ci:
    jobs:
      - build-test-lint
`,
			expectedError: false,
			expectedJobs:  1,
		},
		{
			name: "Invalid configuration",
			yamlData: `
version: 0.1

jobs:
  build-test-lint:
    docker: cimg/node:18.20-browsers
    node: siphon
    steps:
      - invalid_step:  # This is an invalid step format
          name: Invalid step
`,
			expectedError: true,
			expectedJobs:  0,
		},
		{
			name: "Only checkout step",
			yamlData: `
version: 0.1

jobs:
  build-test-lint:
    docker: cimg/node:18.20-browsers
    node: siphon
    steps:
      - checkout
`,
			expectedError: false,
			expectedJobs:  1,
		},
		{
			name: "Multiple steps with mixed formats",
			yamlData: `
version: 0.1

jobs:
  build-test-lint:
    docker: cimg/node:18.20-browsers
    node: siphon
    steps:
      - checkout
      - run:
          name: Test step
          command: echo "Hello"
      - restore_cache:
          name: Restore cache
          keys:
            - cache-key
`,
			expectedError: false,
			expectedJobs:  1,
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
		})
	}
}
