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
import { AddRounded, CloseRounded } from '@mui/icons-material';
import { toast } from 'sonner';

import {
    useCancelStudentLeaveRequestMutation,
    useCreateStudentLeaveRequestMutation,
    useGetMyStudentLeaveRequestsQuery,
    useGetStudentTimetableQuery,
} from '@/api/studentPortalApi';
import PageHeader from '@/components/common/PageHeader';
import { getApiErrorMessage } from '@/utils/apiError';

const leaveTypes = [
    { value: 'LEAVE', label: 'Xin nghỉ' },
    { value: 'LATE', label: 'Xin đi muộn' },
    { value: 'EARLY', label: 'Xin về sớm' },
];

export default function StudentLeavesPage() {
    const [statusFilter, setStatusFilter] = useState('');
    const [classFilter, setClassFilter] = useState('');
    const [createOpen, setCreateOpen] = useState(false);
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

    const { data, error, isFetching, refetch } = useGetMyStudentLeaveRequestsQuery({
        status: statusFilter || undefined,
        class_id: classFilter || undefined,
    });
    const { data: timetableResponse } = useGetStudentTimetableQuery();
    const [createLeaveRequest, { isLoading: isCreating }] = useCreateStudentLeaveRequestMutation();
    const [cancelLeaveRequest, { isLoading: isCancelling }] = useCancelStudentLeaveRequestMutation();

    const requests = data?.data?.requests ?? [];
    const classes = useMemo(() => {
        const map = new Map<string, { id: string; code: string; name: string }>();
        (timetableResponse?.data?.lessons ?? []).forEach((lesson) => {
            if (!map.has(lesson.class_id)) {
                map.set(lesson.class_id, {
                    id: lesson.class_id,
                    code: lesson.class_code,
                    name: lesson.class_name,
                });
            }
        });
        requests.forEach((request) => {
            if (request.class && !map.has(request.class.id)) {
                map.set(request.class.id, request.class);
            }
        });
        return Array.from(map.values()).sort((a, b) => a.name.localeCompare(b.name));
    }, [requests, timetableResponse]);

    const pendingCount = useMemo(() => requests.filter((item) => item.status === 'PENDING').length, [requests]);

    const handleCreate = async () => {
        try {
            await createLeaveRequest({
                leave_type: form.leave_type,
                apply_date: new Date(form.apply_date).toISOString(),
                class_id: form.class_id || undefined,
                subject: form.subject.trim() || undefined,
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
                title="Đơn xin phép"
                subtitle="Tạo, theo dõi và hủy đơn xin nghỉ hoặc xin đi muộn của bạn theo đúng luồng học sinh."
                actions={(
                    <Button variant="contained" startIcon={<AddRounded />} onClick={() => setCreateOpen(true)}>
                        Tạo đơn
                    </Button>
                )}
                breadcrumbs={[
                    { label: 'Cổng học sinh', path: '/app/student/overview' },
                    { label: 'Đơn xin phép' },
                ]}
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
                        sx={{ minWidth: 240 }}
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
                    {getApiErrorMessage(error, 'Không tải được danh sách đơn xin phép của bạn.')}
                </Alert>
            ) : null}

            <TableContainer component={Paper} variant="outlined" sx={{ borderRadius: 3 }}>
                <Table>
                    <TableHead>
                        <TableRow>
                            <TableCell sx={{ fontWeight: 700 }}>Loại đơn</TableCell>
                            <TableCell sx={{ fontWeight: 700 }}>Ngày áp dụng</TableCell>
                            <TableCell sx={{ fontWeight: 700 }}>Lớp</TableCell>
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
                                        {request.leave_type}
                                    </Typography>
                                    <Typography variant="caption" color="text.secondary">
                                        {request.subject || 'Không có chủ đề'}
                                    </Typography>
                                </TableCell>
                                <TableCell>{new Date(request.apply_date).toLocaleString('vi-VN')}</TableCell>
                                <TableCell>{request.class?.name || 'Không gắn lớp'}</TableCell>
                                <TableCell>
                                    <Typography variant="body2">{request.reason}</Typography>
                                    {request.rejection_reason ? (
                                        <Typography variant="caption" color="error">
                                            Lý do từ chối: {request.rejection_reason}
                                        </Typography>
                                    ) : null}
                                </TableCell>
                                <TableCell>
                                    {request.documents?.length ? (
                                        <Stack spacing={0.5}>
                                            {request.documents.map((doc) => (
                                                <Typography key={doc} variant="caption" color="primary">
                                                    {doc}
                                                </Typography>
                                            ))}
                                        </Stack>
                                    ) : (
                                        <Typography variant="body2" color="text.secondary">
                                            Không có
                                        </Typography>
                                    )}
                                </TableCell>
                                <TableCell>
                                    <Chip
                                        size="small"
                                        label={request.status}
                                        color={
                                            request.status === 'APPROVED'
                                                ? 'success'
                                                : request.status === 'REJECTED'
                                                    ? 'error'
                                                    : request.status === 'PENDING'
                                                        ? 'warning'
                                                        : 'default'
                                        }
                                        variant="outlined"
                                    />
                                </TableCell>
                                <TableCell>
                                    {request.status === 'PENDING' ? (
                                        <Button
                                            size="small"
                                            color="error"
                                            startIcon={<CloseRounded />}
                                            onClick={() => handleCancel(request.id)}
                                            disabled={isCancelling}
                                        >
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

                        {(form.leave_type === 'LATE' || form.leave_type === 'EARLY') ? (
                            <TextField
                                label={form.leave_type === 'LATE' ? 'Số phút đi muộn' : 'Số phút về sớm'}
                                type="number"
                                value={form.leave_type === 'LATE' ? form.late_minutes : form.early_minutes}
                                onChange={(event) => setForm((current) => ({
                                    ...current,
                                    late_minutes: form.leave_type === 'LATE' ? Number(event.target.value) || 0 : current.late_minutes,
                                    early_minutes: form.leave_type === 'EARLY' ? Number(event.target.value) || 0 : current.early_minutes,
                                }))}
                                fullWidth
                            />
                        ) : null}

                        <TextField
                            label="Tài liệu đính kèm (URL, cách nhau bằng dấu phẩy)"
                            value={form.documents}
                            onChange={(event) => setForm((current) => ({ ...current, documents: event.target.value }))}
                            fullWidth
                            multiline
                            minRows={2}
                        />
                    </Stack>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setCreateOpen(false)}>Đóng</Button>
                    <Button
                        variant="contained"
                        onClick={handleCreate}
                        disabled={!form.apply_date || !form.reason.trim() || isCreating}
                    >
                        Tạo đơn
                    </Button>
                </DialogActions>
            </Dialog>
        </Stack>
    );
}
