import React from 'react';
import { cn } from '@/lib/utils';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';

interface AuthCardProps {
    title?: string;
    subtitle?: string;
    children: React.ReactNode;
    className?: string;
}

export const AuthCard: React.FC<AuthCardProps> = ({ title, subtitle, children, className }) => {
    return (
        <Card className={cn("w-full shadow-2xl border-muted/40 bg-white/95 dark:bg-slate-900/95 backdrop-blur-md transition-all duration-300", className)}>
            {(title || subtitle) && (
                <CardHeader className="space-y-1.5 pb-6 text-center">
                    {title && (
                        <CardTitle className="text-3xl font-extrabold tracking-tight bg-gradient-to-r from-slate-900 to-slate-700 dark:from-white dark:to-slate-300 bg-clip-text text-transparent">
                            {title}
                        </CardTitle>
                    )}
                    {subtitle && (
                        <CardDescription className="text-sm font-medium text-muted-foreground/80">
                            {subtitle}
                        </CardDescription>
                    )}
                </CardHeader>
            )}
            <CardContent>
                {children}
            </CardContent>
        </Card>
    );
};
