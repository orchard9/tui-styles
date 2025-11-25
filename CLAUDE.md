# TUI Styles: Filesystem-First P2P Collaborative Project Management

We are building TUI Styles, a peer-to-peer collaborative project management platform that treats the filesystem as the source of truth. Users create local `.tui-styles/` workspaces with tasks as Markdown files, run a local agent daemon, and sync in real-time via WebRTC P2P connections using CRDTs for conflict-free convergence.

**Tech Stack**:

-   **Agent/CLI**: Rust (Tokio async runtime, `webrtc-rs`, `notify`, `rusqlite` FTS5, `axum` API, `clap` CLI)
-   **Registry**: Go (WebSocket signaling server, Gorilla WebSocket, Redis ephemeral state)
-   **Desktop**: Tauri 2.0 (Rust backend + React frontend with TypeScript, Radix UI, Tailwind CSS)
-   **Mobile**: React Native (`react-native-webrtc`, AsyncStorage, React Navigation)
-   **Protocol**: MessagePack (binary serialization), Vector Clocks (causality), Operation-based CRDTs (LWW-Register)
-   **Development**: Rust workspace (tui-styles-agent, tui-styles-cli, tui-styles-lib), Go module (tui-styles-registry), Node.js monorepo (tui-styles-desktop, tui-styles-mobile)

**Development**: Multi-component development with separate workflows:

-   Agent/CLI: `cargo build` (shared workspace), `cargo run --bin tui-styles-agent`, `cargo test`
-   Registry: `go run main.go` (WebSocket signaling server on port 8080)
-   Desktop: `npm run tauri dev` (Tauri + React hot reload)
-   Mobile: React Native CLI with Xcode/Android Studio live reload
-   Integration tests: Docker Compose for multi-agent sync testing

**Core User Journey**:

**Solo Developer**: `tui-styles init` → Create tasks via CLI/filesystem → `tui-styles start` (agent watches files) → Work offline → Files sync via Git

**Team Collaboration**: Share secret → `tui-styles connect <secret>` → Agent establishes WebRTC P2P → Real-time CRDT sync → Desktop/mobile UI renders agent-driven schema → Collaborate without server

**Architecture Overview**

**Core Principle**: Filesystem as database - `.tui-styles/` directory contains YAML frontmatter Markdown files for tasks, SQLite for search index/operation log (gitignored), WebRTC P2P mesh for real-time sync, CRDTs for conflict resolution, thin UI clients that render schemas from agents.

**System Components**:

1. **tui-styles-agent (Rust daemon)**: File watching, CRDT sync engine, WebRTC peer connections, SQLite FTS5 search, local HTTP API
2. **tui-styles-cli (Rust)**: Command-line task management, agent control, search, sharing
3. **tui-styles-desktop (Tauri)**: Thin UI client connecting to agents via WebRTC, renders UISchema
4. **tui-styles-mobile (React Native)**: Mobile UI with offline queue, background sync, push notifications
5. **tui-styles-registry (Go)**: Optional WebRTC signaling server, peer discovery, relay fallback

**Backend Architecture**: Local-first agent daemon + optional signaling registry

-   **Agent Service**: Watches `.tui-styles/`, generates CRDT operations, maintains vector clock, indexes with FTS5
-   **Sync Engine**: Processes operation log, resolves conflicts via LWW, broadcasts to peers
-   **WebRTC Manager**: Establishes P2P DataChannels, handles ICE negotiation, maintains connections
-   **Local API**: Axum HTTP server for CLI/desktop communication
-   **Registry**: Go WebSocket server for peer discovery (knows nothing about task data)

**Frontend Architecture**:

-   **Desktop**: Tauri Rust backend (WebRTC client, local storage) + React UI (renders UISchema from agent)
-   **Mobile**: React Native with on-device WebRTC, AsyncStorage offline queue, platform-specific background sync
-   **Agent-Driven UI**: Agents send UISchema (columns, fields, capabilities) → Clients render dynamically

**Data Flow**:

-   **Local Write**: Edit `task.md` → File watcher detects → Generate `OpKind::TaskUpdate` → Apply to DB → Update vector clock → Broadcast to peers via WebRTC
-   **Remote Write**: Receive `OpKind::TaskUpdate` → Validate causality → Apply to DB → Materialize file → Atomic write to `.tui-styles/`
-   **Sync**: Connect via WebRTC → Exchange `Hello` + vector clocks → Calculate delta → Send missing operations → Apply in causal order → Converge

## SDLC Workflow Commands

TUI Styles provides slash commands for AI SDLC workflows (adapted from Masquerade project-agnostic patterns):

-   **`/project-status`** - Analyze roadmap progress and recommend next actions
    -   Counts tasks by status, identifies blockers, recommends next 2-3 tasks
    -   Calculates velocity and estimates completion timeline
    -   Use daily for standup or when deciding what to work on

-   **`/plan-milestone`** - Create new milestone with phases and tasks
    -   Evolutionary design: quick-wins → core-features → polish
    -   Generates structured task files with dependencies
    -   Use when adding new features or major work packages

-   **`/preplan`** - Pre-plan pending tasks to achieve 70%+ confidence
    -   Researches technical approaches, validates task size, assigns confidence scores
    -   Moves tasks to ready (≥70%) or blocked (<70% with flagged comments)
    -   Use before starting implementation on any pending tasks

-   **`/verify`** - Verify completed work meets acceptance criteria
    -   Generates verification checklist and test evidence capture
    -   Creates temporary `.verify/` directory for manual review
    -   Use after implementing features, before marking complete

-   **`/align-milestone`** - Validate alignment and ENFORCE corrections
    -   Analyzes coverage + calculates alignment score (0-100%, threshold: 80%)
    -   **Takes action**: blocks ambiguous tasks, updates index.md, creates missing tasks
    -   **Saves report**: Creates `roadmap/milestone-{N}/alignment.md` for reference
    -   Creates ALIGNMENT_DECISIONS_NEEDED.md if human choices required
    -   Use after planning (validate + fix plan) or after delivery (confirm goals shipped)

-   **`/ship-milestone`** - Autonomously execute entire milestone using coordinated specialist agents
    -   Follows 11-state FSM workflow: assess → align → preplan → implement → verify → goal-delivery → iterate
    -   Alignment gates: Early (validate plan) + Late (confirm delivery)
    -   Spawns agents in parallel (respects dependencies)
    -   Human checkpoints: alignment decisions, blocker resolution, verification review
    -   Use when you want autonomous end-to-end milestone execution

-   **`/handoff-milestone`** - Transition from completed milestone to next with gap analysis
    -   Validates current milestone completion (all tasks done, quality passed)
    -   Identifies gaps: missing scripts, docs, configs, tests needed for next milestone
    -   Creates gap-filling tasks in current milestone
    -   Extracts decisions and patterns, creates HANDOFF_NOTES.md
    -   Use when milestone reaches 100% before starting next milestone

-   **`/update-project-state`** - Update CLAUDE.md with current development phase and priorities
    -   Analyzes entire roadmap to determine development phase
    -   Extracts recent achievements (last 2 weeks) from completed milestones
    -   Identifies current focus and next priorities
    -   Use weekly, after milestone completion, or when phase transitions

**Daily Workflow**: `/project-status` → pick ready task → implement → `/verify` → complete

**Planning Workflow**: `/plan-milestone` → `/align-milestone` → adjust scope → `/preplan`

**Autonomous Workflow**: `/ship-milestone milestone-6` → [11-state FSM with alignment gates] → milestone complete

**Milestone Completion Workflow**: milestone reaches 100% → `/handoff-milestone milestone-N milestone-N+1` → fill gaps → `/update-project-state` → start next milestone

## tui-styles-cli: Project Management CLI

**Installation**: One-time setup (no sudo required)
```bash
cd tools/tui-styles-cli && make install
# Installs to ~/.local/bin (add to PATH if needed)
```

**Purpose**: All roadmap operations (task status updates, milestone management, comments) MUST use tui-styles-cli for consistency and file safety.

**Core Operations**:

**Task Management**:
- `tui-styles-cli task get <id>` - View task details
- `tui-styles-cli task update <id> --status <status>` - Change status (auto-renames file)
- `tui-styles-cli task list [--filter <status>]` - List all tasks
- `tui-styles-cli task search <query>` - Full-text search
- `tui-styles-cli task create <milestone> <phase> <number> <title>` - Create with template
- `tui-styles-cli task edit <id>` - Open in $EDITOR

**Milestone Management**:
- `tui-styles-cli milestone list` - Show all milestones with progress
- `tui-styles-cli milestone info <name>` - View milestone details
- `tui-styles-cli milestone tasks <name>` - List tasks by phase
- `tui-styles-cli milestone create <name> <title>` - Create with 3 phases
- `tui-styles-cli milestone complete <name>` - Mark as complete

**Comments** (for blockers/decisions):
- `tui-styles-cli comment create <id> "<text>" --author <name> [--needs-addressing]`
- `tui-styles-cli comment list <id>` - View all comments

**Project Status**:
- `tui-styles-cli project status` - Overall dashboard with completion %

**Valid Statuses**: pending → ready → in_progress → needs_review → needs_testing → needs_human_verification → complete (or blocked at any stage)

**When to Use**:
- ✅ Agents ALWAYS use CLI for task/milestone/comment operations
- ✅ Use for file renaming (status changes), creates, updates
- ✅ Provides file locking, validation, atomic operations
- ❌ Don't use for reading task content (use Read tool directly)
- ❌ Don't use for searching codebase (use Grep/Glob tools)

**Why CLI**: Ensures proper file naming, YAML frontmatter, status validation, prevents race conditions. All slash command agents rely on CLI for roadmap mutations.

**See**: tools/tui-styles-cli/README.md for complete documentation

## Project Guidelines (Read by AI Agents)

**CRITICAL**: All TUI Styles agents automatically read project guidelines to ensure consistency.

**Guidelines Files**:

-   **CODING_GUIDELINES.md** - Rust/Go/TypeScript patterns, testing standards, security requirements (MANDATORY reading)
-   **specs/wsp.md** - Complete protocol specification (CRDTs, WebRTC, message formats)
-   **plan/*.md** - Implementation plans for each component (agent, CLI, registry, desktop, mobile)

**How Agents Use Guidelines**:

-   Alignment agents: Validate tasks follow architectural patterns, flag violations
-   Pre-planning agents: Reference guidelines for library choices, block tasks that contradict standards
-   Implementation agents: Cite guidelines in code, follow established patterns

**Project-Specific Guidelines**:

-   **Rust Code**: All async operations use Tokio. Block filesystem/SQLite in `spawn_blocking`. Use `Arc<RwLock>` for shared state. Never block async runtime.
-   **CRDT Operations**: All operations must include vector clock. LWW resolution uses timestamp + peer_id tiebreaker. Comments are append-only (never LWW).
-   **WebRTC P2P**: DataChannels must be ordered. Implement keep-alive ping/pong every 30s. Monitor buffered_amount for backpressure.
-   **File Safety**: Atomic writes (write .tmp → fsync → rename). Debounce file watcher (200ms). Self-change ignore to prevent loops.
-   **Desktop/Mobile**: Thin clients render UISchema from agents. Never duplicate business logic. Agent is source of truth.

See plan/ and specs/ for detailed specifications.

## Quality Standards

All code must pass automated quality checks before committing.

**Pre-commit Validation**: Run quality checks via `make test` and `cargo clippy` (Rust), `golangci-lint` (Go), `npm run lint && npm test` (TypeScript).

Component-specific quality commands:

-   **Agent/CLI**: `cargo clippy -- -D warnings && cargo test && cargo fmt --check`
-   **Registry**: `golangci-lint run && go test ./...`
-   **Desktop**: `cd tui-styles-desktop && npm run lint && npm test`
-   **Mobile**: `cd tui-styles-mobile && npm run lint && npm test`

Quality checks include:

-   **Linting**: Zero warnings required
    -   Rust: `clippy` with deny warnings
    -   Go: `golangci-lint` with standard ruleset
    -   TypeScript: ESLint with strict rules

-   **Formatting**: Consistent code formatting enforced
    -   Rust: `rustfmt`
    -   Go: `gofmt`, `goimports`
    -   TypeScript: Prettier

-   **Tests**: Comprehensive test coverage
    -   Unit tests for CRDT logic (property-based tests with `proptest`)
    -   Integration tests for agent-to-agent sync
    -   E2E tests for CLI commands
    -   Performance benchmarks for search indexing (<500ms for 10k files)

-   **Security**: Automated security scanning
    -   No secrets in code (use environment variables)
    -   Dependencies scanned (`cargo audit`, `npm audit`, `go mod verify`)
    -   Capability verification on all agent operations
    -   HKDF secret derivation for access control

-   **Performance**: Performance budgets enforced
    -   Agent startup: <2s
    -   Search queries: <500ms for 10k files
    -   WebRTC connection: <3s establishment
    -   File watching: <500ms detection latency

-   **Code Quality**: Maintainability standards
    -   Keep functions under 50 lines where possible
    -   Avoid deeply nested logic (max 3 levels)
    -   Remove dead code and unused imports
    -   Document all public APIs with rustdoc/godoc/JSDoc

## Naming Conventions

**Folders/Files**: Use `lowercase-with-hyphens` for all directories and files (e.g., `tui-styles-agent`, `crdt-sync.rs`, `task-card.tsx`). Types/structs use `PascalCase` in code (e.g., `struct VectorClock`, `function TaskCard()`).

## Repository Structure

```
tui-styles/
├── tui-styles-agent/          - Rust daemon: file watch, CRDT sync, WebRTC, search, HTTP API
├── tui-styles-cli/            - Rust CLI: task management, agent control, sharing
├── tui-styles-lib/            - Rust shared library: protocol, CRDTs, types
├── tui-styles-registry/       - Go signaling server: WebSocket relay, peer discovery
├── tui-styles-desktop/        - Tauri + React: thin UI client
│   ├── src-tauri/       - Rust backend (WebRTC client)
│   └── src/             - React frontend (UISchema renderer)
├── tui-styles-mobile/         - React Native: iOS/Android UI client
├── specs/               - Protocol specifications
│   └── wsp.md           - Complete TUI Styles spec
├── plan/                - Implementation plans
│   ├── 00-master-plan.md
│   ├── 01-protocol-specification.md
│   ├── 02-agent-implementation.md
│   ├── 03-cli-implementation.md
│   ├── 04-registry-implementation.md
│   ├── 05-desktop-implementation.md
│   └── 06-mobile-implementation.md
├── docs/
│   ├── reference/       - Architecture decisions, roadmap structure
│   └── developer/       - Development setup, workflows
├── roadmap/             - Milestone and task structure
├── tools/
│   └── tui-styles-cli/        - Go CLI tool for roadmap management
├── reference-code/      - Production-tested code patterns
│   ├── next-app/        - Next.js 16 + React 19 (for tui-styles-desktop React frontend)
│   ├── go-api/          - Go HTTP API with Chi router (for tui-styles-registry)
│   └── go-worker/       - Go clean architecture worker (patterns for tui-styles-agent)
└── CODING_GUIDELINES.md
```

## Reference Code

The `reference-code/` directory contains production-tested code from real projects to serve as implementation references:

**next-app/** (Masquerade creator-studio-web)
- Next.js 16 with App Router, React 19, TypeScript
- shadcn/ui component architecture (Radix UI + Tailwind CSS)
- React Hook Form + Zod validation
- Atomic component organization (atoms/molecules/organisms)
- Comprehensive testing (Jest + React Testing Library)
- Quality tooling (ESLint 9, Prettier, knip, madge)
- **Use for**: tui-styles-desktop React frontend structure

**go-api/** (Masquerade creator-api)
- Go 1.25 HTTP service with Chi router
- OpenAPI specification-first design
- Structured logging with zerolog
- CORS, health checks, middleware patterns
- **Use for**: tui-styles-registry WebSocket server structure

**go-worker/** (Peach email-worker)
- Clean architecture (hexagonal/ports & adapters)
- Domain-driven design
- Repository pattern, use cases, dependency injection
- Comprehensive error handling
- **Use for**: Architectural patterns to adapt to Rust for tui-styles-agent

See `reference-code/README.md` for detailed usage guidance.

## Core Concepts

**Milestones**: Directories named `milestone-{N}/` containing phases and tasks. Each has `index.md` with goals, scope, success criteria.

**Phases**: Subdirectories within milestones (`phase-1-quick-wins/`, `phase-2-core-features/`, `phase-3-polish/`) grouping related tasks.

**Tasks**: Markdown files named `{number}_{name}_{status}.md` with structured content (purpose, acceptance criteria, technical approach, dependencies).

**Status Transitions**: File renames represent state changes detectable via file watching, updates broadcast in real-time.

**Metadata Extraction**: Parse frontmatter or structured markdown sections for dependencies, assigned agent, acceptance criteria checkboxes.

**AI-Readable Format**: Structured markdown designed for AI agents to read directly for project context without API calls.

## Before Committing to Git

**Commit Checklist**:

1. **Quality checks passed**:
    - Rust: `cargo clippy -- -D warnings && cargo test && cargo fmt --check`
    - Go: `golangci-lint run && go test ./...`
    - TypeScript: `npm run lint && npm test`

2. **Tests added**:
    - Unit tests for new CRDT operations
    - Integration tests for agent-to-agent sync
    - Property-based tests for conflict resolution
    - Performance benchmarks for indexing/search

3. **Security validated**:
    - No secrets committed (check .env files)
    - Dependencies scanned (`cargo audit`, `npm audit`)
    - Capability verification on operations
    - HKDF derivation for access control

4. **Performance verified**:
    - Agent startup <2s
    - Search <500ms for 10k files
    - WebRTC establishment <3s
    - File watching <500ms latency

5. **Documentation updated**:
    - Update specs/ if protocol changed
    - Update plan/ if architecture changed
    - Add rustdoc/godoc for public APIs
    - Update CHANGELOG.md for user-facing changes

## Rules

**Security**:

-   Never commit secrets (API keys, certificates, private keys, .env files with real values)
-   Encrypt secrets at rest (use OS keychain for desktop, Keychain/Keystore for mobile)
-   All WebRTC DataChannels use DTLS (built-in encryption)
-   Capability tokens (JWT) signed with HMAC-SHA256 from HKDF-derived keys
-   Verify capabilities on agent before applying operations

**Code Quality**:

-   Zero clippy/lint warnings
-   Keep functions under 50 lines where possible
-   Avoid deeply nested logic (max 3 levels)
-   Keep files under 500 lines - split into focused modules
-   Remove dead code and unused imports
-   Document all public APIs with rustdoc/godoc/JSDoc

**Testing**:

-   Write tests for all CRDT operations (property-based tests with `proptest`)
-   Integration tests for multi-agent sync
-   Unit tests for file watching, search indexing
-   Performance benchmarks for critical paths

**Performance**:

-   Never block Tokio async runtime (use `spawn_blocking` for sync I/O)
-   Debounce file watcher events (200ms window)
-   Use MessagePack for efficient serialization (<1ms encode/decode)
-   SQLite FTS5 with incremental indexing (<100ms per file update)

**Development Workflow**:

-   Update CHANGELOG.md for user-facing changes
-   Update specs/ if protocol changes
-   Use semantic commit messages (feat:, fix:, docs:, refactor:, test:, perf:)
-   Small, focused commits with clear descriptions

## Where to Find Help

**Protocol & Architecture**:

-   Complete specification: specs/wsp.md
-   Implementation plans: plan/*.md
-   CRDT theory: Section 5 of specs/wsp.md

**Technical Documentation**:

-   Agent architecture: plan/02-agent-implementation.md
-   CLI design: plan/03-cli-implementation.md
-   Registry server: plan/04-registry-implementation.md
-   Desktop app: plan/05-desktop-implementation.md
-   Mobile app: plan/06-mobile-implementation.md

**Reference Code** (Production-Tested Patterns):

-   Desktop React patterns: reference-code/next-app/ (Next.js 16 + React 19 + shadcn/ui)
-   Go HTTP/WebSocket server: reference-code/go-api/ (Chi router + OpenAPI)
-   Clean architecture patterns: reference-code/go-worker/ (ports/adapters + DDD)
-   Usage guide: reference-code/README.md

**Development**:

-   Getting started: README.md (when created)
-   Coding standards: CODING_GUIDELINES.md
-   Design system: DESIGN_SYSTEM.md
-   Roadmap structure: docs/reference/roadmap-structure.md (when created)

**Rust-Specific**:

-   Agent subsystems: tui-styles-agent/src/ (when created)
-   CRDT implementation: tui-styles-lib/src/crdt/ (when created)
-   WebRTC client: tui-styles-lib/src/webrtc/ (when created)
-   CLI commands: tui-styles-cli/src/ (when created)

**Workflow & Process**:

-   Roadmap planning: Use `/plan-milestone` command
-   Task pre-planning: Use `/preplan` command
-   Milestone execution: Use `/ship-milestone` command
-   Verification workflow: Use `/verify` command
