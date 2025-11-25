# WebSocket Client for Masquerade

A robust, production-ready WebSocket client for real-time event streaming in Masquerade applications. Features automatic reconnection, type-safe event handling, and connection state management.

## Features

- **Automatic Reconnection**: Exponential backoff with configurable delays (1s → 2s → 4s → 8s → max 30s)
- **Type-Safe Events**: Full TypeScript support with typed event handlers
- **Connection State Tracking**: Monitor connection states (connecting, connected, disconnecting, disconnected)
- **Heartbeat/Keepalive**: Automatic heartbeat messages to maintain connection
- **Event Listeners**: Clean API for registering, removing, and one-time event handlers
- **Auth Token Refresh**: Automatic token refresh on reconnection
- **Cross-Platform**: Works in browser and Node.js environments
- **Error Handling**: Comprehensive error events and recovery

## Installation

The WebSocket client is included in the Creator Studio web application. No additional installation required.

For Node.js environments, install the `ws` package:

```bash
npm install ws
```

## Quick Start

### Basic Usage

```typescript
import { WebSocketClient } from '@/lib/websocket-client';

// Create client instance
const client = new WebSocketClient(
  'ws://localhost:34076/ws',
  async () => {
    // Fetch JWT from your auth system
    return localStorage.getItem('jwt_token') || '';
  }
);

// Listen for events
client.on('image.generation.completed', (event) => {
  console.log('Image ready:', event.payload.image_url);
  console.log('Job ID:', event.payload.job_id);
});

// Connect to server
await client.connect();
```

### React Hook Example

```typescript
import { useEffect, useState } from 'react';
import { WebSocketClient } from '@/lib/websocket-client';
import type { ConnectionState } from '@/lib/websocket-types';

export function useWebSocket(url: string, getToken: () => Promise<string>) {
  const [client] = useState(() => new WebSocketClient(url, getToken));
  const [connectionState, setConnectionState] = useState<ConnectionState>('disconnected');

  useEffect(() => {
    // Listen for state changes
    client.on('state:change', (event) => {
      setConnectionState(event.newState);
    });

    // Connect
    client.connect();

    // Cleanup on unmount
    return () => {
      client.disconnect();
    };
  }, [client]);

  return { client, connectionState };
}

// Usage in component
function MyComponent() {
  const { client, connectionState } = useWebSocket(
    'ws://localhost:34076/ws',
    async () => localStorage.getItem('jwt_token') || ''
  );

  useEffect(() => {
    client.on('image.generation.completed', (event) => {
      console.log('Image ready:', event.payload.image_url);
    });
  }, [client]);

  return <div>Connection: {connectionState}</div>;
}
```

## Configuration Options

```typescript
const client = new WebSocketClient(url, getToken, {
  // Initial reconnection delay (default: 1000ms)
  initialReconnectDelay: 1000,

  // Maximum reconnection delay (default: 30000ms)
  maxReconnectDelay: 30000,

  // Heartbeat interval (default: 30000ms)
  heartbeatInterval: 30000,

  // Enable automatic reconnection (default: true)
  autoReconnect: true,

  // Enable debug logging (default: false)
  debug: true,
});
```

## Event Types

### Server Events

All server events implement the `BaseEvent` interface:

```typescript
interface BaseEvent {
  type: string;
  payload: Record<string, unknown>;
  timestamp: string;
}
```

#### Image Generation Events

```typescript
// Image generation started
client.on('image.generation.started', (event) => {
  console.log('Job started:', event.payload.job_id);
  console.log('User:', event.payload.user_id);
  console.log('Prompt:', event.payload.prompt);
});

// Image generation progress
client.on('image.generation.progress', (event) => {
  console.log('Progress:', event.payload.progress + '%');
  console.log('Stage:', event.payload.stage);
});

// Image generation completed
client.on('image.generation.completed', (event) => {
  console.log('Image URL:', event.payload.image_url);
  console.log('Metadata:', event.payload.metadata);
});

// Image generation failed
client.on('image.generation.failed', (event) => {
  console.error('Error:', event.payload.error);
  console.log('Retry allowed:', event.payload.retry_allowed);
});
```

#### System Events

```typescript
// System notifications
client.on('system.notification', (event) => {
  console.log('Severity:', event.payload.severity); // 'info' | 'warning' | 'error'
  console.log('Message:', event.payload.message);
  if (event.payload.action_url) {
    console.log('Action URL:', event.payload.action_url);
  }
});
```

### Client Events

```typescript
// Connection established
client.on('connected', () => {
  console.log('WebSocket connected');
});

// Connection closed
client.on('disconnected', () => {
  console.log('WebSocket disconnected');
});

// State changes
client.on('state:change', (event) => {
  console.log(`State: ${event.oldState} → ${event.newState}`);
});

// Errors
client.on('error', (event) => {
  console.error('Error:', event.error.message);
  console.log('Context:', event.context);
});
```

## API Reference

### Constructor

```typescript
new WebSocketClient(
  url: string,
  getToken: () => Promise<string>,
  options?: WebSocketClientOptions
)
```

### Methods

#### `connect(): Promise<void>`

Connect to the WebSocket server. Fetches auth token via `getToken()` and establishes connection.

```typescript
await client.connect();
```

#### `disconnect(): void`

Disconnect from the WebSocket server. Cancels any pending reconnections.

```typescript
client.disconnect();
```

#### `on<K>(eventType: K, handler: EventHandler<EventTypeMap[K]>): void`

Register an event listener.

```typescript
client.on('image.generation.completed', (event) => {
  console.log(event.payload.image_url);
});
```

#### `off<K>(eventType: K, handler: EventHandler<EventTypeMap[K]>): void`

Remove an event listener.

```typescript
const handler = (event) => console.log(event);
client.on('image.generation.completed', handler);
client.off('image.generation.completed', handler);
```

#### `once<K>(eventType: K, handler: EventHandler<EventTypeMap[K]>): void`

Register a one-time event listener (automatically removed after first call).

```typescript
client.once('connected', () => {
  console.log('Connected!');
});
```

#### `send(event: BaseEvent): void`

Send a message to the server. Throws error if not connected.

```typescript
client.send({
  type: 'custom.event',
  payload: { data: 'value' },
  timestamp: new Date().toISOString(),
});
```

### Properties

#### `state: ConnectionState`

Current connection state (read-only).

```typescript
console.log(client.state); // 'connected' | 'connecting' | 'disconnecting' | 'disconnected'
```

## Reconnection Behavior

The client automatically reconnects on connection loss using exponential backoff:

1. **First attempt**: 1 second delay
2. **Second attempt**: 2 seconds delay
3. **Third attempt**: 4 seconds delay
4. **Fourth attempt**: 8 seconds delay
5. **Fifth+ attempts**: 30 seconds delay (capped at `maxReconnectDelay`)

Reconnection attempts reset to 0 on successful connection.

### Disabling Auto-Reconnect

```typescript
const client = new WebSocketClient(url, getToken, {
  autoReconnect: false,
});
```

### Manual Reconnection

```typescript
// Disconnect will cancel any pending reconnections
client.disconnect();

// Manual reconnect
await client.connect();
```

## Auth Token Refresh

The client automatically refreshes auth tokens on each connection attempt:

```typescript
const client = new WebSocketClient(
  'ws://localhost:34076/ws',
  async () => {
    // Fetch fresh token from your auth system
    const response = await fetch('/api/auth/token');
    const data = await response.json();
    return data.token;
  }
);
```

This ensures expired tokens are replaced with fresh ones during reconnection.

## Error Handling

### Listening for Errors

```typescript
client.on('error', (event) => {
  console.error('WebSocket error:', event.error.message);
  console.log('Context:', event.context); // 'connection' | 'send' | 'parse'
  console.log('Timestamp:', event.timestamp);
});
```

### Common Error Scenarios

#### Connection Failures

```typescript
client.on('error', (event) => {
  if (event.context === 'connection') {
    // Server unreachable, invalid token, etc.
    console.error('Connection failed:', event.error);
  }
});
```

#### Send Failures

```typescript
try {
  client.send(message);
} catch (error) {
  console.error('Failed to send message:', error);
}
```

#### Malformed Messages

```typescript
client.on('error', (event) => {
  if (event.context === 'parse') {
    // Received invalid JSON from server
    console.error('Parse error:', event.error);
  }
});
```

## Node.js Usage

For Node.js environments, install the `ws` package and use a polyfill:

```bash
npm install ws
```

```typescript
import { WebSocket } from 'ws';
import { WebSocketClient } from './lib/websocket-client';

// Polyfill global WebSocket
(global as any).WebSocket = WebSocket;

// Use client normally
const client = new WebSocketClient(
  'ws://localhost:34076/ws',
  async () => process.env.JWT_TOKEN || ''
);

await client.connect();
```

## Best Practices

### Memory Management

Always remove event listeners when components unmount to prevent memory leaks:

```typescript
useEffect(() => {
  const handler = (event) => console.log(event);
  client.on('image.generation.completed', handler);

  return () => {
    client.off('image.generation.completed', handler);
  };
}, [client]);
```

### Single Client Instance

Create a single WebSocket client instance and share it across your application:

```typescript
// lib/websocket.ts
export const wsClient = new WebSocketClient(
  process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:34076/ws',
  async () => {
    // Your token fetching logic
    return localStorage.getItem('jwt_token') || '';
  }
);

// Usage in components
import { wsClient } from '@/lib/websocket';

wsClient.on('image.generation.completed', handler);
```

### Error Boundaries

Wrap WebSocket operations in error boundaries for graceful degradation:

```typescript
try {
  await client.connect();
} catch (error) {
  console.error('Failed to establish WebSocket connection');
  // Fall back to polling or show offline UI
}
```

### Heartbeat Configuration

Adjust heartbeat interval based on your infrastructure:

```typescript
// Shorter interval for aggressive keepalive
const client = new WebSocketClient(url, getToken, {
  heartbeatInterval: 15000, // 15 seconds
});

// Longer interval to reduce bandwidth
const client = new WebSocketClient(url, getToken, {
  heartbeatInterval: 60000, // 60 seconds
});
```

## Debugging

Enable debug logging to troubleshoot connection issues:

```typescript
const client = new WebSocketClient(url, getToken, {
  debug: true,
});
```

Debug logs include:

- Connection state transitions
- Reconnection attempts with delays
- Heartbeat messages
- Message send/receive events
- Error details

## Examples

### Complete React Component

```typescript
import { useEffect, useState } from 'react';
import { WebSocketClient } from '@/lib/websocket-client';
import type { ImageGenerationCompletedEvent } from '@/lib/websocket-types';

export function ImageGenerationMonitor() {
  const [client] = useState(
    () =>
      new WebSocketClient(
        'ws://localhost:34076/ws',
        async () => localStorage.getItem('jwt_token') || ''
      )
  );
  const [images, setImages] = useState<string[]>([]);
  const [connectionState, setConnectionState] = useState('disconnected');

  useEffect(() => {
    // Connection state listener
    client.on('state:change', (event) => {
      setConnectionState(event.newState);
    });

    // Image completion listener
    const handleImageComplete = (event: ImageGenerationCompletedEvent) => {
      setImages((prev) => [...prev, event.payload.image_url]);
    };
    client.on('image.generation.completed', handleImageComplete);

    // Error listener
    client.on('error', (event) => {
      console.error('WebSocket error:', event.error);
    });

    // Connect
    client.connect();

    // Cleanup
    return () => {
      client.off('image.generation.completed', handleImageComplete);
      client.disconnect();
    };
  }, [client]);

  return (
    <div>
      <div>Status: {connectionState}</div>
      <div>Images: {images.length}</div>
      {images.map((url, i) => (
        <img key={i} src={url} alt={`Generated ${i}`} />
      ))}
    </div>
  );
}
```

## Browser Compatibility

- Chrome 43+
- Firefox 11+
- Safari 5.1+
- Edge (all versions)
- Opera 12.1+

WebSocket is supported in all modern browsers. IE11 and older browsers are not supported.

## License

See the main Masquerade repository for license information.
