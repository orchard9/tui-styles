# Flux Protocol Reference Code

This directory contains reference implementations from proven production codebases to serve as templates for Flux Protocol components.

## Structure

### next-app/
**Source**: Masquerade creator-studio-web  
**Technology**: Next.js 16, React 19, TypeScript, Radix UI, Tailwind CSS  
**Purpose**: Reference for flux-desktop (Tauri React frontend)

**Key Patterns**:
- Modern Next.js 16 with App Router
- shadcn/ui component architecture (Radix UI primitives)
- Tailwind CSS with design tokens
- React Hook Form + Zod validation
- Storybook for component documentation
- Comprehensive testing (Jest + React Testing Library)
- Quality tooling (ESLint, Prettier, TypeScript, knip for dead code)

**Use For**:
- Desktop React frontend structure
- Component organization (atoms/molecules/organisms)
- Form patterns with validation
- Theme system (light/dark mode)
- Testing patterns

### go-api/
**Source**: Masquerade creator-api  
**Technology**: Go 1.25, Chi router, OpenAPI, zerolog  
**Purpose**: Reference for flux-registry (WebSocket signaling server)

**Key Patterns**:
- Clean API structure with Chi router
- OpenAPI specification-first design
- Structured logging with zerolog
- CORS handling
- Health check endpoints
- Docker containerization
- Makefile automation

**Use For**:
- Go HTTP/WebSocket server structure
- API endpoint organization
- Middleware patterns
- Testing patterns
- Build and deployment scripts

### go-worker/
**Source**: Peach email-worker  
**Technology**: Go, Clean Architecture (ports & adapters), domain-driven design  
**Purpose**: Reference for flux-agent background processing patterns

**Key Patterns**:
- Hexagonal architecture (ports/adapters)
- Domain-driven design
- Repository pattern
- Use cases / service layer
- Worker/background job patterns
- Configuration management
- Comprehensive error handling
- Integration testing

**Use For**:
- flux-agent daemon architecture
- Background processing patterns
- Clean separation of concerns
- Domain modeling
- Repository patterns for SQLite

## How to Use This Reference Code

### For Desktop (Tauri + React)

1. **Component Structure**: Use next-app/ component organization
   ```
   components/
   ├── ui/              # Radix UI wrappers (from shadcn)
   ├── atoms/           # Basic components
   ├── molecules/       # Composed components
   └── organisms/       # Complex compositions
   ```

2. **Forms**: Adapt React Hook Form + Zod patterns
3. **Theme**: Use Tailwind CSS token system
4. **Testing**: Follow Jest + React Testing Library patterns

### For Registry (Go WebSocket Server)

1. **Project Structure**: Use go-api/ layout
   ```
   cmd/server/          # Entry point
   internal/api/        # HTTP handlers
   internal/config/     # Configuration
   pkg/models/          # Shared types
   ```

2. **Router**: Chi router with middleware chain
3. **Docs**: OpenAPI specification
4. **Logging**: Structured logging with zerolog

### For Agent (Rust Daemon)

1. **Architecture**: Adapt go-worker/ clean architecture to Rust
   ```
   Domain Layer (pure logic)
     ↓
   Use Cases (application logic)
     ↓
   Adapters (file system, WebRTC, SQLite)
   ```

2. **Patterns**:
   - Repository trait for SQLite access
   - Use case traits for business logic
   - Clean separation of concerns
   - Error handling with Result types

## Important Notes

- **Don't Copy Blindly**: These are references, not templates to copy-paste
- **Adapt to Rust**: flux-agent is Rust, not Go - adapt patterns, don't translate
- **Thin Client**: Desktop app should be simpler than next-app (agent-driven UI)
- **Protocol-First**: Registry should be simpler than go-api (just signaling)

## Version Information

- **Next.js**: 16.0.3 (App Router, React 19)
- **Go**: 1.25.4 (creator-api), email-worker Go version
- **Radix UI**: Latest primitives
- **Tailwind CSS**: v4

## Dependencies of Note

### Next.js App
- **UI**: Radix UI primitives, Lucide icons, Tailwind CSS
- **Forms**: React Hook Form, Zod validation
- **Testing**: Jest, React Testing Library, Playwright
- **Quality**: ESLint 9, Prettier, TypeScript, knip, madge

### Go API
- **Router**: Chi v5
- **Logging**: zerolog
- **CORS**: rs/cors
- **Testing**: testify

### Go Worker
- (Check go.mod when copied)

## Reference Documentation

- **Masquerade CLAUDE.md**: Project structure and workflows
- **Masquerade DESIGN_SYSTEM.md**: UI component patterns
- **Peach ARCHITECTURE.md**: Clean architecture principles

## Maintenance

This reference code is a snapshot. Do not modify it. If patterns need updating:
1. Document learnings in Flux implementation
2. Update Flux CODING_GUIDELINES.md
3. Do not backport to reference code
