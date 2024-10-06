/* prettier-ignore-start */

/* eslint-disable */

// @ts-nocheck

// noinspection JSUnusedGlobalSymbols

// This file is auto-generated by TanStack Router

import { createFileRoute } from '@tanstack/react-router'

// Import Routes

import { Route as rootRoute } from './routes/__root'
import { Route as PipelinesWorkflowsWorkflowIdIndexImport } from './routes/pipelines/workflows/$workflowId/index'

// Create Virtual Routes

const SettingsLazyImport = createFileRoute('/settings')()
const ProjectsLazyImport = createFileRoute('/projects')()
const IndexLazyImport = createFileRoute('/')()
const PipelinesIndexLazyImport = createFileRoute('/pipelines/')()
const PipelinesWorkflowsWorkflowIdJobsJobIdLazyImport = createFileRoute(
  '/pipelines/workflows/$workflowId/jobs/$jobId',
)()

// Create/Update Routes

const SettingsLazyRoute = SettingsLazyImport.update({
  path: '/settings',
  getParentRoute: () => rootRoute,
} as any).lazy(() => import('./routes/settings.lazy').then((d) => d.Route))

const ProjectsLazyRoute = ProjectsLazyImport.update({
  path: '/projects',
  getParentRoute: () => rootRoute,
} as any).lazy(() => import('./routes/projects.lazy').then((d) => d.Route))

const IndexLazyRoute = IndexLazyImport.update({
  path: '/',
  getParentRoute: () => rootRoute,
} as any).lazy(() => import('./routes/index.lazy').then((d) => d.Route))

const PipelinesIndexLazyRoute = PipelinesIndexLazyImport.update({
  path: '/pipelines/',
  getParentRoute: () => rootRoute,
} as any).lazy(() =>
  import('./routes/pipelines/index.lazy').then((d) => d.Route),
)

const PipelinesWorkflowsWorkflowIdIndexRoute =
  PipelinesWorkflowsWorkflowIdIndexImport.update({
    path: '/pipelines/workflows/$workflowId/',
    getParentRoute: () => rootRoute,
  } as any)

const PipelinesWorkflowsWorkflowIdJobsJobIdLazyRoute =
  PipelinesWorkflowsWorkflowIdJobsJobIdLazyImport.update({
    path: '/pipelines/workflows/$workflowId/jobs/$jobId',
    getParentRoute: () => rootRoute,
  } as any).lazy(() =>
    import('./routes/pipelines/workflows/$workflowId/jobs/$jobId.lazy').then(
      (d) => d.Route,
    ),
  )

// Populate the FileRoutesByPath interface

declare module '@tanstack/react-router' {
  interface FileRoutesByPath {
    '/': {
      id: '/'
      path: '/'
      fullPath: '/'
      preLoaderRoute: typeof IndexLazyImport
      parentRoute: typeof rootRoute
    }
    '/projects': {
      id: '/projects'
      path: '/projects'
      fullPath: '/projects'
      preLoaderRoute: typeof ProjectsLazyImport
      parentRoute: typeof rootRoute
    }
    '/settings': {
      id: '/settings'
      path: '/settings'
      fullPath: '/settings'
      preLoaderRoute: typeof SettingsLazyImport
      parentRoute: typeof rootRoute
    }
    '/pipelines/': {
      id: '/pipelines/'
      path: '/pipelines'
      fullPath: '/pipelines'
      preLoaderRoute: typeof PipelinesIndexLazyImport
      parentRoute: typeof rootRoute
    }
    '/pipelines/workflows/$workflowId/': {
      id: '/pipelines/workflows/$workflowId/'
      path: '/pipelines/workflows/$workflowId'
      fullPath: '/pipelines/workflows/$workflowId'
      preLoaderRoute: typeof PipelinesWorkflowsWorkflowIdIndexImport
      parentRoute: typeof rootRoute
    }
    '/pipelines/workflows/$workflowId/jobs/$jobId': {
      id: '/pipelines/workflows/$workflowId/jobs/$jobId'
      path: '/pipelines/workflows/$workflowId/jobs/$jobId'
      fullPath: '/pipelines/workflows/$workflowId/jobs/$jobId'
      preLoaderRoute: typeof PipelinesWorkflowsWorkflowIdJobsJobIdLazyImport
      parentRoute: typeof rootRoute
    }
  }
}

// Create and export the route tree

export interface FileRoutesByFullPath {
  '/': typeof IndexLazyRoute
  '/projects': typeof ProjectsLazyRoute
  '/settings': typeof SettingsLazyRoute
  '/pipelines': typeof PipelinesIndexLazyRoute
  '/pipelines/workflows/$workflowId': typeof PipelinesWorkflowsWorkflowIdIndexRoute
  '/pipelines/workflows/$workflowId/jobs/$jobId': typeof PipelinesWorkflowsWorkflowIdJobsJobIdLazyRoute
}

export interface FileRoutesByTo {
  '/': typeof IndexLazyRoute
  '/projects': typeof ProjectsLazyRoute
  '/settings': typeof SettingsLazyRoute
  '/pipelines': typeof PipelinesIndexLazyRoute
  '/pipelines/workflows/$workflowId': typeof PipelinesWorkflowsWorkflowIdIndexRoute
  '/pipelines/workflows/$workflowId/jobs/$jobId': typeof PipelinesWorkflowsWorkflowIdJobsJobIdLazyRoute
}

export interface FileRoutesById {
  __root__: typeof rootRoute
  '/': typeof IndexLazyRoute
  '/projects': typeof ProjectsLazyRoute
  '/settings': typeof SettingsLazyRoute
  '/pipelines/': typeof PipelinesIndexLazyRoute
  '/pipelines/workflows/$workflowId/': typeof PipelinesWorkflowsWorkflowIdIndexRoute
  '/pipelines/workflows/$workflowId/jobs/$jobId': typeof PipelinesWorkflowsWorkflowIdJobsJobIdLazyRoute
}

export interface FileRouteTypes {
  fileRoutesByFullPath: FileRoutesByFullPath
  fullPaths:
    | '/'
    | '/projects'
    | '/settings'
    | '/pipelines'
    | '/pipelines/workflows/$workflowId'
    | '/pipelines/workflows/$workflowId/jobs/$jobId'
  fileRoutesByTo: FileRoutesByTo
  to:
    | '/'
    | '/projects'
    | '/settings'
    | '/pipelines'
    | '/pipelines/workflows/$workflowId'
    | '/pipelines/workflows/$workflowId/jobs/$jobId'
  id:
    | '__root__'
    | '/'
    | '/projects'
    | '/settings'
    | '/pipelines/'
    | '/pipelines/workflows/$workflowId/'
    | '/pipelines/workflows/$workflowId/jobs/$jobId'
  fileRoutesById: FileRoutesById
}

export interface RootRouteChildren {
  IndexLazyRoute: typeof IndexLazyRoute
  ProjectsLazyRoute: typeof ProjectsLazyRoute
  SettingsLazyRoute: typeof SettingsLazyRoute
  PipelinesIndexLazyRoute: typeof PipelinesIndexLazyRoute
  PipelinesWorkflowsWorkflowIdIndexRoute: typeof PipelinesWorkflowsWorkflowIdIndexRoute
  PipelinesWorkflowsWorkflowIdJobsJobIdLazyRoute: typeof PipelinesWorkflowsWorkflowIdJobsJobIdLazyRoute
}

const rootRouteChildren: RootRouteChildren = {
  IndexLazyRoute: IndexLazyRoute,
  ProjectsLazyRoute: ProjectsLazyRoute,
  SettingsLazyRoute: SettingsLazyRoute,
  PipelinesIndexLazyRoute: PipelinesIndexLazyRoute,
  PipelinesWorkflowsWorkflowIdIndexRoute:
    PipelinesWorkflowsWorkflowIdIndexRoute,
  PipelinesWorkflowsWorkflowIdJobsJobIdLazyRoute:
    PipelinesWorkflowsWorkflowIdJobsJobIdLazyRoute,
}

export const routeTree = rootRoute
  ._addFileChildren(rootRouteChildren)
  ._addFileTypes<FileRouteTypes>()

/* prettier-ignore-end */

/* ROUTE_MANIFEST_START
{
  "routes": {
    "__root__": {
      "filePath": "__root.tsx",
      "children": [
        "/",
        "/projects",
        "/settings",
        "/pipelines/",
        "/pipelines/workflows/$workflowId/",
        "/pipelines/workflows/$workflowId/jobs/$jobId"
      ]
    },
    "/": {
      "filePath": "index.lazy.tsx"
    },
    "/projects": {
      "filePath": "projects.lazy.tsx"
    },
    "/settings": {
      "filePath": "settings.lazy.tsx"
    },
    "/pipelines/": {
      "filePath": "pipelines/index.lazy.tsx"
    },
    "/pipelines/workflows/$workflowId/": {
      "filePath": "pipelines/workflows/$workflowId/index.tsx"
    },
    "/pipelines/workflows/$workflowId/jobs/$jobId": {
      "filePath": "pipelines/workflows/$workflowId/jobs/$jobId.lazy.tsx"
    }
  }
}
ROUTE_MANIFEST_END */
