import { useEffect, useMemo, useState } from 'react';
import {
    Alert,
    Box,
    Button,
    Chip,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    MenuItem,
    Paper,
    Select,
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
import { AssignmentTurnedInRounded, SaveRounded } from '@mui/icons-material';
import { toast } from 'sonner';

import {
    useGetLessonAttendanceQuery,
    useUpsertLessonAttendanceMutation,
} from '@/api/lessonApi';
import { useLazyFindMakeupSpotsQuery, type MakeupSpot } from '@/api/schedulingApi';
import { getApiErrorMessage } from '@/utils/apiError';

const ATTENDANCE_OPTIONS = [
    { value: 1, label: 'Có mặt' },
    { value: 2, label: 'Vắng' },
    { value: 3, label: 'Có phép' },
    { value: 4, label: 'Muộn' },
    { value: 5, label: 'Về sớm' },
];

type EditableAttendanceRow = {
    student_id: string;
    student_code: string;
    student_name: string;
    status: number;
    note: string;
};

interface LessonAttendanceManagerProps {
    lessonId: string;
    title?: string;
    subtitle?: string;
}

export default function LessonAttendanceManager({
    lessonId,
    title = 'Điểm danh buổi học',
    subtitle = 'Cập nhật trạng thái chuyên cần cho từng học sinh trong buổi học này.',
}: LessonAttendanceManagerProps) {
    const { data, isLoading, isFetching, error, refetch } = useGetLessonAttendanceQuery(lessonId, {
        skip: !lessonId,
    });
    const [upsertAttendance, { isLoading: isSaving }] = useUpsertLessonAttendanceMutation();
    const [rows, setRows] = useState<EditableAttendanceRow[]>([]);
    const [selectedStudent, setSelectedStudent] = useState<EditableAttendanceRow | null>(null);
    const [findMakeupSpots, { data: makeupData, isFetching: isFindingMakeup, error: makeupError }] =
        useLazyFindMakeupSpotsQuery();

    useEffect(() => {
        const nextRows = (data?.data?.records || []).map((record) => ({
            student_id: record.student.id,
            student_code: record.student.code,
            student_name: record.student.full_name,
            status: record.status,
            note: record.note || '',
        }));
        setRows(nextRows);
    }, [data]);

    const stats = useMemo(() => {
        const counters = {
            total: rows.length,
            present: 0,
            absent: 0,
            excused: 0,
            late: 0,
            early: 0,
            unmarked: 0,
        };
        rows.forEach((row) => {
            switch (row.status) {
                case 1:
                    counters.present += 1;
                    break;
                case 2:
                    counters.absent += 1;
                    break;
                case 3:
                    counters.excused += 1;
                    break;
                case 4:
                    counters.late += 1;
                    break;
                case 5:
                    counters.early += 1;
                    break;
                default:
                    counters.unmarked += 1;
            }
        });
        return counters;
    }, [rows]);

    const handleSave = async () => {
        if (!rows.length) {
            toast.error('Buổi học này chưa có học sinh để điểm danh.');
            return;
        }
        if (rows.some((row) => row.status === 0)) {
            toast.error('Vẫn còn học sinh chưa chọn trạng thái điểm danh.');
            return;
        }

        try {
            await upsertAttendance({
                id: lessonId,
                body: {
                    records: rows.map((row) => ({
                        student_id: row.student_id,
                        status: row.status,
                        note: row.note?.trim() || '',
                    })),
                },
            }).unwrap();
            toast.success('Đã lưu điểm danh buổi học.');
            refetch();
        } catch (err) {
            toast.error(getApiErrorMessage(err, 'Không lưu được điểm danh.'));
        }
    };

    const formatDateTime = (value: string) => {
        const parsed = new Date(value);
        if (Number.isNaN(parsed.getTime())) {
            return value;
        }
        return parsed.toLocaleString('vi-VN', {
            hour: '2-digit',
            minute: '2-digit',
            day: '2-digit',
            month: '2-digit',
            year: 'numeric',
        });
    };

    const matchTypeLabel = (value: string) => {
        switch (value) {
            case 'same_course':
                return 'Cùng course';
            case 'same_subject_grade':
                return 'Cùng môn/khối';
            default:
                return value;
        }
    };

    const handleFindMakeup = async (row: EditableAttendanceRow) => {
        setSelectedStudent(row);
        try {
            await findMakeupSpots({ lessonId, studentId: row.student_id, limit: 8 }).unwrap();
        } catch {
            // Error is surfaced in the dialog body.
        }
    };

    const makeupSpots = makeupData?.data?.spots || [];

    return (
        <>
            <Paper variant="outlined" sx={{ p: 3, borderRadius: 3 }}>
                <Stack spacing={2.5}>
                <Stack
                    direction={{ xs: 'column', md: 'row' }}
                    justifyContent="space-between"
                    alignItems={{ xs: 'flex-start', md: 'center' }}
                    spacing={1.5}
                >
                    <Stack spacing={0.5}>
                        <Stack direction="row" spacing={1} alignItems="center">
                            <AssignmentTurnedInRounded color="primary" />
                            <Typography variant="h6" sx={{ fontWeight: 700 }}>
                                {title}
                            </Typography>
                        </Stack>
                        <Typography variant="body2" color="text.secondary">
                            {subtitle}
                        </Typography>
                    </Stack>

                    <Button
                        variant="contained"
                        startIcon={<SaveRounded />}
                        onClick={handleSave}
                        disabled={isSaving || isLoading || isFetching || !rows.length}
                    >
                        {isSaving ? 'Đang lưu...' : 'Lưu điểm danh'}
                    </Button>
                </Stack>

                <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                    <Chip size="small" label={`${stats.total} học sinh`} variant="outlined" />
                    <Chip size="small" label={`Có mặt ${stats.present}`} color="success" variant="outlined" />
                    <Chip size="small" label={`Vắng ${stats.absent}`} color="error" variant="outlined" />
                    <Chip size="small" label={`Có phép ${stats.excused}`} color="info" variant="outlined" />
                    <Chip size="small" label={`Muộn ${stats.late}`} color="warning" variant="outlined" />
                    <Chip size="small" label={`Về sớm ${stats.early}`} color="secondary" variant="outlined" />
                    {stats.unmarked > 0 ? (
                        <Chip size="small" label={`Chưa chấm ${stats.unmarked}`} color="default" variant="outlined" />
                    ) : null}
                </Stack>

                {error ? (
                    <Alert severity="error">
                        {getApiErrorMessage(error, 'Không tải được dữ liệu điểm danh của buổi học.')}
                    </Alert>
                ) : null}

                    <TableContainer>
                        <Table size="small">
                            <TableHead>
                                <TableRow>
                                    <TableCell sx={{ fontWeight: 700 }}>Học sinh</TableCell>
                                    <TableCell sx={{ fontWeight: 700, width: 180 }}>Trạng thái</TableCell>
                                    <TableCell sx={{ fontWeight: 700 }}>Ghi chú</TableCell>
                                    <TableCell sx={{ fontWeight: 700, width: 150 }}>Học bù</TableCell>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {rows.map((row, index) => (
                                    <TableRow key={row.student_id} hover>
                                        <TableCell>
                                            <Typography variant="body2" sx={{ fontWeight: 600 }}>
                                                {row.student_name}
                                            </Typography>
                                            <Typography variant="caption" color="text.secondary">
                                                {row.student_code}
                                            </Typography>
                                        </TableCell>
                                        <TableCell>
                                            <Select
                                                size="small"
                                                fullWidth
                                                value={String(row.status)}
                                                displayEmpty
                                                onChange={(event) => {
                                                    const nextStatus = Number(event.target.value);
                                                    setRows((current) =>
                                                        current.map((item, currentIndex) =>
                                                            currentIndex === index ? { ...item, status: nextStatus } : item,
                                                        ),
                                                    );
                                                }}
                                            >
                                                <MenuItem value="0">Chọn trạng thái</MenuItem>
                                                {ATTENDANCE_OPTIONS.map((option) => (
                                                    <MenuItem key={option.value} value={String(option.value)}>
                                                        {option.label}
                                                    </MenuItem>
                                                ))}
                                            </Select>
                                        </TableCell>
                                        <TableCell>
                                            <TextField
                                                size="small"
                                                fullWidth
                                                placeholder="Ghi chú thêm nếu cần"
                                                value={row.note}
                                                onChange={(event) => {
                                                    const nextNote = event.target.value;
                                                    setRows((current) =>
                                                        current.map((item, currentIndex) =>
                                                            currentIndex === index ? { ...item, note: nextNote } : item,
                                                        ),
                                                    );
                                                }}
                                            />
                                        </TableCell>
                                        <TableCell>
                                            <Button size="small" onClick={() => void handleFindMakeup(row)}>
                                                Tìm slot
                                            </Button>
                                        </TableCell>
                                    </TableRow>
                                ))}
                                {!rows.length && !isLoading ? (
                                    <TableRow>
                                        <TableCell colSpan={4}>
                                            <Typography variant="body2" color="text.secondary">
                                                Buổi học này chưa có dữ liệu học viên để điểm danh.
                                            </Typography>
                                        </TableCell>
                                    </TableRow>
                                ) : null}
                            </TableBody>
                        </Table>
                    </TableContainer>
                </Stack>
            </Paper>

            <Dialog open={!!selectedStudent} onClose={() => setSelectedStudent(null)} fullWidth maxWidth="md">
                <DialogTitle>Tìm slot học bù</DialogTitle>
                <DialogContent>
                    <Stack spacing={2.5} sx={{ pt: 1 }}>
                        {selectedStudent ? (
                            <Alert severity="info">
                                Đang tra cứu chỗ trống cho <strong>{selectedStudent.student_name}</strong>. Hệ thống lọc theo cùng course hoặc cùng môn/khối, rồi chặn các slot trùng lịch, quá sức chứa hoặc không đủ thời gian di chuyển.
                            </Alert>
                        ) : null}

                        {makeupError ? (
                            <Alert severity="error">
                                {getApiErrorMessage(makeupError, 'Không tìm được slot học bù phù hợp.')}
                            </Alert>
                        ) : null}

                        {!makeupError && !isFindingMakeup && !makeupSpots.length ? (
                            <Alert severity="warning">
                                Chưa tìm thấy slot học bù phù hợp cho học viên này trong cửa sổ tìm kiếm hiện tại.
                            </Alert>
                        ) : null}

                        {isFindingMakeup ? (
                            <Typography variant="body2" color="text.secondary">
                                Đang quét các slot học bù tương thích...
                            </Typography>
                        ) : null}

                        {makeupSpots.map((spot: MakeupSpot) => (
                            <Paper key={spot.lesson_id} variant="outlined" sx={{ p: 2, borderRadius: 2 }}>
                                <Stack spacing={1.25}>
                                    <Stack direction={{ xs: 'column', sm: 'row' }} justifyContent="space-between" spacing={1}>
                                        <Box>
                                            <Typography variant="subtitle2" sx={{ fontWeight: 700 }}>
                                                {spot.class_name} ({spot.class_code})
                                            </Typography>
                                            <Typography variant="body2" color="text.secondary">
                                                {formatDateTime(spot.start_time)} - {formatDateTime(spot.end_time)}
                                            </Typography>
                                        </Box>
                                        <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                                            <Chip size="small" color="primary" label={matchTypeLabel(spot.match_type)} />
                                            <Chip
                                                size="small"
                                                color={spot.remaining_capacity > 1 ? 'success' : 'warning'}
                                                label={`Còn ${spot.remaining_capacity}/${spot.capacity_limit} chỗ`}
                                            />
                                        </Stack>
                                    </Stack>

                                    <Typography variant="body2">
                                        {spot.teacher_name || 'Chưa có giáo viên'} • {spot.room_name || 'Chưa xếp phòng'}
                                    </Typography>
                                    <Typography variant="caption" color="text.secondary">
                                        Utilization: {(spot.capacity_utilization * 100).toFixed(0)}%
                                    </Typography>

                                    {spot.reasons.length ? (
                                        <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                                            {spot.reasons.map((reason) => (
                                                <Chip key={`${spot.lesson_id}-${reason}`} size="small" variant="outlined" label={reason} />
                                            ))}
                                        </Stack>
                                    ) : null}
                                </Stack>
                            </Paper>
                        ))}

                        <Alert severity="info">
                            Bản hiện tại hỗ trợ tra cứu slot học bù để giáo vụ quyết định nhanh. Bước “xác nhận nhét học viên vào một lesson đơn lẻ” sẽ cần model attendance makeup riêng để không làm lệch roster lớp gốc.
                        </Alert>
                    </Stack>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setSelectedStudent(null)}>Đóng</Button>
                </DialogActions>
            </Dialog>
        </>
    );
}
