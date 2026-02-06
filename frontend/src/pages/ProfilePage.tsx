import { useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
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
    CircularProgress
} from '@mui/material';
import { AccountCircle, Email, VpnKey, Logout as LogoutIcon } from '@mui/icons-material';

export const ProfilePage = () => {
    const { user, logout } = useAuth();
    const navigate = useNavigate();
    const [loading, setLoading] = useState(false);

    // If we wanted to fetch fresh profile data, we could do it here
    // For now, use context user

    const handleLogout = async () => {
        setLoading(true);
        try {
            await logout();
            navigate('/login');
        } catch (e) {
            console.error(e);
        } finally {
            setLoading(false);
        }
    }

    if (!user) {
        return (
            <Container>
                <Box display="flex" justifyContent="center" mt={4}>
                    <CircularProgress />
                </Box>
            </Container>
        )
    }

    return (
        <Container component="main" maxWidth="md">
            <Box sx={{ mt: 4, mb: 4 }}>
                <Paper elevation={3} sx={{ p: 4, borderRadius: 2 }}>
                    <Box sx={{ display: 'flex', flexDirection: { xs: 'column', sm: 'row' }, alignItems: 'center', gap: 3, mb: 4 }}>
                        <Avatar
                            sx={{ width: 100, height: 100, bgcolor: 'primary.main', fontSize: 40 }}
                        >
                            {user.full_name?.charAt(0) || 'U'}
                        </Avatar>
                        <Box sx={{ textAlign: { xs: 'center', sm: 'left' } }}>
                            <Typography variant="h4" gutterBottom>
                                {user.full_name}
                            </Typography>
                            <Chip
                                label={user.role || 'Member'}
                                color={user.role === 'ADMIN' ? 'error' : 'default'}
                                size="small"
                            />
                        </Box>
                    </Box>

                    <Divider sx={{ mb: 2 }} />

                    <Typography variant="h6" gutterBottom color="text.secondary">
                        Account Information
                    </Typography>

                    <List>
                        <ListItem>
                            <Box sx={{ mr: 2, color: 'text.secondary' }}> <Email /> </Box>
                            <ListItemText
                                primary="Email Address"
                                secondary={user.email}
                            />
                        </ListItem>
                        <ListItem>
                            <Box sx={{ mr: 2, color: 'text.secondary' }}> <AccountCircle /> </Box>
                            <ListItemText
                                primary="User ID"
                                secondary={user.id}
                            />
                        </ListItem>
                        <ListItem>
                            <Box sx={{ mr: 2, color: 'text.secondary' }}> <VpnKey /> </Box>
                            <ListItemText
                                primary="User Code"
                                secondary={user.code}
                            />
                        </ListItem>
                    </List>

                    <Box sx={{ mt: 4, display: 'flex', gap: 2, flexWrap: 'wrap', justifyContent: { xs: 'center', sm: 'flex-start' } }}>
                        <Button
                            variant="outlined"
                            startIcon={<VpnKey />}
                            onClick={() => navigate('/change-password')}
                        >
                            Change Password
                        </Button>
                        <Button
                            variant="contained"
                            color="error"
                            startIcon={<LogoutIcon />}
                            onClick={handleLogout}
                            disabled={loading}
                        >
                            {loading ? 'Logging out...' : 'Sign Out'}
                        </Button>
                    </Box>
                </Paper>
            </Box>
        </Container>
    );
};
