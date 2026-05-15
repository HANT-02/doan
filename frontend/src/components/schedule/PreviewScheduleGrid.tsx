import { Alert, Box, Chip, Paper, Stack, Typography } from '@mui/material';
import { format, isSameDay, parseISO, startOfWeek, addDays } from 'date-fns';
import { vi } from 'date-fns/locale';
import type { ReactNode } from 'react';

const dayLabels = ['Thứ Hai', 'Thứ Ba', 'Thứ Tư', 'Thứ Năm', 'Thứ Sáu', 'Thứ Bảy', 'Chủ Nhật'];

type ScheduleLikeItem = {
    id: string;
    date_start: string;
    date_end: string;
    shift?: {
        id?: string;
        name?: string;
        start_time?: string;
        end_time?: string;
    } | null;
};

type PreviewScheduleGridProps<T extends ScheduleLikeItem> = {
    title: string;
    weekStart: Date;
    items: T[];
    isFetching?: boolean;
    emptyLabel?: string;
    renderItem: (item: T) => ReactNode;
};

type CalendarRow = {
    key: string;
    shiftLabel: string;
    timeLabel: string;
};

function buildCalendarRowKey(item: ScheduleLikeItem) {
    if (item.shift?.id) {
        return `shift:${item.shift.id}`;
    }

    const start = format(parseISO(item.date_start), 'HH:mm');
    const end = format(parseISO(item.date_end), 'HH:mm');
    return `time:${start}-${end}`;
}

function buildCalendarRows<T extends ScheduleLikeItem>(items: T[]) {
    const seen = new Map<string, CalendarRow>();

    items.forEach((item) => {
        const key = buildCalendarRowKey(item);
        if (seen.has(key)) {
            return;
        }

        const shiftLabel = item.shift?.name || 'Ca học';
        const timeLabel = item.shift?.start_time && item.shift?.end_time
            ? `${item.shift.start_time} - ${item.shift.end_time}`
            : `${format(parseISO(item.date_start), 'HH:mm')} - ${format(parseISO(item.date_end), 'HH:mm')}`;

        seen.set(key, { key, shiftLabel, timeLabel });
    });

    return Array.from(seen.values()).sort((left, right) => left.timeLabel.localeCompare(right.timeLabel));
}

export default function PreviewScheduleGrid<T extends ScheduleLikeItem>({
    title,
    weekStart,
    items,
    isFetching = false,
    emptyLabel = 'Chưa có buổi học trong tuần này.',
    renderItem,
}: PreviewScheduleGridProps<T>) {
    const week = Array.from({ length: 7 }, (_, index) => addDays(startOfWeek(weekStart, { weekStartsOn: 1 }), index));
    const rows = buildCalendarRows(items);

    return (
        <Paper variant="outlined" sx={{ borderRadius: 3, overflow: 'hidden' }}>
            <Box sx={{ px: 2, py: 1.5, backgroundColor: '#f8fafc', borderBottom: '1px solid rgba(15,23,42,0.08)' }}>
                <Stack direction="row" justifyContent="space-between" alignItems="center" spacing={1}>
                    <Typography variant="subtitle2" sx={{ fontWeight: 700 }}>
                        {title}
                    </Typography>
                    {isFetching ? <Chip size="small" label="Đang tải..." color="info" variant="outlined" /> : null}
                </Stack>
            </Box>

            {rows.length === 0 ? (
                <Box sx={{ p: 2 }}>
                    <Alert severity="info">{isFetching ? 'Đang tải lịch...' : emptyLabel}</Alert>
                </Box>
            ) : (
                <Box sx={{ overflowX: 'auto' }}>
                    <Box
                        sx={{
                            minWidth: 1120,
                            display: 'grid',
                            gridTemplateColumns: '160px repeat(7, minmax(0, 1fr))',
                        }}
                    >
                        <Box sx={{ p: 1.5, borderBottom: '1px solid rgba(15,23,42,0.08)', backgroundColor: '#f8fafc' }}>
                            <Typography variant="subtitle2" sx={{ fontWeight: 700 }}>
                                Ca học
                            </Typography>
                        </Box>
                        {week.map((day, dayIndex) => (
                            <Box
                                key={`header-${day.toISOString()}`}
                                sx={{
                                    p: 1.5,
                                    borderLeft: '1px solid rgba(15,23,42,0.08)',
                                    borderBottom: '1px solid rgba(15,23,42,0.08)',
                                    backgroundColor: '#f8fafc',
                                }}
                            >
                                <Typography variant="subtitle2" sx={{ fontWeight: 700 }}>
                                    {dayLabels[dayIndex]}
                                </Typography>
                                <Typography variant="caption" color="text.secondary">
                                    {format(day, 'dd/MM/yyyy', { locale: vi })}
                                </Typography>
                            </Box>
                        ))}

                        {rows.map((row) => (
                            <Box key={row.key} sx={{ display: 'contents' }}>
                                <Box
                                    sx={{
                                        p: 1.5,
                                        borderBottom: '1px solid rgba(15,23,42,0.08)',
                                        backgroundColor: '#fcfcfd',
                                    }}
                                >
                                    <Typography variant="body2" sx={{ fontWeight: 700 }}>
                                        {row.shiftLabel}
                                    </Typography>
                                    <Typography variant="caption" color="text.secondary">
                                        {row.timeLabel}
                                    </Typography>
                                </Box>
                                {week.map((day) => {
                                    const rowItems = items.filter((item) =>
                                        isSameDay(parseISO(item.date_start), day) && buildCalendarRowKey(item) === row.key,
                                    );

                                    return (
                                        <Box
                                            key={`${row.key}-${day.toISOString()}`}
                                            sx={{
                                                p: 1,
                                                minHeight: 132,
                                                borderLeft: '1px solid rgba(15,23,42,0.08)',
                                                borderBottom: '1px solid rgba(15,23,42,0.08)',
                                            }}
                                        >
                                            <Stack spacing={0.75}>
                                                {rowItems.map((item) => renderItem(item))}
                                            </Stack>
                                        </Box>
                                    );
                                })}
                            </Box>
                        ))}
                    </Box>
                </Box>
            )}
        </Paper>
    );
}
