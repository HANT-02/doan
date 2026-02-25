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
import type { Student } from '@/api/studentApi';

const studentSchema = z.object({
    code: z.string().min(1, 'Mã học sinh không được để trống'),
    full_name: z.string().min(1, 'Họ tên không được để trống'),
    email: z.string().email('Email không hợp lệ').optional().or(z.literal('')),
    phone: z.string().optional(),
    guardian_phone: z.string().min(1, 'Số điện thoại phụ huynh không được để trống'),
    grade_level: z.string().min(1, 'Khối lớp không được để trống'),
    status: z.enum(['ACTIVE', 'INACTIVE']),
    date_of_birth: z.string().optional(),
    gender: z.string().optional(),
    address: z.string().optional(),
});

type StudentFormValues = z.infer<typeof studentSchema>;

interface StudentDialogProps {
    open: boolean;
    onClose: () => void;
    onSubmit: (data: StudentFormValues) => Promise<void>;
    student?: Student | null;
    isLoading?: boolean;
}

const StudentDialog = ({ open, onClose, onSubmit, student, isLoading }: StudentDialogProps) => {
    const {
        register,
        handleSubmit,
        reset,
        formState: { errors },
    } = useForm<StudentFormValues>({
        resolver: zodResolver(studentSchema) as any,
        defaultValues: {
            code: '',
            full_name: '',
            email: '',
            phone: '',
            guardian_phone: '',
            grade_level: '',
            status: 'ACTIVE',
            date_of_birth: '',
            gender: 'MALE',
            address: '',
        },
    });

    useEffect(() => {
        if (student) {
            reset({
                code: student.code,
                full_name: student.full_name,
                email: student.email || '',
                phone: student.phone || '',
                guardian_phone: student.guardian_phone || '',
                grade_level: student.grade_level || '',
                status: student.status,
                date_of_birth: student.date_of_birth ? student.date_of_birth.split('T')[0] : '',
                gender: student.gender || 'MALE',
                address: student.address || '',
            });
        } else {
            reset({
                code: '',
                full_name: '',
                email: '',
                phone: '',
                guardian_phone: '',
                grade_level: '',
                status: 'ACTIVE',
                date_of_birth: '',
                gender: 'MALE',
                address: '',
            });
        }
    }, [student, reset, open]);

    const handleFormSubmit = async (data: StudentFormValues) => {
        const submissionData = {
            ...data,
            date_of_birth: data.date_of_birth ? `${data.date_of_birth}T00:00:00Z` : null,
        };
        await onSubmit(submissionData as any);
        onClose();
    };

    return (
        <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
            <DialogTitle>
                {student ? 'Chỉnh sửa thông tin học sinh' : 'Thêm học sinh mới'}
            </DialogTitle>
            <DialogContent dividers>
                <Grid container spacing={2} sx={{ mt: 0.5 }}>
                    <Grid size={4}>
                        <TextField
                            fullWidth
                            label="Mã học sinh"
                            {...register('code')}
                            error={!!errors.code}
                            helperText={errors.code?.message}
                            required
                        />
                    </Grid>
                    <Grid size={8}>
                        <TextField
                            fullWidth
                            label="Họ và tên"
                            {...register('full_name')}
                            error={!!errors.full_name}
                            helperText={errors.full_name?.message}
                            required
                        />
                    </Grid>
                    <Grid size={6}>
                        <TextField
                            fullWidth
                            label="Email"
                            {...register('email')}
                            error={!!errors.email}
                            helperText={errors.email?.message}
                        />
                    </Grid>
                    <Grid size={6}>
                        <TextField
                            fullWidth
                            label="Số điện thoại"
                            {...register('phone')}
                            error={!!errors.phone}
                            helperText={errors.phone?.message}
                        />
                    </Grid>
                    <Grid size={6}>
                        <TextField
                            fullWidth
                            label="SĐT Phụ huynh"
                            {...register('guardian_phone')}
                            error={!!errors.guardian_phone}
                            helperText={errors.guardian_phone?.message}
                            required
                        />
                    </Grid>
                    <Grid size={6}>
                        <TextField
                            fullWidth
                            label="Khối lớp"
                            {...register('grade_level')}
                            error={!!errors.grade_level}
                            helperText={errors.grade_level?.message}
                            required
                        />
                    </Grid>
                    <Grid size={4}>
                        <TextField
                            fullWidth
                            label="Ngày sinh"
                            type="date"
                            {...register('date_of_birth')}
                            InputLabelProps={{ shrink: true }}
                            error={!!errors.date_of_birth}
                            helperText={errors.date_of_birth?.message}
                        />
                    </Grid>
                    <Grid size={4}>
                        <TextField
                            fullWidth
                            select
                            label="Giới tính"
                            {...register('gender')}
                            error={!!errors.gender}
                            helperText={errors.gender?.message}
                        >
                            <MenuItem value="MALE">Nam</MenuItem>
                            <MenuItem value="FEMALE">Nữ</MenuItem>
                            <MenuItem value="OTHER">Khác</MenuItem>
                        </TextField>
                    </Grid>
                    <Grid size={4}>
                        <TextField
                            fullWidth
                            select
                            label="Trạng thái"
                            {...register('status')}
                            error={!!errors.status}
                            helperText={errors.status?.message}
                        >
                            <MenuItem value="ACTIVE">Đang học</MenuItem>
                            <MenuItem value="INACTIVE">Nghỉ học</MenuItem>
                        </TextField>
                    </Grid>
                    <Grid size={12}>
                        <TextField
                            fullWidth
                            label="Địa chỉ"
                            multiline
                            rows={2}
                            {...register('address')}
                            error={!!errors.address}
                            helperText={errors.address?.message}
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
                    {student ? 'Cập nhật' : 'Thêm mới'}
                </Button>
            </DialogActions>
        </Dialog>
    );
};

export default StudentDialog;
