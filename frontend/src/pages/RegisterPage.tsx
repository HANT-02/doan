import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { Link, useNavigate } from 'react-router-dom';
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
import { Visibility, VisibilityOff, PersonAdd } from '@mui/icons-material';

const registerSchema = z.object({
    full_name: z.string().min(2, 'Full name must be at least 2 characters'),
    email: z.string().email('Invalid email address'),
    password: z.string().min(6, 'Password must be at least 6 characters'),
    confirm_password: z.string(),
}).refine((data) => data.password === data.confirm_password, {
    message: "Passwords don't match",
    path: ["confirm_password"],
});

type RegisterFormValues = z.infer<typeof registerSchema>;

export const RegisterPage = () => {
    const navigate = useNavigate();
    const [loading, setLoading] = useState(false);
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
        setLoading(true);
        setErrorMsg(null);
        try {
            await authApi.register({
                email: data.email,
                full_name: data.full_name,
                password_enc: data.password // Backend expects already encrypted or plain via TLS, using plain for now as per previous code
            });
            // Show success via Snackbar or Toast? For now just navigate with state or simple alert
            // Ideally we replace sonner with MUI Snackbar, but let's stick to page logic
            navigate('/login');
        } catch (error: any) {
            console.error(error);
            setErrorMsg(error.response?.data?.message || 'Registration failed.');
        } finally {
            setLoading(false);
        }
    };

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
                        bgcolor: 'secondary.main',
                        color: 'white',
                        p: 1.5,
                        borderRadius: '50%',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center'
                    }}>
                        <PersonAdd />
                    </Box>
                    <Typography component="h1" variant="h5" sx={{ mb: 2 }}>
                        Sign Up
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
                            id="full_name"
                            label="Full Name"
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
                            label="Email Address"
                            autoComplete="email"
                            error={!!errors.email}
                            helperText={errors.email?.message}
                            {...register('email')}
                        />
                        <TextField
                            margin="normal"
                            required
                            fullWidth
                            label="Password"
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
                            label="Confirm Password"
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
                            color="secondary"
                            sx={{ mt: 3, mb: 2, py: 1.2 }}
                            disabled={loading}
                        >
                            {loading ? <CircularProgress size={24} color="inherit" /> : 'Sign Up'}
                        </Button>
                        <Box sx={{ display: 'flex', justifyContent: 'center', mt: 2 }}>
                            <Link to="/login" style={{ textDecoration: 'none' }}>
                                <Typography variant="body2" color="primary">
                                    Already have an account? Sign in
                                </Typography>
                            </Link>
                        </Box>
                    </Box>
                </Paper>
            </Box>
        </Container>
    );
};
