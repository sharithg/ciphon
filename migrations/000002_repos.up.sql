CREATE TABLE github_repos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    repo_id BIGINT NOT NULL,
    -- GitHub's unique identifier for the repository
    name VARCHAR(255) NOT NULL,
    -- Repository name
    owner VARCHAR(255) NOT NULL,
    -- Owner's username or organization name
    description TEXT,
    -- Repository description
    url TEXT NOT NULL,
    -- URL to the repository
    repo_created_at TIMESTAMP,
    -- When the repository was created on GitHub
    raw_data JSONB NOT NULL,
    -- Raw JSON data from GitHub API response
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE(repo_id) -- Ensure unique GitHub repositories by their ID
);