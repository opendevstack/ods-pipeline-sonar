package e2e

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	ott "github.com/opendevstack/ods-pipeline/pkg/odstasktest"
	ttr "github.com/opendevstack/ods-pipeline/pkg/tektontaskrun"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	namespaceConfig *ttr.NamespaceConfig
	rootPath        = "../.."
	taskName        = "ods-pipeline-sonar-scan"
)

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	cc, err := ttr.StartKinDCluster(
		ttr.LoadImage(ttr.ImageBuildConfig{
			Dockerfile: "build/images/Dockerfile.sonar-scan",
			ContextDir: rootPath,
		}),
	)
	if err != nil {
		log.Fatal("Could not start KinD cluster: ", err)
	}
	nc, cleanup, err := ttr.SetupTempNamespace(
		cc,
		ott.StartSonarQube(),
		ott.InstallODSPipeline(),
		installSonarQubeConfigMapAndSecret(),
		ttr.InstallTaskFromPath(
			filepath.Join(rootPath, "build/tasks/scan.yaml"),
			nil,
		),
	)
	if err != nil {
		log.Fatal("Could not setup temporary namespace: ", err)
	}
	defer cleanup()
	namespaceConfig = nc
	return m.Run()
}

func runTask(opts ...ttr.TaskRunOpt) error {
	return ttr.RunTask(append([]ttr.TaskRunOpt{
		ttr.InNamespace(namespaceConfig.Name),
		ttr.UsingTask(taskName),
	}, opts...)...)
}

// installSonarQubeConfigMapAndSecret installs the task prerequisites.
func installSonarQubeConfigMapAndSecret() ttr.NamespaceOpt {
	return func(cc *ttr.ClusterConfig, nc *ttr.NamespaceConfig) error {
		k8sClient, err := newKubeClient()
		if err != nil {
			return fmt.Errorf("create K8s client: %s", err)
		}
		url, err := readSonarURL()
		if err != nil {
			return fmt.Errorf("read SonarQube URL: %s", err)
		}
		_, err = k8sClient.CoreV1().ConfigMaps(nc.Name).Create(context.Background(), &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name: "ods-sonar",
			},
			Data: map[string]string{
				"url":     url,
				"edition": "community",
			},
		}, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("create ConfigMap ods-sonar: %s", err)
		}
		username, password, err := readSonarAuth()
		if err != nil {
			return fmt.Errorf("read SonarQube authentication: %s", err)
		}
		_, err = k8sClient.CoreV1().Secrets(nc.Name).Create(context.Background(), &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: "ods-sonar-auth",
			},
			StringData: map[string]string{
				"password": password,
				"username": username,
			},
		}, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("create Secret ods-sonar-auth: %s", err)
		}
		return nil
	}
}

// readSonarAuth reads the Sonar credentials from the central location.
// Eventually, launching the SonarQube service should be moved to this repository instead.
func readSonarAuth() (username, password string, err error) {
	sonarAuth, err := getTrimmedFileContent("/tmp/ods-pipeline/kind-values/sonar-auth")
	if err != nil {
		return
	}
	username, password, found := strings.Cut(sonarAuth, ":")
	if !found {
		err = errors.New("did not find expected sonar auth string")
	}
	return
}

// readSonarURL reads the Sonar URL from the central location.
// Eventually, launching the SonarQube service should be moved to this repository instead.
func readSonarURL() (string, error) {
	return getTrimmedFileContent("/tmp/ods-pipeline/kind-values/sonar-http")
}

func getTrimmedFileContent(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}
