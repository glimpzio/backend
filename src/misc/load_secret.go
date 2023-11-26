package misc

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	_ "github.com/lib/pq"
)

// Load a secret into the interface
func LoadSecret(secretName string, store interface{}) error {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return err
	}

	svc := secretsmanager.NewFromConfig(cfg)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := svc.GetSecretValue(context.Background(), input)
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(*result.SecretString), store); err != nil {
		return err
	}

	return nil
}

type databaseSecret struct {
	DbInstanceIdentifier string `json:"dbInstanceIdentifier"`
	Password             string `json:"password"`
	Engine               string `json:"engine"`
	Port                 int    `json:"port"`
	Host                 string `json:"host"`
	Username             string `json:"username"`
	Database             string `json:"database"`
	SSLMode              string `json:"sslmode"`
}

// Load database from secret
func LoadDatabaseFromSecret(secretName string) (*sql.DB, error) {
	dbEnv := &databaseSecret{}
	if err := LoadSecret(secretName, dbEnv); err != nil {
		return nil, err
	}

	dbString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s", dbEnv.Username, dbEnv.Password, dbEnv.Host, dbEnv.Port, dbEnv.Database, dbEnv.SSLMode)
	db, err := sql.Open("postgres", dbString)
	if err != nil {
		return nil, err
	}

	return db, nil
}
