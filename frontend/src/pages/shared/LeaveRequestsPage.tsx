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
import { AddRounded, CheckRounded, CloseRounded } from '@mui/icons-material';
import { toast } from 'sonner';

import {
    useApproveLeaveRequestMutation,
    useCancelLeaveRequestMutation,
    useCreateLeaveRequestMutation,
    useGetLeaveRequestsQuery,
    useRejectLeaveRequestMutation,
} from '@/api/leaveApi';
import { useGetClassesQuery } from '@/api/classApi';
import PageHeader from '@/components/common/PageHeader';
import { getApiErrorMessage } from '@/utils/apiError';

type Mode = 'student' | 'staff';

interface LeaveRequestsPageProps {
    mode: Mode;
    title: string;
    subtitle: string;
}

const leaveTypes = [
    { value: 'LEAVE', label: 'Xin nghỉ' },
    { value: 'LATE', label: 'Xin đi muộn' },
    { value: 'EARLY', label: 'Xin về sớm' },
];

export default function LeaveRequestsPage({ mode, title, subtitle }: LeaveRequestsPageProps) {
    const [statusFilter, setStatusFilter] = useState('');
    const [classFilter, setClassFilter] = useState('');
    const [createOpen, setCreateOpen] = useState(false);
    const [rejectingId, setRejectingId] = useState<string | null>(null);
    const [rejectionReason, setRejectionReason] = useState('');
    const [form, setForm] = useState({
        leave_type: 'LEAVE',
        apply_date: '',
        class_id: '',
        subject: '',
        reason: '',
        late_minutes: 0,
        early_minutes: 0,
        documents: '',
    });

    const { data, error, isFetching, refetch } = useGetLeaveRequestsQuery({
        status: statusFilter || undefined,
        class_id: classFilter || undefined,
    });
    const { data: classesResponse } = useGetClassesQuery({ limit: 200 });
    const [createLeaveRequest, { isLoading: isCreating }] = useCreateLeaveRequestMutation();
    const [approveLeaveRequest, { isLoading: isApproving }] = useApproveLeaveRequestMutation();
    const [rejectLeaveRequest, { isLoading: isRejecting }] = useRejectLeaveRequestMutation();
    const [cancelLeaveRequest, { isLoading: isCancelling }] = useCancelLeaveRequestMutation();

    const requests = data?.data?.requests || [];
    const classes = classesResponse?.data?.classes || [];
    const pendingCount = useMemo(() => requests.filter((item) => item.status === 'PENDING').length, [requests]);

    const handleCreate = async () => {
        try {
            await createLeaveRequest({
                leave_type: form.leave_type,
                apply_date: new Date(form.apply_date).toISOString(),
                class_id: form.class_id || undefined,
                subject: form.subject,
                reason: form.reason,
                late_minutes: Number(form.late_minutes) || 0,
                early_minutes: Number(form.early_minutes) || 0,
                documents: form.documents
                    .split(',')
                    .map((item) => item.trim())
                    .filter(Boolean),
            }).unwrap();
            toast.success('Đã tạo đơn xin phép.');
            setCreateOpen(false);
            setForm({
                leave_type: 'LEAVE',
                apply_date: '',
                class_id: '',
                subject: '',
                reason: '',
                late_minutes: 0,
                early_minutes: 0,
                documents: '',
            });
            refetch();
        } catch (err) {
            toast.error(getApiErrorMessage(err, 'Không tạo được đơn xin phép.'));
        }
    };

    const handleApprove = async (id: string) => {
        try {
            await approveLeaveRequest(id).unwrap();
            toast.success('Đã duyệt đơn xin phép.');
            refetch();
        } catch (err) {
            toast.error(getApiErrorMessage(err, 'Không duyệt được đơn xin phép.'));
        }
    };

    const handleReject = async () => {
        if (!rejectingId) return;
        try {
            await rejectLeaveRequest({ id: rejectingId, rejection_reason: rejectionReason }).unwrap();
            toast.success('Đã từ chối đơn xin phép.');
            setRejectingId(null);
            setRejectionReason('');
            refetch();
        } catch (err) {
            toast.error(getApiErrorMessage(err, 'Không từ chối được đơn xin phép.'));
        }
    };

    const handleCancel = async (id: string) => {
        try {
            await cancelLeaveRequest(id).unwrap();
            toast.success('Đã hủy đơn xin phép.');
            refetch();
        } catch (err) {
            toast.error(getApiErrorMessage(err, 'Không hủy được đơn xin phép.'));
        }
    };

    return (
        <Stack sx={{ p: { xs: 2, md: 4 } }} spacing={3}>
            <PageHeader
                title={title}
                subtitle={subtitle}
                actions={
                    mode === 'student' ? (
                        <Button variant="contained" startIcon={<AddRounded />} onClick={() => setCreateOpen(true)}>
                            Tạo đơn
                        </Button>
                    ) : null
                }
            />

            <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 3 }}>
                <Stack direction={{ xs: 'column', md: 'row' }} spacing={2}>
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
                        <MenuItem value="CANCELLED">CANCELLED</MenuItem>
                    </TextField>
                    <TextField
                        select
                        label="Lớp"
                        value={classFilter}
                        onChange={(event) => setClassFilter(event.target.value)}
                        sx={{ minWidth: 220 }}
                    >
                        <MenuItem value="">Tất cả lớp</MenuItem>
                        {classes.map((item) => (
                            <MenuItem key={item.id} value={item.id}>
                                {item.name} ({item.code})
                            </MenuItem>
                        ))}
                    </TextField>
                    <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap alignItems="center">
                        <Chip size="small" label={`${requests.length} đơn`} variant="outlined" />
                        <Chip size="small" label={`Pending ${pendingCount}`} color="warning" variant="outlined" />
                    </Stack>
                </Stack>
            </Paper>

            {error ? (
                <Alert severity="error">
                    {getApiErrorMessage(error, 'Không tải được danh sách đơn xin phép.')}
                </Alert>
            ) : null}

            <TableContainer component={Paper} variant="outlined" sx={{ borderRadius: 3 }}>
                <Table>
                    <TableHead>
                        <TableRow>
                            <TableCell sx={{ fontWeight: 700 }}>Người gửi</TableCell>
                            <TableCell sx={{ fontWeight: 700 }}>Loại đơn</TableCell>
                            <TableCell sx={{ fontWeight: 700 }}>Ngày áp dụng</TableCell>
                            <TableCell sx={{ fontWeight: 700 }}>Lớp</TableCell>
                            <TableCell sx={{ fontWeight: 700 }}>Lý do</TableCell>
                            <TableCell sx={{ fontWeight: 700 }}>Trạng thái</TableCell>
                            <TableCell sx={{ fontWeight: 700 }}>Thao tác</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {requests.map((request) => (
                            <TableRow key={request.id} hover>
                                <TableCell>
                                    <Typography variant="body2" sx={{ fontWeight: 600 }}>
                                        {request.student?.full_name || 'Học sinh'}
                                    </Typography>
                                    <Typography variant="caption" color="text.secondary">
                                        {request.student?.code}
                                    </Typography>
                                </TableCell>
                                <TableCell>{request.leave_type}</TableCell>
                                <TableCell>{new Date(request.apply_date).toLocaleString('vi-VN')}</TableCell>
                                <TableCell>{request.class?.name || 'Không gắn lớp'}</TableCell>
                                <TableCell>{request.reason}</TableCell>
                                <TableCell>
                                    <Chip size="small" label={request.status} color={request.status === 'APPROVED' ? 'success' : request.status === 'REJECTED' ? 'error' : request.status === 'PENDING' ? 'warning' : 'default'} variant="outlined" />
                                </TableCell>
                                <TableCell>
                                    {mode === 'staff' && request.status === 'PENDING' ? (
                                        <Stack direction="row" spacing={1}>
                                            <Button size="small" startIcon={<CheckRounded />} onClick={() => handleApprove(request.id)} disabled={isApproving}>
                                                Duyệt
                                            </Button>
                                            <Button size="small" color="error" startIcon={<CloseRounded />} onClick={() => setRejectingId(request.id)}>
                                                Từ chối
                                            </Button>
                                        </Stack>
                                    ) : null}
                                    {mode === 'student' && request.status === 'PENDING' ? (
                                        <Button size="small" color="error" startIcon={<CloseRounded />} onClick={() => handleCancel(request.id)} disabled={isCancelling}>
                                            Hủy
                                        </Button>
                                    ) : null}
                                </TableCell>
                            </TableRow>
                        ))}
                        {!requests.length && !isFetching ? (
                            <TableRow>
                                <TableCell colSpan={7}>
                                    <Typography variant="body2" color="text.secondary">
                                        Chưa có đơn xin phép nào trong phạm vi bộ lọc hiện tại.
                                    </Typography>
                                </TableCell>
                            </TableRow>
                        ) : null}
                    </TableBody>
                </Table>
            </TableContainer>

            <Dialog open={createOpen} onClose={() => setCreateOpen(false)} maxWidth="sm" fullWidth>
                <DialogTitle>Tạo đơn xin phép</DialogTitle>
                <DialogContent>
                    <Stack spacing={2} sx={{ mt: 1 }}>
                        <TextField
                            select
                            label="Loại đơn"
                            value={form.leave_type}
                            onChange={(event) => setForm((current) => ({ ...current, leave_type: event.target.value }))}
                            fullWidth
                        >
                            {leaveTypes.map((item) => (
                                <MenuItem key={item.value} value={item.value}>
                                    {item.label}
                                </MenuItem>
                            ))}
                        </TextField>
                        <TextField
                            label="Ngày áp dụng"
                            type="datetime-local"
                            value={form.apply_date}
                            onChange={(event) => setForm((current) => ({ ...current, apply_date: event.target.value }))}
                            fullWidth
                            InputLabelProps={{ shrink: true }}
                        />
                        <TextField
                            select
                            label="Lớp"
                            value={form.class_id}
                            onChange={(event) => setForm((current) => ({ ...current, class_id: event.target.value }))}
                            fullWidth
                        >
                            <MenuItem value="">Không gắn lớp cụ thể</MenuItem>
                            {classes.map((item) => (
                                <MenuItem key={item.id} value={item.id}>
                                    {item.name} ({item.code})
                                </MenuItem>
                            ))}
                        </TextField>
                        <TextField
                            label="Chủ đề"
                            value={form.subject}
                            onChange={(event) => setForm((current) => ({ ...current, subject: event.target.value }))}
                            fullWidth
                        />
                        <TextField
                            label="Lý do"
                            value={form.reason}
                            onChange={(event) => setForm((current) => ({ ...current, reason: event.target.value }))}
                            fullWidth
                            multiline
                            minRows={3}
                        />
                        <TextField
                            label="Tài liệu đính kèm (cách nhau bằng dấu phẩy)"
                            value={form.documents}
                            onChange={(event) => setForm((current) => ({ ...current, documents: event.target.value }))}
                            fullWidth
                        />
                    </Stack>
                </DialogContent>
                <DialogActions sx={{ p: 2 }}>
                    <Button onClick={() => setCreateOpen(false)}>Hủy</Button>
                    <Button variant="contained" onClick={handleCreate} disabled={isCreating}>
                        {isCreating ? 'Đang tạo...' : 'Tạo đơn'}
                    </Button>
                </DialogActions>
            </Dialog>

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
