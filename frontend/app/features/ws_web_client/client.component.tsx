import React from 'react';
import ws from './ws_client';

// WebSocket Client Component manages WebSocket connection lifecycle
function WsWebClient() {
    const [isConnected, setIsConnected] = React.useState<boolean>(false);
    const [logs, setBlogs] = React.useState<string[]>([]);
    const consoleRef = React.useRef<HTMLPreElement>(null);

    const [eventName, setEventName] = React.useState<string>('');
    const [eventPayload, setEventPayload] = React.useState<string>('');

    const addLog = React.useCallback(
        (log: string, isError: boolean = false) => {
            setBlogs((prevLogs) => [
                ...prevLogs,
                `${isError ? '❌ ERROR:' : ''}[${new Date().toLocaleTimeString()}] ${log}`,
            ]);
        },
        []
    );

    const clearLogs = React.useCallback(() => {
        setBlogs([]);
    }, []);

    const onEventNameChange = React.useCallback(
        (e: React.ChangeEvent<HTMLInputElement>) => {
            setEventName(e.target.value);
        },
        []
    );

    const onEventPayloadChange = React.useCallback(
        (e: React.ChangeEvent<HTMLTextAreaElement>) => {
            setEventPayload(e.target.value);
        },
        []
    );

    const handleSubmit = React.useCallback(
        (e: React.FormEvent) => {
            e.preventDefault();
            let payload: any = null;

            if (!eventName) {
                addLog('Event name is required', true);
                return;
            }

            try {
                payload = eventPayload ? JSON.parse(eventPayload) : null;
            } catch (error) {
                addLog('Invalid JSON payload', true);
                return;
            }

            ws.emit(eventName, payload);
            addLog(
                `Sent event - Name: ${eventName}, Payload: ${JSON.stringify(
                    payload
                )}`
            );
        },
        [eventName, eventPayload, addLog]
    );

    React.useEffect(() => {
        if (!ws.isConnected()) {
            setIsConnected(true);
            ws.connect();
        } else {
            setIsConnected(false);
        }

        const handleConnection = () => {
            console.log('WebSocket connected (from component)');
            addLog('WebSocket connected');
            setIsConnected(true);
        };

        const handleDisconnection = () => {
            console.log('WebSocket disconnected (from component)');
            addLog('WebSocket disconnected');
            setIsConnected(false);
        };

        const handleMessage = (data: { event: string; payload: any }) => {
            console.log('WebSocket message received:', data);
            addLog(
                `Message received - Event: ${data.event}, Payload: ${JSON.stringify(data.payload)}`
            );
        };

        ws.on('$connection', handleConnection);
        ws.on('$disconnection', handleDisconnection);
        ws.on('$message', handleMessage);

        return () => {
            ws.off('$connection', handleConnection);
            ws.off('$disconnection', handleDisconnection);
            ws.off('$message', handleMessage);
            setIsConnected(false);

            ws.disconnect();
        };
    }, []);

    React.useEffect(() => {
        if (consoleRef.current) {
            consoleRef.current.scrollTop = consoleRef.current.scrollHeight;
            console.log(consoleRef.current.scrollHeight);
        }
    }, [logs]);

    const handleToggleConnection = React.useCallback(() => {
        if (isConnected) {
            ws.disconnect();
        } else {
            ws.connect();
        }
    }, [isConnected]);

    return (
        <div className="w-full max-h-screen h-screen grid grid-cols-2">
            <div className="flex flex-col">
                <div>
                    <h1 className="mt-2 ml-2 text-2xl">
                        WebSocket Server Test
                    </h1>
                    <p className="ml-2">
                        Status:{' '}
                        <span
                            className={
                                isConnected ? 'text-green-500' : 'text-red-500'
                            }
                        >
                            {isConnected ? 'Connected' : 'Disconnected'}
                        </span>
                    </p>
                    <div className="mx-2 my-1 flex gap-1">
                        <button
                            onClick={handleToggleConnection}
                            className={
                                'px-2 py-1.5 text-white bg-gray-900 rounded cursor-pointer hover:bg-gray-800 disabled:opacity-50 disabled:cursor-not-allowed ' +
                                (isConnected
                                    ? 'bg-red-600 hover:bg-red-500'
                                    : 'bg-green-600 hover:bg-green-500')
                            }
                        >
                            {isConnected ? 'Disconnect' : 'Connect'}
                        </button>
                    </div>
                </div>
                <div className="flex flex-col grow">
                    <div className="flex w-full p-2">
                        <input
                            onChange={onEventNameChange}
                            value={eventName}
                            type="text"
                            className="w-full border rounded px-2 py-1"
                            placeholder="Type a event name ..."
                        ></input>
                        <button
                            onClick={handleSubmit}
                            disabled={!isConnected || !eventName}
                            className="text-nowrap ml-1 px-2 py-1.5 text-white bg-blue-600 rounded cursor-pointer hover:bg-blue-500"
                        >
                            Send Event
                        </button>
                    </div>
                    <div className="flex flex-col grow p-2 gap-2">
                        <textarea
                            onChange={onEventPayloadChange}
                            value={eventPayload}
                            className="grow w-full p-2 border rounded resize-none"
                            placeholder="Enter event payload..."
                        ></textarea>
                    </div>
                </div>
            </div>
            <div className="flex flex-col border-l p-2 overflow-hidden">
                <h2>Console Output</h2>
                <pre
                    className="size-full max-h-full overflow-y-scroll"
                    ref={consoleRef}
                >
                    {logs.map((log, index) => (
                        <div key={index} className="mb-1">
                            {log}
                        </div>
                    ))}
                </pre>
            </div>
        </div>
    );
}
export default WsWebClient;
