import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
    Container,
    Button,
    Alert,
    Box,
} from '@mui/material';
import { Add } from '@mui/icons-material';
import { useGetTeachersQuery, useDeleteTeacherMutation, type Teacher, type ListTeachersParams } from '@/api/teacherApi';
import { TeacherListTable } from '@/components/teacher/TeacherListTable';
import { TeacherFilters, type TeacherFiltersState } from '@/components/teacher/TeacherFilters';
import { DeleteTeacherDialog } from '@/components/teacher/DeleteTeacherDialog';
import { useAppSelector } from '@/store';
import { toast } from 'sonner';
import PageHeader from '@/components/common/PageHeader';

export const TeachersPage = () => {
    const navigate = useNavigate();
    const { user } = useAppSelector(state => state.auth);
    const isAdmin = user?.role === 'admin' || user?.role === 'super_admin';

    const [paginationModel, setPaginationModel] = useState({
        page: 0, // DataGrid uses 0-based page index
        pageSize: 10,
    });

    const [filters, setFilters] = useState<TeacherFiltersState>({
        search: '',
        status: '',
        employment_type: '',
    });

    const queryParams: ListTeachersParams = {
        page: paginationModel.page + 1, // API is 1-based
        limit: paginationModel.pageSize,
        ...filters,
    };

    const { data: responseData, isLoading, error, refetch } = useGetTeachersQuery(queryParams);
    const [deleteTeacher, { isLoading: isDeleting }] = useDeleteTeacherMutation();

    const [deleteDialog, setDeleteDialog] = useState<{
        open: boolean;
        teacher: Teacher | null;
    }>({
        open: false,
        teacher: null,
    });

    const handleFilterChange = (newFilters: TeacherFiltersState) => {
        setFilters(newFilters);
        setPaginationModel(prev => ({ ...prev, page: 0 }));
    };

    const handleDeleteClick = (teacher: Teacher) => {
        setDeleteDialog({ open: true, teacher });
    };

    const handleDeleteConfirm = async () => {
        if (!deleteDialog.teacher) return;

        try {
            await deleteTeacher(deleteDialog.teacher.id).unwrap();
            toast.success('Xóa giáo viên thành công');
            setDeleteDialog({ open: false, teacher: null });
            refetch();
        } catch (err: any) {
            toast.error(err?.data?.message || 'Xóa giáo viên thất bại');
        }
    };

    return (
        <Container maxWidth="xl" sx={{ py: 4 }}>
            <PageHeader
                title="Danh sách Giáo viên"
                subtitle="Quản lý và theo dõi thông tin đội ngũ giáo viên"
                actions={
                    isAdmin && (
                        <Button
                            variant="contained"
                            startIcon={<Add />}
                            onClick={() => navigate('/app/admin/teachers/new')}
                            sx={{ borderRadius: 2 }}
                        >
                            Thêm Giáo viên
                        </Button>
                    )
                }
            />

            {error && (
                <Alert severity="error" sx={{ mb: 3 }}>
                    Không thể tải danh sách giáo viên. Vui lòng thử lại sau.
                </Alert>
            )}

            <Box sx={{ mb: 3 }}>
                <TeacherFilters filters={filters} onChange={handleFilterChange} />
            </Box>

            <TeacherListTable
                teachers={responseData?.data?.teachers || []}
                loading={isLoading}
                onView={(teacher) => navigate(`/app/admin/teachers/${teacher.id}`)}
                onEdit={(teacher) => navigate(`/app/admin/teachers/${teacher.id}/edit`)}
                onDelete={handleDeleteClick}
                showActions
                isAdmin={isAdmin}
                paginationModel={paginationModel}
                onPaginationModelChange={setPaginationModel}
                rowCount={responseData?.data?.pagination?.total_items || 0}
            />

            <DeleteTeacherDialog
                open={deleteDialog.open}
                teacherName={deleteDialog.teacher?.full_name || ''}
                onClose={() => setDeleteDialog({ open: false, teacher: null })}
                onConfirm={handleDeleteConfirm}
                loading={isDeleting}
            />
        </Container>
    );
};
