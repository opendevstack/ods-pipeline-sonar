# ods-pipeline-sonar

Tekton task for use with [ODS Pipeline](https://github.com/opendevstack/ods-pipeline) to run the SonarQube scanner.

## Usage

```yaml
tasks:
- name: build
  taskRef:
    resolver: git
    params:
    - { name: url, value: https://github.com/opendevstack/ods-pipeline-sonar.git }
    - { name: revision, value: main }
    - { name: pathInRepo, value: tasks/scan.yaml }
    workspaces:
    - { name: source, workspace: shared-workspace }
```

See the [documentation](https://github.com/opendevstack/ods-pipeline-sonar/blob/main/docs/scan.adoc) for prerequisites, details and available parameters.

## About this repository

`docs` and `tasks` are generated directories from recipes located in `build`. See the `Makefile` target for how everything fits together.
