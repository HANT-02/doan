import {
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Paper,
    IconButton,
    Chip,
    Typography,
    Box,
    Skeleton,
} from '@mui/material';
import { Visibility, Edit, Delete } from '@mui/icons-material';
import { type Teacher } from '@/api/teacherApi';
import { getStatusColor, getEmploymentTypeLabel, formatDate } from '@/utils/teacherHelpers';

interface TeacherListTableProps {
    teachers: Teacher[];
    loading?: boolean;
    onView?: (teacher: Teacher) => void;
    onEdit?: (teacher: Teacher) => void;
    onDelete?: (teacher: Teacher) => void;
    showActions?: boolean;
    isAdmin?: boolean;
}

export const TeacherListTable = ({
    teachers,
    loading = false,
    onView,
    onEdit,
    onDelete,
    showActions = true,
    isAdmin = false,
}: TeacherListTableProps) => {
    if (loading) {
        return (
            <TableContainer component={Paper}>
                <Table>
                    <TableHead>
                        <TableRow>
                            <TableCell>Code</TableCell>
                            <TableCell>Name</TableCell>
                            <TableCell>Email</TableCell>
                            <TableCell>Phone</TableCell>
                            <TableCell>Employment</TableCell>
                            <TableCell>Status</TableCell>
                            {showActions && <TableCell align="right">Actions</TableCell>}
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {[1, 2, 3, 4, 5].map((i) => (
                            <TableRow key={i}>
                                <TableCell><Skeleton /></TableCell>
                                <TableCell><Skeleton /></TableCell>
                                <TableCell><Skeleton /></TableCell>
                                <TableCell><Skeleton /></TableCell>
                                <TableCell><Skeleton /></TableCell>
                                <TableCell><Skeleton /></TableCell>
                                {showActions && <TableCell><Skeleton /></TableCell>}
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </TableContainer>
        );
    }

    if (teachers.length === 0) {
        return (
            <Paper sx={{ p: 4, textAlign: 'center' }}>
                <Typography variant="h6" color="text.secondary">
                    No teachers found
                </Typography>
                <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                    Try adjusting your filters or add a new teacher
                </Typography>
            </Paper>
        );
    }

    return (
        <TableContainer component={Paper}>
            <Table>
                <TableHead>
                    <TableRow>
                        <TableCell>Code</TableCell>
                        <TableCell>Name</TableCell>
                        <TableCell>Email</TableCell>
                        <TableCell>Phone</TableCell>
                        <TableCell>Employment</TableCell>
                        <TableCell>Status</TableCell>
                        <TableCell>Created</TableCell>
                        {showActions && <TableCell align="right">Actions</TableCell>}
                    </TableRow>
                </TableHead>
                <TableBody>
                    {teachers.map((teacher) => (
                        <TableRow key={teacher.id} hover>
                            <TableCell>{teacher.code || '-'}</TableCell>
                            <TableCell>
                                <Typography variant="body2" fontWeight="medium">
                                    {teacher.full_name}
                                </Typography>
                                {teacher.is_school_teacher && (
                                    <Typography variant="caption" color="text.secondary">
                                        {teacher.school_name || 'School Teacher'}
                                    </Typography>
                                )}
                            </TableCell>
                            <TableCell>{teacher.email || '-'}</TableCell>
                            <TableCell>{teacher.phone || '-'}</TableCell>
                            <TableCell>{getEmploymentTypeLabel(teacher.employment_type)}</TableCell>
                            <TableCell>
                                <Chip
                                    label={teacher.status}
                                    color={getStatusColor(teacher.status)}
                                    size="small"
                                />
                            </TableCell>
                            <TableCell>{formatDate(teacher.created_at)}</TableCell>
                            {showActions && (
                                <TableCell align="right">
                                    <Box sx={{ display: 'flex', gap: 0.5, justifyContent: 'flex-end' }}>
                                        {onView && (
                                            <IconButton
                                                size="small"
                                                onClick={() => onView(teacher)}
                                                title="View details"
                                            >
                                                <Visibility fontSize="small" />
                                            </IconButton>
                                        )}
                                        {isAdmin && onEdit && (
                                            <IconButton
                                                size="small"
                                                onClick={() => onEdit(teacher)}
                                                title="Edit teacher"
                                            >
                                                <Edit fontSize="small" />
                                            </IconButton>
                                        )}
                                        {isAdmin && onDelete && (
                                            <IconButton
                                                size="small"
                                                onClick={() => onDelete(teacher)}
                                                title="Delete teacher"
                                                color="error"
                                            >
                                                <Delete fontSize="small" />
                                            </IconButton>
                                        )}
                                    </Box>
                                </TableCell>
                            )}
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </TableContainer>
    );
};
