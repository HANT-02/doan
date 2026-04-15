import {
    AccessTime,
    AssignmentOutlined,
    AssignmentTurnedIn,
    Book,
    CalendarMonth,
    Dashboard,
    DatasetLinked,
    Description,
    EventNote,
    Group,
    ManageAccounts,
    MeetingRoom,
    Person,
    School,
    UploadFile,
} from '@mui/icons-material';
import { hasAnyRole } from '@/utils/roles';

export interface NavItem {
    key: string;
    path?: string;
    labelVi: string;
    labelEn?: string;
    icon: any;
    roles: string[];
    children?: NavItem[];
}

export const NAV_ITEMS: NavItem[] = [
    {
        key: 'admin-overview',
        labelVi: 'Tổng quan',
        labelEn: 'Overview',
        path: '/app/admin/overview',
        icon: Dashboard,
        roles: ['ADMIN', 'SUPER_ADMIN'],
    },
    {
        key: 'admin-basic-data',
        labelVi: 'Thông tin cơ bản',
        icon: DatasetLinked,
        roles: ['ADMIN', 'SUPER_ADMIN'],
        children: [
            {
                key: 'admin-teachers',
                labelVi: 'Giáo viên',
                path: '/app/admin/teachers',
                icon: Group,
                roles: ['ADMIN', 'SUPER_ADMIN'],
            },
            {
                key: 'admin-students',
                labelVi: 'Học sinh',
                path: '/app/admin/students',
                icon: Person,
                roles: ['ADMIN', 'SUPER_ADMIN'],
            },
            {
                key: 'admin-classes',
                labelVi: 'Lớp học',
                path: '/app/admin/classes',
                icon: School,
                roles: ['ADMIN', 'SUPER_ADMIN'],
            },
            {
                key: 'admin-rooms',
                labelVi: 'Phòng học',
                path: '/app/admin/rooms',
                icon: MeetingRoom,
                roles: ['ADMIN', 'SUPER_ADMIN'],
            },
            {
                key: 'admin-shifts',
                labelVi: 'Ca học',
                path: '/app/admin/shifts',
                icon: AccessTime,
                roles: ['ADMIN', 'SUPER_ADMIN'],
            },
            {
                key: 'admin-programs',
                labelVi: 'Chương trình đào tạo',
                path: '/app/admin/programs',
                icon: Book,
                roles: ['ADMIN', 'SUPER_ADMIN'],
            },
            {
                key: 'admin-courses',
                labelVi: 'Khóa học',
                path: '/app/admin/courses',
                icon: Book,
                roles: ['ADMIN', 'SUPER_ADMIN'],
            },
        ],
    },
    {
        key: 'admin-lessons',
        labelVi: 'Quản lý buổi học',
        path: '/app/admin/lessons',
        icon: EventNote,
        roles: ['ADMIN', 'SUPER_ADMIN'],
    },
    {
        key: 'admin-scheduling',
        labelVi: 'Xếp lịch (CSP)',
        path: '/app/admin/scheduling',
        icon: CalendarMonth,
        roles: ['ADMIN', 'SUPER_ADMIN'],
    },
    {
        key: 'admin-predictive',
        labelVi: 'Cảnh báo AT_RISK',
        path: '/app/admin/predictive',
        icon: AssignmentTurnedIn,
        roles: ['ADMIN', 'SUPER_ADMIN'],
    },
    {
        key: 'admin-accounts',
        labelVi: 'Quản lý tài khoản',
        path: '/app/admin/accounts',
        icon: ManageAccounts,
        roles: ['ADMIN', 'SUPER_ADMIN'],
    },
    {
        key: 'admin-leaves',
        labelVi: 'Đơn xin phép',
        path: '/app/admin/leaves',
        icon: Description,
        roles: ['ADMIN', 'SUPER_ADMIN'],
    },

    {
        key: 'teacher-schedule',
        labelVi: 'Lịch giảng dạy',
        path: '/app/teacher/schedule',
        icon: CalendarMonth,
        roles: ['TEACHER'],
    },
    {
        key: 'teacher-attendance',
        labelVi: 'Điểm danh',
        path: '/app/teacher/attendance',
        icon: AssignmentTurnedIn,
        roles: ['TEACHER'],
    },
    {
        key: 'teacher-journal',
        labelVi: 'Sổ đầu bài',
        path: '/app/teacher/journal',
        icon: Book,
        roles: ['TEACHER'],
    },
    {
        key: 'teacher-documents',
        labelVi: 'Tài liệu giảng dạy',
        path: '/app/teacher/documents',
        icon: UploadFile,
        roles: ['TEACHER'],
    },
    {
        key: 'teacher-leaves',
        labelVi: 'Duyệt đơn phép',
        path: '/app/teacher/leaves',
        icon: Description,
        roles: ['TEACHER'],
    },

    {
        key: 'student-timetable',
        labelVi: 'Thời khóa biểu',
        path: '/app/student/timetable',
        icon: AccessTime,
        roles: ['STUDENT', 'PARENT'],
    },
    {
        key: 'student-results',
        labelVi: 'Kết quả học tập',
        path: '/app/student/results',
        icon: School,
        roles: ['STUDENT', 'PARENT'],
    },
    {
        key: 'student-leaves',
        labelVi: 'Đơn xin nghỉ',
        path: '/app/student/leaves',
        icon: Description,
        roles: ['STUDENT', 'PARENT'],
    },

    {
        key: 'compliance-approvals',
        labelVi: 'Tài liệu cần duyệt',
        path: '/app/compliance/approvals',
        icon: AssignmentOutlined,
        roles: ['COMPLIANCE'],
    },
    {
        key: 'compliance-history',
        labelVi: 'Lịch sử kiểm duyệt',
        path: '/app/compliance/history',
        icon: Description,
        roles: ['COMPLIANCE'],
    },
];

export const getNavItemsByRole = (role?: string) => {
    if (!role) return [];

    const filterItems = (items: NavItem[]): NavItem[] =>
        items
            .map((item) => {
                const children = item.children ? filterItems(item.children) : undefined;

                if (children && children.length > 0) {
                    return { ...item, children };
                }

                if (hasAnyRole(role, item.roles)) {
                    return { ...item, children: undefined };
                }

                return null;
            })
            .filter((item): item is NavItem => item !== null);

    return filterItems(NAV_ITEMS);
};
