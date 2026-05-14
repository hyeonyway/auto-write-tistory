import type { PostStatus } from '@/features/posts/types';

const config: Record<PostStatus, { label: string; className: string }> = {
    PUBLISHED: {
        label: '발행됨',
        className: 'bg-success-bg text-success-text',
    },
    DRAFT: {
        label: '초안',
        className: 'bg-tag text-tag-text',
    },
    FAILED: {
        label: '실패',
        className: 'bg-error-bg text-error-text',
    },
};

export default function StatusBadge({ status }: { status: PostStatus }) {
    const { label, className } = config[status];
    return (
        <span
            className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-[11px] font-medium ${className}`}
        >
            {label}
        </span>
    );
}
