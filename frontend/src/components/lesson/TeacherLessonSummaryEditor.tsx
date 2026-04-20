import { useEffect, useMemo, useState } from 'react';
import {
    Alert,
    Button,
    Chip,
    Paper,
    Stack,
    TextField,
    Typography,
} from '@mui/material';
import { BookRounded, SaveRounded } from '@mui/icons-material';
import { format, parseISO } from 'date-fns';
import { vi } from 'date-fns/locale';
import { toast } from 'sonner';

import {
    useGetTeacherLessonSummaryQuery,
    useUpsertTeacherLessonSummaryMutation,
} from '@/api/teacherPortalApi';
import { getApiErrorMessage } from '@/utils/apiError';

interface TeacherLessonSummaryEditorProps {
    lessonId: string;
    title?: string;
    subtitle?: string;
}

type SummaryFormState = {
    topic: string;
    lesson_content: string;
    class_feedback: string;
    homework: string;
    homework_deadline: string;
    teacher_notes: string;
};

const emptyForm: SummaryFormState = {
    topic: '',
    lesson_content: '',
    class_feedback: '',
    homework: '',
    homework_deadline: '',
    teacher_notes: '',
};

const toDateTimeLocalInput = (iso?: string) => {
    if (!iso) {
        return '';
    }
    const parsed = new Date(iso);
    if (Number.isNaN(parsed.getTime())) {
        return '';
    }
    const year = parsed.getFullYear();
    const month = String(parsed.getMonth() + 1).padStart(2, '0');
    const day = String(parsed.getDate()).padStart(2, '0');
    const hours = String(parsed.getHours()).padStart(2, '0');
    const minutes = String(parsed.getMinutes()).padStart(2, '0');
    return `${year}-${month}-${day}T${hours}:${minutes}`;
};

export default function TeacherLessonSummaryEditor({
    lessonId,
    title = 'Tổng kết buổi học',
    subtitle = 'Ghi lại nội dung dạy, phản hồi lớp và bài tập sau buổi học.',
}: TeacherLessonSummaryEditorProps) {
    const { data, isLoading, isFetching, error, refetch } = useGetTeacherLessonSummaryQuery(lessonId, {
        skip: !lessonId,
    });
    const [upsertSummary, { isLoading: isSaving }] = useUpsertTeacherLessonSummaryMutation();
    const [form, setForm] = useState<SummaryFormState>(emptyForm);

    useEffect(() => {
        const summary = data?.data?.summary;
        setForm({
            topic: summary?.topic || '',
            lesson_content: summary?.lesson_content || '',
            class_feedback: summary?.class_feedback || '',
            homework: summary?.homework || '',
            homework_deadline: toDateTimeLocalInput(summary?.homework_deadline),
            teacher_notes: summary?.teacher_notes || '',
        });
    }, [data]);

    const completionCount = useMemo(() => {
        return [
            form.topic,
            form.lesson_content,
            form.class_feedback,
            form.homework,
            form.teacher_notes,
        ].filter((value) => value.trim().length > 0).length;
    }, [form]);

    const handleSave = async () => {
        try {
            await upsertSummary({
                lessonId,
                body: {
                    topic: form.topic.trim(),
                    lesson_content: form.lesson_content.trim(),
                    class_feedback: form.class_feedback.trim(),
                    homework: form.homework.trim(),
                    homework_deadline: form.homework_deadline ? new Date(form.homework_deadline).toISOString() : null,
                    teacher_notes: form.teacher_notes.trim(),
                },
            }).unwrap();
            toast.success('Đã lưu sổ đầu bài cho buổi học.');
            refetch();
        } catch (err) {
            toast.error(getApiErrorMessage(err, 'Không lưu được tổng kết buổi học.'));
        }
    };

    return (
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
                            <BookRounded color="primary" />
                            <Typography variant="h6" sx={{ fontWeight: 700 }}>
                                {title}
                            </Typography>
                        </Stack>
                        <Typography variant="body2" color="text.secondary">
                            {subtitle}
                        </Typography>
                    </Stack>

                    <Button
                        variant="contained"
                        startIcon={<SaveRounded />}
                        onClick={handleSave}
                        disabled={isSaving || isLoading || isFetching}
                    >
                        {isSaving ? 'Đang lưu...' : 'Lưu tổng kết'}
                    </Button>
                </Stack>

                <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                    <Chip
                        size="small"
                        variant="outlined"
                        color={data?.data?.summary ? 'success' : 'default'}
                        label={data?.data?.summary ? 'Đã có tổng kết' : 'Chưa có tổng kết'}
                    />
                    <Chip
                        size="small"
                        color="primary"
                        variant="outlined"
                        label={`${completionCount}/5 trường chính đã điền`}
                    />
                    {data?.data?.summary?.updated_at ? (
                        <Chip
                            size="small"
                            variant="outlined"
                            label={`Cập nhật ${format(parseISO(data.data.summary.updated_at), 'dd/MM/yyyy HH:mm', { locale: vi })}`}
                        />
                    ) : null}
                </Stack>

                {error ? (
                    <Alert severity="error">
                        {getApiErrorMessage(error, 'Không tải được dữ liệu tổng kết buổi học.')}
                    </Alert>
                ) : null}

                <TextField
                    label="Chủ đề buổi học"
                    value={form.topic}
                    onChange={(event) => setForm((current) => ({ ...current, topic: event.target.value }))}
                    fullWidth
                />
                <TextField
                    label="Nội dung đã dạy"
                    value={form.lesson_content}
                    onChange={(event) => setForm((current) => ({ ...current, lesson_content: event.target.value }))}
                    fullWidth
                    multiline
                    minRows={3}
                />
                <TextField
                    label="Phản hồi lớp học"
                    value={form.class_feedback}
                    onChange={(event) => setForm((current) => ({ ...current, class_feedback: event.target.value }))}
                    fullWidth
                    multiline
                    minRows={2}
                />
                <TextField
                    label="Bài tập / giao việc"
                    value={form.homework}
                    onChange={(event) => setForm((current) => ({ ...current, homework: event.target.value }))}
                    fullWidth
                    multiline
                    minRows={2}
                />
                <TextField
                    label="Hạn nộp bài tập"
                    type="datetime-local"
                    value={form.homework_deadline}
                    onChange={(event) => setForm((current) => ({ ...current, homework_deadline: event.target.value }))}
                    fullWidth
                    InputLabelProps={{ shrink: true }}
                />
                <TextField
                    label="Ghi chú của giáo viên"
                    value={form.teacher_notes}
                    onChange={(event) => setForm((current) => ({ ...current, teacher_notes: event.target.value }))}
                    fullWidth
                    multiline
                    minRows={3}
                />
            </Stack>
        </Paper>
    );
}
