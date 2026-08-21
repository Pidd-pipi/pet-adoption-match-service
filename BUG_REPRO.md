# Bug 复现说明

## Bug 是什么

限流器（backend/internal/middleware/rate_limiter.go）和请求日志计数（backend/internal/middleware/logger.go）对共享的 per-IP 计数做并发读写时未加锁，存在 data race：并发请求下计数互相覆盖丢失，本应被限流的请求会被漏过。

## 如何触发

```bash
cd backend
go test -race ./internal/middleware -run '^TestConcurrentRateLimitAndRequestCount$' -count=1
```

## 真实错误信息

```
WARNING: DATA RACE
Write at 0x00c00009e000 by goroutine 9:
	internal/middleware/rate_limiter.go:38
Previous read at 0x00c00009e000 by goroutine 12:
	internal/middleware/rate_limiter.go:35
```
