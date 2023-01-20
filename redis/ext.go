// Package redis tool
package redis

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/assembly-hub/basics/util"
)

// GetLock 获取redis锁，成功返回nil，否则返回对应error
func (r *Redis) GetLock(key string, lockVal *string, exp int) error {
	val := ""
	if lockVal != nil {
		val = *lockVal
	} else {
		val = fmt.Sprintf("%v-%08v", time.Now().UnixNano(), rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(100000000))
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()

	ok, err := r.SetNxSec(ctx, key, val, exp)
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("add lock is failed")
	}
	return nil
}

// TryLock 尝试获取redis锁，指定重试次数和重试间隔，成功返回nil，否则返回对应error
func (r *Redis) TryLock(key string, lockVal *string, exp int, maxTryTime int, intervalMs int64) error {
	val := ""
	if lockVal != nil {
		val = *lockVal
	} else {
		val = fmt.Sprintf("%v-%08v", time.Now().UnixNano(), rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(100000000))
	}

	if maxTryTime < 0 {
		maxTryTime = 0
	}

	if intervalMs < 10 {
		intervalMs = 100
	}

	err := r.GetLock(key, &val, exp)
	if err != nil && maxTryTime > 0 {
		for i := 0; i < maxTryTime; i++ {
			time.Sleep(time.Duration(intervalMs) * time.Millisecond)
			err = r.GetLock(key, &val, exp)
			if err == nil {
				break
			}
		}
	}

	return err
}

// WithLock 尝试获取redis锁，指定重试次数和重试间隔，获取成功之后执行 fun，否则不执行
func (r *Redis) WithLock(key string, lockVal *string, exp int, maxTryTime int, intervalMs int64, fun func()) {
	err := r.TryLock(key, lockVal, exp, maxTryTime, intervalMs)
	if err != nil {
		panic(err)
	}
	defer func() {
		e := r.FreeLock(key)
		if e != nil {
			log.Println(e)
		}
	}()

	fun()
}

// FreeLock 释放指定的redis锁
func (r *Redis) FreeLock(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()

	_, err := r.Del(ctx, key)
	if err != nil {
		return err
	}

	return nil
}

// Register 注册HA 主服务标识，注册成功返回true，否则返回false
func (r *Redis) Register(key string, val string, exp int) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	_, err := r.SetNxSec(ctx, key, val, exp)
	cancel()
	if err != nil {
		log.Println(err)
		return false
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	redisVal, err := r.Get(ctx, key)
	cancel()
	if err != nil {
		log.Println(err)
		return false
	}

	if redisVal == val {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
		_, err = r.SetXxSec(ctx, key, val, exp)
		cancel()
		if err != nil {
			log.Println(err)
		}
		return true
	}
	return false
}

// QuotaLimit 资源量限制，被限制返回true，否则返回false
func (r *Redis) QuotaLimit(limitKay string, maxCount int, exp int, addCount int, addLock bool) bool {
	if maxCount < addCount {
		return true
	}

	innerFunc := func(red *Redis) bool {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
		keyTime, err := red.TTL(ctx, limitKay)
		cancel()
		if err != nil {
			log.Println(err)
			return true
		}

		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
		limitCount, err := red.Get(ctx, limitKay)
		cancel()
		if err != nil {
			log.Println(err)
			return true
		}

		if limitCount != "" && keyTime > 0 {
			lc, err := util.Str2Int[int](limitCount)
			if err != nil {
				log.Println(err)
				return true
			}
			if lc >= maxCount {
				return true
			}

			if keyTime > 0 {
				exp = int(keyTime)
			}
			addCount += lc
		}

		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
		err = red.SetEx(ctx, limitKay, util.IntToStr(addCount), exp)
		cancel()
		if err != nil {
			log.Println(err)
			return true
		}

		return false
	}

	if addLock {
		lockKey := "redis_lock_" + limitKay
		err := r.TryLock(lockKey, nil, 3, 5, 200)
		if err != nil {
			return true
		}
		defer func() {
			e := r.FreeLock(lockKey)
			if e != nil {
				log.Println(e)
			}
		}()

		return innerFunc(r)
	}
	return innerFunc(r)
}
