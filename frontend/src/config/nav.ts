import {
    Group,
    Description,
    Book,
    School,
    CalendarMonth,
    Warning,
    BarChart,
    AccessTime,
    AssignmentTurnedIn,
    UploadFile,
    ManageAccounts,
    ChatBubble,
    Security,
    AssignmentOutlined,
    MeetingRoom,
    Person,
    Dashboard,
    Build
} from '@mui/icons-material';

export interface NavItem {
    key: string;
    path: string;
    labelVi: string;
    labelEn?: string;
    icon: any;
    roles: string[];
}

export const NAV_ITEMS: NavItem[] = [
    // --- Admin Module ---
    {
        key: "admin-overview",
        labelVi: "Tổng quan",
        labelEn: "Overview",
        path: "/app/admin/overview",
        icon: Dashboard,
        roles: ["ADMIN", "SUPER_ADMIN"]
    },
    {
        key: "admin-teachers",
        labelVi: "Quản lý giáo viên",
        path: "/app/admin/teachers",
        icon: Group,
        roles: ["ADMIN", "SUPER_ADMIN"]
    },
    {
        key: "admin-students",
        labelVi: "Quản lý học sinh",
        path: "/app/admin/students",
        icon: Person,
        roles: ["ADMIN", "SUPER_ADMIN"]
    },
    {
        key: "admin-classes",
        labelVi: "Quản lý lớp học",
        path: "/app/admin/classes",
        icon: School,
        roles: ["ADMIN", "SUPER_ADMIN"]
    },
    {
        key: "admin-rooms",
        labelVi: "Quản lý phòng học",
        path: "/app/admin/rooms",
        icon: MeetingRoom,
        roles: ["ADMIN", "SUPER_ADMIN"]
    },
    {
        key: "admin-programs",
        labelVi: "Chương trình / Khóa học",
        path: "/app/admin/programs",
        icon: Book,
        roles: ["ADMIN", "SUPER_ADMIN"]
    },
    {
        key: "admin-scheduling",
        labelVi: "Xếp lịch (CSP)",
        path: "/app/admin/scheduling",
        icon: CalendarMonth,
        roles: ["ADMIN", "SUPER_ADMIN"]
    },
    {
        key: "admin-ai-audit",
        labelVi: "Kiểm duyệt tài liệu (AI)",
        path: "/app/admin/audit",
        icon: Security,
        roles: ["ADMIN", "SUPER_ADMIN"]
    },
    {
        key: "admin-accounts",
        labelVi: "Quản lý tài khoản",
        path: "/app/admin/accounts",
        icon: ManageAccounts,
        roles: ["ADMIN", "SUPER_ADMIN"]
    },
    {
        key: "admin-devtools",
        labelVi: "Công cụ kiểm thử (DevTools)",
        path: "/app/admin/devtools",
        icon: Build,
        roles: ["ADMIN", "SUPER_ADMIN"]
    },

    // --- Teacher Module ---
    {
        key: "teacher-schedule",
        labelVi: "Lịch giảng dạy",
        path: "/app/teacher/schedule",
        icon: CalendarMonth,
        roles: ["TEACHER"]
    },
    {
        key: "teacher-attendance",
        labelVi: "Điểm danh",
        path: "/app/teacher/attendance",
        icon: AssignmentTurnedIn,
        roles: ["TEACHER"]
    },
    {
        key: "teacher-journal",
        labelVi: "Sổ đầu bài",
        path: "/app/teacher/journal",
        icon: Book,
        roles: ["TEACHER"]
    },
    {
        key: "teacher-documents",
        labelVi: "Tài liệu giảng dạy",
        path: "/app/teacher/documents",
        icon: UploadFile,
        roles: ["TEACHER"]
    },

    // --- Student/Parent Module ---
    {
        key: "student-timetable",
        labelVi: "Thời khóa biểu",
        path: "/app/student/timetable",
        icon: AccessTime,
        roles: ["STUDENT", "PARENT"]
    },
    {
        key: "student-results",
        labelVi: "Kết quả học tập",
        path: "/app/student/results",
        icon: School,
        roles: ["STUDENT", "PARENT"]
    },
    {
        key: "student-leaves",
        labelVi: "Đơn xin nghỉ",
        path: "/app/student/leaves",
        icon: Description,
        roles: ["STUDENT", "PARENT"]
    },

    // --- Compliance Module ---
    {
        key: "compliance-approvals",
        labelVi: "Tài liệu cần duyệt",
        path: "/app/compliance/approvals",
        icon: AssignmentOutlined,
        roles: ["COMPLIANCE"]
    },
    {
        key: "compliance-history",
        labelVi: "Lịch sử kiểm duyệt",
        path: "/app/compliance/history",
        icon: Description,
        roles: ["COMPLIANCE"]
    }
];

export const getNavItemsByRole = (role?: string) => {
    if (!role) return [];
    return NAV_ITEMS.filter(item => item.roles.includes(role));
};
