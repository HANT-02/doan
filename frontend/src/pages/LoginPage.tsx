import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { useNavigate, useLocation } from 'react-router-dom';
import { useLoginMutation } from '@/api/authApi';
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
import { Visibility, VisibilityOff, School } from '@mui/icons-material';
import FormCard from '@/components/common/FormCard';

const loginSchema = z.object({
    username: z.string().email('Địa chỉ email không hợp lệ'),
    password: z.string().min(6, 'Mật khẩu phải có ít nhất 6 ký tự'),
});

type LoginFormValues = z.infer<typeof loginSchema>;

export const LoginPage = () => {
    const navigate = useNavigate();
    const location = useLocation();

    // Replace Context with RTK Mutation
    const [loginMutation, { isLoading }] = useLoginMutation();

    const [showPassword, setShowPassword] = useState(false);
    const [errorMsg, setErrorMsg] = useState<string | null>(null);

    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<LoginFormValues>({
        resolver: zodResolver(loginSchema),
        defaultValues: {
            username: '',
            password: '',
        }
    });

    const onSubmit = async (data: LoginFormValues) => {
        setErrorMsg(null);
        try {
            const result = await loginMutation(data).unwrap();

            if (result.success) {
                const from = location.state?.from?.pathname || '/app';
                navigate(from, { replace: true });
            }
        } catch (error: any) {
            console.error(error);
            setErrorMsg(error.data?.message || 'Đăng nhập thất bại. Vui lòng kiểm tra lại thông tin.');
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
                    title="Đăng nhập"
                    subtitle="Chào mừng quay trở lại hệ thống quản lý"
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
                            id="email"
                            label="Địa chỉ Email"
                            autoComplete="email"
                            autoFocus
                            error={!!errors.username}
                            helperText={errors.username?.message}
                            {...register('username')}
                        />
                        <TextField
                            margin="normal"
                            required
                            fullWidth
                            label="Mật khẩu"
                            type={showPassword ? 'text' : 'password'}
                            id="password"
                            autoComplete="current-password"
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
                        <Button
                            type="submit"
                            fullWidth
                            variant="contained"
                            size="large"
                            sx={{ mt: 4, mb: 2, height: 48 }}
                            disabled={isLoading}
                        >
                            {isLoading ? <CircularProgress size={24} color="inherit" /> : 'Đăng nhập'}
                        </Button>

                        <Stack direction="row" justifyContent="flex-end" sx={{ mt: 2 }}>
                            {/* Disabled OTP/Email flows for Demo */}
                        </Stack>
                    </Box>
                </FormCard>

                <Typography variant="body2" color="text.secondary" sx={{ mt: 4, textAlign: 'center' }}>
                    &copy; {new Date().getFullYear()} EduCenter Management System.
                </Typography>
            </Box>
        </Container>
    );
};
