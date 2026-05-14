import { useEffect, useState } from 'react';
import Modal from '@/shared/components/Modal';
import { settingsApi } from '@/features/settings/api';
import type { TistoryCategoryOption } from '../types';
import type { TistoryLoginInput } from '@/features/posts/types';

interface TistoryCategoryModalProps {
    open: boolean;
    blogUrl: string;
    selectedCategoryId: string;
    onClose: () => void;
    onSelect: (category: TistoryCategoryOption) => void;
}

export default function TistoryCategoryModal({
    open,
    blogUrl,
    selectedCategoryId,
    onClose,
    onSelect,
}: TistoryCategoryModalProps) {
    const [loginId, setLoginId] = useState('');
    const [loginPassword, setLoginPassword] = useState('');
    const [items, setItems] = useState<TistoryCategoryOption[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (!open) return;
        setLoginId('');
        setLoginPassword('');
        setItems([]);
        setError(null);
        setLoading(false);
    }, [open]);

    const handleFetch = async () => {
        const trimmedLoginId = loginId.trim();
        if (!blogUrl.trim()) {
            setError('먼저 Blog URL을 입력해주세요.');
            return;
        }
        if (!trimmedLoginId || !loginPassword) {
            setError('카카오 로그인 ID와 비밀번호를 입력해주세요.');
            return;
        }

        setLoading(true);
        setError(null);
        setItems([]);
        try {
            const req = {
                blogUrl,
                login: {
                    id: trimmedLoginId,
                    password: loginPassword,
                } satisfies TistoryLoginInput,
            };
            const res = await settingsApi.fetchCategories(req);
            setItems(res.items);
            if (res.items.length === 0) {
                setError('카테고리를 찾지 못했습니다.');
            }
        } catch (e) {
            setError(e instanceof Error ? e.message : '카테고리를 불러오지 못했습니다.');
        } finally {
            setLoading(false);
        }
    };

    const handleSelect = (item: TistoryCategoryOption) => {
        onSelect(item);
        onClose();
    };

    const handleClose = () => {
        if (loading) return;
        onClose();
    };

    return (
        <Modal open={open} onClose={handleClose}>
            <div className="p-6">
                <h2 className="text-[17px] font-semibold text-text">카테고리 불러오기</h2>
                <p className="text-sm text-text-sec mt-1">
                    Tistory 글쓰기 화면의 카테고리 목록을 Selenium으로 읽어옵니다.
                </p>
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
                <button
                    type="button"
                    onClick={handleFetch}
                    disabled={loading}
                    className="btn-primary"
                >
                    {loading ? '불러오는 중...' : '카테고리 불러오기'}
                </button>
                <p className="text-[11px] text-text-hint">
                    입력한 로그인 정보는 요청 처리 후 저장하지 않습니다.
                </p>
            </div>

            {items.length > 0 && (
                <div className="border-t border-border px-6 py-4">
                    <p className="text-[11px] font-medium text-text-sec mb-2">카테고리 선택</p>
                    <div className="max-h-[320px] overflow-y-auto space-y-1 pr-1">
                        {items.map((item) => {
                            const isSelected = item.categoryId === selectedCategoryId;
                            return (
                                <button
                                    key={`${item.categoryId}:${item.label}`}
                                    type="button"
                                    onClick={() => handleSelect(item)}
                                    className={[
                                        'w-full text-left rounded-xl border px-3 py-2 transition',
                                        isSelected
                                            ? 'border-primary bg-primary-soft'
                                            : 'border-border bg-surface hover:border-primary',
                                    ].join(' ')}
                                >
                                    <div className="flex items-center justify-between gap-3">
                                        <span className="text-sm font-medium text-text">
                                            {item.label}
                                        </span>
                                        <span className="text-[11px] text-text-hint">
                                            {item.categoryId}
                                        </span>
                                    </div>
                                </button>
                            );
                        })}
                    </div>
                </div>
            )}

            {error && (
                <div className="mx-6 mb-4 px-4 py-3 bg-error-bg rounded-lg text-error-text text-sm">
                    {error}
                </div>
            )}

            <div className="border-t border-border px-6 py-4 flex justify-end">
                <button
                    type="button"
                    onClick={handleClose}
                    className="btn-secondary"
                    disabled={loading}
                >
                    닫기
                </button>
            </div>
        </Modal>
    );
}
