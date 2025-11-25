/**
 * WebSocket Client Tests
 */

/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable @typescript-eslint/no-unused-vars */
/* eslint-disable max-lines-per-function */

import { WebSocketClient } from './websocket-client';
import type {
  ConnectionState,
  ImageGenerationCompletedEvent,
  StateChangeEvent,
  ErrorEvent,
} from './websocket-types';

// Mock WebSocket
class MockWebSocket {
  static CONNECTING = 0;
  static OPEN = 1;
  static CLOSING = 2;
  static CLOSED = 3;

  readyState = MockWebSocket.CONNECTING;
  onopen: ((event: Event) => void) | null = null;
  onclose: ((event: CloseEvent) => void) | null = null;
  onmessage: ((event: MessageEvent) => void) | null = null;
  onerror: ((event: Event) => void) | null = null;

  constructor(public url: string) {
    // Simulate async connection - immediate for testing
    Promise.resolve().then(() => {
      this.readyState = MockWebSocket.OPEN;
      if (this.onopen) {
        this.onopen(new Event('open'));
      }
    });
  }

  send(_data: string): void {
    if (this.readyState !== MockWebSocket.OPEN) {
      throw new Error('WebSocket is not open');
    }
  }

  close(): void {
    this.readyState = MockWebSocket.CLOSING;
    Promise.resolve().then(() => {
      this.readyState = MockWebSocket.CLOSED;
      if (this.onclose) {
        this.onclose(new CloseEvent('close', { code: 1000, reason: 'Normal closure' }));
      }
    });
  }

  // Test helper to simulate receiving a message
  simulateMessage(data: string): void {
    if (this.onmessage) {
      this.onmessage(new MessageEvent('message', { data }));
    }
  }

  // Test helper to simulate an error
  simulateError(): void {
    if (this.onerror) {
      this.onerror(new Event('error'));
    }
  }
}

// Replace global WebSocket with mock
(global as any).WebSocket = MockWebSocket;

describe('WebSocketClient', () => {
  let client: WebSocketClient;
  let mockGetToken: jest.Mock<Promise<string>>;

  beforeEach(() => {
    mockGetToken = jest.fn().mockResolvedValue('test-token');
    client = new WebSocketClient('ws://localhost:34076/ws', mockGetToken, {
      debug: false,
    });
  });

  afterEach(() => {
    client.disconnect();
  });

  describe('Connection Lifecycle', () => {
    it('should start in disconnected state', () => {
      expect(client.state).toBe('disconnected');
    });

    it('should transition to connecting then connected state', async () => {
      const stateChanges: ConnectionState[] = [];
      client.on('state:change', (event) => {
        stateChanges.push(event.newState);
      });

      await client.connect();
      // Wait for connection to complete
      await new Promise((resolve) => setTimeout(resolve, 10));

      expect(stateChanges).toContain('connecting');
      expect(stateChanges).toContain('connected');
      expect(client.state).toBe('connected');
    });

    it('should call getToken before connecting', async () => {
      await client.connect();
      await new Promise((resolve) => setTimeout(resolve, 10));
      expect(mockGetToken).toHaveBeenCalled();
    });

    it('should emit connected event when connection succeeds', async () => {
      const connectedHandler = jest.fn();
      client.on('connected', connectedHandler);

      await client.connect();
      await new Promise((resolve) => setTimeout(resolve, 10));

      expect(connectedHandler).toHaveBeenCalled();
    });

    it('should disconnect cleanly', async () => {
      await client.connect();
      await new Promise((resolve) => setTimeout(resolve, 10));

      const disconnectedHandler = jest.fn();
      client.on('disconnected', disconnectedHandler);

      client.disconnect();
      await new Promise((resolve) => setTimeout(resolve, 10));

      expect(client.state).toBe('disconnected');
      expect(disconnectedHandler).toHaveBeenCalled();
    });
  });

  describe('Event Handling', () => {
    it('should register and call event listeners', async () => {
      await client.connect();
      await new Promise((resolve) => setTimeout(resolve, 10));

      const handler = jest.fn();
      client.on('image.generation.completed', handler);

      // Get the WebSocket instance and simulate a message
      const ws = (client as any).ws as MockWebSocket;
      const event: ImageGenerationCompletedEvent = {
        type: 'image.generation.completed',
        payload: {
          job_id: 'job-123',
          image_url: 'https://example.com/image.jpg',
        },
        timestamp: new Date().toISOString(),
      };
      ws.simulateMessage(JSON.stringify(event));

      expect(handler).toHaveBeenCalledWith(event);
    });

    it('should support multiple listeners for same event', async () => {
      await client.connect();
      await new Promise((resolve) => setTimeout(resolve, 10));

      const handler1 = jest.fn();
      const handler2 = jest.fn();
      client.on('image.generation.completed', handler1);
      client.on('image.generation.completed', handler2);

      const ws = (client as any).ws as MockWebSocket;
      const event: ImageGenerationCompletedEvent = {
        type: 'image.generation.completed',
        payload: {
          job_id: 'job-123',
          image_url: 'https://example.com/image.jpg',
        },
        timestamp: new Date().toISOString(),
      };
      ws.simulateMessage(JSON.stringify(event));

      expect(handler1).toHaveBeenCalledWith(event);
      expect(handler2).toHaveBeenCalledWith(event);
    });

    it('should remove event listeners with off()', async () => {
      await client.connect();
      await new Promise((resolve) => setTimeout(resolve, 10));

      const handler = jest.fn();
      client.on('image.generation.completed', handler);
      client.off('image.generation.completed', handler);

      const ws = (client as any).ws as MockWebSocket;
      const event: ImageGenerationCompletedEvent = {
        type: 'image.generation.completed',
        payload: {
          job_id: 'job-123',
          image_url: 'https://example.com/image.jpg',
        },
        timestamp: new Date().toISOString(),
      };
      ws.simulateMessage(JSON.stringify(event));

      expect(handler).not.toHaveBeenCalled();
    });

    it('should support once() for one-time listeners', async () => {
      await client.connect();
      await new Promise((resolve) => setTimeout(resolve, 10));

      const handler = jest.fn();
      client.once('image.generation.completed', handler);

      const ws = (client as any).ws as MockWebSocket;
      const event: ImageGenerationCompletedEvent = {
        type: 'image.generation.completed',
        payload: {
          job_id: 'job-123',
          image_url: 'https://example.com/image.jpg',
        },
        timestamp: new Date().toISOString(),
      };

      // Send event twice
      ws.simulateMessage(JSON.stringify(event));
      ws.simulateMessage(JSON.stringify(event));

      // Handler should only be called once
      expect(handler).toHaveBeenCalledTimes(1);
    });

    it('should emit error event on malformed message', async () => {
      await client.connect();
      await new Promise((resolve) => setTimeout(resolve, 10));

      const errorHandler = jest.fn();
      client.on('error', errorHandler);

      const ws = (client as any).ws as MockWebSocket;
      ws.simulateMessage('invalid json');

      expect(errorHandler).toHaveBeenCalled();
      const errorEvent = errorHandler.mock.calls[0][0] as ErrorEvent;
      expect(errorEvent.error.message).toContain('parse');
    });
  });

  describe('Sending Messages', () => {
    it('should send messages when connected', async () => {
      await client.connect();
      await new Promise((resolve) => setTimeout(resolve, 10));

      const ws = (client as any).ws as MockWebSocket;
      const sendSpy = jest.spyOn(ws, 'send');

      client.send({
        type: 'test.event',
        payload: { data: 'test' },
        timestamp: new Date().toISOString(),
      });

      expect(sendSpy).toHaveBeenCalled();
    });

    it('should throw error when sending while disconnected', () => {
      expect(() => {
        client.send({
          type: 'test.event',
          payload: {},
          timestamp: new Date().toISOString(),
        });
      }).toThrow('not connected');
    });
  });

  describe('State Change Events', () => {
    it('should emit state change events', async () => {
      const stateChanges: StateChangeEvent[] = [];
      client.on('state:change', (event) => {
        stateChanges.push(event);
      });

      await client.connect();
      await new Promise((resolve) => setTimeout(resolve, 10));

      client.disconnect();
      await new Promise((resolve) => setTimeout(resolve, 10));

      expect(stateChanges.length).toBeGreaterThan(0);
      expect(stateChanges[0].oldState).toBe('disconnected');
      expect(stateChanges[0].newState).toBe('connecting');
    });
  });

  describe('Configuration Options', () => {
    it('should accept custom configuration', () => {
      const customClient = new WebSocketClient('ws://localhost:34076/ws', mockGetToken, {
        initialReconnectDelay: 500,
        maxReconnectDelay: 10000,
        heartbeatInterval: 15000,
        autoReconnect: false,
        debug: true,
      });

      expect(customClient).toBeDefined();
      customClient.disconnect();
    });
  });
});
