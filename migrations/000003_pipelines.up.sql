CREATE TYPE job_status_enum AS ENUM (
    'running',
    'failed',
    'success',
    'not_started',
    'pending',
    'queued',
    'cancelled'
);

CREATE TYPE workflow_status_enum AS ENUM (
    'running',
    'failed',
    'success',
    'not_started',
    'pending',
    'queued',
    'cancelled'
);

CREATE TYPE step_status_enum AS ENUM (
    'running',
    'failed',
    'success',
    'not_started',
    'pending',
    'queued',
    'cancelled'
);

CREATE TABLE pipeline_refs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    commit_sha VARCHAR(255) NOT NULL,
    config_file TEXT NOT NULL,
    repo_id BIGINT NOT NULL REFERENCES github_repos(repo_id),
    branch TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE workflow_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    status workflow_status_enum DEFAULT 'not_started' NOT NULL,
    pipeline_ref_id UUID REFERENCES pipeline_refs(id) ON DELETE CASCADE,
    duration FLOAT,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE job_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id UUID REFERENCES workflow_runs(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    status job_status_enum DEFAULT 'not_started' NOT NULL,
    docker VARCHAR(255) NOT NULL,
    node VARCHAR(255),
    requires UUID [],
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE step_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID REFERENCES job_runs(id) ON DELETE CASCADE,
    step_order INT NOT NULL,
    type VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    status step_status_enum DEFAULT 'not_started' NOT NULL,
    command TEXT,
    keys TEXT [],
    paths TEXT [],
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE command_output (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    step_id UUID REFERENCES step_runs(id) ON DELETE CASCADE,
    stdout TEXT NOT NULL,
    type VARCHAR(255) NOT NULL,
    output_order INT,
    created_at TIMESTAMPTZ DEFAULT now()
);