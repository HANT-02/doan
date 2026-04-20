import { useEffect, useMemo, useState } from 'react';
import {
    Alert,
    Chip,
    List,
    ListItemButton,
    ListItemText,
    MenuItem,
    Paper,
    Stack,
    TextField,
    Typography,
} from '@mui/material';
import { CheckCircleRounded, MenuBookRounded, PendingActionsRounded } from '@mui/icons-material';
import { format, parseISO } from 'date-fns';
import { vi } from 'date-fns/locale';
import { useSearchParams } from 'react-router-dom';

import {
    type TeacherLesson,
    useGetTeacherLessonSummaryQuery,
    useGetTeacherLessonsQuery,
} from '@/api/teacherPortalApi';
import PageHeader from '@/components/common/PageHeader';
import TeacherLessonSummaryEditor from '@/components/lesson/TeacherLessonSummaryEditor';

function TeacherLessonSummaryStatusItem({
    lesson,
    selected,
    onSelect,
}: {
    lesson: TeacherLesson;
    selected: boolean;
    onSelect: (lessonId: string) => void;
}) {
    const summaryQuery = useGetTeacherLessonSummaryQuery(lesson.id);
    const hasSummary = Boolean(summaryQuery.data?.data?.summary);

    return (
        <ListItemButton
            selected={selected}
            onClick={() => onSelect(lesson.id)}
            sx={{
                borderRadius: 2,
                alignItems: 'flex-start',
                border: '1px solid',
                borderColor: selected ? 'primary.main' : 'divider',
                mb: 1,
            }}
        >
            <ListItemText
                primary={(
                    <Stack direction="row" justifyContent="space-between" alignItems="center" spacing={1}>
                        <Typography variant="subtitle2" sx={{ fontWeight: 700 }}>
                            {format(parseISO(lesson.date_start), 'dd/MM/yyyy HH:mm', { locale: vi })}
                        </Typography>
                        <Chip
                            size="small"
                            color={hasSummary ? 'success' : 'warning'}
                            icon={hasSummary ? <CheckCircleRounded /> : <PendingActionsRounded />}
                            label={hasSummary ? 'Đã tổng kết' : 'Chưa tổng kết'}
                        />
                    </Stack>
                )}
                secondary={(
                    <Stack spacing={0.5} sx={{ mt: 0.75 }}>
                        <Typography variant="body2" color="text.primary">
                            {lesson.class_name}
                        </Typography>
                        <Typography variant="caption" color="text.secondary">
                            {lesson.shift?.name || 'Chưa gắn ca học'} • {lesson.room_name || 'Chưa xếp phòng'}
                        </Typography>
                    </Stack>
                )}
            />
        </ListItemButton>
    );
}

export default function TeacherJournalPage() {
    const [searchParams] = useSearchParams();
    const preselectedLessonId = searchParams.get('lessonId') || '';
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
    const [hasAppliedPreselection, setHasAppliedPreselection] = useState(false);

    useEffect(() => {
        if (!lessons.length) {
            return;
        }

        if (!hasAppliedPreselection && preselectedLessonId) {
            const preselectedLesson = lessons.find((lesson) => lesson.id === preselectedLessonId);
            setHasAppliedPreselection(true);
            if (preselectedLesson) {
                setSelectedClassId(preselectedLesson.class_id);
                setSelectedLessonId(preselectedLesson.id);
                return;
            }
        }

        if (!selectedClassId) {
            setSelectedClassId(lessons[0].class_id);
        }
    }, [hasAppliedPreselection, lessons, preselectedLessonId, selectedClassId]);

    const filteredLessons = useMemo(
        () =>
            lessons
                .filter((lesson) => !selectedClassId || lesson.class_id === selectedClassId)
                .sort((left, right) => new Date(left.date_start).getTime() - new Date(right.date_start).getTime()),
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
        () => filteredLessons.find((lesson) => lesson.id === selectedLessonId) || null,
        [filteredLessons, selectedLessonId],
    );

    return (
        <Stack sx={{ p: { xs: 2, md: 4 } }} spacing={3}>
            <PageHeader
                title="Sổ đầu bài"
                subtitle="Chọn lớp, chọn buổi học và ghi tổng kết ngay trên nhánh teacher portal."
                icon={<MenuBookRounded />}
                breadcrumbs={[
                    { label: 'Trang giáo viên', path: '/app/teacher/overview' },
                    { label: 'Sổ đầu bài' },
                ]}
            />

            <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 3 }}>
                <Stack spacing={2}>
                    <Typography variant="subtitle1" sx={{ fontWeight: 700 }}>
                        Chọn lớp và buổi học
                    </Typography>
                    <Stack direction={{ xs: 'column', md: 'row' }} spacing={2}>
                        <TextField
                            select
                            label="Lớp học"
                            value={selectedClassId}
                            onChange={(event) => setSelectedClassId(event.target.value)}
                            fullWidth
                            disabled={isLoadingLessons || !classes.length}
                        >
                            {classes.map((item) => (
                                <MenuItem key={item.id} value={item.id}>
                                    {item.code} - {item.name}
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
                                    {format(parseISO(lesson.date_start), 'dd/MM/yyyy HH:mm', { locale: vi })} - {lesson.shift?.name || lesson.class_name}
                                </MenuItem>
                            ))}
                        </TextField>
                    </Stack>
                    <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                        <Chip label={`${classes.length} lớp`} variant="outlined" />
                        <Chip label={`${filteredLessons.length} buổi trong lớp`} variant="outlined" />
                        <Chip label={selectedLesson ? 'Đã chọn 1 buổi để ghi sổ' : 'Chưa chọn buổi'} color="primary" variant="outlined" />
                    </Stack>
                </Stack>
            </Paper>

            {!lessons.length && !isLoadingLessons ? (
                <Alert severity="info">Chưa có buổi học nào được phân công cho giáo viên này.</Alert>
            ) : null}

            {selectedLesson ? (
                <Stack direction={{ xs: 'column', xl: 'row' }} spacing={3} alignItems="stretch">
                    <Paper variant="outlined" sx={{ p: 2, borderRadius: 3, width: { xs: '100%', xl: 360 }, flexShrink: 0 }}>
                        <Stack spacing={1.5}>
                            <Typography variant="subtitle1" sx={{ fontWeight: 700 }}>
                                Điều hướng theo buổi học
                            </Typography>
                            <Typography variant="body2" color="text.secondary">
                                Danh sách dưới đây cho biết buổi nào đã có tổng kết và buổi nào còn thiếu.
                            </Typography>
                            <List disablePadding>
                                {filteredLessons.map((lesson) => (
                                    <TeacherLessonSummaryStatusItem
                                        key={lesson.id}
                                        lesson={lesson}
                                        selected={lesson.id === selectedLessonId}
                                        onSelect={setSelectedLessonId}
                                    />
                                ))}
                            </List>
                        </Stack>
                    </Paper>

                    <Stack spacing={3} sx={{ flex: 1, minWidth: 0 }}>
                        <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 3 }}>
                            <Stack spacing={1}>
                                <Typography variant="h6" sx={{ fontWeight: 700 }}>
                                    {selectedLesson.class_name}
                                </Typography>
                                <Typography variant="body2" color="text.secondary">
                                    {selectedLesson.class_code} • {format(parseISO(selectedLesson.date_start), 'EEEE, dd/MM/yyyy HH:mm', { locale: vi })} - {format(parseISO(selectedLesson.date_end), 'HH:mm')}
                                </Typography>
                                <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                                    <Chip label={selectedLesson.shift?.name || 'Chưa gắn ca học'} color="primary" variant="outlined" />
                                    <Chip label={selectedLesson.room_name || 'Chưa xếp phòng'} variant="outlined" />
                                </Stack>
                                {selectedLesson.notes ? (
                                    <Typography variant="body2" color="text.secondary">
                                        {selectedLesson.notes}
                                    </Typography>
                                ) : null}
                            </Stack>
                        </Paper>

                        <TeacherLessonSummaryEditor
                            lessonId={selectedLessonId}
                            title="Tổng kết buổi học của giáo viên"
                            subtitle="Lưu chủ đề đã dạy, phản hồi lớp, bài tập và ghi chú sau buổi học đang chọn."
                        />
                    </Stack>
                </Stack>
            ) : null}
        </Stack>
    );
}
