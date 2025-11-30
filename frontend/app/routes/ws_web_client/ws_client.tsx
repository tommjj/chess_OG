import type { Route } from './+types/ws_client';

import WsWebClient from '~/features/ws_web_client/client.component';

export function meta({}: Route.MetaArgs) {
    return [
        { title: 'WebSocket Client' },
        { name: 'description', content: 'WebSocket Client Test Page' },
    ];
}

export default function WsClientPage() {
    return <WsWebClient />;
}
