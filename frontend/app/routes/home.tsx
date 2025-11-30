import { Link } from 'react-router';
import type { Route } from './+types/home';

export function meta({}: Route.MetaArgs) {
    return [
        { title: 'New React Router App' },
        { name: 'description', content: 'Welcome to React Router!' },
    ];
}

export default function Home() {
    return (
        <div>
            <h1>Hello</h1>
            <p>Welcome to your new React Router app.</p>

            <hr />
            <nav className="p-2">
                <ul>
                    <li>
                        <Link
                            className="px-2 py-1.5 rounded bg-black text-white "
                            to="/ws-client"
                        >
                            WebSocket Client Page
                        </Link>
                    </li>
                </ul>
            </nav>
        </div>
    );
}
