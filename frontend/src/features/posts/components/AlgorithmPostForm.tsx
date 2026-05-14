import { useEffect } from 'react';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import TagInput from '@/shared/components/TagInput';
import type { AlgorithmInput } from '../types';

const schema = z.object({
    platform: z.string().min(1, '플랫폼을 선택해주세요'),
    problemTitle: z.string().min(1, '문제 제목을 입력해주세요'),
    problemUrl: z.string().url('올바른 URL을 입력해주세요'),
    language: z.string().min(1, '언어를 선택해주세요'),
    categories: z.array(z.string()).default([]),
    approach: z.string().min(1, '접근 방식을 입력해주세요'),
    code: z.string().min(1, '코드를 입력해주세요'),
    timeComplexity: z.string().default(''),
    spaceComplexity: z.string().default(''),
    review: z.string().default(''),
    tags: z.array(z.string()).default([]),
});

type FormValues = z.infer<typeof schema>;

const PLATFORMS = ['Programmers', 'Baekjoon', 'LeetCode', 'Softeer', 'SWEA'];
const LANGUAGES = [
    'Java',
    'Python',
    'Go',
    'TypeScript',
    'JavaScript',
    'C++',
    'C',
    'Kotlin',
    'Rust',
];

interface AlgorithmPostFormProps {
    onPreview: (input: AlgorithmInput) => void;
    onSave: (input: AlgorithmInput) => void;
    onPublish: (input: AlgorithmInput) => void;
    initialValue?: AlgorithmInput;
    isPreviewLoading?: boolean;
    isSaveLoading?: boolean;
    saveLabel?: string;
}

export default function AlgorithmPostForm({
    onPreview,
    onSave,
    onPublish,
    initialValue,
    isPreviewLoading,
    isSaveLoading,
    saveLabel = '초안 저장',
}: AlgorithmPostFormProps) {
    const {
        register,
        control,
        handleSubmit,
        getValues,
        reset,
        formState: { errors },
    } = useForm<FormValues>({
        resolver: zodResolver(schema),
        defaultValues: initialValue ?? { categories: [], tags: [] },
    });

    useEffect(() => {
        if (initialValue) {
            reset(initialValue);
        }
    }, [initialValue, reset]);

    const getInput = (): AlgorithmInput => getValues() as AlgorithmInput;

    return (
        <form className="flex flex-col h-full">
            <div className="flex-1 overflow-y-auto p-4 space-y-4">
                {/* Platform */}
                <div>
                    <label className="label">플랫폼 *</label>
                    <select {...register('platform')} className="input-field">
                        <option value="">선택</option>
                        {PLATFORMS.map((p) => (
                            <option key={p} value={p}>
                                {p}
                            </option>
                        ))}
                    </select>
                    {errors.platform && (
                        <p className="text-error text-xs mt-1">{errors.platform.message}</p>
                    )}
                </div>

                {/* Problem title */}
                <div>
                    <label className="label">문제 제목 *</label>
                    <input
                        {...register('problemTitle')}
                        className="input-field"
                        placeholder="등굣길"
                    />
                    {errors.problemTitle && (
                        <p className="text-error text-xs mt-1">{errors.problemTitle.message}</p>
                    )}
                </div>

                {/* Problem URL */}
                <div>
                    <label className="label">문제 링크 *</label>
                    <input
                        {...register('problemUrl')}
                        type="url"
                        className="input-field"
                        placeholder="https://school.programmers.co.kr/..."
                    />
                    {errors.problemUrl && (
                        <p className="text-error text-xs mt-1">{errors.problemUrl.message}</p>
                    )}
                </div>

                {/* Language */}
                <div>
                    <label className="label">언어 *</label>
                    <select {...register('language')} className="input-field">
                        <option value="">선택</option>
                        {LANGUAGES.map((l) => (
                            <option key={l} value={l}>
                                {l}
                            </option>
                        ))}
                    </select>
                    {errors.language && (
                        <p className="text-error text-xs mt-1">{errors.language.message}</p>
                    )}
                </div>

                {/* Categories */}
                <div>
                    <label className="label">문제 유형</label>
                    <Controller
                        name="categories"
                        control={control}
                        render={({ field }) => (
                            <TagInput
                                value={field.value}
                                onChange={field.onChange}
                                placeholder="DP, 그리디... (Enter로 추가)"
                            />
                        )}
                    />
                </div>

                {/* Approach */}
                <div>
                    <label className="label">접근 방식 *</label>
                    <textarea
                        {...register('approach')}
                        rows={4}
                        className="input-field resize-none"
                        placeholder="문제 풀이 접근 방식을 설명해주세요."
                    />
                    {errors.approach && (
                        <p className="text-error text-xs mt-1">{errors.approach.message}</p>
                    )}
                </div>

                {/* Code */}
                <div>
                    <label className="label">풀이 코드 *</label>
                    <textarea
                        {...register('code')}
                        rows={8}
                        className="w-full px-3 py-2 bg-code-bg text-code-text border border-transparent rounded-lg text-xs font-mono resize-none focus:outline-none focus:ring-2 focus:ring-primary/30"
                        placeholder="// 코드를 입력하세요"
                        spellCheck={false}
                    />
                    {errors.code && (
                        <p className="text-error text-xs mt-1">{errors.code.message}</p>
                    )}
                </div>

                {/* Complexity */}
                <div className="grid grid-cols-2 gap-3">
                    <div>
                        <label className="label">시간복잡도</label>
                        <input
                            {...register('timeComplexity')}
                            className="input-field"
                            placeholder="O(NM)"
                        />
                    </div>
                    <div>
                        <label className="label">공간복잡도</label>
                        <input
                            {...register('spaceComplexity')}
                            className="input-field"
                            placeholder="O(NM)"
                        />
                    </div>
                </div>

                {/* Review */}
                <div>
                    <label className="label">회고</label>
                    <textarea
                        {...register('review')}
                        rows={3}
                        className="input-field resize-none"
                        placeholder="풀고 나서 느낀 점, 다른 방법 등을 자유롭게 적어주세요."
                    />
                </div>

                {/* Tags */}
                <div>
                    <label className="label">Tistory 태그</label>
                    <Controller
                        name="tags"
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
            </div>

            {/* Action bar */}
            <div className="border-t border-border px-4 py-3 flex items-center gap-2 flex-shrink-0">
                <button
                    type="button"
                    onClick={() => onPreview(getInput())}
                    disabled={isPreviewLoading}
                    className="btn-secondary text-xs px-3 py-1.5"
                >
                    {isPreviewLoading ? '생성 중...' : '미리보기 생성'}
                </button>
                <button
                    type="button"
                    onClick={handleSubmit(() => onSave(getInput()))}
                    disabled={isSaveLoading}
                    className="btn-secondary text-xs px-3 py-1.5"
                >
                    {isSaveLoading ? '저장 중...' : saveLabel}
                </button>
                <button
                    type="button"
                    onClick={handleSubmit(() => onPublish(getInput()))}
                    className="btn-primary text-xs px-3 py-1.5 ml-auto"
                >
                    Tistory 비공개 발행
                </button>
            </div>
        </form>
    );
}
