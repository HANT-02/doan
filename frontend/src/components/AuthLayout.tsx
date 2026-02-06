import React from 'react';
import { Outlet } from 'react-router-dom';
import { Toaster } from 'sonner';
import { Box, CssBaseline, ThemeProvider, createTheme } from '@mui/material';

// Create a default theme
const defaultTheme = createTheme({
    palette: {
        mode: 'light',
        primary: {
            main: '#1976d2',
        },
        secondary: {
            main: '#dc004e',
        },
    },
});

export const AuthLayout: React.FC = () => {
    return (
        <ThemeProvider theme={defaultTheme}>
            <CssBaseline />
            <Box
                sx={{
                    minHeight: '100vh',
                    bgcolor: 'grey.100',
                    display: 'flex',
                    flexDirection: 'column',
                    // alignItems: 'center', // Let pages control alignment 
                    // justifyContent: 'center'
                }}
            >
                <Outlet />
                <Toaster richColors />
            </Box>
        </ThemeProvider>
    );
};
