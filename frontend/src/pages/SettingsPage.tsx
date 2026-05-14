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

    const updateMutation = useMutation({
        mutationFn: (req: UpdateSettingsRequest) => settingsApi.update(req),
        onSuccess: (res) => {
            queryClient.setQueryData(['settings'], res);
            alert('설정이 저장되었습니다.');
        },
    });

    return (
        <div className="flex h-screen bg-surface">
            <Sidebar />

            <div className="flex flex-col flex-1 ml-60 min-h-0">
                <Header title="설정" />

                <main className="flex-1 overflow-y-auto p-6">
                    <div className="card p-6 max-w-2xl">
                        <h2 className="text-base font-semibold text-text">Tistory 발행 설정</h2>
                        <p className="text-sm text-text-sec mt-1 mb-6">
                            Selenium 기반 자동 발행을 위한 기본 설정입니다.
                        </p>

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
