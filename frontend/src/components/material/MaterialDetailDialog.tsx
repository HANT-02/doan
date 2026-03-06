import {
    Alert,
    Box,
    Chip,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    Divider,
    Button,
    Stack,
    Tab,
    Tabs,
    Typography,
} from '@mui/material';
import { useState } from 'react';
import type { MaterialItem } from '@/api/materialApi';

interface MaterialDetailDialogProps {
    open: boolean;
    material: MaterialItem | null;
    onClose: () => void;
}

const severityColorMap: Record<string, 'success' | 'warning' | 'error' | 'info'> = {
    SAFE: 'success',
    WARNING: 'warning',
    DANGER: 'error',
};

export default function MaterialDetailDialog({ open, material, onClose }: MaterialDetailDialogProps) {
    const [tab, setTab] = useState(0);

    return (
        <Dialog open={open} onClose={onClose} fullWidth maxWidth="md">
            <DialogTitle>
                <Typography variant="h6" sx={{ fontWeight: 800 }}>
                    Chi tiết tài liệu
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    {material?.title || 'Tài liệu giảng dạy'}
                </Typography>
            </DialogTitle>
            <DialogContent dividers>
                {!material ? null : (
                    <Stack spacing={2.5}>
                        <Box>
                            <Stack direction="row" spacing={1} flexWrap="wrap">
                                <Chip label={material.status} color={material.status === 'APPROVED' ? 'success' : material.status === 'REJECTED' ? 'error' : 'warning'} />
                                {material.latest_label ? (
                                    <Chip
                                        label={`${material.latest_label.name} • ${material.latest_label.severity}`}
                                        color={severityColorMap[material.latest_label.severity] || 'info'}
                                        variant="outlined"
                                    />
                                ) : null}
                            </Stack>
                        </Box>

                        <Tabs value={tab} onChange={(_, value) => setTab(value)}>
                            <Tab label="Tổng quan" />
                            <Tab label="AI Audit" />
                            <Tab label="Phê duyệt" />
                        </Tabs>

                        {tab === 0 ? (
                            <Stack spacing={1.5}>
                                <Typography variant="body2"><strong>Tên file:</strong> {material.file_name}</Typography>
                                <Typography variant="body2"><strong>Loại file:</strong> {material.file_type}</Typography>
                                <Typography variant="body2"><strong>Đường dẫn:</strong> {material.file_path}</Typography>
                                <Typography variant="body2"><strong>Mô tả:</strong> {material.description || 'Chưa có mô tả'}</Typography>
                                <Typography variant="body2"><strong>Thời điểm tải lên:</strong> {new Date(material.uploaded_at).toLocaleString('vi-VN')}</Typography>
                            </Stack>
                        ) : null}

                        {tab === 1 ? (
                            <Stack spacing={2}>
                                {material.latest_audit ? (
                                    <>
                                        <Alert severity={severityColorMap[material.latest_label?.severity || 'WARNING'] || 'warning'}>
                                            {material.latest_audit.reasoning}
                                        </Alert>
                                        <Typography variant="body2"><strong>Độ tin cậy:</strong> {(material.latest_audit.confidence_score * 100).toFixed(0)}%</Typography>
                                        <Typography variant="body2"><strong>Provider:</strong> {material.latest_audit.provider}</Typography>
                                        <Typography variant="body2"><strong>Detected issues:</strong> {material.latest_audit.detected_issues.join(', ') || 'Không có'}</Typography>
                                        <Divider />
                                        <Typography variant="subtitle2" sx={{ fontWeight: 700 }}>
                                            Raw OCR text
                                        </Typography>
                                        <Box sx={{ p: 2, borderRadius: 2, bgcolor: 'grey.50', maxHeight: 220, overflow: 'auto' }}>
                                            <Typography variant="body2" sx={{ whiteSpace: 'pre-wrap' }}>
                                                {material.latest_audit.raw_ocr_text || 'Chưa có nội dung OCR'}
                                            </Typography>
                                        </Box>
                                    </>
                                ) : (
                                    <Alert severity="info">Tài liệu này chưa có kết quả AI audit.</Alert>
                                )}
                            </Stack>
                        ) : null}

                        {tab === 2 ? (
                            material.last_decision ? (
                                <Stack spacing={1.5}>
                                    <Typography variant="body2"><strong>Kết quả:</strong> {material.last_decision.approved ? 'Phê duyệt' : 'Từ chối'}</Typography>
                                    <Typography variant="body2"><strong>Officer:</strong> {material.last_decision.compliance_officer_id}</Typography>
                                    <Typography variant="body2"><strong>Ghi chú:</strong> {material.last_decision.notes || 'Không có'}</Typography>
                                    <Typography variant="body2"><strong>Lý do từ chối:</strong> {material.last_decision.reject_reason || 'Không có'}</Typography>
                                    <Typography variant="body2"><strong>Thời điểm:</strong> {new Date(material.last_decision.decided_at).toLocaleString('vi-VN')}</Typography>
                                </Stack>
                            ) : (
                                <Alert severity="info">Chưa có quyết định phê duyệt cho tài liệu này.</Alert>
                            )
                        ) : null}
                    </Stack>
                )}
            </DialogContent>
            <DialogActions sx={{ px: 3, py: 2 }}>
                <Button onClick={onClose}>Đóng</Button>
            </DialogActions>
        </Dialog>
    );
}
