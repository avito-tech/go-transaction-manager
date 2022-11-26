package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// readonlyFuncWithoutTxDecorator calls readonly commands outside of a transaction.
type readonlyFuncWithoutTxDecorator struct {
	Cmdable
	readOnly redis.Cmdable
}

// ReadOnlyFuncWithoutTxDecorator is decorator, which calls readonly commands outside of the Transaction.
//
//nolint:ireturn,nolintlint
func ReadOnlyFuncWithoutTxDecorator(write Cmdable, readOnly redis.Cmdable) Cmdable {
	return &readonlyFuncWithoutTxDecorator{Cmdable: write, readOnly: readOnly}
}

func (r *readonlyFuncWithoutTxDecorator) Dump(ctx context.Context, key string) *redis.StringCmd {
	return r.readOnly.Dump(ctx, key)
}

func (r *readonlyFuncWithoutTxDecorator) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	return r.readOnly.Exists(ctx, keys...)
}

func (r *readonlyFuncWithoutTxDecorator) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {
	return r.readOnly.Keys(ctx, pattern)
}

func (r *readonlyFuncWithoutTxDecorator) PTTL(ctx context.Context, key string) *redis.DurationCmd {
	return r.readOnly.PTTL(ctx, key)
}

func (r *readonlyFuncWithoutTxDecorator) RandomKey(ctx context.Context) *redis.StringCmd {
	return r.readOnly.RandomKey(ctx)
}

func (r *readonlyFuncWithoutTxDecorator) Touch(ctx context.Context, keys ...string) *redis.IntCmd {
	return r.readOnly.Touch(ctx, keys...)
}

func (r *readonlyFuncWithoutTxDecorator) TTL(ctx context.Context, key string) *redis.DurationCmd {
	return r.readOnly.TTL(ctx, key)
}

func (r *readonlyFuncWithoutTxDecorator) Type(ctx context.Context, key string) *redis.StatusCmd {
	return r.readOnly.Type(ctx, key)
}

func (r *readonlyFuncWithoutTxDecorator) Get(ctx context.Context, key string) *redis.StringCmd {
	return r.readOnly.Get(ctx, key)
}

func (r *readonlyFuncWithoutTxDecorator) GetRange(ctx context.Context, key string, start, end int64) *redis.StringCmd {
	return r.readOnly.GetRange(ctx, key, start, end)
}

func (r *readonlyFuncWithoutTxDecorator) MGet(ctx context.Context, keys ...string) *redis.SliceCmd {
	return r.readOnly.MGet(ctx, keys...)
}

func (r *readonlyFuncWithoutTxDecorator) StrLen(ctx context.Context, key string) *redis.IntCmd {
	return r.readOnly.StrLen(ctx, key)
}

func (r *readonlyFuncWithoutTxDecorator) GetBit(ctx context.Context, key string, offset int64) *redis.IntCmd {
	return r.readOnly.GetBit(ctx, key, offset)
}

func (r *readonlyFuncWithoutTxDecorator) BitCount(ctx context.Context, key string, bitCount *redis.BitCount) *redis.IntCmd {
	return r.readOnly.BitCount(ctx, key, bitCount)
}

func (r *readonlyFuncWithoutTxDecorator) BitPos(ctx context.Context, key string, bit int64, pos ...int64) *redis.IntCmd {
	return r.readOnly.BitPos(ctx, key, bit, pos...)
}

func (r *readonlyFuncWithoutTxDecorator) Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd {
	return r.readOnly.Scan(ctx, cursor, match, count)
}

func (r *readonlyFuncWithoutTxDecorator) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return r.readOnly.SScan(ctx, key, cursor, match, count)
}

func (r *readonlyFuncWithoutTxDecorator) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return r.readOnly.HScan(ctx, key, cursor, match, count)
}

func (r *readonlyFuncWithoutTxDecorator) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return r.readOnly.ZScan(ctx, key, cursor, match, count)
}

func (r *readonlyFuncWithoutTxDecorator) HExists(ctx context.Context, key, field string) *redis.BoolCmd {
	return r.readOnly.HExists(ctx, key, field)
}

func (r *readonlyFuncWithoutTxDecorator) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	return r.readOnly.HGet(ctx, key, field)
}

func (r *readonlyFuncWithoutTxDecorator) HGetAll(ctx context.Context, key string) *redis.StringStringMapCmd {
	return r.readOnly.HGetAll(ctx, key)
}

func (r *readonlyFuncWithoutTxDecorator) HKeys(ctx context.Context, key string) *redis.StringSliceCmd {
	return r.readOnly.HKeys(ctx, key)
}

func (r *readonlyFuncWithoutTxDecorator) HLen(ctx context.Context, key string) *redis.IntCmd {
	return r.readOnly.HLen(ctx, key)
}

func (r *readonlyFuncWithoutTxDecorator) HMGet(ctx context.Context, key string, fields ...string) *redis.SliceCmd {
	return r.readOnly.HMGet(ctx, key, fields...)
}

func (r *readonlyFuncWithoutTxDecorator) HVals(ctx context.Context, key string) *redis.StringSliceCmd {
	return r.readOnly.HVals(ctx, key)
}

func (r *readonlyFuncWithoutTxDecorator) HRandField(ctx context.Context, key string, count int, withValues bool) *redis.StringSliceCmd {
	return r.readOnly.HRandField(ctx, key, count, withValues)
}

func (r *readonlyFuncWithoutTxDecorator) LIndex(ctx context.Context, key string, index int64) *redis.StringCmd {
	return r.readOnly.LIndex(ctx, key, index)
}

func (r *readonlyFuncWithoutTxDecorator) LLen(ctx context.Context, key string) *redis.IntCmd {
	return r.readOnly.LLen(ctx, key)
}

func (r *readonlyFuncWithoutTxDecorator) LPos(ctx context.Context, key string, value string, args redis.LPosArgs) *redis.IntCmd {
	return r.readOnly.LPos(ctx, key, value, args)
}

func (r *readonlyFuncWithoutTxDecorator) LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	return r.readOnly.LRange(ctx, key, start, stop)
}

func (r *readonlyFuncWithoutTxDecorator) SCard(ctx context.Context, key string) *redis.IntCmd {
	return r.readOnly.SCard(ctx, key)
}

func (r *readonlyFuncWithoutTxDecorator) SDiff(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	return r.readOnly.SDiff(ctx, keys...)
}

func (r *readonlyFuncWithoutTxDecorator) SInter(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	return r.readOnly.SInter(ctx, keys...)
}

func (r *readonlyFuncWithoutTxDecorator) SIsMember(ctx context.Context, key string, member interface{}) *redis.BoolCmd {
	return r.readOnly.SIsMember(ctx, key, member)
}

func (r *readonlyFuncWithoutTxDecorator) SMIsMember(ctx context.Context, key string, members ...interface{}) *redis.BoolSliceCmd {
	return r.readOnly.SMIsMember(ctx, key, members...)
}

func (r *readonlyFuncWithoutTxDecorator) SMembers(ctx context.Context, key string) *redis.StringSliceCmd {
	return r.readOnly.SMembers(ctx, key)
}

func (r *readonlyFuncWithoutTxDecorator) SRandMember(ctx context.Context, key string) *redis.StringCmd {
	return r.readOnly.SRandMember(ctx, key)
}

func (r *readonlyFuncWithoutTxDecorator) SUnion(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	return r.readOnly.SUnion(ctx, keys...)
}

func (r *readonlyFuncWithoutTxDecorator) XLen(ctx context.Context, stream string) *redis.IntCmd {
	return r.readOnly.XLen(ctx, stream)
}

func (r *readonlyFuncWithoutTxDecorator) XRange(ctx context.Context, stream, start, stop string) *redis.XMessageSliceCmd {
	return r.readOnly.XRange(ctx, stream, start, stop)
}

func (r *readonlyFuncWithoutTxDecorator) XRevRange(ctx context.Context, stream string, start, stop string) *redis.XMessageSliceCmd {
	return r.readOnly.XRevRange(ctx, stream, start, stop)
}

func (r *readonlyFuncWithoutTxDecorator) XRead(ctx context.Context, a *redis.XReadArgs) *redis.XStreamSliceCmd {
	return r.readOnly.XRead(ctx, a)
}

func (r *readonlyFuncWithoutTxDecorator) XPending(ctx context.Context, stream, group string) *redis.XPendingCmd {
	return r.readOnly.XPending(ctx, stream, group)
}

func (r *readonlyFuncWithoutTxDecorator) ZCard(ctx context.Context, key string) *redis.IntCmd {
	return r.readOnly.ZCard(ctx, key)
}

func (r *readonlyFuncWithoutTxDecorator) ZCount(ctx context.Context, key, min, max string) *redis.IntCmd {
	return r.readOnly.ZCount(ctx, key, min, max)
}

func (r *readonlyFuncWithoutTxDecorator) ZLexCount(ctx context.Context, key, min, max string) *redis.IntCmd {
	return r.readOnly.ZLexCount(ctx, key, min, max)
}

func (r *readonlyFuncWithoutTxDecorator) ZInter(ctx context.Context, store *redis.ZStore) *redis.StringSliceCmd {
	return r.readOnly.ZInter(ctx, store)
}

func (r *readonlyFuncWithoutTxDecorator) ZMScore(ctx context.Context, key string, members ...string) *redis.FloatSliceCmd {
	return r.readOnly.ZMScore(ctx, key, members...)
}

func (r *readonlyFuncWithoutTxDecorator) ZRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	return r.readOnly.ZRange(ctx, key, start, stop)
}

func (r *readonlyFuncWithoutTxDecorator) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return r.readOnly.ZRangeByScore(ctx, key, opt)
}

func (r *readonlyFuncWithoutTxDecorator) ZRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return r.readOnly.ZRangeByLex(ctx, key, opt)
}

func (r *readonlyFuncWithoutTxDecorator) ZRank(ctx context.Context, key, member string) *redis.IntCmd {
	return r.readOnly.ZRank(ctx, key, member)
}

func (r *readonlyFuncWithoutTxDecorator) ZRevRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	return r.readOnly.ZRevRange(ctx, key, start, stop)
}

func (r *readonlyFuncWithoutTxDecorator) ZRevRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return r.readOnly.ZRevRangeByScore(ctx, key, opt)
}

func (r *readonlyFuncWithoutTxDecorator) ZRevRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return r.readOnly.ZRevRangeByLex(ctx, key, opt)
}

func (r *readonlyFuncWithoutTxDecorator) ZRevRank(ctx context.Context, key, member string) *redis.IntCmd {
	return r.readOnly.ZRevRank(ctx, key, member)
}

func (r *readonlyFuncWithoutTxDecorator) ZScore(ctx context.Context, key, member string) *redis.FloatCmd {
	return r.readOnly.ZScore(ctx, key, member)
}

func (r *readonlyFuncWithoutTxDecorator) ZUnion(ctx context.Context, store redis.ZStore) *redis.StringSliceCmd {
	return r.readOnly.ZUnion(ctx, store)
}

func (r *readonlyFuncWithoutTxDecorator) ZRandMember(ctx context.Context, key string, count int, withScores bool) *redis.StringSliceCmd {
	return r.readOnly.ZRandMember(ctx, key, count, withScores)
}

func (r *readonlyFuncWithoutTxDecorator) ZDiff(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	return r.readOnly.ZDiff(ctx, keys...)
}

func (r *readonlyFuncWithoutTxDecorator) PFCount(ctx context.Context, keys ...string) *redis.IntCmd {
	return r.readOnly.PFCount(ctx, keys...)
}

func (r *readonlyFuncWithoutTxDecorator) DBSize(ctx context.Context) *redis.IntCmd {
	return r.readOnly.DBSize(ctx)
}

func (r *readonlyFuncWithoutTxDecorator) GeoPos(ctx context.Context, key string, members ...string) *redis.GeoPosCmd {
	return r.readOnly.GeoPos(ctx, key, members...)
}

func (r *readonlyFuncWithoutTxDecorator) GeoSearch(ctx context.Context, key string, q *redis.GeoSearchQuery) *redis.StringSliceCmd {
	return r.readOnly.GeoSearch(ctx, key, q)
}

func (r *readonlyFuncWithoutTxDecorator) GeoDist(ctx context.Context, key string, member1, member2, unit string) *redis.FloatCmd {
	return r.readOnly.GeoDist(ctx, key, member1, member2, unit)
}

func (r *readonlyFuncWithoutTxDecorator) GeoHash(ctx context.Context, key string, members ...string) *redis.StringSliceCmd {
	return r.readOnly.GeoHash(ctx, key, members...)
}
