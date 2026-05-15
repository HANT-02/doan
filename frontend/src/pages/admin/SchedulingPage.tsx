import { useEffect, useMemo, useState, type ReactNode } from 'react';
import { Controller, useForm, useWatch } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import {
    Autocomplete,
    Box,
    Button,
    Chip,
    MenuItem,
    Paper,
    Skeleton,
    Stack,
    Step,
    StepLabel,
    Stepper,
    TextField,
    Tooltip,
    Typography,
} from '@mui/material';
import {
    CalendarMonthRounded,
    EditCalendarRounded,
    PlayArrowRounded,
    RefreshRounded,
    RuleRounded,
    TuneRounded,
} from '@mui/icons-material';
import { addDays, eachWeekOfInterval, endOfWeek, format, parseISO, startOfWeek } from 'date-fns';
import { vi } from 'date-fns/locale';
import { toast } from 'sonner';

import { useGetClassesQuery } from '@/api/classApi';
import {
    useCommitSchedulingPreviewMutation,
    useLazyGetLatestSchedulingPreviewQuery,
    useLazyGetSchedulingPreviewQuery,
    usePreviewSchedulingMutation,
    type SchedulingAssignment,
    type SchedulingCandidateOption,
    type SchedulingConflict,
    type SchedulingExistingLesson,
    type SchedulingPreview,
} from '@/api/schedulingApi';
import { useGetTeachersQuery } from '@/api/teacherApi';
import PageHeader from '@/components/common/PageHeader';
import { getApiErrorMessage } from '@/utils/apiError';

const schedulingSchema = z.object({
    mode: z.enum(['class', 'teacher']),
    expected_start_date: z.string().optional(),
    class_ids: z.array(z.string()).optional(),
    teacher_ids: z.array(z.string()).optional(),
}).refine((values) => {
    if (values.mode === 'class') {
        return (values.class_ids?.length || 0) > 0;
    }
    return true;
}, {
    message: 'Vui lòng chọn ít nhất một lớp để xếp lịch',
    path: ['class_ids'],
}).refine((values) => {
    if (values.mode === 'teacher') {
        return (values.teacher_ids?.length || 0) > 0;
    }
    return true;
}, {
    message: 'Vui lòng chọn ít nhất một giáo viên để lấy các lớp phụ trách',
    path: ['teacher_ids'],
});

type SchedulingFormValues = z.infer<typeof schedulingSchema>;

type SessionDraft = {
    variableId: string;
    classId: string;
    classCode: string;
    className: string;
    sessionIndex: number;
    sessionTotal: number;
    teacherId?: string;
    teacherLabel?: string;
    baseAssignment?: SchedulingAssignment;
    resolvedAssignment?: SchedulingAssignment;
    candidateOptions: SchedulingCandidateOption[];
    baseConflicts: SchedulingConflict[];
    derivedConflictMessages: string[];
    conflictTags: string[];
};

type DerivedPreviewState = {
    sessions: SessionDraft[];
    calendarAssignments: SchedulingAssignment[];
    existingLessons: SchedulingExistingLesson[];
    existingLessonConflictTags: Record<string, string[]>;
    unresolvedSessions: SessionDraft[];
    remainingConflicts: SchedulingConflict[];
    summary: {
        requestedClasses: number;
        requestedSessions: number;
        scheduledLessons: number;
        unscheduledLessons: number;
        conflictCount: number;
        baseSoftScore: number;
        softScore: number;
        scoreDelta: number;
        manualAdjustmentCount: number;
        manualAdjustmentLimit: number;
        excessiveManualAdjustment: boolean;
    };
    status: 'FAILED' | 'PARTIAL' | 'COMPLETED';
    manualAssignmentPayload: Array<{
        variable_id: string;
        option_key: string;
    }>;
};

type CalendarRow = {
    key: string;
    shiftLabel: string;
    timeLabel: string;
    startTime: string;
    endTime: string;
};

const previewStatusLabel = (value?: string) => {
    switch (value) {
        case 'FAILED':
            return 'Thất bại';
        case 'PARTIAL':
            return 'Một phần';
        case 'COMPLETED':
            return 'Hoàn tất';
        default:
            return value || 'Chưa rõ';
    }
};

const previewSteps = ['Chọn đối tượng', 'Xác nhận cấu hình', 'Xem trước và xử lý'];
const dayLabels = ['Thứ Hai', 'Thứ Ba', 'Thứ Tư', 'Thứ Năm', 'Thứ Sáu', 'Thứ Bảy', 'Chủ Nhật'];

const formatDateInput = (value: Date) => {
    const year = value.getFullYear();
    const month = `${value.getMonth() + 1}`.padStart(2, '0');
    const day = `${value.getDate()}`.padStart(2, '0');
    return `${year}-${month}-${day}`;
};

const fallbackDateFrom = formatDateInput(new Date());

type MinimalScheduleClass = {
    start_date?: string;
    end_date?: string;
};

const deriveDefaultStartDate = (classes: MinimalScheduleClass[]) => {
    const dates = classes
        .map((item) => item.start_date)
        .filter((value): value is string => !!value)
        .sort();
    return dates[0]?.slice(0, 10) || fallbackDateFrom;
};

const formatDateTime = (value: string) =>
    format(parseISO(value), 'dd/MM/yyyy HH:mm', { locale: vi });

const formatCompactDateTime = (value: string) =>
    format(parseISO(value), 'dd/MM HH:mm', { locale: vi });

const getConflictSeverity = (type: string): 'error' | 'warning' | 'info' => {
    switch (type) {
        case 'NO_CLASS_INPUT':
        case 'NO_ACTIVE_ROOM':
        case 'NO_ACTIVE_SHIFT':
        case 'NO_VALID_DATE_RANGE':
        case 'SKILL_MISMATCH':
            return 'error';
        case 'INSUFFICIENT_ENROLLMENT':
        case 'EXCESSIVE_MANUAL_ADJUSTMENT':
        case 'ASSIGNMENT_CONFLICT':
        case 'SYSTEM_LESSON_CONFLICT':
        case 'NO_DOMAIN':
        case 'MISSING_CLASS_SCHEDULE':
        case 'INSUFFICIENT_SCHEDULE_SLOTS':
        case 'ROOM_CAPACITY_BLOCK':
        case 'PREFERRED_ROOM_UNAVAILABLE':
        case 'CLASS_SCHEDULE_ROOM_UNAVAILABLE':
        case 'CLASS_SCHEDULE_NO_SLOT':
            return 'warning';
        default:
            return 'info';
    }
};

const getConflictScopeLabel = (classCode?: string, className?: string, sessionIndex?: number, sessionTotal?: number) => {
    const sessionLabel = sessionIndex && sessionTotal ? ` • Buổi ${sessionIndex}/${sessionTotal}` : '';
    if (classCode) {
        return className ? `${classCode} - ${className}${sessionLabel}` : `${classCode}${sessionLabel}`;
    }

    return className ? `${className}${sessionLabel}` : 'Tổng quan lịch xem trước';
};

const uniqueTags = (items: string[]) => Array.from(new Set(items));

const buildConflictTags = (type: string, message: string) => {
    const content = `${type} ${message}`.toLowerCase();
    const tags: string[] = [];

    if (
        type === 'ASSIGNMENT_CONFLICT' ||
        type === 'SYSTEM_LESSON_CONFLICT' ||
        content.includes('trùng')
    ) {
        tags.push('Trùng ngày/ca');
    }
    if (content.includes('trùng lớp')) {
        tags.push('Trùng lớp');
    }
    if (content.includes('trùng giáo viên')) {
        tags.push('Trùng giáo viên');
    }
    if (content.includes('trùng học sinh')) {
        tags.push('Trùng học sinh');
    }
    if (content.includes('trùng phòng')) {
        tags.push('Trùng phòng');
    }
    if (type === 'INSUFFICIENT_SCHEDULE_SLOTS' || content.includes('chưa đủ') || content.includes('không có slot')) {
        tags.push('Thiếu slot');
    }
    if (type === 'MISSING_CLASS_SCHEDULE' || content.includes('lịch tuần')) {
        tags.push('Thiếu lịch tuần');
    }
    if (type === 'INSUFFICIENT_ENROLLMENT' || content.includes('sĩ số') || content.includes('điều kiện mở lớp')) {
        tags.push('Thiếu sĩ số');
    }
    if (type === 'EXCESSIVE_MANUAL_ADJUSTMENT' || content.includes('chỉnh tay') || content.includes('vượt ngưỡng')) {
        tags.push('Vượt ngưỡng chỉnh tay');
    }
    if (type === 'SKILL_MISMATCH' || content.includes('kỹ năng') || content.includes('chứng chỉ')) {
        tags.push('Sai kỹ năng');
    }

    return uniqueTags(tags);
};

const getAllowedManualAdjustmentLimit = (requestedSessions: number) =>
    Math.max(3, Math.ceil(requestedSessions * 0.35));

const renderConflictChips = (tags: string[]) => {
    if (tags.length === 0) {
        return null;
    }

    return (
        <Stack direction="row" spacing={0.75} flexWrap="wrap" useFlexGap>
            {tags.map((tag) => (
                <Chip
                    key={tag}
                    size="small"
                    color={tag === 'Thiếu slot' || tag === 'Thiếu lịch tuần' ? 'warning' : 'error'}
                    variant="outlined"
                    label={tag}
                />
            ))}
        </Stack>
    );
};

const collectSessionIssueMessages = (session: SessionDraft) => uniqueTags([
    ...session.baseConflicts.map((conflict) => conflict.message),
    ...session.derivedConflictMessages,
]);

const summarizeIssueMessages = (messages: string[]) => ({
    primary: messages[0] || '',
    remainingCount: Math.max(0, messages.length - 1),
});

const getPreviewFromCommitError = (error: unknown): SchedulingPreview | null => {
    if (typeof error !== 'object' || !error || !('data' in error)) {
        return null;
    }

    const apiError = error as { data?: { data?: SchedulingPreview } };
    const previewData = apiError.data?.data;
    return previewData?.run_id ? previewData : null;
};

const getOptionKeyFromAssignment = (assignment: SchedulingAssignment) =>
    `${assignment.start_time}|${assignment.end_time}|${assignment.shift_id || ''}|${assignment.room_id}`;

const buildAssignmentFromOption = (
    session: SessionDraft,
    option: SchedulingCandidateOption,
): SchedulingAssignment => ({
    variable_id: session.variableId,
    class_id: session.classId,
    class_code: session.classCode,
    class_name: session.className,
    session_index: session.sessionIndex,
    session_total: session.sessionTotal,
    teacher_id: session.teacherId || session.baseAssignment?.teacher_id || '',
    teacher_label: session.teacherLabel || session.baseAssignment?.teacher_label || 'Chưa rõ giáo viên',
    room_id: option.room_id,
    room_name: option.room_name,
    room_capacity: option.room_capacity,
    shift_id: option.shift_id,
    shift_code: option.shift_code,
    shift_name: option.shift_name,
    shift_type: option.shift_type,
    start_time: option.start_time,
    end_time: option.end_time,
    constraint_fit: 'MANUAL_OVERRIDE',
});

const overlaps = (left: SchedulingAssignment, right: SchedulingAssignment) =>
    left.start_time < right.end_time && right.start_time < left.end_time;

const overlapsExistingLesson = (assignment: SchedulingAssignment, lesson: SchedulingExistingLesson) =>
    assignment.start_time < lesson.end_time && lesson.start_time < assignment.end_time;

const getCandidateConflictMessages = (
    session: SessionDraft,
    option: SchedulingCandidateOption,
    previewState: DerivedPreviewState,
) => {
    const assignment = buildAssignmentFromOption(session, option);
    const messages: string[] = [];

    previewState.calendarAssignments
        .filter((item) => item.variable_id !== session.variableId)
        .forEach((item) => {
            if (!overlaps(assignment, item)) {
                return;
            }

            if (assignment.class_id && assignment.class_id === item.class_id) {
                messages.push(`trùng lớp với ${item.class_code} buổi ${item.session_index}/${item.session_total}`);
            }
            if (assignment.teacher_id && assignment.teacher_id === item.teacher_id) {
                messages.push(`trùng giáo viên với ${item.class_code} buổi ${item.session_index}/${item.session_total}`);
            }
            if (assignment.room_id && assignment.room_id === item.room_id) {
                messages.push(`phòng ${option.room_name} đã bận bởi ${item.class_code}`);
            }
        });

    previewState.existingLessons.forEach((lesson) => {
        if (!overlapsExistingLesson(assignment, lesson)) {
            return;
        }

        if (assignment.class_id && assignment.class_id === lesson.class_id) {
            messages.push(`trùng lớp với buổi đã lưu ${lesson.class_code}`);
        }
        if (assignment.teacher_id && lesson.teacher_id && assignment.teacher_id === lesson.teacher_id) {
            messages.push(`trùng giáo viên với buổi đã lưu ${lesson.class_code}`);
        }
        if (assignment.room_id && lesson.room_id && assignment.room_id === lesson.room_id) {
            messages.push(`phòng ${option.room_name} đã có buổi ${lesson.class_code}`);
        }
    });

    return uniqueTags(messages);
};

const filterAvailableCandidateOptions = (
    session: SessionDraft,
    previewState: DerivedPreviewState,
) => session.candidateOptions.filter((option) => getCandidateConflictMessages(session, option, previewState).length === 0);

const scoreAssignments = (assignments: SchedulingAssignment[]) => {
    const items = [...assignments].sort((left, right) => left.start_time.localeCompare(right.start_time));
    let score = 0;
    for (let index = 0; index < items.length - 1; index += 1) {
        const current = items[index];
        const next = items[index + 1];
        if (current.teacher_id !== next.teacher_id) {
            continue;
        }

        if (current.end_time.slice(0, 10) !== next.start_time.slice(0, 10)) {
            continue;
        }

        const currentEnd = parseISO(current.end_time);
        const nextStart = parseISO(next.start_time);

        const gapMinutes = (nextStart.getTime() - currentEnd.getTime()) / (1000 * 60);
        if (gapMinutes >= 0 && gapMinutes <= 30) {
            score += 10;
        } else if (gapMinutes > 120) {
            score -= Math.floor(gapMinutes / 60) * 3;
        }
    }

    return score;
};

const formatScoreDelta = (delta: number) => (delta > 0 ? `+${delta}` : `${delta}`);

const getScoreDeltaTone = (delta: number): 'default' | 'success' | 'warning' | 'error' => {
    if (delta > 0) {
        return 'success';
    }
    if (delta < -10) {
        return 'error';
    }
    if (delta < 0) {
        return 'warning';
    }
    return 'default';
};

const replaceAssignmentForScoring = (
    session: SessionDraft,
    option: SchedulingCandidateOption,
    assignments: SchedulingAssignment[],
) => {
    const nextAssignment = buildAssignmentFromOption(session, option);
    const nextAssignments = assignments.filter((assignment) => assignment.variable_id !== session.variableId);
    nextAssignments.push(nextAssignment);
    return nextAssignments;
};

const getOptionScoreDelta = (
    session: SessionDraft,
    option: SchedulingCandidateOption,
    previewState: DerivedPreviewState,
) => scoreAssignments(replaceAssignmentForScoring(session, option, previewState.calendarAssignments)) - previewState.summary.softScore;

type ScoreContribution = {
    key: string;
    score: number;
    gapMinutes: number;
    side: 'before' | 'after';
};

type ScoreDeltaInsight = {
    delta: number;
    positives: string[];
    negatives: string[];
    neutral: string[];
};

const getNeighborContributions = (
    assignment: SchedulingAssignment,
    assignments: SchedulingAssignment[],
) => {
    if (!assignment.teacher_id) {
        return [] as ScoreContribution[];
    }

    const assignmentStartPrefix = assignment.start_time.slice(0, 10);
    const sameTeacherSameDay = assignments
        .filter((item) =>
            item.variable_id !== assignment.variable_id &&
            item.teacher_id === assignment.teacher_id &&
            item.start_time.startsWith(assignmentStartPrefix),
        )
        .sort((left, right) => left.start_time.localeCompare(right.start_time));

    const assignmentStart = parseISO(assignment.start_time);
    const previousCandidates = sameTeacherSameDay
        .filter((item) => parseISO(item.start_time) < assignmentStart);
    const previous = previousCandidates.length > 0 ? previousCandidates[previousCandidates.length - 1] : undefined;
    const next = sameTeacherSameDay.find((item) => parseISO(item.start_time) > assignmentStart);

    const contributions: ScoreContribution[] = [];
    if (previous) {
        const gapMinutes = Math.round((assignmentStart.getTime() - parseISO(previous.end_time).getTime()) / (1000 * 60));
        if (gapMinutes >= 0 && gapMinutes <= 30) {
            contributions.push({ key: `before:${previous.variable_id}`, score: 10, gapMinutes, side: 'before' });
        } else if (gapMinutes > 120) {
            contributions.push({
                key: `before:${previous.variable_id}`,
                score: -Math.floor(gapMinutes / 60) * 3,
                gapMinutes,
                side: 'before',
            });
        }
    }

    if (next) {
        const gapMinutes = Math.round((parseISO(next.start_time).getTime() - parseISO(assignment.end_time).getTime()) / (1000 * 60));
        if (gapMinutes >= 0 && gapMinutes <= 30) {
            contributions.push({ key: `after:${next.variable_id}`, score: 10, gapMinutes, side: 'after' });
        } else if (gapMinutes > 120) {
            contributions.push({
                key: `after:${next.variable_id}`,
                score: -Math.floor(gapMinutes / 60) * 3,
                gapMinutes,
                side: 'after',
            });
        }
    }

    return contributions;
};

const formatGapMinutes = (gapMinutes: number) => {
    const hours = Math.floor(gapMinutes / 60);
    const minutes = gapMinutes % 60;
    if (hours > 0 && minutes > 0) {
        return `${hours} giờ ${minutes} phút`;
    }
    if (hours > 0) {
        return `${hours} giờ`;
    }
    return `${minutes} phút`;
};

const describeContributionChange = (
    oldItem: ScoreContribution | undefined,
    newItem: ScoreContribution | undefined,
    diff: number,
) => {
    const active = newItem || oldItem;
    if (!active) {
        return null;
    }

    const sideLabel = active.side === 'before' ? 'trước buổi này' : 'sau buổi này';
    const diffLabel = diff > 0 ? `+${diff}` : `${diff}`;

    if (diff > 0) {
        if (newItem && newItem.score > 0) {
            return {
                tone: 'positive' as const,
                message: `Điểm cộng ${diffLabel}: lịch giáo viên gọn hơn ${sideLabel}, chỉ nghỉ ${formatGapMinutes(newItem.gapMinutes)}.`,
            };
        }
        if (oldItem && oldItem.score < 0) {
            return {
                tone: 'positive' as const,
                message: `Điểm cộng ${diffLabel}: giảm khoảng trống dài ${sideLabel} từ ${formatGapMinutes(oldItem.gapMinutes)}.`,
            };
        }
    }

    if (diff < 0) {
        if (newItem && newItem.score < 0) {
            return {
                tone: 'negative' as const,
                message: `Điểm trừ ${diffLabel}: tạo khoảng trống dài ${sideLabel} ${formatGapMinutes(newItem.gapMinutes)} trong lịch giáo viên.`,
            };
        }
        if (oldItem && oldItem.score > 0) {
            return {
                tone: 'negative' as const,
                message: `Điểm trừ ${diffLabel}: mất lợi thế xếp liền ca ${sideLabel}, trước đó chỉ nghỉ ${formatGapMinutes(oldItem.gapMinutes)}.`,
            };
        }
    }

    return null;
};

const getOptionScoreInsight = (
    session: SessionDraft,
    option: SchedulingCandidateOption,
    previewState: DerivedPreviewState,
): ScoreDeltaInsight => {
    const delta = getOptionScoreDelta(session, option, previewState);
    const currentAssignment = previewState.calendarAssignments.find((item) => item.variable_id === session.variableId)
        || session.baseAssignment
        || null;
    const candidateAssignment = buildAssignmentFromOption(session, option);
    const otherAssignments = previewState.calendarAssignments.filter((item) => item.variable_id !== session.variableId);

    const oldMap = new Map((currentAssignment ? getNeighborContributions(currentAssignment, otherAssignments) : []).map((item) => [item.key, item]));
    const newMap = new Map(getNeighborContributions(candidateAssignment, otherAssignments).map((item) => [item.key, item]));
    const keys = new Set([...oldMap.keys(), ...newMap.keys()]);

    const positives: string[] = [];
    const negatives: string[] = [];
    keys.forEach((key) => {
        const oldItem = oldMap.get(key);
        const newItem = newMap.get(key);
        const diff = (newItem?.score || 0) - (oldItem?.score || 0);
        if (diff === 0) {
            return;
        }
        const detail = describeContributionChange(oldItem, newItem, diff);
        if (!detail) {
            return;
        }
        if (detail.tone === 'positive') {
            positives.push(detail.message);
            return;
        }
        negatives.push(detail.message);
    });

    const neutral = [
        'Phương án này đã vượt qua kiểm tra trùng lớp, trùng giáo viên và trùng phòng trong thời điểm hiện tại.',
    ];
    if (delta === 0 && positives.length === 0 && negatives.length === 0) {
        neutral.push('Điểm mềm không đổi vì khoảng nghỉ giữa các ca của giáo viên gần như giữ nguyên.');
    }

    return {
        delta,
        positives,
        negatives,
        neutral,
    };
};

const renderScoreInsightTooltip = (insight: ScoreDeltaInsight): ReactNode => (
    <Stack spacing={0.75} sx={{ maxWidth: 320, py: 0.5 }}>
        <Typography variant="subtitle2" sx={{ fontWeight: 700 }}>
            Ảnh hưởng điểm mềm {formatScoreDelta(insight.delta)}
        </Typography>
        {insight.positives.map((item) => (
            <Typography key={`pos-${item}`} variant="caption" sx={{ color: '#bbf7d0' }}>
                • {item}
            </Typography>
        ))}
        {insight.negatives.map((item) => (
            <Typography key={`neg-${item}`} variant="caption" sx={{ color: '#fecaca' }}>
                • {item}
            </Typography>
        ))}
        {insight.neutral.map((item) => (
            <Typography key={`neu-${item}`} variant="caption" sx={{ color: 'rgba(255,255,255,0.86)' }}>
                • {item}
            </Typography>
        ))}
    </Stack>
);

const buildWeeks = (preview: SchedulingPreview) => {
    const dateFrom = preview.filters.date_from || fallbackDateFrom;
    const dateTo = preview.filters.date_to || dateFrom;
    const start = startOfWeek(parseISO(dateFrom), { weekStartsOn: 1 });
    const end = endOfWeek(parseISO(dateTo), { weekStartsOn: 1 });
    return eachWeekOfInterval({ start, end }, { weekStartsOn: 1 }).map((weekStart) =>
        Array.from({ length: 7 }, (_, index) => addDays(weekStart, index)),
    );
};

const getCalendarRowKey = (item: {
    shift_id?: string;
    start_time: string;
    end_time: string;
}) => item.shift_id || `${item.start_time.slice(11, 16)}-${item.end_time.slice(11, 16)}`;

const buildCalendarRows = (
    derivedPreview: DerivedPreviewState,
    activeOptions: SchedulingCandidateOption[],
): CalendarRow[] => {
    const rows = new Map<string, CalendarRow>();

    const registerRow = (entry: {
        shift_id?: string;
        shift_name?: string;
        shift_code?: string;
        start_time: string;
        end_time: string;
    }) => {
        const key = getCalendarRowKey(entry);
        if (rows.has(key)) {
            return;
        }

        rows.set(key, {
            key,
            shiftLabel: entry.shift_name || entry.shift_code || 'Ca học',
            timeLabel: `${entry.start_time.slice(11, 16)} - ${entry.end_time.slice(11, 16)}`,
            startTime: entry.start_time,
            endTime: entry.end_time,
        });
    };

    derivedPreview.calendarAssignments.forEach(registerRow);
    derivedPreview.existingLessons.forEach(registerRow);
    activeOptions.forEach(registerRow);

    return Array.from(rows.values()).sort((left, right) => {
        if (left.startTime === right.startTime) {
            return left.endTime.localeCompare(right.endTime);
        }
        return left.startTime.localeCompare(right.startTime);
    });
};

export const SchedulingPage = () => {
    const [preview, setPreview] = useState<SchedulingPreview | null>(null);
    const [wizardStep, setWizardStep] = useState(0);
    const [manualSelections, setManualSelections] = useState<Record<string, string>>({});
    const [activeSessionId, setActiveSessionId] = useState<string | null>(null);
    const [pendingOptionKey, setPendingOptionKey] = useState('');
    const [selectedShiftId, setSelectedShiftId] = useState<string>('');
    const [selectedRoomId, setSelectedRoomId] = useState<string>('');

    const { data: classesData, isLoading: isLoadingClasses } = useGetClassesQuery({ page: 1, limit: 200, status: 'OPEN' });
    const { data: teachersData, isLoading: isLoadingTeachers } = useGetTeachersQuery({ page: 1, limit: 200, status: 'ACTIVE' });

    const [previewScheduling, { isLoading: isPreviewing }] = usePreviewSchedulingMutation();
    const [commitSchedulingPreview, { isLoading: isCommitting }] = useCommitSchedulingPreviewMutation();
    const [loadPreview, { isFetching: isLoadingPreview }] = useLazyGetSchedulingPreviewQuery();
    const [loadLatestPreview, { isFetching: isLoadingLatest }] = useLazyGetLatestSchedulingPreviewQuery();

    const classes = classesData?.data?.classes || [];
    const teachers = teachersData?.data?.teachers || [];

    const {
        control,
        handleSubmit,
        trigger,
        formState: { errors },
    } = useForm<SchedulingFormValues>({
        resolver: zodResolver(schedulingSchema),
        defaultValues: {
            mode: 'class',
            expected_start_date: '',
            class_ids: [],
            teacher_ids: [],
        },
    });

    const selectedMode = useWatch({ control, name: 'mode' });
    const selectedClassIds = useWatch({ control, name: 'class_ids' }) || [];
    const selectedTeacherIds = useWatch({ control, name: 'teacher_ids' }) || [];
    const expectedStartDate = useWatch({ control, name: 'expected_start_date' }) || '';

    const selectedClasses = useMemo(
        () => classes.filter((item) => selectedClassIds.includes(item.id)),
        [classes, selectedClassIds],
    );

    const schedulingClasses = useMemo(() => {
        if (selectedMode === 'class') {
            return selectedClasses;
        }
        return classes.filter((item) => item.teacher_id && selectedTeacherIds.includes(item.teacher_id));
    }, [classes, selectedClasses, selectedMode, selectedTeacherIds]);

    const derivedDateFrom = expectedStartDate || deriveDefaultStartDate(schedulingClasses);
    const selectedTeacherLabels = teachers
        .filter((item) => selectedTeacherIds.includes(item.id))
        .map((item) => `${item.full_name} (${item.code})`);

    const derivedPreview = useMemo<DerivedPreviewState | null>(() => {
        if (!preview) {
            return null;
        }

        const teacherByClassId = new Map<string, { teacherId?: string; teacherLabel?: string }>();
        preview.assignments.forEach((assignment) => {
            teacherByClassId.set(assignment.class_id, {
                teacherId: assignment.teacher_id,
                teacherLabel: assignment.teacher_label,
            });
        });

        const sessionMap = new Map<string, SessionDraft>();
        preview.assignments.forEach((assignment) => {
            sessionMap.set(assignment.variable_id, {
                variableId: assignment.variable_id,
                classId: assignment.class_id,
                classCode: assignment.class_code,
                className: assignment.class_name,
                sessionIndex: assignment.session_index,
                sessionTotal: assignment.session_total,
                teacherId: assignment.teacher_id,
                teacherLabel: assignment.teacher_label,
                baseAssignment: assignment,
                resolvedAssignment: assignment,
                candidateOptions: preview.candidate_options?.[assignment.variable_id] || [],
                baseConflicts: [],
                derivedConflictMessages: [],
                conflictTags: [],
            });
        });

        preview.conflicts.forEach((conflict) => {
            if (!conflict.variable_id) {
                return;
            }

            const existing = sessionMap.get(conflict.variable_id);
            if (existing) {
                existing.baseConflicts.push(conflict);
                return;
            }

            const teacherInfo = teacherByClassId.get(conflict.class_id);
            sessionMap.set(conflict.variable_id, {
                variableId: conflict.variable_id,
                classId: conflict.class_id,
                classCode: conflict.class_code,
                className: conflict.class_name,
                sessionIndex: conflict.session_index,
                sessionTotal: conflict.session_total,
                teacherId: teacherInfo?.teacherId,
                teacherLabel: teacherInfo?.teacherLabel,
                candidateOptions: preview.candidate_options?.[conflict.variable_id] || [],
                baseConflicts: [conflict],
                derivedConflictMessages: [],
                conflictTags: buildConflictTags(conflict.type, conflict.message),
            });
        });

        Object.entries(preview.candidate_options || {}).forEach(([variableId, options]) => {
            if (sessionMap.has(variableId)) {
                return;
            }

            sessionMap.set(variableId, {
                variableId,
                classId: '',
                classCode: variableId,
                className: 'Ca cần xử lý',
                sessionIndex: 0,
                sessionTotal: 0,
                candidateOptions: options,
                baseConflicts: [],
                derivedConflictMessages: [],
                conflictTags: [],
            });
        });

        sessionMap.forEach((session) => {
            const selectedKey = manualSelections[session.variableId];
            if (!selectedKey) {
                return;
            }

            const option = session.candidateOptions.find((candidate) => candidate.key === selectedKey);
            if (!option) {
                return;
            }

            session.resolvedAssignment = buildAssignmentFromOption(session, option);
        });

        const sessions = Array.from(sessionMap.values()).sort((left, right) => {
            const leftTime = left.resolvedAssignment?.start_time || left.baseAssignment?.start_time || left.candidateOptions[0]?.start_time || '';
            const rightTime = right.resolvedAssignment?.start_time || right.baseAssignment?.start_time || right.candidateOptions[0]?.start_time || '';
            return leftTime.localeCompare(rightTime);
        });

        const resolvedAssignments = sessions
            .map((session) => session.resolvedAssignment)
            .filter((assignment): assignment is SchedulingAssignment => !!assignment);

        const overlapMessages = new Map<string, Set<string>>();
        const overlapTags = new Map<string, Set<string>>();
        const existingLessonConflictTags = new Map<string, Set<string>>();
        for (let i = 0; i < resolvedAssignments.length; i += 1) {
            for (let j = i + 1; j < resolvedAssignments.length; j += 1) {
                const left = resolvedAssignments[i];
                const right = resolvedAssignments[j];
                if (!overlaps(left, right)) {
                    continue;
                }

                if (!overlapTags.has(left.variable_id)) overlapTags.set(left.variable_id, new Set(['Trùng ngày/ca']));
                if (!overlapTags.has(right.variable_id)) overlapTags.set(right.variable_id, new Set(['Trùng ngày/ca']));
                overlapTags.get(left.variable_id)?.add('Trùng ngày/ca');
                overlapTags.get(right.variable_id)?.add('Trùng ngày/ca');

                if (left.class_id === right.class_id) {
                    if (!overlapMessages.has(left.variable_id)) overlapMessages.set(left.variable_id, new Set());
                    if (!overlapMessages.has(right.variable_id)) overlapMessages.set(right.variable_id, new Set());
                    overlapTags.get(left.variable_id)?.add('Trùng lớp');
                    overlapTags.get(right.variable_id)?.add('Trùng lớp');
                    overlapMessages.get(left.variable_id)?.add(`trùng lớp với ${right.class_code} buổi ${right.session_index}/${right.session_total} lúc ${formatCompactDateTime(right.start_time)}`);
                    overlapMessages.get(right.variable_id)?.add(`trùng lớp với ${left.class_code} buổi ${left.session_index}/${left.session_total} lúc ${formatCompactDateTime(left.start_time)}`);
                }
                if (left.teacher_id && left.teacher_id === right.teacher_id) {
                    if (!overlapMessages.has(left.variable_id)) overlapMessages.set(left.variable_id, new Set());
                    if (!overlapMessages.has(right.variable_id)) overlapMessages.set(right.variable_id, new Set());
                    overlapTags.get(left.variable_id)?.add('Trùng giáo viên');
                    overlapTags.get(right.variable_id)?.add('Trùng giáo viên');
                    overlapMessages.get(left.variable_id)?.add(`trùng giáo viên với ${right.class_code} buổi ${right.session_index}/${right.session_total}`);
                    overlapMessages.get(right.variable_id)?.add(`trùng giáo viên với ${left.class_code} buổi ${left.session_index}/${left.session_total}`);
                }
                if (left.room_id && left.room_id === right.room_id) {
                    if (!overlapMessages.has(left.variable_id)) overlapMessages.set(left.variable_id, new Set());
                    if (!overlapMessages.has(right.variable_id)) overlapMessages.set(right.variable_id, new Set());
                    overlapTags.get(left.variable_id)?.add('Trùng phòng');
                    overlapTags.get(right.variable_id)?.add('Trùng phòng');
                    overlapMessages.get(left.variable_id)?.add(`trùng phòng với ${right.class_code} buổi ${right.session_index}/${right.session_total}`);
                    overlapMessages.get(right.variable_id)?.add(`trùng phòng với ${left.class_code} buổi ${left.session_index}/${left.session_total}`);
                }
            }
        }

        for (const assignment of resolvedAssignments) {
            for (const lesson of preview.existing_lessons || []) {
                if (!overlapsExistingLesson(assignment, lesson)) {
                    continue;
                }

                const reasons: string[] = [];
                if (assignment.class_id === lesson.class_id) {
                    reasons.push('trùng lớp');
                }
                if (assignment.teacher_id && lesson.teacher_id && assignment.teacher_id === lesson.teacher_id) {
                    reasons.push('trùng giáo viên');
                }
                if (assignment.room_id && lesson.room_id && assignment.room_id === lesson.room_id) {
                    reasons.push('trùng phòng');
                }
                if (reasons.length === 0) {
                    continue;
                }

                if (!overlapMessages.has(assignment.variable_id)) {
                    overlapMessages.set(assignment.variable_id, new Set());
                }
                if (!overlapTags.has(assignment.variable_id)) {
                    overlapTags.set(assignment.variable_id, new Set());
                }
                if (!existingLessonConflictTags.has(lesson.lesson_id)) {
                    existingLessonConflictTags.set(lesson.lesson_id, new Set());
                }
                overlapTags.get(assignment.variable_id)?.add('Trùng ngày/ca');
                existingLessonConflictTags.get(lesson.lesson_id)?.add('Trùng ngày/ca');
                overlapMessages.get(assignment.variable_id)?.add(
                    `${reasons.join(', ')} với buổi đã lưu ${lesson.class_code} lúc ${formatCompactDateTime(lesson.start_time)}`,
                );
                reasons.forEach((reason) => {
                    const label = reason.replace(/^trùng /, 'Trùng ');
                    overlapTags.get(assignment.variable_id)?.add(label);
                    existingLessonConflictTags.get(lesson.lesson_id)?.add(label);
                });
            }
        }

        sessions.forEach((session) => {
            session.derivedConflictMessages = Array.from(overlapMessages.get(session.variableId) || []);
            session.conflictTags = uniqueTags([
                ...session.baseConflicts.flatMap((conflict) => buildConflictTags(conflict.type, conflict.message)),
                ...Array.from(overlapTags.get(session.variableId) || []),
            ]);
        });

        const unresolvedSessions = sessions.filter((session) => !session.resolvedAssignment);
        const unresolvedConflictEntries = sessions
            .filter((session) => !session.resolvedAssignment)
            .flatMap((session) => {
                if (session.baseConflicts.length > 0) {
                    return session.baseConflicts;
                }
                return [{
                    variable_id: session.variableId,
                    class_id: session.classId,
                    class_code: session.classCode,
                    class_name: session.className,
                    session_index: session.sessionIndex,
                    session_total: session.sessionTotal,
                    type: 'NO_DOMAIN',
                    message: 'Ca học này chưa có phương án xếp lịch hợp lệ.',
                }];
            });

        const overlapConflictEntries = sessions
            .filter((session) => session.derivedConflictMessages.length > 0)
            .map((session) => ({
                variable_id: session.variableId,
                class_id: session.classId,
                class_code: session.classCode,
                class_name: session.className,
                session_index: session.sessionIndex,
                session_total: session.sessionTotal,
                type: 'ASSIGNMENT_CONFLICT',
                message: `Buổi học đang ${session.derivedConflictMessages.join('; ')}`,
            }));

        const fixedVariableIds = new Set(
            sessions.filter((session) => session.resolvedAssignment).map((session) => session.variableId),
        );
        const immutableConflicts = preview.conflicts.filter((conflict) => !conflict.variable_id || !fixedVariableIds.has(conflict.variable_id));

        const visibleConflicts = [
            ...immutableConflicts.filter((conflict) => !conflict.variable_id),
            ...unresolvedConflictEntries,
            ...overlapConflictEntries,
        ];

        const manualAssignmentPayload = Object.entries(manualSelections)
            .map(([variable_id, option_key]) => {
                const session = sessionMap.get(variable_id);
                if (!session) {
                    return null;
                }
                if (session.baseAssignment && getOptionKeyFromAssignment(session.baseAssignment) === option_key) {
                    return null;
                }
                return { variable_id, option_key };
            })
            .filter((item): item is { variable_id: string; option_key: string } => !!item);

        const requestedClasses = new Set(
            sessions.map((session) => session.classId).filter((value) => value),
        ).size || preview.summary.requested_classes;

        const requestedSessions = Math.max(sessions.length, preview.summary.requested_sessions);
        const manualAdjustmentCount = manualAssignmentPayload.length;
        const manualAdjustmentLimit = getAllowedManualAdjustmentLimit(requestedSessions);
        const baseSoftScore = preview.summary.soft_score;
        const softScore = scoreAssignments(resolvedAssignments);

        const summary = {
            requestedClasses,
            requestedSessions,
            scheduledLessons: resolvedAssignments.length,
            unscheduledLessons: unresolvedSessions.length,
            conflictCount: visibleConflicts.length,
            baseSoftScore,
            softScore,
            scoreDelta: softScore - baseSoftScore,
            manualAdjustmentCount,
            manualAdjustmentLimit,
            excessiveManualAdjustment: manualAdjustmentCount > manualAdjustmentLimit,
        };

        const status =
            summary.scheduledLessons === 0
                ? 'FAILED'
                : summary.unscheduledLessons === 0 && summary.conflictCount === 0 && !summary.excessiveManualAdjustment
                    ? 'COMPLETED'
                    : 'PARTIAL';

        return {
            sessions,
            calendarAssignments: [...resolvedAssignments].sort((left, right) => left.start_time.localeCompare(right.start_time)),
            existingLessons: preview.existing_lessons || [],
            unresolvedSessions,
            remainingConflicts: visibleConflicts,
            summary,
            status,
            manualAssignmentPayload,
            existingLessonConflictTags: Object.fromEntries(
                Array.from(existingLessonConflictTags.entries()).map(([lessonId, tags]) => [lessonId, Array.from(tags)]),
            ),
        };
    }, [manualSelections, preview]);

    const weeks = useMemo(() => (preview ? buildWeeks(preview) : []), [preview]);

    const filteredSessions = useMemo(() => derivedPreview?.sessions || [], [derivedPreview]);

    const classPreviewSummaries = useMemo(() => {
        if (!derivedPreview) {
            return [];
        }

        const byClass = new Map<string, {
            classCode: string;
            className: string;
            teacherLabel: string;
            totalSessions: number;
            scheduledSessions: number;
            conflictTags: string[];
            firstTime?: string;
            lastTime?: string;
        }>();

        derivedPreview.sessions.forEach((session) => {
            const key = session.classId || session.classCode;
            const current = byClass.get(key) || {
                classCode: session.classCode,
                className: session.className,
                teacherLabel: session.teacherLabel || 'Chưa rõ giáo viên',
                totalSessions: session.sessionTotal || 0,
                scheduledSessions: 0,
                conflictTags: [],
            };
            current.totalSessions = Math.max(current.totalSessions, session.sessionTotal || 0);
            if (session.resolvedAssignment) {
                current.scheduledSessions += 1;
                const start = session.resolvedAssignment.start_time;
                const end = session.resolvedAssignment.end_time;
                if (!current.firstTime || start < current.firstTime) current.firstTime = start;
                if (!current.lastTime || end > current.lastTime) current.lastTime = end;
            }
            current.conflictTags = uniqueTags([...current.conflictTags, ...session.conflictTags]);
            byClass.set(key, current);
        });

        return Array.from(byClass.values()).sort((left, right) => left.classCode.localeCompare(right.classCode));
    }, [derivedPreview]);

    const activeSession = useMemo(
        () => derivedPreview?.sessions.find((session) => session.variableId === activeSessionId) || null,
        [activeSessionId, derivedPreview],
    );
    const actionableSessions = useMemo(
        () => filteredSessions.filter((session) => !session.resolvedAssignment || session.derivedConflictMessages.length > 0),
        [filteredSessions],
    );
    const activeSessionIssueSummary = useMemo(() => {
        if (!activeSession) {
            return null;
        }
        return summarizeIssueMessages(collectSessionIssueMessages(activeSession));
    }, [activeSession]);
    const activeSessionAvailableOptions = useMemo(() => {
        if (!activeSession || !derivedPreview) {
            return [];
        }
        return filterAvailableCandidateOptions(activeSession, derivedPreview);
    }, [activeSession, derivedPreview]);
    const activeShiftOptions = useMemo(() => {
        const seen = new Map<string, string>();
        activeSessionAvailableOptions.forEach((option) => {
            const key = option.shift_id || '__no_shift__';
            const label = option.shift_name || option.shift_code || 'Ca chưa đặt tên';
            if (!seen.has(key)) {
                seen.set(key, label);
            }
        });
        return Array.from(seen.entries()).map(([value, label]) => ({ value, label }));
    }, [activeSessionAvailableOptions]);
    const activeRoomOptions = useMemo(() => {
        const seen = new Map<string, string>();
        activeSessionAvailableOptions
            .filter((option) => !selectedShiftId || (option.shift_id || '__no_shift__') === selectedShiftId)
            .forEach((option) => {
                if (!seen.has(option.room_id)) {
                    seen.set(option.room_id, option.room_name);
                }
            });
        return Array.from(seen.entries()).map(([value, label]) => ({ value, label }));
    }, [activeSessionAvailableOptions, selectedShiftId]);
    const filteredActiveSessionOptions = useMemo(() => (
        activeSessionAvailableOptions.filter((option) => {
            const optionShiftId = option.shift_id || '__no_shift__';
            if (selectedShiftId && optionShiftId !== selectedShiftId) {
                return false;
            }
            if (selectedRoomId && option.room_id !== selectedRoomId) {
                return false;
            }
            return true;
        })
    ), [activeSessionAvailableOptions, selectedRoomId, selectedShiftId]);
    useEffect(() => {
        if (!activeSession) {
            return;
        }
        if (filteredActiveSessionOptions.some((option) => option.key === pendingOptionKey)) {
            return;
        }
        setPendingOptionKey(filteredActiveSessionOptions[0]?.key || '');
    }, [activeSession, filteredActiveSessionOptions, pendingOptionKey]);

    const activeStep = wizardStep;
    const canCommitPreview = !!derivedPreview && derivedPreview.status === 'COMPLETED' && !derivedPreview.summary.excessiveManualAdjustment;
    const calendarRows = useMemo(
        () => (derivedPreview ? buildCalendarRows(derivedPreview, filteredActiveSessionOptions) : []),
        [derivedPreview, filteredActiveSessionOptions],
    );

    const getAvailableOptionsForSession = (session: SessionDraft) =>
        derivedPreview ? filterAvailableCandidateOptions(session, derivedPreview) : [];

    const resetManualSelections = () => {
        setManualSelections({});
        setActiveSessionId(null);
        setPendingOptionKey('');
        setSelectedShiftId('');
        setSelectedRoomId('');
    };

    const syncPreviewState = (nextPreview: SchedulingPreview) => {
        setPreview(nextPreview);
        setWizardStep(2);
        resetManualSelections();
    };

    const handleContinueConfig = async () => {
        const isValid = await trigger();
        if (isValid) {
            setWizardStep(1);
        }
    };

    const handleRunPreview = handleSubmit(async (values) => {
        try {
            const dateFrom = values.expected_start_date || deriveDefaultStartDate(schedulingClasses);
            const response = await previewScheduling({
                date_from: dateFrom,
                mode: 'replan_with_published_lock',
                class_ids: values.mode === 'class' ? values.class_ids : [],
                teacher_ids: values.mode === 'teacher' ? values.teacher_ids : [],
            }).unwrap();
            syncPreviewState(response.data);
            toast.success('Đã chạy xem trước xếp lịch');
        } catch (error) {
            toast.error(getApiErrorMessage(error, 'Không thể chạy xem trước xếp lịch'));
        }
    });

    const handleLoadLatest = async () => {
        try {
            const response = await loadLatestPreview().unwrap();
            syncPreviewState(response.data);
            toast.success('Đã tải lịch xem trước mới nhất');
        } catch (error) {
            toast.error(getApiErrorMessage(error, 'Chưa có lịch xem trước nào để tải'));
        }
    };

    const handleRefreshCurrent = async () => {
        if (!preview?.run_id) {
            await handleLoadLatest();
            return;
        }

        try {
            const response = await loadPreview(preview.run_id).unwrap();
            syncPreviewState(response.data);
            toast.success('Đã làm mới kết quả lịch xem trước');
        } catch (error) {
            toast.error(getApiErrorMessage(error, 'Không thể tải lại lịch xem trước hiện tại'));
        }
    };

    const toggleEditMode = (variableId: string) => {
        if (activeSessionId === variableId) {
            setActiveSessionId(null);
            setPendingOptionKey('');
            return;
        }

        const session = derivedPreview?.sessions.find((item) => item.variableId === variableId);
        if (!session || !derivedPreview) {
            return;
        }

        const availableOptions = filterAvailableCandidateOptions(session, derivedPreview);
        const selectedKey = manualSelections[variableId];
        const currentKey = session.baseAssignment ? getOptionKeyFromAssignment(session.baseAssignment) : '';
        const initialKey =
            (selectedKey && availableOptions.some((option) => option.key === selectedKey) ? selectedKey : '') ||
            (currentKey && availableOptions.some((option) => option.key === currentKey) ? currentKey : '') ||
            availableOptions[0]?.key ||
            '';
        const initialOption = availableOptions.find((option) => option.key === initialKey);
        setActiveSessionId(variableId);
        setSelectedShiftId(initialOption?.shift_id || '__no_shift__');
        setSelectedRoomId(initialOption?.room_id || '');
        setPendingOptionKey(initialKey);
    };

    const handleApplyManualOption = (variableId: string, optionKey: string) => {
        setManualSelections((current) => ({
            ...current,
            [variableId]: optionKey,
        }));
        setActiveSessionId(null);
        setPendingOptionKey('');
        setSelectedShiftId('');
        setSelectedRoomId('');
    };

    const handleCommitPreview = async () => {
        if (!preview?.run_id || !derivedPreview) {
            toast.error('Chưa có bản xem trước nào để lưu');
            return;
        }

        try {
            const response = await commitSchedulingPreview({
                run_id: preview.run_id,
                manual_assignments: derivedPreview.manualAssignmentPayload,
            }).unwrap();
            toast.success(response.data.message || `Đã xác nhận lưu ${response.data.scheduled_lessons} buổi học`);
        } catch (error) {
            const previewFromError = getPreviewFromCommitError(error);
            if (previewFromError) {
                syncPreviewState(previewFromError);
            } else {
                try {
                    const refreshed = await loadPreview(preview.run_id).unwrap();
                    syncPreviewState(refreshed.data);
                } catch {
                    // Keep the current client-side state if backend preview cannot be reloaded.
                }
            }
            toast.error(getApiErrorMessage(error, 'Không thể lưu bản xem trước xếp lịch'));
        }
    };

    const renderPreviewSkeleton = () => (
        <Stack spacing={1.5}>
            <Skeleton variant="rounded" height={80} />
            <Skeleton variant="rounded" height={520} />
        </Stack>
    );

    return (
        <Box>
            <PageHeader
                title="Xếp lịch"
                breadcrumbs={[
                    { label: 'Tổng quan', path: '/app/admin/overview' },
                    { label: 'Xếp lịch' },
                ]}
                actions={
                    <Stack direction="row" spacing={1}>
                        <Button
                            variant="outlined"
                            startIcon={<RefreshRounded />}
                            onClick={() => void handleLoadLatest()}
                            disabled={isLoadingLatest}
                        >
                            Tải lịch học xem trước mới nhất
                        </Button>
                    </Stack>
                }
            />

            <Stack spacing={3}>
                <Paper variant="outlined" sx={{ p: 3, borderRadius: 4 }}>
                    <Stepper activeStep={activeStep} alternativeLabel>
                        {previewSteps.map((label) => (
                            <Step key={label}>
                                <StepLabel>{label}</StepLabel>
                            </Step>
                        ))}
                    </Stepper>
                </Paper>

                {wizardStep === 0 ? (
                    <Paper variant="outlined" sx={{ p: 3, borderRadius: 4 }}>
                        <Stack spacing={2.5}>
                            <Box>
                                <Typography variant="h6" sx={{ fontWeight: 700 }}>
                                    Chọn đối tượng cần xếp lịch
                                </Typography>
                            </Box>

                            <Controller
                                control={control}
                                name="mode"
                                render={({ field }) => (
                                    <TextField {...field} select label="Xếp lịch theo" fullWidth>
                                        <MenuItem value="class">Theo lớp học</MenuItem>
                                        <MenuItem value="teacher">Theo giáo viên phụ trách</MenuItem>
                                    </TextField>
                                )}
                            />

                            {selectedMode === 'class' ? (
                                <Controller
                                    control={control}
                                    name="class_ids"
                                    render={({ field }) => (
                                        <Autocomplete
                                            multiple
                                            options={classes}
                                            loading={isLoadingClasses}
                                            value={classes.filter((item) => field.value?.includes(item.id))}
                                            onChange={(_, values) => field.onChange(values.map((item) => item.id))}
                                            getOptionLabel={(option) => `${option.name} (${option.code})`}
                                            renderInput={(params) => (
                                                <TextField
                                                    {...params}
                                                    label="Chọn lớp"
                                                    error={!!errors.class_ids}
                                                    helperText={errors.class_ids?.message}
                                                />
                                            )}
                                        />
                                    )}
                                />
                            ) : (
                                <Controller
                                    control={control}
                                    name="teacher_ids"
                                    render={({ field }) => (
                                        <Autocomplete
                                            multiple
                                            options={teachers}
                                            loading={isLoadingTeachers}
                                            value={teachers.filter((item) => field.value?.includes(item.id))}
                                            onChange={(_, values) => field.onChange(values.map((item) => item.id))}
                                            getOptionLabel={(option) => `${option.full_name} (${option.code})`}
                                            renderInput={(params) => (
                                                <TextField
                                                    {...params}
                                                    label="Chọn giáo viên"
                                                    error={!!errors.teacher_ids}
                                                    helperText={errors.teacher_ids?.message}
                                                />
                                            )}
                                        />
                                    )}
                                />
                            )}

                            <Controller
                                control={control}
                                name="expected_start_date"
                                render={({ field }) => (
                                    <TextField
                                        {...field}
                                        type="date"
                                        label="Ngày bắt đầu dự kiến"
                                        fullWidth
                                        InputLabelProps={{ shrink: true }}
                                        helperText={`Để trống sẽ dùng từ ${format(parseISO(deriveDefaultStartDate(schedulingClasses)), 'dd/MM/yyyy', { locale: vi })}.`}
                                    />
                                )}
                            />

                            <Stack direction="row" justifyContent="flex-end">
                                <Button variant="contained" onClick={() => void handleContinueConfig()}>
                                    Tiếp tục
                                </Button>
                            </Stack>
                        </Stack>
                    </Paper>
                ) : null}

                {wizardStep === 1 ? (
                    <Paper variant="outlined" sx={{ p: 3, borderRadius: 4 }}>
                        <Stack spacing={2.5}>
                            <Box>
                                <Typography variant="h6" sx={{ fontWeight: 700 }}>
                                    Xác nhận cấu hình trước khi chạy
                                </Typography>
                            </Box>

                            <Stack direction={{ xs: 'column', md: 'row' }} spacing={2}>
                                <Paper variant="outlined" sx={{ p: 2, borderRadius: 3, flex: 1 }}>
                                    <Typography variant="subtitle2" sx={{ fontWeight: 700 }}>
                                        Đối tượng
                                    </Typography>
                                    <Typography variant="body2" color="text.secondary">
                                        {selectedMode === 'class'
                                            ? `${selectedClasses.length} lớp được chọn`
                                            : `${selectedTeacherLabels.length} giáo viên, ${schedulingClasses.length} lớp phụ trách`}
                                    </Typography>
                                    <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap sx={{ mt: 1 }}>
                                        {(selectedMode === 'class'
                                            ? selectedClasses.map((item) => `${item.name} (${item.code})`)
                                            : selectedTeacherLabels
                                        ).map((label) => (
                                            <Chip key={label} label={label} size="small" variant="outlined" />
                                        ))}
                                    </Stack>
                                </Paper>

                                <Paper variant="outlined" sx={{ p: 2, borderRadius: 3, flex: 1 }}>
                                    <Typography variant="subtitle2" sx={{ fontWeight: 700 }}>
                                        Ngày bắt đầu dự kiến
                                    </Typography>
                                    <Typography variant="body2" color="text.secondary">
                                        {format(parseISO(derivedDateFrom), 'dd/MM/yyyy', { locale: vi })}
                                    </Typography>
                                </Paper>
                            </Stack>

                            <Stack direction="row" justifyContent="space-between">
                                <Button variant="outlined" onClick={() => setWizardStep(0)}>
                                    Quay lại
                                </Button>
                                <Button
                                    variant="contained"
                                    startIcon={<PlayArrowRounded />}
                                    onClick={() => void handleRunPreview()}
                                    disabled={isPreviewing || schedulingClasses.length === 0}
                                >
                                    {isPreviewing ? 'Đang chạy...' : 'Chạy xếp lịch'}
                                </Button>
                            </Stack>
                        </Stack>
                    </Paper>
                ) : null}

                {isPreviewing || isLoadingPreview || isLoadingLatest ? renderPreviewSkeleton() : null}

                {!isPreviewing && !isLoadingPreview && !isLoadingLatest && !preview && wizardStep === 2 ? (
                    <Paper
                        variant="outlined"
                        sx={{ p: 6, borderRadius: 4, borderStyle: 'dashed', textAlign: 'center' }}
                    >
                        <RuleRounded sx={{ fontSize: 48, color: 'text.secondary', mb: 2 }} />
                        <Typography variant="h6" sx={{ fontWeight: 700 }}>
                            Chưa có bản xem trước nào
                        </Typography>
                        <Typography variant="body2" color="text.secondary" sx={{ mt: 1, mb: 3 }}>
                            Bấm `Chạy xếp lịch` để tạo bản xem trước đầu tiên hoặc tải lần chạy gần nhất từ hệ thống.
                        </Typography>
                        <Stack direction="row" spacing={1} justifyContent="center">
                            <Button variant="contained" startIcon={<PlayArrowRounded />} onClick={() => void handleRunPreview()}>
                                Chạy xếp lịch
                            </Button>
                            <Button variant="outlined" startIcon={<RefreshRounded />} onClick={() => void handleLoadLatest()}>
                                Tải lịch xem trước mới nhất
                            </Button>
                        </Stack>
                    </Paper>
                ) : null}

                {preview && derivedPreview && wizardStep === 2 ? (
                    <Stack spacing={2.5}>
                        <Paper
                            variant="outlined"
                            sx={{
                                p: 3,
                                borderRadius: 4,
                                background: 'linear-gradient(135deg, rgba(15,118,110,0.08), rgba(13,148,136,0.02))',
                            }}
                        >
                            <Stack direction={{ xs: 'column', md: 'row' }} spacing={2} justifyContent="space-between">
                                <Box>
                                    <Typography variant="h6" sx={{ fontWeight: 800 }}>
                                        Trạng thái xem trước: {previewStatusLabel(derivedPreview.status)}
                                    </Typography>
                                    <Typography variant="body2" color="text.secondary">
                                        Cập nhật lúc {formatDateTime(preview.generated_at)}
                                    </Typography>
                                </Box>
                                <Stack direction="row" spacing={1} flexWrap="wrap">
                                    <Chip label={`${derivedPreview.summary.requestedClasses} lớp`} color="default" variant="outlined" />
                                    <Chip label={`${derivedPreview.summary.requestedSessions} buổi cần xếp`} color="default" variant="outlined" />
                                    <Chip label={`${derivedPreview.summary.scheduledLessons} buổi đã xếp`} color="success" variant="outlined" />
                                    <Chip label={`${derivedPreview.summary.unscheduledLessons} buổi chưa xếp`} color="warning" variant="outlined" />
                                    <Chip label={`${derivedPreview.summary.conflictCount} trùng`} color="error" variant="outlined" />
                                    <Chip label={`Điểm mềm ${derivedPreview.summary.softScore}`} color="primary" variant="outlined" />
                                    <Chip
                                        label={`Chỉnh tay ${derivedPreview.summary.manualAdjustmentCount}/${derivedPreview.summary.manualAdjustmentLimit}`}
                                        color={derivedPreview.summary.excessiveManualAdjustment ? 'error' : 'secondary'}
                                        variant="outlined"
                                    />
                                </Stack>
                            </Stack>
                        </Paper>

                        <Stack direction={{ xs: 'column', lg: 'row' }} spacing={2}>
                            <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 4, flex: 1 }}>
                                <Typography variant="subtitle2" sx={{ fontWeight: 800, mb: 0.75 }}>
                                    Tác động khi xếp lại
                                </Typography>
                                <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                                    <Chip
                                        label={`${preview.summary.schedule_change_count} ca đổi lịch`}
                                        color={preview.summary.schedule_change_count > 0 ? 'warning' : 'success'}
                                        variant="outlined"
                                    />
                                    <Chip
                                        label={`${preview.summary.teacher_change_count} ca đổi giáo viên`}
                                        color={preview.summary.teacher_change_count > 0 ? 'warning' : 'success'}
                                        variant="outlined"
                                    />
                                    <Chip
                                        label={`${preview.summary.room_change_count} ca đổi phòng`}
                                        color={preview.summary.room_change_count > 0 ? 'warning' : 'success'}
                                        variant="outlined"
                                    />
                                    <Chip
                                        label={`Lấp đầy ${(preview.summary.average_capacity_utilization * 100).toFixed(0)}%`}
                                        color="default"
                                        variant="outlined"
                                    />
                                </Stack>
                            </Paper>
                        </Stack>

                        <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 4 }}>
                            <Typography variant="h6" sx={{ fontWeight: 700, mb: 1.5 }}>
                                Thông tin lớp trong bản xem trước
                            </Typography>
                            <Stack direction={{ xs: 'column', lg: 'row' }} spacing={1.5} flexWrap="wrap" useFlexGap>
                                {classPreviewSummaries.map((item) => (
                                    <Paper
                                        key={item.classCode}
                                        variant="outlined"
                                        sx={{ p: 1.75, borderRadius: 3, flex: '1 1 280px', borderColor: item.conflictTags.length ? 'warning.light' : 'success.light' }}
                                    >
                                        <Stack spacing={1}>
                                            <Box>
                                                <Typography variant="subtitle2" sx={{ fontWeight: 800 }}>
                                                    {item.classCode} - {item.className}
                                                </Typography>
                                                <Typography variant="caption" color="text.secondary">
                                                    {item.teacherLabel}
                                                </Typography>
                                            </Box>
                                            <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                                                <Chip size="small" label={`${item.scheduledSessions}/${item.totalSessions || item.scheduledSessions} buổi đã xếp`} color="primary" variant="outlined" />
                                                {item.firstTime ? (
                                                    <Chip
                                                        size="small"
                                                        label={`${formatCompactDateTime(item.firstTime)} - ${item.lastTime ? formatCompactDateTime(item.lastTime) : ''}`}
                                                        variant="outlined"
                                                    />
                                                ) : null}
                                            </Stack>
                                            {renderConflictChips(item.conflictTags)}
                                        </Stack>
                                    </Paper>
                                ))}
                            </Stack>
                        </Paper>

                        <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 4 }}>
                            <Stack
                                direction={{ xs: 'column', lg: 'row' }}
                                spacing={1.5}
                                justifyContent="space-between"
                                alignItems={{ xs: 'stretch', lg: 'center' }}
                                sx={{ mb: 2 }}
                            >
                                <Stack direction="row" spacing={1}>
                                    <Button
                                        variant="outlined"
                                        onClick={() => setWizardStep(1)}
                                    >
                                        Quay lại cấu hình
                                    </Button>
                                    <Button
                                        variant="outlined"
                                        startIcon={<RefreshRounded />}
                                        onClick={() => void handleRefreshCurrent()}
                                        disabled={isLoadingPreview || isLoadingLatest}
                                    >
                                        Làm mới
                                    </Button>
                                    <Button
                                        variant="outlined"
                                        startIcon={<TuneRounded />}
                                        onClick={resetManualSelections}
                                        disabled={derivedPreview.manualAssignmentPayload.length === 0}
                                    >
                                        Bỏ chỉnh tay
                                    </Button>
                                    <Button
                                        variant="contained"
                                        color={canCommitPreview ? 'primary' : 'inherit'}
                                        onClick={() => void handleCommitPreview()}
                                        disabled={!canCommitPreview || isCommitting}
                                    >
                                        {isCommitting ? 'Đang xác nhận...' : 'Xác nhận lưu'}
                                    </Button>
                                </Stack>
                            </Stack>

                            <Stack spacing={2.5}>
                                <Stack direction="row" spacing={1} alignItems="center">
                                    <CalendarMonthRounded color="primary" />
                                    <Box>
                                        <Typography variant="h6" sx={{ fontWeight: 700 }}>
                                            Xem trước lịch
                                        </Typography>
                                    </Box>
                                </Stack>

                                <Stack spacing={2}>
                                    {weeks.map((week, weekIndex) => (
                                        <Paper key={`${week[0].toISOString()}-${weekIndex}`} variant="outlined" sx={{ borderRadius: 3, overflow: 'hidden' }}>
                                            <Box sx={{ px: 2, py: 1.5, backgroundColor: '#f8fafc' }}>
                                                <Typography variant="subtitle2" sx={{ fontWeight: 700 }}>
                                                    Tuần {weekIndex + 1}: {format(week[0], 'dd/MM', { locale: vi })} - {format(week[6], 'dd/MM/yyyy', { locale: vi })}
                                                </Typography>
                                            </Box>
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

                                                    {calendarRows.map((row) => (
                                                        <Box
                                                            key={`calendar-row-${row.key}`}
                                                            sx={{ display: 'contents' }}
                                                        >
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
                                                                const dayString = format(day, 'yyyy-MM-dd');
                                                                const rowAssignments = derivedPreview.calendarAssignments.filter((assignment) =>
                                                                    assignment.start_time.startsWith(dayString) && getCalendarRowKey(assignment) === row.key,
                                                                );
                                                                const rowExistingLessons = derivedPreview.existingLessons.filter((lesson) =>
                                                                    lesson.start_time.startsWith(dayString) && getCalendarRowKey(lesson) === row.key,
                                                                );
                                                                const rowActiveOptions = activeSession
                                                                    ? filteredActiveSessionOptions.filter((option) =>
                                                                        option.start_time.startsWith(dayString) && getCalendarRowKey(option) === row.key,
                                                                    )
                                                                    : [];

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
                                                                            {rowActiveOptions.map((option) => {
                                                                                const insight = getOptionScoreInsight(activeSession!, option, derivedPreview);
                                                                                const isSelected = pendingOptionKey === option.key;
                                                                                return (
                                                                                    <Tooltip
                                                                                        key={`${activeSession!.variableId}-${option.key}`}
                                                                                        title={renderScoreInsightTooltip(insight)}
                                                                                        placement="top"
                                                                                        arrow
                                                                                    >
                                                                                        <Paper
                                                                                            variant="outlined"
                                                                                            onClick={() => handleApplyManualOption(activeSession!.variableId, option.key)}
                                                                                            sx={{
                                                                                                p: 1,
                                                                                                borderRadius: 2,
                                                                                                cursor: 'pointer',
                                                                                                borderStyle: 'dashed',
                                                                                                borderWidth: isSelected ? 2 : 1,
                                                                                                borderColor: isSelected ? 'primary.main' : 'secondary.main',
                                                                                                backgroundColor: isSelected ? 'rgba(14,165,233,0.08)' : 'rgba(168,85,247,0.04)',
                                                                                            }}
                                                                                        >
                                                                                            <Stack spacing={0.5}>
                                                                                                <Typography variant="body2" sx={{ fontWeight: 700, color: isSelected ? 'primary.dark' : 'secondary.dark' }}>
                                                                                                    Buổi {activeSession!.sessionIndex}/{activeSession!.sessionTotal}
                                                                                                </Typography>
                                                                                                <Typography variant="caption" color="text.secondary">
                                                                                                    {option.room_name}
                                                                                                </Typography>
                                                                                                <Chip
                                                                                                    size="small"
                                                                                                    color={getScoreDeltaTone(insight.delta)}
                                                                                                    variant={isSelected ? 'filled' : 'outlined'}
                                                                                                    label={`Điểm mềm ${formatScoreDelta(insight.delta)}`}
                                                                                                />
                                                                                            </Stack>
                                                                                        </Paper>
                                                                                    </Tooltip>
                                                                                );
                                                                            })}

                                                                            {rowExistingLessons.map((lesson) => (
                                                                                <Tooltip
                                                                                    key={lesson.lesson_id}
                                                                                    title={`${formatDateTime(lesson.start_time)} - ${format(parseISO(lesson.end_time), 'HH:mm', { locale: vi })} • ${lesson.room_name || 'Chưa có phòng'} • ${lesson.teacher_label || 'Chưa có giáo viên'}`}
                                                                                    arrow
                                                                                >
                                                                                    <Paper
                                                                                        variant="outlined"
                                                                                        sx={{
                                                                                            p: 1,
                                                                                            borderRadius: 2,
                                                                                            backgroundColor: 'rgba(15,23,42,0.04)',
                                                                                            borderColor: derivedPreview.existingLessonConflictTags[lesson.lesson_id]?.length
                                                                                                ? 'error.light'
                                                                                                : 'rgba(15,23,42,0.12)',
                                                                                        }}
                                                                                    >
                                                                                        <Stack spacing={0.5}>
                                                                                            <Typography variant="body2" sx={{ fontWeight: 700 }}>
                                                                                                {lesson.class_code} • Đã lưu
                                                                                            </Typography>
                                                                                            <Typography variant="caption" color="text.secondary">
                                                                                                {lesson.room_name || 'Chưa có phòng'}
                                                                                            </Typography>
                                                                                            {renderConflictChips(derivedPreview.existingLessonConflictTags[lesson.lesson_id] || [])}
                                                                                        </Stack>
                                                                                    </Paper>
                                                                                </Tooltip>
                                                                            ))}

                                                                            {rowAssignments.map((assignment) => {
                                                                                const session = derivedPreview.sessions.find((item) => item.variableId === assignment.variable_id);
                                                                                const isManual = !!manualSelections[assignment.variable_id];
                                                                                const isEditingThis = activeSession?.variableId === assignment.variable_id;
                                                                                return (
                                                                                    <Tooltip
                                                                                        key={assignment.variable_id}
                                                                                        title={`${formatDateTime(assignment.start_time)} - ${format(parseISO(assignment.end_time), 'HH:mm', { locale: vi })} • ${assignment.room_name} • ${assignment.shift_name || 'Ca tự do'}`}
                                                                                        arrow
                                                                                    >
                                                                                        <Paper
                                                                                            variant="outlined"
                                                                                            sx={{
                                                                                                p: 1,
                                                                                                borderRadius: 2,
                                                                                                backgroundColor: isEditingThis ? 'rgba(14,165,233,0.12)' : isManual ? 'rgba(56,189,248,0.08)' : 'rgba(15,118,110,0.06)',
                                                                                                borderColor: isEditingThis ? 'primary.main' : session?.derivedConflictMessages.length ? 'warning.main' : 'rgba(15,118,110,0.2)',
                                                                                                borderWidth: isEditingThis ? 2 : 1,
                                                                                            }}
                                                                                        >
                                                                                            <Stack spacing={0.5}>
                                                                                                <Typography variant="body2" sx={{ fontWeight: 700 }}>
                                                                                                    {assignment.class_code} • Buổi {assignment.session_index}/{assignment.session_total}
                                                                                                </Typography>
                                                                                                <Typography variant="caption" color="text.secondary">
                                                                                                    {assignment.room_name}
                                                                                                </Typography>
                                                                                                <Stack direction="row" spacing={0.75} flexWrap="wrap">
                                                                                                    <Chip
                                                                                                        size="small"
                                                                                                        color={assignment.constraint_fit === 'MANUAL_OVERRIDE' ? 'secondary' : 'success'}
                                                                                                        label={assignment.constraint_fit === 'MANUAL_OVERRIDE' ? 'Chỉnh tay' : 'Theo hệ thống'}
                                                                                                    />
                                                                                                </Stack>
                                                                                                {renderConflictChips(session?.conflictTags || [])}
                                                                                                <Button
                                                                                                    size="small"
                                                                                                    color={isEditingThis ? 'error' : 'primary'}
                                                                                                    startIcon={<EditCalendarRounded />}
                                                                                                    onClick={() => toggleEditMode(assignment.variable_id)}
                                                                                                >
                                                                                                    {isEditingThis ? 'Hủy đổi' : 'Đổi chỗ'}
                                                                                                </Button>
                                                                                            </Stack>
                                                                                        </Paper>
                                                                                    </Tooltip>
                                                                                );
                                                                            })}
                                                                        </Stack>
                                                                    </Box>
                                                                );
                                                            })}
                                                        </Box>
                                                    ))}
                                                </Box>
                                            </Box>
                                        </Paper>
                                    ))}
                                </Stack>
                            </Stack>
                        </Paper>

                        <Stack direction={{ xs: 'column', xl: 'row' }} spacing={2.5}>
                            <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 4, flex: 1.1 }}>
                                <Stack direction="row" spacing={1} alignItems="center" sx={{ mb: 1.5 }}>
                                    <Typography variant="h6" sx={{ fontWeight: 700 }}>
                                        Ca cần xử lý
                                    </Typography>
                                    <Chip size="small" label={`${actionableSessions.length} ca`} variant="outlined" />
                                </Stack>
                                {activeSession ? (
                                    <Stack spacing={1.25} sx={{ mb: 2 }}>
                                        <Stack direction={{ xs: 'column', md: 'row' }} spacing={1} alignItems={{ xs: 'flex-start', md: 'center' }}>
                                            <Typography variant="body2" sx={{ fontWeight: 700 }}>
                                                Đang chọn vị trí cho {activeSession.classCode} • Buổi {activeSession.sessionIndex}/{activeSession.sessionTotal}
                                            </Typography>
                                            <Chip size="small" label={`${filteredActiveSessionOptions.length} vị trí hợp lệ`} color="primary" variant="outlined" />
                                        </Stack>
                                        {activeSessionIssueSummary?.primary ? (
                                            <Typography variant="caption" color="text.secondary">
                                                {activeSessionIssueSummary.primary}
                                                {activeSessionIssueSummary.remainingCount > 0 ? ` • +${activeSessionIssueSummary.remainingCount} vấn đề` : ''}
                                            </Typography>
                                        ) : null}
                                        <Stack direction={{ xs: 'column', md: 'row' }} spacing={1.25}>
                                            <TextField
                                                select
                                                label="Chọn ca"
                                                size="small"
                                                value={selectedShiftId}
                                                onChange={(event) => {
                                                    setSelectedShiftId(event.target.value);
                                                    setSelectedRoomId('');
                                                }}
                                                sx={{ minWidth: 220 }}
                                            >
                                                {activeShiftOptions.map((option) => (
                                                    <MenuItem key={option.value} value={option.value}>
                                                        {option.label}
                                                    </MenuItem>
                                                ))}
                                            </TextField>
                                            <TextField
                                                select
                                                label="Chọn phòng"
                                                size="small"
                                                value={selectedRoomId}
                                                onChange={(event) => setSelectedRoomId(event.target.value)}
                                                sx={{ minWidth: 220 }}
                                            >
                                                <MenuItem value="">Tất cả phòng hợp lệ</MenuItem>
                                                {activeRoomOptions.map((option) => (
                                                    <MenuItem key={option.value} value={option.value}>
                                                        {option.label}
                                                    </MenuItem>
                                                ))}
                                            </TextField>
                                        </Stack>
                                    </Stack>
                                ) : null}
                                <Stack spacing={1.25}>
                                    {actionableSessions.length > 0 ? (
                                        actionableSessions.map((session) => {
                                                const availableOptionCount = getAvailableOptionsForSession(session).length;
                                                const issueSummary = summarizeIssueMessages(collectSessionIssueMessages(session));
                                                return (
                                                    <Paper
                                                        key={session.variableId}
                                                        variant="outlined"
                                                        sx={{ p: 1.5, borderRadius: 3, borderColor: 'warning.light' }}
                                                    >
                                                        <Stack spacing={1}>
                                                            <Stack direction={{ xs: 'column', md: 'row' }} justifyContent="space-between" spacing={1}>
                                                                <Box>
                                                                    <Typography variant="body2" sx={{ fontWeight: 700 }}>
                                                                        {session.classCode} - {session.className}
                                                                    </Typography>
                                                                    <Typography variant="caption" color="text.secondary">
                                                                        Buổi {session.sessionIndex}/{session.sessionTotal}
                                                                    </Typography>
                                                                </Box>
                                                                <Stack direction="row" spacing={1} flexWrap="wrap">
                                                                    <Chip
                                                                        size="small"
                                                                        color={session.resolvedAssignment ? 'warning' : 'error'}
                                                                        label={session.resolvedAssignment ? 'Đã xếp nhưng còn xung đột' : 'Chưa xếp được'}
                                                                    />
                                                                    <Chip
                                                                        size="small"
                                                                        variant="outlined"
                                                                        label={`${availableOptionCount} phương án trống`}
                                                                    />
                                                                </Stack>
                                                            </Stack>
                                                            {renderConflictChips(session.conflictTags)}
                                                            {session.resolvedAssignment ? (
                                                                <Typography variant="caption" color="text.secondary">
                                                                    Đang ở: {formatCompactDateTime(session.resolvedAssignment.start_time)} • {session.resolvedAssignment.room_name} • {session.resolvedAssignment.shift_name || 'Ca tự do'}
                                                                </Typography>
                                                            ) : null}
                                                            {issueSummary.primary ? (
                                                                <Stack direction="row" spacing={1} alignItems="center" flexWrap="wrap" useFlexGap>
                                                                    <Typography variant="caption" sx={{ color: 'warning.main' }}>
                                                                        {issueSummary.primary}
                                                                    </Typography>
                                                                    {issueSummary.remainingCount > 0 ? (
                                                                        <Chip size="small" variant="outlined" label={`+${issueSummary.remainingCount} vấn đề`} />
                                                                    ) : null}
                                                                </Stack>
                                                            ) : null}
                                                            <Button
                                                                variant={activeSessionId === session.variableId ? 'contained' : 'outlined'}
                                                                color={activeSessionId === session.variableId ? 'error' : 'primary'}
                                                                startIcon={<EditCalendarRounded />}
                                                                onClick={() => toggleEditMode(session.variableId)}
                                                                disabled={availableOptionCount === 0}
                                                            >
                                                                {availableOptionCount > 0
                                                                    ? activeSessionId === session.variableId
                                                                        ? 'Đóng chọn'
                                                                        : 'Chọn vị trí'
                                                                    : 'Không còn phương án trống'}
                                                            </Button>
                                                        </Stack>
                                                    </Paper>
                                                );
                                            })
                                    ) : (
                                        <Typography variant="body2" color="text.secondary">
                                            Không còn ca học nào cần xử lý.
                                        </Typography>
                                    )}
                                </Stack>
                            </Paper>

                            <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 4, flex: 1 }}>
                                <Stack direction="row" spacing={1} alignItems="center" sx={{ mb: 1.5 }}>
                                    <Typography variant="h6" sx={{ fontWeight: 700 }}>
                                        Lịch còn trùng
                                    </Typography>
                                    <Chip size="small" label={`${derivedPreview.remainingConflicts.length} mục`} variant="outlined" />
                                </Stack>
                                <Stack spacing={1.25}>
                                    {derivedPreview.remainingConflicts.length > 0 ? (
                                        derivedPreview.remainingConflicts.map((conflict, index) => (
                                            <Paper
                                                key={`${conflict.variable_id || 'global'}-${conflict.type}-${index}`}
                                                variant="outlined"
                                                sx={{
                                                    p: 1.5,
                                                    borderRadius: 3,
                                                    borderColor: getConflictSeverity(conflict.type) === 'error' ? 'error.light' : 'warning.light',
                                                }}
                                            >
                                                <Stack direction="row" spacing={1} alignItems="center" flexWrap="wrap" useFlexGap sx={{ mb: 0.75 }}>
                                                    <Typography variant="body2" sx={{ fontWeight: 700 }}>
                                                        {getConflictScopeLabel(conflict.class_code, conflict.class_name, conflict.session_index, conflict.session_total)}
                                                    </Typography>
                                                    {renderConflictChips(buildConflictTags(conflict.type, conflict.message))}
                                                </Stack>
                                                <Typography variant="caption" color="text.secondary">
                                                    {conflict.message}
                                                </Typography>
                                            </Paper>
                                        ))
                                    ) : (
                                        <Typography variant="body2" color="text.secondary">
                                            Không còn lịch trùng.
                                        </Typography>
                                    )}
                                </Stack>
                            </Paper>
                        </Stack>
                    </Stack>
                ) : null}
            </Stack>

        </Box>
    );
};
