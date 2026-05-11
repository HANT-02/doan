import { useMemo, useState } from 'react';
import {
    Alert,
    Box,
    Button,
    Chip,
    Paper,
    Stack,
    Typography,
} from '@mui/material';
import {
    ArrowForwardRounded,
    CalendarMonthRounded,
    ChevronLeftRounded,
    ChevronRightRounded,
    EventAvailableRounded,
    TodayRounded,
} from '@mui/icons-material';
import { addDays, addWeeks, endOfWeek, format, isSameDay, parseISO, startOfWeek } from 'date-fns';
import { vi } from 'date-fns/locale';
import { useNavigate } from 'react-router-dom';

import { useGetTeacherLessonsQuery } from '@/api/teacherPortalApi';
import type { TeacherLesson } from '@/api/teacherPortalApi';
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

function getLessonTone(lesson: TeacherLesson) {
    const seed = lesson.class_id
        .split('')
        .reduce((sum, char) => sum + char.charCodeAt(0), 0);
    return palette[seed % palette.length];
}

export default function TeacherSchedulePage() {
    const navigate = useNavigate();
    const [weekAnchor, setWeekAnchor] = useState(() => startOfWeek(new Date(), { weekStartsOn: 1 }));

    const weekStart = useMemo(() => startOfWeek(weekAnchor, { weekStartsOn: 1 }), [weekAnchor]);
    const weekEnd = useMemo(() => endOfWeek(weekAnchor, { weekStartsOn: 1 }), [weekAnchor]);

    const { data, isLoading, isFetching, error } = useGetTeacherLessonsQuery({
        from: format(weekStart, 'yyyy-MM-dd'),
        to: format(weekEnd, 'yyyy-MM-dd'),
    });

    const lessons = data?.data?.lessons ?? [];

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

    return (
        <Stack sx={{ p: { xs: 2, md: 4 } }} spacing={3}>
            <PageHeader
                title="Lịch giảng dạy"
                subtitle="Theo dõi các buổi học trong tuần theo dạng calendar và chuyển thẳng sang màn chi tiết từng buổi."
                icon={<EventAvailableRounded />}
                breadcrumbs={[
                    { label: 'Trang giáo viên', path: '/app/teacher/overview' },
                    { label: 'Lịch giảng dạy' },
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
                <Chip label={`${uniqueClasses} lớp đang dạy`} color="primary" variant="outlined" />
                <Chip
                    icon={<CalendarMonthRounded />}
                    label={`${format(weekStart, 'dd/MM', { locale: vi })} - ${format(weekEnd, 'dd/MM/yyyy', { locale: vi })}`}
                    color="secondary"
                    variant="outlined"
                />
            </Stack>

            {error ? (
                <Alert severity="error">
                    {getApiErrorMessage(error, 'Không tải được lịch giảng dạy của giáo viên.')}
                </Alert>
            ) : null}

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
                                Lich tuần giảng dạy
                            </Typography>
                            <Typography variant="body2" color="text.secondary">
                                Mỗi card là một buổi học. Bấm vào card hoặc nút chi tiết để mở lesson tương ứng.
                            </Typography>
                        </Box>
                        {isFetching ? (
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
                                minHeight: 260,
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

                                {day.lessons.length > 0 ? (
                                    <Stack spacing={1}>
                                        {day.lessons.map((lesson) => {
                                            const tone = getLessonTone(lesson);
                                            return (
                                                <Paper
                                                    key={lesson.id}
                                                    variant="outlined"
                                                    sx={{
                                                        p: 1.25,
                                                        borderRadius: 2.5,
                                                        cursor: 'pointer',
                                                        borderLeft: `4px solid ${tone.border}`,
                                                        backgroundColor: tone.bg,
                                                        transition: 'transform 0.18s ease, box-shadow 0.18s ease',
                                                        '&:hover': {
                                                            transform: 'translateY(-1px)',
                                                            boxShadow: 2,
                                                        },
                                                    }}
                                                    onClick={() => navigate(`/app/teacher/lessons/${lesson.id}`)}
                                                >
                                                    <Stack spacing={0.75}>
                                                        <Typography variant="body2" sx={{ fontWeight: 700 }}>
                                                            {lesson.class_name}
                                                        </Typography>
                                                        <Typography variant="caption" color="text.secondary">
                                                            {lesson.class_code} • {format(parseISO(lesson.date_start), 'HH:mm')} - {format(parseISO(lesson.date_end), 'HH:mm')}
                                                        </Typography>
                                                        <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                                                            <Chip size="small" label={lesson.shift?.name || 'Chưa gắn ca học'} color="primary" variant="outlined" />
                                                            <Chip size="small" label={lesson.room_name || 'Chưa xếp phòng'} variant="outlined" />
                                                        </Stack>
                                                        {lesson.notes ? (
                                                            <Typography variant="caption" color="text.secondary" noWrap>
                                                                {lesson.notes}
                                                            </Typography>
                                                        ) : null}
                                                        <Button
                                                            variant="text"
                                                            size="small"
                                                            endIcon={<ArrowForwardRounded />}
                                                            onClick={(event) => {
                                                                event.stopPropagation();
                                                                navigate(`/app/teacher/lessons/${lesson.id}`);
                                                            }}
                                                            sx={{ alignSelf: 'flex-start', px: 0 }}
                                                        >
                                                            Xem chi tiết
                                                        </Button>
                                                    </Stack>
                                                </Paper>
                                            );
                                        })}
                                    </Stack>
                                ) : (
                                    <Alert severity="info">
                                        {isLoading || isFetching ? 'Đang tải lịch giảng dạy...' : 'Chưa có buổi học nào trong ngày này.'}
                                    </Alert>
                                )}
                            </Stack>
                        </Box>
                    ))}
                </Box>
            </Paper>
        </Stack>
    );
}
