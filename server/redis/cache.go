package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"shs/app"
	"shs/app/models"
	"shs/config"
	"time"

	"github.com/redis/go-redis/v9"
)

const keyPrefix = "shs:"

const (
	accountSessionTokenTtlDays = 60
)

type Cache struct {
	client *redis.Client
}

func New() *Cache {
	return &Cache{
		client: redis.NewClient(&redis.Options{
			Addr:     config.Env().Cache.Host,
			Password: config.Env().Cache.Password,
			DB:       0,
		}),
	}
}

func accountTokenKey(sessionToken string) string {
	return fmt.Sprintf("%saccount-session-token:%s", keyPrefix, sessionToken)
}

func (c *Cache) SetAuthenticatedAccount(sessionToken string, account models.Account) error {
	accountJson, err := json.Marshal(account)
	if err != nil {
		return err
	}

	return c.client.Set(context.Background(), accountTokenKey(sessionToken), string(accountJson), accountSessionTokenTtlDays*time.Hour*24).Err()
}

func (c *Cache) GetAuthenticatedAccount(sessionToken string) (models.Account, error) {
	res := c.client.Get(context.Background(), accountTokenKey(sessionToken))
	if res == nil {
		return models.Account{}, &app.ErrNotFound{
			ResourceName: "account",
		}
	}
	value, err := res.Result()
	if err == redis.Nil {
		return models.Account{}, &app.ErrNotFound{
			ResourceName: "account",
		}
	} else if err != nil {
		return models.Account{}, err
	}

	var account models.Account
	err = json.Unmarshal([]byte(value), &account)
	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}

func (c *Cache) InvalidateAuthenticatedAccount(sessionToken string) error {
	err := c.client.Del(context.Background(), accountTokenKey(sessionToken)).Err()
	if err != nil {
		return err
	}

	return nil
}
