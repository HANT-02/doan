import { useEffect, useMemo, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import {
    Alert,
    Button,
    Chip,
    Paper,
    Stack,
    TextField,
    Typography,
} from '@mui/material';
import { CachedRounded, PersonSearchRounded } from '@mui/icons-material';
import { format, isAfter, parseISO } from 'date-fns';
import { vi } from 'date-fns/locale';

import { useGetTeacherLessonsQuery } from '@/api/teacherPortalApi';
import PageHeader from '@/components/common/PageHeader';
import { SuggestSubstituteModal } from '@/components/schedule/SuggestSubstituteModal';
import { getApiErrorMessage } from '@/utils/apiError';

export const TeacherSubstitutePage = () => {
    const [searchParams] = useSearchParams();
    const preselectedLessonId = searchParams.get('lessonId') || '';
    const [selectedLessonId, setSelectedLessonId] = useState<string | null>(null);
    const [search, setSearch] = useState('');

    const { data: lessonsData, isLoading, isFetching, error, refetch } = useGetTeacherLessonsQuery();
    const lessons = lessonsData?.data?.lessons || [];

    const upcomingLessons = useMemo(() => {
        const keyword = search.trim().toLowerCase();
        return lessons
            .filter((lesson) => isAfter(parseISO(lesson.date_start), new Date()))
            .filter((lesson) => {
                if (!keyword) {
                    return true;
                }
                return [lesson.class_name, lesson.class_code, lesson.room_name]
                    .filter((value): value is string => !!value)
                    .some((value) => value.toLowerCase().includes(keyword));
            })
            .sort((left, right) => left.date_start.localeCompare(right.date_start));
    }, [lessons, search]);

    useEffect(() => {
        if (!preselectedLessonId || selectedLessonId) {
            return;
        }
        const preselected = lessons.find((lesson) => lesson.id === preselectedLessonId);
        if (preselected) {
            setSelectedLessonId(preselected.id);
        }
    }, [lessons, preselectedLessonId, selectedLessonId]);

    const selectedLesson = useMemo(
        () => lessons.find((lesson) => lesson.id === selectedLessonId) || null,
        [lessons, selectedLessonId],
    );

    return (
        <Stack sx={{ p: { xs: 2, md: 4 } }} spacing={3}>
            <PageHeader
                title="Điều phối dạy thay"
                subtitle="Chọn một ca dạy sắp tới để hệ thống gợi ý giáo viên thay thế dựa trên kỹ năng, lịch trống, tải dạy và khả năng di chuyển."
            />

            <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 3 }}>
                <Stack spacing={2}>
                    <Stack direction={{ xs: 'column', md: 'row' }} spacing={1.5} justifyContent="space-between">
                        <TextField
                            value={search}
                            onChange={(event) => setSearch(event.target.value)}
                            label="Tìm theo lớp hoặc phòng"
                            placeholder="Ví dụ: Toán 9A, P.201..."
                            fullWidth
                        />
                        <Button
                            variant="outlined"
                            startIcon={<CachedRounded />}
                            onClick={() => void refetch()}
                            disabled={isFetching}
                        >
                            Làm mới
                        </Button>
                    </Stack>

                    {error ? (
                        <Alert severity="error">
                            {getApiErrorMessage(error, 'Không tải được danh sách buổi học sắp tới.')}
                        </Alert>
                    ) : null}

                    {!error && !isLoading && !upcomingLessons.length ? (
                        <Alert severity="info">
                            Không có ca dạy sắp tới phù hợp với bộ lọc hiện tại.
                        </Alert>
                    ) : null}

                    {upcomingLessons.map((lesson) => (
                        <Paper key={lesson.id} variant="outlined" sx={{ p: 2, borderRadius: 3 }}>
                            <Stack direction={{ xs: 'column', lg: 'row' }} spacing={2} justifyContent="space-between">
                                <Stack spacing={0.75}>
                                    <Typography variant="subtitle1" sx={{ fontWeight: 800 }}>
                                        {lesson.class_name} ({lesson.class_code})
                                    </Typography>
                                    <Typography variant="body2" color="text.secondary">
                                        {format(parseISO(lesson.date_start), 'EEEE, dd/MM/yyyy HH:mm', { locale: vi })} - {format(parseISO(lesson.date_end), 'HH:mm', { locale: vi })}
                                    </Typography>
                                    <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                                        <Chip size="small" label={lesson.shift?.name || 'Chưa gắn ca học'} variant="outlined" />
                                        <Chip size="small" label={lesson.room_name || 'Chưa xếp phòng'} variant="outlined" />
                                    </Stack>
                                </Stack>

                                <Button
                                    variant="contained"
                                    startIcon={<PersonSearchRounded />}
                                    onClick={() => setSelectedLessonId(lesson.id)}
                                >
                                    Đề xuất dạy thay
                                </Button>
                            </Stack>
                        </Paper>
                    ))}
                </Stack>
            </Paper>

            <SuggestSubstituteModal
                isOpen={!!selectedLessonId}
                lessonId={selectedLessonId}
                lessonLabel={selectedLesson ? `${selectedLesson.class_name} (${selectedLesson.class_code})` : undefined}
                onClose={() => setSelectedLessonId(null)}
            />
        </Stack>
    );
};
