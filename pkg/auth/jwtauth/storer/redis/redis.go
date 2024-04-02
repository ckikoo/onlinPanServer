package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type Config struct {
	Addr      string
	DB        int
	Password  string
	Keyprefix string
}

type Store struct {
	cli    *redis.Client
	prefix string
}

func NewStore(cfg *Config) *Store {
	cli := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := cli.Ping().Err(); err != nil {
		panic(err)
	}
	return &Store{
		cli:    cli,
		prefix: cfg.Keyprefix,
	}
}

func NewStpreWithClient(cli *redis.Client, keyPrefix string) *Store {
	return &Store{
		cli:    cli,
		prefix: keyPrefix,
	}
}

func (s *Store) wrapperKey(key string) string {
	return fmt.Sprintf("%s%s", s.prefix, key)
}

func (s *Store) Set(ctx context.Context, tokenString string, expireation time.Duration) error {
	cmd := s.cli.Set(s.wrapperKey(tokenString), "1", expireation)
	return cmd.Err()
}

func (s *Store) Delete(ctx context.Context, tokenString string) (bool, error) {
	cmd := s.cli.Del(s.wrapperKey(tokenString))
	if err := cmd.Err(); err != nil {
		return false, err
	}

	return true, nil
}

func (s *Store) Check(ctx context.Context, tokenString string) (bool, error) {
	cmd := s.cli.Exists(s.wrapperKey(tokenString))
	if err := cmd.Err(); err != nil {
		return false, err
	}

	return cmd.Val() > 0, nil
}

func (s *Store) Close() error {
	return s.cli.Close()
}
