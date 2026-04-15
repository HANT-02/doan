import { useEffect, useMemo, useState } from 'react';
import {
    Alert,
    Button,
    Checkbox,
    Chip,
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
import { SchoolRounded, SaveRounded, VerifiedRounded } from '@mui/icons-material';
import { toast } from 'sonner';

import {
    useFinalizeLessonAcademicRecordsMutation,
    useGetLessonAcademicRecordsQuery,
    useUpsertLessonAcademicRecordsMutation,
} from '@/api/lessonApi';
import { getApiErrorMessage } from '@/utils/apiError';

type EditableRecordRow = {
    student_id: string;
    student_code: string;
    student_name: string;
    homework_completed: boolean;
    homework_score: number;
    attitude_rating: number;
    participation_score: number;
    personal_comment: string;
    total_score: number;
    is_completed: boolean;
};

interface LessonAcademicRecordManagerProps {
    lessonId: string;
    title?: string;
    subtitle?: string;
}

export default function LessonAcademicRecordManager({
    lessonId,
    title = 'Kết quả học tập',
    subtitle = 'Nhập bài tập, thái độ, mức độ tham gia và chốt kết quả theo từng học sinh.',
}: LessonAcademicRecordManagerProps) {
    const { data, isLoading, isFetching, error, refetch } = useGetLessonAcademicRecordsQuery(lessonId, {
        skip: !lessonId,
    });
    const [upsertRecords, { isLoading: isSaving }] = useUpsertLessonAcademicRecordsMutation();
    const [finalizeRecords, { isLoading: isFinalizing }] = useFinalizeLessonAcademicRecordsMutation();
    const [rows, setRows] = useState<EditableRecordRow[]>([]);

    useEffect(() => {
        const nextRows = (data?.data?.records || []).map((item) => ({
            student_id: item.student.id,
            student_code: item.student.code,
            student_name: item.student.full_name,
            homework_completed: item.record?.homework_completed || false,
            homework_score: item.record?.homework_score || 0,
            attitude_rating: item.record?.attitude_rating || 0,
            participation_score: item.record?.participation_score || 0,
            personal_comment: item.record?.personal_comment || '',
            total_score: item.record?.total_score || 0,
            is_completed: item.record?.is_completed || false,
        }));
        setRows(nextRows);
    }, [data]);

    const completionStats = useMemo(() => {
        const total = rows.length;
        const completed = rows.filter((row) => row.is_completed).length;
        return { total, completed };
    }, [rows]);

    const recalcTotalScore = (homeworkScore: number, participationScore: number, attitudeRating: number) => {
        const total = (homeworkScore + participationScore + Math.min(Math.max(attitudeRating, 0), 5) * 2) / 3;
        return Math.round(total * 100) / 100;
    };

    const handleSave = async () => {
        try {
            await upsertRecords({
                id: lessonId,
                body: {
                    records: rows.map((row) => ({
                        student_id: row.student_id,
                        homework_completed: row.homework_completed,
                        homework_score: Number(row.homework_score) || 0,
                        attitude_rating: Number(row.attitude_rating) || 0,
                        participation_score: Number(row.participation_score) || 0,
                        personal_comment: row.personal_comment || '',
                    })),
                },
            }).unwrap();
            toast.success('Đã lưu kết quả học tập của buổi học.');
            refetch();
        } catch (err) {
            toast.error(getApiErrorMessage(err, 'Không lưu được kết quả học tập.'));
        }
    };

    const handleFinalize = async () => {
        try {
            await finalizeRecords(lessonId).unwrap();
            toast.success('Đã chốt kết quả học tập của buổi học.');
            refetch();
        } catch (err) {
            toast.error(getApiErrorMessage(err, 'Không chốt được kết quả học tập.'));
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
                            <SchoolRounded color="primary" />
                            <Typography variant="h6" sx={{ fontWeight: 700 }}>
                                {title}
                            </Typography>
                        </Stack>
                        <Typography variant="body2" color="text.secondary">
                            {subtitle}
                        </Typography>
                    </Stack>

                    <Stack direction="row" spacing={1}>
                        <Button
                            variant="outlined"
                            startIcon={<SaveRounded />}
                            onClick={handleSave}
                            disabled={isSaving || isLoading || isFetching || !rows.length}
                        >
                            {isSaving ? 'Đang lưu...' : 'Lưu kết quả'}
                        </Button>
                        <Button
                            variant="contained"
                            startIcon={<VerifiedRounded />}
                            onClick={handleFinalize}
                            disabled={isFinalizing || isLoading || isFetching || !rows.length}
                        >
                            {isFinalizing ? 'Đang chốt...' : 'Chốt kết quả'}
                        </Button>
                    </Stack>
                </Stack>

                <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                    <Chip size="small" label={`${completionStats.total} học sinh`} variant="outlined" />
                    <Chip size="small" label={`Đã chốt ${completionStats.completed}`} color="success" variant="outlined" />
                </Stack>

                {error ? (
                    <Alert severity="error">
                        {getApiErrorMessage(error, 'Không tải được dữ liệu kết quả học tập.')}
                    </Alert>
                ) : null}

                <TableContainer>
                    <Table size="small">
                        <TableHead>
                            <TableRow>
                                <TableCell sx={{ fontWeight: 700 }}>Học sinh</TableCell>
                                <TableCell sx={{ fontWeight: 700, width: 100 }}>Hoàn thành BT</TableCell>
                                <TableCell sx={{ fontWeight: 700, width: 110 }}>Điểm BT</TableCell>
                                <TableCell sx={{ fontWeight: 700, width: 110 }}>Thái độ</TableCell>
                                <TableCell sx={{ fontWeight: 700, width: 120 }}>Tham gia</TableCell>
                                <TableCell sx={{ fontWeight: 700 }}>Nhận xét</TableCell>
                                <TableCell sx={{ fontWeight: 700, width: 100 }}>Tổng</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {rows.map((row, index) => (
                                <TableRow key={row.student_id} hover>
                                    <TableCell>
                                        <Typography variant="body2" sx={{ fontWeight: 600 }}>
                                            {row.student_name}
                                        </Typography>
                                        <Typography variant="caption" color="text.secondary">
                                            {row.student_code}
                                        </Typography>
                                    </TableCell>
                                    <TableCell padding="checkbox">
                                        <Checkbox
                                            checked={row.homework_completed}
                                            onChange={(event) => {
                                                const checked = event.target.checked;
                                                setRows((current) =>
                                                    current.map((item, currentIndex) =>
                                                        currentIndex === index ? { ...item, homework_completed: checked } : item,
                                                    ),
                                                );
                                            }}
                                        />
                                    </TableCell>
                                    <TableCell>
                                        <TextField
                                            size="small"
                                            type="number"
                                            value={row.homework_score}
                                            inputProps={{ min: 0, max: 10, step: 0.1 }}
                                            onChange={(event) => {
                                                const nextValue = Number(event.target.value);
                                                setRows((current) =>
                                                    current.map((item, currentIndex) => {
                                                        if (currentIndex !== index) return item;
                                                        return {
                                                            ...item,
                                                            homework_score: nextValue,
                                                            total_score: recalcTotalScore(nextValue, item.participation_score, item.attitude_rating),
                                                        };
                                                    }),
                                                );
                                            }}
                                        />
                                    </TableCell>
                                    <TableCell>
                                        <TextField
                                            size="small"
                                            type="number"
                                            value={row.attitude_rating}
                                            inputProps={{ min: 0, max: 5, step: 1 }}
                                            onChange={(event) => {
                                                const nextValue = Number(event.target.value);
                                                setRows((current) =>
                                                    current.map((item, currentIndex) => {
                                                        if (currentIndex !== index) return item;
                                                        return {
                                                            ...item,
                                                            attitude_rating: nextValue,
                                                            total_score: recalcTotalScore(item.homework_score, item.participation_score, nextValue),
                                                        };
                                                    }),
                                                );
                                            }}
                                        />
                                    </TableCell>
                                    <TableCell>
                                        <TextField
                                            size="small"
                                            type="number"
                                            value={row.participation_score}
                                            inputProps={{ min: 0, max: 10, step: 0.1 }}
                                            onChange={(event) => {
                                                const nextValue = Number(event.target.value);
                                                setRows((current) =>
                                                    current.map((item, currentIndex) => {
                                                        if (currentIndex !== index) return item;
                                                        return {
                                                            ...item,
                                                            participation_score: nextValue,
                                                            total_score: recalcTotalScore(item.homework_score, nextValue, item.attitude_rating),
                                                        };
                                                    }),
                                                );
                                            }}
                                        />
                                    </TableCell>
                                    <TableCell>
                                        <TextField
                                            size="small"
                                            fullWidth
                                            value={row.personal_comment}
                                            onChange={(event) => {
                                                const nextValue = event.target.value;
                                                setRows((current) =>
                                                    current.map((item, currentIndex) =>
                                                        currentIndex === index ? { ...item, personal_comment: nextValue } : item,
                                                    ),
                                                );
                                            }}
                                        />
                                    </TableCell>
                                    <TableCell>
                                        <Stack spacing={0.25}>
                                            <Typography variant="body2" sx={{ fontWeight: 700 }}>
                                                {row.total_score.toFixed(2)}
                                            </Typography>
                                            {row.is_completed ? (
                                                <Chip size="small" label="Đã chốt" color="success" variant="outlined" />
                                            ) : (
                                                <Chip size="small" label="Tạm lưu" variant="outlined" />
                                            )}
                                        </Stack>
                                    </TableCell>
                                </TableRow>
                            ))}
                            {!rows.length && !isLoading ? (
                                <TableRow>
                                    <TableCell colSpan={7}>
                                        <Typography variant="body2" color="text.secondary">
                                            Buổi học này chưa có học sinh để nhập kết quả.
                                        </Typography>
                                    </TableCell>
                                </TableRow>
                            ) : null}
                        </TableBody>
                    </Table>
                </TableContainer>
            </Stack>
        </Paper>
    );
}
