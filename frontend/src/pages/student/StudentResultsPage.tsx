import { useMemo, useState } from 'react';
import {
    Alert,
    Chip,
    LinearProgress,
    MenuItem,
    Paper,
    Stack,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    TextField,
    Typography,
} from '@mui/material';
import { BarChartRounded, SchoolRounded, WarningAmberRounded } from '@mui/icons-material';
import { format, parseISO } from 'date-fns';

import { useGetStudentAcademicRecordsQuery, useGetStudentAtRiskPredictionQuery } from '@/api/studentPortalApi';
import PageHeader from '@/components/common/PageHeader';
import { getApiErrorMessage } from '@/utils/apiError';

function clampProgress(score: number, max: number) {
    if (max <= 0) return 0;
    return Math.max(0, Math.min(100, (score / max) * 100));
}

export default function StudentResultsPage() {
    const [classFilter, setClassFilter] = useState('');
    const { data, isLoading, error } = useGetStudentAcademicRecordsQuery({
        class_id: classFilter || undefined,
    });
    const { data: atRiskResponse } = useGetStudentAtRiskPredictionQuery();

    const classSummaries = data?.data?.class_summaries ?? [];
    const records = data?.data?.records ?? [];
    const prediction = atRiskResponse?.data?.prediction ?? null;

    const availableClasses = useMemo(() => {
        const map = new Map<string, { id: string; code: string; name: string }>();

        classSummaries.forEach((item) => {
            map.set(item.class_id, { id: item.class_id, code: item.class_code, name: item.class_name });
        });
        records.forEach((item) => {
            if (!map.has(item.lesson.class_id)) {
                map.set(item.lesson.class_id, {
                    id: item.lesson.class_id,
                    code: item.lesson.class_code,
                    name: item.lesson.class_name,
                });
            }
        });

        return Array.from(map.values()).sort((a, b) => a.name.localeCompare(b.name));
    }, [classSummaries, records]);

    const selectedSummary = useMemo(
        () => classSummaries.find((item) => item.class_id === classFilter) ?? null,
        [classFilter, classSummaries],
    );

    const completedCount = records.filter((record) => record.is_completed).length;
    const averageScore =
        records.length > 0
            ? records.reduce((sum, record) => sum + record.total_score, 0) / records.length
            : 0;

    return (
        <Stack sx={{ p: { xs: 2, md: 4 } }} spacing={3}>
            <PageHeader
                title="Kết quả học tập"
                subtitle="Theo dõi điểm theo từng buổi, xem lớp đang học và nắm nhanh xu hướng tiến độ học tập của bạn."
                icon={<SchoolRounded />}
                breadcrumbs={[
                    { label: 'Cổng học sinh', path: '/app/student/overview' },
                    { label: 'Kết quả học tập' },
                ]}
            />

            {error ? (
                <Alert severity="error">
                    {getApiErrorMessage(error, 'Không tải được kết quả học tập của bạn.')}
                </Alert>
            ) : null}

            <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 3 }}>
                <Stack direction={{ xs: 'column', md: 'row' }} spacing={2} alignItems={{ xs: 'stretch', md: 'center' }}>
                    <TextField
                        select
                        label="Chọn lớp"
                        value={classFilter}
                        onChange={(event) => setClassFilter(event.target.value)}
                        sx={{ minWidth: 280 }}
                    >
                        <MenuItem value="">Tất cả lớp</MenuItem>
                        {availableClasses.map((item) => (
                            <MenuItem key={item.id} value={item.id}>
                                {item.name} ({item.code})
                            </MenuItem>
                        ))}
                    </TextField>
                    <Typography variant="body2" color="text.secondary">
                        Chọn một lớp để xem riêng kết quả theo từng buổi học và thống kê tổng hợp của lớp đó.
                    </Typography>
                </Stack>
            </Paper>

            <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                <Chip size="small" label={`${records.length} bản ghi`} variant="outlined" />
                <Chip size="small" label={`Đã chốt ${completedCount}`} color="success" variant="outlined" />
                <Chip size="small" label={`Điểm TB ${averageScore.toFixed(2)}`} color="primary" variant="outlined" />
                <Chip size="small" label={`${classSummaries.length} lớp có dữ liệu`} variant="outlined" />
            </Stack>

            {prediction ? (
                <Alert
                    severity={prediction.risk_label === 'AT_RISK' ? 'warning' : 'info'}
                    icon={<WarningAmberRounded fontSize="inherit" />}
                >
                    <strong>
                        {prediction.risk_label === 'AT_RISK'
                            ? `Cảnh báo sớm: nguy cơ học kém ${Math.round(prediction.risk_score * 100)}%`
                            : `Theo dõi học tập: mức nguy cơ hiện tại ${Math.round(prediction.risk_score * 100)}%`}
                    </strong>
                    {` • ${prediction.primary_reason}`}
                </Alert>
            ) : null}

            <Stack direction={{ xs: 'column', lg: 'row' }} spacing={2}>
                <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 3, flex: 1 }}>
                    <Stack spacing={2}>
                        <Stack direction="row" spacing={1} alignItems="center">
                            <BarChartRounded color="primary" />
                            <Typography variant="h6" sx={{ fontWeight: 700 }}>
                                Tổng hợp lớp
                            </Typography>
                        </Stack>

                        {selectedSummary ? (
                            <Stack spacing={1.5}>
                                <Typography variant="subtitle1" sx={{ fontWeight: 700 }}>
                                    {selectedSummary.class_name}
                                </Typography>
                                <Typography variant="body2" color="text.secondary">
                                    {selectedSummary.class_code} • {selectedSummary.completed_count}/{selectedSummary.records_count} bản ghi đã chốt
                                </Typography>
                                <Chip
                                    label={`Điểm trung bình ${selectedSummary.average_total_score.toFixed(2)}`}
                                    color="primary"
                                    variant="outlined"
                                    sx={{ width: 'fit-content' }}
                                />
                                <Stack spacing={1}>
                                    <Typography variant="caption" color="text.secondary">
                                        Tiến độ chốt kết quả
                                    </Typography>
                                    <LinearProgress
                                        variant="determinate"
                                        value={selectedSummary.records_count ? (selectedSummary.completed_count / selectedSummary.records_count) * 100 : 0}
                                        sx={{ height: 10, borderRadius: 999 }}
                                    />
                                </Stack>
                            </Stack>
                        ) : (
                            <Stack spacing={1.5}>
                                {classSummaries.map((summary) => (
                                    <Paper key={summary.class_id} variant="outlined" sx={{ p: 1.5, borderRadius: 2 }}>
                                        <Stack spacing={0.5}>
                                            <Typography variant="body2" sx={{ fontWeight: 700 }}>
                                                {summary.class_name}
                                            </Typography>
                                            <Typography variant="caption" color="text.secondary">
                                                {summary.class_code}
                                            </Typography>
                                            <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                                                <Chip size="small" label={`TB ${summary.average_total_score.toFixed(2)}`} color="primary" variant="outlined" />
                                                <Chip size="small" label={`Đã chốt ${summary.completed_count}/${summary.records_count}`} variant="outlined" />
                                            </Stack>
                                        </Stack>
                                    </Paper>
                                ))}
                                {!classSummaries.length && !isLoading ? (
                                    <Typography variant="body2" color="text.secondary">
                                        Chưa có dữ liệu tổng hợp theo lớp.
                                    </Typography>
                                ) : null}
                            </Stack>
                        )}
                    </Stack>
                </Paper>

                <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 3, width: { xs: '100%', lg: 360 } }}>
                    <Stack spacing={2}>
                        <Typography variant="h6" sx={{ fontWeight: 700 }}>
                            Tiến độ học tập
                        </Typography>
                        <Stack spacing={1.5}>
                            <Typography variant="caption" color="text.secondary">
                                Điểm tổng trung bình
                            </Typography>
                            <LinearProgress
                                variant="determinate"
                                value={clampProgress(selectedSummary?.average_total_score ?? averageScore, 10)}
                                sx={{ height: 10, borderRadius: 999 }}
                            />
                            <Typography variant="body2" sx={{ fontWeight: 600 }}>
                                {(selectedSummary?.average_total_score ?? averageScore).toFixed(2)} / 10
                            </Typography>
                        </Stack>
                        <Stack spacing={1.5}>
                            <Typography variant="caption" color="text.secondary">
                                Tỷ lệ hoàn thành bản ghi
                            </Typography>
                            <LinearProgress
                                color="success"
                                variant="determinate"
                                value={records.length ? (completedCount / records.length) * 100 : 0}
                                sx={{ height: 10, borderRadius: 999 }}
                            />
                            <Typography variant="body2" sx={{ fontWeight: 600 }}>
                                {records.length ? Math.round((completedCount / records.length) * 100) : 0}%
                            </Typography>
                        </Stack>
                    </Stack>
                </Paper>
            </Stack>

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
                            <TableRow key={record.record_id} hover>
                                <TableCell>
                                    <Typography variant="body2" sx={{ fontWeight: 600 }}>
                                        {record.lesson.class_name}
                                    </Typography>
                                    <Typography variant="caption" color="text.secondary" display="block">
                                        {record.lesson.class_code} • {format(parseISO(record.lesson.date_start), 'dd/MM/yyyy HH:mm')}
                                    </Typography>
                                    <Typography variant="caption" color="text.secondary" display="block">
                                        {record.lesson.summary.topic || 'Chưa có chủ đề'} • {record.lesson.teacher?.full_name || 'Chưa phân công'}
                                    </Typography>
                                </TableCell>
                                <TableCell>
                                    <Typography variant="body2">
                                        {record.homework_completed ? `${record.homework_score.toFixed(1)}/10` : 'Chưa hoàn thành'}
                                    </Typography>
                                    <Typography variant="caption" color="text.secondary">
                                        {record.lesson.summary.homework || 'Không có bài tập về nhà'}
                                    </Typography>
                                </TableCell>
                                <TableCell>{record.attitude_rating}/5</TableCell>
                                <TableCell>{record.participation_score.toFixed(1)}/10</TableCell>
                                <TableCell>
                                    <Chip
                                        size="small"
                                        label={record.total_score.toFixed(2)}
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
                                        Chưa có kết quả học tập nào được công bố trong phạm vi hiện tại.
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
