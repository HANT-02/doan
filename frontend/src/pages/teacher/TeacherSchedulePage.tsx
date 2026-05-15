import { useMemo, useState } from 'react';
import { Alert, Button, Chip, Stack } from '@mui/material';
import {
    CalendarMonthRounded,
    ChevronLeftRounded,
    ChevronRightRounded,
    EventAvailableRounded,
    TodayRounded,
} from '@mui/icons-material';
import { addWeeks, endOfWeek, format, parseISO, startOfWeek } from 'date-fns';
import { vi } from 'date-fns/locale';
import { useNavigate } from 'react-router-dom';

import { useGetTeacherLessonsQuery } from '@/api/teacherPortalApi';
import type { TeacherLesson } from '@/api/teacherPortalApi';
import PageHeader from '@/components/common/PageHeader';
import PreviewScheduleGrid from '@/components/schedule/PreviewScheduleGrid';
import ScheduleCardShell from '@/components/schedule/ScheduleCardShell';
import { getApiErrorMessage } from '@/utils/apiError';

export default function TeacherSchedulePage() {
    const navigate = useNavigate();
    const [weekAnchor, setWeekAnchor] = useState(() => startOfWeek(new Date(), { weekStartsOn: 1 }));

    const weekStart = useMemo(() => startOfWeek(weekAnchor, { weekStartsOn: 1 }), [weekAnchor]);
    const weekEnd = useMemo(() => endOfWeek(weekAnchor, { weekStartsOn: 1 }), [weekAnchor]);

    const { data, isFetching, error } = useGetTeacherLessonsQuery({
        from: format(weekStart, 'yyyy-MM-dd'),
        to: format(weekEnd, 'yyyy-MM-dd'),
    });

    const lessons = data?.data?.lessons ?? [];

    const uniqueClasses = useMemo(() => new Set(lessons.map((lesson) => lesson.class_id)).size, [lessons]);

    return (
        <Stack sx={{ p: { xs: 2, md: 4 } }} spacing={3}>
            <PageHeader
                title="Lịch giảng dạy"
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

            <PreviewScheduleGrid
                title="Lịch tuần giảng dạy"
                weekStart={weekStart}
                items={lessons}
                isFetching={isFetching}
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
                        actionLabel="Xem chi tiết"
                        onActionClick={() => navigate(`/app/teacher/lessons/${lesson.id}`)}
                        onClick={() => navigate(`/app/teacher/lessons/${lesson.id}`)}
                    />
                )}
            />
        </Stack>
    );
}
