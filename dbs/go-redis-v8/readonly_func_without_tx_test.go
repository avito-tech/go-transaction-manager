//go:build go1.14
// +build go1.14

package go_redis_v8

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
	redismock "github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/require"
)

const (
	key1 = "key1"
	key2 = "key2"

	start = 1
	end   = 10
	stop  = 20

	offset = 1

	cursor = 1
	match  = "match"
	count  = 3

	field = "field"

	index = 11

	value = "value"

	member = "member"
	unit   = "unit"

	stream = "stream"

	rangeStart = "start"
	rangeStop  = "stop"

	group = "group"

	min = "min"
	max = "max"
)

var (
	keys  = []string{key1, key2}
	store = &redis.ZStore{}
	opt   = &redis.ZRangeBy{}
)

// nolint:ireturn
func newReadOnlyFuncWithoutTxDecorator() (Cmdable, redismock.ClientMock) {
	db, mock := redismock.NewClientMock()

	return ReadOnlyFuncWithoutTxDecorator(nil, db), mock
}

func check(t *testing.T, cmd redis.Cmder, mock redismock.ClientMock) {
	require.NoError(t, cmd.Err())
	require.NoError(t, mock.ExpectationsWereMet())
}

func Test_readonlyFuncWithoutTxDecorator_Dump(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectDump(key1).SetVal(OK)

	cmd := cmdable.Dump(context.Background(), key1)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_Exists(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectExists(key1).SetVal(0)

	cmd := cmdable.Exists(context.Background(), key1)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_Keys(t *testing.T) {
	t.Parallel()

	const pattern = "pattern"

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectKeys(pattern).SetVal(nil)

	cmd := cmdable.Keys(context.Background(), pattern)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_PTTL(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectPTTL(key1).SetVal(0)

	cmd := cmdable.PTTL(context.Background(), key1)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_RandomKey(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectRandomKey().SetVal(OK)

	cmd := cmdable.RandomKey(context.Background())

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_Touch(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectTouch(keys...).SetVal(0)

	cmd := cmdable.Touch(context.Background(), keys...)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_TTL(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectTTL(key1).SetVal(0)

	cmd := cmdable.TTL(context.Background(), key1)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_Type(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectType(key1).SetVal(OK)

	cmd := cmdable.Type(context.Background(), key1)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_Get(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectGet(key1).SetVal(OK)

	cmd := cmdable.Get(context.Background(), key1)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_GetRange(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectGetRange(key1, start, end).SetVal(OK)

	cmd := cmdable.GetRange(context.Background(), key1, start, end)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_MGet(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectMGet(keys...).SetVal(nil)

	cmd := cmdable.MGet(context.Background(), keys...)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_StrLen(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectStrLen(key1).SetVal(0)

	cmd := cmdable.StrLen(context.Background(), key1)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_GetBit(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectGetBit(key1, offset).SetVal(0)

	cmd := cmdable.GetBit(context.Background(), key1, offset)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_BitCount(t *testing.T) {
	t.Parallel()

	bitCount := &redis.BitCount{Start: start, End: end}

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectBitCount(key1, bitCount).SetVal(0)

	cmd := cmdable.BitCount(context.Background(), key1, bitCount)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_BitPos(t *testing.T) {
	t.Parallel()

	const (
		bit = 1
		pos = 2
	)

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectBitPos(key1, bit, pos).SetVal(0)

	cmd := cmdable.BitPos(context.Background(), key1, bit, pos)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_Scan(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectScan(cursor, match, count).SetVal(nil, cursor)

	cmd := cmdable.Scan(context.Background(), cursor, match, count)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_SScan(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectSScan(key1, cursor, match, count).SetVal(nil, cursor)

	cmd := cmdable.SScan(context.Background(), key1, cursor, match, count)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_HScan(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectHScan(key1, cursor, match, count).SetVal(nil, cursor)

	cmd := cmdable.HScan(context.Background(), key1, cursor, match, count)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_ZScan(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectZScan(key1, cursor, match, count).SetVal(nil, cursor)

	cmd := cmdable.ZScan(context.Background(), key1, cursor, match, count)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_HExists(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectHExists(key1, field).SetVal(false)

	cmd := cmdable.HExists(context.Background(), key1, field)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_HGet(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectHGet(key1, field).SetVal(OK)

	cmd := cmdable.HGet(context.Background(), key1, field)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_HGetAll(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectHGetAll(key1).SetVal(nil)

	cmd := cmdable.HGetAll(context.Background(), key1)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_HKeys(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectHKeys(key1).SetVal(nil)

	cmd := cmdable.HKeys(context.Background(), key1)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_HLen(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectHLen(key1).SetVal(0)

	cmd := cmdable.HLen(context.Background(), key1)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_HMGet(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectHMGet(key1, field).SetVal(nil)

	cmd := cmdable.HMGet(context.Background(), key1, field)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_HVals(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectHVals(key1).SetVal(nil)

	cmd := cmdable.HVals(context.Background(), key1)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_HRandField(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectHRandField(key1, count, false).SetVal(nil)

	cmd := cmdable.HRandField(context.Background(), key1, count, false)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_LIndex(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectLIndex(key1, index).SetVal(OK)

	cmd := cmdable.LIndex(context.Background(), key1, index)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_LLen(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectLLen(key1).SetVal(0)

	cmd := cmdable.LLen(context.Background(), key1)

	check(t, cmd, mock)
}

var lPostArgs = redis.LPosArgs{}

func Test_readonlyFuncWithoutTxDecorator_LPos(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectLPos(key1, value, lPostArgs).SetVal(0)

	cmd := cmdable.LPos(context.Background(), key1, value, lPostArgs)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_LRange(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectLRange(key1, start, stop).SetVal(nil)

	cmd := cmdable.LRange(context.Background(), key1, start, stop)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_SCard(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectSCard(key1).SetVal(0)

	cmd := cmdable.SCard(context.Background(), key1)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_SDiff(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectSDiff(keys...).SetVal(nil)

	cmd := cmdable.SDiff(context.Background(), keys...)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_SInter(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectSInter(keys...).SetVal(nil)

	cmd := cmdable.SInter(context.Background(), keys...)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_SIsMember(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectSIsMember(key1, member).SetVal(false)

	cmd := cmdable.SIsMember(context.Background(), key1, member)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_SMIsMember(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectSMIsMember(key1, member).SetVal(nil)

	cmd := cmdable.SMIsMember(context.Background(), key1, member)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_SMembers(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectSMembers(key1).SetVal(nil)

	cmd := cmdable.SMembers(context.Background(), key1)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_SRandMember(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectSRandMember(key1).SetVal(OK)

	cmd := cmdable.SRandMember(context.Background(), key1)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_SUnion(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectSUnion(keys...).SetVal(nil)

	cmd := cmdable.SUnion(context.Background(), keys...)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_XLen(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectXLen(stream).SetVal(0)

	cmd := cmdable.XLen(context.Background(), stream)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_XRange(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectXRange(stream, rangeStart, rangeStop).SetVal(nil)

	cmd := cmdable.XRange(context.Background(), stream, rangeStart, rangeStop)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_XRevRange(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectXRevRange(stream, rangeStart, rangeStop).SetVal(nil)

	cmd := cmdable.XRevRange(context.Background(), stream, rangeStart, rangeStop)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_XRead(t *testing.T) {
	t.Parallel()

	a := &redis.XReadArgs{}

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectXRead(a).SetVal(nil)

	cmd := cmdable.XRead(context.Background(), a)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_XPending(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectXPending(stream, group).SetVal(&redis.XPending{})

	cmd := cmdable.XPending(context.Background(), stream, group)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_ZCard(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectZCard(key1).SetVal(0)

	cmd := cmdable.ZCard(context.Background(), key1)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_ZCount(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectZCount(key1, min, max).SetVal(0)

	cmd := cmdable.ZCount(context.Background(), key1, min, max)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_ZLexCount(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectZLexCount(key1, min, max).SetVal(0)

	cmd := cmdable.ZLexCount(context.Background(), key1, min, max)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_ZInter(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectZInter(store).SetVal(nil)

	cmd := cmdable.ZInter(context.Background(), store)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_ZMScore(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectZMScore(key1, member).SetVal(nil)

	cmd := cmdable.ZMScore(context.Background(), key1, member)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_ZRange(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectZRange(key1, start, stop).SetVal(nil)

	cmd := cmdable.ZRange(context.Background(), key1, start, stop)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_ZRangeByScore(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectZRangeByScore(key1, opt).SetVal(nil)

	cmd := cmdable.ZRangeByScore(context.Background(), key1, opt)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_ZRangeByLex(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectZRangeByLex(key1, opt).SetVal(nil)

	cmd := cmdable.ZRangeByLex(context.Background(), key1, opt)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_ZRank(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectZRank(key1, member).SetVal(0)

	cmd := cmdable.ZRank(context.Background(), key1, member)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_ZRevRange(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectZRevRange(key1, start, stop).SetVal(nil)

	cmd := cmdable.ZRevRange(context.Background(), key1, start, stop)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_ZRevRangeByScore(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectZRevRangeByScore(key1, opt).SetVal(nil)

	cmd := cmdable.ZRevRangeByScore(context.Background(), key1, opt)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_ZRevRangeByLex(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectZRevRangeByLex(key1, opt).SetVal(nil)

	cmd := cmdable.ZRevRangeByLex(context.Background(), key1, opt)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_ZRevRank(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectZRevRank(key1, member).SetVal(0)

	cmd := cmdable.ZRevRank(context.Background(), key1, member)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_ZScore(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectZScore(key1, member).SetVal(0)

	cmd := cmdable.ZScore(context.Background(), key1, member)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_ZUnion(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectZUnion(*store).SetVal(nil)

	cmd := cmdable.ZUnion(context.Background(), *store)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_ZRandMember(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectZRandMember(key1, count, true).SetVal(nil)

	cmd := cmdable.ZRandMember(context.Background(), key1, count, true)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_ZDiff(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectZDiff(keys...).SetVal(nil)

	cmd := cmdable.ZDiff(context.Background(), keys...)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_PFCount(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectPFCount(keys...).SetVal(0)

	cmd := cmdable.PFCount(context.Background(), keys...)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_DBSize(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectDBSize().SetVal(0)

	cmd := cmdable.DBSize(context.Background())

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_GeoPos(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectGeoPos(key1, member).SetVal(nil)

	cmd := cmdable.GeoPos(context.Background(), key1, member)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_GeoSearch(t *testing.T) {
	t.Parallel()

	q := &redis.GeoSearchQuery{}

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectGeoSearch(key1, q).SetVal(nil)

	cmd := cmdable.GeoSearch(context.Background(), key1, q)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_GeoDist(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectGeoDist(key1, member, member, unit).SetVal(0)

	cmd := cmdable.GeoDist(context.Background(), key1, member, member, unit)

	check(t, cmd, mock)
}

func Test_readonlyFuncWithoutTxDecorator_GeoHash(t *testing.T) {
	t.Parallel()

	cmdable, mock := newReadOnlyFuncWithoutTxDecorator()
	mock.ExpectGeoHash(key1, member).SetVal(nil)

	cmd := cmdable.GeoHash(context.Background(), key1, member)

	check(t, cmd, mock)
}
