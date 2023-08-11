package sonar

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/bix-digital/ods-pipeline-sonar/internal/file"
	"github.com/opendevstack/ods-pipeline/pkg/pipelinectxt"
)

// GenerateReports generates SonarQube reports using cnesreport.
// See https://github.com/cnescatlab/sonar-cnes-report.
func (c *Client) GenerateReports(sonarProject, author, branch, rootPath, artifactPrefix string) error {
	reportParams := append(
		c.javaSystemProperties(),
		"-jar", "/usr/local/cnes/cnesreport.jar",
		"-s", c.clientConfig.BaseURL,
		"-t", c.clientConfig.APIToken,
		"-p", sonarProject,
		"-a", author,
		branch,
	)
	cmd := exec.Command("java", reportParams...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf(
			"report generation failed: %w", err,
		)
	}

	artifactsPath := filepath.Join(rootPath, pipelinectxt.SonarAnalysisPath)
	err := copyReportFiles(sonarProject, artifactsPath, artifactPrefix)
	if err != nil {
		return fmt.Errorf("copying report to artifacts failed: %w", err)
	}

	return nil
}

func copyReportFiles(project, destinationDir, artifactPrefix string) error {
	analysisReportFile := fmt.Sprintf(
		"%s-%s-analysis-report.md",
		currentDate(),
		project,
	)
	err := file.Copy(
		analysisReportFile,
		filepath.Join(destinationDir, artifactPrefix+"analysis-report.md"),
	)
	if err != nil {
		return fmt.Errorf("copying %s failed: %w", analysisReportFile, err)
	}

	issuesReportFile := fmt.Sprintf(
		"%s-%s-issues-report.csv",
		currentDate(),
		project,
	)
	err = file.Copy(
		issuesReportFile,
		filepath.Join(destinationDir, artifactPrefix+"issues-report.csv"),
	)
	if err != nil {
		return fmt.Errorf("copying %s failed: %w", issuesReportFile, err)
	}
	return nil
}

// currentDate returns the current date as YYYY-MM-DD
func currentDate() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02")
}
