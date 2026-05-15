import { useEffect, useMemo, useState } from 'react';
import {
    Alert,
    Box,
    Button,
    Chip,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    Paper,
    Stack,
    TextField,
    Typography,
} from '@mui/material';
import { PersonSearchRounded } from '@mui/icons-material';
import { toast } from 'sonner';

import {
    useAssignSubstituteMutation,
    useLazySuggestSubstituteQuery,
    type SubstituteSuggestion,
} from '@/api/schedulingApi';
import { getApiErrorMessage } from '@/utils/apiError';

interface SuggestSubstituteModalProps {
    isOpen: boolean;
    onClose: () => void;
    lessonId: string | null;
    lessonLabel?: string;
}

export const SuggestSubstituteModal = ({
    isOpen,
    onClose,
    lessonId,
    lessonLabel,
}: SuggestSubstituteModalProps) => {
    const [suggest, { data: suggestionsData, isLoading, isFetching, error }] = useLazySuggestSubstituteQuery();
    const [assignSubstitute, { isLoading: isAssigning }] = useAssignSubstituteMutation();
    const [reason, setReason] = useState('Thay thế giáo viên đột xuất');

    useEffect(() => {
        if (isOpen && lessonId) {
            void suggest(lessonId);
        }
    }, [isOpen, lessonId, suggest]);

    const suggestions = useMemo(() => suggestionsData?.data || [], [suggestionsData?.data]);

    const handleAssign = async (teacher: SubstituteSuggestion) => {
        if (!lessonId) {
            return;
        }
        try {
            await assignSubstitute({
                lessonId,
                data: {
                    new_teacher_id: teacher.teacher_id,
                    reason: reason.trim() || 'Thay thế giáo viên đột xuất',
                },
            }).unwrap();
            toast.success(`Đã gán ${teacher.teacher_name} làm giáo viên dạy thay`);
            onClose();
        } catch (assignError) {
            toast.error(getApiErrorMessage(assignError, 'Không thể gán giáo viên dạy thay'));
        }
    };

    return (
        <Dialog open={isOpen} onClose={onClose} fullWidth maxWidth="md">
            <DialogTitle>Đề xuất giáo viên dạy thay</DialogTitle>
            <DialogContent>
                <Stack spacing={2.5} sx={{ pt: 1 }}>
                    <Alert severity="info">
                        {lessonLabel
                            ? `Đang tìm giáo viên thay cho ${lessonLabel}.`
                            : 'Hệ thống sẽ ưu tiên giáo viên đúng kỹ năng, đang trống lịch, ít quá tải và di chuyển kịp giữa các cơ sở.'}
                    </Alert>

                    <TextField
                        label="Lý do thay thế"
                        value={reason}
                        onChange={(event) => setReason(event.target.value)}
                        fullWidth
                        helperText="Lý do này sẽ được lưu cùng thao tác dạy thay để tiện audit."
                    />

                    {error ? (
                        <Alert severity="error">
                            {getApiErrorMessage(error, 'Không tải được danh sách giáo viên dạy thay.')}
                        </Alert>
                    ) : null}

                    {isLoading || isFetching ? (
                        <Typography variant="body2" color="text.secondary">
                            Đang quét giáo viên phù hợp...
                        </Typography>
                    ) : null}

                    {!isLoading && !isFetching && !suggestions.length && !error ? (
                        <Alert severity="warning">
                            Không tìm thấy giáo viên nào thỏa đồng thời kỹ năng, lịch trống và travel gap cho ca này.
                        </Alert>
                    ) : null}

                    {suggestions.map((teacher) => (
                        <Paper key={teacher.teacher_id} variant="outlined" sx={{ p: 2, borderRadius: 3 }}>
                            <Stack direction={{ xs: 'column', md: 'row' }} spacing={2} justifyContent="space-between">
                                <Box>
                                    <Stack direction="row" spacing={1} alignItems="center" flexWrap="wrap" useFlexGap>
                                        <Typography variant="subtitle1" sx={{ fontWeight: 800 }}>
                                            {teacher.teacher_name}
                                        </Typography>
                                        <Chip size="small" label={teacher.teacher_code} variant="outlined" />
                                        <Chip
                                            size="small"
                                            color={teacher.score >= 80 ? 'success' : teacher.score >= 60 ? 'warning' : 'default'}
                                            label={`${teacher.score} điểm`}
                                        />
                                    </Stack>
                                    <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap sx={{ mt: 1 }}>
                                        {teacher.match_reasons.map((item) => (
                                            <Chip
                                                key={`${teacher.teacher_id}-${item}`}
                                                size="small"
                                                variant="outlined"
                                                color={item.toLowerCase().includes('thiếu') ? 'error' : 'info'}
                                                label={item}
                                            />
                                        ))}
                                    </Stack>
                                </Box>

                                <Button
                                    variant="contained"
                                    startIcon={<PersonSearchRounded />}
                                    disabled={isAssigning || !teacher.is_available}
                                    onClick={() => void handleAssign(teacher)}
                                >
                                    {teacher.is_available ? 'Gán dạy thay' : 'Không khả dụng'}
                                </Button>
                            </Stack>
                        </Paper>
                    ))}
                </Stack>
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose}>Đóng</Button>
            </DialogActions>
        </Dialog>
    );
};
