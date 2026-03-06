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
    TextField,
    Typography,
} from '@mui/material';
import {
    RefreshRounded,
    SearchRounded,
    UploadFileRounded,
    VisibilityOutlined,
} from '@mui/icons-material';
import { DataGrid } from '@mui/x-data-grid';
import type { GridColDef, GridRenderCellParams } from '@mui/x-data-grid';
import { toast } from 'sonner';

import { useGetMaterialsQuery, useUploadMaterialMutation, type MaterialItem } from '@/api/materialApi';
import { useGetTeachersQuery } from '@/api/teacherApi';
import PageHeader from '@/components/common/PageHeader';
import MaterialDetailDialog from '@/components/material/MaterialDetailDialog';

const uploadSchema = z.object({
    teacher_id: z.string().min(1, 'Vui lòng chọn giáo viên'),
    title: z.string().min(1, 'Vui lòng nhập tiêu đề'),
    description: z.string().optional(),
});

type UploadFormValues = z.infer<typeof uploadSchema>;

const statusColorMap: Record<string, 'success' | 'warning' | 'error' | 'default'> = {
    AI_REVIEWED: 'warning',
    APPROVED: 'success',
    REJECTED: 'error',
    SCANNING: 'default',
    UPLOADED: 'default',
};

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

export const TeacherDocumentsPage = () => {
    const [selectedTeacherId, setSelectedTeacherId] = useState('');
    const [search, setSearch] = useState('');
    const [selectedMaterial, setSelectedMaterial] = useState<MaterialItem | null>(null);
    const [selectedFile, setSelectedFile] = useState<File | null>(null);

    const { data: teachersData, isLoading: isLoadingTeachers } = useGetTeachersQuery({ page: 1, limit: 200, status: 'ACTIVE' });
    const { data: materialsData, isLoading, isFetching, refetch } = useGetMaterialsQuery(
        { teacher_id: selectedTeacherId || undefined },
        { skip: !selectedTeacherId }
    );
    const [uploadMaterial, { isLoading: isUploading }] = useUploadMaterialMutation();

    const teachers = teachersData?.data?.teachers || [];
    const materials = useMemo(() => materialsData?.data?.materials || [], [materialsData?.data?.materials]);

    const {
        control,
        handleSubmit,
        formState: { errors },
    } = useForm<UploadFormValues>({
        resolver: zodResolver(uploadSchema),
        defaultValues: {
            teacher_id: '',
            title: '',
            description: '',
        },
    });

    const filteredMaterials = useMemo(() => {
        const keyword = search.trim().toLowerCase();
        if (!keyword) {
            return materials;
        }
        return materials.filter((item) =>
            [item.title, item.file_name, item.status, item.latest_label?.name]
                .filter(Boolean)
                .some((value) => value!.toLowerCase().includes(keyword)),
        );
    }, [materials, search]);

    const onSubmit = handleSubmit(async (values) => {
        if (!selectedFile) {
            toast.error('Vui lòng chọn file để tải lên');
            return;
        }

        try {
            const formData = new FormData();
            formData.append('teacher_id', values.teacher_id);
            formData.append('title', values.title);
            formData.append('description', values.description || '');
            formData.append('file', selectedFile);
            await uploadMaterial(formData).unwrap();
            setSelectedTeacherId(values.teacher_id);
            setSelectedFile(null);
            toast.success('Đã tải tài liệu và chạy AI audit giả lập');
        } catch (error) {
            toast.error(getErrorMessage(error, 'Không thể tải tài liệu'));
        }
    });

    const columns: GridColDef<MaterialItem>[] = [
        { field: 'title', headerName: 'Tiêu đề', minWidth: 220, flex: 1.2 },
        { field: 'file_name', headerName: 'Tệp', minWidth: 180, flex: 1 },
        {
            field: 'latest_label',
            headerName: 'Nhãn AI',
            minWidth: 160,
            renderCell: (params: GridRenderCellParams<MaterialItem>) => (
                params.row.latest_label ? (
                    <Chip
                        label={`${params.row.latest_label.name} (${params.row.latest_label.severity})`}
                        color={params.row.latest_label.severity === 'SAFE' ? 'success' : params.row.latest_label.severity === 'DANGER' ? 'error' : 'warning'}
                        variant="outlined"
                        size="small"
                    />
                ) : (
                    <Chip label="Chưa gán nhãn" size="small" variant="outlined" />
                )
            ),
        },
        {
            field: 'status',
            headerName: 'Trạng thái',
            minWidth: 140,
            renderCell: (params: GridRenderCellParams<MaterialItem>) => (
                <Chip label={params.row.status} color={statusColorMap[params.row.status] || 'default'} size="small" />
            ),
        },
        {
            field: 'uploaded_at',
            headerName: 'Tải lên lúc',
            minWidth: 180,
            renderCell: (params: GridRenderCellParams<MaterialItem>) => (
                <Typography variant="body2">{new Date(params.row.uploaded_at).toLocaleString('vi-VN')}</Typography>
            ),
        },
        {
            field: 'actions',
            headerName: '',
            width: 90,
            sortable: false,
            renderCell: (params: GridRenderCellParams<MaterialItem>) => (
                <Button startIcon={<VisibilityOutlined />} onClick={() => setSelectedMaterial(params.row)}>
                    Xem
                </Button>
            ),
        },
    ];

    return (
        <Box>
            <PageHeader
                title="Tài liệu giảng dạy"
                subtitle="Tải file, theo dõi nhãn AI và xem trạng thái phê duyệt cho demo giáo viên."
                breadcrumbs={[
                    { label: 'Teacher Dashboard', path: '/app/teacher/overview' },
                    { label: 'Tài liệu giảng dạy' },
                ]}
                actions={
                    <Button variant="outlined" startIcon={<RefreshRounded />} onClick={() => void refetch()} disabled={!selectedTeacherId || isFetching}>
                        Làm mới
                    </Button>
                }
            />

            <Stack spacing={3}>
                <Paper variant="outlined" sx={{ p: 3, borderRadius: 4 }}>
                    <Stack spacing={2.5}>
                        <Typography variant="h6" sx={{ fontWeight: 700 }}>
                            Upload tài liệu mới
                        </Typography>

                        <Controller
                            control={control}
                            name="teacher_id"
                            render={({ field }) => (
                                <Autocomplete
                                    options={teachers}
                                    loading={isLoadingTeachers}
                                    value={teachers.find((item) => item.id === field.value) || null}
                                    onChange={(_, value) => {
                                        const teacherId = value?.id || '';
                                        field.onChange(teacherId);
                                        setSelectedTeacherId(teacherId);
                                    }}
                                    getOptionLabel={(option) => `${option.full_name} (${option.code})`}
                                    renderInput={(params) => (
                                        <TextField
                                            {...params}
                                            label="Giáo viên"
                                            error={!!errors.teacher_id}
                                            helperText={errors.teacher_id?.message || 'Demo hiện cho phép chọn giáo viên trực tiếp'}
                                        />
                                    )}
                                />
                            )}
                        />

                        <Stack direction={{ xs: 'column', md: 'row' }} spacing={2}>
                            <Controller
                                control={control}
                                name="title"
                                render={({ field }) => (
                                    <TextField
                                        {...field}
                                        fullWidth
                                        label="Tiêu đề tài liệu"
                                        error={!!errors.title}
                                        helperText={errors.title?.message}
                                    />
                                )}
                            />
                            <Button component="label" variant="outlined" startIcon={<UploadFileRounded />} sx={{ minWidth: 220 }}>
                                {selectedFile ? selectedFile.name : 'Chọn file'}
                                <input
                                    hidden
                                    type="file"
                                    onChange={(event) => setSelectedFile(event.target.files?.[0] || null)}
                                />
                            </Button>
                        </Stack>

                        <Controller
                            control={control}
                            name="description"
                            render={({ field }) => (
                                <TextField
                                    {...field}
                                    fullWidth
                                    multiline
                                    minRows={3}
                                    label="Mô tả"
                                    helperText="Có thể ghi mục tiêu bài giảng hoặc phạm vi tài liệu"
                                />
                            )}
                        />

                        <Stack direction="row" spacing={1}>
                            <Button variant="contained" startIcon={<UploadFileRounded />} onClick={() => void onSubmit()} disabled={isUploading}>
                                {isUploading ? 'Đang tải lên...' : 'Tải lên và audit'}
                            </Button>
                            <Button variant="text" onClick={() => setSelectedFile(null)}>
                                Xóa file đã chọn
                            </Button>
                        </Stack>
                    </Stack>
                </Paper>

                <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 4 }}>
                    <Stack direction={{ xs: 'column', lg: 'row' }} spacing={1.5} justifyContent="space-between" sx={{ mb: 2 }}>
                        <TextField
                            size="small"
                            value={search}
                            onChange={(event) => setSearch(event.target.value)}
                            placeholder="Tìm theo tiêu đề, file, trạng thái..."
                            sx={{ minWidth: { xs: '100%', md: 320 } }}
                            InputProps={{
                                startAdornment: (
                                    <InputAdornment position="start">
                                        <SearchRounded fontSize="small" />
                                    </InputAdornment>
                                ),
                            }}
                        />
                        <Chip label={`${filteredMaterials.length} tài liệu`} color="primary" variant="outlined" />
                    </Stack>

                    {!selectedTeacherId ? (
                        <Alert severity="info">Chọn giáo viên ở form upload để xem danh sách tài liệu tương ứng.</Alert>
                    ) : null}
                    {isLoading ? (
                        <Stack spacing={1.5}>
                            <Skeleton variant="rounded" height={64} />
                            <Skeleton variant="rounded" height={320} />
                        </Stack>
                    ) : null}
                    {!isLoading && selectedTeacherId && filteredMaterials.length === 0 ? (
                        <Alert severity="info">Chưa có tài liệu nào cho giáo viên đang chọn.</Alert>
                    ) : null}
                    {!isLoading && filteredMaterials.length > 0 ? (
                        <DataGrid
                            rows={filteredMaterials}
                            columns={columns}
                            autoHeight
                            disableRowSelectionOnClick
                            getRowId={(row) => row.id}
                            pageSizeOptions={[5, 10, 25]}
                            initialState={{
                                pagination: {
                                    paginationModel: {
                                        page: 0,
                                        pageSize: 10,
                                    },
                                },
                            }}
                            sx={{
                                border: 'none',
                                '& .MuiDataGrid-columnHeaders': {
                                    backgroundColor: '#f8fafc',
                                },
                            }}
                        />
                    ) : null}
                </Paper>
            </Stack>

            <MaterialDetailDialog open={!!selectedMaterial} material={selectedMaterial} onClose={() => setSelectedMaterial(null)} />
        </Box>
    );
};
