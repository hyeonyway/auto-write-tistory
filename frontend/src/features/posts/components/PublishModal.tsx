import { useState } from 'react';
import Modal from '@/shared/components/Modal';
import type { AlgorithmInput } from '../types';

interface PublishModalProps {
    open: boolean;
    onClose: () => void;
    onConfirm: (categoryId: string, tags: string[]) => Promise<void>;
    input: AlgorithmInput | null;
    defaultCategoryId?: string;
    defaultTags?: string[];
}

export default function PublishModal({
    open,
    onClose,
    onConfirm,
    input,
    defaultCategoryId = '',
    defaultTags = [],
}: PublishModalProps) {
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const tags = Array.from(new Set([...(defaultTags ?? []), ...(input?.tags ?? [])]));
    const title = input ? `[${input.platform}] ${input.problemTitle} - ${input.language}` : '';

    const handleClose = () => {
        if (loading) return;
        setError(null);
        onClose();
    };

    const handleConfirm = async () => {
        setError(null);
        setLoading(true);
        try {
            await onConfirm(defaultCategoryId, tags);
        } catch (e) {
            setError(e instanceof Error ? e.message : '발행 중 오류가 발생했습니다.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <Modal open={open} onClose={handleClose}>
            <div className="p-6">
                <h2 className="text-[17px] font-semibold text-text">Tistory 발행 확인</h2>
                <p className="text-sm text-text-sec mt-1">
                    아래 설정으로 Tistory 비공개 글을 발행합니다.
                </p>
            </div>

            <div className="border-t border-border px-6 py-4 space-y-4">
                <div>
                    <p className="text-[11px] font-medium text-text-sec">글 제목</p>
                    <p className="text-sm font-medium text-text mt-1">{title}</p>
                </div>
                <div>
                    <p className="text-[11px] font-medium text-text-sec">공개 설정</p>
                    <p className="text-sm font-medium text-text mt-1">비공개 (visibility=0)</p>
                </div>
                {tags.length > 0 && (
                    <div>
                        <p className="text-[11px] font-medium text-text-sec">태그</p>
                        <div className="flex flex-wrap gap-1.5 mt-1.5">
                            {tags.map((t) => (
                                <span
                                    key={t}
                                    className="px-2 py-0.5 bg-tag text-tag-text text-[11px] font-medium rounded-full"
                                >
                                    {t}
                                </span>
                            ))}
                        </div>
                    </div>
                )}
            </div>

            <div className="mx-6 mb-4 px-4 py-3 bg-[#141520] rounded-lg">
                <p className="text-[12px] text-warn font-medium">Tistory 세션 확인</p>
                <p className="text-[11px] text-nav mt-0.5">
                    발행 전 Settings에서 Tistory 로그인 세션을 연결해야 합니다.
                </p>
            </div>

            {error && (
                <div className="mx-6 mb-4 px-4 py-3 bg-error-bg rounded-lg text-error-text text-sm">
                    {error}
                </div>
            )}

            <div className="border-t border-border px-6 py-4 flex items-center justify-between">
                <button
                    type="button"
                    onClick={handleClose}
                    className="btn-secondary"
                    disabled={loading}
                >
                    취소
                </button>
                <button
                    type="button"
                    onClick={handleConfirm}
                    disabled={loading}
                    className="btn-primary"
                >
                    {loading ? '발행 중...' : '발행하기'}
                </button>
            </div>
        </Modal>
    );
}
