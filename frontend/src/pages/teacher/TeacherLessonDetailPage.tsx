import { Alert, Button, Chip, Paper, Stack, Tab, Tabs, Typography } from '@mui/material';
import {
    ArrowBackRounded,
    AssignmentTurnedInRounded,
    BookRounded,
    MenuBookRounded,
    PersonSearchRounded,
} from '@mui/icons-material';
import { parseISO, format } from 'date-fns';
import { vi } from 'date-fns/locale';
import { skipToken } from '@reduxjs/toolkit/query';
import { useMemo, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';

import {
    useGetTeacherLessonAttendanceQuery,
    useGetTeacherLessonSummaryQuery,
} from '@/api/teacherPortalApi';
import PageHeader from '@/components/common/PageHeader';
import TeacherLessonAcademicRecordManager from '@/components/lesson/TeacherLessonAcademicRecordManager';
import { getApiErrorMessage } from '@/utils/apiError';

export default function TeacherLessonDetailPage() {
    const navigate = useNavigate();
    const { lessonId } = useParams<{ lessonId: string }>();
    const [activeTab, setActiveTab] = useState<'overview' | 'records'>('overview');

    const attendanceQuery = useGetTeacherLessonAttendanceQuery(lessonId ?? skipToken);
    const summaryQuery = useGetTeacherLessonSummaryQuery(lessonId ?? skipToken);

    const lesson = attendanceQuery.data?.data?.lesson || summaryQuery.data?.data?.lesson;
    const summary = summaryQuery.data?.data?.summary;
    const attendanceRecords = attendanceQuery.data?.data?.records ?? [];

    const attendanceStats = useMemo(() => {
        return attendanceRecords.reduce(
            (acc, record) => {
                if (record.status === 1) acc.present += 1;
                if (record.status === 2) acc.late += 1;
                if (record.status === 0) acc.absent += 1;
                if (record.status === 3) acc.excused += 1;
                return acc;
            },
            { present: 0, late: 0, absent: 0, excused: 0 },
        );
    }, [attendanceRecords]);

    return (
        <Stack sx={{ p: { xs: 2, md: 4 } }} spacing={3}>
            <PageHeader
                title="Chi tiết buổi học"
                subtitle="Từ đây có thể đi tiếp sang điểm danh và sổ đầu bài của đúng buổi học đang chọn."
                icon={<MenuBookRounded />}
                breadcrumbs={[
                    { label: 'Trang giáo viên', path: '/app/teacher/overview' },
                    { label: 'Lịch giảng dạy', path: '/app/teacher/schedule' },
                    { label: 'Chi tiết buổi học' },
                ]}
                actions={(
                    <Button variant="outlined" startIcon={<ArrowBackRounded />} onClick={() => navigate('/app/teacher/schedule')}>
                        Quay lại lịch
                    </Button>
                )}
            />

            {attendanceQuery.error || summaryQuery.error ? (
                <Alert severity="error">
                    {getApiErrorMessage(attendanceQuery.error || summaryQuery.error, 'Không tải được chi tiết buổi học.')}
                </Alert>
            ) : null}

            {!lesson && !attendanceQuery.isLoading && !summaryQuery.isLoading ? (
                <Alert severity="warning">Không tìm thấy dữ liệu buổi học hoặc bạn không có quyền truy cập buổi học này.</Alert>
            ) : null}

            {lesson ? (
                <>
                    <Paper variant="outlined" sx={{ p: 3, borderRadius: 3 }}>
                        <Tabs
                            value={activeTab}
                            onChange={(_event, value) => setActiveTab(value)}
                            sx={{ mb: 2 }}
                        >
                            <Tab value="overview" label="Tổng quan buổi học" />
                            <Tab value="records" label="Kết quả học tập" />
                        </Tabs>

                        {activeTab === 'overview' ? (
                            <Stack spacing={2}>
                                <Stack direction={{ xs: 'column', md: 'row' }} justifyContent="space-between" spacing={2}>
                                    <Stack spacing={0.75}>
                                        <Typography variant="h5" sx={{ fontWeight: 700 }}>
                                            {lesson.class_name}
                                        </Typography>
                                        <Typography variant="body1" color="text.secondary">
                                            {lesson.class_code} • {format(parseISO(lesson.date_start), "EEEE, dd/MM/yyyy HH:mm", { locale: vi })} - {format(parseISO(lesson.date_end), 'HH:mm')}
                                        </Typography>
                                    </Stack>
                                    <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                                        <Chip label={lesson.shift?.name || 'Chưa gắn ca học'} color="primary" variant="outlined" />
                                        <Chip label={lesson.room_name || 'Chưa xếp phòng'} variant="outlined" />
                                    </Stack>
                                </Stack>

                                {lesson.notes ? (
                                    <Typography variant="body2" color="text.secondary">
                                        {lesson.notes}
                                    </Typography>
                                ) : null}

                                <Stack direction={{ xs: 'column', sm: 'row' }} spacing={1.5}>
                                    <Button
                                        variant="contained"
                                        startIcon={<AssignmentTurnedInRounded />}
                                        onClick={() => navigate(`/app/teacher/attendance?lessonId=${lesson.id}`)}
                                    >
                                        Mở điểm danh
                                    </Button>
                                    <Button
                                        variant="outlined"
                                        startIcon={<BookRounded />}
                                        onClick={() => navigate(`/app/teacher/journal?lessonId=${lesson.id}`)}
                                    >
                                        Mở sổ đầu bài
                                    </Button>
                                    <Button
                                        variant="outlined"
                                        startIcon={<PersonSearchRounded />}
                                        onClick={() => navigate(`/app/teacher/substitute?lessonId=${lesson.id}`)}
                                    >
                                        Tìm dạy thay
                                    </Button>
                                </Stack>
                            </Stack>
                        ) : (
                            <TeacherLessonAcademicRecordManager
                                lessonId={lesson.id}
                                title="Kết quả học tập của buổi học"
                                subtitle="Nhập điểm bài tập, thái độ, mức độ tham gia và chốt kết quả ngay trong trang chi tiết buổi học."
                            />
                        )}
                    </Paper>

                    {activeTab === 'overview' ? (
                        <>
                            <Paper variant="outlined" sx={{ p: 3, borderRadius: 3 }}>
                                <Stack spacing={1.5}>
                                    <Typography variant="h6" sx={{ fontWeight: 700 }}>
                                        Tình trạng điểm danh
                                    </Typography>
                                    <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                                        <Chip label={`Học sinh ${attendanceRecords.length}`} variant="outlined" />
                                        <Chip label={`Có mặt ${attendanceStats.present}`} color="success" variant="outlined" />
                                        <Chip label={`Muộn ${attendanceStats.late}`} color="warning" variant="outlined" />
                                        <Chip label={`Vắng ${attendanceStats.absent}`} color="error" variant="outlined" />
                                        <Chip label={`Xin phép ${attendanceStats.excused}`} color="info" variant="outlined" />
                                    </Stack>
                                </Stack>
                            </Paper>

                            <Paper variant="outlined" sx={{ p: 3, borderRadius: 3 }}>
                                <Stack spacing={1.5}>
                                    <Typography variant="h6" sx={{ fontWeight: 700 }}>
                                        Tóm tắt sổ đầu bài
                                    </Typography>
                                    {summary ? (
                                        <Stack spacing={1}>
                                            <Typography variant="body2">
                                                <strong>Chủ đề:</strong> {summary.topic || 'Chưa cập nhật'}
                                            </Typography>
                                            <Typography variant="body2">
                                                <strong>Nội dung buổi học:</strong> {summary.lesson_content || 'Chưa cập nhật'}
                                            </Typography>
                                            <Typography variant="body2">
                                                <strong>Bài tập:</strong> {summary.homework || 'Chưa cập nhật'}
                                            </Typography>
                                            <Typography variant="body2">
                                                <strong>Phản hồi lớp:</strong> {summary.class_feedback || 'Chưa cập nhật'}
                                            </Typography>
                                            <Typography variant="body2">
                                                <strong>Ghi chú giáo viên:</strong> {summary.teacher_notes || 'Chưa cập nhật'}
                                            </Typography>
                                        </Stack>
                                    ) : (
                                        <Alert severity="info">Buổi học này chưa có tổng kết. Bạn có thể mở sổ đầu bài để cập nhật.</Alert>
                                    )}
                                </Stack>
                            </Paper>
                        </>
                    ) : null}
                </>
            ) : null}
        </Stack>
    );
}
