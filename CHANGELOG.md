# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Note that changes which ONLY affect documentation or the testsuite will not be
listed in the changelog.

## [Unreleased]

## [0.2.0] - 2023-10-09

### Changed

- Migrate from Tekton v1beta1 resources to v1 ([#3](https://github.com/opendevstack/ods-pipeline-sonar/pull/3))

## [0.1.0] - 2023-09-29

NOTE: This version is based on v0.13.2 of the task step `sonar-scan`, which was included in `ods-build-go`, `ods-build-gradle`, `ods-build-npm` and `ods-build-python`. The step has been extracted into its own task for easier maintenance and more flexible use. Listed below are all updates compared to v0.13.2 of the task step.

# Fixed

- sonar-scanner invocations stderr not captured ([#719](https://github.com/opendevstack/ods-pipeline/issues/719))

- sonar-scanner does not start properly: java is lacking tzdb.dat ([#723](https://github.com/opendevstack/ods-pipeline/issues/723))

- update sonar-scanner and cnes-report ([#725](https://github.com/opendevstack/ods-pipeline/issues/725))

- SonarQube doesn't scan FE-related code ([#716](https://github.com/opendevstack/ods-pipeline/issues/716))
