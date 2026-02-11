import ConfirmDialog from '@/components/common/ConfirmDialog';

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
        <ConfirmDialog
            open={open}
            title="Xác nhận xóa Giáo viên"
            message={`Bạn có chắc chắn muốn xóa giáo viên ${teacherName} không? Hành động này sẽ xóa mềm giáo viên. Giáo viên sẽ không còn xuất hiện trong danh sách nhưng có thể được khôi phục bởi quản trị viên nếu cần.`}
            onConfirm={onConfirm}
            onClose={onClose}
            confirmText="Xóa"
            cancelText="Hủy"
            loading={loading}
        />
    );
};
