import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Users, BookOpen, GraduationCap, AlertTriangle } from 'lucide-react';
import {
    BarChart,
    Bar,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    Legend,
    ResponsiveContainer,
    PieChart,
    Pie,
    Cell
} from 'recharts';

export const AdminOverview: React.FC = () => {
    // Mock Statistics
    const stats = [
        { label: "Tổng số học viên", value: "1,234", icon: Users, color: "text-blue-600" },
        { label: "Khóa học đang mở", value: "42", icon: BookOpen, color: "text-green-600" },
        { label: "Lớp học hoạt động", value: "156", icon: GraduationCap, color: "text-purple-600" },
        { label: "Cảnh báo hệ thống", value: "3", icon: AlertTriangle, color: "text-red-600" },
    ];

    // Mock Data for Bar Chart
    const revenueData = [
        { name: 'T1', Lợi_nhuận: 4000, Doanh_thu: 2400 },
        { name: 'T2', Lợi_nhuận: 3000, Doanh_thu: 1398 },
        { name: 'T3', Lợi_nhuận: 2000, Doanh_thu: 9800 },
        { name: 'T4', Lợi_nhuận: 2780, Doanh_thu: 3908 },
        { name: 'T5', Lợi_nhuận: 1890, Doanh_thu: 4800 },
        { name: 'T6', Lợi_nhuận: 2390, Doanh_thu: 3800 },
        { name: 'T7', Lợi_nhuận: 3490, Doanh_thu: 4300 },
    ];

    // Mock Data for Pie Chart
    const courseDistribution = [
        { name: 'Toán', value: 400 },
        { name: 'Tiếng Anh', value: 300 },
        { name: 'Vật Lý', value: 300 },
        { name: 'Hóa Học', value: 200 },
    ];

    const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042'];

    return (
        <div className="space-y-6 animate-in fade-in-50 duration-500">
            <h1 className="text-3xl font-bold tracking-tight text-gray-900">Tổng quan hệ thống</h1>

            {/* Top Stat Cards */}
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                {stats.map((stat) => {
                    const Icon = stat.icon;
                    return (
                        <Card key={stat.label}>
                            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                                <CardTitle className="text-sm font-medium">
                                    {stat.label}
                                </CardTitle>
                                <Icon className={`h-4 w-4 ${stat.color}`} />
                            </CardHeader>
                            <CardContent>
                                <div className="text-2xl font-bold">{stat.value}</div>
                                <p className="text-xs text-muted-foreground">
                                    +2.1% so với tháng trước
                                </p>
                            </CardContent>
                        </Card>
                    );
                })}
            </div>

            {/* Charts Section */}
            <div className="grid gap-4 md:grid-cols-2">
                {/* Bar Chart */}
                <Card>
                    <CardHeader>
                        <CardTitle>Biểu đồ Doanh thu (Fake Data)</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="h-[300px] w-full">
                            <ResponsiveContainer width="100%" height="100%">
                                <BarChart
                                    data={revenueData}
                                    margin={{
                                        top: 5,
                                        right: 30,
                                        left: 20,
                                        bottom: 5,
                                    }}
                                >
                                    <CartesianGrid strokeDasharray="3 3" />
                                    <XAxis dataKey="name" />
                                    <YAxis />
                                    <Tooltip />
                                    <Legend />
                                    <Bar dataKey="Doanh_thu" fill="#8884d8" name="Doanh Thu" />
                                    <Bar dataKey="Lợi_nhuận" fill="#82ca9d" name="Lợi Nhuận" />
                                </BarChart>
                            </ResponsiveContainer>
                        </div>
                    </CardContent>
                </Card>

                {/* Pie Chart */}
                <Card>
                    <CardHeader>
                        <CardTitle>Phân bổ Học viên theo Môn học (Fake Data)</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="h-[300px] w-full flex justify-center items-center">
                            <ResponsiveContainer width="100%" height="100%">
                                <PieChart>
                                    <Pie
                                        data={courseDistribution}
                                        cx="50%"
                                        cy="50%"
                                        labelLine={false}
                                        label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                                        outerRadius={100}
                                        fill="#8884d8"
                                        dataKey="value"
                                    >
                                        {courseDistribution.map((entry, index) => (
                                            <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                                        ))}
                                    </Pie>
                                    <Tooltip />
                                </PieChart>
                            </ResponsiveContainer>
                        </div>
                    </CardContent>
                </Card>
            </div>
        </div>
    );
};
