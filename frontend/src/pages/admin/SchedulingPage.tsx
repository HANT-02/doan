import { useMemo, useState } from 'react';
import { Controller, useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import {
    Alert,
    Autocomplete,
    Box,
    Button,
    Chip,
    InputAdornment,
    Paper,
    Skeleton,
    Stack,
    Step,
    StepLabel,
    Stepper,
    TextField,
    Typography,
} from '@mui/material';
import {
    PlayArrowRounded,
    RefreshRounded,
    RuleRounded,
    SearchRounded,
} from '@mui/icons-material';
import { DataGrid } from '@mui/x-data-grid';
import type { GridColDef, GridRenderCellParams } from '@mui/x-data-grid';
import { toast } from 'sonner';

import { useGetClassesQuery } from '@/api/classApi';
import { useGetRoomsQuery } from '@/api/roomApi';
import {
    useCommitSchedulingPreviewMutation,
    useLazyGetLatestSchedulingPreviewQuery,
    useLazyGetSchedulingPreviewQuery,
    usePreviewSchedulingMutation,
    type SchedulingAssignment,
    type SchedulingPreview,
} from '@/api/schedulingApi';
import { useGetTeachersQuery } from '@/api/teacherApi';
import PageHeader from '@/components/common/PageHeader';
import { getApiErrorMessage } from '@/utils/apiError';

const schedulingSchema = z.object({
    date_from: z.string().min(1, 'Vui lòng chọn ngày bắt đầu'),
    date_to: z.string().min(1, 'Vui lòng chọn ngày kết thúc'),
    class_ids: z.array(z.string()).optional(),
    teacher_ids: z.array(z.string()).optional(),
    room_ids: z.array(z.string()).optional(),
}).refine((values) => values.date_to >= values.date_from, {
    message: 'Ngày kết thúc phải lớn hơn hoặc bằng ngày bắt đầu',
    path: ['date_to'],
});

type SchedulingFormValues = z.infer<typeof schedulingSchema>;

const previewSteps = ['Cấu hình đầu vào', 'Chạy CSP preview', 'Rà soát kết quả'];
const formatDateInput = (value: Date) => {
    const year = value.getFullYear();
    const month = `${value.getMonth() + 1}`.padStart(2, '0');
    const day = `${value.getDate()}`.padStart(2, '0');
    return `${year}-${month}-${day}`;
};

const defaultDateFrom = formatDateInput(new Date());
const defaultDateTo = formatDateInput(new Date(Date.now() + 7 * 24 * 60 * 60 * 1000));

const formatDateTime = (value: string) => new Date(value).toLocaleString('vi-VN');

const getConflictSeverity = (type: string): 'error' | 'warning' | 'info' => {
    switch (type) {
        case 'NO_CLASS_INPUT':
        case 'NO_ACTIVE_ROOM':
        case 'NO_VALID_DATE_RANGE':
            return 'error';
        case 'MISSING_TEACHER':
        case 'PREFERRED_ROOM_UNAVAILABLE':
        case 'ROOM_CAPACITY_BLOCK':
        case 'NO_SLOT_IN_RANGE':
        case 'NO_DOMAIN':
        case 'MISSING_COURSE':
        case 'INVALID_COURSE_SESSION_COUNT':
        case 'INVALID_COURSE_DURATION':
        case 'CLASS_SCHEDULE_NO_SLOT':
        case 'CLASS_SCHEDULE_ROOM_UNAVAILABLE':
            return 'warning';
        default:
            return 'info';
    }
};

const getConflictActionHint = (type: string) => {
    switch (type) {
        case 'MISSING_TEACHER':
            return 'Gợi ý: vào Quản lý lớp học để gán giáo viên phụ trách rồi chạy lại preview.';
        case 'NO_CLASS_INPUT':
            return 'Gợi ý: kiểm tra bộ lọc lớp, trạng thái OPEN hoặc giáo viên đã chọn.';
        case 'NO_ACTIVE_ROOM':
            return 'Gợi ý: bỏ lọc phòng hoặc bổ sung phòng học khả dụng trước khi chạy lại.';
        case 'PREFERRED_ROOM_UNAVAILABLE':
            return 'Gợi ý: đổi phòng đang gán cho lớp hoặc nới lại bộ lọc phòng.';
        case 'ROOM_CAPACITY_BLOCK':
            return 'Gợi ý: chọn phòng sức chứa lớn hơn hoặc điều chỉnh sĩ số tối đa của lớp.';
        case 'NO_SLOT_IN_RANGE':
            return 'Gợi ý: nới khoảng ngày preview để tạo thêm khung giờ khả dụng.';
        case 'NO_DOMAIN':
            return 'Gợi ý: giảm bớt bộ lọc và kiểm tra lại dữ liệu lớp/phòng để solver có thêm miền giá trị.';
        case 'MISSING_COURSE':
            return 'Gợi ý: gắn khóa học cho lớp trước khi chạy preview để hệ thống sinh đúng số buổi học.';
        case 'INVALID_COURSE_SESSION_COUNT':
            return 'Gợi ý: cập nhật `session_count` của khóa học thành số buổi hợp lệ.';
        case 'INVALID_COURSE_DURATION':
            return 'Gợi ý: cập nhật `session_duration_minutes` của khóa học để preview dùng đúng thời lượng buổi học.';
        case 'CLASS_SCHEDULE_NO_SLOT':
            return 'Gợi ý: kiểm tra lịch mẫu của lớp và bảo đảm khung giờ mẫu đủ chứa thời lượng buổi học từ khóa học.';
        case 'CLASS_SCHEDULE_ROOM_UNAVAILABLE':
            return 'Gợi ý: kiểm tra `room_id` trong lịch mẫu của lớp hoặc nới lại bộ lọc phòng để giữ đúng phòng cố định theo lịch mẫu.';
        case 'NO_VALID_DATE_RANGE':
            return 'Gợi ý: chọn ngày kết thúc lớn hơn hoặc bằng ngày bắt đầu.';
        default:
            return 'Gợi ý: rà lại dữ liệu lớp, giáo viên, phòng và bộ lọc preview.';
    }
};

const getConflictScopeLabel = (classCode?: string, className?: string, sessionIndex?: number, sessionTotal?: number) => {
    const sessionLabel = sessionIndex && sessionTotal ? ` • Buổi ${sessionIndex}/${sessionTotal}` : '';
    if (classCode) {
        return className ? `${classCode} - ${className}${sessionLabel}` : `${classCode}${sessionLabel}`;
    }

    return className ? `${className}${sessionLabel}` : 'Tổng quan preview';
};

type ClassPreviewProgress = {
    classId: string;
    classCode: string;
    className: string;
    scheduledSessions: number;
    totalSessions: number;
    conflictCount: number;
};

export const SchedulingPage = () => {
    const [preview, setPreview] = useState<SchedulingPreview | null>(null);
    const [search, setSearch] = useState('');

    const { data: classesData, isLoading: isLoadingClasses } = useGetClassesQuery({ page: 1, limit: 200, status: 'OPEN' });
    const { data: teachersData, isLoading: isLoadingTeachers } = useGetTeachersQuery({ page: 1, limit: 200, status: 'ACTIVE' });
    const { data: roomsData, isLoading: isLoadingRooms } = useGetRoomsQuery({ page: 1, limit: 200 });

    const [previewScheduling, { isLoading: isPreviewing }] = usePreviewSchedulingMutation();
    const [commitSchedulingPreview, { isLoading: isCommitting }] = useCommitSchedulingPreviewMutation();
    const [loadPreview, { isFetching: isLoadingPreview }] = useLazyGetSchedulingPreviewQuery();
    const [loadLatestPreview, { isFetching: isLoadingLatest }] = useLazyGetLatestSchedulingPreviewQuery();

    const classes = classesData?.data?.classes || [];
    const teachers = teachersData?.data?.teachers || [];
    const rooms = roomsData?.data?.rooms || [];

    const {
        control,
        handleSubmit,
        formState: { errors },
    } = useForm<SchedulingFormValues>({
        resolver: zodResolver(schedulingSchema),
        defaultValues: {
            date_from: defaultDateFrom,
            date_to: defaultDateTo,
            class_ids: [],
            teacher_ids: [],
            room_ids: [],
        },
    });

    const filteredAssignments = useMemo(() => {
        const rows = preview?.assignments || [];
        const keyword = search.trim().toLowerCase();
        if (!keyword) {
            return rows;
        }

        return rows.filter((row) =>
            [row.class_code, row.class_name, row.teacher_label, row.room_name, `buoi ${row.session_index}`]
                .filter(Boolean)
                .some((value) => value.toLowerCase().includes(keyword)),
        );
    }, [preview?.assignments, search]);

    const conflictSeveritySummary = useMemo(() => {
        const summary = {
            error: 0,
            warning: 0,
            info: 0,
        };

        for (const conflict of preview?.conflicts || []) {
            summary[getConflictSeverity(conflict.type)] += 1;
        }

        return summary;
    }, [preview?.conflicts]);

    const classProgressRows = useMemo<ClassPreviewProgress[]>(() => {
        if (!preview) {
            return [];
        }

        const progressMap = new Map<string, ClassPreviewProgress>();

        for (const assignment of preview.assignments) {
            const current = progressMap.get(assignment.class_id);
            if (current) {
                current.scheduledSessions += 1;
                current.totalSessions = Math.max(current.totalSessions, assignment.session_total);
                continue;
            }

            progressMap.set(assignment.class_id, {
                classId: assignment.class_id,
                classCode: assignment.class_code,
                className: assignment.class_name,
                scheduledSessions: 1,
                totalSessions: assignment.session_total,
                conflictCount: 0,
            });
        }

        for (const conflict of preview.conflicts) {
            if (!conflict.class_id) {
                continue;
            }

            const current = progressMap.get(conflict.class_id);
            if (current) {
                current.conflictCount += 1;
                current.totalSessions = Math.max(current.totalSessions, conflict.session_total || 0);
                continue;
            }

            progressMap.set(conflict.class_id, {
                classId: conflict.class_id,
                classCode: conflict.class_code,
                className: conflict.class_name,
                scheduledSessions: 0,
                totalSessions: conflict.session_total || 0,
                conflictCount: 1,
            });
        }

        return Array.from(progressMap.values()).sort((left, right) => left.classCode.localeCompare(right.classCode));
    }, [preview]);

    const activeStep = preview ? 2 : 0;
    const canCommitPreview = !!preview && preview.status === 'COMPLETED' && preview.summary.unscheduled_lessons === 0;

    const handleRunPreview = handleSubmit(async (values) => {
        try {
            const response = await previewScheduling(values).unwrap();
            setPreview(response.data);
            toast.success('Đã chạy preview xếp lịch');
        } catch (error) {
            toast.error(getApiErrorMessage(error, 'Không thể chạy preview xếp lịch'));
        }
    });

    const handleLoadLatest = async () => {
        try {
            const response = await loadLatestPreview().unwrap();
            setPreview(response.data);
            toast.success('Đã tải preview mới nhất');
        } catch (error) {
            toast.error(getApiErrorMessage(error, 'Chưa có preview nào để tải'));
        }
    };

    const handleRefreshCurrent = async () => {
        if (!preview?.run_id) {
            await handleLoadLatest();
            return;
        }

        try {
            const response = await loadPreview(preview.run_id).unwrap();
            setPreview(response.data);
            toast.success('Đã làm mới kết quả preview');
        } catch (error) {
            toast.error(getApiErrorMessage(error, 'Không thể tải lại preview hiện tại'));
        }
    };

    const handleCommitPreview = async () => {
        if (!preview?.run_id) {
            toast.error('Chưa có preview nào để commit');
            return;
        }

        try {
            const response = await commitSchedulingPreview({ run_id: preview.run_id }).unwrap();
            toast.success(response.data.message || `Đã commit ${response.data.scheduled_lessons} buổi học`);
        } catch (error) {
            toast.error(getApiErrorMessage(error, 'Không thể commit preview xếp lịch'));
        }
    };

    const columns: GridColDef<SchedulingAssignment>[] = [
        {
            field: 'class_name',
            headerName: 'Lớp học',
            minWidth: 220,
            flex: 1.2,
            renderCell: (params: GridRenderCellParams<SchedulingAssignment>) => (
                <Box>
                    <Typography variant="body2" sx={{ fontWeight: 700 }}>
                        {params.row.class_name}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                        {params.row.class_code} • Buổi {params.row.session_index}/{params.row.session_total}
                    </Typography>
                </Box>
            ),
        },
        {
            field: 'teacher_label',
            headerName: 'Giáo viên',
            minWidth: 180,
            flex: 1,
        },
        {
            field: 'room_name',
            headerName: 'Phòng học',
            minWidth: 160,
            flex: 0.9,
            renderCell: (params: GridRenderCellParams<SchedulingAssignment>) => (
                <Box>
                    <Typography variant="body2">{params.row.room_name}</Typography>
                    <Typography variant="caption" color="text.secondary">
                        Sức chứa {params.row.room_capacity}
                    </Typography>
                </Box>
            ),
        },
        {
            field: 'start_time',
            headerName: 'Thời gian',
            minWidth: 220,
            flex: 1,
            renderCell: (params: GridRenderCellParams<SchedulingAssignment>) => (
                <Box>
                    <Typography variant="body2">{formatDateTime(params.row.start_time)}</Typography>
                    <Typography variant="caption" color="text.secondary">
                        Kết thúc: {formatDateTime(params.row.end_time)}
                    </Typography>
                </Box>
            ),
        },
        {
            field: 'constraint_fit',
            headerName: 'Kết quả',
            width: 160,
            renderCell: (params: GridRenderCellParams<SchedulingAssignment>) => (
                <Chip
                    size="small"
                    color={params.row.constraint_fit === 'HARD_OK' ? 'success' : 'warning'}
                    label={params.row.constraint_fit === 'HARD_OK' ? 'Đạt hard constraints' : 'Partial preview'}
                />
            ),
        },
    ];

    const renderPreviewSkeleton = () => (
        <Stack spacing={1.5}>
            <Skeleton variant="rounded" height={80} />
            <Skeleton variant="rounded" height={360} />
        </Stack>
    );

    return (
        <Box>
            <PageHeader
                title="Xếp lịch (CSP)"
                subtitle="Trigger thuật toán CSP, xem preview xếp lịch và rà soát conflict trước khi commit."
                breadcrumbs={[
                    { label: 'Tổng quan', path: '/app/admin/overview' },
                    { label: 'Xếp lịch (CSP)' },
                ]}
                actions={
                    <Stack direction="row" spacing={1}>
                        <Button
                            variant="outlined"
                            startIcon={<RefreshRounded />}
                            onClick={() => void handleLoadLatest()}
                            disabled={isLoadingLatest}
                        >
                            Tải preview mới nhất
                        </Button>
                        <Button
                            variant="contained"
                            startIcon={<PlayArrowRounded />}
                            onClick={() => void handleRunPreview()}
                            disabled={isPreviewing}
                        >
                            {isPreviewing ? 'Đang chạy...' : 'Chạy xếp lịch'}
                        </Button>
                    </Stack>
                }
            />

            <Stack spacing={3}>
                <Paper variant="outlined" sx={{ p: 3, borderRadius: 4 }}>
                    <Stepper activeStep={activeStep} alternativeLabel>
                        {previewSteps.map((label) => (
                            <Step key={label}>
                                <StepLabel>{label}</StepLabel>
                            </Step>
                        ))}
                    </Stepper>
                </Paper>

                <Paper variant="outlined" sx={{ p: 3, borderRadius: 4 }}>
                    <Stack spacing={2.5}>
                        <Box>
                            <Typography variant="h6" sx={{ fontWeight: 700 }}>
                                Bộ lọc chạy preview
                            </Typography>
                            <Typography variant="body2" color="text.secondary">
                                Chọn khoảng ngày và bộ lọc đầu vào tối thiểu để test hard constraints nền.
                            </Typography>
                        </Box>

                        <Stack direction={{ xs: 'column', md: 'row' }} spacing={2}>
                            <Controller
                                control={control}
                                name="date_from"
                                render={({ field }) => (
                                    <TextField
                                        {...field}
                                        type="date"
                                        label="Từ ngày"
                                        fullWidth
                                        InputLabelProps={{ shrink: true }}
                                        error={!!errors.date_from}
                                        helperText={errors.date_from?.message}
                                    />
                                )}
                            />
                            <Controller
                                control={control}
                                name="date_to"
                                render={({ field }) => (
                                    <TextField
                                        {...field}
                                        type="date"
                                        label="Đến ngày"
                                        fullWidth
                                        InputLabelProps={{ shrink: true }}
                                        error={!!errors.date_to}
                                        helperText={errors.date_to?.message}
                                    />
                                )}
                            />
                        </Stack>

                        <Controller
                            control={control}
                            name="class_ids"
                            render={({ field }) => (
                                <Autocomplete
                                    multiple
                                    options={classes}
                                    loading={isLoadingClasses}
                                    value={classes.filter((item) => field.value?.includes(item.id))}
                                    onChange={(_, values) => field.onChange(values.map((item) => item.id))}
                                    getOptionLabel={(option) => `${option.name} (${option.code})`}
                                    renderInput={(params) => (
                                        <TextField
                                            {...params}
                                            label="Lọc theo lớp"
                                            helperText="Để trống nếu muốn chạy tất cả lớp đang mở"
                                        />
                                    )}
                                />
                            )}
                        />

                        <Stack direction={{ xs: 'column', lg: 'row' }} spacing={2}>
                            <Controller
                                control={control}
                                name="teacher_ids"
                                render={({ field }) => (
                                    <Autocomplete
                                        multiple
                                        options={teachers}
                                        loading={isLoadingTeachers}
                                        value={teachers.filter((item) => field.value?.includes(item.id))}
                                        onChange={(_, values) => field.onChange(values.map((item) => item.id))}
                                        getOptionLabel={(option) => `${option.full_name} (${option.code})`}
                                        renderInput={(params) => (
                                            <TextField
                                                {...params}
                                                label="Lọc theo giáo viên"
                                                helperText="Chỉ lấy lớp thuộc các giáo viên đã chọn"
                                            />
                                        )}
                                    />
                                )}
                            />
                            <Controller
                                control={control}
                                name="room_ids"
                                render={({ field }) => (
                                    <Autocomplete
                                        multiple
                                        options={rooms}
                                        loading={isLoadingRooms}
                                        value={rooms.filter((item) => field.value?.includes(item.id))}
                                        onChange={(_, values) => field.onChange(values.map((item) => item.id))}
                                        getOptionLabel={(option) => `${option.name} (Sức chứa ${option.capacity})`}
                                        renderInput={(params) => (
                                            <TextField
                                                {...params}
                                                label="Lọc theo phòng"
                                                helperText="Dùng để test capacity và room conflict"
                                            />
                                        )}
                                    />
                                )}
                            />
                        </Stack>
                    </Stack>
                </Paper>

                {isPreviewing || isLoadingPreview || isLoadingLatest ? renderPreviewSkeleton() : null}

                {!isPreviewing && !isLoadingPreview && !isLoadingLatest && !preview ? (
                    <Paper
                        variant="outlined"
                        sx={{ p: 6, borderRadius: 4, borderStyle: 'dashed', textAlign: 'center' }}
                    >
                        <RuleRounded sx={{ fontSize: 48, color: 'text.secondary', mb: 2 }} />
                        <Typography variant="h6" sx={{ fontWeight: 700 }}>
                            Chưa có preview nào
                        </Typography>
                        <Typography variant="body2" color="text.secondary" sx={{ mt: 1, mb: 3 }}>
                            Bấm `Chạy xếp lịch` để tạo preview đầu tiên hoặc tải preview gần nhất từ backend.
                        </Typography>
                        <Stack direction="row" spacing={1} justifyContent="center">
                            <Button variant="contained" startIcon={<PlayArrowRounded />} onClick={() => void handleRunPreview()}>
                                Chạy preview
                            </Button>
                            <Button variant="outlined" startIcon={<RefreshRounded />} onClick={() => void handleLoadLatest()}>
                                Tải preview mới nhất
                            </Button>
                        </Stack>
                    </Paper>
                ) : null}

                {preview ? (
                    <Stack spacing={2.5}>
                        <Paper
                            variant="outlined"
                            sx={{
                                p: 3,
                                borderRadius: 4,
                                background: 'linear-gradient(135deg, rgba(15,118,110,0.08), rgba(13,148,136,0.02))',
                            }}
                        >
                            <Stack direction={{ xs: 'column', md: 'row' }} spacing={2} justifyContent="space-between">
                                <Box>
                                    <Typography variant="h6" sx={{ fontWeight: 800 }}>
                                        Preview run: {preview.run_id}
                                    </Typography>
                                    <Typography variant="body2" color="text.secondary">
                                        Sinh lúc {formatDateTime(preview.generated_at)} • Trạng thái {preview.status}
                                    </Typography>
                                </Box>
                                <Stack direction="row" spacing={1} flexWrap="wrap">
                                    <Chip label={`${preview.summary.requested_classes} lớp`} color="default" variant="outlined" />
                                    <Chip label={`${preview.summary.requested_sessions} buổi cần xếp`} color="default" variant="outlined" />
                                    <Chip label={`${preview.summary.scheduled_lessons} buổi đã xếp`} color="success" variant="outlined" />
                                    <Chip label={`${preview.summary.unscheduled_lessons} buổi chưa xếp`} color="warning" variant="outlined" />
                                    <Chip label={`${preview.summary.conflict_count} conflict`} color="error" variant="outlined" />
                                    <Chip label={`Soft score ${preview.summary.soft_score}`} color="primary" variant="outlined" />
                                </Stack>
                            </Stack>
                        </Paper>

                        {preview.conflicts.length > 0 ? (
                            <Stack spacing={1.25}>
                                <Alert severity={preview.assignments.length > 0 ? 'warning' : 'error'}>
                                    Preview ghi nhận {preview.conflicts.length} vấn đề cần xử lý. Mỗi conflict bên dưới đều kèm gợi ý để bạn biết nên sửa dữ liệu nào trước khi chạy lại.
                                </Alert>
                                <Stack direction="row" spacing={1} flexWrap="wrap">
                                    <Chip label={`${conflictSeveritySummary.error} lỗi chặn`} color="error" variant="outlined" />
                                    <Chip label={`${conflictSeveritySummary.warning} cảnh báo`} color="warning" variant="outlined" />
                                    <Chip label={`${conflictSeveritySummary.info} thông tin`} color="info" variant="outlined" />
                                </Stack>
                            </Stack>
                        ) : (
                            <Alert severity="success">Preview hiện không có conflict hard constraint.</Alert>
                        )}

                        <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 4 }}>
                            <Typography variant="h6" sx={{ fontWeight: 700, mb: 1.5 }}>
                                Tiến độ preview theo lớp
                            </Typography>
                            <Stack spacing={1.25} sx={{ mb: 2.5 }}>
                                {classProgressRows.map((row) => (
                                    <Paper
                                        key={row.classId}
                                        variant="outlined"
                                        sx={{
                                            p: 1.5,
                                            borderRadius: 3,
                                            display: 'flex',
                                            flexDirection: { xs: 'column', md: 'row' },
                                            justifyContent: 'space-between',
                                            gap: 1,
                                        }}
                                    >
                                        <Box>
                                            <Typography variant="body2" sx={{ fontWeight: 700 }}>
                                                {row.classCode} - {row.className}
                                            </Typography>
                                            <Typography variant="caption" color="text.secondary">
                                                Đã xếp {row.scheduledSessions}/{row.totalSessions || '?'} buổi
                                            </Typography>
                                        </Box>
                                        <Stack direction="row" spacing={1} flexWrap="wrap">
                                            <Chip
                                                size="small"
                                                color={row.totalSessions > 0 && row.scheduledSessions === row.totalSessions ? 'success' : 'warning'}
                                                variant="outlined"
                                                label={`${row.scheduledSessions}/${row.totalSessions || '?'} buổi`}
                                            />
                                            <Chip
                                                size="small"
                                                color={row.conflictCount > 0 ? 'warning' : 'success'}
                                                variant="outlined"
                                                label={row.conflictCount > 0 ? `${row.conflictCount} conflict` : 'Không conflict'}
                                            />
                                        </Stack>
                                    </Paper>
                                ))}
                            </Stack>

                            <Stack
                                direction={{ xs: 'column', lg: 'row' }}
                                spacing={1.5}
                                justifyContent="space-between"
                                alignItems={{ xs: 'stretch', lg: 'center' }}
                                sx={{ mb: 2 }}
                            >
                                <TextField
                                    size="small"
                                    value={search}
                                    onChange={(event) => setSearch(event.target.value)}
                                    placeholder="Tìm lớp, giáo viên, phòng..."
                                    sx={{ minWidth: { xs: '100%', md: 320 } }}
                                    InputProps={{
                                        startAdornment: (
                                            <InputAdornment position="start">
                                                <SearchRounded fontSize="small" />
                                            </InputAdornment>
                                        ),
                                    }}
                                />
                                <Stack direction="row" spacing={1}>
                                    <Button
                                        variant="outlined"
                                        startIcon={<RefreshRounded />}
                                        onClick={() => void handleRefreshCurrent()}
                                        disabled={isLoadingPreview || isLoadingLatest}
                                    >
                                        Làm mới
                                    </Button>
                                    <Button
                                        variant="outlined"
                                        color={canCommitPreview ? 'primary' : 'inherit'}
                                        onClick={() => void handleCommitPreview()}
                                        disabled={!canCommitPreview || isCommitting}
                                    >
                                        {isCommitting ? 'Đang commit...' : 'Commit xuống lesson'}
                                    </Button>
                                </Stack>
                            </Stack>

                            {canCommitPreview ? (
                                <Alert severity="success" sx={{ mb: 2 }}>
                                    Preview đã đạt trạng thái `COMPLETED`. Bạn có thể commit để tạo `lesson` thật; hệ thống sẽ chặn nếu phát hiện trùng lịch với lesson đang tồn tại.
                                </Alert>
                            ) : (
                                <Alert severity="info" sx={{ mb: 2 }}>
                                    Commit chỉ mở khi preview `COMPLETED` và không còn buổi chưa xếp. Hãy xử lý conflict rồi chạy lại preview trước khi ghi lịch thật xuống `lesson`.
                                </Alert>
                            )}

                            <DataGrid
                                rows={filteredAssignments}
                                columns={columns}
                                autoHeight
                                disableRowSelectionOnClick
                                getRowId={(row) => row.variable_id}
                                pageSizeOptions={[5, 10, 25]}
                                initialState={{
                                    pagination: {
                                        paginationModel: {
                                            page: 0,
                                            pageSize: 10,
                                        },
                                    },
                                }}
                                localeText={{ noRowsLabel: 'Không có assignment nào trong preview hiện tại' }}
                                sx={{
                                    border: 'none',
                                    '& .MuiDataGrid-columnHeaders': {
                                        backgroundColor: '#f8fafc',
                                    },
                                }}
                            />
                        </Paper>

                        {preview.conflicts.length > 0 ? (
                            <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 4 }}>
                                <Typography variant="h6" sx={{ fontWeight: 700, mb: 1.5 }}>
                                    Danh sách conflict
                                </Typography>
                                <Stack spacing={1.25}>
                                    {preview.conflicts.map((conflict) => (
                                        <Alert
                                            key={`${conflict.variable_id || 'global'}-${conflict.type}`}
                                            severity={getConflictSeverity(conflict.type)}
                                        >
                                            <strong>{getConflictScopeLabel(conflict.class_code, conflict.class_name, conflict.session_index, conflict.session_total)}</strong>
                                            {`: ${conflict.message} `}
                                            <Typography component="span" variant="body2" sx={{ fontWeight: 600 }}>
                                                {getConflictActionHint(conflict.type)}
                                            </Typography>
                                        </Alert>
                                    ))}
                                </Stack>
                            </Paper>
                        ) : null}
                    </Stack>
                ) : null}
            </Stack>
        </Box>
    );
};
