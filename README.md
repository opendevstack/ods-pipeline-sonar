# ods-pipeline-sonar

[![Tests](https://github.com/opendevstack/ods-pipeline-sonar/actions/workflows/main.yaml/badge.svg)](https://github.com/opendevstack/ods-pipeline-sonar/actions/workflows/main.yaml)

Tekton task for use with [ODS Pipeline](https://github.com/opendevstack/ods-pipeline) to run the SonarQube scanner.

## Usage

```yaml
tasks:
- name: build
  taskRef:
    resolver: git
    params:
    - { name: url, value: https://github.com/opendevstack/ods-pipeline-sonar.git }
    - { name: revision, value: v0.2.0 }
    - { name: pathInRepo, value: tasks/scan.yaml }
    workspaces:
    - { name: source, workspace: shared-workspace }
```

See the [documentation](https://github.com/opendevstack/ods-pipeline-sonar/blob/main/docs/scan.adoc) for prerequisites, details and available parameters.

## About this repository

`docs` and `tasks` are generated directories from recipes located in `build`. See the `Makefile` target for how everything fits together.
