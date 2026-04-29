import { useState } from 'react';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';
import { useVerifyOtpMutation } from '@/api/authApi';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import * as z from 'zod';
import {
    Alert,
    Box,
    Button,
    CircularProgress,
    Container,
    Stack,
    TextField,
    Typography,
} from '@mui/material';
import { School } from '@mui/icons-material';
import { toast } from 'sonner';
import FormCard from '@/components/common/FormCard';

const verifyOtpSchema = z.object({
    otp: z.string().length(6, 'OTP phải gồm đúng 6 chữ số'),
});

type VerifyOtpFormValues = z.infer<typeof verifyOtpSchema>;

export const VerifyOTPPage = () => {
    const navigate = useNavigate();
    const [searchParams] = useSearchParams();
    const userId = searchParams.get('userId') ?? '';
    const email = searchParams.get('email') ?? '';
    const [verifyOtp, { isLoading }] = useVerifyOtpMutation();
    const [errorMsg, setErrorMsg] = useState<string | null>(null);

    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<VerifyOtpFormValues>({
        resolver: zodResolver(verifyOtpSchema),
    });

    const onSubmit = async (data: VerifyOtpFormValues) => {
        if (!userId) {
            setErrorMsg('Thiếu thông tin đăng ký. Vui lòng đăng ký lại để nhận OTP mới.');
            return;
        }
        setErrorMsg(null);
        try {
            await verifyOtp({
                user_id: userId,
                otp: data.otp,
            }).unwrap();
            toast.success('Xác thực tài khoản thành công. Bạn có thể đăng nhập ngay bây giờ.');
            navigate('/login');
        } catch (error: any) {
            console.error(error);
            setErrorMsg(error?.data?.message || 'Xác thực OTP thất bại. Vui lòng kiểm tra lại mã.');
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
                    title="Xác thực tài khoản"
                    subtitle="Nhập mã OTP đã được gửi tới email của bạn để hoàn tất đăng ký"
                    sx={{ width: '100%' }}
                >
                    {!userId && (
                        <Alert severity="warning" sx={{ mb: 3 }}>
                            Không tìm thấy thông tin người dùng chờ xác thực. Vui lòng quay lại màn đăng ký.
                        </Alert>
                    )}

                    {email && (
                        <Alert severity="info" sx={{ mb: 3 }}>
                            OTP đã được gửi tới <strong>{email}</strong>.
                        </Alert>
                    )}

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
                            id="otp"
                            label="Mã OTP"
                            autoFocus
                            error={!!errors.otp}
                            helperText={errors.otp?.message ?? 'Mã OTP có hiệu lực trong 5 phút'}
                            {...register('otp')}
                        />

                        <Button
                            type="submit"
                            fullWidth
                            variant="contained"
                            size="large"
                            sx={{ mt: 4, mb: 2, height: 48, borderRadius: 2 }}
                            disabled={isLoading || !userId}
                        >
                            {isLoading ? <CircularProgress size={24} color="inherit" /> : 'Xác thực OTP'}
                        </Button>

                        <Box sx={{ display: 'flex', justifyContent: 'space-between', mt: 2 }}>
                            <Link to="/register" style={{ textDecoration: 'none' }}>
                                <Typography variant="body2" color="primary" sx={{ fontWeight: 600 }}>
                                    Quay lại đăng ký
                                </Typography>
                            </Link>
                            <Link to="/login" style={{ textDecoration: 'none' }}>
                                <Typography variant="body2" color="primary" sx={{ fontWeight: 600 }}>
                                    Về đăng nhập
                                </Typography>
                            </Link>
                        </Box>
                    </Box>
                </FormCard>
            </Box>
        </Container>
    );
};
