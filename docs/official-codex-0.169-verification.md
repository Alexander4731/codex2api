# Codex 0.169.0 应用层身份核对

核对标签 `openai/codex rust-v0.169.0`，commit `d32063f96c2dff1dc96af839b3952abb6acfeda4`。

- `codex-rs/core/src/responses_metadata.rs:337` 的 client_metadata 包含安装 ID、session_id、thread_id、窗口 ID；`compatibility_headers:390` 不生成独立安装 ID 头，并从兼容 turn metadata 头移除工具清单。现有网关的这些字段处理保持一致。
- `codex-rs/runtime-agent-instance/src/installation_id.rs:19` 在 Codex 目录内持久化并复用 UUIDv4；本服务按登录账号固定安装身份，是部署策略，非官方自动按账号隔离的机制。
- `codex-rs/model-provider-info/src/lib.rs:585` 使用构建版本生成 Version 头。网关通过账号版本模板跟随同步值，避免 User-Agent/Version 冲突。
- HTTP/compact 实际出站捕获测试覆盖记录 1/2/3 与不同下游设备输入；WebSocket 实际握手/帧捕获测试覆盖两个登录、三个工作区及不同下游调用方，并验证会话字段一致。
- 不声称 Rust 与 Go 的 TLS/HTTP2 指纹、所有可选功能和整个请求逐字节相同。生产推理不截取 TLS 密文后的完整业务请求；安装身份完整字段由本地 mock 上游的实际出站请求验证，服务器配置和真实模型请求另行核对。

本次巡检暂停时段改动不改变上述字段协议。
