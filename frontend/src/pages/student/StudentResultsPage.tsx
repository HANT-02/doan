import {
    Alert,
    Chip,
    Paper,
    Stack,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Typography,
} from '@mui/material';

import { useGetMyAcademicRecordsQuery } from '@/api/academicApi';
import PageHeader from '@/components/common/PageHeader';
import { getApiErrorMessage } from '@/utils/apiError';

export default function StudentResultsPage() {
    const { data, isLoading, error } = useGetMyAcademicRecordsQuery();
    const records = data?.data?.records || [];

    return (
        <Stack sx={{ p: { xs: 2, md: 4 } }} spacing={3}>
            <PageHeader
                title="Kết quả học tập"
                subtitle="Theo dõi điểm bài tập, mức độ tham gia và đánh giá giáo viên theo từng buổi học."
            />

            {error ? (
                <Alert severity="error">
                    {getApiErrorMessage(error, 'Không tải được kết quả học tập của bạn.')}
                </Alert>
            ) : null}

            <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 3 }}>
                <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                    <Chip size="small" label={`${records.length} bản ghi`} variant="outlined" />
                    <Chip
                        size="small"
                        label={`Đã chốt ${records.filter((record) => record.is_completed).length}`}
                        color="success"
                        variant="outlined"
                    />
                </Stack>
            </Paper>

            <TableContainer component={Paper} variant="outlined" sx={{ borderRadius: 3 }}>
                <Table>
                    <TableHead>
                        <TableRow>
                            <TableCell sx={{ fontWeight: 700 }}>Lớp / buổi học</TableCell>
                            <TableCell sx={{ fontWeight: 700 }}>Bài tập</TableCell>
                            <TableCell sx={{ fontWeight: 700 }}>Thái độ</TableCell>
                            <TableCell sx={{ fontWeight: 700 }}>Tham gia</TableCell>
                            <TableCell sx={{ fontWeight: 700 }}>Tổng</TableCell>
                            <TableCell sx={{ fontWeight: 700 }}>Nhận xét</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {records.map((record) => (
                            <TableRow key={record.id} hover>
                                <TableCell>
                                    <Typography variant="body2" sx={{ fontWeight: 600 }}>
                                        {record.lesson_summary?.lesson?.class?.name || 'Lớp học'}
                                    </Typography>
                                    <Typography variant="caption" color="text.secondary">
                                        {record.lesson_summary?.lesson?.date_start
                                            ? new Date(record.lesson_summary.lesson.date_start).toLocaleString('vi-VN')
                                            : 'Không rõ thời gian'}
                                    </Typography>
                                </TableCell>
                                <TableCell>{record.homework_completed ? `${record.homework_score}/10` : 'Chưa hoàn thành'}</TableCell>
                                <TableCell>{record.attitude_rating}/5</TableCell>
                                <TableCell>{record.participation_score}/10</TableCell>
                                <TableCell>
                                    <Chip
                                        size="small"
                                        label={record.total_score?.toFixed(2) || '0.00'}
                                        color={record.is_completed ? 'success' : 'default'}
                                        variant="outlined"
                                    />
                                </TableCell>
                                <TableCell>{record.personal_comment || 'Không có nhận xét'}</TableCell>
                            </TableRow>
                        ))}
                        {!records.length && !isLoading ? (
                            <TableRow>
                                <TableCell colSpan={6}>
                                    <Typography variant="body2" color="text.secondary">
                                        Chưa có kết quả học tập nào được công bố.
                                    </Typography>
                                </TableCell>
                            </TableRow>
                        ) : null}
                    </TableBody>
                </Table>
            </TableContainer>
        </Stack>
    );
}
