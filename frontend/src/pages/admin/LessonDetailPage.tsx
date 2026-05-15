import { useParams, useNavigate } from 'react-router-dom';
import {
    Alert,
    Box,
    Button,
    Chip,
    Grid,
    Paper,
    Skeleton,
    Stack,
    Typography,
} from '@mui/material';
import { RefreshRounded, ArrowBackRounded, ClassOutlined, PersonOutline, AccessTimeOutlined, MeetingRoomOutlined, NotesOutlined, PersonSearchRounded } from '@mui/icons-material';
import { format } from 'date-fns';
import { vi } from 'date-fns/locale';
import { useState } from 'react';

import { useGetLessonByIdQuery } from '@/api/lessonApi';
import PageHeader from '@/components/common/PageHeader';
import LessonAttendanceManager from '@/components/lesson/LessonAttendanceManager';
import LessonSummaryEditor from '@/components/lesson/LessonSummaryEditor';
import LessonAcademicRecordManager from '@/components/lesson/LessonAcademicRecordManager';
import { SuggestSubstituteModal } from '@/components/schedule/SuggestSubstituteModal';

export default function LessonDetailPage() {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const [isSubstituteModalOpen, setIsSubstituteModalOpen] = useState(false);

    const { data: lessonResponse, isLoading, isError, refetch } = useGetLessonByIdQuery(id!, {
        skip: !id,
    });

    const lesson = lessonResponse?.data;

    const formatTime = (iso: string) => {
        try {
            return format(new Date(iso), 'HH:mm dd/MM/yyyy', { locale: vi });
        } catch {
            return iso;
        }
    };

    if (isLoading) {
        return (
            <Stack sx={{ p: { xs: 2, md: 4 } }} spacing={3}>
                <Skeleton variant="text" width={300} height={60} />
                <Skeleton variant="rounded" height={200} />
            </Stack>
        );
    }

    if (isError || !lesson) {
        return (
            <Stack sx={{ p: { xs: 2, md: 4 } }} spacing={3}>
                <Alert severity="error" action={<Button color="inherit" size="small" onClick={() => navigate('/app/admin/lessons')}>Quay lại</Button>}>
                    Không tìm thấy thông tin buổi học hoặc có lỗi xảy ra.
                </Alert>
            </Stack>
        );
    }

    return (
        <Stack sx={{ p: { xs: 2, md: 4 } }} spacing={3}>
            <PageHeader
                title={`Buổi học lớp ${lesson.class?.name || lesson.class_id}`}
                subtitle="Chi tiết lịch chuẩn bị, thời gian, giáo viên và phòng học."
                icon={<ClassOutlined />}
                actions={
                    <Stack direction="row" spacing={2}>
                        <Button
                            startIcon={<ArrowBackRounded />}
                            variant="outlined"
                            onClick={() => navigate('/app/admin/lessons')}
                        >
                            Quay lại
                        </Button>
                        <Button
                            startIcon={<RefreshRounded />}
                            variant="contained"
                            onClick={() => void refetch()}
                        >
                            Làm mới
                        </Button>
                    </Stack>
                }
            />

            <Grid container spacing={3}>
                <Grid size={{ xs: 12, md: 8 }}>
                    <Paper variant="outlined" sx={{ p: 3, borderRadius: 3, height: '100%' }}>
                        <Typography variant="h6" sx={{ mb: 3, fontWeight: 700 }}>Thông tin thời gian & địa điểm</Typography>
                        <Stack spacing={2.5}>
                            <Stack direction="row" spacing={2} alignItems="flex-start">
                                <AccessTimeOutlined sx={{ color: 'text.secondary', mt: 0.5 }} />
                                <Box>
                                    <Typography variant="body2" color="text.secondary">Thời gian bắt đầu</Typography>
                                    <Typography variant="body1" sx={{ fontWeight: 600 }}>{formatTime(lesson.date_start)}</Typography>
                                </Box>
                            </Stack>
                            <Stack direction="row" spacing={2} alignItems="flex-start">
                                <AccessTimeOutlined sx={{ color: 'text.secondary', mt: 0.5 }} />
                                <Box>
                                    <Typography variant="body2" color="text.secondary">Thời gian kết thúc</Typography>
                                    <Typography variant="body1" sx={{ fontWeight: 600 }}>{formatTime(lesson.date_end)}</Typography>
                                </Box>
                            </Stack>
                            <Stack direction="row" spacing={2} alignItems="flex-start">
                                <MeetingRoomOutlined sx={{ color: 'text.secondary', mt: 0.5 }} />
                                <Box>
                                    <Typography variant="body2" color="text.secondary">Phòng học</Typography>
                                    <Typography variant="body1" sx={{ fontWeight: 600 }}>{lesson.room?.name || 'Chưa xếp phòng'}</Typography>
                                </Box>
                            </Stack>
                        </Stack>
                    </Paper>
                </Grid>

                <Grid size={{ xs: 12, md: 4 }}>
                    <Stack spacing={3} sx={{ height: '100%' }}>
                        <Paper variant="outlined" sx={{ p: 3, borderRadius: 3, flex: 1 }}>
                            <Typography variant="h6" sx={{ mb: 3, fontWeight: 700 }}>Giáo viên</Typography>
                            <Stack spacing={2}>
                                <Stack direction="row" spacing={2} alignItems="center">
                                    <PersonOutline sx={{ color: 'primary.main', fontSize: 32 }} />
                                    <Box>
                                        <Typography variant="body1" sx={{ fontWeight: 600 }}>{lesson.teacher?.full_name || 'Chưa phân công'}</Typography>
                                        {lesson.teacher?.code && (
                                            <Typography variant="body2" color="text.secondary">Mã GV: {lesson.teacher.code}</Typography>
                                        )}
                                    </Box>
                                </Stack>
                                <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                                    <Chip
                                        size="small"
                                        label={lesson.class?.code || 'Buổi học'}
                                        variant="outlined"
                                    />
                                    <Button
                                        variant="outlined"
                                        size="small"
                                        startIcon={<PersonSearchRounded />}
                                        onClick={() => setIsSubstituteModalOpen(true)}
                                    >
                                        Đề xuất dạy thay
                                    </Button>
                                </Stack>
                            </Stack>
                        </Paper>
                        
                        <Paper variant="outlined" sx={{ p: 3, borderRadius: 3, flex: 1 }}>
                            <Typography variant="h6" sx={{ mb: 3, fontWeight: 700 }}>Ghi chú học vụ</Typography>
                            <Stack direction="row" spacing={2} alignItems="flex-start">
                                <NotesOutlined sx={{ color: 'text.secondary', mt: 0.5 }} />
                                <Typography variant="body2" sx={{ whiteSpace: 'pre-line' }}>{lesson.notes || 'Không có ghi chú'}</Typography>
                            </Stack>
                        </Paper>
                    </Stack>
                </Grid>
            </Grid>

            <LessonAttendanceManager lessonId={id!} />
            <LessonSummaryEditor lessonId={id!} />
            <LessonAcademicRecordManager lessonId={id!} />
            <SuggestSubstituteModal
                isOpen={isSubstituteModalOpen}
                lessonId={lesson.id}
                lessonLabel={lesson.class?.name || lesson.class_id}
                onClose={() => setIsSubstituteModalOpen(false)}
            />
        </Stack>
    );
}
