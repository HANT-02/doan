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
    AccessTimeRounded,
    AddRounded,
    DeleteOutlineRounded,
    EditOutlined,
    MoreVert,
    RefreshRounded,
    SearchRounded,
} from '@mui/icons-material';
import { DataGrid, type GridColDef, type GridRenderCellParams } from '@mui/x-data-grid';
import { format } from 'date-fns';
import { toast } from 'sonner';

import {
    useCreateShiftMutation,
    useDeleteShiftMutation,
    useGetShiftsQuery,
    useUpdateShiftMutation,
    type Shift,
    type ShiftSessionType,
} from '@/api/shiftApi';
import ShiftDialog from '@/components/admin/ShiftDialog';
import ConfirmDialog from '@/components/common/ConfirmDialog';
import PageHeader from '@/components/common/PageHeader';
import { getApiErrorMessage } from '@/utils/apiError';

type ShiftFormData = {
    code: string;
    name: string;
    start_time: string;
    end_time: string;
    duration_minutes: number;
    session_type: ShiftSessionType;
    is_active: boolean;
    notes?: string;
};

const shiftTypeLabel: Record<ShiftSessionType, string> = {
    MORNING: 'Buổi sáng',
    AFTERNOON: 'Buổi chiều',
    EVENING: 'Buổi tối',
    CUSTOM: 'Tùy chỉnh',
};

export const ShiftsPage = () => {
    const [page, setPage] = useState(0);
    const [pageSize, setPageSize] = useState(10);
    const [search, setSearch] = useState('');

    const [isDialogOpen, setIsDialogOpen] = useState(false);
    const [selectedShift, setSelectedShift] = useState<Shift | null>(null);
    const [shiftToDelete, setShiftToDelete] = useState<Shift | null>(null);
    const [isConfirmOpen, setIsConfirmOpen] = useState(false);

    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    const [menuShift, setMenuShift] = useState<Shift | null>(null);

    const queryParams = useMemo(
        () => ({
            page: page + 1,
            limit: pageSize,
            search: search || undefined,
        }),
        [page, pageSize, search],
    );

    const { data, isLoading, isFetching, isError, refetch } = useGetShiftsQuery(queryParams);
    const [createShift, { isLoading: isCreating }] = useCreateShiftMutation();
    const [updateShift, { isLoading: isUpdating }] = useUpdateShiftMutation();
    const [deleteShift, { isLoading: isDeleting }] = useDeleteShiftMutation();

    const shifts = data?.data?.shifts || [];

    const handleOpenMenu = (event: MouseEvent<HTMLElement>, shift: Shift) => {
        setAnchorEl(event.currentTarget);
        setMenuShift(shift);
    };

    const handleCloseMenu = () => {
        setAnchorEl(null);
        setMenuShift(null);
    };

    const handleAdd = () => {
        setSelectedShift(null);
        setIsDialogOpen(true);
    };

    const handleEdit = (shift?: Shift | null) => {
        const target = shift || menuShift;
        if (!target) {
            return;
        }

        setSelectedShift(target);
        setIsDialogOpen(true);
        handleCloseMenu();
    };

    const handleDeleteClick = () => {
        if (!menuShift) {
            return;
        }

        setShiftToDelete(menuShift);
        setIsConfirmOpen(true);
        handleCloseMenu();
    };

    const handleConfirmDelete = async () => {
        if (!shiftToDelete) {
            return;
        }

        try {
            await deleteShift(shiftToDelete.id).unwrap();
            toast.success('Xóa ca học thành công');
            setShiftToDelete(null);
        } catch (error) {
            toast.error(getApiErrorMessage(error, 'Xóa ca học thất bại'));
        }
    };

    const handleFormSubmit = async (formData: ShiftFormData) => {
        try {
            if (selectedShift) {
                await updateShift({ id: selectedShift.id, body: formData }).unwrap();
                toast.success('Cập nhật ca học thành công');
            } else {
                await createShift(formData).unwrap();
                toast.success('Tạo ca học thành công');
            }
        } catch (error) {
            toast.error(getApiErrorMessage(error, 'Không thể lưu thông tin ca học'));
            throw error;
        }
    };

    const columns: GridColDef<Shift>[] = [
        {
            field: 'code',
            headerName: 'Mã ca',
            width: 120,
            renderCell: (params: GridRenderCellParams<Shift>) => (
                <Typography variant="body2" sx={{ fontWeight: 700, color: 'primary.main' }}>
                    {params.row.code}
                </Typography>
            ),
        },
        {
            field: 'name',
            headerName: 'Tên ca học',
            minWidth: 260,
            flex: 1.2,
            renderCell: (params: GridRenderCellParams<Shift>) => (
                <Stack direction="row" spacing={1} alignItems="center" sx={{ minWidth: 0, height: '100%' }}>
                    <AccessTimeRounded sx={{ color: 'primary.main', fontSize: 18 }} />
                    <Box sx={{ minWidth: 0, display: 'flex', flexDirection: 'column', justifyContent: 'center', py: 1 }}>
                        <Typography variant="body2" sx={{ fontWeight: 700, lineHeight: 1.35 }} noWrap>
                            {params.row.name}
                        </Typography>
                        <Typography variant="caption" color="text.secondary" noWrap sx={{ display: 'block', lineHeight: 1.4, mt: 0.25 }}>
                            {params.row.notes || 'Ca học chuẩn dùng cho scheduling và lịch mẫu lớp học'}
                        </Typography>
                    </Box>
                </Stack>
            ),
        },
        {
            field: 'time_range',
            headerName: 'Khung giờ',
            width: 160,
            sortable: false,
            renderCell: (params: GridRenderCellParams<Shift>) => (
                <Typography variant="body2">
                    {params.row.start_time} - {params.row.end_time}
                </Typography>
            ),
        },
        {
            field: 'duration_minutes',
            headerName: 'Thời lượng',
            width: 120,
            align: 'center',
            headerAlign: 'center',
            renderCell: (params: GridRenderCellParams<Shift>) => (
                <Chip size="small" label={`${params.row.duration_minutes} phút`} color="primary" variant="outlined" />
            ),
        },
        {
            field: 'session_type',
            headerName: 'Loại ca',
            width: 140,
            renderCell: (params: GridRenderCellParams<Shift>) => (
                <Typography variant="body2">{shiftTypeLabel[params.row.session_type]}</Typography>
            ),
        },
        {
            field: 'is_active',
            headerName: 'Trạng thái',
            width: 130,
            renderCell: (params: GridRenderCellParams<Shift>) => (
                <Chip
                    size="small"
                    label={params.row.is_active ? 'Đang hoạt động' : 'Tạm ngưng'}
                    color={params.row.is_active ? 'success' : 'default'}
                    variant={params.row.is_active ? 'filled' : 'outlined'}
                />
            ),
        },
        {
            field: 'updated_at',
            headerName: 'Cập nhật',
            width: 130,
            align: 'center',
            headerAlign: 'center',
            renderCell: (params: GridRenderCellParams<Shift>) => (
                <Typography variant="body2">
                    {params.row.updated_at ? format(new Date(params.row.updated_at), 'dd/MM/yyyy') : '-'}
                </Typography>
            ),
        },
        {
            field: 'actions',
            headerName: '',
            width: 64,
            sortable: false,
            align: 'center',
            renderCell: (params: GridRenderCellParams<Shift>) => (
                <IconButton size="small" onClick={(event) => handleOpenMenu(event, params.row)}>
                    <MoreVert />
                </IconButton>
            ),
        },
    ];

    return (
        <Box>
            <PageHeader
                title="Quản lý ca học"
                subtitle="Chuẩn hóa các ca học làm đầu vào thời gian cho lịch mẫu lớp học và benchmark scheduling."
                breadcrumbs={[
                    { label: 'Tổng quan', path: '/app/admin/overview' },
                    { label: 'Quản lý ca học' },
                ]}
                actions={
                    <Stack direction="row" spacing={1}>
                        <Button variant="outlined" startIcon={<RefreshRounded />} onClick={() => void refetch()} disabled={isFetching}>
                            Làm mới
                        </Button>
                        <Button variant="contained" startIcon={<AddRounded />} onClick={handleAdd}>
                            Thêm ca học
                        </Button>
                    </Stack>
                }
            />

            {isLoading && shifts.length === 0 ? (
                <Stack spacing={1.5}>
                    <Skeleton variant="rounded" height={72} />
                    <Skeleton variant="rounded" height={420} />
                </Stack>
            ) : null}

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
                    Không thể tải danh sách ca học. Kiểm tra backend rồi thử lại.
                </Alert>
            ) : null}

            {!isLoading && !isError ? (
                <Paper
                    variant="outlined"
                    sx={{
                        borderRadius: 4,
                        overflow: 'hidden',
                        p: 2,
                    }}
                >
                    <Stack spacing={2}>
                        <TextField
                            value={search}
                            onChange={(event) => {
                                setSearch(event.target.value);
                                setPage(0);
                            }}
                            placeholder="Tìm theo mã ca hoặc tên ca học"
                            InputProps={{
                                startAdornment: (
                                    <InputAdornment position="start">
                                        <SearchRounded />
                                    </InputAdornment>
                                ),
                            }}
                        />

                        <DataGrid
                            autoHeight
                            disableRowSelectionOnClick
                            rows={shifts}
                            columns={columns}
                            getRowId={(row) => row.id}
                            rowCount={data?.data?.pagination.total_items || 0}
                            paginationMode="server"
                            paginationModel={{ page, pageSize }}
                            onPaginationModelChange={(model) => {
                                setPage(model.page);
                                setPageSize(model.pageSize);
                            }}
                            loading={isFetching}
                            pageSizeOptions={[10, 20, 50]}
                            sx={{
                                border: 0,
                                '& .MuiDataGrid-columnHeaders': {
                                    backgroundColor: 'background.default',
                                },
                            }}
                        />
                    </Stack>
                </Paper>
            ) : null}

            <ShiftDialog
                open={isDialogOpen}
                onClose={() => setIsDialogOpen(false)}
                onSubmit={handleFormSubmit}
                shift={selectedShift}
                isLoading={isCreating || isUpdating}
            />

            <ConfirmDialog
                open={isConfirmOpen}
                title="Xóa ca học"
                message={`Bạn có chắc muốn xóa ca học "${shiftToDelete?.name || ''}" không?`}
                confirmText="Xóa"
                cancelText="Hủy"
                isDanger
                loading={isDeleting}
                onClose={() => setIsConfirmOpen(false)}
                onConfirm={() => void handleConfirmDelete()}
            />

            <Menu anchorEl={anchorEl} open={Boolean(anchorEl)} onClose={handleCloseMenu}>
                <MenuItem onClick={() => handleEdit()}>
                    <EditOutlined fontSize="small" sx={{ mr: 1 }} />
                    Chỉnh sửa
                </MenuItem>
                <MenuItem onClick={handleDeleteClick} sx={{ color: 'error.main' }}>
                    <DeleteOutlineRounded fontSize="small" sx={{ mr: 1 }} />
                    Xóa
                </MenuItem>
            </Menu>
        </Box>
    );
};
