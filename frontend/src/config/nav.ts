
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
    Person
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
        roles: ["ADMIN", "SUPER_ADMIN"]
    },
    {
        label: "Legal & Profiles",
        path: "/app/admin/legal",
        icon: Description,
        roles: ["ADMIN", "SUPER_ADMIN"]
    },
    {
        label: "Teachers",
        path: "/app/admin/teachers",
        icon: Group, // Using Group for teachers as well
        roles: ["ADMIN", "SUPER_ADMIN"]
    },
    {
        label: "Programs & Courses",
        path: "/app/admin/programs",
        icon: Book,
        roles: ["ADMIN", "SUPER_ADMIN"]
    },
    {
        label: "Classes",
        path: "/app/admin/classes",
        icon: School,
        roles: ["ADMIN", "SUPER_ADMIN"]
    },
    {
        label: "Students",
        path: "/app/admin/students",
        icon: Person,
        roles: ["ADMIN", "SUPER_ADMIN"]
    },
    {
        label: "Rooms",
        path: "/app/admin/rooms",
        icon: MeetingRoom,
        roles: ["ADMIN", "SUPER_ADMIN"]
    },
    {
        label: "Auto Scheduling",
        path: "/app/admin/scheduling",
        icon: CalendarMonth,
        roles: ["ADMIN", "SUPER_ADMIN"]
    },
    {
        label: "Conflict Resolution",
        path: "/app/admin/conflicts",
        icon: Warning,
        roles: ["ADMIN", "SUPER_ADMIN"]
    },
    {
        label: "Reports & Analytics",
        path: "/app/admin/reports",
        icon: BarChart,
        roles: ["ADMIN", "SUPER_ADMIN"]
    },

    // Teacher Module
    {
        label: "My Schedule",
        path: "/app/teacher/schedule",
        icon: CalendarMonth,
        roles: ["TEACHER"]
    },
    {
        label: "Attendance",
        path: "/app/teacher/attendance",
        icon: AssignmentTurnedIn,
        roles: ["TEACHER"]
    },
    {
        label: "Lesson Journal",
        path: "/app/teacher/journal",
        icon: Book,
        roles: ["TEACHER"]
    },
    {
        label: "Upload Documents",
        path: "/app/teacher/documents",
        icon: UploadFile,
        roles: ["TEACHER"]
    },
    {
        label: "Substitute Request",
        path: "/app/teacher/substitute",
        icon: ManageAccounts,
        roles: ["TEACHER"]
    },

    // Student/Parent Module
    {
        label: "My Timetable",
        path: "/app/student/timetable",
        icon: AccessTime,
        roles: ["STUDENT", "PARENT"]
    },
    {
        label: "Learning Results",
        path: "/app/student/results",
        icon: School,
        roles: ["STUDENT", "PARENT"]
    },
    {
        label: "Leave Requests",
        path: "/app/student/leaves",
        icon: Description,
        roles: ["STUDENT", "PARENT"]
    },
    {
        label: "Course Consulting",
        path: "/app/student/consulting",
        icon: ChatBubble,
        roles: ["STUDENT", "PARENT"]
    },
    {
        label: "AI Assistant",
        path: "/app/student/ai-chat",
        icon: ChatBubble, // Using ChatBubble for AI Assistant as well
        roles: ["STUDENT", "PARENT"]
    },

    // Compliance Module
    {
        label: "Content Alerts",
        path: "/app/compliance/alerts",
        icon: Security,
        roles: ["COMPLIANCE"]
    },
    {
        label: "Approvals",
        path: "/app/compliance/approvals",
        icon: AssignmentOutlined,
        roles: ["COMPLIANCE"]
    }
];

export const getNavItemsByRole = (role?: string) => {
    if (!role) return [];
    return NAV_ITEMS.filter(item => item.roles.includes(role));
};
