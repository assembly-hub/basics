// Package redis
package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	v8 "github.com/go-redis/redis/v8"

	"github.com/assembly-hub/basics/util"
)

func SimpleNewRedis() {
	opts := DefaultOptions()
	opts.Addr = "127.0.0.1:6379"
	opts.DB = 0

	r := NewRedis(&opts)
	defer r.Close()

	time.Sleep(5 * time.Second)

	ret := r.Ping(context.Background())

	fmt.Println(ret)
}

func SimpleSet() {
	opts := DefaultOptions()
	opts.Addr = "127.0.0.1:6379"
	opts.DB = 0

	r := NewRedis(&opts)
	defer r.Close()

	err := r.Set(context.Background(), "test_key", "val_123")
	if err != nil {
		panic(err)
	}
	fmt.Println("ok")
}

func SimpleSetEx() {
	opts := DefaultOptions()
	opts.Addr = "127.0.0.1:6379"
	opts.DB = 0

	r := NewRedis(&opts)
	defer r.Close()

	err := r.SetEx(context.Background(), "test_key", "val_123", 10)
	if err != nil {
		panic(err)
	}
	fmt.Println("ok")
}

func SimpleGet() {
	opts := DefaultOptions()
	opts.Addr = "127.0.0.1:6379"
	opts.DB = 0

	r := NewRedis(&opts)
	defer r.Close()

	val, err := r.Get(context.Background(), "test_key")
	if err != nil {
		panic(err)
	}
	fmt.Println("val: ", val)
}

func SimpleDel() {
	opts := DefaultOptions()
	opts.Addr = "127.0.0.1:6379"
	opts.DB = 0

	r := NewRedis(&opts)
	defer r.Close()

	ret, err := r.Del(context.Background(), "test_key", "k2", "k3")
	if err != nil {
		panic(err)
	}

	fmt.Println(ret)
}

func SimpleUnlink() {
	opts := DefaultOptions()
	opts.Addr = "127.0.0.1:6379"
	opts.DB = 0

	r := NewRedis(&opts)
	defer r.Close()

	ret, err := r.Unlink(context.Background(), "test_key")
	if err != nil {
		panic(err)
	}

	fmt.Println(ret)
}

func SimpleExists() {
	opts := DefaultOptions()
	opts.Addr = "127.0.0.1:6379"
	opts.DB = 0

	r := NewRedis(&opts)
	defer r.Close()

	ret, err := r.Exists(context.Background(), "test_key", "123")
	if err != nil {
		panic(err)
	}

	fmt.Println(ret)
}

func SimpleGetLock() {
	opts := DefaultOptions()
	opts.Addr = "127.0.0.1:6379"
	opts.DB = 0

	r := NewRedis(&opts)
	defer r.Close()

	lockKey := "test111"

	err := r.GetLock(lockKey, nil, 100)
	if err != nil {
		panic(err)
	}
	defer func() {
		e := r.FreeLock(lockKey)
		if e != nil {
			log.Println(e)
		}
	}()

	fmt.Println("add lock ok")
}

func SimpleRegister() {
	opts := DefaultOptions()
	opts.Addr = "127.0.0.1:6379"
	opts.DB = 0

	r := NewRedis(&opts)
	defer r.Close()

	lockKey := "test111"
	b := r.Register(lockKey, "1234123", 5)

	fmt.Println("ret； ", b)
}

func SimpleQuotaLimit() {
	opts := DefaultOptions()
	opts.Addr = "127.0.0.1:6379"
	opts.DB = 0

	r := NewRedis(&opts)
	defer r.Close()

	limitLock := "test111"
	_, _ = r.Del(context.Background(), limitLock)

	for {
		b := r.QuotaLimit(limitLock, 20, 1000, 1, true)
		fmt.Println(b)
		if b {
			break
		}
	}
	fmt.Println("ok")
}

func SimpleStreamAdd() {
	opts := DefaultOptions()
	opts.Addr = "127.0.0.1:6379"
	opts.DB = 0

	r := NewRedis(&opts)
	defer r.Close()

	cli, ctx := r.Raw()

	streamKey := "test_stream"

	arg := v8.XAddArgs{
		Stream: streamKey,
		ID:     "*",
		Values: map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		},
	}
	ret := cli.XAdd(ctx, &arg)
	result, err := ret.Result()
	if err != nil {
		panic(err)
	}

	// 消息ID
	fmt.Println(result)
}

func SimpleDelAll() {
	opts := DefaultOptions()
	opts.Addr = "127.0.0.1:6379"
	opts.DB = 0

	r := NewRedis(&opts)
	defer r.Close()

	current := uint64(0)
	for {
		scan, u, err := r.Scan(context.Background(), current, "", 100)
		if err != nil {
			panic(err)
		}

		if len(scan) <= 0 {
			break
		}

		current = u

		_, _ = r.Del(context.Background(), scan...)
	}
}

func SimpleReadStreamsAndAck() {
	opts := DefaultOptions()
	opts.Addr = "127.0.0.1:6379"
	opts.DB = 0

	r := NewRedis(&opts)
	defer r.Close()

	cli, ctx := r.Raw()

	streamKey := "test_stream"
	groupName := "test_group"
	consumer := "test_consumer"

	_, err := cli.XGroupCreate(ctx, streamKey, groupName, "0").Result()
	if err != nil {
		if !util.EndWith(err.Error(), "Group name already exists", true) {
			panic(err)
		}
	}

	cli.XGroupSetID(ctx, streamKey, groupName, "$")

	aaa := v8.XPendingExtArgs{
		Stream:   streamKey,
		Group:    groupName,
		Start:    "-",
		End:      "+",
		Count:    10,
		Consumer: consumer,
	}
	p, err := cli.XPendingExt(ctx, &aaa).Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("pending: ", p)

	innerFunc := func(index int) {
		a := v8.XReadGroupArgs{
			Streams:  []string{streamKey, ">"},
			Consumer: consumer,
			Group:    groupName,
			Count:    1,
			// Block: 100 * time.Millisecond,
			// NoAck: true,
		}

		for {
			ret := cli.XReadGroup(ctx, &a)
			result, err := ret.Result()
			if err != nil {
				panic(err)
			}

			fmt.Println("1111111")

			for _, val := range result {
				// key := val.Stream
				msgList := val.Messages
				for _, msg := range msgList {
					fmt.Println(msg.ID, msg.Values)
					//ret, err := cli.XAck(ctx, key, groupName, msg.ID).Result()
					//if err != nil {
					//	panic(err)
					//}
					//fmt.Println("index: ", index, ", ret: ", ret)
				}
			}
		}
	}

	go innerFunc(1)
	go innerFunc(2)
	go innerFunc(3)

	time.Sleep(100 * time.Second)
}

func SimpleReadStreamsAndAck2() {
	opts := DefaultOptions()
	opts.Addr = "127.0.0.1:6379"
	opts.DB = 0

	r := NewRedis(&opts)
	defer r.Close()

	cli, ctx := r.Raw()

	streamKey := "test_stream"
	groupName := "test_group"
	consumer := "test_consumer"

	_, err := cli.XGroupCreate(ctx, streamKey, groupName, "0").Result()
	if err != nil {
		if !util.EndWith(err.Error(), "Group name already exists", true) {
			panic(err)
		}
	}

	cli.XGroupSetID(ctx, streamKey, groupName, "0")

	aaa := v8.XPendingExtArgs{
		Stream:   streamKey,
		Group:    groupName,
		Start:    "-",
		End:      "+",
		Count:    10,
		Consumer: consumer,
	}
	pendingData, err := cli.XPendingExt(ctx, &aaa).Result()
	if err != nil {
		panic(err)
	}

	for _, val := range pendingData {
		fmt.Println(val.ID, val.Idle)
	}

	a := v8.XReadGroupArgs{
		Streams:  []string{streamKey, ">"},
		Consumer: consumer,
		Group:    groupName,
		Count:    1,
		// Block: 100 * time.Millisecond,
		// NoAck: true,
	}

	innerFunc := func(i int) {
		for {
			ret := cli.XReadGroup(ctx, &a)
			result, err := ret.Result()
			if err != nil {
				panic(err)
			}

			for _, val := range result {
				key := val.Stream
				msgList := val.Messages
				for _, msg := range msgList {
					fmt.Println(key, msg.ID, msg.Values)
					ret, err := cli.XAck(ctx, key, groupName, msg.ID).Result()
					if err != nil {
						panic(err)
					}
					fmt.Println("index: ", i, ", ret: ", ret)
				}
			}
		}
	}

	go innerFunc(1)
	go innerFunc(2)
	go innerFunc(3)

	time.Sleep(100 * time.Second)
}

func SimpleZAdd() {
	opts := DefaultOptions()
	opts.Addr = "127.0.0.1:6379"
	opts.DB = 0

	r := NewRedis(&opts)
	defer r.Close()

	zsetKKey := "zset_key"
	_, _ = r.ZAdd(context.Background(), zsetKKey, "test1", 1)
	_, _ = r.ZAdd(context.Background(), zsetKKey, "test2", 2)
	_, _ = r.ZAdd(context.Background(), zsetKKey, "test3", 3)

	memList, err := r.ZRangeByScore(context.Background(), zsetKKey, "0", "3.01", 0, 10)
	if err != nil {
		panic(err)
	}

	fmt.Println(memList)
}

func SimpleZRangeByScore() {
	opts := DefaultOptions()
	opts.Addr = "127.0.0.1:6379"
	opts.DB = 0

	r := NewRedis(&opts)
	defer r.Close()

	zsetKKey := "zset_key"

	ret, err := r.ZRemRangeByScore(context.Background(), zsetKKey, "0", "2")
	if err != nil {
		panic(err)
	}

	fmt.Println(ret)
}

func SimpleZRem() {
	opts := DefaultOptions()
	opts.Addr = "127.0.0.1:6379"
	opts.DB = 0

	r := NewRedis(&opts)
	defer r.Close()

	zsetKKey := "zset_key"

	ret, err := r.ZRem(context.Background(), zsetKKey, "test3", "test4")
	if err != nil {
		panic(err)
	}

	fmt.Println(ret)
}

func SimpleWithLock() {
	opts := DefaultOptions()
	opts.Addr = "127.0.0.1:6379"
	opts.DB = 0

	r := NewRedis(&opts)
	defer r.Close()

	lockKey := "test_key"
	r.WithLock(lockKey, nil, 10, 3, 500, func() {
		fmt.Println("ok")
	})
}
