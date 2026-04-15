import LeaveRequestsPage from '@/pages/shared/LeaveRequestsPage';

export default function TeacherLeavesPage() {
    return (
        <LeaveRequestsPage
            mode="staff"
            title="Duyệt đơn xin phép"
            subtitle="Xem, duyệt hoặc từ chối các đơn xin phép liên quan đến lớp bạn phụ trách."
        />
    );
}
