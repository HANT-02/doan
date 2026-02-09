
import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

interface PlaceholderPageProps {
    title: string;
    description?: string;
}

export const PlaceholderPage: React.FC<PlaceholderPageProps> = ({ title, description }) => {
    return (
        <div className="space-y-6 animate-in fade-in-50 duration-500">
            <div className="flex flex-col gap-2">
                <h1 className="text-3xl font-bold tracking-tight text-gray-900">{title}</h1>
                {description && <p className="text-gray-500">{description}</p>}
            </div>

            <Card className="border-dashed border-2">
                <CardHeader>
                    <CardTitle className="text-lg font-medium text-gray-700">Content Placeholder</CardTitle>
                </CardHeader>
                <CardContent className="h-96 flex flex-col items-center justify-center text-gray-400 gap-4">
                    <div className="h-16 w-16 rounded-full bg-gray-100 flex items-center justify-center">
                        <span className="text-2xl">🚧</span>
                    </div>
                    <p> This module is currently under development.</p>
                </CardContent>
            </Card>
        </div>
    );
};
