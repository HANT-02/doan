import React from 'react';
import { Outlet } from 'react-router-dom';
import { Box } from '@mui/material';

export const AuthLayout: React.FC = () => {
    return (
        <Box
            sx={{
                minHeight: '100vh',
                background: 'linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%)',
                display: 'flex',
                flexDirection: 'column',
                justifyContent: 'center',
                alignItems: 'center',
                p: 2
            }}
        >
            <Outlet />
        </Box>
    );
};
