# TUI Styles Claude Configuration

This directory contains specialized AI agents for the TUI Styles project. These agents follow the patterns established in the protocol specifications and implementation plans.

## Agent Directory

All agents are prefixed with `{{DOMAIN}}-` (WorkStream Protocol, the original codename for Flux) to distinguish them from general-purpose agents.

### Core Protocol & Architecture Agents

#### {{DOMAIN}}-protocol-architect (Martin Kleppmann Persona)
**Use When**: Designing CRDT operations, implementing vector clocks, creating state synchronization algorithms, or resolving distributed conflicts.

**Expertise**: Distributed systems theory, conflict-free replicated data types, causality tracking, operation-based CRDTs, eventual consistency.

**Examples**:
- Implementing task synchronization between peers with concurrent updates
- Designing vector clock comparison and operation delta algorithms
- Handling network partitions and state merge logic
- Resolving conflicts in LWW (Last-Write-Wins) registers

#### {{DOMAIN}}-security-architect
**Use When**: Implementing capability-based security, designing JWT token systems, implementing HKDF secret derivation, or creating cryptographic access control.

**Expertise**: Applied cryptography, capability systems, HKDF key derivation, JWT verification, end-to-end encryption.

**Examples**:
- Deriving capability-specific secrets from root secret using HKDF
- Verifying JWT tokens and checking permissions for operations
- Implementing end-to-end encryption for relay messages
- Designing access control rules in `access.yaml`

#### {{DOMAIN}}-filesystem-architect (Chris Palmer Persona)
**Use When**: Designing workspace directory structure, creating file formats with YAML frontmatter, ensuring Git compatibility, or implementing atomic file operations.

**Expertise**: Filesystem-based data persistence, human-readable formats, atomic operations, YAML/Markdown parsing, Git integration.

**Examples**:
- Structuring task files with YAML frontmatter and Markdown content
- Implementing atomic write-rename pattern for file safety
- Making workspace Git-friendly with proper .gitignore patterns
- Handling editor temp files (e.g., vim swap files)

### Backend Implementation Agents

#### {{DOMAIN}}-rust-systems-engineer (Jon Gjengset Persona)
**Use When**: Building the tui-styles-agent daemon, implementing async runtime architecture, designing concurrent state management, or optimizing Rust performance.

**Expertise**: Systems programming with Tokio, async Rust, concurrent task spawning, file watching with `notify`, SQLite integration, Axum HTTP API.

**Examples**:
- Structuring agent daemon with file watching, WebRTC, and HTTP API concurrently
- Implementing notify-based file watching with debouncing
- Managing shared state (operation log, vector clock) with Arc and RwLock
- Optimizing async runtime performance

#### {{DOMAIN}}-go-specialist (Mat Ryer Persona)
**Use When**: Implementing the registry server in Go, designing concurrent goroutine architecture, implementing WebSocket handlers, or optimizing Go performance.

**Expertise**: Goroutines, WebSocket servers with Gorilla, Redis integration, concurrent message routing, channel patterns.

**Examples**:
- Designing goroutine-based architecture for thousands of concurrent WebSocket connections
- Tracking online peers per project in Redis with TTLs and pub/sub
- Optimizing registry performance for 10k concurrent connections
- Implementing connection pooling and message batching

#### {{DOMAIN}}-webrtc-specialist
**Use When**: Implementing WebRTC peer-to-peer connections, configuring DataChannels, handling ICE negotiation, or designing signaling protocols.

**Expertise**: Real-time communication, WebRTC DataChannels, STUN/TURN servers, NAT traversal, SDP exchange.

**Examples**:
- Setting up WebRTC DataChannels for agent-to-agent communication
- Implementing ICE candidate gathering and TURN relay fallback
- Designing WebSocket signaling protocol for offer/answer exchange
- Handling connection state changes and reconnection logic

### Frontend Implementation Agents

#### {{DOMAIN}}-tauri-desktop-architect (Daniel Thompson-Yvetot Persona)
**Use When**: Building tui-styles-desktop with Tauri, implementing WebRTC client in Rust backend, designing React frontend components, or integrating native OS features.

**Expertise**: Tauri IPC, Rust-to-React communication, WebRTC in desktop context, offline queue, event emission from Rust to React.

**Examples**:
- Organizing Tauri backend to handle WebRTC and communicate with React frontend
- Implementing event emission from Rust to React for real-time updates
- Designing offline queue with SQLite and automatic sync
- Integrating native OS features (file pickers, notifications, system tray)

#### {{DOMAIN}}-react-native-architect (Evan Bacon Persona)
**Use When**: Building tui-styles-mobile with React Native, implementing WebRTC on iOS/Android, handling background connections, or designing mobile-optimized UI.

**Expertise**: Cross-platform mobile development, react-native-webrtc, background tasks, push notifications, AsyncStorage.

**Examples**:
- Integrating WebRTC in React Native for iOS and Android
- Keeping WebRTC connection alive when app goes to background on iOS
- Designing swipeable column layout with mobile gestures
- Implementing offline queue with AsyncStorage

### UI & UX Agents

#### {{DOMAIN}}-ui-schema-designer
**Use When**: Designing the agent-driven UI schema protocol, creating versioned schema formats, implementing dynamic UI rendering, or building extensibility systems.

**Expertise**: Schema-driven architecture, UI abstraction, version management, feature detection, plugin systems.

**Examples**:
- Designing UISchema message format for agent-to-client communication
- Supporting custom workflows per project (columns, fields, validation)
- Implementing schema versioning with backward compatibility
- Creating extension points for UI plugins

#### {{DOMAIN}}-cli-ux-designer (Simon Willison Persona)
**Use When**: Designing CLI command structure, creating user-friendly output formatting, implementing interactive prompts, or writing help documentation.

**Expertise**: Command-line interface usability, table formatting with `tabled`, colored output, progressive disclosure, Clap CLI patterns.

**Examples**:
- Designing intuitive command structure with consistent flag patterns
- Formatting output of `tui-styles task list` for readability
- Designing interactive wizards with sensible defaults
- Creating comprehensive help text and error messages

#### {{DOMAIN}}-realtime-collaboration-specialist
**Use When**: Implementing real-time presence, designing collaborative cursors, handling concurrent edits, or building awareness systems.

**Expertise**: Multiplayer collaboration features, presence tracking, cursor positions, ephemeral state broadcasting.

**Examples**:
- Showing which users are viewing each task in real-time
- Displaying when someone else is typing a comment or editing a task
- Implementing cursor position broadcasting with interpolation
- Handling transient state (typing indicators, viewport state)

### Quality & Testing Agents

#### {{DOMAIN}}-testing-strategist
**Use When**: Designing test strategies for distributed systems, implementing property-based tests for CRDTs, creating integration tests for agent sync, or building performance benchmarks.

**Expertise**: Property-based testing with `proptest`, integration testing, performance benchmarking, CRDT convergence testing.

**Examples**:
- Testing that concurrent operations converge to the same state
- Designing integration tests with in-memory agent instances
- Creating performance benchmarks with realistic workload simulation
- Testing WebRTC connection handling without running multiple processes

#### {{DOMAIN}}-search-indexing-specialist
**Use When**: Implementing full-text search with SQLite FTS5, designing code-aware indexing, creating incremental index updates, or optimizing search query performance.

**Expertise**: SQLite FTS5, tokenization, code parsing with tree-sitter, symbol extraction, performance optimization.

**Examples**:
- Setting up FTS5 to index both task content and code files
- Indexing function definitions and class names separately from regular text
- Optimizing FTS5 query performance and tokenization
- Implementing incremental indexing on file changes

### Planning & Coordination Agents

#### {{DOMAIN}}-planner
**Use When**: Breaking down TUI Styles implementation into phases, identifying dependencies between components, creating development roadmaps, or coordinating multi-agent work.

**Expertise**: Strategic planning, project decomposition, dependency analysis, MVP scoping, parallel development coordination.

**Examples**:
- Creating phased implementation plan with dependencies identified
- Identifying component boundaries for parallel work streams
- Defining MVP features and minimal viable implementation
- Coordinating agent, CLI, registry, desktop, and mobile development

#### {{DOMAIN}}-api-designer
**Use When**: Designing REST API endpoints, creating request/response schemas, implementing error handling, or defining API authentication.

**Expertise**: REST API design, HTTP semantics, error response formats, API versioning, OpenAPI specifications.

**Examples**:
- Designing RESTful API for local agent HTTP endpoints
- Creating standardized error response format with problem details
- Recommending API versioning strategy based on evolution needs
- Defining resource modeling and status code selection

## Agent Usage Patterns

### Protocol Implementation
1. **Design Phase**: Start with `{{DOMAIN}}-protocol-architect` for CRDT operations
2. **Security**: Use `{{DOMAIN}}-security-architect` for capability systems
3. **Testing**: Engage `{{DOMAIN}}-testing-strategist` for convergence tests

### Backend Development
1. **Agent Daemon**: `{{DOMAIN}}-rust-systems-engineer` for concurrent architecture
2. **Registry Server**: `{{DOMAIN}}-go-specialist` for WebSocket handling
3. **WebRTC**: `{{DOMAIN}}-webrtc-specialist` for P2P connections

### Frontend Development
1. **Desktop**: `{{DOMAIN}}-tauri-desktop-architect` for Tauri + React
2. **Mobile**: `{{DOMAIN}}-react-native-architect` for React Native
3. **UI Schema**: `{{DOMAIN}}-ui-schema-designer` for dynamic rendering

### Planning & Coordination
1. **Roadmap**: `{{DOMAIN}}-planner` for implementation strategy
2. **CLI Design**: `{{DOMAIN}}-cli-ux-designer` for command structure
3. **Testing Strategy**: `{{DOMAIN}}-testing-strategist` for test plans

## Configuration Files

### settings.local.json
Defines permissions for Claude Code:
- `allow`: Commands that can run without asking
- `deny`: Commands that are blocked
- `ask`: Commands that require user confirmation

## Workflow Integration

These agents work seamlessly with the Flux CLI tool (`tools/tui-styles-cli`) and follow the patterns defined in:
- `CLAUDE.md` - Project overview and development workflow
- `CODING_GUIDELINES.md` - Rust/Go/TypeScript standards
- `specs/wsp.md` - Complete protocol specification
- `plan/*.md` - Detailed implementation plans per component

## Best Practices

1. **Read Guidelines First**: All agents automatically read `CODING_GUIDELINES.md` and relevant specs
2. **Follow Patterns**: Agents cite established patterns from plan/ documents
3. **Use CLI Tool**: All roadmap operations go through `tui-styles-cli` for consistency
4. **Parallel Work**: Agents can work in parallel on independent components
5. **Dependencies**: Agents respect component boundaries and interfaces defined in plans

## Adding New Agents

When adding new specialized agents:
1. Create agent file in `.claude/agents/` with `{{DOMAIN}}-{name}.md` naming
2. Define clear "Use When" criteria and expertise areas
3. Provide concrete examples of tasks the agent handles
4. Document in this README
5. Ensure agent reads relevant guidelines and specs

## Support

For questions about agent usage or configuration:
- Review `CLAUDE.md` for overall project structure
- Check `specs/wsp.md` for protocol details
- Consult `plan/*.md` for component-specific guidance
