# Bug 复现说明

## Bug 是什么

机构查询 nil 解引用 panic

## 如何触发

机构未命中时返回 nil,nil 丢失 sentinel，service 判空失效后 nil 解引用 panic。

```bash
cd backend
go test ./internal/service -run '^TestOrgGetMissingNoPanic$' -count=1
go test ./internal/service -run '^TestOrgReviewMissingNoPanic$' -count=1
```

```
panic: runtime error: invalid memory address or nil pointer dereference
```
