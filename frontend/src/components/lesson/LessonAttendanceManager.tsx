import { useEffect, useMemo, useState } from 'react';
import {
    Alert,
    Button,
    Chip,
    MenuItem,
    Paper,
    Select,
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
import { AssignmentTurnedInRounded, SaveRounded } from '@mui/icons-material';
import { toast } from 'sonner';

import {
    useGetLessonAttendanceQuery,
    useUpsertLessonAttendanceMutation,
} from '@/api/lessonApi';
import { getApiErrorMessage } from '@/utils/apiError';

const ATTENDANCE_OPTIONS = [
    { value: 1, label: 'Có mặt' },
    { value: 2, label: 'Vắng' },
    { value: 3, label: 'Có phép' },
    { value: 4, label: 'Muộn' },
    { value: 5, label: 'Về sớm' },
];

type EditableAttendanceRow = {
    student_id: string;
    student_code: string;
    student_name: string;
    status: number;
    note: string;
};

interface LessonAttendanceManagerProps {
    lessonId: string;
    title?: string;
    subtitle?: string;
}

export default function LessonAttendanceManager({
    lessonId,
    title = 'Điểm danh buổi học',
    subtitle = 'Cập nhật trạng thái chuyên cần cho từng học sinh trong buổi học này.',
}: LessonAttendanceManagerProps) {
    const { data, isLoading, isFetching, error, refetch } = useGetLessonAttendanceQuery(lessonId, {
        skip: !lessonId,
    });
    const [upsertAttendance, { isLoading: isSaving }] = useUpsertLessonAttendanceMutation();
    const [rows, setRows] = useState<EditableAttendanceRow[]>([]);

    useEffect(() => {
        const nextRows = (data?.data?.records || []).map((record) => ({
            student_id: record.student.id,
            student_code: record.student.code,
            student_name: record.student.full_name,
            status: record.status,
            note: record.note || '',
        }));
        setRows(nextRows);
    }, [data]);

    const stats = useMemo(() => {
        const counters = {
            total: rows.length,
            present: 0,
            absent: 0,
            excused: 0,
            late: 0,
            early: 0,
            unmarked: 0,
        };
        rows.forEach((row) => {
            switch (row.status) {
                case 1:
                    counters.present += 1;
                    break;
                case 2:
                    counters.absent += 1;
                    break;
                case 3:
                    counters.excused += 1;
                    break;
                case 4:
                    counters.late += 1;
                    break;
                case 5:
                    counters.early += 1;
                    break;
                default:
                    counters.unmarked += 1;
            }
        });
        return counters;
    }, [rows]);

    const handleSave = async () => {
        if (!rows.length) {
            toast.error('Buổi học này chưa có học sinh để điểm danh.');
            return;
        }
        if (rows.some((row) => row.status === 0)) {
            toast.error('Vẫn còn học sinh chưa chọn trạng thái điểm danh.');
            return;
        }

        try {
            await upsertAttendance({
                id: lessonId,
                body: {
                    records: rows.map((row) => ({
                        student_id: row.student_id,
                        status: row.status,
                        note: row.note?.trim() || '',
                    })),
                },
            }).unwrap();
            toast.success('Đã lưu điểm danh buổi học.');
            refetch();
        } catch (err) {
            toast.error(getApiErrorMessage(err, 'Không lưu được điểm danh.'));
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
                            <AssignmentTurnedInRounded color="primary" />
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
                        disabled={isSaving || isLoading || isFetching || !rows.length}
                    >
                        {isSaving ? 'Đang lưu...' : 'Lưu điểm danh'}
                    </Button>
                </Stack>

                <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                    <Chip size="small" label={`${stats.total} học sinh`} variant="outlined" />
                    <Chip size="small" label={`Có mặt ${stats.present}`} color="success" variant="outlined" />
                    <Chip size="small" label={`Vắng ${stats.absent}`} color="error" variant="outlined" />
                    <Chip size="small" label={`Có phép ${stats.excused}`} color="info" variant="outlined" />
                    <Chip size="small" label={`Muộn ${stats.late}`} color="warning" variant="outlined" />
                    <Chip size="small" label={`Về sớm ${stats.early}`} color="secondary" variant="outlined" />
                    {stats.unmarked > 0 ? (
                        <Chip size="small" label={`Chưa chấm ${stats.unmarked}`} color="default" variant="outlined" />
                    ) : null}
                </Stack>

                {error ? (
                    <Alert severity="error">
                        {getApiErrorMessage(error, 'Không tải được dữ liệu điểm danh của buổi học.')}
                    </Alert>
                ) : null}

                <TableContainer>
                    <Table size="small">
                        <TableHead>
                            <TableRow>
                                <TableCell sx={{ fontWeight: 700 }}>Học sinh</TableCell>
                                <TableCell sx={{ fontWeight: 700, width: 180 }}>Trạng thái</TableCell>
                                <TableCell sx={{ fontWeight: 700 }}>Ghi chú</TableCell>
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
                                    <TableCell>
                                        <Select
                                            size="small"
                                            fullWidth
                                            value={String(row.status)}
                                            displayEmpty
                                            onChange={(event) => {
                                                const nextStatus = Number(event.target.value);
                                                setRows((current) =>
                                                    current.map((item, currentIndex) =>
                                                        currentIndex === index ? { ...item, status: nextStatus } : item,
                                                    ),
                                                );
                                            }}
                                        >
                                            <MenuItem value="0">Chọn trạng thái</MenuItem>
                                            {ATTENDANCE_OPTIONS.map((option) => (
                                                <MenuItem key={option.value} value={String(option.value)}>
                                                    {option.label}
                                                </MenuItem>
                                            ))}
                                        </Select>
                                    </TableCell>
                                    <TableCell>
                                        <TextField
                                            size="small"
                                            fullWidth
                                            placeholder="Ghi chú thêm nếu cần"
                                            value={row.note}
                                            onChange={(event) => {
                                                const nextNote = event.target.value;
                                                setRows((current) =>
                                                    current.map((item, currentIndex) =>
                                                        currentIndex === index ? { ...item, note: nextNote } : item,
                                                    ),
                                                );
                                            }}
                                        />
                                    </TableCell>
                                </TableRow>
                            ))}
                            {!rows.length && !isLoading ? (
                                <TableRow>
                                    <TableCell colSpan={3}>
                                        <Typography variant="body2" color="text.secondary">
                                            Buổi học này chưa có dữ liệu roster để điểm danh.
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
