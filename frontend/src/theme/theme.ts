import { createTheme } from '@mui/material/styles';
import type { DataGridComponents } from '@mui/x-data-grid/themeAugmentation';

declare module '@mui/material/styles' {
    interface Palette {
        neutral: Palette['primary'];
    }
    interface PaletteOptions {
        neutral?: PaletteOptions['primary'];
    }
    interface Components extends DataGridComponents { }
}

export const theme = createTheme({
    palette: {
        primary: {
            main: '#2563ea', // Modern blue
            light: '#60a5fa',
            dark: '#1e40af',
            contrastText: '#fff',
        },
        secondary: {
            main: '#4f46e5', // Indigo
            light: '#818cf8',
            dark: '#3730a3',
            contrastText: '#fff',
        },
        neutral: {
            main: '#64748b',
            light: '#94a3b8',
            dark: '#334155',
            contrastText: '#fff',
        },
        background: {
            default: '#f8fafc',
            paper: '#ffffff',
        },
    },
    typography: {
        fontFamily: '"Inter", "Roboto", "Helvetica", "Arial", sans-serif',
        h1: { fontSize: '2.5rem', fontWeight: 600 },
        h2: { fontSize: '2rem', fontWeight: 600 },
        h3: { fontSize: '1.75rem', fontWeight: 600 },
        h4: { fontSize: '1.5rem', fontWeight: 600 },
        h5: { fontSize: '1.25rem', fontWeight: 600 },
        h6: { fontSize: '1rem', fontWeight: 600 },
        subtitle1: { fontSize: '1rem', fontWeight: 500 },
        subtitle2: { fontSize: '0.875rem', fontWeight: 500 },
        body1: { fontSize: '1rem' },
        body2: { fontSize: '0.875rem' },
        button: { textTransform: 'none', fontWeight: 500 },
    },
    shape: {
        borderRadius: 8,
    },
    components: {
        MuiButton: {
            styleOverrides: {
                root: {
                    borderRadius: 8,
                    textTransform: 'none',
                    boxShadow: 'none',
                    '&:hover': {
                        boxShadow: 'none',
                    },
                },
            },
        },
        MuiPaper: {
            styleOverrides: {
                root: {
                    backgroundImage: 'none',
                    boxShadow: '0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1)',
                },
                elevation1: {
                    boxShadow: '0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1)',
                },
            },
        },
        MuiAppBar: {
            styleOverrides: {
                root: {
                    backgroundColor: '#ffffff',
                    color: '#0f172a',
                    boxShadow: '0 1px 2px 0 rgb(0 0 0 / 0.05)',
                },
            },
        },
        MuiTextField: {
            styleOverrides: {
                root: {
                    '& .MuiOutlinedInput-root': {
                        borderRadius: 8,
                    },
                },
            },
        },
        MuiDataGrid: {
            styleOverrides: {
                root: {
                    border: 'none',
                    '& .MuiDataGrid-cell': {
                        borderBottom: '1px solid #f1f5f9',
                    },
                    '& .MuiDataGrid-columnHeaders': {
                        backgroundColor: '#f8fafc',
                        borderBottom: '1px solid #e2e8f0',
                    },
                    '& .MuiDataGrid-row:hover': {
                        backgroundColor: '#f1f5f9',
                    },
                },
            },
        },
    },
});
