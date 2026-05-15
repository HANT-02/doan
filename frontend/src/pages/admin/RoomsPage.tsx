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
    ApartmentRounded,
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
    useCreateRoomMutation,
    useDeleteRoomMutation,
    useGetRoomsQuery,
    useUpdateRoomMutation,
    type Room,
    type RoomStatus,
} from '@/api/roomApi';
import RoomDialog from '@/components/admin/RoomDialog';
import ConfirmDialog from '@/components/common/ConfirmDialog';
import PageHeader from '@/components/common/PageHeader';
import { getApiErrorMessage } from '@/utils/apiError';

type RoomFormData = {
    name: string;
    capacity: number;
    location?: string;
    status: RoomStatus;
};

export const RoomsPage = () => {
    const [page, setPage] = useState(0);
    const [pageSize, setPageSize] = useState(10);
    const [search, setSearch] = useState('');

    const [isDialogOpen, setIsDialogOpen] = useState(false);
    const [selectedRoom, setSelectedRoom] = useState<Room | null>(null);
    const [roomToDelete, setRoomToDelete] = useState<Room | null>(null);
    const [isConfirmOpen, setIsConfirmOpen] = useState(false);

    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    const [menuRoom, setMenuRoom] = useState<Room | null>(null);

    const queryParams = useMemo(
        () => ({
            page: page + 1,
            limit: pageSize,
            search: search || undefined,
        }),
        [page, pageSize, search],
    );

    const {
        data,
        isLoading,
        isFetching,
        isError,
        refetch,
    } = useGetRoomsQuery(queryParams);

    const [createRoom, { isLoading: isCreating }] = useCreateRoomMutation();
    const [updateRoom, { isLoading: isUpdating }] = useUpdateRoomMutation();
    const [deleteRoom, { isLoading: isDeleting }] = useDeleteRoomMutation();

    const rooms = data?.data?.rooms || [];

    const handleOpenMenu = (event: MouseEvent<HTMLElement>, room: Room) => {
        setAnchorEl(event.currentTarget);
        setMenuRoom(room);
    };

    const handleCloseMenu = () => {
        setAnchorEl(null);
        setMenuRoom(null);
    };

    const handleAdd = () => {
        setSelectedRoom(null);
        setIsDialogOpen(true);
    };

    const handleEdit = (room?: Room | null) => {
        const target = room || menuRoom;
        if (!target) {
            return;
        }

        setSelectedRoom(target);
        setIsDialogOpen(true);
        handleCloseMenu();
    };

    const handleDeleteClick = () => {
        if (!menuRoom) {
            return;
        }

        setRoomToDelete(menuRoom);
        setIsConfirmOpen(true);
        handleCloseMenu();
    };

    const handleConfirmDelete = async () => {
        if (!roomToDelete) {
            return;
        }

        try {
            await deleteRoom(roomToDelete.id).unwrap();
            toast.success('Xóa phòng học thành công');
            setRoomToDelete(null);
        } catch (error) {
            toast.error(getApiErrorMessage(error, 'Xóa phòng học thất bại'));
        }
    };

    const handleFormSubmit = async (formData: RoomFormData) => {
        try {
            if (selectedRoom) {
                await updateRoom({ id: selectedRoom.id, body: formData }).unwrap();
                toast.success('Cập nhật phòng học thành công');
            } else {
                await createRoom(formData).unwrap();
                toast.success('Tạo phòng học thành công');
            }
        } catch (error) {
            toast.error(getApiErrorMessage(error, 'Không thể lưu thông tin phòng học'));
            throw error;
        }
    };

    const columns: GridColDef<Room>[] = [
        {
            field: 'code',
            headerName: 'Mã phòng',
            width: 140,
            align: 'left',
            headerAlign: 'left',
            renderCell: (params: GridRenderCellParams<Room>) => (
                <Typography variant="body2" sx={{ fontWeight: 700, color: 'primary.main' }}>
                    {params.row.code || 'Chưa có mã'}
                </Typography>
            ),
        },
        {
            field: 'name',
            headerName: 'Tên phòng',
            minWidth: 260,
            flex: 1.3,
            align: 'left',
            headerAlign: 'left',
            renderCell: (params: GridRenderCellParams<Room>) => (
                <Stack direction="row" spacing={1} alignItems="center" sx={{ minWidth: 0, height: '100%' }}>
                    <ApartmentRounded sx={{ color: 'primary.main', fontSize: 18 }} />
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
                            {params.row.address || params.row.location || 'Chưa cập nhật địa chỉ'}
                        </Typography>
                        <Typography
                            variant="caption"
                            color="text.secondary"
                            noWrap
                            sx={{ display: 'block', lineHeight: 1.4 }}
                        >
                            {params.row.campus?.name || 'Chưa gán campus'}
                        </Typography>
                    </Box>
                </Stack>
            ),
        },
        {
            field: 'capacity',
            headerName: 'Sức chứa',
            width: 120,
            align: 'center',
            headerAlign: 'center',
            renderCell: (params: GridRenderCellParams<Room>) => (
                <Chip size="small" label={`${params.row.capacity} chỗ`} color="primary" variant="outlined" />
            ),
        },
        {
            field: 'address',
            headerName: 'Địa chỉ',
            minWidth: 240,
            flex: 1,
            align: 'left',
            headerAlign: 'left',
            renderCell: (params: GridRenderCellParams<Room>) => (
                <Typography variant="body2" color="text.secondary" noWrap>
                    {params.row.address || params.row.location || 'Chưa cập nhật'}
                </Typography>
            ),
        },
        {
            field: 'updated_at',
            headerName: 'Cập nhật',
            width: 130,
            align: 'center',
            headerAlign: 'center',
            renderCell: (params: GridRenderCellParams<Room>) => (
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
            renderCell: (params: GridRenderCellParams<Room>) => (
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
            <ApartmentRounded sx={{ fontSize: 48, color: 'text.secondary', mb: 2 }} />
            <Typography variant="h6" sx={{ fontWeight: 700, mb: 1 }}>
                Chưa có phòng học nào phù hợp
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                Kiểm tra lại bộ lọc hoặc tạo phòng học mới để bắt đầu xếp lớp và xếp lịch.
            </Typography>
            <Button variant="contained" startIcon={<AddRounded />} onClick={handleAdd}>
                Tạo phòng học đầu tiên
            </Button>
        </Paper>
    );

    return (
        <Box>
            <PageHeader
                title="Quản lý phòng học"
                subtitle="Kết nối trực tiếp dữ liệu phòng học từ backend, sẵn sàng cho luồng lớp học và xếp lịch."
                breadcrumbs={[
                    { label: 'Tổng quan', path: '/app/admin/overview' },
                    { label: 'Quản lý phòng học' },
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
                            Thêm phòng mới
                        </Button>
                    </Stack>
                }
            />

            {isLoading && rooms.length === 0 ? renderLoadingState() : null}
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
                    Không thể tải danh sách phòng học. Kiểm tra backend rồi thử lại.
                </Alert>
            ) : null}

            {!isLoading && !isError ? (
                rooms.length === 0 ? (
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
                            <TextField
                                value={search}
                                onChange={(event) => {
                                    setSearch(event.target.value);
                                    setPage(0);
                                }}
                                size="small"
                                placeholder="Tìm theo tên phòng"
                                sx={{ minWidth: { xs: '100%', md: 280 } }}
                                InputProps={{
                                    startAdornment: (
                                        <InputAdornment position="start">
                                            <SearchRounded fontSize="small" />
                                        </InputAdornment>
                                    ),
                                }}
                            />

                            <Chip
                                label={`Hiển thị ${rooms.length}/${data?.data?.pagination?.total_items || rooms.length} phòng`}
                                color="primary"
                                variant="outlined"
                            />
                        </Stack>

                        <DataGrid
                            rows={rooms}
                            columns={columns}
                            loading={isFetching}
                            paginationMode="server"
                            rowCount={data?.data?.pagination?.total_items || 0}
                            paginationModel={{ page, pageSize }}
                            onPaginationModelChange={(model) => {
                                setPage(model.page);
                                setPageSize(model.pageSize);
                            }}
                            pageSizeOptions={[10, 25, 50]}
                            disableRowSelectionOnClick
                            autoHeight
                            getRowId={(row) => row.id}
                            localeText={{ noRowsLabel: 'Không có dữ liệu phòng học' }}
                            sx={{
                                border: 'none',
                                '& .MuiDataGrid-columnHeaders': {
                                    backgroundColor: '#f8fafc',
                                },
                            }}
                        />
                    </Paper>
                )
            ) : null}

            <RoomDialog
                open={isDialogOpen}
                onClose={() => setIsDialogOpen(false)}
                onSubmit={handleFormSubmit}
                room={selectedRoom}
                isLoading={isCreating || isUpdating}
            />

            <ConfirmDialog
                open={isConfirmOpen}
                title="Xóa phòng học"
                message={
                    roomToDelete
                        ? `Bạn có chắc chắn muốn xóa phòng "${roomToDelete.name}"? Hành động này không thể hoàn tác.`
                        : ''
                }
                confirmText="Xóa phòng"
                isDanger
                onClose={() => {
                    setIsConfirmOpen(false);
                    setRoomToDelete(null);
                }}
                onConfirm={() => void handleConfirmDelete()}
                loading={isDeleting}
            />

            <Menu anchorEl={anchorEl} open={Boolean(anchorEl)} onClose={handleCloseMenu}>
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
