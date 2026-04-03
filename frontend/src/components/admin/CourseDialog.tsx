import { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import {
    Button,
    CircularProgress,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    Grid,
    MenuItem,
    TextField,
} from '@mui/material';

import type { Course } from '@/api/courseApi';

const courseSchema = z.object({
    code: z.string().min(1, 'Mã khóa học không được để trống'),
    name: z.string().min(1, 'Tên khóa học không được để trống'),
    description: z.string().optional(),
    subject: z.string().optional(),
    grade_level: z.string().optional(),
    session_count: z.coerce.number().min(0, 'Số buổi không hợp lệ'),
    session_duration_minutes: z.coerce.number().min(0, 'Thời lượng không hợp lệ'),
    total_hours: z.coerce.number().min(0, 'Tổng giờ không hợp lệ'),
    price: z.coerce.number().min(0, 'Học phí không hợp lệ'),
    status: z.enum(['ACTIVE', 'INACTIVE']),
});

type CourseFormValues = z.infer<typeof courseSchema>;

interface CourseDialogProps {
    open: boolean;
    onClose: () => void;
    onSubmit: (data: Partial<Course>) => Promise<void>;
    course?: Course | null;
    isLoading?: boolean;
}

const CourseDialog = ({ open, onClose, onSubmit, course, isLoading }: CourseDialogProps) => {
    const {
        register,
        handleSubmit,
        reset,
        formState: { errors },
    } = useForm<CourseFormValues>({
        resolver: zodResolver(courseSchema) as any,
        defaultValues: {
            code: '',
            name: '',
            description: '',
            subject: '',
            grade_level: '',
            session_count: 0,
            session_duration_minutes: 0,
            total_hours: 0,
            price: 0,
            status: 'ACTIVE',
        },
    });

    useEffect(() => {
        if (course) {
            reset({
                code: course.code || '',
                name: course.name || '',
                description: course.description || '',
                subject: course.subject || '',
                grade_level: course.grade_level || '',
                session_count: course.session_count || 0,
                session_duration_minutes: course.session_duration_minutes || 0,
                total_hours: course.total_hours || 0,
                price: course.price || 0,
                status: course.status === 'ACTIVE' ? 'ACTIVE' : 'INACTIVE',
            });
            return;
        }

        reset({
            code: '',
            name: '',
            description: '',
            subject: '',
            grade_level: '',
            session_count: 0,
            session_duration_minutes: 0,
            total_hours: 0,
            price: 0,
            status: 'ACTIVE',
        });
    }, [open, course, reset]);

    const handleFormSubmit = async (values: CourseFormValues) => {
        await onSubmit(values);
        onClose();
    };

    return (
        <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
            <DialogTitle>{course ? 'Chỉnh sửa khóa học' : 'Thêm khóa học mới'}</DialogTitle>
            <DialogContent dividers>
                <Grid container spacing={2} sx={{ mt: 0.5 }}>
                    <Grid size={{ xs: 12, md: 4 }}>
                        <TextField
                            fullWidth
                            label="Mã khóa học"
                            {...register('code')}
                            error={!!errors.code}
                            helperText={errors.code?.message}
                        />
                    </Grid>
                    <Grid size={{ xs: 12, md: 8 }}>
                        <TextField
                            fullWidth
                            label="Tên khóa học"
                            {...register('name')}
                            error={!!errors.name}
                            helperText={errors.name?.message}
                        />
                    </Grid>
                    <Grid size={{ xs: 12, md: 6 }}>
                        <TextField
                            fullWidth
                            label="Môn học"
                            {...register('subject')}
                            error={!!errors.subject}
                            helperText={errors.subject?.message}
                        />
                    </Grid>
                    <Grid size={{ xs: 12, md: 6 }}>
                        <TextField
                            fullWidth
                            label="Khối lớp"
                            {...register('grade_level')}
                            error={!!errors.grade_level}
                            helperText={errors.grade_level?.message}
                        />
                    </Grid>
                    <Grid size={{ xs: 12, md: 4 }}>
                        <TextField
                            fullWidth
                            type="number"
                            label="Số buổi"
                            {...register('session_count')}
                            error={!!errors.session_count}
                            helperText={errors.session_count?.message}
                        />
                    </Grid>
                    <Grid size={{ xs: 12, md: 4 }}>
                        <TextField
                            fullWidth
                            type="number"
                            label="Phút / buổi"
                            {...register('session_duration_minutes')}
                            error={!!errors.session_duration_minutes}
                            helperText={errors.session_duration_minutes?.message}
                        />
                    </Grid>
                    <Grid size={{ xs: 12, md: 4 }}>
                        <TextField
                            fullWidth
                            select
                            label="Trạng thái"
                            defaultValue="ACTIVE"
                            {...register('status')}
                            error={!!errors.status}
                            helperText={errors.status?.message}
                        >
                            <MenuItem value="ACTIVE">Hoạt động</MenuItem>
                            <MenuItem value="INACTIVE">Tạm dừng</MenuItem>
                        </TextField>
                    </Grid>
                    <Grid size={{ xs: 12, md: 6 }}>
                        <TextField
                            fullWidth
                            type="number"
                            label="Tổng giờ"
                            {...register('total_hours')}
                            error={!!errors.total_hours}
                            helperText={errors.total_hours?.message}
                        />
                    </Grid>
                    <Grid size={{ xs: 12, md: 6 }}>
                        <TextField
                            fullWidth
                            type="number"
                            label="Học phí"
                            {...register('price')}
                            error={!!errors.price}
                            helperText={errors.price?.message}
                        />
                    </Grid>
                    <Grid size={12}>
                        <TextField
                            fullWidth
                            multiline
                            rows={3}
                            label="Mô tả khóa học"
                            {...register('description')}
                            error={!!errors.description}
                            helperText={errors.description?.message}
                        />
                    </Grid>
                </Grid>
            </DialogContent>
            <DialogActions sx={{ px: 3, py: 2 }}>
                <Button onClick={onClose} disabled={isLoading}>
                    Hủy
                </Button>
                <Button
                    variant="contained"
                    onClick={handleSubmit(handleFormSubmit)}
                    disabled={isLoading}
                    startIcon={isLoading ? <CircularProgress size={20} color="inherit" /> : null}
                >
                    {course ? 'Cập nhật' : 'Thêm mới'}
                </Button>
            </DialogActions>
        </Dialog>
    );
};

export default CourseDialog;
