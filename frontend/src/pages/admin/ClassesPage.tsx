import { type MouseEvent, useMemo, useState } from 'react';
import {
    Alert,
    Box,
    Button,
    Chip,
    IconButton,
    InputAdornment,
    Menu,
    MenuItem,
    Paper,
    Skeleton,
    Stack,
    TextField,
    Typography,
} from '@mui/material';
import {
    AddRounded,
    DeleteOutlineRounded,
    EditOutlined,
    MoreVert,
    RefreshRounded,
    SchoolOutlined,
    SearchRounded,
    VisibilityOutlined,
} from '@mui/icons-material';
import { DataGrid } from '@mui/x-data-grid';
import type { GridColDef, GridRenderCellParams } from '@mui/x-data-grid';
import { format } from 'date-fns';
import { toast } from 'sonner';

import {
    useCreateClassMutation,
    useDeleteClassMutation,
    useGetClassesQuery,
    useUpdateClassMutation,
    type Class,
} from '@/api/classApi';
import { useGetTeachersQuery } from '@/api/teacherApi';
import ClassDialog from '@/components/admin/ClassDialog';
import ClassDetailDialog from '@/components/admin/ClassDetailDialog';
import ConfirmDialog from '@/components/common/ConfirmDialog';
import PageHeader from '@/components/common/PageHeader';

const getErrorMessage = (error: unknown, fallback: string) => {
    if (typeof error === 'object' && error && 'data' in error) {
        const apiError = error as { data?: { message?: string } };
        return apiError.data?.message || fallback;
    }

    if (error instanceof Error) {
        return error.message;
    }

    return fallback;
};

const statusMeta: Record<string, { label: string; color: 'success' | 'warning' | 'error' | 'default' }> = {
    OPEN: { label: 'Đang mở', color: 'success' },
    CLOSED: { label: 'Đã đóng', color: 'default' },
    CANCELLED: { label: 'Đã hủy', color: 'error' },
};

export const ClassesPage = () => {
    const [page, setPage] = useState(0);
    const [pageSize, setPageSize] = useState(10);
    const [search, setSearch] = useState('');
    const [statusFilter, setStatusFilter] = useState('');

    const [isDialogOpen, setIsDialogOpen] = useState(false);
    const [selectedClass, setSelectedClass] = useState<Class | null>(null);
    const [selectedClassId, setSelectedClassId] = useState<string | null>(null);
    const [isConfirmOpen, setIsConfirmOpen] = useState(false);
    const [classToDelete, setClassToDelete] = useState<Class | null>(null);

    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    const [menuClass, setMenuClass] = useState<Class | null>(null);

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
        refetch,
    } = useGetClassesQuery(queryParams);
    const { data: teachersData } = useGetTeachersQuery({ page: 1, limit: 200 });

    const [createClass, { isLoading: isCreating }] = useCreateClassMutation();
    const [updateClass, { isLoading: isUpdating }] = useUpdateClassMutation();
    const [deleteClass, { isLoading: isDeleting }] = useDeleteClassMutation();

    const classes = data?.data?.classes || [];
    const teacherMap = new Map((teachersData?.data?.teachers || []).map((teacher) => [teacher.id, teacher]));

    const handleOpenMenu = (event: MouseEvent<HTMLElement>, classData: Class) => {
        setAnchorEl(event.currentTarget);
        setMenuClass(classData);
    };

    const handleCloseMenu = () => {
        setAnchorEl(null);
        setMenuClass(null);
    };

    const handleAdd = () => {
        setSelectedClass(null);
        setIsDialogOpen(true);
    };

    const handleViewDetails = (classData?: Class | null) => {
        const target = classData || menuClass;
        if (!target) {
            return;
        }

        setSelectedClassId(target.id);
        handleCloseMenu();
    };

    const handleEdit = (classData?: Class | null) => {
        const target = classData || menuClass;
        if (!target) {
            return;
        }

        setSelectedClass(target);
        setIsDialogOpen(true);
        handleCloseMenu();
    };

    const handleDeleteClick = () => {
        if (!menuClass) {
            return;
        }

        setClassToDelete(menuClass);
        setIsConfirmOpen(true);
        handleCloseMenu();
    };

    const handleConfirmDelete = async () => {
        if (!classToDelete) {
            return;
        }

        try {
            await deleteClass(classToDelete.id).unwrap();
            toast.success('Xóa lớp học thành công');
            setClassToDelete(null);
        } catch (error) {
            toast.error(getErrorMessage(error, 'Xóa lớp học thất bại'));
        }
    };

    const handleFormSubmit = async (formData: Partial<Class>) => {
        try {
            const payload = {
                ...formData,
                start_date: formData.start_date ? new Date(formData.start_date).toISOString() : undefined,
                end_date: formData.end_date ? new Date(formData.end_date).toISOString() : undefined,
            };

            if (selectedClass) {
                await updateClass({ id: selectedClass.id, body: payload }).unwrap();
                toast.success('Cập nhật lớp học thành công');
            } else {
                await createClass(payload).unwrap();
                toast.success('Tạo lớp học thành công');
            }
        } catch (error) {
            toast.error(getErrorMessage(error, 'Không thể lưu thông tin lớp học'));
            throw error;
        }
    };

    const columns: GridColDef<Class>[] = [
        {
            field: 'code',
            headerName: 'Mã lớp',
            width: 140,
            align: 'left',
            headerAlign: 'left',
            renderCell: (params: GridRenderCellParams<Class>) => (
                <Typography variant="body2" sx={{ fontWeight: 700, color: 'primary.main' }}>
                    {params.row.code}
                </Typography>
            ),
        },
        {
            field: 'name',
            headerName: 'Tên lớp học',
            minWidth: 250,
            flex: 1.5,
            align: 'left',
            headerAlign: 'left',
            renderCell: (params: GridRenderCellParams<Class>) => (
                <Stack direction="row" spacing={1} alignItems="center" sx={{ minWidth: 0, height: '100%' }}>
                    <SchoolOutlined sx={{ color: 'primary.main', fontSize: 18 }} />
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
                            {params.row.notes || 'Chưa có ghi chú'}
                        </Typography>
                    </Box>
                </Stack>
            ),
        },
        {
            field: 'teacher_id',
            headerName: 'Giáo viên phụ trách',
            width: 220,
            align: 'left',
            headerAlign: 'left',
            renderCell: (params: GridRenderCellParams<Class>) => {
                const teacher = params.row.teacher_id ? teacherMap.get(params.row.teacher_id) : null;
                return (
                    <Typography variant="body2" color={teacher ? 'text.primary' : 'text.secondary'}>
                        {teacher ? teacher.full_name : 'Chưa gán giáo viên'}
                    </Typography>
                );
            },
        },
        {
            field: 'start_date',
            headerName: 'Bắt đầu',
            width: 120,
            align: 'center',
            headerAlign: 'center',
            renderCell: (params: GridRenderCellParams<Class>) => (
                <Typography variant="body2">{format(new Date(params.row.start_date), 'dd/MM/yyyy')}</Typography>
            ),
        },
        {
            field: 'max_students',
            headerName: 'Sĩ số tối đa',
            width: 120,
            align: 'center',
            headerAlign: 'center',
        },
        {
            field: 'status',
            headerName: 'Trạng thái',
            width: 140,
            align: 'center',
            headerAlign: 'center',
            renderCell: (params: GridRenderCellParams<Class>) => {
                const meta = statusMeta[params.row.status] || { label: params.row.status, color: 'default' as const };
                return <Chip size="small" label={meta.label} color={meta.color} variant="outlined" />;
            },
        },
        {
            field: 'actions',
            headerName: '',
            width: 64,
            sortable: false,
            align: 'center',
            renderCell: (params: GridRenderCellParams<Class>) => (
                <IconButton size="small" onClick={(event) => handleOpenMenu(event, params.row)}>
                    <MoreVert />
                </IconButton>
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
            <SchoolOutlined sx={{ fontSize: 48, color: 'text.secondary', mb: 2 }} />
            <Typography variant="h6" sx={{ fontWeight: 700, mb: 1 }}>
                Chưa có lớp học nào phù hợp bộ lọc
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                Thử đổi từ khóa tìm kiếm hoặc tạo lớp học mới để bắt đầu quản lý học viên.
            </Typography>
            <Button variant="contained" startIcon={<AddRounded />} onClick={handleAdd}>
                Tạo lớp học đầu tiên
            </Button>
        </Paper>
    );

    return (
        <Box>
            <PageHeader
                title="Quản lý lớp học"
                breadcrumbs={[
                    { label: 'Tổng quan', path: '/app/admin/overview' },
                    { label: 'Quản lý lớp học' },
                ]}
                actions={
                    <Stack direction="row" spacing={1}>
                        <Button
                            variant="outlined"
                            startIcon={<RefreshRounded />}
                            onClick={() => void refetch()}
                            disabled={isFetching}
                        >
                            Làm mới
                        </Button>
                        <Button variant="contained" startIcon={<AddRounded />} onClick={handleAdd}>
                            Mở lớp mới
                        </Button>
                    </Stack>
                }
            />

            {isLoading && classes.length === 0 ? renderLoadingState() : null}
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
                    Không thể tải danh sách lớp học. Kiểm tra backend rồi thử lại.
                </Alert>
            ) : null}

            {!isLoading && !isError ? (
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
                                placeholder="Tìm theo mã lớp hoặc tên lớp"
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
                                <MenuItem value="OPEN">Đang mở</MenuItem>
                                <MenuItem value="CLOSED">Đã đóng</MenuItem>
                                <MenuItem value="CANCELLED">Đã hủy</MenuItem>
                            </TextField>
                        </Stack>

                        <Chip
                            label={`Hiển thị ${classes.length}/${data?.data?.pagination?.total_items || classes.length} lớp`}
                            color="primary"
                            variant="outlined"
                        />
                    </Stack>

                    {classes.length === 0 ? (
                        renderEmptyState()
                    ) : (
                        <DataGrid
                            rows={classes}
                            columns={columns}
                            loading={isFetching}
                            paginationMode="server"
                            rowCount={data?.data?.pagination?.total_items || 0}
                            paginationModel={{ page, pageSize }}
                            onPaginationModelChange={(model) => {
                                setPage(model.page);
                                setPageSize(model.pageSize);
                            }}
                            onRowClick={(params) => handleViewDetails(params.row)}
                            pageSizeOptions={[10, 25, 50]}
                            disableRowSelectionOnClick
                            autoHeight
                            getRowId={(row) => row.id}
                            localeText={{ noRowsLabel: 'Không có dữ liệu lớp học' }}
                            sx={{
                                border: 'none',
                                '& .MuiDataGrid-columnHeaders': {
                                    backgroundColor: '#f8fafc',
                                },
                                '& .MuiDataGrid-row': {
                                    cursor: 'pointer',
                                },
                            }}
                        />
                    )}
                </Paper>
            ) : null}

            <ClassDialog
                open={isDialogOpen}
                onClose={() => setIsDialogOpen(false)}
                onSubmit={handleFormSubmit}
                classData={selectedClass}
                isLoading={isCreating || isUpdating}
            />

            <ClassDetailDialog
                open={!!selectedClassId}
                classId={selectedClassId}
                onClose={() => setSelectedClassId(null)}
                onEdit={(classData) => {
                    setSelectedClassId(null);
                    handleEdit(classData);
                }}
            />

            <ConfirmDialog
                open={isConfirmOpen}
                title="Xóa lớp học"
                message={
                    classToDelete
                        ? `Bạn có chắc chắn muốn xóa lớp "${classToDelete.name}"? Hành động này không thể hoàn tác.`
                        : ''
                }
                confirmText="Xóa lớp"
                isDanger
                onClose={() => {
                    setIsConfirmOpen(false);
                    setClassToDelete(null);
                }}
                onConfirm={() => void handleConfirmDelete()}
                loading={isDeleting}
            />

            <Menu anchorEl={anchorEl} open={Boolean(anchorEl)} onClose={handleCloseMenu}>
                <MenuItem onClick={() => handleViewDetails()}>
                    <VisibilityOutlined fontSize="small" sx={{ mr: 1 }} /> Xem chi tiết
                </MenuItem>
                <MenuItem onClick={() => handleEdit()}>
                    <EditOutlined fontSize="small" sx={{ mr: 1 }} /> Chỉnh sửa
                </MenuItem>
                <MenuItem onClick={handleDeleteClick} sx={{ color: 'error.main' }}>
                    <DeleteOutlineRounded fontSize="small" sx={{ mr: 1 }} /> Xóa
                </MenuItem>
            </Menu>
        </Box>
    );
};
