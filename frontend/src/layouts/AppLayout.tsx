
import React, { useState } from 'react';
import { Outlet } from 'react-router-dom';
import { AppSidebar } from '@/components/layout/AppSidebar';
import { AppHeader } from '@/components/layout/AppHeader';
import { AppFooter } from '@/components/layout/AppFooter';

export const AppLayout: React.FC = () => {
    const [sidebarOpen, setSidebarOpen] = useState(false);

    return (
        <div className="flex min-h-screen bg-gray-50/50">
            {/* Sidebar */}
            <AppSidebar
                isMobileOpen={sidebarOpen}
                onMobileClose={() => setSidebarOpen(false)}
            />

            {/* Main Content Area */}
            <div className="flex flex-col flex-1 min-w-0 transition-all duration-300 ease-in-out md:ml-0">
                <AppHeader onMobileOpen={() => setSidebarOpen(true)} />

                <main className="flex-1 p-4 md:p-6 lg:p-8 overflow-y-auto">
                    <Outlet />
                </main>

                <AppFooter />
            </div>
        </div>
    );
};
