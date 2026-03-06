import { useState } from 'react';
import {
    Box,
    Typography,
    Button,
    Card,
    CardContent,
    IconButton,
    Chip,
    Tooltip,
} from '@mui/material';
import {
    Add as AddIcon,
    Edit as EditIcon,
    Delete as DeleteIcon,
} from '@mui/icons-material';
import { DataGrid, type GridColDef, type GridRenderCellParams, type GridPaginationModel } from '@mui/x-data-grid';
import { toast } from 'sonner';

import {
    useGetCoursesQuery,
    useDeleteCourseMutation,
    type Course,
} from '@/api/courseApi';

export const CoursePage = () => {
    const [paginationModel, setPaginationModel] = useState<GridPaginationModel>({
        page: 0,
        pageSize: 10,
    });

    const { data, isLoading, isFetching } = useGetCoursesQuery({
        page: paginationModel.page + 1,
        limit: paginationModel.pageSize,
    });

    const [deleteCourse] = useDeleteCourseMutation();

    const handleEdit = (_course: Course) => {
        toast.info('Chức năng sửa khóa học đang được phát triển');
    };

    const handleDelete = async (id: string) => {
        if (window.confirm('Bạn có chắc chắn muốn xóa khóa học này?')) {
            try {
                await deleteCourse(id).unwrap();
                toast.success('Xóa khóa học thành công');
            } catch (error) {
                toast.error('Lỗi khi xóa khóa học');
            }
        }
    };

    const courses = data?.data?.courses || [];
    const totalItems = data?.data?.pagination?.total_items || 0;

    const columns: GridColDef[] = [
        { field: 'code', headerName: 'Mã KH', flex: 1, minWidth: 100 },
        { field: 'name', headerName: 'Tên Khóa Học', flex: 2, minWidth: 200 },
        { field: 'subject', headerName: 'Môn Học', flex: 1, minWidth: 150 },
        { field: 'grade_level', headerName: 'Khối Lớp', flex: 1, minWidth: 100 },
        {
            field: 'total_hours',
            headerName: 'Tổng Giờ',
            width: 100,
            renderCell: (params: GridRenderCellParams) => {
                return params.row.total_hours || 0;
            },
        },
        {
            field: 'status',
            headerName: 'Trạng Thái',
            width: 120,
            renderCell: (params: GridRenderCellParams) => {
                const isAct = params.value === 'ACTIVE';
                return (
                    <Chip
                        size="small"
                        label={isAct ? 'Hoạt động' : 'Tạm dừng'}
                        color={isAct ? 'success' : 'default'}
                    />
                );
            },
        },
        {
            field: 'actions',
            headerName: 'Thao tác',
            width: 120,
            sortable: false,
            renderCell: (params: GridRenderCellParams) => {
                const course = params.row as Course;
                return (
                    <Box sx={{ display: 'flex', gap: 1 }}>
                        <Tooltip title="Chỉnh sửa">
                            <IconButton size="small" color="primary" onClick={() => handleEdit(course)}>
                                <EditIcon fontSize="small" />
                            </IconButton>
                        </Tooltip>
                        <Tooltip title="Xóa">
                            <IconButton size="small" color="error" onClick={() => handleDelete(course.id)}>
                                <DeleteIcon fontSize="small" />
                            </IconButton>
                        </Tooltip>
                    </Box>
                );
            },
        },
    ];

    return (
        <Box sx={{ p: 3 }}>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 3 }}>
                <Typography variant="h5" component="h1">
                    Quản Lý Khóa Học
                </Typography>
                <Button
                    variant="contained"
                    startIcon={<AddIcon />}
                    onClick={() => toast.info('Chức năng thêm mới đang được phát triển')}
                >
                    Thêm Khóa Học
                </Button>
            </Box>

            <Card>
                <CardContent sx={{ p: 0 }}>
                    <DataGrid
                        rows={courses}
                        columns={columns}
                        rowCount={totalItems}
                        loading={isLoading || isFetching}
                        paginationMode="server"
                        paginationModel={paginationModel}
                        onPaginationModelChange={setPaginationModel}
                        pageSizeOptions={[10, 25, 50]}
                        disableRowSelectionOnClick
                        autoHeight
                        getRowId={(row) => row.id}
                        sx={{
                            border: 0,
                            '& .MuiDataGrid-cell:focus': { outline: 'none' },
                            '& .MuiDataGrid-row:hover': { backgroundColor: 'action.hover' },
                        }}
                    />
                </CardContent>
            </Card>
        </Box>
    );
};
