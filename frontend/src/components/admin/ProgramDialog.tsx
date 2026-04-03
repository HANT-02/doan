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

import type { Program } from '@/api/programApi';

const programSchema = z.object({
    code: z.string().min(1, 'Mã chương trình không được để trống'),
    name: z.string().min(1, 'Tên chương trình không được để trống'),
    track: z.enum(['BASIC', 'ADVANCED', 'SUPPORT']),
    effective_from: z.string().optional(),
    effective_to: z.string().optional(),
    approval_note: z.string().optional(),
});

type ProgramFormValues = z.infer<typeof programSchema>;

interface ProgramDialogProps {
    open: boolean;
    onClose: () => void;
    onSubmit: (data: Partial<Program>) => Promise<void>;
    program?: Program | null;
    isLoading?: boolean;
}

const ProgramDialog = ({ open, onClose, onSubmit, program, isLoading }: ProgramDialogProps) => {
    const {
        register,
        handleSubmit,
        reset,
        formState: { errors },
    } = useForm<ProgramFormValues>({
        resolver: zodResolver(programSchema) as any,
        defaultValues: {
            code: '',
            name: '',
            track: 'BASIC',
            effective_from: '',
            effective_to: '',
            approval_note: '',
        },
    });

    useEffect(() => {
        if (program) {
            reset({
                code: program.code || '',
                name: program.name || '',
                track: (program.track as ProgramFormValues['track']) || 'BASIC',
                effective_from: program.effective_from ? program.effective_from.split('T')[0] : '',
                effective_to: program.effective_to ? program.effective_to.split('T')[0] : '',
                approval_note: program.approval_note || '',
            });
            return;
        }

        reset({
            code: '',
            name: '',
            track: 'BASIC',
            effective_from: '',
            effective_to: '',
            approval_note: '',
        });
    }, [open, program, reset]);

    const handleFormSubmit = async (values: ProgramFormValues) => {
        await onSubmit({
            code: values.code,
            name: values.name,
            track: values.track,
            effective_from: values.effective_from ? `${values.effective_from}T00:00:00Z` : undefined,
            effective_to: values.effective_to ? `${values.effective_to}T00:00:00Z` : undefined,
            approval_note: values.approval_note || '',
        });
        onClose();
    };

    return (
        <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
            <DialogTitle>{program ? 'Chỉnh sửa chương trình' : 'Thêm chương trình mới'}</DialogTitle>
            <DialogContent dividers>
                <Grid container spacing={2} sx={{ mt: 0.5 }}>
                    <Grid size={{ xs: 12, md: 4 }}>
                        <TextField
                            fullWidth
                            label="Mã chương trình"
                            {...register('code')}
                            error={!!errors.code}
                            helperText={errors.code?.message}
                        />
                    </Grid>
                    <Grid size={{ xs: 12, md: 8 }}>
                        <TextField
                            fullWidth
                            label="Tên chương trình"
                            {...register('name')}
                            error={!!errors.name}
                            helperText={errors.name?.message}
                        />
                    </Grid>
                    <Grid size={{ xs: 12, md: 4 }}>
                        <TextField
                            fullWidth
                            select
                            label="Hệ đào tạo"
                            defaultValue="BASIC"
                            {...register('track')}
                            error={!!errors.track}
                            helperText={errors.track?.message}
                        >
                            <MenuItem value="BASIC">Cơ bản</MenuItem>
                            <MenuItem value="ADVANCED">Nâng cao</MenuItem>
                            <MenuItem value="SUPPORT">Bổ trợ</MenuItem>
                        </TextField>
                    </Grid>
                    <Grid size={{ xs: 12, md: 4 }}>
                        <TextField
                            fullWidth
                            type="date"
                            label="Hiệu lực từ"
                            InputLabelProps={{ shrink: true }}
                            {...register('effective_from')}
                            error={!!errors.effective_from}
                            helperText={errors.effective_from?.message}
                        />
                    </Grid>
                    <Grid size={{ xs: 12, md: 4 }}>
                        <TextField
                            fullWidth
                            type="date"
                            label="Hiệu lực đến"
                            InputLabelProps={{ shrink: true }}
                            {...register('effective_to')}
                            error={!!errors.effective_to}
                            helperText={errors.effective_to?.message}
                        />
                    </Grid>
                    <Grid size={12}>
                        <TextField
                            fullWidth
                            multiline
                            rows={3}
                            label="Ghi chú phê duyệt"
                            {...register('approval_note')}
                            error={!!errors.approval_note}
                            helperText={errors.approval_note?.message || 'Có thể để trống nếu chưa có ghi chú'}
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
                    {program ? 'Cập nhật' : 'Thêm mới'}
                </Button>
            </DialogActions>
        </Dialog>
    );
};

export default ProgramDialog;
