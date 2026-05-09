package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"shs/actions"
	"shs/app"
	"shs/config"
	"time"

	"github.com/redis/go-redis/v9"
)

const keyPrefix = "shs:"

const (
	accountSessionTokenTtlDays = 60
	redirectPathTtlMinutes     = 30
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

func accountIdToTokenKey(accountId uint) string {
	return fmt.Sprintf("%saccount-id-to-token:%d", keyPrefix, accountId)
}

func (c *Cache) SetAuthenticatedAccount(sessionToken string, account actions.Account) error {
	accountJson, err := json.Marshal(account)
	if err != nil {
		return err
	}

	err = c.client.Set(context.Background(), accountIdToTokenKey(account.Id), sessionToken, accountSessionTokenTtlDays*time.Hour*24).Err()
	if err != nil {
		return err
	}

	return c.client.Set(context.Background(), accountTokenKey(sessionToken), string(accountJson), accountSessionTokenTtlDays*time.Hour*24).Err()
}

func (c *Cache) GetAuthenticatedAccount(sessionToken string) (actions.Account, error) {
	res := c.client.Get(context.Background(), accountTokenKey(sessionToken))
	if res == nil {
		return actions.Account{}, &app.ErrNotFound{
			ResourceName: "account",
		}
	}
	value, err := res.Result()
	if err == redis.Nil {
		return actions.Account{}, &app.ErrNotFound{
			ResourceName: "account",
		}
	} else if err != nil {
		return actions.Account{}, err
	}

	var account actions.Account
	err = json.Unmarshal([]byte(value), &account)
	if err != nil {
		return actions.Account{}, err
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

func (c *Cache) InvalidateAuthenticatedAccountById(accountId uint) error {
	sessionToken, err := c.client.Get(context.Background(), accountIdToTokenKey(accountId)).Result()
	if err != nil && err != redis.Nil {
		return err
	}

	err = c.client.Del(context.Background(), accountIdToTokenKey(accountId)).Err()
	if err != nil {
		return err
	}

	// ignored in the case of expiration
	_ = c.client.Del(context.Background(), accountTokenKey(sessionToken)).Err()

	return nil
}

func redirectPathKey(clientHash string) string {
	return fmt.Sprintf("%sredirect-path:%s", keyPrefix, clientHash)
}

func (c *Cache) SetRedirectPath(clientHash, path string) error {
	return c.client.Set(context.Background(), redirectPathKey(clientHash), path, redirectPathTtlMinutes*time.Minute).Err()
}

func (c *Cache) GetRedirectPath(clientHash string) (string, error) {
	value, err := c.client.Get(context.Background(), redirectPathKey(clientHash)).Result()
	if err == redis.Nil {
		return "", errors.New("oopsie")
	} else if err != nil {
		return "", err
	}

	return value, nil
}

func (c *Cache) FlushAll() error {
	return c.client.FlushAll(context.Background()).Err()
}
