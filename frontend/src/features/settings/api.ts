import { api } from '@/shared/api/client';
import type {
    FetchTistoryCategoriesRequest,
    FetchTistoryCategoriesResponse,
    SettingsResponse,
    UpdateSettingsRequest,
} from './types';

export const settingsApi = {
    get: () => api.get<SettingsResponse>('/settings'),
    update: (req: UpdateSettingsRequest) => api.put<SettingsResponse>('/settings', req),
    fetchCategories: (req: FetchTistoryCategoriesRequest) =>
        api.post<FetchTistoryCategoriesResponse>('/tistory/categories/fetch', req),
};
