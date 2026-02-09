
import {
    Users,
    FileText,
    BookOpen,
    GraduationCap,
    Calendar,
    AlertTriangle,
    BarChart3,
    Clock,
    ClipboardCheck,
    Upload,
    UserCog,
    MessageSquare,
    ShieldAlert,
    FileCheck
} from 'lucide-react';

export interface NavItem {
    label: string;
    path: string;
    icon: any;
    roles: string[];
    children?: NavItem[];
}

export const NAV_ITEMS: NavItem[] = [
    // Admin Module
    {
        label: "Accounts & Roles",
        path: "/app/admin/accounts",
        icon: Users,
        roles: ["admin", "super_admin"]
    },
    {
        label: "Legal & Profiles",
        path: "/app/admin/legal",
        icon: FileText,
        roles: ["admin", "super_admin"]
    },
    {
        label: "Teachers",
        path: "/app/admin/teachers",
        icon: Users,
        roles: ["admin", "super_admin"]
    },
    {
        label: "Programs & Courses",
        path: "/app/admin/programs",
        icon: BookOpen,
        roles: ["admin", "super_admin"]
    },
    {
        label: "Classes",
        path: "/app/admin/classes",
        icon: GraduationCap,
        roles: ["admin", "super_admin"]
    },
    {
        label: "Auto Scheduling",
        path: "/app/admin/scheduling",
        icon: Calendar,
        roles: ["admin", "super_admin"]
    },
    {
        label: "Conflict Resolution",
        path: "/app/admin/conflicts",
        icon: AlertTriangle,
        roles: ["admin", "super_admin"]
    },
    {
        label: "Reports & Analytics",
        path: "/app/admin/reports",
        icon: BarChart3,
        roles: ["admin", "super_admin"]
    },

    // Teacher Module
    {
        label: "My Schedule",
        path: "/app/teacher/schedule",
        icon: Calendar,
        roles: ["teacher"]
    },
    {
        label: "Attendance",
        path: "/app/teacher/attendance",
        icon: ClipboardCheck,
        roles: ["teacher"]
    },
    {
        label: "Lesson Journal",
        path: "/app/teacher/journal",
        icon: BookOpen,
        roles: ["teacher"]
    },
    {
        label: "Upload Documents",
        path: "/app/teacher/documents",
        icon: Upload,
        roles: ["teacher"]
    },
    {
        label: "Substitute Request",
        path: "/app/teacher/substitute",
        icon: UserCog,
        roles: ["teacher"]
    },

    // Student/Parent Module
    {
        label: "My Timetable",
        path: "/app/student/timetable",
        icon: Clock,
        roles: ["student", "parent"]
    },
    {
        label: "Learning Results",
        path: "/app/student/results",
        icon: GraduationCap,
        roles: ["student", "parent"]
    },
    {
        label: "Leave Requests",
        path: "/app/student/leaves",
        icon: FileText,
        roles: ["student", "parent"]
    },
    {
        label: "Course Consulting",
        path: "/app/student/consulting",
        icon: MessageSquare,
        roles: ["student", "parent"]
    },
    {
        label: "AI Assistant",
        path: "/app/student/ai-chat",
        icon: MessageSquare,
        roles: ["student", "parent"]
    },

    // Compliance Module
    {
        label: "Content Alerts",
        path: "/app/compliance/alerts",
        icon: ShieldAlert,
        roles: ["compliance"]
    },
    {
        label: "Approvals",
        path: "/app/compliance/approvals",
        icon: FileCheck,
        roles: ["compliance"]
    }
];

export const getNavItemsByRole = (role?: string) => {
    if (!role) return [];
    return NAV_ITEMS.filter(item => item.roles.includes(role));
};
