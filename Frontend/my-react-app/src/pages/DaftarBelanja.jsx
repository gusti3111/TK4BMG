import React, { useState, useEffect, useCallback } from 'react';
import { Loader, AlertCircle, Plus, Trash2, Edit, X, Save, Tag } from 'lucide-react';

// API Base URL (untuk development lokal)
// Pastikan backend Go Anda berjalan di port 8080 dan memiliki CORS
const API_BASE_URL = 'http://localhost:8080/api/v1';

/**
 * Halaman Daftar Belanja (sesuai TK2 Hal. 22)
 * Mengelola (CRUD) item belanja mingguan.
 */
const DaftarBelanja = () => {
    // State untuk daftar item dan kategori
    const [items, setItems] = useState([]);
    const [kategoriList, setKategoriList] = useState([]);
    
    // State untuk form "Tambah Item"
    const [namaItem, setNamaItem] = useState('');
    const [idKategori, setIdKategori] = useState('');
    const [jumlah, setJumlah] = useState(1);
    const [harga, setHarga] = useState('');

    // State untuk modal Edit
    const [isEditing, setIsEditing] = useState(false);
    const [editingItem, setEditingItem] = useState(null);

    // State UI
    const [isLoading, setIsLoading] = useState(true);
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [error, setError] = useState(null);

    /**
     * Fungsi helper untuk mengambil data dengan otentikasi (JWT Token)
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
            throw new Error(errorData.message || errorData.error || `Error ${response.status}: ${response.statusText}`);
        }
        
        const text = await response.text();
        return text ? JSON.parse(text) : null;
    }, []);

    /**
     * Mengambil data awal (item belanja dan kategori)
     */
    const fetchData = useCallback(async () => {
        setIsLoading(true);
        setError(null);
        try {
            // 1. Ambil daftar kategori (untuk dropdown)
            // Sesuai 'category_repository.go' -> GET /api/v1/kategori
            const kategoriData = await fetchWithAuth(`${API_BASE_URL}/kategori`);
            setKategoriList(kategoriData?.data || []);

            // 2. Ambil daftar item belanja
            // Sesuai 'item_repository.go' -> GET /api/v1/items
            const itemsData = await fetchWithAuth(`${API_BASE_URL}/items`);
            setItems(itemsData?.data || []);

        } catch (err) {
            setError(err.message);
            console.error("Gagal mengambil data:", err);
        } finally {
            setIsLoading(false);
        }
    }, [fetchWithAuth]);

    // Mengambil data saat komponen dimuat
    useEffect(() => {
        fetchData();
    }, [fetchData]);

    /**
     * Menangani submit form "Tambah Item Belanja"
     */
    const handleAddItem = async (e) => {
        e.preventDefault();
        if (!namaItem || !idKategori || jumlah <= 0 || harga <= 0) {
            setError("Semua field wajib diisi dan harus valid.");
            return;
        }

        setIsSubmitting(true);
        setError(null);

        try {
            // Sesuai 'item_repository.go' -> POST /api/v1/items
            await fetchWithAuth(`${API_BASE_URL}/items`, {
                method: 'POST',
                body: JSON.stringify({
                    nama_item: namaItem,
                    id_kategori: parseInt(idKategori, 10),
                    jumlah_item: parseInt(jumlah, 10),
                    harga_satuan: parseFloat(harga)
                    // total_harga akan dihitung di backend (Handler/Service)
                })
            });

            // Reset form dan ambil ulang data
            setNamaItem('');
            setIdKategori('');
            setJumlah(1);
            setHarga('');
            await fetchData(); // Refresh tabel

        } catch (err) {
            setError(err.message);
            console.error("Gagal menambah item:", err);
        } finally {
            setIsSubmitting(false);
        }
    };
    
    /**
     * Menangani Hapus Item
     */
    const handleDeleteItem = async (itemID) => {
        if (!window.confirm("Apakah Anda yakin ingin menghapus item ini?")) {
            return;
        }
        
        try {
            // Sesuai 'item_repository.go' -> DELETE /api/v1/items/:id
            await fetchWithAuth(`${API_BASE_URL}/items/${itemID}`, {
                method: 'DELETE'
            });
            await fetchData(); // Refresh tabel
        } catch (err) {
            setError(err.message);
            console.error("Gagal menghapus item:", err);
        }
    };

    /**
     * Membuka modal Edit
     */
    const openEditModal = (item) => {
        setEditingItem({
            ...item,
            id_kategori: item.id_kategori || '', // Handle null category
            jumlah_item: item.jumlah_item || 1,
            harga_satuan: item.harga_satuan || 0
        });
        setIsEditing(true);
    };

    /**
     * Menangani submit form Edit
     */
    const handleUpdateItem = async (e) => {
        e.preventDefault();
        setIsSubmitting(true);
        setError(null);

        try {
            // Sesuai 'item_repository.go' -> PUT /api/v1/items/:id
            await fetchWithAuth(`${API_BASE_URL}/items/${editingItem.id_item}`, {
                method: 'PUT',
                body: JSON.stringify({
                    nama_item: editingItem.nama_item,
                    id_kategori: parseInt(editingItem.id_kategori, 10),
                    jumlah_item: parseInt(editingItem.jumlah_item, 10),
                    harga_satuan: parseFloat(editingItem.harga_satuan)
                })
            });

            setIsEditing(false);
            setEditingItem(null);
            await fetchData(); // Refresh tabel

        } catch (err) {
            setError(err.message);
            console.error("Gagal update item:", err);
        } finally {
            setIsSubmitting(false);
        }
    };

    // Format mata uang Rupiah
    const formatIDR = (value) => 
        new Intl.NumberFormat('id-ID', { 
            style: 'currency', 
            currency: 'IDR', 
            minimumFractionDigits: 0 
        }).format(value || 0);

    // Mencari nama kategori berdasarkan ID
    const getCategoryName = (kategoriId) => {
        const kategori = kategoriList.find(k => k.id_kategori === kategoriId);
        return kategori ? kategori.nama_kategori : 'Tanpa Kategori';
    };

    // Tampilan Loading
    if (isLoading) {
        return (
            <div className="flex justify-center items-center h-64">
                <Loader className="animate-spin h-12 w-12 text-indigo-600" />
            </div>
        );
    }

    // Tampilan Utama (sesuai mockup TK2 Hal. 22)
    return (
        <div className="space-y-6">
            <h1 className="text-3xl font-extrabold text-gray-900">Daftar Belanja</h1>
            
            {error && (
                <div className="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 rounded-lg" role="alert">
                    <p className="font-bold">Gagal memproses permintaan</p>
                    <p>{error}</p>
                </div>
            )}

            {/* 1. Form "Tambah Item Belanja" (sesuai mockup) */}
            <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-200">
                <h2 className="text-xl font-semibold text-gray-800 mb-4 border-b pb-3">
                    Tambah Item Belanja
                </h2>
                <form onSubmit={handleAddItem} className="grid grid-cols-1 md:grid-cols-4 gap-4 items-end">
                    {/* Nama Barang */}
                    <div className="md:col-span-2">
                        <label htmlFor="namaItem" className="block text-sm font-medium text-gray-700">Nama Barang</label>
                        <input
                            type="text"
                            id="namaItem"
                            value={namaItem}
                            onChange={(e) => setNamaItem(e.target.value)}
                            className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                            placeholder="Contoh: Beras"
                            required
                        />
                    </div>
                    {/* Kategori */}
                    <div>
                        <label htmlFor="idKategori" className="block text-sm font-medium text-gray-700">Kategori</label>
                        <select
                            id="idKategori"
                            value={idKategori}
                            onChange={(e) => setIdKategori(e.target.value)}
                            className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                            required
                        >
                            <option value="" disabled>Pilih Kategori</option>
                            {kategoriList.map(k => (
                                <option key={k.id_kategori} value={k.id_kategori}>{k.nama_kategori}</option>
                            ))}
                        </select>
                    </div>
                    {/* Jumlah */}
                    <div>
                        <label htmlFor="jumlah" className="block text-sm font-medium text-gray-700">Jumlah</label>
                        <input
                            type="number"
                            id="jumlah"
                            value={jumlah}
                            onChange={(e) => setJumlah(e.target.value)}
                            className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                            min="1"
                            required
                        />
                    </div>
                    {/* Harga Satuan */}
                    <div className="md:col-span-2">
                        <label htmlFor="harga" className="block text-sm font-medium text-gray-700">Harga Satuan (Rp)</label>
                        <input
                            type="number"
                            id="harga"
                            value={harga}
                            onChange={(e) => setHarga(e.target.value)}
                            className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                            placeholder="15000"
                            min="0"
                            required
                        />
                    </div>
                    {/* Tombol Submit */}
                    <div className="md:col-span-2 flex justify-end">
                        <button
                            type="submit"
                            disabled={isSubmitting}
                            className="flex items-center justify-center w-full md:w-auto py-2 px-6 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:bg-gray-400"
                        >
                            {isSubmitting ? <Loader className="animate-spin h-5 w-5 mr-2" /> : <Plus className="h-5 w-5 mr-2" />}
                            {isSubmitting ? 'Menyimpan...' : 'Tambah'}
                        </button>
                    </div>
                </form>
            </div>

            {/* 2. Tabel "Daftar Belanja Mingguan" (sesuai mockup) */}
            <div className="bg-white rounded-xl shadow-lg border border-gray-200 overflow-hidden">
                <h2 className="text-xl font-semibold text-gray-800 p-6 border-b">
                    Daftar Belanja Mingguan
                </h2>
                <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                            <tr>
                                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Nama Barang</th>
                                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Kategori</th>
                                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Jumlah</th>
                                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Harga Satuan</th>
                                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Total Harga</th>
                                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Aksi</th>
                            </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                            {items.length === 0 ? (
                                <tr>
                                    <td colSpan="6" className="px-6 py-4 text-center text-gray-500">
                                        Belum ada data belanja.
                                    </td>
                                </tr>
                            ) : (
                                items.map((item) => (
                                    <tr key={item.id_item}>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{item.nama_item}</td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                            <span className="flex items-center">
                                                <Tag className="h-4 w-4 mr-1.5 text-gray-400" />
                                                {getCategoryName(item.id_kategori)}
                                            </span>
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{item.jumlah_item}</td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{formatIDR(item.harga_satuan)}</td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-gray-700">{formatIDR(item.total_harga)}</td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium space-x-2">
                                            <button onClick={() => openEditModal(item)} className="text-indigo-600 hover:text-indigo-900 transition duration-150 p-1 rounded-full hover:bg-indigo-100">
                                                <Edit className="h-5 w-5" />
                                            </button>
                                            <button onClick={() => handleDeleteItem(item.id_item)} className="text-red-600 hover:text-red-900 transition duration-150 p-1 rounded-full hover:bg-red-100">
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

            {/* Modal untuk Edit Item */}
            {isEditing && editingItem && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
                    <div className="bg-white rounded-lg shadow-xl w-full max-w-lg p-6 space-y-4">
                        <div className="flex justify-between items-center">
                            <h3 className="text-xl font-semibold">Edit Item Belanja</h3>
                            <button onClick={() => setIsEditing(false)} className="text-gray-400 hover:text-gray-600">
                                <X className="h-6 w-6" />
                            </button>
                        </div>
                        <form onSubmit={handleUpdateItem} className="space-y-4">
                            {/* Form fields... sama seperti form tambah */}
                            <div>
                                <label htmlFor="editNamaItem" className="block text-sm font-medium text-gray-700">Nama Barang</label>
                                <input
                                    type="text"
                                    id="editNamaItem"
                                    value={editingItem.nama_item}
                                    onChange={(e) => setEditingItem({...editingItem, nama_item: e.target.value})}
                                    className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                                    required
                                />
                            </div>
                            <div>
                                <label htmlFor="editIdKategori" className="block text-sm font-medium text-gray-700">Kategori</label>
                                <select
                                    id="editIdKategori"
                                    value={editingItem.id_kategori}
                                    onChange={(e) => setEditingItem({...editingItem, id_kategori: e.target.value})}
                                    className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                                    required
                                >
                                    <option value="" disabled>Pilih Kategori</option>
                                    {kategoriList.map(k => (
                                        <option key={k.id_kategori} value={k.id_kategori}>{k.nama_kategori}</option>
                                    ))}
                                </select>
                            </div>
                            <div>
                                <label htmlFor="editJumlah" className="block text-sm font-medium text-gray-700">Jumlah</label>
                                <input
                                    type="number"
                                    id="editJumlah"
                                    value={editingItem.jumlah_item}
                                    onChange={(e) => setEditingItem({...editingItem, jumlah_item: e.target.value})}
                                    className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                                    min="1"
                                    required
                                />
                            </div>
                            <div>
                                <label htmlFor="editHarga" className="block text-sm font-medium text-gray-700">Harga Satuan (Rp)</label>
                                <input
                                    type="number"
                                    id="editHarga"
                                    value={editingItem.harga_satuan}
                                    onChange={(e) => setEditingItem({...editingItem, harga_satuan: e.target.value})}
                                    className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                                    min="0"
                                    required
                                />
                            </div>
                            <div className="flex justify-end pt-4 border-t">
                                <button
                                    type="button"
                                    onClick={() => setIsEditing(false)}
                                    className="py-2 px-4 border border-gray-300 rounded-lg shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 mr-3"
                                >
                                    Batal
                                </button>
                                <button
                                    type="submit"
                                    disabled={isSubmitting}
                                    className="flex items-center justify-center py-2 px-6 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-indigo-500 disabled:bg-gray-400"
                                >
                                    {isSubmitting ? <Loader className="animate-spin h-5 w-5 mr-2" /> : <Save className="h-5 w-5 mr-2" />}
                                    Simpan Perubahan
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
};

export default DaftarBelanja;

