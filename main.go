package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	log "github.com/sirupsen/logrus"
)

const (
	awsSecret = "SECRET_NAME"
	port      = ":8080"
)

var secrets map[string]string

func getSecret(name string) error {
	sess := session.Must(session.NewSession(&aws.Config{}))

	sm := secretsmanager.New(sess)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(name),
	}

	s, err := sm.GetSecretValue(input)
	if err != nil {
		return err
	}
	log.Debugf("secret value: %s", *s.SecretString)

	var j interface{}
	err = json.Unmarshal([]byte(*s.SecretString), &j)
	if err != nil {
		return err
	}

	for k, v := range j.(map[string]interface{}) {
		log.Debugf("variable %s=%s", k, v)
		secrets[k] = v.(string)
	}
	return nil
}

func getEnv(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/x-sh; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	for k, v := range secrets {
		fmt.Fprintf(w, "export %s=\"%s\"\n", k, v)
	}
}

func getJSON(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/x-sh; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(secrets)
}

func get(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	k := req.URL.Query().Get("key")
	if v, ok := secrets[k]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "no value for given key found\n")
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, v)
	}
}

func main() {
	if os.Getenv("DEBUG") == "true" {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	secrets = make(map[string]string)

	name := os.Getenv(awsSecret)
	if name == "" {
		log.Fatalf("error: no environment variable %s set", awsSecret)
	}

	err := getSecret(name)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	http.HandleFunc("/env", getEnv)
	http.HandleFunc("/json", getJSON)
	http.HandleFunc("/get", get)
	log.Infof("http server is listening on %s", port)
	http.ListenAndServe(port, nil)
}
