/**
 * WebSocket Client for Masquerade Real-Time Events
 *
 * Provides a robust WebSocket client with automatic reconnection,
 * type-safe event handling, and connection state management.
 *
 * @example
 * ```typescript
 * const client = new WebSocketClient(
 *   'ws://localhost:34076/ws',
 *   async () => localStorage.getItem('jwt_token') || ''
 * );
 *
 * client.on('image.generation.completed', (event) => {
 *   console.log('Image ready:', event.payload.image_url);
 * });
 *
 * await client.connect();
 * ```
 */

import type {
  BaseEvent,
  ConnectionState,
  EventHandler,
  EventTypeMap,
  StateChangeEvent,
  ErrorEvent,
} from './websocket-types';

/**
 * WebSocket client options
 */
export interface WebSocketClientOptions {
  /**
   * Initial reconnection delay in milliseconds (default: 1000ms)
   */
  initialReconnectDelay?: number;

  /**
   * Maximum reconnection delay in milliseconds (default: 30000ms)
   */
  maxReconnectDelay?: number;

  /**
   * Heartbeat interval in milliseconds (default: 30000ms)
   */
  heartbeatInterval?: number;

  /**
   * Enable automatic reconnection (default: true)
   */
  autoReconnect?: boolean;

  /**
   * Enable debug logging (default: false)
   */
  debug?: boolean;
}

/**
 * WebSocket client for real-time event streaming
 */
export class WebSocketClient {
  private ws: WebSocket | null = null;
  private url: string;
  private getToken: () => Promise<string>;
  private reconnectDelay: number;
  private maxReconnectDelay: number;
  private reconnectAttempts = 0;
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  private heartbeatTimer: ReturnType<typeof setInterval> | null = null;
  private heartbeatInterval: number;
  private autoReconnect: boolean;
  private debug: boolean;
  private manualDisconnect = false;

  private _state: ConnectionState = 'disconnected';
  private listeners = new Map<string, Set<EventHandler>>();

  /**
   * Create a new WebSocket client
   *
   * @param url - WebSocket server URL (e.g., 'ws://localhost:34076/ws')
   * @param getToken - Async function that returns the JWT auth token
   * @param options - Optional configuration
   */
  constructor(
    url: string,
    getToken: () => Promise<string>,
    options: WebSocketClientOptions = {}
  ) {
    this.url = url;
    this.getToken = getToken;
    this.reconnectDelay = options.initialReconnectDelay ?? 1000;
    this.maxReconnectDelay = options.maxReconnectDelay ?? 30000;
    this.heartbeatInterval = options.heartbeatInterval ?? 30000;
    this.autoReconnect = options.autoReconnect ?? true;
    this.debug = options.debug ?? false;
  }

  /**
   * Current connection state
   */
  get state(): ConnectionState {
    return this._state;
  }

  private setState(newState: ConnectionState): void {
    const oldState = this._state;
    if (oldState === newState) return;

    this._state = newState;
    this.log(`State changed: ${oldState} → ${newState}`);

    const event: StateChangeEvent = {
      oldState,
      newState,
      timestamp: new Date().toISOString(),
    };
    this.emit('state:change', event);
  }

  /**
   * Connect to the WebSocket server
   */
  async connect(): Promise<void> {
    if (this.state === 'connected' || this.state === 'connecting') {
      this.log('Already connected or connecting');
      return;
    }

    this.setState('connecting');
    this.manualDisconnect = false;

    try {
      const token = await this.getToken();
      const wsUrl = `${this.url}?token=${encodeURIComponent(token)}`;

      this.ws = new WebSocket(wsUrl);

      this.ws.onopen = () => {
        this.log('WebSocket connected');
        this.setState('connected');
        this.reconnectAttempts = 0;
        this.startHeartbeat();
        this.emit('connected', undefined);
      };

      this.ws.onmessage = (event) => {
        this.handleMessage(event.data);
      };

      this.ws.onerror = (error) => {
        this.log('WebSocket error:', error);
        this.emitError(new Error('WebSocket error'), 'connection');
      };

      this.ws.onclose = (event) => {
        this.log(`WebSocket closed (code: ${event.code}, reason: ${event.reason})`);
        this.stopHeartbeat();
        this.setState('disconnected');
        this.emit('disconnected', undefined);

        // Only reconnect if not manually disconnected and auto-reconnect enabled
        if (!this.manualDisconnect && this.autoReconnect) {
          this.scheduleReconnect();
        }
      };
    } catch (error) {
      this.log('Failed to connect:', error);
      this.setState('disconnected');
      this.emitError(
        error instanceof Error ? error : new Error(String(error)),
        'connect'
      );

      // Attempt reconnection on connection failure
      if (this.autoReconnect) {
        this.scheduleReconnect();
      }
    }
  }

  /**
   * Disconnect from the WebSocket server
   */
  disconnect(): void {
    this.log('Manual disconnect requested');
    this.manualDisconnect = true;

    // Cancel any pending reconnection
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }

    this.stopHeartbeat();

    if (this.ws) {
      this.setState('disconnecting');
      this.ws.close();
      this.ws = null;
    }

    this.setState('disconnected');
  }

  /**
   * Register an event listener
   *
   * @param eventType - Event type to listen for
   * @param handler - Callback function
   */
  on<K extends keyof EventTypeMap>(
    eventType: K,
    handler: (event: EventTypeMap[K]) => void
  ): void {
    if (!this.listeners.has(eventType)) {
      this.listeners.set(eventType, new Set());
    }
    this.listeners.get(eventType)!.add(handler as EventHandler);
  }

  /**
   * Unregister an event listener
   *
   * @param eventType - Event type
   * @param handler - Callback function to remove
   */
  off<K extends keyof EventTypeMap>(
    eventType: K,
    handler: (event: EventTypeMap[K]) => void
  ): void {
    const handlers = this.listeners.get(eventType);
    if (handlers) {
      handlers.delete(handler as EventHandler);
      if (handlers.size === 0) {
        this.listeners.delete(eventType);
      }
    }
  }

  /**
   * Register a one-time event listener
   *
   * @param eventType - Event type to listen for
   * @param handler - Callback function
   */
  once<K extends keyof EventTypeMap>(
    eventType: K,
    handler: (event: EventTypeMap[K]) => void
  ): void {
    const wrappedHandler = (event: EventTypeMap[K]) => {
      handler(event);
      this.off(eventType, wrappedHandler);
    };
    this.on(eventType, wrappedHandler);
  }

  /**
   * Send a message to the server
   *
   * @param event - Event to send
   */
  send(event: BaseEvent): void {
    if (this.state !== 'connected' || !this.ws) {
      this.log('Cannot send message: not connected');
      throw new Error('WebSocket is not connected');
    }

    try {
      this.ws.send(JSON.stringify(event));
    } catch (error) {
      this.log('Failed to send message:', error);
      this.emitError(
        error instanceof Error ? error : new Error(String(error)),
        'send'
      );
      throw error;
    }
  }

  private handleMessage(data: string): void {
    try {
      const event = JSON.parse(data) as BaseEvent;
      this.log(`Received event: ${event.type}`);
      // Emit the event with the dynamic type
      const handlers = this.listeners.get(event.type);
      if (handlers) {
        handlers.forEach((handler) => {
          try {
            // Type assertion needed for dynamic event handling
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            handler(event as any);
          } catch (error) {
            this.log(`Error in event handler for ${event.type}:`, error);
          }
        });
      }
    } catch (error) {
      this.log('Failed to parse message:', error);
      this.emitError(
        new Error('Failed to parse WebSocket message'),
        'parse'
      );
    }
  }

  private emit<K extends keyof EventTypeMap>(
    eventType: K,
    event: EventTypeMap[K]
  ): void {
    const handlers = this.listeners.get(eventType);
    if (handlers) {
      handlers.forEach((handler) => {
        try {
          // Type assertion needed for dynamic event handling
          // eslint-disable-next-line @typescript-eslint/no-explicit-any
          handler(event as any);
        } catch (error) {
          this.log(`Error in event handler for ${eventType}:`, error);
        }
      });
    }
  }

  private emitError(error: Error, context: string): void {
    const errorEvent: ErrorEvent = {
      error,
      timestamp: new Date().toISOString(),
      context,
    };
    this.emit('error', errorEvent);
  }

  private scheduleReconnect(): void {
    if (this.reconnectTimer) {
      return; // Reconnection already scheduled
    }

    this.reconnectAttempts++;
    const delay = Math.min(
      this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1),
      this.maxReconnectDelay
    );

    this.log(`Scheduling reconnection attempt ${this.reconnectAttempts} in ${delay}ms`);

    this.reconnectTimer = setTimeout(() => {
      this.reconnectTimer = null;
      this.connect();
    }, delay);
  }

  private startHeartbeat(): void {
    this.stopHeartbeat();

    this.heartbeatTimer = setInterval(() => {
      if (this.state === 'connected') {
        try {
          this.send({
            type: 'heartbeat',
            payload: {},
            timestamp: new Date().toISOString(),
          });
          this.log('Heartbeat sent');
        } catch (error) {
          this.log('Failed to send heartbeat:', error);
        }
      }
    }, this.heartbeatInterval);
  }

  private stopHeartbeat(): void {
    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer);
      this.heartbeatTimer = null;
    }
  }

  private log(...args: unknown[]): void {
    if (this.debug) {
      // Debug logging is intentional in this client library
      // eslint-disable-next-line no-console
      console.log('[WebSocketClient]', ...args);
    }
  }
}
