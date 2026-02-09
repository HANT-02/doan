
import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Calendar, Users, FileText, Clock } from 'lucide-react';

export const TeacherOverview: React.FC = () => {
    const stats = [
        { label: "Today's Classes", value: "4", icon: Calendar, color: "text-blue-600" },
        { label: "Total Students", value: "128", icon: Users, color: "text-green-600" },
        { label: "Pending Reports", value: "2", icon: FileText, color: "text-orange-600" },
        { label: "Teaching Hours", value: "24h", icon: Clock, color: "text-purple-600" },
    ];

    return (
        <div className="space-y-6 animate-in fade-in-50 duration-500">
            <h1 className="text-3xl font-bold tracking-tight text-gray-900">Teacher Dashboard</h1>

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
                                    Updated just now
                                </p>
                            </CardContent>
                        </Card>
                    );
                })}
            </div>
            <Card>
                <CardHeader>
                    <CardTitle>Upcoming Schedule</CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="text-sm text-gray-500">
                        You have no upcoming classes today.
                    </div>
                </CardContent>
            </Card>
        </div>
    );
};
