import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { useNavigate, useSearchParams, Link } from 'react-router-dom';
import { authApi } from '@/api/authApi';
import {
    Container,
    Box,
    Typography,
    TextField,
    Button,
    InputAdornment,
    IconButton,
    CircularProgress,
    Alert,
    Stack
} from '@mui/material';
import { Visibility, VisibilityOff, ErrorOutline, School } from '@mui/icons-material';
import { toast } from 'sonner';
import FormCard from '@/components/common/FormCard';

const schema = z.object({
    password: z.string().min(6, 'Mật khẩu phải có ít nhất 6 ký tự'),
    confirm_password: z.string(),
}).refine((data) => data.password === data.confirm_password, {
    message: "Mật khẩu không khớp",
    path: ["confirm_password"],
});

type FormValues = z.infer<typeof schema>;

export const ResetPasswordPage = () => {
    const [searchParams] = useSearchParams();
    const token = searchParams.get('token');
    const navigate = useNavigate();
    const [loading, setLoading] = useState(false);
    const [showPassword, setShowPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);
    const [errorMsg, setErrorMsg] = useState<string | null>(null);

    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<FormValues>({
        resolver: zodResolver(schema),
    });

    const onSubmit = async (data: FormValues) => {
        if (!token) return;
        setLoading(true);
        setErrorMsg(null);
        try {
            await authApi.resetPassword({
                token: token,
                new_password_enc: data.password
            });
            toast.success('Đặt lại mật khẩu thành công');
            navigate('/login');
        } catch (error: any) {
            console.error(error);
            setErrorMsg(error.response?.data?.message || 'Không thể đặt lại mật khẩu. Liên kết có thể đã hết hạn.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <Container maxWidth="xs">
            <Box sx={{ py: 8, display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
                <Stack direction="row" spacing={1} sx={{ mb: 4, alignItems: 'center' }}>
                    <School color="primary" sx={{ fontSize: 40 }} />
                    <Typography variant="h4" component="h1" sx={{ fontWeight: 800, color: 'primary.main' }}>
                        EduCenter
                    </Typography>
                </Stack>

                {!token ? (
                    <FormCard
                        title="Liên kết không hợp lệ"
                        sx={{ width: '100%', textAlign: 'center' }}
                    >
                        <ErrorOutline color="error" sx={{ fontSize: 60, mb: 2 }} />
                        <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                            Liên kết đặt lại mật khẩu này không hợp lệ hoặc đã hết hạn. Vui lòng yêu cầu một liên kết mới.
                        </Typography>
                        <Button
                            fullWidth
                            variant="contained"
                            component={Link}
                            to="/forgot-password"
                        >
                            Yêu cầu liên kết mới
                        </Button>
                    </FormCard>
                ) : (
                    <FormCard
                        title="Đặt lại mật khẩu"
                        subtitle="Tạo mật khẩu mới cho tài khoản của bạn"
                        sx={{ width: '100%' }}
                    >
                        {errorMsg && (
                            <Alert severity="error" sx={{ mb: 3 }}>
                                {errorMsg}
                            </Alert>
                        )}

                        <Box component="form" onSubmit={handleSubmit(onSubmit)} noValidate>
                            <TextField
                                margin="normal"
                                required
                                fullWidth
                                label="Mật khẩu mới"
                                type={showPassword ? 'text' : 'password'}
                                id="password"
                                error={!!errors.password}
                                helperText={errors.password?.message}
                                {...register('password')}
                                InputProps={{
                                    endAdornment: (
                                        <InputAdornment position="end">
                                            <IconButton
                                                aria-label="toggle password visibility"
                                                onClick={() => setShowPassword(!showPassword)}
                                                edge="end"
                                            >
                                                {showPassword ? <VisibilityOff /> : <Visibility />}
                                            </IconButton>
                                        </InputAdornment>
                                    ),
                                }}
                            />
                            <TextField
                                margin="normal"
                                required
                                fullWidth
                                label="Xác nhận mật khẩu mới"
                                type={showConfirmPassword ? 'text' : 'password'}
                                id="confirm_password"
                                error={!!errors.confirm_password}
                                helperText={errors.confirm_password?.message}
                                {...register('confirm_password')}
                                InputProps={{
                                    endAdornment: (
                                        <InputAdornment position="end">
                                            <IconButton
                                                aria-label="toggle password visibility"
                                                onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                                                edge="end"
                                            >
                                                {showConfirmPassword ? <VisibilityOff /> : <Visibility />}
                                            </IconButton>
                                        </InputAdornment>
                                    ),
                                }}
                            />
                            <Button
                                type="submit"
                                fullWidth
                                variant="contained"
                                size="large"
                                sx={{ mt: 4, mb: 2, height: 48 }}
                                disabled={loading}
                            >
                                {loading ? <CircularProgress size={24} color="inherit" /> : 'Đặt lại mật khẩu'}
                            </Button>

                            <Box sx={{ display: 'flex', justifyContent: 'center', mt: 2 }}>
                                <Link to="/login" style={{ textDecoration: 'none' }}>
                                    <Typography variant="body2" color="primary" sx={{ fontWeight: 600 }}>
                                        Quay lại Đăng nhập
                                    </Typography>
                                </Link>
                            </Box>
                        </Box>
                    </FormCard>
                )}
            </Box>
        </Container>
    );
};
