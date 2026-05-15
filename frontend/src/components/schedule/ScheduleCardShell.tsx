import { Box, Button, Chip, Paper, Stack, Typography, type ChipProps } from '@mui/material';
import type { ReactNode } from 'react';

type ScheduleCardChip = {
    label: string;
    color?: ChipProps['color'];
    variant?: ChipProps['variant'];
};

const palette = [
    { bg: 'rgba(14, 116, 144, 0.10)', border: '#0ea5e9' },
    { bg: 'rgba(5, 150, 105, 0.10)', border: '#10b981' },
    { bg: 'rgba(234, 88, 12, 0.10)', border: '#f97316' },
    { bg: 'rgba(190, 24, 93, 0.10)', border: '#ec4899' },
    { bg: 'rgba(79, 70, 229, 0.10)', border: '#6366f1' },
    { bg: 'rgba(161, 98, 7, 0.10)', border: '#eab308' },
];

function getTone(seed: string) {
    const numericSeed = seed
        .split('')
        .reduce((sum, char) => sum + char.charCodeAt(0), 0);
    return palette[numericSeed % palette.length];
}

type ScheduleCardShellProps = {
    seed: string;
    title: string;
    subtitle?: string;
    lines?: string[];
    chips?: ScheduleCardChip[];
    note?: string;
    actionLabel?: string;
    onActionClick?: () => void;
    onClick?: () => void;
    muted?: boolean;
    warning?: boolean;
    children?: ReactNode;
};

export default function ScheduleCardShell({
    seed,
    title,
    subtitle,
    lines = [],
    chips = [],
    note,
    actionLabel,
    onActionClick,
    onClick,
    muted = false,
    warning = false,
    children,
}: ScheduleCardShellProps) {
    const tone = getTone(seed);

    return (
        <Paper
            variant="outlined"
            onClick={onClick}
            sx={{
                p: 1.25,
                borderRadius: 2.5,
                cursor: onClick ? 'pointer' : 'default',
                borderLeft: `4px solid ${muted ? 'rgba(15,23,42,0.24)' : tone.border}`,
                backgroundColor: muted ? 'rgba(15,23,42,0.04)' : tone.bg,
                borderColor: warning ? 'warning.light' : undefined,
                transition: onClick ? 'transform 0.18s ease, box-shadow 0.18s ease' : undefined,
                '&:hover': onClick
                    ? {
                        transform: 'translateY(-1px)',
                        boxShadow: 2,
                    }
                    : undefined,
            }}
        >
            <Stack spacing={0.75}>
                <Box sx={{ minWidth: 0 }}>
                    <Typography variant="body2" sx={{ fontWeight: 700 }} noWrap>
                        {title}
                    </Typography>
                    {subtitle ? (
                        <Typography variant="caption" color="text.secondary">
                            {subtitle}
                        </Typography>
                    ) : null}
                </Box>

                {lines.map((line) => (
                    <Typography key={line} variant="caption" color="text.secondary" noWrap>
                        {line}
                    </Typography>
                ))}

                {chips.length > 0 ? (
                    <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                        {chips.map((chip) => (
                            <Chip
                                key={`${title}-${chip.label}`}
                                size="small"
                                label={chip.label}
                                color={chip.color}
                                variant={chip.variant || 'outlined'}
                            />
                        ))}
                    </Stack>
                ) : null}

                {note ? (
                    <Typography variant="caption" color="text.secondary" noWrap>
                        {note}
                    </Typography>
                ) : null}

                {children}

                {actionLabel && onActionClick ? (
                    <Button
                        variant="text"
                        size="small"
                        onClick={(event) => {
                            event.stopPropagation();
                            onActionClick();
                        }}
                        sx={{ alignSelf: 'flex-start', px: 0 }}
                    >
                        {actionLabel}
                    </Button>
                ) : null}
            </Stack>
        </Paper>
    );
}
