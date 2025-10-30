import React, { useState } from 'react';
import { 
    LayoutDashboard, 
    BarChart3, 
    LogOut, 
    Settings, 
    Menu, 
    X,
    ChevronDown,
    DollarSign,
    ShoppingCart, // <-- Ikon Baru
    Target,       // <-- Ikon Baru
    Tags,         // <-- Ikon Baru
    Users         // <-- Ikon Baru
} from 'lucide-react';

/**
 * Komponen Layout Utama (Responsif) - Sesuai Spek TK2
 * Menyediakan kerangka aplikasi (Sidebar, Header Mobile, Konten)
 *
 * @param {object} props
 * @param {React.Node} props.children - Komponen halaman
 * @param {string} props.currentPage - Halaman yang sedang aktif
 * @param {function} props.setCurrentPage - Fungsi untuk mengubah halaman aktif
 * @param {function} props.onLogout - Fungsi untuk menangani logout
 * @param {string} props.userRole - Peran pengguna ('user' atau 'admin')
 */
const Layout = ({ children, currentPage, setCurrentPage, onLogout, userRole = 'user' }) => {
    // State untuk mengelola visibilitas menu di mobile
    const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

    /**
     * Komponen NavLink (Tombol Menu)
     */
    const NavLink = ({ icon: Icon, label, isActive, onClick, isHidden = false }) => {
        if (isHidden) return null;
        return (
            <button
                onClick={onClick}
                className={`
                    flex items-center w-full px-4 py-3 rounded-lg transition-colors duration-150
                    ${isActive 
                        ? 'bg-indigo-700 text-white shadow-lg' 
                        : 'text-gray-300 hover:bg-indigo-500 hover:text-white'
                    }
                `}
            >
                <Icon className="w-5 h-5 mr-3 flex-shrink-0" />
                <span className="font-medium">{label}</span>
            </button>
        );
    };

    /**
     * Komponen SidebarContent
     * Berisi semua item menu, dipisahkan agar bisa di-reuse
     */
    const SidebarContent = () => (
        <div className="flex flex-col flex-grow p-4">
            {/* Logo/Judul (Sesuai TK2 - BMG) */}
            <div className="flex items-center justify-center h-16 mb-6">
                 <DollarSign className="w-10 h-10 text-white mr-2" />
                <span className="text-2xl font-bold text-white">
                    BMG
                </span>
            </div>
            
            {/* Navigasi Utama (Disesuaikan dengan Peran) */}
            <nav className="flex-1 space-y-3">
                <NavLink
                    icon={LayoutDashboard}
                    label="Dasbor"
                    isActive={currentPage === 'dashboard'}
                    onClick={() => {
                        setCurrentPage('dashboard');
                        setIsMobileMenuOpen(false);
                    }}
                />
                
                {/* --- MENU KHUSUS END USER (Sesuai TK2) --- */}
                <NavLink
                    icon={ShoppingCart}
                    label="Daftar Belanja"
                    isActive={currentPage === 'daftarbelanja'}
                    onClick={() => {
                        setCurrentPage('daftarbelanja');
                        setIsMobileMenuOpen(false);
                    }}
                    isHidden={userRole !== 'user'} // Sembunyikan jika admin
                />
                <NavLink
                    icon={Target}
                    label="Set Budget"
                    isActive={currentPage === 'setBudget'}
                    onClick={() => {
                        setCurrentPage('setBudget');
                        setIsMobileMenuOpen(false);
                    }}
                    isHidden={userRole !== 'user'} // Sembunyikan jika admin
                />
                <NavLink
                    icon={Tags}
                    label="Referensi Belanja"
                    isActive={currentPage === 'refrensi_belanja'}
                    onClick={() => {
                        setCurrentPage('refrensi_belanja');
                        setIsMobileMenuOpen(false);
                    }}
                    isHidden={userRole !== 'user'} // Sembunyikan jika admin
                />

                {/* --- MENU KHUSUS ADMIN (Sesuai TK2) --- */}
                <NavLink
                    icon={Users}
                    label="Kelola User"
                    isActive={currentPage === 'kelola_user'}
                    onClick={() => {
                        setCurrentPage('kelola_user');
                        setIsMobileMenuOpen(false);
                    }}
                    isHidden={userRole !== 'admin'} // Sembunyikan jika bukan admin
                />

                {/* --- MENU BERSAMA --- */}
                <NavLink
                    icon={BarChart3}
                    label="Laporan"
                    isActive={currentPage === 'reports'}
                    onClick={() => {
                        setCurrentPage('reports');
                        setIsMobileMenuOpen(false);
                    }}
                />
            </nav>
            
            {/* Navigasi Bawah (Pengaturan/Logout) */}
            <div className="mt-6">
                <NavLink
                    icon={Settings}
                    label="Pengaturan" // Sesuai Mockup TK2
                    isActive={currentPage === 'settings'}
                    onClick={() => { 
                        setCurrentPage('settings');
                        setIsMobileMenuOpen(false);
                    }}
                />
                <NavLink
                    icon={LogOut}
                    label="Keluar"
                    isActive={false}
                    onClick={onLogout}
                />
            </div>

            {/* Profil Pengguna */}
            <div className="border-t border-indigo-400 mt-6 pt-4">
                <div className="flex items-center">
                    <img 
                        className="w-10 h-10 rounded-full" 
                        src="https://placehold.co/100x100/6366f1/white?text=A" 
                        alt="Avatar Pengguna" 
                    />
                    <div className="ml-3">
                        <p className="text-sm font-medium text-white">User BMG</p>
                        {/* Menampilkan Peran (Role) */}
                        <p className="text-xs text-indigo-200 capitalize">{userRole}</p>
                    </div>
                    <ChevronDown className="w-4 h-4 text-indigo-200 ml-auto" />
                </div>
            </div>
        </div>
    );

    // ... (Sisa kode Layout tetap sama: Sidebar Mobile, Sidebar Desktop, Konten Utama) ...
    
    return (
        <div className="flex min-h-screen bg-gray-100 font-sans">
            
            {/* --- SIDEBAR MOBILE (Overlay) --- */}
            {isMobileMenuOpen && (
                <div 
                    className="fixed inset-0 z-30 bg-black/60 lg:hidden" 
                    onClick={() => setIsMobileMenuOpen(false)}
                    aria-hidden="true"
                ></div>
            )}
            <div className={`
                fixed inset-y-0 left-0 z-40 flex flex-col
                w-64 bg-indigo-600 shadow-2xl
                transform transition-transform duration-300 ease-in-out
                lg:hidden
                ${isMobileMenuOpen ? 'translate-x-0' : '-translate-x-full'}
            `}>
                <SidebarContent />
                <button 
                    onClick={() => setIsMobileMenuOpen(false)} 
                    className="absolute top-4 right-4 text-indigo-200 hover:text-white lg:hidden"
                >
                    <X className="w-6 h-6" />
                </button>
            </div>

            {/* --- SIDEBAR DESKTOP --- */}
            <div className="hidden lg:flex lg:flex-shrink-0 lg:flex-col w-64 bg-indigo-600 shadow-lg">
                <SidebarContent />
            </div>

            {/* --- KONTEN UTAMA & HEADER MOBILE --- */}
            <div className="flex-1 flex flex-col">
                
                {/* Header Mobile (Hanya terlihat di mobile) */}
                <header className="lg:hidden sticky top-0 z-20 bg-white shadow-md">
                    <div className="flex items-center justify-between px-4 py-4">
                        {/* Tombol Hamburger Menu */}
                        <button
                            onClick={() => setIsMobileMenuOpen(true)}
                            className="text-gray-700 focus:outline-none"
                        >
                            <Menu className="w-6 h-6" />
                        </button>
                        
                        {/* Judul Halaman (Contoh) */}
                        <h1 className="text-xl font-bold text-indigo-600 capitalize">
                            {/* --- PERBAIKAN DI SINI --- */}
                            {/* Menambahkan pemeriksaan untuk mencegah error jika currentPage undefined */}
                            {currentPage ? currentPage.replace('_', ' ') : 'Dasbor'}
                        </h1>
                        
                        {/* Ikon User (Contoh) */}
                        <img 
                            className="w-8 h-8 rounded-full" 
                            src="https://placehold.co/100x100/6366f1/white?text=A" 
                            alt="Avatar" 
                        />
                    </div>
                </header>

                {/* Konten Halaman (Dashboard.jsx atau Reports.jsx) */}
                <main className="flex-1 p-4 sm:p-6 lg:p-8">
                    {children}
                </main>
            </div>
        </div>
    );
};

export default Layout;

