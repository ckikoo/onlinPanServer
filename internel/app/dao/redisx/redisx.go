package redisx

import (
	"context"
	"fmt"
	"onlineCLoud/internel/app/config"
	"onlineCLoud/pkg/contextx"

	"runtime"
	"time"

	"github.com/go-redis/redis"
)

var (
	lock = "if redis.call('exists',KEYS[1])==0 or redis.call('hexists',KEYS[1],ARGV[1])==1 then redis.call('hincrby',KEYS[1],ARGV[1],1) redis.call('expire',KEYS[1],ARGV[2]) return 1 else return 0 end"
)
var (
	unlock = "if redis.call('hexists',KEYS[1],ARGV[1])==0 then return -1 elseif redis.call('hincrby',KEYS[1],ARGV[1],-1)==0 then return redis.call('del',KEYS[1]) else return 0 end"
)

var (
	keepAlive = "if redis.call('hexists',KEYS[1],ARGV[1])==1 then return redis.call('expire',KEYS[1],ARGV[2]) else return 0 end"
)

var (
	DownLoadkeepAlive = "if redis.call('exists',KEYS[1],ARGV[1])==1 then return redis.call('expire',KEYS[1],ARGV[2]) else return 0 end"
)

type Store interface {
	Set(ctx context.Context, key string, val string, expired time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Check(ctx context.Context, key string) (bool, error)
	Delete(ctx context.Context, key string, val string, expired time.Duration) error
}

type Locker interface {
	Lock(ctx context.Context, key string) error
	Unlock(ctx context.Context, key string) error
}

type Redisx struct {
	cli     *redis.Client
	expires int64
	cancel  context.CancelFunc
	maxWait time.Duration
}

func (cli *Redisx) KeepAliveWorker(ctx context.Context, key string, val string, exp time.Duration) {
	waiter := time.Tick(exp / 2)
	for {
		select {
		case <-waiter:
			res := cli.cli.Eval(DownLoadkeepAlive, []string{key}, val, exp)
			if res.Err() != nil {
				fmt.Printf("res: %v\n", res.Err())
				panic(res)
			}

			if res.Val().(int64) == 0 {
				return
			} else {
			}
		case <-ctx.Done():

			return
		}

	}
}

func (cli *Redisx) Lock(ctx context.Context, key string) error {
	var maxBackoff = 16
	backoff := 1
	for {
		res := cli.cli.Eval(lock, []string{key}, contextx.FromUUID(ctx), 10)
		if res.Err() != nil {
			return res.Err()
		}
		val, ok := res.Val().(int64)
		if ok && val == 1 {
			if cli.cancel == nil {
				ctxx, canc := context.WithCancel(ctx)
				cli.cancel = canc
				go cli.KeepAliveWorker(ctxx, key, contextx.FromUUID(ctxx), 5)
			}
			return nil
		}
		for i := 0; i < backoff; i++ {
			runtime.Gosched()
		}
		if backoff < maxBackoff {
			backoff <<= 1
		}
	}
}

func (cli *Redisx) Unlock(ctx context.Context, key string) error {
	for {
		res := cli.cli.Eval(unlock, []string{key}, contextx.FromUUID(ctx))
		if res.Err() != nil {
			fmt.Println(res.Err())
			return res.Err()
		}
		fmt.Printf("释放锁: %v 结构%v\n", contextx.FromUUID(ctx), res.Err())
		val := res.Val().(int64)
		if val == 1 {
			fmt.Printf("释放锁成功%s\n", contextx.FromUUID(ctx))
			cli.cancel()
			return nil
		} else {
			fmt.Printf("还没解锁完%s", contextx.FromUUID(ctx))
			continue
		}

	}
}
func NewClient() *Redisx {
	rd := &Redisx{
		cli: redis.NewClient(&redis.Options{
			DB:       0,
			Addr:     config.C.Redis.Addr,
			Password: config.C.Redis.Password,
		}),
	}
	if err := rd.cli.Ping().Err(); err != nil {
		panic(err)
	}

	return rd
}

func (r *Redisx) Set(ctx context.Context, key string, val interface{}, expired time.Duration) error {
	cmd := r.cli.Set(key, val, expired)
	return cmd.Err()
}

func (r *Redisx) ZsetWithTimestamps(ctx context.Context, key string, members []string, batchSize int, expire time.Duration) error {
	timestamp := time.Now().Unix()

	for i := 0; i < len(members); i += batchSize {
		end := i + batchSize
		if end > len(members) {
			end = len(members)
		}

		batchMembers := members[i:end]
		// 准备 ZADD 命令的参数
		zMembers := make([]redis.Z, len(batchMembers))
		for j, member := range batchMembers {
			zMembers[j] = redis.Z{Score: float64(timestamp), Member: member}
		}

		cmd := r.cli.ZAdd(key, zMembers...)

		if cmd.Err() != nil {
			return cmd.Err()
		}
	}
	// 设置键的过期时间
	cmd := r.cli.Expire(key, expire)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (r *Redisx) GetRedisByPage(ctx context.Context, key string, page, pageSize int64) ([]string, error) {
	start := int64((page - 1) * pageSize)
	end := int64(page*pageSize - 1)

	cmd := r.cli.ZRange("products", start, end)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	return cmd.Val(), cmd.Err()
}
func (r *Redisx) Get(ctx context.Context, key string) (string, error) {
	str, err := r.cli.Get(key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return str, err
}
func (r *Redisx) Check(ctx context.Context, key string) (bool, error) {
	val, err := r.cli.Exists(key).Result()

	if err != nil {
		return false, err
	}
	return val > 0, nil

}
func (r *Redisx) Delete(ctx context.Context, key string) (bool, error) {
	cmd := r.cli.Del(key)
	if err := cmd.Err(); err != nil {
		return false, err
	}

	return true, nil
}

func (r *Redisx) HMGet(ctx context.Context, key string, fields ...string) ([]interface{}, error) {
	result, err := r.cli.HMGet(key, fields...).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	return result, nil
}
func (r *Redisx) HGet(ctx context.Context, key string, fields string) (string, error) {
	return r.cli.HGet(key, fields).Result()

}
func (r *Redisx) HCheck(ctx context.Context, key string, fields string) (bool, error) {
	result, err := r.cli.HExists(key, fields).Result()
	if err != nil {
		fmt.Printf("redisx err: %v\n", err)
		return false, err
	}

	return result, err

}

func (r *Redisx) Zrem(ctx context.Context, key string, fields []string) (int64, error) {

	return r.cli.ZRem(key, fields).Result()
}

func (r *Redisx) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.cli.HGetAll(key).Result()
}
func (r *Redisx) HGetWithFile(ctx context.Context, key string, file ...string) ([]interface{}, error) {
	return r.cli.HMGet(key, file...).Result()
}

func (r *Redisx) Hexists(ctx context.Context, key string, field string) (bool, error) {
	return r.cli.HExists(key, field).Result()
}

func (r *Redisx) HMset(ctx context.Context, key string, m map[string]interface{}) (string, error) {
	return r.cli.HMSet(key, m).Result()
}

func (r *Redisx) HSet(ctx context.Context, key string, field string, val interface{}) (bool, error) {
	return r.cli.HSet(key, field, val).Result()
}

func (r *Redisx) HDel(ctx context.Context, key string, field []string) (int64, error) {
	return r.cli.HDel(key, field...).Result()
}

func (r *Redisx) Close() error {
	return r.cli.Close()
}

func NewClientWithClient(ctx context.Context, rd *redis.Client) *Redisx {
	return &Redisx{
		cli:     rd,
		maxWait: 0,
	}
}
