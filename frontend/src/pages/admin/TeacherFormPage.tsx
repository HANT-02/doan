import { useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import {
    Container,
    Paper,
    Box,
    Typography,
    Button,
    TextField,
    Grid,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
    FormControlLabel,
    Switch,
    Alert,
    Skeleton,
    Stack
} from '@mui/material';
import { ArrowBack, Save } from '@mui/icons-material';
import {
    useGetTeacherByIdQuery,
    useCreateTeacherMutation,
    useUpdateTeacherMutation
} from '@/api/teacherApi';
import { createTeacherSchema, updateTeacherSchema } from '@/schemas/teacherSchema';
import { toast } from 'sonner';

export const TeacherFormPage = () => {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const isEditMode = Boolean(id);

    const { data: responseData, isLoading: isFetching, error: fetchError } = useGetTeacherByIdQuery(id!, { skip: !isEditMode });
    const [createTeacher, { isLoading: isCreating }] = useCreateTeacherMutation();
    const [updateTeacher, { isLoading: isUpdating }] = useUpdateTeacherMutation();

    const {
        control,
        handleSubmit,
        reset,
        formState: { errors },
    } = useForm<any>({
        resolver: zodResolver(isEditMode ? updateTeacherSchema : createTeacherSchema),
        defaultValues: {
            full_name: '',
            email: '',
            phone: '',
            code: '',
            is_school_teacher: false,
            school_name: '',
            employment_type: 'FULL_TIME',
            status: 'ACTIVE',
            notes: '',
        },
    });

    useEffect(() => {
        if (isEditMode && responseData?.data) {
            const teacher = responseData.data;
            reset({
                full_name: teacher.full_name,
                email: teacher.email || '',
                phone: teacher.phone || '',
                code: teacher.code || '',
                is_school_teacher: teacher.is_school_teacher || false,
                school_name: teacher.school_name || '',
                employment_type: teacher.employment_type || 'FULL_TIME',
                status: teacher.status || 'ACTIVE',
                notes: teacher.notes || '',
            });
        }
    }, [responseData, isEditMode, reset]);

    const onSubmit = async (data: any) => {
        try {
            if (isEditMode && id) {
                await updateTeacher({ id, ...data }).unwrap();
                toast.success('Cập nhật giáo viên thành công');
                navigate(`/app/admin/teachers/${id}`);
            } else {
                const result = await createTeacher(data).unwrap();
                toast.success('Thêm giáo viên thành công');
                navigate(`/app/admin/teachers/${result.data.id}`);
            }
        } catch (err: any) {
            toast.error(err?.data?.message || 'Lỗi khi lưu thông tin giáo viên');
        }
    };

    if (isFetching) {
        return (
            <Container maxWidth="md" sx={{ py: 4 }}>
                <Skeleton variant="rectangular" height={600} sx={{ borderRadius: 2 }} />
            </Container>
        );
    }

    if (fetchError) {
        return (
            <Container maxWidth="md" sx={{ py: 4 }}>
                <Alert severity="error">Không thể tải thông tin giáo viên</Alert>
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

    const submitting = isCreating || isUpdating;

    return (
        <Container maxWidth="md" sx={{ py: 4 }}>
            <Box sx={{ mb: 3 }}>
                <Button
                    startIcon={<ArrowBack />}
                    onClick={() => navigate(isEditMode ? `/app/admin/teachers/${id}` : '/app/admin/teachers')}
                    sx={{ mb: 2 }}
                >
                    Quay lại
                </Button>

                <Typography variant="h4" component="h1" sx={{ fontWeight: 700 }} gutterBottom>
                    {isEditMode ? 'Chỉnh sửa Giáo viên' : 'Thêm Giáo viên mới'}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    {isEditMode ? 'Cập nhật thông tin chi tiết của giáo viên' : 'Nhập thông tin để tạo hồ sơ giáo viên mới'}
                </Typography>
            </Box>

            <Paper sx={{ p: 4, borderRadius: 3 }}>
                <form onSubmit={handleSubmit(onSubmit)}>
                    <Grid container spacing={3}>
                        <Grid size={12}>
                            <Controller
                                name="full_name"
                                control={control}
                                render={({ field }) => (
                                    <TextField
                                        {...field}
                                        label="Họ và tên"
                                        fullWidth
                                        required
                                        error={!!errors.full_name}
                                        helperText={errors.full_name?.message as string}
                                    />
                                )}
                            />
                        </Grid>

                        <Grid size={{ xs: 12, sm: 6 }}>
                            <Controller
                                name="email"
                                control={control}
                                render={({ field }) => (
                                    <TextField
                                        {...field}
                                        label="Email"
                                        type="email"
                                        fullWidth
                                        error={!!errors.email}
                                        helperText={errors.email?.message as string}
                                    />
                                )}
                            />
                        </Grid>

                        <Grid size={{ xs: 12, sm: 6 }}>
                            <Controller
                                name="phone"
                                control={control}
                                render={({ field }) => (
                                    <TextField
                                        {...field}
                                        label="Số điện thoại"
                                        fullWidth
                                        error={!!errors.phone}
                                        helperText={errors.phone?.message as string}
                                    />
                                )}
                            />
                        </Grid>

                        <Grid size={{ xs: 12, sm: 6 }}>
                            <Controller
                                name="code"
                                control={control}
                                render={({ field }) => (
                                    <TextField
                                        {...field}
                                        label="Mã giáo viên"
                                        fullWidth
                                        error={!!errors.code}
                                        helperText={errors.code?.message as string}
                                    />
                                )}
                            />
                        </Grid>

                        <Grid size={{ xs: 12, sm: 6 }}>
                            <Controller
                                name="employment_type"
                                control={control}
                                render={({ field }) => (
                                    <FormControl fullWidth error={!!errors.employment_type}>
                                        <InputLabel>Hình thức làm việc</InputLabel>
                                        <Select {...field} label="Hình thức làm việc">
                                            <MenuItem value="FULL_TIME">Toàn thời gian</MenuItem>
                                            <MenuItem value="PART_TIME">Bán thời gian</MenuItem>
                                            <MenuItem value="CONTRACTOR">Hợp đồng</MenuItem>
                                        </Select>
                                    </FormControl>
                                )}
                            />
                        </Grid>

                        <Grid size={{ xs: 12, sm: 6 }}>
                            <Controller
                                name="status"
                                control={control}
                                render={({ field }) => (
                                    <FormControl fullWidth error={!!errors.status}>
                                        <InputLabel>Trạng thái</InputLabel>
                                        <Select {...field} label="Trạng thái">
                                            <MenuItem value="ACTIVE">Đang hoạt động</MenuItem>
                                            <MenuItem value="INACTIVE">Ngừng hoạt động</MenuItem>
                                        </Select>
                                    </FormControl>
                                )}
                            />
                        </Grid>

                        <Grid size={12}>
                            <Controller
                                name="is_school_teacher"
                                control={control}
                                render={({ field }) => (
                                    <FormControlLabel
                                        control={<Switch {...field} checked={field.value} />}
                                        label="Là giáo viên trường"
                                    />
                                )}
                            />
                        </Grid>

                        <Grid size={12}>
                            <Controller
                                name="school_name"
                                control={control}
                                render={({ field }) => (
                                    <TextField
                                        {...field}
                                        label="Tên trường (nếu có)"
                                        fullWidth
                                        error={!!errors.school_name}
                                        helperText={errors.school_name?.message as string}
                                    />
                                )}
                            />
                        </Grid>

                        <Grid size={12}>
                            <Controller
                                name="notes"
                                control={control}
                                render={({ field }) => (
                                    <TextField
                                        {...field}
                                        label="Ghi chú"
                                        fullWidth
                                        multiline
                                        rows={4}
                                        error={!!errors.notes}
                                        helperText={errors.notes?.message as string}
                                    />
                                )}
                            />
                        </Grid>

                        <Grid size={12}>
                            <Stack direction="row" spacing={2} justifyContent="flex-end">
                                <Button
                                    variant="outlined"
                                    onClick={() =>
                                        navigate(isEditMode ? `/app/admin/teachers/${id}` : '/app/admin/teachers')
                                    }
                                    disabled={submitting}
                                >
                                    Hủy
                                </Button>
                                <Button
                                    type="submit"
                                    variant="contained"
                                    startIcon={<Save />}
                                    disabled={submitting}
                                >
                                    {submitting ? 'Đang lưu...' : isEditMode ? 'Cập nhật' : 'Thêm mới'}
                                </Button>
                            </Stack>
                        </Grid>
                    </Grid>
                </form>
            </Paper>
        </Container>
    );
};
