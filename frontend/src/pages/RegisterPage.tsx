import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { Link, useNavigate } from 'react-router-dom';
import { useRegisterMutation } from '@/api/authApi';
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
import { toast } from 'sonner';

const registerSchema = z.object({
    full_name: z.string().min(2, 'Họ và tên ít nhất 2 ký tự'),
    email: z.string().email('Địa chỉ email không hợp lệ'),
    password: z.string().min(6, 'Mật khẩu phải có ít nhất 6 ký tự'),
    confirm_password: z.string(),
}).refine((data) => data.password === data.confirm_password, {
    message: "Mật khẩu không khớp",
    path: ["confirm_password"],
});

type RegisterFormValues = z.infer<typeof registerSchema>;

export const RegisterPage = () => {
    const navigate = useNavigate();
    const [registerUser, { isLoading: loading }] = useRegisterMutation();
    const [showPassword, setShowPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);
    const [errorMsg, setErrorMsg] = useState<string | null>(null);

    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<RegisterFormValues>({
        resolver: zodResolver(registerSchema),
    });

    const onSubmit = async (data: RegisterFormValues) => {
        setErrorMsg(null);
        try {
            await registerUser({
                email: data.email,
                full_name: data.full_name,
                password_enc: data.password
            }).unwrap();
            toast.success('Đăng ký tài khoản thành công! Vui lòng đăng nhập.');
            navigate('/login');
        } catch (error: any) {
            console.error(error);
            setErrorMsg(error?.data?.message || 'Đăng ký thất bại. Vui lòng thử lại sau.');
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
                    title="Đăng ký tài khoản"
                    subtitle="Bắt đầu quản lý trung tâm của bạn"
                    sx={{ width: '100%' }}
                >
                    {errorMsg && (
                        <Alert severity="error" sx={{ width: '100%', mb: 2 }}>
                            {errorMsg}
                        </Alert>
                    )}

                    <Box component="form" onSubmit={handleSubmit(onSubmit)} noValidate sx={{ mt: 1, width: '100%' }}>
                        <TextField
                            margin="normal"
                            required
                            fullWidth
                            id="full_name"
                            label="Họ và tên"
                            autoFocus
                            error={!!errors.full_name}
                            helperText={errors.full_name?.message}
                            {...register('full_name')}
                        />
                        <TextField
                            margin="normal"
                            required
                            fullWidth
                            id="email"
                            label="Địa chỉ Email"
                            autoComplete="email"
                            error={!!errors.email}
                            helperText={errors.email?.message}
                            {...register('email')}
                        />
                        <TextField
                            margin="normal"
                            required
                            fullWidth
                            label="Mật khẩu"
                            type={showPassword ? 'text' : 'password'}
                            id="password"
                            autoComplete="new-password"
                            error={!!errors.password}
                            helperText={errors.password?.message}
                            {...register('password')}
                            InputProps={{
                                endAdornment: (
                                    <InputAdornment position="end">
                                        <IconButton
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
                            label="Xác nhận mật khẩu"
                            type={showConfirmPassword ? 'text' : 'password'}
                            id="confirm_password"
                            error={!!errors.confirm_password}
                            helperText={errors.confirm_password?.message}
                            {...register('confirm_password')}
                            InputProps={{
                                endAdornment: (
                                    <InputAdornment position="end">
                                        <IconButton
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
                            sx={{ mt: 3, mb: 2, height: 48, borderRadius: 2 }}
                            disabled={loading}
                        >
                            {loading ? <CircularProgress size={24} color="inherit" /> : 'Đăng ký ngay'}
                        </Button>
                        <Box sx={{ display: 'flex', justifyContent: 'center', mt: 2 }}>
                            <Link to="/login" style={{ textDecoration: 'none' }}>
                                <Typography variant="body2" color="primary" sx={{ fontWeight: 600 }}>
                                    Đã có tài khoản? Đăng nhập
                                </Typography>
                            </Link>
                        </Box>
                    </Box>
                </FormCard>
            </Box>
        </Container>
    );
};
