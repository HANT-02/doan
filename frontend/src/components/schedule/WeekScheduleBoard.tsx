import { Alert, Box, Chip, Paper, Stack, Typography } from '@mui/material';
import { format, isSameDay } from 'date-fns';
import { vi } from 'date-fns/locale';
import type { ReactNode } from 'react';

const weekLabels = ['Thứ Hai', 'Thứ Ba', 'Thứ Tư', 'Thứ Năm', 'Thứ Sáu', 'Thứ Bảy', 'Chủ Nhật'];

export type WeekScheduleDay<T> = {
    date: Date;
    items: T[];
};

type WeekScheduleBoardProps<T> = {
    title: string;
    subtitle?: string;
    days: WeekScheduleDay<T>[];
    isFetching?: boolean;
    emptyLabel?: string;
    minDayHeight?: number;
    renderItem: (item: T) => ReactNode;
};

export default function WeekScheduleBoard<T>({
    title,
    subtitle,
    days,
    isFetching = false,
    emptyLabel = '',
    minDayHeight = 260,
    renderItem,
}: WeekScheduleBoardProps<T>) {
    return (
        <Paper variant="outlined" sx={{ borderRadius: 3, overflow: 'hidden' }}>
            <Box
                sx={{
                    px: 2.5,
                    py: 2,
                    borderBottom: '1px solid rgba(15,23,42,0.08)',
                    background: 'linear-gradient(135deg, rgba(14,165,233,0.08), rgba(16,185,129,0.05))',
                }}
            >
                <Stack
                    direction={{ xs: 'column', md: 'row' }}
                    justifyContent="space-between"
                    alignItems={{ xs: 'flex-start', md: 'center' }}
                    spacing={1}
                >
                    <Box>
                        <Typography variant="h6" sx={{ fontWeight: 700 }}>
                            {title}
                        </Typography>
                        {subtitle ? (
                            <Typography variant="body2" color="text.secondary">
                                {subtitle}
                            </Typography>
                        ) : null}
                    </Box>
                    {isFetching ? (
                        <Chip size="small" label="Đang tải..." color="info" variant="outlined" />
                    ) : null}
                </Stack>
            </Box>

            <Box
                sx={{
                    display: 'grid',
                    gridTemplateColumns: { xs: '1fr', lg: 'repeat(7, minmax(0, 1fr))' },
                    borderBottom: '1px solid rgba(15,23,42,0.08)',
                    backgroundColor: '#f8fafc',
                }}
            >
                {weekLabels.map((label) => (
                    <Box key={label} sx={{ px: 1.5, py: 1.25, borderLeft: { lg: '1px solid rgba(15,23,42,0.08)' } }}>
                        <Typography variant="subtitle2" sx={{ fontWeight: 700 }}>
                            {label}
                        </Typography>
                    </Box>
                ))}
            </Box>

            <Box
                sx={{
                    display: 'grid',
                    gridTemplateColumns: { xs: '1fr', lg: 'repeat(7, minmax(0, 1fr))' },
                }}
            >
                {days.map((day, index) => (
                    <Box
                        key={day.date.toISOString()}
                        sx={{
                            minHeight: minDayHeight,
                            p: 1.25,
                            borderLeft: { xs: 'none', lg: index === 0 ? 'none' : '1px solid rgba(15,23,42,0.08)' },
                            borderTop: '1px solid rgba(15,23,42,0.08)',
                        }}
                    >
                        <Stack spacing={1}>
                            <Stack
                                direction={{ xs: 'column', xl: 'row' }}
                                justifyContent="space-between"
                                alignItems={{ xs: 'flex-start', xl: 'center' }}
                                spacing={0.5}
                            >
                                <Box>
                                    <Typography
                                        variant="body2"
                                        sx={{ fontWeight: 700, color: isSameDay(day.date, new Date()) ? 'primary.main' : 'text.primary' }}
                                    >
                                        {format(day.date, 'dd/MM', { locale: vi })}
                                    </Typography>
                                </Box>
                                <Chip size="small" variant="outlined" label={`${day.items.length} buổi`} />
                            </Stack>

                            {day.items.length > 0 ? (
                                <Stack spacing={1}>
                                    {day.items.map((item) => renderItem(item))}
                                </Stack>
                            ) : emptyLabel ? (
                                <Alert severity="info">{isFetching ? 'Đang tải lịch...' : emptyLabel}</Alert>
                            ) : (
                                <Box sx={{ minHeight: 28 }} />
                            )}
                        </Stack>
                    </Box>
                ))}
            </Box>
        </Paper>
    );
}
