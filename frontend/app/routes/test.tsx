import type { Route } from './+types/test';

export function meta({}: Route.MetaArgs) {
    return [
        { title: 'Test Page' },
        { name: 'description', content: 'This is the test page.' },
    ];
}

export default function Test({ params: { id } }: Route.ComponentProps) {
    return <h1>Test Page: {id}</h1>;
}
