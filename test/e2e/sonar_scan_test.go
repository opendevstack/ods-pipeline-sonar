package e2e

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"testing"

	"github.com/bix-digital/ods-pipeline-sonar/pkg/sonar"
	ott "github.com/opendevstack/ods-pipeline/pkg/odstasktest"
	"github.com/opendevstack/ods-pipeline/pkg/pipelinectxt"
	ttr "github.com/opendevstack/ods-pipeline/pkg/tektontaskrun"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func TestSonarScanTask(t *testing.T) {
	if err := runTask(
		ott.WithGitSourceWorkspace(t, "../testdata/workspaces/go-sample-app"),
		ttr.AfterRun(func(config *ttr.TaskRunConfig, run *tekton.TaskRun) {
			wsDir, ctxt := ott.GetSourceWorkspaceContext(t, config)
			ott.AssertFilesExist(t, wsDir,
				filepath.Join(pipelinectxt.SonarAnalysisPath, "analysis-report.md"),
				filepath.Join(pipelinectxt.SonarAnalysisPath, "issues-report.csv"),
				filepath.Join(pipelinectxt.SonarAnalysisPath, "quality-gate.json"),
			)
			kubeClient, err := newKubeClient()
			if err != nil {
				t.Fatal(err)
			}
			sonarClient, err := newSonarClient(kubeClient, config.Namespace)
			if err != nil {
				t.Fatal(err)
			}
			sonarProject := sonar.ProjectKey(ctxt, "")
			assertSonarQualityGate(t, sonarClient, sonarProject, "OK")
		}),
	); err != nil {
		t.Fatal(err)
	}
}

func newSonarClient(c *kubernetes.Clientset, namespace string) (*sonar.Client, error) {
	sonarToken, err := getSecretKey(c, namespace, "ods-sonar-auth", "password")
	if err != nil {
		return nil, fmt.Errorf("could not get SonarQube token: %s", err)
	}

	return sonar.NewClient(&sonar.ClientConfig{
		APIToken:      sonarToken,
		BaseURL:       "http://localhost:9000", // use localhost instead of ods-test-sonarqube.kind!
		ServerEdition: "community",
	})
}

func assertSonarQualityGate(t *testing.T, c *sonar.Client, sonarProject string, want string) {
	qualityGateResult, err := c.QualityGateGet(
		sonar.QualityGateGetParams{ProjectKey: sonarProject},
	)
	if err != nil {
		t.Fatal(err)
	}
	got := qualityGateResult.ProjectStatus.Status
	if got != want {
		t.Fatalf("Got: %s, want: %s", got, want)
	}
}

func getSecretKey(clientset *kubernetes.Clientset, namespace, secretName, key string) (string, error) {

	log.Printf("Get secret %s", secretName)

	secret, err := clientset.CoreV1().
		Secrets(namespace).
		Get(context.TODO(), secretName, metav1.GetOptions{})

	if err != nil {
		return "", err
	}

	v, ok := secret.Data[key]
	if !ok {
		return "", fmt.Errorf("key %s not found", key)
	}

	return string(v), err
}

func newKubeClient() (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags(
		"", filepath.Join(homedir.HomeDir(), ".kube", "config"),
	)
	if err != nil {
		return nil, fmt.Errorf("build config from kubeconfig filepath: %s", err)
	}
	return kubernetes.NewForConfig(config)
}
