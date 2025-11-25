# Flux Protocol Coding Guidelines & Best Practices

**Status:** Active
**Last Updated:** November 20, 2025

This document defines the technical standards and architectural best practices for the Flux Protocol project. All contributors must adhere to these guidelines to ensure performance, security, and maintainability.

---

## 1. Rust Core (Agent & CLI)

### 1.1 Async Runtime (Tokio)
*   **Blocking Operations**: **NEVER** block the async thread.
    *   **Bad**: `std::thread::sleep`, `std::fs::read`, `rusqlite::Connection::open` inside an `async fn`.
    *   **Good**: Use `tokio::time::sleep`, `tokio::fs`, or wrap synchronous CPU/IO intensive work in `tokio::task::spawn_blocking`.
*   **Error Handling**: Use `thiserror` for library/module errors and `anyhow` for the top-level application/CLI.
    *   Implement `IntoResponse` for custom errors in Axum to return semantic HTTP status codes.
*   **Shutdown**: Implement graceful shutdown. Listen for `tokio::signal::ctrl_c` and broadcast a shutdown signal (via `tokio::sync::broadcast`) to all worker tasks (File Watcher, Sync Engine, etc.) to clean up resources.

### 1.2 Web Framework (Axum)
*   **State Management**: Use `Arc<AppState>` with `axum::extract::State`.
    *   Avoid global mutable state (`static mut`).
*   **Middleware**: Use `tower-http` for common middleware:
    *   `TraceLayer` for structured logging.
    *   `CompressionLayer` for Gzip/Brotli (essential for large JSON payloads).
*   **Validation**: Use the `validator` crate on deserialized structs to enforce data integrity at the API boundary.

### 1.3 CLI Design (Clap v4)
*   **Derive API**: Use the `derive` pattern (`#[derive(Parser)]`) for type-safe argument parsing.
*   **UX Patterns**:
    *   **Output**: Use `tabled` for list commands. Always provide a `--json` flag for scripting compatibility.
    *   **Feedback**: Use `indicatif` spinners for any operation taking > 200ms.
    *   **Colors**: Respect `NO_COLOR` env var. Use `colored` crate for semantic coloring (Green=Success, Red=Error, Yellow=Warning).
*   **Structure**: Split large subcommands into separate modules (`src/commands/task.rs`, `src/commands/peer.rs`).

### 1.4 SQLite & FTS5
*   **Connection Pooling**: SQLite is synchronous. Use `deadpool-sqlite` or manage a `r2d2` pool within `spawn_blocking`.
*   **FTS5 Optimization**:
    *   Use `VACUUM` periodically to reclaim space and optimize B-Trees.
    *   For search queries, utilize `prefix` queries (`*`) carefully as they can be expensive.
*   **Migrations**: Use `rusqlite_migration` to manage schema evolution embedded in the binary.

---

## 2. Distributed Systems & CRDTs

### 2.1 Vector Clocks
*   **Structure**: Use `HashMap<PeerId, u64>`.
*   **Causality**: Always check causality before applying an operation.
    *   If `op.clock <= local.clock`, it's a duplicate.
    *   If `op.clock > local.clock` (with gaps), buffer it or request missing ops (if implementing strict causal delivery). For Flux MVP, we accept "out of order" but ensure strict LWW application.
*   **Transmission**: Compress vector clocks during sync if the peer count grows large (e.g., send only the delta/diff of the clock).

### 2.2 WebRTC (webrtc-rs)
*   **Data Channels**: Use `ordered: true` for the sync stream.
*   **Keep-Alives**: Implement application-level Ping/Pong every 30s if the transport doesn't provide it reliably.
*   **Buffer Management**: Monitor `buffered_amount()` on the DataChannel. If the buffer fills (backpressure), pause sending to avoid OOM crashes.

---

## 3. Go Registry Service

### 3.1 Concurrency
*   **Goroutine-per-Client**: Spawn a read pump and a write pump goroutine for each WebSocket connection.
*   **Channel Safety**: Never write to a closed channel. Use a `quit` channel pattern to signal shutdown to writer goroutines.
*   **Resource Limits**: Set `MaxMessageSize` to prevent memory DoS attacks.

### 3.2 Security
*   **Origin Check**: **ALWAYS** validate the `Origin` header in `upgrader.CheckOrigin`. Do not return `true` unconditionally in production.
*   **Timeouts**: Set `ReadDeadline` and `WriteDeadline` to detect dead clients and prevent file descriptor leaks.

---

## 4. Frontend (Tauri & React)

### 4.1 Tauri Architecture
*   **Security Boundary**: Treat the Frontend as untrusted. Validate all inputs in the Rust backend commands.
*   **IPC**:
    *   Use `tauri::command` for Request/Response.
    *   Use `window.emit()` for server-push events (Sync updates).
*   **State Sync**: Do not duplicate complex state logic in Redux/Zustand. The Rust backend is the source of truth. The Frontend is a **projection**.

### 4.2 React Best Practices
*   **Performance**: Use `React.memo` for list items (Tasks) to prevent re-rendering the entire board on a single task update.
*   **Virtualization**: Use `tanstack/react-virtual` or similar if rendering > 100 tasks in a column.
*   **Modularity**: Isolate "Smart Components" (data fetching) from "Dumb Components" (rendering).

---

## 5. Mobile (React Native)

### 5.1 WebRTC Mobile
*   **Backgrounding**:
    *   iOS: Connections **will** die in background. Handle `RTCPeerConnection.onConnectionStateChange` -> `failed` gracefully.
    *   Android: Use a Foreground Service if persistent sync is required (battery drain warning).
*   **Codecs**: Prefer **H.264** on iOS (hardware accelerated). VP8/VP9 may use software encoding and heat up the device.

### 5.2 UX & Gestures
*   **Native Feel**: Use `react-native-gesture-handler` and `reanimated` for 60fps interactions. Avoid JS-driven animations.
*   **Lists**: Use `FlashList` (from Shopify) instead of `FlatList` for complex task lists to minimize blank areas during fast scrolling.

---

## 6. Git Workflow & Version Control

*   **Commits**: Use Conventional Commits (`feat:`, `fix:`, `chore:`).
*   **Branches**:
    *   `main`: Stable, deployable.
    *   `feature/xyz`: Short-lived feature branches.
*   **PRs**: Must pass CI (clippy, tests, fmt) before merge.
