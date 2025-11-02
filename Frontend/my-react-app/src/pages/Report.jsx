import React, { useState, useEffect, useCallback } from 'react';
import { 
    BarChart, Tooltip, ResponsiveContainer, Bar as RechartsBar, 
    XAxis, YAxis, CartesianGrid, Legend 
} from 'recharts';
import { BarChart3, Loader, AlertCircle, Download } from 'lucide-react'; 

const API_BASE_URL = 'http://localhost:8080/api/v1';

// Komponen ChartCard
const ChartCard = ({ title, children, icon: Icon }) => (
  <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-200 h-full flex flex-col">
    <div className="flex items-center justify-between border-b pb-3 mb-4">
        <h3 className="text-xl font-semibold text-gray-800">{title}</h3>
        {Icon && <Icon className="w-6 h-6 text-indigo-500" />}
    </div>
    <div className="flex-grow w-full h-72"> 
      {children}
    </div>
  </div>
);

const BarIcon = BarChart; // Alias

/**
 * Halaman Laporan Pengeluaran
 */
const Reports = () => {
    const [reportData, setReportData] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState(null);
    
    const [isDownloading, setIsDownloading] = useState(false);
    const [downloadError, setDownloadError] = useState(null);

    // fetchWithAuth (Sudah Benar)
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

    // fetchData (Sudah Benar)
    // Fungsi ini sudah benar mengambil data dari { "data": { "pie_chart": [...] } }
    // sesuai dengan dash_handler.go
    const fetchData = useCallback(async () => {
        setIsLoading(true);
        setError(null);
        try {
            // Memanggil endpoint charts
            const chartsData = await fetchWithAuth(`${API_BASE_URL}/dashboard/charts`);
            // Membaca data yang dibungkus
            if (chartsData && chartsData.data && chartsData.data.pie_chart) {
                setReportData(chartsData.data.pie_chart);
            } else {
                setReportData([]);
            }
        } catch (err) {
            setError(err.message);
            console.error("Gagal mengambil data laporan:", err);
        } finally {
            setIsLoading(false);
        }
    }, [fetchWithAuth]);

    // useEffect (Sudah Benar)
    useEffect(() => {
        fetchData();
    }, [fetchData]);

    
    // handleDownloadReport (Sudah Benar)
    const handleDownloadReport = async () => {
        setIsDownloading(true);
        setDownloadError(null);
        
        try {
            const token = localStorage.getItem('authToken');
            const headers = {};
            if (token) {
                headers['Authorization'] = `Bearer ${token}`;
            }

            const response = await fetch(`${API_BASE_URL}/reports/download?type=excel`, {
                method: 'GET',
                headers: headers
            });

            if (!response.ok) {
                const errorData = await response.json().catch(() => ({})); 
                throw new Error(errorData.message || errorData.error || 'Gagal mengunduh laporan');
            }

            const blob = await response.blob();
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = `Laporan_Mingguan_${new Date().toISOString().split('T')[0]}.xlsx`; 
            document.body.appendChild(a);
            a.click();
            a.remove();
            window.URL.revokeObjectURL(url);

        } catch (err) {
            setDownloadError(err.message);
            console.error("Gagal mengunduh file:", err);
        } finally {
            setIsDownloading(false);
        }
    };


    // totalPengeluaran dan formatCurrency (Sudah Benar)
    const totalPengeluaran = reportData.reduce((sum, item) => sum + item.value, 0);
    const formatCurrency = (amount) => {
        return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(amount);
    };

    // Tampilan Loading dan Error (Sudah Benar)
    if (isLoading) {
        // ... (kode loading)
    }
    if (error) {
        // ... (kode error)
    }

    return (
        <div className="p-4 sm:p-6 lg:p-8 bg-gray-100 min-h-screen">
            {/* Tombol Download (Sudah Benar) */}
            <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-8">
                <h2 className="text-3xl font-extrabold text-gray-900 flex items-center mb-4 sm:mb-0">
                    <BarChart3 className="w-8 h-8  text-indigo-600" />
                    Laporan Pengeluaran Bulanan
                </h2>
                <button
                    onClick={handleDownloadReport}
                    disabled={isDownloading}
                    className="flex items-center justify-center w-full sm:w-auto py-2 px-6 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:bg-gray-400"
                >
                    {isDownloading ? (
                        <Loader className="animate-spin h-5 w-5 mr-2" />
                    ) : (
                        <Download className="h-5 w-5 mr-2" />
                    )}
                    {isDownloading ? 'Memproses...' : 'Unduh Laporan (Excel)'}
                </button>
            </div>
            
            {/* Error download (Sudah Benar) */}
            {downloadError && (
                 <div className="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 rounded-lg mb-6" role="alert">
                    <p className="font-bold">Gagal Mengunduh</p>
                    <p>{downloadError}</p>
                </div>
            )}

            {/* Total Pengeluaran (Sudah Benar) */}
            <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-200 mb-8">
                <p className="text-lg font-medium text-gray-500">Total Pengeluaran Bulan Ini</p>
                <p className="text-4xl font-extrabold text-red-600 mt-1">
                    {formatCurrency(totalPengeluaran)}
                </p>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                
                {/* ======================================================
                  =============== PERBAIKAN 1 DI SINI ==================
                  ======================================================
                */}
                <div className="lg:col-span-2">
                    <ChartCard title="Distribusi Pengeluaran per Kategori" icon={BarIcon}>
                        <ResponsiveContainer width="100%" height="100%">
                            <BarChart data={reportData} margin={{ top: 20, right: 20, left: 20, bottom: 5 }}>
                                <CartesianGrid strokeDasharray="3 3" />
                                {/* Ganti dataKey="category" menjadi "name" */}
                                <XAxis dataKey="name" />
                                <YAxis />
                                <Tooltip formatter={(value) => formatCurrency(value)} />
                                <Legend />
                                <RechartsBar dataKey="value" fill="#6366F1" name="Pengeluaran" />
                            </BarChart>
                        </ResponsiveContainer>
                    </ChartCard>
                </div>
                
                {/* ======================================================
                  =============== PERBAIKAN 2 DI SINI ==================
                  ======================================================
                */}
                <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-200">
                    <h3 className="text-xl font-semibold text-gray-800 mb-4">Rincian Kategori</h3>              
                    <table className="w-full table-auto">
                        <thead>
                            <tr className="bg-gray-100">
                                <th className="px-4 py-2 text-left text-sm font-medium text-red-600">Kategori</th>
                                <th className="px-4 py-2 text-right text-sm font-medium text-red-600">Jumlah Pengeluaran</th>
                            </tr>
                        </thead>
                        <tbody>
                            {reportData.map((item, index) => (
                                <tr key={index} className={index % 2 === 0 ? 'bg-red' : 'bg-gray-50'}>
                                    {/* Ganti item.nama_item menjadi item.name */}
                                    <td className="px-4 py-2 text-sm text-gray-700">{item.name}</td>
                                    <td className="px-4 py-2 text-sm text-gray-700 text-right">{formatCurrency(item.value)}</td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
                
            </div>
        </div>
    );
}

export default Reports;