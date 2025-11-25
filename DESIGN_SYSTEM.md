# TUI Styles Design System

A comprehensive design system for the TUI Styles Protocol platform, built on modern React patterns for both desktop (Tauri) and mobile (React Native). This system ensures consistency across all UI clients while maintaining platform-specific best practices.

## Design Philosophy

**Core Principle**: Agent-driven UI that adapts to schema. Thin clients render dynamic layouts based on UISchema from agents. Design must work across desktop (Tauri) and mobile (React Native) with shared rendering logic but platform-specific interactions.

**Key Tenets**:
- **Agent-Driven**: UI structure determined by agent's UISchema (columns, fields, capabilities)
- **Schema-First**: Components render based on received configuration, not hardcoded structures
- **Cross-Platform**: Shared business logic, platform-specific UI primitives
- **Accessibility by Default**: WCAG 2.1 AA compliance, keyboard navigation, screen readers
- **Performance-Conscious**: 60fps on mobile, <100ms interactions, efficient re-renders
- **Offline-First**: Queue operations when disconnected, optimistic UI updates

## Technology Stack

### Desktop (Tauri + React)
- **Framework**: Tauri 2.0 (Rust backend + web frontend)
- **UI**: React 18 with TypeScript, Server Components not used (SPA architecture)
- **Component Library**: Radix UI primitives (headless, accessible)
- **Styling**: Tailwind CSS for utility-first styling
- **State**: Zustand for UI state (agent connection, schema, local preferences)
- **Icons**: Lucide React (primary icon system)
- **Animation**: Framer Motion for transitions
- **Drag & Drop**: @dnd-kit for task movement

### Mobile (React Native)
- **Framework**: React Native 0.72+ (CLI workflow, native modules)
- **Navigation**: React Navigation 6+ (Native Stack)
- **UI**: React Native Paper (Material Design 3) or custom components
- **State**: Zustand (shared with desktop logic)
- **Icons**: react-native-vector-icons or Lucide React Native
- **Animation**: React Native Reanimated for 60fps animations
- **Gestures**: React Native Gesture Handler for native feel
- **Lists**: FlashList for performance with large task lists

### Shared Logic
- **Schema Renderer**: Shared TypeScript code to interpret UISchema
- **Field Mapping**: `FieldType` → Component mapping works on both platforms
- **Validation**: Zod schemas work identically on desktop/mobile
- **State Management**: Zustand stores can be reused

## Design Tokens

### Color System

Based on HSL for programmatic dark mode and theme variations.

**Brand Colors**:
```typescript
// shared/design-tokens/colors.ts (used by both desktop and mobile)
export const colors = {
  brand: {
    primary: 'hsl(262, 83%, 58%)',    // Vibrant purple (collaboration)
    secondary: 'hsl(200, 98%, 39%)',  // Electric blue (real-time sync)
    accent: 'hsl(142, 76%, 36%)',     // Green (success, completion)
  },

  // Semantic colors
  success: 'hsl(142, 76%, 36%)',
  warning: 'hsl(38, 92%, 50%)',
  error: 'hsl(0, 84%, 60%)',
  info: 'hsl(199, 89%, 48%)',

  // Neutral grays (light mode)
  gray: {
    50: 'hsl(0, 0%, 98%)',
    100: 'hsl(0, 0%, 96%)',
    200: 'hsl(0, 0%, 90%)',
    300: 'hsl(0, 0%, 83%)',
    400: 'hsl(0, 0%, 64%)',
    500: 'hsl(0, 0%, 45%)',
    600: 'hsl(0, 0%, 32%)',
    700: 'hsl(0, 0%, 25%)',
    800: 'hsl(0, 0%, 15%)',
    900: 'hsl(0, 0%, 9%)',
    950: 'hsl(0, 0%, 4%)',
  },
}
```

**Usage**:
- Primary: Active column, selected tasks, key CTAs
- Secondary: Real-time sync indicators, peer presence
- Accent: Task completion, success states
- Semantic: Status feedback, alerts, confirmations

### Typography

**Font Families**:
```typescript
// shared/design-tokens/typography.ts
export const fonts = {
  // Desktop (Tauri - web fonts)
  desktop: {
    sans: ['Inter', 'system-ui', 'sans-serif'],
    display: ['Cal Sans', 'Inter', 'sans-serif'],
    mono: ['JetBrains Mono', 'Consolas', 'monospace'],
  },

  // Mobile (React Native - system fonts)
  mobile: {
    ios: {
      sans: 'System',               // San Francisco
      display: 'System',
      mono: 'Courier',
    },
    android: {
      sans: 'Roboto',
      display: 'Roboto',
      mono: 'monospace',
    },
  },
}
```

**Type Scale** (1.250 - Major Third):
```typescript
export const fontSizes = {
  xs: 12,      // Captions, metadata
  sm: 14,      // Secondary text, labels
  base: 16,    // Body text
  lg: 20,      // Large body, subheadings
  xl: 24,      // H4
  '2xl': 30,   // H3
  '3xl': 36,   // H2
  '4xl': 48,   // H1
}

// React Native uses numeric values directly
// Tailwind CSS converts to rem
```

### Spacing System

Based on 4px grid for consistency.

```typescript
export const spacing = {
  0: 0,
  1: 4,
  2: 8,
  3: 12,
  4: 16,
  5: 20,
  6: 24,
  8: 32,
  10: 40,
  12: 48,
  16: 64,
  20: 80,
  24: 96,
}
```

**Component Spacing**:
- **Micro**: 4-8px - Icon-to-text, inline elements
- **Small**: 12-16px - Component internal padding
- **Medium**: 24-32px - Card padding, component margins
- **Large**: 48-64px - Section spacing, layout gaps

## Component Architecture

### Desktop (Tauri) Structure

```
tui-styles-desktop/
├── src-tauri/                 # Rust backend
│   ├── main.rs               # Tauri entry point
│   ├── webrtc.rs             # WebRTC client (shared with tui-styles-agent)
│   ├── protocol.rs           # Message handling
│   └── commands.rs           # Tauri commands (IPC)
│
├── src/                      # React frontend
│   ├── App.tsx
│   ├── components/
│   │   ├── ui/               # Radix UI wrappers
│   │   │   ├── button.tsx
│   │   │   ├── card.tsx
│   │   │   ├── dialog.tsx
│   │   │   └── input.tsx
│   │   ├── atoms/            # Basic building blocks
│   │   │   ├── Icon.tsx
│   │   │   ├── Badge.tsx
│   │   │   └── Avatar.tsx
│   │   ├── molecules/        # Composed components
│   │   │   ├── TaskCard.tsx
│   │   │   ├── ColumnHeader.tsx
│   │   │   └── SearchBar.tsx
│   │   ├── organisms/        # Complex compositions
│   │   │   ├── Board.tsx             # Main board view
│   │   │   ├── Column.tsx            # Droppable column
│   │   │   ├── TaskDetail.tsx        # Task modal
│   │   │   └── ConnectionStatus.tsx  # WebRTC status
│   │   └── schema/           # Schema-driven components
│   │       ├── SchemaRenderer.tsx    # Interprets UISchema
│   │       ├── FieldRenderer.tsx     # Renders fields based on type
│   │       └── DynamicForm.tsx       # Form from schema
│   ├── stores/
│   │   ├── connection.ts     # WebRTC connection state
│   │   ├── schema.ts         # UISchema from agent
│   │   └── tasks.ts          # Task data
│   └── lib/
│       ├── tauri.ts          # Tauri command wrappers
│       └── messagepack.ts    # MessagePack encoding
```

### Mobile (React Native) Structure

```
tui-styles-mobile/
├── ios/                      # iOS native code
├── android/                  # Android native code
├── src/
│   ├── App.tsx
│   ├── navigation/
│   │   └── RootNavigator.tsx
│   ├── screens/
│   │   ├── BoardScreen.tsx         # Main board (swipeable columns)
│   │   ├── TaskDetailScreen.tsx    # Task detail view
│   │   ├── SearchScreen.tsx        # Search interface
│   │   └── SettingsScreen.tsx      # Connection settings
│   ├── components/
│   │   ├── atoms/
│   │   │   ├── Icon.tsx
│   │   │   ├── Badge.tsx
│   │   │   └── Avatar.tsx
│   │   ├── molecules/
│   │   │   ├── TaskCard.tsx        # Shared with desktop (styles differ)
│   │   │   ├── ColumnSwiper.tsx    # Swipeable column view
│   │   │   └── SearchBar.tsx
│   │   └── schema/
│   │       ├── SchemaRenderer.tsx  # Shared logic from desktop
│   │       ├── FieldRenderer.tsx   # Platform-specific components
│   │       └── DynamicForm.tsx
│   ├── stores/               # Shared with desktop (Zustand)
│   ├── hooks/
│   │   ├── useWebRTC.ts            # WebRTC connection hook
│   │   ├── useOfflineQueue.ts      # Offline operation queue
│   │   └── useBackgroundSync.ts    # Background sync handling
│   └── native/
│       ├── WebRTCModule.ts         # Native WebRTC bridge
│       └── NotificationModule.ts   # Push notifications
```

### Shared Schema Components

Both desktop and mobile use the same schema rendering logic:

```typescript
// shared/schema/types.ts
export interface UISchema {
  version: string
  protocol_version: string
  project: ProjectInfo
  flows: FlowColumn[]
  capabilities: UserCapabilities
  features: string[]
  custom_fields: CustomField[]
  theme?: ThemeConfig
}

export interface FlowColumn {
  id: string
  name: string
  order: number
  color?: string
  icon?: string
  max_tasks?: number
}

export interface CustomField {
  id: string
  name: string
  field_type: FieldType
  required: boolean
  options?: string[]
  default?: string
}

export type FieldType = 'text' | 'select' | 'number' | 'date' | 'checkbox' | 'multi_select'
```

**Schema Renderer Pattern**:

```typescript
// shared/schema/SchemaRenderer.tsx
export function SchemaRenderer({ schema }: { schema: UISchema }) {
  // Sort columns by order
  const sortedColumns = useMemo(
    () => [...schema.flows].sort((a, b) => a.order - b.order),
    [schema.flows]
  )

  return (
    <BoardLayout>
      {sortedColumns.map((column) => (
        <Column
          key={column.id}
          id={column.id}
          name={column.name}
          color={column.color}
          icon={column.icon}
          maxTasks={column.max_tasks}
          canDrop={schema.capabilities.can_move_tasks}
          canCreate={schema.capabilities.can_create_tasks}
        />
      ))}
    </BoardLayout>
  )
}

// Desktop uses div + Tailwind
// Mobile uses View + StyleSheet
```

**Field Renderer Pattern**:

```typescript
// shared/schema/FieldRenderer.tsx
export function FieldRenderer({ field, value, onChange }: FieldRendererProps) {
  switch (field.field_type) {
    case 'text':
      return <TextInput value={value} onChange={onChange} />
    case 'select':
      return <Select options={field.options} value={value} onChange={onChange} />
    case 'number':
      return <NumberInput value={value} onChange={onChange} />
    case 'date':
      return <DatePicker value={value} onChange={onChange} />
    case 'checkbox':
      return <Checkbox checked={value} onChange={onChange} />
    case 'multi_select':
      return <MultiSelect options={field.options} value={value} onChange={onChange} />
    default:
      return <TextInput value={value} onChange={onChange} />
  }
}

// Desktop: Radix UI components
// Mobile: React Native Paper or custom native components
```

## Desktop-Specific Patterns

### Drag & Drop with @dnd-kit

```typescript
// tui-styles-desktop/src/components/organisms/Board.tsx
import { DndContext, closestCenter, PointerSensor, useSensor, useSensors } from '@dnd-kit/core'
import { SortableContext, verticalListSortingStrategy } from '@dnd-kit/sortable'

export function Board({ schema, tasks }: BoardProps) {
  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 8, // Prevent accidental drags
      },
    })
  )

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event
    if (!over) return

    const taskId = active.id as string
    const toColumn = over.id as string

    // Send to agent via Tauri command
    invoke('move_task', { taskId, toColumn })
  }

  return (
    <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
      <div className="flex gap-4 overflow-x-auto p-4">
        {schema.flows.map((column) => (
          <Column key={column.id} column={column} tasks={tasks.filter((t) => t.column === column.id)} />
        ))}
      </div>
    </DndContext>
  )
}
```

### Tauri IPC Commands

```typescript
// tui-styles-desktop/src/lib/tauri.ts
import { invoke } from '@tauri-apps/api/tauri'
import { listen } from '@tauri-apps/api/event'

// Request UISchema from agent
export async function getSchema(): Promise<UISchema> {
  return await invoke('get_schema')
}

// Send task operation to agent
export async function moveTask(taskId: string, toColumn: string): Promise<void> {
  await invoke('move_task', { taskId, toColumn })
}

// Listen for real-time updates from agent
export function onStateUpdate(callback: (data: any) => void) {
  return listen('state_update', (event) => {
    callback(event.payload)
  })
}
```

## Mobile-Specific Patterns

### Swipeable Column View

```typescript
// tui-styles-mobile/src/components/molecules/ColumnSwiper.tsx
import { useState } from 'react'
import { Dimensions, FlatList, View } from 'react-native'
import { Gesture, GestureDetector } from 'react-native-gesture-handler'
import Animated, { useSharedValue, useAnimatedStyle, withSpring } from 'react-native-reanimated'

const { width: SCREEN_WIDTH } = Dimensions.get('window')

export function ColumnSwiper({ columns, tasks }: ColumnSwiperProps) {
  const [currentIndex, setCurrentIndex] = useState(0)
  const translateX = useSharedValue(0)

  const panGesture = Gesture.Pan()
    .onUpdate((event) => {
      translateX.value = event.translationX - currentIndex * SCREEN_WIDTH
    })
    .onEnd((event) => {
      const shouldMoveRight = event.translationX < -SCREEN_WIDTH / 3
      const shouldMoveLeft = event.translationX > SCREEN_WIDTH / 3

      if (shouldMoveRight && currentIndex < columns.length - 1) {
        setCurrentIndex(currentIndex + 1)
      } else if (shouldMoveLeft && currentIndex > 0) {
        setCurrentIndex(currentIndex - 1)
      }

      translateX.value = withSpring(-currentIndex * SCREEN_WIDTH)
    })

  const animatedStyle = useAnimatedStyle(() => ({
    transform: [{ translateX: translateX.value }],
  }))

  return (
    <GestureDetector gesture={panGesture}>
      <Animated.View style={[{ flexDirection: 'row' }, animatedStyle]}>
        {columns.map((column) => (
          <View key={column.id} style={{ width: SCREEN_WIDTH }}>
            <ColumnView column={column} tasks={tasks.filter((t) => t.column === column.id)} />
          </View>
        ))}
      </Animated.View>
    </GestureDetector>
  )
}
```

### Offline Queue Pattern

```typescript
// tui-styles-mobile/src/hooks/useOfflineQueue.ts
import { useEffect } from 'react'
import AsyncStorage from '@react-native-async-storage/async-storage'
import { useWebRTC } from './useWebRTC'

interface QueuedOperation {
  id: string
  timestamp: number
  message: TUI StylesMessage
  retries: number
}

const QUEUE_KEY = '@tui-styles:offline_queue'

export function useOfflineQueue() {
  const { isConnected, sendMessage } = useWebRTC()

  // When connection restores, flush queue
  useEffect(() => {
    if (isConnected) {
      flushQueue()
    }
  }, [isConnected])

  const addToQueue = async (message: TUI StylesMessage) => {
    const queue = await getQueue()
    const operation: QueuedOperation = {
      id: generateId(),
      timestamp: Date.now(),
      message,
      retries: 0,
    }
    queue.push(operation)
    await AsyncStorage.setItem(QUEUE_KEY, JSON.stringify(queue))
  }

  const flushQueue = async () => {
    const queue = await getQueue()
    const remaining: QueuedOperation[] = []

    for (const op of queue) {
      try {
        await sendMessage(op.message)
      } catch (error) {
        op.retries++
        if (op.retries < 5) {
          remaining.push(op)
        }
      }
    }

    await AsyncStorage.setItem(QUEUE_KEY, JSON.stringify(remaining))
  }

  const getQueue = async (): Promise<QueuedOperation[]> => {
    const data = await AsyncStorage.getItem(QUEUE_KEY)
    return data ? JSON.parse(data) : []
  }

  return { addToQueue, flushQueue }
}
```

## Theme System

### Dark Mode (Desktop)

```typescript
// tui-styles-desktop/src/App.tsx
import { useEffect, useState } from 'react'

export function App() {
  const [theme, setTheme] = useState<'light' | 'dark'>('light')

  useEffect(() => {
    // Listen for system theme changes
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
    setTheme(mediaQuery.matches ? 'dark' : 'light')

    const handler = (e: MediaQueryListEvent) => {
      setTheme(e.matches ? 'dark' : 'light')
    }

    mediaQuery.addEventListener('change', handler)
    return () => mediaQuery.removeEventListener('change', handler)
  }, [])

  return (
    <div className={theme}>
      <div className="min-h-screen bg-background text-foreground">{/* App content */}</div>
    </div>
  )
}
```

### Dark Mode (Mobile)

```typescript
// tui-styles-mobile/src/App.tsx
import { useColorScheme } from 'react-native'
import { MD3DarkTheme, MD3LightTheme, PaperProvider } from 'react-native-paper'

export function App() {
  const colorScheme = useColorScheme() // 'light' or 'dark'
  const theme = colorScheme === 'dark' ? MD3DarkTheme : MD3LightTheme

  return <PaperProvider theme={theme}>{/* App content */}</PaperProvider>
}
```

## Performance Optimization

### Desktop

- **React.memo** for TaskCard to prevent re-rendering entire board on single task update
- **useMemo** for sorted columns and filtered tasks
- **Virtual scrolling** for large task lists (react-window)
- **Code splitting** with React.lazy for modals and detail views

### Mobile

- **FlashList** instead of FlatList for better blank area handling
- **React Native Reanimated** for 60fps animations (runs on UI thread)
- **Optimistic UI** updates with rollback on error
- **Image caching** with react-native-fast-image for avatars
- **AsyncStorage** batching for offline queue writes

## Accessibility

### Desktop (WCAG 2.1 AA)

- Keyboard navigation (Tab, Enter, Escape, Arrow keys)
- Focus indicators on all interactive elements
- ARIA labels for screen readers
- Semantic HTML (headings, landmarks, lists)
- Color contrast ratio >4.5:1 for text

### Mobile

- VoiceOver (iOS) and TalkBack (Android) support
- Accessible touch targets (minimum 44×44 points)
- Dynamic type support (respect system font size)
- Reduce motion support (disable animations if enabled)
- Haptic feedback for interactions

## Testing

### Desktop (Vitest + Playwright)

```typescript
// tui-styles-desktop/src/components/molecules/__tests__/TaskCard.test.tsx
import { render, screen } from '@testing-library/react'
import { TaskCard } from '../TaskCard'

describe('TaskCard', () => {
  it('renders task information', () => {
    render(<TaskCard task={{ id: '1', title: 'Test Task', column: 'backlog' }} />)
    expect(screen.getByText('Test Task')).toBeInTheDocument()
  })
})
```

### Mobile (Jest + React Native Testing Library)

```typescript
// tui-styles-mobile/src/components/molecules/__tests__/TaskCard.test.tsx
import { render } from '@testing-library/react-native'
import { TaskCard } from '../TaskCard'

describe('TaskCard', () => {
  it('renders task information', () => {
    const { getByText } = render(<TaskCard task={{ id: '1', title: 'Test Task', column: 'backlog' }} />)
    expect(getByText('Test Task')).toBeTruthy()
  })
})
```

## Resources

### Documentation
- [Tauri Documentation](https://tauri.app/v1/guides/)
- [Radix UI Primitives](https://www.radix-ui.com/primitives)
- [React Native Paper](https://reactnativepaper.com/)
- [Lucide Icons](https://lucide.dev/icons)
- [Framer Motion](https://www.framer.com/motion/)
- [React Native Reanimated](https://docs.swmansion.com/react-native-reanimated/)

### Tools
- [Realtime Colors](https://realtimecolors.com/) - Preview color schemes
- [Contrast Checker](https://webaim.org/resources/contrastchecker/) - WCAG validation
- [React DevTools](https://react.dev/learn/react-developer-tools) - Performance profiling
