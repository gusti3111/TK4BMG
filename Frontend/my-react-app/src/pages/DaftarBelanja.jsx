import React, { useState, useEffect, useCallback } from 'react';
import { Plus, Trash2, Edit, Loader, AlertCircle } from 'lucide-react';

// API Base URL (menggunakan proxy Nginx dari Docker)
const API_BASE_URL = 'http://localhost:8080/api/v1';

/**
 * Komponen Halaman DaftarBelanja (sesuai TK2 - Hal 22)
 * Mengelola CRUD untuk item belanja mingguan.
 */
const DaftarBelanja = () => {
    // State untuk data dari backend
    const [items, setItems] = useState([]);
    const [kategoriList, setKategoriList] = useState([]);
    
    // State untuk formulir
     const [kategoriId, setKategoriId] = useState('');
     const [idUser, setIdUser] = useState(1); // Asumsi user ID 1 untuk demo
    const [namaBarang, setNamaBarang] = useState('');
  
    const [harga, setHarga] = useState('');

    // State untuk UI
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState(null);
    const [submitError, setSubmitError] = useState(null);
    const [isSubmitting, setIsSubmitting] = useState(false);

    /**
     * Fungsi helper untuk mengambil data dengan otentikasi
     */
    const fetchWithAuth = useCallback(async (url, options = {}) => {
        const token = localStorage.getItem('authToken');
        const headers = {
            'Content-Type': 'application/json',
            ...options.headers,
        };

        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }

        const response = await fetch(url, { ...options, headers });

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            throw new Error(errorData.message || errorData.error || `Error: ${response.status}`);
        }
        
        // Menangani respons 204 No Content (misalnya untuk DELETE)
        if (response.status === 204) {
            return null;
        }

        return response.json();
    }, []);

    /**
     * Mengambil semua data (item dan kategori) saat komponen dimuat
     */
    const fetchData = useCallback(async () => {
        setIsLoading(true);
        setError(null);
        try {
            // 1. Ambil daftar item belanja (Asumsi endpoint /items)
            const itemsData = await fetchWithAuth(`${API_BASE_URL}/items`);
            setItems(itemsData.data || []); // Sesuaikan .data jika perlu

            // 2. Ambil daftar kategori (Asumsi endpoint /referensi)
            const kategoriData = await fetchWithAuth(`${API_BASE_URL}/referensi`);
            setKategoriList(kategoriData.data || []); // Sesuaikan .data jika perlu
            
            // Set kategori default jika ada
            if (kategoriData.data && kategoriData.data.length > 0) {
                setKategoriId(kategoriData.data[0].id_kategori); // Sesuaikan .id_kategori
            }

        } catch (err) {
            setError(err.message);
        } finally {
            setIsLoading(false);
        }
    }, [fetchWithAuth]);

    // Hook untuk menjalankan fetchData() saat komponen dimuat
    useEffect(() => {
        fetchData();
    }, [fetchData]);

    /**
     * Menangani submit formulir "Tambah Item Belanja"
     */
    const handleSubmit = async (e) => {
        e.preventDefault();
        setIsSubmitting(true);
        setSubmitError(null);

        try {
            // Sesuai ERD TK2: butuh nama_item, id_kategori, harga_satuan
            // (Asumsi jumlah_item default 1 dan backend menghitung total_harga)
            const newItem = {
               
                id_kategori: parseInt(kategoriId, 10),
                id_user: idUser, // Asumsi user ID 1 untuk demo
                nama_item: namaBarang,
                jumlah_item: jumlahItem, // Asumsi default 1
                harga_satuan: parseFloat(harga),
                total_harga: parseFloat(harga), // Asumsi total sama dengan harga satuan untuk 1 item
                purchase_date: new Date().toISOString().split('T')[0], // Tanggal hari ini
               
            };
            
            // Kirim item baru ke backend (POST /api/v1/items)
            const createdItem = await fetchWithAuth(`${API_BASE_URL}/items`, {
                method: 'POST',
                body: JSON.stringify(newItem),
            });

            // Tambahkan item baru ke state (Optimistic UI)
            setItems(prevItems => [...prevItems, createdItem.data]); // Sesuaikan .data

            // Reset formulir
            setNamaBarang('');
            setHarga('');
            
        } catch (err) {
            setSubmitError(err.message);
        } finally {
            setIsSubmitting(false);
        }
    };

    /**
     * Menangani penghapusan item
     */
    const handleDelete = async (itemId) => {
        // Konfirmasi sebelum menghapus
        if (!window.confirm('Apakah Anda yakin ingin menghapus item ini?')) {
            return;
        }

        try {
            // Kirim request DELETE (DELETE /api/v1/items/:id)
            await fetchWithAuth(`${API_BASE_URL}/items/${itemId}`, {
                method: 'DELETE',
            });

            // Hapus item dari state (Optimistic UI)
            setItems(prevItems => prevItems.filter(item => item.id_item !== itemId)); // Sesuaikan .id_item

        } catch (err) {
            setError(`Gagal menghapus item: ${err.message}`);
        }
    };

    // Tampilan Loading
    if (isLoading) {
        return (
            <div className="flex justify-center items-center h-64">
                <Loader className="animate-spin h-12 w-12 text-indigo-600" />
            </div>
        );
    }
    
    // Tampilan Error Utama
    if (error) {
        return (
            <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded-lg" role="alert">
                <strong className="font-bold">Error: </strong>
                <span className="block sm:inline">{error}</span>
            </div>
        );
    }

    // Tampilan Utama (sesuai mockup TK2)
    return (
        <div className="space-y-8">
            <h1 className="text-3xl font-extrabold text-gray-900">Daftar Belanja</h1>

            {/* 1. Formulir Tambah Item Belanja */}
            <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-200">
                <h2 className="text-xl font-semibold text-gray-800 mb-4">Tambah Item Belanja</h2>
                
                {submitError && (
                    <div className="bg-red-100 text-red-700 px-4 py-3 rounded-lg mb-4" role="alert">
                        {submitError}
                    </div>
                )}

                <form onSubmit={handleSubmit} className="grid grid-cols-1 md:grid-cols-4 gap-4 items-end">
                    {/* Nama Barang */}
                    <div className="md:col-span-2">
                        <label htmlFor="namaBarang" className="block text-sm font-medium text-gray-700">Nama Barang</label>
                        <input
                            type="text"
                            id="namaBarang"
                            value={namaBarang}
                            onChange={(e) => setNamaBarang(e.target.value)}
                            className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                            placeholder="cth: Apel Fuji"
                            required
                        />
                    </div>
                    
                    {/* Kategori */}
                    <div>
                        <label htmlFor="kategori" className="block text-sm font-medium text-gray-700">Kategori</label>
                        <select
                            id="kategori"
                            value={kategoriId}
                            onChange={(e) => setKategoriId(e.target.value)}
                            className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 bg-white"
                            required
                        >
                            <option value="" disabled>Pilih Kategori</option>
                            {kategoriList.map(kat => (
                                <option key={kat.id_kategori} value={kat.id_kategori}>
                                    {kat.nama_kategori}
                                </option>
                            ))}
                        </select>
                    </div>

                    {/* Harga */}
                    <div>
                        <label htmlFor="harga" className="block text-sm font-medium text-gray-700">Harga</label>
                        <input
                            type="number"
                            id="harga"
                            value={harga}
                            onChange={(e) => setHarga(e.target.value)}
                            className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                            placeholder="cth: 50000"
                            required
                        />
                    </div>

                    {/* Tombol Submit */}
                    <div className="md:col-start-4">
                        <button
                            type="submit"
                            disabled={isSubmitting}
                            className="w-full flex justify-center items-center px-4 py-2 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:bg-gray-400"
                        >
                            {isSubmitting ? (
                                <Loader className="animate-spin h-5 w-5 mr-2" />
                            ) : (
                                <Plus className="h-5 w-5 mr-2" />
                            )}
                            Tambah
                        </button>
                    </div>
                </form>
            </div>

            {/* 2. Tabel Daftar Belanja Mingguan */}
            <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-200">
                <h2 className="text-xl font-semibold text-gray-800 mb-4">Daftar Belanja Mingguan</h2>
                
                <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                            <tr>
                                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Nama Barang</th>
                                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Kategori</th>
                                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Harga</th>
                                <th scope="col" className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Aksi</th>
                            </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                            {items.length === 0 ? (
                                <tr>
                                    <td colSpan="4" className="px-6 py-4 text-center text-gray-500">
                                        Belum ada item belanja.
                                    </td>
                                </tr>
                            ) : (
                                items.map((item) => (
                                    <tr key={item.id_item}>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{item.nama_item}</td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                            {/* Menampilkan nama kategori berdasarkan ID */}
                                            {kategoriList.find(k => k.id_kategori === item.id_kategori)?.nama_kategori || 'N/A'}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                            {/* Format harga ke Rupiah */}
                                            {new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(item.harga_satuan)}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium space-x-2">
                                            <button 
                                                // onClick={() => handleEdit(item.id_item)} 
                                                className="text-indigo-600 hover:text-indigo-900"
                                                title="Edit"
                                            >
                                                <Edit className="h-5 w-5" />
                                            </button>
                                            <button 
                                                onClick={() => handleDelete(item.id_item)} 
                                                className="text-red-600 hover:text-red-900"
                                                title="Hapus"
                                            >
                                                <Trash2 className="h-5 w-5" />
                                            </button>
                                        </td>
                                    </tr>
                                ))
                            )}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    );
};

export default DaftarBelanja;
