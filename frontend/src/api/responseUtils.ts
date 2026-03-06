type PaginationCandidate = {
    items_per_page?: number;
    total_items?: number;
    current_page?: number;
    total_pages?: number;
    limit?: number;
    ItemsPerPage?: number;
    TotalItems?: number;
    CurrentPage?: number;
    TotalPages?: number;
};

type Envelope = {
    success?: boolean;
    message?: string;
    data?: Record<string, unknown> | null;
};

export interface NormalizedPagination {
    items_per_page: number;
    total_items: number;
    current_page: number;
    total_pages: number;
}

export const normalizePagination = (pagination?: PaginationCandidate | null): NormalizedPagination => ({
    items_per_page: pagination?.items_per_page ?? pagination?.ItemsPerPage ?? pagination?.limit ?? 10,
    total_items: pagination?.total_items ?? pagination?.TotalItems ?? 0,
    current_page: pagination?.current_page ?? pagination?.CurrentPage ?? 1,
    total_pages: pagination?.total_pages ?? pagination?.TotalPages ?? 1,
});

export const normalizeListEnvelope = <T>(
    response: Envelope,
    keys: string[],
    paginationKeys: string[] = ['pagination', 'Pagination'],
) => {
    const data = response.data ?? {};
    const items = keys.reduce<T[]>((result, key) => {
        if (result.length > 0) {
            return result;
        }

        const value = data[key];
        return Array.isArray(value) ? (value as T[]) : result;
    }, []);

    const pagination = paginationKeys.reduce<PaginationCandidate | null>((result, key) => {
        if (result) {
            return result;
        }

        const value = data[key];
        return value && typeof value === 'object' ? (value as PaginationCandidate) : null;
    }, null);

    return {
        success: response.success ?? false,
        message: response.message ?? '',
        items,
        pagination: normalizePagination(pagination),
    };
};

export const normalizeDetailEnvelope = <T>(
    response: Envelope,
    keys: string[] = [],
): {
    success: boolean;
    message: string;
    data: T | null;
} => {
    const data = response.data;

    if (data && typeof data === 'object') {
        for (const key of keys) {
            const value = data[key];
            if (value && typeof value === 'object') {
                return {
                    success: response.success ?? false,
                    message: response.message ?? '',
                    data: value as T,
                };
            }
        }
    }

    return {
        success: response.success ?? false,
        message: response.message ?? '',
        data: (data as T | null) ?? null,
    };
};
