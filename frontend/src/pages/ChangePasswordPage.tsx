import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { useNavigate } from 'react-router-dom';
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
import { Visibility, VisibilityOff, Key } from '@mui/icons-material';
import { toast } from 'sonner';

const changePasswordSchema = z.object({
    old_password: z.string().min(1, 'Current password is required'),
    new_password: z.string().min(6, 'Password must be at least 6 characters'),
    confirm_password: z.string(),
}).refine((data) => data.new_password === data.confirm_password, {
    message: "New passwords don't match",
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
            toast.success('Password changed successfully');
            navigate('/profile'); // Redirect to profile instead of home
        } catch (error: any) {
            console.error(error);
            setErrorMsg(error.response?.data?.message || 'Failed to change password');
        } finally {
            setLoading(false);
        }
    };

    return (
        <Container component="main" maxWidth="md">
            <Box sx={{ mt: 4, display: 'flex', justifyContent: 'center' }}>
                <Paper elevation={3} sx={{ p: 4, width: '100%', maxWidth: 500, display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
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
                        <Key />
                    </Box>
                    <Typography component="h1" variant="h5" sx={{ mb: 3 }}>
                        Change Password
                    </Typography>

                    {errorMsg && (
                        <Alert severity="error" sx={{ width: '100%', mb: 2 }}>
                            {errorMsg}
                        </Alert>
                    )}

                    <Box component="form" onSubmit={handleSubmit(onSubmit)} noValidate sx={{ width: '100%' }}>
                        <TextField
                            margin="normal"
                            required
                            fullWidth
                            id="old_password"
                            label="Current Password"
                            type={showOldPass ? 'text' : 'password'}
                            error={!!errors.old_password}
                            helperText={errors.old_password?.message}
                            {...register('old_password')}
                            InputProps={{
                                endAdornment: (
                                    <InputAdornment position="end">
                                        <IconButton
                                            aria-label="toggle password visibility"
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
                            margin="normal"
                            required
                            fullWidth
                            id="new_password"
                            label="New Password"
                            type={showNewPass ? 'text' : 'password'}
                            error={!!errors.new_password}
                            helperText={errors.new_password?.message}
                            {...register('new_password')}
                            InputProps={{
                                endAdornment: (
                                    <InputAdornment position="end">
                                        <IconButton
                                            aria-label="toggle password visibility"
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
                            margin="normal"
                            required
                            fullWidth
                            id="confirm_password"
                            label="Confirm New Password"
                            type={showConfirmPass ? 'text' : 'password'}
                            error={!!errors.confirm_password}
                            helperText={errors.confirm_password?.message}
                            {...register('confirm_password')}
                            InputProps={{
                                endAdornment: (
                                    <InputAdornment position="end">
                                        <IconButton
                                            aria-label="toggle password visibility"
                                            onClick={() => setShowConfirmPass(!showConfirmPass)}
                                            edge="end"
                                        >
                                            {showConfirmPass ? <VisibilityOff /> : <Visibility />}
                                        </IconButton>
                                    </InputAdornment>
                                ),
                            }}
                        />
                        <Box sx={{ display: 'flex', gap: 2, mt: 3 }}>
                            <Button
                                fullWidth
                                variant="outlined"
                                onClick={() => navigate('/profile')}
                                disabled={loading}
                            >
                                Cancel
                            </Button>
                            <Button
                                type="submit"
                                fullWidth
                                variant="contained"
                                disabled={loading}
                            >
                                {loading ? <CircularProgress size={24} color="inherit" /> : 'Update Password'}
                            </Button>
                        </Box>
                    </Box>
                </Paper>
            </Box>
        </Container>
    );
};
