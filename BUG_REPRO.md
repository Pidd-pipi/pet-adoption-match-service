# Bug 复现说明

## Bug 是什么

宠物列表 context 未传播

## 如何触发

宠物列表查询忽略 ctx 取消与 deadline，客户端断开后仍继续查库。

```bash
cd backend
go test ./internal/service -run '^TestPetListContextCancellation$' -count=1
```
