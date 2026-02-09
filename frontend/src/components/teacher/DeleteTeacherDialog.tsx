import {
    Dialog,
    DialogTitle,
    DialogContent,
    DialogContentText,
    DialogActions,
    Button,
} from '@mui/material';
import { Warning } from '@mui/icons-material';

interface DeleteTeacherDialogProps {
    open: boolean;
    teacherName: string;
    onClose: () => void;
    onConfirm: () => void;
    loading?: boolean;
}

export const DeleteTeacherDialog = ({
    open,
    teacherName,
    onClose,
    onConfirm,
    loading = false,
}: DeleteTeacherDialogProps) => {
    return (
        <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
            <DialogTitle sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <Warning color="warning" />
                Delete Teacher
            </DialogTitle>
            <DialogContent>
                <DialogContentText>
                    Are you sure you want to delete teacher <strong>{teacherName}</strong>?
                </DialogContentText>
                <DialogContentText sx={{ mt: 1 }}>
                    This action will soft delete the teacher. The teacher will no longer appear in the
                    list but can be restored by an administrator if needed.
                </DialogContentText>
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose} disabled={loading}>
                    Cancel
                </Button>
                <Button onClick={onConfirm} color="error" variant="contained" disabled={loading}>
                    {loading ? 'Deleting...' : 'Delete'}
                </Button>
            </DialogActions>
        </Dialog>
    );
};
