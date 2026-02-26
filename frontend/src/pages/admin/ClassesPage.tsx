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
import { Add, School, MoreVert, CalendarToday, Edit, Delete } from '@mui/icons-material';
import { DataGrid } from '@mui/x-data-grid';
import type { GridColDef, GridRenderCellParams } from '@mui/x-data-grid';
import {
    useGetClassesQuery,
    useCreateClassMutation,
    useUpdateClassMutation,
    useDeleteClassMutation
} from '@/api/classApi';
import PageHeader from '@/components/common/PageHeader';
import ClassDialog from '@/components/admin/ClassDialog';
import ConfirmDialog from '@/components/common/ConfirmDialog';
import { Menu, MenuItem as MuiMenuItem } from '@mui/material';
import { format } from 'date-fns';

export const ClassesPage = () => {
    const [page, setPage] = useState(0);
    const [pageSize, setPageSize] = useState(10);

    // Dialog states
    const [isDialogOpen, setIsDialogOpen] = useState(false);
    const [selectedClass, setSelectedClass] = useState<any>(null);
    const [isConfirmOpen, setIsConfirmOpen] = useState(false);
    const [classToDelete, setClassToDelete] = useState<string | null>(null);

    // Menu state
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    const [menuClass, setMenuClass] = useState<any>(null);

    const { data, isLoading } = useGetClassesQuery({
        page: page + 1,
        limit: pageSize,
    });

    const [createClass, { isLoading: isCreating }] = useCreateClassMutation();
    const [updateClass, { isLoading: isUpdating }] = useUpdateClassMutation();
    const [deleteClass, { isLoading: isDeleting }] = useDeleteClassMutation();

    const handleOpenMenu = (event: React.MouseEvent<HTMLElement>, classData: any) => {
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

    const handleEdit = () => {
        setSelectedClass(menuClass);
        setIsDialogOpen(true);
        handleCloseMenu();
    };

    const handleDeleteClick = () => {
        setClassToDelete(menuClass.id);
        setIsConfirmOpen(true);
        handleCloseMenu();
    };

    const handleConfirmDelete = async () => {
        if (classToDelete) {
            try {
                await deleteClass(classToDelete).unwrap();
            } catch (error) {
                console.error('Failed to delete class:', error);
            } finally {
                setIsConfirmOpen(false);
                setClassToDelete(null);
            }
        }
    };

    const handleFormSubmit = async (formData: any) => {
        try {
            if (selectedClass) {
                await updateClass({ id: selectedClass.id, body: formData }).unwrap();
            } else {
                await createClass(formData).unwrap();
            }
        } catch (error) {
            console.error('Failed to save class:', error);
        }
    };

    const columns: GridColDef[] = [
        {
            field: 'code',
            headerName: 'Mã lớp',
            width: 120,
            renderCell: (params: GridRenderCellParams) => (
                <Typography variant="body2" sx={{ fontWeight: 700, color: 'primary.main' }}>{params.value}</Typography>
            )
        },
        {
            field: 'name',
            headerName: 'Tên lớp học',
            flex: 1.5,
            renderCell: (params: GridRenderCellParams) => (
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <School color="action" fontSize="small" />
                    <Typography variant="body2" sx={{ fontWeight: 500 }}>{params.value}</Typography>
                </Box>
            )
        },
        {
            field: 'start_date',
            headerName: 'Ngày bắt đầu',
            width: 150,
            renderCell: (params: GridRenderCellParams) => (
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                    <CalendarToday sx={{ fontSize: 16, color: 'text.secondary' }} />
                    <Typography variant="body2">{format(new Date(params.value as string), 'dd/MM/yyyy')}</Typography>
                </Box>
            )
        },
        {
            field: 'max_students',
            headerName: 'Sĩ số',
            width: 100,
            align: 'center',
            headerAlign: 'center',
        },
        {
            field: 'status',
            headerName: 'Trạng thái',
            width: 120,
            renderCell: (params: GridRenderCellParams) => {
                const status = params.value as string;
                let color: 'success' | 'warning' | 'error' = 'success';
                if (status === 'CLOSED') color = 'error';
                if (status === 'CANCELLED') color = 'warning';

                return <Chip label={status} size="small" color={color} variant="outlined" sx={{ fontWeight: 600 }} />;
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
                title="Quản lý Lớp học"
                subtitle="Theo dõi và quản lý các lớp học đang diễn ra"
                actions={
                    <Button
                        variant="contained"
                        startIcon={<Add />}
                        sx={{ borderRadius: 2 }}
                        onClick={handleAdd}
                    >
                        Mở lớp mới
                    </Button>
                }
            />

            <Breadcrumbs sx={{ mb: 3 }}>
                <Link underline="hover" color="inherit" href="/app/admin/overview">Dashboard</Link>
                <Typography color="text.primary">Lớp học</Typography>
            </Breadcrumbs>

            <Paper elevation={0} sx={{ p: 2, borderRadius: 3, border: '1px solid #e2e8f0' }}>
                <DataGrid
                    rows={data?.data?.classes || []}
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
            <ClassDialog
                open={isDialogOpen}
                onClose={() => setIsDialogOpen(false)}
                onSubmit={handleFormSubmit}
                classData={selectedClass}
                isLoading={isCreating || isUpdating}
            />

            <ConfirmDialog
                open={isConfirmOpen}
                title="Xác nhận xóa"
                message={`Bạn có chắc chắn muốn xóa lớp học "${menuClass?.name || ''}"? Hành động này không thể hoàn tác.`}
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
