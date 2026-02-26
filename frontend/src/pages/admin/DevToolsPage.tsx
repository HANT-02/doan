import { useState } from 'react';
import {
    Container,
    Typography,
    Box,
    Paper,
    Button,
    Grid,
    Alert,
    AlertTitle,
} from '@mui/material';
import { useAppSelector } from '@/store';
import PageHeader from '@/components/common/PageHeader';
import { baseApi } from '@/api/baseApi';
import { useAppDispatch } from '@/store';

interface PingResult {
    module: string;
    status: 'success' | 'error' | 'pending';
    latency?: number;
    errorText?: string;
}

export const DevToolsPage = () => {
    const { user, isAuthenticated, role } = useAppSelector(state => state.auth);
    const dispatch = useAppDispatch();

    const [results, setResults] = useState<Record<string, PingResult>>({});
    const [lastError, setLastError] = useState<any>(null);

    const pingModule = async (module: string, endpoint: string) => {
        setResults(prev => ({
            ...prev,
            [module]: { module, status: 'pending' }
        }));

        const startTime = Date.now();
        try {
            // Using RTK Query's fetchBaseQuery via a manual dispatch
            const result = await dispatch(
                baseApi.endpoints.pingCustom?.initiate?.({ endpoint }) ||
                baseApi.internalActions.internal_fetchBaseQuery({
                    url: endpoint,
                    method: 'GET'
                })
            );

            const latency = Date.now() - startTime;

            if (result.error) {
                setResults(prev => ({
                    ...prev,
                    [module]: { module, status: 'error', latency, errorText: JSON.stringify(result.error) }
                }));
                setLastError(result.error);
            } else {
                setResults(prev => ({
                    ...prev,
                    [module]: { module, status: 'success', latency }
                }));
            }
        } catch (err: any) {
            const latency = Date.now() - startTime;
            setResults(prev => ({
                ...prev,
                [module]: { module, status: 'error', latency, errorText: err.message || 'Unknown error' }
            }));
            setLastError(err);
        }
    };

    const modulesToTest = [
        { name: 'Teachers', endpoint: '/v1/teachers' },
        { name: 'Rooms', endpoint: '/v1/rooms' },
        { name: 'Classes', endpoint: '/v1/classes' },
        { name: 'Students', endpoint: '/v1/students' },
        { name: 'Courses', endpoint: '/v1/courses' },
        { name: 'Programs', endpoint: '/v1/programs' },
    ];

    return (
        <Container maxWidth="xl" sx={{ py: 4 }}>
            <PageHeader
                title="Dev Tools Panel"
                subtitle="Rapid API Testing & Diagnostic Tool"
            />

            <Grid container spacing={4}>
                {/* Auth State Panel */}
                <Grid item xs={12} md={6}>
                    <Paper sx={{ p: 3, height: '100%', borderRadius: 3 }}>
                        <Typography variant="h6" gutterBottom>Auth State</Typography>
                        <Box sx={{ mt: 2, p: 2, bgcolor: 'grey.50', borderRadius: 2 }}>
                            <Typography variant="body2" sx={{ mb: 1 }}>
                                <strong>Status:</strong> {isAuthenticated ? <span style={{ color: 'green' }}>Authenticated</span> : <span style={{ color: 'red' }}>Not Authenticated</span>}
                            </Typography>
                            <Typography variant="body2" sx={{ mb: 1 }}>
                                <strong>Role:</strong> <span style={{ color: '#1976d2', fontWeight: 'bold' }}>{role || 'None'}</span>
                            </Typography>
                            <Typography variant="body2" sx={{ mb: 1 }}>
                                <strong>User ID:</strong> {user?.id || 'N/A'}
                            </Typography>
                            <Typography variant="body2" sx={{ mb: 1 }}>
                                <strong>User Email:</strong> {user?.email || 'N/A'}
                            </Typography>
                            <Typography variant="caption" color="text.secondary" sx={{ display: 'block', mt: 2 }}>
                                * Tokens are hidden for security reasons.
                            </Typography>
                        </Box>
                    </Paper>
                </Grid>

                {/* API Ping Panel */}
                <Grid item xs={12} md={6}>
                    <Paper sx={{ p: 3, height: '100%', borderRadius: 3 }}>
                        <Typography variant="h6" gutterBottom>API Ping Tests</Typography>
                        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 2 }}>
                            {modulesToTest.map(mod => (
                                <Box key={mod.name} sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                                    <Button
                                        variant="outlined"
                                        onClick={() => pingModule(mod.name, mod.endpoint)}
                                        disabled={results[mod.name]?.status === 'pending'}
                                    >
                                        Ping {mod.name}
                                    </Button>

                                    <Box sx={{ flex: 1, ml: 3 }}>
                                        {results[mod.name]?.status === 'pending' && <Typography variant="body2" color="text.secondary">Testing...</Typography>}
                                        {results[mod.name]?.status === 'success' && <Typography variant="body2" color="success.main">OK ({results[mod.name].latency}ms)</Typography>}
                                        {results[mod.name]?.status === 'error' && <Typography variant="body2" color="error.main">Failed ({results[mod.name].latency}ms)</Typography>}
                                    </Box>
                                </Box>
                            ))}
                        </Box>
                    </Paper>
                </Grid>

                {/* Last Error Display */}
                {lastError && (
                    <Grid item xs={12}>
                        <Alert severity="error" sx={{ borderRadius: 2 }}>
                            <AlertTitle>Last API Error</AlertTitle>
                            <Box sx={{ mt: 1, p: 2, bgcolor: 'rgba(211,47,47,0.05)', borderRadius: 1, overflowX: 'auto' }}>
                                <pre style={{ margin: 0, fontSize: '0.875rem' }}>
                                    {JSON.stringify(lastError, null, 2)}
                                </pre>
                            </Box>
                        </Alert>
                    </Grid>
                )}
            </Grid>
        </Container>
    );
};
