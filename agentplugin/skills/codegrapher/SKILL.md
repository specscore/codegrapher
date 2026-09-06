---
name: codegrapher
description: Use CodeGrapher to index a repository and trace symbol relationships before changing code.
---

# CodeGrapher

Use CodeGrapher when you need evidence about code relationships in the current
repository: symbol definitions, callers, callees, or the impact of a change.

Start by checking whether the repository is indexed:

```sh
codegrapher status --format json
```

If it is not initialized, index it once from the repository root:

```sh
codegrapher init
```

Refresh a previously initialized index after source changes:

```sh
codegrapher sync
```

Use focused queries before editing:

```sh
codegrapher query "<symbol>" --format json
codegrapher callers "<symbol>" --format json
codegrapher callees "<symbol>" --format json
codegrapher impact "<symbol>" --format json
```

Treat the graph as investigation evidence. Read the referenced source before
making a behavior claim, and refresh the index if the relevant files changed.
