import { useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useNavigate, useParams } from 'react-router-dom';
import Sidebar from '@/shared/components/Sidebar';
import Header from '@/shared/components/Header';
import AlgorithmPostForm from '@/features/posts/components/AlgorithmPostForm';
import PostPreview from '@/features/posts/components/PostPreview';
import PublishModal from '@/features/posts/components/PublishModal';
import { postsApi } from '@/features/posts/api';
import { settingsApi } from '@/features/settings/api';
import type { AlgorithmInput } from '@/features/posts/types';

export default function NewAlgorithmPostPage() {
    const navigate = useNavigate();
    const { id } = useParams();
    const queryClient = useQueryClient();
    const editingPostId = id ? Number(id) : null;

    const [previewMarkdown, setPreviewMarkdown] = useState('');
    const [publishOpen, setPublishOpen] = useState(false);
    const [pendingInput, setPendingInput] = useState<AlgorithmInput | null>(null);
    const [savedPostId, setSavedPostId] = useState<number | null>(null);
    const [successUrl, setSuccessUrl] = useState<string | null>(null);

    const { data: editingPost, isLoading: isEditingLoading } = useQuery({
        queryKey: ['post', editingPostId],
        queryFn: () => postsApi.get(editingPostId as number),
        enabled: Boolean(editingPostId),
    });

    const { data: settings } = useQuery({
        queryKey: ['settings'],
        queryFn: settingsApi.get,
    });

    const previewMutation = useMutation({
        mutationFn: (input: AlgorithmInput) => postsApi.preview({ type: 'ALGORITHM', input }),
        onSuccess: (res) => setPreviewMarkdown(res.contentMarkdown),
    });

    const saveMutation = useMutation({
        mutationFn: (input: AlgorithmInput) =>
            editingPostId
                ? postsApi.update(editingPostId, { type: 'ALGORITHM', input })
                : postsApi.create({ type: 'ALGORITHM', input }),
        onSuccess: (res) => {
            setSavedPostId(res.id);
            queryClient.invalidateQueries({ queryKey: ['posts'] });
            queryClient.invalidateQueries({ queryKey: ['post', res.id] });
            alert(`초안이 저장되었습니다. (ID: ${res.id})`);
        },
    });

    const publishMutation = useMutation({
        mutationFn: ({
            id,
            categoryId,
            tags,
        }: {
            id: number;
            categoryId: string;
            tags: string[];
        }) => postsApi.publishTistory(id, { visibility: 0, categoryId, tags }),
        onSuccess: (res) => {
            queryClient.invalidateQueries({ queryKey: ['posts'] });
            setPublishOpen(false);
            if (res.externalUrl) {
                setSuccessUrl(res.externalUrl);
            }
        },
    });

    const handlePreview = (input: AlgorithmInput) => {
        previewMutation.mutate(input);
    };

    const handleSave = (input: AlgorithmInput) => {
        saveMutation.mutate(input);
    };

    const handlePublishOpen = (input: AlgorithmInput) => {
        setPendingInput(input);
        setPublishOpen(true);
    };

    const handlePublishConfirm = async (categoryId: string, tags: string[]) => {
        let postId = editingPostId ?? savedPostId;
        if (postId && pendingInput) {
            await postsApi.update(postId, { type: 'ALGORITHM', input: pendingInput });
            queryClient.invalidateQueries({ queryKey: ['post', postId] });
            queryClient.invalidateQueries({ queryKey: ['posts'] });
        }
        if (!postId && pendingInput) {
            const created = await postsApi.create({ type: 'ALGORITHM', input: pendingInput });
            postId = created.id;
            setSavedPostId(created.id);
            queryClient.invalidateQueries({ queryKey: ['posts'] });
        }
        if (!postId) throw new Error('글 저장에 실패했습니다.');
        await publishMutation.mutateAsync({ id: postId, categoryId, tags });
    };

    const initialInput = editingPost?.input;

    return (
        <div className="flex h-screen bg-surface">
            <Sidebar />

            <div className="flex flex-col flex-1 ml-60 min-h-0">
                <Header title={editingPostId ? '알고리즘 풀이 수정' : '새 알고리즘 풀이 작성'} />

                {successUrl && (
                    <div className="mx-6 mt-4 px-4 py-3 bg-success-bg text-success-text text-sm rounded-lg flex items-center justify-between">
                        <span>✅ Tistory 발행 완료!</span>
                        <a
                            href={successUrl}
                            target="_blank"
                            rel="noreferrer"
                            className="underline font-medium"
                        >
                            {successUrl}
                        </a>
                        <button
                            onClick={() => navigate('/posts')}
                            className="ml-4 text-success-text underline text-xs"
                        >
                            대시보드로 →
                        </button>
                    </div>
                )}

                <main className="flex-1 min-h-0 flex gap-3 p-4 pt-3">
                    {/* Form panel */}
                    <div className="card flex flex-col flex-shrink-0 w-[564px] min-h-0 overflow-hidden">
                        {isEditingLoading ? (
                            <div className="p-6 text-sm text-text-sec">불러오는 중...</div>
                        ) : (
                            <AlgorithmPostForm
                                onPreview={handlePreview}
                                onSave={handleSave}
                                onPublish={handlePublishOpen}
                                initialValue={initialInput}
                                isPreviewLoading={previewMutation.isPending}
                                isSaveLoading={saveMutation.isPending}
                                saveLabel={editingPostId ? '수정 저장' : '초안 저장'}
                            />
                        )}
                    </div>

                    {/* Preview panel */}
                    <div className="card flex flex-col flex-1 min-h-0 overflow-hidden">
                        <div className="px-4 py-3 border-b border-border flex-shrink-0">
                            <h2 className="text-xs font-semibold text-text-sec">
                                Markdown 미리보기
                            </h2>
                        </div>
                        <div className="flex-1 overflow-y-auto p-4">
                            <PostPreview markdown={previewMarkdown} />
                        </div>
                    </div>
                </main>
            </div>

            <PublishModal
                open={publishOpen}
                onClose={() => setPublishOpen(false)}
                onConfirm={handlePublishConfirm}
                input={pendingInput}
                defaultCategoryId={settings?.tistory.categoryId}
                defaultTags={settings?.tistory.defaultTags}
            />
        </div>
    );
}
