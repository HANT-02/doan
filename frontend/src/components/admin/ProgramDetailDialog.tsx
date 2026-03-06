import {
    Alert,
    Box,
    Button,
    Chip,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    Divider,
    Paper,
    Skeleton,
    Stack,
    Tab,
    Tabs,
    Typography,
} from '@mui/material';
import { AutoStoriesOutlined, InfoOutlined, LinkRounded, RefreshRounded } from '@mui/icons-material';
import { DataGrid, type GridColDef, type GridRenderCellParams } from '@mui/x-data-grid';
import { format } from 'date-fns';
import { useMemo, useState } from 'react';
import { toast } from 'sonner';

import { useGetProgramByIdQuery } from '@/api/programApi';
import type { Course } from '@/api/courseApi';

interface ProgramDetailDialogProps {
    open: boolean;
    programId: string | null;
    onClose: () => void;
}

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

const infoField = (label: string, value: string) => (
    <Stack spacing={0.5}>
        <Typography variant="caption" color="text.secondary">
            {label}
        </Typography>
        <Typography variant="body2" sx={{ fontWeight: 600 }}>
            {value || '-'}
        </Typography>
    </Stack>
);

export default function ProgramDetailDialog({ open, programId, onClose }: ProgramDetailDialogProps) {
    const [activeTab, setActiveTab] = useState(0);

    const {
        data,
        isLoading,
        isFetching,
        isError,
        error,
        refetch,
    } = useGetProgramByIdQuery(programId || '', {
        skip: !open || !programId,
    });

    const program = data?.data || null;
    const courses = useMemo(() => program?.courses || [], [program]);

    const columns: GridColDef<Course>[] = [
        {
            field: 'code',
            headerName: 'Mã khóa học',
            width: 140,
            align: 'left',
            headerAlign: 'left',
            renderCell: (params: GridRenderCellParams<Course>) => (
                <Typography variant="body2" sx={{ fontWeight: 700, color: 'primary.main' }}>
                    {params.row.code}
                </Typography>
            ),
        },
        {
            field: 'name',
            headerName: 'Tên khóa học',
            minWidth: 240,
            flex: 1.2,
            align: 'left',
            headerAlign: 'left',
            renderCell: (params: GridRenderCellParams<Course>) => (
                <Box sx={{ minWidth: 0 }}>
                    <Typography variant="body2" sx={{ fontWeight: 700 }} noWrap>
                        {params.row.name}
                    </Typography>
                    <Typography variant="caption" color="text.secondary" noWrap>
                        {params.row.subject || 'Chưa phân môn'}
                    </Typography>
                </Box>
            ),
        },
        {
            field: 'grade_level',
            headerName: 'Khối lớp',
            width: 120,
            align: 'center',
            headerAlign: 'center',
            renderCell: (params: GridRenderCellParams<Course>) => (
                <Typography variant="body2">{params.row.grade_level || '-'}</Typography>
            ),
        },
        {
            field: 'total_hours',
            headerName: 'Tổng giờ',
            width: 100,
            align: 'center',
            headerAlign: 'center',
            renderCell: (params: GridRenderCellParams<Course>) => (
                <Typography variant="body2">{params.row.total_hours || 0}</Typography>
            ),
        },
        {
            field: 'status',
            headerName: 'Trạng thái',
            width: 130,
            align: 'center',
            headerAlign: 'center',
            renderCell: (params: GridRenderCellParams<Course>) => (
                <Chip
                    size="small"
                    color={params.row.status === 'ACTIVE' ? 'success' : 'default'}
                    variant="outlined"
                    label={params.row.status === 'ACTIVE' ? 'Hoạt động' : params.row.status || 'Chưa rõ'}
                />
            ),
        },
    ];

    const renderLoading = () => (
        <Stack spacing={2}>
            <Skeleton variant="rounded" height={84} />
            <Skeleton variant="rounded" height={300} />
        </Stack>
    );

    const renderProgramInfo = () => {
        if (!program) {
            return (
                <Paper
                    variant="outlined"
                    sx={{ p: 4, borderRadius: 3, textAlign: 'center', borderStyle: 'dashed' }}
                >
                    <InfoOutlined sx={{ fontSize: 42, color: 'text.secondary', mb: 1.5 }} />
                    <Typography variant="subtitle1" sx={{ fontWeight: 700, mb: 1 }}>
                        Không có dữ liệu chương trình
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        Chọn lại một chương trình từ danh sách để tải chi tiết.
                    </Typography>
                </Paper>
            );
        }

        return (
            <Stack spacing={2}>
                <Alert severity="info">
                    Contract backend hiện cung cấp `code`, `name`, `track`, `effective_from`, `effective_to`, `approval_note` và `courses`. Chưa có field `description/status`.
                </Alert>

                <Paper variant="outlined" sx={{ p: 3, borderRadius: 3 }}>
                    <Stack spacing={2.5}>
                        <Stack direction="row" spacing={1} alignItems="center" justifyContent="space-between">
                            <Box>
                                <Typography variant="h6" sx={{ fontWeight: 700 }}>
                                    {program.name}
                                </Typography>
                                <Typography variant="body2" color="text.secondary">
                                    Mã chương trình: {program.code}
                                </Typography>
                            </Box>
                            <Chip label={program.track || 'Chưa phân hệ'} color="primary" variant="outlined" />
                        </Stack>

                        <Divider />

                        <Box
                            sx={{
                                display: 'grid',
                                gridTemplateColumns: { xs: '1fr', sm: 'repeat(2, minmax(0, 1fr))' },
                                gap: 2,
                            }}
                        >
                            {infoField('Mã chương trình', program.code)}
                            {infoField('Tên chương trình', program.name)}
                            {infoField('Hệ đào tạo', program.track || '-')}
                            {infoField('Số khóa học liên kết', `${courses.length}`)}
                            {infoField(
                                'Hiệu lực từ',
                                program.effective_from ? format(new Date(program.effective_from), 'dd/MM/yyyy') : '-',
                            )}
                            {infoField(
                                'Hiệu lực đến',
                                program.effective_to ? format(new Date(program.effective_to), 'dd/MM/yyyy') : '-',
                            )}
                        </Box>

                        <Box>
                            <Typography variant="caption" color="text.secondary">
                                Ghi chú phê duyệt
                            </Typography>
                            <Typography variant="body2" sx={{ mt: 0.75 }}>
                                {program.approval_note || 'Chưa có ghi chú phê duyệt'}
                            </Typography>
                        </Box>
                    </Stack>
                </Paper>
            </Stack>
        );
    };

    const renderCoursesTab = () => {
        if (!program) {
            return null;
        }

        if (courses.length === 0) {
            return (
                <Paper
                    variant="outlined"
                    sx={{ p: 4, borderRadius: 3, textAlign: 'center', borderStyle: 'dashed' }}
                >
                    <LinkRounded sx={{ fontSize: 42, color: 'text.secondary', mb: 1.5 }} />
                    <Typography variant="subtitle1" sx={{ fontWeight: 700, mb: 1 }}>
                        Chưa có khóa học nào được liên kết
                    </Typography>
                    <Typography variant="body2" color="text.secondary" sx={{ mb: 2.5 }}>
                        Backend detail đã trả trường `courses`, hiện đang rỗng cho chương trình này.
                    </Typography>
                    <Button
                        variant="contained"
                        startIcon={<LinkRounded />}
                        onClick={() => toast.info('Liên kết khóa học được thực hiện từ màn quản lý Program/Course')}
                    >
                        Liên kết khóa học
                    </Button>
                </Paper>
            );
        }

        return (
            <Paper variant="outlined" sx={{ p: 2, borderRadius: 3 }}>
                <DataGrid
                    rows={courses}
                    columns={columns}
                    autoHeight
                    disableRowSelectionOnClick
                    hideFooter
                    getRowId={(row) => row.id}
                    localeText={{ noRowsLabel: 'Chưa có khóa học liên kết' }}
                    sx={{
                        border: 'none',
                        '& .MuiDataGrid-columnHeaders': {
                            backgroundColor: '#f8fafc',
                        },
                    }}
                />
            </Paper>
        );
    };

    return (
        <Dialog open={open} onClose={onClose} fullWidth maxWidth="lg">
            <DialogTitle sx={{ pb: 1.5 }}>
                <Stack direction="row" justifyContent="space-between" alignItems="center" spacing={2}>
                    <Box>
                        <Typography variant="h6" sx={{ fontWeight: 700 }}>
                            Chi tiết chương trình đào tạo
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                            Gọi trực tiếp `GET /api/v1/programs/:id` để hiển thị dữ liệu thật.
                        </Typography>
                    </Box>
                    <Button
                        variant="outlined"
                        startIcon={<RefreshRounded />}
                        onClick={() => void refetch()}
                        disabled={isFetching || !programId}
                    >
                        Tải lại
                    </Button>
                </Stack>
            </DialogTitle>

            <DialogContent dividers sx={{ p: 3 }}>
                {isLoading ? renderLoading() : null}

                {!isLoading && isError ? (
                    <Alert
                        severity="error"
                        action={
                            <Button color="inherit" size="small" startIcon={<RefreshRounded />} onClick={() => void refetch()}>
                                Tải lại
                            </Button>
                        }
                    >
                        {getErrorMessage(error, 'Không thể tải chi tiết chương trình')}
                    </Alert>
                ) : null}

                {!isLoading && !isError ? (
                    <Stack spacing={2}>
                        <Tabs value={activeTab} onChange={(_event, value) => setActiveTab(value)}>
                            <Tab icon={<InfoOutlined fontSize="small" />} iconPosition="start" label="Thông tin chương trình" />
                            <Tab icon={<AutoStoriesOutlined fontSize="small" />} iconPosition="start" label="Khóa học thuộc chương trình" />
                        </Tabs>

                        {activeTab === 0 ? renderProgramInfo() : null}
                        {activeTab === 1 ? renderCoursesTab() : null}
                    </Stack>
                ) : null}
            </DialogContent>

            <DialogActions sx={{ px: 3, py: 2 }}>
                <Button onClick={onClose} color="inherit">
                    Đóng
                </Button>
            </DialogActions>
        </Dialog>
    );
}
