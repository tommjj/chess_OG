import WSClient, { WSDefaultsEvents } from '~/common/ws/client';

const ws = new WSClient('ws://localhost:8080/ws');

ws.on(WSDefaultsEvents.Connection, () => {
    console.log('WebSocket connection opened');
});

ws.on(WSDefaultsEvents.Disconnection, () => {
    console.log('WebSocket connection closed');
});

ws.on(WSDefaultsEvents.Error, (error) => {
    console.error('WebSocket error:', error);
});

export default ws;
