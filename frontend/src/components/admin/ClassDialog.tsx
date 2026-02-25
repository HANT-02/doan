import { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import {
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    Button,
    TextField,
    MenuItem,
    Grid,
    CircularProgress,
} from '@mui/material';
import type { Class } from '@/api/classApi';
import { useGetTeachersQuery } from '@/api/teacherApi';

const classSchema = z.object({
    code: z.string().min(1, 'Mã lớp không được để trống'),
    name: z.string().min(1, 'Tên lớp không được để trống'),
    start_date: z.string().min(1, 'Ngày bắt đầu không được để trống'),
    end_date: z.string().optional(),
    max_students: z.coerce.number().min(1, 'Sĩ số tối đa phải lớn hơn 0'),
    status: z.enum(['OPEN', 'CLOSED', 'CANCELLED']),
    price: z.coerce.number().min(0, 'Học phí không được âm'),
    teacher_id: z.string().optional(),
    program_id: z.string().optional(),
    course_id: z.string().optional(),
    notes: z.string().optional(),
});

type ClassFormValues = z.infer<typeof classSchema>;

interface ClassDialogProps {
    open: boolean;
    onClose: () => void;
    onSubmit: (data: ClassFormValues) => Promise<void>;
    classData?: Class | null;
    isLoading?: boolean;
}

const ClassDialog = ({ open, onClose, onSubmit, classData, isLoading }: ClassDialogProps) => {
    const { data: teachersData } = useGetTeachersQuery({ limit: 100 });
    // Note: Program and Course APIs are missing in current repo evidence, 
    // using placeholder selectors for now or simple text IDs.

    const {
        register,
        handleSubmit,
        reset,
        formState: { errors },
    } = useForm<ClassFormValues>({
        resolver: zodResolver(classSchema) as any,
        defaultValues: {
            code: '',
            name: '',
            start_date: new Date().toISOString().split('T')[0],
            max_students: 30,
            status: 'OPEN',
            price: 0,
        },
    });

    useEffect(() => {
        if (classData) {
            reset({
                code: classData.code,
                name: classData.name,
                start_date: classData.start_date.split('T')[0],
                end_date: classData.end_date?.split('T')[0] || '',
                max_students: classData.max_students,
                status: classData.status,
                price: classData.price,
                teacher_id: classData.teacher_id || '',
                program_id: classData.program_id || '',
                course_id: classData.course_id || '',
                notes: classData.notes || '',
            });
        } else {
            reset({
                code: '',
                name: '',
                start_date: new Date().toISOString().split('T')[0],
                end_date: '',
                max_students: 30,
                status: 'OPEN',
                price: 0,
                teacher_id: '',
                program_id: '',
                course_id: '',
                notes: '',
            });
        }
    }, [classData, reset, open]);

    const handleFormSubmit = async (data: ClassFormValues) => {
        // Clean up empty strings for optional IDs
        const cleanedData = {
            ...data,
            teacher_id: data.teacher_id || null,
            program_id: data.program_id || null,
            course_id: data.course_id || null,
        };
        await onSubmit(cleanedData as any);
        onClose();
    };

    return (
        <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
            <DialogTitle>
                {classData ? 'Chỉnh sửa lớp học' : 'Thêm lớp học mới'}
            </DialogTitle>
            <DialogContent dividers>
                <Grid container spacing={2} sx={{ mt: 0.5 }}>
                    <Grid size={4}>
                        <TextField
                            fullWidth
                            label="Mã lớp"
                            {...register('code')}
                            error={!!errors.code}
                            helperText={errors.code?.message}
                        />
                    </Grid>
                    <Grid size={8}>
                        <TextField
                            fullWidth
                            label="Tên lớp"
                            {...register('name')}
                            error={!!errors.name}
                            helperText={errors.name?.message}
                        />
                    </Grid>
                    <Grid size={6}>
                        <TextField
                            fullWidth
                            label="Ngày bắt đầu"
                            type="date"
                            InputLabelProps={{ shrink: true }}
                            {...register('start_date')}
                            error={!!errors.start_date}
                            helperText={errors.start_date?.message}
                        />
                    </Grid>
                    <Grid size={6}>
                        <TextField
                            fullWidth
                            label="Ngày kết thúc"
                            type="date"
                            InputLabelProps={{ shrink: true }}
                            {...register('end_date')}
                            error={!!errors.end_date}
                            helperText={errors.end_date?.message}
                        />
                    </Grid>
                    <Grid size={4}>
                        <TextField
                            fullWidth
                            label="Sĩ số tối đa"
                            type="number"
                            {...register('max_students')}
                            error={!!errors.max_students}
                            helperText={errors.max_students?.message}
                        />
                    </Grid>
                    <Grid size={4}>
                        <TextField
                            fullWidth
                            label="Học phí"
                            type="number"
                            {...register('price')}
                            error={!!errors.price}
                            helperText={errors.price?.message}
                        />
                    </Grid>
                    <Grid size={4}>
                        <TextField
                            fullWidth
                            select
                            label="Trạng thái"
                            defaultValue="OPEN"
                            {...register('status')}
                            error={!!errors.status}
                            helperText={errors.status?.message}
                        >
                            <MenuItem value="OPEN">Đang mở</MenuItem>
                            <MenuItem value="CLOSED">Đã đóng</MenuItem>
                            <MenuItem value="CANCELLED">Hủy bỏ</MenuItem>
                        </TextField>
                    </Grid>
                    <Grid size={12}>
                        <TextField
                            fullWidth
                            select
                            label="Giáo viên chủ nhiệm"
                            {...register('teacher_id')}
                            defaultValue=""
                        >
                            <MenuItem value=""><em>Chưa gán</em></MenuItem>
                            {teachersData?.data?.teachers?.map((t: any) => (
                                <MenuItem key={t.id} value={t.id}>
                                    {t.user?.full_name || 'N/A'} ({t.teacher_code || 'N/A'})
                                </MenuItem>
                            ))}
                        </TextField>
                    </Grid>
                    <Grid size={12}>
                        <TextField
                            fullWidth
                            label="Ghi chú"
                            multiline
                            rows={2}
                            {...register('notes')}
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
                    {classData ? 'Cập nhật' : 'Thêm mới'}
                </Button>
            </DialogActions>
        </Dialog>
    );
};

export default ClassDialog;
