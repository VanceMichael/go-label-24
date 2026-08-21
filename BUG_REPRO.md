# Bug Reproduction

## 包的性质

当前 test_model_fix 保存的是被测模型修复后的结果源码，不是初始含 Bug 源码。要复现原始缺陷，必须检出下面固定的 parent SHA；不要在当前修复结果源码上期待重新出现修复前失败。生成系统使用的可信验证补丁和完整验证日志仅在本地留存，不提交到结果分支。

## 问题现象

维修主管一次提交两台挤奶设备的完工报告，第二份因为处置说明为空被拒绝。接口虽然返回失败，第一张工单却已经变成 completed，对应设备也提前恢复 operational 并增加了版本；再次补齐报告提交时随即遇到版本冲突。请修复批量完工流程，让整批先完成校验和状态推演，任何一项失败都保持所有工单与设备原样，全部有效时才一次性发布本批结果。

## 含 Bug 版本

- 仓库：VanceMichael/go-label-24
- 仓库地址：https://github.com/VanceMichael/go-label-24.git
- parent SHA：69d81ae54b9354805a168419bc2bfe7771b76412

## 复现步骤

```bash
git clone -- https://github.com/VanceMichael/go-label-24.git bug-repro
cd bug-repro
git checkout --detach 69d81ae54b9354805a168419bc2bfe7771b76412
go test ./internal/maintenance -run ^TestCompleteBatchRejectsPartialMaintenanceState$ -count=1
```

## 双架构完整错误信息

### linux/amd64

- 容器内复现预期退出码：1
- 容器内复现实际退出码：1

stdout：

```text
$ go test ./internal/maintenance -run ^TestCompleteBatchRejectsPartialMaintenanceState$ -count=1
--- FAIL: TestCompleteBatchRejectsPartialMaintenanceState (0.00s)
    batch_completion_test.go:35: failed batch exposed partial state: orders=[{ID:order-first TenantID:farm AssetID:asset-first RequestedBy: AssignedTo:tech Priority:0 Summary: Status:completed MeterAtOpen:100 ScheduledWindow:{StartsAt:0001-01-01 00:00:00 +0000 UTC EndsAt:0001-01-01 00:00:00 +0000 UTC} OpenedAt:0001-01-01 00:00:00 +0000 UTC StartedAt:2026-08-21 14:41:21.725557791 +0000 UTC CompletedAt:2026-08-21 15:41:21.725557791 +0000 UTC Version:5} {ID:order-second TenantID:farm AssetID:asset-second RequestedBy: AssignedTo:tech Priority:0 Summary: Status:in_progress MeterAtOpen:100 ScheduledWindow:{StartsAt:0001-01-01 00:00:00 +0000 UTC EndsAt:0001-01-01 00:00:00 +0000 UTC} OpenedAt:0001-01-01 00:00:00 +0000 UTC StartedAt:2026-08-21 14:41:21.725557791 +0000 UTC CompletedAt:<nil> Version:4}] assets=[{ID:asset-first TenantID:farm BarnID:barn Name:first Category:milking Status:operational MeterHours:101 ServiceEvery:200 LastServiceHour:101 Version:8} {ID:asset-second TenantID:farm BarnID:barn Name:second Category:milking Status:servicing MeterHours:100 ServiceEvery:200 LastServiceHour:80 Version:7}]
FAIL
FAIL	go-base/internal/maintenance	0.031s
FAIL

```

stderr：

```text
(empty)
```

### linux/arm64

- 容器内复现预期退出码：1
- 容器内复现实际退出码：1

stdout：

```text
$ go test ./internal/maintenance -run ^TestCompleteBatchRejectsPartialMaintenanceState$ -count=1
--- FAIL: TestCompleteBatchRejectsPartialMaintenanceState (0.00s)
    batch_completion_test.go:35: failed batch exposed partial state: orders=[{ID:order-first TenantID:farm AssetID:asset-first RequestedBy: AssignedTo:tech Priority:0 Summary: Status:completed MeterAtOpen:100 ScheduledWindow:{StartsAt:0001-01-01 00:00:00 +0000 UTC EndsAt:0001-01-01 00:00:00 +0000 UTC} OpenedAt:0001-01-01 00:00:00 +0000 UTC StartedAt:2026-08-21 14:41:54.171864751 +0000 UTC CompletedAt:2026-08-21 15:41:54.171864751 +0000 UTC Version:5} {ID:order-second TenantID:farm AssetID:asset-second RequestedBy: AssignedTo:tech Priority:0 Summary: Status:in_progress MeterAtOpen:100 ScheduledWindow:{StartsAt:0001-01-01 00:00:00 +0000 UTC EndsAt:0001-01-01 00:00:00 +0000 UTC} OpenedAt:0001-01-01 00:00:00 +0000 UTC StartedAt:2026-08-21 14:41:54.171864751 +0000 UTC CompletedAt:<nil> Version:4}] assets=[{ID:asset-first TenantID:farm BarnID:barn Name:first Category:milking Status:operational MeterHours:101 ServiceEvery:200 LastServiceHour:101 Version:8} {ID:asset-second TenantID:farm BarnID:barn Name:second Category:milking Status:servicing MeterHours:100 ServiceEvery:200 LastServiceHour:80 Version:7}]
FAIL
FAIL	go-base/internal/maintenance	0.001s
FAIL

```

stderr：

```text
(empty)
```

## 通过条件

当第一份维修报告有效而第二份缺少处置说明时，CompleteBatch 必须返回第二份报告的校验错误，且两张工单的状态、完成时间和版本以及两台设备的状态、计量值和版本均与调用前完全一致；两份报告都有效时才允许整批完成。TestCompleteBatchRejectsPartialMaintenanceState 要由红转绿，maintenance 既有用例、其余包回归和 go build ./... 继续通过，不得通过跳过后续报告、改写输入测试或放宽深度状态比较来规避失败。
