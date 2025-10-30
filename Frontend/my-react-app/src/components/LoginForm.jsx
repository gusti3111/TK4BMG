import React, { useState } from 'react';
import { Mail, Key, Loader } from 'lucide-react';

// Asumsi URL dasar backend Go Anda
const API_BASE_URL = 'http://localhost:8080/api/v1';

/**
 * Komponen Formulir Login (LoginForm)
 * Mengelola state login dan menangani submit ke backend.
 *
 * @param {object} props
 * @param {function} props.onLogin - Fungsi callback yang dipanggil setelah login berhasil.
 * @param {function} props.onSwitchToRegister - Fungsi untuk beralih ke tampilan Register.
 */
const LoginForm = ({ onLogin, onSwitchToRegister }) => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState('');

    /**
     * Menangani submit formulir login.
     */
    const handleSubmit = async (e) => {
        e.preventDefault();
        setIsLoading(true);
        setError('');

        // --- INTEGRASI BACKEND GO DIMULAI ---
        try {
            const response = await fetch(`${API_BASE_URL}/login`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ username, password })
            });

            const data = await response.json();

            if (!response.ok) {
                // Ambil pesan error dari backend
                throw new Error(data.message || data.error || 'Login gagal. Periksa email dan password.');
            }

            // Asumsi backend mengembalikan token JWT
            if (data.token) {
                // Simpan token ke localStorage untuk sesi berikutnya
                localStorage.setItem('authToken', data.token);
                
                // Panggil callback onLogin (dari App.jsx) untuk mengubah state isAuthenticated
                if(onLogin) {
                    onLogin(); 
                }
            } else {
                throw new Error('Token tidak diterima dari server.');
            }

        } catch (err) {
            // Menampilkan error dari backend atau network error
            setError(err.message || 'Terjadi kesalahan. Silakan coba lagi.');
        } finally {
            setIsLoading(false);
        }
        // --- AKHIR INTEGRASI BACKEND ---
    };

    // --- PERUBAHAN DIMULAI: Wrapper untuk memusatkan ---
    // Div ini ditambahkan agar form berada di tengah saat dipratinjau terpisah
    return (
        <div className="min-h-screen w-full flex items-center justify-center bg-gray-50 font-sans p-4">
            <div className="max-w-md w-full bg-white shadow-xl rounded-2xl p-8 border border-gray-200">
                <h2 className="text-3xl font-extrabold text-center text-gray-900 mb-2">
                    Selamat Datang Kembali
                </h2>
                <p className="text-center text-sm text-gray-600 mb-6">
                    Masuk ke akun Anda untuk melanjutkan.
                </p>
                
                {error && (
                    <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded-lg relative mb-4" role="alert">
                        <span className="block sm:inline">{error}</span>
                    </div>
                )}

                <form className="space-y-6" onSubmit={handleSubmit}>
                    <div>
                        <label htmlFor="username" className="text-sm font-medium text-gray-700">
                            Username
                        </label>
                        <div className="mt-1 relative rounded-lg shadow-sm">
                            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                <Mail className="h-5 w-5 text-gray-400" />
                            </div>
                            <input
                                id="username"
                                name="username"
                                type="text"
                                autoComplete="username"
                                required
                                className="appearance-none block w-full px-3 py-3 pl-10 border border-gray-300 rounded-lg placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                                placeholder="Masukkan username Anda"
                                value={username}
                                onChange={(e) => setUsername(e.target.value)}
                                disabled={isLoading}
                            />
                        </div>
                    </div>

                    <div>
                        <label htmlFor="password" className="text-sm font-medium text-gray-700">
                            Password
                        </label>
                        <div className="mt-1 relative rounded-lg shadow-sm">
                            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                <Key className="h-5 w-5 text-gray-400" />
                            </div>
                            <input
                                id="password"
                                name="password"
                                type="password"
                                autoComplete="current-password"
                                required
                                className="appearance-none block w-full px-3 py-3 pl-10 border border-gray-300 rounded-lg placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                                placeholder="Masukkan password Anda"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                disabled={isLoading}
                            />
                        </div>
                    </div>

                    <div>
                        <button
                            type="submit"
                            disabled={isLoading}
                            className="w-full flex justify-center py-3 px-4 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:bg-gray-400 disabled:cursor-not-allowed transition duration-150"
                        >
                            {isLoading ? (
                                <Loader className="animate-spin h-5 w-5 text-white" />
                            ) : (
                                'Masuk'
                            )}
                        </button>
                    </div>
                </form>

                <p className="mt-8 text-center text-sm text-gray-600">
                    Belum punya akun?{' '}
                    <button
                        onClick={onSwitchToRegister}
                        disabled={isLoading}
                        className="font-medium text-indigo-600 hover:text-indigo-500 focus:outline-none disabled:text-gray-400"
                    >
                        Daftar di sini
                    </button>
                </p>
            </div>
        </div>
    );
    // --- AKHIR PERUBAHAN ---
};

export default LoginForm;

