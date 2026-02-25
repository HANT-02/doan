import { useMemo } from 'react';
import {
    DataGrid,
    type GridColDef,
    type GridRenderCellParams,
    GridActionsCellItem,
} from '@mui/x-data-grid';
import {
    Chip,
    Typography,
    Box,
    Paper,
} from '@mui/material';
import {
    Visibility,
    Edit,
    Delete,
} from '@mui/icons-material';
import { type Teacher } from '@/api/teacherApi';
import { getStatusColor, getEmploymentTypeLabel } from '@/utils/teacherHelpers';

interface TeacherListTableProps {
    teachers: Teacher[];
    loading?: boolean;
    onView?: (teacher: Teacher) => void;
    onEdit?: (teacher: Teacher) => void;
    onDelete?: (teacher: Teacher) => void;
    showActions?: boolean;
    isAdmin?: boolean;
    paginationModel: { page: number; pageSize: number };
    onPaginationModelChange: (model: { page: number; pageSize: number }) => void;
    rowCount: number;
}

export const TeacherListTable = ({
    teachers,
    loading = false,
    onView,
    onEdit,
    onDelete,
    showActions = true,
    isAdmin = false,
    paginationModel,
    onPaginationModelChange,
    rowCount,
}: TeacherListTableProps) => {

    const columns = useMemo<GridColDef<Teacher>[]>(() => {
        const cols: GridColDef<Teacher>[] = [
            {
                field: 'code',
                headerName: 'Mã số',
                width: 100,
                renderCell: (params) => params.value || '-'
            },
            {
                field: 'full_name',
                headerName: 'Họ và tên',
                flex: 1,
                minWidth: 200,
                renderCell: (params: GridRenderCellParams<Teacher>) => (
                    <Box sx={{ py: 1 }}>
                        <Typography variant="body2" sx={{ fontWeight: 600 }}>
                            {params.row.full_name}
                        </Typography>
                        {params.row.is_school_teacher && (
                            <Typography variant="caption" color="text.secondary" sx={{ display: 'block' }}>
                                {params.row.school_name || 'Giáo viên trường'}
                            </Typography>
                        )}
                    </Box>
                )
            },
            {
                field: 'email',
                headerName: 'Email',
                width: 200,
                renderCell: (params) => params.value || '-'
            },
            {
                field: 'phone',
                headerName: 'Điện thoại',
                width: 150,
                renderCell: (params) => params.value || '-'
            },
            {
                field: 'employment_type',
                headerName: 'Loại hình',
                width: 150,
                renderCell: (params) => (
                    <Typography variant="body2">
                        {getEmploymentTypeLabel(params.value)}
                    </Typography>
                )
            },
            {
                field: 'status',
                headerName: 'Trạng thái',
                width: 150,
                renderCell: (params) => (
                    <Chip
                        label={params.value}
                        color={getStatusColor(params.value)}
                        size="small"
                        sx={{ fontWeight: 500 }}
                    />
                )
            }
        ];

        if (showActions) {
            cols.push({
                field: 'actions',
                type: 'actions',
                headerName: 'Hành động',
                width: 120,
                getActions: (params) => {
                    const actions = [];
                    if (onView) {
                        actions.push(
                            <GridActionsCellItem
                                icon={<Visibility />}
                                label="Xem"
                                onClick={() => onView(params.row)}
                            />
                        );
                    }
                    if (isAdmin && onEdit) {
                        actions.push(
                            <GridActionsCellItem
                                icon={<Edit />}
                                label="Sửa"
                                onClick={() => onEdit(params.row)}
                            />
                        );
                    }
                    if (isAdmin && onDelete) {
                        actions.push(
                            <GridActionsCellItem
                                icon={<Delete color="error" />}
                                label="Xóa"
                                onClick={() => onDelete(params.row)}
                            />
                        );
                    }
                    return actions;
                }
            });
        }

        return cols;
    }, [onView, onEdit, onDelete, isAdmin, showActions]);

    return (
        <Paper sx={{ height: 600, width: '100%', borderRadius: 3, overflow: 'hidden' }}>
            <DataGrid
                rows={teachers}
                columns={columns}
                loading={loading}
                rowCount={rowCount}
                paginationMode="server"
                paginationModel={paginationModel}
                onPaginationModelChange={onPaginationModelChange}
                pageSizeOptions={[10, 25, 50]}
                disableRowSelectionOnClick
                sx={{
                    border: 'none',
                    '& .MuiDataGrid-columnHeaderTitle': {
                        fontWeight: 700,
                    },
                }}
            />
        </Paper>
    );
};
