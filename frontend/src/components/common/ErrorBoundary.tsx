import { Component, type ErrorInfo, type ReactNode } from 'react';
import { Box, Typography, Button, Paper, Alert, AlertTitle } from '@mui/material';
import { Refresh } from '@mui/icons-material';

interface Props {
    children?: ReactNode;
    fallback?: ReactNode;
}

interface State {
    hasError: boolean;
    error: Error | null;
}

class ErrorBoundary extends Component<Props, State> {
    public state: State = {
        hasError: false,
        error: null
    };

    public static getDerivedStateFromError(error: Error): State {
        return { hasError: true, error };
    }

    public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
        console.error('ErrorBoundary caught an error:', error, errorInfo);
    }

    private handleReset = () => {
        this.setState({ hasError: false, error: null });
        window.location.assign('/app/overview');
    };

    public render() {
        if (this.state.hasError) {
            if (this.props.fallback) {
                return this.props.fallback;
            }

            return (
                <Box
                    sx={{
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        p: 4,
                        minHeight: '200px'
                    }}
                >
                    <Paper
                        elevation={0}
                        sx={{
                            p: 4,
                            maxWidth: 600,
                            borderRadius: 4,
                            border: '1px solid #fee2e2',
                            background: '#fffbfb'
                        }}
                    >
                        <Alert severity="error" variant="outlined" sx={{ mb: 3 }}>
                            <AlertTitle sx={{ fontWeight: 700 }}>Đã xảy ra lỗi hệ thống</AlertTitle>
                            Không thể hiển thị nội dung này do lỗi xử lý dữ liệu.
                        </Alert>
                        <Typography variant="body2" color="text.secondary" sx={{ mb: 3, fontFamily: 'monospace' }}>
                            {this.state.error?.message}
                        </Typography>
                        <Button
                            variant="contained"
                            color="primary"
                            startIcon={<Refresh />}
                            onClick={this.handleReset}
                            sx={{ borderRadius: 2 }}
                        >
                            Thử tải lại trang
                        </Button>
                    </Paper>
                </Box>
            );
        }

        return this.props.children;
    }
}

export default ErrorBoundary;
