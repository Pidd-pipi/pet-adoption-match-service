# Bug 复现说明

## Bug 是什么

回访逾期状态机错位

## 如何触发

逾期回访任务提交状态转换缺口、MarkOverdue 漏 pending 过滤、旧状态回写。

```bash
cd backend
go test ./internal/repository -run '^TestMarkOverdueOnlyPending$' -count=1
go test ./internal/service -run '^TestReviewOverdueSubmit$' -count=1
```
