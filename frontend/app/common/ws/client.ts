export type WSEventHandler<T> = (data: T) => void;

export type WSEvents<N extends string, P extends any> = Record<N, P>;

export interface WSEventDefaults {
    $connection: null;
    $disconnection: null;
    $error: { message: string; details?: any };
    $message: {
        event: string;
        payload: any;
    }; // raw message
}

export enum WSDefaultsEvents {
    Connection = '$connection',
    Disconnection = '$disconnection',
    Error = '$error',
    Message = '$message', // all raw messages
}

/**
 * WSClient - A WebSocket client with typed events
 *
 * events format:
 * {
 *   event: string;
 *   payload: any;
 * }
 *
 * @template ClientEvent - Events emitted by the client
 * @template ServerEvent - Events emitted by the server
 */
class WSClient<
    ClientEvent extends WSEvents<string, any>,
    ServerEvent extends WSEvents<string, any>,
> {
    private socket: WebSocket | null = null;
    private url: string;

    private events = new Map<string, Map<WSEventHandler<any>, null>>();

    /**
     * constructor - Create a new WSClient instance
     *
     * @param url - WebSocket server URL e.g. ws://example.com/socket
     * @param headers - Optional headers to include as query parameters in the WebSocket URL e.g. { token: 'abc123' }
     */
    constructor(
        url: string,
        headers?: Record<string, string | number | boolean>
    ) {
        const merged = {
            ...(headers ?? {}),
        };

        const query = new URLSearchParams(
            Object.entries(merged).map(([k, v]) => [k, String(v)])
        ).toString();

        this.url = `${url}?${query}`;
        this.events = new Map();
    }

    /**
     * on - Subscribe to a WebSocket event
     *
     * @param event - Event name
     * @param handler - Event handler function
     * @returns - Unsubscribe function
     */
    on<K extends keyof (ServerEvent & WSEventDefaults)>(
        event: K,
        handler: WSEventHandler<(ServerEvent & WSEventDefaults)[K]>
    ) {
        if (!this.events.has(event as string)) {
            this.events.set(event as string, new Map());
        }
        this.events.get(event as string)!.set(handler, null);

        return () => this.off(event, handler); // Return unsubscribe function
    }

    /**
     * off - Unsubscribe from a WebSocket event
     *
     * @param event - Event name
     * @param handler - Event handler function
     */
    off<K extends keyof (ServerEvent & WSEventDefaults)>(
        event: K,
        handler: WSEventHandler<(ServerEvent & WSEventDefaults)[K]>
    ) {
        this.events.get(event as string)?.delete(handler);
    }

    /**
     * emit - Emit a WebSocket event to the server
     *
     * @param event - Event name
     * @param data - Event payload
     */
    emit<K extends keyof ClientEvent>(event: K, data: ClientEvent[K]) {
        const message = JSON.stringify({ event, payload: data });
        this.sendMessage(message);
    }

    /**
     * connect - Establish the WebSocket connection
     */
    connect() {
        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            console.warn('WebSocket is already connected.');
            return;
        }

        this.socket = new WebSocket(this.url);

        this.socket.onopen = () => {
            this.events
                .get('$connection')
                ?.forEach((_, handler) => handler(null));
        };

        this.socket.onmessage = (event) => {
            try {
                const msg = JSON.parse(event.data);

                if (!msg.event) throw new Error('Missing event field');
                const ev = msg.event;

                this.events
                    .get('$message')
                    ?.forEach((_, handler) => handler(msg)); // Emit raw message event

                this.events
                    .get(ev)
                    ?.forEach((_, handler) => handler(msg.payload));
            } catch (err) {
                console.warn('Invalid WS message:', event.data);
                this.disconnect(); // Disconnect on invalid message (protocol violation)
            }
        };

        this.socket.onclose = () => {
            this.events
                .get('$disconnection')
                ?.forEach((_, handler) => handler(null));
        };

        this.socket.onerror = (error) => {
            this.events.get('$error')?.forEach((_, handler) =>
                handler({
                    message: 'WebSocket error occurred',
                    details: error,
                })
            );
        };
    }

    /**
     * sendMessage - Send a raw message through the WebSocket
     * @param message
     */
    private sendMessage(message: string) {
        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            this.socket.send(message);
        } else {
            console.error('WebSocket is not open. Unable to send message.');
        }
    }

    /**
     * disconnect - Close the WebSocket connection
     */
    disconnect() {
        if (this.socket) {
            this.socket.close();
            this.socket = null;
        }
    }

    isConnected() {
        return this.socket?.readyState === WebSocket.OPEN;
    }

    isDisconnected() {
        return (
            this.socket?.readyState === WebSocket.CLOSED ||
            this.socket?.readyState === WebSocket.CLOSING ||
            this.socket === null
        );
    }
}

export default WSClient;
