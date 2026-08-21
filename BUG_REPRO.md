# Bug 复现说明

## Bug 是什么

评论事务 defer 吞错

## 如何触发

评论事务 defer 提交覆盖主错误、错误分支未回滚，事务错误被吞。

```bash
cd backend
go test ./internal/repository -run '^TestPostCommentRepoPreservesError$' -count=1
go test ./internal/service -run '^TestCommentCreatePreservesError$' -count=1
go test ./internal/service -run '^TestCommentDeletePreservesError$' -count=1
```
