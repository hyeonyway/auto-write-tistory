import { useState } from 'react';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import TagInput from '@/shared/components/TagInput';
import type { TistorySettings, UpdateSettingsRequest } from '../types';
import TistoryCategoryModal from './TistoryCategoryModal';

const schema = z.object({
    blogUrl: z.string().url('올바른 URL을 입력해주세요'),
    categoryId: z.string().default(''),
    defaultVisibility: z.number().default(0),
    defaultTags: z.array(z.string()).default([]),
});

type FormValues = z.infer<typeof schema>;

interface TistorySettingsFormProps {
    settings: TistorySettings;
    onSave: (values: UpdateSettingsRequest['tistory']) => Promise<unknown>;
    isSaving?: boolean;
}

export default function TistorySettingsForm({
    settings,
    onSave,
    isSaving,
}: TistorySettingsFormProps) {
    const [categoryModalOpen, setCategoryModalOpen] = useState(false);
    const [selectedCategoryLabel, setSelectedCategoryLabel] = useState('');
    const {
        register,
        control,
        watch,
        setValue,
        handleSubmit,
        formState: { errors },
    } = useForm<FormValues>({
        resolver: zodResolver(schema),
        defaultValues: {
            blogUrl: settings.blogUrl,
            categoryId: settings.categoryId,
            defaultVisibility: settings.defaultVisibility,
            defaultTags: settings.defaultTags,
        },
    });
    const blogUrl = watch('blogUrl');
    const categoryId = watch('categoryId');

    const handleSave = (values: FormValues) => {
        const payload: UpdateSettingsRequest['tistory'] = {
            blogUrl: values.blogUrl,
            categoryId: values.categoryId,
            defaultVisibility: values.defaultVisibility,
            defaultTags: values.defaultTags,
        };
        return onSave(payload);
    };

    return (
        <>
            <form onSubmit={handleSubmit(handleSave)} className="space-y-5">
                {/* Blog URL */}
                <div>
                    <label className="label">Blog URL</label>
                    <input
                        {...register('blogUrl')}
                        type="url"
                        className="input-field"
                        placeholder="https://yourblog.tistory.com"
                    />
                    {errors.blogUrl && (
                        <p className="text-error text-xs mt-1">{errors.blogUrl.message}</p>
                    )}
                </div>

                {/* Category ID */}
                <div>
                    <div className="flex items-center justify-between gap-3">
                        <label className="label">카테고리</label>
                        <button
                            type="button"
                            onClick={() => setCategoryModalOpen(true)}
                            className="text-xs font-medium text-primary hover:underline"
                        >
                            카테고리 불러오기
                        </button>
                    </div>
                    <input
                        {...register('categoryId')}
                        className="input-field"
                        placeholder="카테고리를 불러와 선택하세요"
                        readOnly
                    />
                    {selectedCategoryLabel ? (
                        <p className="text-text-hint text-[11px] mt-1">
                            선택됨: {selectedCategoryLabel} ({categoryId || '미선택'})
                        </p>
                    ) : (
                        <p className="text-text-hint text-[11px] mt-1">
                            Tistory 글쓰기 화면의 카테고리 목록을 불러와 선택합니다.
                        </p>
                    )}
                </div>

                <div className="border border-border rounded-xl p-4">
                    <p className="text-sm font-semibold text-text">카카오 로그인 정보</p>
                    <p className="text-[11px] text-text-hint mt-1">
                        로그인 ID/PW는 저장하지 않습니다. Tistory 발행 버튼을 누른 뒤 확인 모달에서
                        1회 입력합니다.
                    </p>
                </div>

                {/* Default visibility */}
                <div>
                    <label className="label">기본 공개 설정</label>
                    <div className="flex gap-6">
                        {[
                            { label: '비공개 (기본값)', value: 0 },
                            { label: '공개', value: 3 },
                        ].map(({ label, value }) => (
                            <label key={value} className="flex items-center gap-2 cursor-pointer">
                                <input
                                    {...register('defaultVisibility', { valueAsNumber: true })}
                                    type="radio"
                                    value={value}
                                    className="accent-primary"
                                />
                                <span className="text-sm text-text">{label}</span>
                            </label>
                        ))}
                    </div>
                </div>

                {/* Default tags */}
                <div>
                    <label className="label">기본 태그</label>
                    <Controller
                        name="defaultTags"
                        control={control}
                        render={({ field }) => (
                            <TagInput
                                value={field.value}
                                onChange={field.onChange}
                                placeholder="알고리즘, Java... (Enter로 추가)"
                            />
                        )}
                    />
                </div>

                <button type="submit" disabled={isSaving} className="btn-primary">
                    {isSaving ? '저장 중...' : '설정 저장'}
                </button>
            </form>

            <TistoryCategoryModal
                open={categoryModalOpen}
                blogUrl={blogUrl}
                selectedCategoryId={categoryId}
                onClose={() => setCategoryModalOpen(false)}
                onSelect={(category) => {
                    setValue('categoryId', category.categoryId, {
                        shouldDirty: true,
                        shouldValidate: true,
                    });
                    setSelectedCategoryLabel(category.label);
                }}
            />
        </>
    );
}
