import { useEffect, useMemo, useState } from 'react';
import {
    Alert,
    Chip,
    MenuItem,
    Paper,
    Select,
    Stack,
    Tab,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Tabs,
    TextField,
    Typography,
    Button,
} from '@mui/material';
import { AssignmentTurnedInRounded, SaveRounded, TimelineRounded } from '@mui/icons-material';
import { skipToken } from '@reduxjs/toolkit/query';
import { isSameDay, parseISO } from 'date-fns';
import { useSearchParams } from 'react-router-dom';
import { toast } from 'sonner';

import {
    useGetTeacherAttendanceSummaryQuery,
    useGetTeacherLeaveRequestsQuery,
    useGetTeacherLessonAttendanceQuery,
    useGetTeacherLessonsQuery,
    useSubmitTeacherLessonAttendanceMutation,
} from '@/api/teacherPortalApi';
import PageHeader from '@/components/common/PageHeader';
import { getApiErrorMessage } from '@/utils/apiError';

const ATTENDANCE_OPTIONS = [
    { value: 1, label: 'Có mặt' },
    { value: 0, label: 'Vắng' },
    { value: 2, label: 'Muộn' },
    { value: 3, label: 'Xin phép' },
];

type EditableAttendanceRow = {
    student_id: string;
    student_code: string;
    student_name: string;
    status: number;
    note: string;
};

type AttendanceTab = 'marking' | 'summary';

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

export default function TeacherAttendancePage() {
    const [searchParams] = useSearchParams();
    const preselectedLessonId = searchParams.get('lessonId') || '';

    const [activeTab, setActiveTab] = useState<AttendanceTab>('marking');
    const { data: lessonsResponse, isLoading: isLoadingLessons } = useGetTeacherLessonsQuery();
    const lessons = lessonsResponse?.data?.lessons ?? [];

    const classes = useMemo(() => {
        const map = new Map<string, { id: string; code: string; name: string }>();
        lessons.forEach((lesson) => {
            if (!map.has(lesson.class_id)) {
                map.set(lesson.class_id, {
                    id: lesson.class_id,
                    code: lesson.class_code,
                    name: lesson.class_name,
                });
            }
        });
        return Array.from(map.values());
    }, [lessons]);

    const [selectedClassId, setSelectedClassId] = useState('');
    const [selectedLessonId, setSelectedLessonId] = useState('');
    const [rows, setRows] = useState<EditableAttendanceRow[]>([]);

    useEffect(() => {
        if (!lessons.length) return;

        const preselectedLesson = lessons.find((lesson) => lesson.id === preselectedLessonId);
        if (preselectedLesson) {
            setSelectedClassId(preselectedLesson.class_id);
            setSelectedLessonId(preselectedLesson.id);
            return;
        }

        if (!selectedClassId) {
            setSelectedClassId(lessons[0].class_id);
        }
    }, [lessons, preselectedLessonId, selectedClassId]);

    const filteredLessons = useMemo(
        () => lessons.filter((lesson) => !selectedClassId || lesson.class_id === selectedClassId),
        [lessons, selectedClassId],
    );

    useEffect(() => {
        if (!filteredLessons.length) {
            setSelectedLessonId('');
            return;
        }

        if (selectedLessonId && filteredLessons.some((lesson) => lesson.id === selectedLessonId)) {
            return;
        }

        setSelectedLessonId(filteredLessons[0].id);
    }, [filteredLessons, selectedLessonId]);

    const selectedLesson = useMemo(
        () => lessons.find((lesson) => lesson.id === selectedLessonId) || null,
        [lessons, selectedLessonId],
    );

    const attendanceQuery = useGetTeacherLessonAttendanceQuery(selectedLessonId || skipToken);
    const summaryQuery = useGetTeacherAttendanceSummaryQuery(selectedClassId || skipToken);
    const leaveRequestsQuery = useGetTeacherLeaveRequestsQuery(
        selectedClassId
            ? {
                class_id: selectedClassId,
                status: 'PENDING',
            }
            : skipToken,
    );
    const [submitAttendance, { isLoading: isSaving }] = useSubmitTeacherLessonAttendanceMutation();

    useEffect(() => {
        const nextRows = (attendanceQuery.data?.data?.records ?? []).map((record) => ({
            student_id: record.student.id,
            student_code: record.student.code,
            student_name: record.student.full_name,
            status: record.status,
            note: record.note || '',
        }));
        setRows(nextRows);
    }, [attendanceQuery.data]);

    const attendanceStats = useMemo(() => {
        const counters = {
            total: rows.length,
            present: 0,
            absent: 0,
            late: 0,
            excused: 0,
            unmarked: 0,
        };

        rows.forEach((row) => {
            switch (row.status) {
                case 1:
                    counters.present += 1;
                    break;
                case 0:
                    counters.absent += 1;
                    break;
                case 2:
                    counters.late += 1;
                    break;
                case 3:
                    counters.excused += 1;
                    break;
                default:
                    counters.unmarked += 1;
            }
        });

        return counters;
    }, [rows]);

    const relatedLeaveRequests = useMemo(() => {
        const requests = leaveRequestsQuery.data?.data?.requests ?? [];
        const map = new Map<string, typeof requests[number]>();

        requests.forEach((request) => {
            if (!selectedLesson) return;
            const sameStudent = request.student.id;
            const sameLesson = request.lesson?.id === selectedLesson.id;
            const sameDay = isSameDay(parseISO(request.apply_date), parseISO(selectedLesson.date_start));
            if (!sameLesson && !sameDay) {
                return;
            }
            const existing = map.get(sameStudent);
            if (!existing || (existing.status !== 'PENDING' && request.status === 'PENDING')) {
                map.set(sameStudent, request);
            }
        });

        return map;
    }, [leaveRequestsQuery.data, selectedLesson]);

    const handleSave = async () => {
        if (!selectedLessonId) {
            toast.error('Hãy chọn buổi học trước khi lưu điểm danh.');
            return;
        }
        if (!rows.length) {
            toast.error('Buổi học này chưa có học sinh để điểm danh.');
            return;
        }

        try {
            await submitAttendance({
                lessonId: selectedLessonId,
                records: rows.map((row) => ({
                    student_id: row.student_id,
                    status: row.status,
                    note: row.note?.trim() || '',
                })),
            }).unwrap();
            toast.success('Đã lưu điểm danh buổi học.');
            attendanceQuery.refetch();
            summaryQuery.refetch();
        } catch (error) {
            toast.error(getApiErrorMessage(error, 'Không lưu được điểm danh của giáo viên.'));
        }
    };

    return (
        <Stack sx={{ p: { xs: 2, md: 4 } }} spacing={3}>
            <PageHeader
                title="Điểm danh"
                subtitle="Chọn lớp, chọn buổi học và cập nhật chuyên cần cho học sinh. Có thể chuyển sang tab tổng hợp để theo dõi tình hình chuyên cần của cả lớp."
            />

            <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 3 }}>
                <Stack spacing={2}>
                    <Stack direction={{ xs: 'column', md: 'row' }} spacing={2}>
                        <TextField
                            select
                            label="Lớp"
                            value={selectedClassId}
                            onChange={(event) => setSelectedClassId(event.target.value)}
                            fullWidth
                            disabled={isLoadingLessons || !classes.length}
                        >
                            {classes.map((item) => (
                                <MenuItem key={item.id} value={item.id}>
                                    {item.name} ({item.code})
                                </MenuItem>
                            ))}
                        </TextField>
                        <TextField
                            select
                            label="Buổi học"
                            value={selectedLessonId}
                            onChange={(event) => setSelectedLessonId(event.target.value)}
                            fullWidth
                            disabled={isLoadingLessons || !filteredLessons.length}
                        >
                            {filteredLessons.map((lesson) => (
                                <MenuItem key={lesson.id} value={lesson.id}>
                                    {lesson.class_name} - {new Date(lesson.date_start).toLocaleString('vi-VN')}
                                </MenuItem>
                            ))}
                        </TextField>
                    </Stack>

                    <Tabs value={activeTab} onChange={(_, nextValue) => setActiveTab(nextValue)} sx={{ minHeight: 0 }}>
                        <Tab value="marking" label="Điểm danh buổi học" />
                        <Tab value="summary" label="Theo dõi chuyên cần" />
                    </Tabs>
                </Stack>
            </Paper>

            {activeTab === 'marking' ? (
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
                                        Điểm danh theo buổi học
                                    </Typography>
                                </Stack>
                                <Typography variant="body2" color="text.secondary">
                                    Giáo viên có thể chọn trạng thái Có mặt / Vắng / Muộn / Xin phép cho từng học sinh, rồi lưu một lần cho cả buổi.
                                </Typography>
                            </Stack>

                            <Button
                                variant="contained"
                                startIcon={<SaveRounded />}
                                onClick={handleSave}
                                disabled={isSaving || attendanceQuery.isLoading || !rows.length}
                            >
                                {isSaving ? 'Đang lưu...' : 'Lưu điểm danh'}
                            </Button>
                        </Stack>

                        <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                            <Chip size="small" label={`${attendanceStats.total} học sinh`} variant="outlined" />
                            <Chip size="small" label={`Có mặt ${attendanceStats.present}`} color="success" variant="outlined" />
                            <Chip size="small" label={`Vắng ${attendanceStats.absent}`} color="error" variant="outlined" />
                            <Chip size="small" label={`Muộn ${attendanceStats.late}`} color="warning" variant="outlined" />
                            <Chip size="small" label={`Xin phép ${attendanceStats.excused}`} color="info" variant="outlined" />
                            {attendanceStats.unmarked > 0 ? (
                                <Chip size="small" label={`Chưa chấm ${attendanceStats.unmarked}`} variant="outlined" />
                            ) : null}
                        </Stack>

                        {attendanceQuery.error ? (
                            <Alert severity="error">
                                {getApiErrorMessage(attendanceQuery.error, 'Không tải được dữ liệu điểm danh của buổi học.')}
                            </Alert>
                        ) : null}

                        <TableContainer>
                            <Table size="small">
                                <TableHead>
                                    <TableRow>
                                        <TableCell sx={{ fontWeight: 700 }}>Học sinh</TableCell>
                                        <TableCell sx={{ fontWeight: 700, width: 220 }}>Trạng thái</TableCell>
                                        <TableCell sx={{ fontWeight: 700, width: 220 }}>Đơn xin phép liên quan</TableCell>
                                        <TableCell sx={{ fontWeight: 700 }}>Ghi chú</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {rows.map((row, index) => {
                                        const relatedLeave = relatedLeaveRequests.get(row.student_id);

                                        return (
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
                                                        onChange={(event) => {
                                                            const nextStatus = Number(event.target.value);
                                                            setRows((current) =>
                                                                current.map((item, currentIndex) =>
                                                                    currentIndex === index ? { ...item, status: nextStatus } : item,
                                                                ),
                                                            );
                                                        }}
                                                    >
                                                        {ATTENDANCE_OPTIONS.map((option) => (
                                                            <MenuItem key={option.value} value={String(option.value)}>
                                                                {option.label}
                                                            </MenuItem>
                                                        ))}
                                                    </Select>
                                                </TableCell>
                                                <TableCell>
                                                    {relatedLeave ? (
                                                        <Stack spacing={0.5}>
                                                            <Chip
                                                                size="small"
                                                                label={`${relatedLeave.status} • ${leaveTypeLabel(relatedLeave.leave_type)}`}
                                                                color={relatedLeave.status === 'PENDING' ? 'warning' : 'info'}
                                                                variant="outlined"
                                                            />
                                                            <Typography variant="caption" color="text.secondary">
                                                                {relatedLeave.reason}
                                                            </Typography>
                                                        </Stack>
                                                    ) : (
                                                        <Typography variant="caption" color="text.secondary">
                                                            Không có đơn liên quan
                                                        </Typography>
                                                    )}
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
                                            </TableRow>
                                        );
                                    })}
                                    {!rows.length && !attendanceQuery.isLoading ? (
                                        <TableRow>
                                            <TableCell colSpan={4}>
                                                <Typography variant="body2" color="text.secondary">
                                                    Buổi học này chưa có dữ liệu roster để điểm danh.
                                                </Typography>
                                            </TableCell>
                                        </TableRow>
                                    ) : null}
                                </TableBody>
                            </Table>
                        </TableContainer>
                    </Stack>
                </Paper>
            ) : (
                <Paper variant="outlined" sx={{ p: 3, borderRadius: 3 }}>
                    <Stack spacing={2.5}>
                        <Stack direction="row" spacing={1} alignItems="center">
                            <TimelineRounded color="primary" />
                            <Typography variant="h6" sx={{ fontWeight: 700 }}>
                                Theo dõi chuyên cần theo lớp
                            </Typography>
                        </Stack>

                        {summaryQuery.error ? (
                            <Alert severity="error">
                                {getApiErrorMessage(summaryQuery.error, 'Không tải được thống kê chuyên cần theo lớp.')}
                            </Alert>
                        ) : null}

                        <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                            <Chip label={`Tổng số buổi ${summaryQuery.data?.data?.total_lessons ?? 0}`} variant="outlined" />
                            <Chip label={`Số học sinh ${summaryQuery.data?.data?.students?.length ?? 0}`} color="primary" variant="outlined" />
                        </Stack>

                        <TableContainer>
                            <Table size="small">
                                <TableHead>
                                    <TableRow>
                                        <TableCell sx={{ fontWeight: 700 }}>Học sinh</TableCell>
                                        <TableCell sx={{ fontWeight: 700 }}>Có mặt</TableCell>
                                        <TableCell sx={{ fontWeight: 700 }}>Vắng</TableCell>
                                        <TableCell sx={{ fontWeight: 700 }}>Muộn</TableCell>
                                        <TableCell sx={{ fontWeight: 700 }}>Xin phép</TableCell>
                                        <TableCell sx={{ fontWeight: 700 }}>Chưa chấm</TableCell>
                                        <TableCell sx={{ fontWeight: 700 }}>Tỷ lệ chuyên cần</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {(summaryQuery.data?.data?.students ?? []).map((item) => (
                                        <TableRow key={item.student.id} hover>
                                            <TableCell>
                                                <Typography variant="body2" sx={{ fontWeight: 600 }}>
                                                    {item.student.full_name}
                                                </Typography>
                                                <Typography variant="caption" color="text.secondary">
                                                    {item.student.code}
                                                </Typography>
                                            </TableCell>
                                            <TableCell>{item.present_count}</TableCell>
                                            <TableCell>{item.absent_count}</TableCell>
                                            <TableCell>{item.late_count}</TableCell>
                                            <TableCell>{item.excused_count}</TableCell>
                                            <TableCell>{item.unmarked_count}</TableCell>
                                            <TableCell>{Math.round(item.attendance_rate * 100)}%</TableCell>
                                        </TableRow>
                                    ))}
                                    {!summaryQuery.data?.data?.students?.length && !summaryQuery.isLoading ? (
                                        <TableRow>
                                            <TableCell colSpan={7}>
                                                <Typography variant="body2" color="text.secondary">
                                                    Chưa có thống kê chuyên cần cho lớp đang chọn.
                                                </Typography>
                                            </TableCell>
                                        </TableRow>
                                    ) : null}
                                </TableBody>
                            </Table>
                        </TableContainer>
                    </Stack>
                </Paper>
            )}
        </Stack>
    );
}
