export type PostStatus = 'DRAFT' | 'PUBLISHED' | 'FAILED';

export interface Post {
    id: number;
    title: string;
    type: string;
    sourcePlatform: string | null;
    sourceUrl: string | null;
    language: string | null;
    contentMarkdown: string;
    status: PostStatus;
    externalProvider: string | null;
    externalUrl: string | null;
    errorMessage: string | null;
    publishedAt: string | null;
    createdAt: string;
    updatedAt: string;
    input?: AlgorithmInput;
}

export interface PostListResponse {
    items: Post[];
    page: number;
    size: number;
    total: number;
}

export interface AlgorithmInput {
    platform: string;
    problemTitle: string;
    problemUrl: string;
    language: string;
    categories: string[];
    approach: string;
    code: string;
    timeComplexity: string;
    spaceComplexity: string;
    review: string;
    tags: string[];
}

export interface PreviewRequest {
    type: 'ALGORITHM';
    input: AlgorithmInput;
}

export interface PreviewResponse {
    title: string;
    contentMarkdown: string;
}

export interface CreatePostRequest {
    type: 'ALGORITHM';
    input: AlgorithmInput;
}

export interface CreatePostResponse {
    id: number;
    title: string;
    status: PostStatus;
}

export interface UpdatePostRequest {
    type: 'ALGORITHM';
    input: AlgorithmInput;
}

export interface UpdatePostResponse {
    id: number;
    title: string;
    status: PostStatus;
}

export interface TistoryLoginInput {
    id: string;
    password: string;
}

export interface PublishTistoryRequest {
    visibility: number;
    categoryId: string;
    tags: string[];
    login: TistoryLoginInput;
}

export interface PublishTistoryResponse {
    success: boolean;
    externalUrl: string | null;
}
