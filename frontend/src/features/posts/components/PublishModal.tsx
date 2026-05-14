import { useState } from 'react';
import Modal from '@/shared/components/Modal';
import type { AlgorithmInput, TistoryLoginInput } from '../types';

interface PublishModalProps {
    open: boolean;
    onClose: () => void;
    onConfirm: (categoryId: string, tags: string[], login: TistoryLoginInput) => Promise<void>;
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
    const [loginId, setLoginId] = useState('');
    const [loginPassword, setLoginPassword] = useState('');

    const tags = Array.from(new Set([...(defaultTags ?? []), ...(input?.tags ?? [])]));
    const title = input ? `[${input.platform}] ${input.problemTitle} - ${input.language}` : '';

    const handleClose = () => {
        if (loading) return;
        setError(null);
        setLoginId('');
        setLoginPassword('');
        onClose();
    };

    const handleConfirm = async () => {
        setError(null);
        const trimmedLoginId = loginId.trim();
        if (!trimmedLoginId || !loginPassword) {
            setError('카카오 로그인 ID와 비밀번호를 입력해주세요.');
            return;
        }
        setLoading(true);
        try {
            await onConfirm(defaultCategoryId, tags, {
                id: trimmedLoginId,
                password: loginPassword,
            });
            setLoginId('');
            setLoginPassword('');
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

            <div className="border-t border-border px-6 py-4 space-y-3">
                <div>
                    <label className="label">카카오 로그인 ID</label>
                    <input
                        value={loginId}
                        onChange={(e) => setLoginId(e.target.value)}
                        className="input-field"
                        autoComplete="off"
                        placeholder="kakao-account@example.com"
                        disabled={loading}
                    />
                </div>
                <div>
                    <label className="label">카카오 로그인 비밀번호</label>
                    <input
                        value={loginPassword}
                        onChange={(e) => setLoginPassword(e.target.value)}
                        type="password"
                        className="input-field"
                        autoComplete="new-password"
                        disabled={loading}
                    />
                </div>
                <p className="text-[11px] text-text-hint">
                    입력한 로그인 정보는 이번 발행 요청에만 사용하고 저장하지 않습니다.
                </p>
            </div>

            {/* Warning */}
            <div className="mx-6 mb-4 px-4 py-3 bg-[#141520] rounded-lg">
                <p className="text-[12px] text-warn font-medium">
                    ⚠ Selenium이 Tistory 화면을 자동 제어합니다.
                </p>
                <p className="text-[11px] text-nav mt-0.5">
                    발행이 시작되면 브라우저를 조작하지 마세요.
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
                    {loading ? 'Selenium 실행 중...' : '발행하기'}
                </button>
            </div>
        </Modal>
    );
}
