import { api } from '@/shared/api/client';
import type {
    FetchTistoryCategoriesRequest,
    FetchTistoryCategoriesResponse,
    SettingsResponse,
    TistorySessionInfo,
    UpdateSettingsRequest,
} from './types';

export const settingsApi = {
    get: () => api.get<SettingsResponse>('/settings'),
    update: (req: UpdateSettingsRequest) => api.put<SettingsResponse>('/settings', req),
    startSession: () => api.post<TistorySessionInfo>('/tistory/session/start', {}),
    confirmSession: () => api.post<TistorySessionInfo>('/tistory/session/confirm', {}),
    sessionStatus: () => api.get<TistorySessionInfo>('/tistory/session/status'),
    deleteSession: () => api.delete<{ deleted: boolean }>('/tistory/session'),
    fetchCategories: (req: FetchTistoryCategoriesRequest) =>
        api.post<FetchTistoryCategoriesResponse>('/tistory/categories/fetch', req),
};
