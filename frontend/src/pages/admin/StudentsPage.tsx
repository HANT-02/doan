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
    Avatar,
} from '@mui/material';
import { Add, MoreVert, Edit, Delete } from '@mui/icons-material';
import { DataGrid } from '@mui/x-data-grid';
import type { GridColDef, GridRenderCellParams } from '@mui/x-data-grid';
import {
    useGetStudentsQuery,
    useCreateStudentMutation,
    useUpdateStudentMutation,
    useDeleteStudentMutation
} from '@/api/studentApi';
import PageHeader from '@/components/common/PageHeader';
import StudentDialog from '@/components/admin/StudentDialog';
import ConfirmDialog from '@/components/common/ConfirmDialog';
import { Menu, MenuItem as MuiMenuItem } from '@mui/material';

export const StudentsPage = () => {
    const [page, setPage] = useState(0);
    const [pageSize, setPageSize] = useState(10);

    // Dialog states
    const [isDialogOpen, setIsDialogOpen] = useState(false);
    const [selectedStudent, setSelectedStudent] = useState<any>(null);
    const [isConfirmOpen, setIsConfirmOpen] = useState(false);
    const [studentToDelete, setStudentToDelete] = useState<string | null>(null);

    // Menu state
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    const [menuStudent, setMenuStudent] = useState<any>(null);

    const { data, isLoading } = useGetStudentsQuery({
        page: page + 1,
        limit: pageSize,
    });

    const [createStudent, { isLoading: isCreating }] = useCreateStudentMutation();
    const [updateStudent, { isLoading: isUpdating }] = useUpdateStudentMutation();
    const [deleteStudent, { isLoading: isDeleting }] = useDeleteStudentMutation();

    const handleOpenMenu = (event: React.MouseEvent<HTMLElement>, student: any) => {
        setAnchorEl(event.currentTarget);
        setMenuStudent(student);
    };

    const handleCloseMenu = () => {
        setAnchorEl(null);
        setMenuStudent(null);
    };

    const handleAdd = () => {
        setSelectedStudent(null);
        setIsDialogOpen(true);
    };

    const handleEdit = () => {
        setSelectedStudent(menuStudent);
        setIsDialogOpen(true);
        handleCloseMenu();
    };

    const handleDeleteClick = () => {
        setStudentToDelete(menuStudent.id);
        setIsConfirmOpen(true);
        handleCloseMenu();
    };

    const handleConfirmDelete = async () => {
        if (studentToDelete) {
            try {
                await deleteStudent(studentToDelete).unwrap();
            } catch (error) {
                console.error('Failed to delete student:', error);
            } finally {
                setIsConfirmOpen(false);
                setStudentToDelete(null);
            }
        }
    };

    const handleFormSubmit = async (formData: any) => {
        try {
            if (selectedStudent) {
                await updateStudent({ id: selectedStudent.id, body: formData }).unwrap();
            } else {
                await createStudent(formData).unwrap();
            }
        } catch (error) {
            console.error('Failed to save student:', error);
        }
    };

    const columns: GridColDef[] = [
        {
            field: 'code',
            headerName: 'Mã HS',
            width: 120,
            renderCell: (params: GridRenderCellParams) => (
                <Typography variant="body2" sx={{ fontWeight: 600, color: 'primary.main' }}>
                    {params.value}
                </Typography>
            )
        },
        {
            field: 'full_name',
            headerName: 'Họ và tên',
            flex: 1.5,
            renderCell: (params: GridRenderCellParams) => (
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5 }}>
                    <Avatar sx={{ width: 32, height: 32, bgcolor: 'primary.light', fontSize: '0.875rem' }}>
                        {(params.value as string).charAt(0)}
                    </Avatar>
                    <Typography variant="body2" sx={{ fontWeight: 500 }}>{params.value}</Typography>
                </Box>
            )
        },
        { field: 'grade_level', headerName: 'Khối lớp', width: 100 },
        { field: 'phone', headerName: 'Số điện thoại', width: 130 },
        { field: 'guardian_phone', headerName: 'SĐT Phụ huynh', width: 130 },
        {
            field: 'status',
            headerName: 'Trạng thái',
            width: 130,
            renderCell: (params: GridRenderCellParams) => {
                const status = params.value as string;
                const color = status === 'ACTIVE' ? 'success' : 'error';
                const label = status === 'ACTIVE' ? 'Đang học' : 'Nghỉ học';

                return <Chip label={label} size="small" color={color} sx={{ fontWeight: 600 }} />;
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
                title="Quản lý Học sinh"
                subtitle="Danh sách học sinh đang theo học tại trung tâm"
                actions={
                    <Button
                        variant="contained"
                        startIcon={<Add />}
                        sx={{ borderRadius: 2 }}
                        onClick={handleAdd}
                    >
                        Thêm học sinh
                    </Button>
                }
            />

            <Breadcrumbs sx={{ mb: 3 }}>
                <Link underline="hover" color="inherit" href="/app/admin/overview">Dashboard</Link>
                <Typography color="text.primary">Học sinh</Typography>
            </Breadcrumbs>

            <Paper elevation={0} sx={{ p: 2, borderRadius: 3, border: '1px solid #e2e8f0' }}>
                <DataGrid
                    rows={data?.data?.students || []}
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

            <StudentDialog
                open={isDialogOpen}
                onClose={() => setIsDialogOpen(false)}
                onSubmit={handleFormSubmit}
                student={selectedStudent}
                isLoading={isCreating || isUpdating}
            />

            <ConfirmDialog
                open={isConfirmOpen}
                title="Xác nhận xóa"
                message={`Bạn có chắc chắn muốn xóa học sinh "${menuStudent?.full_name || ''}"? Hành động này không thể hoàn tác.`}
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
