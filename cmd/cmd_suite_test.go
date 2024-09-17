package cmd_test

import (
	"database/sql"
	"fmt"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/spf13/viper"
	"gourd/internal/common"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	db             *sql.DB
	pool           *dockertest.Pool
	resource       *dockertest.Resource
	cmd            *exec.Cmd
	client         *http.Client
	testConfig     common.Config
	testAdminToken = "127f3733-61cd-4b0a-b8bd-79786501357b"
	testUserToken  = "a9cbd4c0-1296-4290-9239-5fee6b10082e"
)

func waitForServerReady(url string, timeout time.Duration) error {
	timeoutChan := time.After(timeout)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeoutChan:
			return fmt.Errorf("server did not become ready within %v", timeout)
		case <-ticker.C:
			resp, err := http.Get(url)
			if err == nil && resp.StatusCode == http.StatusOK {
				resp.Body.Close()
				return nil
			}
		}
	}
}

var _ = BeforeSuite(func() {
	os.Setenv("RUNNING_TEST", "true")
	viper.SetConfigName("../config_test")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./")
	err := viper.ReadInConfig()
	Expect(err).ToNot(HaveOccurred())
	err = viper.Unmarshal(&testConfig)
	Expect(err).ToNot(HaveOccurred())

	pool, err = dockertest.NewPool("")
	Expect(err).ToNot(HaveOccurred())

	resource, err = pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			fmt.Sprintf("POSTGRES_USER=%s", testConfig.DB.User),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", testConfig.DB.Password),
			fmt.Sprintf("POSTGRES_DB=%s", testConfig.DB.Name),
		},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432/tcp": {{HostIP: "", HostPort: "5433"}},
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	Expect(err).ToNot(HaveOccurred())

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		testConfig.DB.Host, testConfig.DB.Port, testConfig.DB.User, testConfig.DB.Password, testConfig.DB.Name)
	err = pool.Retry(func() error {
		var err error
		db, err = sql.Open("postgres", psqlInfo)
		if err != nil {
			return err
		}
		return db.Ping()
	})
	Expect(err).ToNot(HaveOccurred())

	client = &http.Client{
		Timeout: 5 * time.Second,
	}

	cmd = exec.Command("go", "run", "../main.go", "serve", "--config=../config_test.toml")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	Expect(err).ToNot(HaveOccurred())
	err = waitForServerReady(fmt.Sprintf("http://localhost:%d/", testConfig.ServerPort), 5*time.Second)
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	if cmd != nil && cmd.Process != nil {
		cmd.Process.Kill()
		cmd.Wait()
	}
	err := pool.Purge(resource)
	Expect(err).ToNot(HaveOccurred())

	if db != nil {
		db.Close()
	}
})

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmd Suite")
}
