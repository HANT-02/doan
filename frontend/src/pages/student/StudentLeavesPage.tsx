import LeaveRequestsPage from '@/pages/shared/LeaveRequestsPage';

export default function StudentLeavesPage() {
    return (
        <LeaveRequestsPage
            mode="student"
            title="Đơn xin phép"
            subtitle="Tạo, theo dõi và hủy đơn xin nghỉ hoặc xin đi muộn của bạn."
        />
    );
}
