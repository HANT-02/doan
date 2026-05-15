import { useEffect, useMemo } from 'react';
import { Controller, useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import {
    Autocomplete,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    Button,
    Chip,
    TextField,
    MenuItem,
    Grid,
    CircularProgress,
    Stack,
    Typography,
} from '@mui/material';
import type { Class } from '@/api/classApi';
import { useGetSkillCatalogQuery, useGetTeachersQuery } from '@/api/teacherApi';
import { useGetProgramsQuery } from '@/api/programApi';
import { useGetCoursesQuery } from '@/api/courseApi';

const classSchema = z.object({
    code: z.string().min(1, 'Mã lớp không được để trống'),
    name: z.string().min(1, 'Tên lớp không được để trống'),
    start_date: z.string().min(1, 'Ngày bắt đầu không được để trống'),
    end_date: z.string().optional(),
    max_students: z.number().min(1, 'Sĩ số tối đa phải lớn hơn 0'),
    status: z.enum(['OPEN', 'CLOSED', 'CANCELLED']),
    price: z.number().min(0, 'Học phí không được âm'),
    teacher_id: z.string().optional(),
    program_id: z.string().optional(),
    course_id: z.string().optional(),
    teacher_skill_codes: z.array(z.string()).optional(),
    notes: z.string().optional(),
});

type ClassFormValues = z.infer<typeof classSchema>;
type ClassDialogSubmitValues = Omit<ClassFormValues, 'teacher_id' | 'program_id' | 'course_id' | 'teacher_skill_codes'> & {
    teacher_id?: string;
    program_id?: string;
    course_id?: string;
};

interface ClassDialogProps {
    open: boolean;
    onClose: () => void;
    onSubmit: (data: ClassDialogSubmitValues) => Promise<void>;
    classData?: Class | null;
    isLoading?: boolean;
}

const ClassDialog = ({ open, onClose, onSubmit, classData, isLoading }: ClassDialogProps) => {
    const { data: programsData } = useGetProgramsQuery({ limit: 100 });
    const programs = programsData?.data?.programs || [];

    const { data: coursesData } = useGetCoursesQuery({ limit: 200 });
    const courses = coursesData?.data?.courses || [];
    const { data: skillCatalogResponse, isFetching: isFetchingSkillCatalog } = useGetSkillCatalogQuery({ limit: 200 });
    const skillCatalog = skillCatalogResponse?.data?.skills || [];

    const {
        control,
        handleSubmit,
        reset,
        setValue,
        watch,
        formState: { errors },
    } = useForm<ClassFormValues>({
        resolver: zodResolver(classSchema),
        defaultValues: {
            code: '',
            name: '',
            start_date: new Date().toISOString().split('T')[0],
            max_students: 30,
            status: 'OPEN',
            price: 0,
            teacher_skill_codes: [],
        },
    });
    const selectedCourseId = watch('course_id');
    const selectedSkillCodes = watch('teacher_skill_codes') || [];
    const selectedCourse = useMemo(
        () => courses.find((course) => course.id === selectedCourseId) || null,
        [courses, selectedCourseId],
    );
    const { data: teachersData } = useGetTeachersQuery({ limit: 100, course_id: selectedCourseId || undefined });
    const teachers = useMemo(() => {
        const source = teachersData?.data?.teachers || [];
        if (selectedSkillCodes.length === 0) {
            return source;
        }

        return source.filter((teacher) =>
            selectedSkillCodes.every((skillCode) => teacher.skills?.includes(skillCode)),
        );
    }, [selectedSkillCodes, teachersData?.data?.teachers]);

    useEffect(() => {
        if (classData) {
            reset({
                code: classData.code,
                name: classData.name,
                start_date: classData.start_date.split('T')[0],
                end_date: classData.end_date?.split('T')[0] || '',
                max_students: classData.max_students,
                status: classData.status,
                price: classData.price,
                teacher_id: classData.teacher_id || '',
                program_id: classData.program_id || '',
                course_id: classData.course_id || '',
                teacher_skill_codes: selectedCourse?.required_skills || [],
                notes: classData.notes || '',
            });
        } else {
            reset({
                code: '',
                name: '',
                start_date: new Date().toISOString().split('T')[0],
                end_date: '',
                max_students: 30,
                status: 'OPEN',
                price: 0,
                teacher_id: '',
                program_id: '',
                course_id: '',
                teacher_skill_codes: [],
                notes: '',
            });
        }
    }, [classData, reset, open]);

    useEffect(() => {
        setValue('teacher_skill_codes', selectedCourse?.required_skills || [], { shouldDirty: false });
    }, [selectedCourse?.id, selectedCourse?.required_skills, setValue]);

    const handleFormSubmit = async (data: ClassFormValues) => {
        // Clean up empty strings for optional IDs
        const cleanedData = {
            ...data,
            teacher_id: data.teacher_id || undefined,
            program_id: data.program_id || undefined,
            course_id: data.course_id || undefined,
        };
        delete (cleanedData as Partial<ClassFormValues>).teacher_skill_codes;
        await onSubmit(cleanedData);
        onClose();
    };

    return (
        <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
            <DialogTitle>
                {classData ? 'Chỉnh sửa lớp học' : 'Thêm lớp học mới'}
            </DialogTitle>
            <DialogContent dividers>
                <Grid container spacing={2} sx={{ mt: 0.5 }}>
                    <Grid size={4}>
                        <Controller
                            name="code"
                            control={control}
                            render={({ field }) => (
                                <TextField
                                    {...field}
                                    fullWidth
                                    label="Mã lớp"
                                    error={!!errors.code}
                                    helperText={errors.code?.message}
                                />
                            )}
                        />
                    </Grid>
                    <Grid size={8}>
                        <Controller
                            name="name"
                            control={control}
                            render={({ field }) => (
                                <TextField
                                    {...field}
                                    fullWidth
                                    label="Tên lớp"
                                    error={!!errors.name}
                                    helperText={errors.name?.message}
                                />
                            )}
                        />
                    </Grid>
                    <Grid size={6}>
                        <Controller
                            name="start_date"
                            control={control}
                            render={({ field }) => (
                                <TextField
                                    {...field}
                                    fullWidth
                                    label="Ngày bắt đầu"
                                    type="date"
                                    InputLabelProps={{ shrink: true }}
                                    error={!!errors.start_date}
                                    helperText={errors.start_date?.message}
                                />
                            )}
                        />
                    </Grid>
                    <Grid size={6}>
                        <Controller
                            name="end_date"
                            control={control}
                            render={({ field }) => (
                                <TextField
                                    {...field}
                                    fullWidth
                                    label="Ngày kết thúc"
                                    type="date"
                                    InputLabelProps={{ shrink: true }}
                                    error={!!errors.end_date}
                                    helperText={errors.end_date?.message}
                                />
                            )}
                        />
                    </Grid>
                    <Grid size={4}>
                        <Controller
                            name="max_students"
                            control={control}
                            render={({ field }) => (
                                <TextField
                                    {...field}
                                    fullWidth
                                    label="Sĩ số tối đa"
                                    type="number"
                                    onChange={(event) => field.onChange(Number(event.target.value))}
                                    error={!!errors.max_students}
                                    helperText={errors.max_students?.message}
                                />
                            )}
                        />
                    </Grid>
                    <Grid size={4}>
                        <Controller
                            name="price"
                            control={control}
                            render={({ field }) => (
                                <TextField
                                    {...field}
                                    fullWidth
                                    label="Học phí"
                                    type="number"
                                    onChange={(event) => field.onChange(Number(event.target.value))}
                                    error={!!errors.price}
                                    helperText={errors.price?.message}
                                />
                            )}
                        />
                    </Grid>
                    <Grid size={4}>
                        <Controller
                            name="status"
                            control={control}
                            render={({ field }) => (
                                <TextField
                                    {...field}
                                    fullWidth
                                    select
                                    label="Trạng thái"
                                    error={!!errors.status}
                                    helperText={errors.status?.message}
                                >
                                    <MenuItem value="OPEN">Đang mở</MenuItem>
                                    <MenuItem value="CLOSED">Đã đóng</MenuItem>
                                    <MenuItem value="CANCELLED">Hủy bỏ</MenuItem>
                                </TextField>
                            )}
                        />
                    </Grid>
                    <Grid size={6}>
                        <Controller
                            name="program_id"
                            control={control}
                            render={({ field }) => (
                                <TextField
                                    {...field}
                                    fullWidth
                                    select
                                    label="Chương trình học"
                                >
                                    <MenuItem value=""><em>Chưa gán</em></MenuItem>
                                    {programs.map((p) => (
                                        <MenuItem key={p.id} value={p.id}>
                                            {p.name} ({p.code})
                                        </MenuItem>
                                    ))}
                                </TextField>
                            )}
                        />
                    </Grid>
                    <Grid size={6}>
                        <Controller
                            name="course_id"
                            control={control}
                            render={({ field }) => (
                                <TextField
                                    {...field}
                                    fullWidth
                                    select
                                    label="Khóa học"
                                >
                                    <MenuItem value=""><em>Chưa gán</em></MenuItem>
                                    {courses.map((c) => (
                                        <MenuItem key={c.id} value={c.id}>
                                            {c.name} ({c.code})
                                        </MenuItem>
                                    ))}
                                </TextField>
                            )}
                        />
                    </Grid>
                    <Grid size={12}>
                        <Controller
                            name="teacher_skill_codes"
                            control={control}
                            render={({ field }) => (
                                <Autocomplete<string, true, false, false>
                                    multiple
                                    options={skillCatalog.map((item) => item.code)}
                                    loading={isFetchingSkillCatalog}
                                    value={field.value || []}
                                    onChange={(_, value) => field.onChange(value)}
                                    filterSelectedOptions
                                    disableCloseOnSelect
                                    getOptionLabel={(option) => {
                                        const matchedSkill = skillCatalog.find((item) => item.code === option);
                                        return matchedSkill ? `${matchedSkill.name} (${matchedSkill.code})` : option;
                                    }}
                                    renderTags={(value, getTagProps) =>
                                        value.map((option, index) => {
                                            const matchedSkill = skillCatalog.find((item) => item.code === option);
                                            return (
                                                <Chip
                                                    {...getTagProps({ index })}
                                                    key={option}
                                                    label={matchedSkill?.name || option}
                                                    size="small"
                                                    color="primary"
                                                    variant="outlined"
                                                />
                                            );
                                        })
                                    }
                                    renderInput={(params) => (
                                        <TextField
                                            {...params}
                                            label="Lọc giáo viên theo kỹ năng"
                                            helperText="Tự gợi ý theo khóa học. Bạn có thể bỏ bớt kỹ năng để mở rộng danh sách giáo viên."
                                        />
                                    )}
                                />
                            )}
                        />
                    </Grid>
                    <Grid size={12}>
                        <Controller
                            name="teacher_id"
                            control={control}
                            render={({ field }) => (
                                <Autocomplete
                                    options={teachers}
                                    value={teachers.find((teacher) => teacher.id === field.value) || null}
                                    onChange={(_, value) => field.onChange(value?.id || '')}
                                    getOptionLabel={(option) => `${option.full_name || 'N/A'} (${option.code || 'N/A'})`}
                                    isOptionEqualToValue={(option, value) => option.id === value.id}
                                    renderOption={(props, option) => (
                                        <Stack component="li" {...props} spacing={0.25} sx={{ py: 1 }}>
                                            <Typography variant="body2" sx={{ fontWeight: 700 }}>
                                                {option.full_name || 'Chưa có tên'}
                                            </Typography>
                                            <Typography variant="caption" color="text.secondary">
                                                {option.code || 'N/A'}
                                                {option.skills?.length ? ` • ${option.skills.join(', ')}` : ''}
                                            </Typography>
                                        </Stack>
                                    )}
                                    renderInput={(params) => (
                                        <TextField
                                            {...params}
                                            label="Giáo viên chủ nhiệm"
                                            helperText={selectedSkillCodes.length > 0 && teachers.length === 0
                                                ? 'Không có giáo viên phù hợp với nhóm kỹ năng đang lọc.'
                                                : undefined}
                                        />
                                    )}
                                />
                            )}
                        />
                    </Grid>
                    <Grid size={12}>
                        <Controller
                            name="notes"
                            control={control}
                            render={({ field }) => (
                                <TextField
                                    {...field}
                                    fullWidth
                                    label="Ghi chú"
                                    multiline
                                    rows={2}
                                />
                            )}
                        />
                    </Grid>
                </Grid>
            </DialogContent>
            <DialogActions sx={{ px: 3, py: 2 }}>
                <Button onClick={onClose} disabled={isLoading}>
                    Hủy
                </Button>
                <Button
                    variant="contained"
                    onClick={handleSubmit(handleFormSubmit)}
                    disabled={isLoading}
                    startIcon={isLoading ? <CircularProgress size={20} color="inherit" /> : null}
                >
                    {classData ? 'Cập nhật' : 'Thêm mới'}
                </Button>
            </DialogActions>
        </Dialog>
    );
};

export default ClassDialog;
