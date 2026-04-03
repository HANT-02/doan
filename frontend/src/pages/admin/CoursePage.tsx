import { useMemo, useState } from 'react';
import {
    Alert,
    Box,
    Button,
    Chip,
    IconButton,
    InputAdornment,
    MenuItem,
    Paper,
    Skeleton,
    Stack,
    TextField,
    Tooltip,
    Typography,
} from '@mui/material';
import {
    AddRounded,
    DeleteOutlineRounded,
    EditOutlined,
    MenuBookRounded,
    RefreshRounded,
    SearchRounded,
} from '@mui/icons-material';
import { DataGrid, type GridColDef, type GridRenderCellParams } from '@mui/x-data-grid';
import { toast } from 'sonner';

import {
    useCreateCourseMutation,
    useDeleteCourseMutation,
    useGetCoursesQuery,
    useUpdateCourseMutation,
    type Course,
} from '@/api/courseApi';
import CourseDialog from '@/components/admin/CourseDialog';
import ConfirmDialog from '@/components/common/ConfirmDialog';
import PageHeader from '@/components/common/PageHeader';
import { getApiErrorMessage } from '@/utils/apiError';

const statusMeta: Record<string, { label: string; color: 'success' | 'default' }> = {
    ACTIVE: { label: 'Hoạt động', color: 'success' },
    INACTIVE: { label: 'Tạm dừng', color: 'default' },
};

const getCourseSubtitle = (course: Course) => {
    if (course.description?.trim()) {
        return course.description;
    }

    const subject = course.subject?.trim() || 'Chưa phân môn';
    const gradeLevel = course.grade_level?.trim() || 'Chưa khai báo khối lớp';
    return `${subject} • ${gradeLevel}`;
};

export const CoursePage = () => {
    const [page, setPage] = useState(0);
    const [pageSize, setPageSize] = useState(10);
    const [search, setSearch] = useState('');
    const [statusFilter, setStatusFilter] = useState('');
    const [isDialogOpen, setIsDialogOpen] = useState(false);
    const [selectedCourse, setSelectedCourse] = useState<Course | null>(null);
    const [deleteTarget, setDeleteTarget] = useState<Course | null>(null);

    const queryParams = useMemo(
        () => ({
            page: page + 1,
            limit: pageSize,
            search: search || undefined,
            status: statusFilter || undefined,
        }),
        [page, pageSize, search, statusFilter],
    );

    const {
        data,
        isLoading,
        isFetching,
        isError,
        error,
        refetch,
    } = useGetCoursesQuery(queryParams);

    const [createCourse, { isLoading: isCreating }] = useCreateCourseMutation();
    const [updateCourse, { isLoading: isUpdating }] = useUpdateCourseMutation();
    const [deleteCourse, { isLoading: isDeleting }] = useDeleteCourseMutation();

    const courses = data?.data?.courses || [];
    const totalItems = data?.data?.pagination?.total_items || 0;

    const handleCloseDialog = () => {
        setIsDialogOpen(false);
        setSelectedCourse(null);
    };

    const handleOpenCreate = () => {
        setSelectedCourse(null);
        setIsDialogOpen(true);
    };

    const handleOpenEdit = (course: Course) => {
        setSelectedCourse(course);
        setIsDialogOpen(true);
    };

    const handleSubmitCourse = async (formData: Partial<Course>) => {
        try {
            if (selectedCourse) {
                await updateCourse({ id: selectedCourse.id, ...formData }).unwrap();
                toast.success('Cập nhật khóa học thành công');
            } else {
                await createCourse(formData).unwrap();
                toast.success('Tạo khóa học thành công');
            }
        } catch (submitError) {
            toast.error(getApiErrorMessage(submitError, 'Không thể lưu khóa học'));
            throw submitError;
        }
    };

    const handleDelete = async () => {
        if (!deleteTarget) {
            return;
        }

        try {
            await deleteCourse(deleteTarget.id).unwrap();
            toast.success('Xóa khóa học thành công');
            setDeleteTarget(null);
        } catch (deleteError) {
            toast.error(getApiErrorMessage(deleteError, 'Lỗi khi xóa khóa học'));
        }
    };

    const columns: GridColDef<Course>[] = [
        {
            field: 'code',
            headerName: 'Mã KH',
            width: 130,
            align: 'left',
            headerAlign: 'left',
            renderCell: (params: GridRenderCellParams<Course>) => (
                <Typography variant="body2" sx={{ fontWeight: 700, color: 'primary.main' }}>
                    {params.row.code}
                </Typography>
            ),
        },
        {
            field: 'name',
            headerName: 'Tên khóa học',
            minWidth: 260,
            flex: 1.5,
            align: 'left',
            headerAlign: 'left',
            renderCell: (params: GridRenderCellParams<Course>) => (
                <Stack direction="row" spacing={1} alignItems="center" sx={{ minWidth: 0, height: '100%' }}>
                    <MenuBookRounded sx={{ color: 'primary.main', fontSize: 18 }} />
                    <Box
                        sx={{
                            minWidth: 0,
                            display: 'flex',
                            flexDirection: 'column',
                            justifyContent: 'center',
                            py: 1,
                        }}
                    >
                        <Typography variant="body2" sx={{ fontWeight: 700, lineHeight: 1.35 }} noWrap>
                            {params.row.name}
                        </Typography>
                        <Typography
                            variant="caption"
                            color="text.secondary"
                            noWrap
                            sx={{ display: 'block', lineHeight: 1.4, mt: 0.25 }}
                        >
                            {getCourseSubtitle(params.row)}
                        </Typography>
                    </Box>
                </Stack>
            ),
        },
        {
            field: 'subject',
            headerName: 'Môn học',
            minWidth: 150,
            flex: 0.8,
            align: 'left',
            headerAlign: 'left',
            renderCell: (params: GridRenderCellParams<Course>) => (
                <Typography variant="body2">{params.row.subject || '-'}</Typography>
            ),
        },
        {
            field: 'grade_level',
            headerName: 'Khối lớp',
            width: 120,
            align: 'center',
            headerAlign: 'center',
            renderCell: (params: GridRenderCellParams<Course>) => (
                <Typography variant="body2">{params.row.grade_level || '-'}</Typography>
            ),
        },
        {
            field: 'total_hours',
            headerName: 'Tổng giờ',
            width: 110,
            align: 'center',
            headerAlign: 'center',
            renderCell: (params: GridRenderCellParams<Course>) => (
                <Typography variant="body2">{params.row.total_hours || 0}</Typography>
            ),
        },
        {
            field: 'status',
            headerName: 'Trạng thái',
            width: 130,
            align: 'center',
            headerAlign: 'center',
            renderCell: (params: GridRenderCellParams<Course>) => {
                const meta = statusMeta[params.row.status] || { label: params.row.status || 'Chưa rõ', color: 'default' as const };
                return <Chip size="small" label={meta.label} color={meta.color} variant="outlined" />;
            },
        },
        {
            field: 'actions',
            headerName: '',
            width: 104,
            sortable: false,
            align: 'center',
            renderCell: (params: GridRenderCellParams<Course>) => (
                <Stack direction="row" spacing={0.5}>
                    <Tooltip title="Chỉnh sửa">
                        <IconButton size="small" color="primary" onClick={() => handleOpenEdit(params.row)}>
                            <EditOutlined fontSize="small" />
                        </IconButton>
                    </Tooltip>
                    <Tooltip title="Xóa">
                        <IconButton size="small" color="error" onClick={() => setDeleteTarget(params.row)}>
                            <DeleteOutlineRounded fontSize="small" />
                        </IconButton>
                    </Tooltip>
                </Stack>
            ),
        },
    ];

    const renderLoadingState = () => (
        <Stack spacing={1.5}>
            <Skeleton variant="rounded" height={72} />
            <Skeleton variant="rounded" height={420} />
        </Stack>
    );

    const renderEmptyState = () => (
        <Paper
            variant="outlined"
            sx={{
                p: 6,
                borderRadius: 4,
                textAlign: 'center',
                borderStyle: 'dashed',
            }}
        >
            <MenuBookRounded sx={{ fontSize: 48, color: 'text.secondary', mb: 2 }} />
            <Typography variant="h6" sx={{ fontWeight: 700, mb: 1 }}>
                Chưa có khóa học nào phù hợp bộ lọc
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                Thử đổi bộ lọc hiện tại hoặc tạo khóa học mới để liên kết vào chương trình đào tạo.
            </Typography>
            <Button variant="contained" startIcon={<AddRounded />} onClick={handleOpenCreate}>
                Thêm khóa học mới
            </Button>
        </Paper>
    );

    return (
        <Box>
            <PageHeader
                title="Quản lý khóa học"
                subtitle="Theo dõi danh sách khóa học, chỉnh sửa thông tin và chuẩn bị dữ liệu để liên kết vào các chương trình đào tạo."
                breadcrumbs={[
                    { label: 'Tổng quan', path: '/app/admin/overview' },
                    { label: 'Khóa học' },
                ]}
                actions={
                    <Stack direction="row" spacing={1}>
                        <Button variant="outlined" startIcon={<RefreshRounded />} onClick={() => void refetch()} disabled={isFetching}>
                            Làm mới
                        </Button>
                        <Button variant="contained" startIcon={<AddRounded />} onClick={handleOpenCreate}>
                            Thêm khóa học
                        </Button>
                    </Stack>
                }
            />

            {isLoading && courses.length === 0 ? renderLoadingState() : null}
            {!isLoading && isError ? (
                <Alert
                    severity="error"
                    action={
                        <Button color="inherit" size="small" startIcon={<RefreshRounded />} onClick={() => void refetch()}>
                            Tải lại
                        </Button>
                    }
                    sx={{ mb: 3 }}
                >
                    {getApiErrorMessage(error, 'Không thể tải danh sách khóa học')}
                </Alert>
            ) : null}

            {!isLoading && !isError ? (
                courses.length === 0 ? (
                    renderEmptyState()
                ) : (
                    <Paper elevation={0} sx={{ p: 2.5, borderRadius: 4, border: '1px solid #e2e8f0' }}>
                        <Stack
                            direction={{ xs: 'column', lg: 'row' }}
                            spacing={1.5}
                            justifyContent="space-between"
                            alignItems={{ xs: 'stretch', lg: 'center' }}
                            sx={{ mb: 2 }}
                        >
                            <Stack direction={{ xs: 'column', md: 'row' }} spacing={1.5}>
                                <TextField
                                    value={search}
                                    onChange={(event) => {
                                        setSearch(event.target.value);
                                        setPage(0);
                                    }}
                                    size="small"
                                    placeholder="Tìm theo mã, tên hoặc môn học"
                                    sx={{ minWidth: { xs: '100%', md: 280 } }}
                                    InputProps={{
                                        startAdornment: (
                                            <InputAdornment position="start">
                                                <SearchRounded fontSize="small" />
                                            </InputAdornment>
                                        ),
                                    }}
                                />
                                <TextField
                                    select
                                    size="small"
                                    label="Trạng thái"
                                    value={statusFilter}
                                    onChange={(event) => {
                                        setStatusFilter(event.target.value);
                                        setPage(0);
                                    }}
                                    sx={{ minWidth: { xs: '100%', md: 180 } }}
                                >
                                    <MenuItem value="">Tất cả trạng thái</MenuItem>
                                    <MenuItem value="ACTIVE">Hoạt động</MenuItem>
                                    <MenuItem value="INACTIVE">Tạm dừng</MenuItem>
                                </TextField>
                            </Stack>

                            <Chip
                                label={`Hiển thị ${courses.length}/${totalItems || courses.length} khóa học`}
                                color="primary"
                                variant="outlined"
                            />
                        </Stack>

                        <DataGrid
                            rows={courses}
                            columns={columns}
                            rowCount={totalItems}
                            loading={isFetching}
                            paginationMode="server"
                            paginationModel={{ page, pageSize }}
                            onPaginationModelChange={(model) => {
                                setPage(model.page);
                                setPageSize(model.pageSize);
                            }}
                            pageSizeOptions={[10, 25, 50]}
                            disableRowSelectionOnClick
                            autoHeight
                            getRowId={(row) => row.id}
                            localeText={{ noRowsLabel: 'Không có dữ liệu khóa học' }}
                            sx={{
                                border: 0,
                                '& .MuiDataGrid-columnHeaders': {
                                    backgroundColor: '#f8fafc',
                                },
                                '& .MuiDataGrid-cell:focus': { outline: 'none' },
                            }}
                        />
                    </Paper>
                )
            ) : null}

            <ConfirmDialog
                open={!!deleteTarget}
                title="Xóa khóa học"
                message={
                    deleteTarget
                        ? `Bạn có chắc chắn muốn xóa khóa học "${deleteTarget.name}"? Hành động này không thể hoàn tác.`
                        : ''
                }
                confirmText="Xóa khóa học"
                isDanger
                onClose={() => setDeleteTarget(null)}
                onConfirm={() => void handleDelete()}
                loading={isDeleting}
            />

            <CourseDialog
                open={isDialogOpen}
                onClose={handleCloseDialog}
                onSubmit={handleSubmitCourse}
                course={selectedCourse}
                isLoading={isCreating || isUpdating}
            />
        </Box>
    );
};
