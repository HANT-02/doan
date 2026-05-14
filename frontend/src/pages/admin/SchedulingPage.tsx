import { useEffect, useMemo, useState } from 'react';
import { Controller, useForm, useWatch } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import {
    Alert,
    Autocomplete,
    Box,
    Button,
    Chip,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    Divider,
    InputAdornment,
    MenuItem,
    Paper,
    Skeleton,
    Stack,
    Step,
    StepLabel,
    Stepper,
    TextField,
    Typography,
} from '@mui/material';
import {
    CalendarMonthRounded,
    EditCalendarRounded,
    PlayArrowRounded,
    RefreshRounded,
    RuleRounded,
    SearchRounded,
    TuneRounded,
} from '@mui/icons-material';
import { addDays, eachWeekOfInterval, endOfWeek, format, isSameDay, parseISO, startOfWeek } from 'date-fns';
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

const deriveDefaultEndDate = (classes: MinimalScheduleClass[], dateFrom: string) => {
    const dates = classes
        .map((item) => item.end_date)
        .filter((value): value is string => !!value)
        .sort();
    if (dates.length > 0) {
        return dates[dates.length - 1].slice(0, 10);
    }
    return formatDateInput(addDays(parseISO(dateFrom), 90));
};

const formatDateTime = (value: string) =>
    format(parseISO(value), 'dd/MM/yyyy HH:mm', { locale: vi });

const formatCompactDateTime = (value: string) =>
    format(parseISO(value), 'dd/MM HH:mm', { locale: vi });

const formatWeekday = (value: string) => {
    const date = parseISO(value);
    return dayLabels[(date.getDay() + 6) % 7];
};

const formatOptionTitle = (option: SchedulingCandidateOption) =>
    `${formatWeekday(option.start_time)}, ${format(parseISO(option.start_time), 'dd/MM/yyyy', { locale: vi })}`;

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

const getConflictActionHint = (type: string) => {
    switch (type) {
        case 'ASSIGNMENT_CONFLICT':
            return 'Gợi ý: chọn lại slot/phòng cho một trong các buổi đang trùng ngay trong workspace bên dưới.';
        case 'SYSTEM_LESSON_CONFLICT':
            return 'Gợi ý: xem card lesson đã lưu trên calendar rồi đổi tay sang ca/ngày/phòng khác để tránh đè dữ liệu hiện có.';
        case 'MISSING_TEACHER':
            return 'Gợi ý: vào Quản lý lớp học để gán giáo viên phụ trách rồi chạy lại xem trước.';
        case 'SKILL_MISMATCH':
            return 'Gợi ý: cập nhật kỹ năng/chứng chỉ cho giáo viên hoặc chỉnh `required skills` của khóa học để chỉ các lớp đủ chuẩn mới được đưa vào xếp lịch.';
        case 'NO_CLASS_INPUT':
            return 'Gợi ý: kiểm tra bộ lọc lớp, trạng thái OPEN hoặc giáo viên đã chọn.';
        case 'NO_ACTIVE_ROOM':
            return 'Gợi ý: bổ sung phòng học khả dụng hoặc kiểm tra phòng đã gán trong lịch tuần lớp.';
        case 'PREFERRED_ROOM_UNAVAILABLE':
            return 'Gợi ý: đổi phòng đang gán trong lịch tuần lớp hoặc chọn phòng trống khi chỉnh tay.';
        case 'ROOM_CAPACITY_BLOCK':
            return 'Gợi ý: chọn phòng sức chứa lớn hơn hoặc điều chỉnh sĩ số tối đa của lớp.';
        case 'NO_SLOT_IN_RANGE':
            return 'Gợi ý: đổi ngày bắt đầu dự kiến hoặc bổ sung thêm ngày/ca trong lịch tuần lớp.';
        case 'NO_DOMAIN':
            return 'Gợi ý: thử đổi sang ngày, ca hoặc phòng trống khác trong dialog chỉnh tay.';
        case 'MISSING_COURSE':
            return 'Gợi ý: gắn khóa học cho lớp trước khi chạy xem trước xếp lịch để hệ thống sinh đúng số buổi học.';
        case 'INVALID_COURSE_SESSION_COUNT':
            return 'Gợi ý: cập nhật `số ca` của khóa học thành số buổi hợp lệ.';
        case 'INVALID_COURSE_DURATION':
            return 'Gợi ý: cập nhật `thời lượng ca học` của khóa học để xem trước dùng đúng thời lượng buổi học.';
        case 'CLASS_SCHEDULE_NO_SLOT':
            return 'Gợi ý: kiểm tra lịch mẫu của lớp, `mã ca` đã gán và bảo đảm ca học đủ chứa thời lượng buổi học từ khóa học.';
        case 'CLASS_SCHEDULE_ROOM_UNAVAILABLE':
            return 'Gợi ý: kiểm tra `mã phòng` trong lịch mẫu của lớp hoặc nới lại bộ lọc phòng để giữ đúng phòng cố định theo lịch mẫu.';
        case 'MISSING_CLASS_SCHEDULE':
            return 'Gợi ý: vào Quản lý lớp học, mở tab lịch tuần và cấu hình ngày học cùng ca học cho lớp trước khi chạy xem trước lịch.';
        case 'INSUFFICIENT_SCHEDULE_SLOTS':
            return 'Gợi ý: tăng khoảng ngày xem trước hoặc bổ sung thêm ngày/ca trong lịch tuần lớp để đủ số slot hợp lệ.';
        case 'INSUFFICIENT_ENROLLMENT':
            return 'Gợi ý: chỉ đưa lớp vào xếp lịch khi lớp đã đạt ngưỡng sĩ số tối thiểu hoặc điều chỉnh lại chính sách mở lớp.';
        case 'EXCESSIVE_MANUAL_ADJUSTMENT':
            return 'Gợi ý: số ca cần sửa tay đang quá nhiều. Hãy tối ưu dữ liệu lớp, phòng, ca và lịch tuần trước khi chạy lại solver.';
        case 'NO_ACTIVE_SHIFT':
            return 'Gợi ý: vào Quản lý ca học để tạo hoặc kích hoạt ít nhất một `Ca học` rồi chạy lại xem trước.';
        case 'NO_VALID_DATE_RANGE':
            return 'Gợi ý: chọn ngày kết thúc lớn hơn hoặc bằng ngày bắt đầu.';
        default:
            return 'Gợi ý: rà lại dữ liệu lớp, giáo viên, phòng và bộ lọc preview.';
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
    parseISO(left.start_time) < parseISO(right.end_time) &&
    parseISO(right.start_time) < parseISO(left.end_time);

const overlapsExistingLesson = (assignment: SchedulingAssignment, lesson: SchedulingExistingLesson) =>
    parseISO(assignment.start_time) < parseISO(lesson.end_time) &&
    parseISO(lesson.start_time) < parseISO(assignment.end_time);

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
            messages.push(`trùng lớp với lesson đã lưu ${lesson.class_code}`);
        }
        if (assignment.teacher_id && lesson.teacher_id && assignment.teacher_id === lesson.teacher_id) {
            messages.push(`trùng giáo viên với lesson đã lưu ${lesson.class_code}`);
        }
        if (assignment.room_id && lesson.room_id && assignment.room_id === lesson.room_id) {
            messages.push(`phòng ${option.room_name} đã có lesson ${lesson.class_code}`);
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

        const currentEnd = parseISO(current.end_time);
        const nextStart = parseISO(next.start_time);
        if (!isSameDay(currentEnd, nextStart)) {
            continue;
        }

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

const buildWeeks = (preview: SchedulingPreview) => {
    const dateFrom = preview.filters.date_from || fallbackDateFrom;
    const dateTo = preview.filters.date_to || dateFrom;
    const start = startOfWeek(parseISO(dateFrom), { weekStartsOn: 1 });
    const end = endOfWeek(parseISO(dateTo), { weekStartsOn: 1 });
    return eachWeekOfInterval({ start, end }, { weekStartsOn: 1 }).map((weekStart) =>
        Array.from({ length: 7 }, (_, index) => addDays(weekStart, index)),
    );
};

export const SchedulingPage = () => {
    const [preview, setPreview] = useState<SchedulingPreview | null>(null);
    const [wizardStep, setWizardStep] = useState(0);
    const [search, setSearch] = useState('');
    const [manualSelections, setManualSelections] = useState<Record<string, string>>({});
    const [activeSessionId, setActiveSessionId] = useState<string | null>(null);
    const [pendingOptionKey, setPendingOptionKey] = useState('');
    const [draggedSessionId, setDraggedSessionId] = useState<string | null>(null);

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
    const derivedDateTo = deriveDefaultEndDate(schedulingClasses, derivedDateFrom);
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
                    `${reasons.join(', ')} với lesson đã lưu ${lesson.class_code} lúc ${formatCompactDateTime(lesson.start_time)}`,
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

    const filteredSessions = useMemo(() => {
        if (!derivedPreview) {
            return [];
        }

        const keyword = search.trim().toLowerCase();
        if (!keyword) {
            return derivedPreview.sessions;
        }

        return derivedPreview.sessions.filter((session) => {
            const bag = [
                session.classCode,
                session.className,
                session.teacherLabel,
                session.resolvedAssignment?.room_name,
                session.resolvedAssignment?.shift_name,
            ]
                .filter((value): value is string => typeof value === 'string' && value.length > 0)
                .join(' ')
                .toLowerCase();
            return bag.includes(keyword);
        });
    }, [derivedPreview, search]);

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
    const activeSessionAvailableOptions = useMemo(() => {
        if (!activeSession || !derivedPreview) {
            return [];
        }
        return filterAvailableCandidateOptions(activeSession, derivedPreview);
    }, [activeSession, derivedPreview]);
    useEffect(() => {
        if (!activeSession) {
            return;
        }
        if (activeSessionAvailableOptions.some((option) => option.key === pendingOptionKey)) {
            return;
        }
        setPendingOptionKey(activeSessionAvailableOptions[0]?.key || '');
    }, [activeSession, activeSessionAvailableOptions, pendingOptionKey]);

    const activeStep = wizardStep;
    const canCommitPreview = !!derivedPreview && derivedPreview.status === 'COMPLETED' && !derivedPreview.summary.excessiveManualAdjustment;

    const getAvailableOptionsForSession = (session: SessionDraft) =>
        derivedPreview ? filterAvailableCandidateOptions(session, derivedPreview) : [];

    const resetManualSelections = () => {
        setManualSelections({});
        setActiveSessionId(null);
        setPendingOptionKey('');
        setDraggedSessionId(null);
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
                date_to: deriveDefaultEndDate(schedulingClasses, dateFrom),
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

    const openEditDialog = (variableId: string) => {
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
        setActiveSessionId(variableId);
        setPendingOptionKey(initialKey);
    };

    const handleSaveManualSelection = () => {
        if (!activeSessionId || !pendingOptionKey) {
            setActiveSessionId(null);
            setPendingOptionKey('');
            return;
        }
        if (!activeSessionAvailableOptions.some((option) => option.key === pendingOptionKey)) {
            toast.error('Phương án đã chọn không còn trống. Hãy chọn lại ngày, ca hoặc phòng khác.');
            return;
        }

        setManualSelections((current) => ({
            ...current,
            [activeSessionId]: pendingOptionKey,
        }));
        setActiveSessionId(null);
        setPendingOptionKey('');
    };

    const handleClearSessionSelection = () => {
        if (!activeSessionId) {
            return;
        }

        setManualSelections((current) => {
            const next = { ...current };
            delete next[activeSessionId];
            return next;
        });
        setActiveSessionId(null);
        setPendingOptionKey('');
    };

    const handleApplyManualOption = (variableId: string, optionKey: string) => {
        setManualSelections((current) => ({
            ...current,
            [variableId]: optionKey,
        }));
        setDraggedSessionId(null);
    };

    const draggedSession = useMemo(
        () => derivedPreview?.sessions.find((session) => session.variableId === draggedSessionId) || null,
        [derivedPreview, draggedSessionId],
    );

    const handleCommitPreview = async () => {
        if (!preview?.run_id || !derivedPreview) {
            toast.error('Chưa có preview nào để commit');
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
            toast.error(getApiErrorMessage(error, 'Không thể xác nhận xem trước xếp lịch'));
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
                subtitle="Chạy xem trước, xem lịch theo dạng lịch và chỉnh tay các ca học bị trùng trước khi xác nhận lưu."
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
                                <Typography variant="body2" color="text.secondary">
                                    Chỉ cần chọn lớp hoặc giáo viên. Phòng, ngày học, ca học và giáo viên phụ trách sẽ lấy từ lịch tuần đã cấu hình trong từng lớp.
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
                                                    helperText={errors.class_ids?.message || 'Có thể chọn một hoặc nhiều lớp OPEN cần xếp lịch.'}
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
                                                    helperText={errors.teacher_ids?.message || 'Hệ thống sẽ lấy các lớp OPEN mà giáo viên đang phụ trách.'}
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
                                        helperText={`Để trống sẽ dùng ngày bắt đầu lớp: ${format(parseISO(deriveDefaultStartDate(schedulingClasses)), 'dd/MM/yyyy', { locale: vi })}. Ngày kết thúc tự lấy theo cài đặt lớp.`}
                                    />
                                )}
                            />

                            <Alert severity="info">
                                Bộ lọc phòng và khoảng thời gian đã được bỏ khỏi màn này để tránh nhập trùng dữ liệu. Nếu muốn đổi phòng/ca, hãy sửa lịch tuần của lớp hoặc chỉnh tay trong màn preview.
                            </Alert>

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
                                <Typography variant="body2" color="text.secondary">
                                    Hệ thống sẽ sinh số buổi cần xếp dựa trên khóa học, lịch tuần và khoảng ngày của lớp.
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
                                        Khoảng ngày hệ thống dùng
                                    </Typography>
                                    <Typography variant="body2" color="text.secondary">
                                        Từ {format(parseISO(derivedDateFrom), 'dd/MM/yyyy', { locale: vi })} đến {format(parseISO(derivedDateTo), 'dd/MM/yyyy', { locale: vi })}
                                    </Typography>
                                    <Typography variant="caption" color="text.secondary">
                                        Ngày bắt đầu lấy từ lựa chọn dự kiến hoặc start date của lớp; ngày kết thúc lấy từ end date của lớp.
                                    </Typography>
                                </Paper>
                            </Stack>

                            <Alert severity={schedulingClasses.length > 0 ? 'success' : 'warning'}>
                                {schedulingClasses.length > 0
                                    ? 'Cấu hình đã sẵn sàng. Preview sẽ dùng lịch tuần lớp để chọn ngày, ca và phòng.'
                                    : 'Chưa tìm thấy lớp OPEN phù hợp. Bạn có thể quay lại đổi lớp hoặc giáo viên.'}
                            </Alert>

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
                            Chưa có preview nào
                        </Typography>
                        <Typography variant="body2" color="text.secondary" sx={{ mt: 1, mb: 3 }}>
                            Bấm `Chạy xếp lịch` để tạo preview đầu tiên hoặc tải preview gần nhất từ backend.
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
                                        Chạy xem trước: {preview.run_id}
                                    </Typography>
                                    <Typography variant="body2" color="text.secondary">
                                        Sinh lúc {formatDateTime(preview.generated_at)} • Trạng thái hiện tại {derivedPreview.status}
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
                                        label={`Delta ${formatScoreDelta(derivedPreview.summary.scoreDelta)}`}
                                        color={getScoreDeltaTone(derivedPreview.summary.scoreDelta)}
                                        variant="outlined"
                                    />
                                    <Chip
                                        label={`Chỉnh tay ${derivedPreview.summary.manualAdjustmentCount}/${derivedPreview.summary.manualAdjustmentLimit}`}
                                        color={derivedPreview.summary.excessiveManualAdjustment ? 'error' : 'secondary'}
                                        variant="outlined"
                                    />
                                </Stack>
                            </Stack>
                        </Paper>

                        <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 4 }}>
                            <Typography variant="h6" sx={{ fontWeight: 700, mb: 1.5 }}>
                                Thông tin lớp trong preview
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

                        {derivedPreview.remainingConflicts.length > 0 ? (
                            <Alert severity={derivedPreview.calendarAssignments.length > 0 ? 'warning' : 'error'}>
                                Lịch xem trước còn {derivedPreview.remainingConflicts.length} vấn đề cần xử lý. Bạn có thể đổi chỗ/phòng cho từng ca học ngay hoặc quay lại sửa dữ liệu lớp nếu trùng lịch là do thiếu giáo viên/lịch tuần.
                            </Alert>
                        ) : (
                            <Alert severity="success">Lịch xem trước hiện không còn lịch trùng.</Alert>
                        )}

                        {derivedPreview.summary.excessiveManualAdjustment ? (
                            <Alert severity="warning">
                                Số ca chỉnh tay hiện tại là {derivedPreview.summary.manualAdjustmentCount}, vượt ngưỡng cho phép {derivedPreview.summary.manualAdjustmentLimit}. Hãy tối ưu lại dữ liệu đầu vào thay vì vá tay quá nhiều.
                            </Alert>
                        ) : null}

                        {derivedPreview.summary.scoreDelta < 0 ? (
                            <Alert severity={derivedPreview.summary.scoreDelta < -10 ? 'warning' : 'info'}>
                                Điểm luật mềm đang thấp hơn {Math.abs(derivedPreview.summary.scoreDelta)} điểm so với kết quả solver ban đầu ({derivedPreview.summary.baseSoftScore}). Nếu đây không phải phương án bắt buộc, nên chọn slot khác ít ảnh hưởng chất lượng lịch hơn.
                            </Alert>
                        ) : null}

                        <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 4 }}>
                            <Stack
                                direction={{ xs: 'column', lg: 'row' }}
                                spacing={1.5}
                                justifyContent="space-between"
                                alignItems={{ xs: 'stretch', lg: 'center' }}
                                sx={{ mb: 2 }}
                            >
                                <TextField
                                    size="small"
                                    value={search}
                                    onChange={(event) => setSearch(event.target.value)}
                                    placeholder="Tìm lớp, giáo viên, phòng, ca học..."
                                    sx={{ minWidth: { xs: '100%', md: 320 } }}
                                    InputProps={{
                                        startAdornment: (
                                            <InputAdornment position="start">
                                                <SearchRounded fontSize="small" />
                                            </InputAdornment>
                                        ),
                                    }}
                                />
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

                            {canCommitPreview ? (
                                <Alert severity="success" sx={{ mb: 2 }}>
                                    Lịch hiện đã đạt trạng thái `COMPLETED`. Bạn có thể commit để tạo `lesson` thật; backend sẽ kiểm tra lần cuối trước khi ghi dữ liệu.
                                </Alert>
                            ) : (
                                <Alert severity="info" sx={{ mb: 2 }}>
                                    Xác nhận lưu chỉ mở khi tất cả ca học đã được xếp, không còn trùng lịch và số chỉnh tay không vượt ngưỡng cho phép. Hãy xử lý các ca học chưa thỏa mãn ở bảng và lịch bên dưới.
                                </Alert>
                            )}

                            <Stack spacing={2.5}>
                                <Stack direction="row" spacing={1} alignItems="center">
                                    <CalendarMonthRounded color="primary" />
                                    <Box>
                                        <Typography variant="h6" sx={{ fontWeight: 700 }}>
                                            Xem trước lịch
                                        </Typography>
                                        <Typography variant="body2" color="text.secondary">
                                            Các buổi đã xếp được hiển thị theo tuần. Bấm vào từng card để đổi chỗ/phòng ngay tại chỗ.
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
                                            <Box
                                                sx={{
                                                    display: 'grid',
                                                    gridTemplateColumns: { xs: '1fr', lg: 'repeat(7, minmax(0, 1fr))' },
                                                }}
                                            >
                                                {week.map((day, dayIndex) => {
                                                    const dayAssignments = derivedPreview.calendarAssignments.filter((assignment) =>
                                                        isSameDay(parseISO(assignment.start_time), day),
                                                    );
                                                    const dayExistingLessons = derivedPreview.existingLessons.filter((lesson) =>
                                                        isSameDay(parseISO(lesson.start_time), day),
                                                    );
                                                    const dayDropOptions = draggedSession
                                                        ? getAvailableOptionsForSession(draggedSession).filter((option) =>
                                                            isSameDay(parseISO(option.start_time), day),
                                                        )
                                                        : [];

                                                    return (
                                                        <Box
                                                            key={day.toISOString()}
                                                            sx={{
                                                                p: 1.5,
                                                                minHeight: 180,
                                                                borderLeft: { xs: 'none', lg: dayIndex === 0 ? 'none' : '1px solid rgba(15,23,42,0.08)' },
                                                                borderTop: '1px solid rgba(15,23,42,0.08)',
                                                            }}
                                                        >
                                                            <Stack spacing={1}>
                                                                <Box>
                                                                    <Typography variant="subtitle2" sx={{ fontWeight: 700 }}>
                                                                        {dayLabels[dayIndex]}
                                                                    </Typography>
                                                                    <Typography variant="caption" color="text.secondary">
                                                                        {format(day, 'dd/MM/yyyy', { locale: vi })}
                                                                    </Typography>
                                                                </Box>

                                                                {draggedSession && dayDropOptions.length > 0 ? (
                                                                    <Paper
                                                                        variant="outlined"
                                                                        sx={{
                                                                            p: 1,
                                                                            borderRadius: 2,
                                                                            borderStyle: 'dashed',
                                                                            borderColor: 'secondary.main',
                                                                            backgroundColor: 'rgba(168,85,247,0.04)',
                                                                        }}
                                                                        onDragOver={(event) => event.preventDefault()}
                                                                    >
                                                                        <Stack spacing={1}>
                                                                            <Typography variant="caption" sx={{ fontWeight: 700 }}>
                                                                                Thả {draggedSession.classCode} • Buổi {draggedSession.sessionIndex}/{draggedSession.sessionTotal} vào một slot trống:
                                                                            </Typography>
                                                                            <Stack direction="row" spacing={0.75} flexWrap="wrap" useFlexGap>
                                                                                {dayDropOptions.map((option) => {
                                                                                    const delta = getOptionScoreDelta(draggedSession, option, derivedPreview);
                                                                                    return (
                                                                                        <Chip
                                                                                            key={`${draggedSession.variableId}-${option.key}`}
                                                                                            clickable
                                                                                            color={getScoreDeltaTone(delta)}
                                                                                            variant="outlined"
                                                                                            label={`${option.shift_name || option.shift_code || 'Ca'} • ${format(parseISO(option.start_time), 'HH:mm', { locale: vi })} • ${option.room_name} • Δ ${formatScoreDelta(delta)}`}
                                                                                            onClick={() => handleApplyManualOption(draggedSession.variableId, option.key)}
                                                                                            onDragOver={(event) => event.preventDefault()}
                                                                                            onDrop={(event) => {
                                                                                                event.preventDefault();
                                                                                                handleApplyManualOption(draggedSession.variableId, option.key);
                                                                                            }}
                                                                                            sx={{ cursor: 'copy' }}
                                                                                        />
                                                                                    );
                                                                                })}
                                                                            </Stack>
                                                                        </Stack>
                                                                    </Paper>
                                                                ) : null}

                                                                {dayExistingLessons.length > 0 ? (
                                                                    <Stack spacing={1}>
                                                                        {dayExistingLessons.map((lesson) => (
                                                                            <Paper
                                                                                key={lesson.lesson_id}
                                                                                variant="outlined"
                                                                                sx={{
                                                                                    p: 1.25,
                                                                                    borderRadius: 2.5,
                                                                                    backgroundColor: 'rgba(15,23,42,0.04)',
                                                                                    borderColor: derivedPreview.existingLessonConflictTags[lesson.lesson_id]?.length
                                                                                        ? 'error.light'
                                                                                        : 'rgba(15,23,42,0.12)',
                                                                                }}
                                                                            >
                                                                                <Stack spacing={0.5}>
                                                                                    <Typography variant="body2" sx={{ fontWeight: 700 }}>
                                                                                        {lesson.class_code} • Lesson đã lưu
                                                                                    </Typography>
                                                                                    <Typography variant="caption" color="text.secondary">
                                                                                        {formatCompactDateTime(lesson.start_time)} - {format(parseISO(lesson.end_time), 'HH:mm', { locale: vi })}
                                                                                    </Typography>
                                                                                    <Typography variant="caption" color="text.secondary">
                                                                                        {lesson.room_name || 'Chưa có phòng'} • {lesson.teacher_label || 'Chưa có giáo viên'}
                                                                                    </Typography>
                                                                                    {renderConflictChips(derivedPreview.existingLessonConflictTags[lesson.lesson_id] || [])}
                                                                                </Stack>
                                                                            </Paper>
                                                                        ))}
                                                                    </Stack>
                                                                ) : null}

                                                                {dayAssignments.length > 0 ? (
                                                                    <Stack spacing={1}>
                                                                        {dayAssignments.map((assignment) => {
                                                                            const session = derivedPreview.sessions.find((item) => item.variableId === assignment.variable_id);
                                                                            const isManual = !!manualSelections[assignment.variable_id];
                                                                            return (
                                                                                <Paper
                                                                                    key={assignment.variable_id}
                                                                                    variant="outlined"
                                                                                    sx={{
                                                                                        p: 1.25,
                                                                                        borderRadius: 2.5,
                                                                                        backgroundColor: isManual ? 'rgba(56,189,248,0.08)' : 'rgba(15,118,110,0.06)',
                                                                                        borderColor: session?.derivedConflictMessages.length ? 'warning.main' : 'rgba(15,118,110,0.2)',
                                                                                    }}
                                                                                >
                                                                                    <Stack spacing={0.75}>
                                                                                        <Typography variant="body2" sx={{ fontWeight: 700 }}>
                                                                                            {assignment.class_code} • Buổi {assignment.session_index}/{assignment.session_total}
                                                                                        </Typography>
                                                                                        <Typography variant="caption" color="text.secondary">
                                                                                            {formatCompactDateTime(assignment.start_time)} - {format(parseISO(assignment.end_time), 'HH:mm', { locale: vi })}
                                                                                        </Typography>
                                                                                        <Typography variant="caption" color="text.secondary">
                                                                                            {assignment.room_name} • {assignment.shift_name || 'Ca tự do'}
                                                                                        </Typography>
                                                                                        <Stack direction="row" spacing={1} flexWrap="wrap">
                                                                                            <Chip
                                                                                                size="small"
                                                                                                color={assignment.constraint_fit === 'MANUAL_OVERRIDE' ? 'secondary' : 'success'}
                                                                                                label={assignment.constraint_fit === 'MANUAL_OVERRIDE' ? 'Chỉnh tay' : 'Theo hệ thống'}
                                                                                            />
                                                                                        </Stack>
                                                                                        {renderConflictChips(session?.conflictTags || [])}
                                                                                        <Button
                                                                                            size="small"
                                                                                            startIcon={<EditCalendarRounded />}
                                                                                            onClick={() => openEditDialog(assignment.variable_id)}
                                                                                        >
                                                                                            Đổi chỗ
                                                                                        </Button>
                                                                                        <Chip
                                                                                            size="small"
                                                                                            draggable
                                                                                            label="Kéo thả để đổi ca"
                                                                                            onDragStart={() => setDraggedSessionId(assignment.variable_id)}
                                                                                            onDragEnd={() => setDraggedSessionId(null)}
                                                                                            sx={{ alignSelf: 'flex-start', cursor: 'grab' }}
                                                                                        />
                                                                                    </Stack>
                                                                                </Paper>
                                                                            );
                                                                        })}
                                                                    </Stack>
                                                                ) : !dayExistingLessons.length ? (
                                                                    <Typography variant="caption" color="text.secondary">
                                                                        Không có buổi nào được xếp trong ngày này.
                                                                    </Typography>
                                                                ) : null}
                                                            </Stack>
                                                        </Box>
                                                    );
                                                })}
                                            </Box>
                                        </Paper>
                                    ))}
                                </Stack>
                            </Stack>
                        </Paper>

                        <Stack direction={{ xs: 'column', xl: 'row' }} spacing={2.5}>
                            <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 4, flex: 1.1 }}>
                                <Typography variant="h6" sx={{ fontWeight: 700, mb: 1.5 }}>
                                    Ca học cần xử lý
                                </Typography>
                                <Stack spacing={1.25}>
                                    {filteredSessions.filter((session) => !session.resolvedAssignment || session.derivedConflictMessages.length > 0).length > 0 ? (
                                        filteredSessions
                                            .filter((session) => !session.resolvedAssignment || session.derivedConflictMessages.length > 0)
                                            .map((session) => {
                                                const availableOptionCount = getAvailableOptionsForSession(session).length;
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
                                                                        label={session.resolvedAssignment ? 'Đã xếp nhưng còn conflict' : 'Chưa xếp được'}
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
                                                                    Hiện tại: {formatCompactDateTime(session.resolvedAssignment.start_time)} • {session.resolvedAssignment.room_name} • {session.resolvedAssignment.shift_name || 'Ca tự do'}
                                                                </Typography>
                                                            ) : null}
                                                            {session.baseConflicts.map((conflict) => (
                                                                <Alert key={`${conflict.variable_id}-${conflict.type}`} severity={getConflictSeverity(conflict.type)}>
                                                                    {conflict.message}
                                                                </Alert>
                                                            ))}
                                                            {session.derivedConflictMessages.map((message) => (
                                                                <Alert key={`${session.variableId}-${message}`} severity="warning">
                                                                    {message}
                                                                </Alert>
                                                            ))}
                                                            <Button
                                                                variant="outlined"
                                                                startIcon={<EditCalendarRounded />}
                                                                onClick={() => openEditDialog(session.variableId)}
                                                                disabled={availableOptionCount === 0}
                                                            >
                                                                {availableOptionCount > 0 ? 'Chọn ngày/ca/phòng trống' : 'Không còn phương án trống'}
                                                            </Button>
                                                            {availableOptionCount > 0 ? (
                                                                <Chip
                                                                    size="small"
                                                                    draggable
                                                                    label="Kéo vào lịch tuần để đổi ca"
                                                                    onDragStart={() => setDraggedSessionId(session.variableId)}
                                                                    onDragEnd={() => setDraggedSessionId(null)}
                                                                    sx={{ alignSelf: 'flex-start', cursor: 'grab' }}
                                                                />
                                                            ) : null}
                                                        </Stack>
                                                    </Paper>
                                                );
                                            })
                                    ) : (
                                        <Alert severity="success">
                                            Không còn ca học nào cần xử lý trong bộ lọc hiện tại.
                                        </Alert>
                                    )}
                                </Stack>
                            </Paper>

                            <Paper variant="outlined" sx={{ p: 2.5, borderRadius: 4, flex: 1 }}>
                                <Typography variant="h6" sx={{ fontWeight: 700, mb: 1.5 }}>
                                    Danh sách trùng lịch
                                </Typography>
                                <Stack spacing={1.25}>
                                    {derivedPreview.remainingConflicts.length > 0 ? (
                                        derivedPreview.remainingConflicts.map((conflict, index) => (
                                            <Alert
                                                key={`${conflict.variable_id || 'global'}-${conflict.type}-${index}`}
                                                severity={getConflictSeverity(conflict.type)}
                                            >
                                                {renderConflictChips(buildConflictTags(conflict.type, conflict.message))}
                                                <strong>{getConflictScopeLabel(conflict.class_code, conflict.class_name, conflict.session_index, conflict.session_total)}</strong>
                                                {`: ${conflict.message} `}
                                                <Typography component="span" variant="body2" sx={{ fontWeight: 600 }}>
                                                    {getConflictActionHint(conflict.type)}
                                                </Typography>
                                            </Alert>
                                        ))
                                    ) : (
                                        <Alert severity="success">
                                            Không còn lịch trùng nào sau các chỉnh sửa hiện tại.
                                        </Alert>
                                    )}
                                </Stack>
                            </Paper>
                        </Stack>
                    </Stack>
                ) : null}
            </Stack>

            <Dialog open={!!activeSession} onClose={() => setActiveSessionId(null)} fullWidth maxWidth="sm">
                <DialogTitle>Chỉnh chỗ cho ca học</DialogTitle>
                <DialogContent dividers>
                    {activeSession ? (
                        <Stack spacing={2}>
                            <Box>
                                <Typography variant="subtitle1" sx={{ fontWeight: 700 }}>
                                    {activeSession.classCode} - {activeSession.className}
                                </Typography>
                                <Typography variant="body2" color="text.secondary">
                                    Buổi {activeSession.sessionIndex}/{activeSession.sessionTotal}
                                </Typography>
                                {activeSession.resolvedAssignment ? (
                                    <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                                        Hiện tại: {formatDateTime(activeSession.resolvedAssignment.start_time)} • {activeSession.resolvedAssignment.room_name} • {activeSession.resolvedAssignment.shift_name || 'Ca tự do'}
                                    </Typography>
                                ) : null}
                            </Box>

                            <Divider />

                            <TextField
                                select
                                label="Chọn ngày, ca và phòng trống"
                                value={pendingOptionKey}
                                onChange={(event) => setPendingOptionKey(event.target.value)}
                                fullWidth
                                helperText={
                                    activeSessionAvailableOptions.length > 0
                                        ? 'Danh sách chỉ hiển thị các phương án chưa trùng lớp, giáo viên hoặc phòng trong thời điểm đó.'
                                        : 'Tất cả phương án hiện đang trùng lịch. Hãy nới khoảng ngày hoặc bổ sung ca/phòng khác.'
                                }
                            >
                                {activeSessionAvailableOptions.map((option) => {
                                    const delta = derivedPreview ? getOptionScoreDelta(activeSession, option, derivedPreview) : 0;
                                    return (
                                        <MenuItem key={option.key} value={option.key}>
                                            <Stack
                                                direction="row"
                                                spacing={1}
                                                alignItems="center"
                                                justifyContent="space-between"
                                                sx={{ width: '100%' }}
                                            >
                                                <Box>
                                                    <Typography variant="body2" sx={{ fontWeight: 600 }}>
                                                        {formatOptionTitle(option)}
                                                    </Typography>
                                                    <Typography variant="caption" color="text.secondary">
                                                        {option.shift_name || option.shift_code || 'Ca tự do'} • {format(parseISO(option.start_time), 'HH:mm', { locale: vi })} - {format(parseISO(option.end_time), 'HH:mm', { locale: vi })} • {option.room_name} trống • Sức chứa {option.room_capacity}
                                                    </Typography>
                                                </Box>
                                                <Chip
                                                    size="small"
                                                    color={getScoreDeltaTone(delta)}
                                                    variant="outlined"
                                                    label={`Δ ${formatScoreDelta(delta)}`}
                                                />
                                            </Stack>
                                        </MenuItem>
                                    );
                                })}
                            </TextField>
                        </Stack>
                    ) : null}
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setActiveSessionId(null)}>Đóng</Button>
                    <Button onClick={handleClearSessionSelection} disabled={!manualSelections[activeSessionId || '']}>
                        Bỏ chỉnh ca này
                    </Button>
                    <Button
                        variant="contained"
                        onClick={handleSaveManualSelection}
                        disabled={!activeSession || !pendingOptionKey || activeSessionAvailableOptions.length === 0}
                    >
                        Áp dụng
                    </Button>
                </DialogActions>
            </Dialog>
        </Box>
    );
};
