# Bug 复现说明

## Bug 是什么

首页聚合 nil map 写入 panic

## 如何触发

首页聚合计数（internal/service/stats_service.go 与 internal/handler/home_handler.go）写路径未初始化内层 map，首次写入触发 panic。

```bash
cd backend
go test ./internal/service -run '^TestStatsNilMapSafety$' -count=1
go test ./internal/handler -run '^TestHomeHelpersNilMapSafety$' -count=1
```

```
panic: assignment to entry in nil map
```
