import { useEffect, useState } from 'react';
import {
    Alert,
    MenuItem,
    Paper,
    Stack,
    TextField,
    Typography,
} from '@mui/material';
import { useSearchParams } from 'react-router-dom';

import { useGetTeacherLessonsQuery } from '@/api/teacherPortalApi';
import PageHeader from '@/components/common/PageHeader';
import LessonSummaryEditor from '@/components/lesson/LessonSummaryEditor';
import LessonAcademicRecordManager from '@/components/lesson/LessonAcademicRecordManager';

export default function TeacherJournalPage() {
    const [searchParams] = useSearchParams();
    const preselectedLessonId = searchParams.get('lessonId') || '';
    const { data: lessonsResponse, isLoading: isLoadingLessons } = useGetTeacherLessonsQuery();

    const lessons = lessonsResponse?.data?.lessons || [];
    const [selectedLessonId, setSelectedLessonId] = useState('');

    useEffect(() => {
        if (preselectedLessonId && lessons.some((lesson) => lesson.id === preselectedLessonId)) {
            setSelectedLessonId(preselectedLessonId);
            return;
        }
        if (!selectedLessonId && lessons.length > 0) {
            setSelectedLessonId(lessons[0].id);
        }
    }, [lessons, selectedLessonId, preselectedLessonId]);

    return (
        <Stack sx={{ p: { xs: 2, md: 4 } }} spacing={3}>
            <PageHeader
                title="Sổ đầu bài"
                subtitle="Ghi nhận nội dung đã dạy, phản hồi lớp và bài tập sau từng buổi học."
            />

            <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 3 }}>
                <Stack spacing={2}>
                    <Typography variant="subtitle1" sx={{ fontWeight: 700 }}>
                        Chọn buổi học để ghi sổ đầu bài
                    </Typography>
                    <TextField
                        select
                        label="Buổi học"
                        value={selectedLessonId}
                        onChange={(event) => setSelectedLessonId(event.target.value)}
                        fullWidth
                        disabled={isLoadingLessons || !lessons.length}
                    >
                        {lessons.map((lesson) => (
                            <MenuItem key={lesson.id} value={lesson.id}>
                                {lesson.class_name} - {new Date(lesson.date_start).toLocaleString('vi-VN')}
                            </MenuItem>
                        ))}
                    </TextField>
                </Stack>
            </Paper>

            {selectedLessonId ? (
                <Stack spacing={3}>
                    <LessonSummaryEditor
                        lessonId={selectedLessonId}
                        title="Tổng kết buổi học của giáo viên"
                        subtitle="Lưu chủ đề đã dạy, phản hồi lớp và bài tập được giao cho buổi học đang chọn."
                    />
                    <LessonAcademicRecordManager
                        lessonId={selectedLessonId}
                        title="Kết quả học tập của buổi học"
                        subtitle="Nhập điểm bài tập, thái độ, tham gia và chốt kết quả cho học sinh trong buổi học."
                    />
                </Stack>
            ) : (
                <Alert severity="info">
                    {isLoadingLessons
                        ? 'Đang tải danh sách buổi học...'
                        : 'Chưa có buổi học nào được phân công cho giáo viên này.'}
                </Alert>
            )}
        </Stack>
    );
}
