import React, { useState, useEffect, useCallback } from 'react';
import { Loader, AlertCircle, Plus, Edit, Trash2, X } from 'lucide-react';

const API_BASE_URL = 'http://localhost:8080/api/v1';

/**
 * Halaman Referensi Belanja (sesuai TK2 Hal. 22)
 * Mengelola kategori belanja (CRUD).
 */
const ReferensiBelanja = () => {
    const [kategoriList, setKategoriList] = useState([]);
    const [namaKategori, setNamaKategori] = useState('');
    const [editingId, setEditingId] = useState(null); // ID kategori yang sedang diedit
    
    // Pisahkan state loading untuk list dan form
    const [isLoadingList, setIsLoadingList] = useState(true);
    const [isSubmitting, setIsSubmitting] = useState(false);
    
    const [error, setError] = useState(null); // Error untuk fetch list
    const [submitError, setSubmitError] = useState(null); // Error untuk form submit

    /**
     * Fungsi helper untuk mengambil data dengan otentikasi (JWT Token)
     */
    const fetchWithAuth = useCallback(async (url, options = {}) => {
        const token = localStorage.getItem('authToken');
        const headers = {
            'Content-Type': 'application/json',
            ...options.headers,
        };
        // Tambahkan token ke header jika ada
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }

        const response = await fetch(url, { ...options, headers });

        if (!response.ok) {
            // Coba parsing error JSON dari backend Go
            const errorData = await response.json().catch(() => ({})); 
            // Prioritaskan pesan error dari backend
            throw new Error(errorData.message || errorData.error || `Error ${response.status}: ${response.statusText}`);
        }
        
        // Handle body kosong
        const text = await response.text();
        return text ? JSON.parse(text) : null; 
    }, []);

    /**
     * Mengambil daftar kategori dari backend
     */
    const fetchKategori = useCallback(async () => {
        setIsLoadingList(true);
        setError(null);
        try {
            // Asumsi endpoint: GET /api/v1/kategori
            const data = await fetchWithAuth(`${API_BASE_URL}/kategori`);
            if (data && data.data) {
                setKategoriList(data.data);
            } else {
                setKategoriList([]); // Set array kosong jika data null
            }
        } catch (err) {
            setError(err.message);
            console.error("Gagal mengambil kategori:", err);
        } finally {
            setIsLoadingList(false);
        }
    }, [fetchWithAuth]);

    // Mengambil data saat komponen dimuat
    useEffect(() => {
        fetchKategori();
    }, [fetchKategori]);

    /**
     * Menangani submit form (Tambah Kategori Baru atau Update)
     */
    const handleSubmit = async (e) => {
        e.preventDefault();
        if (!namaKategori.trim()) {
            setSubmitError("Nama kategori tidak boleh kosong.");
            return;
        }
        
        setIsSubmitting(true); 
        setSubmitError(null);

        const url = editingId 
            ? `${API_BASE_URL}/kategori/${editingId}` // Endpoint UPDATE (PUT)
            : `${API_BASE_URL}/kategori`; // Endpoint CREATE (POST)
        
        const method = editingId ? 'PUT' : 'POST';

        try {
            await fetchWithAuth(url, {
                method: method,
                body: JSON.stringify({ nama_kategori: namaKategori })
            });
            
            // Sukses: reset form dan ambil ulang data
            setNamaKategori('');
            setEditingId(null);
            await fetchKategori(); // Refresh daftar

        } catch (err) {
            setSubmitError(err.message);
            console.error("Gagal submit kategori:", err);
        } finally {
            setIsSubmitting(false); // Sembunyikan loader di tombol
        }
    };

    /**
     * Menangani klik tombol Hapus
     */
    const handleDelete = async (id) => {
        // Tampilkan konfirmasi dialog
        if (window.confirm('Apakah Anda yakin ingin menghapus kategori ini? Item belanja terkait mungkin akan terpengaruh.')) {
            setIsLoadingList(true); // Tampilkan loader umum
            setError(null);
            try {
                // Asumsi endpoint: DELETE /api/v1/kategori/{id}
                await fetchWithAuth(`${API_BASE_URL}/kategori/${id}`, {
                    method: 'DELETE'
                });
                await fetchKategori(); // Refresh daftar
            } catch (err) {
                setError(err.message);
                console.error("Gagal menghapus kategori:", err);
            } finally {
                setIsLoadingList(false);
            }
        }
    };

    /**
     * Menangani klik tombol Edit (Mengisi form)
     */
    const handleEdit = (kategori) => {
        setEditingId(kategori.id_kategori);
        setNamaKategori(kategori.nama_kategori);
        setSubmitError(null);
        // Scroll ke atas ke form
        window.scrollTo({ top: 0, behavior: 'smooth' }); 
    };
    
    /**
     * Membatalkan mode edit
     */
    const handleCancelEdit = () => {
        setEditingId(null);
        setNamaKategori('');
        setSubmitError(null);
    };

    // Tampilan Utama
    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <h1 className="text-3xl font-extrabold text-gray-900">Referensi Belanja (Kategori)</h1>
            </div>
            <p className="text-gray-600">Kelola kategori untuk item belanja Anda. (Sesuai TK2 Hal. 22)</p>

            {/* Form Tambah/Edit Kategori Baru (sesuai mockup) */}
            <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-200">
                <h2 className="text-xl font-semibold text-gray-800 mb-4 border-b pb-3">
                    {editingId ? 'Edit Kategori' : 'Tambah Kategori Baru'}
                </h2>
                
                <form onSubmit={handleSubmit} className="flex flex-col sm:flex-row sm:items-end gap-4">
                    <div className="flex-grow">
                        <label htmlFor="namaKategori" className="block text-sm font-medium text-gray-700">Nama Kategori</label>
                        <input
                            type="text"
                            id="namaKategori"
                            value={namaKategori}
                            onChange={(e) => setNamaKategori(e.target.value)}
                            className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                            placeholder="Contoh: Makanan Pokok, Transportasi, Hiburan..."
                            required
                        />
                    </div>
                    <div className="flex items-center gap-2 flex-shrink-0">
                        {editingId && (
                            <button
                                type="button"
                                onClick={handleCancelEdit}
                                className="py-2 px-4 rounded-lg font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 transition duration-150 flex items-center"
                            >
                                <X className="w-4 h-4 mr-1" /> Batal
                            </button>
                        )}
                        <button
                            type="submit"
                            disabled={isSubmitting}
                            className="flex items-center justify-center py-2 px-4 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:bg-gray-400 disabled:cursor-not-allowed transition duration-150"
                        >
                            {isSubmitting ? (
                                <Loader className="animate-spin h-5 w-5 mr-2" /> 
                            ) : (
                                editingId ? <Edit className="h-5 w-5 mr-1" /> : <Plus className="h-5 w-5 mr-1" />
                            )}
                            {editingId ? 'Simpan Perubahan' : 'Tambahkan'}
                        </button>
                    </div>
                </form>
                {submitError && (
                    <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded-lg relative mt-4" role="alert">
                        {submitError}
                    </div>
                )}
            </div>

            {/* Daftar Kategori (sesuai mockup) */}
            <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-200">
                <h2 className="text-xl font-semibold text-gray-800 mb-4 border-b pb-3">Daftar Kategori</h2>
                
                {isLoadingList && (
                    <div className="flex justify-center items-center h-40">
                        <Loader className="animate-spin h-8 w-8 text-indigo-600" />
                    </div>
                )}

                {error && (
                    <div className="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 rounded-lg" role="alert">
                        <div className="flex">
                            <div className="py-1"><AlertCircle className="h-6 w-6 text-red-500 mr-3" /></div>
                            <div>
                                <p className="font-bold">Gagal mengambil data</p>
                                <p className="text-sm">{error}</p>
                            </div>
                        </div>
                    </div>
                )}

                {!isLoadingList && !error && kategoriList.length === 0 && (
                    <p className="text-center text-gray-500 py-10">Belum ada kategori yang ditambahkan.</p>
                )}

                {kategoriList.length > 0 && (
                    <div className="overflow-x-auto">
                        <table className="min-w-full divide-y divide-gray-200">
                            <thead className="bg-gray-50">
                                <tr>
                                    <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                        Nama Kategori
                                    </th>
                                    <th scope="col" className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                                        Aksi
                                    </th>
                                </tr>
                            </thead>
                            <tbody className="bg-white divide-y divide-gray-200">
                                {kategoriList.map((kategori) => (
                                    <tr key={kategori.id_kategori}>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                                            {kategori.nama_kategori}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium space-x-2">
                                            <button 
                                                onClick={() => handleEdit(kategori)}
                                                className="text-indigo-600 hover:text-indigo-900"
                                                title="Edit"
                                                disabled={isSubmitting}
                                            >
                                                <Edit className="w-5 h-5" />
                                            </button>
                                            <button 
                                                onClick={() => handleDelete(kategori.id_kategori)}
                                                className="text-red-600 hover:text-red-900"
                                                title="Hapus"
                                                disabled={isSubmitting}
                                            >
                                                <Trash2 className="w-5 h-5" />
                                            </button>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )}
            </div>
        </div>
    );
};

export default ReferensiBelanja;


