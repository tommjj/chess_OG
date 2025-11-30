import { type RouteConfig, index, route } from '@react-router/dev/routes';

export default [
    index('routes/home.tsx'),
    route('/test/:id', 'routes/test.tsx'),
    route('/ws-client', 'routes/ws_web_client/ws_client.tsx'),
] satisfies RouteConfig;
