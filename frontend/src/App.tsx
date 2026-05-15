import { BrowserRouter as Router, Routes, Route, Navigate, Outlet, useLocation } from 'react-router-dom';
import { AuthProvider } from '@/contexts/AuthContext';
import { PrivateRoute } from '@/routes/PrivateRoute';
import { RoleRoute } from '@/routes/RoleRoute';
import { AuthLayout } from '@/components/AuthLayout';
import { AppLayout } from '@/layouts/AppLayout';
import { Provider } from 'react-redux';
import { store } from '@/store';

// Public Pages
import { LoginPage } from '@/pages/LoginPage';
import { RegisterPage } from '@/pages/RegisterPage';
import { ForgotPasswordPage } from '@/pages/ForgotPasswordPage';
import { ResetPasswordPage } from '@/pages/ResetPasswordPage';
import { VerifyOTPPage } from '@/pages/VerifyOTPPage';

// Protected Pages
import { ProfilePage } from '@/pages/ProfilePage';
import { ChangePasswordPage } from '@/pages/ChangePasswordPage';

// Role Overview Pages
import { AdminOverview } from '@/pages/admin/AdminOverview';
import AdminLeaveRequestsPage from '@/pages/admin/AdminLeaveRequestsPage';
import { TeacherOverview } from '@/pages/teacher/TeacherOverview';
import TeacherLessonDetailPage from '@/pages/teacher/TeacherLessonDetailPage';
import TeacherSchedulePage from '@/pages/teacher/TeacherSchedulePage';
import TeacherAttendancePage from '@/pages/teacher/TeacherAttendancePage';
import TeacherJournalPage from '@/pages/teacher/TeacherJournalPage';
import TeacherLeavesPage from '@/pages/teacher/TeacherLeavesPage';
import { TeacherSubstitutePage } from '@/pages/teacher/TeacherSubstitutePage';
import { StudentOverview } from '@/pages/student/StudentOverview';
import StudentTimetablePage from '@/pages/student/StudentTimetablePage';
import StudentResultsPage from '@/pages/student/StudentResultsPage';
import StudentLeavesPage from '@/pages/student/StudentLeavesPage';
import { ComplianceOverview } from '@/pages/compliance/ComplianceOverview';
import { TeacherDocumentsPage } from '@/pages/teacher/TeacherDocumentsPage';
import { ComplianceQueuePage } from '@/pages/compliance/ComplianceQueuePage';

// Admin Teacher Management Pages
import { TeachersPage } from '@/pages/admin/TeachersPage';
import { TeacherDetailPage } from '@/pages/admin/TeacherDetailPage';
import { TeacherFormPage } from '@/pages/admin/TeacherFormPage';
import { RoomsPage } from '@/pages/admin/RoomsPage';
import { ShiftsPage } from '@/pages/admin/ShiftsPage';
import { ClassesPage } from '@/pages/admin/ClassesPage';
import { StudentsPage } from '@/pages/admin/StudentsPage';
import { ProgramPage } from '@/pages/admin/ProgramPage';
import { CoursePage } from '@/pages/admin/CoursePage';
import { SchedulingPage } from '@/pages/admin/SchedulingPage';
import { PredictiveAlertsPage } from '@/pages/admin/PredictiveAlertsPage';
import { AccountsPage } from '@/pages/admin/AccountsPage';
import LessonsPage from '@/pages/admin/LessonsPage';
import LessonDetailPage from '@/pages/admin/LessonDetailPage';

// Placeholder & Error Pages
import { PlaceholderPage } from '@/components/layout/PlaceholderPage';
import ErrorBoundary from '@/components/common/ErrorBoundary';
import { ForbiddenPage } from '@/pages/ForbiddenPage';
import { NotFoundPage } from '@/pages/NotFoundPage';

import { Toaster } from 'sonner';
import { useAuth } from '@/contexts/AuthContext';
import { normalizeRole } from '@/utils/roles';

const DashboardRedirect = () => {
  const { user } = useAuth();

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  // Redirect based on role (standardized to uppercase)
  switch (normalizeRole(user.role)) {
    case 'ADMIN':
    case 'SUPER_ADMIN':
      return <Navigate to="/app/admin/overview" replace />;
    case 'TEACHER':
      return <Navigate to="/app/teacher/overview" replace />;
    case 'STUDENT':
    case 'PARENT':
      return <Navigate to="/app/student/overview" replace />;
    case 'COMPLIANCE':
      return <Navigate to="/app/compliance/overview" replace />;
    default:
      console.warn('[DEV] Unknown Role for Redirect:', user.role);
      return <Navigate to="/403" replace />;
  }
};

const ErrorBoundaryShell = () => (
  <ErrorBoundary>
    <Outlet />
  </ErrorBoundary>
);

const AppRoutes = () => {
  const location = useLocation();

  return (
    <Routes location={location} key={location.pathname}>
      {/* Public Routes wrapped in AuthLayout */}
      <Route element={<AuthLayout />}>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route path="/verify-otp" element={<VerifyOTPPage />} />
        <Route path="/forgot-password" element={<ForgotPasswordPage />} />
        <Route path="/reset-password" element={<ResetPasswordPage />} />
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
          <Route path="admin" element={<RoleRoute allowedRoles={['ADMIN', 'SUPER_ADMIN']} />}>
            <Route element={<ErrorBoundaryShell />}>
              <Route index element={<Navigate to="overview" replace />} />
              <Route path="overview" element={<AdminOverview />} />
              <Route path="accounts" element={<AccountsPage />} />
              <Route path="teachers" element={<TeachersPage />} />
              <Route path="teachers/new" element={<TeacherFormPage />} />
              <Route path="teachers/:id" element={<TeacherDetailPage />} />
              <Route path="teachers/:id/edit" element={<TeacherFormPage />} />
              <Route path="programs" element={<ProgramPage />} />
              <Route path="courses" element={<CoursePage />} />
              <Route path="classes" element={<ClassesPage />} />
              <Route path="students" element={<StudentsPage />} />
              <Route path="rooms" element={<RoomsPage />} />
              <Route path="shifts" element={<ShiftsPage />} />
              <Route path="scheduling" element={<SchedulingPage />} />
              <Route path="predictive" element={<PredictiveAlertsPage />} />
              <Route path="lessons" element={<LessonsPage />} />
              <Route path="lessons/:id" element={<LessonDetailPage />} />
              <Route path="leaves" element={<AdminLeaveRequestsPage />} />
              <Route path="conflicts" element={<PlaceholderPage title="Conflict Resolution" />} />
              <Route path="reports" element={<PlaceholderPage title="Reports & Analytics" />} />
            </Route>
          </Route>

          {/* Teacher Routes */}
            <Route path="teacher" element={<RoleRoute allowedRoles={['TEACHER']} />}>
              <Route index element={<Navigate to="overview" replace />} />
              <Route path="overview" element={<TeacherOverview />} />
              <Route path="schedule" element={<TeacherSchedulePage />} />
              <Route path="lessons/:lessonId" element={<TeacherLessonDetailPage />} />
            <Route path="attendance" element={<TeacherAttendancePage />} />
            <Route path="journal" element={<TeacherJournalPage />} />
            <Route path="leaves" element={<TeacherLeavesPage />} />
            <Route path="documents" element={<TeacherDocumentsPage />} />
            <Route path="substitute" element={<TeacherSubstitutePage />} />
          </Route>

          {/* Student/Parent Routes */}
          <Route path="student" element={<RoleRoute allowedRoles={['STUDENT', 'PARENT']} />}>
            <Route index element={<Navigate to="overview" replace />} />
            <Route path="overview" element={<StudentOverview />} />
            <Route path="timetable" element={<StudentTimetablePage />} />
            <Route path="results" element={<StudentResultsPage />} />
            <Route path="leaves" element={<StudentLeavesPage />} />
            <Route path="consulting" element={<PlaceholderPage title="Course Consulting" />} />
            <Route path="ai-chat" element={<PlaceholderPage title="AI Assistant" />} />
          </Route>

          {/* Compliance Routes */}
          <Route path="compliance" element={<RoleRoute allowedRoles={['COMPLIANCE']} />}>
            <Route index element={<Navigate to="overview" replace />} />
            <Route path="overview" element={<ComplianceOverview />} />
            <Route path="alerts" element={<PlaceholderPage title="Content Alerts" />} />
            <Route path="approvals" element={<ComplianceQueuePage />} />
          </Route>
        </Route>
      </Route>

      {/* Root Redirect - Check role or go login */}
      <Route path="/" element={<Navigate to="/app" replace />} />

      {/* Error Pages */}
      <Route path="/403" element={<ForbiddenPage />} />
      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  );
};

function App() {
  return (
    <Provider store={store}>
      <Router>
        <AuthProvider>
          <AppRoutes />
          <Toaster />
        </AuthProvider>
      </Router>
    </Provider>
  );
}

export default App;
