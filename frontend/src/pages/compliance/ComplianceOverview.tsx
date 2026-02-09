
import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { ShieldAlert, FileCheck, AlertTriangle, UserCheck } from 'lucide-react';

export const ComplianceOverview: React.FC = () => {
    const stats = [
        { label: "Pending Reviews", value: "12", icon: FileCheck, color: "text-blue-600" },
        { label: "Policy Violations", value: "2", icon: ShieldAlert, color: "text-red-600" },
        { label: "Flagged Content", value: "5", icon: AlertTriangle, color: "text-orange-600" },
        { label: "Audits Completed", value: "18", icon: UserCheck, color: "text-green-600" },
    ];

    return (
        <div className="space-y-6 animate-in fade-in-50 duration-500">
            <h1 className="text-3xl font-bold tracking-tight text-gray-900">Compliance Dashboard</h1>

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
                                    Need Action
                                </p>
                            </CardContent>
                        </Card>
                    );
                })}
            </div>
            <Card>
                <CardHeader>
                    <CardTitle>Priority Alerts</CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="text-sm text-gray-500">
                        All systems nominal. No critical alerts.
                    </div>
                </CardContent>
            </Card>
        </div>
    );
};
