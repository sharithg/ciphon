-- name: GetAllNodes :many
SELECT id,
    host,
    name,
    username,
    status,
    convert_from(decode(pem_file, 'base64'), 'UTF8') as pem_file,
    agent_token
FROM nodes;

-- name: CreateNode :one
INSERT INTO nodes (
        host,
        username,
        name,
        pem_file,
        port,
        agent_token
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;

-- name: GetNodeById :one
SELECT id,
    host,
    name,
    username,
    status,
    pem_file,
    agent_token
FROM nodes
WHERE id = $1;

-- name: UpdateNodeStatus :exec
UPDATE nodes
SET status = $1,
    updated_at = now()
WHERE id = $2;

-- name: CreatePipelineRun :one
INSERT INTO pipeline_runs (commit_sha, repo_id, config_file, branch, status)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: CreateRepo :one
INSERT INTO github_repos (
        repo_id,
        name,
        owner,
        description,
        url,
        repo_created_at,
        raw_data
    )
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: GetAllRepos :many
SELECT repo_id,
    name,
    owner,
    description,
    url,
    repo_created_at
FROM github_repos;

-- name: CreateJobRun :one
INSERT INTO job_runs (
        id,
        workflow_id,
        name,
        docker,
        node,
        requires
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;

-- name: GetJobsByWorkflowId :many
SELECT id,
    name,
    status
FROM job_runs
WHERE workflow_id = $1;

-- name: UpdateJobRunStatus :exec
UPDATE job_runs
SET status = $1
WHERE id = $2;

-- name: CreateStepRun :one
INSERT INTO step_runs (
        job_id,
        step_order,
        type,
        name,
        command,
        keys,
        paths
    )
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: GetStepsByJobId :many
SELECT type,
    s.id,
    s.name,
    s.command,
    s.status
FROM step_runs s
    JOIN job_runs j ON s.job_id = j.id
WHERE j.id = $1
ORDER BY step_order;

-- name: UpdateStepRunStatus :exec
UPDATE step_runs
SET status = $1
WHERE id = $2;

-- name: CreateCommandOutput :one
INSERT INTO command_output (step_id, stdout, type)
VALUES ($1, $2, $3)
RETURNING id;

-- name: GetCommandOutputsByStepId :many
SELECT id,
    step_id,
    stdout,
    type,
    created_at
FROM command_output
WHERE step_id = $1;

-- name: CreateUser :one
INSERT INTO users (username, email, external_id, auth_type)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: CreateGitHubUserInfo :exec
INSERT INTO github_user_info (user_id, data)
VALUES ($1, $2);

-- name: GetUserByExternalId :one
SELECT id,
    username,
    email
FROM users
WHERE external_id = $1;

-- name: GetUserById :one
SELECT u.id,
    u.username,
    u.email,
    (gh.data->>'avatar_url')::text AS avatar_url
FROM users u
    JOIN github_user_info gh ON u.id = gh.user_id
WHERE u.id = $1;

-- name: CreateWorkflowRun :one
INSERT INTO workflow_runs (name, pipeline_run_id)
VALUES ($1, $2)
RETURNING id;

-- name: GetWorkflowRuns :many
SELECT pr.commit_sha,
    r.name as repo_name,
    pr.id as pipeline_id,
    w.id as workflow_id,
    w.status,
    w.name as workflow_name,
    pr.branch,
    pr.created_at,
    w.duration
FROM workflow_runs w
    JOIN pipeline_runs pr ON pr.id = w.pipeline_run_id
    JOIN github_repos r ON r.repo_id = pr.repo_id
ORDER BY w.created_at DESC
LIMIT 20;

-- name: GetWorkflowRunById :many
SELECT j.id as job_id,
    s.id as step_id,
    s.command,
    s.type,
    s.keys,
    s.paths,
    s.step_order,
    r.url,
    r.name as repo_name,
    pr.commit_sha,
    pr.branch,
    j.docker
FROM workflow_runs w
    JOIN pipeline_runs pr ON pr.id = w.pipeline_run_id
    JOIN github_repos r ON r.repo_id = pr.repo_id
    JOIN job_runs j ON j.workflow_id = w.id
    JOIN step_runs s ON s.job_id = j.id
WHERE w.id = $1
ORDER BY s.step_order;

-- name: UpdateWorkflowRunStatus :exec
UPDATE workflow_runs
SET status = $1
WHERE id = $2;

-- name: UpdateWorkflowRunDuration :exec
UPDATE workflow_runs
SET duration = $1
WHERE id = $2;

-- name: UpdateWorkflowRunStatusNull :exec
UPDATE workflow_runs
SET status = NULL
WHERE id = $1;

-- name: UpdateJobRunStatusNull :exec
UPDATE job_runs
SET status = NULL
WHERE workflow_id = $1;

-- name: UpdateStepRunStatusNull :exec
UPDATE step_runs
SET status = NULL
WHERE job_id IN (
        SELECT id
        FROM job_runs
        WHERE workflow_id = $1
    );

-- name: DeleteCommandOutputByWorkflowId :exec
DELETE FROM command_output
WHERE step_id IN (
        SELECT id
        FROM step_runs
        WHERE job_id IN (
                SELECT id
                FROM job_runs
                WHERE workflow_id = $1
            )
    );