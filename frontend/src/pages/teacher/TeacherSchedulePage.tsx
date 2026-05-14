import { useMemo, useState } from 'react';
import { Alert, Button, Chip, Stack } from '@mui/material';
import {
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
import ScheduleCardShell from '@/components/schedule/ScheduleCardShell';
import WeekScheduleBoard from '@/components/schedule/WeekScheduleBoard';
import { getApiErrorMessage } from '@/utils/apiError';

function buildWeekDays(weekStart: Date) {
    return Array.from({ length: 7 }, (_, index) => addDays(weekStart, index));
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
                items: lessons
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
                subtitle="Theo dõi các buổi học trong tuần theo cùng ngôn ngữ hiển thị với màn preview scheduling."
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

            <WeekScheduleBoard
                title="Lịch tuần giảng dạy"
                subtitle="Mỗi buổi học là một card với chip ca học và phòng học, giống cách đọc trên preview."
                days={calendarDays}
                isFetching={isFetching}
                emptyLabel={isLoading || isFetching ? 'Đang tải lịch giảng dạy...' : 'Chưa có buổi học nào trong ngày này.'}
                renderItem={(lesson: TeacherLesson) => (
                    <ScheduleCardShell
                        key={lesson.id}
                        seed={lesson.class_id}
                        title={lesson.class_name}
                        subtitle={`${lesson.class_code} • ${format(parseISO(lesson.date_start), 'HH:mm')} - ${format(parseISO(lesson.date_end), 'HH:mm')}`}
                        chips={[
                            { label: lesson.shift?.name || 'Chưa gắn ca học', color: 'primary' },
                            { label: lesson.room_name || 'Chưa xếp phòng' },
                        ]}
                        note={lesson.notes || undefined}
                        actionLabel="Xem chi tiết"
                        onActionClick={() => navigate(`/app/teacher/lessons/${lesson.id}`)}
                        onClick={() => navigate(`/app/teacher/lessons/${lesson.id}`)}
                    />
                )}
            />
        </Stack>
    );
}
