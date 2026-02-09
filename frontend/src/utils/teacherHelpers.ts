export const getEmploymentTypeLabel = (type: string): string => {
    switch (type) {
        case 'PART_TIME':
            return 'Part Time';
        case 'FULL_TIME':
            return 'Full Time';
        default:
            return type;
    }
};

export const getStatusLabel = (status: string): string => {
    switch (status) {
        case 'ACTIVE':
            return 'Active';
        case 'INACTIVE':
            return 'Inactive';
        default:
            return status;
    }
};

export const getStatusColor = (status: string): 'success' | 'default' | 'error' => {
    switch (status) {
        case 'ACTIVE':
            return 'success';
        case 'INACTIVE':
            return 'default';
        default:
            return 'default';
    }
};

export const formatDate = (dateString: string): string => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
    });
};

export const formatDateTime = (dateString: string): string => {
    const date = new Date(dateString);
    return date.toLocaleString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
    });
};
