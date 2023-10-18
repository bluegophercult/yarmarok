//go:build local

package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"

	_ "github.com/kaznasho/yarmarok"

	"github.com/kaznasho/yarmarok/function"
	"github.com/kaznasho/yarmarok/testinfra/firestore"
)

type testEnv struct {
	cleanups []func()
	tmpDirs  []string
	oldEnvs  map[string]string
}

func (t *testEnv) TempDir() string {
	tmp := filepath.Join(os.TempDir(), "yarmarok", strconv.FormatInt(time.Now().UnixMilli(), 10))
	t.tmpDirs = append(t.tmpDirs, tmp)
	return tmp
}

func (t *testEnv) Cleanup(f func()) {
	t.cleanups = append(t.cleanups, f)
}

func (t *testEnv) Fatalf(format string, args ...any) {
	log.Fatalf(format, args...)
}

func (t *testEnv) Log(args ...any) {
	log.Println(args...)
}

func (t *testEnv) Setenv(key, value string) {
	old, ok := os.LookupEnv(key)
	if ok {
		t.oldEnvs[key] = old
	}

	err := os.Setenv(key, value)
	if err != nil {
		log.Println("Failed to set env:", key, value)
	}
}

func (t *testEnv) runCleanup() {
	for _, c := range t.cleanups {
		c()
	}
	for _, d := range t.tmpDirs {
		err := os.RemoveAll(d)
		if err != nil {
			log.Println("Failed to remove:", d, err)
		}
	}
	for k, v := range t.oldEnvs {
		err := os.Setenv(k, v)
		if err != nil {
			log.Println("Failed to set env back:", k, v, err)
		}
	}
}

func main() {
	t := &testEnv{}

	firestoreInstance, err := firestore.RunInstance(t)
	if err != nil {
		t.Fatalf("Run firestore: %s", err)
	}

	t.Setenv(function.ProjectIDEnvVar, firestoreInstance.ProjectID())
	t.Setenv("FUNCTION_TARGET", "Entrypoint")

	port := "8081"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, os.Kill)

	go func() {
		t.Log("Starting on port", port)
		if err = funcframework.Start(port); err != nil {
			t.Fatalf("Start func: %s", err)
		}
	}()

	<-sigs
	t.runCleanup()
}
