import {
    Alert,
    Autocomplete,
    Box,
    Button,
    Chip,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    Divider,
    IconButton,
    Paper,
    Skeleton,
    Stack,
    Tab,
    Tabs,
    TextField,
    Typography,
} from '@mui/material';
import {
    AutoStoriesOutlined,
    DeleteOutlineRounded,
    InfoOutlined,
    LinkRounded,
    RefreshRounded,
} from '@mui/icons-material';
import { DataGrid, type GridColDef, type GridRenderCellParams } from '@mui/x-data-grid';
import { format } from 'date-fns';
import { useEffect, useMemo, useState } from 'react';
import { toast } from 'sonner';

import { useGetCoursesQuery, type Course } from '@/api/courseApi';
import {
    useAddCoursesToProgramMutation,
    useGetProgramByIdQuery,
    useRemoveCoursesFromProgramMutation,
} from '@/api/programApi';
import { getApiErrorMessage } from '@/utils/apiError';

interface ProgramDetailDialogProps {
    open: boolean;
    programId: string | null;
    onClose: () => void;
}

const infoField = (label: string, value: string) => (
    <Stack spacing={0.5}>
        <Typography variant="caption" color="text.secondary">
            {label}
        </Typography>
        <Typography variant="body2" sx={{ fontWeight: 600 }}>
            {value || '-'}
        </Typography>
    </Stack>
);

const trackLabelMap: Record<string, string> = {
    BASIC: 'Cơ bản',
    ADVANCED: 'Nâng cao',
    SUPPORT: 'Bổ trợ',
};

const getTrackLabel = (track?: string) => {
    if (!track) {
        return 'Chưa phân hệ';
    }

    return trackLabelMap[track] || track;
};

const getProgramSummary = (effectiveFrom?: string, effectiveTo?: string) => {
    if (!effectiveFrom && !effectiveTo) {
        return 'Theo dõi thông tin chương trình và danh sách khóa học đang liên kết.';
    }

    const from = effectiveFrom ? format(new Date(effectiveFrom), 'dd/MM/yyyy') : '--';
    const to = effectiveTo ? format(new Date(effectiveTo), 'dd/MM/yyyy') : '--';
    return `Hiệu lực áp dụng: ${from} - ${to}.`;
};

const getCourseOptionLabel = (course: Course) => `${course.code} - ${course.name}`;

const matchesCourseSearch = (course: Course, query: string) => {
    const normalizedQuery = query.trim().toLowerCase();
    if (!normalizedQuery) {
        return true;
    }

    return [course.code, course.name, course.subject, course.grade_level]
        .filter(Boolean)
        .some((value) => value?.toLowerCase().includes(normalizedQuery));
};

export default function ProgramDetailDialog({ open, programId, onClose }: ProgramDetailDialogProps) {
    const [activeTab, setActiveTab] = useState(0);
    const [isLinkDialogOpen, setIsLinkDialogOpen] = useState(false);
    const [selectedCourseIds, setSelectedCourseIds] = useState<string[]>([]);

    const {
        data,
        isLoading,
        isFetching,
        isError,
        error,
        refetch,
    } = useGetProgramByIdQuery(programId || '', {
        skip: !open || !programId,
    });

    const {
        data: coursesData,
        isFetching: isFetchingAvailableCourses,
        error: availableCoursesError,
    } = useGetCoursesQuery(
        {
            page: 1,
            limit: 500,
        },
        {
            skip: !open || !programId || !isLinkDialogOpen,
        },
    );

    const [addCoursesToProgram, { isLoading: isAddingCourses }] = useAddCoursesToProgramMutation();
    const [removeCoursesFromProgram, { isLoading: isRemovingCourses }] = useRemoveCoursesFromProgramMutation();

    const program = data?.data || null;
    const courses = useMemo(() => program?.courses || [], [program]);
    const linkedCourseIds = useMemo(() => new Set(courses.map((course) => course.id)), [courses]);
    const availableCourses = useMemo(
        () => (coursesData?.data?.courses || []).filter((course) => !linkedCourseIds.has(course.id)),
        [coursesData, linkedCourseIds],
    );
    const selectedCourses = useMemo(
        () => availableCourses.filter((course) => selectedCourseIds.includes(course.id)),
        [availableCourses, selectedCourseIds],
    );

    useEffect(() => {
        if (!open) {
            setActiveTab(0);
            setIsLinkDialogOpen(false);
            setSelectedCourseIds([]);
        }
    }, [open]);

    useEffect(() => {
        if (!isLinkDialogOpen) {
            setSelectedCourseIds([]);
        }
    }, [isLinkDialogOpen]);

    const handleAddCourses = async () => {
        if (!programId || selectedCourseIds.length === 0) {
            return;
        }

        try {
            await addCoursesToProgram({ programId, courseIds: selectedCourseIds }).unwrap();
            toast.success('Liên kết khóa học thành công');
            setIsLinkDialogOpen(false);
            await refetch();
        } catch (addError) {
            toast.error(getApiErrorMessage(addError, 'Không thể liên kết khóa học vào chương trình'));
        }
    };

    const handleRemoveCourse = async (courseId: string) => {
        if (!programId) {
            return;
        }

        try {
            await removeCoursesFromProgram({ programId, courseIds: [courseId] }).unwrap();
            toast.success('Đã gỡ khóa học khỏi chương trình');
            await refetch();
        } catch (removeError) {
            toast.error(getApiErrorMessage(removeError, 'Không thể gỡ khóa học khỏi chương trình'));
        }
    };

    const columns: GridColDef<Course>[] = [
        {
            field: 'code',
            headerName: 'Mã khóa học',
            width: 140,
            align: 'left',
            headerAlign: 'left',
            renderCell: (params: GridRenderCellParams<Course>) => (
                <Typography variant="body2" sx={{ fontWeight: 700, color: 'primary.main' }}>
                    {params.row.code}
                </Typography>
            ),
        },
        {
            field: 'name',
            headerName: 'Tên khóa học',
            minWidth: 240,
            flex: 1.2,
            align: 'left',
            headerAlign: 'left',
            renderCell: (params: GridRenderCellParams<Course>) => (
                <Box sx={{ minWidth: 0 }}>
                    <Typography variant="body2" sx={{ fontWeight: 700 }} noWrap>
                        {params.row.name}
                    </Typography>
                    <Typography variant="caption" color="text.secondary" noWrap>
                        {params.row.subject || 'Chưa phân môn'}
                    </Typography>
                </Box>
            ),
        },
        {
            field: 'grade_level',
            headerName: 'Khối lớp',
            width: 120,
            align: 'center',
            headerAlign: 'center',
            renderCell: (params: GridRenderCellParams<Course>) => (
                <Typography variant="body2">{params.row.grade_level || '-'}</Typography>
            ),
        },
        {
            field: 'total_hours',
            headerName: 'Tổng giờ',
            width: 100,
            align: 'center',
            headerAlign: 'center',
            renderCell: (params: GridRenderCellParams<Course>) => (
                <Typography variant="body2">{params.row.total_hours || 0}</Typography>
            ),
        },
        {
            field: 'status',
            headerName: 'Trạng thái',
            width: 130,
            align: 'center',
            headerAlign: 'center',
            renderCell: (params: GridRenderCellParams<Course>) => (
                <Chip
                    size="small"
                    color={params.row.status === 'ACTIVE' ? 'success' : 'default'}
                    variant="outlined"
                    label={params.row.status === 'ACTIVE' ? 'Hoạt động' : params.row.status || 'Chưa rõ'}
                />
            ),
        },
        {
            field: 'actions',
            headerName: '',
            width: 86,
            sortable: false,
            align: 'center',
            renderCell: (params: GridRenderCellParams<Course>) => (
                <IconButton
                    size="small"
                    color="error"
                    disabled={isRemovingCourses}
                    onClick={() => void handleRemoveCourse(params.row.id)}
                >
                    <DeleteOutlineRounded fontSize="small" />
                </IconButton>
            ),
        },
    ];

    const renderLoading = () => (
        <Stack spacing={2}>
            <Skeleton variant="rounded" height={84} />
            <Skeleton variant="rounded" height={300} />
        </Stack>
    );

    const renderProgramInfo = () => {
        if (!program) {
            return (
                <Paper
                    variant="outlined"
                    sx={{ p: 4, borderRadius: 3, textAlign: 'center', borderStyle: 'dashed' }}
                >
                    <InfoOutlined sx={{ fontSize: 42, color: 'text.secondary', mb: 1.5 }} />
                    <Typography variant="subtitle1" sx={{ fontWeight: 700, mb: 1 }}>
                        Không có dữ liệu chương trình
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        Chọn lại một chương trình từ danh sách để tải chi tiết.
                    </Typography>
                </Paper>
            );
        }

        return (
            <Stack spacing={2}>
                <Paper variant="outlined" sx={{ p: 3, borderRadius: 3 }}>
                    <Stack spacing={2.5}>
                        <Stack direction="row" spacing={1} alignItems="center" justifyContent="space-between">
                            <Box>
                                <Typography variant="h6" sx={{ fontWeight: 700 }}>
                                    {program.name}
                                </Typography>
                                <Typography variant="body2" color="text.secondary">
                                    {getProgramSummary(program.effective_from, program.effective_to)}
                                </Typography>
                            </Box>
                            <Chip label={getTrackLabel(program.track)} color="primary" variant="outlined" />
                        </Stack>

                        <Divider />

                        <Box
                            sx={{
                                display: 'grid',
                                gridTemplateColumns: { xs: '1fr', sm: 'repeat(2, minmax(0, 1fr))' },
                                gap: 2,
                            }}
                        >
                            {infoField('Mã chương trình', program.code)}
                            {infoField('Tên chương trình', program.name)}
                            {infoField('Hệ đào tạo', getTrackLabel(program.track))}
                            {infoField('Số khóa học liên kết', `${courses.length}`)}
                            {infoField(
                                'Hiệu lực từ',
                                program.effective_from ? format(new Date(program.effective_from), 'dd/MM/yyyy') : '-',
                            )}
                            {infoField(
                                'Hiệu lực đến',
                                program.effective_to ? format(new Date(program.effective_to), 'dd/MM/yyyy') : '-',
                            )}
                        </Box>

                        <Box>
                            <Typography variant="caption" color="text.secondary">
                                Ghi chú phê duyệt
                            </Typography>
                            <Typography variant="body2" sx={{ mt: 0.75 }}>
                                {program.approval_note || 'Chưa có ghi chú phê duyệt'}
                            </Typography>
                        </Box>
                    </Stack>
                </Paper>
            </Stack>
        );
    };

    const renderCoursesTab = () => {
        if (!program) {
            return null;
        }

        if (courses.length === 0) {
            return (
                <Paper
                    variant="outlined"
                    sx={{ p: 4, borderRadius: 3, textAlign: 'center', borderStyle: 'dashed' }}
                >
                    <LinkRounded sx={{ fontSize: 42, color: 'text.secondary', mb: 1.5 }} />
                    <Typography variant="subtitle1" sx={{ fontWeight: 700, mb: 1 }}>
                        Chưa có khóa học nào được liên kết
                    </Typography>
                    <Typography variant="body2" color="text.secondary" sx={{ mb: 2.5 }}>
                        Mở danh sách khóa học khả dụng để chọn và liên kết trực tiếp vào chương trình này.
                    </Typography>
                    <Button variant="contained" startIcon={<LinkRounded />} onClick={() => setIsLinkDialogOpen(true)}>
                        Liên kết khóa học
                    </Button>
                </Paper>
            );
        }

        return (
            <Paper variant="outlined" sx={{ p: 2, borderRadius: 3 }}>
                <Stack
                    direction={{ xs: 'column', sm: 'row' }}
                    spacing={1.5}
                    justifyContent="space-between"
                    alignItems={{ xs: 'stretch', sm: 'center' }}
                    sx={{ mb: 2 }}
                >
                    <Stack spacing={0.5}>
                        <Typography variant="subtitle2" sx={{ fontWeight: 700 }}>
                            Khóa học đang liên kết
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                            Hiện có {courses.length} khóa học thuộc chương trình này.
                        </Typography>
                    </Stack>
                    <Button variant="contained" startIcon={<LinkRounded />} onClick={() => setIsLinkDialogOpen(true)}>
                        Liên kết khóa học
                    </Button>
                </Stack>

                <DataGrid
                    rows={courses}
                    columns={columns}
                    autoHeight
                    disableRowSelectionOnClick
                    hideFooter
                    getRowId={(row) => row.id}
                    localeText={{ noRowsLabel: 'Chưa có khóa học liên kết' }}
                    sx={{
                        border: 'none',
                        '& .MuiDataGrid-columnHeaders': {
                            backgroundColor: '#f8fafc',
                        },
                    }}
                />
            </Paper>
        );
    };

    return (
        <>
            <Dialog open={open} onClose={onClose} fullWidth maxWidth="lg">
                <DialogTitle sx={{ pb: 1.5 }}>
                    <Stack direction="row" justifyContent="space-between" alignItems="center" spacing={2}>
                        <Box>
                            <Typography variant="h6" sx={{ fontWeight: 700 }}>
                                Chi tiết chương trình đào tạo
                            </Typography>
                            <Typography variant="body2" color="text.secondary">
                                Theo dõi thông tin chi tiết và danh sách khóa học thuộc chương trình.
                            </Typography>
                        </Box>
                        <Button
                            variant="outlined"
                            startIcon={<RefreshRounded />}
                            onClick={() => void refetch()}
                            disabled={isFetching || !programId}
                        >
                            Tải lại
                        </Button>
                    </Stack>
                </DialogTitle>

                <DialogContent dividers sx={{ p: 3 }}>
                    {isLoading ? renderLoading() : null}

                    {!isLoading && isError ? (
                        <Alert
                            severity="error"
                            action={
                                <Button color="inherit" size="small" startIcon={<RefreshRounded />} onClick={() => void refetch()}>
                                    Tải lại
                                </Button>
                            }
                        >
                            {getApiErrorMessage(error, 'Không thể tải chi tiết chương trình')}
                        </Alert>
                    ) : null}

                    {!isLoading && !isError ? (
                        <Stack spacing={2}>
                            <Tabs value={activeTab} onChange={(_event, value) => setActiveTab(value)}>
                                <Tab icon={<InfoOutlined fontSize="small" />} iconPosition="start" label="Thông tin chương trình" />
                                <Tab
                                    icon={<AutoStoriesOutlined fontSize="small" />}
                                    iconPosition="start"
                                    label="Khóa học thuộc chương trình"
                                />
                            </Tabs>

                            {activeTab === 0 ? renderProgramInfo() : null}
                            {activeTab === 1 ? renderCoursesTab() : null}
                        </Stack>
                    ) : null}
                </DialogContent>

                <DialogActions sx={{ px: 3, py: 2 }}>
                    <Button onClick={onClose} color="inherit">
                        Đóng
                    </Button>
                </DialogActions>
            </Dialog>

            <Dialog
                open={isLinkDialogOpen}
                onClose={() => setIsLinkDialogOpen(false)}
                fullWidth
                maxWidth="md"
            >
                <DialogTitle>Liên kết khóa học vào chương trình</DialogTitle>
                <DialogContent dividers>
                    <Stack spacing={2}>
                        <Typography variant="body2" color="text.secondary">
                            Tìm trực tiếp trong danh sách khóa học khả dụng và chọn nhiều khóa học trong một lần lưu.
                        </Typography>

                        {availableCoursesError ? (
                            <Alert severity="error">
                                {getApiErrorMessage(availableCoursesError, 'Không thể tải danh sách khóa học khả dụng')}
                            </Alert>
                        ) : null}

                        <Autocomplete
                            multiple
                            disableCloseOnSelect
                            options={availableCourses}
                            loading={isFetchingAvailableCourses}
                            value={selectedCourses}
                            onChange={(_event, value) => setSelectedCourseIds(value.map((course) => course.id))}
                            getOptionLabel={getCourseOptionLabel}
                            isOptionEqualToValue={(option, value) => option.id === value.id}
                            filterOptions={(options, state) =>
                                options.filter((course) => matchesCourseSearch(course, state.inputValue))
                            }
                            noOptionsText={
                                isFetchingAvailableCourses
                                    ? 'Đang tải khóa học...'
                                    : 'Không còn khóa học khả dụng để liên kết'
                            }
                            renderOption={(props, option) => (
                                <Box component="li" {...props}>
                                    <Stack spacing={0.25}>
                                        <Typography variant="body2" sx={{ fontWeight: 600 }}>
                                            {getCourseOptionLabel(option)}
                                        </Typography>
                                        <Typography variant="caption" color="text.secondary">
                                            {(option.subject || 'Chưa phân môn') + ' • ' + (option.grade_level || 'Chưa khai báo khối lớp')}
                                        </Typography>
                                    </Stack>
                                </Box>
                            )}
                            renderInput={(params) => (
                                <TextField
                                    {...params}
                                    label="Danh sách khóa học khả dụng"
                                    placeholder="Tìm theo mã, tên, môn học hoặc khối lớp"
                                    helperText="Ô chọn này đã tích hợp tìm kiếm, không cần thêm bộ lọc riêng."
                                />
                            )}
                        />

                        {!isFetchingAvailableCourses && availableCourses.length === 0 ? (
                            <Alert severity="info">
                                Tất cả khóa học hiện có đã được liên kết hoặc chưa có dữ liệu khóa học để chọn.
                            </Alert>
                        ) : null}
                    </Stack>
                </DialogContent>
                <DialogActions sx={{ px: 3, py: 2 }}>
                    <Button onClick={() => setIsLinkDialogOpen(false)} disabled={isAddingCourses}>
                        Hủy
                    </Button>
                    <Button
                        variant="contained"
                        onClick={() => void handleAddCourses()}
                        disabled={isAddingCourses || selectedCourseIds.length === 0}
                    >
                        {isAddingCourses ? 'Đang lưu...' : 'Lưu liên kết'}
                    </Button>
                </DialogActions>
            </Dialog>
        </>
    );
}
