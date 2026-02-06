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
    Paper,
    InputAdornment,
    IconButton,
    CircularProgress,
    Alert
} from '@mui/material';
import { Visibility, VisibilityOff, LockReset, ErrorOutline } from '@mui/icons-material';
import { toast } from 'sonner';

const schema = z.object({
    password: z.string().min(6, 'Password must be at least 6 characters'),
    confirm_password: z.string(),
}).refine((data) => data.password === data.confirm_password, {
    message: "Passwords don't match",
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
                new_password_enc: data.password // Backend handles plain/enc
            });
            toast.success('Password reset successfully');
            navigate('/login');
        } catch (error: any) {
            console.error(error);
            setErrorMsg(error.response?.data?.message || 'Failed to reset password. The link may be expired.');
        } finally {
            setLoading(false);
        }
    };

    if (!token) {
        return (
            <Container component="main" maxWidth="xs">
                <Box sx={{ mt: 8, display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
                    <Paper elevation={3} sx={{ p: 4, width: '100%', display: 'flex', flexDirection: 'column', alignItems: 'center', textAlign: 'center' }}>
                        <ErrorOutline color="error" sx={{ fontSize: 60, mb: 2 }} />
                        <Typography variant="h6" gutterBottom color="error">
                            Invalid Link
                        </Typography>
                        <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                            This password reset link is invalid or has expired.
                        </Typography>
                        <Link to="/forgot-password" style={{ textDecoration: 'none' }}>
                            <Button variant="outlined">Request a new one</Button>
                        </Link>
                    </Paper>
                </Box>
            </Container>
        );
    }

    return (
        <Container component="main" maxWidth="xs">
            <Box
                sx={{
                    marginTop: 8,
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: 'center',
                }}
            >
                <Paper elevation={3} sx={{ p: 4, width: '100%', display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
                    <Box sx={{
                        m: 1,
                        bgcolor: 'primary.main',
                        color: 'white',
                        p: 1.5,
                        borderRadius: '50%',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center'
                    }}>
                        <LockReset />
                    </Box>
                    <Typography component="h1" variant="h5" sx={{ mb: 2 }}>
                        Reset Password
                    </Typography>

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
                            label="New Password"
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
                            label="Confirm New Password"
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
                            sx={{ mt: 3, mb: 2, py: 1.2 }}
                            disabled={loading}
                        >
                            {loading ? <CircularProgress size={24} color="inherit" /> : 'Reset Password'}
                        </Button>
                        <Box sx={{ display: 'flex', justifyContent: 'center', mt: 2 }}>
                            <Link to="/login" style={{ textDecoration: 'none' }}>
                                <Typography variant="body2" color="primary">
                                    Back to Login
                                </Typography>
                            </Link>
                        </Box>
                    </Box>
                </Paper>
            </Box>
        </Container>
    );
};
