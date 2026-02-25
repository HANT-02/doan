import { useState } from 'react';
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
import { useGetTeacherByIdQuery, useDeleteTeacherMutation } from '@/api/teacherApi';
import { DeleteTeacherDialog } from '@/components/teacher/DeleteTeacherDialog';
import { useAppSelector } from '@/store';
import { toast } from 'sonner';
import {
    getEmploymentTypeLabel,
    getStatusColor,
    formatDateTime,
} from '@/utils/teacherHelpers';

export const TeacherDetailPage = () => {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const { user } = useAppSelector((state) => state.auth);
    const isAdmin = user?.role === 'admin' || user?.role === 'super_admin';

    const { data: responseData, isLoading, error } = useGetTeacherByIdQuery(id!);
    const [deleteTeacher, { isLoading: isDeleting }] = useDeleteTeacherMutation();

    const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

    const handleDelete = async () => {
        if (!id) return;

        try {
            await deleteTeacher(id).unwrap();
            toast.success('Xóa giáo viên thành công');
            navigate('/app/admin/teachers');
        } catch (err: any) {
            toast.error(err?.data?.message || 'Lỗi khi xóa giáo viên');
        } finally {
            setDeleteDialogOpen(false);
        }
    };

    if (isLoading) {
        return (
            <Container maxWidth="lg" sx={{ py: 4 }}>
                <Skeleton variant="rectangular" height={400} sx={{ borderRadius: 2 }} />
            </Container>
        );
    }

    if (error || !responseData?.data) {
        return (
            <Container maxWidth="lg" sx={{ py: 4 }}>
                <Alert severity="error">{(error as any)?.data?.message || 'Không tìm thấy giáo viên'}</Alert>
                <Button
                    startIcon={<ArrowBack />}
                    onClick={() => navigate('/app/admin/teachers')}
                    sx={{ mt: 2 }}
                >
                    Quay lại danh sách
                </Button>
            </Container>
        );
    }

    const teacher = responseData.data;

    return (
        <Container maxWidth="lg" sx={{ py: 4 }}>
            <Box sx={{ mb: 3 }}>
                <Button
                    startIcon={<ArrowBack />}
                    onClick={() => navigate('/app/admin/teachers')}
                    sx={{ mb: 2 }}
                >
                    Quay lại danh sách
                </Button>

                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', flexWrap: 'wrap', gap: 2 }}>
                    <Box>
                        <Typography variant="h4" component="h1" sx={{ fontWeight: 700 }} gutterBottom>
                            {teacher.full_name}
                        </Typography>
                        <Stack direction="row" spacing={1}>
                            <Chip label={teacher.status} color={getStatusColor(teacher.status)} size="small" />
                            <Chip label={getEmploymentTypeLabel(teacher.employment_type)} size="small" />
                            {teacher.is_school_teacher && <Chip label="Giáo viên trường" size="small" variant="outlined" color="primary" />}
                        </Stack>
                    </Box>

                    {isAdmin && (
                        <Stack direction="row" spacing={1}>
                            <Button
                                variant="contained"
                                startIcon={<Edit />}
                                onClick={() => navigate(`/app/admin/teachers/${teacher.id}/edit`)}
                            >
                                Chỉnh sửa
                            </Button>
                            <Button
                                variant="outlined"
                                color="error"
                                startIcon={<Delete />}
                                onClick={() => setDeleteDialogOpen(true)}
                            >
                                Xóa
                            </Button>
                        </Stack>
                    )}
                </Box>
            </Box>

            <Paper sx={{ p: 4, borderRadius: 3 }}>
                <Typography variant="h6" sx={{ fontWeight: 600 }} gutterBottom>
                    Thông tin cơ bản
                </Typography>
                <Divider sx={{ mb: 3 }} />

                <Grid container spacing={3}>
                    <Grid size={{ xs: 12, sm: 6 }}>
                        <Typography variant="caption" color="text.secondary">
                            Mã giáo viên
                        </Typography>
                        <Typography variant="body1" sx={{ fontWeight: 500 }}>{teacher.code || '-'}</Typography>
                    </Grid>

                    <Grid size={{ xs: 12, sm: 6 }}>
                        <Typography variant="caption" color="text.secondary">
                            Họ và tên
                        </Typography>
                        <Typography variant="body1" sx={{ fontWeight: 500 }}>{teacher.full_name}</Typography>
                    </Grid>

                    <Grid size={{ xs: 12, sm: 6 }}>
                        <Typography variant="caption" color="text.secondary">
                            Email
                        </Typography>
                        <Typography variant="body1" sx={{ fontWeight: 500 }}>{teacher.email || '-'}</Typography>
                    </Grid>

                    <Grid size={{ xs: 12, sm: 6 }}>
                        <Typography variant="caption" color="text.secondary">
                            Số điện thoại
                        </Typography>
                        <Typography variant="body1" sx={{ fontWeight: 500 }}>{teacher.phone || '-'}</Typography>
                    </Grid>

                    <Grid size={{ xs: 12, sm: 6 }}>
                        <Typography variant="caption" color="text.secondary">
                            Hình thức làm việc
                        </Typography>
                        <Typography variant="body1" sx={{ fontWeight: 500 }}>
                            {getEmploymentTypeLabel(teacher.employment_type)}
                        </Typography>
                    </Grid>

                    <Grid size={{ xs: 12, sm: 6 }}>
                        <Typography variant="caption" color="text.secondary">
                            Trạng thái
                        </Typography>
                        <Typography variant="body1" sx={{ fontWeight: 500 }}>{teacher.status}</Typography>
                    </Grid>

                    {teacher.is_school_teacher && (
                        <Grid size={12}>
                            <Typography variant="caption" color="text.secondary">
                                Tên trường
                            </Typography>
                            <Typography variant="body1" sx={{ fontWeight: 500 }}>{teacher.school_name || '-'}</Typography>
                        </Grid>
                    )}

                    {teacher.notes && (
                        <Grid size={12}>
                            <Typography variant="caption" color="text.secondary">
                                Ghi chú
                            </Typography>
                            <Typography variant="body1" sx={{ whiteSpace: 'pre-wrap', fontWeight: 500 }}>
                                {teacher.notes}
                            </Typography>
                        </Grid>
                    )}

                    <Grid size={{ xs: 12, sm: 6 }}>
                        <Typography variant="caption" color="text.secondary">
                            Ngày tạo
                        </Typography>
                        <Typography variant="body1" sx={{ fontWeight: 500 }}>{formatDateTime(teacher.created_at)}</Typography>
                    </Grid>

                    <Grid size={{ xs: 12, sm: 6 }}>
                        <Typography variant="caption" color="text.secondary">
                            Cập nhật gần nhất
                        </Typography>
                        <Typography variant="body1" sx={{ fontWeight: 500 }}>{formatDateTime(teacher.updated_at)}</Typography>
                    </Grid>
                </Grid>
            </Paper>

            <DeleteTeacherDialog
                open={deleteDialogOpen}
                teacherName={teacher.full_name}
                onClose={() => setDeleteDialogOpen(false)}
                onConfirm={handleDelete}
                loading={isDeleting}
            />
        </Container>
    );
};
