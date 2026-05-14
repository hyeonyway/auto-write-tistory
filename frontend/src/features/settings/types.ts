import type { TistoryLoginInput } from '@/features/posts/types';

export interface TistorySettings {
    blogUrl: string;
    categoryId: string;
    defaultVisibility: number;
    defaultTags: string[];
}

export interface SettingsResponse {
    tistory: TistorySettings;
}

export interface UpdateSettingsRequest {
    tistory: TistorySettings;
}

export interface TistoryCategoryOption {
    categoryId: string;
    label: string;
}

export interface FetchTistoryCategoriesRequest {
    blogUrl: string;
    login: TistoryLoginInput;
}

export interface FetchTistoryCategoriesResponse {
    items: TistoryCategoryOption[];
}
