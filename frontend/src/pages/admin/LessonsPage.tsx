import { useMemo, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import {
    Alert,
    Box,
    Button,
    ButtonGroup,
    Chip,
    MenuItem,
    Paper,
    Skeleton,
    Stack,
    TextField,
    Typography,
} from '@mui/material';
import { DatePicker, LocalizationProvider } from '@mui/x-date-pickers';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import {
    CalendarMonthRounded,
    ChevronLeftRounded,
    ChevronRightRounded,
    RefreshRounded,
    TodayRounded,
    VisibilityOutlined,
    ViewWeekRounded,
} from '@mui/icons-material';
import {
    addDays,
    addMonths,
    addWeeks,
    endOfMonth,
    endOfWeek,
    format,
    isSameDay,
    isSameMonth,
    parseISO,
    startOfMonth,
    startOfWeek,
} from 'date-fns';
import { vi } from 'date-fns/locale';

import { useGetLessonsQuery, type Lesson } from '@/api/lessonApi';
import { useGetClassesQuery } from '@/api/classApi';
import { useGetTeachersQuery } from '@/api/teacherApi';
import PageHeader from '@/components/common/PageHeader';
import { getApiErrorMessage } from '@/utils/apiError';

type CalendarViewMode = 'week' | 'month';

type CalendarDay = {
    date: Date;
    lessons: Lesson[];
};

const weekLabels = ['Thứ Hai', 'Thứ Ba', 'Thứ Tư', 'Thứ Năm', 'Thứ Sáu', 'Thứ Bảy', 'Chủ Nhật'];

const formatDateParam = (value: Date) => format(value, 'yyyy-MM-dd');

const formatDateTime = (value: string) => format(parseISO(value), 'dd/MM/yyyy HH:mm', { locale: vi });

const formatTimeOnly = (value: string) => format(parseISO(value), 'HH:mm', { locale: vi });

const palette = [
    { bg: 'rgba(14, 116, 144, 0.10)', border: '#0ea5e9' },
    { bg: 'rgba(5, 150, 105, 0.10)', border: '#10b981' },
    { bg: 'rgba(234, 88, 12, 0.10)', border: '#f97316' },
    { bg: 'rgba(190, 24, 93, 0.10)', border: '#ec4899' },
    { bg: 'rgba(79, 70, 229, 0.10)', border: '#6366f1' },
    { bg: 'rgba(161, 98, 7, 0.10)', border: '#eab308' },
];

function buildVisibleRange(anchorDate: Date, viewMode: CalendarViewMode) {
    if (viewMode === 'week') {
        return {
            start: startOfWeek(anchorDate, { weekStartsOn: 1 }),
            end: endOfWeek(anchorDate, { weekStartsOn: 1 }),
        };
    }

    const monthStart = startOfMonth(anchorDate);
    const monthEnd = endOfMonth(anchorDate);
    return {
        start: startOfWeek(monthStart, { weekStartsOn: 1 }),
        end: endOfWeek(monthEnd, { weekStartsOn: 1 }),
    };
}

function buildCalendarDays(anchorDate: Date, viewMode: CalendarViewMode, lessons: Lesson[]) {
    const { start, end } = buildVisibleRange(anchorDate, viewMode);
    const days: CalendarDay[] = [];

    for (let current = start; current <= end; current = addDays(current, 1)) {
        days.push({
            date: current,
            lessons: lessons
                .filter((lesson) => isSameDay(parseISO(lesson.date_start), current))
                .sort((left, right) => left.date_start.localeCompare(right.date_start)),
        });
    }

    return days;
}

function getLessonTone(lesson: Lesson) {
    const seed = lesson.class_id
        .split('')
        .reduce((sum, char) => sum + char.charCodeAt(0), 0);
    return palette[seed % palette.length];
}

export default function LessonsPage() {
    const navigate = useNavigate();
    const [searchParams] = useSearchParams();
    const [viewMode, setViewMode] = useState<CalendarViewMode>('month');
    const [anchorDate, setAnchorDate] = useState(new Date());
    const [classId, setClassId] = useState(searchParams.get('class_id') || '');
    const [teacherId, setTeacherId] = useState('');

    const visibleRange = useMemo(
        () => buildVisibleRange(anchorDate, viewMode),
        [anchorDate, viewMode],
    );

    const queryParams = useMemo(
        () => ({
            page: 1,
            limit: 500,
            class_id: classId || undefined,
            teacher_id: teacherId || undefined,
            date_from: formatDateParam(visibleRange.start),
            date_to: formatDateParam(visibleRange.end),
            sortBy: 'date_start',
            sortOrder: 'asc',
        }),
        [classId, teacherId, visibleRange.end, visibleRange.start],
    );

    const { data: lessonsResponse, isLoading, isFetching, isError, error, refetch } = useGetLessonsQuery(queryParams);
    const { data: classesResponse } = useGetClassesQuery({ limit: 200 });
    const { data: teachersResponse } = useGetTeachersQuery({ limit: 200 });

    const lessons = lessonsResponse?.data?.lessons || [];
    const totalItems = lessonsResponse?.data?.pagination?.total_items || 0;
    const classes = classesResponse?.data?.classes || [];
    const teachers = teachersResponse?.data?.teachers || [];
    const hasFilters = !!(classId || teacherId);

    const calendarDays = useMemo(
        () => buildCalendarDays(anchorDate, viewMode, lessons),
        [anchorDate, lessons, viewMode],
    );

    const groupedWeeks = useMemo(() => {
        const weeks: CalendarDay[][] = [];
        for (let index = 0; index < calendarDays.length; index += 7) {
            weeks.push(calendarDays.slice(index, index + 7));
        }
        return weeks;
    }, [calendarDays]);

    const lessonsWithoutRoom = useMemo(
        () => lessons.filter((lesson) => !lesson.room?.name).length,
        [lessons],
    );

    const lessonsWithoutTeacher = useMemo(
        () => lessons.filter((lesson) => !lesson.teacher?.full_name).length,
        [lessons],
    );

    const headerLabel =
        viewMode === 'week'
            ? `${format(visibleRange.start, 'dd/MM', { locale: vi })} - ${format(visibleRange.end, 'dd/MM/yyyy', { locale: vi })}`
            : format(anchorDate, 'MMMM yyyy', { locale: vi });

    const shiftCalendar = (direction: 'prev' | 'next') => {
        setAnchorDate((current) =>
            viewMode === 'week'
                ? direction === 'prev'
                    ? addWeeks(current, -1)
                    : addWeeks(current, 1)
                : direction === 'prev'
                    ? addMonths(current, -1)
                    : addMonths(current, 1),
        );
    };

    const renderLessonCard = (lesson: Lesson) => {
        const tone = getLessonTone(lesson);
        return (
            <Paper
                key={lesson.id}
                variant="outlined"
                onClick={() => navigate(`/app/admin/lessons/${lesson.id}`)}
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
            >
                <Stack spacing={0.75}>
                    <Stack direction="row" justifyContent="space-between" alignItems="flex-start" spacing={1}>
                        <Box sx={{ minWidth: 0 }}>
                            <Typography variant="body2" sx={{ fontWeight: 700 }} noWrap>
                                {lesson.class?.name || lesson.class_id}
                            </Typography>
                            <Typography variant="caption" color="text.secondary">
                                {lesson.class?.code || 'Không có mã lớp'}
                            </Typography>
                        </Box>
                        <Button
                            size="small"
                            startIcon={<VisibilityOutlined />}
                            onClick={(event) => {
                                event.stopPropagation();
                                navigate(`/app/admin/lessons/${lesson.id}`);
                            }}
                            sx={{ minWidth: 'auto', px: 1 }}
                        >
                            Xem
                        </Button>
                    </Stack>

                    <Typography variant="caption" sx={{ fontWeight: 700 }}>
                        {formatTimeOnly(lesson.date_start)} - {formatTimeOnly(lesson.date_end)}
                    </Typography>

                    <Typography variant="caption" color="text.secondary" noWrap>
                        {lesson.teacher?.full_name || 'Chưa phân công giáo viên'}
                    </Typography>

                    <Typography variant="caption" color="text.secondary" noWrap>
                        {lesson.room?.name || 'Chưa xếp phòng'}
                    </Typography>

                    {lesson.notes ? (
                        <Typography variant="caption" color="text.secondary" sx={{ display: 'block' }} noWrap>
                            {lesson.notes}
                        </Typography>
                    ) : null}
                </Stack>
            </Paper>
        );
    };

    return (
        <LocalizationProvider dateAdapter={AdapterDateFns} adapterLocale={vi}>
            <Stack sx={{ p: { xs: 2, md: 4 } }} spacing={3}>
                <PageHeader
                    title="Quản lý buổi học"
                    subtitle="Theo dõi lesson theo dạng lịch, lọc theo lớp hoặc giáo viên và đi thẳng vào chi tiết từng buổi."
                    icon={<CalendarMonthRounded />}
                    actions={(
                        <Button
                            startIcon={<RefreshRounded />}
                            variant="outlined"
                            onClick={() => void refetch()}
                            disabled={isFetching}
                        >
                            Làm mới
                        </Button>
                    )}
                />

                <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 3 }}>
                    <Stack spacing={2.5}>
                        <Stack
                            direction={{ xs: 'column', xl: 'row' }}
                            spacing={2}
                            justifyContent="space-between"
                            alignItems={{ xs: 'stretch', xl: 'center' }}
                        >
                            <Stack direction={{ xs: 'column', md: 'row' }} spacing={2} sx={{ flex: 1 }}>
                                <TextField
                                    select
                                    label="Lọc theo lớp"
                                    value={classId}
                                    onChange={(event) => setClassId(event.target.value)}
                                    size="small"
                                    fullWidth
                                >
                                    <MenuItem value="">Tất cả lớp</MenuItem>
                                    {classes.map((lessonClass) => (
                                        <MenuItem key={lessonClass.id} value={lessonClass.id}>
                                            {lessonClass.name} ({lessonClass.code})
                                        </MenuItem>
                                    ))}
                                </TextField>

                                <TextField
                                    select
                                    label="Lọc theo giáo viên"
                                    value={teacherId}
                                    onChange={(event) => setTeacherId(event.target.value)}
                                    size="small"
                                    fullWidth
                                >
                                    <MenuItem value="">Tất cả giáo viên</MenuItem>
                                    {teachers.map((teacher) => (
                                        <MenuItem key={teacher.id} value={teacher.id}>
                                            {teacher.full_name}
                                        </MenuItem>
                                    ))}
                                </TextField>
                            </Stack>

                            <Stack direction="row" spacing={1} alignItems="center" flexWrap="wrap">
                                <ButtonGroup variant="outlined" size="small">
                                    <Button
                                        variant={viewMode === 'month' ? 'contained' : 'outlined'}
                                        onClick={() => setViewMode('month')}
                                        startIcon={<CalendarMonthRounded />}
                                    >
                                        Tháng
                                    </Button>
                                    <Button
                                        variant={viewMode === 'week' ? 'contained' : 'outlined'}
                                        onClick={() => setViewMode('week')}
                                        startIcon={<ViewWeekRounded />}
                                    >
                                        Tuần
                                    </Button>
                                </ButtonGroup>
                                <DatePicker
                                    label="Đi tới ngày"
                                    value={anchorDate}
                                    onChange={(value) => {
                                        if (value) {
                                            setAnchorDate(value);
                                        }
                                    }}
                                    slotProps={{ textField: { size: 'small' } }}
                                />
                            </Stack>
                        </Stack>

                        <Stack
                            direction={{ xs: 'column', lg: 'row' }}
                            justifyContent="space-between"
                            alignItems={{ xs: 'stretch', lg: 'center' }}
                            spacing={1.5}
                        >
                            <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                                <Chip
                                    size="small"
                                    color="primary"
                                    variant="outlined"
                                    label={hasFilters ? 'Đang lọc dữ liệu' : 'Toàn bộ lesson'}
                                />
                                <Chip
                                    size="small"
                                    variant="outlined"
                                    label={`${totalItems} buổi trong khung hiển thị`}
                                />
                                <Chip
                                    size="small"
                                    variant="outlined"
                                    color={lessonsWithoutRoom > 0 ? 'warning' : 'default'}
                                    label={`${lessonsWithoutRoom} buổi chưa có phòng`}
                                />
                                <Chip
                                    size="small"
                                    variant="outlined"
                                    color={lessonsWithoutTeacher > 0 ? 'warning' : 'default'}
                                    label={`${lessonsWithoutTeacher} buổi chưa có giáo viên`}
                                />
                            </Stack>

                            <Stack direction="row" spacing={1}>
                                <Button
                                    variant="outlined"
                                    startIcon={<ChevronLeftRounded />}
                                    onClick={() => shiftCalendar('prev')}
                                >
                                    {viewMode === 'week' ? 'Tuần trước' : 'Tháng trước'}
                                </Button>
                                <Button
                                    variant="outlined"
                                    startIcon={<TodayRounded />}
                                    onClick={() => setAnchorDate(new Date())}
                                >
                                    Hôm nay
                                </Button>
                                <Button
                                    variant="outlined"
                                    endIcon={<ChevronRightRounded />}
                                    onClick={() => shiftCalendar('next')}
                                >
                                    {viewMode === 'week' ? 'Tuần sau' : 'Tháng sau'}
                                </Button>
                                {hasFilters ? (
                                    <Button
                                        variant="text"
                                        onClick={() => {
                                            setClassId('');
                                            setTeacherId('');
                                        }}
                                    >
                                        Xóa bộ lọc
                                    </Button>
                                ) : null}
                            </Stack>
                        </Stack>
                    </Stack>
                </Paper>

                {isError ? (
                    <Alert
                        severity="error"
                        action={(
                            <Button size="small" onClick={() => void refetch()} startIcon={<RefreshRounded />}>
                                Thử lại
                            </Button>
                        )}
                    >
                        {getApiErrorMessage(error, 'Không thể tải danh sách buổi học. Kiểm tra backend và thử lại.')}
                    </Alert>
                ) : null}

                {isLoading ? (
                    <Stack spacing={1.5}>
                        <Skeleton variant="rounded" height={80} />
                        <Skeleton variant="rounded" height={420} />
                    </Stack>
                ) : (
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
                                    <Typography variant="h6" sx={{ fontWeight: 700, textTransform: 'capitalize' }}>
                                        {headerLabel}
                                    </Typography>
                                    <Typography variant="body2" color="text.secondary">
                                        Hiển thị lịch từ {formatDateTime(queryParams.date_from || '')} đến {formatDateTime(queryParams.date_to || '')}
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
                                gridTemplateColumns: 'repeat(7, minmax(0, 1fr))',
                                borderBottom: '1px solid rgba(15,23,42,0.08)',
                                backgroundColor: '#f8fafc',
                            }}
                        >
                            {weekLabels.map((label) => (
                                <Box key={label} sx={{ px: 1.5, py: 1.25, borderLeft: '1px solid rgba(15,23,42,0.08)' }}>
                                    <Typography variant="subtitle2" sx={{ fontWeight: 700 }}>
                                        {label}
                                    </Typography>
                                </Box>
                            ))}
                        </Box>

                        <Stack spacing={0}>
                            {groupedWeeks.map((week, weekIndex) => (
                                <Box
                                    key={`${weekIndex}-${week[0]?.date.toISOString()}`}
                                    sx={{
                                        display: 'grid',
                                        gridTemplateColumns: 'repeat(7, minmax(0, 1fr))',
                                    }}
                                >
                                    {week.map((day) => (
                                        <Box
                                            key={day.date.toISOString()}
                                            sx={{
                                                minHeight: viewMode === 'month' ? 220 : 320,
                                                p: 1.25,
                                                borderLeft: '1px solid rgba(15,23,42,0.08)',
                                                borderBottom: '1px solid rgba(15,23,42,0.08)',
                                                backgroundColor: isSameMonth(day.date, anchorDate) || viewMode === 'week'
                                                    ? '#ffffff'
                                                    : 'rgba(248,250,252,0.7)',
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
                                                            sx={{
                                                                fontWeight: 700,
                                                                color: isSameDay(day.date, new Date()) ? 'primary.main' : 'text.primary',
                                                            }}
                                                        >
                                                            {format(day.date, 'dd/MM', { locale: vi })}
                                                        </Typography>
                                                        <Typography variant="caption" color="text.secondary">
                                                            {format(day.date, 'EEEE', { locale: vi })}
                                                        </Typography>
                                                    </Box>
                                                    <Chip
                                                        size="small"
                                                        variant="outlined"
                                                        label={`${day.lessons.length} buổi`}
                                                    />
                                                </Stack>

                                                {day.lessons.length > 0 ? (
                                                    <Stack spacing={1}>
                                                        {day.lessons.map(renderLessonCard)}
                                                    </Stack>
                                                ) : (
                                                    <Typography variant="caption" color="text.secondary">
                                                        Không có buổi học trong ngày này.
                                                    </Typography>
                                                )}
                                            </Stack>
                                        </Box>
                                    ))}
                                </Box>
                            ))}
                        </Stack>
                    </Paper>
                )}
            </Stack>
        </LocalizationProvider>
    );
}
