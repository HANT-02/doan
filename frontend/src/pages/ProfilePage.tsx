import { useNavigate } from 'react-router-dom';
import {
    Container,
    Paper,
    Box,
    Typography,
    Avatar,
    Divider,
    List,
    ListItem,
    ListItemText,
    Button,
    Chip,
} from '@mui/material';
import { AccountCircle, Email, VpnKey, Logout as LogoutIcon } from '@mui/icons-material';
import { useAppSelector, useAppDispatch } from '@/store';
import { logout } from '@/store/authSlice';
import { useLogoutAccountMutation } from '@/api/authApi';
import { toast } from 'sonner';

export const ProfilePage = () => {
    const navigate = useNavigate();
    const dispatch = useAppDispatch();
    const { user } = useAppSelector((state) => state.auth);
    const [logoutAccount, { isLoading: isLoggingOut }] = useLogoutAccountMutation();

    const handleLogout = async () => {
        try {
            await logoutAccount(undefined).unwrap();
            dispatch(logout());
            toast.success('Đã đăng xuất');
            navigate('/login');
        } catch (e) {
            console.error(e);
            dispatch(logout()); // Ensure logout even on fail
            navigate('/login');
        }
    }

    if (!user) {
        return null;
    }

    return (
        <Container component="main" maxWidth="md" sx={{ py: 4 }}>
            <Box>
                <Paper elevation={0} sx={{ p: 4, borderRadius: 3, border: '1px solid #e2e8f0' }}>
                    <Box sx={{ display: 'flex', flexDirection: { xs: 'column', sm: 'row' }, alignItems: 'center', gap: 3, mb: 4 }}>
                        <Avatar
                            sx={{ width: 100, height: 100, bgcolor: 'primary.main', fontSize: 40, fontWeight: 700 }}
                        >
                            {user.full_name?.charAt(0) || 'U'}
                        </Avatar>
                        <Box sx={{ textAlign: { xs: 'center', sm: 'left' } }}>
                            <Typography variant="h4" sx={{ fontWeight: 700 }} gutterBottom>
                                {user.full_name}
                            </Typography>
                            <Chip
                                label={user.role?.toUpperCase() || 'MEMBER'}
                                color={user.role === 'admin' ? 'primary' : 'default'}
                                size="small"
                                sx={{ fontWeight: 600 }}
                            />
                        </Box>
                    </Box>

                    <Divider sx={{ mb: 4 }} />

                    <Typography variant="h6" gutterBottom sx={{ fontWeight: 600, color: 'text.secondary', mb: 2 }}>
                        Thông tin tài khoản
                    </Typography>

                    <List sx={{ mb: 4 }}>
                        <ListItem sx={{ px: 0 }}>
                            <Box sx={{ mr: 2, p: 1, color: 'primary.main', bgcolor: 'primary.lighter', borderRadius: 1.5 }}> <Email /> </Box>
                            <ListItemText
                                primary={<Typography variant="caption" color="text.secondary">Địa chỉ Email</Typography>}
                                secondary={<Typography variant="body1" sx={{ fontWeight: 500 }}>{user.email}</Typography>}
                            />
                        </ListItem>
                        <ListItem sx={{ px: 0 }}>
                            <Box sx={{ mr: 2, p: 1, color: 'primary.main', bgcolor: 'primary.lighter', borderRadius: 1.5 }}> <AccountCircle /> </Box>
                            <ListItemText
                                primary={<Typography variant="caption" color="text.secondary">ID Người dùng</Typography>}
                                secondary={<Typography variant="body1" sx={{ fontWeight: 500 }}>{user.id}</Typography>}
                            />
                        </ListItem>
                        <ListItem sx={{ px: 0 }}>
                            <Box sx={{ mr: 2, p: 1, color: 'primary.main', bgcolor: 'primary.lighter', borderRadius: 1.5 }}> <VpnKey /> </Box>
                            <ListItemText
                                primary={<Typography variant="caption" color="text.secondary">Mã nhân viên/giáo viên</Typography>}
                                secondary={<Typography variant="body1" sx={{ fontWeight: 500 }}>{user.code || 'N/A'}</Typography>}
                            />
                        </ListItem>
                    </List>

                    <Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap' }}>
                        <Button
                            variant="outlined"
                            startIcon={<VpnKey />}
                            onClick={() => navigate('/app/change-password')}
                            sx={{ borderRadius: 2 }}
                        >
                            Đổi mật khẩu
                        </Button>
                        <Button
                            variant="contained"
                            color="error"
                            startIcon={<LogoutIcon />}
                            onClick={handleLogout}
                            disabled={isLoggingOut}
                            sx={{ borderRadius: 2 }}
                        >
                            {isLoggingOut ? 'Đang đăng xuất...' : 'Đăng xuất'}
                        </Button>
                    </Box>
                </Paper>
            </Box>
        </Container>
    );
};
