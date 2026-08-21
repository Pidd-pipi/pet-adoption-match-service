# Bug 复现说明

## Bug 是什么

捐款使用明细错误链断裂

## 如何触发

不存在的捐款写使用明细时，repository 用 %v 断链、service 用错误 sentinel 判断，导致 404 被误判成 500。

```bash
cd backend
go test ./internal/service -run '^TestDonationErrorChainAndMapping$' -count=1
go test ./internal/handler -run '^TestDonationUsageHandlerStatus$' -count=1
```
