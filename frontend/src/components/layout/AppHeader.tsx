
import React, { useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { Menu, Bell, LogOut, KeyRound } from 'lucide-react';
import { Link, useNavigate } from 'react-router-dom';

interface AppHeaderProps {
    onMobileOpen: () => void;
}

export const AppHeader: React.FC<AppHeaderProps> = ({ onMobileOpen }) => {
    const { user, logout } = useAuth();
    const navigate = useNavigate();
    const [isProfileOpen, setIsProfileOpen] = useState(false);

    const handleLogout = async () => {
        await logout();
        navigate('/login');
    };

    return (
        <header className="sticky top-0 z-40 flex h-16 shrink-0 items-center gap-x-4 border-b bg-white px-4 shadow-sm sm:gap-x-6 sm:px-6 lg:px-8">
            <button
                type="button"
                className="-m-2.5 p-2.5 text-gray-700 lg:hidden"
                onClick={onMobileOpen}
            >
                <span className="sr-only">Open sidebar</span>
                <Menu className="h-6 w-6" aria-hidden="true" />
            </button>

            <div className="flex flex-1 gap-x-4 self-stretch lg:gap-x-6">
                <form className="relative flex flex-1" action="#" method="GET">
                    {/* Optional Search Bar Placeholder */}
                    <div className="flex w-full items-center">
                        <input
                            id="search-field"
                            className="bg-transparent border-none w-full text-sm focus:ring-0 text-gray-900 placeholder:text-gray-400"
                            placeholder="Type to search..."
                            name="search"
                            type="search"
                        />
                    </div>
                </form>
                <div className="flex items-center gap-x-4 lg:gap-x-6">
                    <button type="button" className="-m-2.5 p-2.5 text-gray-400 hover:text-gray-500">
                        <span className="sr-only">View notifications</span>
                        <Bell className="h-6 w-6" aria-hidden="true" />
                    </button>

                    {/* Separator */}
                    <div className="hidden lg:block lg:h-6 lg:w-px lg:bg-gray-200" aria-hidden="true" />

                    {/* Profile Dropdown */}
                    <div className="relative">
                        <button
                            className="-m-1.5 flex items-center p-1.5 focus:outline-none"
                            onClick={() => setIsProfileOpen(!isProfileOpen)}
                            onBlur={() => setTimeout(() => setIsProfileOpen(false), 200)} // Simple delay to allow click
                        >
                            <span className="sr-only">Open user menu</span>
                            <div className="flex h-8 w-8 items-center justify-center rounded-full bg-blue-100 text-blue-600 font-bold border border-blue-200">
                                {user?.full_name?.charAt(0) || 'U'}
                            </div>
                            <span className="hidden lg:flex lg:items-center">
                                <span className="ml-4 text-sm font-semibold leading-6 text-gray-900" aria-hidden="true">
                                    {user?.full_name || 'User'}
                                </span>
                                <span className="ml-2 text-xs text-gray-500 bg-gray-100 px-2 py-0.5 rounded-full capitalize">
                                    {user?.role}
                                </span>
                            </span>
                        </button>

                        {/* Dropdown Menu */}
                        {isProfileOpen && (
                            <div className="absolute right-0 z-10 mt-2.5 w-48 origin-top-right rounded-md bg-white py-2 shadow-lg ring-1 ring-gray-900/5 focus:outline-none">
                                <div className="px-4 py-2 border-b mb-1 lg:hidden">
                                    <p className="text-sm font-semibold text-gray-900">{user?.full_name}</p>
                                    <p className="text-xs text-gray-500 capitalize">{user?.role}</p>
                                </div>

                                <Link
                                    to="/app/change-password"
                                    className="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-50"
                                >
                                    <KeyRound className="mr-2 h-4 w-4 text-gray-500" />
                                    Change Password
                                </Link>
                                <button
                                    onClick={handleLogout}
                                    className="flex w-full items-center px-4 py-2 text-sm text-red-600 hover:bg-red-50"
                                >
                                    <LogOut className="mr-2 h-4 w-4" />
                                    Sign out
                                </button>
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </header>
    );
};
