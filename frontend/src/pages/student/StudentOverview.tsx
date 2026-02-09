
import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Calendar, GraduationCap, Star, BookOpen } from 'lucide-react';

export const StudentOverview: React.FC = () => {
    const stats = [
        { label: "GPA", value: "3.8", icon: Star, color: "text-yellow-500" },
        { label: "Credits Earned", value: "45", icon: GraduationCap, color: "text-blue-600" },
        { label: "Active Courses", value: "6", icon: BookOpen, color: "text-green-600" },
        { label: "Attendance", value: "98%", icon: Calendar, color: "text-purple-600" },
    ];

    return (
        <div className="space-y-6 animate-in fade-in-50 duration-500">
            <h1 className="text-3xl font-bold tracking-tight text-gray-900">Student Portal</h1>

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
                                    Current Semester
                                </p>
                            </CardContent>
                        </Card>
                    );
                })}
            </div>
            <Card>
                <CardHeader>
                    <CardTitle>My Timetable</CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="text-sm text-gray-500">
                        No classes scheduled for today.
                    </div>
                </CardContent>
            </Card>
        </div>
    );
};
