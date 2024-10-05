CREATE TABLE pipeline_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    commit_sha VARCHAR(255) NOT NULL,
    config_file TEXT NOT NULL,
    repo_id BIGINT NOT NULL REFERENCES github_repos(repo_id),
    status VARCHAR(255) NOT NULL,
    branch TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE workflow_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    status VARCHAR(255),
    pipeline_run_id UUID REFERENCES pipeline_runs(id) ON DELETE CASCADE,
    duration FLOAT,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE job_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id UUID REFERENCES workflow_runs(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(255),
    docker VARCHAR(255),
    node VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);


CREATE TABLE step_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID REFERENCES job_runs(id) ON DELETE CASCADE,
    step_order INT,
    type VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    status VARCHAR(255),
    command TEXT,
    keys TEXT[],
    paths TEXT[],
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

