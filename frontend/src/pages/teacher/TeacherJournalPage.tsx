import { useEffect, useMemo, useState } from 'react';
import {
    Alert,
    MenuItem,
    Paper,
    Stack,
    TextField,
    Typography,
} from '@mui/material';
import { skipToken } from '@reduxjs/toolkit/query';

import { useGetLessonsQuery } from '@/api/lessonApi';
import { useGetTeachersQuery } from '@/api/teacherApi';
import PageHeader from '@/components/common/PageHeader';
import LessonSummaryEditor from '@/components/lesson/LessonSummaryEditor';
import LessonAcademicRecordManager from '@/components/lesson/LessonAcademicRecordManager';
import { useAuth } from '@/contexts/AuthContext';

export default function TeacherJournalPage() {
    const { user } = useAuth();
    const { data: teachersResponse, isLoading: isLoadingTeachers } = useGetTeachersQuery(
        user?.email ? { search: user.email, limit: 20 } : skipToken,
    );

    const currentTeacher = useMemo(() => {
        const teachers = teachersResponse?.data?.teachers || [];
        if (!user?.email) {
            return null;
        }
        return teachers.find((teacher) => teacher.email === user.email) || teachers[0] || null;
    }, [teachersResponse, user?.email]);

    const { data: lessonsResponse, isLoading: isLoadingLessons } = useGetLessonsQuery(
        currentTeacher
            ? {
                teacher_id: currentTeacher.id,
                limit: 200,
                sortBy: 'date_start',
                sortOrder: 'asc',
            }
            : skipToken,
    );

    const lessons = lessonsResponse?.data?.lessons || [];
    const [selectedLessonId, setSelectedLessonId] = useState('');

    useEffect(() => {
        if (!selectedLessonId && lessons.length > 0) {
            setSelectedLessonId(lessons[0].id);
        }
    }, [lessons, selectedLessonId]);

    return (
        <Stack sx={{ p: { xs: 2, md: 4 } }} spacing={3}>
            <PageHeader
                title="Sổ đầu bài"
                subtitle="Ghi nhận nội dung đã dạy, phản hồi lớp và bài tập sau từng buổi học."
            />

            {!user?.email ? (
                <Alert severity="warning">Không tìm thấy email người dùng hiện tại để đối chiếu hồ sơ giáo viên.</Alert>
            ) : null}

            {!isLoadingTeachers && user?.email && !currentTeacher ? (
                <Alert severity="warning">
                    Không tìm thấy hồ sơ giáo viên trùng với email đăng nhập hiện tại. Hãy kiểm tra dữ liệu teacher trong hệ thống.
                </Alert>
            ) : null}

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
                        disabled={isLoadingTeachers || isLoadingLessons || !lessons.length}
                    >
                        {lessons.map((lesson) => (
                            <MenuItem key={lesson.id} value={lesson.id}>
                                {lesson.class?.name || lesson.class_id} - {new Date(lesson.date_start).toLocaleString('vi-VN')}
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
                    {isLoadingTeachers || isLoadingLessons
                        ? 'Đang tải danh sách buổi học...'
                        : 'Chưa có buổi học nào được phân công cho giáo viên này.'}
                </Alert>
            )}
        </Stack>
    );
}
