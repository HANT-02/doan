import { useMemo, useState } from 'react';
import { Controller, useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import {
    Alert,
    Box,
    Button,
    Chip,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    InputAdornment,
    Paper,
    Skeleton,
    Stack,
    TextField,
    Typography,
    MenuItem,
} from '@mui/material';
import {
    CheckCircleRounded,
    RefreshRounded,
    SearchRounded,
    VisibilityOutlined,
} from '@mui/icons-material';
import { DataGrid } from '@mui/x-data-grid';
import type { GridColDef, GridRenderCellParams } from '@mui/x-data-grid';
import { toast } from 'sonner';

import { useGetMaterialsQuery, useReviewMaterialMutation, type MaterialItem } from '@/api/materialApi';
import PageHeader from '@/components/common/PageHeader';
import MaterialDetailDialog from '@/components/material/MaterialDetailDialog';

const reviewSchema = z.object({
    compliance_officer_id: z.string().min(1, 'Vui lòng nhập ID compliance officer'),
    approved: z.enum(['APPROVE', 'REJECT']),
    reject_reason: z.string().optional(),
    notes: z.string().optional(),
});

type ReviewFormValues = z.infer<typeof reviewSchema>;

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

export const ComplianceQueuePage = () => {
    const [search, setSearch] = useState('');
    const [selectedMaterial, setSelectedMaterial] = useState<MaterialItem | null>(null);
    const [reviewTarget, setReviewTarget] = useState<MaterialItem | null>(null);

    const { data, isLoading, isFetching, refetch } = useGetMaterialsQuery({ queue: 'flagged', status: 'AI_REVIEWED' });
    const [reviewMaterial, { isLoading: isReviewing }] = useReviewMaterialMutation();

    const {
        control,
        handleSubmit,
        formState: { errors },
        reset,
    } = useForm<ReviewFormValues>({
        resolver: zodResolver(reviewSchema),
        defaultValues: {
            compliance_officer_id: '',
            approved: 'APPROVE',
            reject_reason: '',
            notes: '',
        },
    });

    const materials = useMemo(() => data?.data?.materials || [], [data?.data?.materials]);
    const filteredMaterials = useMemo(() => {
        const keyword = search.trim().toLowerCase();
        if (!keyword) {
            return materials;
        }
        return materials.filter((item) =>
            [item.title, item.file_name, item.latest_label?.name, item.latest_label?.severity]
                .filter(Boolean)
                .some((value) => value!.toLowerCase().includes(keyword)),
        );
    }, [materials, search]);

    const handleReviewSubmit = handleSubmit(async (values) => {
        if (!reviewTarget) {
            return;
        }

        try {
            await reviewMaterial({
                id: reviewTarget.id,
                compliance_officer_id: values.compliance_officer_id,
                approved: values.approved === 'APPROVE',
                reject_reason: values.approved === 'REJECT' ? values.reject_reason : '',
                notes: values.notes,
            }).unwrap();
            toast.success('Đã cập nhật quyết định kiểm duyệt');
            setReviewTarget(null);
            reset();
        } catch (error) {
            toast.error(getErrorMessage(error, 'Không thể xử lý phê duyệt tài liệu'));
        }
    });

    const columns: GridColDef<MaterialItem>[] = [
        { field: 'title', headerName: 'Tài liệu', minWidth: 220, flex: 1.2 },
        {
            field: 'latest_label',
            headerName: 'Mức độ',
            minWidth: 160,
            renderCell: (params: GridRenderCellParams<MaterialItem>) => (
                params.row.latest_label ? (
                    <Chip
                        label={`${params.row.latest_label.name} (${params.row.latest_label.severity})`}
                        color={params.row.latest_label.severity === 'SAFE' ? 'success' : params.row.latest_label.severity === 'DANGER' ? 'error' : 'warning'}
                        variant="outlined"
                        size="small"
                    />
                ) : (
                    <Chip label="Chưa có nhãn" size="small" />
                )
            ),
        },
        {
            field: 'reasoning',
            headerName: 'AI reasoning',
            minWidth: 260,
            flex: 1,
            renderCell: (params: GridRenderCellParams<MaterialItem>) => (
                <Typography variant="body2" noWrap>
                    {params.row.latest_audit?.reasoning || 'Chưa có reasoning'}
                </Typography>
            ),
        },
        {
            field: 'uploaded_at',
            headerName: 'Tải lên',
            minWidth: 180,
            renderCell: (params: GridRenderCellParams<MaterialItem>) => (
                <Typography variant="body2">{new Date(params.row.uploaded_at).toLocaleString('vi-VN')}</Typography>
            ),
        },
        {
            field: 'actions',
            headerName: '',
            width: 180,
            sortable: false,
            renderCell: (params: GridRenderCellParams<MaterialItem>) => (
                <Stack direction="row" spacing={1}>
                    <Button startIcon={<VisibilityOutlined />} onClick={() => setSelectedMaterial(params.row)}>
                        Xem
                    </Button>
                    <Button variant="contained" color="secondary" onClick={() => setReviewTarget(params.row)}>
                        Duyệt
                    </Button>
                </Stack>
            ),
        },
    ];

    return (
        <Box>
            <PageHeader
                title="Kiểm duyệt tài liệu"
                subtitle="Hàng chờ tài liệu đã qua AI audit và cần compliance officer quyết định cuối cùng."
                breadcrumbs={[
                    { label: 'Compliance Dashboard', path: '/app/compliance/overview' },
                    { label: 'Tài liệu cần duyệt' },
                ]}
                actions={
                    <Button variant="outlined" startIcon={<RefreshRounded />} onClick={() => void refetch()} disabled={isFetching}>
                        Làm mới
                    </Button>
                }
            />

            <Stack spacing={3}>
                <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 4 }}>
                    <Stack direction={{ xs: 'column', md: 'row' }} spacing={1.5} justifyContent="space-between" sx={{ mb: 2 }}>
                        <TextField
                            size="small"
                            value={search}
                            onChange={(event) => setSearch(event.target.value)}
                            placeholder="Tìm theo tài liệu, nhãn, severity..."
                            sx={{ minWidth: { xs: '100%', md: 320 } }}
                            InputProps={{
                                startAdornment: (
                                    <InputAdornment position="start">
                                        <SearchRounded fontSize="small" />
                                    </InputAdornment>
                                ),
                            }}
                        />
                        <Chip label={`${filteredMaterials.length} tài liệu cần duyệt`} color="warning" variant="outlined" />
                    </Stack>

                    {isLoading ? (
                        <Stack spacing={1.5}>
                            <Skeleton variant="rounded" height={64} />
                            <Skeleton variant="rounded" height={320} />
                        </Stack>
                    ) : null}
                    {!isLoading && filteredMaterials.length === 0 ? (
                        <Alert severity="success">Hiện không có tài liệu nào trong hàng chờ kiểm duyệt.</Alert>
                    ) : null}
                    {!isLoading && filteredMaterials.length > 0 ? (
                        <DataGrid
                            rows={filteredMaterials}
                            columns={columns}
                            autoHeight
                            disableRowSelectionOnClick
                            getRowId={(row) => row.id}
                            pageSizeOptions={[5, 10, 25]}
                            initialState={{
                                pagination: {
                                    paginationModel: {
                                        page: 0,
                                        pageSize: 10,
                                    },
                                },
                            }}
                            sx={{
                                border: 'none',
                                '& .MuiDataGrid-columnHeaders': {
                                    backgroundColor: '#f8fafc',
                                },
                            }}
                        />
                    ) : null}
                </Paper>
            </Stack>

            <MaterialDetailDialog open={!!selectedMaterial} material={selectedMaterial} onClose={() => setSelectedMaterial(null)} />

            <Dialog open={!!reviewTarget} onClose={() => setReviewTarget(null)} fullWidth maxWidth="sm">
                <DialogTitle>
                    <Typography variant="h6" sx={{ fontWeight: 800 }}>
                        Phê duyệt tài liệu
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        {reviewTarget?.title}
                    </Typography>
                </DialogTitle>
                <DialogContent dividers>
                    <Stack spacing={2}>
                        <Alert severity="warning">
                            Trang này đang dùng AI audit stub. Quyết định approve/reject vẫn được lưu xuống backend scaffold để demo workflow.
                        </Alert>
                        <Controller
                            control={control}
                            name="compliance_officer_id"
                            render={({ field }) => (
                                <TextField
                                    {...field}
                                    label="Compliance officer ID"
                                    error={!!errors.compliance_officer_id}
                                    helperText={errors.compliance_officer_id?.message || 'Nhập ID người duyệt để lưu audit trail'}
                                />
                            )}
                        />
                        <Controller
                            control={control}
                            name="approved"
                            render={({ field }) => (
                                <TextField
                                    {...field}
                                    select
                                    label="Quyết định"
                                    helperText="Approve nếu nội dung đạt yêu cầu, Reject nếu cần trả lại giáo viên"
                                >
                                    <MenuItem value="APPROVE">Approve</MenuItem>
                                    <MenuItem value="REJECT">Reject</MenuItem>
                                </TextField>
                            )}
                        />
                        <Controller
                            control={control}
                            name="notes"
                            render={({ field }) => (
                                <TextField
                                    {...field}
                                    multiline
                                    minRows={3}
                                    label="Ghi chú"
                                />
                            )}
                        />
                        <Controller
                            control={control}
                            name="reject_reason"
                            render={({ field }) => (
                                <TextField
                                    {...field}
                                    multiline
                                    minRows={2}
                                    label="Lý do từ chối"
                                    helperText="Chỉ cần nhập khi chọn Reject"
                                />
                            )}
                        />
                    </Stack>
                </DialogContent>
                <DialogActions sx={{ px: 3, py: 2 }}>
                    <Button onClick={() => setReviewTarget(null)}>Hủy</Button>
                    <Button variant="contained" startIcon={<CheckCircleRounded />} onClick={() => void handleReviewSubmit()} disabled={isReviewing}>
                        {isReviewing ? 'Đang lưu...' : 'Lưu quyết định'}
                    </Button>
                </DialogActions>
            </Dialog>
        </Box>
    );
};
