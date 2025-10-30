import React, { useState } from 'react';
import { User, Mail, Key, Loader } from 'lucide-react';

// Asumsi URL dasar backend Go Anda
const API_BASE_URL = 'http://localhost:8080/api/v1';

/**
 * Komponen Formulir Pendaftaran (RegisterForm)
 * Mengelola state pendaftaran dan menangani submit ke backend.
 * * @param {object} props
 * @param {function} props.onSwitchToLogin - Fungsi untuk beralih ke tampilan Login.
 */
const RegisterForm = ({ onSwitchToLogin }) => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [nama, setName] = useState('');
    const [email, setEmail] = useState('');
    
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState('');

    /**
     * Menangani submit formulir pendaftaran.
     */
    const handleSubmit = async (e) => {
        e.preventDefault();
        setIsLoading(true);
        setError('');

        // 1. Validasi input
        if (!username || !password ||!nama || !email ) {
            setError('Semua field wajib diisi.');
            setIsLoading(false);
            return;
        }

        try {
            // --- INTEGRASI BACKEND GO DIMULAI ---
            const response = await fetch(`${API_BASE_URL}/register`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ username,password, nama, email  })
            });

            const data = await response.json();

            if (!response.ok) {
                // Ambil pesan error dari backend (asumsi backend mengirim format { "message": "..." })
                throw new Error(data.message || data.error || 'Pendaftaran gagal');
            }
            
            console.log('Pendaftaran Berhasil:', data);

            // Jika berhasil, panggil onSwitchToLogin untuk mengarahkan pengguna ke form Login
            if(onSwitchToLogin) {
                onSwitchToLogin();
            }
            // --- AKHIR INTEGRASI BACKEND ---

        } catch (err) {
            // Menampilkan error dari backend atau network error
            setError(err.message || 'Terjadi kesalahan. Silakan coba lagi.');
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="max-w-md w-full bg-white shadow-xl rounded-2xl p-8 border border-gray-200 font-sans">
            <h2 className="text-3xl font-extrabold text-center text-gray-900 mb-2">
                Buat Akun Baru
            </h2>
            <p className="text-center text-sm text-gray-600 mb-6">
                Silakan isi detail Anda untuk memulai.
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
                            <User className="h-5 w-5 text-gray-400" />
                        </div>
                        <input
                            id="username"
                            name="username"
                            type="text"
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
                    <label htmlFor="password-register" className="text-sm font-medium text-gray-700">
                        Password
                    </label>
                    <div className="mt-1 relative rounded-lg shadow-sm">
                        <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                            <Key className="h-5 w-5 text-gray-400" />
                        </div>
                        <input
                            id="password-register"
                            name="password"
                            type="password"
                            autoComplete="new-password"
                            required
                            className="appearance-none block w-full px-3 py-3 pl-10 border border-gray-300 rounded-lg placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                            placeholder="Buat password Anda"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            disabled={isLoading}
                        />
                    </div>
                </div>

                <div>
                    <label htmlFor="name" className="text-sm font-medium text-gray-700">
                        Nama Lengkap
                    </label>
                    
                    <div className="mt-1 relative rounded-lg shadow-sm">
                        <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                            <User className="h-5 w-5 text-gray-400" />
                        </div>
                        <input
                            id="nama"
                            name="nama"
                            type="text"
                            required
                            className="appearance-none block w-full px-3 py-3 pl-10 border border-gray-300 rounded-lg placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                            placeholder="Masukkan nama lengkap Anda"
                            value={nama}
                            onChange={(e) => setName(e.target.value)}
                            disabled={isLoading}
                        />
                    </div>
                </div>

                <div>
                    <label htmlFor="email-register" className="text-sm font-medium text-gray-700">
                        Alamat Email
                    </label>
                    <div className="mt-1 relative rounded-lg shadow-sm">
                        <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                            <Mail className="h-5 w-5 text-gray-400" />
                        </div>
                        <input
                            id="email-register"
                            name="email"
                            type="email"
                            autoComplete="email"
                            required
                            className="appearance-none block w-full px-3 py-3 pl-10 border border-gray-300 rounded-lg placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                            placeholder="email@example.com"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
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
                            'Daftar Sekarang'
                        )}
                    </button>
                </div>
            </form>

            <p className="mt-8 text-center text-sm text-gray-600">
                Sudah punya akun?{' '}
                <button
                    onClick={onSwitchToLogin}
                    disabled={isLoading}
                    className="font-medium text-indigo-600 hover:text-indigo-500 focus:outline-none disabled:text-gray-400"
                >
                    Masuk di sini
                </button>
            </p>
        </div>
    );
};

export default RegisterForm;

