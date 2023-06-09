// Package redis
package redis

import (
	"context"
	"fmt"
	"time"

	v8 "github.com/go-redis/redis/v8"
)

type Redis struct {
	opt *Options
	cli *v8.Client
	ctx context.Context
}

func NewRedis(opt *Options) *Redis {
	return NewRedisWithCtx(context.Background(), opt)
}

func NewRedisWithCtx(ctx context.Context, opt *Options) *Redis {
	redisConn := v8.NewClient(&opt.Options)
	r := new(Redis)
	r.cli = redisConn
	r.opt = opt
	r.ctx = ctx
	return r
}

func FromV8(cli *v8.Client) *Redis {
	r := new(Redis)
	r.cli = cli
	r.opt = &Options{
		*cli.Options(),
	}
	r.ctx = context.Background()
	return r
}

func (r *Redis) RawRedis() *v8.Client {
	return r.cli
}

func (r *Redis) RawCtx() context.Context {
	return r.ctx
}

func (r *Redis) Raw() (*v8.Client, context.Context) {
	return r.cli, r.ctx
}

func (r *Redis) Close() error {
	return r.cli.Close()
}

func (r *Redis) Ping(ctx context.Context) string {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}
	cmd := r.cli.Ping(ctxObj)
	if cmd == nil {
		return ""
	}

	ret, err := cmd.Result()
	if err != nil {
		return ""
	}

	return ret
}

func (r *Redis) Set(ctx context.Context, key, val string) error {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}
	cmd := r.cli.Set(ctxObj, key, val, 0)
	if cmd == nil {
		return fmt.Errorf("redis client error")
	}

	_, err := cmd.Result()
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) SetEx(ctx context.Context, key, val string, expSecond int) error {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.SetEX(ctxObj, key, val, time.Duration(expSecond)*time.Second)
	if cmd == nil {
		return fmt.Errorf("redis client error")
	}

	_, err := cmd.Result()
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) setNX(ctx context.Context, key, val string, exp time.Duration) (bool, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.SetNX(ctxObj, key, val, exp)
	if cmd == nil {
		return false, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return false, err
	}

	return ret, nil
}

func (r *Redis) SetNx(ctx context.Context, key, val string) (bool, error) {
	return r.setNX(ctx, key, val, -1)
}

func (r *Redis) SetNxSec(ctx context.Context, key, val string, expSecond int) (bool, error) {
	return r.setNX(ctx, key, val, time.Second*time.Duration(expSecond))
}

func (r *Redis) SetNxMs(ctx context.Context, key, val string, expMill int) (bool, error) {
	return r.setNX(ctx, key, val, time.Millisecond*time.Duration(expMill))
}

func (r *Redis) setXX(ctx context.Context, key, val string, exp time.Duration) (bool, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.SetXX(ctxObj, key, val, exp)
	if cmd == nil {
		return false, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return false, err
	}

	return ret, nil
}

func (r *Redis) SetXx(ctx context.Context, key, val string) (bool, error) {
	return r.setNX(ctx, key, val, -1)
}

func (r *Redis) SetXxSec(ctx context.Context, key, val string, expSecond int) (bool, error) {
	return r.setXX(ctx, key, val, time.Second*time.Duration(expSecond))
}

func (r *Redis) SetXxMs(ctx context.Context, key, val string, expMill int) (bool, error) {
	return r.setXX(ctx, key, val, time.Millisecond*time.Duration(expMill))
}

func (r *Redis) Del(ctx context.Context, key ...string) (int64, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.Del(ctxObj, key...)
	if cmd == nil {
		return 0, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *Redis) Unlink(ctx context.Context, key ...string) (int64, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.Unlink(ctxObj, key...)
	if cmd == nil {
		return 0, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *Redis) Exists(ctx context.Context, key ...string) (int64, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.Exists(ctxObj, key...)
	if cmd == nil {
		return 0, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.Get(ctxObj, key)
	if cmd == nil {
		return "", fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		if err == v8.Nil {
			return "", nil
		}
		return "", err
	}

	return ret, nil
}

func (r *Redis) SetRange(ctx context.Context, key string, offset int64, value string) (int64, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.SetRange(ctxObj, key, offset, value)
	if cmd == nil {
		return 0, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *Redis) StrLen(ctx context.Context, key string) (int64, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.StrLen(ctxObj, key)
	if cmd == nil {
		return 0, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *Redis) GetBit(ctx context.Context, key string, offset int64) (int64, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.GetBit(ctxObj, key, offset)
	if cmd == nil {
		return 0, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *Redis) SetBit(ctx context.Context, key string, offset int64, value int) (int64, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.SetBit(ctxObj, key, offset, value)
	if cmd == nil {
		return 0, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *Redis) BitCount(ctx context.Context, key string, bitStart, bitEnd int64) (int64, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.BitCount(ctxObj, key, &v8.BitCount{
		Start: bitStart,
		End:   bitEnd,
	})
	if cmd == nil {
		return 0, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *Redis) BitCountAll(ctx context.Context, key string) (int64, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.BitCount(ctxObj, key, nil)
	if cmd == nil {
		return 0, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *Redis) BitOpAnd(ctx context.Context, destKey string, keys ...string) (int64, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.BitOpAnd(ctxObj, destKey, keys...)
	if cmd == nil {
		return 0, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *Redis) BitOpOr(ctx context.Context, destKey string, keys ...string) (int64, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.BitOpOr(ctxObj, destKey, keys...)
	if cmd == nil {
		return 0, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *Redis) BitOpXor(ctx context.Context, destKey string, keys ...string) (int64, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.BitOpXor(ctxObj, destKey, keys...)
	if cmd == nil {
		return 0, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *Redis) BitOpNot(ctx context.Context, destKey string, key string) (int64, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.BitOpNot(ctxObj, destKey, key)
	if cmd == nil {
		return 0, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *Redis) BitPos(ctx context.Context, key string, bit int64, pos ...int64) (int64, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.BitPos(ctxObj, key, bit, pos...)
	if cmd == nil {
		return 0, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *Redis) BitField(ctx context.Context, key string, args ...interface{}) ([]int64, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.BitField(ctxObj, key, args...)
	if cmd == nil {
		return nil, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *Redis) Scan(ctx context.Context, cursor uint64, match string, count int64) (keys []string, cur uint64, err error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.Scan(ctxObj, cursor, match, count)
	if cmd == nil {
		return nil, 0, fmt.Errorf("redis client error")
	}

	result, cur, err := cmd.Result()
	if err != nil {
		return nil, 0, err
	}

	return result, cur, nil
}

func (r *Redis) ScanType(ctx context.Context, cursor uint64, match string, count int64, keyType string) (keys []string, cur uint64, err error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.ScanType(ctxObj, cursor, match, count, keyType)
	if cmd == nil {
		return nil, 0, fmt.Errorf("redis client error")
	}

	result, cur, err := cmd.Result()
	if err != nil {
		return nil, 0, err
	}

	return result, cur, nil
}

func (r *Redis) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) (keys []string, cur uint64, err error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.SScan(ctxObj, key, cursor, match, count)
	if cmd == nil {
		return nil, 0, fmt.Errorf("redis client error")
	}

	result, cur, err := cmd.Result()
	if err != nil {
		return nil, 0, err
	}

	return result, cur, nil
}

func (r *Redis) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) (keys []string, cur uint64, err error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.HScan(ctxObj, key, cursor, match, count)
	if cmd == nil {
		return nil, 0, fmt.Errorf("redis client error")
	}

	result, cur, err := cmd.Result()
	if err != nil {
		return nil, 0, err
	}

	return result, cur, nil
}

func (r *Redis) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) (keys []string, cur uint64, err error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.ZScan(ctxObj, key, cursor, match, count)
	if cmd == nil {
		return nil, 0, fmt.Errorf("redis client error")
	}

	result, cur, err := cmd.Result()
	if err != nil {
		return nil, 0, err
	}

	return result, cur, nil
}

func (r *Redis) PTTL(ctx context.Context, key string) (int64, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.PTTL(ctxObj, key)
	if cmd == nil {
		return 0, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return 0, err
	}

	return int64(ret / time.Millisecond), nil
}

func (r *Redis) TTL(ctx context.Context, key string) (int64, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.TTL(ctxObj, key)
	if cmd == nil {
		return 0, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return 0, err
	}

	return int64(ret / time.Second), nil
}

func (r *Redis) SAdd(ctx context.Context, key string, members ...interface{}) (int64, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.SAdd(ctxObj, key, members...)
	if cmd == nil {
		return 0, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *Redis) SCard(ctx context.Context, key string) (int64, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.SCard(ctxObj, key)
	if cmd == nil {
		return 0, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *Redis) SDiff(ctx context.Context, keys ...string) ([]string, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.SDiff(ctxObj, keys...)
	if cmd == nil {
		return nil, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *Redis) SDiffStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.SDiffStore(ctxObj, destination, keys...)
	if cmd == nil {
		return 0, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *Redis) SInter(ctx context.Context, keys ...string) ([]string, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.SInter(ctxObj, keys...)
	if cmd == nil {
		return nil, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *Redis) SInterStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.SInterStore(ctxObj, destination, keys...)
	if cmd == nil {
		return 0, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *Redis) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.SIsMember(ctxObj, key, member)
	if cmd == nil {
		return false, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return false, err
	}

	return ret, nil
}

func (r *Redis) SMIsMember(ctx context.Context, key string, members ...interface{}) ([]bool, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.SMIsMember(ctxObj, key, members...)
	if cmd == nil {
		return nil, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *Redis) SMembers(ctx context.Context, key string) ([]string, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.SMembers(ctxObj, key)
	if cmd == nil {
		return nil, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *Redis) SMembersMap(ctx context.Context, key string) (map[string]struct{}, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.SMembersMap(ctxObj, key)
	if cmd == nil {
		return nil, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *Redis) SMove(ctx context.Context, source, destination string, member interface{}) (bool, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.SMove(ctxObj, source, destination, member)
	if cmd == nil {
		return false, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return false, err
	}

	return ret, nil
}

func (r *Redis) SPop(ctx context.Context, key string) (string, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.SPop(ctxObj, key)
	if cmd == nil {
		return "", fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return "", err
	}

	return ret, nil
}

func (r *Redis) SPopN(ctx context.Context, key string, count int64) ([]string, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.SPopN(ctxObj, key, count)
	if cmd == nil {
		return nil, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *Redis) SRandMember(ctx context.Context, key string) (string, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}

	cmd := r.cli.SRandMember(ctxObj, key)
	if cmd == nil {
		return "", fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return "", err
	}

	return ret, nil
}

func (r *Redis) SRandMemberN(ctx context.Context, key string, count int64) ([]string, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}
	cmd := r.cli.SRandMemberN(ctxObj, key, count)
	if cmd == nil {
		return nil, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *Redis) SRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}
	cmd := r.cli.SRem(ctxObj, key, members)
	if cmd == nil {
		return 0, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *Redis) SUnion(ctx context.Context, keys ...string) ([]string, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}
	cmd := r.cli.SUnion(ctxObj, keys...)
	if cmd == nil {
		return nil, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *Redis) SUnionStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}
	cmd := r.cli.SUnionStore(ctxObj, destination, keys...)
	if cmd == nil {
		return 0, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return 0, err
	}

	return ret, nil
}

// zset funcs
func (r *Redis) ZAdd(ctx context.Context, key string, member interface{}, score float64) (int64, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}
	z := v8.Z{
		Member: member,
		Score:  score,
	}
	cmd := r.cli.ZAdd(ctxObj, key, &z)
	if cmd == nil {
		return 0, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return 0, err
	}

	return ret, nil
}

// ZAddList memList: {member: m1, score: 1}
func (r *Redis) ZAddList(ctx context.Context, key string, memList ...map[string]interface{}) (int64, error) {
	count := len(memList)
	if count <= 0 {
		return 0, nil
	}

	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}
	zList := make([]*v8.Z, count)
	for i, mem := range memList {
		zList[i] = &v8.Z{
			Member: mem["member"],
			Score:  mem["score"].(float64),
		}
	}
	cmd := r.cli.ZAdd(ctxObj, key, zList...)
	if cmd == nil {
		return 0, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *Redis) ZRem(ctx context.Context, key string, memberList ...interface{}) (int64, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}
	cmd := r.cli.ZRem(ctxObj, key, memberList...)
	if cmd == nil {
		return 0, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *Redis) ZRemRangeByScore(ctx context.Context, key, min, max string) (int64, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}
	cmd := r.cli.ZRemRangeByScore(ctxObj, key, min, max)
	if cmd == nil {
		return 0, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *Redis) ZRangeByScore(ctx context.Context, key, min, max string, offset, count int64) ([]string, error) {
	ctxObj := r.ctx
	if ctx != nil {
		ctxObj = ctx
	}
	opt := v8.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: offset,
		Count:  count,
	}
	cmd := r.cli.ZRangeByScore(ctxObj, key, &opt)
	if cmd == nil {
		return nil, fmt.Errorf("redis client error")
	}

	ret, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	return ret, nil
}
