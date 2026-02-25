import { useState, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { Link } from 'react-router-dom';
import { useForgotPasswordMutation } from '@/api/authApi';
import {
    Container,
    Box,
    Typography,
    TextField,
    Button,
    CircularProgress,
    Alert,
    Stack
} from '@mui/material';
import { ArrowBack, School } from '@mui/icons-material';
import FormCard from '@/components/common/FormCard';

const forgotPasswordSchema = z.object({
    email: z.string().email('Địa chỉ email không hợp lệ'),
});

type ForgotPasswordFormValues = z.infer<typeof forgotPasswordSchema>;

export const ForgotPasswordPage = () => {
    const [forgotPassword, { isLoading: loading }] = useForgotPasswordMutation();
    const [success, setSuccess] = useState(false);
    const [cooldown, setCooldown] = useState(0);

    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<ForgotPasswordFormValues>({
        resolver: zodResolver(forgotPasswordSchema),
    });

    useEffect(() => {
        let interval: ReturnType<typeof setInterval>;
        if (cooldown > 0) {
            interval = setInterval(() => {
                setCooldown((prev) => prev - 1);
            }, 1000);
        }
        return () => clearInterval(interval);
    }, [cooldown]);

    const onSubmit = async (data: ForgotPasswordFormValues) => {
        try {
            await forgotPassword(data.email).unwrap();
            setCooldown(60);
            setSuccess(true);
        } catch (error) {
            console.error(error);
            // Even if it fails, we usually show success for security or handle specific errors
            setSuccess(true);
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

                <FormCard
                    title="Quên mật khẩu"
                    subtitle="Nhập email để nhận hướng dẫn đặt lại mật khẩu"
                    sx={{ width: '100%' }}
                >
                    {success ? (
                        <Box>
                            <Alert severity="success" sx={{ mb: 3 }}>
                                Nếu tài khoản tồn tại, chúng tôi đã gửi hướng dẫn đặt lại mật khẩu tới email của bạn.
                            </Alert>
                            {cooldown > 0 && (
                                <Typography variant="body2" color="text.secondary" align="center" sx={{ mb: 2 }}>
                                    Bạn có thể yêu cầu gửi lại sau <strong>{cooldown}s</strong>
                                </Typography>
                            )}
                            <Button
                                fullWidth
                                variant="outlined"
                                onClick={() => setSuccess(false)}
                                disabled={cooldown > 0}
                            >
                                Thử với email khác
                            </Button>
                        </Box>
                    ) : (
                        <Box component="form" onSubmit={handleSubmit(onSubmit)} noValidate>
                            <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                                Một liên kết bảo mật sẽ được gửi đến email của bạn để giúp bạn đặt lại mật khẩu.
                            </Typography>

                            <TextField
                                margin="normal"
                                required
                                fullWidth
                                id="email"
                                label="Địa chỉ Email"
                                autoComplete="email"
                                autoFocus
                                error={!!errors.email}
                                helperText={errors.email?.message}
                                {...register('email')}
                            />
                            <Button
                                type="submit"
                                fullWidth
                                variant="contained"
                                size="large"
                                sx={{ mt: 4, mb: 2, height: 48, borderRadius: 2 }}
                                disabled={loading || cooldown > 0}
                            >
                                {loading ? <CircularProgress size={24} color="inherit" /> : 'Gửi liên kết đặt lại'}
                            </Button>
                        </Box>
                    )}

                    <Box sx={{ display: 'flex', justifyContent: 'center', mt: 3 }}>
                        <Link to="/login" style={{ textDecoration: 'none' }}>
                            <Box sx={{ display: 'flex', alignItems: 'center', color: 'primary.main' }}>
                                <ArrowBack sx={{ fontSize: 18, mr: 1 }} />
                                <Typography variant="body2" sx={{ fontWeight: 600 }}>
                                    Quay lại Đăng nhập
                                </Typography>
                            </Box>
                        </Link>
                    </Box>
                </FormCard>
            </Box>
        </Container>
    );
};
