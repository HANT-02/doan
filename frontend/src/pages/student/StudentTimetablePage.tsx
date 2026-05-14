import { useMemo, useState } from 'react';
import {
    Alert,
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
import ScheduleCardShell from '@/components/schedule/ScheduleCardShell';
import WeekScheduleBoard from '@/components/schedule/WeekScheduleBoard';
import { getApiErrorMessage } from '@/utils/apiError';

function isRenderableEmptyStateError(error: unknown) {
    if (typeof error !== 'object' || !error || !('status' in error)) {
        return false;
    }

    const apiError = error as { status?: number | string };
    return apiError.status === 403 || apiError.status === 404;
}

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

    const { data: leaveResponse, error: leaveError } = useGetMyStudentLeaveRequestsQuery({
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
    const canRenderTimetableShell = !timetableError || isRenderableEmptyStateError(timetableError);
    const canRenderAttendanceShell = !attendanceError || isRenderableEmptyStateError(attendanceError);
    const canRenderLeaveShell = !leaveError || isRenderableEmptyStateError(leaveError);

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
                items: lessons
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
                subtitle="Theo dõi lịch học theo tuần ở dạng calendar, đồng bộ cách đọc card với màn preview scheduling."
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

            {!canRenderTimetableShell && timetableError ? (
                <Alert severity="error">{getApiErrorMessage(timetableError, 'Không tải được thời khóa biểu của bạn.')}</Alert>
            ) : null}
            {!canRenderLeaveShell && leaveError ? (
                <Alert severity="warning">{getApiErrorMessage(leaveError, 'Không tải được dữ liệu đơn xin phép liên quan.')}</Alert>
            ) : null}
            {!canRenderAttendanceShell && attendanceError ? (
                <Alert severity="error">{getApiErrorMessage(attendanceError, 'Không tải được dữ liệu điểm danh của bạn.')}</Alert>
            ) : null}
            {(isRenderableEmptyStateError(timetableError) || isRenderableEmptyStateError(attendanceError)) ? (
                <Alert severity="info">
                    Tài khoản này hiện chưa có dữ liệu lịch học khả dụng. Giao diện vẫn được mở để bạn theo dõi tuần học, và sẽ tự có dữ liệu ngay khi học viên được gán lớp hoặc phát sinh buổi học.
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
                        <Chip label={`Tổng buổi ${canRenderAttendanceShell ? attendanceSummary?.total_lessons ?? 0 : 0}`} variant="outlined" />
                        <Chip label={`Có mặt ${canRenderAttendanceShell ? attendanceSummary?.present_count ?? 0 : 0}`} color="success" variant="outlined" />
                        <Chip label={`Vắng ${canRenderAttendanceShell ? attendanceSummary?.absent_count ?? 0 : 0}`} color="error" variant="outlined" />
                        <Chip label={`Muộn ${canRenderAttendanceShell ? attendanceSummary?.late_count ?? 0 : 0}`} color="warning" variant="outlined" />
                        <Chip label={`Xin phép ${canRenderAttendanceShell ? attendanceSummary?.excused_count ?? 0 : 0}`} color="info" variant="outlined" />
                        <Chip label={`Chưa chấm ${canRenderAttendanceShell ? attendanceSummary?.unmarked_count ?? 0 : 0}`} variant="outlined" />
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

            <WeekScheduleBoard
                title="Lịch tuần"
                subtitle="Các buổi học hiển thị cùng phong cách card/chip như preview, kèm trạng thái điểm danh và đơn xin phép."
                days={calendarDays}
                isFetching={isFetchingTimetable}
                minDayHeight={280}
                emptyLabel={isLoadingTimetable || isFetchingTimetable ? 'Đang tải thời khóa biểu...' : 'Chưa có buổi học nào trong ngày này.'}
                renderItem={(lesson: StudentTimetableLesson) => {
                    const status = getLessonStatus(lesson);
                    const linkedLeave = findMatchingLeave(lesson, leaveRequests);
                    const attendanceRecord = attendanceRecords.find((record) => record.lesson.id === lesson.id);
                    const attendanceStatus = getAttendanceStatus(attendanceRecord?.status);

                    return (
                        <ScheduleCardShell
                            key={lesson.id}
                            seed={lesson.class_id}
                            title={lesson.class_name}
                            subtitle={`${lesson.class_code} • ${format(parseISO(lesson.date_start), 'HH:mm')} - ${format(parseISO(lesson.date_end), 'HH:mm')}`}
                            lines={[`Giáo viên: ${lesson.teacher?.full_name || 'Chưa phân công'}`]}
                            chips={[
                                { label: lesson.shift?.name || 'Chưa gắn ca học', color: 'primary' },
                                { label: lesson.room_name || 'Chưa xếp phòng' },
                                { label: status.label, color: status.color },
                                { label: attendanceStatus.label, color: attendanceStatus.color },
                                ...(linkedLeave ? [{
                                    label: `Đơn phép: ${linkedLeave.status}`,
                                    color: linkedLeave.status === 'APPROVED' ? 'success' as const : linkedLeave.status === 'PENDING' ? 'warning' as const : 'default' as const,
                                }] : []),
                            ]}
                            warning={!!linkedLeave}
                        >
                            {linkedLeave ? (
                                <Alert severity={linkedLeave.status === 'APPROVED' ? 'success' : linkedLeave.status === 'REJECTED' ? 'error' : 'warning'}>
                                    {linkedLeave.subject || linkedLeave.reason || 'Đã ghi nhận đơn xin phép cho buổi học này.'}
                                </Alert>
                            ) : null}
                        </ScheduleCardShell>
                    );
                }}
            />

            <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 3 }}>
                <Stack spacing={2}>
                    <Typography variant="h6" sx={{ fontWeight: 700 }}>
                        Lịch sử điểm danh gần đây
                    </Typography>
                    <TableContainer>
                        <Table size="small">
                            <TableHead>
                                <TableRow>
                                    <TableCell>Ngày học</TableCell>
                                    <TableCell>Lớp</TableCell>
                                    <TableCell>Ca học</TableCell>
                                    <TableCell>Trạng thái</TableCell>
                                    <TableCell>Ghi chú</TableCell>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {sortedAttendanceRecords.map((record) => {
                                    const attendanceStatus = getAttendanceStatus(record.status);
                                    return (
                                        <TableRow key={record.lesson.id}>
                                            <TableCell>{format(parseISO(record.lesson.date_start), 'dd/MM/yyyy HH:mm')}</TableCell>
                                            <TableCell>{record.lesson.class_name}</TableCell>
                                            <TableCell>{record.lesson.shift?.name || 'Chưa gắn ca học'}</TableCell>
                                            <TableCell>
                                                <Chip size="small" label={attendanceStatus.label} color={attendanceStatus.color} variant="outlined" />
                                            </TableCell>
                                            <TableCell>{record.note || '-'}</TableCell>
                                        </TableRow>
                                    );
                                })}
                                {!isLoadingAttendance && sortedAttendanceRecords.length === 0 ? (
                                    <TableRow>
                                        <TableCell colSpan={5}>
                                            <Typography variant="body2" color="text.secondary">
                                                Chưa có dữ liệu điểm danh cho bộ lọc hiện tại.
                                            </Typography>
                                        </TableCell>
                                    </TableRow>
                                ) : null}
                            </TableBody>
                        </Table>
                    </TableContainer>
                </Stack>
            </Paper>
        </Stack>
    );
}
