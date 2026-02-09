import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
    Box,
    Container,
    Typography,
    Button,
    Pagination,
    Alert,
    Stack,
} from '@mui/material';
import { Add } from '@mui/icons-material';
import { teacherApi, type Teacher, type ListTeachersParams } from '@/api/teacherApi';
import { TeacherListTable } from '@/components/teacher/TeacherListTable';
import { TeacherFilters, type TeacherFiltersState } from '@/components/teacher/TeacherFilters';
import { DeleteTeacherDialog } from '@/components/teacher/DeleteTeacherDialog';
import { useAuth } from '@/contexts/AuthContext';
import { toast } from 'sonner';

export const TeachersPage = () => {
    const navigate = useNavigate();
    const { user } = useAuth();
    const isAdmin = user?.role === 'admin' || user?.role === 'super_admin';

    const [teachers, setTeachers] = useState<Teacher[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);
    const [totalItems, setTotalItems] = useState(0);
    const [filters, setFilters] = useState<TeacherFiltersState>({
        search: '',
        status: '',
        employment_type: '',
    });

    const [deleteDialog, setDeleteDialog] = useState<{
        open: boolean;
        teacher: Teacher | null;
    }>({
        open: false,
        teacher: null,
    });
    const [deleting, setDeleting] = useState(false);

    const fetchTeachers = async () => {
        try {
            setLoading(true);
            setError(null);

            const params: ListTeachersParams = {
                page,
                limit: 10,
                ...filters,
            };

            // Remove empty filters
            Object.keys(params).forEach((key) => {
                if (params[key as keyof ListTeachersParams] === '') {
                    delete params[key as keyof ListTeachersParams];
                }
            });

            const response = await teacherApi.list(params);
            setTeachers(response.teachers || []);
            setTotalPages(response.pagination.total_pages);
            setTotalItems(response.pagination.total_items);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to load teachers');
            toast.error('Failed to load teachers');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchTeachers();
    }, [page, filters]);

    const handleFilterChange = (newFilters: TeacherFiltersState) => {
        setFilters(newFilters);
        setPage(1); // Reset to first page when filters change
    };

    const handleView = (teacher: Teacher) => {
        navigate(`/app/admin/teachers/${teacher.id}`);
    };

    const handleEdit = (teacher: Teacher) => {
        navigate(`/app/admin/teachers/${teacher.id}/edit`);
    };

    const handleDeleteClick = (teacher: Teacher) => {
        setDeleteDialog({ open: true, teacher });
    };

    const handleDeleteConfirm = async () => {
        if (!deleteDialog.teacher) return;

        try {
            setDeleting(true);
            await teacherApi.delete(deleteDialog.teacher.id);
            toast.success('Teacher deleted successfully');
            setDeleteDialog({ open: false, teacher: null });
            fetchTeachers(); // Refresh list
        } catch (err) {
            toast.error(err instanceof Error ? err.message : 'Failed to delete teacher');
        } finally {
            setDeleting(false);
        }
    };

    const handleDeleteCancel = () => {
        setDeleteDialog({ open: false, teacher: null });
    };

    return (
        <Container maxWidth="xl" sx={{ py: 4 }}>
            <Box sx={{ mb: 4, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <Box>
                    <Typography variant="h4" component="h1" gutterBottom>
                        Teachers
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        Manage teachers and their information
                    </Typography>
                </Box>
                {isAdmin && (
                    <Button
                        variant="contained"
                        startIcon={<Add />}
                        onClick={() => navigate('/app/admin/teachers/new')}
                    >
                        Add Teacher
                    </Button>
                )}
            </Box>

            {error && (
                <Alert severity="error" sx={{ mb: 3 }} onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            <TeacherFilters filters={filters} onChange={handleFilterChange} />

            <TeacherListTable
                teachers={teachers}
                loading={loading}
                onView={handleView}
                onEdit={handleEdit}
                onDelete={handleDeleteClick}
                showActions
                isAdmin={isAdmin}
            />

            {!loading && totalPages > 1 && (
                <Stack spacing={2} alignItems="center" sx={{ mt: 3 }}>
                    <Typography variant="body2" color="text.secondary">
                        Showing {teachers.length} of {totalItems} teachers
                    </Typography>
                    <Pagination
                        count={totalPages}
                        page={page}
                        onChange={(_, value) => setPage(value)}
                        color="primary"
                        size="large"
                    />
                </Stack>
            )}

            <DeleteTeacherDialog
                open={deleteDialog.open}
                teacherName={deleteDialog.teacher?.full_name || ''}
                onClose={handleDeleteCancel}
                onConfirm={handleDeleteConfirm}
                loading={deleting}
            />
        </Container>
    );
};
