import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import Sidebar from '@/shared/components/Sidebar';
import Header from '@/shared/components/Header';
import StatusBadge from '@/shared/components/StatusBadge';
import { postsApi } from '@/features/posts/api';
import type { Post, PostStatus } from '@/features/posts/types';

function StatCard({ label, value, color }: { label: string; value: number; color: string }) {
    return (
        <div className="card px-5 py-4 min-w-[160px]">
            <div className={`text-2xl font-bold ${color}`}>{value}</div>
            <div className="text-xs text-text-sec mt-2">{label}</div>
        </div>
    );
}

const statusCount = (posts: Post[], status: PostStatus) =>
    posts.filter((p) => p.status === status).length;

export default function DashboardPage() {
    const queryClient = useQueryClient();
    const { data, isLoading, isError } = useQuery({
        queryKey: ['posts'],
        queryFn: () => postsApi.list(),
    });

    const deleteMutation = useMutation({
        mutationFn: postsApi.delete,
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['posts'] }),
    });

    const posts = data?.items ?? [];
    const total = data?.total ?? 0;

    return (
        <div className="flex h-screen bg-surface">
            <Sidebar />

            <div className="flex flex-col flex-1 ml-60 min-h-0">
                <Header
                    title="작성 이력"
                    action={
                        <Link to="/posts/new/algorithm" className="btn-primary text-sm">
                            + 새 글 만들기
                        </Link>
                    }
                />

                <main className="flex-1 overflow-y-auto p-6 space-y-5">
                    {/* Stat cards */}
                    <div className="flex gap-4">
                        <StatCard label="전체 글" value={total} color="text-text" />
                        <StatCard
                            label="발행됨"
                            value={statusCount(posts, 'PUBLISHED')}
                            color="text-success"
                        />
                        <StatCard
                            label="초안"
                            value={statusCount(posts, 'DRAFT')}
                            color="text-primary"
                        />
                        <StatCard
                            label="실패"
                            value={statusCount(posts, 'FAILED')}
                            color="text-error"
                        />
                    </div>

                    {/* Posts table */}
                    <div className="card overflow-hidden">
                        {isLoading && (
                            <div className="py-20 text-center text-text-hint text-sm">
                                불러오는 중...
                            </div>
                        )}
                        {isError && (
                            <div className="py-20 text-center text-error text-sm">
                                데이터를 불러오지 못했습니다.
                            </div>
                        )}
                        {!isLoading && !isError && (
                            <table className="w-full text-sm">
                                <thead>
                                    <tr className="bg-surface border-b border-border">
                                        <th className="text-left px-4 py-3 text-[11px] font-semibold text-text-sec">
                                            제목
                                        </th>
                                        <th className="text-left px-4 py-3 text-[11px] font-semibold text-text-sec w-20">
                                            유형
                                        </th>
                                        <th className="text-left px-4 py-3 text-[11px] font-semibold text-text-sec w-28">
                                            플랫폼
                                        </th>
                                        <th className="text-left px-4 py-3 text-[11px] font-semibold text-text-sec w-24">
                                            상태
                                        </th>
                                        <th className="text-left px-4 py-3 text-[11px] font-semibold text-text-sec w-28">
                                            생성일
                                        </th>
                                        <th className="text-left px-4 py-3 text-[11px] font-semibold text-text-sec w-44">
                                            Tistory URL
                                        </th>
                                        <th className="w-20" />
                                    </tr>
                                </thead>
                                <tbody>
                                    {posts.length === 0 ? (
                                        <tr>
                                            <td
                                                colSpan={7}
                                                className="py-20 text-center text-text-hint"
                                            >
                                                아직 작성한 글이 없습니다.{' '}
                                                <Link
                                                    to="/posts/new/algorithm"
                                                    className="text-primary underline"
                                                >
                                                    첫 글 작성하기 →
                                                </Link>
                                            </td>
                                        </tr>
                                    ) : (
                                        posts.map((post) => (
                                            <tr
                                                key={post.id}
                                                className="border-t border-border hover:bg-surface/60"
                                            >
                                                <td className="px-4 py-3 font-medium text-text truncate max-w-xs">
                                                    {post.title}
                                                </td>
                                                <td className="px-4 py-3">
                                                    <span className="px-2 py-0.5 bg-tag text-tag-text text-[10px] font-semibold rounded">
                                                        ALGO
                                                    </span>
                                                </td>
                                                <td className="px-4 py-3 text-text-sec">
                                                    {post.sourcePlatform ?? '-'}
                                                </td>
                                                <td className="px-4 py-3">
                                                    <StatusBadge status={post.status} />
                                                </td>
                                                <td className="px-4 py-3 text-text-sec">
                                                    {post.createdAt.slice(0, 10)}
                                                </td>
                                                <td className="px-4 py-3">
                                                    {post.externalUrl ? (
                                                        <a
                                                            href={post.externalUrl}
                                                            target="_blank"
                                                            rel="noreferrer"
                                                            className="text-primary text-xs hover:underline truncate block max-w-[160px]"
                                                        >
                                                            {post.externalUrl.replace(
                                                                'https://',
                                                                '',
                                                            )}
                                                        </a>
                                                    ) : (
                                                        <span className="text-text-hint text-xs">
                                                            -
                                                        </span>
                                                    )}
                                                </td>
                                                <td className="px-2 py-3">
                                                    <Link
                                                        to={`/posts/${post.id}/edit`}
                                                        className="text-primary hover:underline text-xs mr-3"
                                                    >
                                                        수정
                                                    </Link>
                                                    <button
                                                        onClick={() => {
                                                            if (
                                                                confirm('이 글을 삭제하시겠습니까?')
                                                            ) {
                                                                deleteMutation.mutate(post.id);
                                                            }
                                                        }}
                                                        className="text-text-hint hover:text-error text-xs transition-colors"
                                                    >
                                                        삭제
                                                    </button>
                                                </td>
                                            </tr>
                                        ))
                                    )}
                                </tbody>
                            </table>
                        )}
                    </div>
                </main>
            </div>
        </div>
    );
}
