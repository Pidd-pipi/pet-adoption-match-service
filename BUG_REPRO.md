# Bug 复现说明

## Bug 是什么

收藏冲突错误判错

## 如何触发

重复收藏时 sentinel 未正确包装、service 用错误比较方式，冲突被误判成系统错误。

```bash
cd backend
go test ./internal/service -run '^TestFavoriteAddDuplicateConflict$' -count=1
go test ./internal/service -run '^TestFavoriteRemoveMissing$' -count=1
```
