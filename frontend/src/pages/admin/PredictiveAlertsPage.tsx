import { useMemo, useState } from 'react';
import {
    Alert,
    Box,
    Breadcrumbs,
    Button,
    Chip,
    Link,
    Paper,
    Stack,
    Switch,
    TextField,
    Typography,
} from '@mui/material';
import { PsychologyRounded, RefreshRounded, SearchRounded, WarningAmberRounded } from '@mui/icons-material';
import { DataGrid } from '@mui/x-data-grid';
import type { GridColDef, GridRenderCellParams } from '@mui/x-data-grid';

import { useGetAtRiskPredictionsQuery } from '@/api/predictiveApi';
import PageHeader from '@/components/common/PageHeader';
import { getApiErrorMessage } from '@/utils/apiError';

const formatPercent = (value: number) => `${Math.round(value * 100)}%`;
const formatDateTime = (value: string) => new Date(value).toLocaleString('vi-VN');

export const PredictiveAlertsPage = () => {
    const [page, setPage] = useState(0);
    const [pageSize, setPageSize] = useState(10);
    const [search, setSearch] = useState('');
    const [searchInput, setSearchInput] = useState('');
    const [onlyAtRisk, setOnlyAtRisk] = useState(true);
    const [refreshKey, setRefreshKey] = useState(0);

    const { data, isLoading, isFetching, error } = useGetAtRiskPredictionsQuery({
        page: page + 1,
        limit: pageSize,
        search,
        only_at_risk: onlyAtRisk,
        refresh: refreshKey > 0,
    });

    const rows = data?.data?.items || [];
    const summary = data?.data?.summary;
    const modelMetadata = data?.data?.model_metadata;

    const columns = useMemo<GridColDef[]>(() => [
        {
            field: 'student_name',
            headerName: 'Học sinh',
            minWidth: 220,
            flex: 1.1,
            renderCell: (params: GridRenderCellParams) => (
                <Stack spacing={0.25} sx={{ py: 1 }}>
                    <Typography variant="body2" sx={{ fontWeight: 700 }}>
                        {params.row.student_name}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                        {params.row.student_code} • Khối {params.row.grade_level || 'N/A'}
                    </Typography>
                </Stack>
            ),
        },
        {
            field: 'class_name',
            headerName: 'Lớp học',
            minWidth: 180,
            flex: 1,
            renderCell: (params: GridRenderCellParams) => (
                <Stack spacing={0.25} sx={{ py: 1 }}>
                    <Typography variant="body2" sx={{ fontWeight: 600 }}>
                        {params.row.class_name}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                        {params.row.class_code}
                    </Typography>
                </Stack>
            ),
        },
        {
            field: 'label',
            headerName: 'Kết quả',
            width: 135,
            renderCell: (params: GridRenderCellParams) => (
                <Chip
                    label={params.value === 'AT_RISK' ? 'AT_RISK' : 'Ổn định'}
                    color={params.value === 'AT_RISK' ? 'error' : 'success'}
                    size="small"
                    sx={{ fontWeight: 700 }}
                />
            ),
        },
        {
            field: 'risk_score',
            headerName: 'Risk Score',
            width: 125,
            renderCell: (params: GridRenderCellParams) => (
                <Stack spacing={0.25} sx={{ py: 1 }}>
                    <Typography variant="body2" sx={{ fontWeight: 700 }}>
                        {formatPercent(params.row.risk_score)}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                        {params.row.risk_band}
                    </Typography>
                </Stack>
            ),
        },
        {
            field: 'primary_reason',
            headerName: 'Insight chính',
            minWidth: 320,
            flex: 1.7,
            renderCell: (params: GridRenderCellParams) => (
                <Stack spacing={0.25} sx={{ py: 1 }}>
                    <Typography variant="body2">
                        {params.row.primary_reason}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                        Chuyên cần {formatPercent(params.row.feature_summary.attendance_rate_28d)} • Điểm TB {params.row.feature_summary.average_total_score_28d.toFixed(2)} • Homework {formatPercent(params.row.feature_summary.homework_completion_rate_28d)}
                    </Typography>
                </Stack>
            ),
        },
        {
            field: 'snapshot_at',
            headerName: 'Snapshot',
            width: 180,
            renderCell: (params: GridRenderCellParams) => (
                <Typography variant="body2">
                    {formatDateTime(params.value as string)}
                </Typography>
            ),
        },
    ], []);

    return (
        <Box>
            <PageHeader
                title="Cảnh báo học sinh nguy cơ học kém"
                subtitle="Danh sách AT_RISK được suy luận từ chuyên cần, điểm số và tín hiệu vận hành 28 ngày gần nhất."
                actions={(
                    <Button
                        variant="contained"
                        startIcon={<RefreshRounded />}
                        onClick={() => setRefreshKey((current) => current + 1)}
                        disabled={isFetching}
                    >
                        Train lại từ DB
                    </Button>
                )}
            />

            <Breadcrumbs sx={{ mb: 3 }}>
                <Link underline="hover" color="inherit" href="/app/admin/overview">Dashboard</Link>
                <Typography color="text.primary">Predictive Analytics</Typography>
            </Breadcrumbs>

            {error && (
                <Alert severity="error" sx={{ mb: 3 }}>
                    {getApiErrorMessage(error, 'Không tải được predictive analytics.')}
                </Alert>
            )}

            <Stack direction={{ xs: 'column', lg: 'row' }} spacing={2} sx={{ mb: 3 }}>
                <Paper elevation={0} sx={{ p: 2.5, borderRadius: 3, border: '1px solid #fee2e2', flex: 1 }}>
                    <Stack direction="row" spacing={1.5} alignItems="center" sx={{ mb: 1 }}>
                        <WarningAmberRounded color="error" />
                        <Typography variant="subtitle2" sx={{ fontWeight: 700 }}>Số học sinh AT_RISK</Typography>
                    </Stack>
                    <Typography variant="h4" sx={{ fontWeight: 800 }}>
                        {summary?.at_risk_count ?? 0}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        Trên tổng {summary?.total_students_evaluated ?? 0} học sinh/lớp được đánh giá.
                    </Typography>
                </Paper>

                <Paper elevation={0} sx={{ p: 2.5, borderRadius: 3, border: '1px solid #dbeafe', flex: 1 }}>
                    <Stack direction="row" spacing={1.5} alignItems="center" sx={{ mb: 1 }}>
                        <PsychologyRounded color="primary" />
                        <Typography variant="subtitle2" sx={{ fontWeight: 700 }}>Model hiện dùng</Typography>
                    </Stack>
                    <Typography variant="h5" sx={{ fontWeight: 800 }}>
                        {modelMetadata?.model_name || 'N/A'}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        Version {modelMetadata?.version || 'N/A'} • Train size {modelMetadata?.train_size ?? 0} • Test size {modelMetadata?.test_size ?? 0}
                    </Typography>
                </Paper>
            </Stack>

            <Paper elevation={0} sx={{ p: 2, borderRadius: 3, border: '1px solid #e2e8f0', mb: 3 }}>
                <Stack direction={{ xs: 'column', md: 'row' }} spacing={2} alignItems={{ md: 'center' }}>
                    <TextField
                        label="Tìm học sinh / lớp"
                        size="small"
                        value={searchInput}
                        onChange={(event) => setSearchInput(event.target.value)}
                        onKeyDown={(event) => {
                            if (event.key === 'Enter') {
                                setPage(0);
                                setSearch(searchInput.trim());
                            }
                        }}
                        InputProps={{
                            startAdornment: <SearchRounded sx={{ mr: 1, color: 'text.secondary' }} />,
                        }}
                        sx={{ minWidth: 280 }}
                    />
                    <Button
                        variant="outlined"
                        onClick={() => {
                            setPage(0);
                            setSearch(searchInput.trim());
                        }}
                    >
                        Áp dụng
                    </Button>
                    <Stack direction="row" spacing={1} alignItems="center">
                        <Switch
                            checked={onlyAtRisk}
                            onChange={(_, checked) => {
                                setPage(0);
                                setOnlyAtRisk(checked);
                            }}
                        />
                        <Typography variant="body2">Chỉ hiện AT_RISK</Typography>
                    </Stack>
                    <Chip
                        label={`Train lúc: ${modelMetadata?.trained_at ? formatDateTime(modelMetadata.trained_at) : 'N/A'}`}
                        variant="outlined"
                        sx={{ ml: { md: 'auto' } }}
                    />
                </Stack>
            </Paper>

            <Paper elevation={0} sx={{ p: 2, borderRadius: 3, border: '1px solid #e2e8f0' }}>
                <DataGrid
                    rows={rows}
                    columns={columns}
                    loading={isLoading || isFetching}
                    paginationMode="server"
                    rowCount={data?.data?.pagination?.total_items || 0}
                    paginationModel={{ page, pageSize }}
                    onPaginationModelChange={(model) => {
                        setPage(model.page);
                        setPageSize(model.pageSize);
                    }}
                    pageSizeOptions={[10, 25, 50]}
                    disableRowSelectionOnClick
                    autoHeight
                    getRowId={(row) => `${row.student_id}-${row.class_id}`}
                    sx={{
                        border: 'none',
                        '& .MuiDataGrid-columnHeaders': {
                            backgroundColor: '#f8fafc',
                        },
                        '& .MuiDataGrid-cell:focus': {
                            outline: 'none',
                        },
                    }}
                />
            </Paper>
        </Box>
    );
};
