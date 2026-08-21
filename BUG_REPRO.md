# Bug 复现说明

## Bug 是什么

帖子过滤切片别名污染

## 如何触发

帖子列表用 s[:0] 原地压缩共享底层数组，二次过滤污染原列表。

```bash
cd backend
go test ./internal/service -run '^TestPostFiltersNoAlias$' -count=1
go test ./internal/handler -run '^TestPostHelpersNoAlias$' -count=1
```
