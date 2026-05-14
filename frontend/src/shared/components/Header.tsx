import type { ReactNode } from 'react';

interface HeaderProps {
    title: string;
    action?: ReactNode;
}

export default function Header({ title, action }: HeaderProps) {
    return (
        <header className="h-16 bg-white border-b border-border flex items-center justify-between px-8 flex-shrink-0">
            <h1 className="text-xl font-semibold text-text">{title}</h1>
            {action && <div>{action}</div>}
        </header>
    );
}
