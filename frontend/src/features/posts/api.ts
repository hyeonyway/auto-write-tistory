import { api } from '@/shared/api/client';
import type {
    PostListResponse,
    Post,
    PreviewRequest,
    PreviewResponse,
    CreatePostRequest,
    CreatePostResponse,
    UpdatePostRequest,
    UpdatePostResponse,
    PublishTistoryRequest,
    PublishTistoryResponse,
} from './types';

export const postsApi = {
    list: (page = 1, size = 20) => api.get<PostListResponse>(`/posts?page=${page}&size=${size}`),

    get: (id: number) => api.get<Post>(`/posts/${id}`),

    delete: (id: number) => api.delete<void>(`/posts/${id}`),

    preview: (req: PreviewRequest) => api.post<PreviewResponse>('/posts/preview', req),

    create: (req: CreatePostRequest) => api.post<CreatePostResponse>('/posts', req),

    update: (id: number, req: UpdatePostRequest) =>
        api.put<UpdatePostResponse>(`/posts/${id}`, req),

    publishTistory: (id: number, req: PublishTistoryRequest) =>
        api.post<PublishTistoryResponse>(`/posts/${id}/publish/tistory`, req),
};
