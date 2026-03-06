import { useMemo, useState } from 'react';
import {
    Alert,
    Box,
    Button,
    Chip,
    IconButton,
    InputAdornment,
    MenuItem,
    Paper,
    Skeleton,
    Stack,
    TextField,
    Tooltip,
    Typography,
} from '@mui/material';
import {
    AddRounded,
    AutoStoriesOutlined,
    DeleteOutlineRounded,
    EditOutlined,
    RefreshRounded,
    SearchRounded,
    VisibilityOutlined,
} from '@mui/icons-material';
import { DataGrid, type GridColDef, type GridRenderCellParams } from '@mui/x-data-grid';
import { format } from 'date-fns';
import { toast } from 'sonner';

import {
    useDeleteProgramMutation,
    useGetProgramsQuery,
    type Program,
} from '@/api/programApi';
import ProgramDetailDialog from '@/components/admin/ProgramDetailDialog';
import ConfirmDialog from '@/components/common/ConfirmDialog';
import PageHeader from '@/components/common/PageHeader';

const getErrorMessage = (error: unknown, fallback: string) => {
    if (typeof error === 'object' && error && 'data' in error) {
        const apiError = error as { data?: { message?: string } };
        return apiError.data?.message || fallback;
    }

    if (error instanceof Error) {
        return error.message;
    }

    return fallback;
};

const trackMeta: Record<string, { label: string; color: 'info' | 'secondary' | 'success' | 'default' }> = {
    BASIC: { label: 'Cơ bản', color: 'info' },
    ADVANCED: { label: 'Nâng cao', color: 'secondary' },
    SUPPORT: { label: 'Bổ trợ', color: 'success' },
};

export const ProgramPage = () => {
    const [page, setPage] = useState(0);
    const [pageSize, setPageSize] = useState(10);
    const [search, setSearch] = useState('');
    const [trackFilter, setTrackFilter] = useState('');
    const [selectedProgramId, setSelectedProgramId] = useState<string | null>(null);
    const [deleteTarget, setDeleteTarget] = useState<Program | null>(null);

    const queryParams = useMemo(
        () => ({
            page: page + 1,
            limit: pageSize,
            search: search || undefined,
            track: trackFilter || undefined,
        }),
        [page, pageSize, search, trackFilter],
    );

    const {
        data,
        isLoading,
        isFetching,
        isError,
        error,
        refetch,
    } = useGetProgramsQuery(queryParams);

    const [deleteProgram, { isLoading: isDeleting }] = useDeleteProgramMutation();

    const programs = data?.data?.programs || [];
    const totalItems = data?.data?.pagination?.total_items || 0;

    const handleDelete = async () => {
        if (!deleteTarget) {
            return;
        }

        try {
            await deleteProgram(deleteTarget.id).unwrap();
            toast.success('Xóa chương trình thành công');
            setDeleteTarget(null);
        } catch (deleteError) {
            toast.error(getErrorMessage(deleteError, 'Lỗi khi xóa chương trình'));
        }
    };

    const columns: GridColDef<Program>[] = [
        {
            field: 'code',
            headerName: 'Mã CT',
            width: 130,
            align: 'left',
            headerAlign: 'left',
            renderCell: (params: GridRenderCellParams<Program>) => (
                <Typography variant="body2" sx={{ fontWeight: 700, color: 'primary.main' }}>
                    {params.row.code}
                </Typography>
            ),
        },
        {
            field: 'name',
            headerName: 'Tên chương trình',
            minWidth: 240,
            flex: 1.4,
            align: 'left',
            headerAlign: 'left',
            renderCell: (params: GridRenderCellParams<Program>) => (
                <Stack direction="row" spacing={1} alignItems="center" sx={{ minWidth: 0 }}>
                    <AutoStoriesOutlined sx={{ color: 'primary.main', fontSize: 18 }} />
                    <Box sx={{ minWidth: 0 }}>
                        <Typography variant="body2" sx={{ fontWeight: 700 }} noWrap>
                            {params.row.name}
                        </Typography>
                        <Typography variant="caption" color="text.secondary" noWrap>
                            {params.row.approval_note || 'Chưa có ghi chú phê duyệt'}
                        </Typography>
                    </Box>
                </Stack>
            ),
        },
        {
            field: 'track',
            headerName: 'Hệ đào tạo',
            width: 140,
            align: 'center',
            headerAlign: 'center',
            renderCell: (params: GridRenderCellParams<Program>) => {
                const meta = trackMeta[params.row.track] || { label: params.row.track || 'Chưa rõ', color: 'default' as const };
                return <Chip size="small" label={meta.label} color={meta.color} variant="outlined" />;
            },
        },
        {
            field: 'effective_from',
            headerName: 'Hiệu lực từ',
            width: 130,
            align: 'center',
            headerAlign: 'center',
            renderCell: (params: GridRenderCellParams<Program>) => (
                <Typography variant="body2">
                    {params.row.effective_from ? format(new Date(params.row.effective_from), 'dd/MM/yyyy') : '-'}
                </Typography>
            ),
        },
        {
            field: 'courses_count',
            headerName: 'Số khóa học',
            width: 120,
            align: 'center',
            headerAlign: 'center',
            sortable: false,
            renderCell: (params: GridRenderCellParams<Program>) => (
                <Typography variant="body2">{params.row.courses?.length || 0}</Typography>
            ),
        },
        {
            field: 'actions',
            headerName: '',
            width: 116,
            sortable: false,
            align: 'center',
            renderCell: (params: GridRenderCellParams<Program>) => (
                <Stack direction="row" spacing={0.5}>
                    <Tooltip title="Xem chi tiết">
                        <IconButton size="small" onClick={() => setSelectedProgramId(params.row.id)}>
                            <VisibilityOutlined fontSize="small" />
                        </IconButton>
                    </Tooltip>
                    <Tooltip title="Chỉnh sửa">
                        <IconButton
                            size="small"
                            color="primary"
                            onClick={() => toast.info('Form chỉnh sửa chương trình đang được phát triển')}
                        >
                            <EditOutlined fontSize="small" />
                        </IconButton>
                    </Tooltip>
                    <Tooltip title="Xóa">
                        <IconButton size="small" color="error" onClick={() => setDeleteTarget(params.row)}>
                            <DeleteOutlineRounded fontSize="small" />
                        </IconButton>
                    </Tooltip>
                </Stack>
            ),
        },
    ];

    const renderLoadingState = () => (
        <Stack spacing={1.5}>
            <Skeleton variant="rounded" height={72} />
            <Skeleton variant="rounded" height={420} />
        </Stack>
    );

    const renderEmptyState = () => (
        <Paper
            variant="outlined"
            sx={{
                p: 6,
                borderRadius: 4,
                textAlign: 'center',
                borderStyle: 'dashed',
            }}
        >
            <AutoStoriesOutlined sx={{ fontSize: 48, color: 'text.secondary', mb: 2 }} />
            <Typography variant="h6" sx={{ fontWeight: 700, mb: 1 }}>
                Chưa có chương trình nào phù hợp bộ lọc
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                Dữ liệu đang lấy trực tiếp từ `GET /api/v1/programs`. Thử đổi bộ lọc hoặc tạo chương trình mới.
            </Typography>
            <Button
                variant="contained"
                startIcon={<AddRounded />}
                onClick={() => toast.info('Form thêm chương trình đang được phát triển')}
            >
                Tạo chương trình mới
            </Button>
        </Paper>
    );

    return (
        <Box>
            <PageHeader
                title="Quản lý chương trình đào tạo"
                subtitle="Danh sách chương trình được đồng bộ từ API thật, có xem chi tiết và khóa học liên kết theo từng chương trình."
                breadcrumbs={[
                    { label: 'Tổng quan', path: '/app/admin/overview' },
                    { label: 'Chương trình đào tạo' },
                ]}
                actions={
                    <Stack direction="row" spacing={1}>
                        <Button
                            variant="outlined"
                            startIcon={<RefreshRounded />}
                            onClick={() => void refetch()}
                            disabled={isFetching}
                        >
                            Làm mới
                        </Button>
                        <Button
                            variant="contained"
                            startIcon={<AddRounded />}
                            onClick={() => toast.info('Form thêm chương trình đang được phát triển')}
                        >
                            Thêm chương trình
                        </Button>
                    </Stack>
                }
            />

            {isLoading && programs.length === 0 ? renderLoadingState() : null}
            {!isLoading && isError ? (
                <Alert
                    severity="error"
                    action={
                        <Button color="inherit" size="small" startIcon={<RefreshRounded />} onClick={() => void refetch()}>
                            Tải lại
                        </Button>
                    }
                    sx={{ mb: 3 }}
                >
                    {getErrorMessage(error, 'Không thể tải danh sách chương trình')}
                </Alert>
            ) : null}

            {!isLoading && !isError ? (
                programs.length === 0 ? (
                    renderEmptyState()
                ) : (
                    <Paper elevation={0} sx={{ p: 2.5, borderRadius: 4, border: '1px solid #e2e8f0' }}>
                        <Stack
                            direction={{ xs: 'column', lg: 'row' }}
                            spacing={1.5}
                            justifyContent="space-between"
                            alignItems={{ xs: 'stretch', lg: 'center' }}
                            sx={{ mb: 2 }}
                        >
                            <Stack direction={{ xs: 'column', md: 'row' }} spacing={1.5}>
                                <TextField
                                    value={search}
                                    onChange={(event) => {
                                        setSearch(event.target.value);
                                        setPage(0);
                                    }}
                                    size="small"
                                    placeholder="Tìm theo mã hoặc tên chương trình"
                                    sx={{ minWidth: { xs: '100%', md: 280 } }}
                                    InputProps={{
                                        startAdornment: (
                                            <InputAdornment position="start">
                                                <SearchRounded fontSize="small" />
                                            </InputAdornment>
                                        ),
                                    }}
                                />
                                <TextField
                                    select
                                    size="small"
                                    label="Hệ đào tạo"
                                    value={trackFilter}
                                    onChange={(event) => {
                                        setTrackFilter(event.target.value);
                                        setPage(0);
                                    }}
                                    sx={{ minWidth: { xs: '100%', md: 180 } }}
                                >
                                    <MenuItem value="">Tất cả hệ đào tạo</MenuItem>
                                    <MenuItem value="BASIC">Cơ bản</MenuItem>
                                    <MenuItem value="ADVANCED">Nâng cao</MenuItem>
                                    <MenuItem value="SUPPORT">Bổ trợ</MenuItem>
                                </TextField>
                            </Stack>

                            <Chip
                                label={`Hiển thị ${programs.length}/${totalItems || programs.length} chương trình`}
                                color="primary"
                                variant="outlined"
                            />
                        </Stack>

                        <DataGrid
                            rows={programs}
                            columns={columns}
                            rowCount={totalItems}
                            loading={isFetching}
                            paginationMode="server"
                            paginationModel={{ page, pageSize }}
                            onPaginationModelChange={(model) => {
                                setPage(model.page);
                                setPageSize(model.pageSize);
                            }}
                            pageSizeOptions={[10, 25, 50]}
                            disableRowSelectionOnClick
                            autoHeight
                            getRowId={(row) => row.id}
                            localeText={{ noRowsLabel: 'Không có dữ liệu chương trình' }}
                            sx={{
                                border: 0,
                                '& .MuiDataGrid-columnHeaders': {
                                    backgroundColor: '#f8fafc',
                                },
                                '& .MuiDataGrid-cell:focus': { outline: 'none' },
                            }}
                        />
                    </Paper>
                )
            ) : null}

            <ProgramDetailDialog
                key={selectedProgramId || 'program-detail-closed'}
                open={!!selectedProgramId}
                programId={selectedProgramId}
                onClose={() => setSelectedProgramId(null)}
            />

            <ConfirmDialog
                open={!!deleteTarget}
                title="Xóa chương trình"
                message={
                    deleteTarget
                        ? `Bạn có chắc chắn muốn xóa chương trình "${deleteTarget.name}"? Hành động này không thể hoàn tác.`
                        : ''
                }
                confirmText="Xóa chương trình"
                isDanger
                onClose={() => setDeleteTarget(null)}
                onConfirm={() => void handleDelete()}
                loading={isDeleting}
            />
        </Box>
    );
};
