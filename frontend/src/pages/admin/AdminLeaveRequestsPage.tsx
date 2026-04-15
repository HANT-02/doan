import LeaveRequestsPage from '@/pages/shared/LeaveRequestsPage';

export default function AdminLeaveRequestsPage() {
    return (
        <LeaveRequestsPage
            mode="staff"
            title="Quản lý đơn xin phép"
            subtitle="Admin theo dõi toàn bộ đơn xin phép và hỗ trợ duyệt / từ chối khi cần."
        />
    );
}
