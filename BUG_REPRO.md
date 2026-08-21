# Bug 复现说明

## Bug 是什么

领养申请状态机错位

## 如何触发

领养申请状态机转换表漏边、非法迁移被静默放行、宠物状态回写旧常量。

```bash
cd backend
go test ./internal/constants -run '^TestConfirmedRejectedTransition$' -count=1
go test ./internal/service -run '^TestApplicationStatusMachine$' -count=1
```
