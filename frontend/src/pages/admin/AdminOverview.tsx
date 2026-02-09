
import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Users, BookOpen, GraduationCap, AlertTriangle } from 'lucide-react';

export const AdminOverview: React.FC = () => {
    const stats = [
        { label: "Total Users", value: "1,234", icon: Users, color: "text-blue-600" },
        { label: "Active Courses", value: "42", icon: BookOpen, color: "text-green-600" },
        { label: "Classes", value: "156", icon: GraduationCap, color: "text-purple-600" },
        { label: "System Alerts", value: "3", icon: AlertTriangle, color: "text-red-600" },
    ];

    return (
        <div className="space-y-6 animate-in fade-in-50 duration-500">
            <h1 className="text-3xl font-bold tracking-tight text-gray-900">Admin Dashboard</h1>

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
                                    +2.1% from last month
                                </p>
                            </CardContent>
                        </Card>
                    );
                })}
            </div>

            <Card>
                <CardHeader>
                    <CardTitle>Recent Activity</CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="text-sm text-gray-500">
                        No recent activity to display.
                    </div>
                </CardContent>
            </Card>
        </div>
    );
};
