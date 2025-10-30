import React, { useState } from 'react';

// Import komponen kerangka (layout)
import Layout from '../src/components/Layout'; 
import '../src/pages/SetBudget';
// Import halaman aplikasi
import '../src/pages/RefrensiBelanja';
import DaftarBelanja from './pages/DaftarBelanja';
import Dashboard from '../src/pages/Dashboard'; 
import Reports from '../src/pages/Report'; 
// Import formulir autentikasi yang sudah dipisahkan
import LoginForm from './components/LoginForm'; 
import RegisterForm from './components/RegisterForm'; 
import SetBudget from '../src/pages/SetBudget';
import ReferensiBelanja from './pages/RefrensiBelanja';

/**
 * Komponen Utama Aplikasi (Akar)
 * Menangani state global (autentikasi, halaman aktif) dan routing sederhana.
 */
const App = () => {
    // State untuk status autentikasi pengguna
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    // State untuk navigasi antar halaman (nilai: 'dashboard' atau 'reports')
    const [currentPage, setCurrentPage] = useState('dashboard');
    // State untuk memilih tampilan auth (nilai: 'login' atau 'register')
    const [authView, setAuthView] = useState('login'); 

    // Fungsi yang dipanggil ketika login berhasil
    const handleLogin = () => {
        // Dalam proyek nyata, ini akan menyimpan token JWT ke LocalStorage
        setIsAuthenticated(true);
        setCurrentPage('dashboard'); 
    };

    // Fungsi yang dipanggil ketika tombol logout ditekan
    const handleLogout = () => {
        // Dalam proyek nyata, ini akan menghapus token dari LocalStorage
        setIsAuthenticated(false);
        setAuthView('login');
        console.log("Pengguna berhasil logout.");
    };

    // Logika untuk menampilkan komponen Halaman yang sesuai
    const renderPage = () => {
        switch (currentPage) {
            case 'dashboard':
                return <Dashboard onLogout={handleLogout} />;
            case 'reports':
                return <Reports onLogout={handleLogout} />;
            case'daftarbelanja':
                return <DaftarBelanja onLogout={handleLogout} />;
            case 'setBudget':
                return <SetBudget onLogout={handleLogout} />;
            case 'refrensi_belanja':
                return <ReferensiBelanja onLogout={handleLogout} />;
            default:
                return <Dashboard onLogout={handleLogout} />;

               
        }
    };

    // --- LOGIKA UTAMA RENDER ---

    // 1. TAMPILKAN LOGIN/REGISTER JIKA BELUM AUTHENTICATED
    if (!isAuthenticated) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-gray-50 font-sans">
                {authView === 'login' ? (
                    <LoginForm 
                        onLogin={handleLogin} 
                        onSwitchToRegister={() => setAuthView('register')}
                    />
                ) : (
                    <RegisterForm 
                        onSwitchToLogin={() => setAuthView('login')}
                    />
                )}
            </div>
        );
    }

    // 2. TAMPILKAN LAYOUT UTAMA JIKA SUDAH AUTHENTICATED
    return (
        <Layout 
            currentPage={currentPage} 
            setCurrentPage={setCurrentPage}
            onLogout={handleLogout} 
        >
            {renderPage()}
        </Layout>
    );
};

export default App;
