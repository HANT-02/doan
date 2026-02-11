import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { useNavigate } from 'react-router-dom';
import { authApi } from '@/api/authApi';
import {
    Box,
    TextField,
    Button,
    InputAdornment,
    IconButton,
    CircularProgress,
    Alert,
    Stack
} from '@mui/material';
import { Visibility, VisibilityOff } from '@mui/icons-material';
import { toast } from 'sonner';
import FormCard from '@/components/common/FormCard';
import PageHeader from '@/components/common/PageHeader';

const changePasswordSchema = z.object({
    old_password: z.string().min(1, 'Mật khẩu hiện tại là bắt buộc'),
    new_password: z.string().min(6, 'Mật khẩu mới phải có ít nhất 6 ký tự'),
    confirm_password: z.string(),
}).refine((data) => data.new_password === data.confirm_password, {
    message: "Mật khẩu mới không khớp",
    path: ["confirm_password"],
});

type ChangePasswordFormValues = z.infer<typeof changePasswordSchema>;

export const ChangePasswordPage = () => {
    const navigate = useNavigate();
    const [loading, setLoading] = useState(false);
    const [showOldPass, setShowOldPass] = useState(false);
    const [showNewPass, setShowNewPass] = useState(false);
    const [showConfirmPass, setShowConfirmPass] = useState(false);
    const [errorMsg, setErrorMsg] = useState<string | null>(null);

    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<ChangePasswordFormValues>({
        resolver: zodResolver(changePasswordSchema),
    });

    const onSubmit = async (data: ChangePasswordFormValues) => {
        setLoading(true);
        setErrorMsg(null);
        try {
            await authApi.changePassword({
                old_password_enc: data.old_password,
                new_password_enc: data.new_password
            });
            toast.success('Đổi mật khẩu thành công');
            navigate('/app/profile');
        } catch (error: any) {
            console.error(error);
            setErrorMsg(error.response?.data?.message || 'Không thể đổi mật khẩu. Vui lòng kiểm tra lại mật khẩu cũ.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <Box>
            <PageHeader
                title="Đổi mật khẩu"
                subtitle="Cập nhật mật khẩu định kỳ để bảo vệ tài khoản của bạn"
                breadcrumbs={[
                    { label: 'Hệ thống', path: '/app' },
                    { label: 'Cá nhân', path: '/app/profile' },
                    { label: 'Đổi mật khẩu' }
                ]}
            />

            <Box sx={{ maxWidth: 600 }}>
                <FormCard>
                    {errorMsg && (
                        <Alert severity="error" sx={{ mb: 3 }}>
                            {errorMsg}
                        </Alert>
                    )}

                    <Box component="form" onSubmit={handleSubmit(onSubmit)} noValidate>
                        <Stack spacing={2}>
                            <TextField
                                required
                                fullWidth
                                id="old_password"
                                label="Mật khẩu hiện tại"
                                type={showOldPass ? 'text' : 'password'}
                                error={!!errors.old_password}
                                helperText={errors.old_password?.message}
                                {...register('old_password')}
                                InputProps={{
                                    endAdornment: (
                                        <InputAdornment position="end">
                                            <IconButton
                                                onClick={() => setShowOldPass(!showOldPass)}
                                                edge="end"
                                            >
                                                {showOldPass ? <VisibilityOff /> : <Visibility />}
                                            </IconButton>
                                        </InputAdornment>
                                    ),
                                }}
                            />

                            <TextField
                                required
                                fullWidth
                                id="new_password"
                                label="Mật khẩu mới"
                                type={showNewPass ? 'text' : 'password'}
                                error={!!errors.new_password}
                                helperText={errors.new_password?.message}
                                {...register('new_password')}
                                InputProps={{
                                    endAdornment: (
                                        <InputAdornment position="end">
                                            <IconButton
                                                onClick={() => setShowNewPass(!showNewPass)}
                                                edge="end"
                                            >
                                                {showNewPass ? <VisibilityOff /> : <Visibility />}
                                            </IconButton>
                                        </InputAdornment>
                                    ),
                                }}
                            />

                            <TextField
                                required
                                fullWidth
                                id="confirm_password"
                                label="Xác nhận mật khẩu mới"
                                type={showConfirmPass ? 'text' : 'password'}
                                error={!!errors.confirm_password}
                                helperText={errors.confirm_password?.message}
                                {...register('confirm_password')}
                                InputProps={{
                                    endAdornment: (
                                        <InputAdornment position="end">
                                            <IconButton
                                                onClick={() => setShowConfirmPass(!showConfirmPass)}
                                                edge="end"
                                            >
                                                {showConfirmPass ? <VisibilityOff /> : <Visibility />}
                                            </IconButton>
                                        </InputAdornment>
                                    ),
                                }}
                            />
                        </Stack>

                        <Stack direction="row" spacing={2} sx={{ mt: 4 }}>
                            <Button
                                variant="outlined"
                                onClick={() => navigate(-1)}
                                fullWidth
                                sx={{ height: 48 }}
                            >
                                Hủy
                            </Button>
                            <Button
                                type="submit"
                                variant="contained"
                                fullWidth
                                disabled={loading}
                                sx={{ height: 48 }}
                            >
                                {loading ? <CircularProgress size={24} color="inherit" /> : 'Cập nhật mật khẩu'}
                            </Button>
                        </Stack>
                    </Box>
                </FormCard>
            </Box>
        </Box>
    );
};
