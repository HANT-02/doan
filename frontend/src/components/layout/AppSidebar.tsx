import React, { useState, useEffect, useMemo } from 'react';
import { useLocation } from 'react-router-dom';
import { useAuth } from '@/contexts/AuthContext';
import { getNavItemsByRole, type NavItem } from '@/config/nav';
import {
    Drawer,
    List,
    ListItemButton,
    ListItemIcon,
    ListItemText,
    Box,
    Typography,
    IconButton,
    useTheme,
    useMediaQuery,
    Divider,
    Toolbar,
    Tooltip,
    alpha,
    Collapse
} from '@mui/material';
import {
    ChevronRight as ChevronRightIcon,
    ChevronLeft as ChevronLeftIcon,
    ExpandLess as ExpandLessIcon,
    ExpandMore as ExpandMoreIcon,
    Close as CloseIcon,
    School
} from '@mui/icons-material';

interface AppSidebarProps {
    isMobileOpen: boolean;
    onMobileClose: () => void;
}

const drawerWidth = 260;
const collapsedWidth = 80;

export const AppSidebar: React.FC<AppSidebarProps> = ({ isMobileOpen, onMobileClose }) => {
    const { user } = useAuth();
    const location = useLocation();
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('md'));
    const [collapsed, setCollapsed] = useState(false);
    const [openGroups, setOpenGroups] = useState<Record<string, boolean>>({});

    const navItems = useMemo(() => getNavItemsByRole(user?.role), [user?.role]);

    const isItemActive = (item: NavItem): boolean => {
        if (item.path && location.pathname.startsWith(item.path)) {
            return true;
        }

        return item.children?.some((child) => isItemActive(child)) ?? false;
    };

    useEffect(() => {
        if (isMobile) {
            onMobileClose();
        }
    }, [location.pathname, isMobile, onMobileClose]);

    useEffect(() => {
        setOpenGroups((prev) => {
            let hasChanges = false;
            const next = { ...prev };

            navItems.forEach((item) => {
                if (item.children?.length && isItemActive(item) && !prev[item.key]) {
                    next[item.key] = true;
                    hasChanges = true;
                }
            });

            return hasChanges ? next : prev;
        });
    }, [location.pathname, navItems]);

    const renderNavItem = (item: NavItem, depth = 0) => {
        const Icon = item.icon;
        const isActive = isItemActive(item);
        const hasChildren = !!item.children?.length;
        const isGroupOpen = openGroups[item.key] ?? isActive;

        if (hasChildren) {
            const itemContent = (
                <ListItemButton
                    onClick={() => {
                        if (collapsed && !isMobile) {
                            setCollapsed(false);
                            setOpenGroups((prev) => ({ ...prev, [item.key]: true }));
                            return;
                        }

                        setOpenGroups((prev) => ({ ...prev, [item.key]: !isGroupOpen }));
                    }}
                    sx={{
                        borderRadius: 2,
                        mb: 1,
                        py: 1.2,
                        px: collapsed && !isMobile ? 0 : 2,
                        justifyContent: collapsed && !isMobile ? 'center' : 'initial',
                        bgcolor: isActive ? alpha(theme.palette.primary.main, 0.06) : 'transparent',
                        '&:hover': {
                            bgcolor: alpha(theme.palette.action.hover, 0.5),
                        },
                    }}
                >
                    <ListItemIcon
                        sx={{
                            minWidth: 0,
                            mr: collapsed && !isMobile ? 0 : 2,
                            justifyContent: 'center',
                            color: isActive ? 'primary.main' : 'text.secondary',
                        }}
                    >
                        <Icon sx={{ fontSize: 22 }} />
                    </ListItemIcon>
                    {(!collapsed || isMobile) && (
                        <>
                            <ListItemText
                                primary={item.labelVi}
                                primaryTypographyProps={{
                                    variant: 'body2',
                                    sx: { fontWeight: isActive ? 700 : 600 },
                                }}
                            />
                            {isGroupOpen ? <ExpandLessIcon fontSize="small" /> : <ExpandMoreIcon fontSize="small" />}
                        </>
                    )}
                </ListItemButton>
            );

            return (
                <React.Fragment key={item.key}>
                    {collapsed && !isMobile ? (
                        <Tooltip title={item.labelVi} placement="right">
                            {itemContent}
                        </Tooltip>
                    ) : (
                        itemContent
                    )}
                    {(!collapsed || isMobile) && (
                        <Collapse in={isGroupOpen} timeout="auto" unmountOnExit>
                            <List disablePadding sx={{ mb: 1 }}>
                                {item.children?.map((child) => renderNavItem(child, depth + 1))}
                            </List>
                        </Collapse>
                    )}
                </React.Fragment>
            );
        }

        const itemContent = (
            <ListItemButton
                component="a"
                href={item.path || '/app'}
                selected={isActive}
                onClick={() => {
                    if (isMobile) {
                        onMobileClose();
                    }
                }}
                sx={{
                    borderRadius: 2,
                    mb: 1,
                    py: 1.1,
                    px: collapsed && !isMobile ? 0 : 2,
                    pl: collapsed && !isMobile ? 0 : 2 + depth * 2,
                    justifyContent: collapsed && !isMobile ? 'center' : 'initial',
                    transition: 'all 0.2s',
                    '&.Mui-selected': {
                        bgcolor: alpha(theme.palette.primary.main, 0.08),
                        color: 'primary.main',
                        '& .MuiListItemIcon-root': {
                            color: 'primary.main',
                        },
                        '& .MuiTypography-root': {
                            fontWeight: 700,
                        }
                    },
                    '&:hover': {
                        bgcolor: alpha(theme.palette.action.hover, 0.5),
                    }
                }}
            >
                <ListItemIcon
                    sx={{
                        minWidth: 0,
                        mr: collapsed && !isMobile ? 0 : 2,
                        justifyContent: 'center',
                        color: isActive ? 'primary.main' : 'text.secondary'
                    }}
                >
                    <Icon sx={{ fontSize: depth > 0 ? 20 : 22 }} />
                </ListItemIcon>
                {(!collapsed || isMobile) && (
                    <ListItemText
                        primary={item.labelVi}
                        primaryTypographyProps={{
                            variant: 'body2',
                            sx: { fontWeight: isActive ? 700 : depth > 0 ? 500 : 600 }
                        }}
                    />
                )}
            </ListItemButton>
        );

        if (collapsed && !isMobile) {
            return (
                <Tooltip key={item.key} title={item.labelVi} placement="right">
                    {itemContent}
                </Tooltip>
            );
        }

        return <React.Fragment key={item.key}>{itemContent}</React.Fragment>;
    };

    const drawerContent = (
        <Box sx={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
            <Toolbar sx={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: collapsed && !isMobile ? 'center' : 'flex-start',
                px: 2,
                gap: 1.5,
                minHeight: '70px !important'
            }}>
                <School color="primary" sx={{ fontSize: 32 }} />
                {(!collapsed || isMobile) && (
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
                )}
                {isMobile && (
                    <IconButton onClick={onMobileClose} sx={{ ml: 'auto' }}>
                        <CloseIcon />
                    </IconButton>
                )}
            </Toolbar>

            <List sx={{ flexGrow: 1, px: 2, py: 2 }}>
                {navItems.map((item) => renderNavItem(item))}
            </List>

            <Divider sx={{ opacity: 0.6 }} />

            <Box sx={{ p: 2 }}>
                <ListItemButton
                    onClick={() => setCollapsed(!collapsed)}
                    sx={{
                        borderRadius: 2,
                        justifyContent: collapsed && !isMobile ? 'center' : 'initial',
                        color: 'text.secondary',
                        py: 1,
                        px: collapsed && !isMobile ? 0 : 2,
                    }}
                >
                    <ListItemIcon sx={{
                        minWidth: 0,
                        mr: collapsed && !isMobile ? 0 : 2,
                        justifyContent: 'center'
                    }}>
                        {collapsed ? <ChevronRightIcon /> : <ChevronLeftIcon />}
                    </ListItemIcon>
                    {(!collapsed || isMobile) && (
                        <ListItemText
                            primary="Thu gọn"
                            primaryTypographyProps={{ variant: 'body2', sx: { fontWeight: 500 } }}
                        />
                    )}
                </ListItemButton>
            </Box>
        </Box>
    );

    return (
        <Box
            component="nav"
            sx={{
                width: { md: collapsed ? collapsedWidth : drawerWidth },
                flexShrink: { md: 0 },
                transition: theme.transitions.create('width', {
                    easing: theme.transitions.easing.sharp,
                    duration: theme.transitions.duration.shorter,
                }),
            }}
        >
            <Drawer
                variant={isMobile ? "temporary" : "permanent"}
                open={isMobile ? isMobileOpen : true}
                onClose={onMobileClose}
                ModalProps={{ keepMounted: true }}
                sx={{
                    '& .MuiDrawer-paper': {
                        boxSizing: 'border-box',
                        width: isMobile ? drawerWidth : (collapsed ? collapsedWidth : drawerWidth),
                        transition: theme.transitions.create('width', {
                            easing: theme.transitions.easing.sharp,
                            duration: theme.transitions.duration.shorter,
                        }),
                        borderRight: '1px solid',
                        borderColor: 'divider',
                        backgroundColor: 'background.paper',
                        backgroundImage: 'none',
                    },
                }}
            >
                {drawerContent}
            </Drawer>
        </Box>
    );
};
