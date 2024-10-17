// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package repository

import (
	"context"
	"time"

	github "github.com/google/go-github/v65/github"
	"github.com/google/uuid"
	auth "github.com/sharithg/siphon/internal/auth"
)

const createCommandOutput = `-- name: CreateCommandOutput :one
INSERT INTO command_output (step_id, stdout, type)
VALUES ($1, $2, $3)
RETURNING id
`

type CreateCommandOutputParams struct {
	StepID uuid.UUID `json:"stepId"`
	Stdout string    `json:"stdout"`
	Type   string    `json:"type"`
}

func (q *Queries) CreateCommandOutput(ctx context.Context, arg CreateCommandOutputParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createCommandOutput, arg.StepID, arg.Stdout, arg.Type)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const createGitHubUserInfo = `-- name: CreateGitHubUserInfo :exec
INSERT INTO github_user_info (user_id, data)
VALUES ($1, $2)
`

type CreateGitHubUserInfoParams struct {
	UserID uuid.UUID       `json:"userId"`
	Data   auth.GitHubUser `json:"data"`
}

func (q *Queries) CreateGitHubUserInfo(ctx context.Context, arg CreateGitHubUserInfoParams) error {
	_, err := q.db.Exec(ctx, createGitHubUserInfo, arg.UserID, arg.Data)
	return err
}

const createJobRun = `-- name: CreateJobRun :one
INSERT INTO job_runs (
        id,
        workflow_id,
        name,
        docker,
        node,
        requires
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id
`

type CreateJobRunParams struct {
	ID         uuid.UUID   `json:"id"`
	WorkflowID uuid.UUID   `json:"workflowId"`
	Name       string      `json:"name"`
	Docker     string      `json:"docker"`
	Node       *string     `json:"node"`
	Requires   []uuid.UUID `json:"requires"`
}

func (q *Queries) CreateJobRun(ctx context.Context, arg CreateJobRunParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createJobRun,
		arg.ID,
		arg.WorkflowID,
		arg.Name,
		arg.Docker,
		arg.Node,
		arg.Requires,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const createNode = `-- name: CreateNode :one
INSERT INTO nodes (
        host,
        username,
        name,
        pem_file,
        port,
        agent_token
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id
`

type CreateNodeParams struct {
	Host       string `json:"host"`
	Username   string `json:"username"`
	Name       string `json:"name"`
	PemFile    string `json:"pemFile"`
	Port       int32  `json:"port"`
	AgentToken string `json:"agentToken"`
}

func (q *Queries) CreateNode(ctx context.Context, arg CreateNodeParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createNode,
		arg.Host,
		arg.Username,
		arg.Name,
		arg.PemFile,
		arg.Port,
		arg.AgentToken,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const createPipelineRun = `-- name: CreatePipelineRun :one
INSERT INTO pipeline_runs (commit_sha, repo_id, config_file, branch, status)
VALUES ($1, $2, $3, $4, $5)
RETURNING id
`

type CreatePipelineRunParams struct {
	CommitSha  string `json:"commitSha"`
	RepoID     int64  `json:"repoId"`
	ConfigFile string `json:"configFile"`
	Branch     string `json:"branch"`
	Status     string `json:"status"`
}

func (q *Queries) CreatePipelineRun(ctx context.Context, arg CreatePipelineRunParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createPipelineRun,
		arg.CommitSha,
		arg.RepoID,
		arg.ConfigFile,
		arg.Branch,
		arg.Status,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const createRepo = `-- name: CreateRepo :one
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
RETURNING id
`

type CreateRepoParams struct {
	RepoID        int64             `json:"repoId"`
	Name          string            `json:"name"`
	Owner         string            `json:"owner"`
	Description   *string           `json:"description"`
	Url           string            `json:"url"`
	RepoCreatedAt time.Time         `json:"repoCreatedAt"`
	RawData       github.Repository `json:"rawData"`
}

func (q *Queries) CreateRepo(ctx context.Context, arg CreateRepoParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createRepo,
		arg.RepoID,
		arg.Name,
		arg.Owner,
		arg.Description,
		arg.Url,
		arg.RepoCreatedAt,
		arg.RawData,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const createStepRun = `-- name: CreateStepRun :one
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
RETURNING id
`

type CreateStepRunParams struct {
	JobID     uuid.UUID `json:"jobId"`
	StepOrder int32     `json:"stepOrder"`
	Type      string    `json:"type"`
	Name      *string   `json:"name"`
	Command   *string   `json:"command"`
	Keys      []string  `json:"keys"`
	Paths     []string  `json:"paths"`
}

func (q *Queries) CreateStepRun(ctx context.Context, arg CreateStepRunParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createStepRun,
		arg.JobID,
		arg.StepOrder,
		arg.Type,
		arg.Name,
		arg.Command,
		arg.Keys,
		arg.Paths,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (username, email, external_id, auth_type)
VALUES ($1, $2, $3, $4)
RETURNING id
`

type CreateUserParams struct {
	Username   string  `json:"username"`
	Email      string  `json:"email"`
	ExternalID string  `json:"externalId"`
	AuthType   *string `json:"authType"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Username,
		arg.Email,
		arg.ExternalID,
		arg.AuthType,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const createWorkflowRun = `-- name: CreateWorkflowRun :one
INSERT INTO workflow_runs (name, pipeline_run_id)
VALUES ($1, $2)
RETURNING id
`

type CreateWorkflowRunParams struct {
	Name          string    `json:"name"`
	PipelineRunID uuid.UUID `json:"pipelineRunId"`
}

func (q *Queries) CreateWorkflowRun(ctx context.Context, arg CreateWorkflowRunParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createWorkflowRun, arg.Name, arg.PipelineRunID)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const deleteCommandOutputByWorkflowId = `-- name: DeleteCommandOutputByWorkflowId :exec
DELETE FROM command_output
WHERE step_id IN (
        SELECT id
        FROM step_runs
        WHERE job_id IN (
                SELECT id
                FROM job_runs
                WHERE workflow_id = $1
            )
    )
`

func (q *Queries) DeleteCommandOutputByWorkflowId(ctx context.Context, workflowID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteCommandOutputByWorkflowId, workflowID)
	return err
}

const getAllNodes = `-- name: GetAllNodes :many
SELECT id,
    host,
    name,
    username,
    status,
    convert_from(decode(pem_file, 'base64'), 'UTF8') as pem_file,
    agent_token
FROM nodes
`

type GetAllNodesRow struct {
	ID         uuid.UUID `json:"id"`
	Host       string    `json:"host"`
	Name       string    `json:"name"`
	Username   string    `json:"username"`
	Status     string    `json:"status"`
	PemFile    string    `json:"pemFile"`
	AgentToken string    `json:"agentToken"`
}

func (q *Queries) GetAllNodes(ctx context.Context) ([]GetAllNodesRow, error) {
	rows, err := q.db.Query(ctx, getAllNodes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllNodesRow
	for rows.Next() {
		var i GetAllNodesRow
		if err := rows.Scan(
			&i.ID,
			&i.Host,
			&i.Name,
			&i.Username,
			&i.Status,
			&i.PemFile,
			&i.AgentToken,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllRepos = `-- name: GetAllRepos :many
SELECT repo_id,
    name,
    owner,
    description,
    url,
    repo_created_at
FROM github_repos
`

type GetAllReposRow struct {
	RepoID        int64     `json:"repoId"`
	Name          string    `json:"name"`
	Owner         string    `json:"owner"`
	Description   *string   `json:"description"`
	Url           string    `json:"url"`
	RepoCreatedAt time.Time `json:"repoCreatedAt"`
}

func (q *Queries) GetAllRepos(ctx context.Context) ([]GetAllReposRow, error) {
	rows, err := q.db.Query(ctx, getAllRepos)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllReposRow
	for rows.Next() {
		var i GetAllReposRow
		if err := rows.Scan(
			&i.RepoID,
			&i.Name,
			&i.Owner,
			&i.Description,
			&i.Url,
			&i.RepoCreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCommandOutputsByStepId = `-- name: GetCommandOutputsByStepId :many
SELECT id,
    step_id,
    stdout,
    type,
    created_at
FROM command_output
WHERE step_id = $1
`

type GetCommandOutputsByStepIdRow struct {
	ID        uuid.UUID `json:"id"`
	StepID    uuid.UUID `json:"stepId"`
	Stdout    string    `json:"stdout"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
}

func (q *Queries) GetCommandOutputsByStepId(ctx context.Context, stepID uuid.UUID) ([]GetCommandOutputsByStepIdRow, error) {
	rows, err := q.db.Query(ctx, getCommandOutputsByStepId, stepID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCommandOutputsByStepIdRow
	for rows.Next() {
		var i GetCommandOutputsByStepIdRow
		if err := rows.Scan(
			&i.ID,
			&i.StepID,
			&i.Stdout,
			&i.Type,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getJobsAndStepsByWorkflowId = `-- name: GetJobsAndStepsByWorkflowId :many
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
    j.docker,
    j.requires
FROM workflow_runs w
    JOIN pipeline_runs pr ON pr.id = w.pipeline_run_id
    JOIN github_repos r ON r.repo_id = pr.repo_id
    JOIN job_runs j ON j.workflow_id = w.id
    JOIN step_runs s ON s.job_id = j.id
WHERE w.id = $1
ORDER BY s.step_order
`

type GetJobsAndStepsByWorkflowIdRow struct {
	JobID     uuid.UUID   `json:"jobId"`
	StepID    uuid.UUID   `json:"stepId"`
	Command   *string     `json:"command"`
	Type      string      `json:"type"`
	Keys      []string    `json:"keys"`
	Paths     []string    `json:"paths"`
	StepOrder int32       `json:"stepOrder"`
	Url       string      `json:"url"`
	RepoName  string      `json:"repoName"`
	CommitSha string      `json:"commitSha"`
	Branch    string      `json:"branch"`
	Docker    string      `json:"docker"`
	Requires  []uuid.UUID `json:"requires"`
}

func (q *Queries) GetJobsAndStepsByWorkflowId(ctx context.Context, id uuid.UUID) ([]GetJobsAndStepsByWorkflowIdRow, error) {
	rows, err := q.db.Query(ctx, getJobsAndStepsByWorkflowId, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetJobsAndStepsByWorkflowIdRow
	for rows.Next() {
		var i GetJobsAndStepsByWorkflowIdRow
		if err := rows.Scan(
			&i.JobID,
			&i.StepID,
			&i.Command,
			&i.Type,
			&i.Keys,
			&i.Paths,
			&i.StepOrder,
			&i.Url,
			&i.RepoName,
			&i.CommitSha,
			&i.Branch,
			&i.Docker,
			&i.Requires,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getJobsByWorkflowId = `-- name: GetJobsByWorkflowId :many
SELECT id,
    name,
    status,
    requires
FROM job_runs
WHERE workflow_id = $1
`

type GetJobsByWorkflowIdRow struct {
	ID       uuid.UUID   `json:"id"`
	Name     string      `json:"name"`
	Status   *string     `json:"status"`
	Requires []uuid.UUID `json:"requires"`
}

func (q *Queries) GetJobsByWorkflowId(ctx context.Context, workflowID uuid.UUID) ([]GetJobsByWorkflowIdRow, error) {
	rows, err := q.db.Query(ctx, getJobsByWorkflowId, workflowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetJobsByWorkflowIdRow
	for rows.Next() {
		var i GetJobsByWorkflowIdRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Status,
			&i.Requires,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNodeById = `-- name: GetNodeById :one
SELECT id,
    host,
    name,
    username,
    status,
    pem_file,
    agent_token
FROM nodes
WHERE id = $1
`

type GetNodeByIdRow struct {
	ID         uuid.UUID `json:"id"`
	Host       string    `json:"host"`
	Name       string    `json:"name"`
	Username   string    `json:"username"`
	Status     string    `json:"status"`
	PemFile    string    `json:"pemFile"`
	AgentToken string    `json:"agentToken"`
}

func (q *Queries) GetNodeById(ctx context.Context, id uuid.UUID) (GetNodeByIdRow, error) {
	row := q.db.QueryRow(ctx, getNodeById, id)
	var i GetNodeByIdRow
	err := row.Scan(
		&i.ID,
		&i.Host,
		&i.Name,
		&i.Username,
		&i.Status,
		&i.PemFile,
		&i.AgentToken,
	)
	return i, err
}

const getStepsByJobId = `-- name: GetStepsByJobId :many
SELECT type,
    s.id,
    s.name,
    s.command,
    s.status
FROM step_runs s
    JOIN job_runs j ON s.job_id = j.id
WHERE j.id = $1
ORDER BY step_order
`

type GetStepsByJobIdRow struct {
	Type    string    `json:"type"`
	ID      uuid.UUID `json:"id"`
	Name    *string   `json:"name"`
	Command *string   `json:"command"`
	Status  *string   `json:"status"`
}

func (q *Queries) GetStepsByJobId(ctx context.Context, id uuid.UUID) ([]GetStepsByJobIdRow, error) {
	rows, err := q.db.Query(ctx, getStepsByJobId, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetStepsByJobIdRow
	for rows.Next() {
		var i GetStepsByJobIdRow
		if err := rows.Scan(
			&i.Type,
			&i.ID,
			&i.Name,
			&i.Command,
			&i.Status,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByExternalId = `-- name: GetUserByExternalId :one
SELECT id,
    username,
    email
FROM users
WHERE external_id = $1
`

type GetUserByExternalIdRow struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

func (q *Queries) GetUserByExternalId(ctx context.Context, externalID string) (GetUserByExternalIdRow, error) {
	row := q.db.QueryRow(ctx, getUserByExternalId, externalID)
	var i GetUserByExternalIdRow
	err := row.Scan(&i.ID, &i.Username, &i.Email)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT u.id,
    u.username,
    u.email,
    (gh.data->>'avatar_url')::text AS avatar_url
FROM users u
    JOIN github_user_info gh ON u.id = gh.user_id
WHERE u.id = $1
`

type GetUserByIdRow struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	AvatarUrl string    `json:"avatarUrl"`
}

func (q *Queries) GetUserById(ctx context.Context, id uuid.UUID) (GetUserByIdRow, error) {
	row := q.db.QueryRow(ctx, getUserById, id)
	var i GetUserByIdRow
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.AvatarUrl,
	)
	return i, err
}

const getWorkflowRuns = `-- name: GetWorkflowRuns :many
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
LIMIT 20
`

type GetWorkflowRunsRow struct {
	CommitSha    string    `json:"commitSha"`
	RepoName     string    `json:"repoName"`
	PipelineID   uuid.UUID `json:"pipelineId"`
	WorkflowID   uuid.UUID `json:"workflowId"`
	Status       *string   `json:"status"`
	WorkflowName string    `json:"workflowName"`
	Branch       string    `json:"branch"`
	CreatedAt    time.Time `json:"createdAt"`
	Duration     *float64  `json:"duration"`
}

func (q *Queries) GetWorkflowRuns(ctx context.Context) ([]GetWorkflowRunsRow, error) {
	rows, err := q.db.Query(ctx, getWorkflowRuns)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetWorkflowRunsRow
	for rows.Next() {
		var i GetWorkflowRunsRow
		if err := rows.Scan(
			&i.CommitSha,
			&i.RepoName,
			&i.PipelineID,
			&i.WorkflowID,
			&i.Status,
			&i.WorkflowName,
			&i.Branch,
			&i.CreatedAt,
			&i.Duration,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateJobRunStatus = `-- name: UpdateJobRunStatus :exec
UPDATE job_runs
SET status = $1
WHERE id = $2
`

type UpdateJobRunStatusParams struct {
	Status *string   `json:"status"`
	ID     uuid.UUID `json:"id"`
}

func (q *Queries) UpdateJobRunStatus(ctx context.Context, arg UpdateJobRunStatusParams) error {
	_, err := q.db.Exec(ctx, updateJobRunStatus, arg.Status, arg.ID)
	return err
}

const updateJobRunStatusNull = `-- name: UpdateJobRunStatusNull :exec
UPDATE job_runs
SET status = NULL
WHERE workflow_id = $1
`

func (q *Queries) UpdateJobRunStatusNull(ctx context.Context, workflowID uuid.UUID) error {
	_, err := q.db.Exec(ctx, updateJobRunStatusNull, workflowID)
	return err
}

const updateNodeStatus = `-- name: UpdateNodeStatus :exec
UPDATE nodes
SET status = $1,
    updated_at = now()
WHERE id = $2
`

type UpdateNodeStatusParams struct {
	Status string    `json:"status"`
	ID     uuid.UUID `json:"id"`
}

func (q *Queries) UpdateNodeStatus(ctx context.Context, arg UpdateNodeStatusParams) error {
	_, err := q.db.Exec(ctx, updateNodeStatus, arg.Status, arg.ID)
	return err
}

const updateStepRunStatus = `-- name: UpdateStepRunStatus :exec
UPDATE step_runs
SET status = $1
WHERE id = $2
`

type UpdateStepRunStatusParams struct {
	Status *string   `json:"status"`
	ID     uuid.UUID `json:"id"`
}

func (q *Queries) UpdateStepRunStatus(ctx context.Context, arg UpdateStepRunStatusParams) error {
	_, err := q.db.Exec(ctx, updateStepRunStatus, arg.Status, arg.ID)
	return err
}

const updateStepRunStatusNull = `-- name: UpdateStepRunStatusNull :exec
UPDATE step_runs
SET status = NULL
WHERE job_id IN (
        SELECT id
        FROM job_runs
        WHERE workflow_id = $1
    )
`

func (q *Queries) UpdateStepRunStatusNull(ctx context.Context, workflowID uuid.UUID) error {
	_, err := q.db.Exec(ctx, updateStepRunStatusNull, workflowID)
	return err
}

const updateWorkflowRunDuration = `-- name: UpdateWorkflowRunDuration :exec
UPDATE workflow_runs
SET duration = $1
WHERE id = $2
`

type UpdateWorkflowRunDurationParams struct {
	Duration *float64  `json:"duration"`
	ID       uuid.UUID `json:"id"`
}

func (q *Queries) UpdateWorkflowRunDuration(ctx context.Context, arg UpdateWorkflowRunDurationParams) error {
	_, err := q.db.Exec(ctx, updateWorkflowRunDuration, arg.Duration, arg.ID)
	return err
}

const updateWorkflowRunStatus = `-- name: UpdateWorkflowRunStatus :exec
UPDATE workflow_runs
SET status = $1
WHERE id = $2
`

type UpdateWorkflowRunStatusParams struct {
	Status *string   `json:"status"`
	ID     uuid.UUID `json:"id"`
}

func (q *Queries) UpdateWorkflowRunStatus(ctx context.Context, arg UpdateWorkflowRunStatusParams) error {
	_, err := q.db.Exec(ctx, updateWorkflowRunStatus, arg.Status, arg.ID)
	return err
}

const updateWorkflowRunStatusNull = `-- name: UpdateWorkflowRunStatusNull :exec
UPDATE workflow_runs
SET status = NULL
WHERE id = $1
`

func (q *Queries) UpdateWorkflowRunStatusNull(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, updateWorkflowRunStatusNull, id)
	return err
}
