import request from '@/utils/http';
import type { ApiResponse } from '@/types/api';
import type { Provider } from '@/types/cmdb';

const anntApi = {
  ProviderList: '/cmdb/providers',
  ProviderDetail: '/cmdb/provider',
};

export function listProviders(params?: any): Promise<ApiResponse<Provider[]>> {
  return request.get({
    url: anntApi.ProviderList,
    params,
  });
}

export function createProvider(data: Omit<Provider, 'id' | 'created_at' | 'updated_at'>): Promise<ApiResponse> {
  return request.post({
    url: anntApi.ProviderDetail,
    data,
  });
}

export function updateProvider(id: number, data: Partial<Provider>): Promise<ApiResponse> {
  return request.put({
    url: `${anntApi.ProviderDetail}/${id}`,
    data,
  });
}

export function deleteProvider(id: number): Promise<ApiResponse> {
  return request.delete({
    url: `${anntApi.ProviderDetail}/${id}`,
  });
}

export function syncProvider(id: number): Promise<ApiResponse> {
  return request.post({
    url: `${anntApi.ProviderDetail}/${id}/sync`,
  });
} 