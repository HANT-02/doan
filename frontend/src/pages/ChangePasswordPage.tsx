import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { useNavigate } from 'react-router-dom';
import { useChangePasswordMutation, useConfirmChangePasswordOtpMutation } from '@/api/authApi';
import {
    Box,
    TextField,
    Button,
    InputAdornment,
    IconButton,
    CircularProgress,
    Alert,
    Stack,
    Container
} from '@mui/material';
import { Visibility, VisibilityOff, ArrowBack } from '@mui/icons-material';
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

const confirmOtpSchema = z.object({
    otp: z.string().length(6, 'OTP phải gồm đúng 6 chữ số'),
});

type ConfirmOtpFormValues = z.infer<typeof confirmOtpSchema>;

export const ChangePasswordPage = () => {
    const navigate = useNavigate();
    const [changePassword, { isLoading: isRequestingOtp }] = useChangePasswordMutation();
    const [confirmChangePasswordOtp, { isLoading: isConfirmingOtp }] = useConfirmChangePasswordOtpMutation();
    const [showOldPass, setShowOldPass] = useState(false);
    const [showNewPass, setShowNewPass] = useState(false);
    const [showConfirmPass, setShowConfirmPass] = useState(false);
    const [errorMsg, setErrorMsg] = useState<string | null>(null);
    const [otpRequested, setOtpRequested] = useState(false);

    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<ChangePasswordFormValues>({
        resolver: zodResolver(changePasswordSchema),
    });

    const {
        register: registerOtp,
        handleSubmit: handleSubmitOtp,
        formState: { errors: otpErrors },
    } = useForm<ConfirmOtpFormValues>({
        resolver: zodResolver(confirmOtpSchema),
    });

    const onSubmit = async (data: ChangePasswordFormValues) => {
        setErrorMsg(null);
        try {
            await changePassword({
                old_password_enc: data.old_password,
                new_password_enc: data.new_password
            }).unwrap();
            setOtpRequested(true);
            toast.success('Mã OTP đã được gửi tới email của bạn. Vui lòng xác nhận để hoàn tất đổi mật khẩu.');
        } catch (error: any) {
            console.error(error);
            setErrorMsg(error?.data?.message || 'Không thể đổi mật khẩu. Vui lòng kiểm tra lại mật khẩu cũ.');
        }
    };

    const onConfirmOtp = async (data: ConfirmOtpFormValues) => {
        setErrorMsg(null);
        try {
            await confirmChangePasswordOtp({ otp: data.otp }).unwrap();
            toast.success('Đổi mật khẩu thành công');
            navigate('/app/profile');
        } catch (error: any) {
            console.error(error);
            setErrorMsg(error?.data?.message || 'Không thể xác nhận OTP. Vui lòng kiểm tra lại mã.');
        }
    };

    return (
        <Container maxWidth="md" sx={{ py: 4 }}>
            <PageHeader
                title="Đổi mật khẩu"
                subtitle="Cập nhật mật khẩu định kỳ để bảo vệ tài khoản của bạn"
            />

            <Box sx={{ maxWidth: 600, mt: 3 }}>
                <Button
                    startIcon={<ArrowBack />}
                    onClick={() => navigate(-1)}
                    sx={{ mb: 2 }}
                >
                    Quay lại
                </Button>

                <FormCard>
                    {errorMsg && (
                        <Alert severity="error" sx={{ mb: 3 }}>
                            {errorMsg}
                        </Alert>
                    )}

                    {!otpRequested ? (
                        <Box component="form" onSubmit={handleSubmit(onSubmit)} noValidate>
                            <Stack spacing={3}>
                                <TextField
                                    required
                                    fullWidth
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
                                    sx={{ height: 48, borderRadius: 2 }}
                                    disabled={isRequestingOtp}
                                >
                                    Hủy
                                </Button>
                                <Button
                                    type="submit"
                                    variant="contained"
                                    fullWidth
                                    disabled={isRequestingOtp}
                                    sx={{ height: 48, borderRadius: 2 }}
                                >
                                    {isRequestingOtp ? <CircularProgress size={24} color="inherit" /> : 'Gửi OTP xác nhận'}
                                </Button>
                            </Stack>
                        </Box>
                    ) : (
                        <Box component="form" onSubmit={handleSubmitOtp(onConfirmOtp)} noValidate>
                            <Alert severity="info" sx={{ mb: 3 }}>
                                Mã OTP đã được gửi đến email của bạn. Nhập OTP để hoàn tất đổi mật khẩu.
                            </Alert>

                            <TextField
                                required
                                fullWidth
                                label="Mã OTP"
                                error={!!otpErrors.otp}
                                helperText={otpErrors.otp?.message ?? 'OTP có hiệu lực trong 5 phút'}
                                {...registerOtp('otp')}
                            />

                            <Stack direction="row" spacing={2} sx={{ mt: 4 }}>
                                <Button
                                    variant="outlined"
                                    onClick={() => setOtpRequested(false)}
                                    fullWidth
                                    sx={{ height: 48, borderRadius: 2 }}
                                    disabled={isConfirmingOtp}
                                >
                                    Quay lại
                                </Button>
                                <Button
                                    type="submit"
                                    variant="contained"
                                    fullWidth
                                    disabled={isConfirmingOtp}
                                    sx={{ height: 48, borderRadius: 2 }}
                                >
                                    {isConfirmingOtp ? <CircularProgress size={24} color="inherit" /> : 'Xác nhận OTP'}
                                </Button>
                            </Stack>
                        </Box>
                    )}
                </FormCard>
            </Box>
        </Container>
    );
};
