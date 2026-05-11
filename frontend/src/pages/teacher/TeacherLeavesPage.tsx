import { useMemo, useState } from 'react';
import {
    Alert,
    Button,
    Chip,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    MenuItem,
    Paper,
    Stack,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    TextField,
    Typography,
} from '@mui/material';
import { CheckRounded, CloseRounded, DescriptionRounded } from '@mui/icons-material';
import { toast } from 'sonner';

import {
    useApproveTeacherLeaveRequestMutation,
    useGetTeacherLeaveRequestsQuery,
    useGetTeacherLessonsQuery,
    useRejectTeacherLeaveRequestMutation,
} from '@/api/teacherPortalApi';
import PageHeader from '@/components/common/PageHeader';
import { getApiErrorMessage } from '@/utils/apiError';

function leaveTypeLabel(value: string) {
    switch (value) {
        case 'LEAVE':
            return 'Xin nghỉ';
        case 'LATE':
            return 'Xin đi muộn';
        case 'EARLY':
            return 'Xin về sớm';
        default:
            return value;
    }
}

function statusColor(value: string): 'success' | 'error' | 'warning' | 'default' {
    switch (value) {
        case 'APPROVED':
            return 'success';
        case 'REJECTED':
            return 'error';
        case 'PENDING':
            return 'warning';
        default:
            return 'default';
    }
}

export default function TeacherLeavesPage() {
    const [statusFilter, setStatusFilter] = useState('');
    const [classFilter, setClassFilter] = useState('');
    const [studentFilter, setStudentFilter] = useState('');
    const [rejectingId, setRejectingId] = useState<string | null>(null);
    const [rejectionReason, setRejectionReason] = useState('');

    const lessonsQuery = useGetTeacherLessonsQuery();
    const requestsQuery = useGetTeacherLeaveRequestsQuery({
        class_id: classFilter || undefined,
        status: statusFilter || undefined,
        student_id: studentFilter || undefined,
    });
    const requestOptionsQuery = useGetTeacherLeaveRequestsQuery(
        classFilter
            ? {
                class_id: classFilter,
            }
            : undefined,
    );

    const [approveLeaveRequest, { isLoading: isApproving }] = useApproveTeacherLeaveRequestMutation();
    const [rejectLeaveRequest, { isLoading: isRejecting }] = useRejectTeacherLeaveRequestMutation();

    const requests = requestsQuery.data?.data?.requests ?? [];
    const optionRequests = requestOptionsQuery.data?.data?.requests ?? [];
    const classes = useMemo(() => {
        const map = new Map<string, { id: string; code: string; name: string }>();
        (lessonsQuery.data?.data?.lessons ?? []).forEach((lesson) => {
            if (!map.has(lesson.class_id)) {
                map.set(lesson.class_id, {
                    id: lesson.class_id,
                    code: lesson.class_code,
                    name: lesson.class_name,
                });
            }
        });
        return Array.from(map.values());
    }, [lessonsQuery.data]);

    const students = useMemo(() => {
        const map = new Map<string, { id: string; code: string; full_name: string }>();
        optionRequests.forEach((request) => {
            if (!map.has(request.student.id)) {
                map.set(request.student.id, request.student);
            }
        });
        return Array.from(map.values()).sort((left, right) => left.full_name.localeCompare(right.full_name, 'vi'));
    }, [optionRequests]);

    const pendingCount = useMemo(() => requests.filter((item) => item.status === 'PENDING').length, [requests]);

    const handleApprove = async (id: string) => {
        try {
            await approveLeaveRequest(id).unwrap();
            toast.success('Đã duyệt đơn xin phép.');
            requestsQuery.refetch();
            requestOptionsQuery.refetch();
        } catch (err) {
            toast.error(getApiErrorMessage(err, 'Không duyệt được đơn xin phép.'));
        }
    };

    const handleReject = async () => {
        if (!rejectingId) {
            return;
        }

        try {
            await rejectLeaveRequest({
                id: rejectingId,
                rejection_reason: rejectionReason,
            }).unwrap();
            toast.success('Đã từ chối đơn xin phép.');
            setRejectingId(null);
            setRejectionReason('');
            requestsQuery.refetch();
            requestOptionsQuery.refetch();
        } catch (err) {
            toast.error(getApiErrorMessage(err, 'Không từ chối được đơn xin phép.'));
        }
    };

    return (
        <Stack sx={{ p: { xs: 2, md: 4 } }} spacing={3}>
            <PageHeader
                title="Đơn xin phép"
                subtitle="Xem, lọc, duyệt hoặc từ chối các đơn xin phép thuộc lớp bạn đang phụ trách."
                icon={<DescriptionRounded />}
            />

            <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 3 }}>
                <Stack direction={{ xs: 'column', lg: 'row' }} spacing={2}>
                    <TextField
                        select
                        label="Trạng thái"
                        value={statusFilter}
                        onChange={(event) => setStatusFilter(event.target.value)}
                        sx={{ minWidth: 180 }}
                    >
                        <MenuItem value="">Tất cả</MenuItem>
                        <MenuItem value="PENDING">PENDING</MenuItem>
                        <MenuItem value="APPROVED">APPROVED</MenuItem>
                        <MenuItem value="REJECTED">REJECTED</MenuItem>
                    </TextField>
                    <TextField
                        select
                        label="Lớp"
                        value={classFilter}
                        onChange={(event) => {
                            setClassFilter(event.target.value);
                            setStudentFilter('');
                        }}
                        sx={{ minWidth: 220 }}
                    >
                        <MenuItem value="">Tất cả lớp</MenuItem>
                        {classes.map((item) => (
                            <MenuItem key={item.id} value={item.id}>
                                {item.name} ({item.code})
                            </MenuItem>
                        ))}
                    </TextField>
                    <TextField
                        select
                        label="Học sinh"
                        value={studentFilter}
                        onChange={(event) => setStudentFilter(event.target.value)}
                        sx={{ minWidth: 240 }}
                    >
                        <MenuItem value="">Tất cả học sinh</MenuItem>
                        {students.map((item) => (
                            <MenuItem key={item.id} value={item.id}>
                                {item.full_name} ({item.code})
                            </MenuItem>
                        ))}
                    </TextField>
                    <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap alignItems="center">
                        <Chip size="small" label={`${requests.length} đơn`} variant="outlined" />
                        <Chip size="small" label={`Pending ${pendingCount}`} color="warning" variant="outlined" />
                    </Stack>
                </Stack>
            </Paper>

            {requestsQuery.error ? (
                <Alert severity="error">
                    {getApiErrorMessage(requestsQuery.error, 'Không tải được danh sách đơn xin phép của giáo viên.')}
                </Alert>
            ) : null}

            <TableContainer component={Paper} variant="outlined" sx={{ borderRadius: 3 }}>
                <Table>
                    <TableHead>
                        <TableRow>
                            <TableCell sx={{ fontWeight: 700 }}>Học sinh</TableCell>
                            <TableCell sx={{ fontWeight: 700 }}>Lớp</TableCell>
                            <TableCell sx={{ fontWeight: 700 }}>Loại đơn</TableCell>
                            <TableCell sx={{ fontWeight: 700 }}>Ngày xin</TableCell>
                            <TableCell sx={{ fontWeight: 700 }}>Lý do</TableCell>
                            <TableCell sx={{ fontWeight: 700 }}>Tài liệu</TableCell>
                            <TableCell sx={{ fontWeight: 700 }}>Trạng thái</TableCell>
                            <TableCell sx={{ fontWeight: 700 }}>Thao tác</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {requests.map((request) => (
                            <TableRow key={request.id} hover>
                                <TableCell>
                                    <Typography variant="body2" sx={{ fontWeight: 600 }}>
                                        {request.student.full_name}
                                    </Typography>
                                    <Typography variant="caption" color="text.secondary">
                                        {request.student.code}
                                    </Typography>
                                </TableCell>
                                <TableCell>
                                    <Typography variant="body2">{request.class?.name || 'Không gắn lớp'}</Typography>
                                    <Typography variant="caption" color="text.secondary">
                                        {request.class?.code || ''}
                                    </Typography>
                                </TableCell>
                                <TableCell>{leaveTypeLabel(request.leave_type)}</TableCell>
                                <TableCell>
                                    <Typography variant="body2">
                                        {new Date(request.apply_date).toLocaleString('vi-VN')}
                                    </Typography>
                                    {request.lesson ? (
                                        <Typography variant="caption" color="text.secondary">
                                            Buổi học {new Date(request.lesson.date_start).toLocaleString('vi-VN')}
                                        </Typography>
                                    ) : null}
                                </TableCell>
                                <TableCell sx={{ maxWidth: 260 }}>
                                    <Typography variant="body2">{request.reason}</Typography>
                                    {request.rejection_reason ? (
                                        <Typography variant="caption" color="error.main">
                                            Từ chối: {request.rejection_reason}
                                        </Typography>
                                    ) : null}
                                </TableCell>
                                <TableCell>
                                    {request.documents?.length ? (
                                        <Stack spacing={0.5}>
                                            {request.documents.map((documentUrl) => (
                                                <Typography
                                                    key={documentUrl}
                                                    component="a"
                                                    href={documentUrl}
                                                    target="_blank"
                                                    rel="noreferrer"
                                                    variant="caption"
                                                    sx={{ color: 'primary.main', textDecoration: 'none' }}
                                                >
                                                    {documentUrl}
                                                </Typography>
                                            ))}
                                        </Stack>
                                    ) : (
                                        <Typography variant="caption" color="text.secondary">
                                            Không có
                                        </Typography>
                                    )}
                                </TableCell>
                                <TableCell>
                                    <Chip size="small" label={request.status} color={statusColor(request.status)} variant="outlined" />
                                </TableCell>
                                <TableCell>
                                    {request.status === 'PENDING' ? (
                                        <Stack direction="row" spacing={1}>
                                            <Button
                                                size="small"
                                                startIcon={<CheckRounded />}
                                                onClick={() => handleApprove(request.id)}
                                                disabled={isApproving}
                                            >
                                                Duyệt
                                            </Button>
                                            <Button
                                                size="small"
                                                color="error"
                                                startIcon={<CloseRounded />}
                                                onClick={() => setRejectingId(request.id)}
                                            >
                                                Từ chối
                                            </Button>
                                        </Stack>
                                    ) : (
                                        <Typography variant="caption" color="text.secondary">
                                            Không còn thao tác
                                        </Typography>
                                    )}
                                </TableCell>
                            </TableRow>
                        ))}
                        {!requests.length && !requestsQuery.isFetching ? (
                            <TableRow>
                                <TableCell colSpan={8}>
                                    <Typography variant="body2" color="text.secondary">
                                        Chưa có đơn xin phép nào trong phạm vi bộ lọc hiện tại.
                                    </Typography>
                                </TableCell>
                            </TableRow>
                        ) : null}
                    </TableBody>
                </Table>
            </TableContainer>

            <Dialog open={!!rejectingId} onClose={() => setRejectingId(null)} maxWidth="sm" fullWidth>
                <DialogTitle>Từ chối đơn xin phép</DialogTitle>
                <DialogContent>
                    <TextField
                        autoFocus
                        margin="dense"
                        label="Lý do từ chối"
                        value={rejectionReason}
                        onChange={(event) => setRejectionReason(event.target.value)}
                        fullWidth
                        multiline
                        minRows={3}
                    />
                </DialogContent>
                <DialogActions sx={{ p: 2 }}>
                    <Button onClick={() => setRejectingId(null)}>Hủy</Button>
                    <Button variant="contained" color="error" onClick={handleReject} disabled={isRejecting}>
                        {isRejecting ? 'Đang xử lý...' : 'Xác nhận từ chối'}
                    </Button>
                </DialogActions>
            </Dialog>
        </Stack>
    );
}
