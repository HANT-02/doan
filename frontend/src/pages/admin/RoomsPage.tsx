import { useState } from 'react';
import {
    Box,
    Typography,
    Button,
    Paper,
    Breadcrumbs,
    Link,
    Chip,
    IconButton,
} from '@mui/material';
import { Add, Room as RoomIcon, MoreVert, Edit, Delete } from '@mui/icons-material';
import { DataGrid } from '@mui/x-data-grid';
import type { GridColDef, GridRenderCellParams } from '@mui/x-data-grid';
import {
    useGetRoomsQuery,
    useCreateRoomMutation,
    useUpdateRoomMutation,
    useDeleteRoomMutation
} from '@/api/roomApi';
import PageHeader from '@/components/common/PageHeader';
import RoomDialog from '@/components/admin/RoomDialog';
import ConfirmDialog from '@/components/common/ConfirmDialog';
import { Menu, MenuItem as MuiMenuItem } from '@mui/material';

export const RoomsPage = () => {
    const [page, setPage] = useState(0);
    const [pageSize, setPageSize] = useState(10);

    // Dialog states
    const [isDialogOpen, setIsDialogOpen] = useState(false);
    const [selectedRoom, setSelectedRoom] = useState<any>(null);
    const [isConfirmOpen, setIsConfirmOpen] = useState(false);
    const [roomToDelete, setRoomToDelete] = useState<string | null>(null);

    // Menu state
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    const [menuRoom, setMenuRoom] = useState<any>(null);

    const { data, isLoading } = useGetRoomsQuery({
        page: page + 1,
        limit: pageSize,
    });

    const [createRoom, { isLoading: isCreating }] = useCreateRoomMutation();
    const [updateRoom, { isLoading: isUpdating }] = useUpdateRoomMutation();
    const [deleteRoom, { isLoading: isDeleting }] = useDeleteRoomMutation();

    const handleOpenMenu = (event: React.MouseEvent<HTMLElement>, room: any) => {
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

    const handleEdit = () => {
        setSelectedRoom(menuRoom);
        setIsDialogOpen(true);
        handleCloseMenu();
    };

    const handleDeleteClick = () => {
        setRoomToDelete(menuRoom.id);
        setIsConfirmOpen(true);
        handleCloseMenu();
    };

    const handleConfirmDelete = async () => {
        if (roomToDelete) {
            try {
                await deleteRoom(roomToDelete).unwrap();
            } catch (error) {
                console.error('Failed to delete room:', error);
            } finally {
                setIsConfirmOpen(false);
                setRoomToDelete(null);
            }
        }
    };

    const handleFormSubmit = async (formData: any) => {
        try {
            if (selectedRoom) {
                await updateRoom({ id: selectedRoom.id, body: formData }).unwrap();
            } else {
                await createRoom(formData).unwrap();
            }
        } catch (error) {
            console.error('Failed to save room:', error);
        }
    };

    const columns: GridColDef[] = [
        {
            field: 'name',
            headerName: 'Tên phòng',
            flex: 1,
            renderCell: (params: GridRenderCellParams) => (
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <RoomIcon color="primary" fontSize="small" />
                    <Typography variant="body2" sx={{ fontWeight: 600 }}>{params.value}</Typography>
                </Box>
            )
        },
        { field: 'capacity', headerName: 'Sức chứa', width: 120 },
        { field: 'location', headerName: 'Vị trí', flex: 1 },
        {
            field: 'status',
            headerName: 'Trạng thái',
            width: 150,
            renderCell: (params: GridRenderCellParams) => {
                const status = params.value as string;
                let color: 'success' | 'warning' | 'error' = 'success';
                if (status === 'MAINTENANCE') color = 'warning';
                if (status === 'INACTIVE') color = 'error';

                return <Chip label={status} size="small" color={color} sx={{ fontWeight: 600 }} />;
            }
        },
        {
            field: 'actions',
            headerName: '',
            width: 50,
            sortable: false,
            renderCell: (params: GridRenderCellParams) => (
                <IconButton size="small" onClick={(e) => handleOpenMenu(e, params.row)}>
                    <MoreVert />
                </IconButton>
            )
        }
    ];

    return (
        <Box>
            <PageHeader
                title="Quản lý Phòng học"
                subtitle="Danh sách các phòng học và cơ sở vật chất"
                actions={
                    <Button
                        variant="contained"
                        startIcon={<Add />}
                        sx={{ borderRadius: 2 }}
                        onClick={handleAdd}
                    >
                        Thêm phòng mới
                    </Button>
                }
            />

            <Breadcrumbs sx={{ mb: 3 }}>
                <Link underline="hover" color="inherit" href="/app/admin/overview">Dashboard</Link>
                <Typography color="text.primary">Phòng học</Typography>
            </Breadcrumbs>

            <Paper elevation={0} sx={{ p: 2, borderRadius: 3, border: '1px solid #e2e8f0' }}>
                <DataGrid
                    rows={data?.data?.rooms || []}
                    columns={columns}
                    loading={isLoading}
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
                    getRowId={(row) => row.id || row.key || Math.random().toString()}
                    sx={{
                        border: 'none',
                        '& .MuiDataGrid-columnHeaders': {
                            backgroundColor: '#f8fafc',
                            borderRadius: 0,
                        },
                        '& .MuiDataGrid-cell:focus': {
                            outline: 'none',
                        },
                    }}
                />
            </Paper>
            <RoomDialog
                open={isDialogOpen}
                onClose={() => setIsDialogOpen(false)}
                onSubmit={handleFormSubmit}
                room={selectedRoom}
                isLoading={isCreating || isUpdating}
            />

            <ConfirmDialog
                open={isConfirmOpen}
                title="Xác nhận xóa"
                message={`Bạn có chắc chắn muốn xóa phòng học "${menuRoom?.name || ''}"? Hành động này không thể hoàn tác.`}
                onClose={() => setIsConfirmOpen(false)}
                onConfirm={handleConfirmDelete}
                loading={isDeleting}
            />

            <Menu
                anchorEl={anchorEl}
                open={Boolean(anchorEl)}
                onClose={handleCloseMenu}
            >
                <MuiMenuItem onClick={handleEdit}>
                    <Edit fontSize="small" sx={{ mr: 1 }} /> Chỉnh sửa
                </MuiMenuItem>
                <MuiMenuItem onClick={handleDeleteClick} sx={{ color: 'error.main' }}>
                    <Delete fontSize="small" sx={{ mr: 1 }} /> Xóa
                </MuiMenuItem>
            </Menu>
        </Box>
    );
};
