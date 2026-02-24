import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from '@/contexts/AuthContext';
import { PrivateRoute } from '@/routes/PrivateRoute';
import { RoleRoute } from '@/routes/RoleRoute';
import { AuthLayout } from '@/components/AuthLayout';
import { AppLayout } from '@/layouts/AppLayout';

// Public Pages
import { LoginPage } from '@/pages/LoginPage';
// import { RegisterPage } from '@/pages/RegisterPage';
// import { ForgotPasswordPage } from '@/pages/ForgotPasswordPage';
// import { ResetPasswordPage } from '@/pages/ResetPasswordPage';

// Protected Pages
import { ProfilePage } from '@/pages/ProfilePage';
import { ChangePasswordPage } from '@/pages/ChangePasswordPage';

// Role Overview Pages
import { AdminOverview } from '@/pages/admin/AdminOverview';
import { TeacherOverview } from '@/pages/teacher/TeacherOverview';
import { StudentOverview } from '@/pages/student/StudentOverview';
import { ComplianceOverview } from '@/pages/compliance/ComplianceOverview';

// Admin Teacher Management Pages
import { TeachersPage } from '@/pages/admin/TeachersPage';
import { TeacherDetailPage } from '@/pages/admin/TeacherDetailPage';
import { TeacherFormPage } from '@/pages/admin/TeacherFormPage';

// Placeholder & Error Pages
import { PlaceholderPage } from '@/components/layout/PlaceholderPage';
import { ForbiddenPage } from '@/pages/ForbiddenPage';
import { NotFoundPage } from '@/pages/NotFoundPage';

import { Toaster } from 'sonner';
import { useAuth } from '@/contexts/AuthContext';

const DashboardRedirect = () => {
  const { user } = useAuth();

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  // Redirect based on role
  switch (user.role) {
    case 'admin':
    case 'super_admin':
      return <Navigate to="/app/admin/overview" replace />;
    case 'teacher':
      return <Navigate to="/app/teacher/overview" replace />;
    case 'student':
    case 'parent':
      return <Navigate to="/app/student/overview" replace />;
    case 'compliance':
      return <Navigate to="/app/compliance/overview" replace />;
    default:
      return <Navigate to="/403" replace />;
  }
};

function App() {
  return (
    <Router>
      <AuthProvider>
        <Routes>
          {/* Public Routes wrapped in AuthLayout */}
          <Route element={<AuthLayout />}>
            <Route path="/login" element={<LoginPage />} />
            {/* Disabled for Demo */}
            {/* <Route path="/register" element={<RegisterPage />} />
            <Route path="/forgot-password" element={<ForgotPasswordPage />} />
            <Route path="/reset-password" element={<ResetPasswordPage />} /> */}
          </Route>

          {/* Protected Dashboard Routes */}
          <Route path="/app" element={<PrivateRoute />}>
            <Route element={<AppLayout />}>
              {/* Dashboard Root Redirect */}
              <Route index element={<DashboardRedirect />} />

              {/* Common Protected Routes */}
              <Route path="profile" element={<ProfilePage />} />
              <Route path="change-password" element={<ChangePasswordPage />} />

              {/* Admin Routes */}
              <Route element={<RoleRoute allowedRoles={['admin', 'super_admin']} />}>
                <Route path="admin/overview" element={<AdminOverview />} />
                <Route path="admin/accounts" element={<PlaceholderPage title="Accounts & Roles" />} />
                <Route path="admin/legal" element={<PlaceholderPage title="Legal & Profiles" />} />
                <Route path="admin/teachers" element={<TeachersPage />} />
                <Route path="admin/teachers/new" element={<TeacherFormPage />} />
                <Route path="admin/teachers/:id" element={<TeacherDetailPage />} />
                <Route path="admin/teachers/:id/edit" element={<TeacherFormPage />} />
                <Route path="admin/programs" element={<PlaceholderPage title="Programs & Courses" />} />
                <Route path="admin/classes" element={<PlaceholderPage title="Classes" />} />
                <Route path="admin/scheduling" element={<PlaceholderPage title="Auto Scheduling" />} />
                <Route path="admin/conflicts" element={<PlaceholderPage title="Conflict Resolution" />} />
                <Route path="admin/reports" element={<PlaceholderPage title="Reports & Analytics" />} />
              </Route>

              {/* Teacher Routes */}
              <Route element={<RoleRoute allowedRoles={['teacher']} />}>
                <Route path="teacher/overview" element={<TeacherOverview />} />
                <Route path="teacher/schedule" element={<PlaceholderPage title="My Schedule" />} />
                <Route path="teacher/attendance" element={<PlaceholderPage title="Attendance" />} />
                <Route path="teacher/journal" element={<PlaceholderPage title="Lesson Journal" />} />
                <Route path="teacher/documents" element={<PlaceholderPage title="Upload Documents" />} />
                <Route path="teacher/substitute" element={<PlaceholderPage title="Substitute Request" />} />
              </Route>

              {/* Student/Parent Routes */}
              <Route element={<RoleRoute allowedRoles={['student', 'parent']} />}>
                <Route path="student/overview" element={<StudentOverview />} />
                <Route path="student/timetable" element={<PlaceholderPage title="My Timetable" />} />
                <Route path="student/results" element={<PlaceholderPage title="Learning Results" />} />
                <Route path="student/leaves" element={<PlaceholderPage title="Leave Requests" />} />
                <Route path="student/consulting" element={<PlaceholderPage title="Course Consulting" />} />
                <Route path="student/ai-chat" element={<PlaceholderPage title="AI Assistant" />} />
              </Route>

              {/* Compliance Routes */}
              <Route element={<RoleRoute allowedRoles={['compliance']} />}>
                <Route path="compliance/overview" element={<ComplianceOverview />} />
                <Route path="compliance/alerts" element={<PlaceholderPage title="Content Alerts" />} />
                <Route path="compliance/approvals" element={<PlaceholderPage title="Approvals" />} />
              </Route>
            </Route>
          </Route>

          {/* Root Redirect - Check role or go login */}
          <Route path="/" element={<Navigate to="/app" replace />} />

          {/* Error Pages */}
          <Route path="/403" element={<ForbiddenPage />} />
          <Route path="*" element={<NotFoundPage />} />
        </Routes>
        <Toaster />
      </AuthProvider>
    </Router>
  );
}

export default App;
