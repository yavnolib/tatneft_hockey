package postgre

import (
	"context"
	"errors"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/url"
	"time"
)

type postgresConfig struct {
	Host            string        `env:"DB_HOST"`
	Port            string        `env:"DB_PORT"`
	User            string        `env:"DB_USER"`
	Password        string        `env:"DB_PASSWORD"`
	DBName          string        `env:"DB_NAME"`
	MaxConnLifetime time.Duration `env:"DB_MAX_CONN_LIFETIME"`
	SSLMode         string        `env:"DB_SSL_MODE"`
	Timeout         time.Duration `yaml:"timeout" env-default:"5"`
}

func (cfg *postgresConfig) toConnectionString() string {
	return fmt.Sprintf("%s://%s:%s@web-pgdb-1/%s?sslmode=disable&connect_timeout=%d",
		"postgres",
		url.QueryEscape(cfg.User),
		url.QueryEscape(cfg.Password),
		//cfg.Host,
		//cfg.Port,
		cfg.DBName,
		cfg.Timeout,
	)
}

// TODO свои ошибки

func loadPostgresConfig() (*postgresConfig, error) {
	log.Println("loading DB config")

	cfg := &postgresConfig{}

	if err := env.Parse(cfg); err != nil {
		log.Printf("[ERROR] %+v\n", err)
		return nil, err
	}

	t := postgresConfig{}
	if *cfg == t {
		log.Printf("[ERROR] %+v\n", cfg)
		return nil, errors.New("EMPTY config")
	}
	log.Printf("[INFO] setup DB config: \n")

	return cfg, nil
}

func LoadPgxPool() (*pgxpool.Pool, error) {
	cfg, err := loadPostgresConfig()
	if err != nil {
		log.Printf("[ERROR] %+v\n", err)
		return nil, err
	}
	poolCfg, err := pgxpool.ParseConfig(cfg.toConnectionString())
	if err != nil {
		log.Printf("[ERROR] %+v\n", err)
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		log.Printf("[ERROR] %+v\n", err)
		return nil, err
	}
	if err = TestPing(pool); err != nil {
		return nil, err
	}
	log.Printf("[INFO] setuo pg pool config \n")
	return pool, nil
}

func TestPing(pool *pgxpool.Pool) error {
	ctx := context.Background()
	if err := pool.Ping(ctx); err != nil {
		log.Println("[ERROR] ping not ok")
		return err
	}
	log.Println("[INFO] ping ok")
	return nil
}
