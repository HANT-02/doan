
import React from 'react';

export const AppFooter: React.FC = () => {
    // Should verify if VITE_APP_VERSION is available, otherwise fallback
    const version = import.meta.env.VITE_APP_VERSION || '1.0.0';

    return (
        <footer className="bg-white border-t border-gray-100 py-4 px-6 md:px-8 mt-auto">
            <div className="flex flex-col md:flex-row justify-between items-center text-xs text-gray-500 gap-2">
                <p>&copy; {new Date().getFullYear()} EduCenter Management System. All rights reserved.</p>
                <p>Version {version}</p>
            </div>
        </footer>
    );
};
