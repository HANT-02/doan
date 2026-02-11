
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
    AssignmentOutlined
} from '@mui/icons-material';

export interface NavItem {
    label: string;
    path: string;
    icon: any; // Type remains 'any' or can be refined to React.ElementType<SvgIconProps>
    roles: string[];
    children?: NavItem[];
}

export const NAV_ITEMS: NavItem[] = [
    // Admin Module
    {
        label: "Accounts & Roles",
        path: "/app/admin/accounts",
        icon: Group,
        roles: ["admin", "super_admin"]
    },
    {
        label: "Legal & Profiles",
        path: "/app/admin/legal",
        icon: Description,
        roles: ["admin", "super_admin"]
    },
    {
        label: "Teachers",
        path: "/app/admin/teachers",
        icon: Group, // Using Group for teachers as well
        roles: ["admin", "super_admin"]
    },
    {
        label: "Programs & Courses",
        path: "/app/admin/programs",
        icon: Book,
        roles: ["admin", "super_admin"]
    },
    {
        label: "Classes",
        path: "/app/admin/classes",
        icon: School,
        roles: ["admin", "super_admin"]
    },
    {
        label: "Auto Scheduling",
        path: "/app/admin/scheduling",
        icon: CalendarMonth,
        roles: ["admin", "super_admin"]
    },
    {
        label: "Conflict Resolution",
        path: "/app/admin/conflicts",
        icon: Warning,
        roles: ["admin", "super_admin"]
    },
    {
        label: "Reports & Analytics",
        path: "/app/admin/reports",
        icon: BarChart,
        roles: ["admin", "super_admin"]
    },

    // Teacher Module
    {
        label: "My Schedule",
        path: "/app/teacher/schedule",
        icon: CalendarMonth,
        roles: ["teacher"]
    },
    {
        label: "Attendance",
        path: "/app/teacher/attendance",
        icon: AssignmentTurnedIn,
        roles: ["teacher"]
    },
    {
        label: "Lesson Journal",
        path: "/app/teacher/journal",
        icon: Book,
        roles: ["teacher"]
    },
    {
        label: "Upload Documents",
        path: "/app/teacher/documents",
        icon: UploadFile,
        roles: ["teacher"]
    },
    {
        label: "Substitute Request",
        path: "/app/teacher/substitute",
        icon: ManageAccounts,
        roles: ["teacher"]
    },

    // Student/Parent Module
    {
        label: "My Timetable",
        path: "/app/student/timetable",
        icon: AccessTime,
        roles: ["student", "parent"]
    },
    {
        label: "Learning Results",
        path: "/app/student/results",
        icon: School,
        roles: ["student", "parent"]
    },
    {
        label: "Leave Requests",
        path: "/app/student/leaves",
        icon: Description,
        roles: ["student", "parent"]
    },
    {
        label: "Course Consulting",
        path: "/app/student/consulting",
        icon: ChatBubble,
        roles: ["student", "parent"]
    },
    {
        label: "AI Assistant",
        path: "/app/student/ai-chat",
        icon: ChatBubble, // Using ChatBubble for AI Assistant as well
        roles: ["student", "parent"]
    },

    // Compliance Module
    {
        label: "Content Alerts",
        path: "/app/compliance/alerts",
        icon: Security,
        roles: ["compliance"]
    },
    {
        label: "Approvals",
        path: "/app/compliance/approvals",
        icon: AssignmentOutlined,
        roles: ["compliance"]
    }
];

export const getNavItemsByRole = (role?: string) => {
    if (!role) return [];
    return NAV_ITEMS.filter(item => item.roles.includes(role));
};
