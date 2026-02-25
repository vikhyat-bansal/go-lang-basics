package db

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Config stores PostgreSQL connection settings.
type Config struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

// Client executes SQL commands through psql.
type Client struct {
	cfg Config
}

func NewConfigFromEnv() Config {
	return Config{
		Host: getEnv("DB_HOST", "localhost"),
		Port: getEnv("DB_PORT", "5433"),
		User: getEnv("DB_USER", "local"),
		Pass: getEnv("DB_PASS", "local"),
		Name: getEnv("DB_NAME", "todo"),
	}
}

func NewClient(cfg Config) *Client {
	return &Client{cfg: cfg}
}

func (c *Client) Init() error {
	if _, err := c.QueryValue("SELECT 1;"); err != nil {
		return err
	}
	return c.ensureSchema()
}

func (c *Client) QueryValue(query string) (string, error) {
	cmd := c.command("-c", query)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("psql query failed: %w: %s", err, strings.TrimSpace(string(output)))
	}
	return strings.TrimSpace(string(output)), nil
}

func (c *Client) Exec(query string) error {
	_, err := c.QueryValue(query)
	return err
}

func (c *Client) command(args ...string) *exec.Cmd {
	baseArgs := []string{
		"-X",
		"-q",
		"-t",
		"-A",
		"-h", c.cfg.Host,
		"-p", c.cfg.Port,
		"-U", c.cfg.User,
		"-d", c.cfg.Name,
	}
	baseArgs = append(baseArgs, args...)

	cmd := exec.Command("psql", baseArgs...)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+c.cfg.Pass)
	return cmd
}

func (c *Client) ensureSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS todos (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT NOT NULL DEFAULT '',
		completed BOOLEAN NOT NULL DEFAULT FALSE,
		user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`
	return c.Exec(schema)
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func EscapeLiteral(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}
