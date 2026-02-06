import { useState, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { Link } from 'react-router-dom';
import { authApi } from '@/api/authApi';
import {
    Container,
    Box,
    Typography,
    TextField,
    Button,
    Paper,
    CircularProgress,
    Alert
} from '@mui/material';
import { LockReset, ArrowBack } from '@mui/icons-material';

const forgotPasswordSchema = z.object({
    email: z.string().email('Invalid email address'),
});

type ForgotPasswordFormValues = z.infer<typeof forgotPasswordSchema>;

export const ForgotPasswordPage = () => {
    const [loading, setLoading] = useState(false);
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
        let interval: any;
        if (cooldown > 0) {
            interval = setInterval(() => {
                setCooldown((prev) => prev - 1);
            }, 1000);
        }
        return () => clearInterval(interval);
    }, [cooldown]);

    const onSubmit = async (data: ForgotPasswordFormValues) => {
        setLoading(true);
        try {
            await authApi.forgotPassword(data.email);
            // Start cooldown
            setCooldown(60);
        } catch (error) {
            console.error(error);
            // Always show success to prevent enumeration
        } finally {
            setLoading(false);
            setSuccess(true);
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
                        bgcolor: 'warning.main',
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
                        Forgot Password
                    </Typography>

                    <Typography variant="body2" color="text.secondary" align="center" sx={{ mb: 3 }}>
                        Enter your email address and we'll send you a link to reset your password.
                    </Typography>

                    {success ? (
                        <Alert severity="success" sx={{ width: '100%', mb: 2 }}>
                            If an account exists for that email, we have sent password reset instructions.
                            {cooldown > 0 && <Box component="span" display="block" mt={1} fontSize="0.875rem">Resend available in {cooldown}s</Box>}
                        </Alert>
                    ) : (
                        <Box component="form" onSubmit={handleSubmit(onSubmit)} noValidate sx={{ mt: 1, width: '100%' }}>
                            <TextField
                                margin="normal"
                                required
                                fullWidth
                                id="email"
                                label="Email Address"
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
                                color="primary"
                                sx={{ mt: 3, mb: 2, py: 1.2 }}
                                disabled={loading || cooldown > 0}
                            >
                                {loading ? <CircularProgress size={24} color="inherit" /> : 'Send Reset Link'}
                            </Button>
                        </Box>
                    )}
                    <Box sx={{ display: 'flex', justifyContent: 'center', mt: 2 }}>
                        <Link to="/login" style={{ textDecoration: 'none' }}>
                            <Box sx={{ display: 'flex', alignItems: 'center' }}>
                                <ArrowBack sx={{ fontSize: 16, mr: 0.5 }} color="primary" />
                                <Typography variant="body2" color="primary">
                                    Back to Login
                                </Typography>
                            </Box>
                        </Link>
                    </Box>
                </Paper>
            </Box>
        </Container>
    );
};
