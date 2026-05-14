/** @type {import('tailwindcss').Config} */
export default {
    content: ['./index.html', './src/**/*.{ts,tsx}'],
    theme: {
        extend: {
            colors: {
                sidebar: '#0F1729',
                'sidebar-item': '#253040',
                primary: {
                    DEFAULT: '#6366F1',
                    dark: '#4F46E5',
                },
                surface: '#F8F9FA',
                text: {
                    DEFAULT: '#1E293B',
                    sec: '#64748B',
                    hint: '#94A3B8',
                },
                border: '#E2E8F0',
                success: {
                    DEFAULT: '#10B981',
                    bg: '#E0F9F0',
                    text: '#0A8159',
                },
                error: {
                    DEFAULT: '#EF4444',
                    bg: '#FDEBEB',
                    text: '#C92C2C',
                },
                tag: {
                    DEFAULT: '#EDEEFD',
                    text: '#4F52D5',
                },
                code: {
                    bg: '#1D2130',
                    text: '#BFD9EF',
                },
                nav: '#BDC9DC',
                warn: '#F59E0B',
            },
            boxShadow: {
                card: '0 2px 8px rgba(0,0,0,0.06)',
                modal: '0 20px 40px rgba(0,0,0,0.25)',
            },
        },
    },
    plugins: [],
};
