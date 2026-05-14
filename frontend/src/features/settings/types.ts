export interface TistorySettings {
    blogUrl: string;
    categoryId: string;
    categoryLabel: string;
    defaultVisibility: number;
    defaultTags: string[];
    categories: TistoryCategoryOption[];
    categoriesSyncedAt: string;
    session: TistorySessionInfo;
}

export interface SettingsResponse {
    tistory: TistorySettings;
}

export interface UpdateSettingsRequest {
    tistory: TistorySettingsInput;
}

export interface TistorySettingsInput {
    blogUrl: string;
    categoryId: string;
    defaultVisibility: number;
    defaultTags: string[];
}

export interface TistoryCategoryOption {
    categoryId: string;
    label: string;
}

export interface FetchTistoryCategoriesRequest {
    blogUrl: string;
}

export interface FetchTistoryCategoriesResponse {
    items: TistoryCategoryOption[];
}

export interface TistorySessionInfo {
    connected: boolean;
    message: string;
}
