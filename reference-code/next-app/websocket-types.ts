/**
 * WebSocket Client Type Definitions
 *
 * Type definitions for the Masquerade WebSocket client, including
 * connection states, event types, and message payloads.
 */

/**
 * Connection state of the WebSocket client
 */
export type ConnectionState =
  | 'connecting'
  | 'connected'
  | 'disconnecting'
  | 'disconnected';

/**
 * Base event structure for all WebSocket messages
 */
export interface BaseEvent {
  type: string;
  payload: Record<string, unknown>;
  timestamp: string;
}

/**
 * Event handler callback function
 */
export type EventHandler<T extends BaseEvent = BaseEvent> = (event: T) => void;

/**
 * Image generation started event
 */
export interface ImageGenerationStartedEvent extends BaseEvent {
  type: 'image.generation.started';
  payload: {
    job_id: string;
    user_id: string;
    prompt?: string;
  };
}

/**
 * Image generation progress event
 */
export interface ImageGenerationProgressEvent extends BaseEvent {
  type: 'image.generation.progress';
  payload: {
    job_id: string;
    progress: number; // 0-100
    stage?: string;
  };
}

/**
 * Image generation completed event
 */
export interface ImageGenerationCompletedEvent extends BaseEvent {
  type: 'image.generation.completed';
  payload: {
    job_id: string;
    image_url: string;
    metadata?: Record<string, unknown>;
  };
}

/**
 * Image generation failed event
 */
export interface ImageGenerationFailedEvent extends BaseEvent {
  type: 'image.generation.failed';
  payload: {
    job_id: string;
    error: string;
    retry_allowed?: boolean;
  };
}

/**
 * System notification event
 */
export interface SystemNotificationEvent extends BaseEvent {
  type: 'system.notification';
  payload: {
    severity: 'info' | 'warning' | 'error';
    message: string;
    action_url?: string;
  };
}

/**
 * Heartbeat event
 */
export interface HeartbeatEvent extends BaseEvent {
  type: 'heartbeat';
  payload: Record<string, never>;
}

/**
 * Connection state change event
 */
export interface StateChangeEvent {
  oldState: ConnectionState;
  newState: ConnectionState;
  timestamp: string;
}

/**
 * Error event
 */
export interface ErrorEvent {
  error: Error;
  timestamp: string;
  context?: string;
}

/**
 * Union type of all known event types
 */
export type KnownEvent =
  | ImageGenerationStartedEvent
  | ImageGenerationProgressEvent
  | ImageGenerationCompletedEvent
  | ImageGenerationFailedEvent
  | SystemNotificationEvent
  | HeartbeatEvent;

/**
 * Event type to payload mapping for type-safe event handlers
 */
export interface EventTypeMap {
  'image.generation.started': ImageGenerationStartedEvent;
  'image.generation.progress': ImageGenerationProgressEvent;
  'image.generation.completed': ImageGenerationCompletedEvent;
  'image.generation.failed': ImageGenerationFailedEvent;
  'system.notification': SystemNotificationEvent;
  'heartbeat': HeartbeatEvent;
  'state:change': StateChangeEvent;
  'error': ErrorEvent;
  'connected': void;
  'disconnected': void;
}
