'use client'
import {useEffect} from 'react';

export default function UserLogger() {
    useEffect(() => {
        async function run() {
            try {
                const result = await fetch("/api/go/users", { cache: 'no-store' });
                if (!result.ok) {
                    console.error("Failed to Fetch All Users", result.status, await result.text());
                    return;
                }
                const allUsersList = await result.json();
                console.log("This is the list of all Users", allUsersList);
            }
            catch (err) {
                console.error("Error", err)
            }
        }
        run();
    }, [])

    return null;
}