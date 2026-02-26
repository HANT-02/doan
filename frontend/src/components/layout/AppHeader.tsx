import React, { useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { useLogoutAccountMutation } from '@/api/authApi';
import { useNavigate } from 'react-router-dom';
import {
    AppBar,
    Toolbar,
    IconButton,
    Box,
    Typography,
    Avatar,
    Menu,
    MenuItem,
    ListItemIcon,
    ListItemText,
    useTheme
} from '@mui/material';
import {
    Menu as MenuIcon,
    Logout as LogoutIcon,
    Key as KeyIcon,
    School,
    PersonOutline
} from '@mui/icons-material';

interface AppHeaderProps {
    onMobileOpen: () => void;
}

export const AppHeader: React.FC<AppHeaderProps> = ({ onMobileOpen }) => {
    const { user } = useAuth();
    const [logoutMutation] = useLogoutAccountMutation();
    const navigate = useNavigate();
    const theme = useTheme();
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

    const handleMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
        setAnchorEl(event.currentTarget);
    };

    const handleMenuClose = () => {
        setAnchorEl(null);
    };

    const handleLogout = async () => {
        handleMenuClose();
        try {
            await logoutMutation({}).unwrap();
        } catch (e) { /* Redux slice logout handled in api */ }
        navigate('/login');
    };

    const handleChangePassword = () => {
        handleMenuClose();
        navigate('/app/change-password');
    };

    return (
        <AppBar
            position="sticky"
            elevation={0}
            sx={{
                backgroundColor: 'background.paper',
                borderBottom: '1px solid',
                borderColor: 'divider',
                zIndex: (theme) => theme.zIndex.drawer + 1
            }}
        >
            <Toolbar sx={{ justifyContent: 'space-between', px: { xs: 2, md: 3 } }}>
                <Box sx={{ display: 'flex', alignItems: 'center' }}>
                    <IconButton
                        color="inherit"
                        aria-label="open drawer"
                        edge="start"
                        onClick={onMobileOpen}
                        sx={{ mr: 2, display: { md: 'none' }, color: 'text.primary' }}
                    >
                        <MenuIcon />
                    </IconButton>

                    <Box sx={{ display: { xs: 'none', md: 'flex' }, alignItems: 'center', gap: 1 }}>
                        <School color="primary" sx={{ fontSize: 32 }} />
                        <Typography
                            variant="h6"
                            noWrap
                            sx={{
                                fontWeight: 800,
                                color: 'primary.main',
                                letterSpacing: '-0.5px'
                            }}
                        >
                            EduCenter
                        </Typography>
                    </Box>
                </Box>

                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <Box
                        onClick={handleMenuOpen}
                        sx={{
                            display: 'flex',
                            alignItems: 'center',
                            gap: 1.5,
                            ml: 1,
                            p: 0.5,
                            pr: 1.5,
                            borderRadius: '24px',
                            cursor: 'pointer',
                            transition: 'all 0.2s',
                            '&:hover': {
                                backgroundColor: 'action.hover'
                            }
                        }}
                    >
                        <Avatar
                            sx={{
                                width: 36,
                                height: 36,
                                bgcolor: 'primary.main',
                                fontSize: '0.875rem',
                                fontWeight: 600
                            }}
                        >
                            {user?.full_name?.charAt(0) || 'U'}
                        </Avatar>
                        <Box sx={{ display: { xs: 'none', lg: 'block' } }}>
                            <Typography variant="subtitle2" sx={{ lineHeight: 1.2, fontWeight: 700, color: 'text.primary' }}>
                                {user?.full_name || 'Hệ thống'}
                            </Typography>
                            <Typography variant="caption" sx={{ color: 'text.secondary', textTransform: 'capitalize' }}>
                                {user?.role || 'Quản trị viên'}
                            </Typography>
                        </Box>
                    </Box>

                    <Menu
                        anchorEl={anchorEl}
                        open={Boolean(anchorEl)}
                        onClose={handleMenuClose}
                        onClick={handleMenuClose}
                        PaperProps={{
                            elevation: 0,
                            sx: {
                                mt: 1.5,
                                minWidth: 200,
                                borderRadius: 2,
                                border: '1px solid',
                                borderColor: 'divider',
                                boxShadow: '0px 4px 20px rgba(0,0,0,0.08)',
                                '& .MuiMenuItem-root': {
                                    px: 2,
                                    py: 1,
                                    borderRadius: 1,
                                    mx: 1,
                                    mb: 0.5,
                                },
                            },
                        }}
                        transformOrigin={{ horizontal: 'right', vertical: 'top' }}
                        anchorOrigin={{ horizontal: 'right', vertical: 'bottom' }}
                    >
                        <Box sx={{ px: 2, py: 1.5, mb: 1, borderBottom: '1px solid', borderColor: 'divider' }}>
                            <Typography variant="subtitle2" sx={{ fontWeight: 700 }}>
                                {user?.full_name}
                            </Typography>
                            <Typography variant="caption" color="text.secondary">
                                {user?.email}
                            </Typography>
                        </Box>

                        <MenuItem onClick={() => navigate('/app/profile')}>
                            <ListItemIcon>
                                <PersonOutline fontSize="small" />
                            </ListItemIcon>
                            <ListItemText>Trang cá nhân</ListItemText>
                        </MenuItem>

                        <MenuItem onClick={handleChangePassword}>
                            <ListItemIcon>
                                <KeyIcon fontSize="small" />
                            </ListItemIcon>
                            <ListItemText>Đổi mật khẩu</ListItemText>
                        </MenuItem>

                        <MenuItem onClick={handleLogout} sx={{ color: 'error.main' }}>
                            <ListItemIcon>
                                <LogoutIcon fontSize="small" color="error" />
                            </ListItemIcon>
                            <ListItemText>Đăng xuất</ListItemText>
                        </MenuItem>
                    </Menu>
                </Box>
            </Toolbar>
        </AppBar >
    );
};
