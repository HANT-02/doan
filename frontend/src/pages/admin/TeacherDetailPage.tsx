import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
    Container,
    Paper,
    Box,
    Typography,
    Button,
    Grid,
    Chip,
    Divider,
    Alert,
    Skeleton,
    Stack,
} from '@mui/material';
import { ArrowBack, Edit, Delete } from '@mui/icons-material';
import { teacherApi, type Teacher } from '@/api/teacherApi';
import { DeleteTeacherDialog } from '@/components/teacher/DeleteTeacherDialog';
import { useAuth } from '@/contexts/AuthContext';
import { toast } from 'sonner';
import {
    getEmploymentTypeLabel,
    getStatusColor,
    formatDateTime,
} from '@/utils/teacherHelpers';

export const TeacherDetailPage = () => {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const { user } = useAuth();
    const isAdmin = user?.role === 'admin' || user?.role === 'super_admin';

    const [teacher, setTeacher] = useState<Teacher | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [deleteDialog, setDeleteDialog] = useState(false);
    const [deleting, setDeleting] = useState(false);

    useEffect(() => {
        if (!id) return;

        const fetchTeacher = async () => {
            try {
                setLoading(true);
                setError(null);
                const data = await teacherApi.getById(id);
                setTeacher(data);
            } catch (err) {
                setError(err instanceof Error ? err.message : 'Failed to load teacher');
                toast.error('Failed to load teacher details');
            } finally {
                setLoading(false);
            }
        };

        fetchTeacher();
    }, [id]);

    const handleDelete = async () => {
        if (!teacher) return;

        try {
            setDeleting(true);
            await teacherApi.delete(teacher.id);
            toast.success('Teacher deleted successfully');
            navigate('/app/admin/teachers');
        } catch (err) {
            toast.error(err instanceof Error ? err.message : 'Failed to delete teacher');
        } finally {
            setDeleting(false);
            setDeleteDialog(false);
        }
    };

    if (loading) {
        return (
            <Container maxWidth="lg" sx={{ py: 4 }}>
                <Skeleton variant="rectangular" height={400} />
            </Container>
        );
    }

    if (error || !teacher) {
        return (
            <Container maxWidth="lg" sx={{ py: 4 }}>
                <Alert severity="error">{error || 'Teacher not found'}</Alert>
                <Button
                    startIcon={<ArrowBack />}
                    onClick={() => navigate('/app/admin/teachers')}
                    sx={{ mt: 2 }}
                >
                    Back to Teachers
                </Button>
            </Container>
        );
    }

    return (
        <Container maxWidth="lg" sx={{ py: 4 }}>
            <Box sx={{ mb: 3 }}>
                <Button
                    startIcon={<ArrowBack />}
                    onClick={() => navigate('/app/admin/teachers')}
                    sx={{ mb: 2 }}
                >
                    Back to Teachers
                </Button>

                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                    <Box>
                        <Typography variant="h4" component="h1" gutterBottom>
                            {teacher.full_name}
                        </Typography>
                        <Stack direction="row" spacing={1}>
                            <Chip label={teacher.status} color={getStatusColor(teacher.status)} size="small" />
                            <Chip label={getEmploymentTypeLabel(teacher.employment_type)} size="small" />
                            {teacher.is_school_teacher && <Chip label="School Teacher" size="small" />}
                        </Stack>
                    </Box>

                    {isAdmin && (
                        <Stack direction="row" spacing={1}>
                            <Button
                                variant="contained"
                                startIcon={<Edit />}
                                onClick={() => navigate(`/app/admin/teachers/${teacher.id}/edit`)}
                            >
                                Edit
                            </Button>
                            <Button
                                variant="outlined"
                                color="error"
                                startIcon={<Delete />}
                                onClick={() => setDeleteDialog(true)}
                            >
                                Delete
                            </Button>
                        </Stack>
                    )}
                </Box>
            </Box>

            <Paper sx={{ p: 3 }}>
                <Typography variant="h6" gutterBottom>
                    Basic Information
                </Typography>
                <Divider sx={{ mb: 3 }} />

                <Grid container spacing={3}>
                    <Grid item xs={12} sm={6}>
                        <Typography variant="caption" color="text.secondary">
                            Code
                        </Typography>
                        <Typography variant="body1">{teacher.code || '-'}</Typography>
                    </Grid>

                    <Grid item xs={12} sm={6}>
                        <Typography variant="caption" color="text.secondary">
                            Full Name
                        </Typography>
                        <Typography variant="body1">{teacher.full_name}</Typography>
                    </Grid>

                    <Grid item xs={12} sm={6}>
                        <Typography variant="caption" color="text.secondary">
                            Email
                        </Typography>
                        <Typography variant="body1">{teacher.email || '-'}</Typography>
                    </Grid>

                    <Grid item xs={12} sm={6}>
                        <Typography variant="caption" color="text.secondary">
                            Phone
                        </Typography>
                        <Typography variant="body1">{teacher.phone || '-'}</Typography>
                    </Grid>

                    <Grid item xs={12} sm={6}>
                        <Typography variant="caption" color="text.secondary">
                            Employment Type
                        </Typography>
                        <Typography variant="body1">
                            {getEmploymentTypeLabel(teacher.employment_type)}
                        </Typography>
                    </Grid>

                    <Grid item xs={12} sm={6}>
                        <Typography variant="caption" color="text.secondary">
                            Status
                        </Typography>
                        <Typography variant="body1">{teacher.status}</Typography>
                    </Grid>

                    {teacher.is_school_teacher && (
                        <Grid item xs={12}>
                            <Typography variant="caption" color="text.secondary">
                                School Name
                            </Typography>
                            <Typography variant="body1">{teacher.school_name || '-'}</Typography>
                        </Grid>
                    )}

                    {teacher.notes && (
                        <Grid item xs={12}>
                            <Typography variant="caption" color="text.secondary">
                                Notes
                            </Typography>
                            <Typography variant="body1" sx={{ whiteSpace: 'pre-wrap' }}>
                                {teacher.notes}
                            </Typography>
                        </Grid>
                    )}

                    <Grid item xs={12} sm={6}>
                        <Typography variant="caption" color="text.secondary">
                            Created At
                        </Typography>
                        <Typography variant="body1">{formatDateTime(teacher.created_at)}</Typography>
                    </Grid>

                    <Grid item xs={12} sm={6}>
                        <Typography variant="caption" color="text.secondary">
                            Updated At
                        </Typography>
                        <Typography variant="body1">{formatDateTime(teacher.updated_at)}</Typography>
                    </Grid>
                </Grid>
            </Paper>

            <DeleteTeacherDialog
                open={deleteDialog}
                teacherName={teacher.full_name}
                onClose={() => setDeleteDialog(false)}
                onConfirm={handleDelete}
                loading={deleting}
            />
        </Container>
    );
};
