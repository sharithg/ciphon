export type TGithubRepoResponse = {
  id: number;
  name: string;
  description: string;
  lastUpdated: string;
  owner: string;
};

export type TConnectRepoRequest = {
  name: string;
  owner: string;
};

export type TNode = {
  id: string;
  host: string;
  name: string;
  user: string;
  port: number;
  status: string;
};

export type TTokenPair = {
  accessToken: string;
  refreshToken: string;
};

export type TJobs = {
  id: string;
  name: string;
  status: string;
};

export type TListRepo = {
  repoId: number;
  name: string;
  owner: string;
  description: string;
  url: string;
  repoCreatedAt: string;
};

export type TSteps = {
  type: string;
  id: string;
  name: string;
  command: string;
  status: string;
};

export type TCommandOutput = {
  id: string;
  step_id: string;
  stdout: string;
  type: string;
  created_at: string;
};

export type TUserDisplay = {
  id: string;
  username: string;
  email: string;
  avatarUrl: string;
};

export type TWorkflowRunInfo = {
  commitSha: string;
  repoName: string;
  workflowName: string;
  pipelineId: string;
  workflowId: string;
  status: string | null;
  branch: string;
  createdAt: string;
  duration: number | null;
};

