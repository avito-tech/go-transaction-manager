# Result


Output of the command `go build -v -gcflags="-m -l" ./...`.

```
./key_in_ctx.go:26:44: leaking param: wg
./key_in_ctx.go:17:2: moved to heap: wg
./key_in_ctx.go:15:25: key escapes to heap
./key_in_ctx.go:15:25: idKey escapes to heap
./tx_in_ctx.go:29:43: leaking param: wg
./tx_in_ctx.go:20:2: moved to heap: wg
./tx_in_ctx.go:15:8: &sqlx.Tx{} escapes to heap
./tx_in_ctx.go:18:25: key escapes to heap
```

## Conclusion

context.Context is placed in stack because ctx is not moved in heap.