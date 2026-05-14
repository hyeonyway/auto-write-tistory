import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import Sidebar from '@/shared/components/Sidebar';
import Header from '@/shared/components/Header';
import TistorySettingsForm from '@/features/settings/components/TistorySettingsForm';
import { settingsApi } from '@/features/settings/api';
import type { UpdateSettingsRequest } from '@/features/settings/types';

export default function SettingsPage() {
    const queryClient = useQueryClient();

    const { data, isLoading, isError } = useQuery({
        queryKey: ['settings'],
        queryFn: settingsApi.get,
    });
    const { data: sessionStatus } = useQuery({
        queryKey: ['tistory-session-status'],
        queryFn: settingsApi.sessionStatus,
    });

    const updateMutation = useMutation({
        mutationFn: (req: UpdateSettingsRequest) => settingsApi.update(req),
        onSuccess: (res) => {
            queryClient.setQueryData(['settings'], res);
            alert('설정이 저장되었습니다.');
        },
    });

    const startSessionMutation = useMutation({
        mutationFn: settingsApi.startSession,
        onSuccess: (res) => {
            alert(res.message);
            queryClient.invalidateQueries({ queryKey: ['tistory-session-status'] });
        },
    });

    const confirmSessionMutation = useMutation({
        mutationFn: settingsApi.confirmSession,
        onSuccess: (res) => {
            alert(res.message);
            queryClient.invalidateQueries({ queryKey: ['tistory-session-status'] });
        },
    });

    const deleteSessionMutation = useMutation({
        mutationFn: settingsApi.deleteSession,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['tistory-session-status'] });
            alert('Tistory 세션을 삭제했습니다.');
        },
    });

    const currentSession = sessionStatus ?? data?.tistory.session;

    return (
        <div className="flex h-screen bg-surface">
            <Sidebar />

            <div className="flex flex-col flex-1 ml-60 min-h-0">
                <Header title="설정" />

                <main className="flex-1 overflow-y-auto p-6">
                    <div className="card p-6 max-w-2xl">
                        <h2 className="text-base font-semibold text-text">Tistory 발행 설정</h2>
                        <p className="text-sm text-text-sec mt-1 mb-6">
                            Playwright 로그인 세션과 Tistory 기본 발행 설정입니다.
                        </p>

                        {data && currentSession && (
                            <div className="mb-6 rounded-xl border border-border p-4">
                                <div className="flex items-start justify-between gap-4">
                                    <div>
                                        <p className="text-sm font-semibold text-text">
                                            Tistory 로그인 세션
                                        </p>
                                        <p className="text-xs text-text-sec mt-1">
                                            {currentSession.message}
                                        </p>
                                    </div>
                                    <span
                                        className={[
                                            'text-[11px] font-semibold px-2 py-1 rounded',
                                            currentSession.connected
                                                ? 'bg-success-bg text-success-text'
                                                : 'bg-tag text-tag-text',
                                        ].join(' ')}
                                    >
                                        {currentSession.connected ? '연결됨' : '연결 필요'}
                                    </span>
                                </div>
                                <div className="flex flex-wrap gap-2 mt-4">
                                    <button
                                        type="button"
                                        onClick={() => startSessionMutation.mutate()}
                                        className="btn-secondary text-xs"
                                        disabled={startSessionMutation.isPending}
                                    >
                                        로그인 브라우저 열기
                                    </button>
                                    <button
                                        type="button"
                                        onClick={() => confirmSessionMutation.mutate()}
                                        className="btn-secondary text-xs"
                                        disabled={confirmSessionMutation.isPending}
                                    >
                                        인증 완료 확인
                                    </button>
                                    <button
                                        type="button"
                                        onClick={() => deleteSessionMutation.mutate()}
                                        className="btn-secondary text-xs"
                                        disabled={deleteSessionMutation.isPending}
                                    >
                                        세션 연결 해제
                                    </button>
                                </div>
                            </div>
                        )}

                        <div className="border-t border-border pt-6">
                            {isLoading && (
                                <div className="py-12 text-center text-text-hint text-sm">
                                    불러오는 중...
                                </div>
                            )}
                            {isError && (
                                <div className="py-12 text-center text-error text-sm">
                                    설정을 불러오지 못했습니다.
                                </div>
                            )}
                            {data && (
                                <TistorySettingsForm
                                    settings={data.tistory}
                                    onSave={(values) =>
                                        updateMutation.mutateAsync({ tistory: values })
                                    }
                                    isSaving={updateMutation.isPending}
                                />
                            )}
                        </div>
                    </div>
                </main>
            </div>
        </div>
    );
}
