# ods-pipeline-sonar

Tekton task for use with [ODS Pipeline](https://github.com/opendevstack/ods-pipeline) to run the SonarQube scanner.

## Usage

```yaml
tasks:
- name: build
  taskRef:
    resolver: git
    params:
    - { name: url, value: https://github.com/bix-digital/ods-pipeline-sonar.git }
    - { name: revision, value: latest }
    - { name: pathInRepo, value: tasks/scan.yaml }
    workspaces:
    - { name: source, workspace: shared-workspace }
```

See the [documentation](https://github.com/BIX-Digital/ods-pipeline-sonar/blob/main/docs/deploy.adoc) for details and available parameters.

## About this repository

`docs` and `tasks` are generated directories from recipes located in `build`. See the `Makefile` target for how everything fits together.