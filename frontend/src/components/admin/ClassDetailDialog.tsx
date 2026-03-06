import { useEffect, useMemo, useState } from 'react';
import { Controller, useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import {
    Alert,
    Autocomplete,
    Avatar,
    Box,
    Button,
    Chip,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    Divider,
    Grid,
    Paper,
    Skeleton,
    Stack,
    Tab,
    Tabs,
    TextField,
    Typography,
} from '@mui/material';
import {
    AddRounded,
    DeleteOutlineRounded,
    EditOutlined,
    InfoOutlined,
    PersonOutline,
    RefreshRounded,
    SchoolOutlined,
} from '@mui/icons-material';
import { DataGrid } from '@mui/x-data-grid';
import type { GridColDef, GridRenderCellParams } from '@mui/x-data-grid';
import { toast } from 'sonner';

import {
    useAssignTeacherMutation,
    useEnrollStudentsMutation,
    useGetClassByIdQuery,
    useGetClassRosterQuery,
    useRemoveStudentsMutation,
    type Class,
} from '@/api/classApi';
import { useGetStudentsQuery, type Student } from '@/api/studentApi';
import { useGetTeachersQuery, type Teacher } from '@/api/teacherApi';
import ConfirmDialog from '@/components/common/ConfirmDialog';

const enrollStudentsSchema = z.object({
    student_ids: z.array(z.string()).min(1, 'Vui lòng chọn ít nhất một học sinh'),
});

const assignTeacherSchema = z.object({
    teacher_id: z.string().min(1, 'Vui lòng chọn giáo viên phụ trách'),
});

type EnrollStudentsValues = z.infer<typeof enrollStudentsSchema>;
type AssignTeacherValues = z.infer<typeof assignTeacherSchema>;

interface ClassDetailDialogProps {
    open: boolean;
    classId: string | null;
    onClose: () => void;
    onEdit: (classData: Class) => void;
}

const getErrorMessage = (error: unknown, fallback: string) => {
    if (typeof error === 'object' && error && 'data' in error) {
        const apiError = error as { data?: { message?: string } };
        return apiError.data?.message || fallback;
    }

    if (error instanceof Error) {
        return error.message;
    }

    return fallback;
};

const formatDate = (value?: string) => {
    if (!value) {
        return 'Chưa cập nhật';
    }

    const date = new Date(value);
    if (Number.isNaN(date.getTime())) {
        return value;
    }

    return date.toLocaleDateString('vi-VN');
};

const statusChip = (status?: string) => {
    switch (status) {
        case 'OPEN':
            return { label: 'Đang mở', color: 'success' as const };
        case 'CLOSED':
            return { label: 'Đã đóng', color: 'default' as const };
        case 'CANCELLED':
            return { label: 'Đã hủy', color: 'error' as const };
        default:
            return { label: status || 'Không rõ', color: 'default' as const };
    }
};

export default function ClassDetailDialog({
    open,
    classId,
    onClose,
    onEdit,
}: ClassDetailDialogProps) {
    const [tab, setTab] = useState(0);
    const [isEnrollDialogOpen, setIsEnrollDialogOpen] = useState(false);
    const [rosterSearch, setRosterSearch] = useState('');
    const [studentSearch, setStudentSearch] = useState('');
    const [studentToRemove, setStudentToRemove] = useState<Student | null>(null);

    const {
        data: classResponse,
        isLoading: isLoadingClass,
        isFetching: isFetchingClass,
        isError: isClassError,
        refetch: refetchClass,
    } = useGetClassByIdQuery(classId || '', {
        skip: !open || !classId,
    });
    const {
        data: rosterResponse,
        isLoading: isLoadingRoster,
        isFetching: isFetchingRoster,
        isError: isRosterError,
        refetch: refetchRoster,
    } = useGetClassRosterQuery(classId || '', {
        skip: !open || !classId,
    });
    const { data: teachersResponse, isFetching: isFetchingTeachers } = useGetTeachersQuery(
        { page: 1, limit: 200, status: 'ACTIVE' },
        { skip: !open },
    );
    const { data: studentsResponse, isFetching: isFetchingStudents } = useGetStudentsQuery(
        { page: 1, limit: 200, status: 'ACTIVE', search: studentSearch },
        { skip: !open || !isEnrollDialogOpen },
    );

    const [enrollStudents, { isLoading: isEnrolling }] = useEnrollStudentsMutation();
    const [removeStudents, { isLoading: isRemoving }] = useRemoveStudentsMutation();
    const [assignTeacher, { isLoading: isAssigning }] = useAssignTeacherMutation();

    const teachers = useMemo(() => teachersResponse?.data?.teachers || [], [teachersResponse?.data?.teachers]);
    const classData = classResponse?.data || null;
    const roster = rosterResponse?.data;

    const teacherForm = useForm<AssignTeacherValues>({
        resolver: zodResolver(assignTeacherSchema),
        defaultValues: { teacher_id: '' },
    });
    const enrollForm = useForm<EnrollStudentsValues>({
        resolver: zodResolver(enrollStudentsSchema),
        defaultValues: { student_ids: [] },
    });

    const currentTeacher = useMemo(
        () => teachers.find((teacher) => teacher.id === classData?.teacher_id) || null,
        [teachers, classData?.teacher_id],
    );

    useEffect(() => {
        teacherForm.reset({ teacher_id: classData?.teacher_id || '' });
    }, [classData?.teacher_id, teacherForm]);

    const rosterStudents = useMemo(() => {
        const keyword = rosterSearch.trim().toLowerCase();
        const rows = roster?.students || [];

        if (!keyword) {
            return rows;
        }

        return rows.filter((student) =>
            [student.full_name, student.code, student.grade_level, student.phone]
                .filter(Boolean)
                .some((value) => value.toLowerCase().includes(keyword)),
        );
    }, [roster?.students, rosterSearch]);

    const availableStudents = useMemo(() => {
        const enrolledIds = new Set((roster?.students || []).map((student) => student.id));
        return (studentsResponse?.data?.students || []).filter((student) => !enrolledIds.has(student.id));
    }, [roster?.students, studentsResponse?.data?.students]);

    const capacityLimit = roster?.capacity_limit || classData?.max_students || 0;
    const currentCount = roster?.current_count || 0;
    const capacityExceeded = capacityLimit > 0 && currentCount > capacityLimit;
    const capacityReached = capacityLimit > 0 && currentCount === capacityLimit;
    const teacherChip = statusChip(classData?.status);

    const handleRefresh = async () => {
        await Promise.all([refetchClass(), refetchRoster()]);
    };

    const handleEnrollSubmit = enrollForm.handleSubmit(async (values) => {
        if (!classId) {
            return;
        }

        try {
            await enrollStudents({ id: classId, student_ids: values.student_ids }).unwrap();
            toast.success('Đã thêm học sinh vào lớp');
            setIsEnrollDialogOpen(false);
            enrollForm.reset({ student_ids: [] });
        } catch (error) {
            toast.error(getErrorMessage(error, 'Không thể thêm học sinh vào lớp'));
        }
    });

    const handleAssignTeacherSubmit = teacherForm.handleSubmit(async (values) => {
        if (!classId) {
            return;
        }

        try {
            await assignTeacher({ id: classId, teacher_id: values.teacher_id }).unwrap();
            toast.success('Đã cập nhật giáo viên phụ trách');
        } catch (error) {
            toast.error(getErrorMessage(error, 'Không thể cập nhật giáo viên phụ trách'));
        }
    });

    const handleRemoveStudent = async () => {
        if (!classId || !studentToRemove) {
            return;
        }

        try {
            await removeStudents({ id: classId, student_ids: [studentToRemove.id] }).unwrap();
            toast.success(`Đã xóa ${studentToRemove.full_name} khỏi lớp`);
            setStudentToRemove(null);
        } catch (error) {
            toast.error(getErrorMessage(error, 'Không thể xóa học sinh khỏi lớp'));
        }
    };

    const rosterColumns: GridColDef<Student>[] = [
        {
            field: 'full_name',
            headerName: 'Học sinh',
            flex: 1.3,
            minWidth: 220,
            renderCell: (params: GridRenderCellParams<Student>) => (
                <Stack direction="row" spacing={1.5} alignItems="center" sx={{ minWidth: 0 }}>
                    <Avatar sx={{ width: 34, height: 34, bgcolor: 'primary.light', color: 'primary.dark' }}>
                        {params.row.full_name.charAt(0)}
                    </Avatar>
                    <Box sx={{ minWidth: 0 }}>
                        <Typography variant="body2" sx={{ fontWeight: 700 }} noWrap>
                            {params.row.full_name}
                        </Typography>
                        <Typography variant="caption" color="text.secondary" noWrap>
                            {params.row.code}
                        </Typography>
                    </Box>
                </Stack>
            ),
        },
        { field: 'grade_level', headerName: 'Khối lớp', width: 120 },
        { field: 'phone', headerName: 'Số điện thoại', width: 140 },
        { field: 'guardian_phone', headerName: 'SĐT phụ huynh', width: 150 },
        {
            field: 'status',
            headerName: 'Trạng thái',
            width: 130,
            renderCell: (params: GridRenderCellParams<Student>) => (
                <Chip
                    size="small"
                    color={params.row.status === 'ACTIVE' ? 'success' : 'default'}
                    label={params.row.status === 'ACTIVE' ? 'Đang học' : 'Không hoạt động'}
                    variant={params.row.status === 'ACTIVE' ? 'filled' : 'outlined'}
                />
            ),
        },
        {
            field: 'actions',
            headerName: '',
            width: 86,
            sortable: false,
            align: 'center',
            renderCell: (params: GridRenderCellParams<Student>) => (
                <Button
                    color="error"
                    size="small"
                    startIcon={<DeleteOutlineRounded />}
                    onClick={() => setStudentToRemove(params.row)}
                >
                    Xóa
                </Button>
            ),
        },
    ];

    const renderLoading = () => (
        <Stack spacing={2}>
            <Skeleton variant="rounded" height={80} />
            <Skeleton variant="rounded" height={220} />
            <Skeleton variant="rounded" height={220} />
        </Stack>
    );

    const renderError = () => (
        <Alert
            severity="error"
            action={
                <Button color="inherit" size="small" startIcon={<RefreshRounded />} onClick={() => void handleRefresh()}>
                    Tải lại
                </Button>
            }
        >
            Không thể tải chi tiết lớp học. Kiểm tra backend và thử lại.
        </Alert>
    );

    const renderInfoTab = () => {
        if (!classData || !roster) {
            return null;
        }

        return (
            <Stack spacing={2.5}>
                <Grid container spacing={2}>
                    <Grid size={{ xs: 12, md: 4 }}>
                        <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 3 }}>
                            <Typography variant="overline" color="text.secondary">
                                Mã lớp
                            </Typography>
                            <Typography variant="h6" sx={{ mt: 0.5, fontWeight: 700 }}>
                                {classData.code}
                            </Typography>
                        </Paper>
                    </Grid>
                    <Grid size={{ xs: 12, md: 4 }}>
                        <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 3 }}>
                            <Typography variant="overline" color="text.secondary">
                                Trạng thái
                            </Typography>
                            <Box sx={{ mt: 1 }}>
                                <Chip label={teacherChip.label} color={teacherChip.color} />
                            </Box>
                        </Paper>
                    </Grid>
                    <Grid size={{ xs: 12, md: 4 }}>
                        <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 3 }}>
                            <Typography variant="overline" color="text.secondary">
                                Sĩ số hiện tại
                            </Typography>
                            <Typography variant="h6" sx={{ mt: 0.5, fontWeight: 700 }}>
                                {currentCount}/{capacityLimit || classData.max_students}
                            </Typography>
                        </Paper>
                    </Grid>
                </Grid>

                <Paper variant="outlined" sx={{ p: 3, borderRadius: 3 }}>
                    <Stack spacing={2}>
                        <Box>
                            <Typography variant="h6" sx={{ fontWeight: 700 }}>
                                {classData.name}
                            </Typography>
                            <Typography variant="body2" color="text.secondary">
                                {classData.notes || 'Chưa có ghi chú cho lớp học này.'}
                            </Typography>
                        </Box>
                        <Divider />
                        <Grid container spacing={2}>
                            <Grid size={{ xs: 12, md: 6 }}>
                                <Typography variant="body2" color="text.secondary">
                                    Ngày bắt đầu
                                </Typography>
                                <Typography variant="body1" sx={{ fontWeight: 600 }}>
                                    {formatDate(classData.start_date)}
                                </Typography>
                            </Grid>
                            <Grid size={{ xs: 12, md: 6 }}>
                                <Typography variant="body2" color="text.secondary">
                                    Ngày kết thúc
                                </Typography>
                                <Typography variant="body1" sx={{ fontWeight: 600 }}>
                                    {formatDate(classData.end_date)}
                                </Typography>
                            </Grid>
                            <Grid size={{ xs: 12, md: 6 }}>
                                <Typography variant="body2" color="text.secondary">
                                    Học phí
                                </Typography>
                                <Typography variant="body1" sx={{ fontWeight: 600 }}>
                                    {new Intl.NumberFormat('vi-VN').format(classData.price || 0)} đ
                                </Typography>
                            </Grid>
                            <Grid size={{ xs: 12, md: 6 }}>
                                <Typography variant="body2" color="text.secondary">
                                    Mốc sức chứa đang áp dụng
                                </Typography>
                                <Typography variant="body1" sx={{ fontWeight: 600 }}>
                                    {capacityLimit || classData.max_students} học sinh
                                </Typography>
                            </Grid>
                        </Grid>
                    </Stack>
                </Paper>
            </Stack>
        );
    };

    const renderRosterTab = () => (
        <Stack spacing={2}>
            {capacityExceeded ? (
                <Alert severity="error">
                    Lớp đang vượt sức chứa cho phép ({currentCount}/{capacityLimit}). Cần điều chỉnh roster hoặc sức chứa phòng/lớp.
                </Alert>
            ) : null}
            {!capacityExceeded && capacityReached ? (
                <Alert severity="warning">
                    Lớp đã đầy ({currentCount}/{capacityLimit}). Chỉ nên thêm học sinh khi đã tăng sức chứa.
                </Alert>
            ) : null}

            <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 3 }}>
                <Stack
                    direction={{ xs: 'column', md: 'row' }}
                    spacing={1.5}
                    justifyContent="space-between"
                    alignItems={{ xs: 'stretch', md: 'center' }}
                    sx={{ mb: 2 }}
                >
                    <TextField
                        value={rosterSearch}
                        onChange={(event) => setRosterSearch(event.target.value)}
                        placeholder="Tìm học sinh theo tên, mã, khối lớp..."
                        size="small"
                        sx={{ minWidth: { xs: '100%', md: 320 } }}
                    />
                    <Stack direction="row" spacing={1}>
                        <Button
                            variant="outlined"
                            startIcon={<RefreshRounded />}
                            onClick={() => void refetchRoster()}
                            disabled={isFetchingRoster}
                        >
                            Làm mới
                        </Button>
                        <Button
                            variant="contained"
                            startIcon={<AddRounded />}
                            onClick={() => setIsEnrollDialogOpen(true)}
                        >
                            Thêm học sinh
                        </Button>
                    </Stack>
                </Stack>

                <DataGrid
                    autoHeight
                    disableRowSelectionOnClick
                    rows={rosterStudents}
                    columns={rosterColumns}
                    loading={isLoadingRoster || isFetchingRoster}
                    getRowId={(row) => row.id}
                    pageSizeOptions={[5, 10, 20]}
                    initialState={{
                        pagination: {
                            paginationModel: {
                                pageSize: 5,
                                page: 0,
                            },
                        },
                    }}
                    localeText={{ noRowsLabel: 'Chưa có học sinh trong lớp này' }}
                    sx={{
                        border: 'none',
                        '& .MuiDataGrid-columnHeaders': {
                            backgroundColor: '#f8fafc',
                        },
                    }}
                />
            </Paper>
        </Stack>
    );

    const renderTeacherTab = () => (
        <Stack spacing={2}>
            <Paper variant="outlined" sx={{ p: 3, borderRadius: 3 }}>
                <Stack spacing={2}>
                    <Box>
                        <Typography variant="h6" sx={{ fontWeight: 700 }}>
                            Giáo viên hiện tại
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                            Gán giáo viên phụ trách để dùng cho demo quản trị và luồng xếp lịch tiếp theo.
                        </Typography>
                    </Box>

                    {currentTeacher ? (
                        <Stack direction="row" spacing={1.5} alignItems="center">
                            <Avatar sx={{ bgcolor: 'secondary.light', color: 'secondary.dark' }}>
                                {currentTeacher.full_name.charAt(0)}
                            </Avatar>
                            <Box>
                                <Typography variant="body1" sx={{ fontWeight: 700 }}>
                                    {currentTeacher.full_name}
                                </Typography>
                                <Typography variant="body2" color="text.secondary">
                                    {currentTeacher.code} • {currentTeacher.employment_type}
                                </Typography>
                            </Box>
                        </Stack>
                    ) : (
                        <Alert severity="info">Lớp này chưa có giáo viên phụ trách.</Alert>
                    )}

                    <Box component="form" onSubmit={handleAssignTeacherSubmit}>
                        <Stack spacing={2}>
                            <Controller
                                control={teacherForm.control}
                                name="teacher_id"
                                render={({ field, fieldState }) => (
                                    <Autocomplete<Teacher, false, false, false>
                                        options={teachers}
                                        loading={isFetchingTeachers}
                                        value={teachers.find((teacher) => teacher.id === field.value) || null}
                                        onChange={(_, value) => field.onChange(value?.id || '')}
                                        getOptionLabel={(option) => `${option.full_name} (${option.code})`}
                                        renderInput={(params) => (
                                            <TextField
                                                {...params}
                                                label="Chọn giáo viên phụ trách"
                                                error={!!fieldState.error}
                                                helperText={fieldState.error?.message || 'Có thể đổi giáo viên nhiều lần để demo nghiệp vụ'}
                                            />
                                        )}
                                    />
                                )}
                            />

                            <Stack direction="row" spacing={1}>
                                <Button type="submit" variant="contained" disabled={isAssigning}>
                                    {isAssigning ? 'Đang lưu...' : 'Lưu giáo viên phụ trách'}
                                </Button>
                                <Button
                                    variant="outlined"
                                    startIcon={<EditOutlined />}
                                    onClick={() => classData && onEdit(classData)}
                                    disabled={!classData}
                                >
                                    Chỉnh sửa thông tin lớp
                                </Button>
                            </Stack>
                        </Stack>
                    </Box>
                </Stack>
            </Paper>
        </Stack>
    );

    return (
        <>
            <Dialog open={open} onClose={onClose} fullWidth maxWidth="lg">
                <DialogTitle sx={{ pb: 1.5 }}>
                    <Stack direction="row" justifyContent="space-between" alignItems="flex-start" spacing={2}>
                        <Box>
                            <Typography variant="h5" sx={{ fontWeight: 800 }}>
                                Chi tiết lớp học
                            </Typography>
                            <Typography variant="body2" color="text.secondary">
                                {classData ? `${classData.code} • ${classData.name}` : 'Đang tải dữ liệu lớp học'}
                            </Typography>
                        </Box>
                        <Button
                            variant="outlined"
                            startIcon={<RefreshRounded />}
                            onClick={() => void handleRefresh()}
                            disabled={isFetchingClass || isFetchingRoster}
                        >
                            Làm mới
                        </Button>
                    </Stack>
                </DialogTitle>
                <DialogContent dividers sx={{ p: 3 }}>
                    {isLoadingClass || isLoadingRoster ? renderLoading() : null}
                    {!isLoadingClass && !isLoadingRoster && (isClassError || isRosterError) ? renderError() : null}
                    {!isLoadingClass && !isLoadingRoster && !isClassError && !isRosterError ? (
                        <Stack spacing={3}>
                            <Paper
                                variant="outlined"
                                sx={{
                                    p: 2.5,
                                    borderRadius: 3,
                                    background: 'linear-gradient(135deg, rgba(14,116,144,0.08), rgba(8,145,178,0.02))',
                                }}
                            >
                                <Stack direction={{ xs: 'column', md: 'row' }} spacing={2} justifyContent="space-between">
                                    <Stack direction="row" spacing={1.5} alignItems="center">
                                        <Avatar sx={{ bgcolor: 'primary.main' }}>
                                            <SchoolOutlined />
                                        </Avatar>
                                        <Box>
                                            <Typography variant="h6" sx={{ fontWeight: 700 }}>
                                                {classData?.name}
                                            </Typography>
                                            <Typography variant="body2" color="text.secondary">
                                                Theo dõi roster, sĩ số và giáo viên phụ trách trong một màn hình.
                                            </Typography>
                                        </Box>
                                    </Stack>
                                    <Stack direction="row" spacing={1} alignItems="center">
                                        <Chip icon={<PersonOutline />} label={`${currentCount} học sinh`} variant="outlined" />
                                        <Chip icon={<InfoOutlined />} label={`Giới hạn ${capacityLimit || 0}`} variant="outlined" />
                                    </Stack>
                                </Stack>
                            </Paper>

                            <Box>
                                <Tabs value={tab} onChange={(_, nextTab) => setTab(nextTab)} sx={{ mb: 2 }}>
                                    <Tab label="Thông tin lớp" />
                                    <Tab label="Danh sách học sinh" />
                                    <Tab label="Giáo viên phụ trách" />
                                </Tabs>

                                {tab === 0 ? renderInfoTab() : null}
                                {tab === 1 ? renderRosterTab() : null}
                                {tab === 2 ? renderTeacherTab() : null}
                            </Box>
                        </Stack>
                    ) : null}
                </DialogContent>
                <DialogActions sx={{ px: 3, py: 2 }}>
                    <Button onClick={onClose}>Đóng</Button>
                </DialogActions>
            </Dialog>

            <Dialog open={isEnrollDialogOpen} onClose={() => setIsEnrollDialogOpen(false)} fullWidth maxWidth="sm">
                <DialogTitle>Thêm học sinh vào lớp</DialogTitle>
                <DialogContent dividers>
                    <Stack spacing={2} sx={{ pt: 0.5 }}>
                        <TextField
                            label="Tìm học sinh"
                            value={studentSearch}
                            onChange={(event) => setStudentSearch(event.target.value)}
                            helperText="Lọc theo tên, mã hoặc trạng thái hoạt động"
                        />

                        <Controller
                            control={enrollForm.control}
                            name="student_ids"
                            render={({ field, fieldState }) => (
                                <Autocomplete<Student, true, false, false>
                                    multiple
                                    options={availableStudents}
                                    loading={isFetchingStudents}
                                    value={availableStudents.filter((student) => field.value.includes(student.id))}
                                    onChange={(_, values) => field.onChange(values.map((student) => student.id))}
                                    getOptionLabel={(option) => `${option.full_name} (${option.code})`}
                                    renderInput={(params) => (
                                        <TextField
                                            {...params}
                                            label="Danh sách học sinh khả dụng"
                                            error={!!fieldState.error}
                                            helperText={fieldState.error?.message || 'Có thể chọn nhiều học sinh trong một lần lưu'}
                                        />
                                    )}
                                />
                            )}
                        />

                        {availableStudents.length === 0 && !isFetchingStudents ? (
                            <Alert severity="info">Không còn học sinh khả dụng để thêm vào lớp này.</Alert>
                        ) : null}
                    </Stack>
                </DialogContent>
                <DialogActions sx={{ px: 3, py: 2 }}>
                    <Button onClick={() => setIsEnrollDialogOpen(false)} color="inherit">
                        Hủy
                    </Button>
                    <Button variant="contained" onClick={() => void handleEnrollSubmit()} disabled={isEnrolling}>
                        {isEnrolling ? 'Đang thêm...' : 'Lưu roster'}
                    </Button>
                </DialogActions>
            </Dialog>

            <ConfirmDialog
                open={!!studentToRemove}
                title="Xóa học sinh khỏi lớp"
                message={
                    studentToRemove
                        ? `Bạn có chắc muốn xóa "${studentToRemove.full_name}" khỏi roster của lớp này không?`
                        : ''
                }
                confirmText="Xóa khỏi lớp"
                isDanger
                loading={isRemoving}
                onClose={() => setStudentToRemove(null)}
                onConfirm={() => void handleRemoveStudent()}
            />
        </>
    );
}
