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
import type { Room } from '@/api/roomApi';

const roomSchema = z.object({
    name: z.string().min(1, 'Tên phòng không được để trống'),
    capacity: z.coerce.number().min(1, 'Sức chứa phải lớn hơn 0'),
    location: z.string().optional(),
    status: z.enum(['ACTIVE', 'MAINTENANCE', 'INACTIVE']),
});

type RoomFormValues = z.infer<typeof roomSchema>;

interface RoomDialogProps {
    open: boolean;
    onClose: () => void;
    onSubmit: (data: RoomFormValues) => Promise<void>;
    room?: Room | null;
    isLoading?: boolean;
}

const RoomDialog = ({ open, onClose, onSubmit, room, isLoading }: RoomDialogProps) => {
    const {
        register,
        handleSubmit,
        reset,
        formState: { errors },
    } = useForm<RoomFormValues>({
        resolver: zodResolver(roomSchema) as any,
        defaultValues: {
            name: '',
            capacity: 30,
            location: '',
            status: 'ACTIVE',
        },
    });

    useEffect(() => {
        if (room) {
            reset({
                name: room.name,
                capacity: room.capacity,
                location: room.location || '',
                status: room.status,
            });
        } else {
            reset({
                name: '',
                capacity: 30,
                location: '',
                status: 'ACTIVE',
            });
        }
    }, [room, reset, open]);

    const handleFormSubmit = async (data: RoomFormValues) => {
        await onSubmit(data);
        onClose();
    };

    return (
        <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
            <DialogTitle>
                {room ? 'Chỉnh sửa phòng học' : 'Thêm phòng học mới'}
            </DialogTitle>
            <DialogContent dividers>
                <Grid container spacing={2} sx={{ mt: 0.5 }}>
                    <Grid size={12}>
                        <TextField
                            fullWidth
                            label="Tên phòng"
                            {...register('name')}
                            error={!!errors.name}
                            helperText={errors.name?.message}
                        />
                    </Grid>
                    <Grid size={6}>
                        <TextField
                            fullWidth
                            label="Sức chứa"
                            type="number"
                            {...register('capacity')}
                            error={!!errors.capacity}
                            helperText={errors.capacity?.message}
                        />
                    </Grid>
                    <Grid size={6}>
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
                            <MenuItem value="MAINTENANCE">Bảo trì</MenuItem>
                            <MenuItem value="INACTIVE">Ngừng hoạt động</MenuItem>
                        </TextField>
                    </Grid>
                    <Grid size={12}>
                        <TextField
                            fullWidth
                            label="Vị trí"
                            multiline
                            rows={2}
                            {...register('location')}
                            error={!!errors.location}
                            helperText={errors.location?.message}
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
                    {room ? 'Cập nhật' : 'Thêm mới'}
                </Button>
            </DialogActions>
        </Dialog>
    );
};

export default RoomDialog;
