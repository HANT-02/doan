
import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { FileQuestion } from 'lucide-react';

export const NotFoundPage: React.FC = () => {
    const navigate = useNavigate();

    return (
        <div className="flex flex-col items-center justify-center min-h-[80vh] text-center gap-6 px-4">
            <div className="bg-blue-50 p-6 rounded-full">
                <FileQuestion className="h-16 w-16 text-blue-600" />
            </div>
            <div className="space-y-2">
                <h1 className="text-4xl font-extrabold tracking-tight text-gray-900 sm:text-5xl">404</h1>
                <h2 className="text-2xl font-semibold text-gray-800">Page Not Found</h2>
                <p className="max-w-md mx-auto text-gray-600 md:text-lg">
                    Sorry, we couldn't find the page you're looking for.
                    It might have been removed or renamed.
                </p>
            </div>
            <div className="flex gap-4">
                <Button variant="outline" onClick={() => navigate(-1)}>
                    Go Back
                </Button>
                <Button onClick={() => navigate('/app')}>
                    Go to Dashboard
                </Button>
            </div>
        </div>
    );
};
