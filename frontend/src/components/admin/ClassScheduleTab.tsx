import { useState } from 'react';
import {
    Autocomplete,
    Box,
    Button,
    Grid,
    Paper,
    Stack,
    TextField,
    Typography,
    Chip,
} from '@mui/material';
import { DataGrid } from '@mui/x-data-grid';
import type { GridColDef, GridRenderCellParams } from '@mui/x-data-grid';
import { AddRounded, DeleteOutlineRounded, EventNote } from '@mui/icons-material';
import { toast } from 'sonner';

import {
    useGetClassSchedulesQuery,
    useCreateClassScheduleMutation,
    useDeleteClassScheduleMutation,
    type ClassSchedule,
} from '@/api/classApi';
import { useGetShiftsQuery } from '@/api/shiftApi';
import { useGetRoomsQuery } from '@/api/roomApi';
import ConfirmDialog from '@/components/common/ConfirmDialog';

interface ClassScheduleTabProps {
    classId: string;
}

const DAYS_OF_WEEK = [
    { value: 'MONDAY', label: 'Thứ Hai' },
    { value: 'TUESDAY', label: 'Thứ Ba' },
    { value: 'WEDNESDAY', label: 'Thứ Tư' },
    { value: 'THURSDAY', label: 'Thứ Năm' },
    { value: 'FRIDAY', label: 'Thứ Sáu' },
    { value: 'SATURDAY', label: 'Thứ Bảy' },
    { value: 'SUNDAY', label: 'Chủ Nhật' },
];

export default function ClassScheduleTab({ classId }: ClassScheduleTabProps) {
    const { data: scheduleResponse, isLoading: isLoadingSchedules } = useGetClassSchedulesQuery(classId);
    const { data: shiftsResponse, isLoading: isLoadingShifts } = useGetShiftsQuery({ limit: 100 });
    const { data: roomsResponse, isLoading: isLoadingRooms } = useGetRoomsQuery({ limit: 100 });

    const [createSchedule, { isLoading: isCreating }] = useCreateClassScheduleMutation();
    const [deleteSchedule, { isLoading: isDeleting }] = useDeleteClassScheduleMutation();

    const [selectedDay, setSelectedDay] = useState<string | null>(null);
    const [selectedShift, setSelectedShift] = useState<string | null>(null);
    const [selectedRoom, setSelectedRoom] = useState<string | null>(null);

    const [scheduleToRemove, setScheduleToRemove] = useState<ClassSchedule | null>(null);

    const schedules = scheduleResponse?.data?.schedules || [];
    const shifts = shiftsResponse?.data?.shifts || [];
    const rooms = roomsResponse?.data?.rooms || [];

    const handleCreate = async () => {
        if (!selectedDay || !selectedShift) {
            toast.error('Vui lòng chọn ít nhất Thứ và Ca học');
            return;
        }

        try {
            await createSchedule({
                classId,
                day_of_week: selectedDay,
                shift_id: selectedShift,
                room_id: selectedRoom || undefined,
            }).unwrap();

            toast.success('Thêm lịch tuần thành công');
            setSelectedDay(null);
            setSelectedShift(null);
            setSelectedRoom(null);
        } catch (err: any) {
            toast.error(err?.data?.message || 'Không thể thêm lịch tuần');
        }
    };

    const handleDelete = async () => {
        if (!scheduleToRemove) return;
        try {
            await deleteSchedule({ classId, scheduleId: scheduleToRemove.id }).unwrap();
            toast.success('Xóa lịch tuần thành công');
            setScheduleToRemove(null);
        } catch (err: any) {
            toast.error(err?.data?.message || 'Không thể xóa lịch tuần');
        }
    };

    const columns: GridColDef<ClassSchedule>[] = [
        {
            field: 'day_of_week',
            headerName: 'Ngày trong tuần',
            width: 150,
            renderCell: (params: GridRenderCellParams<ClassSchedule>) => {
                const day = DAYS_OF_WEEK.find(d => d.value === params.row.day_of_week);
                return (
                    <Chip
                        icon={<EventNote />}
                        label={day?.label || params.row.day_of_week}
                        color="primary"
                        variant="outlined"
                        size="small"
                    />
                );
            },
        },
        {
            field: 'shift',
            headerName: 'Ca học',
            flex: 1,
            renderCell: (params: GridRenderCellParams<ClassSchedule>) => {
                const shift = params.row.shift;
                if (!shift) return <Typography variant="body2" color="text.secondary">N/A</Typography>;
                return (
                    <Typography variant="body2" sx={{ fontWeight: 600 }}>
                        {shift.name} ({shift.start_time.substring(0, 5)} - {shift.end_time.substring(0, 5)})
                    </Typography>
                );
            },
        },
        {
            field: 'room',
            headerName: 'Phòng học',
            width: 200,
            renderCell: (params: GridRenderCellParams<ClassSchedule>) => {
                const room = params.row.room;
                if (!room) return <Typography variant="body2" color="text.secondary">Không xếp phòng</Typography>;
                return (
                    <Typography variant="body2">
                        {room.name} {room.code ? `(${room.code})` : ''}
                    </Typography>
                );
            },
        },
        {
            field: 'actions',
            headerName: '',
            width: 86,
            sortable: false,
            align: 'center',
            renderCell: (params: GridRenderCellParams<ClassSchedule>) => (
                <Button
                    color="error"
                    size="small"
                    startIcon={<DeleteOutlineRounded />}
                    onClick={() => setScheduleToRemove(params.row)}
                >
                    Xóa
                </Button>
            ),
        },
    ];

    return (
        <Stack spacing={3}>
            <Paper variant="outlined" sx={{ p: 3, borderRadius: 3 }}>
                <Box mb={2}>
                    <Typography variant="h6" sx={{ fontWeight: 700 }}>Theo dõi lịch tuần</Typography>
                    <Typography variant="body2" color="text.secondary">
                        Thêm các ca học cố định hàng tuần cho lớp học. Khi tính năng Xếp Lịch Tự Động chạy, các ca học này sẽ được tự động xếp lịch nếu phù hợp tính sẵn sàng.
                    </Typography>
                </Box>
                <Grid container spacing={2} alignItems="center">
                    <Grid size={{ xs: 12, md: 3 }}>
                        <Autocomplete
                            options={DAYS_OF_WEEK}
                            getOptionLabel={(option) => option.label}
                            value={DAYS_OF_WEEK.find(d => d.value === selectedDay) || null}
                            onChange={(_, newValue) => setSelectedDay(newValue?.value || null)}
                            renderInput={(params) => <TextField {...params} label="Ngày trong tuần *" size="small" />}
                        />
                    </Grid>
                    <Grid size={{ xs: 12, md: 3 }}>
                        <Autocomplete
                            options={shifts}
                            getOptionLabel={(option) => `${option.name} (${option.start_time.substring(0, 5)} - ${option.end_time.substring(0, 5)})`}
                            value={shifts.find(s => s.id === selectedShift) || null}
                            onChange={(_, newValue) => setSelectedShift(newValue?.id || null)}
                            renderInput={(params) => <TextField {...params} label="Ca học *" size="small" />}
                            loading={isLoadingShifts}
                        />
                    </Grid>
                    <Grid size={{ xs: 12, md: 3 }}>
                        <Autocomplete
                            options={rooms}
                            getOptionLabel={(option) => `${option.name} ${option.code ? `(${option.code})` : ''}`}
                            value={rooms.find(r => r.id === selectedRoom) || null}
                            onChange={(_, newValue) => setSelectedRoom(newValue?.id || null)}
                            renderInput={(params) => <TextField {...params} label="Phòng học (tùy chọn)" size="small" />}
                            loading={isLoadingRooms}
                        />
                    </Grid>
                    <Grid size={{ xs: 12, md: 3 }}>
                        <Button
                            variant="contained"
                            fullWidth
                            startIcon={<AddRounded />}
                            onClick={handleCreate}
                            disabled={!selectedDay || !selectedShift || isCreating}
                        >
                            Thêm lịch tuần
                        </Button>
                    </Grid>
                </Grid>
            </Paper>

            <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 3 }}>
                <DataGrid
                    autoHeight
                    disableRowSelectionOnClick
                    rows={schedules}
                    columns={columns}
                    loading={isLoadingSchedules}
                    getRowId={(row) => row.id}
                    hideFooter
                    localeText={{ noRowsLabel: 'Chưa có lịch tuần nào được lên cho lớp này' }}
                    sx={{
                        border: 'none',
                        '& .MuiDataGrid-columnHeaders': {
                            backgroundColor: '#f8fafc',
                        },
                    }}
                />
            </Paper>

            <ConfirmDialog
                open={!!scheduleToRemove}
                title="Xóa lịch tuần"
                message={
                    scheduleToRemove
                        ? `Bạn có chắc muốn xóa ca lịch tuần ${DAYS_OF_WEEK.find(d => d.value === scheduleToRemove.day_of_week)?.label || ''} khỏi lớp không?`
                        : ''
                }
                confirmText="Xóa lịch"
                isDanger
                loading={isDeleting}
                onClose={() => setScheduleToRemove(null)}
                onConfirm={handleDelete}
            />
        </Stack>
    );
}
