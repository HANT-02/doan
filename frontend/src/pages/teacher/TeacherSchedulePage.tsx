import { useMemo, useState } from 'react';
import {
    Alert,
    Button,
    Chip,
    Paper,
    Stack,
    Typography,
} from '@mui/material';
import {
    ArrowForwardRounded,
    ChevronLeftRounded,
    ChevronRightRounded,
    EventAvailableRounded,
} from '@mui/icons-material';
import { addDays, addWeeks, endOfWeek, format, isSameDay, parseISO, startOfWeek } from 'date-fns';
import { vi } from 'date-fns/locale';
import { useNavigate } from 'react-router-dom';

import { useGetTeacherLessonsQuery } from '@/api/teacherPortalApi';
import type { TeacherLesson } from '@/api/teacherPortalApi';
import PageHeader from '@/components/common/PageHeader';
import { getApiErrorMessage } from '@/utils/apiError';

type DailyGroup = {
    label: string;
    date: Date;
    lessons: TeacherLesson[];
};

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

    const dailyGroups = useMemo<DailyGroup[]>(() => {
        const weekDays = buildWeekDays(weekStart);
        return weekDays.map((day) => ({
            label: format(day, "EEEE, dd/MM", { locale: vi }),
            date: day,
            lessons: lessons.filter((lesson) => isSameDay(parseISO(lesson.date_start), day)),
        }));
    }, [lessons, weekStart]);

    const uniqueClasses = useMemo(() => new Set(lessons.map((lesson) => lesson.class_id)).size, [lessons]);

    return (
        <Stack sx={{ p: { xs: 2, md: 4 } }} spacing={3}>
            <PageHeader
                title="Lịch giảng dạy"
                subtitle="Theo dõi các buổi học trong tuần, xem nhanh ca học và chuyển thẳng sang màn chi tiết từng buổi."
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
                        <Button variant="outlined" onClick={() => setWeekAnchor(startOfWeek(new Date(), { weekStartsOn: 1 }))}>
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

            <Stack spacing={2}>
                {dailyGroups.map((group) => (
                    <Paper key={group.label} variant="outlined" sx={{ p: 2.5, borderRadius: 3 }}>
                        <Stack spacing={2}>
                            <Stack
                                direction={{ xs: 'column', md: 'row' }}
                                justifyContent="space-between"
                                alignItems={{ xs: 'flex-start', md: 'center' }}
                                spacing={1}
                            >
                                <Typography variant="h6" sx={{ fontWeight: 700, textTransform: 'capitalize' }}>
                                    {group.label}
                                </Typography>
                                <Chip label={`${group.lessons.length} buổi`} size="small" variant="outlined" />
                            </Stack>

                            {group.lessons.length ? (
                                <Stack spacing={1.5}>
                                    {group.lessons.map((lesson) => (
                                        <Paper
                                            key={lesson.id}
                                            variant="outlined"
                                            sx={{
                                                p: 2,
                                                borderRadius: 2.5,
                                                cursor: 'pointer',
                                                transition: 'all 0.2s ease',
                                                '&:hover': {
                                                    borderColor: 'primary.main',
                                                    boxShadow: 2,
                                                },
                                            }}
                                            onClick={() => navigate(`/app/teacher/lessons/${lesson.id}`)}
                                        >
                                            <Stack
                                                direction={{ xs: 'column', md: 'row' }}
                                                justifyContent="space-between"
                                                spacing={2}
                                            >
                                                <Stack spacing={0.75}>
                                                    <Typography variant="subtitle1" sx={{ fontWeight: 700 }}>
                                                        {lesson.class_name}
                                                    </Typography>
                                                    <Typography variant="body2" color="text.secondary">
                                                        {lesson.class_code} • {format(parseISO(lesson.date_start), 'HH:mm')} - {format(parseISO(lesson.date_end), 'HH:mm')}
                                                    </Typography>
                                                    <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                                                        <Chip
                                                            size="small"
                                                            label={lesson.shift?.name || 'Chưa gắn ca học'}
                                                            color="primary"
                                                            variant="outlined"
                                                        />
                                                        <Chip
                                                            size="small"
                                                            label={lesson.room_name || 'Chưa xếp phòng'}
                                                            variant="outlined"
                                                        />
                                                    </Stack>
                                                    {lesson.notes ? (
                                                        <Typography variant="body2" color="text.secondary">
                                                            {lesson.notes}
                                                        </Typography>
                                                    ) : null}
                                                </Stack>
                                                <Button
                                                    variant="text"
                                                    endIcon={<ArrowForwardRounded />}
                                                    onClick={(event) => {
                                                        event.stopPropagation();
                                                        navigate(`/app/teacher/lessons/${lesson.id}`);
                                                    }}
                                                >
                                                    Xem chi tiết
                                                </Button>
                                            </Stack>
                                        </Paper>
                                    ))}
                                </Stack>
                            ) : (
                                <Alert severity="info">
                                    {isLoading || isFetching
                                        ? 'Đang tải lịch giảng dạy...'
                                        : 'Chưa có buổi học nào trong ngày này.'}
                                </Alert>
                            )}
                        </Stack>
                    </Paper>
                ))}
            </Stack>
        </Stack>
    );
}
