# Streamforge CDC

Streamforge CDC是纯Go数据库变更捕获与可靠复制服务，将PostgreSQL逻辑复制和MySQL binlog抽象为统一事务事件。

## 当前能力

- 数据源注册、校验和状态管理，支持PostgreSQL/MySQL模型。
- PostgreSQL复制消息和MySQL binlog事件的边界安全解码器。
- 统一Event/Transaction模型，稳定事件ID支持幂等重放。
- Schema版本比较、strict/compatible/raw策略基础实现。
- 内存队列、文件日志、配额和三阶段检查点模型。
- HTTP Webhook、Kafka兼容接口和内存接收器。
- REST控制面：/healthz、/readyz、/metrics、数据源、管道、事件游标和模拟事务接口。
- `web/` 提供一个可直接部署的静态控制台页面，用于查看服务健康状态并跳转到运维端点。

## 运行

```sh
go run ./cmd/streamforge-cdc
```

默认监听`:8092`，管理API需要`Authorization: Bearer dev-key`。外部数据库不可用时，使用`POST /api/v1/simulate/transaction`验证流程。执行`sh scripts/smoke.sh`可完成冒烟。静态页面可由任意 Web 服务器托管；如果与 Go 服务同源部署，页面会直接请求`/healthz`。

## 可靠性语义

来源读取、持久化和下游确认检查点分开记录。只有下游确认成功才推进confirmed位置。崩溃恢复允许重复事件，不允许丢失已确认变更；普通HTTP目标不宣称精确一次，消费者应按事件ID幂等。

## 验证

```sh
gofmt -w .
go vet ./...
go test -race ./...
go build ./...
```
