import { NavLink } from 'react-router-dom';

const navItems = [
    { icon: '📋', label: '대시보드', to: '/posts' },
    { icon: '✏️', label: '새 글 작성', to: '/posts/new/algorithm' },
    { icon: '⚙️', label: '설정', to: '/settings' },
];

export default function Sidebar() {
    return (
        <aside className="fixed top-0 left-0 h-screen w-60 bg-sidebar flex flex-col z-20">
            {/* Logo */}
            <div className="px-5 pt-6 pb-4">
                <div className="text-white font-bold text-[17px]">DevLog Studio</div>
                <div className="text-text-hint text-[11px] mt-0.5">개발 블로그 자동화</div>
            </div>

            <div className="mx-4 h-px bg-sidebar-item" />

            {/* Nav */}
            <nav className="flex-1 px-2 py-3 space-y-1">
                {navItems.map(({ icon, label, to }) => (
                    <NavLink
                        key={to}
                        to={to}
                        end={to === '/posts'}
                        className={({ isActive }) =>
                            [
                                'flex items-center gap-3 px-3 py-3 rounded-lg text-[13px] transition-colors',
                                isActive
                                    ? 'bg-sidebar-item text-white font-semibold'
                                    : 'text-nav hover:bg-sidebar-item/60',
                            ].join(' ')
                        }
                    >
                        <span>{icon}</span>
                        <span>{label}</span>
                    </NavLink>
                ))}
            </nav>

            {/* Footer */}
            <div className="px-5 py-4 text-[11px] text-text-hint">v0.1.0 · Phase 1</div>
        </aside>
    );
}
