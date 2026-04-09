import { useEffect } from 'react';
import { type Resolver, useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import {
    Button,
    CircularProgress,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    FormControlLabel,
    Grid,
    MenuItem,
    Switch,
    TextField,
} from '@mui/material';

import type { Shift, ShiftSessionType } from '@/api/shiftApi';

const shiftSchema = z.object({
    code: z.string().min(1, 'Mã ca học không được để trống'),
    name: z.string().min(1, 'Tên ca học không được để trống'),
    start_time: z.string().min(1, 'Giờ bắt đầu không được để trống'),
    end_time: z.string().min(1, 'Giờ kết thúc không được để trống'),
    duration_minutes: z.coerce.number().min(1, 'Thời lượng phải lớn hơn 0'),
    session_type: z.enum(['MORNING', 'AFTERNOON', 'EVENING', 'CUSTOM']),
    is_active: z.boolean(),
    notes: z.string().optional(),
});

type ShiftFormValues = z.infer<typeof shiftSchema>;

interface ShiftDialogProps {
    open: boolean;
    onClose: () => void;
    onSubmit: (data: ShiftFormValues) => Promise<void>;
    shift?: Shift | null;
    isLoading?: boolean;
}

const ShiftDialog = ({ open, onClose, onSubmit, shift, isLoading }: ShiftDialogProps) => {
    const {
        register,
        handleSubmit,
        reset,
        watch,
        setValue,
        formState: { errors },
    } = useForm<ShiftFormValues>({
        resolver: zodResolver(shiftSchema) as Resolver<ShiftFormValues>,
        defaultValues: {
            code: '',
            name: '',
            start_time: '18:00',
            end_time: '20:00',
            duration_minutes: 120,
            session_type: 'EVENING',
            is_active: true,
            notes: '',
        },
    });

    useEffect(() => {
        if (shift) {
            reset({
                code: shift.code,
                name: shift.name,
                start_time: shift.start_time,
                end_time: shift.end_time,
                duration_minutes: shift.duration_minutes,
                session_type: shift.session_type as ShiftSessionType,
                is_active: shift.is_active,
                notes: shift.notes || '',
            });
            return;
        }

        reset({
            code: '',
            name: '',
            start_time: '18:00',
            end_time: '20:00',
            duration_minutes: 120,
            session_type: 'EVENING',
            is_active: true,
            notes: '',
        });
    }, [shift, reset, open]);

    const handleFormSubmit = async (data: ShiftFormValues) => {
        await onSubmit(data);
        onClose();
    };

    return (
        <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
            <DialogTitle>{shift ? 'Chỉnh sửa ca học' : 'Thêm ca học mới'}</DialogTitle>
            <DialogContent dividers>
                <Grid container spacing={2} sx={{ mt: 0.5 }}>
                    <Grid size={{ xs: 12, sm: 4 }}>
                        <TextField
                            fullWidth
                            label="Mã ca"
                            {...register('code')}
                            error={!!errors.code}
                            helperText={errors.code?.message}
                        />
                    </Grid>
                    <Grid size={{ xs: 12, sm: 8 }}>
                        <TextField
                            fullWidth
                            label="Tên ca học"
                            {...register('name')}
                            error={!!errors.name}
                            helperText={errors.name?.message}
                        />
                    </Grid>
                    <Grid size={{ xs: 12, sm: 4 }}>
                        <TextField
                            fullWidth
                            label="Giờ bắt đầu"
                            placeholder="18:00"
                            {...register('start_time')}
                            error={!!errors.start_time}
                            helperText={errors.start_time?.message}
                        />
                    </Grid>
                    <Grid size={{ xs: 12, sm: 4 }}>
                        <TextField
                            fullWidth
                            label="Giờ kết thúc"
                            placeholder="20:00"
                            {...register('end_time')}
                            error={!!errors.end_time}
                            helperText={errors.end_time?.message}
                        />
                    </Grid>
                    <Grid size={{ xs: 12, sm: 4 }}>
                        <TextField
                            fullWidth
                            label="Thời lượng (phút)"
                            type="number"
                            {...register('duration_minutes')}
                            error={!!errors.duration_minutes}
                            helperText={errors.duration_minutes?.message}
                        />
                    </Grid>
                    <Grid size={{ xs: 12, sm: 6 }}>
                        <TextField
                            fullWidth
                            select
                            label="Loại ca"
                            defaultValue="EVENING"
                            {...register('session_type')}
                            error={!!errors.session_type}
                            helperText={errors.session_type?.message}
                        >
                            <MenuItem value="MORNING">Buổi sáng</MenuItem>
                            <MenuItem value="AFTERNOON">Buổi chiều</MenuItem>
                            <MenuItem value="EVENING">Buổi tối</MenuItem>
                            <MenuItem value="CUSTOM">Tùy chỉnh</MenuItem>
                        </TextField>
                    </Grid>
                    <Grid size={{ xs: 12, sm: 6 }} sx={{ display: 'flex', alignItems: 'center' }}>
                        <FormControlLabel
                            control={
                                <Switch
                                    checked={watch('is_active')}
                                    onChange={(_event, checked) => setValue('is_active', checked)}
                                />
                            }
                            label="Ca học đang hoạt động"
                        />
                    </Grid>
                    <Grid size={12}>
                        <TextField
                            fullWidth
                            label="Ghi chú"
                            multiline
                            rows={3}
                            {...register('notes')}
                            error={!!errors.notes}
                            helperText={errors.notes?.message || 'Mô tả ngắn mục đích sử dụng của ca học nếu cần'}
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
                    {shift ? 'Cập nhật' : 'Thêm mới'}
                </Button>
            </DialogActions>
        </Dialog>
    );
};

export default ShiftDialog;
