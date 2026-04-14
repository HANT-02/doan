import { useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
    Alert,
    Box,
    Button,
    Chip,
    Divider,
    MenuItem,
    Paper,
    Skeleton,
    Stack,
    TextField,
    Typography,
} from '@mui/material';
import { DatePicker, LocalizationProvider } from '@mui/x-date-pickers';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import {
    RefreshRounded,
    ClassOutlined,
    AccessTimeOutlined,
    VisibilityOutlined,
} from '@mui/icons-material';
import { DataGrid } from '@mui/x-data-grid';
import type { GridColDef, GridRenderCellParams } from '@mui/x-data-grid';
import { format } from 'date-fns';
import { vi } from 'date-fns/locale';

import { useGetLessonsQuery, type Lesson } from '@/api/lessonApi';
import { useGetClassesQuery } from '@/api/classApi';
import { useGetTeachersQuery } from '@/api/teacherApi';
import PageHeader from '@/components/common/PageHeader';
import { getApiErrorMessage } from '@/utils/apiError';

export default function LessonsPage() {
    const navigate = useNavigate();
    const [page, setPage] = useState(0);
    const [pageSize, setPageSize] = useState(20);
    const [classId, setClassId] = useState('');
    const [teacherId, setTeacherId] = useState('');
    const [dateFrom, setDateFrom] = useState<Date | null>(null);
    const [dateTo, setDateTo] = useState<Date | null>(null);
    const [sortBy, setSortBy] = useState('date_start');
    const [sortOrder, setSortOrder] = useState('asc');

    const buildDateStr = (d: Date | null) => (d ? format(d, 'yyyy-MM-dd') : undefined);

    const queryParams = useMemo(
        () => ({
            page: page + 1,
            limit: pageSize,
            class_id: classId || undefined,
            teacher_id: teacherId || undefined,
            date_from: buildDateStr(dateFrom),
            date_to: buildDateStr(dateTo),
            sortBy,
            sortOrder,
        }),
        [page, pageSize, classId, teacherId, dateFrom, dateTo, sortBy, sortOrder],
    );

    const { data: lessonsResponse, isLoading, isFetching, isError, error, refetch } = useGetLessonsQuery(queryParams);
    const { data: classesResponse } = useGetClassesQuery({ limit: 200 });
    const { data: teachersResponse } = useGetTeachersQuery({ limit: 200 });

    const lessons = lessonsResponse?.data?.lessons || [];
    const totalItems = lessonsResponse?.data?.pagination?.total_items || 0;
    const classes = classesResponse?.data?.classes || [];
    const teachers = teachersResponse?.data?.teachers || [];
    const hasFilters = !!(classId || teacherId || dateFrom || dateTo);

    const formatTime = (iso: string) => {
        try {
            return format(new Date(iso), 'dd/MM/yyyy HH:mm', { locale: vi });
        } catch {
            return iso;
        }
    };

    const columns: GridColDef<Lesson>[] = [
        {
            field: 'class',
            headerName: 'Lớp học',
            flex: 1.2,
            minWidth: 180,
            renderCell: (params: GridRenderCellParams<Lesson>) => (
                <Stack direction="row" spacing={1} alignItems="center">
                    <ClassOutlined sx={{ color: 'primary.main', fontSize: 18 }} />
                    <Box>
                        <Typography variant="body2" fontWeight={700} noWrap>
                            {params.row.class?.name || params.row.class_id}
                        </Typography>
                        <Typography variant="caption" color="text.secondary">
                            {params.row.class?.code}
                        </Typography>
                    </Box>
                </Stack>
            ),
        },
        {
            field: 'teacher',
            headerName: 'Giáo viên',
            flex: 1,
            minWidth: 150,
            renderCell: (params: GridRenderCellParams<Lesson>) => (
                <Typography variant="body2">
                    {params.row.teacher?.full_name || '—'}
                </Typography>
            ),
        },
        {
            field: 'date_start',
            headerName: 'Bắt đầu',
            width: 170,
            renderCell: (params: GridRenderCellParams<Lesson>) => (
                <Stack direction="row" spacing={0.5} alignItems="center">
                    <AccessTimeOutlined sx={{ fontSize: 15, color: 'text.secondary' }} />
                    <Typography variant="body2">{formatTime(params.row.date_start)}</Typography>
                </Stack>
            ),
        },
        {
            field: 'date_end',
            headerName: 'Kết thúc',
            width: 170,
            renderCell: (params: GridRenderCellParams<Lesson>) => (
                <Typography variant="body2">{formatTime(params.row.date_end)}</Typography>
            ),
        },
        {
            field: 'room',
            headerName: 'Phòng học',
            width: 140,
            renderCell: (params: GridRenderCellParams<Lesson>) => (
                <Typography variant="body2">
                    {params.row.room?.name || '—'}
                </Typography>
            ),
        },
        {
            field: 'notes',
            headerName: 'Ghi chú',
            flex: 1,
            minWidth: 120,
            renderCell: (params: GridRenderCellParams<Lesson>) => (
                <Typography variant="body2" color="text.secondary" noWrap>
                    {params.row.notes || '—'}
                </Typography>
            ),
        },
        {
            field: 'actions',
            headerName: 'Chi tiết',
            width: 130,
            sortable: false,
            filterable: false,
            align: 'center',
            headerAlign: 'center',
            renderCell: (params: GridRenderCellParams<Lesson>) => (
                <Button
                    size="small"
                    startIcon={<VisibilityOutlined />}
                    onClick={() => navigate(`/app/admin/lessons/${params.row.id}`)}
                >
                    Xem
                </Button>
            ),
        },
    ];

    return (
        <LocalizationProvider dateAdapter={AdapterDateFns} adapterLocale={vi}>
            <Stack sx={{ p: { xs: 2, md: 4 } }} spacing={3}>
                <PageHeader
                    title="Danh sách buổi học"
                    subtitle="Xem và lọc toàn bộ buổi học đã được lên lịch trong hệ thống."
                    icon={<ClassOutlined />}
                    actions={
                        <Button
                            startIcon={<RefreshRounded />}
                            variant="outlined"
                            onClick={() => void refetch()}
                            disabled={isFetching}
                        >
                            Làm mới
                        </Button>
                    }
                />

                <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 3 }}>
                    <Box
                        sx={{
                            display: 'grid',
                            gridTemplateColumns: {
                                xs: '1fr',
                                md: 'repeat(2, minmax(0, 1fr))',
                                xl: 'repeat(5, minmax(0, 1fr))',
                            },
                            gap: 2,
                            alignItems: 'start',
                        }}
                    >
                        <TextField
                            select
                            label="Lọc theo lớp"
                            value={classId}
                            onChange={(e) => { setClassId(e.target.value); setPage(0); }}
                            size="small"
                            fullWidth
                        >
                            <MenuItem value="">Tất cả lớp</MenuItem>
                            {classes.map((c) => (
                                <MenuItem key={c.id} value={c.id}>
                                    {c.name} ({c.code})
                                </MenuItem>
                            ))}
                        </TextField>

                        <TextField
                            select
                            label="Lọc theo giáo viên"
                            value={teacherId}
                            onChange={(e) => { setTeacherId(e.target.value); setPage(0); }}
                            size="small"
                            fullWidth
                        >
                            <MenuItem value="">Tất cả giáo viên</MenuItem>
                            {teachers.map((t) => (
                                <MenuItem key={t.id} value={t.id}>
                                    {t.full_name}
                                </MenuItem>
                            ))}
                        </TextField>

                        <DatePicker
                            label="Từ ngày"
                            value={dateFrom}
                            onChange={(d: Date | null) => { setDateFrom(d); setPage(0); }}
                            slotProps={{ textField: { size: 'small', fullWidth: true } }}
                        />

                        <DatePicker
                            label="Đến ngày"
                            value={dateTo}
                            onChange={(d: Date | null) => { setDateTo(d); setPage(0); }}
                            slotProps={{ textField: { size: 'small', fullWidth: true } }}
                        />

                        <TextField
                            select
                            label="Sắp xếp theo"
                            value={`${sortBy}_${sortOrder}`}
                            onChange={(e) => {
                                const [field, order] = e.target.value.split('_');
                                setSortBy(field);
                                setSortOrder(order);
                                setPage(0);
                            }}
                            size="small"
                            fullWidth
                        >
                            <MenuItem value="date_start_asc">Ngày bắt đầu ↑</MenuItem>
                            <MenuItem value="date_start_desc">Ngày bắt đầu ↓</MenuItem>
                            <MenuItem value="date_end_asc">Ngày kết thúc ↑</MenuItem>
                            <MenuItem value="date_end_desc">Ngày kết thúc ↓</MenuItem>
                        </TextField>
                    </Box>

                    <Divider sx={{ my: 2 }} />

                    <Stack
                        direction={{ xs: 'column', md: 'row' }}
                        justifyContent="space-between"
                        alignItems={{ xs: 'flex-start', md: 'center' }}
                        spacing={1.5}
                    >
                        <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                            <Chip
                                size="small"
                                color="primary"
                                variant="outlined"
                                label={hasFilters ? 'Đang lọc dữ liệu' : 'Toàn bộ lesson'}
                            />
                            <Chip
                                size="small"
                                variant="outlined"
                                label={`Sắp xếp: ${sortBy === 'date_end' ? 'Kết thúc' : 'Bắt đầu'} ${sortOrder === 'desc' ? 'giảm dần' : 'tăng dần'}`}
                            />
                        </Stack>

                        {hasFilters ? (
                            <Button
                                variant="text"
                                onClick={() => {
                                    setClassId('');
                                    setTeacherId('');
                                    setDateFrom(null);
                                    setDateTo(null);
                                    setPage(0);
                                }}
                                size="small"
                            >
                                Xóa bộ lọc
                            </Button>
                        ) : null}
                    </Stack>
                </Paper>

                {isError ? (
                    <Alert severity="error" action={
                        <Button size="small" onClick={() => void refetch()} startIcon={<RefreshRounded />}>Thử lại</Button>
                    }>
                        {getApiErrorMessage(error, 'Không thể tải danh sách buổi học. Kiểm tra backend và thử lại.')}
                    </Alert>
                ) : null}

                {isLoading ? (
                    <Stack spacing={1.5}>
                        <Skeleton variant="rounded" height={60} />
                        <Skeleton variant="rounded" height={60} />
                        <Skeleton variant="rounded" height={60} />
                        <Skeleton variant="rounded" height={60} />
                    </Stack>
                ) : (
                    <Paper variant="outlined" sx={{ borderRadius: 3 }}>
                        <Box sx={{ px: 2.5, pt: 2, pb: 1 }}>
                            <Stack direction="row" justifyContent="space-between" alignItems="center">
                                <Typography variant="subtitle2" color="text.secondary">
                                    {totalItems} buổi học
                                    {hasFilters ? ' (đã lọc)' : ' tổng cộng'}
                                </Typography>
                                {isFetching ? (
                                    <Chip size="small" label="Đang tải..." color="info" variant="outlined" />
                                ) : null}
                            </Stack>
                        </Box>
                        <DataGrid
                            autoHeight
                            disableRowSelectionOnClick
                            rows={lessons}
                            columns={columns}
                            rowCount={totalItems}
                            paginationMode="server"
                            paginationModel={{ page, pageSize }}
                            onPaginationModelChange={(model) => {
                                setPage(model.page);
                                setPageSize(model.pageSize);
                            }}
                            pageSizeOptions={[10, 20, 50]}
                            loading={isLoading || isFetching}
                            getRowId={(row) => row.id}
                            onRowClick={(params) => navigate(`/app/admin/lessons/${params.row.id}`)}
                            localeText={{ noRowsLabel: 'Chưa có buổi học nào được tạo. Hãy chạy commit scheduling để sinh buổi học.' }}
                            sx={{
                                border: 'none',
                                '& .MuiDataGrid-columnHeaders': {
                                    backgroundColor: '#f8fafc',
                                },
                                '& .MuiDataGrid-row': {
                                    cursor: 'pointer',
                                },
                            }}
                        />
                    </Paper>
                )}
            </Stack>
        </LocalizationProvider>
    );
}
