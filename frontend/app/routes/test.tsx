import React, { useEffect } from 'react';
import type { Route } from './+types/test';
import { Await } from 'react-router';

export function meta({}: Route.MetaArgs) {
    return [
        { title: 'Test Page' },
        { name: 'description', content: 'This is the test page.' },
    ];
}

export async function clientLoader() {
    const data = await new Promise<string>((resolve) => {
        setTimeout(() => resolve('Client Loader Data Loaded'), 500); // 2s delay
    });

    const w = new Promise<string>((resolve, reject) => {
        setTimeout(() => {
            if (Math.random() < 0.5) {
                resolve('Awaited');
            } else {
                reject('Failed to load awaited value');
            }
        }, 3000); // 5s delay
    });

    return {
        data,
        awaited: w,
    };
}

export default function Test({
    params: { id },
    loaderData,
}: Route.ComponentProps) {
    const [state] = React.useState<string>(windowTest);
    const { data, awaited } = loaderData;

    return (
        <div>
            <h1 className="text-5xl">
                {state} : {data}
            </h1>
            <div>
                <React.Suspense
                    fallback={
                        <h2 className="text-3xl">Loading awaited value...</h2>
                    }
                >
                    <Await
                        resolve={awaited}
                        errorElement={<h2>Error loading awaited value</h2>}
                    >
                        {(value: string) => (
                            <h2 className="text-3xl">Awaited value: {value}</h2>
                        )}
                    </Await>
                </React.Suspense>
            </div>

            <div></div>
        </div>
    );
}

function windowTest() {
    if (typeof window !== 'undefined') {
        return 'Window is defined';
    } else {
        return 'Window is not defined';
    }
}
