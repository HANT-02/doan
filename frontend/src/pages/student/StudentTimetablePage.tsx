import { useMemo, useState } from 'react';
import {
    Alert,
    Box,
    Button,
    Chip,
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
import {
    CalendarMonthRounded,
    ChevronLeftRounded,
    ChevronRightRounded,
    EventAvailableRounded,
    TodayRounded,
} from '@mui/icons-material';
import { addDays, addWeeks, endOfWeek, format, isSameDay, parseISO, startOfWeek } from 'date-fns';
import { vi } from 'date-fns/locale';

import {
    useGetStudentAttendanceQuery,
    useGetMyStudentLeaveRequestsQuery,
    useGetStudentTimetableQuery,
    type StudentAttendanceRecord,
    type StudentPortalLeaveRequest,
    type StudentTimetableLesson,
} from '@/api/studentPortalApi';
import PageHeader from '@/components/common/PageHeader';
import { getApiErrorMessage } from '@/utils/apiError';

const weekLabels = ['Thứ Hai', 'Thứ Ba', 'Thứ Tư', 'Thứ Năm', 'Thứ Sáu', 'Thứ Bảy', 'Chủ Nhật'];

const palette = [
    { bg: 'rgba(14, 116, 144, 0.10)', border: '#0ea5e9' },
    { bg: 'rgba(5, 150, 105, 0.10)', border: '#10b981' },
    { bg: 'rgba(234, 88, 12, 0.10)', border: '#f97316' },
    { bg: 'rgba(190, 24, 93, 0.10)', border: '#ec4899' },
    { bg: 'rgba(79, 70, 229, 0.10)', border: '#6366f1' },
];

function buildWeekDays(weekStart: Date) {
    return Array.from({ length: 7 }, (_, index) => addDays(weekStart, index));
}

function getAttendanceStatus(status?: number) {
    switch (status) {
        case 1:
            return { label: 'Có mặt', color: 'success' as const };
        case 0:
            return { label: 'Vắng', color: 'error' as const };
        case 2:
            return { label: 'Muộn', color: 'warning' as const };
        case 3:
            return { label: 'Xin phép', color: 'info' as const };
        default:
            return { label: 'Chưa điểm danh', color: 'default' as const };
    }
}

function getLessonStatus(lesson: StudentTimetableLesson) {
    const now = new Date();
    const start = parseISO(lesson.date_start);
    const end = parseISO(lesson.date_end);

    if (end < now) {
        return { label: 'Đã qua', color: 'default' as const };
    }
    if (start > now) {
        return { label: 'Sắp tới', color: 'primary' as const };
    }
    return { label: 'Đang diễn ra', color: 'success' as const };
}

function findMatchingLeave(lesson: StudentTimetableLesson, requests: StudentPortalLeaveRequest[]) {
    return requests.find((request) => {
        if (request.lesson?.id && request.lesson.id === lesson.id) {
            return true;
        }

        const requestDate = request.apply_date ? new Date(request.apply_date) : null;
        if (!requestDate) {
            return false;
        }

        return request.class?.id === lesson.class_id && isSameDay(requestDate, parseISO(lesson.date_start));
    });
}

function getLessonTone(lesson: StudentTimetableLesson) {
    const seed = lesson.class_id
        .split('')
        .reduce((sum, char) => sum + char.charCodeAt(0), 0);
    return palette[seed % palette.length];
}

export default function StudentTimetablePage() {
    const [weekAnchor, setWeekAnchor] = useState(() => startOfWeek(new Date(), { weekStartsOn: 1 }));
    const [classFilter, setClassFilter] = useState('');

    const weekStart = useMemo(() => startOfWeek(weekAnchor, { weekStartsOn: 1 }), [weekAnchor]);
    const weekEnd = useMemo(() => endOfWeek(weekAnchor, { weekStartsOn: 1 }), [weekAnchor]);

    const {
        data: timetableResponse,
        isLoading: isLoadingTimetable,
        isFetching: isFetchingTimetable,
        error: timetableError,
    } = useGetStudentTimetableQuery({
        class_id: classFilter || undefined,
        from: format(weekStart, 'yyyy-MM-dd'),
        to: format(weekEnd, 'yyyy-MM-dd'),
    });

    const {
        data: leaveResponse,
        error: leaveError,
    } = useGetMyStudentLeaveRequestsQuery({
        class_id: classFilter || undefined,
    });
    const {
        data: attendanceResponse,
        isLoading: isLoadingAttendance,
        error: attendanceError,
    } = useGetStudentAttendanceQuery({
        class_id: classFilter || undefined,
    });

    const lessons = timetableResponse?.data?.lessons ?? [];
    const leaveRequests = leaveResponse?.data?.requests ?? [];
    const attendanceSummary = attendanceResponse?.data?.summary;
    const attendanceRecords = attendanceResponse?.data?.records ?? [];

    const availableClasses = useMemo(() => {
        const source = attendanceRecords.length ? attendanceRecords.map((record) => record.lesson) : lessons;
        const map = new Map<string, { id: string; code: string; name: string }>();
        source.forEach((lesson) => {
            if (!map.has(lesson.class_id)) {
                map.set(lesson.class_id, {
                    id: lesson.class_id,
                    code: lesson.class_code,
                    name: lesson.class_name,
                });
            }
        });
        return Array.from(map.values()).sort((a, b) => a.name.localeCompare(b.name));
    }, [attendanceRecords, lessons]);

    const sortedAttendanceRecords = useMemo<StudentAttendanceRecord[]>(
        () => [...attendanceRecords].sort((a, b) => new Date(b.lesson.date_start).getTime() - new Date(a.lesson.date_start).getTime()),
        [attendanceRecords],
    );

    const calendarDays = useMemo(
        () =>
            buildWeekDays(weekStart).map((day) => ({
                date: day,
                lessons: lessons
                    .filter((lesson) => isSameDay(parseISO(lesson.date_start), day))
                    .sort((left, right) => left.date_start.localeCompare(right.date_start)),
            })),
        [lessons, weekStart],
    );

    const uniqueClasses = useMemo(() => new Set(lessons.map((lesson) => lesson.class_id)).size, [lessons]);
    const linkedLeaveCount = useMemo(
        () => lessons.filter((lesson) => findMatchingLeave(lesson, leaveRequests)).length,
        [lessons, leaveRequests],
    );

    return (
        <Stack sx={{ p: { xs: 2, md: 4 } }} spacing={3}>
            <PageHeader
                title="Thời khóa biểu"
                subtitle="Theo dõi lịch học theo tuần ở dạng calendar, xem nhanh ca học, giáo viên, phòng học và các đơn xin phép liên quan."
                icon={<EventAvailableRounded />}
                breadcrumbs={[
                    { label: 'Cổng học sinh', path: '/app/student/overview' },
                    { label: 'Thời khóa biểu' },
                ]}
                actions={(
                    <Stack direction="row" spacing={1}>
                        <Button variant="outlined" startIcon={<ChevronLeftRounded />} onClick={() => setWeekAnchor((current) => addWeeks(current, -1))}>
                            Tuần trước
                        </Button>
                        <Button variant="outlined" startIcon={<TodayRounded />} onClick={() => setWeekAnchor(startOfWeek(new Date(), { weekStartsOn: 1 }))}>
                            Tuần này
                        </Button>
                        <Button variant="outlined" endIcon={<ChevronRightRounded />} onClick={() => setWeekAnchor((current) => addWeeks(current, 1))}>
                            Tuần sau
                        </Button>
                    </Stack>
                )}
            />

            <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                <Chip label={`${lessons.length} buổi học trong tuần`} variant="outlined" />
                <Chip label={`${uniqueClasses} lớp đang theo học`} color="primary" variant="outlined" />
                <Chip label={`${linkedLeaveCount} buổi có đơn xin phép`} color="warning" variant="outlined" />
                <Chip
                    label={`Chuyên cần ${attendanceSummary ? `${Math.round(attendanceSummary.attendance_rate * 100)}%` : '--'}`}
                    color="success"
                    variant="outlined"
                />
                <Chip
                    icon={<CalendarMonthRounded />}
                    label={`${format(weekStart, 'dd/MM', { locale: vi })} - ${format(weekEnd, 'dd/MM/yyyy', { locale: vi })}`}
                    color="secondary"
                    variant="outlined"
                />
            </Stack>

            <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 3 }}>
                <Stack direction={{ xs: 'column', md: 'row' }} spacing={2} alignItems={{ xs: 'stretch', md: 'center' }}>
                    <TextField
                        select
                        label="Lọc theo lớp"
                        value={classFilter}
                        onChange={(event) => setClassFilter(event.target.value)}
                        sx={{ minWidth: 260 }}
                    >
                        <MenuItem value="">Tất cả lớp</MenuItem>
                        {availableClasses.map((item) => (
                            <MenuItem key={item.id} value={item.id}>
                                {item.name} ({item.code})
                            </MenuItem>
                        ))}
                    </TextField>
                    <Typography variant="body2" color="text.secondary">
                        Bộ lọc này áp dụng đồng thời cho calendar, đơn xin phép và thống kê chuyên cần.
                    </Typography>
                </Stack>
            </Paper>

            {timetableError ? (
                <Alert severity="error">
                    {getApiErrorMessage(timetableError, 'Không tải được thời khóa biểu của bạn.')}
                </Alert>
            ) : null}

            {leaveError ? (
                <Alert severity="warning">
                    {getApiErrorMessage(leaveError, 'Không tải được dữ liệu đơn xin phép liên quan.')}
                </Alert>
            ) : null}

            {attendanceError ? (
                <Alert severity="error">
                    {getApiErrorMessage(attendanceError, 'Không tải được dữ liệu điểm danh của bạn.')}
                </Alert>
            ) : null}

            {attendanceSummary?.warning ? (
                <Alert severity="warning">
                    {attendanceSummary.warning_message || 'Tỷ lệ vắng đã vượt ngưỡng cảnh báo. Hãy liên hệ giáo viên phụ trách để được hỗ trợ.'}
                </Alert>
            ) : null}

            <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 3 }}>
                <Stack spacing={2}>
                    <Typography variant="h6" sx={{ fontWeight: 700 }}>
                        Tổng hợp chuyên cần
                    </Typography>
                    <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                        <Chip label={`Tổng buổi ${attendanceSummary?.total_lessons ?? 0}`} variant="outlined" />
                        <Chip label={`Có mặt ${attendanceSummary?.present_count ?? 0}`} color="success" variant="outlined" />
                        <Chip label={`Vắng ${attendanceSummary?.absent_count ?? 0}`} color="error" variant="outlined" />
                        <Chip label={`Muộn ${attendanceSummary?.late_count ?? 0}`} color="warning" variant="outlined" />
                        <Chip label={`Xin phép ${attendanceSummary?.excused_count ?? 0}`} color="info" variant="outlined" />
                        <Chip label={`Chưa chấm ${attendanceSummary?.unmarked_count ?? 0}`} variant="outlined" />
                    </Stack>
                    <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                        <Chip
                            label={`Tỷ lệ chuyên cần ${attendanceSummary ? `${Math.round(attendanceSummary.attendance_rate * 100)}%` : '--'}`}
                            color="primary"
                            variant="outlined"
                        />
                        <Chip
                            label={`Tỷ lệ vắng ${attendanceSummary ? `${Math.round(attendanceSummary.absent_rate * 100)}%` : '--'}`}
                            color={attendanceSummary?.warning ? 'warning' : 'default'}
                            variant="outlined"
                        />
                    </Stack>
                </Stack>
            </Paper>

            <Paper variant="outlined" sx={{ borderRadius: 3, overflow: 'hidden' }}>
                <Box
                    sx={{
                        px: 2.5,
                        py: 2,
                        borderBottom: '1px solid rgba(15,23,42,0.08)',
                        background: 'linear-gradient(135deg, rgba(14,165,233,0.08), rgba(16,185,129,0.05))',
                    }}
                >
                    <Stack
                        direction={{ xs: 'column', md: 'row' }}
                        justifyContent="space-between"
                        alignItems={{ xs: 'flex-start', md: 'center' }}
                        spacing={1}
                    >
                        <Box>
                            <Typography variant="h6" sx={{ fontWeight: 700 }}>
                                Calendar tuần học
                            </Typography>
                            <Typography variant="body2" color="text.secondary">
                                Theo dõi các buổi học theo ngày, kèm trạng thái điểm danh và đơn xin phép nếu có.
                            </Typography>
                        </Box>
                        {isFetchingTimetable ? (
                            <Chip size="small" label="Đang tải..." color="info" variant="outlined" />
                        ) : null}
                    </Stack>
                </Box>

                <Box
                    sx={{
                        display: 'grid',
                        gridTemplateColumns: { xs: '1fr', lg: 'repeat(7, minmax(0, 1fr))' },
                        borderBottom: '1px solid rgba(15,23,42,0.08)',
                        backgroundColor: '#f8fafc',
                    }}
                >
                    {weekLabels.map((label) => (
                        <Box key={label} sx={{ px: 1.5, py: 1.25, borderLeft: { lg: '1px solid rgba(15,23,42,0.08)' } }}>
                            <Typography variant="subtitle2" sx={{ fontWeight: 700 }}>
                                {label}
                            </Typography>
                        </Box>
                    ))}
                </Box>

                <Box
                    sx={{
                        display: 'grid',
                        gridTemplateColumns: { xs: '1fr', lg: 'repeat(7, minmax(0, 1fr))' },
                    }}
                >
                    {calendarDays.map((day, index) => (
                        <Box
                            key={day.date.toISOString()}
                            sx={{
                                minHeight: 280,
                                p: 1.25,
                                borderLeft: { xs: 'none', lg: index === 0 ? 'none' : '1px solid rgba(15,23,42,0.08)' },
                                borderTop: '1px solid rgba(15,23,42,0.08)',
                            }}
                        >
                            <Stack spacing={1}>
                                <Stack
                                    direction={{ xs: 'column', xl: 'row' }}
                                    justifyContent="space-between"
                                    alignItems={{ xs: 'flex-start', xl: 'center' }}
                                    spacing={0.5}
                                >
                                    <Box>
                                        <Typography
                                            variant="body2"
                                            sx={{ fontWeight: 700, color: isSameDay(day.date, new Date()) ? 'primary.main' : 'text.primary' }}
                                        >
                                            {format(day.date, 'dd/MM', { locale: vi })}
                                        </Typography>
                                        <Typography variant="caption" color="text.secondary">
                                            {format(day.date, 'EEEE', { locale: vi })}
                                        </Typography>
                                    </Box>
                                    <Chip size="small" variant="outlined" label={`${day.lessons.length} buổi`} />
                                </Stack>

                                {day.lessons.length ? (
                                    <Stack spacing={1}>
                                        {day.lessons.map((lesson) => {
                                            const status = getLessonStatus(lesson);
                                            const linkedLeave = findMatchingLeave(lesson, leaveRequests);
                                            const attendanceRecord = attendanceRecords.find((record) => record.lesson.id === lesson.id);
                                            const attendanceStatus = getAttendanceStatus(attendanceRecord?.status);
                                            const tone = getLessonTone(lesson);

                                            return (
                                                <Paper
                                                    key={lesson.id}
                                                    variant="outlined"
                                                    sx={{
                                                        p: 1.25,
                                                        borderRadius: 2.5,
                                                        borderLeft: `4px solid ${tone.border}`,
                                                        backgroundColor: tone.bg,
                                                        borderColor: linkedLeave ? 'warning.light' : undefined,
                                                    }}
                                                >
                                                    <Stack spacing={0.75}>
                                                        <Typography variant="body2" sx={{ fontWeight: 700 }}>
                                                            {lesson.class_name}
                                                        </Typography>
                                                        <Typography variant="caption" color="text.secondary">
                                                            {lesson.class_code} • {format(parseISO(lesson.date_start), 'HH:mm')} - {format(parseISO(lesson.date_end), 'HH:mm')}
                                                        </Typography>
                                                        <Typography variant="caption" color="text.secondary">
                                                            Giáo viên: {lesson.teacher?.full_name || 'Chưa phân công'}
                                                        </Typography>
                                                        <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                                                            <Chip size="small" label={lesson.shift?.name || 'Chưa gắn ca học'} color="primary" variant="outlined" />
                                                            <Chip size="small" label={lesson.room_name || 'Chưa xếp phòng'} variant="outlined" />
                                                            <Chip size="small" label={status.label} color={status.color} variant="outlined" />
                                                            <Chip size="small" label={attendanceStatus.label} color={attendanceStatus.color} variant="outlined" />
                                                            {linkedLeave ? (
                                                                <Chip
                                                                    size="small"
                                                                    label={`Đơn phép: ${linkedLeave.status}`}
                                                                    color={linkedLeave.status === 'APPROVED' ? 'success' : linkedLeave.status === 'PENDING' ? 'warning' : 'default'}
                                                                />
                                                            ) : null}
                                                        </Stack>
                                                        {lesson.notes ? (
                                                            <Typography variant="caption" color="text.secondary" noWrap>
                                                                {lesson.notes}
                                                            </Typography>
                                                        ) : null}
                                                        {linkedLeave ? (
                                                            <Alert severity={linkedLeave.status === 'APPROVED' ? 'success' : linkedLeave.status === 'REJECTED' ? 'error' : 'warning'}>
                                                                {linkedLeave.subject || 'Đơn xin phép'}: {linkedLeave.reason}
                                                            </Alert>
                                                        ) : null}
                                                    </Stack>
                                                </Paper>
                                            );
                                        })}
                                    </Stack>
                                ) : (
                                    <Alert severity="info">
                                        {isLoadingTimetable || isFetchingTimetable ? 'Đang tải thời khóa biểu...' : 'Chưa có buổi học nào trong ngày này.'}
                                    </Alert>
                                )}
                            </Stack>
                        </Box>
                    ))}
                </Box>
            </Paper>

            <Paper variant="outlined" sx={{ borderRadius: 3 }}>
                <Stack spacing={2} sx={{ p: 2.5 }}>
                    <Typography variant="h6" sx={{ fontWeight: 700 }}>
                        Chi tiết điểm danh từng buổi
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        Theo dõi trạng thái điểm danh, ghi chú giáo viên và thời điểm được cập nhật gần nhất.
                    </Typography>
                </Stack>
                <TableContainer>
                    <Table>
                        <TableHead>
                            <TableRow>
                                <TableCell sx={{ fontWeight: 700 }}>Buổi học</TableCell>
                                <TableCell sx={{ fontWeight: 700 }}>Giáo viên</TableCell>
                                <TableCell sx={{ fontWeight: 700 }}>Ca học</TableCell>
                                <TableCell sx={{ fontWeight: 700 }}>Trạng thái</TableCell>
                                <TableCell sx={{ fontWeight: 700 }}>Ghi chú</TableCell>
                                <TableCell sx={{ fontWeight: 700 }}>Cập nhật</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {sortedAttendanceRecords.map((record) => {
                                const status = getAttendanceStatus(record.status);
                                return (
                                    <TableRow key={record.lesson.id} hover>
                                        <TableCell>
                                            <Typography variant="body2" sx={{ fontWeight: 600 }}>
                                                {record.lesson.class_name}
                                            </Typography>
                                            <Typography variant="caption" color="text.secondary">
                                                {format(parseISO(record.lesson.date_start), 'dd/MM/yyyy HH:mm')} - {format(parseISO(record.lesson.date_end), 'HH:mm')}
                                            </Typography>
                                        </TableCell>
                                        <TableCell>{record.lesson.teacher?.full_name || 'Chưa phân công'}</TableCell>
                                        <TableCell>{record.lesson.shift?.name || 'Chưa gắn ca học'}</TableCell>
                                        <TableCell>
                                            <Chip size="small" label={status.label} color={status.color} variant="outlined" />
                                        </TableCell>
                                        <TableCell>{record.note || 'Không có ghi chú'}</TableCell>
                                        <TableCell>
                                            {record.marked_at ? new Date(record.marked_at).toLocaleString('vi-VN') : 'Chưa cập nhật'}
                                        </TableCell>
                                    </TableRow>
                                );
                            })}
                            {!sortedAttendanceRecords.length && !isLoadingAttendance ? (
                                <TableRow>
                                    <TableCell colSpan={6}>
                                        <Typography variant="body2" color="text.secondary">
                                            Chưa có dữ liệu điểm danh nào trong phạm vi hiện tại.
                                        </Typography>
                                    </TableCell>
                                </TableRow>
                            ) : null}
                        </TableBody>
                    </Table>
                </TableContainer>
            </Paper>
        </Stack>
    );
}
