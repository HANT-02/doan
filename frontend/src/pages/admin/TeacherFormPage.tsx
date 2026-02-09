import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import {
    Container,
    Paper,
    Box,
    Typography,
    Button,
    TextField,
    Grid,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
    FormControlLabel,
    Switch,
    Alert,
    Skeleton,
} from '@mui/material';
import { ArrowBack, Save } from '@mui/icons-material';
import { teacherApi, type CreateTeacherRequest } from '@/api/teacherApi';
import { createTeacherSchema, updateTeacherSchema } from '@/schemas/teacherSchema';
import { toast } from 'sonner';

export const TeacherFormPage = () => {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const isEditMode = Boolean(id);

    const [loading, setLoading] = useState(isEditMode);
    const [submitting, setSubmitting] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const {
        control,
        handleSubmit,
        reset,
        formState: { errors },
    } = useForm<CreateTeacherRequest>({
        resolver: zodResolver(isEditMode ? updateTeacherSchema : createTeacherSchema),
        defaultValues: {
            full_name: '',
            email: '',
            phone: '',
            code: '',
            is_school_teacher: false,
            school_name: '',
            employment_type: 'FULL_TIME',
            status: 'ACTIVE',
            notes: '',
        },
    });

    useEffect(() => {
        if (!isEditMode || !id) return;

        const fetchTeacher = async () => {
            try {
                setLoading(true);
                setError(null);
                const data = await teacherApi.getById(id);
                reset({
                    full_name: data.full_name,
                    email: data.email || '',
                    phone: data.phone || '',
                    code: data.code || '',
                    is_school_teacher: data.is_school_teacher,
                    school_name: data.school_name || '',
                    employment_type: data.employment_type,
                    status: data.status,
                    notes: data.notes || '',
                });
            } catch (err) {
                setError(err instanceof Error ? err.message : 'Failed to load teacher');
                toast.error('Failed to load teacher details');
            } finally {
                setLoading(false);
            }
        };

        fetchTeacher();
    }, [id, isEditMode, reset]);

    const onSubmit = async (data: CreateTeacherRequest) => {
        try {
            setSubmitting(true);

            if (isEditMode && id) {
                await teacherApi.update(id, data);
                toast.success('Teacher updated successfully');
                navigate(`/app/admin/teachers/${id}`);
            } else {
                const newTeacher = await teacherApi.create(data);
                toast.success('Teacher created successfully');
                navigate(`/app/admin/teachers/${newTeacher.id}`);
            }
        } catch (err) {
            toast.error(err instanceof Error ? err.message : 'Failed to save teacher');
        } finally {
            setSubmitting(false);
        }
    };

    if (loading) {
        return (
            <Container maxWidth="md" sx={{ py: 4 }}>
                <Skeleton variant="rectangular" height={600} />
            </Container>
        );
    }

    if (error) {
        return (
            <Container maxWidth="md" sx={{ py: 4 }}>
                <Alert severity="error">{error}</Alert>
                <Button
                    startIcon={<ArrowBack />}
                    onClick={() => navigate('/app/admin/teachers')}
                    sx={{ mt: 2 }}
                >
                    Back to Teachers
                </Button>
            </Container>
        );
    }

    return (
        <Container maxWidth="md" sx={{ py: 4 }}>
            <Box sx={{ mb: 3 }}>
                <Button
                    startIcon={<ArrowBack />}
                    onClick={() => navigate(isEditMode ? `/app/admin/teachers/${id}` : '/app/admin/teachers')}
                    sx={{ mb: 2 }}
                >
                    Back
                </Button>

                <Typography variant="h4" component="h1" gutterBottom>
                    {isEditMode ? 'Edit Teacher' : 'Add New Teacher'}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    {isEditMode ? 'Update teacher information' : 'Fill in the details to create a new teacher'}
                </Typography>
            </Box>

            <Paper sx={{ p: 3 }}>
                <form onSubmit={handleSubmit(onSubmit)}>
                    <Grid container spacing={3}>
                        <Grid item xs={12}>
                            <Controller
                                name="full_name"
                                control={control}
                                render={({ field }) => (
                                    <TextField
                                        {...field}
                                        label="Full Name"
                                        fullWidth
                                        required
                                        error={!!errors.full_name}
                                        helperText={errors.full_name?.message}
                                    />
                                )}
                            />
                        </Grid>

                        <Grid item xs={12} sm={6}>
                            <Controller
                                name="email"
                                control={control}
                                render={({ field }) => (
                                    <TextField
                                        {...field}
                                        label="Email"
                                        type="email"
                                        fullWidth
                                        error={!!errors.email}
                                        helperText={errors.email?.message}
                                    />
                                )}
                            />
                        </Grid>

                        <Grid item xs={12} sm={6}>
                            <Controller
                                name="phone"
                                control={control}
                                render={({ field }) => (
                                    <TextField
                                        {...field}
                                        label="Phone"
                                        fullWidth
                                        error={!!errors.phone}
                                        helperText={errors.phone?.message}
                                    />
                                )}
                            />
                        </Grid>

                        <Grid item xs={12} sm={6}>
                            <Controller
                                name="code"
                                control={control}
                                render={({ field }) => (
                                    <TextField
                                        {...field}
                                        label="Teacher Code"
                                        fullWidth
                                        error={!!errors.code}
                                        helperText={errors.code?.message}
                                    />
                                )}
                            />
                        </Grid>

                        <Grid item xs={12} sm={6}>
                            <Controller
                                name="employment_type"
                                control={control}
                                render={({ field }) => (
                                    <FormControl fullWidth error={!!errors.employment_type}>
                                        <InputLabel>Employment Type</InputLabel>
                                        <Select {...field} label="Employment Type">
                                            <MenuItem value="FULL_TIME">Full Time</MenuItem>
                                            <MenuItem value="PART_TIME">Part Time</MenuItem>
                                        </Select>
                                    </FormControl>
                                )}
                            />
                        </Grid>

                        <Grid item xs={12} sm={6}>
                            <Controller
                                name="status"
                                control={control}
                                render={({ field }) => (
                                    <FormControl fullWidth error={!!errors.status}>
                                        <InputLabel>Status</InputLabel>
                                        <Select {...field} label="Status">
                                            <MenuItem value="ACTIVE">Active</MenuItem>
                                            <MenuItem value="INACTIVE">Inactive</MenuItem>
                                        </Select>
                                    </FormControl>
                                )}
                            />
                        </Grid>

                        <Grid item xs={12}>
                            <Controller
                                name="is_school_teacher"
                                control={control}
                                render={({ field }) => (
                                    <FormControlLabel
                                        control={<Switch {...field} checked={field.value} />}
                                        label="Is School Teacher"
                                    />
                                )}
                            />
                        </Grid>

                        <Grid item xs={12}>
                            <Controller
                                name="school_name"
                                control={control}
                                render={({ field }) => (
                                    <TextField
                                        {...field}
                                        label="School Name"
                                        fullWidth
                                        error={!!errors.school_name}
                                        helperText={errors.school_name?.message}
                                    />
                                )}
                            />
                        </Grid>

                        <Grid item xs={12}>
                            <Controller
                                name="notes"
                                control={control}
                                render={({ field }) => (
                                    <TextField
                                        {...field}
                                        label="Notes"
                                        fullWidth
                                        multiline
                                        rows={4}
                                        error={!!errors.notes}
                                        helperText={errors.notes?.message}
                                    />
                                )}
                            />
                        </Grid>

                        <Grid item xs={12}>
                            <Box sx={{ display: 'flex', gap: 2, justifyContent: 'flex-end' }}>
                                <Button
                                    variant="outlined"
                                    onClick={() =>
                                        navigate(isEditMode ? `/app/admin/teachers/${id}` : '/app/admin/teachers')
                                    }
                                    disabled={submitting}
                                >
                                    Cancel
                                </Button>
                                <Button
                                    type="submit"
                                    variant="contained"
                                    startIcon={<Save />}
                                    disabled={submitting}
                                >
                                    {submitting ? 'Saving...' : isEditMode ? 'Update Teacher' : 'Create Teacher'}
                                </Button>
                            </Box>
                        </Grid>
                    </Grid>
                </form>
            </Paper>
        </Container>
    );
};
